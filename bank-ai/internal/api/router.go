package api

import (
	"net/http"

	"bank-ai-chatbot/internal/api/handlers"
	"bank-ai-chatbot/internal/api/middleware"
	"bank-ai-chatbot/internal/services"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	Health  *handlers.HealthHandler
	Auth    *handlers.AuthHandler
	Me      *handlers.MeHandler
	Chat    *handlers.ChatHandler
	Banking *handlers.BankingHandler
	AuthMW  *middleware.AuthMiddleware
}

func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", deps.Health.Health)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", deps.Auth.Register)
		r.Post("/login", deps.Auth.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(deps.AuthMW.RequireAuth)
		r.Get("/me", deps.Me.Me)

		r.Post("/chat", deps.Chat.CreateChat)
		r.Post("/chat/{chat_id}/message", deps.Chat.SendMessage)
		r.Get("/chat/{chat_id}/history", deps.Chat.GetHistory)

		r.Get("/banking/balance", deps.Banking.GetBalance)
		r.Get("/banking/transactions", deps.Banking.GetTransactions)
	})

	return r
}

func BuildDependencies(
	health *handlers.HealthHandler,
	authService *services.AuthService,
	chatService *services.ChatService,
	bankingService *services.BankingService,
	authMW *middleware.AuthMiddleware,
) Dependencies {
	return Dependencies{
		Health:  health,
		Auth:    handlers.NewAuthHandler(authService),
		Me:      handlers.NewMeHandler(authService),
		Chat:    handlers.NewChatHandler(chatService),
		Banking: handlers.NewBankingHandler(bankingService),
		AuthMW:  authMW,
	}
}
