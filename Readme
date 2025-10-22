# URL Shortener

Мини-сервис для создания сокращённых ссылок.

Быстрый старт
### запуск
```
go run cmd/main.go
```

### или сборка
```
go build -o cmd/main.go
./url-shortener
```

### Флаги

-port адрес и порт сервера. По умолчанию :8080
-repo_type тип хранилища("json", "postgres"). По умолчанию json
-data_source путь к файлу json-хранилища. По умолчанию data.json рядом с исполняемым файлом

### Что поднимается:
- Логирование и recover: middleware.Logger(), middleware.Recover()
- Статика: GET /front/static/* из front/static
- Главная: GET / -> front/templates/index.html

Маршруты API регистрируются в internal/transport.Route(e, app.URLHandler)

### Хранилище

Доступны два типа хранилища
- JSON. Путь до файла-хранилища задаётся флагом -data_source.
Структура имеет вид "short_code" : "original_url". Например:
  {"xpsKsIX": "https://google.com" }
- Postgres. Переменные окружения для доступа к базе хранятся в .env файле в корне проекта.
База имеет одну таблицу ```short_links``` со следующим форматом данных:
 
    **id** (serial) - автоинкремент
 
    **original_url** (text) - оригинальный URL-адрес
 
    **short_code** (text) - сокращённая ссылка

Миграция для таблицы хранится в ```internal/repo/repo_postgres/migrations```
