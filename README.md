# Backend Application

A Go backend API server with PostgreSQL database featuring comprehensive API logging.

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

## API Logging

The application includes comprehensive logging for all API endpoints that captures detailed information about each request and response.

### Logging Features

- **Automatic Request Logging**: All API endpoints are automatically logged via middleware
- **Detailed Log Entries**: Each log entry includes:
  - Timestamp
  - HTTP method (GET, POST, PUT, DELETE)
  - Request path
  - HTTP status code
  - Response time in milliseconds
  - User agent
  - Remote IP address
  - Request body (for POST/PUT requests, limited to 10KB)
  - Response size
  - Error messages (when applicable)

### Log File Location

- Logs are stored in the `logs/` directory
- Log files are named with the format: `api_YYYY-MM-DD.log`
- Each log entry is stored as a JSON object on a separate line

### Log Entry Format

```json
{
  "timestamp": "2026-02-04T10:30:45Z",
  "method": "POST",
  "path": "/api/items",
  "status_code": 201,
  "response_time_ms": 45,
  "user_agent": "curl/7.68.0",
  "remote_addr": "127.0.0.1:54321",
  "request_body": "{\"name\":\"Test Item\",\"description\":\"A test item\"}",
  "response_size": 156
}
```

### Viewing Logs

#### Console Output
The application also displays colored log output in the console during development:
```
2026-02-04 10:30:45 [POST] /api/items - 201 (45ms)
```

#### Log File Analysis
Use the provided log viewer script to analyze logs:

```bash
# View latest 10 log entries
./view_logs.sh latest

# View latest 20 log entries
./view_logs.sh latest 20

# Show only error entries (4xx and 5xx status codes)
./view_logs.sh errors

# Show request statistics
./view_logs.sh stats

# Filter by endpoint path
./view_logs.sh filter /api/items
```

#### Manual Log Analysis
You can also analyze logs manually using standard tools:

```bash
# View latest log file
tail -f logs/api_$(date +%Y-%m-%d).log

# Count requests by status code
grep -o '"status_code":[0-9]*' logs/api_*.log | cut -d: -f2 | sort | uniq -c

# Find all error requests
grep '"status_code":[45][0-9][0-9]' logs/api_*.log

# Calculate average response time
grep -o '"response_time_ms":[0-9]*' logs/api_*.log | cut -d: -f2 | awk '{sum+=$1; count++} END {print "Average:", sum/count "ms"}'
```

### Error Logging

In addition to request/response logging, the application logs detailed error information:
- Database connection errors
- SQL query failures
- JSON parsing errors
- Validation errors
- Any unexpected application errors

Error logs include the specific error message and context to help with debugging.

## Endpoints

### Clients
- `POST /api/clients` - Create new client
- `GET /api/clients/{id}` - Get client by ID
- `PUT /api/clients/{id}` - Update client
- `DELETE /api/clients/{id}` - Delete client

### Banks
- `GET /api/banks` - Get all banks
- `POST /api/banks` - Create new bank
- `GET /api/banks/{id}` - Get bank by ID
- `PUT /api/banks/{id}` - Update bank
- `DELETE /api/banks/{id}` - Delete bank

### Credits
- `GET /api/credits` - Get all credits
- `POST /api/credits` - Create new credit
- `GET /api/credits/{id}` - Get credit by ID
- `PUT /api/credits/{id}` - Update credit
- `DELETE /api/credits/{id}` - Delete credit
- `GET /api/clients/{clientId}/credits` - Get credits by client
- `GET /api/banks/{bankId}/credits` - Get credits by bank

## Example Requests

```

### Create Client
```bash
curl -X POST http://localhost:8080/api/clients \
  -H "Content-Type: application/json" \
  -d '{"full_name":"John Doe","email":"john@example.com","birth_date":"1990-01-01T00:00:00Z","country":"USA"}'
```

### Create Bank
```bash
curl -X POST http://localhost:8080/api/banks \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Bank","type":"PRIVATE"}'
```
