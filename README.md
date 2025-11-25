# Akali 專案部署說明

## 專案概述
Akali 為 Kubernetes (GKE) 上運行的服務，主要結合 Cloud SQL 與外部 API，並提供健康檢查機制與後台管理介面。

---

## 基本資訊

| 項目 | 值 |
|------|----|
| **GKE Cluster 名稱** | `????` |
| **GCP 區域** | `????` |
| **Namespace** | `akali` |
| **服務訪問 IP** | [http://localhost/](http://localhost/) |
| **健康檢查路徑** | [http://localhost/healthy](http://localhost/healthy) |

---

## Cloud SQL 設定

| 項目 | 值 |
|------|----|
| **DB Host (Private IP)** | `XXX.XXX.XXX.XXX` |
| **DB 使用者** | `root` |
| **DB 密碼** | `??????` |
| **連線方式** | Private IP 連線 (同 VPC 內) |

---

## phpMyAdmin 操作說明

可透過 Port Forward 方式連線至 Kubernetes 內部的 phpMyAdmin 服務。

```bash

```

## proto編譯指令
```bash
protoc --go_out=./grpc --go-grpc_out=./grpc grpc/proto/user.proto
```

## pkg專案import指令
```bash
go env -w GOPRIVATE=gitlab.en.mcbwvx.com
```

## pkg列表
* Tools 工具包
    * https://gitlab.en.mcbwvx.com/frame/teemo#
* Curl 工具包
    * https://gitlab.en.mcbwvx.com/frame/ezreal#
* Log 自定義包
    * https://gitlab.en.mcbwvx.com/frame/zilean#

---

## 框架架構說明

### 專案結構

```
akali/
├── main.go                 # 應用程式入口點
├── app/                     # 應用層核心邏輯
│   ├── base.go             # App 初始化與依賴注入
│   ├── controllers/        # 控制器層（處理 HTTP 請求）
│   │   ├── base.go        # 控制器基類
│   │   └── user.go        # 使用者控制器範例
│   ├── services/           # 服務層（業務邏輯）
│   │   ├── base.go        # 服務基類
│   │   └── users.go       # 使用者服務範例
│   ├── repositories/       # 資料存取層（資料庫操作）
│   │   ├── base.go        # 倉儲基類
│   │   └── user.go        # 使用者倉儲範例
│   ├── models/             # 資料模型（GORM 模型）
│   │   └── user.go        # 使用者模型範例
│   ├── requests/           # 請求驗證層
│   │   ├── base.go        # 請求驗證基類（泛型驗證）
│   │   └── user.go        # 使用者請求驗證範例
│   ├── resources/          # 資源轉換層（API 回應格式化）
│   │   ├── base.go        # 資源基類
│   │   └── user.go        # 使用者資源範例
│   ├── servers/             # HTTP 伺服器（Gin）
│   │   ├── server.go      # 伺服器初始化與啟動
│   │   ├── route.go       # 路由註冊
│   │   ├── middleware.go  # 中介層（Init, Recover, Main）
│   │   └── other.go       # 輔助函數（TraceID, 日誌等）
│   └── console/            # 定時任務排程
│       ├── schedule.go    # 排程器管理
│       └── job.go         # 任務介面定義
├── configs/                # 配置管理
│   └── env.go             # 環境變數載入
├── database/               # 資料庫連線
│   ├── mysql/             # MySQL 主從連線
│   │   ├── conn.go        # 資料庫連線初始化
│   │   ├── migrate.go     # 資料庫遷移
│   │   └── seeder.go      # 資料庫種子資料
│   └── redis/              # Redis 連線
│       └── conn.go        # Redis 連線初始化
├── global/                 # 全域定義
│   ├── define.go          # Context Key 定義
│   ├── structs.go         # 全域結構體
│   └── errInfos/          # 錯誤碼定義
│       ├── code.go        # 錯誤碼常數
│       ├── message.go     # 錯誤訊息
│       └── base.go        # 錯誤處理基類
├── grpc/                   # gRPC 服務
│   ├── proto/              # Protocol Buffers 定義
│   └── services/           # gRPC 服務實作
├── libs/                   # 第三方庫封裝
│   ├── mq/                 # RabbitMQ 封裝
│   └── ws/                 # WebSocket 封裝
├── apis/                   # 外部 API 呼叫封裝
└── rpc/                    # RPC 呼叫封裝
```

### 架構模式

本框架採用 **分層架構（Layered Architecture）** 與 **倉儲模式（Repository Pattern）**，支援兩種通訊模式：

#### HTTP RESTful API
1. **Controller 層**：處理 HTTP 請求，負責參數驗證與回應格式化
2. **Service 層**：處理業務邏輯，協調 Repository 與 Resource
3. **Repository 層**：負責資料庫操作，封裝 GORM 查詢
4. **Model 層**：定義資料模型結構（GORM）
5. **Request 層**：請求參數驗證與轉換
6. **Resource 層**：將 Model 轉換為 API 回應格式

#### gRPC 服務
1. **gRPC Service 層**：處理 gRPC 請求，實作 Protocol Buffers 定義的服務
2. **Service 層**：處理業務邏輯，協調 Repository
3. **Repository 層**：負責資料庫操作，封裝 GORM 查詢
4. **Model 層**：定義資料模型結構（GORM）
5. **Protocol Buffers**：定義服務介面與訊息格式

### 資料流程

#### HTTP RESTful API 流程

```
HTTP Request
    ↓
Middleware (TraceID, Body備份, 日誌)
    ↓
Controller (參數驗證 via Request)
    ↓
Request (驗證請求參數)
    ↓
Service (業務邏輯)
    ↓
Repository (資料庫操作)
    ↓
Model (GORM)
    ↓
Service (資料處理)
    ↓
Resource (格式化回應)
    ↓
Controller (JSON 回應)
    ↓
Middleware (記錄日誌)
    ↓
HTTP Response
```

#### gRPC 服務流程

```
gRPC Request
    ↓
InitMiddleware (WaitGroup, RequestTime)
    ↓
RecoverMiddleware (Panic 保護)
    ↓
MainMiddleware (Metadata注入, TraceID, 超時控制)
    ↓
gRPC Service (業務邏輯)
    ↓
Repository (資料庫操作)
    ↓
Model (GORM)
    ↓
gRPC Service (資料處理)
    ↓
gRPC Response
    ↓
MainMiddleware (記錄日誌)
    ↓
gRPC Response
```

---

## 開發規範與注意事項

### 1. 命名規範

- **檔案命名**：使用小寫字母，多單字用底線分隔（如 `user_controller.go`）
- **結構體命名**：使用大寫開頭的駝峰式（如 `UserController`）
- **方法命名**：公開方法使用大寫開頭（如 `Get`, `Create`），私有方法使用小寫開頭
- **常數命名**：全大寫，底線分隔（如 `SYSTEM_ERROR`）

### 2. 程式碼風格

- **縮排**：使用 Tab（`\t`），不使用空格
- **註解**：公開函數必須有註解說明
- **錯誤處理**：所有可能出錯的操作都必須處理錯誤
- **Panic 處理**：使用 `recover` 機制防止服務崩潰

### 3. 資料庫操作規範

#### 主從分離
- **讀取操作**：使用 `app.Mysql.RDB`（從庫）
- **寫入操作**：使用 `app.Mysql.WDB`（主庫）
- **範例**：
```go
// 讀取
user, err := rp.app.Mysql.RDB.WithContext(ctx).Model(&rp.model).Find(&data).Error

// 寫入
err := rp.app.Mysql.WDB.WithContext(ctx).Model(&rp.model).Create(&data).Error
```

#### Context 傳遞
- Repository 層必須接收 `context.Context` 參數
- 從 Gin Context 轉換：使用 `BaseService.dbCtx(ctx)` 方法
- 確保 TraceID 正確傳遞到資料庫層

### 4. 錯誤處理規範

#### 錯誤碼定義規則
- 格式：`專案流水號(2位) + 功能類型(2位) + 流水號(4位)`
- 範例：`110001` = 專案1 + 系統相關(1) + 錯誤001
- 功能類型：
  - `1`：系統相關
  - `2`：資料庫、快取相關
  - `3`：其他相關
  - `4`：使用者相關

#### 錯誤回傳格式
- Controller 層統一使用 `BaseController.JSON()` 方法
- 回傳格式：
```go
ctl.JSON(ctx, global.Ret{
    Status:  http.StatusOK,
    Datas:   datas,
    ErrInfo: errInfo,  // nil 表示成功
    Err:     err,      // 實際錯誤物件
})
```

### 5. 請求驗證規範

- 使用 `requests.Validate[T]()` 泛型函數進行驗證
- GET/DELETE 請求：自動從 Query 參數驗證
- POST/PUT 請求：優先驗證 JSON，失敗則嘗試 Form Data
- 使用 `binding` tag 進行欄位驗證：
```go
type UserCreateRequest struct {
    Name string `json:"name" binding:"required"`
    Ips  []string `json:"ips"`
}
```

### 6. 路由註冊規範

- 在 `app/servers/route.go` 的 `LoadRoutes()` 方法中註冊路由
- 路由結構：
```go
{
    Method:      http.MethodGet,
    Path:        "/user",
    Controller:  s.action.user.Get,
    Middlewares: []gin.HandlerFunc{},  // 可選的中介層
}
```

### 7. Swagger 註解規範

- 所有 API 端點必須添加 Swagger 註解
- 範例：
```go
// @Summary 查詢使用者
// @description
// @Tags User
// @Param Content-Type header string true "Content-Type" default(application/json)
// @Param Tid header string false "TraceID"
// @Param ID query int true "會員ID"
// @Success 200 {object} global.ApiResponse{datas=resources.UserGetResource} "回傳"
// @Router /user [get]
func (ctl *UserController) Get(ctx *gin.Context) {
    // ...
}
```

### 8. gRPC 服務開發規範

#### Protocol Buffers 定義
- 在 `grpc/proto/` 目錄下定義 `.proto` 檔案
- 範例結構：
```protobuf
syntax = "proto3";

package user;
option go_package = "./proto/user";

service UserService {
  rpc Get (GetRequest) returns (GetResponse);
}

message GetRequest {
  int64 ID = 1;
}

message GetResponse {
  int64 Code = 1;
  string Msg = 2;
  GetResponseDatas Datas = 3;
}
```

#### 編譯 Protocol Buffers
```bash
protoc --go_out=./grpc --go-grpc_out=./grpc grpc/proto/user.proto
```

#### gRPC 服務實作
- 在 `grpc/services/` 目錄下實作服務
- 必須嵌入 `BaseService` 和 `UnimplementedXxxServiceServer`
- 範例：
```go
type User struct {
    user.UnimplementedUserServiceServer
    BaseService
    App            *app.App
    UserRepository *repositories.UserRepository
}

func (s *User) Get(ctx context.Context, req *user.GetRequest) (*user.GetResponse, error) {
    data, err := s.UserRepository.Get(ctx, models.User{ID: uint(req.GetID())})
    if err != nil {
        return &user.GetResponse{Code: 100, Msg: err.Error()}, err
    }
    return &user.GetResponse{
        Msg: "OK",
        Datas: &user.GetResponseDatas{
            ID:   int64(data.ID),
            Name: data.Name,
        },
    }, nil
}
```

#### 註冊 gRPC 服務
- 在 `grpc/server.go` 的 `registerServices()` 方法中註冊：
```go
func (s *Grpc) registerServices() {
    user.RegisterUserServiceServer(s.srv, &services.User{
        App:            s.app,
        UserRepository: repositories.NewUserRepository(s.app),
    })
}
```

#### gRPC 中介層
- **InitMiddleware**：初始化 WaitGroup 和 RequestTime
- **RecoverMiddleware**：Panic 保護與錯誤記錄
- **MainMiddleware**：Metadata 注入、TraceID 處理、超時控制、日誌記錄

#### gRPC 超時處理
- 預設超時時間建議：5 秒
- 自動處理 Context 取消與超時

#### gRPC Metadata 與 TraceID
- TraceID 從 Metadata 的 `Tid` 欄位取得，若無則自動生成
- Metadata 會自動注入到 Context 中
- 所有日誌都包含 TraceID 用於追蹤

### 9. 定時任務開發規範

- 實作 `console.Job` 介面：
```go
type Job interface {
    Name() string        // 任務名稱
    Description()        // 任務說明
    Repositories()       // 初始化需要的 Repository
    Handle(string) error // 主程式（參數為 TraceID）
}
```
- 在 `app/console/schedule.go` 的 `loadJobs()` 方法中註冊任務
- 使用 Cron 表達式：`秒 分 時 日 月 星期 * * * * * *`

### 10. 環境變數配置

- 所有配置透過環境變數載入（`.env` 檔案）
- 使用 `github.com/joho/godotenv/autoload` 自動載入
- 在 `configs/env.go` 中定義配置結構
- 必須設定的環境變數請參考 `configs/env.go`

### 11. 日誌記錄規範

- 使用 `gitlab.en.mcbwvx.com/frame/zilean` 進行日誌記錄
- API 請求日誌：自動在 `MainMiddleware` 中記錄
- 資料庫日誌：GORM 自動記錄錯誤 SQL（僅 Debug 模式）
- 排程任務日誌：在 `Scheduler` 中自動記錄
- 所有日誌都包含 TraceID 用於追蹤

### 12. 優雅退出機制

- 所有服務都實作 `Start()` 和 `Stop()` 方法
- 使用 `sync.WaitGroup` 等待所有請求完成
- 在 `main.go` 中統一處理 SIGTERM/SIGINT 信號
- 退出順序：Scheduler → RabbitMQ → WebSocket → Gin → gRPC

### 13. TraceID 追蹤

- 每個請求自動生成或從 Header 取得 TraceID
- TraceID 會傳遞到：
  - Gin Context
  - 資料庫操作 Context
  - 日誌記錄
  - 外部 API 呼叫
- 用於追蹤整個請求生命週期

---

## 開發流程範例

### 新增一個 HTTP RESTful API 端點

1. **定義 Model**（`app/models/xxx.go`）
```go
type Xxx struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"type:varchar(30)"`
}
```

2. **定義 Request**（`app/requests/xxx.go`）
```go
type XxxCreateRequest struct {
    Name string `json:"name" binding:"required"`
}

func (r *XxxRequest) Create(ctx *gin.Context) (*XxxCreateRequest, *errInfos.Res, error) {
    req, err := Validate[XxxCreateRequest](ctx)
    if err != nil {
        return nil, r.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), err
    }
    return req, nil, nil
}
```

3. **定義 Repository**（`app/repositories/xxx.go`）
```go
func (rp *XxxRepository) Create(ctx context.Context, data models.Xxx) (models.Xxx, error) {
    err := rp.app.Mysql.WDB.WithContext(ctx).Model(&rp.model).Create(&data).Error
    return data, err
}
```

4. **定義 Resource**（`app/resources/xxx.go`）
```go
type XxxCreateResource struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

func (rs *XxxResource) Create(ctx *gin.Context, datas models.Xxx) (*XxxCreateResource, error) {
    return &XxxCreateResource{
        ID:   datas.ID,
        Name: datas.Name,
    }, nil
}
```

5. **定義 Service**（`app/services/xxx.go`）
```go
func (s *XxxService) Create(ctx *gin.Context, req *requests.XxxCreateRequest) (any, *errInfos.Res, error) {
    data, err := s.xxxRepository.Create(ctx, models.Xxx{Name: req.Name})
    if err != nil {
        return nil, s.app.Err.New(errInfos.SQL_ERROR), err
    }
    
    response, err := s.xxxResource.Create(ctx, data)
    if err != nil {
        return nil, s.app.Err.New(errInfos.FORMAT_RESOURCE_ERROR), err
    }
    
    return response, nil, nil
}
```

6. **定義 Controller**（`app/controllers/xxx.go`）
```go
// @Summary 新增XXX
// @Tags Xxx
// @Param Params body requests.XxxCreateRequest true "參數"
// @Success 200 {object} global.ApiResponse{datas=resources.XxxCreateResource}
// @Router /xxx [post]
func (ctl *XxxController) Create(ctx *gin.Context) {
    req, eInfo, err := ctl.XxxRequest.Create(ctx)
    if err != nil {
        ctl.JSON(ctx, global.Ret{Status: http.StatusBadRequest, ErrInfo: eInfo, Err: err})
        return
    }
    
    datas, eInfo, err := ctl.XxxService.Create(ctx, req)
    if err != nil {
        ctl.JSON(ctx, global.Ret{Status: http.StatusInternalServerError, ErrInfo: eInfo, Err: err})
        return
    }
    
    ctl.JSON(ctx, global.Ret{Status: http.StatusOK, Datas: datas})
}
```

7. **註冊路由**（`app/servers/route.go`）
```go
func (s *Server) LoadRoutes() {
    s.routes = []route{
        {http.MethodPost, "/xxx", s.action.xxx.Create, []gin.HandlerFunc{}},
    }
}

func (s *Server) NewControllers() {
    s.action.xxx = controllers.NewXxxController(s.app)
}
```

8. **更新 actions 結構**（`app/servers/route.go`）
```go
type actions struct {
    user *controllers.UserController
    xxx  *controllers.XxxController  // 新增
}
```

### 新增一個 gRPC 服務端點

1. **定義 Protocol Buffers**（`grpc/proto/xxx.proto`）
```protobuf
syntax = "proto3";

package xxx;
option go_package = "./proto/xxx";

service XxxService {
  rpc Get (GetRequest) returns (GetResponse);
  rpc Create (CreateRequest) returns (CreateResponse);
}

message GetRequest {
  int64 ID = 1;
}

message GetResponse {
  int64 Code = 1;
  string Msg = 2;
  GetResponseDatas Datas = 3;
}

message GetResponseDatas {
  int64 ID = 1;
  string Name = 2;
}

message CreateRequest {
  string Name = 1;
}

message CreateResponse {
  int64 Code = 1;
  string Msg = 2;
  CreateResponseDatas Datas = 3;
}

message CreateResponseDatas {
  int64 ID = 1;
  string Name = 2;
}
```

2. **編譯 Protocol Buffers**
```bash
protoc --go_out=./grpc --go-grpc_out=./grpc grpc/proto/xxx.proto
```

3. **實作 gRPC 服務**（`grpc/services/xxx.go`）
```go
package services

import (
    "akali/app"
    "akali/app/models"
    "akali/app/repositories"
    "akali/grpc/proto/xxx"
    "context"
    "time"
)

type Xxx struct {
    xxx.UnimplementedXxxServiceServer
    BaseService
    App            *app.App
    XxxRepository *repositories.XxxRepository
}

func (s *Xxx) Get(ctx context.Context, req *xxx.GetRequest) (*xxx.GetResponse, error) {
    data, err := do(func() error {
        var err error
        data, err = s.XxxRepository.Get(ctx, models.Xxx{ID: uint(req.GetID())})
        return err
    })
    if err != nil {
        return &xxx.GetResponse{Code: 20001, Msg: err.Error()}, err
    }

    if data.ID == 0 {
        return &xxx.GetResponse{Code: 40001, Msg: "Not found"}, nil
    }

    return &xxx.GetResponse{
        Code: 0,
        Msg:  "OK",
        Datas: &xxx.GetResponseDatas{
            ID:   int64(data.ID),
            Name: data.Name,
        },
    }, nil
}

func (s *Xxx) Create(ctx context.Context, req *xxx.CreateRequest) (*xxx.CreateResponse, error) {
    data, err := do(func() error {
        var err error
        data, err = s.XxxRepository.Create(ctx, models.Xxx{Name: req.GetName()})
        return err
    })
    if err != nil {
        return &xxx.CreateResponse{Code: 20001, Msg: err.Error()}, err
    }

    return &xxx.CreateResponse{
        Code: 0,
        Msg:  "OK",
        Datas: &xxx.CreateResponseDatas{
            ID:   int64(data.ID),
            Name: data.Name,
        },
    }, nil
}
```

4. **註冊 gRPC 服務**（`grpc/server.go`）
```go
func (s *Grpc) registerServices() {
    user.RegisterUserServiceServer(s.srv, &services.User{...})
    xxx.RegisterXxxServiceServer(s.srv, &services.Xxx{  // 新增
        App:            s.app,
        XxxRepository: repositories.NewXxxRepository(s.app),
    })
}
```

### HTTP RESTful vs gRPC 選擇建議

#### 使用 HTTP RESTful 的場景
- 需要對外提供公開 API
- 需要 Swagger 文檔
- 需要瀏覽器直接訪問
- 需要簡單的 JSON 格式
- 需要跨語言、跨平台的通用性

#### 使用 gRPC 的場景
- 微服務內部通訊
- 需要高效能、低延遲
- 需要強型別檢查
- 需要流式處理（Streaming）
- 需要雙向通訊
- 需要更好的錯誤處理機制

#### 共用 Repository 層
- HTTP RESTful 和 gRPC 可以共用相同的 Repository 層
- 業務邏輯可以共用 Service 層（但需要適配不同的輸入輸出格式）
- 資料模型（Model）完全共用

---

## 重要注意事項

### ⚠️ 資料庫操作
- **務必使用 Context**：所有資料庫操作必須傳遞 `context.Context`
- **主從分離**：讀取用 `RDB`，寫入用 `WDB`
- **連線池設定**：已自動設定，無需手動管理

### ⚠️ 錯誤處理
- **不要直接 panic**：使用錯誤回傳機制
- **統一錯誤碼**：在 `global/errInfos/code.go` 中定義
- **錯誤訊息**：在 `global/errInfos/message.go` 中定義

### ⚠️ 並發安全
- **Context 傳遞**：確保 TraceID 正確傳遞
- **資料庫連線**：GORM 已處理並發安全
- **優雅退出**：使用 WaitGroup 等待所有請求完成

### ⚠️ 效能考量
- **資料庫查詢**：避免 N+1 查詢問題
- **快取使用**：適當使用 Redis 快取
- **日誌記錄**：Debug 模式才記錄詳細日誌

### ⚠️ 安全性
- **參數驗證**：所有輸入都必須驗證
- **SQL 注入**：使用 GORM 參數化查詢，避免拼接 SQL
- **敏感資訊**：不要在日誌中記錄密碼等敏感資訊

### ⚠️ gRPC 服務開發
- **Protocol Buffers 編譯**：修改 `.proto` 檔案後必須重新編譯
- **錯誤碼統一**：gRPC 回應中的錯誤碼應與 HTTP API 保持一致
- **Context 傳遞**：確保 Context 正確傳遞到 Repository 層
- **Metadata 處理**：TraceID 從 Metadata 取得，若無則自動生成
- **服務註冊**：新增服務後必須在 `registerServices()` 中註冊
- **Unimplemented 嵌入**：必須嵌入 `UnimplementedXxxServiceServer` 以保持向後兼容
