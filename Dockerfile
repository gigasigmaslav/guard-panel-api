FROM golang:1.26 AS builder
WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

# Копируем исходный код
COPY . .

# Собираем бинарник с использованием vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o main ./cmd/main.go

# Финальный образ
FROM alpine:3.21
WORKDIR /app

# Устанавливаем ca-certificates для HTTPS и wget для dockerize
RUN apk --no-cache add ca-certificates wget && \
    wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz && \
    tar -xvzf dockerize-linux-amd64-v0.6.1.tar.gz && \
    mv dockerize /usr/local/bin/ && \
    rm dockerize-linux-amd64-v0.6.1.tar.gz

# Копируем бинарник из builder
COPY --from=builder /app/main .

# Переменные окружения
COPY .env .

CMD ["dockerize", "-wait", "tcp://pg_db:5432", "-timeout", "60s", "./main"]
