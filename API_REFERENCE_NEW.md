# API 端點完整對照表

**最後更新**：2026年1月31日
**來源**：swagger.json

---

## Admin

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/admin/centers` | 取得中心列表 |
| POST | `/api/v1/admin/centers` | 新增中心 |
| GET | `/api/v1/admin/centers/{id}/holidays` | 取得假日列表 |
| POST | `/api/v1/admin/centers/{id}/holidays` | 新增假日 |
| POST | `/api/v1/admin/centers/{id}/holidays/bulk` | 批次建立假日 |
| DELETE | `/api/v1/admin/centers/{id}/holidays/{holiday_id}` | 刪除假日 |
| POST | `/api/v1/admin/centers/{id}/invitations` | 邀請老師加入中心 |
| POST | `/api/v1/admin/centers/{id}/offerings/{offering_id}/copy` | 複製班別 |
| GET | `/api/v1/admin/courses` | 取得課程列表 |
| POST | `/api/v1/admin/courses` | 新增課程 |
| GET | `/api/v1/admin/courses/active` | 取得已啟用的課程列表 |
| PUT | `/api/v1/admin/courses/{course_id}` | 更新課程 |
| DELETE | `/api/v1/admin/courses/{course_id}` | 刪除課程 |
| PATCH | `/api/v1/admin/courses/{course_id}/toggle-active` | 切換課程啟用狀態 |
| GET | `/api/v1/admin/offerings` | 取得班別列表 |
| POST | `/api/v1/admin/offerings` | 新增班別 |
| GET | `/api/v1/admin/offerings/active` | 取得啟用的班別列表 |
| PUT | `/api/v1/admin/offerings/{offering_id}` | 更新班別 |
| DELETE | `/api/v1/admin/offerings/{offering_id}` | 刪除班別 |
| PATCH | `/api/v1/admin/offerings/{offering_id}/toggle-active` | 切換班別啟用狀態 |
| GET | `/api/v1/admin/rooms` | 取得教室列表 |
| POST | `/api/v1/admin/rooms` | 新增教室 |
| GET | `/api/v1/admin/rooms/active` | 取得已啟用的教室列表 |
| PUT | `/api/v1/admin/rooms/{room_id}` | 更新教室 |
| PATCH | `/api/v1/admin/rooms/{room_id}/toggle-active` | 切換教室啟用狀態 |
| GET | `/api/v1/admin/teachers` | 取得中心的老師列表 |
| GET | `/api/v1/admin/teachers/{teacher_id}/note` | 取得老師評分與備註 |
| PUT | `/api/v1/admin/teachers/{teacher_id}/note` | 新增或更新老師評分與備註 |
| DELETE | `/api/v1/admin/teachers/{teacher_id}/note` | 刪除老師評分與備註 |
| GET | `/api/v1/admin/templates` | 取得課表模板列表 |
| POST | `/api/v1/admin/templates` | 新增課表模板 |
| PUT | `/api/v1/admin/templates/{template_id}` | 更新課表模板 |
| DELETE | `/api/v1/admin/templates/{template_id}` | 刪除課表模板 |
| POST | `/api/v1/admin/templates/{template_id}/apply` | 套用課表模板 |
| GET | `/api/v1/admin/templates/{template_id}/cells` | 取得模板中的格子 |
| POST | `/api/v1/admin/templates/{template_id}/cells` | 新增模板中的格子 |
| PUT | `/api/v1/admin/templates/{template_id}/cells/reorder` | 重新排序模板中的格子 |
| DELETE | `/api/v1/admin/templates/{template_id}/cells/{cell_id}` | 刪除格子 |
| POST | `/api/v1/admin/templates/{template_id}/validate-apply` | 驗證套用模板 |
| GET | `/api/v1/teachers` | 取得老師列表 |
| DELETE | `/api/v1/teachers/{id}` | 刪除老師 |

## Admin - Dashboard

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/admin/dashboard/today-summary` | 取得管理員後台首頁的今日課表摘要 |

## Admin - Invitations

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/admin/centers/{id}/invitations` | 取得邀請列表 |
| POST | `/api/v1/admin/centers/{id}/invitations/generate-link` | 產生邀請連結 |
| GET | `/api/v1/admin/centers/{id}/invitations/links` | 取得邀請連結列表 |
| GET | `/api/v1/admin/centers/{id}/invitations/stats` | 取得邀請統計資料 |
| DELETE | `/api/v1/admin/invitations/links/{id}` | 撤回邀請連結 |

## Admin - LINE

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/admin/me/line-binding` | 取得 LINE 綁定狀態 |
| POST | `/admin/me/line/bind` | 初始化 LINE 綁定 |
| PATCH | `/admin/me/line/notify-settings` | 更新 LINE 通知設定 |
| DELETE | `/admin/me/line/unbind` | 解除 LINE 綁定 |

## Admin - Management

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/admin/admins` | 取得管理員列表 |
| POST | `/admin/admins` | 建立管理員 |
| POST | `/admin/admins/change-role` | 修改管理員角色 |
| POST | `/admin/admins/reset-password` | 重設管理員密碼 |
| POST | `/admin/admins/toggle-status` | 停用/啟用管理員 |

## Admin - Notifications

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/admin/notifications/queue-stats` | 取得通知佇列統計 |

## Admin - Profile

| Method | Endpoint | Summary |
|:---|:---|:---|
| POST | `/admin/me/change-password` | 修改管理員密碼 |
| GET | `/admin/me/profile` | 取得管理員個人資料 |

## Admin - Scheduling

| Method | Endpoint | Summary |
|:---|:---|:---|
| POST | `/api/v1/admin/scheduling/check-overlap` | 檢查課程時間是否與現有排程衝突 |
| POST | `/api/v1/admin/scheduling/check-room-buffer` | 檢查教室的緩衝時間是否足夠 |
| POST | `/api/v1/admin/scheduling/check-teacher-buffer` | 檢查老師的緩衝時間是否足夠 |
| GET | `/api/v1/admin/scheduling/exceptions` | 取得指定日期範圍內的所有例外申請 |
| GET | `/api/v1/admin/scheduling/exceptions/all` | 取得所有例外申請（可依狀態篩選） |
| GET | `/api/v1/admin/scheduling/exceptions/pending` | 取得所有待審核的例外申請 |
| POST | `/api/v1/admin/scheduling/exceptions/{id}/review` | 審核例外申請（核准/拒絕） |
| POST | `/api/v1/admin/scheduling/expand` | 展開排課規則為具體課程場次 |
| POST | `/api/v1/admin/scheduling/phase-transitions` | 偵測課程序列中的階段轉換點 |
| GET | `/api/v1/admin/scheduling/rules` | 取得中心的所有排課規則 |
| POST | `/api/v1/admin/scheduling/rules` | 建立新的排課規則 |
| POST | `/api/v1/admin/scheduling/rules/check-lock` | 檢查規則是否已超過異動截止日 |
| PUT | `/api/v1/admin/scheduling/rules/{id}` | 更新排課規則 |
| DELETE | `/api/v1/admin/scheduling/rules/{id}` | 刪除排課規則 |
| GET | `/api/v1/admin/scheduling/rules/{id}/exceptions` | 取得指定規則的所有例外申請 |
| POST | `/api/v1/admin/scheduling/validate` | 完整驗證排課（硬衝突 + 緩衝檢查） |

## Auth

| Method | Endpoint | Summary |
|:---|:---|:---|
| POST | `/api/v1/auth/admin/login` | 管理員登入 |
| POST | `/api/v1/auth/logout` | 登出 |
| POST | `/api/v1/auth/refresh` | 刷新 Token |
| POST | `/api/v1/auth/teacher/line/login` | 老師 LINE 登入 |

## Export

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/calendar/subscribe/{token}.ics` | 透過 Token 訂閱課表 |

## Geo

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/geo/cities` | 取得城市列表 |

## Hashtag

| Method | Endpoint | Summary |
|:---|:---|:---|
| POST | `/api/v1/hashtags` | 建立新標籤 |
| GET | `/api/v1/hashtags/search` | 搜尋標籤 |

## Invitations

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/invitations/{token}` | 取得公開邀請資訊 |
| POST | `/api/v1/invitations/{token}/accept` |  |

## LINE

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/admin/me/line/qrcode` | 產生 LINE 官方帳號 QR Code |
| GET | `/api/v1/admin/me/line/qrcode-with-code` | 產生包含驗證碼的 QR Code |

## Notifications

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/notifications` | 取得通知列表 |
| POST | `/notifications/read-all` | 標記所有通知為已讀 |
| POST | `/notifications/test` | 發送測試通知 |
| POST | `/notifications/token` | 設定通知 Token |
| GET | `/notifications/unread-count` | 取得未讀通知數量 |
| POST | `/notifications/{id}/read` | 標記通知為已讀 |

## Smart Matching

| Method | Endpoint | Summary |
|:---|:---|:---|
| POST | `/admin/smart-matching/alternatives` | 取得替代時段建議 |
| POST | `/admin/smart-matching/matches` | 智慧媒合搜尋 |
| GET | `/admin/smart-matching/suggestions` | 取得搜尋建議 |
| POST | `/admin/smart-matching/talent/invite` | 邀請人才合作 |
| GET | `/admin/smart-matching/talent/search` | 人才庫搜尋 |
| GET | `/admin/smart-matching/talent/stats` | 取得人才庫統計資料 |
| GET | `/admin/smart-matching/teachers/{teacher_id}/sessions` | 取得教師課表 |

## Teacher

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/teacher/exceptions` | 老師查看自己的例外申請列表 |
| POST | `/api/v1/teacher/exceptions` | 老師提出停課/改期申請 |
| POST | `/api/v1/teacher/exceptions/{id}/revoke` | 老師撤回待審核的例外申請 |
| GET | `/api/v1/teacher/me/centers` | 取得老師已加入的中心列表 |
| GET | `/api/v1/teacher/me/centers/{center_id}/schedule-rules` | 獲取老師在指定中心的排課規則 |
| GET | `/api/v1/teacher/me/certificates` | 取得老師證照列表 |
| POST | `/api/v1/teacher/me/certificates` | 新增老師證照 |
| POST | `/api/v1/teacher/me/certificates/upload` | 上傳證照檔案 |
| DELETE | `/api/v1/teacher/me/certificates/{id}` | 刪除老師證照 |
| GET | `/api/v1/teacher/me/profile` | 取得老師個人資料 |
| PUT | `/api/v1/teacher/me/profile` | 更新老師個人資料 |
| GET | `/api/v1/teacher/me/schedule` | 取得老師的綜合課表 |
| GET | `/api/v1/teacher/me/skills` | 取得老師技能列表 |
| POST | `/api/v1/teacher/me/skills` | 新增老師技能 |
| PUT | `/api/v1/teacher/me/skills/{id}` | 更新老師技能 |
| DELETE | `/api/v1/teacher/me/skills/{id}` | 刪除老師技能 |
| GET | `/api/v1/teacher/schedules` | 取得老師的課表 |
| POST | `/api/v1/teacher/scheduling/check-rule-lock` | 檢查老師規則鎖定狀態 |
| POST | `/api/v1/teacher/scheduling/delete-recurring` | 刪除循環排課 |
| POST | `/api/v1/teacher/scheduling/edit-recurring` | 編輯循環排課 |
| POST | `/api/v1/teacher/scheduling/preview-recurrence-edit` | 預覽循環編輯影響範圍 |

## Teacher - Events

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/teacher/me/personal-events` | 取得老師個人行程列表 |
| POST | `/api/v1/teacher/me/personal-events` | 新增老師個人行程 |
| DELETE | `/api/v1/teacher/me/personal-events/{id}` | 刪除老師個人行程 |
| PATCH | `/api/v1/teacher/me/personal-events/{id}` | 更新老師個人行程 |
| GET | `/api/v1/teacher/me/personal-events/{id}/note` | 取得個人行程備註 |
| PUT | `/api/v1/teacher/me/personal-events/{id}/note` | 更新個人行程備註 |

## Teacher - Export

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/teacher/me/backgrounds` | 取得自訂背景圖列表 |
| POST | `/api/v1/teacher/me/backgrounds` | 上傳自訂背景圖 |
| DELETE | `/api/v1/teacher/me/backgrounds` | 刪除自訂背景圖 |
| GET | `/api/v1/teacher/me/schedule.ics` | 匯出課表為 iCalendar 格式 |
| GET | `/api/v1/teacher/me/schedule/image` | 匯出課表為圖片 |
| POST | `/api/v1/teacher/me/schedule/subscription` | 建立課表訂閱連結 |
| DELETE | `/api/v1/teacher/me/schedule/subscription` | 取消課表訂閱 |

## Teacher - Invitations

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/teacher/me/invitations` | 取得老師的邀請列表 |
| GET | `/api/v1/teacher/me/invitations/pending-count` | 取得待處理邀請數量 |
| POST | `/api/v1/teacher/me/invitations/respond` | 老師回應邀請 |

## Teacher - Scheduling

| Method | Endpoint | Summary |
|:---|:---|:---|
| POST | `/api/v1/teacher/me/exceptions` | 建立新的例外申請 |

## Teacher - Sessions

| Method | Endpoint | Summary |
|:---|:---|:---|
| GET | `/api/v1/teacher/sessions/note` | 取得課堂筆記 |
| PUT | `/api/v1/teacher/sessions/note` | 新增或更新課堂筆記 |

