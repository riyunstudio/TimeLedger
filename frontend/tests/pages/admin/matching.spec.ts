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
    watch: (fn: () => void) => {},
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

vi.mock('~/types', () => ({
  SKILL_CATEGORIES: {
    INSTRUMENT: { icon: '🎹', color: 'bg-purple-500/20 text-purple-400' },
    VOCAL: { icon: '🎤', color: 'bg-pink-500/20 text-pink-400' },
    THEORY: { icon: '📚', color: 'bg-blue-500/20 text-blue-400' },
    COMPOSITION: { icon: '🎼', color: 'bg-yellow-500/20 text-yellow-400' },
    OTHER: { icon: '✨', color: 'bg-slate-500/20 text-slate-400' },
  },
}))

describe('admin/matching.vue 頁面邏輯', () => {
  // MatchingFormLogic 類別 - 媒合搜尋表單邏輯
  class MatchingFormLogic {
    form: {
      start_time: string
      end_time: string
      room_ids: number[]
      skills: string
    }

    constructor() {
      this.form = {
        start_time: '',
        end_time: '',
        room_ids: [],
        skills: '',
      }
    }

    setTimeRange(start: string, end: string) {
      this.form.start_time = start
      this.form.end_time = end
    }

    toggleRoom(roomId: number) {
      const index = this.form.room_ids.indexOf(roomId)
      if (index === -1) {
        this.form.room_ids.push(roomId)
      } else {
        this.form.room_ids.splice(index, 1)
      }
    }

    isRoomSelected(roomId: number): boolean {
      return this.form.room_ids.includes(roomId)
    }

    setSkills(skills: string) {
      this.form.skills = skills
    }

    getRequiredSkills(): string[] {
      return this.form.skills
        .split(',')
        .map(s => s.trim())
        .filter(Boolean)
    }

    isFormValid(): boolean {
      return Boolean(this.form.start_time && this.form.end_time)
    }

    clearForm() {
      this.form = {
        start_time: '',
        end_time: '',
        room_ids: [],
        skills: '',
      }
    }

    resetForm() {
      this.clearForm()
    }
  }

  // MatchingResultLogic 類別 - 媒合結果邏輯
  class MatchingResultLogic {
    matches: any[]
    hasSearched: boolean
    searching: boolean

    constructor() {
      this.matches = []
      this.hasSearched = false
      this.searching = false
    }

    setMatches(matches: any[]) {
      this.matches = matches
    }

    hasResults(): boolean {
      return this.matches.length > 0
    }

    getMatchCount(): number {
      return this.matches.length
    }

    getTopMatches(limit: number = 3): any[] {
      return [...this.matches]
        .sort((a, b) => b.match_score - a.match_score)
        .slice(0, limit)
    }

    getMatchById(teacherId: number): any | undefined {
      return this.matches.find(m => m.teacher_id === teacherId)
    }

    filterBySkillMatch(minScore: number): any[] {
      return this.matches.filter(m => m.skill_match >= minScore)
    }

    filterByRating(minRating: number): any[] {
      return this.matches.filter(m => (m.rating || 0) >= minRating)
    }

    getAverageMatchScore(): number {
      if (this.matches.length === 0) return 0
      const total = this.matches.reduce((sum, m) => sum + m.match_score, 0)
      return Math.round(total / this.matches.length)
    }

    setSearching(searching: boolean) {
      this.searching = searching
    }

    setHasSearched(hasSearched: boolean) {
      this.hasSearched = hasSearched
    }

    clearResults() {
      this.matches = []
      this.hasSearched = false
    }
  }

  // TalentSearchLogic 類別 - 人才庫搜尋邏輯
  class TalentSearchLogic {
    searchParams: {
      city: string
      skills: string
      hashtags: string
    }
    results: any[]
    searching: boolean

    constructor() {
      this.searchParams = {
        city: '',
        skills: '',
        hashtags: '',
      }
      this.results = []
      this.searching = false
    }

    setCity(city: string) {
      this.searchParams.city = city
    }

    setSkills(skills: string) {
      this.searchParams.skills = skills
    }

    setHashtags(hashtags: string) {
      this.searchParams.hashtags = hashtags
    }

    getSearchParams(): Record<string, string> {
      const params: Record<string, string> = {}
      if (this.searchParams.city) params.city = this.searchParams.city
      if (this.searchParams.skills) params.skills = this.searchParams.skills
      if (this.searchParams.hashtags) params.hashtags = this.searchParams.hashtags
      return params
    }

    hasSearchParams(): boolean {
      return Boolean(
        this.searchParams.city ||
        this.searchParams.skills ||
        this.searchParams.hashtags
      )
    }

    setResults(results: any[]) {
      this.results = results
    }

    getResultCount(): number {
      return this.results.length
    }

    filterByCity(city: string): any[] {
      return this.results.filter(t => t.city?.includes(city))
    }

    filterBySkill(skill: string): any[] {
      return this.results.filter(t =>
        t.skills?.some((s: any) =>
          s.name?.toLowerCase().includes(skill.toLowerCase())
        )
      )
    }

    filterByHashtag(hashtag: string): any[] {
      return this.results.filter(t =>
        t.personal_hashtags?.includes(hashtag)
      )
    }

    setSearching(searching: boolean) {
      this.searching = searching
    }

    clearSearch() {
      this.searchParams = { city: '', skills: '', hashtags: '' }
      this.results = []
    }
  }

  // RoomSelectionLogic 類別 - 教室選擇邏輯
  class RoomSelectionLogic {
    rooms: any[]
    selectedRoomIds: number[]

    constructor() {
      this.rooms = []
      this.selectedRoomIds = []
    }

    setRooms(rooms: any[]) {
      this.rooms = rooms
    }

    getRooms(): any[] {
      return this.rooms
    }

    toggleRoom(roomId: number) {
      const index = this.selectedRoomIds.indexOf(roomId)
      if (index === -1) {
        this.selectedRoomIds.push(roomId)
      } else {
        this.selectedRoomIds.splice(index, 1)
      }
    }

    isRoomSelected(roomId: number): boolean {
      return this.selectedRoomIds.includes(roomId)
    }

    getSelectedRooms(): any[] {
      return this.rooms.filter(r => this.selectedRoomIds.includes(r.id))
    }

    getSelectedCount(): number {
      return this.selectedRoomIds.length
    }

    selectAllRooms() {
      this.selectedRoomIds = this.rooms.map(r => r.id)
    }

    clearSelection() {
      this.selectedRoomIds = []
    }

    hasSelection(): boolean {
      return this.selectedRoomIds.length > 0
    }
  }

  // SkillCategoryLogic 類別 - 技能類別邏輯
  class SkillCategoryLogic {
    categories: Record<string, { icon: string; color: string }>

    constructor() {
      this.categories = {
        INSTRUMENT: { icon: '🎹', color: 'bg-purple-500/20 text-purple-400' },
        VOCAL: { icon: '🎤', color: 'bg-pink-500/20 text-pink-400' },
        THEORY: { icon: '📚', color: 'bg-blue-500/20 text-blue-400' },
        COMPOSITION: { icon: '🎼', color: 'bg-yellow-500/20 text-yellow-400' },
        OTHER: { icon: '✨', color: 'bg-slate-500/20 text-slate-400' },
      }
    }

    getIcon(category: string): string {
      return this.categories[category as keyof typeof this.categories]?.icon || '✨'
    }

    getColor(category: string): string {
      return this.categories[category as keyof typeof this.categories]?.color || 'bg-slate-500/20 text-slate-400'
    }

    getCategoryList(): { id: string; icon: string; color: string }[] {
      return Object.entries(this.categories).map(([id, value]) => ({
        id,
        ...value,
      }))
    }
  }

  describe('MatchingFormLogic 媒合表單邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new MatchingFormLogic()
      expect(logic.form.start_time).toBe('')
      expect(logic.form.end_time).toBe('')
      expect(logic.form.room_ids).toHaveLength(0)
      expect(logic.form.skills).toBe('')
    })

    it('setTimeRange 應該正確設定時間範圍', () => {
      const logic = new MatchingFormLogic()
      logic.setTimeRange('2026-01-20T09:00', '2026-01-20T12:00')
      expect(logic.form.start_time).toBe('2026-01-20T09:00')
      expect(logic.form.end_time).toBe('2026-01-20T12:00')
    })

    it('toggleRoom 應該正確切換教室選擇', () => {
      const logic = new MatchingFormLogic()
      logic.toggleRoom(1)
      expect(logic.isRoomSelected(1)).toBe(true)
      logic.toggleRoom(1)
      expect(logic.isRoomSelected(1)).toBe(false)
      logic.toggleRoom(1)
      expect(logic.isRoomSelected(1)).toBe(true)
    })

    it('getRequiredSkills 應該正確解析技能字串', () => {
      const logic = new MatchingFormLogic()
      logic.form.skills = '鋼琴, 小提琴, 鋼琴'
      const skills = logic.getRequiredSkills()
      expect(skills).toEqual(['鋼琴', '小提琴', '鋼琴'])
    })

    it('getRequiredSkills 應該過濾空白值', () => {
      const logic = new MatchingFormLogic()
      logic.form.skills = '鋼琴, , 小提琴,  '
      const skills = logic.getRequiredSkills()
      expect(skills).toEqual(['鋼琴', '小提琴'])
    })

    it('getRequiredSkills 應該處理空字串', () => {
      const logic = new MatchingFormLogic()
      logic.form.skills = ''
      const skills = logic.getRequiredSkills()
      expect(skills).toEqual([])
    })

    it('isFormValid 應該在有開始和結束時間時返回 true', () => {
      const logic = new MatchingFormLogic()
      expect(logic.isFormValid()).toBe(false)
      logic.setTimeRange('2026-01-20T09:00', '')
      expect(logic.isFormValid()).toBe(false)
      logic.setTimeRange('', '2026-01-20T12:00')
      expect(logic.isFormValid()).toBe(false)
      logic.setTimeRange('2026-01-20T09:00', '2026-01-20T12:00')
      expect(logic.isFormValid()).toBe(true)
    })

    it('clearForm 應該重置表單', () => {
      const logic = new MatchingFormLogic()
      logic.form.start_time = '2026-01-20T09:00'
      logic.form.end_time = '2026-01-20T12:00'
      logic.form.room_ids = [1, 2]
      logic.form.skills = '鋼琴'
      logic.clearForm()
      expect(logic.form.start_time).toBe('')
      expect(logic.form.end_time).toBe('')
      expect(logic.form.room_ids).toHaveLength(0)
      expect(logic.form.skills).toBe('')
    })
  })

  describe('MatchingResultLogic 媒合結果邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new MatchingResultLogic()
      expect(logic.matches).toHaveLength(0)
      expect(logic.hasSearched).toBe(false)
      expect(logic.searching).toBe(false)
    })

    it('setMatches 應該正確設定媒合結果', () => {
      const logic = new MatchingResultLogic()
      const matches = [
        { teacher_id: 1, teacher_name: '張老師', match_score: 85 },
        { teacher_id: 2, teacher_name: '李老師', match_score: 72 },
      ]
      logic.setMatches(matches)
      expect(logic.matches).toHaveLength(2)
    })

    it('hasResults 應該正確判斷是否有結果', () => {
      const logic = new MatchingResultLogic()
      expect(logic.hasResults()).toBe(false)
      logic.setMatches([{ teacher_id: 1 }])
      expect(logic.hasResults()).toBe(true)
    })

    it('getTopMatches 應該返回分數最高的老師', () => {
      const logic = new MatchingResultLogic()
      logic.setMatches([
        { teacher_id: 1, match_score: 85 },
        { teacher_id: 2, match_score: 95 },
        { teacher_id: 3, match_score: 78 },
      ])
      const top = logic.getTopMatches(2)
      expect(top).toHaveLength(2)
      expect(top[0].match_score).toBe(95)
      expect(top[1].match_score).toBe(85)
    })

    it('getMatchById 應該正確取得特定老師的媒合結果', () => {
      const logic = new MatchingResultLogic()
      logic.setMatches([
        { teacher_id: 1, teacher_name: '張老師' },
        { teacher_id: 2, teacher_name: '李老師' },
      ])
      const match = logic.getMatchById(2)
      expect(match?.teacher_name).toBe('李老師')
    })

    it('getMatchById 應該在找不到時返回 undefined', () => {
      const logic = new MatchingResultLogic()
      logic.setMatches([{ teacher_id: 1 }])
      const match = logic.getMatchById(999)
      expect(match).toBeUndefined()
    })

    it('filterBySkillMatch 應該過濾技能匹配分數', () => {
      const logic = new MatchingResultLogic()
      logic.setMatches([
        { teacher_id: 1, skill_match: 90 },
        { teacher_id: 2, skill_match: 75 },
        { teacher_id: 3, skill_match: 85 },
      ])
      const filtered = logic.filterBySkillMatch(80)
      expect(filtered).toHaveLength(2)
      expect(filtered[0].skill_match).toBe(90)
      expect(filtered[1].skill_match).toBe(85)
    })

    it('filterByRating 應該過濾評分', () => {
      const logic = new MatchingResultLogic()
      logic.setMatches([
        { teacher_id: 1, rating: 4.5 },
        { teacher_id: 2, rating: 3.8 },
        { teacher_id: 3, rating: 4.2 },
      ])
      const filtered = logic.filterByRating(4)
      expect(filtered).toHaveLength(2)
    })

    it('getAverageMatchScore 應該計算平均匹配分數', () => {
      const logic = new MatchingResultLogic()
      expect(logic.getAverageMatchScore()).toBe(0)
      logic.setMatches([
        { match_score: 80 },
        { match_score: 90 },
        { match_score: 70 },
      ])
      expect(logic.getAverageMatchScore()).toBe(80)
    })
  })

  describe('TalentSearchLogic 人才庫搜尋邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new TalentSearchLogic()
      expect(logic.searchParams.city).toBe('')
      expect(logic.searchParams.skills).toBe('')
      expect(logic.searchParams.hashtags).toBe('')
      expect(logic.results).toHaveLength(0)
    })

    it('setCity 應該正確設定城市', () => {
      const logic = new TalentSearchLogic()
      logic.setCity('台北市')
      expect(logic.searchParams.city).toBe('台北市')
    })

    it('setSkills 應該正確設定技能', () => {
      const logic = new TalentSearchLogic()
      logic.setSkills('鋼琴')
      expect(logic.searchParams.skills).toBe('鋼琴')
    })

    it('setHashtags 應該正確設定標籤', () => {
      const logic = new TalentSearchLogic()
      logic.setHashtags('古典 兒童')
      expect(logic.searchParams.hashtags).toBe('古典 兒童')
    })

    it('getSearchParams 應該返回非空的搜尋參數', () => {
      const logic = new TalentSearchLogic()
      logic.setCity('台北市')
      logic.setSkills('鋼琴')
      const params = logic.getSearchParams()
      expect(params.city).toBe('台北市')
      expect(params.skills).toBe('鋼琴')
      expect(params.hashtags).toBeUndefined()
    })

    it('hasSearchParams 應該正確判斷是否有搜尋參數', () => {
      const logic = new TalentSearchLogic()
      expect(logic.hasSearchParams()).toBe(false)
      logic.setCity('台北市')
      expect(logic.hasSearchParams()).toBe(true)
    })

    it('setResults 應該正確設定搜尋結果', () => {
      const logic = new TalentSearchLogic()
      const results = [
        { id: 1, name: '張老師' },
        { id: 2, name: '李老師' },
      ]
      logic.setResults(results)
      expect(logic.getResultCount()).toBe(2)
    })

    it('filterByCity 應該過濾城市', () => {
      const logic = new TalentSearchLogic()
      logic.setResults([
        { id: 1, city: '台北市' },
        { id: 2, city: '新北市' },
        { id: 3, city: '台北市' },
      ])
      const filtered = logic.filterByCity('台北市')
      expect(filtered).toHaveLength(2)
    })

    it('filterBySkill 應該過濾技能', () => {
      const logic = new TalentSearchLogic()
      logic.setResults([
        { id: 1, skills: [{ name: '鋼琴' }] },
        { id: 2, skills: [{ name: '小提琴' }] },
        { id: 3, skills: [{ name: '鋼琴' }] },
      ])
      const filtered = logic.filterBySkill('鋼琴')
      expect(filtered).toHaveLength(2)
    })

    it('filterByHashtag 應該過濾標籤', () => {
      const logic = new TalentSearchLogic()
      logic.setResults([
        { id: 1, personal_hashtags: ['古典', '兒童'] },
        { id: 2, personal_hashtags: ['流行'] },
        { id: 3, personal_hashtags: ['古典'] },
      ])
      const filtered = logic.filterByHashtag('古典')
      expect(filtered).toHaveLength(2)
    })

    it('clearSearch 應該清除搜尋', () => {
      const logic = new TalentSearchLogic()
      logic.setCity('台北市')
      logic.setSkills('鋼琴')
      logic.setResults([{ id: 1 }])
      logic.clearSearch()
      expect(logic.hasSearchParams()).toBe(false)
      expect(logic.getResultCount()).toBe(0)
    })
  })

  describe('RoomSelectionLogic 教室選擇邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new RoomSelectionLogic()
      expect(logic.rooms).toHaveLength(0)
      expect(logic.selectedRoomIds).toHaveLength(0)
    })

    it('setRooms 應該正確設定教室列表', () => {
      const logic = new RoomSelectionLogic()
      logic.setRooms([
        { id: 1, name: '教室 A' },
        { id: 2, name: '教室 B' },
      ])
      expect(logic.getRooms()).toHaveLength(2)
    })

    it('toggleRoom 應該正確切換選擇狀態', () => {
      const logic = new RoomSelectionLogic()
      logic.setRooms([{ id: 1 }, { id: 2 }])
      logic.toggleRoom(1)
      expect(logic.isRoomSelected(1)).toBe(true)
      expect(logic.isRoomSelected(2)).toBe(false)
      logic.toggleRoom(2)
      expect(logic.isRoomSelected(1)).toBe(true)
      expect(logic.isRoomSelected(2)).toBe(true)
    })

    it('getSelectedRooms 應該返回選中的教室', () => {
      const logic = new RoomSelectionLogic()
      logic.setRooms([
        { id: 1, name: '教室 A' },
        { id: 2, name: '教室 B' },
        { id: 3, name: '教室 C' },
      ])
      logic.toggleRoom(1)
      logic.toggleRoom(3)
      const selected = logic.getSelectedRooms()
      expect(selected).toHaveLength(2)
      expect(selected.map(r => r.id)).toEqual([1, 3])
    })

    it('getSelectedCount 應該返回選中數量', () => {
      const logic = new RoomSelectionLogic()
      logic.setRooms([{ id: 1 }, { id: 2 }])
      expect(logic.getSelectedCount()).toBe(0)
      logic.toggleRoom(1)
      expect(logic.getSelectedCount()).toBe(1)
    })

    it('selectAllRooms 應該選中所有教室', () => {
      const logic = new RoomSelectionLogic()
      logic.setRooms([{ id: 1 }, { id: 2 }, { id: 3 }])
      logic.selectAllRooms()
      expect(logic.hasSelection()).toBe(true)
      expect(logic.getSelectedCount()).toBe(3)
    })

    it('clearSelection 應該清除選擇', () => {
      const logic = new RoomSelectionLogic()
      logic.setRooms([{ id: 1 }])
      logic.toggleRoom(1)
      logic.clearSelection()
      expect(logic.hasSelection()).toBe(false)
      expect(logic.getSelectedCount()).toBe(0)
    })
  })

  describe('SkillCategoryLogic 技能類別邏輯', () => {
    it('應該正確初始化並包含所有類別', () => {
      const logic = new SkillCategoryLogic()
      expect(Object.keys(logic.categories)).toHaveLength(5)
    })

    it('getIcon 應該返回正確的圖示', () => {
      const logic = new SkillCategoryLogic()
      expect(logic.getIcon('INSTRUMENT')).toBe('🎹')
      expect(logic.getIcon('VOCAL')).toBe('🎤')
      expect(logic.getIcon('THEORY')).toBe('📚')
      expect(logic.getIcon('COMPOSITION')).toBe('🎼')
      expect(logic.getIcon('OTHER')).toBe('✨')
      expect(logic.getIcon('UNKNOWN')).toBe('✨')
    })

    it('getColor 應該返回正確的顏色', () => {
      const logic = new SkillCategoryLogic()
      expect(logic.getColor('INSTRUMENT')).toContain('purple')
      expect(logic.getColor('VOCAL')).toContain('pink')
      expect(logic.getColor('THEORY')).toContain('blue')
      expect(logic.getColor('COMPOSITION')).toContain('yellow')
      expect(logic.getColor('UNKNOWN')).toContain('slate')
    })

    it('getCategoryList 應該返回類別列表', () => {
      const logic = new SkillCategoryLogic()
      const list = logic.getCategoryList()
      expect(list).toHaveLength(5)
      expect(list[0]).toHaveProperty('id')
      expect(list[0]).toHaveProperty('icon')
      expect(list[0]).toHaveProperty('color')
    })
  })

  describe('頁面整合邏輯', () => {
    it('應該能夠完整執行媒合搜尋流程', () => {
      const formLogic = new MatchingFormLogic()
      const resultLogic = new MatchingResultLogic()
      const roomLogic = new RoomSelectionLogic()

      // 設定教室
      roomLogic.setRooms([
        { id: 1, name: '教室 A' },
        { id: 2, name: '教室 B' },
      ])
      roomLogic.toggleRoom(1)

      // 設定表單
      formLogic.setTimeRange('2026-01-20T09:00', '2026-01-20T12:00')
      formLogic.setSkills('鋼琴')

      // 驗證表單
      expect(formLogic.isFormValid()).toBe(true)
      expect(roomLogic.hasSelection()).toBe(true)

      // 模擬搜尋結果
      resultLogic.setMatches([
        {
          teacher_id: 1,
          teacher_name: '張老師',
          match_score: 85,
          skill_match: 90,
          rating: 4.5,
        },
      ])

      // 驗證結果
      expect(resultLogic.hasResults()).toBe(true)
      expect(resultLogic.getMatchCount()).toBe(1)
      expect(resultLogic.getAverageMatchScore()).toBe(85)
    })

    it('應該能夠完整執行人才庫搜尋流程', () => {
      const talentLogic = new TalentSearchLogic()
      const skillLogic = new SkillCategoryLogic()

      // 設定搜尋參數
      talentLogic.setCity('台北市')
      talentLogic.setSkills('鋼琴')
      talentLogic.setHashtags('古典')

      // 驗證搜尋參數
      expect(talentLogic.hasSearchParams()).toBe(true)

      // 設定搜尋結果
      talentLogic.setResults([
        {
          id: 1,
          name: '張老師',
          city: '台北市',
          skills: [{ name: '鋼琴', category: 'INSTRUMENT' }],
          personal_hashtags: ['古典'],
          bio: '鋼琴教學經驗豐富',
        },
      ])

      // 驗證結果
      expect(talentLogic.getResultCount()).toBe(1)

      // 測試技能類別顯示
      const icon = skillLogic.getIcon('INSTRUMENT')
      expect(icon).toBe('🎹')
    })
  })
})
