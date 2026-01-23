# Backend Application

A simple Go backend API server.

## Setup

```bash
go mod download
```

## Run

```bash
go run main.go
```

## Endpoints

- `GET /health` - Health check
- `GET /api/items` - Get all items
- `GET /api/items/{id}` - Get item by ID
- `POST /api/items` - Create new item
