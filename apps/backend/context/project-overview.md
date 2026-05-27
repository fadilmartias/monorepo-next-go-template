# Backend API — monorepo-next-go-template

## Overview

Backend service written in Go (Fiber) that provides REST/GraphQL APIs
and background jobs for the monorepo Next + Go template. Bertanggung
jawab untuk autentikasi, manajemen user, konten (article, banner),
pembayaran, dan integrasi layanan eksternal (S3, Redis, email).

## Primary Goals

1. Stabil: API yang mudah diuji dan dapat di-deploy sebagai layanan
   mandiri.
2. Konsisten: Ikuti pola `controllers → services → repositories`
   yang ada di kodebase.
3. Observable: Logging terstruktur (zap), monitoring, dan metrics.

## Core Flows (ringkas)

1. User registers / authenticates (JWT / passkeys).
2. User membuat resource (article, order) melalui API.
3. Background job memproses task (asynq) untuk pembayaran / notifikasi.
4. File uploads disimpan di S3 (atau storage terkonfigurasi).

## Features

- REST & GraphQL endpoints (gqlgen + Fiber route handlers)
- Background workers via `asynq` dan cron jobs
- Payment gateway integrations (Midtrans, Ipaymu, Digiflazz, dll.)
- Email templates & transactional email
- File uploads (S3 or local fallback)

## Scope

In scope:
- Backend APIs, background jobs, integrations, migrations, seeds.

Out of scope:
- Frontend UI implementation (Next app lives in another workspace folder).

## Success Criteria

1. API endpoints implement usecase contracts in `usecases/` and pass
   basic integration checks.
2. Background tasks enqueue/run reliably in dev and staging.
3. Repository and service layers follow existing patterns in codebase.
