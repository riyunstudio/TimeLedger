/**
 * 循環擴展 Composable
 *
 * 提供循環規則（Recurrence Rule）的展開功能，
 * 將循環事件展開為指定時間範圍內的具體實例。
 */

import type { PersonalEvent, RecurrenceRule } from '~/types'
import { formatDateToString } from './useTaiwanTime'

/**
 * 循環頻率類型
 */
export type RecurrenceFrequency = 'DAILY' | 'WEEKLY' | 'BIWEEKLY' | 'MONTHLY'

/**
 * 循環規則介面
 */
export interface RecurrenceRuleInput {
  type: RecurrenceFrequency
  interval?: number
}

/**
 * 循環事件實例介面
 */
export interface RecurrenceEventInstance {
  id: string
  originalId: number | string
  title: string
  start_at: string
  end_at: string
  is_all_day?: boolean
  color_hex?: string
  recurrence_rule?: RecurrenceRule
  [key: string]: unknown
}

/**
 * 循環展開選項
 */
export interface ExpandRecurrenceOptions {
  /** 事件 ID（數字或字串） */
  eventId: number | string
  /** 事件標題 */
  title: string
  /** 事件開始時間 */
  startAt: string
  /** 事件結束時間 */
  endAt: string
  /** 是否為全天事件 */
  isAllDay?: boolean
  /** 事件顏色 */
  colorHex?: string
  /** 循環規則 */
  recurrenceRule?: RecurrenceRuleInput
  /** 目標範圍開始時間 */
  rangeStart: Date
  /** 目標範圍結束時間 */
  rangeEnd: Date
}

/**
 * 格式化日期時間為 API 格式（台灣時區）
 * @param date - Date 物件
 * @returns API 格式的字串（如：2026-01-15T10:30:00+08:00）
 */
const formatDateTimeForApi = (date: Date): string => {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}+08:00`
}

/**
 * 展開單個循環事件為指定範圍內的所有實例
 *
 * @param options - 循環展開選項
 * @returns 循環事件實例陣列
 *
 * @example
 * ```typescript
 * const instances = expandRecurrenceEvent({
 *   eventId: 123,
 *   title: '每週會議',
 *   startAt: '2026-01-06T14:00:00+08:00',
 *   endAt: '2026-01-06T15:00:00+08:00',
 *   recurrenceRule: { type: 'WEEKLY', interval: 1 },
 *   rangeStart: new Date('2026-01-06'),
 *   rangeEnd: new Date('2026-01-20')
 * })
 * ```
 */
export function expandRecurrenceEvent(options: ExpandRecurrenceOptions): RecurrenceEventInstance[] {
  const {
    eventId,
    title,
    startAt,
    endAt,
    isAllDay = false,
    colorHex,
    recurrenceRule,
    rangeStart,
    rangeEnd,
  } = options

  // 如果沒有循環規則，直接返回單一事件
  if (!recurrenceRule) {
    return [{
      id: String(eventId),
      originalId: eventId,
      title,
      start_at: startAt,
      end_at: endAt,
      is_all_day: isAllDay,
      color_hex: colorHex,
    }]
  }

  const { type, interval = 1 } = recurrenceRule
  const startDate = new Date(startAt)
  const endDate = new Date(endAt)
  const duration = endDate.getTime() - startDate.getTime()

  const instances: RecurrenceEventInstance[] = []
  let currentDate = new Date(startDate)

  while (currentDate <= rangeEnd) {
    // 檢查當前日期是否在目標範圍內
    if (currentDate >= rangeStart) {
      const instanceEnd = new Date(currentDate.getTime() + duration)
      const instanceId = `${eventId}_${formatDateToString(currentDate)}`

      instances.push({
        id: instanceId,
        originalId: eventId,
        title,
        start_at: formatDateTimeForApi(currentDate),
        end_at: formatDateTimeForApi(instanceEnd),
        is_all_day: isAllDay,
        color_hex: colorHex,
      })
    }

    // 根據循環類型前進到下一個日期
    switch (type) {
      case 'DAILY':
        currentDate.setDate(currentDate.getDate() + interval)
        break
      case 'WEEKLY':
        currentDate.setDate(currentDate.getDate() + (7 * interval))
        break
      case 'BIWEEKLY':
        currentDate.setDate(currentDate.getDate() + (14 * interval))
        break
      case 'MONTHLY':
        currentDate.setMonth(currentDate.getMonth() + interval)
        break
      default:
        // 未知的循環類型，提前結束迴圈
        currentDate = new Date(rangeEnd.getTime() + 1)
    }
  }

  return instances
}

/**
 * 展開多個循環事件為指定範圍內的所有實例
 *
 * @param events - 包含循環規則的事件陣列
 * @param rangeStart - 目標範圍開始時間
 * @param rangeEnd - 目標範圍結束時間
 * @returns 展開後的事件實例陣列
 *
 * @example
 * ```typescript
 * const events: PersonalEvent[] = [...]
 * const expanded = expandRecurrenceEvents(events, weekStart, weekEnd)
 * ```
 */
export function expandRecurrenceEvents(
  events: PersonalEvent[],
  rangeStart: Date,
  rangeEnd: Date
): RecurrenceEventInstance[] {
  const expandedEvents: RecurrenceEventInstance[] = []

  events.forEach(event => {
    const instances = expandRecurrenceEvent({
      eventId: event.id,
      title: event.title,
      startAt: event.start_at,
      endAt: event.end_at,
      isAllDay: event.is_all_day,
      colorHex: event.color_hex,
      recurrenceRule: event.recurrence_rule
        ? {
            type: event.recurrence_rule.type as RecurrenceFrequency,
            interval: event.recurrence_rule.interval,
          }
        : undefined,
      rangeStart,
      rangeEnd,
    })
    expandedEvents.push(...instances)
  })

  return expandedEvents
}

/**
 * 計算循環事件的下一個發生日期
 *
 * @param startAt - 事件開始時間
 * @param recurrenceRule - 循環規則
 * @param afterDate - 參考日期（預設為現在）
 * @returns 下一個發生日期，如果沒有則返回 null
 *
 * @example
 * ```typescript
 * const nextDate = getNextRecurrenceDate(
 *   '2026-01-06T14:00:00+08:00',
 *   { type: 'WEEKLY', interval: 1 },
 *   new Date()
 * )
 * ```
 */
export function getNextRecurrenceDate(
  startAt: string,
  recurrenceRule: RecurrenceRuleInput,
  afterDate: Date = new Date()
): Date | null {
  const { type, interval = 1 } = recurrenceRule
  const startDate = new Date(startAt)
  const duration = 0 // 只是為了計算下一次發生

  // 如果開始日期在參考日期之前，需要找到下一個
  if (startDate <= afterDate) {
    let currentDate = new Date(startDate)

    while (currentDate <= afterDate) {
      switch (type) {
        case 'DAILY':
          currentDate.setDate(currentDate.getDate() + interval)
          break
        case 'WEEKLY':
          currentDate.setDate(currentDate.getDate() + (7 * interval))
          break
        case 'BIWEEKLY':
          currentDate.setDate(currentDate.getDate() + (14 * interval))
          break
        case 'MONTHLY':
          currentDate.setMonth(currentDate.getMonth() + interval)
          break
        default:
          return null
      }
    }

    return currentDate
  }

  return startDate
}

/**
 * 檢查循環規則是否有效
 *
 * @param rule - 循環規則
 * @returns 是否有效
 */
export function isValidRecurrenceRule(rule: RecurrenceRuleInput | undefined): boolean {
  if (!rule) return false

  const validTypes: RecurrenceFrequency[] = ['DAILY', 'WEEKLY', 'BIWEEKLY', 'MONTHLY']

  if (!validTypes.includes(rule.type)) {
    return false
  }

  if (rule.interval !== undefined && (rule.interval < 1 || !Number.isInteger(rule.interval))) {
    return false
  }

  return true
}

/**
 * 獲取循環頻率的中文名稱
 *
 * @param frequency - 循環頻率
 * @returns 中文名稱
 */
export function getRecurrenceFrequencyLabel(frequency: RecurrenceFrequency): string {
  const labels: Record<RecurrenceFrequency, string> = {
    DAILY: '每日',
    WEEKLY: '每週',
    BIWEEKLY: '每兩週',
    MONTHLY: '每月',
  }

  return labels[frequency] || frequency
}

/**
 * 組合函數：取得循環事件的顯示標籤
 *
 * @param recurrenceRule - 循環規則
 * @returns 格式化的標籤字串（如：每週、每隔 2 週）
 */
export function getRecurrenceLabel(recurrenceRule: RecurrenceRuleInput | undefined): string {
  if (!recurrenceRule || !isValidRecurrenceRule(recurrenceRule)) {
    return ''
  }

  const { type, interval = 1 } = recurrenceRule
  const frequencyLabel = getRecurrenceFrequencyLabel(frequency)

  if (interval === 1) {
    return frequencyLabel
  }

  return `每隔 ${interval} ${frequencyLabel}`
}
