package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrConflict          = errors.New("resource already exists")
	ErrBadRequest        = errors.New("bad request")
	ErrInternal          = errors.New("internal server error")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrTooManyRequests    = errors.New("too many requests")
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func NotFound(message string) *AppError {
	return NewAppError("NOT_FOUND", message, ErrNotFound)
}

func Unauthorized(message string) *AppError {
	return NewAppError("UNAUTHORIZED", message, ErrUnauthorized)
}

func Forbidden(message string) *AppError {
	return NewAppError("FORBIDDEN", message, ErrForbidden)
}

func Conflict(message string) *AppError {
	return NewAppError("CONFLICT", message, ErrConflict)
}

func BadRequest(message string) *AppError {
	return NewAppError("BAD_REQUEST", message, ErrBadRequest)
}

func Internal(message string, err error) *AppError {
	return NewAppError("INTERNAL_ERROR", message, err)
}

func ServiceUnavailable(message string) *AppError {
	return NewAppError("SERVICE_UNAVAILABLE", message, ErrServiceUnavailable)
}

func TooManyRequests(message string) *AppError {
	return NewAppError("TOO_MANY_REQUESTS", message, ErrTooManyRequests)
}
