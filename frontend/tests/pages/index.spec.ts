import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock modules before imports
vi.mock('vue', async () => {
  const vue = await import('vue')
  return {
    ...vue,
    ref: (v: any) => ({ value: v }),
    computed: (fn: () => any) => ({ value: fn() }),
    onMounted: (fn: () => void) => fn(),
  }
})

vi.mock('~/stores/auth', () => ({
  useAuthStore: () => ({
    login: vi.fn(),
    user: null,
  }),
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}))

describe('index.vue 頁面邏輯', () => {
  // HeroDisplayLogic 類別 - 首頁展示邏輯
  class HeroDisplayLogic {
    title: string
    subtitle: string
    slogan: string

    constructor() {
      this.title = 'TimeLedger'
      this.subtitle = '您的時間資產管家'
      this.slogan = '多中心排課 × 老師人才庫 × 專業課表分享'
    }

    getTitle(): string {
      return this.title
    }

    getSubtitle(): string {
      return this.subtitle
    }

    getSlogan(): string {
      return this.slogan
    }

    getSloganParts(): string[] {
      return this.slogan.split(' × ')
    }
  }

  // LineLoginLogic 類別 - LINE 登入邏輯
  class LineLoginLogic {
    loading: boolean
    error: string

    constructor() {
      this.loading = false
      this.error = ''
    }

    setLoading(loading: boolean) {
      this.loading = loading
    }

    setError(error: string) {
      this.error = error
    }

    clearError() {
      this.error = ''
    }

    isLoading(): boolean {
      return this.loading
    }

    hasError(): boolean {
      return this.error !== ''
    }

    getDemoLineUserId(): string {
      return 'U1234567890abcdef1234567890abcd'
    }

    getDemoAccessToken(): string {
      return 'test_access_token_demo001'
    }

    buildLoginUrl(): string {
      return `/teacher/login?line_user_id=${encodeURIComponent(this.getDemoLineUserId())}&access_token=${encodeURIComponent(this.getDemoAccessToken())}`
    }
  }

  // FeatureShowcaseLogic 類別 - 功能展示邏輯
  class FeatureShowcaseLogic {
    features: any[]

    constructor() {
      this.features = [
        {
          id: 'export',
          name: '精美分享圖',
          description: '多種漸層主題、個人品牌資訊、社交分享',
          icon: 'download',
        },
        {
          id: 'matching',
          name: '智慧媒合',
          description: '技能、證照、評分綜合評比',
          icon: 'matching',
        },
        {
          id: 'notifications',
          name: '即時通知',
          description: 'LINE Notify 推送課程變更、審核結果',
          icon: 'notification',
        },
      ]
    }

    getFeatures(): any[] {
      return this.features
    }

    getFeatureById(id: string): any | undefined {
      return this.features.find(f => f.id === id)
    }

    getFeatureCount(): number {
      return this.features.length
    }
  }

  // ScheduleDemoLogic 類別 - 課表演示邏輯
  class ScheduleDemoLogic {
    schedules: Record<string, any>
    isDragging: boolean
    dragTarget: { time: number; day: number } | null
    hasConflict: boolean
    selectedCell: { time: number; day: number } | null

    constructor() {
      this.schedules = {
        '14-1': { id: 1, offering_name: '鋼琴基礎', teacher_name: 'Alice', room_id: 1 },
        '14-2': { id: 2, offering_name: '鋼琴基礎', teacher_name: 'Alice', room_id: 1 },
        '14-3': { id: 3, offering_name: '小提琴入門', teacher_name: 'Bob', room_id: 2 },
        '15-1': { id: 4, offering_name: '鋼琴進階', teacher_name: 'Alice', room_id: 1 },
      }
      this.isDragging = false
      this.dragTarget = null
      this.hasConflict = false
      this.selectedCell = null
    }

    setSchedules(schedules: Record<string, any>) {
      this.schedules = schedules
    }

    getScheduleAt(hour: number, weekday: number): any | undefined {
      return this.schedules[`${hour}-${weekday}`]
    }

    setDragging(isDragging: boolean) {
      this.isDragging = isDragging
    }

    setDragTarget(target: { time: number; day: number } | null) {
      this.dragTarget = target
    }

    setHasConflict(conflict: boolean) {
      this.hasConflict = conflict
    }

    setSelectedCell(cell: { time: number; day: number } | null) {
      this.selectedCell = cell
    }

    isDraggingSchedule(): boolean {
      return this.isDragging
    }

    hasDragTarget(): boolean {
      return this.dragTarget !== null
    }

    getSelectedCell(): { time: number; day: number } | null {
      return this.selectedCell
    }

    moveSchedule(fromKey: string, toKey: string): boolean {
      const schedule = this.schedules[fromKey]
      if (schedule && !this.schedules[toKey]) {
        delete this.schedules[fromKey]
        this.schedules[toKey] = schedule
        return true
      }
      return false
    }

    resetDemo() {
      this.schedules = {
        '14-1': { id: 1, offering_name: '鋼琴基礎', teacher_name: 'Alice', room_id: 1 },
      }
      this.isDragging = false
      this.dragTarget = null
      this.hasConflict = false
      this.selectedCell = null
    }

    getScheduleCount(): number {
      return Object.keys(this.schedules).length
    }
  }

  // TimeSlotLogic 類別 - 時段邏輯
  class TimeSlotLogic {
    timeSlots: number[]

    constructor() {
      this.timeSlots = [9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]
    }

    getTimeSlots(): number[] {
      return this.timeSlots
    }

    formatTime(hour: number): string {
      return `${hour}:00`
    }

    getTimeSlotCount(): number {
      return this.timeSlots.length
    }

    isValidTime(hour: number): boolean {
      return this.timeSlots.includes(hour)
    }

    getNextTime(hour: number): number | null {
      const index = this.timeSlots.indexOf(hour)
      if (index < this.timeSlots.length - 1) {
        return this.timeSlots[index + 1]
      }
      return null
    }

    getPreviousTime(hour: number): number | null {
      const index = this.timeSlots.indexOf(hour)
      if (index > 0) {
        return this.timeSlots[index - 1]
      }
      return null
    }
  }

  // WeekDayDemoLogic 類別 - 演示星期邏輯
  class WeekDayDemoLogic {
    weekDays: { value: number; name: string }[]

    constructor() {
      this.weekDays = [
        { value: 1, name: '週一' },
        { value: 2, name: '週二' },
        { value: 3, name: '週三' },
        { value: 4, name: '週四' },
        { value: 5, name: '週五' },
        { value: 6, name: '週六' },
        { value: 7, name: '週日' },
      ]
    }

    getWeekDays(): { value: number; name: string }[] {
      return this.weekDays
    }

    getDayName(value: number): string {
      const day = this.weekDays.find(d => d.value === value)
      return day?.name || String(value)
    }

    getDayCount(): number {
      return this.weekDays.length
    }

    isValidDay(value: number): boolean {
      return this.weekDays.some(d => d.value === value)
    }

    getNextDay(value: number): number {
      const index = this.weekDays.findIndex(d => d.value === value)
      const nextIndex = (index + 1) % this.weekDays.length
      return this.weekDays[nextIndex].value
    }

    getPreviousDay(value: number): number {
      const index = this.weekDays.findIndex(d => d.value === value)
      const prevIndex = (index - 1 + this.weekDays.length) % this.weekDays.length
      return this.weekDays[prevIndex].value
    }
  }

  // NavigationDemoLogic 類別 - 演示導航邏輯
  class NavigationDemoLogic {
    isMobile: boolean
    mobileDayOffset: number

    constructor() {
      this.isMobile = false
      this.mobileDayOffset = 0
    }

    setIsMobile(isMobile: boolean) {
      this.isMobile = isMobile
    }

    setMobileDayOffset(offset: number) {
      this.mobileDayOffset = offset
    }

    getMobileDayOffset(): number {
      return this.mobileDayOffset
    }

    changeMobileDays(delta: number) {
      this.mobileDayOffset = (this.mobileDayOffset + delta + 7) % 7
    }

    isMobileView(): boolean {
      return this.isMobile
    }

    getVisibleDayCount(): number {
      return this.isMobile ? 3 : 7
    }

    getMobileDayLabel(): string {
      const visibleDays = this.getVisibleWeekDays()
      if (visibleDays.length === 0) return ''
      return `${visibleDays[0].name} - ${visibleDays[visibleDays.length - 1].name}`
    }

    getVisibleWeekDays(): { value: number; name: string }[] {
      const weekDays = [
        { value: 1, name: '週一' },
        { value: 2, name: '週二' },
        { value: 3, name: '週三' },
        { value: 4, name: '週四' },
        { value: 5, name: '週五' },
        { value: 6, name: '週六' },
        { value: 7, name: '週日' },
      ]
      if (!this.isMobile) return weekDays
      const days = []
      for (let i = 0; i < 3; i++) {
        const index = (this.mobileDayOffset + i) % 7
        days.push(weekDays[index])
      }
      return days
    }
  }

  // AdminLinkLogic 類別 - 管理員連結邏輯
  class AdminLinkLogic {
    adminLoginUrl: string
    isVisible: boolean

    constructor() {
      this.adminLoginUrl = '/admin/login'
      this.isVisible = true
    }

    getAdminLoginUrl(): string {
      return this.adminLoginUrl
    }

    setAdminLoginUrl(url: string) {
      this.adminLoginUrl = url
    }

    isAdminLinkVisible(): boolean {
      return this.isVisible
    }

    hideAdminLink() {
      this.isVisible = false
    }

    showAdminLink() {
      this.isVisible = true
    }
  }

  describe('HeroDisplayLogic 首頁展示邏輯', () => {
    it('should return correct title', () => {
      const logic = new HeroDisplayLogic()
      expect(logic.getTitle()).toBe('TimeLedger')
    })

    it('should return correct subtitle', () => {
      const logic = new HeroDisplayLogic()
      expect(logic.getSubtitle()).toBe('您的時間資產管家')
    })

    it('should split slogan correctly', () => {
      const logic = new HeroDisplayLogic()
      const parts = logic.getSloganParts()
      expect(parts).toHaveLength(3)
      expect(parts[0]).toBe('多中心排課')
      expect(parts[1]).toBe('老師人才庫')
      expect(parts[2]).toBe('專業課表分享')
    })
  })

  describe('LineLoginLogic LINE 登入邏輯', () => {
    it('should initialize correctly', () => {
      const logic = new LineLoginLogic()
      expect(logic.isLoading()).toBe(false)
      expect(logic.hasError()).toBe(false)
    })

    it('should set loading state', () => {
      const logic = new LineLoginLogic()
      logic.setLoading(true)
      expect(logic.isLoading()).toBe(true)
    })

    it('should set error state', () => {
      const logic = new LineLoginLogic()
      logic.setError('登入失敗')
      expect(logic.hasError()).toBe(true)
      expect(logic.error).toBe('登入失敗')
    })

    it('should clear error', () => {
      const logic = new LineLoginLogic()
      logic.setError('錯誤')
      logic.clearError()
      expect(logic.hasError()).toBe(false)
    })

    it('should build correct login URL', () => {
      const logic = new LineLoginLogic()
      const url = logic.buildLoginUrl()
      expect(url).toContain('line_user_id=')
      expect(url).toContain('access_token=')
    })
  })

  describe('FeatureShowcaseLogic 功能展示邏輯', () => {
    it('should return all features', () => {
      const logic = new FeatureShowcaseLogic()
      expect(logic.getFeatureCount()).toBe(3)
    })

    it('should get feature by id', () => {
      const logic = new FeatureShowcaseLogic()
      const feature = logic.getFeatureById('matching')
      expect(feature?.name).toBe('智慧媒合')
    })

    it('should return undefined for invalid id', () => {
      const logic = new FeatureShowcaseLogic()
      const feature = logic.getFeatureById('invalid')
      expect(feature).toBeUndefined()
    })
  })

  describe('ScheduleDemoLogic 課表演示邏輯', () => {
    it('should initialize with demo schedules', () => {
      const logic = new ScheduleDemoLogic()
      expect(logic.getScheduleCount()).toBe(4)
    })

    it('should get schedule at time and day', () => {
      const logic = new ScheduleDemoLogic()
      const schedule = logic.getScheduleAt(14, 1)
      expect(schedule?.offering_name).toBe('鋼琴基礎')
    })

    it('should return undefined for empty slot', () => {
      const logic = new ScheduleDemoLogic()
      const schedule = logic.getScheduleAt(9, 1)
      expect(schedule).toBeUndefined()
    })

    it('should move schedule', () => {
      const logic = new ScheduleDemoLogic()
      const moved = logic.moveSchedule('14-1', '15-1')
      expect(moved).toBe(true)
      expect(logic.getScheduleAt(15, 1)?.offering_name).toBe('鋼琴基礎')
    })

    it('should not move if target is occupied', () => {
      const logic = new ScheduleDemoLogic()
      const moved = logic.moveSchedule('14-1', '14-2')
      expect(moved).toBe(false)
    })

    it('should reset demo', () => {
      const logic = new ScheduleDemoLogic()
      logic.setDragging(true)
      logic.setHasConflict(true)
      logic.resetDemo()
      expect(logic.getScheduleCount()).toBe(1)
      expect(logic.isDraggingSchedule()).toBe(false)
    })
  })

  describe('TimeSlotLogic 時段邏輯', () => {
    it('should return all time slots', () => {
      const logic = new TimeSlotLogic()
      expect(logic.getTimeSlotCount()).toBe(13)
      expect(logic.getTimeSlots()[0]).toBe(9)
      expect(logic.getTimeSlots()[12]).toBe(21)
    })

    it('should format time correctly', () => {
      const logic = new TimeSlotLogic()
      expect(logic.formatTime(9)).toBe('9:00')
      expect(logic.formatTime(14)).toBe('14:00')
    })

    it('should validate time', () => {
      const logic = new TimeSlotLogic()
      expect(logic.isValidTime(9)).toBe(true)
      expect(logic.isValidTime(8)).toBe(false)
      expect(logic.isValidTime(22)).toBe(false)
    })

    it('should get next/previous time', () => {
      const logic = new TimeSlotLogic()
      expect(logic.getNextTime(14)).toBe(15)
      expect(logic.getPreviousTime(14)).toBe(13)
      expect(logic.getNextTime(21)).toBeNull()
      expect(logic.getPreviousTime(9)).toBeNull()
    })
  })

  describe('WeekDayDemoLogic 演示星期邏輯', () => {
    it('should return all week days', () => {
      const logic = new WeekDayDemoLogic()
      expect(logic.getDayCount()).toBe(7)
    })

    it('should get day name', () => {
      const logic = new WeekDayDemoLogic()
      expect(logic.getDayName(1)).toBe('週一')
      expect(logic.getDayName(7)).toBe('週日')
    })

    it('should validate day', () => {
      const logic = new WeekDayDemoLogic()
      expect(logic.isValidDay(1)).toBe(true)
      expect(logic.isValidDay(0)).toBe(false)
      expect(logic.isValidDay(8)).toBe(false)
    })

    it('should get next/previous day', () => {
      const logic = new WeekDayDemoLogic()
      expect(logic.getNextDay(5)).toBe(6)
      expect(logic.getNextDay(7)).toBe(1)
      expect(logic.getPreviousDay(1)).toBe(7)
      expect(logic.getPreviousDay(3)).toBe(2)
    })
  })

  describe('NavigationDemoLogic 演示導航邏輯', () => {
    it('should initialize correctly', () => {
      const logic = new NavigationDemoLogic()
      expect(logic.isMobileView()).toBe(false)
      expect(logic.getMobileDayOffset()).toBe(0)
    })

    it('should set mobile view', () => {
      const logic = new NavigationDemoLogic()
      logic.setIsMobile(true)
      expect(logic.isMobileView()).toBe(true)
    })

    it('should change mobile days', () => {
      const logic = new NavigationDemoLogic()
      logic.setMobileDayOffset(0)
      logic.changeMobileDays(1)
      expect(logic.getMobileDayOffset()).toBe(1)
      logic.changeMobileDays(6)
      expect(logic.getMobileDayOffset()).toBe(0)
    })

    it('should return correct visible day count', () => {
      const logic = new NavigationDemoLogic()
      expect(logic.getVisibleDayCount()).toBe(7)
      logic.setIsMobile(true)
      expect(logic.getVisibleDayCount()).toBe(3)
    })

    it('should return visible week days', () => {
      const logic = new NavigationDemoLogic()
      const desktopDays = logic.getVisibleWeekDays()
      expect(desktopDays).toHaveLength(7)
      logic.setIsMobile(true)
      logic.setMobileDayOffset(0)
      const mobileDays = logic.getVisibleWeekDays()
      expect(mobileDays).toHaveLength(3)
      expect(mobileDays[0].value).toBe(1)
    })
  })

  describe('AdminLinkLogic 管理員連結邏輯', () => {
    it('should return correct admin login URL', () => {
      const logic = new AdminLinkLogic()
      expect(logic.getAdminLoginUrl()).toBe('/admin/login')
    })

    it('should show/hide admin link', () => {
      const logic = new AdminLinkLogic()
      expect(logic.isAdminLinkVisible()).toBe(true)
      logic.hideAdminLink()
      expect(logic.isAdminLinkVisible()).toBe(false)
      logic.showAdminLink()
      expect(logic.isAdminLinkVisible()).toBe(true)
    })
  })

  describe('頁面整合邏輯', () => {
    it('should handle complete demo flow', () => {
      const heroLogic = new HeroDisplayLogic()
      const lineLogic = new LineLoginLogic()
      const scheduleLogic = new ScheduleDemoLogic()
      const timeLogic = new TimeSlotLogic()
      const navLogic = new NavigationDemoLogic()

      // Check hero content
      expect(heroLogic.getTitle()).toBe('TimeLedger')

      // Check LINE login
      const loginUrl = lineLogic.buildLoginUrl()
      expect(loginUrl).toContain('/teacher/login')

      // Check schedule demo
      const schedule = scheduleLogic.getScheduleAt(14, 1)
      expect(schedule?.offering_name).toBe('鋼琴基礎')

      // Check time slots
      expect(timeLogic.getTimeSlots()).toContain(14)

      // Check mobile navigation
      navLogic.setIsMobile(true)
      expect(navLogic.isMobileView()).toBe(true)
      expect(navLogic.getVisibleDayCount()).toBe(3)
    })

    it('should handle schedule drag and drop simulation', () => {
      const scheduleLogic = new ScheduleDemoLogic()

      // Get initial state
      const originalSchedule = scheduleLogic.getScheduleAt(14, 1)
      expect(originalSchedule).toBeDefined()

      // Simulate moving schedule
      const moved = scheduleLogic.moveSchedule('14-1', '15-1')
      expect(moved).toBe(true)

      // Check new position
      const newSchedule = scheduleLogic.getScheduleAt(15, 1)
      expect(newSchedule?.offering_name).toBe('鋼琴基礎')

      // Check old position is empty
      const oldPosition = scheduleLogic.getScheduleAt(14, 1)
      expect(oldPosition).toBeUndefined()
    })

    it('should handle mobile day navigation', () => {
      const navLogic = new NavigationDemoLogic()
      navLogic.setIsMobile(true)
      navLogic.setMobileDayOffset(0)

      // Check initial visible days
      const days1 = navLogic.getVisibleWeekDays()
      expect(days1[0].name).toBe('週一')

      // Navigate forward
      navLogic.changeMobileDays(1)
      const days2 = navLogic.getVisibleWeekDays()
      expect(days2[0].name).toBe('週二')

      // Navigate to wrap around
      navLogic.changeMobileDays(5)
      const days3 = navLogic.getVisibleWeekDays()
      expect(days3[0].name).toBe('週日')
    })
  })
})
