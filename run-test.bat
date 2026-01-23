@echo off
chcp 65001 >nul
echo ╔════════════════════════════════════════════════════════════════╗
echo ║                  TimeLedger 測試腳本                        ║
echo ╚════════════════════════════════════════════════════════════════╝
echo.

if "%1"=="" (
    echo 🚀 運行所有測試...
    go test ./testing/test/... -v -count=1
) else if "%1"=="quick" (
    echo ⚡ 快速測試 (安靜模式)...
    go test ./testing/test/... -count=1 -short
) else if "%1"=="coverage" (
    echo 📊 運行測試並生成覆蓋率報告...
    go test ./testing/test/... -v -count=1 -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    echo ✅ 覆蓋率報告已生成: coverage.html
    if exist coverage.html (
        echo 📝 正在開啟覆蓋率報告...
        start coverage.html
    )
) else if "%1"=="report" (
    echo 📋 生成測試報告...
    go test ./testing/test/... -v -count=1 -xml=test-results.xml
    echo ✅ 測試報告已生成: test-results.xml
) else if "%1"=="clean" (
    echo 🧹 清理測試檔案...
    del /q coverage.out coverage.html test-results.xml test-output.txt 2>nul
    echo ✅ 清理完成
) else if "%1"=="build" (
    echo 🔨 建構專案...
    go build -mod=vendor -o main .
    echo ✅ 建構完成: main.exe
) else if "%1"=="help" (
    echo 📚 使用方式:
    echo   run-test.bat           - 運行所有測試
    echo   run-test.bat quick     - 快速測試 (安靜模式)
    echo   run-test.bat coverage - 生成覆蓋率報告
    echo   run-test.bat report    - 生成測試報告 (XML)
    echo   run-test.bat clean     - 清理測試檔案
    echo   run-test.bat build     - 建構專案
    echo   run-test.bat help      - 顯示幫助
) else (
    echo 🔍 運行特定測試: %1
    go test ./testing/test/... -v -count=1 -run "%1"
)

echo.
pause
