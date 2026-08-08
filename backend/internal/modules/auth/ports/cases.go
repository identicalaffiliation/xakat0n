package ports

import (
	"context"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/dto"
)

type LoginUsecase interface {
	Login(ctx context.Context, rawUsername string) (*dto.LoginResponse, error)
}
