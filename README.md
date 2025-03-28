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

Например, запуск тестов для iter12 из под GitBash для windows. 
```
./shortenertest -test.v -test.run=^TestIteration12$ -binary-path=cmd/shortener/shortener -database-dsn="host=localhost user=myuser password=123456 dbname=ShortURL sslmode=disable"
```

Комманды для ручной проверки сервера (на данной итерации)

Стартовый POST запрос в cURL:

```
curl -v -X POST -H "Content-Type: text/plain" -d www.ya.ru http://localhost:8082
```

POST запрос c JSON в cURL:
```
curl -v -X POST -H "Content-Type: application/json" -d "{\"url\":\"www.ya.ru\"}" http://localhost:8082/api/shorten
```

Get запрос для извлечения ссылки в заголовке location:

```
curl -v GET -H "Content-Type: text/plain" http://localhost:8082/<id>
```

Ping запрос для проверки подключения к базе данных:

```
curl -v GET http://localhost:8082/ping
```

Ручной запрос для проверки множественного POST:
```
curl -v -X POST -H "Content-Type: application/json" -d '[
{"correlation_id":"ID","original_url":"www.ya.ru"},
{"correlation_id":"ID","original_url":"www.dlya.ru"},
{"correlation_id":"ID","original_url":"www.Nya.ru"},
{"correlation_id":"ID","original_url":"www.Qya.ru"},
{"correlation_id":"ID","original_url":"www.Mya.ru"}]' http://localhost:8082/api/shorten/batch
```

POST запрос c установленными cookie:

```
curl -v -X POST -H "Content-Type: text/plain" --cookie "userID=<cookie>" -d www.ya.ru http://localhost:8082
```

POST запрос с JSON и установленными cookie:

```
curl -v -X POST -H "Content-Type: application/json" --cookie "userID=<cookie>" -d "{\"url\":\"www.Nya.ru\"}" http://localhost:8082/api/shorten
```

GET запрос для получения всех URL когда-либо сокращенных пользователем:
```
curl -v GET  --cookie "userID=<cookie>" http://localhost:8082/api/user/urls 
```

Ручной запрос для проверки множественного POST с cookie:
```
curl -v -X POST -H "Content-Type: application/json" --cookie "userID=<cookie>" -d '[
{"correlation_id":"ID","original_url":"www.ya.ru"},
{"correlation_id":"ID","original_url":"www.dlya.ru"},
{"correlation_id":"ID","original_url":"www.Nya.ru"},
{"correlation_id":"ID","original_url":"www.Qya.ru"},
{"correlation_id":"ID","original_url":"www.Mya.ru"}]' http://localhost:8082/api/shorten/batch
```
Удаление URL, ранее сохраненных пользователем
```
curl -v -X DELETE -H "Content-Type: application/json" --cookie "userID=<cookie>" -d '["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]' http://localhost:8082/api/user/urls
```

Запуск основного серверева с соединение к БД, но без сохранений в файлах (для файлов добавить флаг `-f=./storage.csv`):
```
DATABASE_DSN="host=localhost user=myuser password=123456 dbname=ShortURL sslmode=disable" ./shortener.exe -a :8082
```