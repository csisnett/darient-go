# Backend Application

A Go backend API server with PostgreSQL database.

## Setup with Docker (Recommended)

1. Start the application and database:
```bash
docker-compose up -d
```

2. View logs:
```bash
docker-compose logs -f
```

3. Stop the application:
```bash
docker-compose down
```

The application will be available at `http://localhost:8080` and will automatically create the required tables on startup.

## Local Development Setup

1. Start only the database:
```bash
docker-compose up -d db
```

2. Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

3. Install dependencies:
```bash
go mod download
```

4. Run the application locally:
```bash
go run main.go
```

## Endpoints

- `GET /health` - Health check
- `GET /api/items` - Get all items
- `GET /api/items/{id}` - Get item by ID
- `POST /api/items` - Create new item

## Example Request

```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Item","description":"A test item"}'
```
