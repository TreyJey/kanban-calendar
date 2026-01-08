# Dockerfile
# Многостадийная сборка для Go приложения

# ========== СТАДИЯ СБОРКИ ==========
FROM golang:1.25-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git ca-certificates

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main .

# ========== ФИНАЛЬНЫЙ ОБРАЗ ==========
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

# Создаем пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Рабочая директория
WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/main .

# Копируем миграции
COPY migrations ./migrations

# Копируем файлы конфигурации
COPY .env.example ./.env.example

# Меняем владельца файлов
RUN chown -R appuser:appgroup /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Порт приложения
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]