# anubix-server

Backend for Anubix, a zero-knowledge password manager. All encryption and decryption happens on the client — the server never sees plaintext data or encryption keys.

## Requirements

- Go 1.22+
- PostgreSQL 15+
- [Air](https://github.com/air-verse/air) (optional, for hot-reload)
- Docker & Docker Compose (optional, for local DB)

## Setup

```bash
cp .env.example .env
```

Edit `.env`:

```env
PORT=8080
APP_ENV=development
DB_URL=postgres://user:password@localhost:5432/anubix?sslmode=disable
```

## Database

**Option A — Docker (recommended for local dev):**

```bash
docker compose up db -d
```

This spins up a PostgreSQL instance at `localhost:5432` using the credentials defined in `docker-compose.yml`.

**Option B — Local or remote Postgres:**

Point `DB_URL` in your `.env` to your instance.

> Tables are created automatically on server startup (auto-migrate). No manual migrations needed.

## Run

**Development (with hot-reload):**

```bash
air
```

**Without hot-reload:**

```bash
go run ./cmd/api
```

**Build:**

```bash
go build -o ./bin/server ./cmd/api
./bin/server
```

## Docker

```bash
docker compose up --build
```

## Health check

```
GET /health → { "status": "ok" }
```
