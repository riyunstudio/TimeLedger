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
  alertWarning: vi.fn().mockResolvedValue(true),
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

describe('admin/templates.vue 頁面邏輯', () => {
  // TemplateListLogic 類別 - 模板列表邏輯
  class TemplateListLogic {
    templates: any[]

    constructor() {
      this.templates = []
    }

    setTemplates(templates: any[]) {
      this.templates = templates
    }

    getTemplates(): any[] {
      return this.templates
    }

    getTemplateById(id: number): any | undefined {
      return this.templates.find(t => t.id === id)
    }

    getActiveTemplates(): any[] {
      return this.templates.filter(t => t.is_active !== false)
    }

    getInactiveTemplates(): any[] {
      return this.templates.filter(t => t.is_active === false)
    }

    getRoomTemplates(): any[] {
      return this.templates.filter(t => t.row_type === 'ROOM')
    }

    getTeacherTemplates(): any[] {
      return this.templates.filter(t => t.row_type === 'TEACHER')
    }

    hasTemplates(): boolean {
      return this.templates.length > 0
    }

    getTemplateCount(): number {
      return this.templates.length
    }

    addTemplate(template: any) {
      this.templates.push(template)
    }

    removeTemplate(id: number) {
      this.templates = this.templates.filter(t => t.id !== id)
    }

    updateTemplate(id: number, updates: any) {
      const index = this.templates.findIndex(t => t.id === id)
      if (index !== -1) {
        this.templates[index] = { ...this.templates[index], ...updates }
      }
    }
  }

  // TemplateTypeLogic 類別 - 模板類型邏輯
  class TemplateTypeLogic {
    getTypeLabel(rowType: string): string {
      switch (rowType) {
        case 'ROOM':
          return '教室視角'
        case 'TEACHER':
          return '老師視角'
        default:
          return rowType
      }
    }

    getTypeClass(rowType: string): string {
      switch (rowType) {
        case 'ROOM':
          return 'bg-primary-500/20 text-primary-500'
        case 'TEACHER':
          return 'bg-secondary-500/20 text-secondary-500'
        default:
          return 'bg-slate-500/20 text-slate-400'
      }
    }

    getTypeOptions(): { value: string; label: string }[] {
      return [
        { value: 'ROOM', label: '教室視角' },
        { value: 'TEACHER', label: '老師視角' },
      ]
    }

    isValidType(type: string): boolean {
      return type === 'ROOM' || type === 'TEACHER'
    }
  }

  // TemplateFormLogic 類別 - 模板表單邏輯
  class TemplateFormLogic {
    form: {
      name: string
      row_type: string
    }

    constructor() {
      this.form = {
        name: '',
        row_type: 'ROOM',
      }
    }

    setName(name: string) {
      this.form.name = name
    }

    setRowType(rowType: string) {
      this.form.row_type = rowType
    }

    resetForm() {
      this.form = { name: '', row_type: 'ROOM' }
    }

    isValid(): boolean {
      return Boolean(this.form.name && this.form.name.trim())
    }

    getFormData(): any {
      return {
        name: this.form.name.trim(),
        row_type: this.form.row_type,
      }
    }
  }

  // WeekdaySelectionLogic 類別 - 星期選擇邏輯
  class WeekdaySelectionLogic {
    weekdays: { value: number; label: string }[]

    constructor() {
      this.weekdays = [
        { value: 1, label: '週一' },
        { value: 2, label: '週二' },
        { value: 3, label: '週三' },
        { value: 4, label: '週四' },
        { value: 5, label: '週五' },
        { value: 6, label: '週六' },
        { value: 7, label: '週日' },
      ]
    }

    getWeekdays(): { value: number; label: string }[] {
      return this.weekdays
    }

    getWeekdayLabel(value: number): string {
      const weekday = this.weekdays.find(w => w.value === value)
      return weekday?.label || String(value)
    }

    toggleWeekday(selected: number[], value: number): number[] {
      const index = selected.indexOf(value)
      if (index === -1) {
        return [...selected, value]
      } else {
        return selected.filter(v => v !== value)
      }
    }

    isWeekdaySelected(selected: number[], value: number): boolean {
      return selected.includes(value)
    }

    getSelectedCount(selected: number[]): number {
      return selected.length
    }

    hasWeekdaysSelected(selected: number[]): boolean {
      return selected.length > 0
    }

    sortWeekdays(values: number[]): number[] {
      return [...values].sort((a, b) => a - b)
    }

    getWeekdayRange(values: number[]): string {
      if (values.length === 0) return ''
      const sorted = this.sortWeekdays(values)
      const labels = sorted.map(v => this.getWeekdayLabel(v))
      return labels.join('、')
    }
  }

  // ApplyFormLogic 類別 - 套用表單邏輯
  class ApplyFormLogic {
    form: {
      templateId: number
      templateName: string
      offeringId: string
      startDate: string
      endDate: string
      weekdays: number[]
      duration: number
    }

    constructor() {
      this.form = {
        templateId: 0,
        templateName: '',
        offeringId: '',
        startDate: '',
        endDate: '',
        weekdays: [],
        duration: 60,
      }
    }

    setTemplate(template: any) {
      this.form.templateId = template.id
      this.form.templateName = template.name
    }

    setOfferingId(offeringId: string) {
      this.form.offeringId = offeringId
    }

    setDateRange(start: string, end: string) {
      this.form.startDate = start
      this.form.endDate = end
    }

    setWeekdays(weekdays: number[]) {
      this.form.weekdays = weekdays
    }

    setDuration(duration: number) {
      this.form.duration = duration
    }

    resetForm() {
      this.form = {
        templateId: 0,
        templateName: '',
        offeringId: '',
        startDate: '',
        endDate: '',
        weekdays: [],
        duration: 60,
      }
    }

    isValid(): boolean {
      return Boolean(
        this.form.offeringId &&
        this.form.startDate &&
        this.form.endDate &&
        this.form.weekdays.length > 0
      )
    }

    getApplyData(): any {
      return {
        offering_id: Number(this.form.offeringId),
        start_date: this.form.startDate,
        end_date: this.form.endDate,
        weekdays: this.form.weekdays,
        duration: this.form.duration,
      }
    }
  }

  // TemplateCellLogic 類別 - 模板格子邏輯
  class TemplateCellLogic {
    cells: any[]

    constructor() {
      this.cells = []
    }

    setCells(cells: any[]) {
      this.cells = cells
    }

    getCells(): any[] {
      return this.cells
    }

    hasCells(): boolean {
      return this.cells.length > 0
    }

    getCellCount(): number {
      return this.cells.length
    }

    getCellById(id: number): any | undefined {
      return this.cells.find(c => c.id === id)
    }

    addCell(cell: any) {
      this.cells.push(cell)
    }

    removeCell(id: number) {
      this.cells = this.cells.filter(c => c.id !== id)
    }

    updateCell(id: number, updates: any) {
      const index = this.cells.findIndex(c => c.id === id)
      if (index !== -1) {
        this.cells[index] = { ...this.cells[index], ...updates }
      }
    }

    sortCells(): any[] {
      return [...this.cells].sort((a, b) => {
        if (a.row_no !== b.row_no) return a.row_no - b.row_no
        return a.col_no - b.col_no
      })
    }
  }

  describe('TemplateListLogic 模板列表邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new TemplateListLogic()
      expect(logic.templates).toHaveLength(0)
    })

    it('setTemplates 應該正確設定模板列表', () => {
      const logic = new TemplateListLogic()
      logic.setTemplates([
        { id: 1, name: '模板 A' },
        { id: 2, name: '模板 B' },
      ])
      expect(logic.templates).toHaveLength(2)
    })

    it('getTemplateById 應該正確取得特定模板', () => {
      const logic = new TemplateListLogic()
      logic.setTemplates([
        { id: 1, name: '模板 A' },
        { id: 2, name: '模板 B' },
      ])
      const template = logic.getTemplateById(2)
      expect(template?.name).toBe('模板 B')
    })

    it('getTemplateById 應該在找不到時返回 undefined', () => {
      const logic = new TemplateListLogic()
      logic.setTemplates([{ id: 1 }])
      expect(logic.getTemplateById(999)).toBeUndefined()
    })

    it('getActiveTemplates 應該返回活躍的模板', () => {
      const logic = new TemplateListLogic()
      logic.setTemplates([
        { id: 1, is_active: true },
        { id: 2, is_active: false },
        { id: 3, is_active: true },
      ])
      const active = logic.getActiveTemplates()
      expect(active).toHaveLength(2)
    })

    it('getRoomTemplates 應該返回教室視角的模板', () => {
      const logic = new TemplateListLogic()
      logic.setTemplates([
        { id: 1, row_type: 'ROOM' },
        { id: 2, row_type: 'TEACHER' },
      ])
      const room = logic.getRoomTemplates()
      expect(room).toHaveLength(1)
      expect(room[0].row_type).toBe('ROOM')
    })

    it('getTeacherTemplates 應該返回老師視角的模板', () => {
      const logic = new TemplateListLogic()
      logic.setTemplates([
        { id: 1, row_type: 'ROOM' },
        { id: 2, row_type: 'TEACHER' },
      ])
      const teacher = logic.getTeacherTemplates()
      expect(teacher).toHaveLength(1)
      expect(teacher[0].row_type).toBe('TEACHER')
    })

    it('addTemplate 應該新增模板', () => {
      const logic = new TemplateListLogic()
      logic.addTemplate({ id: 1, name: '新模板' })
      expect(logic.templates).toHaveLength(1)
    })

    it('removeTemplate 應該刪除模板', () => {
      const logic = new TemplateListLogic()
      logic.setTemplates([
        { id: 1, name: '模板 A' },
        { id: 2, name: '模板 B' },
      ])
      logic.removeTemplate(1)
      expect(logic.templates).toHaveLength(1)
      expect(logic.getTemplateById(1)).toBeUndefined()
    })
  })

  describe('TemplateTypeLogic 模板類型邏輯', () => {
    it('getTypeLabel 應該返回正確的標籤', () => {
      const logic = new TemplateTypeLogic()
      expect(logic.getTypeLabel('ROOM')).toBe('教室視角')
      expect(logic.getTypeLabel('TEACHER')).toBe('老師視角')
      expect(logic.getTypeLabel('UNKNOWN')).toBe('UNKNOWN')
    })

    it('getTypeClass 應該返回正確的 CSS 類別', () => {
      const logic = new TemplateTypeLogic()
      expect(logic.getTypeClass('ROOM')).toContain('primary')
      expect(logic.getTypeClass('TEACHER')).toContain('secondary')
      expect(logic.getTypeClass('UNKNOWN')).toContain('slate')
    })

    it('getTypeOptions 應該返回類型選項列表', () => {
      const logic = new TemplateTypeLogic()
      const options = logic.getTypeOptions()
      expect(options).toHaveLength(2)
      expect(options[0].value).toBe('ROOM')
      expect(options[1].value).toBe('TEACHER')
    })

    it('isValidType 應該正確驗證類型', () => {
      const logic = new TemplateTypeLogic()
      expect(logic.isValidType('ROOM')).toBe(true)
      expect(logic.isValidType('TEACHER')).toBe(true)
      expect(logic.isValidType('INVALID')).toBe(false)
    })
  })

  describe('TemplateFormLogic 模板表單邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new TemplateFormLogic()
      expect(logic.form.name).toBe('')
      expect(logic.form.row_type).toBe('ROOM')
    })

    it('setName 應該正確設定名稱', () => {
      const logic = new TemplateFormLogic()
      logic.setName('新模板')
      expect(logic.form.name).toBe('新模板')
    })

    it('setRowType 應該正確設定類型', () => {
      const logic = new TemplateFormLogic()
      logic.setRowType('TEACHER')
      expect(logic.form.row_type).toBe('TEACHER')
    })

    it('resetForm 應該重置表單', () => {
      const logic = new TemplateFormLogic()
      logic.form.name = '測試'
      logic.form.row_type = 'TEACHER'
      logic.resetForm()
      expect(logic.form.name).toBe('')
      expect(logic.form.row_type).toBe('ROOM')
    })

    it('isValid 應該在有名稱時返回 true', () => {
      const logic = new TemplateFormLogic()
      expect(logic.isValid()).toBe(false)
      logic.setName('新模板')
      expect(logic.isValid()).toBe(true)
    })

    it('isValid 應該在名稱為空白時返回 false', () => {
      const logic = new TemplateFormLogic()
      logic.setName('   ')
      expect(logic.isValid()).toBe(false)
    })

    it('getFormData 應該返回表單資料', () => {
      const logic = new TemplateFormLogic()
      logic.setName('新模板')
      logic.setRowType('TEACHER')
      const data = logic.getFormData()
      expect(data.name).toBe('新模板')
      expect(data.row_type).toBe('TEACHER')
    })
  })

  describe('WeekdaySelectionLogic 星期選擇邏輯', () => {
    it('應該正確初始化所有星期選項', () => {
      const logic = new WeekdaySelectionLogic()
      expect(logic.getWeekdays()).toHaveLength(7)
    })

    it('getWeekdayLabel 應該返回正確的標籤', () => {
      const logic = new WeekdaySelectionLogic()
      expect(logic.getWeekdayLabel(1)).toBe('週一')
      expect(logic.getWeekdayLabel(7)).toBe('週日')
    })

    it('toggleWeekday 應該正確切換星期選擇', () => {
      const logic = new WeekdaySelectionLogic()
      let selected: number[] = []
      selected = logic.toggleWeekday(selected, 1)
      expect(selected).toContain(1)
      selected = logic.toggleWeekday(selected, 1)
      expect(selected).not.toContain(1)
    })

    it('isWeekdaySelected 應該正確判斷是否選中', () => {
      const logic = new WeekdaySelectionLogic()
      expect(logic.isWeekdaySelected([1, 2], 1)).toBe(true)
      expect(logic.isWeekdaySelected([1, 2], 3)).toBe(false)
    })

    it('getWeekdayRange 應該返回正確的範圍文字', () => {
      const logic = new WeekdaySelectionLogic()
      expect(logic.getWeekdayRange([1, 2, 3])).toBe('週一、週二、週三')
      expect(logic.getWeekdayRange([1])).toBe('週一')
      expect(logic.getWeekdayRange([])).toBe('')
    })

    it('sortWeekdays 應該正確排序星期', () => {
      const logic = new WeekdaySelectionLogic()
      const sorted = logic.sortWeekdays([3, 1, 2])
      expect(sorted).toEqual([1, 2, 3])
    })
  })

  describe('ApplyFormLogic 套用表單邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new ApplyFormLogic()
      expect(logic.form.templateId).toBe(0)
      expect(logic.form.weekdays).toHaveLength(0)
      expect(logic.form.duration).toBe(60)
    })

    it('setTemplate 應該正確設定模板資訊', () => {
      const logic = new ApplyFormLogic()
      logic.setTemplate({ id: 1, name: '模板 A' })
      expect(logic.form.templateId).toBe(1)
      expect(logic.form.templateName).toBe('模板 A')
    })

    it('setWeekdays 應該正確設定星期', () => {
      const logic = new ApplyFormLogic()
      logic.setWeekdays([1, 3, 5])
      expect(logic.form.weekdays).toEqual([1, 3, 5])
    })

    it('isValid 應該在所有必要欄位都有值時返回 true', () => {
      const logic = new ApplyFormLogic()
      expect(logic.isValid()).toBe(false)
      logic.setOfferingId('1')
      expect(logic.isValid()).toBe(false)
      logic.setDateRange('2026-01-01', '2026-01-31')
      expect(logic.isValid()).toBe(false)
      logic.setWeekdays([1])
      expect(logic.isValid()).toBe(true)
    })

    it('getApplyData 應該返回正確的套用資料', () => {
      const logic = new ApplyFormLogic()
      logic.setOfferingId('1')
      logic.setDateRange('2026-01-01', '2026-01-31')
      logic.setWeekdays([1, 2, 3])
      logic.setDuration(90)
      const data = logic.getApplyData()
      expect(data.offering_id).toBe(1)
      expect(data.start_date).toBe('2026-01-01')
      expect(data.end_date).toBe('2026-01-31')
      expect(data.weekdays).toEqual([1, 2, 3])
      expect(data.duration).toBe(90)
    })
  })

  describe('TemplateCellLogic 模板格子邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new TemplateCellLogic()
      expect(logic.cells).toHaveLength(0)
    })

    it('setCells 應該正確設定格子列表', () => {
      const logic = new TemplateCellLogic()
      logic.setCells([
        { id: 1, row_no: 1, col_no: 1 },
        { id: 2, row_no: 1, col_no: 2 },
      ])
      expect(logic.cells).toHaveLength(2)
    })

    it('hasCells 應該正確判斷是否有格子', () => {
      const logic = new TemplateCellLogic()
      expect(logic.hasCells()).toBe(false)
      logic.setCells([{ id: 1 }])
      expect(logic.hasCells()).toBe(true)
    })

    it('sortCells 應該正確排序格子', () => {
      const logic = new TemplateCellLogic()
      logic.setCells([
        { id: 1, row_no: 2, col_no: 1 },
        { id: 2, row_no: 1, col_no: 2 },
        { id: 3, row_no: 1, col_no: 1 },
      ])
      const sorted = logic.sortCells()
      expect(sorted[0].id).toBe(3)
      expect(sorted[1].id).toBe(2)
      expect(sorted[2].id).toBe(1)
    })

    it('addCell 應該新增格子', () => {
      const logic = new TemplateCellLogic()
      logic.addCell({ id: 1 })
      expect(logic.getCellCount()).toBe(1)
    })

    it('removeCell 應該刪除格子', () => {
      const logic = new TemplateCellLogic()
      logic.setCells([{ id: 1 }, { id: 2 }])
      logic.removeCell(1)
      expect(logic.getCellCount()).toBe(1)
      expect(logic.getCellById(1)).toBeUndefined()
    })
  })

  describe('頁面整合邏輯', () => {
    it('應該能夠完整處理模板管理流程', () => {
      const listLogic = new TemplateListLogic()
      const typeLogic = new TemplateTypeLogic()
      const formLogic = new TemplateFormLogic()
      const weekdayLogic = new WeekdaySelectionLogic()
      const applyLogic = new ApplyFormLogic()

      // 設定模板
      listLogic.setTemplates([
        { id: 1, name: '週一模板', row_type: 'ROOM', is_active: true },
        { id: 2, name: '週二模板', row_type: 'TEACHER', is_active: true },
      ])

      // 驗證模板列表
      expect(listLogic.hasTemplates()).toBe(true)
      expect(listLogic.getRoomTemplates()).toHaveLength(1)
      expect(listLogic.getTeacherTemplates()).toHaveLength(1)

      // 驗證類型顯示
      expect(typeLogic.getTypeLabel('ROOM')).toBe('教室視角')
      expect(typeLogic.getTypeClass('ROOM')).toContain('primary')

      // 建立新模板
      formLogic.setName('新模板')
      formLogic.setRowType('TEACHER')
      expect(formLogic.isValid()).toBe(true)
      const formData = formLogic.getFormData()
      expect(formData.name).toBe('新模板')
      expect(formData.row_type).toBe('TEACHER')

      // 選擇星期
      let selectedWeekdays: number[] = []
      selectedWeekdays = weekdayLogic.toggleWeekday(selectedWeekdays, 1)
      selectedWeekdays = weekdayLogic.toggleWeekday(selectedWeekdays, 3)
      selectedWeekdays = weekdayLogic.toggleWeekday(selectedWeekdays, 5)
      expect(weekdayLogic.getSelectedCount(selectedWeekdays)).toBe(3)
      expect(weekdayLogic.getWeekdayRange(selectedWeekdays)).toBe('週一、週三、週五')

      // 套用模板
      applyLogic.setTemplate(listLogic.getTemplates()[0])
      applyLogic.setOfferingId('1')
      applyLogic.setDateRange('2026-01-01', '2026-01-31')
      applyLogic.setWeekdays(selectedWeekdays)
      applyLogic.setDuration(90)
      expect(applyLogic.isValid()).toBe(true)

      const applyData = applyLogic.getApplyData()
      expect(applyData.offering_id).toBe(1)
      expect(applyData.weekdays).toEqual([1, 3, 5])
      expect(applyData.duration).toBe(90)
    })

    it('應該正確處理不同狀態的模板', () => {
      const listLogic = new TemplateListLogic()

      listLogic.setTemplates([
        { id: 1, name: '活躍模板', is_active: true },
        { id: 2, name: '非活躍模板', is_active: false },
        { id: 3, name: '預設活躍', is_active: undefined },
      ])

      expect(listLogic.getActiveTemplates()).toHaveLength(2)
      expect(listLogic.getInactiveTemplates()).toHaveLength(1)
    })
  })
})
