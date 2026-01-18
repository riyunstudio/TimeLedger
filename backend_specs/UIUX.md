# UIUX.md — TimeLedger（RWD Web / MVP）
日期：2026-01-18

> 原則：MVP 只做「中心 ↔ 老師連結」與「時間治理」。不做九宮格拖曳。

---

## 1. 視覺風格
- 乾淨、留白、低學習成本
- 顏色：狀態色固定
  - NORMAL：中性
  - PENDING：灰
  - CANCELLED：警示
  - RESCHEDULED：提醒
- 元件一致性：Primary/Secondary/Danger

---

## 2. Teacher 端（Mobile First）
### 2.1 我的課表（/teacher/schedule）
- 預設週視圖 list
- Tab：中心課表 / 我的行程
- FAB：新增個人行程（modal）
- 行程卡：標題 + 時間 + 來源 tag（Center/Personal）+ 狀態

### 2.2 新增/編輯個人行程（modal）
- start/end datetime pickers
- title（必填）note（選填）
- 儲存前呼叫 validate（僅 personal overlap）
- 衝突提示：ConflictsList

---

## 3. Admin 端（Desktop first / RWD）
### 3.1 Center 列表（/admin/centers）
- 顯示 plan、試用剩餘天數 banner
- CTA：進入管理

### 3.2 老師管理（/admin/centers/:id/teachers）
- 清單：teacher 名稱、狀態(invited/active)、最後登入（選配）
- 邀請：輸入 teacher_id（MVP）
- 啟用：activate

### 3.3 Offering（/admin/centers/:id/offerings）
- 列表 + 新增/編輯（name/duration/status）

### 3.4 Rules（/admin/centers/:id/rules）
- 表格：weekday/start/end/teacher/offering
- 新增/更新前先 call validate
- 顯示衝突清單（ConflictsList）

### 3.5 Exceptions（/admin/centers/:id/exceptions）
- 篩選：Pending/Approved/Rejected
- 卡片：type、原日期、（新日期）、reason
- 操作：Approve/Reject
- approve 前若 RESCHEDULE，validate 失敗需阻擋

### 3.6 Schedule View（/admin/centers/:id/schedule）
- from/to 選擇
- teacher filter
- list 顯示 sessions（含狀態）

---

## 4. 共用 UX Pattern
- ConflictsList：三段式（哪裡衝突/為何/怎麼解）
- PlanBanner：試用剩餘天數 + 到期唯讀提示
- Toast：成功/失敗
