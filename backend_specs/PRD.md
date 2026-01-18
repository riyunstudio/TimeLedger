# PRD｜TimeLedger（中心 ↔ 老師連結版）
版本：v2.0
日期：2026-01-18
狀態：可進入 MVP 開發

---

## 1. 產品概述

### 1.1 目的
TimeLedger 是一套 **「中心 ↔ 老師」課表治理系統**：
- 老師（免費）先管理自己的時間（包含非中心課程）
- 中心（付費/試用）透過老師同意的連結，引用老師時間做排課
- 系統不處理學生報名/點名/金流，只專注「排課衝突防呆 + 例外變動治理」

### 1.2 Out of Scope（本期不做）
- 學生端報名/點名/繳費/金流
- App（iOS/Android）
- LMS（教材/作業）
- 自動排課演算法
- Room/Teacher buffer（MVP 暫不做；僅做 overlap）

### 1.3 目標使用者
- Center Admin（中心管理員/行政）：排課、邀請/啟用老師、審核停課/改期
- Teacher（老師）：維護個人忙碌時間；查看中心排課；提出改期/停課（若中心啟用）

---

## 2. 需求假設與規模
- 初期：單機部署可支援（~50 users / 同時 <10）
- 行為：查課表為主；排課與例外為低頻
- 多中心：center 間資料隔離（Hard Rule）
- 同一位老師可加入多中心，但中心僅能看到「該中心授權範圍」

---

## 3. 名詞定義
- Center：中心/租戶（資料隔離）
- Teacher：老師（全域帳號；可加入多中心）
- Membership：老師加入中心的關係（invited/active）
- Offering：中心的授課單位（課程/班別/場地等）
- Rule：固定排課規則（weekday + start/end + teacher + offering）
- Session：查詢時計算出的結果（Rule 投影後套用例外）
- Exception：停課/改期（PENDING/APPROVED/REJECTED）
- Plan/Trial：中心方案與試用（僅 gating，不做金流）

---

## 4. 權限模型（MVP）
- Teacher：只能操作自己的 personal sessions；只能看自己相關 schedule view
- Admin：只能操作自己 center 的 offerings/rules/exceptions/memberships
- Center scoping：所有 Admin API 必須綁 center_id，跨中心必須回 31001/403

---

## 5. 核心功能（MVP）

### 5.1 Teacher（免費）
- 個人課表（Personal Sessions）CRUD
- 查詢自己的 schedule view（包含中心排課 + 自己行程）
- 基本衝突防呆（自己行程 overlap）

### 5.2 Center Admin（試用/付費）
- Center 管理（可先單中心）
- 邀請/啟用老師（Membership）
- Offering CRUD
- Rule CRUD（排課前可 validate）
- Validate API：檢查候選時段與「老師個人行程 + 中心既有 rules 投影 + 例外」的衝突
- Exceptions：CANCEL/RESCHEDULE 建立、待審列表、approve/reject
- 方案與試用：30 天試用 + 到期降級唯讀；限制老師數（Starter/Pro/Team）
- Audit log：所有寫入記錄 who/what/when

---

## 6. 防呆機制（MVP）
- Hard Conflict：同 teacher 同時間不得 overlap（包含 personal 與 center rules）
- validate 必須同時做：前端 UX 提示 + 後端強制檢查
- 覆寫/緩衝（buffer）本期不做

---

## 7. 例外規格（MVP）
- CANCEL：取消某日投影 session
- RESCHEDULE：原日取消 + 新日時段新增
- 狀態：PENDING/APPROVED/REJECTED
- 唯一性：同一 (center_id, rule_id, original_date) 同時只能有一筆有效例外（PENDING/APPROVED）
- approve 時必須重新 validate 新時段（RESCHEDULE）

---

## 8. NFR
- 安全：token 過期；中心隔離不可被繞過
- 可維運：錯誤碼一致（見 docs/ERR_CODES.md）；audit/event 可追蹤
- 可用性：單機 + MySQL；每日備份（保留 7~14 天）

---

## 9. 里程碑（建議）
- M1：MySQL schema + Teacher 個人課表 + Center/Membership/Offering
- M2：Rules + Validate + Schedule view
- M3：Exceptions + 審核 + Audit/Event + Plan gating + RWD 前端
