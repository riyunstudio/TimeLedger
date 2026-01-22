@echo off
REM TimeLedger Deployment Script for Windows

echo ========================================
echo TimeLedger VPS Deployment Script
echo ========================================
echo.

REM Configuration
set COMPOSE_FILE=docker-compose.yml
set ENV_FILE=.env.production

REM Check if .env.production exists
if not exist "%ENV_FILE%" (
    echo [ERROR] %ENV_FILE% not found!
    echo Please copy .env.production.example to %ENV_FILE% and update the values.
    exit /b 1
)

REM Pull latest code (optional)
echo Pulling latest code...
REM git pull origin main

REM Build and start containers
echo Building Docker images...
docker-compose -f %COMPOSE_FILE% --env-file %ENV_FILE% build

echo Starting containers...
docker-compose -f %COMPOSE_FILE% --env-file %ENV_FILE% up -d

REM Wait for services to be healthy
echo Waiting for services to be healthy...
timeout /t 30 /nobreak

REM Check container status
echo Container status:
docker-compose -f %COMPOSE_FILE% --env-file %ENV_FILE% ps

echo.
echo ========================================
echo Deployment completed successfully!
echo ========================================
echo.
echo Access the application at:
echo   - API: http://your-server-ip:8080
echo   - Swagger UI: http://your-server-ip:8080/swagger/index.html
echo   - RabbitMQ Management: http://your-server-ip:15672
echo.
echo To view logs:
echo   docker-compose -f %COMPOSE_FILE% --env-file %ENV_FILE% logs -f app
echo.
