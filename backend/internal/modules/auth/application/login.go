package application

import (
	"context"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/ports"
)

type LoginUsecase struct {
	users  ports.UserRepository
	issuer ports.TokenIssuer
	logger ports.Logger
}

func NewLoginUsecase(users ports.UserRepository, issuer ports.TokenIssuer, logger ports.Logger) *LoginUsecase {
	return &LoginUsecase{
		users:  users,
		issuer: issuer,
		logger: logger,
	}
}

func (u *LoginUsecase) Login(ctx context.Context, rawUsername string) (*dto.LoginResponse, error) {
	username, err := domain.NewUsername(rawUsername)
	if err != nil {
		return nil, err
	}

	userID, err := u.users.GetOrCreate(ctx, username)
	if err != nil {
		u.logger.Error("failed to get or create user", "username", username.String(), "error", err)
		return nil, domain.ErrInternal
	}

	token, err := u.issuer.Issue(userID, username)
	if err != nil {
		u.logger.Error("failed to issue token", "user_id", userID.String(), "error", err)
		return nil, domain.ErrInternal
	}

	return dto.NewLoginResponse(userID, username.String(), token), nil
}
