package application

import (
	"context"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/dto"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/logging"
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
	ctx = logging.WithUsername(ctx, u.logger, username.String())

	userID, err := u.users.GetOrCreate(ctx, username)
	if err != nil {
		ctx = u.logger.ContextFromError(ctx, err)
		u.logger.ErrorContext(ctx, "failed to get or create user", "error", err)
		return nil, domain.ErrInternal
	}
	ctx = logging.WithUserID(ctx, u.logger, userID)

	token, err := u.issuer.Issue(userID, username)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to issue token", "error", err)
		return nil, domain.ErrInternal
	}

	return dto.NewLoginResponse(userID, username.String(), token), nil
}
