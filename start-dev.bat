@echo off
echo ========================================
echo TimeLedger Docker Development Environment
echo ========================================

set SCRIPT_DIR=%~dp0
cd /d "%SCRIPT_DIR%"

echo.
echo [1/5] Check Docker status...
docker info >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo Docker is running, starting services...
    docker compose up -d mysql redis rabbitmq phpmyadmin
    echo.
    echo Waiting for MySQL to be ready...
    timeout /t 15 /nobreak >nul
    echo MySQL is ready!
) else (
    echo Docker is not running. Please start Docker Desktop first.
    echo Or ensure you have MySQL (127.0.0.1:3306) and Redis (127.0.0.1:6379) running locally.
    pause
    exit /b 1
)

echo.
echo [2/5] Copy environment variables...
if not exist .env.local copy .env.local.example .env.local >nul
echo Done!

echo.
echo [3/5] Start backend service...
echo Please run in another terminal: go run main.go
echo.

echo [4/5] Frontend development (optional)...
echo cd frontend && npm run dev
echo.

echo [5/5] Service ports:
echo ========================================
echo   - Backend API:  http://localhost:8888
echo   - phpMyAdmin:   http://localhost:8080
echo     User: root, Password: timeledger_root_2026
echo   - Redis:        localhost:6379
echo     Password: timeledger_redis_2026
echo   - RabbitMQ:     http://localhost:15672
echo     User: guest, Password: guest
echo ========================================

echo.
echo To run tests:
echo   go test ./testing/test/... -v
echo.

pause
