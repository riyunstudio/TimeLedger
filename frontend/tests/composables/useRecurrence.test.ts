/**
 * 循環擴展單元測試
 *
 * 測試 useRecurrence composable 中的循環擴展邏輯
 */

import { describe, it, expect } from 'vitest'
import type { PersonalEvent } from '~/types'
import {
  expandRecurrenceEvent,
  expandRecurrenceEvents,
  getNextRecurrenceDate,
  isValidRecurrenceRule,
  getRecurrenceFrequencyLabel,
  getRecurrenceLabel,
} from '~/composables/useRecurrence'

describe('useRecurrence', () => {
  // ==================== expandRecurrenceEvent 測試 ====================

  describe('expandRecurrenceEvent', () => {
    it('should return single event when no recurrence rule', () => {
      const result = expandRecurrenceEvent({
        eventId: 123,
        title: '單次會議',
        startAt: '2026-01-06T14:00:00+08:00',
        endAt: '2026-01-06T15:00:00+08:00',
        rangeStart: new Date('2026-01-06'),
        rangeEnd: new Date('2026-01-12'),
      })

      expect(result).toHaveLength(1)
      expect(result[0].id).toBe('123')
      expect(result[0].originalId).toBe(123)
      expect(result[0].title).toBe('單次會議')
    })

    it('should expand daily recurrence correctly', () => {
      const result = expandRecurrenceEvent({
        eventId: 456,
        title: '每日瑜珈',
        startAt: '2026-01-06T09:00:00+08:00',
        endAt: '2026-01-06T10:00:00+08:00',
        recurrenceRule: { frequency: 'DAILY', interval: 1 },
        rangeStart: new Date('2026-01-06'),
        rangeEnd: new Date('2026-01-10'),
      })

      // rangeEnd 為 00:00:00，所以只會到 1/9，共 4 天
      expect(result).toHaveLength(4) // 6, 7, 8, 9
      expect(result[0].id).toBe('456_2026-01-06')
      expect(result[1].id).toBe('456_2026-01-07')
      expect(result[3].id).toBe('456_2026-01-09')
    })

    it('should expand weekly recurrence correctly', () => {
      const result = expandRecurrenceEvent({
        eventId: 789,
        title: '每週會議',
        startAt: '2026-01-06T14:00:00+08:00',
        endAt: '2026-01-06T15:00:00+08:00',
        recurrenceRule: { frequency: 'WEEKLY', interval: 1 },
        rangeStart: new Date('2026-01-01'),
        rangeEnd: new Date('2026-01-31'),
      })

      // rangeEnd 為 00:00:00，所以 1/27 是最後一天，共 4 天
      expect(result).toHaveLength(4) // 1/6, 1/13, 1/20, 1/27
      expect(result[0].id).toBe('789_2026-01-06')
      expect(result[1].id).toBe('789_2026-01-13')
      expect(result[2].id).toBe('789_2026-01-20')
      expect(result[3].id).toBe('789_2026-01-27')
    })

    it('should expand biweekly recurrence correctly', () => {
      const result = expandRecurrenceEvent({
        eventId: 111,
        title: '雙週課程',
        startAt: '2026-01-06T10:00:00+08:00',
        endAt: '2026-01-06T12:00:00+08:00',
        recurrenceRule: { frequency: 'BIWEEKLY', interval: 1 },
        rangeStart: new Date('2026-01-01'),
        rangeEnd: new Date('2026-01-31'),
      })

      // rangeEnd 為 00:00:00，所以 1/20 是最後一天，共 2 天
      expect(result).toHaveLength(2) // 1/6, 1/20
      expect(result[0].id).toBe('111_2026-01-06')
      expect(result[1].id).toBe('111_2026-01-20')
    })

    it('should expand monthly recurrence correctly', () => {
      const result = expandRecurrenceEvent({
        eventId: 222,
        title: '月度報告',
        startAt: '2026-01-15T09:00:00+08:00',
        endAt: '2026-01-15T10:00:00+08:00',
        recurrenceRule: { frequency: 'MONTHLY', interval: 1 },
        rangeStart: new Date('2026-01-01'),
        rangeEnd: new Date('2026-03-31'),
      })

      expect(result).toHaveLength(3) // 1/15, 2/15, 3/15
      expect(result[0].id).toBe('222_2026-01-15')
      expect(result[1].id).toBe('222_2026-02-15')
      expect(result[2].id).toBe('222_2026-03-15')
    })

    it('should respect interval greater than 1', () => {
      const result = expandRecurrenceEvent({
        eventId: 333,
        title: '每 3 天健身',
        startAt: '2026-01-01T08:00:00+08:00',
        endAt: '2026-01-01T09:00:00+08:00',
        recurrenceRule: { frequency: 'DAILY', interval: 3 },
        rangeStart: new Date('2026-01-01'),
        rangeEnd: new Date('2026-01-10'),
      })

      expect(result).toHaveLength(4) // 1/1, 1/4, 1/7, 1/10
      expect(result[0].id).toBe('333_2026-01-01')
      expect(result[1].id).toBe('333_2026-01-04')
      expect(result[2].id).toBe('333_2026-01-07')
      expect(result[3].id).toBe('333_2026-01-10')
    })

    it('should handle range that starts after event start date', () => {
      const result = expandRecurrenceEvent({
        eventId: 444,
        title: '每週例會',
        startAt: '2026-01-01T15:00:00+08:00',
        endAt: '2026-01-01T16:00:00+08:00',
        recurrenceRule: { frequency: 'WEEKLY', interval: 1 },
        rangeStart: new Date('2026-01-15'),
        rangeEnd: new Date('2026-01-31'),
      })

      expect(result).toHaveLength(3) // 1/15, 1/22, 1/29
      expect(result[0].id).toBe('444_2026-01-15')
    })

    it('should preserve event duration when expanding', () => {
      const result = expandRecurrenceEvent({
        eventId: 555,
        title: '2 小時課程',
        startAt: '2026-01-06T10:00:00+08:00',
        endAt: '2026-01-06T12:00:00+08:00',
        recurrenceRule: { frequency: 'DAILY', interval: 1 },
        // 使用當天開始和隔天開始，確保當天被包含
        rangeStart: new Date('2026-01-06'),
        rangeEnd: new Date('2026-01-07'),
      })

      // rangeEnd 00:00:00 + interval = 01/07 00:00:00，超過 01/07 00:00:00
      // 所以只會有 01/06 一天
      expect(result).toHaveLength(1)
      // 10:00 - 12:00 should be 2 hours
      expect(result[0].start_at).toContain('T10:00:00')
      expect(result[0].end_at).toContain('T12:00:00')
    })

    it('should handle event with string ID', () => {
      const result = expandRecurrenceEvent({
        eventId: 'abc_123',
        title: '字串 ID 事件',
        startAt: '2026-01-06T14:00:00+08:00',
        endAt: '2026-01-06T15:00:00+08:00',
        recurrenceRule: { frequency: 'DAILY', interval: 1 },
        // 擴大範圍以確保事件被包含
        rangeStart: new Date('2026-01-06T00:00:00+08:00'),
        rangeEnd: new Date('2026-01-06T23:59:59+08:00'),
      })

      expect(result).toHaveLength(1)
      expect(result[0].id).toBe('abc_123_2026-01-06')
      expect(result[0].originalId).toBe('abc_123')
    })

    it('should preserve optional properties', () => {
      const result = expandRecurrenceEvent({
        eventId: 666,
        title: '全天活動',
        startAt: '2026-01-06T00:00:00+08:00',
        endAt: '2026-01-06T23:59:59+08:00',
        isAllDay: true,
        colorHex: '#FF5733',
        recurrenceRule: { frequency: 'WEEKLY', interval: 1 },
        // 使用包含時區的日期範圍
        rangeStart: new Date('2026-01-06T00:00:00+08:00'),
        rangeEnd: new Date('2026-01-12T23:59:59+08:00'),
      })

      // 1/6 (Tue), 下一週 1/13 超過 rangeEnd
      expect(result).toHaveLength(1) // 只有 1/6
      expect(result[0].is_all_day).toBe(true)
      expect(result[0].color_hex).toBe('#FF5733')
    })
  })

  // ==================== expandRecurrenceEvents 測試 ====================

  describe('expandRecurrenceEvents', () => {
    it('should expand multiple events correctly', () => {
      const events: PersonalEvent[] = [
        {
          id: 1,
          teacher_id: 1,
          title: '每日運動',
          start_at: '2026-01-06T07:00:00+08:00',
          end_at: '2026-01-06T08:00:00+08:00',
          color: '#FF0000',
          created_at: '',
          updated_at: '',
          recurrence_rule: {
            frequency: 'DAILY',
            interval: 1,
          },
        },
        {
          id: 2,
          teacher_id: 1,
          title: '每週會議',
          start_at: '2026-01-06T14:00:00+08:00',
          end_at: '2026-01-06T15:00:00+08:00',
          color: '#00FF00',
          created_at: '',
          updated_at: '',
          recurrence_rule: {
            frequency: 'WEEKLY',
            interval: 1,
          },
        },
      ]

      const result = expandRecurrenceEvents(
        events,
        new Date('2026-01-06'),
        new Date('2026-01-10')
      )

      // Daily: 6, 7, 8, 9, 10 = 5 instances
      // Weekly: 6 only (next week is 13) = 1 instance
      expect(result.length).toBeGreaterThanOrEqual(5)
    })

    it('should handle events without recurrence rule', () => {
      const events: PersonalEvent[] = [
        {
          id: 3,
          teacher_id: 1,
          title: '單次活動',
          start_at: '2026-01-06T10:00:00+08:00',
          end_at: '2026-01-06T11:00:00+08:00',
          color: '#0000FF',
          created_at: '',
          updated_at: '',
        },
      ]

      const result = expandRecurrenceEvents(
        events,
        new Date('2026-01-06'),
        new Date('2026-01-10')
      )

      expect(result).toHaveLength(1)
      expect(result[0].id).toBe('3')
    })

    it('should handle empty events array', () => {
      const result = expandRecurrenceEvents(
        [],
        new Date('2026-01-06'),
        new Date('2026-01-10')
      )

      expect(result).toHaveLength(0)
    })
  })

  // ==================== getNextRecurrenceDate 測試 ====================

  describe('getNextRecurrenceDate', () => {
    it('should return start date when it is in the future', () => {
      const futureDate = new Date()
      futureDate.setDate(futureDate.getDate() + 5)

      const result = getNextRecurrenceDate(
        futureDate.toISOString(),
        { frequency: 'DAILY', interval: 1 }
      )

      expect(result).not.toBeNull()
      expect(result!.getTime()).toBe(futureDate.getTime())
    })

    it('should find next occurrence for daily recurrence', () => {
      const startDate = new Date('2026-01-01T10:00:00+08:00')
      const afterDate = new Date('2026-01-05T10:00:00+08:00')

      const result = getNextRecurrenceDate(
        startDate.toISOString(),
        { frequency: 'DAILY', interval: 1 },
        afterDate
      )

      expect(result).not.toBeNull()
      expect(result!.getDate()).toBe(6) // 1/6
    })

    it('should find next occurrence for weekly recurrence', () => {
      const startDate = new Date('2026-01-06T10:00:00+08:00') // Tuesday
      const afterDate = new Date('2026-01-08T10:00:00+08:00') // Thursday

      const result = getNextRecurrenceDate(
        startDate.toISOString(),
        { frequency: 'WEEKLY', interval: 1 },
        afterDate
      )

      expect(result).not.toBeNull()
      expect(result!.getDate()).toBe(13) // Next Tuesday
    })

    it('should return null for unknown frequency', () => {
      const result = getNextRecurrenceDate(
        '2026-01-06T10:00:00+08:00',
        { frequency: 'UNKNOWN' as any, interval: 1 }
      )

      expect(result).toBeNull()
    })

    it('should respect interval for biweekly recurrence', () => {
      const startDate = new Date('2026-01-06T10:00:00+08:00')
      const afterDate = new Date('2026-01-14T10:00:00+08:00')

      const result = getNextRecurrenceDate(
        startDate.toISOString(),
        { frequency: 'BIWEEKLY', interval: 1 },
        afterDate
      )

      expect(result).not.toBeNull()
      expect(result!.getDate()).toBe(20) // Two weeks later
    })
  })

  // ==================== isValidRecurrenceRule 測試 ====================

  describe('isValidRecurrenceRule', () => {
    it('should return false for undefined rule', () => {
      expect(isValidRecurrenceRule(undefined)).toBe(false)
    })

    it('should return true for valid daily rule', () => {
      expect(isValidRecurrenceRule({ frequency: 'DAILY', interval: 1 })).toBe(true)
    })

    it('should return true for valid weekly rule', () => {
      expect(isValidRecurrenceRule({ frequency: 'WEEKLY', interval: 1 })).toBe(true)
    })

    it('should return true for valid biweekly rule', () => {
      expect(isValidRecurrenceRule({ frequency: 'BIWEEKLY', interval: 1 })).toBe(true)
    })

    it('should return true for valid monthly rule', () => {
      expect(isValidRecurrenceRule({ frequency: 'MONTHLY', interval: 1 })).toBe(true)
    })

    it('should return false for invalid frequency', () => {
      expect(isValidRecurrenceRule({ frequency: 'INVALID' as any, interval: 1 })).toBe(false)
    })

    it('should return false for interval less than 1', () => {
      expect(isValidRecurrenceRule({ frequency: 'DAILY', interval: 0 })).toBe(false)
      expect(isValidRecurrenceRule({ frequency: 'DAILY', interval: -1 })).toBe(false)
    })

    it('should return false for non-integer interval', () => {
      expect(isValidRecurrenceRule({ frequency: 'DAILY', interval: 1.5 })).toBe(false)
    })

    it('should return true when interval is not provided (default to 1)', () => {
      expect(isValidRecurrenceRule({ frequency: 'DAILY' })).toBe(true)
    })
  })

  // ==================== getRecurrenceFrequencyLabel 測試 ====================

  describe('getRecurrenceFrequencyLabel', () => {
    it('should return Chinese label for DAILY', () => {
      expect(getRecurrenceFrequencyLabel('DAILY')).toBe('每日')
    })

    it('should return Chinese label for WEEKLY', () => {
      expect(getRecurrenceFrequencyLabel('WEEKLY')).toBe('每週')
    })

    it('should return Chinese label for BIWEEKLY', () => {
      expect(getRecurrenceFrequencyLabel('BIWEEKLY')).toBe('每兩週')
    })

    it('should return Chinese label for MONTHLY', () => {
      expect(getRecurrenceFrequencyLabel('MONTHLY')).toBe('每月')
    })

    it('should return frequency as-is for unknown value', () => {
      expect(getRecurrenceFrequencyLabel('UNKNOWN' as any)).toBe('UNKNOWN')
    })
  })

  // ==================== getRecurrenceLabel 測試 ====================

  describe('getRecurrenceLabel', () => {
    it('should return empty string for undefined rule', () => {
      expect(getRecurrenceLabel(undefined)).toBe('')
    })

    it('should return frequency label for interval 1', () => {
      expect(getRecurrenceLabel({ frequency: 'WEEKLY', interval: 1 })).toBe('每週')
      expect(getRecurrenceLabel({ frequency: 'DAILY', interval: 1 })).toBe('每日')
      expect(getRecurrenceLabel({ frequency: 'MONTHLY', interval: 1 })).toBe('每月')
    })

    it('should include interval for interval greater than 1', () => {
      expect(getRecurrenceLabel({ frequency: 'WEEKLY', interval: 2 })).toBe('每隔 2 每週')
      expect(getRecurrenceLabel({ frequency: 'DAILY', interval: 3 })).toBe('每隔 3 每日')
    })

    it('should handle invalid rule gracefully', () => {
      expect(getRecurrenceLabel({ frequency: 'INVALID' as any, interval: 1 })).toBe('')
    })
  })
})
