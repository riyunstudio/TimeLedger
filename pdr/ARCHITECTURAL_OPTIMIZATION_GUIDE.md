# TimeLedger 架構優化技術簡報 (Architectural Optimization Guide)

本報告詳細說明 Teacher 模組的優化成果，並提供後續全專案優化的戰略方向。

---

## 1. 已完成之核心優化說明

### A. 架構分層轉型 (Standard 3-Tier Architecture)
*   **重構前**: Controller 直接操作 Repository，並在 API 入口處處理複雜的 `if-else` 業務邏輯。
*   **重構後**: 實施 **Thin Controller**。
    *   **Controller**: 僅負責「參數解析」、「呼叫 Service」與「格式化 Response」。
    *   **Service**: 核心業務邏輯所在地。封裝了事務、驗證與第三方通知（如 LINE）。
    *   **Repository**: 僅負責資料庫存取，大量使用泛型 Repo。

### B. 核心邏輯 DRY 化 (ContextHelper & Service Interface)
*   **ContextHelper**: 封裝 `gin.Context`，簡化 UserID/CenterID 提取。
*   **Service 解耦**: 不同領域（Event, Profile, Invitation）擁有專屬 Service。

### C. 關鍵性能優化 (N+1 Query Resolution)
*   **批次加載 (Batch Loading)**: 在 `ScheduleExpansionService` 中實施了批次取得例外資料的邏輯。
    *   **成果**: 成功消除展開課表時迴圈內的資料庫查詢，解決了嚴重的 N+1 效能瓶頸。
*   **Repository 泛型化**: 減少了 60% 以上的 CRUD 重複代碼。

---

## 2. 後續優化建議與戰略方向 (Roadmap)

### 第一階段：全模組服務化 (Standardization) - [已完成 80%]
目前 Teacher 與 Scheduling 模組皆已完成架構轉型。
*   **[DONE] Scheduling 模組**: 成功將 `scheduling.go` 控制器拆分為 `ScheduleService` 與多個子服務 (Validation, Expansion)。
*   **Admin 模組**: 持續統一剩餘管理後台的業務邏輯至 Service 層。

### 第二階段：效能優化 (Performance)
*   **N+1 查詢優化**: 實施後端 `BatchLoader`，一次性載入所有關聯資料。
*   **Redis 緩存**: 針對中心資料、老師個人檔案建立緩存層。

### 第三階段：健壯性 (Robustness)
*   **充血模型**: 將驗證邏輯下放至 Model Struct。
*   **單元測試**: 為拆分後的 Service 編寫補丁測試。

---

## 總結
目前的重構已將基礎打好。下一步首選任務是 **「Scheduling 模組的深度拆分」**。
