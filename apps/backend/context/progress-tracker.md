# Progress Tracker

Update this file after every meaningful implementation change.

## Current Phase

- Completed — backend utility migration and validation.

## Current Goal

- Make context/ docs reflect the actual backend stack and
  repository patterns so contributors have an accurate guide.

## Completed

- Reviewed and updated `project-overview.md`.
- Reviewed and updated `architecture.md`.
- Reviewed and updated `ui-context.md`.
- Reviewed and updated `code-standards.md`.
- Reviewed and updated `ai-workflow-rules.md`.
- Replaced custom Redis wrapper with direct `*redis.Client` usage.
- Updated Redis wiring across controllers, services, jobs, and CLI.
- Replaced Resty-based HTTP helper with a Fiber v3 client adapter.
- Preserved the existing `Http()` fluent surface for backward compatibility.
- Verified the backend with `go test ./...` after the HTTP client migration.
- Migrated Telegram, Fonnte, and OAuth HTTP flows to the shared `Http()` helper.
- Removed the Resty dependency from module usage with `go mod tidy`.

## In Progress

- None.

## Next Up

- Keep docs in sync when implementing features; open PRs for
  context changes tied to implementation work.

## Open Questions

- Confirm preferred DB: current code uses GORM + MySQL driver.
- Confirm storage defaults (S3 vs local) for deployments.

## Architecture Decisions

- Use Fiber for HTTP routing and middleware; keep handlers thin.

## Session Notes

- After changing architecture or conventions, update the
  relevant context file and note the change here with rationale.
