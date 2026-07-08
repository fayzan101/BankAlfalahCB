package services

import (
	"context"
	"strings"

	"bank-ai-chatbot/internal/dto"
	"bank-ai-chatbot/internal/models"
	"bank-ai-chatbot/internal/repository/postgres"
	"bank-ai-chatbot/internal/security"
	"bank-ai-chatbot/internal/utils"
	apperrors "bank-ai-chatbot/pkg/errors"
	"github.com/google/uuid"
)

type AuthService struct {
	users  *postgres.UserRepository
	tokens *security.TokenManager
}

func NewAuthService(users *postgres.UserRepository, tokens *security.TokenManager) *AuthService {
	return &AuthService{
		users:  users,
		tokens: tokens,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	req.FullName = strings.TrimSpace(req.FullName)
	req.Email = utils.NormalizeEmail(req.Email)

	if err := utils.ValidateRegisterInput(req.FullName, req.Email, req.Password); err != nil {
		return nil, err
	}

	existing, err := s.users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.Internal("failed to check existing user", err)
	}
	if existing != nil {
		return nil, apperrors.Conflict("email already registered")
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.Internal("failed to process password", err)
	}

	user := &models.User{
		ID:           uuid.New(),
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hash,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to create user", err)
	}

	token, err := s.tokens.Generate(user.ID, user.Email)
	if err != nil {
		return nil, apperrors.Internal("failed to generate token", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserSummary{
			ID:       user.ID.String(),
			FullName: user.FullName,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	req.Email = utils.NormalizeEmail(req.Email)

	if err := utils.ValidateLoginInput(req.Email, req.Password); err != nil {
		return nil, err
	}

	user, err := s.users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.Internal("failed to lookup user", err)
	}
	if user == nil {
		return nil, apperrors.Unauthorized("invalid email or password")
	}

	if err := security.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, apperrors.Unauthorized("invalid email or password")
	}

	token, err := s.tokens.Generate(user.ID, user.Email)
	if err != nil {
		return nil, apperrors.Internal("failed to generate token", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserSummary{
			ID:       user.ID.String(),
			FullName: user.FullName,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*dto.UserSummary, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("failed to lookup user", err)
	}
	if user == nil {
		return nil, apperrors.NotFound("user not found")
	}

	return &dto.UserSummary{
		ID:       user.ID.String(),
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}