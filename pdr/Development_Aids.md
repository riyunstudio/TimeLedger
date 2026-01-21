# TimeLedger 開發輔助與維護規範 (Development Aids)

為了確保 TimeLedger 在快速開發 Stage 1~6 時不失去可維護性，建議在開發環境 (Cursor/Claude) 中配置以下 **Agent Skills** 與 **MCP Servers**。

---

## 1. 建議的 Agent Skills (自訂指令)

建議在 `.agent/skills/` (或 `.cursorrules`) 中建立以下自動化規範：

### 1.1 契約優先同步 (Contract-First Sync)
*   **Skill 檔案**：[.agent/skills/contract-sync.md](file:///.agent/skills/contract-sync.md)
*   **功能**：每次修改 `pdr/API.md` 後，強制 AI 立刻更新 Go 的 Request/Response Structs 以及 Nuxt 的 TypeScript Interfaces (`types/api.d.ts`)。
*   **維護性價值**：防止「後端改了名，前端不知道」的低級 Bug，確保前後端永遠對接。

### 1.2 排課引擎 TDD 執法 (Logic Tester)
*   **Skill 檔案**：[.agent/skills/scheduling-validator.md](file:///.agent/skills/scheduling-validator.md)
*   **功能**：針對 `擴大規則 (expander)` 與 `衝突檢查 (validator)` 模組，強制在寫 Logic 之前先寫 Unit Test，且必須涵蓋「跨日課程」、「緩衝邊界」、「跨閏年」等極端案例。
*   **維護性價值**：排課邏輯是系統心臟，透過 TDD (測試驅動開發) 確保未來修改時不發生 Regression。

### 1.3 認證適配器守衛 (Auth Adapter Guard)
*   **Skill 檔案**：[.agent/skills/auth-adapter-guard.md](file:///.agent/skills/auth-adapter-guard.md)
*   **功能**：管理 Mock 登入與真實 LINE 登入之間的抽象層，確保「開發快、切換穩」。
*   **維護性價值**：實作 Mock 優先策略，並在代碼層級強制隔離敏感的內部評價資料。

---

## 2. 推薦的 MCP Servers (工具擴充)

如果您使用的是支援 MCP (Model Context Protocol) 的 AI 客戶端，建議連接以下伺服器：

1.  **MySQL MCP**:
    *   **用途**：讓 AI 直接下 SQL 檢查 Live DB 的 Seed Data 是否正確。
    *   **場景**：開發 Stage 3 時，讓 AI 幫忙 Debug 為什麼某筆排課會衝突，直接查表最準。
2.  **GitHub/Gitlab MCP**:
    *   **用途**：讓 AI 讀取 PR 或 Issue。
    *   **場景**：Stage 6 測試期，讓 AI 自動解析 Bug Report 並定位代碼修改點。
3.  **Memory / Knowledge Graph MCP**:
    *   **用途**：將 `/pdr` 資料夾內容索引進知識圖譜。
    *   **場景**：當專案變大時，AI 常會忘記前期的戰略細節 (如 LINE Only)。透過 Memory MCP 可以確保 AI 始終維持全局視角。

---

## 3. 現有資產優化

您已有的 `backend_specs/_cursor_prompts/backend_module_generator.md` 非常優秀。
**建議調整**：將其路徑提升至 `.cursorrules` 或更顯眼的位置，並在 `MASTER_PROMPT.md` 中將其列為 **「每一行 Backend 代碼都必須遵守」** 的鐵律。

---

---

## 4. AI 前端開發成功策略 (AI-Frontend Synergy)

為了避免 AI 在實作前端時發生「視覺不一致」或「邏輯混亂」，請遵循以下策略：

### 4.1 核心設計標記 (Design Tokens)
*   **規則**：禁止 AI 在組件中寫死 `Hex Code` (如 `#6366f1`)。
*   **執行**：在 `frontend/tailwind.config.ts` 中完成標記定義，AI 必須僅使用 `text-brand-primary` 或 `bg-midnight-glass` 等語義化 Class。

### 4.2 組件庫優先 (Atomic UI)
*   **規範**：優先指導 AI 建立 `BaseButton.vue`, `BaseModal.vue`, `BaseCard.vue`。
*   **AI 提示語**：在開始任何頁面開發前，先要求 AI 閱讀 `pdr/UiUX.md` 並生成一套基礎組件。這能確保後續所有頁面的「WOW 質感」一致。

### 4.3 邏輯與介面分離 (Logic Extraction)
*   **規範**：排課網格的「座標計算」、「時段衝突判定」等純運算邏輯，禁止寫在 `.vue` 的 `<script>` 內。
*   **執行**：要求 AI 將其抽離至 `frontend/utils/scheduling.ts`。這不僅能讓 AI 更精準地撰寫測試，也避免單一 Vue 檔案過大導致 AI 遺忘語境。

### 4.4 模擬環境先行 (LIFF Mocking)
*   **規範**：封裝 `useLiff` composable。在非 LINE 環境下，自動回傳 Mock User 資料。
*   **維護性價值**：讓 AI 能在標準瀏覽器環境下完成 99% 的介面開發，而不需要依賴 LIFF 真實環境。

---

## 5. 維護性黃金法則 (AI-Driven Development)

1.  **拒絕大檔案**：強制單一 `.go` 或 `.vue` 檔案不超過 300 行。
2.  **自帶註解**：要求 AI 在 Service 層標註對應的 `pdr/功能業務邏輯.md` 的章節編號。
3.  **版本化專案檔**：保持 `task.md` 隨時更新，這能讓接手的 AI 快速進入狀況。
4.  **語系一致性**：強制 AI 使用中文與開發團隊溝通。
