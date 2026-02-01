/**
 * useScheduleStore 單元測試
 * 測試重點：Pure functions、座標映射、狀態管理
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import type { CenterMembership, WeekSchedule } from '../../types'

// Mock 依賴
vi.mock('../../composables/useTaiwanTime', () => ({
  formatDateToString: vi.fn((date: Date) => {
    const pad = (n: number) => String(n).padStart(2, '0')
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
  })
}))

vi.mock('../../composables/useRecurrence', () => ({
  expandRecurrenceEvents: vi.fn((events, start, end) => events)
}))

vi.mock('../../utils/loadingHelper', () => ({
  withLoading: vi.fn((loadingState, fn) => fn())
}))

vi.mock('../../composables/useApi', () => ({
  useApi: vi.fn(() => ({
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn()
  }))
}))

// ============================================
// Helper Functions 測試
// ============================================

describe('useScheduleStore - Helper Functions', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('getWeekStart', () => {
    it('應該正確計算週一為起點', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      // 週三測試
      const wedJan15 = new Date(2025, 0, 15) // 2025-01-15 是週三
      const weekStart = store.getWeekStart(wedJan15)

      expect(weekStart.getDay()).toBe(1) // 應該是週一
      expect(weekStart.getDate()).toBe(13) // 2025-01-13 是週一
    })

    it('週一的起始應該是自己', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      const monJan13 = new Date(2025, 0, 13) // 2025-01-13 是週一
      const weekStart = store.getWeekStart(monJan13)

      expect(weekStart.getDay()).toBe(1)
      expect(weekStart.getDate()).toBe(13)
    })

    it('週日的起始應該是上週一', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      const sunJan19 = new Date(2025, 0, 19) // 2025-01-19 是週日
      const weekStart = store.getWeekStart(sunJan19)

      expect(weekStart.getDay()).toBe(1)
      expect(weekStart.getDate()).toBe(13) // 上週一
    })

    it('應該正確處理跨月和跨年的日期', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      // 2024-12-30 是週一
      const dec302024 = new Date(2024, 11, 30)
      const weekStart = store.getWeekStart(dec302024)

      expect(weekStart.getFullYear()).toBe(2024)
      expect(weekStart.getMonth()).toBe(11) // December
      expect(weekStart.getDate()).toBe(30)
    })
  })

  describe('changeWeek', () => {
    it('應該能夠前進一週', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      // 設定初始週
      const initialDate = new Date(2025, 0, 13) // 週一
      store.weekStart = initialDate

      // 前進一週
      store.changeWeek(1)

      expect(store.weekStart.getDate()).toBe(20) // 2025-01-20
      expect(store.weekStart.getDay()).toBe(1)
    })

    it('應該能夠後退一週', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      const initialDate = new Date(2025, 0, 13)
      store.weekStart = initialDate

      store.changeWeek(-1)

      expect(store.weekStart.getDate()).toBe(6) // 2025-01-06
    })

    it('應該能夠前進多週', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      const initialDate = new Date(2025, 0, 13)
      store.weekStart = initialDate

      store.changeWeek(2)

      expect(store.weekStart.getDate()).toBe(27)
    })

    it('應該能夠後退多週', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      const initialDate = new Date(2025, 0, 13)
      store.weekStart = initialDate

      store.changeWeek(-3)

      expect(store.weekStart.getDate()).toBe(23) // 2024-12-23
    })

    it('應該在 weekStart 為 null 時不執行任何操作', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      store.weekStart = null as any

      // 不應該拋出錯誤
      expect(() => store.changeWeek(1)).not.toThrow()
    })
  })

  describe('weekEnd Computed', () => {
    it('應該正確計算週結束日期', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      store.weekStart = new Date(2025, 0, 13) // 週一
      const weekEnd = store.weekEnd

      expect(weekEnd).not.toBeNull()
      expect(weekEnd!.getDate()).toBe(19) // 週日
      expect(weekEnd!.getDay()).toBe(0)
    })

    it('應該在 weekStart 為 null 時返回 null', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      store.weekStart = null as any
      const weekEnd = store.weekEnd

      expect(weekEnd).toBeNull()
    })
  })

  describe('weekLabel Computed', () => {
    it('應該正確生成週標籤', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      store.weekStart = new Date(2025, 0, 13)

      const label = store.weekLabel

      expect(label).toContain('1月13日')
      expect(label).toContain('1月19日')
      expect(label).toContain(' - ')
    })

    it('應該在資料不完整時返回空字串', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      store.weekStart = null as any

      expect(store.weekLabel).toBe('')
    })
  })

  describe('formatDateTimeForApi', () => {
    it('應該正確格式化為 API 所需的 datetime 格式', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      const date = new Date(2025, 0, 15, 14, 30, 45)
      const formatted = store.formatDateTimeForApi(date)

      expect(formatted).toBe('2025-01-15T14:30:45+08:00')
    })

    it('應該正確處理零填充', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      const date = new Date(2025, 0, 5, 9, 5, 5)
      const formatted = store.formatDateTimeForApi(date)

      expect(formatted).toBe('2025-01-05T09:05:05+08:00')
    })

    it('應該正確處理單數月份和日期', () => {
      const { useScheduleStore } = require('../../stores/useScheduleStore')
      const store = useScheduleStore()

      const date = new Date(2025, 5, 3, 8, 3, 0)
      const formatted = store.formatDateTimeForApi(date)

      expect(formatted).toBe('2025-06-03T08:03:00+08:00')
    })
  })
})

// ============================================
// transformToWeekSchedule 測試
// ============================================

describe('useScheduleStore - transformToWeekSchedule', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('應該正確將課表資料轉換為週格式', () => {
    const { useScheduleStore } = require('../../stores/useScheduleStore')
    const store = useScheduleStore()

    // 設定週開始日期
    store.weekStart = new Date(2025, 0, 13)

    const inputItems = [
      {
        id: '1',
        type: 'SCHEDULE_RULE',
        title: '瑜伽課程',
        date: '2025-01-13',
        start_time: '09:00:00',
        end_time: '10:00:00',
        room_id: 1,
        center_id: 1,
        center_name: '台北中心',
        status: 'ACTIVE',
        is_cross_day_part: false
      },
      {
        id: '2',
        type: 'CENTER_SESSION',
        title: '舞蹈課程',
        date: '2025-01-14',
        start_time: '14:00:00',
        end_time: '15:30:00',
        room_id: 2,
        center_id: 1,
        center_name: '台北中心',
        status: 'ACTIVE',
        is_cross_day_part: false
      }
    ] as any

    const result = store.transformToWeekSchedule(inputItems)

    expect(result.week_start).toBe('2025-01-13')
    expect(result.week_end).toBe('2025-01-19')
    expect(result.days.length).toBe(7)
    expect(result.days[0].date).toBe('2025-01-13')
    expect(result.days[0].day_of_week).toBe(1)
    expect(result.days[1].date).toBe('2025-01-14')
    expect(result.days[1].day_of_week).toBe(2)
  })

  it('應該正確處理空課表', () => {
    const { useScheduleStore } = require('../../stores/useScheduleStore')
    const store = useScheduleStore()

    store.weekStart = new Date(2025, 0, 13)

    const result = store.transformToWeekSchedule([])

    expect(result.days.length).toBe(7)
    expect(result.days[0].items.length).toBe(0)
  })

  it('應該正確處理跨天課程', () => {
    const { useScheduleStore } = require('../../stores/useScheduleStore')
    const store = useScheduleStore()

    store.weekStart = new Date(2025, 0, 13)

    const inputItems = [
      {
        id: '1',
        type: 'SCHEDULE_RULE',
        title: '長期課程',
        date: '2025-01-13',
        start_time: '22:00:00',
        end_time: '02:00:00',
        room_id: 1,
        center_id: 1,
        center_name: '台北中心',
        status: 'ACTIVE',
        is_cross_day_part: true
      }
    ] as any

    const result = store.transformToWeekSchedule(inputItems)

    expect(result.days[0].items.length).toBe(1)
    expect(result.days[0].items[0].is_cross_day_part).toBe(true)
  })

  it('應該正確處理 PENDING_CANCEL 狀態的顏色', () => {
    const { useScheduleStore } = require('../../stores/useScheduleStore')
    const store = useScheduleStore()

    store.weekStart = new Date(2025, 0, 13)

    const inputItems = [
      {
        id: '1',
        type: 'SCHEDULE_RULE',
        title: '待取消課程',
        date: '2025-01-13',
        start_time: '10:00:00',
        end_time: '11:00:00',
        room_id: 1,
        center_id: 1,
        center_name: '台北中心',
        status: 'PENDING_CANCEL',
        is_cross_day_part: false
      }
    ] as any

    const result = store.transformToWeekSchedule(inputItems)

    expect(result.days[0].items[0].color).toBe('#F59E0B')
  })

  it('應該正確處理所有類型的課表項目', () => {
    const { useScheduleStore } = require('../../stores/useScheduleStore')
    const store = useScheduleStore()

    store.weekStart = new Date(2025, 0, 13)

    const inputItems = [
      {
        id: '1',
        type: 'SCHEDULE_RULE',
        title: '規則課程',
        date: '2025-01-13',
        start_time: '09:00:00',
        end_time: '10:00:00',
        room_id: 1,
        center_id: 1,
        center_name: '中心A',
        status: 'ACTIVE',
        is_cross_day_part: false
      },
      {
        id: '2',
        type: 'PERSONAL_EVENT',
        title: '私人行程',
        date: '2025-01-14',
        start_time: '11:00:00',
        end_time: '12:00:00',
        room_id: 0,
        center_id: 0,
        center_name: undefined,
        status: 'ACTIVE',
        is_cross_day_part: false
      },
      {
        id: '3',
        type: 'CENTER_SESSION',
        title: '中心課程',
        date: '2025-01-15',
        start_time: '14:00:00',
        end_time: '15:00:00',
        room_id: 2,
        center_id: 1,
        center_name: '中心B',
        status: 'ACTIVE',
        is_cross_day_part: false
      }
    ] as any

    const result = store.transformToWeekSchedule(inputItems)

    // 檢查類型轉換
    expect(result.days[0].items[0].type).toBe('SCHEDULE_RULE')
    expect(result.days[1].items[0].type).toBe('PERSONAL_EVENT')
    expect(result.days[2].items[0].type).toBe('CENTER_SESSION')
  })
})

// ============================================
// Loading States 測試
// ============================================

describe('useScheduleStore - Loading States', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('應該有正確的初始 loading 狀態', () => {
    const { useScheduleStore } = require('../../stores/useScheduleStore')
    const store = useScheduleStore()

    expect(store.isLoading).toBe(false)
    expect(store.isFetching).toBe(false)
    expect(store.isCreating).toBe(false)
    expect(store.isUpdating).toBe(false)
    expect(store.isDeleting).toBe(false)
  })

  it('應該有正確的初始事件 loading 狀態', () => {
    const { useScheduleStore } = require('../../stores/useScheduleStore')
    const store = useScheduleStore()

    expect(store.isCreatingEvent).toBe(false)
    expect(store.isUpdatingEvent).toBe(false)
    expect(store.isDeletingEvent).toBe(false)
    expect(store.isCreatingException).toBe(false)
    expect(store.isRevokingException).toBe(false)
    expect(store.isSavingNote).toBe(false)
    expect(store.isRespondingInvitation).toBe(false)
  })
})

// ============================================
// Schedule Data 測試
// ============================================

describe('useScheduleStore - Schedule Data', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('應該有正確的初始資料狀態', () => {
    const { useScheduleStore } = require('../../stores/useScheduleStore')
    const store = useScheduleStore()

    expect(store.centers).toEqual([])
    expect(store.currentCenter).toBeNull()
    expect(store.schedule).toBeNull()
    expect(store.exceptions).toEqual([])
    expect(store.personalEvents).toEqual([])
    expect(store.sessionNote).toBeNull()
    expect(store.invitations).toEqual([])
    expect(store.pendingInvitationsCount).toBe(0)
  })
})
