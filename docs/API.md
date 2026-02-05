# HTTP API

Base URL: `http://localhost:8080` (or your deployed host). JSON request and response bodies unless noted.

## Public

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness: database and ACA-Py issuer/verifier checks |
| GET | `/api/status` | Simple API banner |
| POST | `/api/login` | Start sign-in; returns QR payload and `proof_request_id` |
| POST | `/api/proof-callback` | Webhook for verified presentation (called by verifier agent / integration). If `PROOF_CALLBACK_SECRET` is set on the server, the same value must be sent in header `X-Webhook-Secret`. |
| GET | `/api/login/status/:proofRequestId` | Poll login completion; returns `pending`, `not_found`, or `completed` with token |

## Authenticated

Send header: `Authorization: Bearer <session_token>`.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/dashboard` | Current user profile |
| POST | `/api/logout` | Invalidate current session (204 No Content) |

## Credential administration (issuer)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/schemas` | Create schema via issuer agent |
| POST | `/api/credential-definitions` | Create credential definition via issuer agent |
