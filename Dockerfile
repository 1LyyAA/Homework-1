# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . .
RUN go build -o main .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# # Создаём директорию для логов
# RUN mkdir -p /app/logs

COPY --from=builder /build/main .

EXPOSE 8080

ENTRYPOINT ["./main"]
