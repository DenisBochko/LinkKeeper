﻿# LinkKeeper
---
<code>git clone https://github.com/DenisBochko/LinkKeeper.git</code>
---
.env берём из примера .env.example вставляем тг токен
<code>echo TOKEN=<ВАШ ТОКЕН> > .env</code>
---
Запускаем postgres из docker (Возможно необходимо будет сменить порт на 5431) (порт лучше не менять)
1. <code>docker-compose up -d</code>
2. <code>docker ps</code>
3. <code>docker-compose down</code>
---
cd database
cd migrations
go install github.com/pressly/goose/v3/cmd/goose@latest
---
потом нужно применить миграции 
- goose postgres "user=postgres dbname=LinkKeeper password=postgres host=localhost sslmode=disable" up

устанавливаем gpt4free
- pip install -U g4f[all]
- 
