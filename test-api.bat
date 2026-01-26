@echo off
echo Testing Admin Login API...
curl -X POST "http://localhost:3005/api/v1/auth/admin/login" -H "Content-Type: application/json" -d "{\"email\":\"admin@timeledger.com\",\"password\":\"admin123\"}"
echo.
echo Testing Get Pending Exceptions API (with mock token)...
curl -X GET "http://localhost:3005/api/v1/admin/scheduling/exceptions/pending" -H "Authorization: Bearer mock-admin-token-1"
echo.
pause
