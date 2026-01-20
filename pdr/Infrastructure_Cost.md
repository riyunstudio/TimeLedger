# 基礎架構與成本評估 (Infrastructure & Cost Analysis)

根據您的需求，本專案將採用 **單體架構 (Monolithic Architecture)** 並部署於 **VPS (Virtual Private Server)**。此方案在早中期能將成本壓至最低，且維護與備份最為單純。

---

## 1. 架構規劃 (Monolithic Architecture)

### 1.1 技術堆疊
*   **Application**: Node.js (Backend) + React/Vue (Frontend SPA) 打包在同一個 Docker Image，或透過 Nginx 反向代理 Host。
*   **Database**: MySQL 8.0 (Docker Container)。
*   **Cache/Lock**: Redis (Docker Container) - 用於 Session Store 與 排課併發鎖 (Distributed Lock)。
*   **Orchestration**: Docker Compose。

### 1.2 部署拓樸 (Topology)
*   **單機部署 (Single VPS)**: 所有服務 (App, MySQL, Redis, Nginx) 跑在同一台機器上。
*   **備份策略**: 每日 Cron Job 匯出 MySQL Dump -> 上傳至 S3/R2 (外部冷儲存)。

---

## 2. 機器成本預估 (Cost Estimation)

我們以 **DigitalOcean / Linode / Vultr** 等標準 VPS 廠商報價估算。

### 2.1 初期/輕量方案 (MVP ~ 50 users)
適合開發期與小規模試營運。
*   **規格**: 1 CPU / 1GB RAM / 25GB SSD
*   **費用**: **$5 ~ $6 USD / 月** (約 NT$ 160 ~ 200)
    *   VPS: $6
    *   S3 Backup (AWS/Cloudflare R2): 免費額度或 < $1
*   **限制**: 1GB RAM 跑 MySQL + Node.js 會比較緊繃，建議設定 Swap Memory (虛擬記憶體) 2GB 以防 OOM (Out of Memory)。

### 2.2 中型/正式營運方案 (Stable ~ 500 users)
適合正式商用，確保排課效能流暢。
*   **規格**: 2 CPU / 4GB RAM / 80GB SSD
*   **費用**: **$20 ~ $24 USD / 月** (約 NT$ 650 ~ 780)
    *   VPS: $24
    *   Domain 網域費: $12 / 年 (平均 $1/月)
*   **優勢**: 4GB RAM 綽綽有餘，可承受較高的並發排課請求。

### 2.3 成長型方案 (High Traffic)
若未來擴張至多據點、數千名老師。
*   **規格**: 4 CPU / 8GB RAM (或者將 DB 拆分到獨立機器)
*   **費用**: **$48 USD / 月**
*   **策略**: 此時建議將 MySQL 移至 Managed Database (由廠商代管)，應用層仍留在 VPS，成本會上升但穩定性更高。

---

## 3. 備份與隱性成本

除了機器月租費，還需考慮：
1.  **S3/Object Storage (圖片儲存)**:
    *   老師證照與匯出圖片。
    *   Cloudflare R2 (推薦): 10GB 內免費，無流量費。
    *   成本估算: **$0** (初期)。
2.  **Domain (網域)**:
    *   `.com` / `.app` / `.tw`: 約 **$10 ~ $25 USD / 年**。
3.  **維護成本 (Man-hour)**:
    *   VPS 需要自己設定 Firewall (UFW), SSL (Certbot), OS Update。
    *   若無專職 DevOps，建議撰寫自動化 Script (Setup Script)。

## 4. 總結建議

**「每月預算 NT$ 200 (開發期) -> NT$ 800 (正式營運)」** 即可搞定。

相較於 Cloud PaaS (Vercel Pro $20 + RDS $30 + Redis $10 = $60+)，VPS 單體架構可節省約 **60% ~ 90%** 的基礎建設費用，非常適合本專案的屬性。

---

## 5. 商業收費策略 (Pricing Strategy - Taiwan Market)

基於台灣補教/才藝市場的特性，我們採取 **「低門檻排課、高價值獵才」** 的定價策略。

| 方案 (Tier) | 目標客群 | 費用 (NTD) | 市場對標與價值 |
|:---|:---|:---|:---|
| **Free (老師版)** | 個人接案老師 | **$0 (永久免費)** | **養魚策略**。讓老師習慣用此平台管理私人與多中心行程，累積流量。 |
| **Starter (工作室)** | 個人工作室/家教 | **$499 /月** | **比訂便當還便宜**。用極低價解決 Google Calendar 無法「防撞」與「審核」的痛點。 |
| **Growth (成長型)** | 社區型才藝教室 | **$1,299 /月** | **省下半個行政人員薪水**。含「智慧代課」，解決週五晚上臨時找不到人的痛苦。 |
| **Pro (獵才版)** | 連鎖/大型機構 | **$2,490 /月** | **比 104 還划算**。104 刊登費動輒數千，且難找才藝專才。這裡能直接搜尋「有空檔 + 有證照」的即戰力。 |

*   **台灣市場優勢**:
    *   **LINE 黏著度**: 善用台灣人離不開 LINE 的特性，強調「排課通知直接推播到 LINE」，這是國外軟體 (Calendly) 做不到的。
    *   **證照迷思**: 台灣家長看重證照 (Yamaha, 英國皇家)。我們的「證照驗證」與「公開履歷」精準打中此需求。

---

## 6. 財務損益預估 (Conservative Financial Projection)

基於 **2026 年** 時空背景 (LINE Notify 已於 2025/3 停止服務)，我們必須採用 **Messaging API (Push Message)** 計算成本，每則約 NT$ 0.2。

### 6.1 變動成本分析 (每位老師)
假設每位老師每月有 20 堂課，需發送 20 則「明日提醒」+ 5 則「異動通知」。
*   **單人月用量**: 25 則訊息。
*   **單人月成本**: 25 * 0.2 = **NT$ 5 /月**。
*   **結論**: 即便全額負擔通知費，佔不到月費 ($499) 的 1%。**毛利極高**。

### 6.2 保守收入模型 (Year 1 目標)
設定目標：**100 間付費中心** (約 300~500 位活躍老師)。

| 項目 | 計算方式 | 金額 (NTD) |
|:---|:---|:---:|
| **(+) 月營收** | 50間 Starter ($499) + 40間 Growth ($1299) + 10間 Pro ($2490) | **$101,810 /月** |
| **(-) 機器成本** | DigitalOcean VPS ($24 USD) | **$780 /月** |
| **(-) LINE 通知費** | 500師 * $5 | **$2,500 /月** |
| **(=) 淨利 (Net Profit)** | | **$98,530 /月** |

*   **損益兩平點 (Break-even)**: 只要有 **2 間** Starter 中心付費，即可覆蓋所有伺服器成本。
*   **結論**: 本專案屬 **「高毛利 微型 SaaS」**，只要累積 100 個客戶，即可創造每月 10 萬被動現金流。機器擴充成本極低，適合長期持有。

---

## 7. 效能與容量預估 (Performance & Capacity)

針對「多少機器能扛多少人」的保守評估。我們的瓶頸主要在 **資料庫排課鎖 (DB Lock)**，而非 CPU。

### 7.1 單機容量對照表 (Monolith VPS)

| 機器規格 | 適用階段 | 並發上線 (CCU) | 日活躍 (DAU) | 估計能容納中心數 |
|:---|:---|:---:|:---:|:---:|
| **Starter ($6/mo)**<br>1 vCPU / 1GB RAM | 開發期 / 試營運 | ~50 人 | ~5,000 人 | 約 50 間 |
| **Pro ($24/mo)**<br>2 vCPU / 4GB RAM | 正式營運 (Year 1-3) | **~300 人** | **~30,000 人** | **約 300 間** |
| **Scale-out ($80+/mo)**<br>DB 分離 + 讀寫分離 | 大型擴張期 | 1000+ 人 | 10萬+ 人 | 1000 間以上 |

### 7.2 評估依據
*   **Go (Gin)**: 在 2 vCPU 下處理純讀取 API (如看課表) 可達 5,000+ QPS。
*   **MySQL (Write Lock)**:
    *   排課寫入 (Validation + Lock) 耗時約 50ms。
    *   理論上限：單一中心每秒約可處理 20 筆排課。
    *   但在真實世界，300 間中心是 **平行運作** (互不鎖定)，因此系統總吞吐量極高。
*   **結論**:
    *   一台 **US$24 (NT$ 780)** 的機器，足以支撐 **300 間中心、3 萬名老師** 的日常運作。
    *   這意味著我們在達到「月營收 300 萬」之前，根本不需要煩惱架構升級的問題。
