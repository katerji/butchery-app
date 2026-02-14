# /verify — Run All Checks

Run the full verification suite for the backend. Execute these steps in order from the `backend/` directory:

1. **Unit tests**: Run `go test ./...` and report results.
2. **E2E tests**: Run `go test ./internal/e2e/...` and report results.
3. **Linter**: Run `golangci-lint run ./...` and report results.

After all three complete, provide a summary:
- Total unit tests passed/failed
- Total e2e tests passed/failed
- Linter issues (if any)

If any step fails, report the failures clearly so they can be fixed.
