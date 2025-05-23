# QuoteBook — Мини-сервис для цитат

**QuoteBook** — это REST API сервис на Go для хранения и управления цитатами. Все данные хранятся в памяти. Поддерживаются добавление, получение, фильтрация, удаление и выбор случайной цитаты.

---

## Быстрый старт

### Локально (Go 1.21+)

```sh
git clone https://github.com/Frozelo/quoteBook.git
cd quoteBook
go run ./cmd/main.go
```

### Через Docker

```sh
docker build -t quote-book .
docker run --rm -p 8080:8080 quote-book
```

---

## Быстрая проверка через curl

Добавить новую цитату:

```sh
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```

Получить все цитаты:

```sh
curl http://localhost:8080/quotes
```

Получить случайную цитату:

```sh
curl http://localhost:8080/quotes/random
```

Фильтрация по автору:

```sh
curl http://localhost:8080/quotes?author=Confucius
```

Удалить цитату по ID:

```sh
curl -X DELETE http://localhost:8080/quotes/1
```

---

## Автоматическая самопроверка (Self-Check)

Можно запустить сервис в режиме автоматической проверки всех функций API.

### Через Docker

```sh
docker run --rm -e SELF_CHECK=1 quote-book
```

Если все запросы прошли успешно — контейнер завершится с кодом `0` и в логах будет видно успешное выполнение self-check.
Если какая-либо операция завершится с ошибкой, контейнер завершится с ненулевым кодом и сообщением об ошибке.

---

## API эндпоинты

| Метод  | Путь             | Описание                  | Тело запроса / Query                |
| ------ | ---------------- | ------------------------- | ----------------------------------- |
| POST   | `/quotes`        | Добавить новую цитату     | `{"author": "...", "quote": "..."}` |
| GET    | `/quotes`        | Получить все цитаты       | (опционально) `?author=...`         |
| GET    | `/quotes/random` | Получить случайную цитату |                                     |
| DELETE | `/quotes/{id}`   | Удалить цитату по ID      |                                     |

---

## Структура проекта

```
.
├── cmd/main.go              # Точка входа (main)
├── internal/
│   ├── app/                 # Инициализация приложения, self-check
│   ├── handlers/            # HTTP-обработчики
│   ├── store/               # In-memory хранилище и логика работы с цитатами
│   ├── middelware/          # Middleware (логирование и др.)
│   └── server/              # Обертка для cервера и graceful shutdown сервера
├── pkg/errors/              # Кастомные ошибки
├── README.md
└── Dockerfile
└── docker-compose.yml
```

---

## Сборка и запуск тестов

```sh
go test ./internal/store
go test ./internal/handlers
```

---

## Работа с Makefile

Добавлены команды для удобства запуска тестов и cборки docker в обычном и self-check режиме

Запуск всех тестов
```sh
make test
```

Сборка docker-образа
```sh
make docker-build
```

Запуск docker-образа в обычном режиме
 ```sh
 make docker-run
 ```

Запуск self-check
```sh
make docker-self-check
```

Выполнить всё попорядку build + self_check
```sh
make check
```


## Примечания

* Данные хранятся только в памяти (in-memory store) в виде слайса.
* Используются только стандартные библиотеки Go и `gorilla/mux`.
* Все curl-примеры и документация доступны в этом README.
* Добавлены unit- и интеграционные тесты.
* Логирование запросов реализовано через Go slog.
* Сервер поддерживает корректное завершение (graceful shutdown).

---

## Автор

Автор: [Ivan Sizov](https://t.me/just_kilmz)
**Жду вашего фидбека**
