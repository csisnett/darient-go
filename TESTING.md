# API Testing Guide

This document explains how to test the different API endpoints in this project.

## Available Test Methods

### 1. Unit Tests

There are two types of unit tests:

#### A. Pure Unit Tests (`internal/handlers/handlers_unit_test.go`)
These tests focus on input validation and error handling without requiring a database.

**Run pure unit tests:**
```bash
go test ./internal/handlers -run "Unit" -v
```

#### B. Database-dependent Tests (`internal/handlers/handlers_unit_test.go`)
These tests require a database connection and will fail if the database is not available.

**Run all unit tests:**
```bash
go test ./internal/handlers -v
```

**Note:** The database-dependent tests will fail with nil pointer errors if no database is connected. Use the integration tests for full database testing.

### 2. Integration Tests (`integration_test.go`)

These tests run against a real database and test the complete API workflow.

**Setup:**
1. Set up a test database (PostgreSQL)
2. Set the `TEST_DATABASE_URL` environment variable:
   ```bash
   export TEST_DATABASE_URL="postgres://username:password@localhost/test_db?sslmode=disable"
   ```

**Run integration tests:**
```bash
go test -v
```

**What it tests:**
- Complete CRUD operations for all entities
- Database interactions
- API endpoint responses
- Data validation

### 3. Manual Testing with Shell Script (`test_api.sh`)

An automated shell script that tests all endpoints using curl.

**Prerequisites:**
- Server must be running on `localhost:8080`
- `curl` must be installed

**Run the script:**
```bash
# Start your server first
go run main.go

# In another terminal, run the test script
./test_api.sh
```

**What it tests:**
- All CRUD operations
- Error cases
- Invalid input validation
- Response formats

### 4. Postman Collection (`postman_collection.json`)

Import this collection into Postman for interactive API testing.

**How to use:**
1. Import `postman_collection.json` into Postman
2. Set the `baseUrl` variable to your server URL (default: `http://localhost:8080`)
3. Run individual requests or the entire collection

## API Endpoints Overview

### Health Check
- `GET /health` - Check if the server is running

### Items
- `GET /api/items` - Get all items
- `POST /api/items` - Create a new item
- `GET /api/items/{id}` - Get item by ID

### Clients
- `POST /api/clients` - Create a new client
- `GET /api/clients/{id}` - Get client by ID
- `PUT /api/clients/{id}` - Update client
- `DELETE /api/clients/{id}` - Delete client

### Banks
- `GET /api/banks` - Get all banks
- `POST /api/banks` - Create a new bank
- `GET /api/banks/{id}` - Get bank by ID
- `PUT /api/banks/{id}` - Update bank
- `DELETE /api/banks/{id}` - Delete bank

### Credits
- `GET /api/credits` - Get all credits
- `POST /api/credits` - Create a new credit
- `GET /api/credits/{id}` - Get credit by ID
- `PUT /api/credits/{id}` - Update credit
- `DELETE /api/credits/{id}` - Delete credit
- `GET /api/clients/{clientId}/credits` - Get credits by client
- `GET /api/banks/{bankId}/credits` - Get credits by bank

## Sample Data Formats

### Client
```json
{
  "full_name": "John Doe",
  "email": "john.doe@example.com",
  "birth_date": "1990-01-01T00:00:00Z",
  "country": "USA"
}
```

### Bank
```json
{
  "name": "Test Bank",
  "type": "PRIVATE"
}
```
Valid types: `PRIVATE`, `GOVERNMENT`

### Credit
```json
{
  "client_id": 1,
  "bank_id": 1,
  "min_payment": 100.0,
  "max_payment": 1000.0,
  "term_months": 12,
  "credit_type": "AUTO",
  "status": "PENDING"
}
```
Valid credit types: `AUTO`, `MORTGAGE`, `COMMERCIAL`
Valid statuses: `PENDING`, `APPROVED`, `REJECTED`

### Item
```json
{
  "name": "Test Item",
  "description": "This is a test item"
}
```

## Running the Server

Before testing, make sure your server is running:

```bash
# Set up environment variables
cp .env.example .env
# Edit .env with your database configuration

# Run the server
go run main.go
```

The server will start on port 8080 by default.

## Database Setup

Make sure you have PostgreSQL running and create the necessary databases:

```sql
-- For development
CREATE DATABASE your_db_name;

-- For testing
CREATE DATABASE your_test_db_name;
```

The application will automatically run migrations when it starts.

## Troubleshooting

### Common Issues

1. **Database connection errors**: Make sure PostgreSQL is running and the connection string is correct
2. **Port already in use**: Change the PORT environment variable or stop other services using port 8080
3. **Permission denied on test_api.sh**: Run `chmod +x test_api.sh` to make it executable

### Test Data Cleanup

The integration tests automatically clean up test data. For manual testing, you can uncomment the cleanup section in `test_api.sh` or manually delete test records from the database.