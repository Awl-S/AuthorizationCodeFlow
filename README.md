
# OAuth Server with Authorization Code Grant Flow

OAuth сервер на основе Go и фреймворка Gin, реализующий поток Authorization Code Grant для выдачи токенов доступа и предоставления доступа к защищённым ресурсам.

## Описание

Сервер реализует поток Authorization Code Grant, что позволяет клиенту сначала запросить код авторизации, а затем обменять его на access token, который может быть использован для доступа к защищённым ресурсам.

### Основные компоненты

- **AuthorizeHandler** — принимает `client_id` и `redirect_uri`, генерирует код авторизации, и перенаправляет пользователя обратно на `redirect_uri` с кодом в URL.
- **TokenHandler** — обрабатывает код авторизации и выдает access token, если клиент предоставляет действительный код.
- **ResourceHandler** — проверяет наличие валидного access token и предоставляет доступ к защищённым данным.

## Структура проекта

1. **`main.go`** — Основной файл с определением обработчиков запросов и логикой OAuth сервера.
2. **`OAuthServer`** — Структура для хранения зарегистрированных клиентов, кодов авторизации и активных токенов.

## Установка

### Требования

- Go 1.18+
- Gin Framework для работы с HTTP запросами.

### Шаги для установки

1. Клонируйте репозиторий:
    ```bash
    git clone https://github.com/yourcompany/oauth-server.git
    ```

2. Перейдите в директорию проекта:
    ```bash
    cd oauth-server
    ```

3. Установите зависимости:
    ```bash
    go mod tidy
    ```

4. Запустите сервер:
    ```bash
    go run main.go
    ```

Сервер будет доступен по адресу `http://localhost:8080`.

## Использование

### 1. Получение Authorization Code

Отправьте GET-запрос на `/authorize`, указав `client_id` и `redirect_uri` в URL:
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/authorize?client_id=my_client_id&redirect_uri=http://localhost:8080/callback" -Method GET
```

Ответ будет перенаправлением на `redirect_uri` с кодом авторизации:
```
http://localhost:8080/callback?code=ABC123XYZ
```

### 2. Получение access token

Используйте код авторизации, чтобы получить access token, отправив POST-запрос на `/token` с `client_id` и `code`:
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/token" -Method POST -Body "client_id=my_client_id&code=ABC123XYZ"
```

Ответ сервера:
```json
{
  "access_token": "FWdD07QTfcA2pIikIA2X",
  "token_type": "Bearer"
}
```

### 3. Доступ к защищённому ресурсу

Используйте access token для доступа к защищённому ресурсу с помощью GET-запроса на `/resource`, добавив заголовок `Authorization`:
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/resource" -Method GET -Headers @{"Authorization" = "FWdD07QTfcA2pIikIA2X"}
```

Ответ сервера:
```json
{
  "data": "secure data"
}
```

### 4. Ошибки доступа

Если токен не предоставлен или является недействительным, сервер вернет ответ с кодом ошибки 401 и сообщением:
```json
{
  "error": "invalid_token"
}
```

## Архитектура

### `OAuthServer`

`OAuthServer` отвечает за хранение зарегистрированных клиентов (по `client_id` и `client_secret`), активных кодов авторизации и токенов доступа. Поток включает:

- **AuthorizeHandler** — генерирует код авторизации и связывает его с `client_id`.
- **TokenHandler** — проверяет, соответствует ли код авторизации зарегистрированному `client_id`, и выдает access token.
- **ResourceHandler** — проверяет access token и предоставляет доступ к защищённым данным.

### Генерация токенов и кодов

Функция `randSeq` используется для генерации случайных кодов и токенов из букв и цифр. Код авторизации создается длиной 10 символов, а токен доступа — 20.

## Безопасность и развертывание

Для развертывания сервера в продакшн необходимо учитывать следующие аспекты:

- Применение HTTPS для защиты данных при передаче.
- Использование более надежного метода хранения токенов, таких как JWT.
- Настройка защиты от атак с использованием слабых токенов и кодов.

## Лицензия

Этот проект распространяется под лицензией MIT.
