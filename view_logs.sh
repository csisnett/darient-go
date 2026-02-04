#!/bin/bash

# API Log Viewer Script
# Usage: ./view_logs.sh [command] [options]

cd "$(dirname "$0")"

if [ ! -d "logs" ]; then
    echo "No logs directory found. Make sure the server has been started and has generated some logs."
    exit 1
fi

# Check if any log files exist
if [ ! "$(ls -A logs/api_*.log 2>/dev/null)" ]; then
    echo "No API log files found in logs/ directory."
    exit 1
fi

# Run the Go log viewer
go run scripts/view_logs.go "$@"