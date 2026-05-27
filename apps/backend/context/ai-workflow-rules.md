# AI Workflow Rules (backend devs)

## Approach

Implement features in small, verifiable increments that map to
units already defined in `context/`. Follow the repository's
existing layering: `http` → `services` → `repositories`.

## Scoping Rules

- Work on one feature unit at a time (one API + related jobs +
  migrations).
- Prefer tiny, testable changes and add automated checks where
  possible (unit tests, integration smoke tests).
- Split cross-cutting changes (API + long-running infra) into
  separate commits and PRs.

## When to Split Work

- UI vs backend: backend changes only for this repo; frontend
  changes belong to frontend repo.
- Multiple unrelated routes or services should be separate PRs.

## Handling Missing Requirements

- Do not invent behaviour. If unclear, add an entry to
  `progress-tracker.md` under Open Questions and discuss.

## Protected Files

- Avoid changing generated files (e.g. gqlgen generated code)
  unless the change follows the generator workflow.

## Keeping Docs in Sync

Update context files when implementation changes architecture,
storage model, or code conventions.

## Before Moving On

1. Unit and integration checks pass locally.
2. `progress-tracker.md` updated with completed items.
3. Any required migrations or seed data are included.
