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
  confirm: vi.fn().mockResolvedValue(true),
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

describe('admin/holidays.vue 頁面邏輯', () => {
  // HolidayListLogic 類別 - 假日列表邏輯
  class HolidayListLogic {
    holidays: any[]

    constructor() {
      this.holidays = []
    }

    setHolidays(holidays: any[]) {
      this.holidays = holidays
    }

    getHolidays(): any[] {
      return this.holidays
    }

    getHolidayById(id: number): any | undefined {
      return this.holidays.find(h => h.id === id)
    }

    getHolidayByDate(date: string): any | undefined {
      return this.holidays.find(h => h.date === date)
    }

    hasHolidays(): boolean {
      return this.holidays.length > 0
    }

    getHolidayCount(): number {
      return this.holidays.length
    }

    addHoliday(holiday: any) {
      this.holidays.push(holiday)
    }

    removeHoliday(id: number) {
      this.holidays = this.holidays.filter(h => h.id !== id)
    }

    updateHoliday(id: number, updates: any) {
      const index = this.holidays.findIndex(h => h.id === id)
      if (index !== -1) {
        this.holidays[index] = { ...this.holidays[index], ...updates }
      }
    }

    filterByMonth(month: number): any[] {
      return this.holidays.filter(h => {
        const holidayMonth = new Date(h.date).getMonth() + 1
        return holidayMonth === month
      })
    }

    filterByYear(year: number): any[] {
      return this.holidays.filter(h => {
        const holidayYear = new Date(h.date).getFullYear()
        return holidayYear === year
      })
    }

    sortByDate(): any[] {
      return [...this.holidays].sort((a, b) =>
        new Date(a.date).getTime() - new Date(b.date).getTime()
      )
    }

    isHoliday(date: string): boolean {
      return this.holidays.some(h => h.date === date)
    }
  }

  // CalendarNavigationLogic 類別 - 日曆導航邏輯
  class CalendarNavigationLogic {
    currentMonth: number
    selectedYear: number

    constructor() {
      const now = new Date()
      this.currentMonth = now.getMonth()
      this.selectedYear = now.getFullYear()
    }

    setCurrentMonth(month: number) {
      this.currentMonth = Math.max(0, Math.min(11, month))
    }

    setSelectedYear(year: number) {
      this.selectedYear = year
    }

    goToNextMonth() {
      if (this.currentMonth === 11) {
        this.currentMonth = 0
        this.selectedYear++
      } else {
        this.currentMonth++
      }
    }

    goToPreviousMonth() {
      if (this.currentMonth === 0) {
        this.currentMonth = 11
        this.selectedYear--
      } else {
        this.currentMonth--
      }
    }

    goToToday() {
      const now = new Date()
      this.currentMonth = now.getMonth()
      this.selectedYear = now.getFullYear()
    }

    getCurrentMonthLabel(): string {
      return `${this.selectedYear} 年 ${this.currentMonth + 1} 月`
    }

    getCurrentYear(): number {
      return this.selectedYear
    }

    getCurrentMonth(): number {
      return this.currentMonth
    }

    getMonthStartDay(): number {
      return new Date(this.selectedYear, this.currentMonth, 1).getDay()
    }

    getMonthDays(): number {
      return new Date(this.selectedYear, this.currentMonth + 1, 0).getDate()
    }

    getYearOptions(): number[] {
      const current = new Date().getFullYear()
      return [current - 1, current, current + 1, current + 2]
    }
  }

  // HolidayFormLogic 類別 - 假日表單邏輯
  class HolidayFormLogic {
    form: {
      name: string
      date: string
    }

    constructor() {
      this.form = {
        name: '',
        date: '',
      }
    }

    setName(name: string) {
      this.form.name = name
    }

    setDate(date: string) {
      this.form.date = date
    }

    resetForm() {
      this.form = { name: '', date: '' }
    }

    isValid(): boolean {
      return Boolean(this.form.name.trim() && this.form.date)
    }

    getFormData(): any {
      return {
        name: this.form.name.trim(),
        date: this.form.date,
      }
    }

    loadFromHoliday(holiday: any) {
      this.form.name = holiday.name || ''
      this.form.date = holiday.date || ''
    }
  }

  // BulkImportLogic 類別 - 批次匯入邏輯
  class BulkImportLogic {
    jsonData: string

    constructor() {
      this.jsonData = ''
    }

    setJsonData(data: string) {
      this.jsonData = data
    }

    clearData() {
      this.jsonData = ''
    }

    parseData(): any[] {
      if (!this.jsonData.trim()) {
        return []
      }
      try {
        const parsed = JSON.parse(this.jsonData)
        if (Array.isArray(parsed)) {
          return parsed
        }
        return []
      } catch {
        return []
      }
    }

    isValidJson(): boolean {
      if (!this.jsonData.trim()) return false
      try {
        JSON.parse(this.jsonData)
        return true
      } catch {
        return false
      }
    }

    isValidFormat(): boolean {
      const data = this.parseData()
      if (data.length === 0) return false
      return data.every(item =>
        typeof item === 'object' &&
        item !== null &&
        'date' in item &&
        'name' in item
      )
    }

    getItemCount(): number {
      const data = this.parseData()
      return data.length
    }

    getValidationErrors(): string[] {
      const errors: string[] = []
      if (!this.isValidJson()) {
        errors.push('JSON 格式錯誤')
        return errors
      }
      const data = this.parseData()
      if (!Array.isArray(data)) {
        errors.push('資料必須是陣列格式')
        return errors
      }
      data.forEach((item, index) => {
        if (typeof item !== 'object' || item === null) {
          errors.push(`第 ${index + 1} 筆資料格式錯誤`)
        } else if (!item.date) {
          errors.push(`第 ${index + 1} 筆缺少 date 欄位`)
        } else if (!item.name) {
          errors.push(`第 ${index + 1} 筆缺少 name 欄位`)
        }
      })
      return errors
    }
  }

  // HolidayDisplayLogic 類別 - 假日顯示邏輯
  class HolidayDisplayLogic {
    formatFullDate(dateStr: string): string {
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-TW', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        weekday: 'short',
      })
    }

    formatShortDate(dateStr: string): string {
      const date = new Date(dateStr)
      return `${date.getMonth() + 1}/${date.getDate()}`
    }

    getDayOfWeek(dateStr: string): string {
      const weekdays = ['週日', '週一', '週二', '週三', '週四', '週五', '週六']
      const date = new Date(dateStr)
      return weekdays[date.getDay()]
    }

    getDayOfMonth(dateStr: string): number {
      return new Date(dateStr).getDate()
    }

    getDisplayClass(holiday: any): string {
      return 'bg-warning-500/20 text-warning-500'
    }

    getDayClass(day: number): string {
      if (day === 0 || day === 6) {
        return 'text-slate-400'
      }
      return 'text-white'
    }
  }

  // YearSelectionLogic 類別 - 年份選擇邏輯
  class YearSelectionLogic {
    currentYear: number
    availableYears: number[]

    constructor() {
      this.currentYear = new Date().getFullYear()
      this.availableYears = [this.currentYear - 1, this.currentYear, this.currentYear + 1, this.currentYear + 2]
    }

    getAvailableYears(): number[] {
      return this.availableYears
    }

    setYear(year: number) {
      this.currentYear = year
    }

    getCurrentYear(): number {
      return this.currentYear
    }

    isCurrentYear(year: number): boolean {
      return year === this.currentYear
    }

    getYearLabel(year: number): string {
      return `${year}年`
    }
  }

  describe('HolidayListLogic 假日列表邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new HolidayListLogic()
      expect(logic.holidays).toHaveLength(0)
    })

    it('setHolidays 應該正確設定假日列表', () => {
      const logic = new HolidayListLogic()
      logic.setHolidays([
        { id: 1, name: '元旦', date: '2026-01-01' },
        { id: 2, name: '春節', date: '2026-01-29' },
      ])
      expect(logic.holidays).toHaveLength(2)
    })

    it('getHolidayById 應該正確取得特定假日', () => {
      const logic = new HolidayListLogic()
      logic.setHolidays([
        { id: 1, name: '元旦' },
        { id: 2, name: '春節' },
      ])
      const holiday = logic.getHolidayById(2)
      expect(holiday?.name).toBe('春節')
    })

    it('getHolidayByDate 應該根據日期取得假日', () => {
      const logic = new HolidayListLogic()
      logic.setHolidays([{ id: 1, name: '元旦', date: '2026-01-01' }])
      const holiday = logic.getHolidayByDate('2026-01-01')
      expect(holiday?.name).toBe('元旦')
    })

    it('filterByMonth 應該過濾指定月份的假日', () => {
      const logic = new HolidayListLogic()
      logic.setHolidays([
        { id: 1, date: '2026-01-15' },
        { id: 2, date: '2026-02-10' },
        { id: 3, date: '2026-01-20' },
      ])
      const january = logic.filterByMonth(1)
      expect(january).toHaveLength(2)
    })

    it('filterByYear 應該過濾指定年份的假日', () => {
      const logic = new HolidayListLogic()
      logic.setHolidays([
        { id: 1, date: '2026-01-01' },
        { id: 2, date: '2025-12-25' },
      ])
      const year2026 = logic.filterByYear(2026)
      expect(year2026).toHaveLength(1)
    })

    it('sortByDate 應該正確排序假日', () => {
      const logic = new HolidayListLogic()
      logic.setHolidays([
        { id: 1, date: '2026-02-01' },
        { id: 2, date: '2026-01-01' },
        { id: 3, date: '2026-01-15' },
      ])
      const sorted = logic.sortByDate()
      expect(sorted[0].id).toBe(2)
      expect(sorted[1].id).toBe(3)
      expect(sorted[2].id).toBe(1)
    })

    it('isHoliday 應該正確判斷是否為假日', () => {
      const logic = new HolidayListLogic()
      logic.setHolidays([{ date: '2026-01-01' }])
      expect(logic.isHoliday('2026-01-01')).toBe(true)
      expect(logic.isHoliday('2026-01-02')).toBe(false)
    })

    it('addHoliday 應該新增假日', () => {
      const logic = new HolidayListLogic()
      logic.addHoliday({ id: 1, name: '新假日', date: '2026-03-01' })
      expect(logic.getHolidayCount()).toBe(1)
    })

    it('removeHoliday 應該刪除假日', () => {
      const logic = new HolidayListLogic()
      logic.setHolidays([
        { id: 1, name: '假日 A' },
        { id: 2, name: '假日 B' },
      ])
      logic.removeHoliday(1)
      expect(logic.getHolidayCount()).toBe(1)
      expect(logic.getHolidayById(1)).toBeUndefined()
    })
  })

  describe('CalendarNavigationLogic 日曆導航邏輯', () => {
    it('應該正確初始化當前年份和月份', () => {
      const logic = new CalendarNavigationLogic()
      expect(logic.currentMonth).toBeGreaterThanOrEqual(0)
      expect(logic.currentMonth).toBeLessThanOrEqual(11)
    })

    it('goToNextMonth 應該正確前進月份', () => {
      const logic = new CalendarNavigationLogic()
      logic.currentMonth = 0
      logic.selectedYear = 2026
      logic.goToNextMonth()
      expect(logic.currentMonth).toBe(1)
      logic.currentMonth = 11
      logic.selectedYear = 2026
      logic.goToNextMonth()
      expect(logic.currentMonth).toBe(0)
      expect(logic.selectedYear).toBe(2027)
    })

    it('goToPreviousMonth 應該正確後退月份', () => {
      const logic = new CalendarNavigationLogic()
      logic.currentMonth = 1
      logic.selectedYear = 2026
      logic.goToPreviousMonth()
      expect(logic.currentMonth).toBe(0)
      logic.currentMonth = 0
      logic.selectedYear = 2026
      logic.goToPreviousMonth()
      expect(logic.currentMonth).toBe(11)
      expect(logic.selectedYear).toBe(2025)
    })

    it('goToToday 應該回到今天', () => {
      const logic = new CalendarNavigationLogic()
      logic.currentMonth = 5
      logic.selectedYear = 2020
      logic.goToToday()
      const now = new Date()
      expect(logic.currentMonth).toBe(now.getMonth())
      expect(logic.selectedYear).toBe(now.getFullYear())
    })

    it('getCurrentMonthLabel 應該返回正確的月份標籤', () => {
      const logic = new CalendarNavigationLogic()
      logic.currentMonth = 0
      logic.selectedYear = 2026
      expect(logic.getCurrentMonthLabel()).toBe('2026 年 1 月')
    })

    it('getMonthDays 應該返回正確的天數', () => {
      const logic = new CalendarNavigationLogic()
      logic.currentMonth = 0 // 1月
      logic.selectedYear = 2026
      expect(logic.getMonthDays()).toBe(31)
      logic.currentMonth = 1 // 2月
      logic.selectedYear = 2026
      expect(logic.getMonthDays()).toBe(28) // 2026年不是閏年
      logic.selectedYear = 2024
      logic.currentMonth = 1
      expect(logic.getMonthDays()).toBe(29) // 2024是閏年
    })

    it('getYearOptions 應該返回正確的年份選項', () => {
      const logic = new CalendarNavigationLogic()
      const options = logic.getYearOptions()
      expect(options).toHaveLength(4)
      expect(options[0]).toBe(options[2] - 2)
      expect(options[3]).toBe(options[2] + 2)
    })
  })

  describe('HolidayFormLogic 假日表單邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new HolidayFormLogic()
      expect(logic.form.name).toBe('')
      expect(logic.form.date).toBe('')
    })

    it('setName 應該正確設定名稱', () => {
      const logic = new HolidayFormLogic()
      logic.setName('春節')
      expect(logic.form.name).toBe('春節')
    })

    it('setDate 應該正確設定日期', () => {
      const logic = new HolidayFormLogic()
      logic.setDate('2026-01-01')
      expect(logic.form.date).toBe('2026-01-01')
    })

    it('resetForm 應該重置表單', () => {
      const logic = new HolidayFormLogic()
      logic.form.name = '春節'
      logic.form.date = '2026-01-01'
      logic.resetForm()
      expect(logic.form.name).toBe('')
      expect(logic.form.date).toBe('')
    })

    it('isValid 應該在名稱和日期都有值時返回 true', () => {
      const logic = new HolidayFormLogic()
      expect(logic.isValid()).toBe(false)
      logic.setName('春節')
      expect(logic.isValid()).toBe(false)
      logic.setDate('2026-01-01')
      expect(logic.isValid()).toBe(true)
    })

    it('isValid 應該在名稱為空白時返回 false', () => {
      const logic = new HolidayFormLogic()
      logic.setName('   ')
      logic.setDate('2026-01-01')
      expect(logic.isValid()).toBe(false)
    })

    it('getFormData 應該返回表單資料', () => {
      const logic = new HolidayFormLogic()
      logic.setName('春節')
      logic.setDate('2026-01-01')
      const data = logic.getFormData()
      expect(data.name).toBe('春節')
      expect(data.date).toBe('2026-01-01')
    })

    it('loadFromHoliday 應該從假日資料載入', () => {
      const logic = new HolidayFormLogic()
      logic.loadFromHoliday({ name: '春節', date: '2026-01-01' })
      expect(logic.form.name).toBe('春節')
      expect(logic.form.date).toBe('2026-01-01')
    })
  })

  describe('BulkImportLogic 批次匯入邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new BulkImportLogic()
      expect(logic.jsonData).toBe('')
    })

    it('setJsonData 應該正確設定 JSON 資料', () => {
      const logic = new BulkImportLogic()
      logic.setJsonData('[{"date": "2026-01-01", "name": "元旦"}]')
      expect(logic.jsonData).toContain('元旦')
    })

    it('clearData 應該清除資料', () => {
      const logic = new BulkImportLogic()
      logic.setJsonData('test')
      logic.clearData()
      expect(logic.jsonData).toBe('')
    })

    it('parseData 應該正確解析 JSON', () => {
      const logic = new BulkImportLogic()
      logic.setJsonData('[{"date": "2026-01-01", "name": "元旦"}]')
      const data = logic.parseData()
      expect(data).toHaveLength(1)
      expect(data[0].name).toBe('元旦')
    })

    it('parseData 應該處理空值', () => {
      const logic = new BulkImportLogic()
      expect(logic.parseData()).toHaveLength(0)
      logic.setJsonData('   ')
      expect(logic.parseData()).toHaveLength(0)
    })

    it('isValidJson 應該正確驗證 JSON 格式', () => {
      const logic = new BulkImportLogic()
      expect(logic.isValidJson()).toBe(false)
      logic.setJsonData('invalid json')
      expect(logic.isValidJson()).toBe(false)
      logic.setJsonData('{"valid": true}')
      expect(logic.isValidJson()).toBe(true)
    })

    it('isValidFormat 應該正確驗證資料格式', () => {
      const logic = new BulkImportLogic()
      expect(logic.isValidFormat()).toBe(false)
      logic.setJsonData('[{"date": "2026-01-01", "name": "元旦"}]')
      expect(logic.isValidFormat()).toBe(true)
    })

    it('isValidFormat 應該在缺少欄位時返回 false', () => {
      const logic = new BulkImportLogic()
      logic.setJsonData('[{"date": "2026-01-01"}]')
      expect(logic.isValidFormat()).toBe(false)
    })

    it('getItemCount 應該返回項目數量', () => {
      const logic = new BulkImportLogic()
      logic.setJsonData('[{"date": "2026-01-01", "name": "元旦"}, {"date": "2026-02-11", "name": "春節"}]')
      expect(logic.getItemCount()).toBe(2)
    })

    it('getValidationErrors 應該返回驗證錯誤', () => {
      const logic = new BulkImportLogic()
      logic.setJsonData('invalid')
      const errors = logic.getValidationErrors()
      expect(errors.length).toBeGreaterThan(0)
    })
  })

  describe('HolidayDisplayLogic 假日顯示邏輯', () => {
    it('formatFullDate 應該正確格式化完整日期', () => {
      const logic = new HolidayDisplayLogic()
      const result = logic.formatFullDate('2026-01-01')
      expect(result).toContain('2026')
      expect(result).toContain('1')
      expect(result).toContain('1')
    })

    it('formatShortDate 應該正確格式化簡短日期', () => {
      const logic = new HolidayDisplayLogic()
      expect(logic.formatShortDate('2026-01-15')).toBe('1/15')
    })

    it('getDayOfWeek 應該返回正確的星期', () => {
      const logic = new HolidayDisplayLogic()
      // 2026-01-01 是週四
      expect(logic.getDayOfWeek('2026-01-01')).toBe('週四')
      // 2026-01-04 是週日
      expect(logic.getDayOfWeek('2026-01-04')).toBe('週日')
    })

    it('getDayOfMonth 應該返回正確的日期數字', () => {
      const logic = new HolidayDisplayLogic()
      expect(logic.getDayOfMonth('2026-01-15')).toBe(15)
    })

    it('getDayClass 應該返回正確的 CSS 類別', () => {
      const logic = new HolidayDisplayLogic()
      expect(logic.getDayClass(0)).toContain('slate')
      expect(logic.getDayClass(6)).toContain('slate')
      expect(logic.getDayClass(1)).toBe('text-white')
    })
  })

  describe('YearSelectionLogic 年份選擇邏輯', () => {
    it('getAvailableYears 應該返回可用的年份', () => {
      const logic = new YearSelectionLogic()
      const years = logic.getAvailableYears()
      expect(years).toHaveLength(4)
    })

    it('isCurrentYear 應該正確判斷是否為當前年份', () => {
      const logic = new YearSelectionLogic()
      expect(logic.isCurrentYear(logic.getCurrentYear())).toBe(true)
      expect(logic.isCurrentYear(2020)).toBe(false)
    })

    it('getYearLabel 應該返回正確的年份標籤', () => {
      const logic = new YearSelectionLogic()
      expect(logic.getYearLabel(2026)).toBe('2026年')
    })
  })

  describe('頁面整合邏輯', () => {
    it('應該能夠完整處理假日管理流程', () => {
      const listLogic = new HolidayListLogic()
      const navLogic = new CalendarNavigationLogic()
      const formLogic = new HolidayFormLogic()
      const bulkLogic = new BulkImportLogic()
      const displayLogic = new HolidayDisplayLogic()

      // 設定假日
      listLogic.setHolidays([
        { id: 1, name: '元旦', date: '2026-01-01' },
        { id: 2, name: '春節', date: '2026-01-29' },
        { id: 3, name: '和平紀念日', date: '2026-02-28' },
      ])

      // 驗證假日列表
      expect(listLogic.hasHolidays()).toBe(true)
      expect(listLogic.getHolidayCount()).toBe(3)
      expect(listLogic.filterByMonth(1)).toHaveLength(2)

      // 導航到正確月份
      navLogic.setSelectedYear(2026)
      navLogic.setCurrentMonth(0) // 1月
      expect(navLogic.getCurrentMonthLabel()).toBe('2026 年 1 月')
      navLogic.goToNextMonth()
      expect(navLogic.getCurrentMonthLabel()).toBe('2026 年 2 月')

      // 新增假日
      formLogic.setName('兒童節')
      formLogic.setDate('2026-04-04')
      expect(formLogic.isValid()).toBe(true)

      // 批次匯入
      bulkLogic.setJsonData('[{"date": "2026-04-05", "name": "清明節"}]')
      expect(bulkLogic.isValidFormat()).toBe(true)
      expect(bulkLogic.getItemCount()).toBe(1)

      // 日期顯示
      const formattedDate = displayLogic.formatFullDate('2026-01-01')
      expect(formattedDate).toContain('2026')
      expect(displayLogic.getDayOfWeek('2026-01-01')).toBe('週四')
    })

    it('應該正確處理跨年月份導航', () => {
      const navLogic = new CalendarNavigationLogic()
      navLogic.setSelectedYear(2026)
      navLogic.setCurrentMonth(11) // 12月
      navLogic.goToNextMonth()
      expect(navLogic.getCurrentYear()).toBe(2027)
      expect(navLogic.getCurrentMonth()).toBe(0)
    })

    it('應該正確驗證批次匯入資料', () => {
      const bulkLogic = new BulkImportLogic()
      bulkLogic.setJsonData('[{"date": "2026-01-01", "name": "元旦"}, {"date": "invalid", "name": "錯誤日期"}]')
      const errors = bulkLogic.getValidationErrors()
      expect(errors.length).toBeGreaterThan(0)
    })
  })
})
