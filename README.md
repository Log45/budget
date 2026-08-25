# Insight

Personal finance app for budgeting, loan simulation, expense tracking, and property management.

## Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) (for PostgreSQL)

## Quick start

1. **Start PostgreSQL**

   ```powershell
   docker compose up -d
   ```

2. **Configure environment**

   Copy the example env file and adjust values as needed:

   ```powershell
   copy .env.example .env
   ```

3. **Run the backend**

   ```powershell
   cd backend
   go run .
   ```

   The server listens on `http://localhost:8080`. Migrations run automatically on startup.

4. **Verify health**

   ```powershell
   curl http://localhost:8080/health
   ```

5. **Run the complete app with Docker**

   ```powershell
   docker compose up --build
   ```

   Open `http://localhost:3000`. The web proxy forwards browser requests under
   `/api` to the backend, so the frontend and API use one origin and do not need
   CORS for the containerized deployment.

## Environment variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `JWT_SECRET` | Yes | Secret key for signing JWT tokens |
| `PLAID_CLIENT_ID` | No | Plaid client ID (Phase 6) |
| `PLAID_SECRET` | No | Plaid secret (Phase 6) |
| `PLAID_ENV` | No | Plaid environment: `sandbox` or `production` |

## API endpoints

| Method | Endpoint | Status |
|--------|----------|--------|
| `GET` | `/health` | Done |
| `POST` | `/budgets` | Stub (Phase 4) |
| `GET` | `/budgets` | Stub (Phase 4) |
| `POST` | `/loans/payment` | Stub (Phase 3) |

See [docs/dev.md](docs/dev.md) for the full planned API and [TODO.md](TODO.md) for the implementation roadmap.

## Project structure

```
backend/
├── api/          HTTP handlers and routes
├── db/           PostgreSQL repositories and migrations
├── finance/      Pure calculation helpers (loans, interest, taxes)
├── models/       Domain structs
└── services/     Business logic
```

## Development

Build and test:

```powershell
cd backend
go build .
go test ./...
```
