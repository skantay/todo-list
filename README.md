# Todo-list

Микросервис для работы со списками задач: создание, удаление и обновление задач.

Используемые технологии:

- MongoDB (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Swagger (для документации API)
- Gin (веб-фреймворк)

Сервис написан с использованием Clean Architecture, что позволяет легко расширять его функциональность и тестировать. Также реализован Graceful Shutdown для корректного завершения работы сервиса.

# Getting started

## Usage

По умолчанию используется порт 7777 (vip kazakh port xD)

Для запуска сервиса выполните команду `make compose-up`.

После запуска сервиса вы сможете просмотреть документацию API по адресу http://localhost:7777/swagger/index.html.

## Примеры

Некоторые примеры запросов:

- [Создание задачи](#create-task)
- [Удаление задачи](#delete-task)
- [Обновление задачи](#update-task)
- [Пометка задачи как завершенной](#mark-task)
- [Получение всех активных задач](#list-active-tasks)
- [Получение всех завершенных задач](#list-done-tasks)

### Создание задачи <a name="create-task"></a>

Request
```curl
curl --location --request POST 'localhost:7777/api/v1/todo-list/tasks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title":"title",
    "activeAt":"2024-04-01"
}'
```

Response
```json
{
    "id": "661fbb485131cd932a981b26"
}
```

### Удаление задачи <a name="delete-task"></a>

Request
```curl
curl --location --request DELETE 'localhost:7777/api/v1/todo-list/tasks/661f23f7f65b382540934424'
```

No body response

### Обновление задачи <a name="update-task"></a>

Request
```curl
curl --location --request PUT 'localhost:7777/api/v1/todo-list/tasks/661fbb485131cd932a981b26' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title":"updated",
    "activeAt":"2024-04-01"
}'
```

No body response

### Пометка задачи как завершенной <a name="mark-task"></a>

Request
```curl
curl --location --request PUT 'localhost:7777/api/v1/todo-list/tasks/661fb7d85131cd932a981b25/done'
```

No body response

### Получение всех активных задач <a name="list-active-tasks"></a>

Request
```curl
curl --location --request GET 'localhost:7777/api/v1/todo-list/tasks?status=active'
```

Response
```json
[
    {
        "id": "661fbb485131cd932a981b26",
        "title": "updated",
        "activeAt": "2024-04-01"
    }
]
```

### Получение всех завершенных задач <a name="list-done-tasks"></a>

Request
```curl
curl --location --request GET 'localhost:7777/api/v1/todo-list/tasks?status=done'
```

Response
```json
[
    {
        "id": "661fbb485131cd932a981b26",
        "title": "updated",
        "activeAt": "2024-04-01"
    }
]
```