# Messenger API

Проект представляет собой RESTful API для обмена сообщениями (мессенджер), реализованный на Go с использованием современного стека технологий.

## Технологический стек

- **Язык**: Go 1.21+
- **Фреймворк**: Gorilla Mux для маршрутизации
- **База данных**: PostgreSQL
- **ORM**: GORM v2
- **Миграции**: Goose
- **Логгирование**: Zap от Uber
- **Тестирование**: Testify
- **Контейнеризация**: Docker + Docker Compose

## Архитектура

Проект построен по принципу **трехуровневой архитектуры (3-tier)**:

1. **Transport Layer** (`transport/`) - обработка HTTP запросов/ответов
2. **Service Layer** (`service/`) - бизнес-логика приложения
3. **Repository Layer** (`repository/`) - работа с базой данных

## Быстрый запуск

### Предварительные требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd messenger
```

### 2. Настройка окружения
Создайте файл .env на основе примера:

```bash
cp .env.example .env
```

Отредактируйте .env файл:

```bash
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=ваш юзер
POSTGRES_PASSWORD=ваш пароль
POSTGRES_DB=messenger
SSLMode=disable

APP_PORT=8080
```

### 3. Запуск через Docker Compose

```bash
docker-compose up --build
```

### 4. Работа с API через Postman

Эндпойнты:

1. http://localhost:8080/chats

* В теле передаем title. Метод POST

**Пример тела**

```bash
{
    "title": "chat1"
}
```

**Пример ответа**

```bash
{
    "id": 6,
    "title": "chat1",
    "created_at": "2026-01-31T23:27:14Z"
}
```

2. http://localhost:8080/chats/1/messages

* В теле передаем text. Метод POST

**Пример тела**

```bash
{
    "text": "hello"
}
```

**Пример ответа**

```bash
{"id":5,"chat_id":1,"text":"hello","created_at":"2026-01-31T23:30:51Z"}
```

3. http://localhost:8080/chats/1?limit=50

* В query параметрах передаем limit, по умолчанию он равен 20. Метод Get

**Пример ответа**

```bash
{"1":[{"text":"hello","time":"2026-01-31T23:30:51Z"}]}
```

4. http://localhost:8080/chats/1

*  Метод Delete

## Запуск unit-тестов

```bash
go test ./internal/transport -v
```

Вывод:

```bash
=== RUN   TestHandler_HendleCreateChat_Success
--- PASS: TestHandler_HendleCreateChat_Success (0.00s)
=== RUN   TestHandler_HendleCreateMessage_InvalidChatID
--- PASS: TestHandler_HendleCreateMessage_InvalidChatID (0.00s)
=== RUN   TestHandler_HendleGetMessages_WithLimit
--- PASS: TestHandler_HendleGetMessages_WithLimit (0.00s)
PASS
ok      messenger/internal/transport    (cached)
```

## Для запуска без docker`a

Первые 2 пункта прийдется пройти в любом случае.
Но потом можно будет восполььзоваться Makefile
