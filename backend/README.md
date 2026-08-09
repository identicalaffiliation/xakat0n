# Backend

Модульный монолит на Go. Обоснование архитектурных решений (БД, атомарность, polling vs
WebSocket и т.д.) — в `../docs/architecture.md`. Этот файл — только про то, как физически
устроен код и куда класть новый код.

## Структура

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # composition root: конфиг, логгер, пул, tx.Manager,
│                                 # конструирует модули, регистрирует их роуты, graceful shutdown
│
├── internal/
│   ├── modules/                 # один каталог — один модуль предметной области
│   │   └── queue/
│   │       ├── domain/          # сущности и бизнес-правила, ничего не знает о HTTP/SQL
│   │       ├── ports/           # интерфейсы, которые нужны application от внешнего мира
│   │       │                    # (репозиторий, TxManager, Logger) — application зависит
│   │       │                    # от этих интерфейсов, а не от конкретной реализации
│   │       ├── dto/             # структуры запроса/ответа usecase'ов (между application и presentation)
│   │       ├── application/     # usecase'ы: оркестрируют domain + ports, транзакционные границы
│   │       ├── infrastructure/
│   │       │   └── postgres/    # реализация ports.*Repository поверх pgx
│   │       ├── presentation/
│   │       │   └── http/        # HTTP-хендлеры: разбирают запрос, зовут usecase, сериализуют ответ
│   │       └── module.go        # собирает репозиторий -> usecase -> хендлер, регистрирует роуты
│   │
│   └── shared/                  # инфраструктура, нужная любому модулю — не бизнес-логика
│       ├── config/              # чтение конфига (yaml + env через cleanenv)
│       ├── postgres/            # пул соединений (pgxpool), IsUniqueViolation
│       ├── tx/                  # DBTX-интерфейс, Manager.WithTx, DBTXFromContext
│       ├── logger/              # Logger-интерфейс + реализация на log/slog
│       ├── httpserver/          # голый http.Server + chi.Router с общими middleware,
│       │                        # ничего не знает о роутах конкретных модулей
│       └── httpx/               # SessionAuth (мидлварь авторизации), EncodeJSON
│
├── tests/
│   └── integrations/            # интеграционные тесты на реальном Postgres (testcontainers)
│
├── Dockerfile
└── go.mod
```

Миграции — в отдельном Go-модуле `../migrator` (свой `go.mod`, свой контейнер в
`docker-compose.yaml`), не в `backend`.

## Слои модуля и направление зависимостей

```
presentation/http  →  application  →  ports  ←  infrastructure/postgres
                            ↓
                          domain
```

Правило одно: зависимости смотрят только внутрь, к `domain`. `application` не импортирует
`presentation` и не импортирует `infrastructure` напрямую — только интерфейсы из `ports`.
Конкретную реализацию (`infrastructure/postgres.QueueRepository`) в `application` подставляет
`module.go` при сборке. Это даёт возможность подменить Postgres на что угодно другое, не трогая
бизнес-логику, и юнит-тестировать `application` без БД (моком порта) — на практике в этом проекте
пока тестируем через реальный Postgres в `tests/integrations`, но интерфейс это не блокирует.

`internal/shared/*` — не модуль и не подчиняется этому правилу: это то, что нужно сразу всем
модулям (соединение с БД, конфиг, логгер, транзакции, HTTP-обвязка). Модуль может свободно
импортировать `internal/shared/*`. Модуль **не должен** импортировать другой модуль напрямую,
кроме его `ports` — например, когда `queue` будет использовать `items` для проверки остатка,
`queue/application` будет зависеть от `items/ports.ItemsRepository`, а не от
`items/infrastructure`.

## Как добавить новую ручку в существующий модуль (`queue`)

На примере: `queue` уже содержит `POST /products/{productId}/queue`. Пусть нужно добавить,
например, `GET /products/{productId}/queue/me`.

1. **`domain/`** — если нужны новые правила/поля на сущности, добавляй сюда. Для чтения статуса
   обычно ничего нового не требуется.
2. **`ports/`** — опиши, что usecase'у нужно от репозитория (`repo.go`) — например,
   `GetByProductAndUser(ctx, productID, userID) (*domain.Queue, error)`.
3. **`infrastructure/postgres/`** — реализуй этот метод в `QueueRepository` (SQL-запрос).
4. **`dto/`** — опиши форму ответа, если она отличается от того, что уже есть.
5. **`application/`** — новый usecase (или метод на существующем): принимает `ports.QueueRepository`
   и остальные порты через конструктор, реализует бизнес-логику поверх интерфейсов.
6. **`presentation/http/`** — новый хендлер `http.HandlerFunc`, парсит path/заголовки, зовёт
   usecase, сериализует ответ через `httpx.EncodeJSON`. Достаёт `user_id`, положенный
   `httpx.SessionAuth`, через `httpx.UserID(r.Context())` — сам заголовок не парсит.
7. **`module.go`** — зарегистрируй роут: `r.With(httpx.SessionAuth).Get(route, handler.GetMyQueue(usecase))`.
   Если новому usecase'у нужны новые зависимости — прокинь их через сигнатуру `queue.New(...)`
   и через вызов в `cmd/api/main.go`.
8. Тесты: юнит — рядом с кодом (`_test.go` в том же пакете); интеграционные, если метод трогает
   БД — в `tests/integrations/`, с реальным Postgres через testcontainers (см. `main_test.go`,
   `helpers.go`).

Ничего за пределами `internal/modules/queue/` трогать не должно понадобиться — `shared/*` уже
есть и переиспользуется.

## Как завести новый модуль (например, `items`)

1. Создай `internal/modules/items/` с тем же набором подпакетов, что у `queue` (только те, что
   реально нужны — если у модуля пока нет HTTP-ручек, `presentation/` и `dto/` не заводи).
2. `domain/` — сущности модуля (`Item`), его собственные ошибки (`ErrItemNotFound`).
3. `ports/` — что нужно от инфраструктуры и от кого-то ещё, если он вызывается снаружи (у `items`
   в первую очередь — `ItemsRepository` с методом `LockStock`, который зовёт `queue`).
4. `infrastructure/postgres/` — реализация портов поверх `internal/shared/tx.DBTX` /
   `internal/shared/postgres`, ровно как у `queue`.
5. Если у модуля есть свои HTTP-ручки — `application/`, `dto/`, `presentation/http/`, `module.go`
   по тому же образцу, что у `queue`.
6. Подключи модуль в `cmd/api/main.go`: сконструируй его (`items.New(...)`) и, если есть роуты —
   зарегистрируй их на общем `router := httpserver.NewRouter()` до `httpserver.New(cfg, router)`.
7. Если модулю нужна новая таблица — миграция в `../migrator/migrations/`, отдельным файлом,
   goose-таймстамп в имени (не переиспользуй и не переименовывай существующие миграции — они
   могут быть уже применены на чьей-то локальной БД, переименование ломает версионирование
   goose).

## Как использовать логгер

Логгер построен на `log/slog`. Поля накапливаются в `context.Context`, а
`logger.HandlerMiddleware` автоматически добавляет их в каждую запись, сделанную через
`DebugContext`, `InfoContext`, `WarnContext` или `ErrorContext`.

### Инициализация

В `cmd/api/main.go` создаётся один объект `logger.Logging` на всё приложение. Внутри него
находятся `*slog.Logger` и `*logctx.LogCtx`:

```go
slogger, err := logger.NewLogger(cfg)
if err != nil {
	log.Fatal(err)
}
```

В конструкторы модулей передаётся только этот объект:

```go
items := itemsModule.New(pool, slogger)
queue := queueModule.New(pool, txManager, slogger, cfg.CheckoutTimer)
```

Модули не зависят от конкретной структуры `logger.Logging`. Они объявляют единый интерфейс
`ports.Logger`, который содержит и методы записи, и методы работы с log context:

```go
type Logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)

	WithField(ctx context.Context, key string, value any) context.Context
	WrapError(ctx context.Context, err error) error
	ContextFromError(ctx context.Context, err error) context.Context
}
```

### Модуль `logctx`

- **Поля добавляются один раз.** Когда `userID`, `itemID` или другой параметр становится
  известен, он сохраняется в контексте и автоматически доступен всем нижележащим вызовам.
- **Бизнес-логика не собирает поля перед каждым логом.** `HandlerMiddleware` самостоятельно
  извлекает их из контекста, поэтому вызову `ErrorContext` достаточно передать сообщение и
  ошибку.
- **Контекст можно поднять наверх.** Обычный `context.Context` распространяется только вниз по
  стеку. `WrapError` сохраняет поля внутри ошибки, а `ContextFromError` восстанавливает их на
  верхнем слое, где выполняется логирование.
- **Ошибку достаточно логировать один раз.** Repository и другие внутренние слои возвращают
  обёрнутую ошибку, а верхний слой создаёт единственную запись со всеми накопленными полями.
- **Сохраняется стандартная работа с ошибками Go.** Обёртка реализует `Unwrap`, поэтому после
  добавления log context продолжают работать `errors.Is` и `errors.As`.
- **Исходный контекст не изменяется.** `WithField` копирует карту полей и возвращает новый
  контекст. Это позволяет безопасно создавать независимые ветки с разными параметрами.
- **Повторное значение не создаёт новый контекст.** Если такой же ключ с таким же значением уже
  существует, `WithField` возвращает исходный `context.Context` без лишнего копирования.
- **Ошибка хранит снимок полей.** `WrapError` копирует текущую карту, поэтому последующие
  изменения другого контекста не меняют параметры уже возвращённой ошибки.
- **Поддерживаются произвольные структурированные поля.** В контекст можно положить UUID,
  строку, число, срез или другое значение, которое поддерживает `slog`.

### Добавление полей

В каждом модуле, которому нужны собственные поля логирования, необходимо создать пакет
`logging`, например `internal/modules/queue/logging`. Добавление полей должно выполняться только
через именованные функции этого пакета:

```go
func WithItemID(ctx context.Context, logger ports.Logger, itemID uuid.UUID) context.Context {
	if logger == nil || itemID == uuid.Nil {
		return ctx
	}
	return logger.WithField(ctx, "itemID", itemID)
}
```

Не вызывай `logger.WithField(ctx, "itemID", itemID)` напрямую в разных слоях. Если строковые
ключи распределить по handlers, usecases и repositories, со временем появятся разные варианты
одного имени, например `itemID`, `itemId` и `item_id`. Именованный метод хранит ключ в одном
месте, задаёт единый формат значения и позволяет централизованно добавить валидацию или
маскирование чувствительных данных.

Поле следует добавлять сразу, как только значение стало известно, и только один раз. Не нужно
ждать места, где произойдёт логирование: все вложенные вызовы должны заранее получить уже
обогащённый контекст. Обычно первая подходящая точка выглядит так:

- `userID` — в JWT middleware после проверки токена;
- `itemID` — в HTTP handler после разбора path-параметра;
- `ticketID` — после успешного декодирования тела запроса;
- созданный `queueID` — в application-слое после создания очереди.

Пример добавления поля сразу после того, как оно стало известно:

```go
itemID, err := uuid.Parse(chi.URLParam(r, ItemIdMuxPattern))
if err != nil {
	// вернуть 400
}

ctx := logging.WithItemID(r.Context(), logger, itemID)
result, err := usecase.GetItem(ctx, itemID)
```

Все названия полей пишутся в camelCase: `userID`, `itemID`, `queueID`, `ticketID`.
Повторно добавлять известное поле в usecase и repository не нужно. Если одинаковое поле всё же
будет передано повторно, `LogCtx.WithField` вернёт исходный контекст без нового копирования.

### Запись лога

Всегда передавай контекст в логгер:

```go
u.logger.ErrorContext(ctx, "failed to get item", "error", err)
```

Поля из контекста добавятся автоматически. Поэтому не нужно повторять их аргументами:

```go
// Не нужно: itemID уже находится в ctx.
u.logger.ErrorContext(ctx, "failed to get item", "itemID", itemID, "error", err)
```

Обычные методы `Error`, `Info`, `Warn` и `Debug` не используются: они не получают контекст и
не смогут автоматически добавить его поля.

### Передача полей наверх через ошибку

Обычный `context.Context` распространяется только вниз: handler передаёт его в usecase, а
usecase — в repository. Если новое поле стало известно внутри repository или другого глубокого
вызова, простое присваивание нового контекста не изменит контекст вызывающего слоя.

Пакет `shared/logctx` решает эту проблему и позволяет развернуть накопленные параметры лога
обратно до верхнего слоя:

1. `WrapError` сохраняет снимок полей текущего контекста внутри возвращаемой ошибки.
2. Ошибка поднимается через остальные слои обычным `return err`.
3. `ContextFromError` на слое, который пишет лог, достаёт поля из ошибки и возвращает
   восстановленный контекст.

Так параметры, известные глубоко в стеке вызовов, доходят до единственной итоговой записи в
логе, хотя сам `context.Context` напрямую наверх не передаётся.

### Ошибки repository-слоя

Repository не пишет ошибку в лог самостоятельно, иначе одна ошибка будет залогирована на
нескольких слоях. Он добавляет техническое описание и сохраняет текущие поля внутри ошибки:

```go
if err != nil {
	return nil, repo.wrapError(ctx, fmt.Errorf("get item: %w", err))
}
```

Где локальный helper безопасно работает и в тестах без логгера:

```go
func (repo *Repository) wrapError(ctx context.Context, err error) error {
	if repo.logger == nil {
		return err
	}
	return repo.logger.WrapError(ctx, err)
}
```

Application-слой восстанавливает контекст перед единственной записью в лог:

```go
item, err := u.repo.GetItemByID(ctx, itemID)
if err != nil {
	ctx = u.logger.ContextFromError(ctx, err)
	u.logger.ErrorContext(ctx, "failed to get item", "error", err)
	return nil, domain.ErrInternal
}
```

`WrapError` реализует `Unwrap`, поэтому проверки через `errors.Is` и `errors.As` продолжают
работать. `WrapError` нужно вызывать только для возвращаемой ошибки; успешный результат
оборачивать не требуется.

## RSA-ключи для JWT

Auth-модуль подписывает JWT алгоритмом RS256. Для локального запуска создай пару RSA-ключей
в каталоге `backend/keys` из корня репозитория:

```sh
mkdir -p backend/keys

openssl genpkey \
  -algorithm RSA \
  -pkeyopt rsa_keygen_bits:2048 \
  -out backend/keys/private.pem

openssl pkey \
  -in backend/keys/private.pem \
  -pubout \
  -out backend/keys/public.pem

chmod 600 backend/keys/private.pem
```

После генерации файлы должны находиться здесь:

```text
backend/keys/private.pem  # auth-сервис подписывает токены; не коммитить
backend/keys/public.pem   # сервисы проверяют подпись; можно коммитить
```

Пути внутри контейнера задаются в `internal/configs/backend.yml` как
`./keys/private.pem` и `./keys/public.pem`. `docker-compose.yaml` монтирует локальные файлы в
`/app/keys/private.pem` и `/app/keys/public.pem`, поэтому дополнительных действий для Docker не
нужно. Приватный ключ уже исключён из Git через корневой `.gitignore`.

После создания ключей пересобери backend:

```sh
docker compose up --build --force-recreate -d backend
```

## Тесты и проверки

- `task backend:vet` / `task backend:test` (с `-race`) / `task backend:build` — по отдельности,
  или все разом `task backend:check`.
- Интеграционные тесты поднимают Postgres в Docker через testcontainers — нужен работающий
  Docker, отдельно поднимать `docker compose` не требуется.
- `task backend:lint` — `golangci-lint`, конфиг в `../.golangci.yaml` (общий для `backend` и
  `migrator`).
