## Установка

Запускаем clickhouse, prometheus, cadvisor и grafana. Датасорс и дашборды настроятся автоматически
```bash
docker compose up -d
```

При необходимости создаём virtualenv (если не было создано автоматически через IDE)
```bash
python3 -m venv venv
source venv/bin/activate
```

Устанавливаем пакеты
```bash
pip3 install -r requirements.txt
```

## Использование

Запускаем тест
```bash
locust --headless --only-summary --host=http://api.example.127.0.0.1.nip.io/ --locustfile locustfiles/api.py
```

Результаты смотрим в графане http://127.0.0.1:3000

## Прочее

Посмотреть список тестов
```bash
./manage.py list_tests
```

Удалить указанный тест
```bash
./manage.py delete_test <test_id>
```

Чтобы удалить вообще все тесты и метрики в prometheus, проще будет удалить volumes контейнеров
```shell
docker compose down -v
```
