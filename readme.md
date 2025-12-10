# CLI ToDo
Консольное приложение для управления задачами (создание, чтение, обновление, удаление, отметка выполнения) с хранением в PostgreSQL.

## Функциональность
- `create` — добавить задачу (name, desc).
- `read` — показать задачу по `id`.
- `update` — изменить `name` и/или `desc`.
- `delete` — удалить по `id`.
- `done` — отметить задачу выполненной.
- `help` — краткая справка по командам.

## Требования и структура
- Go 1.25.2.
- Postgres 15 (по умолчанию в docker-compose).
- Миграции через goose, путь: `backend/pkg/migrations/postgres`.
- Слои: `port` (интерфейсы), `cases` (use case), `adapter/repo/postgres` (репозиторий), `entity` (модели), `app` (инициализация и CLI).

## Настройка окружения
1. Создайте `.env` в корне:
   ```
   POSTGRES_CONNECTION_STRING=postgres://user:password@localhost:5432/todo_db?sslmode=disable
   HTTP_PORT=:8080
   ```
2. Запустите Postgres (вариант):
   - Локально свой экземпляр PostgreSQL, параметры должны соответствовать строке подключения.
   - Или `docker-compose up -d postgres` в корне проекта (порт 5432 наружу).

## Миграции
- Автоматически применяются при запуске CLI (goose, пути `backend/pkg/migrations/postgres`).
- Основная миграция: `00001_init_tasks.sql` создаёт таблицу `tasks` (id, name, description, completed).

## Запуск CLI
Команды выполняются из корня репозитория:
```
go run ./backend/cmd <command> [flags]
```
Примеры:
- Создать: `go run ./backend/cmd create -name "Task" -desc "About"`
- Прочитать: `go run ./backend/cmd read -id 1`
- Обновить: `go run ./backend/cmd update -id 1 -name "New" -desc "New desc"`
- Удалить: `go run ./backend/cmd delete -id 1`
- Отметить выполненной: `go run ./backend/cmd done -id 1`
- Справка: `go run ./backend/cmd help`

### При ошибках подключения проверьте `POSTGRES_CONNECTION_STRING` и что Postgres запущен.