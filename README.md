# World Cup Predict Backend

Go backend for a World Cup prediction app with channel-based competition, JWT auth, and admin-managed events.

## Requirements

- Go 1.22+
- PostgreSQL

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `JWT_SECRET` | Yes | Secret for signing JWT access tokens |
| `PORT` | No | HTTP port (default: `8080`) |
| `ADMIN_BOOTSTRAP_NAME` | No | First admin username (created once if no admin exists) |
| `ADMIN_BOOTSTRAP_PASSWORD` | No | First admin password |

## Local Development

```bash
# Copy example env and edit values
cp .env.example .env

# Run server (migrations run automatically on startup)
go run .
```

Config is loaded from `.env` if present. On Railway, real environment variables are used directly (`.env` is optional).

Health check: `GET /health`

Swagger UI: `GET /swagger/` — server URL is detected automatically from the current host (localhost or Railway).

User API DTO & response dokümantasyonu: [`docs/user-api.md`](docs/user-api.md)

## Seed WC2026 Events

Load group-stage matches and placement events (champion, runner-up, third place):

```bash
go run ./cmd/seed
```

- 3 placement events — deadline: **2026-06-10T19:00:00Z** (1 day before opening match, GMT)
- 72 group matches — deadline: **1 hour before kickoff** (all times stored in UTC/GMT)
- All 48 participating nations included as team options for placement events

Re-run with `-force` to replace existing WC2026 events (e.g. after team name updates).

## API Overview

### Auth (public)
- `POST /api/v1/auth/register` — `{ "name", "password", "channel_code" }`
- `POST /api/v1/auth/login` — `{ "name", "password", "channel_code" }` (admin can omit `channel_code`)

### User (Bearer token)
- `GET /api/v1/me`
- `PATCH /api/v1/me/nickname` — görünen isim güncelle
- `GET /api/v1/leaderboard`
- `GET /api/v1/events?status=open|locked|pending|completed`
- `GET /api/v1/events/{id}`
- `PUT /api/v1/events/{id}/prediction` — `{ "choice": {...} }`
- `GET /api/v1/events/{id}/predictions` — visible after deadline

### Admin (Bearer token + admin role)
- `POST /api/v1/admin/channels` — `{ "code", "name" }`
- `GET /api/v1/admin/channels`
- `POST /api/v1/admin/events` — `{ "type", "title", "metadata", "deadline" }`
- `PATCH /api/v1/admin/events/{id}`
- `POST /api/v1/admin/events/{id}/result` — `{ "result": {...} }`
- `POST /api/v1/admin/events/calculate-scores` — sonucu girilmiş tüm eventlerin puanlarını hesapla
- `POST /api/v1/admin/events/{id}/calculate-scores` — tek event puan hesapla

## Scoring Rules

| Event type | Tahmin (`choice`) | Sonuç (`result`) | Puan |
|------------|-------------------|------------------|------|
| `match_score` | `{"home_score": 2, "away_score": 1}` | aynı format | 1 (sonuç) / 3 (tam skor) |
| `champion` | `{"team": "Brezilya"}` | `{"team": "Brezilya"}` | 10 |
| `runner_up` | `{"team": "Arjantin"}` | `{"team": "Arjantin"}` | 5 |
| `third_place` | `{"team": "Fransa"}` | `{"team": "Fransa"}` | 3 |

## Prediction Examples

**Match score choice:**
```json
{ "home_score": 2, "away_score": 1 }
```

**Champion / runner_up / third_place choice (tek takım seçimi):**
```json
{ "team": "Brezilya" }
```

## Railway Deployment

Railway auto-detects Go projects via Nixpacks. No Dockerfile required.

1. Create a Railway project and add **PostgreSQL**
2. Add a service from this GitHub repo
3. Set environment variables:
   - `JWT_SECRET` — generate a strong random string
   - `ADMIN_BOOTSTRAP_NAME` / `ADMIN_BOOTSTRAP_PASSWORD` — for first admin
   - `DATABASE_URL` — linked automatically from PostgreSQL addon
   - `PORT` — set automatically by Railway
4. Configure build/start (optional, [`railway.toml`](railway.toml) included):
   - **Build:** `go build -o server .`
   - **Start:** `./server`

Migrations run automatically on each server startup via embedded goose migrations.

## Typical Flow

1. Admin logs in (no channel code) and creates channels
2. Users register with a channel code
3. Admin creates events with deadlines
4. Users submit predictions before deadline
5. After deadline, users can see channel predictions
6. Admin enters results and triggers score calculation
7. Users view leaderboard in their channel
