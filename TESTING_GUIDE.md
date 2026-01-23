# TimeLedger 測試指南

## 快速開始

### Windows
```batch
# 運行所有測試
run-test.bat

# 快速測試
run-test.bat quick

# 生成覆蓋率報告
run-test.bat coverage

# 查看幫助
run-test.bat help
```

### Linux/macOS (使用 Makefile)
```bash
# 運行所有測試
make test

# 生成覆蓋率報告
make test-coverage

# 快速測試
make test-quick

# 查看幫助
make help
```

## 測試命令

| 命令 | 說明 |
|:---|:---|
| `make test` | 運行所有測試 |
| `make test-coverage` | 生成測試覆蓋率報告 |
| `make test-report` | 生成 JUnit XML 測試報告 |
| `make test-quick` | 快速測試（安靜模式） |
| `make test file=auth_test` | 運行特定測試檔案 |
| `make test name=TestAuth` | 運行特定測試名稱 |
| `make clean` | 清理測試產物 |

## GitHub Actions 自動測試

專案已配置 GitHub Actions，會在以下情況自動運行測試：

1. **Push 到 main 或 claudecode 分支**
2. **Pull Request 到 main 分支**

測試流程：
1. 啟動 MySQL 8.0 服務
2. 運行所有單元測試
3. 建構專案
4. 執行程式碼檢查 (golangci-lint)

## 資料庫配置

測試使用 MySQL 資料庫，配置如下：

```env
MYSQL_MASTER_HOST=127.0.0.1
MYSQL_MASTER_PORT=3307
MYSQL_MASTER_USER=root
MYSQL_MASTER_PASS=rootpassword
MYSQL_MASTER_NAME=timeledger_test
```

確保 MySQL 服務正在運行，並建立 `timeledger_test` 資料庫。

## 新增測試

### 1. 建立測試檔案
在 `testing/test/` 目錄下建立新的測試檔案，命名規範：`*_test.go`

### 2. 測試結構
```go
package test

import (
	"testing"
)

func TestYourFeature(t *testing.T) {
	// 測試代碼
}
```

### 3. 使用測試輔助函數
```go
// 初始化測試資料庫
db, err := InitializeTestDB()
if err != nil {
	t.Fatalf("Failed to connect to test database: %v", err)
}
defer CloseDB(db)

// 創建測試資料
teacher := &models.Teacher{
	Name: "Test Teacher",
}
db.Create(teacher)
```

## 常見問題

### Q: 測試失敗怎麼辦？
A: 
1. 檢查 MySQL 服務是否運行
2. 確認 `timeledger_test` 資料庫存在
3. 查看詳細錯誤輸出

### Q: 如何排除特定測試？
A: 使用 `-skip` 參數
```bash
go test ./testing/test/... -v -skip="TestAuth"
```

### Q: 測試資料重複怎麼辦？
A: 每次測試使用唯一的資料，或在 `setup` 函數中清理資料

## CI/CD 狀態

![CI/CD](https://github.com/your-org/timeledger/actions/workflows/test.yml/badge.svg)

## 覆蓋率目標

- **核心業務邏輯**: > 80%
- **排課驗證引擎**: > 90%
- **API 端點**: > 70%
