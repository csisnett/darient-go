# Restart Instructions for Logging Functionality

The logging functionality has been implemented, but you need to restart the Docker containers to use it.

## Quick Start

Run this command in your terminal:

```bash
./restart_with_logging.sh
```

This script will:
1. Stop the existing containers
2. Rebuild the Docker image with the new logging code
3. Start the containers
4. Test that everything is working

## Manual Steps (if you prefer)

If you want to do it manually:

```bash
# 1. Stop containers
sudo docker-compose down

# 2. Rebuild and start
sudo docker-compose up -d --build

# 3. Check status
sudo docker-compose ps

# 4. Test the API
curl http://localhost:8080/health
```

## Verify Logging is Working

After restarting, make some API calls:

```bash
# Make a few test requests
curl http://localhost:8080/health
curl http://localhost:8080/api/items
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","description":"Testing logging"}'
```

Then check the logs:

```bash
# View the log file
ls -la logs/
cat logs/api_*.log

# Or use the log viewer
./view_logs.sh latest
```

## Viewing Docker Logs

You can also view the application logs from Docker:

```bash
# View all logs
sudo docker-compose logs -f

# View only app logs
sudo docker-compose logs -f app

# View last 50 lines
sudo docker-compose logs --tail=50 app
```

## What Changed

The following components were added:
- `internal/logger/logger.go` - Logging functionality
- `internal/middleware/logging.go` - HTTP middleware for automatic logging
- `main.go` - Updated to initialize logger and add middleware
- `internal/handlers/handlers.go` - Added error logging
- `docker-compose.yml` - Added volume mount for logs directory

## Troubleshooting

If you don't see logs being created:

1. **Check if the logs directory exists:**
   ```bash
   ls -la logs/
   ```

2. **Check Docker container logs:**
   ```bash
   sudo docker-compose logs app | grep -i log
   ```

3. **Verify the container is running the new code:**
   ```bash
   sudo docker-compose ps
   ```

4. **Rebuild without cache:**
   ```bash
   sudo docker-compose down
   sudo docker-compose build --no-cache
   sudo docker-compose up -d
   ```

## Next Steps

Once the server is restarted with logging:
- All API requests will be automatically logged
- Logs are saved to `logs/api_YYYY-MM-DD.log`
- Use `./view_logs.sh` to analyze logs
- Check the README.md for full logging documentation