package ports

import (
	"context"

	"github.com/identicalaffiliation/xakat0n/backend/internal/dto"
)

type CreateUsecase interface {
	CreateQueue(ctx context.Context, in *dto.CreateQueueRequest) (*dto.QueueResponse, error)
	QuiteQueue(ctx context.Context, in *dto.QuiteQueueRequest) (*dto.QueueResponse, error)
}
