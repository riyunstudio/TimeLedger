/**
 * 循環展開效能基準測試
 *
 * 測試 useRecurrence composable 的循環展開效能，
 * 確保在 500+ 項目時仍能維持在 10ms 以內
 */

import { describe, it, expect, beforeAll } from 'vitest'
import type { PersonalEvent } from '~/types'
import {
  expandRecurrenceEvent,
  expandRecurrenceEvents,
} from '~/composables/useRecurrence'

// ==================== 測試資料生成器 ====================

/**
 * 生成隨機循環事件測試資料
 */
function generateRecurrenceEvent(id: number): PersonalEvent {
  const startHour = 8 + Math.floor(Math.random() * 10) // 8:00 - 18:00
  const duration = 1 + Math.floor(Math.random() * 3) // 1-3 小時

  const startDate = new Date(2026, 0, 6 + Math.floor(Math.random() * 20))
  const startAt = `${startDate.toISOString().split('T')[0]}T${String(startHour).padStart(2, '0')}:00:00+08:00`

  const endDate = new Date(startDate)
  endDate.setHours(startHour + duration)
  const endAt = `${endDate.toISOString().split('T')[0]}T${String(endDate.getHours()).padStart(2, '0')}:00:00+08:00`

  const frequencies: Array<'DAILY' | 'WEEKLY' | 'BIWEEKLY' | 'MONTHLY'> = ['DAILY', 'WEEKLY', 'BIWEEKLY', 'MONTHLY']
  const frequency = frequencies[Math.floor(Math.random() * frequencies.length)]
  const interval = 1 + Math.floor(Math.random() * 2) // 1-2

  return {
    id,
    teacher_id: 1,
    title: `課程 ${id}`,
    start_at: startAt,
    end_at: endAt,
    color: '#FF5733',
    created_at: '',
    updated_at: '',
    recurrence_rule: {
      frequency,
      interval,
    },
  }
}

/**
 * 生成大量循環事件資料
 */
function generateManyEvents(count: number): PersonalEvent[] {
  return Array.from({ length: count }, (_, i) => generateRecurrenceEvent(i + 1))
}

// ==================== 基準測試資料 ====================

// 生成 500 個事件的測試資料（用於大規模測試）
const fiveHundredEvents = generateManyEvents(500)
// 生成 1000 個事件的測試資料（壓力測試）
const thousandEvents = generateManyEvents(1000)

// 測試用的時間範圍（一個月）
const rangeStart = new Date('2026-01-01T00:00:00+08:00')
const rangeEnd = new Date('2026-01-31T23:59:59+08:00')

// 效能閾值（毫秒）
const PERFORMANCE_THRESHOLD_MS = 10

describe('useRecurrence 效能基準測試', () => {
  // ==================== expandRecurrenceEvent 單事件展開 ====================

  describe('expandRecurrenceEvent 單事件展開', () => {
    it('展開每日循環事件（單一事件）', () => {
      const startTime = performance.now()
      expandRecurrenceEvent({
        eventId: 1,
        title: '每日瑜珈',
        startAt: '2026-01-06T09:00:00+08:00',
        endAt: '2026-01-06T10:00:00+08:00',
        recurrenceRule: { frequency: 'DAILY', interval: 1 },
        rangeStart,
        rangeEnd,
      })
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })

    it('展開每週循環事件（單一事件）', () => {
      const startTime = performance.now()
      expandRecurrenceEvent({
        eventId: 2,
        title: '每週會議',
        startAt: '2026-01-06T14:00:00+08:00',
        endAt: '2026-01-06T15:00:00+08:00',
        recurrenceRule: { frequency: 'WEEKLY', interval: 1 },
        rangeStart,
        rangeEnd,
      })
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })

    it('展開雙週循環事件（單一事件）', () => {
      const startTime = performance.now()
      expandRecurrenceEvent({
        eventId: 3,
        title: '雙週課程',
        startAt: '2026-01-06T10:00:00+08:00',
        endAt: '2026-01-06T12:00:00+08:00',
        recurrenceRule: { frequency: 'BIWEEKLY', interval: 1 },
        rangeStart,
        rangeEnd,
      })
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })

    it('展開每月循環事件（單一事件）', () => {
      const startTime = performance.now()
      expandRecurrenceEvent({
        eventId: 4,
        title: '月度報告',
        startAt: '2026-01-15T09:00:00+08:00',
        endAt: '2026-01-15T10:00:00+08:00',
        recurrenceRule: { frequency: 'MONTHLY', interval: 1 },
        rangeStart,
        rangeEnd,
      })
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })

    it('展開高頻率循環事件（每日，間隔 1 天，1 個月）', () => {
      const startTime = performance.now()
      expandRecurrenceEvent({
        eventId: 5,
        title: '每日高頻',
        startAt: '2026-01-01T08:00:00+08:00',
        endAt: '2026-01-01T09:00:00+08:00',
        recurrenceRule: { frequency: 'DAILY', interval: 1 },
        rangeStart: new Date('2026-01-01T00:00:00+08:00'),
        rangeEnd: new Date('2026-01-31T23:59:59+08:00'),
      })
      const duration = performance.now() - startTime
      // 高頻事件需要展開 31 個實例，閾值放寬
      expect(duration).toBeLessThan(5)
    })
  })

  // ==================== expandRecurrenceEvents 多事件展開 ====================

  describe('expandRecurrenceEvents 多事件展開（核心效能測試）', () => {
    it('展開 50 個循環事件', () => {
      const events = generateManyEvents(50)
      const startTime = performance.now()
      expandRecurrenceEvents(events, rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })

    it('展開 100 個循環事件', () => {
      const events = generateManyEvents(100)
      const startTime = performance.now()
      expandRecurrenceEvents(events, rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })

    it('展開 200 個循環事件', () => {
      const events = generateManyEvents(200)
      const startTime = performance.now()
      expandRecurrenceEvents(events, rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })

    it('展開 500 個循環事件（目標：<10ms）', () => {
      const startTime = performance.now()
      expandRecurrenceEvents(fiveHundredEvents, rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      console.log(`500 個循環事件展開時間: ${duration.toFixed(2)}ms`)
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })

    it('展開 1000 個循環事件（壓力測試）', () => {
      const startTime = performance.now()
      expandRecurrenceEvents(thousandEvents, rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      console.log(`1000 個循環事件展開時間: ${duration.toFixed(2)}ms`)
      // 壓力測試閾值放寬到 20ms
      expect(duration).toBeLessThan(20)
    })
  })

  // ==================== 大規模場景測試 ====================

  describe('大規模課表展開場景', () => {
    it('展開 500 個每日事件（模擬月檢視）', () => {
      const dailyEvents = Array.from({ length: 500 }, (_, i) => ({
        id: i + 1,
        teacher_id: 1,
        title: `每日課程 ${i + 1}`,
        start_at: `2026-01-06T0${8 + (i % 10)}:00:00+08:00`,
        end_at: `2026-01-06T0${9 + (i % 10)}:00:00+08:00`,
        color: '#FF5733',
        created_at: '',
        updated_at: '',
        recurrence_rule: {
          frequency: 'DAILY' as const,
          interval: 1,
        },
      }))
      const startTime = performance.now()
      expandRecurrenceEvents(dailyEvents, rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      console.log(`500 個每日事件展開時間: ${duration.toFixed(2)}ms`)
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })

    it('展開混合頻率事件（各類型 125 個，共 500 個）', () => {
      const mixedEvents: PersonalEvent[] = []

      // 每日事件
      for (let i = 0; i < 125; i++) {
        mixedEvents.push({
          id: i + 1,
          teacher_id: 1,
          title: `每日課程 ${i + 1}`,
          start_at: '2026-01-06T09:00:00+08:00',
          end_at: '2026-01-06T10:00:00+08:00',
          color: '#FF5733',
          created_at: '',
          updated_at: '',
          recurrence_rule: { frequency: 'DAILY', interval: 1 },
        })
      }

      // 每週事件
      for (let i = 0; i < 125; i++) {
        mixedEvents.push({
          id: i + 126,
          teacher_id: 1,
          title: `每週課程 ${i + 1}`,
          start_at: '2026-01-06T14:00:00+08:00',
          end_at: '2026-01-06T15:00:00+08:00',
          color: '#00FF00',
          created_at: '',
          updated_at: '',
          recurrence_rule: { frequency: 'WEEKLY', interval: 1 },
        })
      }

      // 雙週事件
      for (let i = 0; i < 125; i++) {
        mixedEvents.push({
          id: i + 251,
          teacher_id: 1,
          title: `雙週課程 ${i + 1}`,
          start_at: '2026-01-06T10:00:00+08:00',
          end_at: '2026-01-06T12:00:00+08:00',
          color: '#0000FF',
          created_at: '',
          updated_at: '',
          recurrence_rule: { frequency: 'BIWEEKLY', interval: 1 },
        })
      }

      // 每月事件
      for (let i = 0; i < 125; i++) {
        mixedEvents.push({
          id: i + 376,
          teacher_id: 1,
          title: `每月課程 ${i + 1}`,
          start_at: '2026-01-15T09:00:00+08:00',
          end_at: '2026-01-15T10:00:00+08:00',
          color: '#FF00FF',
          created_at: '',
          updated_at: '',
          recurrence_rule: { frequency: 'MONTHLY', interval: 1 },
        })
      }

      const startTime = performance.now()
      expandRecurrenceEvents(mixedEvents, rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      console.log(`500 個混合頻率事件展開時間: ${duration.toFixed(2)}ms`)
      expect(duration).toBeLessThan(PERFORMANCE_THRESHOLD_MS)
    })
  })

  // ==================== 邊界條件測試 ====================

  describe('邊界條件效能測試', () => {
    it('展開空陣列', () => {
      const startTime = performance.now()
      expandRecurrenceEvents([], rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(1)
    })

    it('展開無循環規則的單一事件', () => {
      const events: PersonalEvent[] = [{
        id: 1,
        teacher_id: 1,
        title: '單次活動',
        start_at: '2026-01-06T10:00:00+08:00',
        end_at: '2026-01-06T11:00:00+08:00',
        color: '#0000FF',
        created_at: '',
        updated_at: '',
      }]
      const startTime = performance.now()
      expandRecurrenceEvents(events, rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(1)
    })

    it('展開大間隔循環事件（每月，共 3 個月）', () => {
      const events: PersonalEvent[] = [{
        id: 1,
        teacher_id: 1,
        title: '月度課程',
        start_at: '2026-01-15T09:00:00+08:00',
        end_at: '2026-01-15T10:00:00+08:00',
        color: '#FF5733',
        created_at: '',
        updated_at: '',
        recurrence_rule: { frequency: 'MONTHLY', interval: 1 },
      }]
      const startTime = performance.now()
      expandRecurrenceEvents(events, rangeStart, rangeEnd)
      const duration = performance.now() - startTime
      expect(duration).toBeLessThan(1)
    })
  })
})

// ==================== 斷言驗證（確保功能正確性）====================

describe('useRecurrence 正確性驗證', () => {
  it('展開 500 個事件應該產生結果', () => {
    const result = expandRecurrenceEvents(fiveHundredEvents, rangeStart, rangeEnd)
    expect(result.length).toBeGreaterThan(0)
  })

  it('展開結果應該包含正確的 ID 格式', () => {
    const result = expandRecurrenceEvents(fiveHundredEvents, rangeStart, rangeEnd)
    // 檢查第一個結果的 ID 格式
    if (result.length > 0) {
      expect(result[0].id).toMatch(/^\d+_\d{4}-\d{2}-\d{2}$/)
    }
  })

  it('展開結果應該保留原始事件的所有屬性', () => {
    const events: PersonalEvent[] = [{
      id: 1,
      teacher_id: 1,
      title: '測試課程',
      start_at: '2026-01-06T09:00:00+08:00',
      end_at: '2026-01-06T10:00:00+08:00',
      color: '#FF5733',
      created_at: '',
      updated_at: '',
      recurrence_rule: { frequency: 'DAILY', interval: 1 },
    }]

    const result = expandRecurrenceEvents(events, rangeStart, rangeEnd)

    expect(result.length).toBeGreaterThan(0)
    expect(result[0].title).toBe('測試課程')
    expect(result[0].originalId).toBe(1)
  })
})
