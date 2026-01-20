# TimeLedger UI/UX 規劃書

## 1. 設計原則與風格 (Premium Design System)

為了達成令人驚艷 (WOW) 的使用者體驗，全系統採用 **「深藍高質感 (Midnight Premium)」** 結合 **「毛玻璃 (Glassmorphism)」** 風格。

### 1.1 核心視覺規格
- **設計風格**：
  - **Modern Glassmorphism**：卡片與彈出視窗採用 `backdrop-blur` 半透明質感，搭配細微的 `1px` 白色框線 (Border)。
  - **雙模式適配 (Dual-Theme Strategy)**：
    - **深色模式 (Midnight Premium)**：針對 20-35 歲族群，打造高科技與專業感。
    - **淺色模式 (Cloud Glass)**：針對 35-45 歲族群或高亮度環境，提供高對比度與清爽感，維持毛玻璃質感。
- **色彩系統 (The Palette)**：
  - **Primary (品牌藍)**：`#6366F1` (Indigo 500) - 用於主要按鈕、選起狀態。
  - **Secondary (漸層紫)**：`#A855F7` (Purple 500) - 與品牌藍形成補色漸層，用於匯出圖片背景。
  - **Success (翡翠綠)**：`#10B981` (Emerald 500) - 已核准、安全時段。
  - **Critical (野玫瑰紅)**：`#F43F5E` (Rose 500) - 硬衝突、取消申請。
  - **Warning (琥珀橘)**：`#F59E0B` (Amber 500) - 待審核、Buffer 警示。
- **字體 (Typography)**：
  - **主體字體**：`'Outfit', sans-serif` (Google Fonts) - 提供圓潤現代且具科技感的閱讀體驗。
  - **標題字體**：`'Inter', sans-serif` - 確保在各種螢幕大小下的清晰度。
- **動態效果 (Micro-animations)**：
  - **Spring Transitions**：所有行程卡片的彈出與收合使用 `stiffness: 300, damping: 20` 的彈力效果。
  - **Haptic Feedback**：手機版長按、拖曳時具備輕微觸控震動感。
  - **Smooth Inertia**：日曆滑動具備物理慣性與緩動機制。

### 2.0 公開登陸頁 (Official Landing Page - Pre-login) - *New*
**目標**：在 User 登入前，透過震撼的視覺與簡單範例吸引其使用。

#### A. 英雄區 (Hero Section)
- **視覺**：全螢幕的 **Mesh Gradient (動態網格漸層)** 背景。
- **Logo**：採用 [TimeLedger Premium Logo](file:///Users/chenlilin/.gemini/antigravity/brain/0f69a3ea-c23a-497c-b1f7-1458758bf464/timeledger_logo_premium_1768888114927.png) 作為核心識別。
- **標語**：中心顯示標語：「TimeLedger - 您的時間資產管家」。
- **Slogan**：多中心排課 × 老師人才庫 × 專業課表分享。
- **大按鈕 (CTA)**：「立刻使用 LINE 登入 (免費)」，下方搭配「管理員後台入口」。

#### B. 互動範例 (Interactive WOW Demo)
- **結構**：在頁面中央放置一個「迷你排課沙盒 (Mini Sandbox)」。
- **功能**：
  - User 不需要登入就能在網格上拖動一張範例課程卡片。
  - **即時反饋**：拖動到衝突區會顯示 **紅光/毛玻璃警告**，拖動到安全區顯示 **靛藍色吸附動畫**。
  - **設計意圖**：讓 User 10秒內感受到「排課治理」的流暢度與專業感。

#### C. 特色展示 (Feature Showcase)
- 採用 **毛玻璃卡片瀑布流**，展示：
  - **精美分享圖**：顯示 Mesh Gradient 的 IG Story 課表範例。
  - **智慧媒合**：展示 5 星評價老師的搜尋結果卡片。

---

## 2. 核心介面佈局 (Layout)

### 2.1 老師 / 個人端 (Mobile App-like Web)
**目標**：讓使用者能單手操作，快速查看行程與進行變動。

#### A. 登入頁面 (Login)
- **情境一：從 LINE App 開啟 (LIFF)**
  - **無介面 (Silent Login)**：系統自動偵測 LINE 身分並登入，直接進入「主儀表板」。
- **情境二：從外部瀏覽器開啟**
  - TimeLedger Logo
  - **大按鈕**：「使用 LINE 登入」
  - 連結：「切換至中心管理員登入」(跳轉至 Admin Login)

#### B. 主儀表板 (Dashboard - Calendar View)
- **頂部導覽列 (Top Bar)**：
  - 左側：漢堡選單 (切換不同中心 / 個人行程 / 設定)
  - 中間：當前月份/日期 (例：2026年 1月)
    - **導覽控制 (Nav)**：`< (上一週)` 與 `> (下一週)` 按鈕。手機版亦可透過滑動切換。
  - 右側：**匯出圖片按鈕 (Export Icon)**
- **日曆視圖 (Dynamic Calendar Grid)**：
  - **手機版 (Width < 768px)**：
    - **預設視圖**：**「近三天」網格圖 (3-Day View)**。因為手機寬度窄，顯示 7 天會過於擁擠，顯示 3 天能清楚呈現時段與課程名稱。
    - 支援左右滑動切換日期。
  - **平板/桌機版 (Width >= 768px)**：
    - **預設視圖**：**「一週」網格圖 (Week View)**。利用寬螢幕優勢，完整呈現週行程。
- **浮動動作按鈕 (FAB)**：
  - 「+」號按鈕：新增個人行程 (Personal Event)。

#### C. 新增/編輯個人行程 (Personal Event Modal)
- **視覺**：採用 **毛玻璃 (Glassmorphism)** 底色搭配高亮 Border。
- **彈出方式**：手機版由底部滑入 (Spring Bottom Sheet)，桌機版為置中 Modal。
- **欄位**：
  - 標題 (Title)
  - 時間 (Start/End Time)
  - **重複設定 (Recurrence)**：
    - 選項：不重複、每天、每週、每兩週、每月。
    - 自訂循環 (Custom)：設定每 X (日/週) 重複，結束於某日期或無期限。
  - 顏色標籤 (Color Tag)
  - 備註 (Notes)

#### D. 中心課程互動 (Exception Handling)
- 點擊行事曆上的「中心課程」圖塊：
  - 跳出詳情卡片 (Course Detail)。
  - 操作按鈕：
    - **請假/停課 (Cancel)**
    - **改期 (Reschedule)**
  - **教學紀錄 (Teaching Record)**：
    - **內容 (Content)**：輸入框，紀錄今日上課重點。
    - **備課 (Prep)**：輸入框，紀錄下次預計進度。
  - 狀態顯示：若由中心審核中，顯示「待審核 (Pending)」標籤。

#### E. 個人檔案與媒合設定 (Profile & Discovery)
- **基本資料**：頭像、姓名、Bio。
- **求職設定 (Job Seeker Settings)**：
  - **開關**：「開放中心搜尋我的檔案」(Be Discoverable)。
  - **公開資訊**：輸入「自訂聯繫方式 (Line/Phone/Email)」。此欄位僅在開啟媒合時對中心可見。
  - **證照權限**：開啟媒合時，證照將自動對搜尋到的中心可見 (唯讀)。

#### F. 匯出課表圖片 (Export as Image) - *Viral Feature*
- **設計目標**：極具質感的背景漸層 (Mesh Gradient)，讓老師引以為傲的高級感。
- **客製化選項**：
  - **主題選擇**：
    - **Midnight Glow** (深藍紫漸層)
    - **Emerald Mist** (翡翠綠毛玻璃)
    - **Sunset Quartz** (琥珀漸層)
  - **隱私模式 (Privacy Mode)**：
    - 切換開關：「隱藏行程細節」。
    - 開啟後，所有「個人行程」僅顯示「Busy」或「已保留」，中心課程維持顯示。方便老師公告空堂給學生，但保留隱私。
  - **格式**：自動生成適合 IG Story (9:16) 或 一般貼文 (4:3) 的尺寸。
- **操作流程**：點擊匯出 -> 預覽畫面 (調整選項) -> 下載/分享。

---

### 2.2 中心後台 (Center Admin - Desktop Dashboard)
**目標**：高效率的排課與審核，資訊密度高。

#### A. 管理員登入 (Admin Login)
- Email / Password 表單。
- 連結：「老師登入請點此」。

#### B. 排課主工作台 (Scheduling Grid)
- **三欄式佈局**：
  - **左欄 (資源庫)**：待排課程 (Offerings)、老師列表、教室列表。支援拖拉 (Drag & Drop) 或 點選指派。
  - **中欄 (排課網格)**：
    - 顯示全週/全教室/全老師 的時間軸。
    - 即時回饋：拖曳時，系統即時計算 Validate。
      - 🟢 綠色殘影：Safe
      - 🔴 紅色殘影：Hard Overlap (禁止放置)
      - 🟠 橘色殘影：Buffer Conflict (可放置但需 Override/Audit)
  - **右欄 (狀態/衝突詳情)**：
    - 顯示當前選取 Cell 的詳細資訊。
    - 若有衝突，顯示衝突原因 (例：該課程要求 10 分鐘 Room Buffer)。
    - **Override 開關**：若該 Offering 權限允許，顯示「強制排入」Checkbox 與「原因」輸入框。

#### C. 審核中心 (Approval Center)
- 列表式呈現所有 Exception (停課/改期) 申請。
- Before/After 對照：顯示 原本時間 vs 修改後時間。
- 衝突檢測：審核當下 **Re-validate**。若核准後會導致新衝突，需提示管理員 (Override or Reject)。

#### D. 智慧代課/排課媒合 (Smart Substitute Finder) - *Center Highlight*
- **觸發場景**：
  - 收到「請假申請」需要找代課老師時。
  - 點擊空白時段，選擇「尋找合適老師」。
- **介面 (Right Drawer)**：
  - **條件篩選**：自動帶入當前時段、課程技能標籤 (Tag)。
  - **候選名單 (Matching List)**：
    - 依 **Match Score** 排序。
    - 顯示資訊：頭像、姓名、**技能標籤**、**中心評分 (Rating ⭐️)**、**內部備註**。
    - **可用性**：
      - ✅ 可排課 (Available)
      - ⚠️ 有 Buffer 衝突 (可 Override)
      - ⛔️ 已有課 (Busy)
  - **快速操作**：點擊「指派」直接填入 Grid。

#### E. 人才庫搜尋 (Talent Search Page) - *UX Optimized*
- **功能**：主動尋找已開啟「求職媒合」的老師。
- **漸進式揭露 (Progressive Disclosure)**：
  - 預設僅顯示「大類選擇 (Categories)」。
  - 當大類被選取後，才動態展開對應的「子類/Hashtag」標籤雲。防止過多選項一次湧現。
- **篩選器**：
  - **大類選擇 (Categories)**：橫向滾動或下拉 (有氧、瑜伽、舞蹈、兒幼、TRX、空中瑜伽、技擊、樂齡、肌力、器械皮拉提斯)。
  - **地區篩選 (Region)**：下拉選單或行政區標籤 (如：台北市、新北市、大安區)。
  - **子類/Hashtag**：點擊類別後展開的標籤雲 (如：街舞、國標舞) 與熱門 #Hashtag。
  - **關鍵字**：姓名/Bio。
- **結果卡片**：
  - 顯示：頭像、Bio、**服務地區**、**技能分類** (例如 舞蹈 > 街舞)、**Hashtags**、**公開聯繫方式**、證照預覽。
  - 操作：**邀請加入中心** (發送 Invite Link)。

#### G. 管理員帳號與操作紀錄 (Admins & Audit) - *System Control*
- **管理員列表**：
  - 顯示：姓名、Email、權限等級 (Role)、狀態。
  - 操作：[新增管理員] 彈窗、[編輯]、[停用/刪除] 確認。
- **新增/編輯管理員彈窗 (Modal)**：
  - **欄位規格**：姓名 (Input)、Email (Input)、起始密碼 (Password Input)、權限角色 (Select: Owner/Admin/Staff)。
  - **交查邏輯**：Email 必須為唯一、權限分級檢查 (僅 Owner 可增刪 Admin/Staff)。
- **操作紀錄視圖 (Audit Log)**：
  - 列表式呈現：時間、**操作管理員 (Actor)**、動作 (Action)、對象 (Target)、備註 (如 Override 原因)。
  - 支援：按管理員名稱或日期範圍篩選。

#### F. 老師個人檔案編輯 (Profile Edit) - *UX Optimized*
- **區域選擇 (Region Picker)**：提供縣市與行政區的級聯選擇 (Cascading Selector)。
- **大類快速選取 (One-tap selection)**：提供常用大類圖示按鈕，點擊即自動填入類別，減少打字。
- **階層式輸入 (Hierarchical Flow)**：選取大類後，子類輸入框自動獲得焦點 (Focus) 並跳出該類別下的熱門建議。
- **Hashtag 輸入組件**：
  - **功能**：支援選取、刪除與動態新增標籤。
  - **模糊搜尋 (Debounced Search)**：輸入時延遲 300ms 呼叫 `/common/hashtags/search` 並彈出建議清單。
  - **動態新增**：若建議清單無相符項，按下 Enter 或逗號後自動將當前輸入轉為新標籤。
  - **視覺呈現**：以 Chip (標籤塊) 顯示，並提供 [x] 快速刪除，且限制單一標籤長度避免撐破排版。

### 2.3 多中心呈現策略 (Multi-Center Visualization)
當老師同時加入 Center A 與 Center B 時，介面需整合呈現以避免衝突。

1.  **綜合視圖 (Unified View - Default)**：
    *   **預設行為**：Dashboard 直接顯示「所有已加入中心」的課程 ＋ 個人行程。老師不需要切換帳號即可看到整週全貌。
    *   **區分方式**：
        *   **顏色區分**：不同中心分配不同色系的邊框 (Border) 或 標籤 (Badge)。
        *   **卡片內容**：課程卡片第一行顯示 `[Center Name]` 或縮寫。
2.  **狀態呈現 (Status)**：
    *   **審核中 (Pending)**：
        *   卡片背景顯示 **斜線陰影 (Striped)**。
        *   右下角顯示橙色 `⏳` icon。
        *   點開後的詳情頁會顯示「等待 [中心名稱] 管理員審核中」。
    *   **已拒絕 (Rejected)**：
        *   該異動申請顯示為灰色刪除線，並保留三天后自動消失，或需手動歸檔。
    *   **無老師 (Unassigned)**：
        *   卡片背景顯示 **灰色斜線 (Diagonal Stripes)** 或 **中灰色毛玻璃**。
        *   左上角顯示黃色 `⚠️` 警告標籤，並註明「未指派老師」。
        *   **交互行為**：點擊該卡片後，右側詳情面板優先顯示 **[智慧指派老師]** 按鈕，點擊後開啟代課媒合介面。


## 3. 關鍵互動細節 (Micro-interactions) - *UX Delighters*

1. **手機版手勢 (Mobile Gestures)**：
   - **長按空白處 (Long Press)**：快速新增行程 (Quick Add)，自動帶入預設 1 小時，並震動回饋 (Haptic)。
   - **拖曳修改 (Drag & Drop)**：在週/3日視圖中，長按行程後可拖曳至新時段 (僅限個人行程或允許的中心課程)。
   - **滑動切換 (Swipe)**：順暢的左右滑動切換週次，搭配 Inertia (慣性) 動畫。

2. **動態網格切換 (Responsive Grid)**：
   - 使用 CSS Grid 或 Flexbox，配合 Media Queries。
   - 斷點設為 `768px` (iPad Portrait / Large Phone Landscape)。
   - 小於斷點：`grid-template-columns: repeat(3, 1fr)` (顯示當天+未來2天)。
   - 大於斷點：`grid-template-columns: repeat(7, 1fr)` (顯示週日~週六)。

3. **驗證回饋 (Optimistic UI)**：
   - 用戶拖曳放開瞬間，立即吸附並顯示成功動畫。
   - 若後端 `validate` 回傳衝突，則優雅地彈回原位，並以 Toast 提示原因。

4. **匯出圖片 (Smart Export)**：
   - 避免直接 `html2canvas` 整個 DOM，因為手機版寬度太窄，生成的圖片不適合列印或電腦看。
   - **解法**：在背景 (Hidden Container) 渲染一個專門用於匯出的「標準寬度 DOM (例如 1200px 寬)」，填入相同資料，再對該 DOM 進行 `html2canvas` 截圖。確保匯出的課表永遠是整齊的「週視圖」。
## 4. 頁面佈局推薦 (Recommended Layouts)

### 4.1 老師行動端：主儀表板 (Mobile Dashboard)
- **結構**：
  1. **頂部導覽 (Top Nav)**：`h-14` 固定高度。左側選單、右側匯出圖片。
  2. **時間滑動條 (Week Bar)**：緊接導覽列下方，顯示當週日期，點擊可快速切換。
  3. **主內容區 (Scroll Area)**：
     - 上半部：**網格行程表 (3-Day Grid)**，滿版寬度，支援左右滑動。
     - 底部 (選配)：當日行程摘要列表 (List View)，方便快速瀏覽。
  4. **浮動按鈕 (FAB)**：右下角靛藍色圓形「+」按鈕，具備毛玻璃外圈及陰影。

### 4.2 中心後台：排課工作台 (Desktop Scheduling Workspace)
- **佈局核心**：100vh 高度，無網頁捲軸，各區塊獨立捲動。
- **結構**：
  1. **左側邊欄 (Sidebar)**：`w-64` 靛藍色深色底。收納：中心設定、老師管理、人才庫、審核中心。
  2. **上方過濾列 (Filter Bar)**：`h-16` 白色/淺灰毛玻璃。快速切換「全教室」或「指定老師」。
  3. **中央主視圖 (Core Workspace)**：
     - **左三 (資源面板)**：待指派課程與未排課老師清單。
     - **中六 (主網格)**：橫向為時間 (08:00 - 22:00)，縱向為 教室 或 老師。
     - **右三 (詳情與衝突)**：點擊格子後滑出 (Slide-over)，顯示驗證結果。若 Offering 開啟 `allow_buffer_override` 則顯示覆寫控制項。

### 4.3 管理員：人才搜尋頁 (Talent Discovery Layout)
- **結構**：
  1. **搜尋標頭 (Search Header)**：置頂大區塊。包含關鍵字搜尋框與「地區、大類」的大型圖示選取器。
  2. **動態過濾雲 (Tag Cloud)**：搜尋框下方展開的自動推薦 Hashtag 區域。
  3. **瀑布流/網格 (Result Grid)**：
     - 每列 3-4 張卡片。
     - 每張卡片：左側頭像，右側姓名與地區；下方為兩層技能 Chip 與 3 個熱門 Hashtag。
     - 卡片底部為「邀請加入」主要按鈕。
