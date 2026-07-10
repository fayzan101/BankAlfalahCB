package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bank-ai-chatbot/internal/ai"
	"bank-ai-chatbot/internal/api"
	"bank-ai-chatbot/internal/api/handlers"
	"bank-ai-chatbot/internal/api/middleware"
	"bank-ai-chatbot/internal/audit"
	"bank-ai-chatbot/internal/config"
	"bank-ai-chatbot/internal/repository/postgres"
	"bank-ai-chatbot/internal/security"
	"bank-ai-chatbot/internal/services"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/app.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if cfg.Database.URL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if err := os.Setenv("DATABASE_URL", cfg.Database.URL); err != nil {
		log.Fatalf("set DATABASE_URL: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	auditLogger := audit.NewLogger(logger)

	ctx := context.Background()
	db, err := postgres.NewDB(ctx)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	chatRepo := postgres.NewChatRepository(db)
	messageRepo := postgres.NewMessageRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)

	tokenManager := security.NewTokenManager(cfg.JWT.Secret, cfg.JWT.Expiry)
	authService := services.NewAuthService(userRepo, tokenManager)
	bankingService := services.NewBankingService(transactionRepo)
	authMW := middleware.NewAuthMiddleware(tokenManager)

	llmEnabled := cfg.OpenAI.APIKey != ""
	var llmService *services.LLMService
	if llmEnabled {
		openaiClient := ai.NewClient(
			cfg.OpenAI.APIKey,
			cfg.OpenAI.Model,
			cfg.OpenAI.MaxTokens,
			cfg.OpenAI.Timeout,
		)
		llmService = services.NewLLMService(openaiClient, cfg.OpenAI.MaxHistoryMessages)
		logger.Info("openai integration enabled", "model", cfg.OpenAI.Model)
	} else {
		logger.Warn("OPENAI_API_KEY not set; general chat messages will return service unavailable")
	}

	chatService := services.NewChatService(chatRepo, messageRepo, bankingService, llmService, llmEnabled)

	healthHandler := handlers.NewHealthHandler(db.Pool)
	deps := api.BuildDependencies(cfg, logger, auditLogger, healthHandler, authService, chatService, bankingService, authMW)
	router := api.NewRouter(deps)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout + cfg.OpenAI.Timeout,
	}

	go func() {
		logger.Info("server listening", "addr", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}
