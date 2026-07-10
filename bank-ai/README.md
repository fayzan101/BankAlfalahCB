# bank-ai

Banking AI chatbot backend (Go + PostgreSQL + OpenAI).

## Prerequisites

- Go 1.22+
- PostgreSQL
- OpenAI API key (for general chat assistant responses; banking queries work without it)
- `psql` CLI (for migrations)

## Quick Start

### 1) Set environment variables

PowerShell:

```powershell
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/bank_ai?sslmode=disable"
$env:JWT_SECRET="dev-secret-change-in-production"
$env:OPENAI_API_KEY="sk-your-openai-key"
```

### 2) Run migrations

PowerShell:

```powershell
cd bank-ai
Get-ChildItem migrations/*.sql | Sort-Object Name | ForEach-Object { psql $env:DATABASE_URL -v ON_ERROR_STOP=1 -f $_.FullName }
```

### 3) Seed sample data

```powershell
go run ./scripts/seed_data.go
```

Demo user password: `DemoPass123`

### 4) Start the server

```powershell
go run ./cmd/server/main.go
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Health + DB status |
| POST | `/auth/register` | No | Register new user |
| POST | `/auth/login` | No | Login and get JWT |
| GET | `/me` | JWT | Current user profile |
| POST | `/chat` | JWT | Create a new chat |
| POST | `/chat/{chat_id}/message` | JWT | Send message, get AI or banking reply |
| GET | `/chat/{chat_id}/history` | JWT | Fetch chat message history |
| GET | `/banking/balance` | JWT | Get account balance |
| GET | `/banking/transactions?limit=10` | JWT | Get recent transactions |

## Example chat flow

```powershell
# Login
$login = curl -s -X POST http://localhost:8080/auth/login `
  -H "Content-Type: application/json" `
  -d '{"email":"demo+...@bank.local","password":"DemoPass123"}' | ConvertFrom-Json
$token = $login.data.token

# Create chat
$chat = curl -s -X POST http://localhost:8080/chat `
  -H "Authorization: Bearer $token" `
  -H "Content-Type: application/json" `
  -d '{"title":"Account help"}' | ConvertFrom-Json
$chatId = $chat.data.chat.id

# Send general message
curl -s -X POST "http://localhost:8080/chat/$chatId/message" `
  -H "Authorization: Bearer $token" `
  -H "Content-Type: application/json" `
  -d '{"content":"What services do you offer?"}'

# Send banking intent (works without OpenAI)
curl -s -X POST "http://localhost:8080/chat/$chatId/message" `
  -H "Authorization: Bearer $token" `
  -H "Content-Type: application/json" `
  -d '{"content":"What is my balance?"}'

# Direct banking APIs
curl -s http://localhost:8080/banking/balance -H "Authorization: Bearer $token"
curl -s "http://localhost:8080/banking/transactions?limit=5" -H "Authorization: Bearer $token"

# Get history
curl -s "http://localhost:8080/chat/$chatId/history" `
  -H "Authorization: Bearer $token"
```

## Configuration

Default config: `configs/app.yaml`

Environment overrides:

- `CONFIG_PATH` — config file path
- `PORT` — server port
- `DATABASE_URL` — PostgreSQL connection string
- `JWT_SECRET` — JWT signing secret
- `JWT_EXPIRY` — token expiry (e.g. `24h`)
- `OPENAI_API_KEY` — OpenAI API key
- `OPENAI_MODEL` — model name (default `gpt-4o-mini`)
- `CORS_ALLOWED_ORIGINS` — comma-separated origins (default `*`)
- `RATE_LIMIT_IP_PER_MINUTE` — public IP rate limit (default `60`)
- `RATE_LIMIT_USER_PER_MINUTE` — authenticated user rate limit (default `120`)

## Operational features (Phase 6)

- Structured JSON request logs with method, path, status, latency, and `X-Request-ID`
- IP rate limiting on all routes; per-user rate limiting on authenticated routes (`429` when exceeded)
- Audit logs for register/login and banking balance/transaction reads
- Security headers (`X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, `Cache-Control`)
- CORS support for browser clients
