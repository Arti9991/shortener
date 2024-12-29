# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

Комманды для ручной проверки сервера (на данной итерации)

Стартовый POST запрос в cURL:

```
curl -v -X POST -H "Content-Type: text/plain" -d www.ya.ru http://localhost:8080
```

POST запрос c JSON в cURL:
```
curl -v -X POST -H "Content-Type: application/json" -d "{\"url\":\"www.ya.ru\"}" http://localhost:8080/api/shorten
```

Get запрос для извлечения ссылки в заголовке location:

```
curl -v GET -H "Content-Type: text/plain" http://localhost:8080/<id>
```

Ping запрос для проверки подключения к базе данных:

```
curl -v GET http://localhost:8080/ping
```

"{\"correlation_id\":\"ID\",\"url\":\"www.ya.ru\"}"
curl -v -X POST -H "Content-Type: application/json" -d "{\"correlation_id\":\"ID\",\"url\":\"www.ya.ru\"}" http://localhost:8082/api/shorten/batch