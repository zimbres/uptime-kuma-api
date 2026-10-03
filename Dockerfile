FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o uptime-kuma-api .

FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/uptime-kuma-api .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./uptime-kuma-api"]