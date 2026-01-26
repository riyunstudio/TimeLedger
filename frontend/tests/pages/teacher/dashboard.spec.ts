import { describe, it, expect, vi, beforeEach } from 'vitest'

// ============================================
// Mock Setup
// ============================================

// Mock useAlert composable
const mockAlert = {
  error: vi.fn(),
  confirm: vi.fn().mockResolvedValue(true),
  success: vi.fn(),
}

vi.mock('~/composables/useAlert', () => ({
  alertError: mockAlert.error,
  alertConfirm: mockAlert.confirm,
  alertSuccess: mockAlert.success,
}))

// Mock useNotification
vi.mock('~/composables/useNotification', () => ({
  useNotification: () => ({
    show: { value: false },
    close: vi.fn(),
  }),
}))

// Mock useSidebar
const mockSidebarStore = {
  isOpen: { value: false },
  close: vi.fn(),
  open: vi.fn(),
}

vi.mock('~/composables/useSidebar', () => ({
  useSidebar: () => mockSidebarStore,
}))

// ============================================
// Simple Reactive Wrapper
// ============================================

function createRef<T>(value: T) {
  return { value }
}

// ============================================
// Teacher Dashboard Page Logic
// ============================================

class TeacherDashboardPageLogic {
  // 狀態
  viewMode = createRef('grid')
  listCurrentDate = createRef('')
  isMobile = createRef(false)
  gridDayOffset = createRef(0)
  showPersonalEventModal = createRef(false)
  editingEvent = createRef<any>(null)
  showSessionNoteModal = createRef(false)
  selectedScheduleItem = createRef<any>(null)
  isDragging = createRef(false)
  dragTarget = createRef<{ time: number; date: string } | null>(null)
  draggedItem = createRef<any>(null)
  sourceDate = createRef('')
  sourceHour = createRef(0)

  // 時間槽
  timeSlots = [9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]

  // 模擬課表資料
  mockSchedule = {
    days: [
      { date: '2026-01-20', items: [] },
      { date: '2026-01-21', items: [] },
      { date: '2026-01-22', items: [] },
      { date: '2026-01-23', items: [] },
      { date: '2026-01-24', items: [] },
      { date: '2026-01-25', items: [] },
      { date: '2026-01-26', items: [] },
    ],
  }

  // Computed（改用 getter）
  get allWeekDays() {
    return this.mockSchedule.days.map(day => {
      const date = new Date(day.date)
      return {
        date: day.date,
        weekday: date.toLocaleDateString('zh-TW', { weekday: 'short' }),
        day: date.getDate(),
      }
    })
  }

  get displayWeekDays() {
    if (!this.isMobile.value) return this.allWeekDays

    const start = this.gridDayOffset.value
    const days = []
    for (let i = 0; i < 3; i++) {
      const index = (start + i) % 7
      if (this.allWeekDays[index]) {
        days.push(this.allWeekDays[index])
      }
    }
    return days
  }

  get gridColsClass() {
    return this.isMobile.value ? 'grid-cols-[50px_repeat(3,1fr)]' : 'grid-cols-[80px_repeat(7,1fr)]'
  }

  get gridDayLabel() {
    if (!this.displayWeekDays.length) return ''
    const start = this.displayWeekDays[0]
    const end = this.displayWeekDays[this.displayWeekDays.length - 1]
    return `${start.weekday} - ${end.weekday}`
  }

  get currentDayItems() {
    if (!this.listCurrentDate.value) return []

    const day = this.mockSchedule.days.find(d => d.date === this.listCurrentDate.value)
    return day?.items || []
  }

  // 方法
  changeGridDay(delta: number) {
    this.gridDayOffset.value = (this.gridDayOffset.value + delta + 7) % 7
  }

  changeListDay(delta: number, schedule: any) {
    if (!schedule) return

    const currentIndex = schedule.days.findIndex(d => d.date === this.listCurrentDate.value)
    const newIndex = currentIndex + delta

    if (newIndex >= 0 && newIndex < schedule.days.length) {
      this.listCurrentDate.value = schedule.days[newIndex].date
    }
  }

  formatDate(dateStr: string): string {
    const date = new Date(dateStr)
    const today = new Date()
    today.setHours(0, 0, 0, 0)

    const diffDays = Math.floor((date.getTime() - today.getTime()) / (1000 * 60 * 60 * 24))

    if (diffDays === 0) return '今天'
    if (diffDays === 1) return '明天'
    if (diffDays === -1) return '昨天'

    return date.toLocaleDateString('zh-TW', {
      month: 'long',
      day: 'numeric',
      weekday: 'short',
    })
  }

  getGridCellClass(date: string, hour: number): string {
    // 模擬檢查該時段是否有行程
    const hasItems = this.mockSchedule.days.some(day =>
      day.date === date && day.items.some(item => {
        const startHour = parseInt(item.start_time?.split(':')[0] || '0')
        const endHour = parseInt(item.end_time?.split(':')[0] || '0')
        return hour >= startHour && hour < endHour
      })
    )

    if (hasItems) return ''

    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const cellDate = new Date(date)
    const isPast = cellDate < today

    if (isPast) return 'bg-slate-800/50'
    return 'hover:bg-white/5'
  }

  getItemBgClass(item: any): string {
    if (item.type === 'PERSONAL_EVENT') {
      return 'bg-purple-500/30 border border-purple-500/50'
    }

    const data = item.data
    if (data?.has_exception) {
      return 'bg-warning-500/30 border border-warning-500/50'
    }

    return 'bg-success-500/20 border border-success-500/30'
  }

  getItemBorderClass(item: any): string {
    if (item.type === 'PERSONAL_EVENT') {
      return 'border-purple-500/50 bg-purple-500/10'
    }

    const data = item.data
    if (data?.has_exception) {
      return 'border-warning-500/50 bg-warning-500/10'
    }

    return 'border-success-500/50 bg-success-500/10'
  }

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

  isTargetCell(hour: number, date: string): boolean {
    return this.dragTarget.value?.time === hour && this.dragTarget.value?.date === date
  }

  handleDragStart(item: any, hour: number, date: string) {
    this.isDragging.value = true
    this.draggedItem.value = item
    this.sourceHour.value = hour
    this.sourceDate.value = date
  }

  handleDragEnd() {
    this.isDragging.value = false
    this.dragTarget.value = null
    this.draggedItem.value = null
  }

  handleDragEnter(hour: number, date: string) {
    this.dragTarget.value = { time: hour, date }
  }

  handleDragLeave() {
    // 空實作
  }

  openItemDetail(item: any) {
    if (item.type === 'SCHEDULE_RULE' || item.type === 'CENTER_SESSION') {
      this.selectedScheduleItem.value = item
      this.showSessionNoteModal.value = true
    } else if (item.type === 'PERSONAL_EVENT') {
      this.editingEvent.value = item.data
      this.showPersonalEventModal.value = true
    }
  }

  handleNoteModalClose() {
    this.showSessionNoteModal.value = false
    this.selectedScheduleItem.value = null
  }
}

// ============================================
// Test Suites
// ============================================

describe('Teacher Dashboard Logic - View Mode', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該預設為網格視圖', () => {
    const logic = new TeacherDashboardPageLogic()
    expect(logic.viewMode.value).toBe('grid')
  })

  it('應該正確切換到列表視圖', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.viewMode.value = 'grid'

    logic.viewMode.value = 'list'

    expect(logic.viewMode.value).toBe('list')
  })

  it('應該正確切換回網格視圖', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.viewMode.value = 'list'

    logic.viewMode.value = 'grid'

    expect(logic.viewMode.value).toBe('grid')
  })
})

describe('Teacher Dashboard Logic - Week Navigation', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確計算顯示的天數（非行動裝置）', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.isMobile.value = false

    expect(logic.displayWeekDays.length).toBe(7)
  })

  it('應該正確計算顯示的天數（行動裝置，偏移為0）', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.isMobile.value = true
    logic.gridDayOffset.value = 0

    expect(logic.displayWeekDays.length).toBe(3)
  })

  it('應該正確計算顯示的天數（行動裝置，偏移為1）', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.isMobile.value = true
    logic.gridDayOffset.value = 1

    expect(logic.displayWeekDays.length).toBe(3)
  })

  it('應該正確切換網格日期（往後）', () => {
    const logic = new TeacherDashboardPageLogic()
    const initialOffset = logic.gridDayOffset.value

    logic.changeGridDay(1)

    expect(logic.gridDayOffset.value).toBe((initialOffset + 1) % 7)
  })

  it('應該正確切換網格日期（往前）', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.gridDayOffset.value = 0

    logic.changeGridDay(-1)

    expect(logic.gridDayOffset.value).toBe(6) // 應該繞回 6
  })
})

describe('Teacher Dashboard Logic - List View Navigation', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確切換到下一天', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.listCurrentDate.value = '2026-01-20'

    logic.changeListDay(1, logic.mockSchedule)

    expect(logic.listCurrentDate.value).toBe('2026-01-21')
  })

  it('應該正確切換到上一天', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.listCurrentDate.value = '2026-01-21'

    logic.changeListDay(-1, logic.mockSchedule)

    expect(logic.listCurrentDate.value).toBe('2026-01-20')
  })

  it('應該在邊界時停止切換（第一天不能往前）', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.listCurrentDate.value = '2026-01-20'

    logic.changeListDay(-1, logic.mockSchedule)

    expect(logic.listCurrentDate.value).toBe('2026-01-20')
  })

  it('應該在邊界時停止切換（最後一天不能往後）', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.listCurrentDate.value = '2026-01-26'

    logic.changeListDay(1, logic.mockSchedule)

    expect(logic.listCurrentDate.value).toBe('2026-01-26')
  })
})

describe('Teacher Dashboard Logic - Date Formatting', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確格式化今天的日期', () => {
    const logic = new TeacherDashboardPageLogic()
    const today = new Date()
    const todayStr = today.toISOString().split('T')[0]

    expect(logic.formatDate(todayStr)).toBe('今天')
  })

  it('應該正確格式化明天的日期', () => {
    const logic = new TeacherDashboardPageLogic()
    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    const tomorrowStr = tomorrow.toISOString().split('T')[0]

    expect(logic.formatDate(tomorrowStr)).toBe('明天')
  })

  it('應該正確格式化昨天的日期', () => {
    const logic = new TeacherDashboardPageLogic()
    const yesterday = new Date()
    yesterday.setDate(yesterday.getDate() - 1)
    const yesterdayStr = yesterday.toISOString().split('T')[0]

    expect(logic.formatDate(yesterdayStr)).toBe('昨天')
  })

  it('應該正確格式化一般日期', () => {
    const logic = new TeacherDashboardPageLogic()

    const result = logic.formatDate('2026-01-15')

    expect(result).toContain('1')
    expect(result).toContain('15')
    expect(result).toContain('週')
  })
})

describe('Teacher Dashboard Logic - Grid Cell Class', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該在有行程時返回空字串', () => {
    const logic = new TeacherDashboardPageLogic()
    // 模擬有行程的狀況
    logic.mockSchedule.days[0].items = [
      { id: 1, start_time: '10:00', end_time: '11:00' }
    ]

    const result = logic.getGridCellClass('2026-01-20', 10)

    expect(result).toBe('')
  })

  it('應該在過去日期無行程時返回灰色背景', () => {
    const logic = new TeacherDashboardPageLogic()
    // 過去日期
    const pastDate = new Date()
    pastDate.setDate(pastDate.getDate() - 5)
    const pastDateStr = pastDate.toISOString().split('T')[0]

    const result = logic.getGridCellClass(pastDateStr, 10)

    expect(result).toContain('bg-slate-800/50')
  })

  it('應該在未來日期無行程時返回懸停效果', () => {
    const logic = new TeacherDashboardPageLogic()
    // 未來日期
    const futureDate = new Date()
    futureDate.setDate(futureDate.getDate() + 10)
    const futureDateStr = futureDate.toISOString().split('T')[0]

    const result = logic.getGridCellClass(futureDateStr, 10)

    expect(result).toContain('hover:bg-white/5')
  })
})

describe('Teacher Dashboard Logic - Item Styling', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確為個人行程設定樣式', () => {
    const logic = new TeacherDashboardPageLogic()
    const item = { type: 'PERSONAL_EVENT', data: {} }

    expect(logic.getItemBgClass(item)).toContain('purple')
    expect(logic.getItemBorderClass(item)).toContain('purple')
  })

  it('應該正確為有例外的課程設定樣式', () => {
    const logic = new TeacherDashboardPageLogic()
    const item = { type: 'SCHEDULE_RULE', data: { has_exception: true } }

    expect(logic.getItemBgClass(item)).toContain('warning')
    expect(logic.getItemBorderClass(item)).toContain('warning')
  })

  it('應該正確為一般課程設定樣式', () => {
    const logic = new TeacherDashboardPageLogic()
    const item = { type: 'SCHEDULE_RULE', data: {} }

    expect(logic.getItemBgClass(item)).toContain('success')
    expect(logic.getItemBorderClass(item)).toContain('success')
  })
})

describe('Teacher Dashboard Logic - Status Display', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確顯示待審核狀態', () => {
    const logic = new TeacherDashboardPageLogic()
    expect(logic.getStatusText('PENDING')).toBe('待審核')
    expect(logic.getStatusClass('PENDING')).toContain('warning')
  })

  it('應該正確顯示已核准狀態', () => {
    const logic = new TeacherDashboardPageLogic()
    expect(logic.getStatusText('APPROVED')).toBe('已核准')
    expect(logic.getStatusClass('APPROVED')).toContain('success')
  })

  it('應該正確顯示已拒絕狀態', () => {
    const logic = new TeacherDashboardPageLogic()
    expect(logic.getStatusText('REJECTED')).toBe('已拒絕')
    expect(logic.getStatusClass('REJECTED')).toContain('critical')
  })

  it('應該正確處理未知狀態', () => {
    const logic = new TeacherDashboardPageLogic()
    expect(logic.getStatusText('UNKNOWN')).toBe('UNKNOWN')
    expect(logic.getStatusClass('UNKNOWN')).toContain('slate')
  })
})

describe('Teacher Dashboard Logic - Drag and Drop', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確開始拖曳', () => {
    const logic = new TeacherDashboardPageLogic()
    const item = { id: 1, start_time: '10:00', end_time: '11:00' }

    logic.handleDragStart(item, 10, '2026-01-20')

    expect(logic.isDragging.value).toBe(true)
    expect(logic.draggedItem.value).toEqual(item)
    expect(logic.sourceHour.value).toBe(10)
    expect(logic.sourceDate.value).toBe('2026-01-20')
  })

  it('應該正確結束拖曳', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.isDragging.value = true
    logic.dragTarget.value = { time: 10, date: '2026-01-20' }
    logic.draggedItem.value = { id: 1 }

    logic.handleDragEnd()

    expect(logic.isDragging.value).toBe(false)
    expect(logic.dragTarget.value).toBeNull()
    expect(logic.draggedItem.value).toBeNull()
  })

  it('應該正確處理拖曳進入', () => {
    const logic = new TeacherDashboardPageLogic()

    logic.handleDragEnter(14, '2026-01-20')

    expect(logic.dragTarget.value).toEqual({ time: 14, date: '2026-01-20' })
  })

  it('應該正確判斷目標儲存格', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.dragTarget.value = { time: 14, date: '2026-01-20' }

    expect(logic.isTargetCell(14, '2026-01-20')).toBe(true)
    expect(logic.isTargetCell(15, '2026-01-20')).toBe(false)
    expect(logic.isTargetCell(14, '2026-01-21')).toBe(false)
  })
})

describe('Teacher Dashboard Logic - Item Detail', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確打開課程詳細資料', () => {
    const logic = new TeacherDashboardPageLogic()
    const item = { type: 'SCHEDULE_RULE', data: { id: 1 }, id: 1 }

    logic.openItemDetail(item)

    expect(logic.showSessionNoteModal.value).toBe(true)
    expect(logic.selectedScheduleItem.value).toEqual(item)
  })

  it('應該正確打開個人行程編輯', () => {
    const logic = new TeacherDashboardPageLogic()
    const eventData = { id: 1, title: '會議' }
    const item = { type: 'PERSONAL_EVENT', data: eventData }

    logic.openItemDetail(item)

    expect(logic.showPersonalEventModal.value).toBe(true)
    expect(logic.editingEvent.value).toEqual(eventData)
  })

  it('應該正確關閉備註 Modal', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.showSessionNoteModal.value = true
    logic.selectedScheduleItem.value = { id: 1 }

    logic.handleNoteModalClose()

    expect(logic.showSessionNoteModal.value).toBe(false)
    expect(logic.selectedScheduleItem.value).toBeNull()
  })
})

describe('Teacher Dashboard Logic - Grid Columns', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該在行動裝置時返回正確的網格類別', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.isMobile.value = true

    expect(logic.gridColsClass).toBe('grid-cols-[50px_repeat(3,1fr)]')
  })

  it('應該在桌面裝置時返回正確的網格類別', () => {
    const logic = new TeacherDashboardPageLogic()
    logic.isMobile.value = false

    expect(logic.gridColsClass).toBe('grid-cols-[80px_repeat(7,1fr)]')
  })
})
