# 排課碰撞引擎測試矩陣 (Scheduling Engine Test Matrix)

本文件定義了「排課驗證引擎」必須通過的所有核心測試案例，旨在確保系統在併發環境與極端時段下的穩定性。

---

## 1. 硬重疊檢查 (Hard Overlap Tests)
**原則**：同一個資源（老師或教室）在同一時間點絕對不能有兩筆 Active 的行程。

| ID | 測試場景 | 輸入數據 | 預期結果 | 關鍵邏輯 |
|:---|:---|:---|:---|:---|
| T1.1 | **標準重疊** | 已有 [10:00-11:00], 排 [10:30-11:30] | ❌ `E_SCHED_OVERLAP` | `S_old < E_new AND E_old > S_new` |
| T1.2 | **包含重疊** | 已有 [10:00-12:00], 排 [10:30-11:30] | ❌ `E_SCHED_OVERLAP` | 舊包含新 |
| T2.3 | **被包含重疊** | 已有 [10:30-11:30], 排 [10:00-12:00] | ❌ `E_SCHED_OVERLAP` | 新包含舊 |
| T1.4 | **精確重合** | 已有 [10:00-11:00], 排 [10:00-11:00] | ❌ `E_SCHED_OVERLAP` | 完全相同 |
| T1.5 | **臨界觸碰 (Pass)**| 已有 [10:00-11:00], 排 [11:00-12:00] | ✅ `SUCCESS` | 結尾與開始剛好銜接 |
| T1.6 | **跨日重疊** | 已有 [23:00-01:00], 排 [00:30-02:00] | ❌ `E_SCHED_OVERLAP` | 需正確處理 LocalDate 邊界 |
| T1.7 | **週期階段重疊**| 規則 A (1-2月) [週二 10:00], 排 規則 B (2-3月) [週二 10:00] | ❌ `E_SCHED_OVERLAP` | 偵測 `effective_range` 的交集 |

---

## 2. 緩衝時間檢查 (Buffer Violation Tests)
**原則**：資源轉換需要時間（老師轉場、教室清潔）。

| ID | 測試場景 | 輸入數據 | 預期結果 | 關鍵邏輯 |
|:---|:---|:---|:---|:---|
| T2.1 | **老師緩衝不足** | 前課 11:00 完, Buffer 15, 排 11:10 | 🟠 `E_SCHED_BUFFER` | `New.Start - Prev.End < T_Buf` |
| T2.2 | **教室緩衝不足** | 前課 11:00 完, Buffer 10, 排 11:05 | 🟠 `E_SCHED_BUFFER` | `New.Start - Prev.End < R_Buf` |
| T2.3 | **後方銜接緩衝** | 排 11:00-12:00, 後課 12:05 始, Buf 10| 🟠 `E_SCHED_BUFFER` | 需同時檢查 Prev 與 Next |
| T2.4 | **草稿排課 (Draft)**| `teacher_id` 為空, 但教室衝突 | 🟠 `E_SCHED_BUFFER` | 跳過老師檢查，不跳過教室 |
| T2.5 | **覆寫(Override)** | Buffer 衝突 + `override=true` | ✅ `SUCCESS` | 僅管理員或具權限班別可過 |
| T2.6 | **教室容量爆滿** | Room Cap=1, 已排 [10:00], 排 [10:00] | ❌ `E_SCHED_OVERLAP` | 教室資源被完全佔用 |
| T2.7 | **多班共用教室** | Room Cap=20, 班A 10人, 班B 15人 | ❌ `E_ROOM_CAPACITY` | 總人數超過教室負荷 |

---

## 3. 併發與鎖定 (Concurrency & Locking Tests)
**原則**：防止 Race Condition 導致的 Double Booking。

| ID | 測試場景 | 模擬行為 | 預期結果 | 實作驗證 |
|:---|:---|:---|:---|:---|
| T3.1 | **同時搶位測試** | 2 個 Request 同時排同老師同轉段 | 1 成功, 1 失敗 (`E_SCHED_LOCKED` 或 `OVERLAP`) | `FOR UPDATE` 排隊機制 |
| T3.2 | **死鎖預防測試** | A 排 T1->T2, B 排 T2->T1 | 兩筆皆正確完成或超時 | 確保 Lock 獲取順序一致性 |
| T3.3 | **Read-committed**| 排課檢測時不應讀到「已開始但未提交」的衝突 | 隔離等級驗證 | Transaction Isolation Level |

---

## 4. 例外與循環展開 (Exception & Recurrence Tests)
**原則**：例外單 (Exception) 必須能精確覆蓋或修改週期規則 (Rule)。

| ID | 測試場景 | 輸入數據 | 預期結果 |
|:---|:---|:---|:---|
| T4.1 | **停課屏蔽測試** | Rule 每週一, 1/20 有 CANCEL 例外 | 查詢 1/20 課表應不含此場次 |
| T4.2 | **改期取代測試** | Rule 10:00, 1/20 RESCHEDULE 到 14:00 | 10:00 消失, 14:00 出現新場次 |
| T4.3 | **審核併發防撞** | 管理員審核 Reschedule 申請時，新時段已被佔用 | ❌ `E_SCHED_OVERLAP` (禁止核准) |
| T4.4 | **循環終止測試** | Rule 設定至 2026/06/30 | 查詢 2026/07/01 應無資料 |
| T4.5 | **階段切換展開** | 規則 A 至 2/28, 規則 B 自 3/1 始 | 查詢 3/1 後應正確切換到規則 B |
| T4.6 | **假日自動停課** | 中心設 2/14 為假日, 規則 A 每週六有課 | 查詢 2/14 課表應為空 (或標註假日) |
| T4.7 | **異動截止鎖定** | lock_at 已過, 老師嘗試 POST /exceptions | ❌ `E_FORBIDDEN` |
| T4.8 | **管理員豁免鎖定** | lock_at 已過, 管理員 POST /exceptions | ✅ `SUCCESS` |
---

## 6. 全疊整合流程測試 (Full-Stack Integration Flows)
**原則**：驗證跨模組、前端與後端、以及外部服務 (LINE) 的協同工作。

| ID | 整合流程 | 涉及組件 | 驗證重點 |
|:---|:---|:---|:---|
| I1.1 | **邀請至註冊全流程** | Admin UI -> API -> LINE -> Teacher UI | 邀請碼生成 -> LINE 通知接收 -> 點擊連結 -> 自動綁定中心。 |
| I1.2 | **排課至提醒閉環** | Admin Grid -> Validate -> DB -> Worker -> LINE | Admin 排課儲存 -> 背景任務偵測 -> 晚上 20:00 發送準確的 LINE 提醒。 |
| I1.3 | **異動審核與即時通知**| Teacher App -> Exception API -> Admin Review -> LINE| 老師申請 -> 管理員即時收到待審 -> 管理員核准 -> 老師即時收到核准通知。 |
| I1.4 | **下期預排與鎖定期** | Admin Settings -> Rule Phase -> Teacher App | 設定 4/15 截止 -> 4/14 老師可異動 -> 4/16 老師端功能鎖定。 |
| I1.5 | **多中心衝突綜合視圖**| Center A + Center B -> Teacher App | A 中心排課 + B 中心排課 -> 老師 App 顯示統一課表且顏色正確區分。 |

---

## 7. 前後端開發契約驗證 (Frontend-Backend Contract)
*   **端對端 Mocking 策略**：前端使用 `MSW` 或 `Vite Proxy` 模擬後端回傳，確保介面在後端開發中也能並行測試。
*   **資料一致性**：驗證 `effective_start/end` 日期格式在前後端傳輸中無時區偏移 (Timezone Drift)。
*   **錯誤處理回饋**：當後端回傳 `E_SCHED_OVERLAP` 時，前端必須精確顯示「紅框」與該衝突場次的細節，而非通用錯誤。

---

## 8. 邊緣案例與擴充功能測試 (Edge Cases & Extended Features)

### 8.1 課程複製測試 (Course Copy)
| ID | 測試場景 | 輸入數據 | 預期結果 |
|:---|:---|:---|:---|
| T8.1 | 標準複製 | 複製班別，含 5 條規則 | 新班別 ID 不同，規則數量正確 |
| T8.2 | 複製不帶規則 | copy_rules=false | 新班別無規則，需重新建立 |
| T8.3 | 複製並更換老師 | 設定 new_teacher_id | 新班別老師為指定人員 |
| T8.4 | 複製至新日期範圍 | effective_start=2026-03-01 | 新規則的 effective_start 正確 |
| T8.5 | 複製衝突 | 新範圍內已有同班別課程 | 衝突檢測正常運作 |

### 8.2 循環編輯測試 (Recurrence Edit)
| ID | 測試場景 | 輸入數據 | 預期結果 |
|:---|:---|:---|:---|
| T8.6 | 編輯單一場次 | update_mode=SINGLE, start_at 改變 | 產生 CANCEL + ADD 兩筆例外 |
| T8.7 | 編輯未來場次 | update_mode=FUTURE | 原規則截斷，新規則從該場次開始 |
| T8.8 | 編輯整組循環 | update_mode=ALL, recurrence 改變 | recurrence 欄位更新，所有場次受影響 |
| T8.9 | 刪除單一場次 | DELETE, mode=SINGLE | 該場次標記為 CANCEL |
| T8.10 | 刪除未來場次 | DELETE, mode=FUTURE | 原規則截斷，該場次後刪除 |
| T8.11 | 刪除整組循環 | DELETE, mode=ALL | schedule_rules 記錄刪除 |

### 8.3 多中心配色測試 (Multi-Center Colors)
| ID | 測試場景 | 預期結果 |
|:---|:---|:---|
| T8.12 | 同色系中心 | 亮度自動調整，顏色可區分 |
| T8.13 | 中心篩選 | 點擊圖例可 Toggle 顯示/隱藏 |
| T8.14 | 課表卡片標籤 | 顯示中心名稱縮寫 |
| T8.15 | 導出圖片 | 維持原有中心配色 |

### 8.4 分頁與排序測試 (Pagination & Sorting)
| ID | 測試場景 | 預期結果 |
|:---|:---|:---|
| T8.16 | 第一頁 | has_prev=false, has_next=true |
| T8.17 | 中間頁 | has_prev=true, has_next=true |
| T8.18 | 最後一頁 | has_prev=true, has_next=false |
| T8.19 | 排序切換 | 資料依排序欄位正確排列 |
| T8.20 | 超出範圍 page | 回傳空 data, pagination 正確 |

### 8.5 軟刪除測試 (Soft Delete)
| ID | 測試場景 | 預期結果 |
|:---|:---|:---|
| T8.21 | 刪除課程模板 | is_active=false, 列表不再顯示 |
| T8.22 | 恢復課程模板 | PATCH is_active=true, 列表重新顯示 |
| T8.23 | 刪除有關聯的課程 | 回傳 E_COURSE_IN_USE |
| T8.24 | 刪除空白教室 | 軟刪除成功 |
| T8.25 | 查詢時過濾軟刪除 | 查詢自動過濾 is_active=false |

### 8.6 例外撤回測試 (Exception Revoke)
| ID | 測試場景 | 預期結果 |
|:---|:---|:---|
| T8.26 | 撤回待審核例外 | status 變為 REVOKED |
| T8.27 | 撤回已核准例外 | 回傳 E_INVALID_STATUS |
| T8.28 | 撤回他人例外 | 回傳 E_FORBIDDEN |

---

## 5. 使用者體驗細節 (UX/Edge Case Tests)
*   **1 分鐘銜接**：[09:00:00 - 09:59:59] 與 [10:00:00] 是否判定為不重疊？
*   **夏季時制 (若適用)**：台灣目前雖無，但 API 應能正確處理 `UTC+8` 位移。
*   **極速連續點擊**：前端在網速慢時點了兩次「存檔」，後端應具備冪等性 (Idempotency) 或鎖定。
