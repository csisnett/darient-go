#!/bin/bash

echo "=== Restarting Docker containers with logging functionality ==="
echo

# Stop existing containers
echo "1. Stopping existing containers..."
sudo docker-compose down

echo
echo "2. Rebuilding and starting containers..."
sudo docker-compose up -d --build

echo
echo "3. Waiting for services to start..."
sleep 5

# Check if services are running
echo
echo "4. Checking service status..."
sudo docker-compose ps

echo
echo "5. Testing the API..."
sleep 2
curl -s http://localhost:8080/health | jq .

echo
echo
echo "✅ Docker containers restarted with logging functionality!"
echo
echo "📋 To view logs:"
echo "   - Application logs: sudo docker-compose logs -f app"
echo "   - All logs: sudo docker-compose logs -f"
echo
echo "📁 API logs will be saved in the logs/ directory inside the container"
echo "   To access them: sudo docker-compose exec app ls -la logs/"
echo
echo "🧪 Test the logging:"
echo "   ./test_logging.sh"