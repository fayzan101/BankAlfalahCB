package utils

import (
	"net/mail"
	"strings"
	"unicode/utf8"

	apperrors "bank-ai-chatbot/pkg/errors"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
)

func ValidateRegisterInput(fullName, email, password string) error {
	fullName = strings.TrimSpace(fullName)
	email = strings.TrimSpace(strings.ToLower(email))

	if fullName == "" {
		return apperrors.BadRequest("full name is required")
	}
	if len(fullName) > 100 {
		return apperrors.BadRequest("full name must be at most 100 characters")
	}
	if email == "" {
		return apperrors.BadRequest("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return apperrors.BadRequest("invalid email address")
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	return nil
}

func ValidateLoginInput(email, password string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return apperrors.BadRequest("email is required")
	}
	if password == "" {
		return apperrors.BadRequest("password is required")
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return apperrors.BadRequest("password is required")
	}
	if utf8.RuneCountInString(password) < MinPasswordLength {
		return apperrors.BadRequest("password must be at least 8 characters")
	}
	if len(password) > MaxPasswordLength {
		return apperrors.BadRequest("password must be at most 72 characters")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
