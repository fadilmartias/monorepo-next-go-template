# Code Standards — Go backend

## General principles

- Keep packages small and single-purpose. Prefer composition over
  large god-objects.
- Fix root causes; do not layer workarounds. Prefer clear error
  propagation using `error` wrapping (`%w`).
- Follow existing repository patterns: `http/` controllers →
  `services/` → `repositories/`.

## Go idioms and tooling

- Use `context.Context` for request-scoped data, propagation,
  and cancellation.
- Use structured logging with `zap`. Log relevant fields, not
  entire objects.
- Return typed errors from repository/service layers and translate
  to HTTP responses in controllers.
- Run `gofmt` and `go vet` before commits. Keep functions short
  and avoid deep nesting.

## Validation

- Validate and sanitize external input at the request boundary
  (controllers) using `go-playground/validator` or explicit checks.
- Map validated DTOs to domain models in `requests/` and
  `responses/` packages.

## Database and transactions

- All DB access must go through `repositories/`.
- Keep transactions at the service layer when multiple repository
  calls must be atomic.

## Background jobs

- Use `asynq` for background tasks. Handlers in `jobs/` must be
  idempotent and durable.

## APIs and responses

- Return consistent JSON shapes via `responses/` helpers.
- Controllers must be thin: parse/validate input, call service,
  return response or mapped error.

## File organization (repository conventions)

- `app/` — application wiring and bootstrap helpers.
- `http/` — controllers, middleware, route registration.
- `services/` — business logic.
- `repositories/` — DB access and queries.
- `models/` — GORM models.
- `jobs/` — background workers and cron tasks.
- `requests/` & `responses/` — DTOs and API shapes.
