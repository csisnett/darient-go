#!/bin/bash

echo "=== Checking Logging Functionality ==="
echo

# Check if logs directory exists
if [ -d "logs" ]; then
    echo "✅ Logs directory exists"
    
    # Check if there are any log files
    log_count=$(ls -1 logs/api_*.log 2>/dev/null | wc -l)
    if [ "$log_count" -gt 0 ]; then
        echo "✅ Found $log_count log file(s)"
        
        # Show the latest log file
        latest_log=$(ls -t logs/api_*.log 2>/dev/null | head -1)
        echo "📁 Latest log file: $latest_log"
        
        # Count entries in the latest log
        if [ -f "$latest_log" ]; then
            entry_count=$(wc -l < "$latest_log")
            echo "📊 Log entries: $entry_count"
            
            if [ "$entry_count" -gt 0 ]; then
                echo
                echo "📋 Last 5 log entries:"
                echo "===================="
                tail -5 "$latest_log" | while read line; do
                    # Pretty print JSON if jq is available
                    if command -v jq &> /dev/null; then
                        echo "$line" | jq -r '"\(.timestamp) [\(.method)] \(.path) - \(.status_code) (\(.response_time_ms)ms)"'
                    else
                        echo "$line"
                    fi
                done
                echo
                echo "✅ Logging is working!"
            else
                echo "⚠️  Log file is empty. Make some API requests to generate logs."
            fi
        fi
    else
        echo "❌ No log files found in logs/ directory"
        echo "   The server may not have been restarted with the new logging code."
        echo "   Run: ./restart_with_logging.sh"
    fi
else
    echo "❌ Logs directory does not exist"
    echo "   Creating logs directory..."
    mkdir -p logs
    echo "   Run: ./restart_with_logging.sh"
fi

echo
echo "🔍 Checking Docker container..."
if command -v docker-compose &> /dev/null; then
    if sudo docker-compose ps 2>/dev/null | grep -q "backend_app"; then
        echo "✅ Docker container is running"
        
        # Check if logs directory exists in container
        echo
        echo "📦 Checking logs in Docker container..."
        sudo docker-compose exec -T app ls -la logs/ 2>/dev/null || echo "   Could not access container logs directory"
    else
        echo "⚠️  Docker container is not running"
        echo "   Run: ./restart_with_logging.sh"
    fi
else
    echo "⚠️  docker-compose not found or requires sudo"
fi

echo
echo "💡 Quick commands:"
echo "   - Restart with logging: ./restart_with_logging.sh"
echo "   - View logs: ./view_logs.sh latest"
echo "   - Test logging: ./test_logging.sh"