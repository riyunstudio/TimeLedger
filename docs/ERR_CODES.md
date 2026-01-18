# ERR_CODES.md — TimeLedger 錯誤碼規格
版本：v1.0
日期：2026-01-18

> 原則：所有 API 回傳 envelope 固定：
```json
{"code":0,"message":"OK","datas":{}}
```
> code=0 表成功；非 0 為失敗。

---

## 1. 通用錯誤（Global）
| code | message | HTTP | 說明 |
|---:|---|---:|---|
| 30001 | PARAMS_VALIDATE_ERROR | 400 | Request 參數驗證失敗 |
| 30002 | UNAUTHORIZED | 401 | 未授權 / Token 無效 |
| 30003 | FORBIDDEN | 403 | 權限不足（含跨中心存取） |
| 30004 | NOT_FOUND | 404 | 資源不存在 |
| 30005 | RATE_LIMITED | 429 | 請求過於頻繁（選配） |
| 50001 | SQL_ERROR | 500 | DB 錯誤 |
| 50002 | SYSTEM_ERROR | 500 | 系統例外 / 未預期錯誤 |

---

## 2. 業務錯誤（Business）
| code | message | HTTP | 說明 |
|---:|---|---:|---|
| 31001 | CENTER_SCOPE_VIOLATION | 403 | 跨中心資料存取（必須阻擋） |
| 31002 | TRIAL_EXPIRED_READONLY | 403 | 試用到期/方案到期，進入唯讀（禁止寫入） |
| 31003 | PLAN_LIMIT_EXCEEDED | 403 | 方案限制超出（MVP：老師數） |
| 31101 | TIME_CONFLICT | 409 | 排課/改期時段衝突（validate 會回衝突明細） |
| 31102 | EXCEPTION_ALREADY_EXISTS | 409 | 同 rule 在同一日期已存在有效例外（PENDING/APPROVED） |
| 31103 | EXCEPTION_NOT_PENDING | 409 | 僅能審核 PENDING 例外 |
| 31104 | MEMBERSHIP_NOT_ACTIVE | 403 | 老師尚未啟用/未加入中心，中心不可排課/引用 |
| 31105 | INVALID_DATE_RANGE | 400 | 查詢區間不合法（from>to 或超過限制） |

---

## 3. validate 衝突明細 datas 格式（TIME_CONFLICT）
HTTP 409 + code=31101，datas 為陣列：
```json
[
  {
    "type":"OVERLAP",
    "source":"PERSONAL|CENTER_RULE|EXCEPTION",
    "start_at":"2026-01-20 10:00:00",
    "end_at":"2026-01-20 12:00:00",
    "title":"私人教學",
    "meta":{"session_id":123,"center_id":1,"rule_id":55}
  }
]
```

欄位說明：
- type：OVERLAP（MVP 只做 overlap；未來可擴 room/teacher buffer）
- source：衝突來源
- start_at/end_at/title：讓前端能「說人話」
- meta：可選，帶關聯 id 方便點擊跳轉
