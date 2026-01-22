# 測試環境配置說明

## Windows 環境測試指南

### 啟用 CGO (必需)

由於本專案使用 SQLite 作為測試資料庫，必須啟用 CGO 才能運行測試。

#### 方法 1: 環境變數設置
```bash
set CGO_ENABLED=1
go test ./testing/test/... -v
```

#### 方法 2: TCC-GCC (推薦)

下載並安裝 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) 或 [MinGW-w64](https://www.mingw-w64.org/)。

安裝後確認：
```bash
gcc --version
```

#### 方法 3: 使用 Git Bash 或 WSL

如果使用 Git Bash 或 Windows Subsystem for Linux (WSL)，可以直接運行：
```bash
export CGO_ENABLED=1
go test ./testing/test/... -v
```

## 運行特定測試

### 單一測試
```bash
go test ./testing/test/... -v -run TestCenterRepository_CRUD
```

### 所有測試
```bash
go test ./testing/test/... -v
```

### 特定目錄測試
```bash
go test ./testing/test/center_test.go -v
```

## 測試數據庫

測試使用：
- **SQLite (In-memory)**: 作為主資料庫 mock
- **MinRedis**: 作為 Redis mock

測試完成後不會保留任何數據。

## 故障排除

### 錯誤: "Binary was compiled with 'CGO_ENABLED=0'"

**解決方案**: 設置環境變數 `CGO_ENABLED=1`

### 錯誤: "gcc: not found"

**解決方案**: 安裝 GCC 編譯器（TDM-GCC 或 MinGW-w64）

### 錯誤: SQLite 相關錯誤

**解決方案**: 確保已安裝 GCC 並啟用 CGO

## Docker 環境 (替代方案)

如果 Windows 環境配置困難，可以使用 Docker：

```bash
# 構建測試鏡像
docker build -t timeledger-test -f Dockerfile.test .

# 運行測試
docker run --rm timeledger-test
```
