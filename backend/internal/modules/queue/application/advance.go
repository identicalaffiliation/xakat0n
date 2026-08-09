package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/queue/ports"
)

type AdvanceQueueUsecase struct {
	items     ports.ItemsRepository
	queue     ports.QueueRepository
	txManager ports.TxManager
	logger    ports.Logger
}

func NewAdvanceQueueUsecase(
	items ports.ItemsRepository,
	queue ports.QueueRepository,
	txManager ports.TxManager,
	logger ports.Logger,
) *AdvanceQueueUsecase {
	return &AdvanceQueueUsecase{
		items:     items,
		queue:     queue,
		txManager: txManager,
		logger:    logger,
	}
}

// AdvanceQueue — идемпотентная функция продвижения очереди товара: переводит
// просроченные права в EXPIRED, выдаёт освободившиеся места голове очереди
// и завершает очередь SOLD_OUT, если весь сток выкуплен. Целиком выполняется
// в одной транзакции; блокировка строки items — первая операция транзакции
// (единый порядок блокировок, см. architecture.md).
//
// Не возвращает *domain.Item: значение из этой (уже закоммиченной к моменту
// возврата) транзакции нельзя переиспользовать в последующей — каждая
// транзакция, трогающая items, обязана брать собственный LockStock.
func (u *AdvanceQueueUsecase) AdvanceQueue(ctx context.Context, itemID uuid.UUID, ttl time.Duration) error {
	err := u.txManager.WithTx(ctx, func(ctx context.Context) error {
		item, err := u.items.LockStock(ctx, itemID)
		if err != nil {
			return err
		}

		if err := u.queue.ExpireStale(ctx, itemID); err != nil {
			return err
		}

		taken, err := u.queue.CountTaken(ctx, itemID)
		if err != nil {
			return err
		}

		if freeSlots := item.Stock - taken; freeSlots > 0 {
			if _, err := u.queue.PromoteNext(ctx, itemID, freeSlots, ttl); err != nil {
				return err
			}
		}

		return u.queue.MarkSoldOut(ctx, itemID, item.Stock)
	})
	if err != nil {
		u.logger.Error(
			"failed to advance queue",
			"itemId", itemID,
			"error", err,
		)
		return err
	}

	return nil
}
