# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Faraway is a full-stack multiplayer scoreboard app — track scores across players and rounds for any game. Go backend + Vue 3 frontend, PostgreSQL storage.

## Commands

### Backend (`api/`)

```bash
cd api
go test ./...          # Run all tests
go test ./handler/...  # Run tests in a single package
golangci-lint run      # Lint
go build -o server .   # Build binary
go run main.go         # Run locally (requires DB)
```

### Frontend (`front/`)

```bash
cd front
npm install
npm run dev        # Dev server at http://localhost:5173 (proxies /api to :8080)
npm run build      # Production build → dist/
npm run typecheck  # vue-tsc type check (no emit)
npm run lint       # ESLint with typescript-eslint + eslint-plugin-vue
npm run test       # Vitest (jsdom environment)
```

### Full Stack

```bash
docker-compose up  # PostgreSQL on :5432 + API on :8080
```

The frontend dev server proxies `/api` → `http://localhost:8080`, so you can run `docker-compose up` for the backend and `npm run dev` for the frontend simultaneously.

## Architecture

### Backend (Clean Layered — `api/`)

Strict dependency flow: `Handler → Service → Repository → GORM → PostgreSQL`

- **`models/`** — GORM domain entities: `Game`, `Player`, `Round`, `RoundScore`. All embed a base model with timestamps and soft deletes.
- **`repository/`** — Interface-based data access (one interface + implementation per entity). Abstracts GORM from business logic, enabling test mocks.
- **`service/`** — Business rules live here. Key constraints enforced: max 8 players per game, no duplicate player names per game.
- **`handler/`** — Gin HTTP handlers. Parse/validate request → call service → return JSON.
- **`router/`** — Wires all routes. `main.go` instantiates repos → services → handlers → router in dependency order.
- **`config/`** — Viper-based config loaded from `.env` or environment variables.
- **`db/`** — DB connection + auto-migration on startup.

### Frontend (Vue 3 SPA — `front/src/`)

**Plain JavaScript is forbidden. All frontend code must be TypeScript.** Use explicit types wherever inference is insufficient — prefer typed over `any`.

- **`types/api.ts`** — Canonical API types (`Game`, `Player`, `Round`, `RoundScore`) derived from `openapi.yaml`. Update here when the API changes.
- **`api/`** — Service modules (`games.ts`, `players.ts`, `rounds.ts`, `scores.ts`) wrapping a shared `client.ts`. Components never call `fetch` directly.
- **`views/`** — Two pages: `GamesView` (list/create games) and `GameView` (interactive scoreboard).
- **`components/`** — `GameCard`, `CreateGame` modal, plus inline score editing in `GameView`.
- **`router.ts`** — Vue Router with two routes: `/` and `/games/:id`.
- All Vue files use `<script setup lang="ts">`.
- Tailwind CSS for styling; Lucide Vue for icons.

### Data Model

```
Game → has many → Players (max 8, unique names)
Game → has many → Rounds
Round → has many → RoundScores
RoundScore → belongs to → Player + Round (score per player per round)
```

### CI/CD (`.github/workflows/ci.yml`)

Pipeline: lint → test → build multi-platform Docker images (amd64 + arm64) → push to GHCR. Images are built for both platforms in parallel and merged into a manifest list.

## Key Files

- `openapi.yaml` — Full API specification (source of truth for endpoints)
- `api/main.go` — Dependency wiring entry point
- `front/vite.config.ts` — Dev proxy config (`/api` → `:8080`)
