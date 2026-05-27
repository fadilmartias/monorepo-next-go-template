# Architecture — Backend service

## Stack

| Layer        | Technology / package            | Role                                  |
| ------------ | ------------------------------ | ------------------------------------- |
| Language     | Go 1.26.3                        | Application language                   |
| Framework    | github.com/gofiber/fiber/v3     | HTTP routing and middleware           |
| GraphQL      | github.com/99designs/gqlgen     | GraphQL schema & resolvers (optional) |
| ORM          | gorm.io/gorm + gorm.io/driver/mysql | Persistent models (MySQL)        |
| Cache/Queue  | Redis / github.com/hibiken/asynq | Caching & background jobs            |
| Logger       | go.uber.org/zap                | Structured logging                    |
| Storage      | S3 (aws-sdk-go-v2)             | File uploads                          |

## System Boundaries (folders)

- `http/` — Fiber controllers, middleware, route registration.
- `services/` — Business logic and orchestration.
- `repositories/` — Database access (GORM) and query abstractions.
- `models/` — GORM models and DB schemas.
- `jobs/` — Background worker handlers and cronjobs.
- `config/` — App configuration (env, storage, DB, redis).

## Storage Model

- Relational DB (MySQL): canonical source for domain entities
  (users, articles, orders, payment records, settings).
- Redis: caching, rate-limiting, ephemeral data, and Asynq queue.
- S3 (or S3-compatible storage): user uploads, generated files,
  and public assets. Local `tmp/` used for transient files.

## Auth and Access

- Primary auth: JWT access tokens issued by auth service.
- Secondary: password reset tokens, passkeys stored in
  `user_passkeys_model.go` when enabled.
- Authorization: perform ownership and role checks at service
  layer before mutating resources.

## Invariants

1. HTTP handlers must be thin — delegate to `services/` for logic.
2. Repositories own all DB interactions and return typed errors.
3. Context (`context.Context`) is propagated for cancellation/timeouts.
4. No long-running work inside request handlers; use Asynq/cron.
