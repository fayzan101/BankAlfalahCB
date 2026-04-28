# Bank AI Chatbot - Final MVP Implementation

## 1) MVP Goal
Deliver a secure banking chatbot backend in Go that supports:
- User authentication (register/login with JWT)
- Authenticated chat endpoint with stored chat history
- OpenAI-powered assistant responses with banking-safe prompts
- Basic transaction inquiry support (read-only)
- Health checks, structured logging, and rate limiting

This MVP is backend-first and API-driven. UI can be added later.

## 2) In-Scope Features
- Auth APIs: register, login, token validation
- Chat APIs: create chat, send message, fetch chat history
- LLM integration: prompt orchestration + response generation
- Banking service: fetch account balance and latest transactions
- PostgreSQL persistence: users, chats, messages, transactions
- Security baseline: JWT auth middleware, password hashing, audit logs
- Operational baseline: Docker setup, migration script, seed script

## 3) Out-of-Scope for MVP
- Full RBAC admin console
- Real-time streaming responses
- Complex multi-channel integrations (WhatsApp, call center, etc.)
- Advanced fraud detection and compliance automation
- Production-grade RAG pipeline with vector database

## 4) Target Folder Structure
Use this structure as implementation source of truth:

- `cmd/server/main.go`
- `internal/config/*`
- `internal/api/*`
- `internal/services/*`
- `internal/repository/postgres/*`
- `internal/models/*`
- `internal/dto/*`
- `internal/ai/*`
- `internal/security/*`
- `internal/utils/*`
- `pkg/database/postgres.go`
- `pkg/errors/errors.go`
- `pkg/response/response.go`
- `migrations/*.sql`
- `scripts/*`
- `configs/*`
- `docs/*`
- `tests/*`
- `docker/*`

## 5) Implementation Sequence (Final MVP)

### Step 1: Bootstrap Project
- Initialize Go module and base app runner.
- Add config loader (`app.yaml` + environment overrides).
- Add server startup with router wiring.

Acceptance:
- App starts with `go run cmd/server/main.go`.
- `/health` returns `200 OK`.

### Step 2: Database + Migrations
- Implement DB connection utilities.
- Create migrations for users, chats, messages, transactions.
- Add migration script and seed script.

Acceptance:
- Fresh database can be migrated from zero.
- Seed script inserts sample users and transactions.

### Step 3: Auth Module
- Build `auth_service.go`, `auth_handler.go`, and `jwt.go`.
- Add password hashing (bcrypt) and JWT claims.
- Add auth middleware for protected APIs.

Acceptance:
- Register/login endpoints return JWT.
- Protected endpoint denies invalid/expired tokens.

### Step 4: Chat Core
- Build chat repositories and chat service.
- Implement chat handlers for create/send/history.
- Persist all prompts/responses in messages table.

Acceptance:
- User can create chat and retrieve message history.
- Data is isolated per authenticated user.

### Step 5: LLM Integration
- Add OpenAI client wrapper and prompt builder.
- Add system prompts for banking safety behavior.
- Integrate LLM call in chat service.

Acceptance:
- Chat endpoint returns assistant response from LLM.
- Failures return safe error response without leaking internals.

### Step 6: Banking Service Integration
- Implement banking service for read-only intent execution.
- Map simple intents: balance inquiry, latest transactions.
- Keep deterministic banking lookup before LLM free-text response.

Acceptance:
- "What is my balance?" returns account-derived value.
- "Show recent transactions" returns top recent records.

### Step 7: Security + Middleware Hardening
- Add request logging middleware and request IDs.
- Add rate limiting middleware per user/IP.
- Add audit logger for auth and banking actions.

Acceptance:
- Repeated abuse requests are throttled.
- Sensitive actions are auditable.

### Step 8: Docker + Documentation + Tests
- Finalize `Dockerfile` and `docker-compose.yml`.
- Write architecture/API/security docs.
- Add minimum unit + integration tests.

Acceptance:
- App runs via docker compose with DB.
- Core API test suite passes.

## 6) API Endpoints (MVP Contract)

### Public
- `GET /health`
- `POST /auth/register`
- `POST /auth/login`

### Protected (JWT)
- `POST /chat`
- `POST /chat/{chat_id}/message`
- `GET /chat/{chat_id}/history`
- `GET /banking/balance`
- `GET /banking/transactions?limit=10`

## 7) Data Model (Minimum)
- `users`: id, full_name, email, password_hash, created_at
- `chats`: id, user_id, title, created_at
- `messages`: id, chat_id, sender_type, content, created_at
- `transactions`: id, user_id, amount, type, description, created_at

## 8) Security Requirements (MVP Minimum)
- Passwords hashed with bcrypt.
- JWT with expiry and signature validation.
- No plain-text secrets in repo.
- Basic input validation on all DTOs.
- Do not expose internal stack traces in API responses.

## 9) Definition of Done (Final MVP)
MVP is complete when all below are true:
- Core endpoints implemented and tested
- JWT auth + middleware enforced on protected routes
- Chat + message persistence working in Postgres
- OpenAI response integrated via service abstraction
- Banking read-only queries functional and authorized
- Dockerized local setup and migration scripts available
- Essential docs (`architecture`, `api_spec`, `security_model`) finalized

## 10) Recommended Build Order by File Groups
1. `internal/models`, `internal/dto`
2. `pkg/database`, `internal/repository/postgres`
3. `internal/security`, `internal/services/auth_service.go`
4. `internal/api/handlers/auth_handler.go`, `internal/api/middleware/auth.go`
5. `internal/services/chat_service.go`, `internal/api/handlers/chat_handler.go`
6. `internal/ai/*`, `internal/services/llm_service.go`
7. `internal/services/banking_service.go`, transaction repository
8. middleware logger/rate limit + docs/tests/docker

## 11) Next Immediate Action
Start coding with:
1) migrations + DB connection
2) auth service/handler + JWT middleware
3) chat handler/service/repo

These three deliver a usable vertical slice fastest.
