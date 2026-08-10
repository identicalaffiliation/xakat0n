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
│       └── httpx/               # JWTAuth (мидлварь авторизации), EncodeJSON
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

На примере: `queue` уже содержит `POST /items/{itemId}/queue`. Пусть нужно добавить,
например, `GET /items/{itemId}/queue/me`.

1. **`domain/`** — если нужны новые правила/поля на сущности, добавляй сюда. Для чтения статуса
   обычно ничего нового не требуется.
2. **`ports/`** — опиши, что usecase'у нужно от репозитория (`repo.go`) — например,
   `GetByProductAndUser(ctx, itemID, userID) (*domain.Queue, error)`.
3. **`infrastructure/postgres/`** — реализуй этот метод в `QueueRepository` (SQL-запрос).
4. **`dto/`** — опиши форму ответа, если она отличается от того, что уже есть.
5. **`application/`** — новый usecase (или метод на существующем): принимает `ports.QueueRepository`
   и остальные порты через конструктор, реализует бизнес-логику поверх интерфейсов.
6. **`presentation/http/`** — новый хендлер `http.HandlerFunc`, парсит path/заголовки, зовёт
   usecase, сериализует ответ через `httpx.EncodeJSON`. Достаёт `user_id` через
   `httpx.UserID(r.Context())` — сам заголовок не парсит и ничего не знает про JWT.
7. **`module.go`** — зарегистрируй роут на своём под-роутере: `r.With(httpx.Metrics).Get(route,
   handler.GetMyTicket(usecase))`. Авторизацию отдельно вешать не нужно — `httpx.JWTAuth(verifier)`
   применяется один раз глобально на весь роутер в `cmd/api/main.go` (`r.Use(...)`), а не на
   каждый роут по отдельности.
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
