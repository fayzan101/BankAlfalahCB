# bank-ai

Banking AI chatbot backend (Go + PostgreSQL + OpenAI).

## Prerequisites

- Go 1.22+
- PostgreSQL
- OpenAI API key (for chat assistant responses)
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
| POST | `/chat/{chat_id}/message` | JWT | Send message, get AI reply |
| GET | `/chat/{chat_id}/history` | JWT | Fetch chat message history |

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

# Send message
curl -s -X POST "http://localhost:8080/chat/$chatId/message" `
  -H "Authorization: Bearer $token" `
  -H "Content-Type: application/json" `
  -d '{"content":"What services do you offer?"}'

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
