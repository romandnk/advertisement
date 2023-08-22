# Объявления

Приложение для создания объявлений. API приложения предоставляет инструмент для регистрации пользователей, аутентификации с использованием JWT-токенов и создания объявлений с изображениями. Каждое объявление обладает не только текстовой информацией, но и загруженными локально фотографиями, в то время как сами данные хранятся надежно в базе данных PostgreSQL. Кроме того, вы можете получать изображения объявлений, используя уникальные ссылки.

## Методы HTTP

### Объявления:

- `POST /events` - создать объявление
- `GET /events/{id}` - получить объявление по ID
- `DELETE /events/{id}` - удалить объявление по ID

### Пользователь:

- `POST /users/sign-up` - зарегистрировать нового пользователя
- `POST /users/sign-in` - авторизоваться пользователем через email и пароль

### Изображение:

- `GET /images/{id}` - получить изображение по ID

## Используемые технологии, методологии и инструменты:

- Golang
- Chi
- PostgreSQL
- Docker
- Git
- HTTP/API
- JWT авторизация
- RESTful
- Unit-тестирование

## Запуск

1. Склонируйте репозиторий.
2. Зайдите в папку с проектом.
3. Введите `make run`.
