package handlers

import (
	"context"
	"net/http"
	"time"

	"bank-ai-chatbot/pkg/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	pool *pgxpool.Pool
}

func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	dbStatus := "ok"
	httpStatus := http.StatusOK

	pingCtx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.pool.Ping(pingCtx); err != nil {
		dbStatus = "down"
		status = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}

	response.JSON(w, httpStatus, map[string]string{
		"status":   status,
		"database": dbStatus,
	})
}
