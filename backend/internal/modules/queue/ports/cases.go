package ports

import (
	"context"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/dto"
)

type CreateUsecase interface {
	CreateQueue(ctx context.Context, in *dto.CreateQueueRequest) (*dto.QueueResponse, error)
}

type QuitUsecase interface {
	QuitQueue(ctx context.Context, in *dto.QuitQueueRequest) (*dto.QueueResponse, error)
}
