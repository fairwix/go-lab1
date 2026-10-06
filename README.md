# Лабораторная работа №1 — бронирование ресурсов

Небольшой сервис на Go для бронирования переговорной комнаты и рабочего стола. У него два способа обращения — HTTP и gRPC. Оба сервера работают с одним экземпляром `BookingService`, поэтому бронь, созданная через один протокол, видна через другой. Данные хранятся в памяти и после перезапуска исчезают.

## Запуск

Нужен Go 1.27.1. Из папки `lab1`:

```bash
go run ./cmd/service
```

По умолчанию HTTP слушает порт `8080`, gRPC — `9090`. Порты можно поменять переменными окружения:

```bash
HTTP_PORT=8081 GRPC_PORT=9091 go run ./cmd/service
```

При запуске создаются два свободных ресурса: переговорная с ID `1` на 10 человек и стол с ID `2`. Каждый ресурс можно забронировать один раз до перезапуска сервиса. Число участников должно быть положительным; для переговорной оно также не может превышать 10. У стола вместимость не проверяется: параметр нужен общему методу `Reserve`, но сам стол бронируется как одно место.

## HTTP API

```bash
curl http://localhost:8080/health
curl http://localhost:8080/resources
curl -i -X POST http://localhost:8080/resources/book \
  -H 'Content-Type: application/json' \
  -d '{"id":1,"attendees":4}'
curl http://localhost:8080/resources
```

`GET /health` возвращает `{"status":"ok"}`. `GET /resources` возвращает список ID и статусов (`available` или `booked`). При успешном `POST /resources/book` ответ пустой, код — `204 No Content`. Неизвестные поля и лишний JSON в теле запроса отклоняются.

## gRPC

Контракт находится в [`proto/booking.proto`](proto/booking.proto), сгенерированные Go-файлы лежат рядом с ним. На сервере включена reflection, поэтому для ручной проверки подходит `grpcurl`:

```bash
grpcurl -plaintext localhost:9090 booking.BookingService/GetState
grpcurl -plaintext -d '{"id":2,"attendees":1}' localhost:9090 booking.BookingService/Book
grpcurl -plaintext localhost:9090 booking.BookingService/GetState
grpcurl -plaintext localhost:9090 grpc.health.v1.Health/Check
```

После вызова `Book` через gRPC можно вызвать `GET /resources` по HTTP и увидеть тот же статус. Работает и обратная проверка.

## Как устроен код

- `internal/model` содержит общие типы: `Reservable`, статус бронирования, ошибки, `BaseEntity` и `TimeSlot`.
- `internal/model/room` и `internal/model/desk` реализуют `Reservable`. Комната проверяет вместимость, стол — нет. Оба типа встраивают `BaseEntity`: это композиция, а не наследование с переопределением методов базового класса.
- `internal/service` проверяет входные данные, находит ресурс и меняет его состояние. Доступ к общей карте защищён `sync.Mutex`.
- `internal/transport/http` и `internal/transport/grpc` разбирают запросы, вызывают один и тот же сервис и преобразуют результат в ответ своего протокола.
- `cmd/service/main.go` создаёт ресурсы и одновременно запускает оба сервера. По `SIGINT` или `SIGTERM` они останавливаются с общим тайм-аутом в 5 секунд.

`TimeSlot` здесь демонстрирует работу со значением и указателем: методы чтения получают копию, а `ShiftSlot` меняет исходный интервал через указатель. Он создаётся и сдвигается при запуске, но в текущем API время бронирования не задаётся.

Ошибки предметной области возвращаются из модели и сервиса, а транспортные слои сопоставляют их с кодами ответа:

| Ситуация | HTTP | gRPC |
| --- | --- | --- |
| Неверные данные | `400 Bad Request` | `InvalidArgument` |
| Ресурс не найден | `404 Not Found` | `NotFound` |
| Ресурс уже забронирован | `409 Conflict` | `AlreadyExists` |
| Неожиданная ошибка | `500 Internal Server Error` | `Internal` |

В сервисе ошибка бронирования оборачивается через `%w`, а адаптеры находят нужную причину через `errors.Is` и `errors.As`.

## Проверка сборки

```bash
go build ./...
go vet ./...
```
