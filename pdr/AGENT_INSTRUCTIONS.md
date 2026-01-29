# Cursor Agent 優化實作指令集

以下指令設計用於引導 Cursor Agent 逐步完成 `pdr/task.md` 中的優化任務。考慮到 200k token 限制，建議**分段執行**，每一段落完成後請要求 Agent 提供總結。

---

## 任務 1：解決 ExpandRules 中的 N+1 查詢問題 (Task ID: 17)

**指令：**
> 請分析 `app/services/scheduling_expansion.go` 中的 `ExpandRules` 方法。
>
> **目標：** 消除迴圈中的 `s.exceptionRepo.GetByRuleIDAndDateStr` 呼叫。
>
> **步驟：**
> 1. 在 `ExpandRules` 進入日期迴圈前，收集所有傳入 `rules` 的 ID。
> 2. 呼叫 `exceptionRepo` 批次取得該範圍內 (`startDate` 到 `endDate`) 所有相關的例外資料。 (可能需要新增 Repository 方法 `GetByRuleIDsAndDateRange`)。
> 3. 將取得的結果建立成 `map[uint]map[string][]models.ScheduleException` (RuleID -> DateString -> Exceptions)。
> 4. 修改日期迴圈，改從 map 中讀取例外資料。
> 5. 確保邏輯與原先一致（處理 PENDING/APPROVED 狀態），並進行單元測試（若有提供）。
>
> **完成後：** 請提供修改前後的效能預期差異總結。

---

## 任務 2：重構 TeacherController 為多個領域 Controller (Task ID: 20)

**指令：**
> `app/controllers/teacher.go` 目前過於龐大，請進行拆分重構。
>
> **目標：** 將 `TeacherController` 職責拆分為 `Profile`, `Schedule`, `Exception` 三個部分。
>
> **步驟：**
> 1. **建立新檔案：** `teacher_profile.go`, `teacher_schedule.go`, `teacher_exception.go`。
> 2. **搬移邏輯：**
>    - `Profile`: 搬移 `GetProfile`, `UpdateProfile`, `SearchHashtags`, `CreateHashtag`, `GetCenters`, `GetSkills`, `GetCertificates` 相關方法。
>    - `Schedule`: 搬移 `GetSchedule`, `GetSchedules`, `GetCenterScheduleRules` 相關方法。
>    - `Exception`: 搬移 `CreateException`, `RevokeException` 相關方法。
> 3. **清理 `NewTeacherController`：** 確保每個子 Controller 只初始化其必要的依賴。
> 4. **更新路由：** 在 `apis/base.go` 或相關路由註冊處，更新對應的路由綁定。
>
> **注意：** 由於檔案極大，請一次處理一個 sub-controller，完成後向我確認再進行下一個。

---

## 任務 3：實作 DTO 模式與領域邏輯封裝 (Task ID: 13, 14)

**指令：**
> 請針對 `TeacherProfile` 相關的功能實作 DTO 模式，並將邏輯移至 Model。
>
> **目標：** 讓 API 傳輸資料 (`resources`) 與資料庫模型 (`models`) 完全分離。
>
> **步驟：**
> 1. 在 `app/resources/teacher.go` 中定義 `TeacherProfileResponse`。
> 2. 將 `TeacherController` 中手動拼湊 JSON map 的邏輯，改為呼叫 `teacher.ToResource()` 方法。
> 3. 在 `models/teacher.go` 中實作業務邏輯，例如 `CanEditProfile()`。
> 4. 在 `app/requests/teacher.go` 定義專用的 `UpdateProfileRequest`，並使用 `ShouldBindJSON`。
>
> **完成後：** 總結 DTO 帶來的安全性（防止 Mass Assignment）效益。

---

## 任務 4：統一錯誤處理與 AppError (Task ID: 15)

**指令：**
> 請為專案建立統一的錯誤處理架構。
>
> **目標：** 消除 Controller 中重複的 `global.ApiResponse` 拼湊，改用全域 Middleware 處理。
>
> **步驟：**
> 1. 在 `global/errors` 建立 `AppError` 結構（包含 Code, Message, HTTPStatus）。
> 2. 在 `app/middleware` 建立 `ErrorHandlerMiddleware`，捕捉 `ctx.Error` 並統一輸出 JSON。
> 3. 修改 `TeacherController` 中的一個方法作為範例，將 `ctx.JSON(500, ...)` 改為 `ctx.Error(errors.NewSystemError(...))`。
>
> **完成後：** 說明如何新增一個新的業務錯誤碼。

---

## 任務 5：Redis 快取與事件驅動架構 (Task ID: 18, 19)

**指令：**
> 這是一個較大的架構變動，請分兩步執行：
>
> **第一步 (Cache)：**
> 1. 為 `ScheduleQueryService` 實作 Redis 快取邏輯。
> 2. 確保在 `GetTeacherSchedule` 時優先查詢 Redis。
>
> **第二步 (Event Bus)：**
> 1. 在 `libs` 建立一個簡單的 `InternalEventBus`。
> 2. 修改 `ExceptionService`，在申請成功後發布 `EVENT_EXCEPTION_CREATED`。
> 3. 建立訂閱者來處理原本的 LINE 通知發送。
>
> **注意：** 每個步驟完成後請提供代碼片段與重構後的依賴圖概括。
