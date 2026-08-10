# Авито очередь. Кейс 2
Сервис реализован в рамках Авито хакатона.

Работа велась при помощи:
* Jira, kan deck
* New feature -> branch(feature/KAN-№) with pr's.
* Agile метадологии

В рамках проекта был реализован CI пайплайн для автоматизации сборки, линтинга и тестирования продукта.

### Особенность реализации
Сервис представляет собой модульный MVP монолит, реализующий FIFO атомарную очередь при помощи PostgreSQL для честного доступа к лимитированным товарам. Основная задача сервиса — гарантировать единство пользователя на товар и исключить ситуацию множественных покупок одного товара разными пользователями. Атомарность операций достигается за счет использования row-level locking в PostgreSQL: при постановке в очередь или покупке товара выполняется блокировка конкретной строки товара командой SELECT FOR UPDATE, что позволяет конкурентным запросам корректно обрабатываться последовательно, а не создавать гонки данных. Вся бизнес-логика обернута в транзакции, которые гарантируют, что изменения в очереди и остатках товара применяются атомарно. Для демонстрации полного workflow покупки были реализованы мок-сервисы: inventory/catalogue для управления остатками и информацией о товарах, payment для обработки платежей и auth для авторизации пользователей через JWT. Также был сделан осознанный выбор в пользу HTTP-Pooling, а не Web-Socket в рамках простоты реализации для MVP. Для observability была использована Victoria Metrics с custom dashboard в Grafana. Victoria Metrics является более производительным и легкой технологией для мониторинга и сбора метрик, чем Prometheus. Каждый компонент и инфраструктура проекта запускаются в Docker-контейнерах, что обеспечивает полную изоляцию сервисов и их независимость друг от друга. 

### Стек
* Go(chi, pgx, goose, slog, testcontainers, etc ..)
* PostgreSQL
* Docker
* Typescript
* React
* Victoria Metrics
* Grafana
* GolangCI-Lint
* Playwright

### Структура

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

Миграции — в отдельном Go-модуле ../migrator (свой go.mod, свой контейнер в docker-compose.yaml), не в backend.

Используется Hexagonal Architecture(ports & adapters for each module) 

### Запуск
Подготовка переменных окружения:

```zsh
touch .env
cp .env.template .env
```

Для локальной разработки и развертывания были описаны Taskfile таргеты:

| Task | Description |
|---|---|
| `task up` | Build and start the entire application |
| `task down` | Stop all services |
| `task logs` | Follow logs from all services |
| `task backend:test` | Run backend integration and unit tests |
| `task tests:run` | Build and run automated tests |
| `task clean` | Stop all services and remove volumes |

```sh
go install github.com/go-task/task/v3/cmd/task@latest # Установка Taskfile

task --version # check taskfile in path
```

Генерация RSA ключей:
```zsh
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

Приложение полностью конфигурируется при помощи конфига, который расположен по пути:
```zsh
./backend/internal/configs/backend.yml
```


### Тестирование
В рамках проекта ключевые модули системы были покрыты e2e, unit и integration тестами:
* pytest + playwright
* testcontainers
* testify
* testing

Общее покрытие тестами составляет около ~70%.


### Grafana
Метрики и дашборды доступны в Grafana UI:  http://localhost:4000. Метрики собираются с каждой ручки items и queue path: кол-во запросов в секунду, статус коды каждого запроса, urls и тд.

Data:
* Login: admin
* Password: admin


### Вклад каждого участника
#### Даша(DashaSmir)
* Полностью взяла на себя UI/UX дизайн проекта.
* Реализовала форму авторизации, каталог, карточку товара, окнол ожидания очереди на товар.
* Реализовала фильтрацию каталога по категориям.
* Реализовала форму покупки товара.
* Реклизовала граммотную смену статуса на frontend для каждого события в очереди(queued -> offered -> checkout -? soldout/cancelled, etc.. )
* Реализовала таймер очереди при помощи данных от backend.
* Интегрировала взаимодействие frontend с backend.
* project idea brainstorming.
* Написала mock заглушки для дальнейшей работы после MVP в рамках финала хакатона.


#### Даня(bober-17)
- Написал и дорабатывал техническую документацию и API-контракт проекта после общего согласования с командой.
 - Разложил backend на модульный монолит и исправил баг с транзакциями (WithTx не оборачивал запросы в реальную transaction).
 - Реализовал usecase advanceQueue (идемпотентное продвижение очереди) и эндпоинт GET /queue/me для поллинга статуса заявки.
 - Реализовал модуль checkout с нуля: usecase оформления покупки и обработки payment callback.
 - Привёл auth, queue, checkout и каталог к api-контракту: персистентные пользователи (upsert + миграция), JWT issuer/audience/ttl из конфига, публичный каталог с реальным soldOut через batch-запрос в очередь.
 - Реализовал GET /items/{itemId}/similar (подбор похожих товаров по категории) и обновил сиды каталога с реальными фото.
 - Пофиксил критичные UX-баги checkout/queue на frontend: таймер, поллинг очереди, рабочую отмену заказа (слот раньше не освобождался).


#### Влад(idenicalaffiliation)
* refactoring + project idea brainstorming.
* review each pr.
* unit + integration tests with testcontainers.
* e2e auto tests with playwright and pytest.
* Разработал inventory mock service(get item/items для каталога) + seed data.
* Интегрировал и настроил observability for routes with Victoria Metrics and Grafana.
* migrations with goose.
* Реализовал полный workflow становления пользователя в очередь(атомарная вставка без гонки данных при помощи row level locking and constraints).


#### Леша(kotafan1rich)
* Настроил и поднял инфраструктуру при помощи Docker.
* Интегрировал CI пайплайн(build, lint, test) в actions.
* Разработал mock auth сервис(login endpoint, auth middleware, jwt bearer token with rsa key).
* refactoring + project idea brainstorming.
* unit + integration tests with testcontainers.
* Реализовал полный workflow выхода пользователя из очереди(endpoint, usecase, repo), обновление статусов при выходе(QUEUED/OFFERED/CHECKOUT -> CANCELLED).
* review each pr.
