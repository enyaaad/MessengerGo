# messengerTest — API чатов и сообщений (Go)

Сервис для тестового задания: **чаты** и **сообщения**.

## Запуск

```bash
cd messengerTest
docker compose up --build || make up
```

Сервис поднимется на `http://localhost:8080`.

### Переменные окружения / доступы к Postgres

- **По умолчанию postgres не публикуется наружу** (нет `ports: 5432:5432`), доступ к нему есть только внутри нетворка.
- Креды/имя БД задаются через переменные (можно переопределить перед запуском):
  - `POSTGRES_USER` (default `messenger`)
  - `POSTGRES_PASSWORD` (default `messenger`)
  - `POSTGRES_DB` (default `messenger`)

## Swagger / OpenAPI

- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI spec: `http://localhost:8080/openapi.yaml`