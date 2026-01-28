# TimeLedger 教師排課平台

**TimeLedger** 是一款教師中心化多據點排課平台，專為台灣市場設計，深度整合 LINE 生態系。

## 專案特色

- **教師端**：LINE 單一登入（無密碼策略），透過 LIFF 行動網頁授課
- **管理員後台**：完整的排課管理、例外審核、人才庫功能
- **智慧媒合**：AI 輔助教師媒合與替代時段推薦
- **即時通知**：透過 LINE Bot 發送例外申請通知

## 技術堆疊

| 層面 | 技術 |
|:---|:---|
| **後端** | Go (Gin) + MySQL 8.0 + Redis |
| **前端** | Nuxt 3 (SSR) + Tailwind CSS + LINE LIFF |
| **部署** | Docker Compose（單一 VPS 容器化部署） |
| **通訊** | HTTP REST API + LINE Bot Webhook |

## 快速開始

### 前置條件

- Docker & Docker Compose
- Go 1.21+
- Node.js 18+

### 本地開發環境

1. **複製專案**
```bash
git clone https://github.com/your-org/timeLedger.git
cd timeLedger
```

2. **啟動資料庫服務**
```bash
docker-compose up -d mysql redis
```

3. **設定環境變數**
```bash
cp .env.example .env
# 編輯 .env 檔案
```

4. **執行後端**
```bash
go run main.go
```

5. **執行前端**
```bash
cd frontend
npm install
npm run dev
```

### Docker 完整部署

```bash
# 啟動所有服務
docker-compose up -d

# 查看服務狀態
docker-compose ps

# 查看日誌
docker-compose logs -f app
```

## API 文件

啟動伺服器後，可透過以下路徑存取互動式 Swagger API 文件：

```
http://localhost:8888/swagger/index.html
```

### 主要 API 分類

| 分類 | 端點數量 | 說明 |
|:---|:---:|:---|
| 認證（Auth） | 4 | 管理員登入、老師 LINE 登入 |
| 教師（Teacher） | 30+ | 教師檔案、課表、證照 |
| 管理員（Admin） | 40+ | 中心管理、排課、例外審核 |
| 排課（Scheduling） | 15+ | 排課規則、課程、教室 |
| 智慧媒合（Smart Matching） | 7 | 人才庫搜尋、媒合推薦 |
| 通知（Notification） | 8 | 佇列統計、通知管理 |
| 匯出（Export） | 4 | 課表匯出、報告產生 |

## 測試執行

### 執行所有測試

```bash
go test ./testing/test/... -v
```

### 執行特定測試

```bash
# 執行認證相關測試
go test ./testing/test/auth_test.go -v

# 執行管理員相關測試
go test ./testing/test/admin_user_test.go -v

# 執行智慧媒合測試
go test ./testing/test/smart_matching_test.go -v
```

### 測試覆蓋率

```
測試通過率：92.8% (26/28)
- admin_user_test.go：17/17 通過
- auth_test.go：11/11 通過
- smart_matching_test：13/13 通過
- center_invitation_test：6/6 通過
```

## LINE Bot 設定

### Webhook 設定

```
LINE Messaging API Webhook URL：
https://your-domain.com/line/webhook
```

### 支援的關鍵字

| 關鍵字 | 回覆 |
|:---|:---|
| `綁定` | 顯示綁定連結 |
| `幫助` | 使用說明 |
| `狀態` | 查詢綁定狀態 |
| `解除綁定` | 顯示解除綁定連結 |

### 通知類型

- **老師歡迎訊息**：首次登入或受邀請時發送
- **管理員歡迎訊息**：首次登入且未綁定時發送
- **例外申請通知**：老師提交例外時通知所有管理員
- **例外審核結果**：核准/拒絕時通知申請老師

## 專案結構

```
timeLedger/
├── main.go                    # 應用程式入口點
├── app/                       # 後端核心
│   ├── controllers/           # API 控制器
│   ├── services/              # 業務邏輯
│   ├── repositories/          # 資料存取層
│   ├── models/                # 資料模型
│   ├── requests/              # 請求驗證
│   ├── resources/             # API 回應格式化
│   └── servers/               # HTTP 伺服器配置
├── frontend/                  # Nuxt 3 前端
│   ├── pages/                 # 頁面元件
│   ├── components/            # 共用元件
│   └── composables/           # Vue Composables
├── database/                  # 資料庫遷移
│   └── mysql/
├── testing/                   # 測試程式
│   └── test/                  # 單元測試
├── docs/                      # API 文件
├── docker-compose.yml         # Docker 部署配置
└── CLAUDE.md                  # 開發規範文件
```

## 開發規範

請參考 [CLAUDE.md](CLAUDE.md) 了解：

- 分層架構規範
- 命名慣例
- 錯誤處理模式
- 排課驗證引擎
- 權限管控矩陣
- 測試策略

## 環境變數

主要環境變數：

| 變數 | 說明 | 預設值 |
|:---|:---|:---|
| `DB_HOST` | MySQL 主機 | `127.0.0.1` |
| `DB_PORT` | MySQL 連接埠 | `3306` |
| `DB_USER` | MySQL 使用者 | `root` |
| `DB_PASSWORD` | MySQL 密碼 | - |
| `REDIS_ADDR` | Redis 地址 | `127.0.0.1:6379` |
| `JWT_SECRET` | JWT 密鑰 | - |
| `LINE_CHANNEL_ID` | LINE Channel ID | - |
| `LINE_CHANNEL_SECRET` | LINE Channel Secret | - |
| `LINE_BOT_TOKEN` | LINE Bot Access Token | - |

## 監控與健康檢查

- **健康檢查端點**：`GET /healthy`
- **API 服務**：`http://localhost:8888`
- **Swagger UI**：`http://localhost:8888/swagger/index.html`

## 貢獻指南

1. Fork 本專案
2. 建立功能分支：`git checkout -b feature/xxx`
3. 遵循 [CLAUDE.md](CLAUDE.md) 中的開發規範
4. 確保所有測試通過
5. 提交 PR 前更新相關文件

## 授權

本專案為專有軟體，保留所有權利。
