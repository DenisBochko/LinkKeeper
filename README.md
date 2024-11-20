# LinkKeeper

LinkKeeper - телеграм-бот, написанный на golang, для сохранения, управления и анализа ссылок.

## Используемые технологии
1. Golang
2. Python
3. PostreSQL
4. Docker
5. Telegram Bot API

## Архитектура приложения

![Не загрузилось(](/image.png "Архитектура")

## Локальный запуск

1. Клонируем проект 
- ```
  git clone https://github.com/DenisBochko/LinkKeeper.git
  ```
2. Переходим в директорию LinkKeeper и создаём .env по образцу 
- ```
  cd LinkKeeper
  ```
- ```
  echo TOKEN=<ВАШ ТОКЕН> > .env
  ```
3. Запускаем postgres из docker
- ```
  docker-compose up -d
  ```
- ```
  docker ps
  ```
- ```
  docker-compose down
  ```
4. Устанавливаем библиотеку для работы с миграциями 
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```
5. Переходим в директорию с миграциями и применяеим их с помощью goose
- ```
  cd database
  ```
- ```
  cd migrations
  ```
- ```
  goose postgres "user=postgres dbname=LinkKeeper password=postgres host=localhost sslmode=disable" up
  ```
6. Создаём отдельную директорию и активируем в ней venv python и устанавилваем gpt4free
- ```
  cd PythonAIserver
  ```
- ```
  python3 -m venv venv
  ```
- ```
  .\venv\Scripts\activate
  ```
- ```
  pip install -U g4f[all]
  ```
7. Запускаем 3 сервера на 3 разных портах 
- ```
  g4f api --ignored-providers 'RubiksAI' --bind "0.0.0.0:1337"
  ```
- ```
  g4f api --ignored-providers 'RubiksAI' --bind "0.0.0.0:1338"
  ```
- ```
  g4f api --ignored-providers 'RubiksAI' --bind "0.0.0.0:1339"
  ```
8. Запускаем в нашей директории главный файл и радуемся:)
- ```
  go run main.go
  ```