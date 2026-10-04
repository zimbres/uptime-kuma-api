FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o uptime-kuma-api .

FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

COPY --from=builder /app/uptime-kuma-api .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

USER nonroot:nonroot

CMD ["./uptime-kuma-api"]