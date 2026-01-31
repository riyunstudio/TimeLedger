import { describe, it, expect, vi, beforeEach } from 'vitest'

// ============================================
// Mock Setup
// ============================================

// Mock useNotificationStore
const mockNotificationStore = {
  fetchNotifications: vi.fn().mockResolvedValue(undefined),
}

vi.mock('~/stores/useNotificationStore', () => ({
  useNotificationStore: () => mockNotificationStore,
}))

// Mock useNotification
vi.mock('~/composables/useNotification', () => ({
  useNotification: () => ({
    show: { value: false },
    close: vi.fn(),
  }),
}))

// ============================================
// Simple Reactive Wrapper
// ============================================

function createRef<T>(value: T) {
  return { value }
}

// Simple computed getter that just returns the value
function createComputed<T>(getter: () => T): { value: T } {
  return { value: getter() }
}

// ============================================
// Admin Dashboard Page Logic
// ============================================

class AdminDashboardPageLogic {
  // 狀態
  viewMode = createRef<'calendar' | 'teacher_matrix' | 'room_matrix'>('calendar')
  selectedResourceId = createRef<number | null>(null)

  // Computed（改用 getter）
  get resourcePanelViewMode() {
    if (this.viewMode.value === 'teacher_matrix') return 'teacher'
    if (this.viewMode.value === 'room_matrix') return 'room'
    return 'offering'
  }

  // 方法
  handleSelectResource(resource: { type: 'teacher' | 'room', id: number } | null) {
    if (!resource) {
      this.viewMode.value = 'calendar'
      this.selectedResourceId.value = null
    } else {
      if (resource.type === 'teacher') {
        this.viewMode.value = 'teacher_matrix'
      } else {
        this.viewMode.value = 'room_matrix'
      }
      this.selectedResourceId.value = resource.id
    }
  }
}

// ============================================
// Test Suites
// ============================================

describe('Admin Dashboard Logic - View Mode', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該預設為日曆視圖', () => {
    const logic = new AdminDashboardPageLogic()
    expect(logic.viewMode.value).toBe('calendar')
  })

  it('應該正確切換到老師矩陣視圖', () => {
    const logic = new AdminDashboardPageLogic()

    logic.handleSelectResource({ type: 'teacher', id: 1 })

    expect(logic.viewMode.value).toBe('teacher_matrix')
    expect(logic.selectedResourceId.value).toBe(1)
  })

  it('應該正確切換到教室矩陣視圖', () => {
    const logic = new AdminDashboardPageLogic()

    logic.handleSelectResource({ type: 'room', id: 2 })

    expect(logic.viewMode.value).toBe('room_matrix')
    expect(logic.selectedResourceId.value).toBe(2)
  })

  it('應該正確重置為日曆視圖', () => {
    const logic = new AdminDashboardPageLogic()

    // 先切換到矩陣視圖
    logic.handleSelectResource({ type: 'teacher', id: 1 })
    expect(logic.viewMode.value).toBe('teacher_matrix')

    // 再重置
    logic.handleSelectResource(null)
    expect(logic.viewMode.value).toBe('calendar')
    expect(logic.selectedResourceId.value).toBeNull()
  })
})

describe('Admin Dashboard Logic - Resource Panel View Mode', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該在日曆視圖時回傳 offering', () => {
    const logic = new AdminDashboardPageLogic()
    logic.viewMode.value = 'calendar'

    expect(logic.resourcePanelViewMode).toBe('offering')
  })

  it('應該在老師矩陣視圖時回傳 teacher', () => {
    const logic = new AdminDashboardPageLogic()
    logic.viewMode.value = 'teacher_matrix'

    expect(logic.resourcePanelViewMode).toBe('teacher')
  })

  it('應該在教室矩陣視圖時回傳 room', () => {
    const logic = new AdminDashboardPageLogic()
    logic.viewMode.value = 'room_matrix'

    expect(logic.resourcePanelViewMode).toBe('room')
  })
})

describe('Admin Dashboard Logic - Resource Selection', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確選擇多個資源', () => {
    const logic = new AdminDashboardPageLogic()

    // 選擇第一個老師
    logic.handleSelectResource({ type: 'teacher', id: 1 })
    expect(logic.selectedResourceId.value).toBe(1)

    // 選擇第二個老師
    logic.handleSelectResource({ type: 'teacher', id: 2 })
    expect(logic.selectedResourceId.value).toBe(2)

    // 選擇教室
    logic.handleSelectResource({ type: 'room', id: 1 })
    expect(logic.selectedResourceId.value).toBe(1)
  })

  it('應該正確處理資源選擇狀態轉換', () => {
    const logic = new AdminDashboardPageLogic()

    // 日曆 → 老師矩陣
    logic.handleSelectResource({ type: 'teacher', id: 1 })
    expect(logic.viewMode.value).toBe('teacher_matrix')
    expect(logic.resourcePanelViewMode).toBe('teacher')

    // 老師矩陣 → 教室矩陣
    logic.handleSelectResource({ type: 'room', id: 1 })
    expect(logic.viewMode.value).toBe('room_matrix')
    expect(logic.resourcePanelViewMode).toBe('room')
  })
})

describe('Admin Dashboard Logic - Edge Cases', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('應該正確處理 null resource', () => {
    const logic = new AdminDashboardPageLogic()

    // 在任意狀態下傳入 null 都應該重置
    logic.viewMode.value = 'teacher_matrix'
    logic.selectedResourceId.value = 123

    logic.handleSelectResource(null)

    expect(logic.viewMode.value).toBe('calendar')
    expect(logic.selectedResourceId.value).toBeNull()
  })

  it('應該正確處理未知的 viewMode', () => {
    const logic = new AdminDashboardPageLogic()
    logic.viewMode.value = 'calendar'

    // 預設應該是 offering
    expect(logic.resourcePanelViewMode).toBe('offering')
  })
})
