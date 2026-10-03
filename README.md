# Uptime Kuma Monitor API (Go)

A Go REST API for managing Uptime Kuma monitors with MariaDB/MySQL backend.

**Module:** `zimbres/uptime-kuma-api`

## Features

- Full CRUD operations for monitors
- Get monitor last heartbeat status
- Monitor tags management (CRUD + associate/dissociate)
- Maintenance windows management (CRUD + associate/dissociate)
- Pause/Resume monitors
- Compatible with Uptime Kuma's MariaDB schema
- Optional Bearer token authentication (configurable via environment)
- Database connection via environment variables
- Swagger/OpenAPI documentation
- Docker support

## Project Structure

```
Go/
├── main.go                    # Entry point (package main)
├── go.mod                     # Module: zimbres/uptime-kuma-api
├── go.sum                     # Go modules checksums
├── .env.example               # Example environment variables
├── docker-compose.yml         # Docker Compose for local development
├── Dockerfile                 # Docker image definition
├── README.md                  # This file
├── docs/                      # Swagger documentation (auto-generated)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── internal/                  # Private packages
    ├── config/                # package config - Configuration management
    │   └── config.go
    ├── database/              # package database - Database connection
    │   └── database.go
    ├── models/                # package models - Data models & request/response types
    │   └── models.go
    ├── repository/            # package repository - Data access layer
    │   └── repository.go
    ├── handler/               # package handler - HTTP handlers
    │   └── handler.go
    └── middleware/            # package middleware - Authentication middleware
        └── auth.go
```

## Getting Started

### Prerequisites

- Go 1.27+
- MariaDB/MySQL 13.0+

### Local Development

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your database credentials

3. Run the database migrations (the mariadb.sql from the parent directory contains the schema):
   ```bash
   mysql -u root -p uptime < ../mariadb.sql
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

   Or build and run:
   ```bash
   go build -o uptime-kuma-api
   ./uptime-kuma-api
   ```

### Using Docker Compose

```bash
docker-compose up -d
```

This will start:
- Go API server (port 8080)

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | localhost | Database host |
| DB_PORT | 3306 | Database port |
| DB_USER | root | Database user |
| DB_PASSWORD | | Database password |
| DB_NAME | uptime | Database name |
| SERVER_PORT | 8080 | Server port |
| ENABLE_AUTH | false | Enable Bearer token authentication |
| AUTH_TOKEN | | Bearer token (required if ENABLE_AUTH=true) |

## API Endpoints

All endpoints are prefixed with `/api/v1`

### Monitors

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /monitors | Create a new monitor |
| GET | /monitors | List all monitors (paginated) |
| GET | /monitors/:id | Get a specific monitor (includes tags) |
| PUT | /monitors/:id | Update a monitor |
| DELETE | /monitors/:id | Delete a monitor |
| POST | /monitors/:id/pause | Pause (deactivate) a monitor |
| POST | /monitors/:id/resume | Resume (activate) a monitor |
| GET | /monitors/:id/heartbeat | Get monitor's last heartbeat |
| POST | /monitors/:id/tags | Add tag to monitor |
| GET | /monitors/:id/tags | Get monitor's tags |
| DELETE | /monitors/:id/tags/:tagId | Remove tag from monitor |
| POST | /monitors/:id/maintenances | Add maintenance to monitor |
| GET | /monitors/:id/maintenances | Get monitor's maintenances |
| DELETE | /monitors/:id/maintenances/:maintenanceId | Remove maintenance from monitor |

### Tags

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /tags | Create a new tag |
| GET | /tags | List all tags |
| GET | /tags/:id | Get a specific tag |
| PUT | /tags/:id | Update a tag |
| DELETE | /tags/:id | Delete a tag |

### Maintenances

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /maintenances | Create a new maintenance |
| GET | /maintenances | List all maintenances (paginated) |
| GET | /maintenances/:id | Get a specific maintenance |
| PUT | /maintenances/:id | Update a maintenance |
| DELETE | /maintenances/:id | Delete a maintenance |

### Authentication

If `ENABLE_AUTH=true`, include the Bearer token in the Authorization header:

```
Authorization: Bearer your-secret-token
```

### Swagger Documentation

Access the Swagger UI at: `http://localhost:8080/swagger/index.html`

## Monitor Types

The API supports all Uptime Kuma monitor types:
- HTTP, KEYWORD, JSON_QUERY, GRPC_KEYWORD
- PORT, PING, DNS
- DOCKER, REAL_BROWSER, PUSH
- STEAM, GAMEDIG, MQTT
- SQLSERVER, POSTGRES, MYSQL, MONGODB, REDIS
- RADIUS, TAILSCALE_PING
- KAFKA_PRODUCER, RABBITMQ, SNMP, SMTP
- MANUAL, SYSTEM_SERVICE, WEBSOCKET_UPGRADE
- GLOBALPING, SIP_OPTIONS, NTP, ORACLEDB, PM2

## Example Requests

### Create a Monitor

```bash
curl -X POST http://localhost:8080/api/v1/monitors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Google",
    "type": "http",
    "url": "https://google.com",
    "interval": 60
  }'
```

### List Monitors

```bash
curl http://localhost:8080/api/v1/monitors?page=1&limit=10
```

### Get a Monitor (includes tags)

```bash
curl http://localhost:8080/api/v1/monitors/1
```

### Update a Monitor

```bash
curl -X PUT http://localhost:8080/api/v1/monitors/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Google Updated",
    "interval": 120
  }'
```

### Delete a Monitor

```bash
curl -X DELETE http://localhost:8080/api/v1/monitors/1
```

### Pause/Resume Monitor

```bash
# Pause
curl -X POST http://localhost:8080/api/v1/monitors/1/pause

# Resume
curl -X POST http://localhost:8080/api/v1/monitors/1/resume
```

### Get Monitor Heartbeat

```bash
curl http://localhost:8080/api/v1/monitors/1/heartbeat
```

### Add Tag to Monitor

```bash
curl -X POST http://localhost:8080/api/v1/monitors/1/tags \
  -H "Content-Type: application/json" \
  -d '{
    "monitor_id": 1,
    "tag_id": 1,
    "value": "production"
  }'
```

### Add Maintenance to Monitor

```bash
curl -X POST http://localhost:8080/api/v1/monitors/1/maintenances \
  -H "Content-Type: application/json" \
  -d '{
    "monitor_id": 1,
    "maintenance_id": 1
  }'
```

## License

MIT