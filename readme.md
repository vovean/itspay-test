Для хранения курсов я бы использовал Clikchouse, потмоу что он лучше подходит для append-only нагрузки

## Запуск:

```shell
make run-docker-compose
``` 

> Обратите внимание, что там используется команда `docker compose` (без дефиса). Если у вас в системе используется
> `docker-compose`, то обновите Makefile

Эта команда

1. Поднимет нужную инфру
2. Проведет миграции
3. Запустит приложение через docker-compose

## Тесты

Чтоб прогнать тесты надо поднять БД для тестов (отдельный контейнер с postgres в docker-compose.yml). При запуске
`make test` также поднимается вся нужная инфра.

Для создания моков используется [moq](https://github.com/matryer/moq).

```shell
go install github.com/matryer/moq@latest
```

Для тестов используется testify

## Миграции

Для управления схемой БД используется инструмент [go-migrate](https://github.com/golang-migrate/migrate)

```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Кроме того, эти миграции используются динамически в тестах, чтоб реализовать очистку данных после каждого теста - перед
каждым тестам накатываются, после каждого теста откатываются

## go:generate

Генерация моков и гошного кода из protobuf делается с помощью go generate, можно запустить через`make generate`

## Что еще реализовано

1. Метрики: :8081/metrics
2. Healthcheck: :8081/health
3. Readiness check: :8081/ready
4. Pprof: :8081/debug/pprof/
5. Jaeger: localhost:16686/ - настроен сбор трейсов

## Другие мысли

1. Я бы использовл Clickhouse, чтобы хранить данные о курсах, потому что эта БД лучше подходит для хранения append-only
   данных