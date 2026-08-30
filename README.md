# MTGViewer V2

A collection management tool for Magic: The Gathering cards. Track your cards across containers, look up metadata, and manage transfers between containers.

This is a ground-up rebuild of the original MTGViewer, using a Go backend and a Vue 3 frontend.

If you've found this repository, go away!

## Tech Stack

**Backend**
- [Go](https://go.dev/) 1.25 with [Gin](https://gin-gonic.com/) for the HTTP API
- [PostgreSQL](https://www.postgresql.org/) via `lib/pq`
- [MTGJSON SDK](https://github.com/mtgjson/mtgjson-sdk-go) for card metadata and pricing data
- JWT-based authentication (`golang-jwt`)

**Frontend**
- [Vue 3](https://vuejs.org/) (Composition API, `<script setup>`) with TypeScript
- [Vue Router](https://router.vuejs.org/)
- [Vuetify 3](https://vuetifyjs.com/) for UI components
- [Vite](https://vitejs.dev/) for tooling/bundling

**Infrastructure**
- Multi-stage [Docker](https://www.docker.com/) build — one image serving both the compiled Go binary and the built Vue static assets

## Features

- **Card lookup** — search, browse, and view individual card details, including live pricing and metadata sourced from MTGJSON
- **Container management** — organize your collection into named containers (binders, boxes, decks, etc.), with per-container capacity limits
- **Import / withdraw** — add or remove cards from a container's inventory
- **Prune** — identify and clear out low-value cards from a container in bulk
- **Transaction log** — every collection change is recorded, so you can review what moved and when
- **User accounts** — signup/login with role-based access (`user` / `admin`)

## Project Structure

```
.
├── backend/
│   ├── auth/          # JWT auth, CORS middleware, role checks
│   ├── cards/          # Card lookups, pricing, MTGJSON SDK integration
│   ├── containers/     # Container CRUD, deposits/withdrawals, pruning
│   ├── database/       # DB connection + SQL schema/migrations
│   ├── routes/         # HTTP route registration per resource
│   ├── transactions/   # Transaction/change logging
│   └── users/          # User accounts, credentials
└── frontend/
    └── src/
        ├── components/  # Reusable UI components (nav bar, etc.)
        ├── fetch/       # Typed API client wrapper
        ├── plugins/     # Vue Router, Vuetify setup
        └── views/       # Page-level components
```

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 20+
- PostgreSQL
- Docker (optional, for containerized builds)

### Environment Variables

The backend reads its configuration from environment variables (a `.env` file is supported for local development via [`godotenv`](https://github.com/joho/godotenv)):

| Variable         | Description                                              |
|------------------|------------------------------------------------------------|
| `DB_HOST`        | Postgres host                                             |
| `DB_PORT`        | Postgres port                                              |
| `DB_NAME`        | Postgres database name                                     |
| `DB_USER`        | Postgres user                                               |
| `DB_PASSWORD`    | Postgres password                                            |
| `DB_SSLMODE`     | Postgres SSL mode (e.g. `disable`, `require`)                 |
| `TOKEN_KEY`      | Signing key used for issuing/validating JWTs                   |
| `CLIENT_ORIGINS` | Comma-separated list of allowed CORS origins (e.g. frontend URL) |

### Running locally

**Backend**
```bash
cd backend
go mod download
go run .
```
The API starts on `:8080`.

**Frontend**
```bash
cd frontend
npm install
npm run dev
```

### Database setup

Apply the schema in `backend/database/sql/schema.sql` to your Postgres instance before starting the backend for the first time.

### Running with Docker

A single multi-stage `Dockerfile` at the repo root builds both the frontend and backend into one image, serving the compiled Vue app directly from the Go server:

```bash
docker build -t mtgviewer-v2 .
docker run -p 8080:8080 --env-file .env mtgviewer-v2
```

## API Overview

All endpoints are served under `/api`.

| Resource       | Endpoints                                                                 |
|----------------|----------------------------------------------------------------------------|
| `/users`       | `POST /signup`, `POST /login`, `POST /logout`                              |
| `/cards`       | `GET /`, `GET /:card`, `GET /search`, `GET /random`, `POST /import`, `POST /withdraw`, `POST /oracle` |
| `/containers`  | `GET`, `GET /:container`, `GET /:container/cards`, `GET /prune`, `POST /prune` |
| `/transactions`| `GET`, `GET /:group`                                                       |
