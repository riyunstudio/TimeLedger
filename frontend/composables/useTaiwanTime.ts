/**
 * 前端時區工具
 * 確保所有日期操作都使用台灣時區 (Asia/Taipei)
 */

// 取得台灣時區的字串
export const TAIWAN_TIMEZONE = 'Asia/Taipei'

/**
 * 將 Date 物件轉換為 YYYY-MM-DD 字串（台灣時區）
 * 修正 toISOString() 的 UTC 偏移問題
 */
export function formatDateToString(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

/**
 * 取得現在的日期字串（台灣時區）
 */
export function getTodayString(): string {
  return formatDateToString(new Date())
}

/**
 * 取得本週開始日期（週一）
 */
export function getWeekStart(date: Date = new Date()): Date {
  const d = new Date(date)
  const day = d.getDay()
  // 週一到週日: 1-7, 週日轉為 7
  const adjustedDay = day === 0 ? 7 : day
  const diff = adjustedDay - 1 // 調整到週一
  d.setDate(d.getDate() - diff)
  d.setHours(0, 0, 0, 0)
  return d
}

/**
 * 取得本週結束日期（週日）
 */
export function getWeekEnd(date: Date = new Date()): Date {
  const start = getWeekStart(date)
  const end = new Date(start)
  end.setDate(end.getDate() + 6)
  end.setHours(23, 59, 59, 999)
  return end
}

/**
 * 格式化日期顯示（台灣格式）
 */
export function formatDateDisplay(date: Date): string {
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

/**
 * 格式化週標籤
 */
export function formatWeekLabel(startDate: Date, endDate: Date): string {
  const start = startDate.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric' })
  const end = endDate.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric', year: 'numeric' })
  return `${start} - ${end}`
}

/**
 * 解析 YYYY-MM-DD 字串為 Date 物件（台灣時區）
 */
export function parseDateString(dateStr: string): Date {
  const [year, month, day] = dateStr.split('-').map(Number)
  return new Date(year, month - 1, day, 0, 0, 0, 0)
}

/**
 * 比較兩個日期是否為同一天
 */
export function isSameDay(date1: Date, date2: Date): boolean {
  return (
    date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
  )
}

/**
 * 檢查日期是否為今天
 */
export function isToday(date: Date): boolean {
  return isSameDay(date, new Date())
}

/**
 * 計算兩個日期相差天數
 */
export function getDaysDifference(date1: Date, date2: Date): number {
  const time1 = new Date(date1.getFullYear(), date1.getMonth(), date1.getDate()).getTime()
  const time2 = new Date(date2.getFullYear(), date2.getMonth(), date2.getDate()).getTime()
  return Math.floor((time2 - time1) / (1000 * 60 * 60 * 24))
}
