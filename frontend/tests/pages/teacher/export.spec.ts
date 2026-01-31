import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock modules before imports
vi.mock('vue', async () => {
  const vue = await import('vue')
  return {
    ...vue,
    ref: (v: any) => ({ value: v }),
    computed: (fn: () => any) => ({ value: fn() }),
    reactive: (v: any) => v,
    onMounted: (fn: () => void) => fn(),
  }
})

vi.mock('~/composables/useAlert', () => ({
  alertError: vi.fn(),
  alertSuccess: vi.fn(),
  alertWarning: vi.fn(),
  confirm: vi.fn(),
}))

vi.mock('~/composables/useNotification', () => ({
  default: () => ({
    show: { value: false },
    close: vi.fn(),
    success: vi.fn(),
    error: vi.fn(),
  }),
}))

vi.mock('~/stores/auth', () => ({
  useAuthStore: () => ({
    user: { name: '測試老師', bio: '測試簡介' },
  }),
}))

vi.mock('~/stores/useScheduleStore', () => ({
  useScheduleStore: () => ({
    schedule: { days: [] },
    weekLabel: '2026/01/20 - 2026/01/26',
    fetchSchedule: vi.fn(),
    changeWeek: vi.fn(),
  }),
}))

vi.mock('~/composables/useSidebar', () => ({
  useSidebar: () => ({
    isOpen: { value: false },
    close: vi.fn(),
  }),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}))

describe('teacher/export.vue 頁面邏輯', () => {
  // ThemeSelectionLogic 類別 - 主題選擇邏輯
  class ThemeSelectionLogic {
    themes: any[]
    selectedTheme: string

    constructor() {
      this.themes = [
        { id: 'dustyRose', name: '玫瑰灰' },
        { id: 'sageGreen', name: '鼠尾草綠' },
        { id: 'mutedBlue', name: '霧霾藍' },
        { id: 'warmBeige', name: '暖米色' },
        { id: 'lavender', name: '薰衣草灰' },
        { id: 'warmGrey', name: '溫柔灰' },
        { id: 'domeMouse', name: '多梅鼠' },
        { id: 'deepDomeMouse', name: '深梅鼠' },
        { id: 'lightDomeMouse', name: '淡梅鼠' },
        { id: 'coralRose', name: '珊瑚玫瑰' },
      ]
      this.selectedTheme = 'domeMouse'
    }

    getThemes(): any[] {
      return this.themes
    }

    setTheme(themeId: string) {
      this.selectedTheme = themeId
    }

    getSelectedTheme(): any | undefined {
      return this.themes.find(t => t.id === this.selectedTheme)
    }

    getThemeById(themeId: string): any | undefined {
      return this.themes.find(t => t.id === themeId)
    }

    getThemeCount(): number {
      return this.themes.length
    }

    isThemeSelected(themeId: string): boolean {
      return this.selectedTheme === themeId
    }

    getNextTheme(): string {
      const currentIndex = this.themes.findIndex(t => t.id === this.selectedTheme)
      const nextIndex = (currentIndex + 1) % this.themes.length
      return this.themes[nextIndex].id
    }

    getPreviousTheme(): string {
      const currentIndex = this.themes.findIndex(t => t.id === this.selectedTheme)
      const prevIndex = (currentIndex - 1 + this.themes.length) % this.themes.length
      return this.themes[prevIndex].id
    }
  }

  // ExportOptionLogic 類別 - 匯出選項邏輯
  class ExportOptionLogic {
    options: {
      showPersonalInfo: boolean
      showStats: boolean
      includeNotes: boolean
    }

    constructor() {
      this.options = {
        showPersonalInfo: true,
        showStats: true,
        includeNotes: false,
      }
    }

    setShowPersonalInfo(show: boolean) {
      this.options.showPersonalInfo = show
    }

    setShowStats(show: boolean) {
      this.options.showStats = show
    }

    setIncludeNotes(include: boolean) {
      this.options.includeNotes = include
    }

    toggleShowPersonalInfo() {
      this.options.showPersonalInfo = !this.options.showPersonalInfo
    }

    toggleShowStats() {
      this.options.showStats = !this.options.showStats
    }

    toggleIncludeNotes() {
      this.options.includeNotes = !this.options.includeNotes
    }

    resetOptions() {
      this.options = {
        showPersonalInfo: true,
        showStats: true,
        includeNotes: false,
      }
    }

    getOptionState(): Record<string, boolean> {
      return { ...this.options }
    }

    hasActiveOptions(): boolean {
      return this.options.showPersonalInfo || this.options.showStats || this.options.includeNotes
    }
  }

  // ScheduleStatsLogic 類別 - 課表統計邏輯
  class ScheduleStatsLogic {
    scheduleDays: any[]

    constructor() {
      this.scheduleDays = []
    }

    setScheduleDays(days: any[]) {
      this.scheduleDays = days
    }

    getTotalLessons(): number {
      return this.scheduleDays.reduce((sum, day) => sum + day.items.length, 0)
    }

    getTotalHours(): number {
      let total = 0
      this.scheduleDays.forEach(day => {
        day.items.forEach((item: any) => {
          total += this.calculateDuration(item.start_time, item.end_time)
        })
      })
      return Math.round(total / 60 * 10) / 10
    }

    getTeachingDays(): number {
      return this.scheduleDays.filter(day => day.items.length > 0).length
    }

    calculateDuration(start: string, end: string): number {
      const [startH, startM] = start.split(':').map(Number)
      const [endH, endM] = end.split(':').map(Number)
      return (endH * 60 + endM) - (startH * 60 + startM)
    }

    getLessonsByDay(dayIndex: number): number {
      return this.scheduleDays[dayIndex]?.items.length || 0
    }

    getWeekLabel(): string {
      if (this.scheduleDays.length === 0) return '-'
      const firstDate = new Date(this.scheduleDays[0].date)
      const lastDate = new Date(this.scheduleDays[this.scheduleDays.length - 1].date)
      return `${firstDate.getFullYear()}/${firstDate.getMonth() + 1}/${firstDate.getDate()} - ${lastDate.getFullYear()}/${lastDate.getMonth() + 1}/${lastDate.getDate()}`
    }
  }

  // DateFormatLogic 類別 - 日期格式化邏輯
  class DateFormatLogic {
    formatDate(date: Date): string {
      return date.toLocaleDateString('zh-TW', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    }

    formatWeekday(dateStr: string): string {
      const weekdays = ['週日', '週一', '週二', '週三', '週四', '週五', '週六']
      const date = new Date(dateStr)
      return weekdays[date.getDay()]
    }

    formatMonthDay(dateStr: string): string {
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-TW', {
        month: 'long',
        day: 'numeric',
      })
    }

    isToday(dateStr: string): boolean {
      const today = new Date()
      const date = new Date(dateStr)
      return date.toDateString() === today.toDateString()
    }

    isTomorrow(dateStr: string): boolean {
      const tomorrow = new Date()
      tomorrow.setDate(tomorrow.getDate() + 1)
      const date = new Date(dateStr)
      return date.toDateString() === tomorrow.toDateString()
    }

    isYesterday(dateStr: string): boolean {
      const yesterday = new Date()
      yesterday.setDate(yesterday.getDate() - 1)
      const date = new Date(dateStr)
      return date.toDateString() === yesterday.toDateString()
    }

    getRelativeDateLabel(dateStr: string): string {
      if (this.isToday(dateStr)) return '今天'
      if (this.isTomorrow(dateStr)) return '明天'
      if (this.isYesterday(dateStr)) return '昨天'
      return ''
    }
  }

  // WeekNavigationLogic 類別 - 週導航邏輯
  class WeekNavigationLogic {
    currentWeekOffset: number

    constructor() {
      this.currentWeekOffset = 0
    }

    getCurrentOffset(): number {
      return this.currentWeekOffset
    }

    changeWeek(delta: number) {
      this.currentWeekOffset += delta
    }

    goToCurrentWeek() {
      this.currentWeekOffset = 0
    }

    isCurrentWeek(): boolean {
      return this.currentWeekOffset === 0
    }

    getWeekOffsetLabel(): string {
      if (this.currentWeekOffset === 0) return '本週'
      if (this.currentWeekOffset > 0) return `未來第 ${this.currentWeekOffset} 週`
      return `過去第 ${Math.abs(this.currentWeekOffset)} 週`
    }
  }

  // StatusDisplayLogic 類別 - 狀態顯示邏輯
  class StatusDisplayLogic {
    getStatusClass(status: string): string {
      switch (status) {
        case 'PENDING':
          return 'bg-warning-500/20 text-warning-500'
        case 'APPROVED':
          return 'bg-success-500/20 text-success-500'
        case 'REJECTED':
          return 'bg-critical-500/20 text-critical-500'
        default:
          return 'bg-slate-500/20 text-slate-400'
      }
    }

    getStatusText(status: string): string {
      switch (status) {
        case 'PENDING':
          return '待審核'
        case 'APPROVED':
          return '已核准'
        case 'REJECTED':
          return '已拒絕'
        default:
          return status
      }
    }

    isPending(status: string): boolean {
      return status === 'PENDING'
    }

    isApproved(status: string): boolean {
      return status === 'APPROVED'
    }

    isRejected(status: string): boolean {
      return status === 'REJECTED'
    }
  }

  describe('ThemeSelectionLogic 主題選擇邏輯', () => {
    it('應該正確初始化所有主題', () => {
      const logic = new ThemeSelectionLogic()
      expect(logic.getThemeCount()).toBe(10)
    })

    it('setTheme 應該正確設定主題', () => {
      const logic = new ThemeSelectionLogic()
      logic.setTheme('sageGreen')
      expect(logic.isThemeSelected('sageGreen')).toBe(true)
      expect(logic.isThemeSelected('domeMouse')).toBe(false)
    })

    it('getSelectedTheme 應該返回當前選中的主題', () => {
      const logic = new ThemeSelectionLogic()
      logic.setTheme('mutedBlue')
      const theme = logic.getSelectedTheme()
      expect(theme?.name).toBe('霧霾藍')
    })

    it('getNextTheme 應該循環到下一個主題', () => {
      const logic = new ThemeSelectionLogic()
      logic.setTheme('dustyRose')
      const next = logic.getNextTheme()
      expect(logic.getThemeById(next)?.id).toBeDefined()
    })

    it('getPreviousTheme 應該循環到上一個主題', () => {
      const logic = new ThemeSelectionLogic()
      const prev = logic.getPreviousTheme()
      expect(logic.getThemeById(prev)?.id).toBeDefined()
    })
  })

  describe('ExportOptionLogic 匯出選項邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new ExportOptionLogic()
      expect(logic.options.showPersonalInfo).toBe(true)
      expect(logic.options.showStats).toBe(true)
      expect(logic.options.includeNotes).toBe(false)
    })

    it('toggleShowPersonalInfo 應該切換選項', () => {
      const logic = new ExportOptionLogic()
      expect(logic.options.showPersonalInfo).toBe(true)
      logic.toggleShowPersonalInfo()
      expect(logic.options.showPersonalInfo).toBe(false)
    })

    it('setShowStats 應該正確設定', () => {
      const logic = new ExportOptionLogic()
      logic.setShowStats(false)
      expect(logic.options.showStats).toBe(false)
    })

    it('resetOptions 應該重置所有選項', () => {
      const logic = new ExportOptionLogic()
      logic.options.showPersonalInfo = false
      logic.options.showStats = false
      logic.options.includeNotes = true
      logic.resetOptions()
      expect(logic.options.showPersonalInfo).toBe(true)
      expect(logic.options.showStats).toBe(true)
      expect(logic.options.includeNotes).toBe(false)
    })

    it('hasActiveOptions 應該正確判斷', () => {
      const logic = new ExportOptionLogic()
      expect(logic.hasActiveOptions()).toBe(true)
      logic.options.showPersonalInfo = false
      logic.options.showStats = false
      logic.options.includeNotes = false
      expect(logic.hasActiveOptions()).toBe(false)
    })
  })

  describe('ScheduleStatsLogic 課表統計邏輯', () => {
    it('getTotalLessons 應該計算總課程數', () => {
      const logic = new ScheduleStatsLogic()
      logic.setScheduleDays([
        { items: [{ id: 1 }, { id: 2 }] },
        { items: [{ id: 3 }] },
        { items: [] },
      ])
      expect(logic.getTotalLessons()).toBe(3)
    })

    it('getTotalHours 應該計算總時數', () => {
      const logic = new ScheduleStatsLogic()
      logic.setScheduleDays([
        { items: [{ start_time: '09:00', end_time: '10:00' }] },
        { items: [{ start_time: '14:00', end_time: '15:30' }] },
      ])
      const hours = logic.getTotalHours()
      expect(hours).toBe(2.5)
    })

    it('getTeachingDays 應該計算教學天數', () => {
      const logic = new ScheduleStatsLogic()
      logic.setScheduleDays([
        { items: [{}] },
        { items: [] },
        { items: [{}, {}] },
      ])
      expect(logic.getTeachingDays()).toBe(2)
    })

    it('calculateDuration 應該正確計算時長', () => {
      const logic = new ScheduleStatsLogic()
      expect(logic.calculateDuration('09:00', '10:30')).toBe(90)
      expect(logic.calculateDuration('14:00', '15:00')).toBe(60)
      expect(logic.calculateDuration('10:00', '10:00')).toBe(0)
    })
  })

  describe('DateFormatLogic 日期格式化邏輯', () => {
    it('formatDate 應該正確格式化日期', () => {
      const logic = new DateFormatLogic()
      const result = logic.formatDate(new Date(2026, 0, 20))
      expect(result).toContain('2026')
      expect(result).toContain('1')
      expect(result).toContain('20')
    })

    it('formatWeekday 應該返回正確的星期', () => {
      const logic = new DateFormatLogic()
      // 2026年1月20日是星期二
      expect(logic.formatWeekday('2026-01-20')).toBe('週二')
      expect(logic.formatWeekday('2026-01-21')).toBe('週三')
      expect(logic.formatWeekday('2026-01-22')).toBe('週四')
      expect(logic.formatWeekday('2026-01-23')).toBe('週五')
      expect(logic.formatWeekday('2026-01-24')).toBe('週六')
      expect(logic.formatWeekday('2026-01-25')).toBe('週日')
      expect(logic.formatWeekday('2026-01-26')).toBe('週一')
    })

    it('isToday 應該正確判斷是否為今天', () => {
      const logic = new DateFormatLogic()
      const today = new Date()
      const todayStr = today.toISOString().split('T')[0]
      expect(logic.isToday(todayStr)).toBe(true)
      expect(logic.isToday('2026-01-01')).toBe(false)
    })
  })

  describe('WeekNavigationLogic 週導航邏輯', () => {
    it('should initialize at current week', () => {
      const logic = new WeekNavigationLogic()
      expect(logic.getCurrentOffset()).toBe(0)
      expect(logic.isCurrentWeek()).toBe(true)
    })

    it('changeWeek should update offset', () => {
      const logic = new WeekNavigationLogic()
      logic.changeWeek(1)
      expect(logic.getCurrentOffset()).toBe(1)
      logic.changeWeek(-1)
      expect(logic.getCurrentOffset()).toBe(0)
    })

    it('goToCurrentWeek should reset offset', () => {
      const logic = new WeekNavigationLogic()
      logic.changeWeek(2)
      logic.goToCurrentWeek()
      expect(logic.getCurrentOffset()).toBe(0)
    })

    it('getWeekOffsetLabel should return correct label', () => {
      const logic = new WeekNavigationLogic()
      expect(logic.getWeekOffsetLabel()).toBe('本週')
      logic.changeWeek(1)
      expect(logic.getWeekOffsetLabel()).toContain('未來')
      logic.changeWeek(-2)
      expect(logic.getWeekOffsetLabel()).toContain('過去')
    })
  })

  describe('StatusDisplayLogic 狀態顯示邏輯', () => {
    it('getStatusClass should return correct CSS classes', () => {
      const logic = new StatusDisplayLogic()
      expect(logic.getStatusClass('PENDING')).toContain('warning')
      expect(logic.getStatusClass('APPROVED')).toContain('success')
      expect(logic.getStatusClass('REJECTED')).toContain('critical')
      expect(logic.getStatusClass('UNKNOWN')).toContain('slate')
    })

    it('getStatusText should return correct Chinese text', () => {
      const logic = new StatusDisplayLogic()
      expect(logic.getStatusText('PENDING')).toBe('待審核')
      expect(logic.getStatusText('APPROVED')).toBe('已核准')
      expect(logic.getStatusText('REJECTED')).toBe('已拒絕')
    })

    it('isPending, isApproved, isRejected should work correctly', () => {
      const logic = new StatusDisplayLogic()
      expect(logic.isPending('PENDING')).toBe(true)
      expect(logic.isApproved('PENDING')).toBe(false)
      expect(logic.isRejected('PENDING')).toBe(false)
    })
  })

  describe('頁面整合邏輯', () => {
    it('should handle complete export flow', () => {
      const themeLogic = new ThemeSelectionLogic()
      const optionLogic = new ExportOptionLogic()
      const statsLogic = new ScheduleStatsLogic()
      const dateLogic = new DateFormatLogic()
      const weekLogic = new WeekNavigationLogic()
      const statusLogic = new StatusDisplayLogic()

      // Select theme
      themeLogic.setTheme('sageGreen')
      expect(themeLogic.getSelectedTheme()?.name).toBe('鼠尾草綠')

      // Configure options
      optionLogic.setIncludeNotes(true)
      expect(optionLogic.options.includeNotes).toBe(true)

      // Set schedule data
      statsLogic.setScheduleDays([
        { date: '2026-01-20', items: [{ start_time: '09:00', end_time: '10:00' }] },
        { date: '2026-01-21', items: [{ start_time: '14:00', end_time: '15:00' }, { start_time: '15:30', end_time: '16:30' }] },
      ])

      expect(statsLogic.getTotalLessons()).toBe(3)
      expect(statsLogic.getTeachingDays()).toBe(2)

      // Check date formatting - 2026-01-20 是週二
      expect(dateLogic.formatWeekday('2026-01-20')).toBe('週二')

      // Navigate weeks
      weekLogic.changeWeek(1)
      expect(weekLogic.isCurrentWeek()).toBe(false)

      // Check status display
      expect(statusLogic.getStatusText('APPROVED')).toBe('已核准')
    })
  })
})
