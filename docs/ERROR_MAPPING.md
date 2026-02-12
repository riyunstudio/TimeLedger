# TimeLedger 錯誤碼對照表

## 1. 錯誤碼體系

錯誤碼格式：`AppID(1位) + 類型(2位) + 流水號(4位)`

- **AppID**: 專案流水號（預設 1）
- **類型**: 01-12 代表不同錯誤類型
- **流水號**: 0001-9999

---

## 2. 錯誤碼對照

### 2.1 系統相關 (1xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 10001 | 500 | System error. | 系統錯誤 |
| 10002 | 400 | Invalid parameters | 參數格式錯誤 |
| 10003 | 500 | Json encode error. | JSON 編碼錯誤 |
| 10004 | 400 | Json decode error. | JSON 解碼錯誤 |
| 10005 | 500 | Json process error. | JSON 處理錯誤 |
| 10006 | 500 | Format resource error. | 資源格式錯誤 |
| 10007 | 408 | Request timeout. | 請求超時 |
| 10008 | 501 | Feature not implemented | 功能尚未實作 |
| 10009 | 429 | Rate limit exceeded | 請求頻率過高，請稍後再試 |

### 2.2 資料庫/快取相關 (2xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 20001 | 500 | Database operation failed | 資料庫操作失敗 |
| 20002 | 500 | SQL transaction error. | 交易錯誤 |

### 2.3 權限與認證 (3xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 30001 | 401 | Please login first | 請先登入 |
| 30002 | 403 | Permission denied | 權限不足 |
| 30003 | 401 | Token expired | Token 已過期 |
| 30004 | 401 | Invalid token | 無效的 Token |
| 30005 | 400 | Invalid or expired invite | 邀請碼無效或已過期 |

### 2.4 業務資源類 (4xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 40001 | 404 | Resource not found | 找不到資源 |
| 40002 | 409 | Resource already exists | 資源已存在 |
| 40003 | 400 | Invalid hashtag format | 標籤格式錯誤 |
| 40004 | 403 | Plan limit reached | 超過方案配額 |
| 40005 | 409 | Resource is in use | 資源仍在使用中 |
| 40006 | 409 | Course has active offerings | 課程模板仍有關聯班別 |
| 40007 | 409 | Offering has schedule rules | 班別仍有排課規則 |
| 40008 | 409 | Room has active schedules | 教室仍有排課安排 |
| 40009 | 400 | Invalid status transition | 不允許的狀態轉換 |
| 40010 | 400 | Teacher not registered | 老師尚未註冊（需要先完成註冊流程） |

### 2.5 排課核心類 (5xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 50001 | 409 | Time slot occupied | 時段被佔用 |
| 50002 | 409 | Insufficient buffer time | 緩衝時間不足 |
| 50003 | 400 | Cannot book past time | 不能排過去的時間 |
| 50004 | 409 | Slot is locked by another | 時段已被鎖定 |
| 50005 | 400 | Center is closed | 非營業時間 |
| 50006 | 400 | Invalid date range | 日期範圍錯誤 |
| 50007 | 409 | Rule conflict detected | 規則衝突 |
| 50008 | 409 | Exception already exists | 該日期已有例外單 |
| 50009 | 400 | Teacher required | 必須指定老師 |
| 50010 | 400 | Room required | 必須指定教室 |
| 50011 | 404 | Offering not found | 班別不存在 |
| 50012 | 404 | Course not found | 課程模板不存在 |
| 50013 | 400 | Invalid weekday | 無效的星期幾 |
| 50014 | 400 | Invalid duration | 無效的課程時長 |
| 50015 | 400 | Start after end | 開始時間晚於結束時間 |
| 50016 | 400 | Invalid date format | 無效的日期格式 |
| 50017 | 400 | End before start | 結束日期早於開始日期 |
| 50018 | 400 | Duration exceeds limit | 課程時長超過限制 |

### 2.6 例外與審核類 (6xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 60001 | 404 | Exception request not found | 例外申請不存在 |
| 60002 | 400 | Invalid action for current status | 當前狀態不允許此操作 |
| 60003 | 400 | Exception already reviewed | 例外已審核過 |
| 60004 | 400 | Exception was revoked | 例外已撤回 |
| 60005 | 403 | Cannot reject own request | 不能拒絕自己提交的申請 |
| 60006 | 400 | Deadline exceeded | 已超過異動截止日（需提前 14 天申請） |
| 60007 | 403 | Self review forbidden | 不能審核自己提交的申請 |
| 60008 | 400 | Exception already processed | 例外已被處理 |
| 60009 | 409 | Reschedule time conflicts | 調課時間與現有排程衝突 |
| 60010 | 400 | Invalid substitute teacher | 代課老師無效 |
| 60011 | 400 | Cancellation deadline passed | 停課截止日已過 |
| 60012 | 400 | Reschedule requires new time | 調課必須提供新時間 |
| 60013 | 400 | Invalid edit mode | 無效的編輯模式 |
| 60014 | 400 | No affected sessions | 沒有受影響的場次 |
| 60015 | 400 | Future mode requires edit date | FUTURE 模式必須指定編輯日期 |
| 60016 | 400 | Edit date required | 編輯日期為必填 |
| 60017 | 400 | Delete confirmation required | 刪除操作需要確認 |
| 60018 | 400 | Batch limit exceeded | 批量操作超過限制 |

### 2.7 檔案與媒體類 (7xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 70001 | 413 | File size exceeds limit | 檔案超過大小限制 |
| 70002 | 400 | Invalid file type | 不支援的檔案類型 |
| 70003 | 500 | Upload failed | 上傳失敗 |
| 70004 | 404 | Certificate not found | 證照不存在 |

### 2.8 搜尋與媒合類 (8xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 80001 | 400 | Talent search not available | 該老師未開放搜尋 |

### 2.9 LINE Bot 與通知類 (9xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 90001 | 400 | LINE account already bound | LINE 帳號已綁定 |
| 90002 | 400 | LINE account not bound | LINE 帳號未綁定 |
| 90003 | 400 | Invalid binding code | 驗證碼無效 |
| 90004 | 400 | Binding code expired | 驗證碼已過期，請重新產生 |
| 90005 | 500 | Failed to send LINE notification | LINE 通知發送失敗 |

### 2.10 管理員類 (10xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 100001 | 404 | Admin user not found | 管理員不存在 |
| 100002 | 400 | Email already registered | Email 已被註冊 |
| 100003 | 400 | Current password is incorrect | 舊密碼錯誤 |
| 100004 | 400 | Cannot disable yourself | 不能停用自己的帳號 |

### 2.11 資源鎖定與衝突類 (11xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 110001 | 409 | Resource is locked | 資源正在被其他操作修改，請稍後再試 |
| 110002 | 409 | Resource was modified | 資源已被其他請求修改，請重新整理後再試 |
| 110003 | 500 | Transaction failed | 交易執行失敗，請稍後再試 |

### 2.12 交易執行類 (12xxxx)

| 錯誤碼 | HTTP Status | 英文訊息 | 中文訊息 |
|:---:|:---:|:---|:---|
| 120001 | 504 | Transaction timeout | 交易執行超時 |
| 120002 | 500 | Deadlock detected | 偵測到資料庫死鎖 |
| 120003 | 500 | Partial completion | 交易部分完成（部分操作成功） |
| 120004 | 500 | Rollback failed | 交易回滾失敗 |
| 120005 | 400 | Constraint violation | 違反資料庫約束 |
| 120006 | 400 | Foreign key violation | 外鍵約束違反 |
| 120007 | 400 | Unique violation | 唯一約束違反 |
| 120008 | 400 | Check violation | CHECK 約束違反 |
| 120009 | 409 | Serialization failure | 序列化失敗（並發衝突） |
| 120010 | 504 | Lock wait timeout | 鎖等待超時 |
| 120011 | 500 | Share lock failed | 共享鎖獲取失敗 |
| 120012 | 500 | Exclusive lock failed | 排他鎖獲取失敗 |

---

## 3. HTTP Status 對照

| HTTP Status | 說明 | 常用錯誤碼 |
|:---:|:---|:---|
| 200 | OK | 成功 |
| 201 | Created | 建立成功 |
| 204 | No Content | 刪除成功 |
| 400 | Bad Request | 參數錯誤、業務規則錯誤 |
| 401 | Unauthorized | 未登入、Token 過期 |
| 403 | Forbidden | 權限不足 |
| 404 | Not Found | 資源不存在 |
| 408 | Request Timeout | 請求超時 |
| 409 | Conflict | 資料衝突（重疊、鎖定等） |
| 413 | Payload Too Large | 檔案過大 |
| 429 | Too Many Requests | 頻率限制 |
| 500 | Internal Server Error | 系統錯誤 |
| 501 | Not Implemented | 功能未實作 |
| 504 | Gateway Timeout | 逾時 |

---

## 4. 錯誤處理流程

### 4.1 錯誤回應格式

```json
{
  "code": 40001,
  "message": "找不到資源",
  "data": null
}
```

### 4.2 服務層錯誤產生

```go
// 方式 1: 使用錯誤碼
return nil, app.Err.New(errInfos.NOT_FOUND), err

// 方式 2: 使用錯誤碼 + 自訂訊息
return nil, app.Err.New(errInfos.SCHED_OVERLAP), fmt.Errorf("teacher %d has conflict", teacherID)
```

### 4.3 控制器錯誤處理

```go
// 使用 ContextHelper
helper.ErrorWithInfo(errInfo)

// 或手動指定
helper.ErrorWithCode(http.StatusConflict, errInfos.SCHED_OVERLAP, "時段被佔用")
```

---

*本文件基於實際程式碼生成，最後更新：2026-02-12*
