#!/bin/bash

echo "=== Checking Docker Services ==="

# Check if Docker is running
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed or not in PATH"
    exit 1
fi

# Check Docker daemon
if ! docker info &> /dev/null; then
    echo "❌ Docker daemon is not running"
    echo "Please start Docker daemon first"
    exit 1
fi

echo "✅ Docker daemon is running"

# Check if docker-compose is available
if command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
elif docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
else
    echo "❌ Neither docker-compose nor 'docker compose' is available"
    exit 1
fi

echo "✅ Docker Compose is available: $COMPOSE_CMD"

# Check running containers
echo ""
echo "=== Current Docker Containers ==="
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Check if our specific containers are running
echo ""
echo "=== Checking Backend Services ==="

if docker ps --format "{{.Names}}" | grep -q "backend_db"; then
    echo "✅ Database container (backend_db) is running"
else
    echo "❌ Database container (backend_db) is not running"
fi

if docker ps --format "{{.Names}}" | grep -q "backend_app"; then
    echo "✅ App container (backend_app) is running"
else
    echo "❌ App container (backend_app) is not running"
fi

# Check if port 8080 is occupied
echo ""
echo "=== Checking Port 8080 ==="
if ss -tlnp | grep -q ":8080"; then
    echo "✅ Something is listening on port 8080"
    ss -tlnp | grep ":8080"
else
    echo "❌ Nothing is listening on port 8080"
fi

# Test health endpoint
echo ""
echo "=== Testing Health Endpoint ==="
if curl -s -f http://localhost:8080/health > /dev/null; then
    echo "✅ Health endpoint is responding"
    curl -s http://localhost:8080/health
else
    echo "❌ Health endpoint is not responding"
fi

echo ""
echo "=== Recommendations ==="
echo "If containers are not running, try:"
echo "  $COMPOSE_CMD up -d"
echo ""
echo "If you want to rebuild and restart:"
echo "  $COMPOSE_CMD down"
echo "  $COMPOSE_CMD up --build -d"
echo ""
echo "To check logs:"
echo "  $COMPOSE_CMD logs app"
echo "  $COMPOSE_CMD logs db"