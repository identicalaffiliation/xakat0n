package application

import (
	"context"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/ports"
)

func AddSeedData(ctx context.Context, itemsRepo ports.ItemsRepository, logger ports.Logger) {
	all, err := itemsRepo.GetAll(ctx)
	if err != nil {
		logger.Error(
			"failed to get all items",
			"error", err,
		)
	}

	if len(all) != 0 {
		return
	}

	items := []*domain.Item{
		domain.NewItem("Iphone 16 Pro", "Simple and strong", 80000),
		domain.NewItem("Компьютерный стол", "Ширина, Высота", 15000),
		domain.NewItem("Футбольный мяч", "Хорошее состояние", 1500),
		domain.NewItem("Трулик", "10/10", 2500),
		domain.NewItem("Аренда помещения под кебаб", "your best kebab place", 50000),
	}

	for _, item := range items {
		err := itemsRepo.CreateItem(ctx, item)
		if err != nil {
			logger.Error(
				"failed to create item",
				"error", err,
			)
		}
	}
}
