# 專案規模評估與待確認事項 (Project Assessment)

本文件針對 TimeLedger 的開發複雜度進行評級，並列出在進入 Stage 0 開發前建議確認的外部依賴。

---

## 1. 開發前準備清單 (Readiness Checklist)

為了讓 `claude-code` 或其他 AI 順利開發，請務必先準備以下資訊與環境：

### 1.1 環境準備 (Local Env)
*   **Go 1.24+**: 後端運作基礎。
*   **Node.js 20+ & npm**: 前端 Nuxt 3 需要。
*   **Docker Desktop**: 運行 MySQL 與 Redis。
*   **LINE 官方帳號**: 需擁有一個管理員權限的 LINE OA 帳號。

### 1.2 LINE 開發者憑證 (Essential for Stage 2 & 3)
請至 [LINE Developers](https://developers.line.biz/) 建立並取得：
*   **LINE Login Channel**:
    *   `LIFF ID` (用於前端自動登入)。
    *   `Channel ID` / `Channel Secret` (用於後端換取身分)。
*   **Messaging API Channel**:
    *   `Channel Access Token` / `Secret` (用於發送推播通知)。

### 1.3 品牌與設計資產 (Design Assets)
*   **Logo**: 至少一張 512x512 的方形圖檔 (用於 LINE 登入授權頁與 App 圖示)。
*   **主色調 (Primary Color)**: 例如核心藍或專業綠的 Hex Code。
*   **網域名稱 (Domain)**:
    - **正式環境**: 為必備項目 (需 HTTPS)，用於正式對接 LINE Login。
    - **開發階段 (無網域解決方案)**: 可以使用 **Local Tunnel (如 ngrok 或 Cloudflare Tunnel)**。這能將您的本地開發環境對外暴露成一個暫時的 HTTPS 網址，用於測試 LINE Login 回調。
    - **影響**: Stage 1 不受影響；Stage 2 (LINE 登入) 開始需要有測試網址。

### 1.4 雲端主機 (Deploy target)
*   **VPS 帳號 (推薦 DigitalOcean/Linode)**: 用於部署單體 Docker Compose。
*   **S3-Compatible Storage**: 推薦使用 **Cloudflare R2** (免費額度大)，用於存放老師證照圖檔。

---

## 2. 專案規模評估 (Size & Complexity)

### 2.1 評級結論：中型偏大 (Medium-Large / Vertical SaaS)
這**絕對不是**一個「小型」專案 (如 Todo List 或 形象官網)。它屬於 **「垂直領域 SaaS 核心系統」**。

| 維度 | 等級 | 原因分析 |
|:---|:---:|:---|
| **業務邏輯 (Backend)** | **High** | 排課涉及「時間重疊」、「緩衝區計算」、「並發鎖定」、「遞迴規則展開」。這是演算法等級的邏輯，非單純 CRUD。 |
| **互動體驗 (Frontend)** | **High** | 涉及「網格拖曳 (Drag & Drop)」、「手機手勢」、「即時衝突檢測 (Optimistic UI)」。前端狀態管理極為複雜。 |
| **資料關聯 (Database)** | **Medium** | 表數量約 10-15 張，但關聯性強 (多對多、遞迴查表)。 |
| **權限架構 (Security)** | **Medium** | 需處理多租戶 (Multi-tenant) 資料隔離與 RBAC。 |

### 2.2 風險與挑戰
1.  **Stage 3 (排課引擎) 是魔王關**:
    *   這是系統的心臟。若邏輯不對 (例如：沒算到清潔時間)，整個系統就不可用。建議此階段投入 40% 的後端開發精力。
2.  **手機版 Grid 的實作難度**:
    *   在手機狹窄螢幕上做「拖曳排課」極具挑戰。我們採用的「3日視圖」策略能緩解此問題，但仍需大量調校 CSS 與 Touch Event。

### 2.3 給業主的建議
*   **不要急著做全**: 嚴格遵守 **MVP (Minimum Viable Product)** 範疇。
    *   先做「能排課、能防撞」。
    *   「智慧媒合」、「匯出美圖」等亮點功能放在 Stage 4, 5，不要在一開始就卡住核心開發。
*   **測試至上**: 這種系統 **不能有 Bug** (會導致中心營運事故)。Stage 6 的 E2E 測試必須覆蓋所有排課場景。

---

## 3. 技術選型建議 (Frontend Tech Stack Rationale)
針對為何選擇 **Nuxt 3** (Vue Framework) 而非純 SPA (Vite) 或 Next.js 的決策分析：

### 3.1 為何推薦 Nuxt 3？
1.  **社群分享優化 (Link Previews)**: TimeLedger 有「分享行程」的功能。Nuxt 內建 SSR (Server-Side Rendering)，能讓分享到 LINE/FB 時顯示正確的 Open Graph 預覽圖/標題。純 SPA (Vite) 做不到這點。
2.  **開發速度 (Auto-imports)**: Nuxt 自動引入 Components 與 Composables，大幅減少 Boilerplate code。對於「單人/小團隊」開發中大型系統，效率極高。
3.  **路由與佈局 (File-based Routing & Layouts)**: 前台 (Teacher-Mobile) 與 後台 (Admin-Desktop) 需要截然不同的 Layout。Nuxt 的 `layouts/` 機制切換非常直覺。
4.  **Vue vs React**: 若後端是 Go，Vue 的心智模型 (Template-based) 通常對後端開發者更友善，且 Vue 3 Composition API 的邏輯複用能力已不輸 React Hooks。

### 3.2 替代方案比較
| 方案 | 適合場景 | TimeLedger 適用度 | 缺點 |
|:---|:---|:---:|:---|
| **Vite + Vue 3 (SPA)** | 純後台管理系統 (無需 SEO) | Medium | 分享行程連結到 LINE 時抓不到標題/縮圖 (因為是 CSR)。 |
| **Next.js (React)** | 團隊全是 React 專家 | High | 與 Nuxt 同級，但若團隊無 React 強烈偏好，Config 較繁瑣。 |
| **Nuxt 3** | **需兼顧 SEO/Sharing 的 App** | **Best** | 學習曲線稍高 (需懂 SSR 概念)，但帶來的效益值得。

---

## 4. 關鍵戰略決策 (Strategic Decisions)

### 4.1 身份認證評估 (LINE Only vs. Hybrid)
針對「是否需要讓老師綁定帳密」的評估結果：

| 方案 | 描述 | 優點 (Pros) | 缺點 (Cons) | 結論 |
|:---|:---|:---|:---|:---|
| **Hybrid Mode** | LINE 登入後，強制或選填 Email/密碼。 | 當 LINE 當機時仍可登入。 | UX 摩擦力極大。老師容易忘記密碼。增加「忘記密碼/重設密碼」開發成本。 | ❌ 不採用 |
| **LINE Native** | **僅透過 LINE 登入** (含 LIFF Silent Login)。 | **極致順暢 (Zero Friction)**。不需維護密碼資安。符合台灣人使用習慣。 | 若 LINE 帳號遺失需人工由後台重綁。 | **✅ 採用** |

*   **決策關鍵**:
    *   TimeLedger 的核心價值是「快與方便」。任何阻擋老師進入系統的門檻 (如輸入密碼) 都應移除。
    *   **LIFF Silent Login**: 這是最強殺手鐧。老師只要從 LINE 官方帳號點擊選單，**無需任何登入動作** 直接進入課表。這是帳密系統做不到的體驗。

### 4.2 未來 App 開發評估 (Native App Strategy)
針對「是否開發獨立 App 脫離 LINE」的長遠規劃：

*   **現階段 (Phase 1: 0~10k Users)**: **堅決不做 App**。
    *   **原因**: 下載 App 的轉換率極低 (High Friction)。在品牌知名度未打開前，LINE 是最好的「寄生宿主」，能以最低成本獲取用戶。
*   **未來階段 (Phase 2: Stickiness High)**: **考慮開發 App (Flutter/React Native)**。
    *   **觸發點 (Trigger)**:
        1.  **通知費太貴**: 當 LINE Push 費用超過 App 推播維護成本時。
        2.  **功能天花板**: 用戶強烈需求 **桌面小工具 (iOS Widgets)** 或 **離線存取 (Offline Mode)** 時 (這些 Web 做不到)。
        2.  **功能天花板**: 用戶強烈需求 **桌面小工具 (iOS Widgets)** 或 **離線存取 (Offline Mode)** 時 (這些 Web 做不到)。
    *   **移轉策略**: 保持 LINE 為「輕量版」，App 為「專業版 (Pro)」。用「更強的功能 (如 iOS Widget)」誘使高黏著用戶主動下載 App，而非強迫轉移。

### 4.3 殘餘風險盤點 (Hidden Risks & Mitigations)
您提到的「隱憂」與「四處散亂」風險，我們已透過 PDR 文件標準化解決了大半。但仍需注意以下長期風險：

| 風險項目 | 描述 | 緩解策略 (Mitigation) |
|:---|:---|:---|
| **LINE 平台風險** | 若 LINE 修改收費規則或封鎖帳號，我們將失去老師入口。 | **Open ID 設計**: DB 設計不綁死 LINE，保留未來「補登 Email」的欄位。資料都在我們自家 VPS，不怕被綁架。 |
| **單點故障 (VPS)** | 單機部署若掛掉，全台中心都會停擺。 | **S3 異地備份**: 每日自動備份 SQL 到 Cloudflare R2。若主機掛掉，可在一小時內於新機器還原。 |
| **多中心隱私洩露** | A 中心是否會偷看到老師在 B 中心的課？ | **嚴格 API 遮蔽**: 系統層級對非所屬中心的事件，統一僅回傳 `{ status: "BUSY", title: "Busy" }`，絕不回傳學生姓名與筆記。 |

### 4.4 開發策略：Mock Auth 優先週期 (Mock Auth First Strategy)
針對您提到的「初期先 Mock 登入」以避免開發卡點，以下是詳細方案與隱憂評估：

#### **1. 實作方案 (Implementation)**
*   **後端**: 建立 `AuthService` Interface，初期注入 `MockAuthService`。提供一個 `/auth/dev-login` 接口，傳入 `user_id` 即可獲得該用戶的 JWT。
*   **前端**: 依據 `ENV` 切換，若為 `DEV` 則跳過 `liff.init()`，由 API 獲取固定 Token。

#### **2. 潛在隱憂與配套措施 (Risks & Mitigations)**
| 隱憂項目 | 影響 | 配套措施 (Mitigations) |
|:---|:---|:---|
| **Payload 差異** | LINE 傳回的 Claims (如 `sub`, `email`) 可能與 Mock 資料欄位不同。 | **統一 DTO**: 定義嚴格的 `IdentInfo` 物件，Auth Service 必須將 LINE 或 Mock 資料統一轉為此格式，後續邏輯不依賴原始 Payload。 |
| **自動註冊流程** | 本地開發跳過真實 Callback，可能漏掉「新用戶自動建立」的處理邏輯。 | **完整路徑 Mock**: `MockAuthService` 應模擬「查表 -> 若無則 Insert」的完整資料庫行為，而非僅回傳 Token。 |
| **LIFF SDK 依賴** | 前端有些邏輯依賴 LIFF 環境 (如獲取瀏覽器外殼資訊)。 | **抽象封裝**: 前端建立 `useAuth()` Composable，將 `liff` 調用封裝在內，在 Mock 模式下回傳預設值。 |
| **串接整合痛點** | 到 Stage 2 才接真實 LINE 可能發現設定問題 (如白名單、網域)。 | **提前驗證環境**: 雖使用 Mock 開發，但應在 **Stage 1 結束前** 先完成一個簡單的「LINE Connectivity Test」腳本，確保憑證有效。 |

#### **結論**:
**✅ 採用 Mock 優先策略** 是明智的。這能確保大腦 (Logic) 與 軀幹 (UI) 先跑起來，最後再接入 外接設備 (LINE)。只要遵循以上 **「抽象化封裝」** 原則，後續切換將會是「更換一組設定檔」等級的改動。
