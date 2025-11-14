# Async api PoC

Эксперимент с реализацией bff с использованием асинхронных технологий

## Установка

Клонируем репозиторий и регистрируем воркспейс
```shell
git clone <>
cd async-api-poc

elc ws add example $PWD/workspace
elc ws select example
```

Запускаем скрипт инициализации. Он:
* копирует и заполняет env файлы
* создаёт БД для сервисов
* устанавливает зависимости

```shell
./workspace/scripts/init.sh
```
