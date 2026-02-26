# Contributing

## Run locally

- **Full stack:** `docker compose up` from the repository root (see README for ports and agents).
- **Backend only:** Go 1.21+, PostgreSQL running, then from `backend/` run `go run .` with `DATABASE_URL` and other env vars set (see below).
- **Frontend only:** from `frontend/` run `npm install` and `npm start` (expects API at `REACT_APP_API_URL` or the dev proxy).

## Environment variables (backend)

| Variable | Purpose |
|----------|---------|
| `PORT` | HTTP listen port (default `8080`) |
| `DATABASE_URL` | PostgreSQL DSN |
| `ISSUER_AGENT_URL` | ACA-Py issuer admin API |
| `VERIFIER_AGENT_URL` | ACA-Py verifier admin API |
| `LEDGER_URL` | Indy / genesis-related base URL |
| `CREDENTIAL_DEFINITION_ID` | Required for `/api/login` proof requests |
| `CALLBACK_URL` | Default proof callback URL if not passed as query param |
| `VERIFIER_ENDPOINT` | Wallet-reachable verifier base URL (QR / OOB) |
| `CORS_ALLOW_ORIGINS` | Comma-separated allowed origins; unset allows permissive CORS for development |
| `PROOF_CALLBACK_SECRET` | If set, `POST /api/proof-callback` requires matching `X-Webhook-Secret` header |

## Commits

Use [Conventional Commits](https://www.conventionalcommits.org/) style when possible, for example:

- `fix:` bug fixes
- `feat:` new behavior
- `refactor:`, `perf:`, `test:`, `docs:`, `chore:`, `ci:` as appropriate

## Checks

From `backend/`:

```bash
golangci-lint run ./...
go test ./...
go vet ./...
```

From `frontend/`:

```bash
npm run build
```
