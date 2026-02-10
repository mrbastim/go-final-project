# Выполненные задания повышенной сложности

Были выполнены все задания повышенной сложности, включая [Dockerfile](Dockerfile), кроме аутентификации в шаге №8.

# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

## Локальный запуск

1. (Опционально) Создайте файл `.env` в корне проекта (на примере `.env.example`):

```env
TODO_PORT=7540
TODO_DBFILE=./scheduler.db
```

2. Запустите сервер из корня проекта:

```bash
go run ./back
```

Адрес в браузере: http://localhost:7540/ (порт берется из `TODO_PORT`).

## Запуск тестов

1. Убедитесь, что сервер запущен.
2. Проверьте параметры в [tests/settings.go](tests/settings.go):

```go
var Port = 7540          // должен совпадать с TODO_PORT
var DBFile = "../scheduler.db" // путь к файлу БД, который использует сервер
var FullNextDate = true
var Search = true
var Token = ``
```

3. Запустите тесты:

```bash
go test ./tests
```

## Docker: сборка и запуск

Сборка образа:

```bash
docker build -t todo-app .
```

Запуск контейнера с пробросом порта и сохранением БД:

```bash
docker run --rm -p 7540:7540 \
	-e TODO_PORT=7540 \
	-e TODO_DBFILE=/app/scheduler.db \
	-v "$PWD/scheduler.db:/app/scheduler.db" \
	todo-app
```

Адрес в браузере: http://localhost:7540/ (или другой порт, если изменили `TODO_PORT`).