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
		return
	}

	if len(all) != 0 {
		return
	}

	items := []*domain.Item{
		// Сушёные насекомые
		newSeedItem(
			"Сушёные комары, 1 кг",
			"Собраны прошлым летом на даче, естественная сушка без химии. Идеально для коллекции или подкормки.",
			3500, "Сушёные насекомые", true, 1, "dried-mosquitoes.jpg",
		),
		newSeedItem(
			"Сушёные мухи, 1 кг",
			"Отборные экземпляры, без плесени и запаха. Штучный товар, повтора не будет.",
			3200, "Сушёные насекомые", true, 1, "dried-flies.jpg",
		),
		newSeedItem(
			"Сушёные тараканы, 1 кг",
			"Всегда в наличии, поставки не заканчиваются никогда. Оптом дешевле.",
			500, "Сушёные насекомые", false, 100, "dried-cockroaches.jpg",
		),

		// Пакет пакетов
		newSeedItem(
			"Пакет пакетов, ёмкость 4 пакета, семейная реликвия",
			"Достался от бабушки, полностью укомплектован. Не пакует, а хранит историю.",
			150, "Хозтовары", true, 1, "bag-of-bags-4.jpg",
		),
		newSeedItem(
			"Пакет пакетов, ёмкость 30 пакетов, коллекционный",
			"Собирался годами, пакеты разных эпох и магазинов. Раритет для истинных ценителей.",
			400, "Хозтовары", true, 1, "bag-of-bags-30.jpg",
		),
		newSeedItem(
			"Пакет пакетов, ёмкость 60 пакетов, топовый экземпляр",
			"Абсолютный рекорд плотности упаковки. Продаю только из-за переезда.",
			700, "Хозтовары", true, 1, "bag-of-bags-60.jpg",
		),
		newSeedItem(
			"Пакет пакетов обычный, свежий",
			"Новый, ещё пустой — набьёте сами. Всегда в наличии.",
			50, "Хозтовары", false, 50, "bag-of-bags-plain.jpg",
		),

		// Переходник с дизеля на бензин
		newSeedItem(
			"Переходник с дизеля на бензин, стандарт",
			"Кто не шарит, не мешайте вести бизнес. Последний в наличии.",
			800, "Автотовары", true, 1, "diesel-adapter-standard.jpg",
		),
		newSeedItem(
			"Переходник с дизеля на бензин, усиленный",
			"Кто не шарит, не мешайте вести бизнес. Для тех, кому мало одного.",
			1200, "Автотовары", true, 1, "diesel-adapter-reinforced.jpg",
		),

		// Участки на Луне — вставлены последними: GetAll сортирует по created_at DESC,
		// поэтому самая свежая по вставке категория оказывается первой в каталоге.
		newSeedItem(
			"Участок на Луне, Море Спокойствия, 1 га",
			"Историческое место посадки Apollo 11 неподалёку. Свидетельство о праве собственности прилагается.",
			4990, "Недвижимость на Луне", true, 1, "moon-sea-of-tranquility.jpg",
		),
		newSeedItem(
			"Участок на Луне, Море Дождей, 1 га",
			"Крупнейшее лунное море, вид на Землю в полный рост. Соседей не будет ещё лет 500.",
			4990, "Недвижимость на Луне", true, 1, "moon-sea-of-rains.jpg",
		),
		newSeedItem(
			"Участок на Луне, кратер Тихо, 1 га",
			"Один из самых заметных кратеров с Земли — будет видно даже без телескопа, кому показать.",
			5490, "Недвижимость на Луне", true, 1, "moon-tycho-crater.jpg",
		),
		newSeedItem(
			"Участок на Луне, Океан Бурь, 1 га",
			"Самое большое лунное море. Простор для будущей застройки.",
			4790, "Недвижимость на Луне", true, 1, "moon-ocean-of-storms.jpg",
		),
		newSeedItem(
			"Участок на Луне, Море Изобилия, 1 га",
			"Ровный рельеф, удобный подъезд с любой стороны — если долетите.",
			4990, "Недвижимость на Луне", true, 1, "moon-sea-of-fertility.jpg",
		),
		newSeedItem(
			"Участок на Луне, Море Ясности, 1 га",
			"Тихий спальный район лунной поверхности. Идеально для интровертов.",
			4690, "Недвижимость на Луне", true, 1, "moon-sea-of-clarity.jpg",
		),
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

func newSeedItem(
	title string,
	description string,
	price int64,
	category string,
	isLimited bool,
	stock int,
	imagePath string,
) *domain.Item {
	item := domain.NewItem(title, description, price, isLimited)
	item.Category = &category
	item.Stock = stock
	item.ImagePath = &imagePath

	return item
}
