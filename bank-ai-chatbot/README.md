# bank-ai-chatbot

## PostgreSQL Step 2 Quick Start

### 1) Set database URL

PowerShell:

`$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/bank_ai?sslmode=disable"`

### 2) Run migrations

Git Bash / WSL:

`sh scripts/migrate.sh`

PowerShell (no shell script required):

`Get-ChildItem migrations/*.sql | Sort-Object Name | ForEach-Object { psql $env:DATABASE_URL -v ON_ERROR_STOP=1 -f $_.FullName }`

### 3) Seed sample data

`go run ./scripts/seed_data.go`
