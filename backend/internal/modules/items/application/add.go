package application

import (
	"context"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/domain"
	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/items/ports"
)

func AddSeedData(ctx context.Context, itemsRepo ports.ItemsRepository, logger ports.Logger) {
	all, err := itemsRepo.GetAll(ctx)
	if err != nil {
		ctx = logger.ContextFromError(ctx, err)
		logger.ErrorContext(ctx,
			"failed to get all items",
			"error", err,
		)
		return
	}

	if len(all) != 0 {
		return
	}

	items := []*domain.Item{
		newSeedItem(
			"iPhone 16 Pro",
			"Флагманский смартфон Apple с титановым корпусом и экраном 6,3 дюйма.",
			80000,
			"Электроника",
			true,
			1,
		),
		newSeedItem(
			"Компьютерный стол",
			"Письменный стол с широкой столешницей и отделением для системного блока.",
			15000,
			"Мебель",
			true,
			3,
		),
		newSeedItem(
			"Футбольный мяч",
			"Тренировочный футбольный мяч стандартного пятого размера.",
			1500,
			"Спорт",
			false,
			20,
		),
		newSeedItem(
			"Турник настенный",
			"Стальной настенный турник для домашних тренировок, допустимая нагрузка до 150 кг.",
			2500,
			"Спорт",
			true,
			5,
		),
		newSeedItem(
			"Аренда помещения под кебаб",
			"Помещение с отдельным входом и вытяжкой, подготовленное для точки общественного питания.",
			50000,
			"Недвижимость",
			false,
			1,
		),
	}

	for _, item := range items {
		itemCtx := ctx
		itemCtx = logger.WithField(itemCtx, "itemTitle", item.Title)
		err := itemsRepo.CreateItem(itemCtx, item)
		if err != nil {
			itemCtx = logger.ContextFromError(itemCtx, err)
			logger.ErrorContext(itemCtx,
				"failed to create item",
				"error", err,
			)
		}
	}
}

func newSeedItem(
	title string,
	description string,
	price int64,
	category string,
	isLimited bool,
	stock int,
) *domain.Item {
	item := domain.NewItem(title, description, price, isLimited)
	item.Category = &category
	item.Stock = stock

	return item
}
