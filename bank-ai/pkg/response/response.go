package response

import (
	"encoding/json"
	"errors"
	"net/http"

	apperrors "bank-ai-chatbot/pkg/errors"
)

type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Success: status >= 200 && status < 300,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, err error) {
	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		appErr = apperrors.Internal("something went wrong", err)
	}

	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, apperrors.ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(err, apperrors.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, apperrors.ErrConflict):
		status = http.StatusConflict
	case errors.Is(err, apperrors.ErrBadRequest):
		status = http.StatusBadRequest
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		status = http.StatusServiceUnavailable
	case errors.Is(err, apperrors.ErrTooManyRequests):
		status = http.StatusTooManyRequests
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Success: false,
		Error: &ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}
