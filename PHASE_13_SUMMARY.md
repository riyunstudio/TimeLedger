# Stage 13：DTO 與 Resource 映射層實作總結

## 一、概述

本階段成功實作了資料傳輸對象（DTO）模式，將資料庫模型與 API 響應格式進行解耦，達成前後端資料格式的職責分離。透過建立獨立的 Resource 層，系統架構更加清晰，維護性大幅提升。

### 實作目標

1. 建立統一的 API 響應格式轉換層
2. 移除 Service 層中的響應格式邏輯
3. 確保 Controller 層的職責單一化
4. 提升程式碼的可測試性與可維護性

---

## 二、新增 Resource 檔案

### 2.1 RoomResource（教室資源）

**檔案位置**：`app/resources/room.go`

**結構定義**：

```go
type RoomResponse struct {
    ID        uint      `json:"id"`
    CenterID  uint      `json:"center_id"`
    Name      string    `json:"name"`
    Capacity  int       `json:"capacity"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
}
```

**轉換方法**：

- `ToRoomResponse(room models.Room) *RoomResponse` - 單一教室轉換
- `ToRoomResponses(rooms []models.Room) []RoomResponse` - 批量教室轉換

### 2.2 CenterResource（中心資源）

**檔案位置**：`app/resources/center.go`

**結構定義**：

```go
type CenterResponse struct {
    ID        uint           `json:"id"`
    Name      string         `json:"name"`
    PlanLevel string         `json:"plan_level"`
    Settings  CenterSettings `json:"settings"`
    CreatedAt time.Time      `json:"created_at"`
}

type CenterSettings struct {
    AllowPublicRegister bool   `json:"allow_public_register"`
    DefaultLanguage     string `json:"default_language"`
    ExceptionLeadDays   int    `json:"exception_lead_days"`
}
```

**轉換方法**：

- `ToCenterResponse(center models.Center) *CenterResponse` - 中心轉換
- `ToCenterResponses(centers []models.Center) []CenterResponse` - 批量轉換

### 2.3 CourseResource（課程資源）

**檔案位置**：`app/resources/course.go`

**結構定義**：

```go
type CourseResponse struct {
    ID               uint      `json:"id"`
    CenterID         uint      `json:"center_id"`
    Name             string    `json:"name"`
    DefaultDuration  int       `json:"default_duration"`
    ColorHex         string    `json:"color_hex"`
    RoomBufferMin    int       `json:"room_buffer_min"`
    TeacherBufferMin int       `json:"teacher_buffer_min"`
    IsActive         bool      `json:"is_active"`
    CreatedAt        time.Time `json:"created_at"`
}
```

**轉換方法**：

- `ToCourseResponse(course models.Course) *CourseResponse` - 單一課程轉換
- `ToCourseResponses(courses []models.Course) []CourseResponse` - 批量轉換

---

## 三、Service 層重構

### 3.1 變更檔案清單

| 檔案 | 變更內容 |
|:---|:---|
| `app/services/room.go` | 移除 RoomResponse 結構、ToRoomResponse()、ToRoomResponses() |
| `app/services/course.go` | 移除 CourseResponse 結構、ToCourseResponse()、ToCourseResponses() |
| `app/services/center.go` | 無變更（原本無映射方法） |

### 3.2 重構前程式碼

```go
// app/services/room.go（重構前）
type RoomResponse struct {
    ID        uint      `json:"id"`
    CenterID  uint      `json:"center_id"`
    Name      string    `json:"name"`
    Capacity  int       `json:"capacity"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
}

func (s *RoomService) ToRoomResponse(room models.Room) *RoomResponse {
    return &RoomResponse{
        ID:        room.ID,
        CenterID:  room.CenterID,
        Name:      room.Name,
        Capacity:  room.Capacity,
        IsActive:  room.IsActive,
        CreatedAt: room.CreatedAt,
    }
}
```

### 3.3 重構後程式碼

```go
// app/services/room.go（重構後）
func (s *RoomService) CreateRoom(ctx context.Context, req *requests.CreateRoomRequest) (*models.Room, *errInfos.Res, error) {
    // ... 業務邏輯
    return room, nil, nil
}

// 響應格式轉換移至 Resource 層
```

---

## 四、Controller 層整合

### 4.1 新增依賴注入

| 控制器 | 新增依賴 |
|:---|:---|
| `AdminRoomController` | `roomResource *resources.RoomResource` |
| `AdminCenterController` | `centerResource *resources.CenterResource` |
| `AdminCourseController` | `courseResource *resources.CourseResource` |

### 4.2 控制器初始化範例

```go
type AdminRoomController struct {
    app          *app.App
    roomRepo     *repositories.RoomRepository
    roomResource *resources.RoomResource
}

func NewAdminRoomController(app *app.App) *AdminRoomController {
    return &AdminRoomController{
        app:          app,
        roomRepo:     repositories.NewRoomRepository(app),
        roomResource: resources.NewRoomResource(app),
    }
}
```

### 4.3 API 方法更新

| 方法 | 更新內容 |
|:---|:---|
| `GetRooms` | 使用 `roomResource.ToRoomResponses()` 轉換響應 |
| `CreateRoom` | 使用 `roomResource.ToRoomResponse()` 轉換響應 |
| `UpdateRoom` | 使用 `roomResource.ToRoomResponse()` 轉換響應 |
| `GetActiveRooms` | 使用 `roomResource.ToRoomResponses()` 轉換響應 |

---

## 五、架構效益

### 5.1 職責分離

| 層級 | 職責 |
|:---|:---|
| **Controller** | 請求解析與響應格式化 |
| **Service** | 業務邏輯處理 |
| **Resource** | 資料格式轉換 |

### 5.2 具體效益

| 效益類別 | 說明 |
|:---|:---|
| **欄位過濾** | 資料庫模型的所有欄位不直接暴露給前端，可彈性控制回傳欄位（如隱藏敏感資料） |
| **格式標準化** | 所有 API 響應格式由 Resource 層統一處理，時間格式、欄位命名保持一致性 |
| **可維護性** | 未來若有格式變更，只需修改 Resource 層，新增欄位或調整格式不影響 Service 層 |
| **可測試性** | Resource 層可獨立單元測試，確保轉換邏輯正確 |

### 5.3 程式碼品質提升

| 指標 | 改善前 | 改善後 |
|:---|:---:|:---:|
| Service 層職責 | 包含業務邏輯與響應格式 | 僅包含業務邏輯 |
| 重複程式碼 | 多個 Controller 有相似轉換邏輯 | 集中於 Resource 層 |
| 修改影響範圍 | 格式變更需修改多個檔案 | 僅需修改 Resource 層 |

---

## 六、app/resources/ 目錄結構

```
app/resources/
├── base.go              # BaseResource（基礎結構）
├── center.go            # CenterResource（中心資源）← 新增
├── course.go            # CourseResource（課程資源）← 新增
├── invitation.go        # InvitationResource（邀請資源）
├── room.go              # RoomResource（教室資源）← 新增
├── scheduling.go        # ScheduleResource（排課資源）
├── session_note.go      # SessionNoteResource（課程筆記資源）
├── teacher.go           # TeacherProfileResource（教師資源）
└── user.go              # UserResource（用戶資源）
```

---

## 七、編譯驗證

```bash
$ go build -mod=mod ./...
# 輸出：無錯誤，全部編譯成功
```

---

## 八、後續建議

### 8.1 高優先順序

| 項目 | 說明 |
|:---|:---|
| 遷移剩餘控制器 | 將 TeacherProfileController、ScheduleController 等也遷移至 Resource 模式 |
| 單元測試 | 為 Resource 層撰寫單元測試，確保轉換邏輯正確 |

### 8.2 中優先順序

| 項目 | 說明 |
|:---|:---|
| 錯誤響應 Resource | 建立通用錯誤響應 Resource（ErrorResponse） |
| 分頁響應 Resource | 建立標準化分頁響應格式 |

### 8.3 架構擴展方向

| 方向 | 說明 |
|:---|:---|
| API 版本控制 | 支援不同版本的 API 響應格式 |
| 欄位篩選 | 支援前端指定回傳欄位（field filtering） |
| 響應快取 | Resource 層可加入快取機制，提升效能 |

---

## 九、總結

Stage 13 成功實作了 DTO 與 Resource 映射層，達成以下目標：

1. ✅ 建立三個 Resource 檔案（RoomResource、CenterResource、CourseResource）
2. ✅ 完成 Service 層重構，移除響應格式邏輯
3. ✅ 完成 Controller 層整合，注入 Resource 依賴
4. ✅ 編譯驗證通過，系統正常運作

透過本次重構，系統架構更加清晰，各層職責更加明確，為後續的功能擴展與維護奠定良好基礎。

---

**完成日期**：2026年1月31日

**Commit**：`[待填入]`
