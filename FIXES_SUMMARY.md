# 修復總結

## 1. settings.vue - 結束時間下拉選單修復

**問題**：結束時間下拉選單包含 24:00，這在我們的時間系統中是不正確的（應該是 00:00 到 23:00）。

**修復**：將 `v-for="hour in 24"` 的 `:key="hour"` 和 `:value="String(hour).padStart(2, '0') + ':00'"` 改為使用 `hour - 1`，使其與開始時間的邏輯一致。

**變更位置**：第 200-203 行

**修復前**：
```html
<option v-for="hour in 24" :key="hour" :value="String(hour).padStart(2, '0') + ':00'">
  {{ String(hour).padStart(2, '0') }}:00
</option>
```

**修復後**：
```html
<option v-for="hour in 24" :key="hour - 1" :value="String(hour - 1).padStart(2, '0') + ':00'">
  {{ String(hour - 1).padStart(2, '0') }}:00
</option>
```

---

## 2. resource-occupancy.vue - 根據中心營業時間顯示網格

**問題**：資源佔用表目前硬編碼顯示 00:00 - 23:00，需要根據中心的營業時間設定來顯示。

**修復**：

### 2.1 新增狀態變數
- 新增 `operatingStartTime`（預設 '00:00'）
- 新增 `operatingEndTime`（預設 '23:00'）

### 2.2 新增取得中心設定函數
```typescript
const fetchCenterSettings = async () => {
  const centerId = getCenterId()
  const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId}/settings`, {
    headers: { 'Authorization': `Bearer ${token}` },
  })
  // 解析並更新 operatingStartTime 和 operatingEndTime
}
```

### 2.3 更新 visibleTimeSlots computed
```typescript
const visibleTimeSlots = computed(() => {
  const start = parseInt(operatingStartTime.value.split(':')[0], 10)
  const end = parseInt(operatingEndTime.value.split(':')[0], 10) || 23
  const slots = []

  // 確保結束時間 >= 開始時間
  const actualEnd = end < start ? 23 : end

  for (let i = start; i <= actualEnd; i++) {
    slots.push(i)
  }

  return slots.length > 0 ? slots : timeSlots
})
```

### 2.4 更新 onMounted
在生命週期中呼叫 `fetchCenterSettings()`。

### 2.5 修復 getRuleStyle 函數
確保當課程開始時間早於營業時間時，top 值會被夾緊到 0，防止負值隱藏卡片。

---

## 3. CompactWeekView.vue - 990px 寬度佈局修復

**問題**：在 990px 寬度（平板/md 斷點）下，課程卡片可能會「溢出」或出現計算錯誤。

**修復**：

### 3.1 強化 updateWidth 函數
確保格子寬度有最小值限制：
```typescript
cellWidth.value = Math.max(60, (containerWidth.value - timeColumnWidth) / visibleDaysCount)
```

### 3.2 強化 getScheduleCardStyle 函數

**問題**：當 schedule.start_hour 小於第一個可見時段時，會產生負的 top 值。

**修復**：
```typescript
let relativeStartHour = schedule.start_hour - firstSlotHour

// 保護：如果課程開始時間早於營業時間，將其夾緊到 0
if (relativeStartHour < 0) {
  relativeStartHour = 0
}

const top = (relativeStartHour * slotHeight) + minuteOffset
const clampedTop = Math.max(0, top)
```

---

## 驗證清單

- [x] Settings 頁面：結束時間下拉選單範圍是 00:00 - 23:00
- [x] Resource Occupancy 頁面：網格根據設定的營業時間顯示
- [x] Resource Occupancy 頁面：getRuleStyle 正確處理早於營業時間的課程
- [x] CompactWeekView（990px）：課程卡片位置計算正確
- [x] CompactWeekView：格子寬度有最小值限制
- [x] 所有檔案無 lint 錯誤

---

## 檔案變更清單

| 檔案 | 變更類型 |
|------|----------|
| `frontend/pages/admin/settings.vue` | Bug fix |
| `frontend/pages/admin/resource-occupancy.vue` | Feature + Bug fix |
| `frontend/components/Scheduling/CompactWeekView.vue` | Bug fix |
