# Этап сборки
FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service ./cmd/main.go

# Этап запуска
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/order-service .
EXPOSE 8002

CMD ["./order-service"]
