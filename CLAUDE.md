# CLAUDE.md — Butchery App

## Project Overview

A butchery management application built as a **monorepo** containing:
- **Backend**: Go (REST API) with PostgreSQL
- **Frontend**: Next.js + TypeScript + Tailwind CSS + shadcn/ui

---

## Interaction Rules

> These rules are **non-negotiable** and must be followed at all times.

- **NEVER assume requirements.** If anything is ambiguous, unclear, or could be interpreted in multiple ways — stop and ask a clarification question before proceeding.
- **NEVER assume business logic.** Always ask what the expected behavior should be. Do not guess domain rules.
- **NEVER implement a feature without confirming the acceptance criteria** with the user first.
- **Always confirm before**:
  - Creating a new bounded context or aggregate
  - Adding a new database table or modifying an existing schema
  - Introducing a new dependency or third-party library
  - Making architectural decisions that affect multiple layers
  - Deleting or renaming existing code that other parts depend on
- When presenting options, explain trade-offs clearly and let the user decide.
- Prefer asking one focused question over making a wrong assumption.

---

## Architecture Principles

### Clean Architecture

All backend code MUST follow **Clean Architecture** with strict layer separation. Dependencies point **inward only** — outer layers depend on inner layers, never the reverse.

```
[Handlers/Controllers] → [Use Cases/Application] → [Domain] ← [Infrastructure/Repository]
```

**Layers (innermost to outermost):**

1. **Domain Layer** (innermost — zero dependencies)
   - Entities, Value Objects, Aggregates, Domain Events, Domain Services
   - Repository interfaces (ports) — defined here, implemented in infrastructure
   - Contains all business rules and invariants
   - MUST NOT import from any other layer
   - MUST NOT depend on any framework, database driver, or external library

2. **Application Layer** (Use Cases)
   - Application services / use case interactors
   - DTOs (command/query objects) for input and output
   - Orchestrates domain objects to fulfill use cases
   - Depends ONLY on the Domain layer
   - MUST NOT contain business logic — delegate to domain

3. **Infrastructure Layer** (outermost)
   - Repository implementations (PostgreSQL)
   - External service adapters (email, payment, etc.)
   - Database migrations and connection management
   - Implements interfaces defined in the Domain layer

4. **Interface / Presentation Layer** (outermost)
   - HTTP handlers / controllers (REST API)
   - Request/Response mapping and validation
   - Middleware (auth, logging, CORS, etc.)
   - Route definitions

### SOLID Principles

Every piece of code must adhere to SOLID:

- **S — Single Responsibility**: Each struct/function/file has one reason to change.
- **O — Open/Closed**: Extend behavior through interfaces and composition, not modification.
- **L — Liskov Substitution**: Implementations must be substitutable for their interfaces.
- **I — Interface Segregation**: Define small, focused interfaces. Prefer many small interfaces over one large one. Go idiom: accept interfaces, return structs.
- **D — Dependency Inversion**: High-level modules depend on abstractions (interfaces), not concrete implementations. Use dependency injection — wire dependencies in `main` or a composition root.

### Domain-Driven Design (DDD)

DDD is the **primary design methodology** for this project. Every feature must be modeled through DDD concepts.

**Strategic Design:**
- Identify and define **Bounded Contexts** before writing code. Ask the user to confirm context boundaries.
- Map relationships between bounded contexts (Shared Kernel, Anti-Corruption Layer, etc.).
- Use **Ubiquitous Language** — code must reflect the language of the butchery domain. Variable names, function names, types, and comments must use domain terminology (not generic tech jargon).

**Tactical Design:**
- **Entities**: Objects with identity that persists over time. Implement equality by ID.
- **Value Objects**: Immutable objects defined by their attributes. No identity. Implement equality by value. Use for things like Money, Weight, Address, etc.
- **Aggregates**: Cluster of entities/value objects with a single Aggregate Root. All modifications go through the root. Each aggregate is a transactional boundary.
- **Domain Events**: Represent something meaningful that happened in the domain. Use for cross-aggregate communication.
- **Domain Services**: Business logic that doesn't naturally belong to a single entity or value object.
- **Repositories**: One repository per aggregate root. Interface in domain, implementation in infrastructure.
- **Factories**: Use when object creation is complex or involves invariant validation.

**DDD Rules:**
- Always model the domain FIRST, before thinking about persistence or API design.
- Never let infrastructure concerns leak into the domain layer.
- Aggregate roots enforce all invariants for their cluster.
- Reference other aggregates by ID only, not by direct object reference.
- One transaction = one aggregate. Use domain events for eventual consistency across aggregates.

---

## Test-Driven Development (TDD)

TDD is **mandatory** for all backend code. No production code may be written without a failing test first.

### The TDD Cycle (Red-Green-Refactor)

1. **RED**: Write a failing test that defines the expected behavior. Run it — it must fail.
2. **GREEN**: Write the **minimum** production code to make the test pass. No more.
3. **REFACTOR**: Clean up the code while keeping all tests green. Remove duplication, improve naming, extract functions/types.

### TDD Rules

- **Always start by writing a test.** Before touching any production code, write a test that specifies the behavior.
- **One behavior per test.** Each test should assert one logical concept.
- **Tests drive the design.** If something is hard to test, the design needs improvement — not the test.
- **Run tests after every change.** Never assume code works without running the tests.
- **Name tests descriptively**: `TestCreateOrder_WhenInsufficientStock_ReturnsError` — use the pattern `Test<Unit>_<Scenario>_<ExpectedResult>`.

### Testing Strategy by Layer

| Layer          | Test Type           | What to Test                                      |
|----------------|---------------------|---------------------------------------------------|
| Domain         | Unit tests          | Entities, value objects, domain services, aggregates — pure business logic |
| Application    | Unit tests (mocked) | Use cases with mocked repositories/services       |
| Infrastructure | Integration tests   | Repository implementations against real PostgreSQL (use testcontainers) |
| Interface      | Integration tests   | HTTP handlers with full request/response cycle     |
| End-to-End     | E2E tests           | Critical user flows across the full stack          |

### Testing Tools

- **Go standard `testing` package** + **testify** for assertions and mocks
- **testify/assert** for assertions: `assert.Equal`, `assert.NoError`, `assert.ErrorIs`
- **testify/mock** for mocking interfaces
- **testify/suite** for test suites when setup/teardown is needed
- Use **table-driven tests** for testing multiple input/output scenarios
- Use **testcontainers-go** for integration tests requiring a real PostgreSQL instance

---

## Project Structure (Monorepo)

```
butchery-app/
├── CLAUDE.md
├── docker-compose.yml
├── Makefile
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go                 # Composition root, dependency wiring
│   ├── internal/
│   │   ├── domain/                     # Domain layer
│   │   │   ├── <context>/             # One package per bounded context
│   │   │   │   ├── entity.go
│   │   │   │   ├── value_object.go
│   │   │   │   ├── aggregate.go
│   │   │   │   ├── repository.go       # Interface only
│   │   │   │   ├── service.go          # Domain services
│   │   │   │   ├── event.go            # Domain events
│   │   │   │   └── errors.go           # Domain-specific errors
│   │   │   │
│   │   ├── application/                # Application layer (use cases)
│   │   │   ├── <context>/
│   │   │   │   ├── commands/           # Write operations
│   │   │   │   ├── queries/            # Read operations
│   │   │   │   └── dto/                # Data transfer objects
│   │   │   │
│   │   ├── infrastructure/             # Infrastructure layer
│   │   │   ├── persistence/
│   │   │   │   ├── postgres/
│   │   │   │   │   ├── <context>_repository.go
│   │   │   │   │   └── migrations/
│   │   │   │   └── connection.go
│   │   │   └── external/               # External service adapters
│   │   │
│   │   └── interface/                  # Interface layer
│   │       └── http/
│   │           ├── handler/
│   │           │   └── <context>_handler.go
│   │           ├── middleware/
│   │           ├── router.go
│   │           └── dto/                # Request/Response types
│   │
│   ├── pkg/                            # Shared utilities (logging, config, etc.)
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── app/                        # Next.js App Router
│   │   ├── components/
│   │   │   ├── ui/                     # shadcn/ui components
│   │   │   └── features/              # Feature-specific components
│   │   ├── lib/                        # Utilities, API client, types
│   │   └── styles/
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.ts
│   └── next.config.ts
│
└── .github/
    └── workflows/                      # CI/CD pipelines
```

---

## Backend Conventions (Go)

### General
- Use **Go modules** for dependency management.
- Target the latest stable Go version.
- Run `go vet`, `golangci-lint`, and `go test ./...` before considering any code complete.
- Use `context.Context` as the first parameter for any function that does I/O or can be cancelled.
- Handle errors explicitly. Never ignore errors with `_`. Wrap errors with `fmt.Errorf("context: %w", err)` for stack tracing.
- Return domain-specific errors, not generic ones.

### Naming Conventions
- Follow standard Go naming conventions (camelCase for unexported, PascalCase for exported).
- Package names: short, lowercase, singular (e.g., `order`, `product`, `inventory`).
- Interface names: descriptive verbs or nouns (e.g., `OrderRepository`, `PriceCalculator`).
- Avoid stuttering: `order.Order` is fine, but avoid `order.OrderService` — prefer `order.Service`.

### Dependency Injection
- Wire all dependencies in `cmd/api/main.go` (the composition root).
- Use constructor functions: `func NewOrderService(repo OrderRepository) *OrderService`.
- Never use global state or `init()` functions for dependency setup.

---

## Frontend Conventions (Next.js + TypeScript)

- Use the **App Router** (not Pages Router).
- All components in TypeScript (`.tsx`).
- Use **shadcn/ui** for UI components — do not build custom components when shadcn provides one.
- Use **Tailwind CSS** for all styling — no CSS modules, no styled-components.
- Keep components small and focused. Extract reusable logic into custom hooks.
- Use `fetch` or a lightweight client for API calls. Keep API logic in `lib/api/`.
- Validate forms with **React Hook Form** + **Zod** for schema validation.

---

## Database Conventions (PostgreSQL)

- Use database migrations for all schema changes — never modify the database manually.
- Migration files must be versioned and sequential.
- Each migration must be reversible (include both `up` and `down`).
- Use UUIDs for primary keys.
- Table and column names: `snake_case`.
- Always add created_at and updated_at timestamps.
- Write indexes for columns that will be queried frequently.
- Never expose database models directly — always map to/from domain entities.

---

## REST API Conventions

- Follow RESTful resource naming: `/api/v1/<resource>` (plural nouns).
- Use proper HTTP methods: GET (read), POST (create), PUT (full update), PATCH (partial update), DELETE.
- Return consistent JSON response envelopes:
  ```json
  {
    "data": {},
    "error": null,
    "meta": { "page": 1, "total": 100 }
  }
  ```
- Use proper HTTP status codes (200, 201, 204, 400, 401, 403, 404, 409, 422, 500).
- Version all APIs (`/api/v1/...`).
- Validate all incoming requests at the handler level before passing to use cases.

---

## Git Workflow (Git Flow)

- **main**: Production-ready code only. Never commit directly.
- **develop**: Integration branch. All feature branches merge here.
- **feature/<ticket-or-name>**: New features, branched from `develop`.
- **bugfix/<ticket-or-name>**: Bug fixes, branched from `develop`.
- **release/<version>**: Release preparation, branched from `develop`, merged to `main` + `develop`.
- **hotfix/<ticket-or-name>**: Urgent production fixes, branched from `main`, merged to `main` + `develop`.

### Commit Messages
- Use **Conventional Commits**: `type(scope): description`
- Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `ci`
- Examples:
  - `feat(order): add create order use case`
  - `test(product): add unit tests for price calculation`
  - `fix(inventory): correct stock deduction logic`

---

## Docker & CI/CD

- **docker-compose.yml** at the repo root for local development (Go API + PostgreSQL + Frontend).
- Backend Dockerfile: multi-stage build (build stage + minimal runtime image).
- Frontend Dockerfile: multi-stage build with Next.js standalone output.
- CI pipeline must:
  1. Lint (golangci-lint for Go, ESLint for frontend)
  2. Run all tests (unit + integration)
  3. Build Docker images
  4. Run E2E tests (optional, on PR to main)

---

## Workflow for New Features

When implementing any new feature, always follow this order:

1. **Clarify requirements** — ask the user about expected behavior, edge cases, and acceptance criteria.
2. **Identify the bounded context** — confirm which context this belongs to.
3. **Model the domain first** — define entities, value objects, aggregates, and domain events. Get confirmation.
4. **Write domain tests (TDD)** — write failing tests for domain logic, then implement.
5. **Write use case tests (TDD)** — write failing tests for the application layer, then implement.
6. **Implement infrastructure** — repository implementations, migrations.
7. **Write integration tests** — test repositories against a real database.
8. **Implement HTTP handlers** — request/response handling, validation.
9. **Write handler tests** — test the full HTTP flow.
10. **Implement the frontend** — UI components, API integration.
