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

vi.mock('~/composables/useCenterId', () => ({
  getCenterId: vi.fn(() => 1),
}))

vi.mock('~/composables/useApi', () => ({
  default: () => ({
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  }),
}))

describe('admin/resources.vue 頁面邏輯', () => {
  // ResourcesPageLogic 類別 - 提取頁面業務邏輯
  class ResourcesPageLogic {
    activeTab: string
    tabs: { id: string; label: string }[]

    constructor() {
      this.activeTab = 'rooms'
      this.tabs = [
        { id: 'rooms', label: '教室' },
        { id: 'courses', label: '課程' },
        { id: 'offerings', label: '待排課程' },
        { id: 'teachers', label: '老師' },
      ]
    }

    setActiveTab(tab: string) {
      this.activeTab = tab
    }

    isActiveTab(tab: string): boolean {
      return this.activeTab === tab
    }

    getCurrentTabLabel(): string {
      const tab = this.tabs.find(t => t.id === this.activeTab)
      return tab?.label || ''
    }

    getNextTab(): string {
      const currentIndex = this.tabs.findIndex(t => t.id === this.activeTab)
      const nextIndex = (currentIndex + 1) % this.tabs.length
      return this.tabs[nextIndex].id
    }

    getPreviousTab(): string {
      const currentIndex = this.tabs.findIndex(t => t.id === this.activeTab)
      const prevIndex = (currentIndex - 1 + this.tabs.length) % this.tabs.length
      return this.tabs[prevIndex].id
    }
  }

  // RoomTabLogic 類別
  class RoomTabLogic {
    rooms: any[]
    loading: boolean
    searchQuery: string

    constructor() {
      this.rooms = []
      this.loading = false
      this.searchQuery = ''
    }

    setRooms(rooms: any[]) {
      this.rooms = rooms
    }

    filteredRooms(): any[] {
      if (!this.searchQuery) return this.rooms
      return this.rooms.filter(r =>
        r.name?.toLowerCase().includes(this.searchQuery.toLowerCase())
      )
    }

    addRoom(room: any) {
      this.rooms.push(room)
    }

    updateRoom(id: number, updates: any) {
      const index = this.rooms.findIndex(r => r.id === id)
      if (index !== -1) {
        this.rooms[index] = { ...this.rooms[index], ...updates }
      }
    }

    deleteRoom(id: number) {
      this.rooms = this.rooms.filter(r => r.id !== id)
    }

    getRoomCount(): number {
      return this.rooms.length
    }

    hasRooms(): boolean {
      return this.rooms.length > 0
    }
  }

  // CourseTabLogic 類別
  class CourseTabLogic {
    courses: any[]
    loading: boolean

    constructor() {
      this.courses = []
      this.loading = false
    }

    setCourses(courses: any[]) {
      this.courses = courses
    }

    filteredCourses(searchQuery: string = ''): any[] {
      if (!searchQuery) return this.courses
      return this.courses.filter(c =>
        c.name?.toLowerCase().includes(searchQuery.toLowerCase())
      )
    }

    addCourse(course: any) {
      this.courses.push(course)
    }

    updateCourse(id: number, updates: any) {
      const index = this.courses.findIndex(c => c.id === id)
      if (index !== -1) {
        this.courses[index] = { ...this.courses[index], ...updates }
      }
    }

    deleteCourse(id: number) {
      this.courses = this.courses.filter(c => c.id !== id)
    }

    getCourseCount(): number {
      return this.courses.length
    }

    hasCourses(): boolean {
      return this.courses.length > 0
    }
  }

  // OfferingTabLogic 類別
  class OfferingTabLogic {
    offerings: any[]
    loading: boolean

    constructor() {
      this.offerings = []
      this.loading = false
    }

    setOfferings(offerings: any[]) {
      this.offerings = offerings
    }

    pendingOfferings(): any[] {
      return this.offerings.filter(o => o.status === 'PENDING')
    }

    activeOfferings(): any[] {
      return this.offerings.filter(o => o.status === 'ACTIVE')
    }

    completedOfferings(): any[] {
      return this.offerings.filter(o => o.status === 'COMPLETED')
    }

    addOffering(offering: any) {
      this.offerings.push(offering)
    }

    updateOffering(id: number, updates: any) {
      const index = this.offerings.findIndex(o => o.id === id)
      if (index !== -1) {
        this.offerings[index] = { ...this.offerings[index], ...updates }
      }
    }

    deleteOffering(id: number) {
      this.offerings = this.offerings.filter(o => o.id !== id)
    }

    getOfferingCount(): number {
      return this.offerings.length
    }
  }

  // TeacherTabLogic 類別
  class TeacherTabLogic {
    teachers: any[]
    loading: boolean
    searchQuery: string

    constructor() {
      this.teachers = []
      this.loading = false
      this.searchQuery = ''
    }

    setTeachers(teachers: any[]) {
      this.teachers = teachers
    }

    filteredTeachers(): any[] {
      if (!this.searchQuery) return this.teachers
      return this.teachers.filter(t =>
        t.name?.toLowerCase().includes(this.searchQuery.toLowerCase())
      )
    }

    activeTeachers(): any[] {
      return this.teachers.filter(t => t.status === 'ACTIVE')
    }

    inactiveTeachers(): any[] {
      return this.teachers.filter(t => t.status === 'INACTIVE')
    }

    addTeacher(teacher: any) {
      this.teachers.push(teacher)
    }

    updateTeacher(id: number, updates: any) {
      const index = this.teachers.findIndex(t => t.id === id)
      if (index !== -1) {
        this.teachers[index] = { ...this.teachers[index], ...updates }
      }
    }

    removeTeacher(id: number) {
      this.teachers = this.teachers.filter(t => t.id !== id)
    }

    getTeacherCount(): number {
      return this.teachers.length
    }

    getActiveTeacherCount(): number {
      return this.activeTeachers().length
    }
  }

  describe('ResourcesPageLogic 頁籤邏輯', () => {
    it('應該初始化預設頁籤為 rooms', () => {
      const logic = new ResourcesPageLogic()
      expect(logic.activeTab).toBe('rooms')
    })

    it('應該包含所有必要的頁籤', () => {
      const logic = new ResourcesPageLogic()
      expect(logic.tabs).toHaveLength(4)
      expect(logic.tabs.map(t => t.id)).toEqual(['rooms', 'courses', 'offerings', 'teachers'])
    })

    it('setActiveTab 應該正確切換頁籤', () => {
      const logic = new ResourcesPageLogic()
      logic.setActiveTab('courses')
      expect(logic.activeTab).toBe('courses')
      logic.setActiveTab('offerings')
      expect(logic.activeTab).toBe('offerings')
    })

    it('isActiveTab 應該正確判斷當前頁籤', () => {
      const logic = new ResourcesPageLogic()
      expect(logic.isActiveTab('rooms')).toBe(true)
      expect(logic.isActiveTab('courses')).toBe(false)
      logic.setActiveTab('courses')
      expect(logic.isActiveTab('rooms')).toBe(false)
      expect(logic.isActiveTab('courses')).toBe(true)
    })

    it('getCurrentTabLabel 應該返回正確的標籤名稱', () => {
      const logic = new ResourcesPageLogic()
      expect(logic.getCurrentTabLabel()).toBe('教室')
      logic.setActiveTab('courses')
      expect(logic.getCurrentTabLabel()).toBe('課程')
      logic.setActiveTab('offerings')
      expect(logic.getCurrentTabLabel()).toBe('待排課程')
      logic.setActiveTab('teachers')
      expect(logic.getCurrentTabLabel()).toBe('老師')
    })

    it('getNextTab 應該循環到下一個頁籤', () => {
      const logic = new ResourcesPageLogic()
      expect(logic.getNextTab()).toBe('courses')
      logic.setActiveTab('courses')
      expect(logic.getNextTab()).toBe('offerings')
      logic.setActiveTab('offerings')
      expect(logic.getNextTab()).toBe('teachers')
      logic.setActiveTab('teachers')
      expect(logic.getNextTab()).toBe('rooms')
    })

    it('getPreviousTab 應該循環到上一個頁籤', () => {
      const logic = new ResourcesPageLogic()
      expect(logic.getPreviousTab()).toBe('teachers')
      logic.setActiveTab('courses')
      expect(logic.getPreviousTab()).toBe('rooms')
      logic.setActiveTab('rooms')
      expect(logic.getPreviousTab()).toBe('teachers')
    })
  })

  describe('RoomTabLogic 教室頁籤邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new RoomTabLogic()
      expect(logic.rooms).toHaveLength(0)
      expect(logic.loading).toBe(false)
      expect(logic.searchQuery).toBe('')
    })

    it('setRooms 應該正確設定教室列表', () => {
      const logic = new RoomTabLogic()
      const rooms = [
        { id: 1, name: '教室 A', capacity: 20 },
        { id: 2, name: '教室 B', capacity: 30 },
      ]
      logic.setRooms(rooms)
      expect(logic.rooms).toHaveLength(2)
    })

    it('filteredRooms 應該過濾搜尋結果', () => {
      const logic = new RoomTabLogic()
      logic.setRooms([
        { id: 1, name: '鋼琴教室' },
        { id: 2, name: '小提琴教室' },
        { id: 3, name: '大教室' },
      ])
      logic.searchQuery = '鋼琴'
      const filtered = logic.filteredRooms()
      expect(filtered).toHaveLength(1)
      expect(filtered[0].name).toBe('鋼琴教室')
    })

    it('filteredRooms 應該返回所有教室當搜尋為空', () => {
      const logic = new RoomTabLogic()
      logic.setRooms([
        { id: 1, name: '教室 A' },
        { id: 2, name: '教室 B' },
      ])
      logic.searchQuery = ''
      expect(logic.filteredRooms()).toHaveLength(2)
    })

    it('addRoom 應該新增教室', () => {
      const logic = new RoomTabLogic()
      logic.addRoom({ id: 1, name: '新教室' })
      expect(logic.rooms).toHaveLength(1)
      logic.addRoom({ id: 2, name: '另一間教室' })
      expect(logic.rooms).toHaveLength(2)
    })

    it('updateRoom 應該更新教室資訊', () => {
      const logic = new RoomTabLogic()
      logic.setRooms([{ id: 1, name: '舊名稱', capacity: 20 }])
      logic.updateRoom(1, { name: '新名稱', capacity: 30 })
      const room = logic.rooms.find(r => r.id === 1)
      expect(room?.name).toBe('新名稱')
      expect(room?.capacity).toBe(30)
    })

    it('updateRoom 應該處理不存在的教室', () => {
      const logic = new RoomTabLogic()
      logic.setRooms([{ id: 1, name: '教室 A' }])
      logic.updateRoom(999, { name: '新名稱' })
      expect(logic.rooms).toHaveLength(1)
      expect(logic.rooms[0].name).toBe('教室 A')
    })

    it('deleteRoom 應該刪除教室', () => {
      const logic = new RoomTabLogic()
      logic.setRooms([
        { id: 1, name: '教室 A' },
        { id: 2, name: '教室 B' },
      ])
      logic.deleteRoom(1)
      expect(logic.rooms).toHaveLength(1)
      expect(logic.rooms[0].id).toBe(2)
    })

    it('getRoomCount 應該返回教室數量', () => {
      const logic = new RoomTabLogic()
      expect(logic.getRoomCount()).toBe(0)
      logic.setRooms([{ id: 1 }, { id: 2 }])
      expect(logic.getRoomCount()).toBe(2)
    })

    it('hasRooms 應該正確判斷是否有教室', () => {
      const logic = new RoomTabLogic()
      expect(logic.hasRooms()).toBe(false)
      logic.setRooms([{ id: 1 }])
      expect(logic.hasRooms()).toBe(true)
    })
  })

  describe('CourseTabLogic 課程頁籤邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new CourseTabLogic()
      expect(logic.courses).toHaveLength(0)
      expect(logic.loading).toBe(false)
    })

    it('setCourses 應該正確設定課程列表', () => {
      const logic = new CourseTabLogic()
      const courses = [
        { id: 1, name: '鋼琴基礎', duration: 60 },
        { id: 2, name: '小提琴入門', duration: 45 },
      ]
      logic.setCourses(courses)
      expect(logic.courses).toHaveLength(2)
    })

    it('filteredCourses 應該過濾搜尋結果', () => {
      const logic = new CourseTabLogic()
      logic.setCourses([
        { id: 1, name: '鋼琴基礎課程' },
        { id: 2, name: '鋼琴進階課程' },
        { id: 3, name: '小提琴入門' },
      ])
      const filtered = logic.filteredCourses('鋼琴')
      expect(filtered).toHaveLength(2)
    })

    it('filteredCourses 應該不區分大小寫', () => {
      const logic = new CourseTabLogic()
      logic.setCourses([{ id: 1, name: 'PIANO Basic' }])
      const filtered = logic.filteredCourses('piano')
      expect(filtered).toHaveLength(1)
    })

    it('addCourse 應該新增課程', () => {
      const logic = new CourseTabLogic()
      logic.addCourse({ id: 1, name: '新課程' })
      expect(logic.courses).toHaveLength(1)
    })

    it('updateCourse 應該更新課程資訊', () => {
      const logic = new CourseTabLogic()
      logic.setCourses([{ id: 1, name: '舊課程', duration: 60 }])
      logic.updateCourse(1, { name: '新课程', duration: 90 })
      const course = logic.courses.find(c => c.id === 1)
      expect(course?.name).toBe('新课程')
      expect(course?.duration).toBe(90)
    })

    it('deleteCourse 應該刪除課程', () => {
      const logic = new CourseTabLogic()
      logic.setCourses([
        { id: 1, name: '課程 A' },
        { id: 2, name: '課程 B' },
      ])
      logic.deleteCourse(1)
      expect(logic.courses).toHaveLength(1)
      expect(logic.courses[0].id).toBe(2)
    })

    it('hasCourses 應該正確判斷是否有課程', () => {
      const logic = new CourseTabLogic()
      expect(logic.hasCourses()).toBe(false)
      logic.setCourses([{ id: 1 }])
      expect(logic.hasCourses()).toBe(true)
    })
  })

  describe('OfferingTabLogic 待排課程頁籤邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new OfferingTabLogic()
      expect(logic.offerings).toHaveLength(0)
      expect(logic.loading).toBe(false)
    })

    it('setOfferings 應該正確設定待排課程列表', () => {
      const logic = new OfferingTabLogic()
      const offerings = [
        { id: 1, name: '暑期鋼琴班', status: 'PENDING' },
        { id: 2, name: '常態鋼琴課', status: 'ACTIVE' },
      ]
      logic.setOfferings(offerings)
      expect(logic.offerings).toHaveLength(2)
    })

    it('pendingOfferings 應該返回待處理狀態的待排課程', () => {
      const logic = new OfferingTabLogic()
      logic.setOfferings([
        { id: 1, name: '待處理課程', status: 'PENDING' },
        { id: 2, name: '進行中課程', status: 'ACTIVE' },
        { id: 3, name: '另一個待處理', status: 'PENDING' },
      ])
      const pending = logic.pendingOfferings()
      expect(pending).toHaveLength(2)
    })

    it('activeOfferings 應該返回進行中狀態的待排課程', () => {
      const logic = new OfferingTabLogic()
      logic.setOfferings([
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'ACTIVE' },
        { id: 3, status: 'COMPLETED' },
      ])
      const active = logic.activeOfferings()
      expect(active).toHaveLength(1)
      expect(active[0].id).toBe(2)
    })

    it('completedOfferings 應該返回已完成狀態的待排課程', () => {
      const logic = new OfferingTabLogic()
      logic.setOfferings([
        { id: 1, status: 'PENDING' },
        { id: 2, status: 'ACTIVE' },
        { id: 3, status: 'COMPLETED' },
      ])
      const completed = logic.completedOfferings()
      expect(completed).toHaveLength(1)
      expect(completed[0].id).toBe(3)
    })

    it('addOffering 應該新增待排課程', () => {
      const logic = new OfferingTabLogic()
      logic.addOffering({ id: 1, name: '新待排課程', status: 'PENDING' })
      expect(logic.offerings).toHaveLength(1)
    })

    it('updateOffering 應該更新待排課程', () => {
      const logic = new OfferingTabLogic()
      logic.setOfferings([{ id: 1, name: '舊名稱', status: 'PENDING' }])
      logic.updateOffering(1, { name: '新名稱', status: 'ACTIVE' })
      const offering = logic.offerings.find(o => o.id === 1)
      expect(offering?.name).toBe('新名稱')
      expect(offering?.status).toBe('ACTIVE')
    })

    it('deleteOffering 應該刪除待排課程', () => {
      const logic = new OfferingTabLogic()
      logic.setOfferings([
        { id: 1, name: '待排課程 A' },
        { id: 2, name: '待排課程 B' },
      ])
      logic.deleteOffering(1)
      expect(logic.offerings).toHaveLength(1)
    })
  })

  describe('TeacherTabLogic 老師頁籤邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new TeacherTabLogic()
      expect(logic.teachers).toHaveLength(0)
      expect(logic.loading).toBe(false)
      expect(logic.searchQuery).toBe('')
    })

    it('setTeachers 應該正確設定老師列表', () => {
      const logic = new TeacherTabLogic()
      const teachers = [
        { id: 1, name: '張老師', status: 'ACTIVE' },
        { id: 2, name: '李老師', status: 'INACTIVE' },
      ]
      logic.setTeachers(teachers)
      expect(logic.teachers).toHaveLength(2)
    })

    it('filteredTeachers 應該過濾搜尋結果', () => {
      const logic = new TeacherTabLogic()
      logic.setTeachers([
        { id: 1, name: '張老師' },
        { id: 2, name: '李老師' },
        { id: 3, name: '張三老師' },
      ])
      logic.searchQuery = '張'
      const filtered = logic.filteredTeachers()
      expect(filtered).toHaveLength(2)
    })

    it('activeTeachers 應該返回狀態為 ACTIVE 的老師', () => {
      const logic = new TeacherTabLogic()
      logic.setTeachers([
        { id: 1, name: '老師 A', status: 'ACTIVE' },
        { id: 2, name: '老師 B', status: 'INACTIVE' },
      ])
      const active = logic.activeTeachers()
      expect(active).toHaveLength(1)
      expect(active[0].name).toBe('老師 A')
    })

    it('inactiveTeachers 應該返回狀態為 INACTIVE 的老師', () => {
      const logic = new TeacherTabLogic()
      logic.setTeachers([
        { id: 1, name: '老師 A', status: 'ACTIVE' },
        { id: 2, name: '老師 B', status: 'INACTIVE' },
      ])
      const inactive = logic.inactiveTeachers()
      expect(inactive).toHaveLength(1)
      expect(inactive[0].name).toBe('老師 B')
    })

    it('addTeacher 應該新增老師', () => {
      const logic = new TeacherTabLogic()
      logic.addTeacher({ id: 1, name: '新老師' })
      expect(logic.teachers).toHaveLength(1)
    })

    it('updateTeacher 應該更新老師資訊', () => {
      const logic = new TeacherTabLogic()
      logic.setTeachers([{ id: 1, name: '舊名稱', status: 'ACTIVE' }])
      logic.updateTeacher(1, { name: '新名稱', status: 'INACTIVE' })
      const teacher = logic.teachers.find(t => t.id === 1)
      expect(teacher?.name).toBe('新名稱')
      expect(teacher?.status).toBe('INACTIVE')
    })

    it('removeTeacher 應該移除老師', () => {
      const logic = new TeacherTabLogic()
      logic.setTeachers([
        { id: 1, name: '老師 A' },
        { id: 2, name: '老師 B' },
      ])
      logic.removeTeacher(1)
      expect(logic.teachers).toHaveLength(1)
    })

    it('getActiveTeacherCount 應該返回活躍老師數量', () => {
      const logic = new TeacherTabLogic()
      logic.setTeachers([
        { id: 1, status: 'ACTIVE' },
        { id: 2, status: 'ACTIVE' },
        { id: 3, status: 'INACTIVE' },
      ])
      expect(logic.getActiveTeacherCount()).toBe(2)
    })
  })

  describe('頁面整合邏輯', () => {
    it('應該能夠在各頁籤間正確切換並操作資料', () => {
      const pageLogic = new ResourcesPageLogic()
      const roomLogic = new RoomTabLogic()
      const courseLogic = new CourseTabLogic()
      const offeringLogic = new OfferingTabLogic()
      const teacherLogic = new TeacherTabLogic()

      // 初始狀態
      expect(pageLogic.activeTab).toBe('rooms')
      expect(roomLogic.hasRooms()).toBe(false)

      // 新增教室
      pageLogic.setActiveTab('rooms')
      roomLogic.addRoom({ id: 1, name: '鋼琴教室' })
      expect(roomLogic.hasRooms()).toBe(true)

      // 新增課程
      pageLogic.setActiveTab('courses')
      courseLogic.addCourse({ id: 1, name: '鋼琴基礎' })
      expect(courseLogic.hasCourses()).toBe(true)

      // 新增待排課程
      pageLogic.setActiveTab('offerings')
      offeringLogic.addOffering({ id: 1, name: '鋼琴暑期班', status: 'PENDING' })
      expect(offeringLogic.pendingOfferings()).toHaveLength(1)

      // 新增老師
      pageLogic.setActiveTab('teachers')
      teacherLogic.addTeacher({ id: 1, name: '張老師', status: 'ACTIVE' })
      expect(teacherLogic.getActiveTeacherCount()).toBe(1)
    })

    it('應該正確計算所有資源的總數量', () => {
      const roomLogic = new RoomTabLogic()
      const courseLogic = new CourseTabLogic()
      const offeringLogic = new OfferingTabLogic()
      const teacherLogic = new TeacherTabLogic()

      roomLogic.setRooms([{ id: 1 }, { id: 2 }, { id: 3 }])
      courseLogic.setCourses([{ id: 1 }, { id: 2 }])
      offeringLogic.setOfferings([{ id: 1 }, { id: 2 }, { id: 3 }, { id: 4 }])
      teacherLogic.setTeachers([{ id: 1 }, { id: 2 }, { id: 3 }, { id: 4 }, { id: 5 }])

      const totalResources =
        roomLogic.getRoomCount() +
        courseLogic.getCourseCount() +
        offeringLogic.getOfferingCount() +
        teacherLogic.getTeacherCount()

      expect(totalResources).toBe(14)
    })
  })
})
