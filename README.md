# Сервис для агрегации данных об подписках пользователей


## Требования
- Docker и Docker Compose

## Запуск
```bash
docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`, PostgreSQL — на порту `5432`.

## Конфигурация

Переменные окружения (значения по умолчанию заданы в `docker-compose.yml`):

- `APP_PORT`: `8080` - Порт HTTP-сервера
- `DB_PORT`: `5432` - Порт PostgreSQL
- `DB_HOST`: `db` - Хост PostgreSQL
- `DB_NAME`: `subscriptions` - Имя БД
- `DB_USER`: `postgres` - Пользователь БД
- `DB_PASSWORD`: `postgres` - Пароль БД
- `LOG_LEVEL`: `info` - Уровень логов

Для локального запуска без Docker можно создать `.env` или задать переменные вручную; Конфиг читается из `internal/config/config.yaml` с подстановкой переменных окружения.

## API

- **Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Health:** `GET /health` — проверка работоспособности

### Эндпоинты

- POST: `/subscriptions` - создать подписку 
- GET: `/subscriptions` - список подписок (query: `limit`, `offset`)
- GET: `/subscriptions/total-cost` - суммарная стоимость за период (query: `start_date`, `end_date`, опционально `user_id`, `service_name`)
- GET: `/subscriptions/{id}` - подписка по ID
- PUT: `/subscriptions/{id}` - обновить подписку
- DELETE: `/subscriptions/{id}` - удалить подписку

Формат дат: **MM-YYYY** (например, `07-2025`). Стоимость — целое число рублей.

### Примеры запросов

**Создать подписку**

```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

**Получить одну подписку по id**

```bash
curl "http://localhost:8080/subscriptions/480850a7-0c6c-445d-8be6-3ff0b130168b"
```

**Список подписок**

```bash
curl "http://localhost:8080/subscriptions?limit=20&offset=0"
```

**Суммарная стоимость за период (с фильтром по пользователю)**

```bash
curl "http://localhost:8080/subscriptions/total-cost?start_date=01-2025&end_date=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

**Обновить подписку**

```bash
curl -X PUT "http://localhost:8080/subscriptions/480850a7-0c6c-445d-8be6-3ff0b130168b" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus Premium",
    "price": 500,
    "end_date": "06-2026"
  }'
```

**Удалить подписку**

```bash
curl -X DELETE "http://localhost:8080/subscriptions/480850a7-0c6c-445d-8be6-3ff0b130168b"
```


## Сборка и тесты

```bash
go build ./cmd/app/
go test ./...
```

## Структура проекта

- `cmd/app/` — точка входа приложения
- `internal/config/` — конфигурация и логгер
- `internal/handlers/` — HTTP-обработчики
- `internal/models/` — модели и DTO
- `internal/repository/` — работа с PostgreSQL
- `internal/service/` — бизнес-логика и валидация
- `internal/validation/` — кастомные валидаторы (формат дат)
- `migrations/` — SQL-миграции (golang-migrate)
- `docs/` — сгенерированная Swagger-спецификация
