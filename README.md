# Calculator

Full stack calculator. Go REST API for the math, React + TypeScript for the UI. Seven operations: add, subtract, multiply, divide, power, square root, percentage.

## General Design & Architecture

- Test-First strategy: Core domain and HTTP layers were built alongside table-driven test cases.

- No over-engineering: Clean Architecture is ceremony for a simple domain package, so it was deliberately skipped. The only strict boundary is keeping net/http entirely out of the math package.

- Backend dictates operations: The frontend dynamically fetches supported operations and metadata from the backend to prevent configuration drift.

- Single source of truth: Math executes strictly on the server using IEEE 754 float64. The frontend only formats strings for display and never rounds for active calculations.

- Scope cuts: Dropped husky pre-commits and rate limiting. A token bucket rate limiter keyed by IP makes sense for production, but falls outside the current scope.

- Stateless: No server-side persistence or auth. History is capped at 10 items and lives entirely in the browser's localStorage.

## Prerequisites & Setup

- Go 1.27+
- Node 24+

**Run the backend:**

```bash
cd backend
go run .
# listening on :8080
```

**Run the frontend:**

```bash
cd frontend
npm install
npm run dev
# proxying /api to :8080 via Vite
```

**Or run the whole stack in Docker:**

```bash
docker compose up --build
# http://localhost:8080
```

nginx serves the built SPA and proxies `/api/` to the api service, so the browser sees a single origin and there's no `VITE_API_URL` baked into the bundle. Only `web` publishes a port; the API is reachable only on the compose network. `web` waits on `condition: service_healthy`, not plain `depends_on`, so it doesn't come up before the API can actually answer.

On Windows PowerShell, use `curl.exe`, not `curl` (it's aliased to `Invoke-WebRequest`, different flags).

**Tests:**

```bash
cd backend && go test -race -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | tail -1
cd frontend && npm run test:run && npm run lint && npm run build
```

## API Usage

**POST /api/v1/calculate**

```bash
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":12,"b":3}'
# {"operation":"divide","a":12,"b":3,"result":4}
```

`b` is omitted for unary operations:

```bash
curl -s -X POST localhost:8080/api/v1/calculate -H 'Content-Type: application/json' -d '{"operation":"sqrt","a":144}'
# {"operation":"sqrt","a":144,"result":12}
```

**GET /api/v1/operations** - the frontend builds its dropdown from this, so the two sides can't drift apart.

```bash
curl -s localhost:8080/api/v1/operations
# [{"name":"add","symbol":"+","arity":2,"label":"Add"}, ...all 7]
```

**GET /api/v1/health** → `{"status":"ok"}`

**Errors** - always the plural `errors` array, even for one failure. `code` is what the client switches on, `field` is what lets the form mark the right input red.

```json
{
  "errors": [
    {
      "code": "DIVISION_BY_ZERO",
      "message": "cannot divide by zero",
      "field": "b"
    }
  ]
}
```

400 means the body couldn't be parsed. 422 means it parsed fine and is wrong.

| Condition                      | Status | code                       |
| ------------------------------ | ------ | -------------------------- |
| Malformed JSON / unknown field | 400    | `INVALID_JSON`             |
| Body over 1MB                  | 413    | `REQUEST_TOO_LARGE`        |
| Operand missing for the arity  | 422    | `MISSING_OPERAND`          |
| Unknown operation              | 422    | `UNKNOWN_OPERATION`        |
| Divide by zero                 | 422    | `DIVISION_BY_ZERO`         |
| Square root of a negative      | 422    | `NEGATIVE_SQUARE_ROOT`     |
| Result is NaN or Inf           | 422    | `RESULT_NOT_REPRESENTABLE` |
| Unknown path                   | 404    | `NOT_FOUND`                |
| Wrong method                   | 405    | `METHOD_NOT_ALLOWED`       |

**Env vars** (all optional, `backend/.env.example` has the full list): `PORT` (8080), `APP_ENV` (development/production), `CORS_ORIGIN` (http://localhost:5173), `MAX_BODY_BYTES` (1048576), `SHUTDOWN_TIMEOUT` (10s). Config is validated at boot, a bad value fails startup instead of failing lazily on the first request.

## Implementation Steps

### 1. Basic Operations (Go Domain)

Objective: Write basic operations inside `backend/internal/calc/basic.go`. All basic operations must receive float64 parameters and return float64. Divide must return `(float64, error)` to handle division by zero as a value, not a panic.

Tests: Create `basic_test.go` using Go's table-driven pattern: anonymous struct slice, test loop, each scenario as an isolated subtest (`t.Run`), error handling per iteration.

### 2. Power, Sqrt, and Percentage

Objective: Add 3 extra operation methods receiving float64, in `advanced.go`. Power and sqrt can error, so they return `error`. Sentinel errors split into a dedicated `errors.go`.

Tests: Same table-driven pattern.

### 3. Guard Against Non-Operable Values

Objective: A guard at the domain boundary (`finite.go`) so NaN and Inf never reach the JSON encoder, which would otherwise turn a domain failure into an opaque 500.

Tasks: Guard covers NaN and Inf. New `ErrResultNotRepresentable` in `errors.go`.

### 4. Operation Registry and Apply Dispatcher

Objective: Translate a request like `{"operation": "add", "a": 5, "b": 10}` into a function call safely.

Tasks: `Descriptor` ties together operation, symbol, arity, and label. `Operations()` returns a copy, never the live slice. `Apply` checks arity first, then that the operand is present, then dispatches.

Tests: Loop over mock requests covering every operation and the arity failure case.

### 5. JSON Helper & Error Envelope

Objective: `writeJSON` (content-type, status code, encode) and the error envelope (`respond/`).

Tasks: Status set before the body is written, since `w.Write` implicitly sends 200. Error struct always returns the plural `errors` array with `code`, `message`, `field`.

### 6. Health and Operations Endpoints

Objective: `/api/v1/health` and `/api/v1/operations`.

Tasks: Health confirms the server is up. Operations calls `calc.Operations()` to feed the frontend dropdown.

### 7. Router & Fallbacks

Objective: Wire the routes, handle everything not explicitly declared.

Tasks: POST/GET routes for calculate, health, operations. Catch-all `NotFound` and a `MethodNotAllowed` factory, both returning the JSON envelope instead of Go's plain text default.

### 8. Server Config & main.go

Objective: Keep `main.go` tiny, push startup, timeouts, and shutdown into `run.go`.

Tasks: Middleware for 500 JSON recovery, `X-Request-Id`, status-capturing logging, CORS, and a chain composer, in that order (recovery outermost, since a panic otherwise kills the whole process). `.env` / `.env.example`, validated cleanly on boot.

### 9. E2E Testing

Objective: One table test over `httptest` covering the full contract.

Tasks: Status mapping, unknown JSON field, oversized body, unknown path, 404. Assertions on `code` and `field` only, never on message text, since Go 1.27's `encoding/json` v2 backend changed error strings.

### 10. Frontend: Vite & Tailwind

Objective: Modern, fast React shell.

Tasks: Vite's scaffold installs TypeScript@latest by default, which is 7.x, the Go-native compiler. typescript-eslint only supports up to `<6.1.0`, so type-aware linting silently stops working on 7.x with nothing but a console warning. Pinned to the 6.0 line instead. Tailwind v4 (CSS-first, no `tailwind.config.js`). ESLint flat config + Prettier based on the knowledge vault rules.

### 11. API Proxying & Test Tooling

Objective: Connect front to back, get a real test runner in place.

Tasks: Proxy `/api` to the Go service in `vite.config.ts`. Manual check: `go run .`, `npm run dev`, `curl.exe -s http://localhost:5173/api/v1/health`. Vitest + React Testing Library, `setup.ts` wired into Vite for global test config.

### 12. Zod Schemas & API Client

Objective: Strict boundary control over what's sent and what comes back.

Tasks: `lib/schemas.ts`: `operationName` (ended up a plain non-empty string, not an enum, since the operations list comes from the backend and a hardcoded enum of names would fight the discovery endpoint instead of trusting it), `calculateResponse`, `apiError`, `errorBody`. `api.ts`: `ApiErrorResponse` class, `calculate` via `fetch`.

### 13. Component Refactor

Objective: Fix structural debt in the first pass at the form.

Tasks: `validateOperand` moved to `utils/`. `useState` logic extracted into `useCalculatorForm`, a custom hook. Prop types moved to their own file. Shared code moved to the right feature folder instead of sitting in the component.

### 14. Bug Fixes From a Fresh Review

Objective: Audit what #13 actually produced before building on top of it.

Tasks: `CalculatorForm` was concatenating a stale server error with a fresh client validation error instead of replacing it. `useCalculatorForm` had a `safeParse` on the operation dropdown that could never fail, since the value only ever came from the app's own rendered options. `http.ts` hand-rolled its own `ApiError` type instead of importing the one already inferred from the zod schema. `useCalculator` never cleared the previous successful result when a later calculation failed, so an old result and a new error showed on screen together. Fixed all four, with a regression test for each.

### 15. History Panel + localStorage

Objective: A short calculation history, newest first, surviving a reload.

Tasks: Array in `useCalculator`, capped at 10, id generated once at creation time, not during render, so list keys stay stable. `historyStorage.ts` wraps `localStorage` reads/writes in `try/catch` (storage can be full, disabled, or blocked in private browsing) and validates what comes back with the same zod schema used for API responses, so corrupted or outdated data is discarded instead of crashing the app.

Tests: Caught jsdom keeping `localStorage` alive across tests in the same file, unlike React state resetting. Fixed by clearing it globally in the shared test setup instead of patching each file.

### 16. UI Pass

Objective: The app worked but read as unstyled, including a genuinely broken dropdown (white text on a white popup).

Tasks: Root cause of the dropdown was a parent `text-white` class the `<select>` inherited, while the native option popup painted the OS default regardless of that. Fixed with `color-scheme: dark`. Built a Tailwind v4 theme (one accent, one neutral ramp, self-hosted variable fonts), a parallax background using refs and `requestAnimationFrame` instead of `useState` for pointer tracking, and a two-column layout collapsing to one below `md`.

Verification: Opened the running app in a real browser and drove it end to end. Measured contrast ratios in-browser and caught one real WCAG AA failure (3.56:1) that reading the CSS wouldn't have shown, fixed it, re-measured at 5.24:1. Ran the calculate/error/history flows for real and checked 390px width for overflow.

### 17. Docs

Objective: Write this file.

Tasks: Verified every API example against real `curl` output from the running server instead of from memory, since earlier notes had drifted from what actually shipped in a few places (corrected inline above).

## Coverage

Backend: 83.1% statements overall. Domain (`internal/calc`) and transport (`internal/httpapi/*`) are 97-100%; the total is pulled down by `main`/`run.go` at 12.5%, which is signal handling and `ListenAndServe`, hard to cover without standing up and killing a real server.

Frontend: 90.7% statements, 77.5% branches, 96% functions.

```bash
cd backend && go tool cover "-html=coverage.out" # the view worth reading
cd frontend && npm run test:coverage
```

## Assumptions & Limitations

- No server-side persistence, no auth. Every operation is naturally idempotent.
- History is per-browser via `localStorage`, capped at 10, not synced anywhere.
- `float64` throughout. Honest choice for a general calculator; the obligation is documenting the precision, not pretending it's exact. `formatResult` rounds for display only and never feeds back into a calculation.
- `percentage(a, b)` is "b percent of a" (`a * b / 100`), genuinely ambiguous, pinned here.
- Rate limiting, Docker, and CI are not in this repo.
