# TimeLedger 模組合約

## 1. 模組架構模式

### 1.1 分層架構

每個模組遵循標準的分層架構：

```
Controller → Service → Repository → Model
```

### 1.2 通用模式

#### Triple Return Pattern
```go
func (s *Service) Method(ctx context.Context, req *Request) (data any, *errInfos.Res, error) {
    // 1. 查詢/驗證
    // 2. 業務邏輯
    // 3. 回傳
    return result, nil, nil
}
```

#### ContextHelper Pattern
```go
func (c *Controller) Handler(ctx *gin.Context) {
    helper := controllers.NewContextHelper(ctx)
    
    // 參數解析
    id := helper.MustParamUint("id")
    if id == 0 { return }
    
    var req Request
    if !helper.MustBindJSON(&req) { return }
    
    // 呼叫 Service
    result, errInfo, err := c.service.Method(ctx, id, &req)
    if err != nil {
        helper.ErrorWithInfo(errInfo)
        return
    }
    
    helper.Success(result)
}
```

---

## 2. 教師模組合約

### 2.1 TeacherProfileService

| 方法 | 輸入 | 輸出 | 錯誤碼 |
|:---|:---|:---|:---|
| `GetProfile` | teacherID uint | TeacherProfile | SQL_ERROR |
| `UpdateProfile` | teacherID uint, req *UpdateRequest | TeacherProfile | SQL_ERROR |
| `GetCenters` | teacherID uint | []Center | SQL_ERROR |
| `GetSkills` | teacherID uint | []TeacherSkill | SQL_ERROR |
| `CreateSkill` | teacherID uint, req *SkillRequest | TeacherSkill | SQL_ERROR |
| `UpdateSkill` | skillID, teacherID uint, req *SkillRequest | TeacherSkill | NOT_FOUND, FORBIDDEN |
| `DeleteSkill` | skillID, teacherID uint | - | NOT_FOUND, FORBIDDEN |
| `GetCertificates` | teacherID uint | []TeacherCertificate | SQL_ERROR |
| `CreateCertificate` | teacherID uint, req *CertificateRequest | TeacherCertificate | SQL_ERROR |
| `DeleteCertificate` | certID, teacherID uint | - | NOT_FOUND, FORBIDDEN |

### 2.2 TeacherProfileController

| 方法 | 端點 | 認證 | 權限 |
|:---|:---|:---|:---|
| `GetProfile` | GET /teacher/me/profile | ✅ | Teacher |
| `UpdateProfile` | PUT /teacher/me/profile | ✅ | Teacher |
| `GetCenters` | GET /teacher/me/centers | ✅ | Teacher |
| `GetSkills` | GET /teacher/me/skills | ✅ | Teacher |
| `CreateSkill` | POST /teacher/me/skills | ✅ | Teacher |
| `UpdateSkill` | PUT /teacher/me/skills/:id | ✅ | Teacher |
| `DeleteSkill` | DELETE /teacher/me/skills/:id | ✅ | Teacher |
| `GetCertificates` | GET /teacher/me/certificates | ✅ | Teacher |
| `CreateCertificate` | POST /teacher/me/certificates | ✅ | Teacher |
| `DeleteCertificate` | DELETE /teacher/me/certificates/:id | ✅ | Teacher |

---

## 3. 排課模組合約

### 3.1 ScheduleService 介面

| 方法 | 輸入 | 輸出 | 錯誤碼 |
|:---|:---|:---|:---|
| `ValidateOverlap` | centerID, teacherID, roomID, startAt, endAt | ValidationResult | SQL_ERROR |
| `ValidateBuffer` | centerID, teacherID, roomID, startAt | ValidationResult | SQL_ERROR |
| `CreateRule` | centerID uint, req *CreateRuleRequest | ScheduleRule | SQL_ERROR, SCHED_OVERLAP, SCHED_BUFFER |
| `UpdateRule` | ruleID, centerID uint, req *UpdateRuleRequest | ScheduleRule | NOT_FOUND, FORBIDDEN, SCHED_OVERLAP |
| `DeleteRule` | ruleID, centerID uint | - | NOT_FOUND, FORBIDDEN |
| `ExpandRules` | centerID uint, from, to time.Time | []ExpandedSession | SQL_ERROR |
| `GetExceptions` | centerID uint, filters ExceptionFilters | []ScheduleException | SQL_ERROR |
| `CreateException` | centerID, teacherID uint, req *ExceptionRequest | ScheduleException | SQL_ERROR, EXCEPTION_DEADLINE_EXCEEDED |
| `ReviewException` | exceptionID, adminID uint, action ReviewAction | ScheduleException | EXCEPTION_NOT_FOUND, EXCEPTION_INVALID_ACTION |

### 3.2 ValidationResult 結構

```go
type ValidationResult struct {
    Valid     bool                   `json:"valid"`
    Conflicts []ValidationConflict   `json:"conflicts,omitempty"`
}

type ValidationConflict struct {
    Type              string    `json:"type"` // TEACHER_OVERLAP, ROOM_OVERLAP, TEACHER_BUFFER, ROOM_BUFFER
    Message           string    `json:"message"`
    CurrentGapMinutes int       `json:"current_gap_minutes"`
    RequiredMinutes   int       `json:"required_buffer_minutes"`
    CanOverride       bool      `json:"can_override"`
}
```

### 3.3 Exception State Machine

| 狀態 | 可轉換至 | 觸發者 | 動作 |
|:---|:---|:---|:---|
| PENDING | REVOKED | Teacher | 撤回例外 |
| PENDING | APPROVED | Admin | 核准例外 → 執行重驗證 → 更新 Schedule |
| PENDING | REJECTED | Admin | 拒絕例外 → 通知 Teacher |
| APPROVED | CANCELLED | Admin | 取消核准 → 恢復 Schedule |

---

## 4. 智慧媒合模組合約

### 4.1 SmartMatchingService

| 方法 | 輸入 | 輸出 | 錯誤碼 |
|:---|:---|:---|:---|
| `FindMatches` | centerID uint, req *MatchingRequest | []MatchResult | SQL_ERROR |
| `SearchTalent` | centerID uint, filters TalentFilters | []Teacher | SQL_ERROR |
| `GetTalentStats` | centerID uint | TalentStats | SQL_ERROR |
| `InviteTalent` | centerID, adminID uint, req *InviteRequest | InviteResult | SQL_ERROR, TALENT_NOT_OPEN |
| `GetSuggestions` | centerID uint, keyword string | []SearchSuggestion | SQL_ERROR |
| `GetAlternativeSlots` | centerID uint, req *AlternativeRequest | []AlternativeSlot | SQL_ERROR |

### 4.2 MatchResult 結構

```go
type MatchResult struct {
    TeacherID     uint     `json:"teacher_id"`
    TeacherName   string   `json:"teacher_name"`
    Score         float64  `json:"score"`
    Availability  float64  `json:"availability"`
    Rating        float64  `json:"rating"`
    SkillsMatch   float64  `json:"skills_match"`
    RegionMatch   bool     `json:"region_match"`
    IsOpen        bool     `json:"is_open"`
}
```

### 4.3 評分因子

| 因子 | 權重 | 評分邏輯 |
|:---|:---:|:---|
| **Availability** | 40% | 完全空閒 +40分，Buffer 衝突 +15分，Hard Overlap 0分 |
| **Internal Evaluation** | 40% | 星等評分正規化 0~30分，內部備註關鍵字額外 +10分 |
| **Skill & Region Match** | 20% | 技能命中 +10分，標籤命中 +8分，地區命中 +10分 |

---

## 5. 通知模組合約

### 5.1 NotificationService

| 方法 | 輸入 | 輸出 | 錯誤碼 |
|:---|:---|:---|:---|
| `CreateNotification` | userID, userType uint, req *NotificationRequest | Notification | SQL_ERROR |
| `MarkAsRead` | notificationID, userID uint | - | NOT_FOUND |
| `MarkAllAsRead` | userID uint | - | SQL_ERROR |
| `GetUnreadCount` | userID uint | int | SQL_ERROR |
| `QueueLINENotification` | userID uint, msg *LINEMessage | - | SQL_ERROR |

### 5.2 NotificationQueueService

| 方法 | 說明 |
|:---|:---|
| `ProcessQueue` | 處理 Redis 佇列中的通知 |
| `GetQueueStats` | 取得佇列統計 |
| `Enqueue` | 加入通知佇列 |

### 5.3 QueueStats 結構

```go
type QueueStats struct {
    PendingCount  int64 `json:"pending_count"`
    RetryCount    int64 `json:"retry_count"`
    CompletedCount int64 `json:"completed_count"`
    FailedCount   int64 `json:"failed_count"`
    FailureRate   float64 `json:"failure_rate"`
    RedisConnected bool   `json:"redis_connected"`
    WorkerRunning bool   `json:"worker_running"`
}
```

---

## 6. 邀請模組合約

### 6.1 InvitationService

| 方法 | 輸入 | 輸出 | 錯誤碼 |
|:---|:---|:---|:---|
| `InviteTeacher` | centerID, adminID uint, teacherEmail string | CenterInvitation | SQL_ERROR, DUPLICATE |
| `GenerateLink` | centerID, adminID uint, req *LinkRequest | InvitationLink | SQL_ERROR |
| `AcceptInvitation` | token string, teacherID uint | CenterMembership | INVALID_INVITE, DUPLICATE |
| `RevokeInvitation` | invitationID, adminID uint | - | NOT_FOUND, FORBIDDEN |

---

## 7. Repository 通用介面

### 7.1 GenericRepository[T]

| 方法 | 說明 |
|:---|:---|
| `GetByID` | 依 ID 查詢 |
| `GetByIDWithCenterScope` | 依 ID + center_id 查詢 |
| `First` | 第一筆記錄 |
| `FirstWithCenterScope` | 第一筆記錄（含 center_id 過濾） |
| `Find` | 查詢多筆記錄 |
| `FindWithCenterScope` | 查詢多筆記錄（含 center_id 過濾） |
| `FindPaged` | 分頁查詢 |
| `Create` | 建立單筆記錄 |
| `CreateBatch` | 批次建立 |
| `Update` | 更新記錄 |
| `UpdateFields` | 更新特定欄位 |
| `UpdateFieldsWithCenterScope` | 更新特定欄位（含 center_id） |
| `DeleteByID` | 軟刪除 |
| `DeleteByIDWithCenterScope` | 軟刪除（含 center_id） |
| `Exists` | 檢查是否存在 |
| `Count` | 計數 |
| `Transaction` | 交易執行 |

---

## 8. Service 通用基礎

### 8.1 BaseService

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `App` | *app.App | 應用實例 |
| `Logger` | *ServiceLogger | 服務日誌 |

### 8.2 PaginationParams

| 欄位 | 類型 | 說明 |
|:---|:---|:---|
| `Page` | int | 頁碼 |
| `Limit` | int | 每頁筆數 |
| `SortBy` | string | 排序欄位 |
| `SortOrder` | string | 排序方向 |

### 8.3 FilterBuilder

| 方法 | 說明 |
|:---|:---|
| `AddEq` | 等於條件 |
| `AddNe` | 不等於條件 |
| `AddGt` | 大於條件 |
| `AddGte` | 大於等於條件 |
| `AddLt` | 小於條件 |
| `AddLte` | 小於等於條件 |
| `AddLike` | LIKE 條件 |
| `AddIn` | IN 條件 |
| `AddBetween` | 區間條件 |
| `AddCenterScope` | 中心範圍條件 |

---

## 9. 常數定義

### 9.1 排課規則狀態

| 常數 | 值 | 說明 |
|:---|:---|:---|
| `RuleStatusPlanned` | PLANNED | 預計/未開成 |
| `RuleStatusConfirmed` | CONFIRMED | 已開成/正式課 |
| `RuleStatusSuspended` | SUSPENDED | 停課/暫停 |
| `RuleStatusArchived` | ARCHIVED | 已歸檔 |

### 9.2 例外類型

| 常數 | 值 | 說明 |
|:---|:---|:---|
| `ExceptionTypeLeave` | LEAVE | 請假 |
| `ExceptionTypeReschedule` | RESCHEDULE | 調課 |
| `ExceptionTypeSwap` | SWAP | 換課 |
| `ExceptionTypeCancel` | CANCEL | 停課 |

### 9.3 例外狀態

| 常數 | 值 | 說明 |
|:---|:---|:---|
| `ExceptionStatusPending` | PENDING | 待審核 |
| `ExceptionStatusApproved` | APPROVED | 已核准 |
| `ExceptionStatusRejected` | REJECTED | 已拒絕 |
| `ExceptionStatusRevoked` | REVOKED | 已撤回 |
| `ExceptionStatusCancelled` | CANCELLED | 已取消 |

### 9.4 管理員角色

| 常數 | 值 | 說明 |
|:---|:---|:---|
| `RoleOwner` | OWNER | 所有者 |
| `RoleAdmin` | ADMIN | 管理員 |
| `RoleStaff` | STAFF | 員工 |

---

*本文件基於實際程式碼生成，最後更新：2026-02-12*
