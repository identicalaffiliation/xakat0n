package checkout

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/checkout/domain"
	queueapplication "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/application"
	queuedomain "github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/domain"
)

// AdvanceAdapter переводит queue/domain.ErrItemNotFound в checkout/domain.ErrItemNotFound на границе модулей.
// Экспортирован, чтобы integration-тесты могли собрать его напрямую с реальным AdvanceQueueUsecase.
type AdvanceAdapter struct {
	usecase *queueapplication.AdvanceQueueUsecase
}

func NewAdvanceAdapter(usecase *queueapplication.AdvanceQueueUsecase) *AdvanceAdapter {
	return &AdvanceAdapter{usecase: usecase}
}

func (a *AdvanceAdapter) AdvanceQueue(ctx context.Context, itemID uuid.UUID, ttl time.Duration) error {
	err := a.usecase.AdvanceQueue(ctx, itemID, ttl)
	if err != nil {
		if errors.Is(err, queuedomain.ErrItemNotFound) {
			return domain.ErrItemNotFound
		}

		return err
	}

	return nil
}
