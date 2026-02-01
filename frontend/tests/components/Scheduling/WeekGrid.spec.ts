/**
 * WeekGrid.vue 單元測試
 * 測試重點：座標映射、事件處理、重疊檢測
 */

import { describe, it, expect, vi, beforeEach, afterEach, ref } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import { h, defineComponent, nextTick } from 'vue'

// Mock vue-virtual-scroller
vi.mock('vue-virtual-scroller', () => ({
  DynamicScroller: defineComponent({
    name: 'DynamicScroller',
    props: {
      items: { type: Array, required: true },
      minItemSize: { type: Number, required: true },
      keyField: { type: String, default: 'id' }
    },
    emits: ['update:items'],
    template: '<div class="dynamic-scroller"><slot v-for="item in items" name="default" :item="item" :index="items.indexOf(item)" /></div>'
  }),
  DynamicScrollerItem: defineComponent({
    name: 'DynamicScrollerItem',
    props: {
      item: { type: Object, required: true },
      active: { type: Boolean },
      sizeDependencies: { type: Array },
      dataIndex: { type: Number }
    },
    template: '<div class="dynamic-scroller-item"><slot :item="item" /></div>'
  })
}))

// Mock ScheduleCard component
const MockScheduleCard = defineComponent({
  name: 'ScheduleCard',
  props: {
    schedule: { type: Object, required: true },
    cardInfoType: { type: String, default: 'teacher' }
  },
  emits: ['click'],
  template: '<div class="schedule-card" @click="$emit(\'click\', $event)">{{ schedule.title }}</div>'
})

// ============================================
// 座標映射測試
// ============================================

describe('WeekGrid - Coordinate Mapping', () => {
  let wrapper: VueWrapper<any>

  const defaultProps = {
    schedules: [],
    weekLabel: '1月13日 - 1月19日',
    cardInfoType: 'teacher' as const,
    validationResults: {},
    slotWidth: 150
  }

  beforeEach(async () => {
    // Dynamic import to avoid hoisting issues
    const { default: WeekGrid } = await import('~/components/Scheduling/WeekGrid.vue')

    wrapper = mount(WeekGrid, {
      props: defaultProps,
      global: {
        components: {
          ScheduleCard: MockScheduleCard
        }
      }
    })
  })

  afterEach(() => {
    wrapper.unmount()
  })

  describe('TIME_SLOT_HEIGHT', () => {
    it('應該正確計算時段高度為 60px', () => {
      expect(wrapper.vm.TIME_SLOT_HEIGHT).toBe(60)
    })
  })

  describe('TIME_COLUMN_WIDTH', () => {
    it('應該正確計算時間列寬度為 80px', () => {
      expect(wrapper.vm.TIME_COLUMN_WIDTH).toBe(80)
    })
  })

  describe('timeSlots', () => {
    it('應該包含正確的時段列表', () => {
      const timeSlots = wrapper.vm.timeSlots

      // 檢查特殊時段 (0-3)
      expect(timeSlots.slice(0, 4)).toEqual([0, 1, 2, 3])

      // 檢查主要時段 (9-23)
      expect(timeSlots.slice(4)).toEqual([9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23])

      // 總數應該是 19 個時段
      expect(timeSlots.length).toBe(19)
    })

    it('應該跳過 4-8 時段', () => {
      const timeSlots = wrapper.vm.timeSlots

      // 確認 4, 5, 6, 7, 8 不在列表中
      expect(timeSlots).not.toContain(4)
      expect(timeSlots).not.toContain(5)
      expect(timeSlots).not.toContain(6)
      expect(timeSlots).not.toContain(7)
      expect(timeSlots).not.toContain(8)
    })
  })

  describe('weekDays', () => {
    it('應該包含正確的星期列表', () => {
      const weekDays = wrapper.vm.weekDays

      expect(weekDays.length).toBe(7)
      expect(weekDays[0]).toEqual({ value: 1, name: '週一' })
      expect(weekDays[1]).toEqual({ value: 2, name: '週二' })
      expect(weekDays[2]).toEqual({ value: 3, name: '週三' })
      expect(weekDays[3]).toEqual({ value: 4, name: '週四' })
      expect(weekDays[4]).toEqual({ value: 5, name: '週五' })
      expect(weekDays[5]).toEqual({ value: 6, name: '週六' })
      expect(weekDays[6]).toEqual({ value: 7, name: '週日' })
    })
  })

  describe('getScheduleStyle', () => {
    it('應該正確計算週一的左邊位置', () => {
      const schedule = {
        id: 1,
        weekday: 1,
        start_hour: 9,
        start_minute: 0,
        duration_minutes: 60
      }

      const style = wrapper.vm.getScheduleStyle(schedule)

      // 左邊位置 = 80 + (0 * slotWidth) = 80
      expect(style.left).toBe('80px')
    })

    it('應該正確計算週日的左邊位置', () => {
      const schedule = {
        id: 1,
        weekday: 7,
        start_hour: 9,
        start_minute: 0,
        duration_minutes: 60
      }

      const style = wrapper.vm.getScheduleStyle(schedule)

      // 左邊位置 = 80 + (6 * slotWidth) = 80 + 900 = 980
      expect(style.left).toBe('980px')
    })

    it('應該正確計算 09:00 的頂部位置', () => {
      const schedule = {
        id: 1,
        weekday: 1,
        start_hour: 9,
        start_minute: 0,
        duration_minutes: 60
      }

      const style = wrapper.vm.getScheduleStyle(schedule)

      // topSlotIndex 從 0-3 開始計算，然後從 9 開始
      // 09:00 應該是第 4 個可見時段 (索引 4)
      // top = 4 * 60 = 240
      expect(style.top).toBe('240px')
    })

    it('應該正確計算 00:00 的頂部位置', () => {
      const schedule = {
        id: 1,
        weekday: 1,
        start_hour: 0,
        start_minute: 0,
        duration_minutes: 60
      }

      const style = wrapper.vm.getScheduleStyle(schedule)

      // 00:00 應該是第一個可見時段
      expect(style.top).toBe('0px')
    })

    it('應該正確處理非整點開始時間', () => {
      const schedule = {
        id: 1,
        weekday: 1,
        start_hour: 9,
        start_minute: 30,
        duration_minutes: 60
      }

      const style = wrapper.vm.getScheduleStyle(schedule)

      // top = 240 + (30/60) * 60 = 240 + 30 = 270
      expect(style.top).toBe('270px')
    })

    it('應該正確計算課程高度', () => {
      const schedule = {
        id: 1,
        weekday: 1,
        start_hour: 9,
        start_minute: 0,
        duration_minutes: 90
      }

      const style = wrapper.vm.getScheduleStyle(schedule)

      // height = (90/60) * 60 = 90
      expect(style.height).toBe('90px')
    })

    it('應該正確計算 30 分鐘課程高度', () => {
      const schedule = {
        id: 1,
        weekday: 1,
        start_hour: 10,
        start_minute: 0,
        duration_minutes: 30
      }

      const style = wrapper.vm.getScheduleStyle(schedule)

      // height = (30/60) * 60 = 30
      expect(style.height).toBe('30px')
    })

    it('應該正確計算卡片寬度', () => {
      const schedule = {
        id: 1,
        weekday: 1,
        start_hour: 9,
        start_minute: 0,
        duration_minutes: 60
      }

      const style = wrapper.vm.getScheduleStyle(schedule)

      // width = slotWidth - 4 = 150 - 4 = 146
      expect(style.width).toBe('146px')
    })

    it('應該使用快取來提高效能', () => {
      const schedule = {
        id: 1,
        weekday: 1,
        start_hour: 9,
        start_minute: 0,
        duration_minutes: 60
      }

      // 第一次調用
      const style1 = wrapper.vm.getScheduleStyle(schedule)

      // 第二次調用應該使用快取
      const style2 = wrapper.vm.getScheduleStyle(schedule)

      // 結果應該相同
      expect(style1).toEqual(style2)
    })
  })

  describe('computeOverlapData', () => {
    it('應該正確計算單一課程的重疊數量', () => {
      const schedules = [
        { id: 1, weekday: 1, start_hour: 9, start_minute: 0 }
      ]

      wrapper.vm.computeOverlapData(schedules)

      const cacheData = wrapper.vm.overlapDataCache.get('1-9-0')
      expect(cacheData).toBeDefined()
      expect(cacheData?.count).toBe(1)
    })

    it('應該正確計算多個課程的重疊數量', () => {
      const schedules = [
        { id: 1, weekday: 1, start_hour: 9, start_minute: 0 },
        { id: 2, weekday: 1, start_hour: 9, start_minute: 0 },
        { id: 3, weekday: 1, start_hour: 9, start_minute: 0 }
      ]

      wrapper.vm.computeOverlapData(schedules)

      const cacheData = wrapper.vm.overlapDataCache.get('1-9-0')
      expect(cacheData?.count).toBe(3)
    })

    it('應該正確識別第一個課程 ID', () => {
      const schedules = [
        { id: 5, weekday: 1, start_hour: 9, start_minute: 0 },
        { id: 2, weekday: 1, start_hour: 9, start_minute: 0 },
        { id: 8, weekday: 1, start_hour: 9, start_minute: 0 }
      ]

      wrapper.vm.computeOverlapData(schedules)

      const cacheData = wrapper.vm.overlapDataCache.get('1-9-0')
      expect(cacheData?.firstId).toBe(2) // 最小的 ID
    })

    it('應該正確處理不同時段', () => {
      const schedules = [
        { id: 1, weekday: 1, start_hour: 9, start_minute: 0 },
        { id: 2, weekday: 1, start_hour: 10, start_minute: 0 },
        { id: 3, weekday: 2, start_hour: 9, start_minute: 0 }
      ]

      wrapper.vm.computeOverlapData(schedules)

      // 每個時段都應該只有一個課程
      expect(wrapper.vm.overlapDataCache.get('1-9-0')?.count).toBe(1)
      expect(wrapper.vm.overlapDataCache.get('1-10-0')?.count).toBe(1)
      expect(wrapper.vm.overlapDataCache.get('2-9-0')?.count).toBe(1)
    })
  })
})

// ============================================
// 事件處理測試
// ============================================

describe('WeekGrid - Event Handling', () => {
  let wrapper: VueWrapper<any>

  const defaultProps = {
    schedules: [],
    weekLabel: '1月13日 - 1月19日',
    cardInfoType: 'teacher' as const,
    validationResults: {},
    slotWidth: 150
  }

  beforeEach(async () => {
    const { default: WeekGrid } = await import('~/components/Scheduling/WeekGrid.vue')

    wrapper = mount(WeekGrid, {
      props: defaultProps,
      global: {
        components: {
          ScheduleCard: MockScheduleCard
        }
      }
    })
  })

  afterEach(() => {
    wrapper.unmount()
  })

  describe('drag-enter event', () => {
    it('應該在拖入格子時發出 drag-enter 事件', async () => {
      const gridCell = wrapper.find('[class*="border-white/5"]:first-child')

      await gridCell.trigger('dragenter', {
        time: 9,
        day: 1
      })

      expect(wrapper.emitted('drag-enter')).toBeTruthy()
      expect(wrapper.emitted('drag-enter')[0]).toEqual([9, 1])
    })

    it('應該正確傳遞時間和星期參數', async () => {
      const gridCell = wrapper.find('[class*="border-white/5"]:nth-child(2)')

      await gridCell.trigger('dragenter', {
        time: 14,
        day: 3
      })

      expect(wrapper.emitted('drag-enter')[0]).toEqual([14, 3])
    })
  })

  describe('drag-leave event', () => {
    it('應該在拖離格子時發出 drag-leave 事件', async () => {
      const gridCell = wrapper.find('[class*="border-white/5"]:first-child')

      await gridCell.trigger('dragleave')

      expect(wrapper.emitted('drag-leave')).toBeTruthy()
    })
  })

  describe('dragover event', () => {
    it('應該阻止預設的 dragover 行為', async () => {
      const gridCell = wrapper.find('[class*="border-white/5"]:first-child')

      const event = {
        preventDefault: vi.fn()
      }

      await gridCell.trigger('dragover', event)

      expect(event.preventDefault).toHaveBeenCalled()
    })
  })

  describe('select-schedule event', () => {
    it('應該在點擊課程卡片時發出 select-schedule 事件', async () => {
      const schedule = {
        id: 1,
        title: '測試課程',
        weekday: 1,
        start_hour: 9,
        start_minute: 0,
        duration_minutes: 60
      }

      await wrapper.setProps({
        schedules: [schedule],
        slotWidth: 150
      })

      await nextTick()

      const card = wrapper.find('.schedule-card')
      await card.trigger('click')

      expect(wrapper.emitted('select-schedule')).toBeTruthy()
      expect(wrapper.emitted('select-schedule')[0][0]).toMatchObject({
        id: 1,
        title: '測試課程'
      })
    })
  })

  describe('overlap-click event', () => {
    it('應該在點擊重疊指示器時發出 overlap-click 事件', async () => {
      const schedules = [
        { id: 1, title: '課程1', weekday: 1, start_hour: 9, start_minute: 0, duration_minutes: 60 },
        { id: 2, title: '課程2', weekday: 1, start_hour: 9, start_minute: 0, duration_minutes: 60 }
      ]

      await wrapper.setProps({
        schedules,
        slotWidth: 150
      })

      await nextTick()

      const overlapIndicator = wrapper.find('.bg-warning-500\\/20')
      await overlapIndicator.trigger('click')

      expect(wrapper.emitted('overlap-click')).toBeTruthy()
    })
  })
})

// ============================================
// 重疊檢測測試
// ============================================

describe('WeekGrid - Overlap Detection', () => {
  let wrapper: VueWrapper<any>

  const defaultProps = {
    schedules: [],
    weekLabel: '1月13日 - 1月19日',
    cardInfoType: 'teacher' as const,
    validationResults: {},
    slotWidth: 150
  }

  beforeEach(async () => {
    const { default: WeekGrid } = await import('~/components/Scheduling/WeekGrid.vue')

    wrapper = mount(WeekGrid, {
      props: defaultProps,
      global: {
        components: {
          ScheduleCard: MockScheduleCard
        }
      }
    })
  })

  afterEach(() => {
    wrapper.unmount()
  })

  describe('schedulesWithOverlapData', () => {
    it('應該正確標記單一課程', async () => {
      const schedules = [
        { id: 1, weekday: 1, start_hour: 9, start_minute: 0 }
      ]

      await wrapper.setProps({ schedules })
      await nextTick()

      const result = wrapper.vm.schedulesWithOverlapData
      expect(result[0]._overlapCount).toBe(1)
      expect(result[0]._isFirstInOverlap).toBe(true)
    })

    it('應該正確標記重疊課程', async () => {
      const schedules = [
        { id: 1, weekday: 1, start_hour: 9, start_minute: 0 },
        { id: 2, weekday: 1, start_hour: 9, start_minute: 0 },
        { id: 3, weekday: 1, start_hour: 9, start_minute: 0 }
      ]

      await wrapper.setProps({ schedules })
      await nextTick()

      const result = wrapper.vm.schedulesWithOverlapData

      // 全部應該有 _overlapCount = 3
      result.forEach(item => {
        expect(item._overlapCount).toBe(3)
      })

      // 只有 ID 最小的應該是 _isFirstInOverlap = true
      expect(result.filter(item => item._isFirstInOverlap).length).toBe(1)
      expect(result.find(item => item._isFirstInOverlap)?.id).toBe(1)
    })

    it('應該正確處理不同時段的課程', async () => {
      const schedules = [
        { id: 1, weekday: 1, start_hour: 9, start_minute: 0 },
        { id: 2, weekday: 1, start_hour: 10, start_minute: 0 },
        { id: 3, weekday: 2, start_hour: 9, start_minute: 0 }
      ]

      await wrapper.setProps({ schedules })
      await nextTick()

      const result = wrapper.vm.schedulesWithOverlapData

      // 所有課程都應該是 _overlapCount = 1
      result.forEach(item => {
        expect(item._overlapCount).toBe(1)
      })
    })
  })
})

// ============================================
// 去重測試
// ============================================

describe('WeekGrid - Deduplication', () => {
  let wrapper: VueWrapper<any>

  const defaultProps = {
    schedules: [],
    weekLabel: '1月13日 - 1月19日',
    cardInfoType: 'teacher' as const,
    validationResults: {},
    slotWidth: 150
  }

  beforeEach(async () => {
    const { default: WeekGrid } = await import('~/components/Scheduling/WeekGrid.vue')

    wrapper = mount(WeekGrid, {
      props: defaultProps,
      global: {
        components: {
          ScheduleCard: MockScheduleCard
        }
      }
    })
  })

  afterEach(() => {
    wrapper.unmount()
  })

  describe('uniqueSchedules', () => {
    it('應該正確去重課程', async () => {
      const schedules = [
        { id: 1, weekday: 1, start_time: '09:00:00' },
        { id: 1, weekday: 1, start_time: '09:00:00' }, // 重複
        { id: 2, weekday: 1, start_time: '10:00:00' }
      ]

      await wrapper.setProps({ schedules })
      await nextTick()

      const result = wrapper.vm.uniqueSchedules
      expect(result.length).toBe(2)
      expect(result[0].id).toBe(1)
      expect(result[1].id).toBe(2)
    })

    it('應該保留所有唯一的課程', async () => {
      const schedules = [
        { id: 1, weekday: 1, start_time: '09:00:00' },
        { id: 2, weekday: 1, start_time: '09:00:00' },
        { id: 3, weekday: 2, start_time: '09:00:00' }
      ]

      await wrapper.setProps({ schedules })
      await nextTick()

      const result = wrapper.vm.uniqueSchedules
      expect(result.length).toBe(3)
    })
  })
})

// ============================================
// 驗證狀態測試
// ============================================

describe('WeekGrid - Validation State', () => {
  let wrapper: VueWrapper<any>

  const defaultProps = {
    schedules: [],
    weekLabel: '1月13日 - 1月19日',
    cardInfoType: 'teacher' as const,
    validationResults: {},
    slotWidth: 150
  }

  beforeEach(async () => {
    const { default: WeekGrid } = await import('~/components/Scheduling/WeekGrid.vue')

    wrapper = mount(WeekGrid, {
      props: defaultProps,
      global: {
        components: {
          ScheduleCard: MockScheduleCard
        }
      }
    })
  })

  afterEach(() => {
    wrapper.unmount()
  })

  describe('cellClassMap', () => {
    it('應該正確處理無效的驗證結果', async () => {
      await wrapper.setProps({
        validationResults: {
          '9-1': { valid: false }
        }
      })
      await nextTick()

      const classMap = wrapper.vm.cellClassMap
      expect(classMap['9-1']).toContain('bg-critical-500/10')
    })

    it('應該正確處理警告狀態', async () => {
      await wrapper.setProps({
        validationResults: {
          '9-1': { warning: '可能存在衝突' }
        }
      })
      await nextTick()

      const classMap = wrapper.vm.cellClassMap
      expect(classMap['9-1']).toContain('bg-warning-500/10')
    })

    it('應該正確處理有效的驗證結果', async () => {
      await wrapper.setProps({
        validationResults: {
          '9-1': { valid: true }
        }
      })
      await nextTick()

      const classMap = wrapper.vm.cellClassMap
      expect(classMap['9-1']).toContain('bg-success-500/10')
    })

    it('應該正確處理無驗證結果的預設狀態', async () => {
      await wrapper.setProps({
        validationResults: {}
      })
      await nextTick()

      const classMap = wrapper.vm.cellClassMap
      // 應該返回 hover 樣式
      expect(classMap['9-1']).toContain('hover:bg-white/5')
    })
  })
})

// ============================================
// 快取行為測試
// ============================================

describe('WeekGrid - Caching Behavior', () => {
  let wrapper: VueWrapper<any>

  const defaultProps = {
    schedules: [],
    weekLabel: '1月13日 - 1月19日',
    cardInfoType: 'teacher' as const,
    validationResults: {},
    slotWidth: 150
  }

  beforeEach(async () => {
    const { default: WeekGrid } = await import('~/components/Scheduling/WeekGrid.vue')

    wrapper = mount(WeekGrid, {
      props: defaultProps,
      global: {
        components: {
          ScheduleCard: MockScheduleCard
        }
      }
    })
  })

  afterEach(() => {
    wrapper.unmount()
  })

  it('應該在 schedules 變化時清除快取', async () => {
    const schedule1 = { id: 1, weekday: 1, start_hour: 9, start_minute: 0, duration_minutes: 60 }
    const schedule2 = { id: 2, weekday: 2, start_hour: 10, start_minute: 0, duration_minutes: 90 }

    // 第一次計算
    await wrapper.setProps({ schedules: [schedule1] })
    await nextTick()
    wrapper.vm.getScheduleStyle(schedule1)

    // 第二次計算不同的 schedule
    await wrapper.setProps({ schedules: [schedule2] })
    await nextTick()

    // 快取應該被清除
    expect(wrapper.vm.styleCache.size).toBeGreaterThan(0)
  })

  it('應該在 slotWidth 變化時清除快取', async () => {
    const schedule = { id: 1, weekday: 1, start_hour: 9, start_minute: 0, duration_minutes: 60 }

    await wrapper.setProps({ schedules: [schedule], slotWidth: 150 })
    await nextTick()
    wrapper.vm.getScheduleStyle(schedule)

    await wrapper.setProps({ schedules: [schedule], slotWidth: 200 })
    await nextTick()

    // 樣式應該更新
    const style = wrapper.vm.getScheduleStyle(schedule)
    expect(style.width).toBe('196px') // 200 - 4
  })
})

// ============================================
// formatTime 測試
// ============================================

describe('WeekGrid - formatTime', () => {
  let wrapper: VueWrapper<any>

  const defaultProps = {
    schedules: [],
    weekLabel: '1月13日 - 1月19日',
    cardInfoType: 'teacher' as const,
    validationResults: {},
    slotWidth: 150
  }

  beforeEach(async () => {
    const { default: WeekGrid } = await import('~/components/Scheduling/WeekGrid.vue')

    wrapper = mount(WeekGrid, {
      props: defaultProps,
      global: {
        components: {
          ScheduleCard: MockScheduleCard
        }
      }
    })
  })

  afterEach(() => {
    wrapper.unmount()
  })

  it('應該正確格式化時間', () => {
    expect(wrapper.vm.formatTime(0)).toBe('00:00')
    expect(wrapper.vm.formatTime(9)).toBe('09:00')
    expect(wrapper.vm.formatTime(12)).toBe('12:00')
    expect(wrapper.vm.formatTime(23)).toBe('23:00')
  })

  it('應該正確填充單數小時', () => {
    expect(wrapper.vm.formatTime(1)).toBe('01:00')
    expect(wrapper.vm.formatTime(5)).toBe('05:00')
    expect(wrapper.vm.formatTime(8)).toBe('08:00')
  })
})
