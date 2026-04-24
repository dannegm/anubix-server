# anubix-server

Backend for Anubix, a zero-knowledge password manager. All encryption and decryption happens on the client — the server never sees plaintext data or encryption keys.

## Requirements

- Go 1.22+
- PostgreSQL 15+
- [Air](https://github.com/air-verse/air) (optional, for hot-reload)

## Setup

```bash
cp .env.example .env
```

Edit `.env` and set your database connection:

```env
PORT=8080
APP_ENV=development
DB_URL=postgres://user:password@localhost:5432/anubix?sslmode=disable
```

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
