<div align="center">
  <h1 align="center">Song Library API</h1>
  <h3>API for managing the song library</h3>
</div>

<br/>

Song Library API — это API для управления библиотекой песен и получения подробной информации о песнях из внешнего API(add).

## Демо

![Song Library API Postman.jpg](https://sun9-77.userapi.com/s/v1/ig2/1Qb99iOpZg5P6H951I1IMnRP8JUa0FOAso7pCotwJE_SicSkQ50RCCcJMMxL60XNL-Mg57m_NETQ-aKDyy3-TfDI.jpg?quality=95&as=32x17,48x26,72x39,108x58,160x86,240x129,360x193,480x257,540x289,640x343,720x386,1080x579,1280x686,1440x772,1920x1029&from=bu&u=iOLeD52NDWUFBek2JO64bRt52bo9TzAQyHU6KcgYcko&cs=1920x1029)

## Стек

- [Go](https://go.dev/) – Programming language
- PostgreSQL - Database
- [Swagger](https://github.com/swaggo/swag) –  Tool for generating documentation
- [Chi](https://github.com/go-chi/chi) – Framework

## Приступая к работе

### Предварительные требования

Вот что вам нужно для запуска:

- Go (version >= 18)
- PostgreSQL Database

### 1. Склонируйте репозиторий

```shell
git clone https://github.com/aashpv/song-lib
```

### 2. Настройте .env

Измените .env файл, используя [.env](.env) как шаблон. Укажите параметры базы данных и другие настройки.
##### ВАЖНО❗ Не меняйте название БД❗

### 3. Запустите сервер

```shell
cd song-lib\cmd\app
go run main.go
```
##### ВАЖНО❗Создавать БД вручную и запускать миграции не требуется.

### 4. Откройте приложение в своем браузере

Посетите сайт [http://localhost:8080/songs](http://localhost:8080/songs) в своем браузере.

Документация будет доступна по адресу [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

## Основные маршруты API

- **GET /songs** - Получение списка песен с возможностью фильтрации и пагинации.
- **POST /songs** - Добавление новой песни.
- **GET /songs/{id}/text** - Получение текста песни с пагинацией по куплетам.
- **PUT /songs/{id}** - Обновление информации о песне.
- **DELETE /songs/{id}** - Удаление песни по ID.

## Пример использования внешнего API

При добавлении песни вызывается внешнее API, предоставляющее дополнительную информацию о песне:

```go
    group := strings.Replace(req.Group, " ", "+", -1)
    song := strings.Replace(req.Song, " ", "+", -1)
    externalApiUrl := fmt.Sprintf("http://localhost:8081/info?group=%s&song=%s", group, song)
    fmt.Println(externalApiUrl)
    
    response, err := http.Get(externalApiUrl)
    // Логика обработки запроса
```
Этот функционал был эмулирован на тестовом сервере, работающем на порту 8081.

