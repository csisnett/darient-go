#!/bin/bash

# Test script to demonstrate API logging functionality
# This script makes several API calls and then shows the logs

echo "=== Testing API Logging Functionality ==="
echo

# Check if server is running
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ Server is not running. Please start the server first:"
    echo "   docker-compose up -d"
    echo "   OR"
    echo "   go run main.go"
    exit 1
fi

echo "✅ Server is running"
echo

# Make some test API calls
echo "📡 Making test API calls..."

# Health check
echo "1. Health check..."
curl -s http://localhost:8080/health > /dev/null

# Get items (should be empty initially)
echo "2. Getting items..."
curl -s http://localhost:8080/api/items > /dev/null

# Create an item
echo "3. Creating an item..."
curl -s -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Item","description":"A test item for logging"}' > /dev/null

# Get the item we just created
echo "4. Getting items again..."
curl -s http://localhost:8080/api/items > /dev/null

# Try to get a non-existent item (should generate 404)
echo "5. Trying to get non-existent item (will generate 404)..."
curl -s http://localhost:8080/api/items/999 > /dev/null

# Try to create item with invalid JSON (should generate 400)
echo "6. Trying to create item with invalid JSON (will generate 400)..."
curl -s -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"invalid": json}' > /dev/null

echo
echo "✅ Test API calls completed!"
echo

# Wait a moment for logs to be written
sleep 1

# Show the logs
echo "📋 Recent log entries:"
echo "===================="

if [ -f "./view_logs.sh" ]; then
    ./view_logs.sh latest 10
else
    echo "❌ view_logs.sh not found. Showing raw log file instead:"
    if [ -d "logs" ]; then
        latest_log=$(ls -t logs/api_*.log 2>/dev/null | head -1)
        if [ -n "$latest_log" ]; then
            echo "Latest log file: $latest_log"
            echo "Last 10 entries:"
            tail -10 "$latest_log"
        else
            echo "No log files found in logs/ directory"
        fi
    else
        echo "No logs directory found"
    fi
fi

echo
echo "🎉 Logging test completed!"
echo
echo "💡 You can also use these commands to analyze logs:"
echo "   ./view_logs.sh latest     - Show latest 10 entries"
echo "   ./view_logs.sh errors     - Show only errors"
echo "   ./view_logs.sh stats      - Show statistics"
echo "   ./view_logs.sh filter /api/items - Filter by endpoint"