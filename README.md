# L0_task

Сервис обработки заказов: получает данные из Kafka, сохраняет в PostgreSQL и отдаёт через HTTP API.

## Установка

1. Склонировать репозиторий:
   ```bash
   git clone https://github.com/alinakobzar/L0_task.git

2. Поднять kafka
   ```bash
   docker-compose up -d

3. Запустить producer 
    ```bash
    cd cmd/producer
    go run main.go 
   ```

4. Запустить сервис
   ```bash
   cd cmd/api
   go run main.go
   ```
