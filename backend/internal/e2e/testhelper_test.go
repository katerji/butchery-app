package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	admincmd "github.com/katerji/butchery-app/backend/internal/application/admin/commands"
	authcmd "github.com/katerji/butchery-app/backend/internal/application/auth/commands"
	custcmd "github.com/katerji/butchery-app/backend/internal/application/customer/commands"
	infraauth "github.com/katerji/butchery-app/backend/internal/infrastructure/auth"
	pgrepo "github.com/katerji/butchery-app/backend/internal/infrastructure/persistence/postgres"
	apphttp "github.com/katerji/butchery-app/backend/internal/interface/http"
	"github.com/katerji/butchery-app/backend/internal/interface/http/handler"
	"github.com/katerji/butchery-app/backend/internal/interface/http/middleware"
	"github.com/katerji/butchery-app/backend/pkg/httpresponse"
)

const (
	testJWTSecret     = "e2e-test-secret-that-is-long-enough-for-hmac-sha256"
	testAdminEmail    = "admin@butchery.com"
	testAdminPassword = "admin123"
)

func init() {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
}

// testServer wraps an httptest.Server and pgxpool.Pool for e2e tests.
type testServer struct {
	server *httptest.Server
	pool   *pgxpool.Pool
}

// setupTestServer starts a PostgreSQL testcontainer with all migrations,
// wires the full application stack (mirroring cmd/api/main.go), and returns
// a running httptest.Server. The container and server are cleaned up when the
// test finishes.
func setupTestServer(t *testing.T, accessTokenTTL time.Duration) *testServer {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx := context.Background()

	_, currentFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(currentFile),
		"..", "infrastructure", "persistence", "postgres", "migrations")

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("butchery_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithInitScripts(
			filepath.Join(migrationsDir, "V1__create_admins_table.sql"),
			filepath.Join(migrationsDir, "V2__create_customers_table.sql"),
			filepath.Join(migrationsDir, "V3__create_refresh_tokens_table.sql"),
			filepath.Join(migrationsDir, "V4__seed_admin.sql"),
		),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err, "failed to start postgres container")

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "failed to get connection string")

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err, "failed to create connection pool")

	t.Cleanup(func() { pool.Close() })

	// Repositories
	adminRepo := pgrepo.NewAdminRepository(pool)
	customerRepo := pgrepo.NewCustomerRepository(pool)
	refreshTokenRepo := pgrepo.NewRefreshTokenRepository(pool)

	// Infrastructure services
	passwordHasher := infraauth.NewBcryptHasher()
	tokenService := infraauth.NewTokenService(testJWTSecret, accessTokenTTL)

	// Use case handlers
	adminLoginHandler := admincmd.NewAdminLoginHandler(adminRepo, passwordHasher, tokenService, refreshTokenRepo, accessTokenTTL)
	registerCustomerHandler := custcmd.NewRegisterCustomerHandler(customerRepo, passwordHasher)
	customerLoginHandler := custcmd.NewCustomerLoginHandler(customerRepo, passwordHasher, tokenService, refreshTokenRepo, accessTokenTTL)
	refreshTokenHandler := authcmd.NewRefreshTokenHandler(refreshTokenRepo, tokenService, accessTokenTTL)
	logoutHandler := authcmd.NewLogoutHandler(refreshTokenRepo)

	// HTTP handlers
	adminAuthHandler := handler.NewAdminAuthHandler(adminLoginHandler)
	customerAuthHandler := handler.NewCustomerAuthHandler(registerCustomerHandler, customerLoginHandler)
	authHandler := handler.NewAuthHandler(refreshTokenHandler, logoutHandler)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// Router
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	router := apphttp.NewRouter(apphttp.RouterDeps{
		Logger:              logger,
		AuthMiddleware:      authMiddleware,
		AdminAuthHandler:    adminAuthHandler,
		CustomerAuthHandler: customerAuthHandler,
		AuthHandler:         authHandler,
	})

	server := httptest.NewServer(router)
	t.Cleanup(func() { server.Close() })

	return &testServer{server: server, pool: pool}
}

// url returns the full URL for a given API path.
func (ts *testServer) url(path string) string {
	return ts.server.URL + path
}

// postJSON sends a POST request with a JSON body and returns the response.
func (ts *testServer) postJSON(t *testing.T, path string, body any) *http.Response {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	resp, err := http.Post(ts.url(path), "application/json", bytes.NewReader(jsonBody))
	require.NoError(t, err)

	return resp
}

// postJSONWithAuth sends an authenticated POST request with a JSON body.
func (ts *testServer) postJSONWithAuth(t *testing.T, path string, body any, token string) *http.Response {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, ts.url(path), bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

// responseEnvelope matches the httpresponse.Response structure.
type responseEnvelope struct {
	Data  json.RawMessage `json:"data,omitempty"`
	Error any             `json:"error,omitempty"`
}

// Ensure the test knows about the httpresponse package (compile-time check).
var _ = httpresponse.Response{}

// parseJSON decodes the response body into a responseEnvelope, then
// unmarshals the Data field into the provided target.
func parseJSON(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	defer resp.Body.Close()

	var env responseEnvelope
	err := json.NewDecoder(resp.Body).Decode(&env)
	require.NoError(t, err, "failed to decode response body")

	if target != nil && env.Data != nil {
		err = json.Unmarshal(env.Data, target)
		require.NoError(t, err, "failed to unmarshal data field")
	}
}

// parseError decodes the response body and returns the error string.
func parseError(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()

	var env responseEnvelope
	err := json.NewDecoder(resp.Body).Decode(&env)
	require.NoError(t, err, "failed to decode error response")

	if s, ok := env.Error.(string); ok {
		return s
	}
	return ""
}
