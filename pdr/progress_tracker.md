# TimeLedger 實作進度追蹤 (Development Progress Tracker)

> [!IMPORTANT]
> 此文件由 AI 持續維護。每完成一個任務或階段，請在此更新狀態與「上下文恢復快照」。

## 1. 階段性進度表 (Roadmap Status)

| 階段 | 任務說明 | 狀態 | 備註 |
|:---|:---|:---:|:---|
| **Stage 1** | **資料庫架構與種子資料** | `[/] IN_PROGRESS` | 準備開始實作 |
| **Stage 2** | **認證與基礎 Profile API** | `[ ] TODO` | |
| **Stage 3** | **排課引擎 (魔王關)** | `[ ] TODO` | |
| **Stage 4** | **人才搜尋與智慧媒合** | `[ ] TODO` | |
| **Stage 5** | **UI/UX 拋光與匯出功能** | `[ ] TODO` | |
| **Stage 6** | **E2E 測試與部署** | `[ ] TODO` | |

---

## 2. 當前上下文快照 (Context Snapshot)
*   **最後更新**: 2026-01-20
*   **當前目標**: 執行 Stage 1：根據 `pdr/Mysql.md` 建立 Go Models 與 Migration。
*   **關鍵架構決策**: 
    - 採用顆粒化緩衝 (Course-level buffers)。
    - 標籤雙層體系 (Expertise bound to Skill + Personal Branding Tags)。
*   **已上線 Skills**: `contract-sync`, `scheduling-validator`, `auth-adapter-guard`。

---

## 3. 已完成任務詳情 (Completed Tasks - Audit Log)
*   *(尚未開始，等待 Stage 1 啟動)*
