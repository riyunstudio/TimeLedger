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

vi.mock('~/stores/auth', () => ({
  useAuthStore: () => ({
    user: { name: '測試老師', bio: '測試簡介' },
  }),
}))

vi.mock('~/stores/useScheduleStore', () => ({
  useScheduleStore: () => ({
    centers: [],
    fetchCenters: vi.fn(),
  }),
}))

vi.mock('~/composables/useSidebar', () => ({
  useSidebar: () => ({
    isOpen: { value: false },
    close: vi.fn(),
  }),
}))

describe('teacher/profile.vue 頁面邏輯', () => {
  // MembershipStatusLogic 類別 - 會員狀態邏輯
  class MembershipStatusLogic {
    getStatusClass(status: string): string {
      switch (status) {
        case 'ACTIVE':
          return 'bg-success-500/20 text-success-500'
        case 'INVITED':
          return 'bg-warning-500/20 text-warning-500'
        case 'INACTIVE':
          return 'bg-slate-500/20 text-slate-400'
        default:
          return 'bg-slate-500/20 text-slate-400'
      }
    }

    getStatusText(status: string): string {
      switch (status) {
        case 'ACTIVE':
          return '已加入'
        case 'INVITED':
          return '邀請中'
        case 'INACTIVE':
          return '已離開'
        default:
          return status
      }
    }

    isActive(status: string): boolean {
      return status === 'ACTIVE'
    }

    isInvited(status: string): boolean {
      return status === 'INVITED'
    }

    isInactive(status: string): boolean {
      return status === 'INACTIVE'
    }

    canManageMembership(status: string): boolean {
      return status === 'ACTIVE' || status === 'INVITED'
    }
  }

  // CenterMembershipLogic 類別 - 中心會員邏輯
  class CenterMembershipLogic {
    memberships: any[]

    constructor() {
      this.memberships = []
    }

    setMemberships(memberships: any[]) {
      this.memberships = memberships
    }

    getActiveMemberships(): any[] {
      return this.memberships.filter(m => m.status === 'ACTIVE')
    }

    getInvitedMemberships(): any[] {
      return this.memberships.filter(m => m.status === 'INVITED')
    }

    getInactiveMemberships(): any[] {
      return this.memberships.filter(m => m.status === 'INACTIVE')
    }

    getMembershipCount(): number {
      return this.memberships.length
    }

    getActiveCount(): number {
      return this.getActiveMemberships().length
    }

    getMembershipByCenterId(centerId: number): any | undefined {
      return this.memberships.find(m => m.center_id === centerId)
    }

    hasMemberships(): boolean {
      return this.memberships.length > 0
    }

    hasActiveMemberships(): boolean {
      return this.getActiveMemberships().length > 0
    }

    isMemberOfCenter(centerId: number): boolean {
      const membership = this.getMembershipByCenterId(centerId)
      return membership?.status === 'ACTIVE'
    }
  }

  // ModalStateLogic 類別 - Modal 狀態邏輯
  class ModalStateLogic {
    showProfileModal: boolean
    showHiringModal: boolean
    showSkillsModal: boolean
    showExportModal: boolean

    constructor() {
      this.showProfileModal = false
      this.showHiringModal = false
      this.showSkillsModal = false
      this.showExportModal = false
    }

    openProfileModal() {
      this.closeAll()
      this.showProfileModal = true
    }

    openHiringModal() {
      this.closeAll()
      this.showHiringModal = true
    }

    openSkillsModal() {
      this.closeAll()
      this.showSkillsModal = true
    }

    openExportModal() {
      this.closeAll()
      this.showExportModal = true
    }

    closeProfileModal() {
      this.showProfileModal = false
    }

    closeHiringModal() {
      this.showHiringModal = false
    }

    closeSkillsModal() {
      this.showSkillsModal = false
    }

    closeExportModal() {
      this.showExportModal = false
    }

    closeAll() {
      this.showProfileModal = false
      this.showHiringModal = false
      this.showSkillsModal = false
      this.showExportModal = false
    }

    isAnyModalOpen(): boolean {
      return (
        this.showProfileModal ||
        this.showHiringModal ||
        this.showSkillsModal ||
        this.showExportModal
      )
    }

    getOpenModalName(): string | null {
      if (this.showProfileModal) return 'ProfileModal'
      if (this.showHiringModal) return 'HiringModal'
      if (this.showSkillsModal) return 'SkillsModal'
      if (this.showExportModal) return 'ExportModal'
      return null
    }
  }

  // ProfileDisplayLogic 類別 - 個人檔案顯示邏輯
  class ProfileDisplayLogic {
    getAvatarInitial(name: string): string {
      return name?.charAt(0) || 'T'
    }

    getAvatarGradient(index: number = 0): string {
      const gradients = [
        'from-primary-500 to-secondary-500',
        'from-blue-500 to-purple-500',
        'from-green-500 to-teal-500',
        'from-orange-500 to-red-500',
      ]
      return gradients[index % gradients.length]
    }

    formatJoinDate(dateStr: string): string {
      if (!dateStr) return '-'
      return new Date(dateStr).toLocaleDateString('zh-TW')
    }

    truncateBio(bio: string, maxLength: number = 50): string {
      if (!bio) return ''
      if (bio.length <= maxLength) return bio
      return bio.substring(0, maxLength) + '...'
    }
  }

  // NavigationLogic 類別 - 導航邏輯
  class NavigationLogic {
    menuItems: { id: string; label: string; icon: string; route?: string }[]

    constructor() {
      this.menuItems = [
        { id: 'profile', label: '個人檔案', icon: 'user' },
        { id: 'hiring', label: '求職設定', icon: 'briefcase' },
        { id: 'skills', label: '技能與證照', icon: 'skills' },
        { id: 'export', label: '匯出課表', icon: 'download' },
      ]
    }

    getMenuItems() {
      return this.menuItems
    }

    getMenuItemById(id: string) {
      return this.menuItems.find(item => item.id === id)
    }

    getMenuItemCount(): number {
      return this.menuItems.length
    }

    isValidMenuId(id: string): boolean {
      return this.menuItems.some(item => item.id === id)
    }
  }

  // StatsLogic 類別 - 統計資訊邏輯
  class StatsLogic {
    memberships: any[]
    skills: any[]
    certificates: any[]

    constructor() {
      this.memberships = []
      this.skills = []
      this.certificates = []
    }

    setMemberships(memberships: any[]) {
      this.memberships = memberships
    }

    setSkills(skills: any[]) {
      this.skills = skills
    }

    setCertificates(certificates: any[]) {
      this.certificates = certificates
    }

    getTotalCenterCount(): number {
      return this.memberships.length
    }

    getActiveCenterCount(): number {
      return this.memberships.filter(m => m.status === 'ACTIVE').length
    }

    getSkillCount(): number {
      return this.skills.length
    }

    getCertificateCount(): number {
      return this.certificates.length
    }

    getSkillCategories(): string[] {
      return [...new Set(this.skills.map(s => s.category))]
    }

    getSkillsByCategory(category: string): any[] {
      return this.skills.filter(s => s.category === category)
    }
  }

  describe('MembershipStatusLogic 會員狀態邏輯', () => {
    it('getStatusClass 應該返回正確的 CSS 類別', () => {
      const logic = new MembershipStatusLogic()
      expect(logic.getStatusClass('ACTIVE')).toContain('success')
      expect(logic.getStatusClass('INVITED')).toContain('warning')
      expect(logic.getStatusClass('INACTIVE')).toContain('slate')
      expect(logic.getStatusClass('UNKNOWN')).toContain('slate')
    })

    it('getStatusText 應該返回正確的中文文字', () => {
      const logic = new MembershipStatusLogic()
      expect(logic.getStatusText('ACTIVE')).toBe('已加入')
      expect(logic.getStatusText('INVITED')).toBe('邀請中')
      expect(logic.getStatusText('INACTIVE')).toBe('已離開')
      expect(logic.getStatusText('UNKNOWN')).toBe('UNKNOWN')
    })

    it('isActive 應該正確判斷是否為活躍狀態', () => {
      const logic = new MembershipStatusLogic()
      expect(logic.isActive('ACTIVE')).toBe(true)
      expect(logic.isActive('INACTIVE')).toBe(false)
    })

    it('isInvited 應該正確判斷是否為邀請狀態', () => {
      const logic = new MembershipStatusLogic()
      expect(logic.isInvited('INVITED')).toBe(true)
      expect(logic.isInvited('ACTIVE')).toBe(false)
    })

    it('canManageMembership 應該在狀態為 ACTIVE 或 INVITED 時返回 true', () => {
      const logic = new MembershipStatusLogic()
      expect(logic.canManageMembership('ACTIVE')).toBe(true)
      expect(logic.canManageMembership('INVITED')).toBe(true)
      expect(logic.canManageMembership('INACTIVE')).toBe(false)
    })
  })

  describe('CenterMembershipLogic 中心會員邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new CenterMembershipLogic()
      expect(logic.memberships).toHaveLength(0)
    })

    it('setMemberships 應該正確設定會員列表', () => {
      const logic = new CenterMembershipLogic()
      logic.setMemberships([
        { center_id: 1, center_name: '中心 A', status: 'ACTIVE' },
        { center_id: 2, center_name: '中心 B', status: 'INVITED' },
      ])
      expect(logic.memberships).toHaveLength(2)
    })

    it('getActiveMemberships 應該返回活躍的會員', () => {
      const logic = new CenterMembershipLogic()
      logic.setMemberships([
        { center_id: 1, status: 'ACTIVE' },
        { center_id: 2, status: 'INVITED' },
        { center_id: 3, status: 'INACTIVE' },
      ])
      const active = logic.getActiveMemberships()
      expect(active).toHaveLength(1)
      expect(active[0].center_id).toBe(1)
    })

    it('getInvitedMemberships 應該返回邀請中的會員', () => {
      const logic = new CenterMembershipLogic()
      logic.setMemberships([
        { center_id: 1, status: 'ACTIVE' },
        { center_id: 2, status: 'INVITED' },
      ])
      const invited = logic.getInvitedMemberships()
      expect(invited).toHaveLength(1)
      expect(invited[0].center_id).toBe(2)
    })

    it('getMembershipByCenterId 應該正確取得特定中心的會員', () => {
      const logic = new CenterMembershipLogic()
      logic.setMemberships([
        { center_id: 1, center_name: '中心 A' },
        { center_id: 2, center_name: '中心 B' },
      ])
      const membership = logic.getMembershipByCenterId(2)
      expect(membership?.center_name).toBe('中心 B')
    })

    it('getMembershipByCenterId 應該在找不到時返回 undefined', () => {
      const logic = new CenterMembershipLogic()
      logic.setMemberships([{ center_id: 1 }])
      const membership = logic.getMembershipByCenterId(999)
      expect(membership).toBeUndefined()
    })

    it('isMemberOfCenter 應該正確判斷是否為中心會員', () => {
      const logic = new CenterMembershipLogic()
      logic.setMemberships([
        { center_id: 1, status: 'ACTIVE' },
        { center_id: 2, status: 'INVITED' },
      ])
      expect(logic.isMemberOfCenter(1)).toBe(true)
      expect(logic.isMemberOfCenter(2)).toBe(false)
      expect(logic.isMemberOfCenter(3)).toBe(false)
    })

    it('hasActiveMemberships 應該正確判斷是否有活躍會員', () => {
      const logic = new CenterMembershipLogic()
      expect(logic.hasActiveMemberships()).toBe(false)
      logic.setMemberships([{ center_id: 1, status: 'INACTIVE' }])
      expect(logic.hasActiveMemberships()).toBe(false)
      logic.setMemberships([{ center_id: 1, status: 'ACTIVE' }])
      expect(logic.hasActiveMemberships()).toBe(true)
    })
  })

  describe('ModalStateLogic Modal 狀態邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new ModalStateLogic()
      expect(logic.showProfileModal).toBe(false)
      expect(logic.showHiringModal).toBe(false)
      expect(logic.showSkillsModal).toBe(false)
      expect(logic.showExportModal).toBe(false)
    })

    it('openProfileModal 應該開啟個人檔案 Modal 並關閉其他', () => {
      const logic = new ModalStateLogic()
      logic.showHiringModal = true
      logic.openProfileModal()
      expect(logic.showProfileModal).toBe(true)
      expect(logic.showHiringModal).toBe(false)
    })

    it('openHiringModal 應該開啟求職設定 Modal', () => {
      const logic = new ModalStateLogic()
      logic.openHiringModal()
      expect(logic.showHiringModal).toBe(true)
    })

    it('closeAll 應該關閉所有 Modal', () => {
      const logic = new ModalStateLogic()
      logic.showProfileModal = true
      logic.showHiringModal = true
      logic.showSkillsModal = true
      logic.showExportModal = true
      logic.closeAll()
      expect(logic.isAnyModalOpen()).toBe(false)
    })

    it('isAnyModalOpen 應該正確判斷是否有 Modal 開啟', () => {
      const logic = new ModalStateLogic()
      expect(logic.isAnyModalOpen()).toBe(false)
      logic.showProfileModal = true
      expect(logic.isAnyModalOpen()).toBe(true)
    })

    it('getOpenModalName 應該返回正確的 Modal 名稱', () => {
      const logic = new ModalStateLogic()
      expect(logic.getOpenModalName()).toBeNull()
      logic.showProfileModal = true
      expect(logic.getOpenModalName()).toBe('ProfileModal')
    })
  })

  describe('ProfileDisplayLogic 個人檔案顯示邏輯', () => {
    it('getAvatarInitial 應該返回名字的第一個字元', () => {
      const logic = new ProfileDisplayLogic()
      expect(logic.getAvatarInitial('張老師')).toBe('張')
      expect(logic.getAvatarInitial('李老師')).toBe('李')
      expect(logic.getAvatarInitial('')).toBe('T')
      expect(logic.getAvatarInitial(null as any)).toBe('T')
    })

    it('getAvatarGradient 應該返回正確的漸層類別', () => {
      const logic = new ProfileDisplayLogic()
      expect(logic.getAvatarGradient(0)).toContain('primary')
      expect(logic.getAvatarGradient(1)).toContain('blue')
      expect(logic.getAvatarGradient(4)).toContain('primary')
    })

    it('formatJoinDate 應該正確格式化加入日期', () => {
      const logic = new ProfileDisplayLogic()
      const result = logic.formatJoinDate('2026-01-20')
      expect(result).toContain('2026')
      expect(result).toContain('1')
      expect(result).toContain('20')
    })

    it('formatJoinDate 應該處理空值', () => {
      const logic = new ProfileDisplayLogic()
      expect(logic.formatJoinDate('')).toBe('-')
    })

    it('truncateBio 應該正確截斷文字', () => {
      const logic = new ProfileDisplayLogic()
      const shortBio = '簡短簡介'
      expect(logic.truncateBio(shortBio, 50)).toBe(shortBio)

      const longBio = '這是一個很長的簡介，包含了許多的內容，需要被截斷處理'
      const truncated = logic.truncateBio(longBio, 20)
      expect(truncated).toBe('這是一個很長的簡介，...')
      expect(truncated.length).toBe(23)
    })
  })

  describe('NavigationLogic 導航邏輯', () => {
    it('應該正確初始化所有選單項目', () => {
      const logic = new NavigationLogic()
      expect(logic.getMenuItemCount()).toBe(4)
    })

    it('getMenuItems 應該返回所有選單項目', () => {
      const logic = new NavigationLogic()
      const items = logic.getMenuItems()
      expect(items[0].id).toBe('profile')
      expect(items[1].id).toBe('hiring')
      expect(items[2].id).toBe('skills')
      expect(items[3].id).toBe('export')
    })

    it('getMenuItemById 應該返回正確的選單項目', () => {
      const logic = new NavigationLogic()
      const item = logic.getMenuItemById('hiring')
      expect(item?.label).toBe('求職設定')
    })

    it('isValidMenuId 應該正確驗證選單 ID', () => {
      const logic = new NavigationLogic()
      expect(logic.isValidMenuId('profile')).toBe(true)
      expect(logic.isValidMenuId('invalid')).toBe(false)
    })
  })

  describe('StatsLogic 統計資訊邏輯', () => {
    it('應該正確初始化', () => {
      const logic = new StatsLogic()
      expect(logic.memberships).toHaveLength(0)
      expect(logic.skills).toHaveLength(0)
      expect(logic.certificates).toHaveLength(0)
    })

    it('setMemberships 應該正確設定會員資料', () => {
      const logic = new StatsLogic()
      logic.setMemberships([{ center_id: 1 }, { center_id: 2 }])
      expect(logic.getTotalCenterCount()).toBe(2)
    })

    it('setSkills 應該正確設定技能資料', () => {
      const logic = new StatsLogic()
      logic.setSkills([
        { name: '鋼琴', category: 'INSTRUMENT' },
        { name: '小提琴', category: 'INSTRUMENT' },
        { name: '樂理', category: 'THEORY' },
      ])
      expect(logic.getSkillCount()).toBe(3)
    })

    it('getSkillCategories 應該返回不重複的類別列表', () => {
      const logic = new StatsLogic()
      logic.setSkills([
        { category: 'INSTRUMENT' },
        { category: 'INSTRUMENT' },
        { category: 'THEORY' },
      ])
      const categories = logic.getSkillCategories()
      expect(categories).toHaveLength(2)
      expect(categories).toContain('INSTRUMENT')
      expect(categories).toContain('THEORY')
    })

    it('getSkillsByCategory 應該返回指定類別的技能', () => {
      const logic = new StatsLogic()
      logic.setSkills([
        { name: '鋼琴', category: 'INSTRUMENT' },
        { name: '小提琴', category: 'INSTRUMENT' },
        { name: '樂理', category: 'THEORY' },
      ])
      const skills = logic.getSkillsByCategory('INSTRUMENT')
      expect(skills).toHaveLength(2)
      expect(skills[0].name).toBe('鋼琴')
      expect(skills[1].name).toBe('小提琴')
    })
  })

  describe('頁面整合邏輯', () => {
    it('應該能夠完整處理個人檔案頁面流程', () => {
      const membershipLogic = new CenterMembershipLogic()
      const modalLogic = new ModalStateLogic()
      const profileLogic = new ProfileDisplayLogic()
      const statusLogic = new MembershipStatusLogic()

      // 設定會員資料
      membershipLogic.setMemberships([
        {
          center_id: 1,
          center_name: '音樂教室 A',
          status: 'ACTIVE',
          joined_at: '2025-06-15',
        },
        {
          center_id: 2,
          center_name: '藝術中心 B',
          status: 'INVITED',
        },
      ])

      // 驗證會員狀態顯示
      expect(statusLogic.getStatusClass('ACTIVE')).toContain('success')
      expect(statusLogic.getStatusText('ACTIVE')).toBe('已加入')

      // 驗證 Modal 狀態
      expect(modalLogic.isAnyModalOpen()).toBe(false)
      modalLogic.openProfileModal()
      expect(modalLogic.getOpenModalName()).toBe('ProfileModal')

      // 驗證頭像顯示
      const initial = profileLogic.getAvatarInitial('張老師')
      expect(initial).toBe('張')

      // 驗證會員數量
      expect(membershipLogic.getActiveCount()).toBe(1)
      expect(membershipLogic.hasActiveMemberships()).toBe(true)

      // 關閉 Modal
      modalLogic.closeProfileModal()
      expect(modalLogic.isAnyModalOpen()).toBe(false)
    })

    it('應該正確處理多個中心的會員狀態', () => {
      const membershipLogic = new CenterMembershipLogic()
      const statusLogic = new MembershipStatusLogic()

      membershipLogic.setMemberships([
        { center_id: 1, status: 'ACTIVE' },
        { center_id: 2, status: 'ACTIVE' },
        { center_id: 3, status: 'INVITED' },
        { center_id: 4, status: 'INACTIVE' },
      ])

      // 測試活躍會員
      const activeMemberships = membershipLogic.getActiveMemberships()
      expect(activeMemberships).toHaveLength(2)

      // 測試邀請中會員
      const invitedMemberships = membershipLogic.getInvitedMemberships()
      expect(invitedMemberships).toHaveLength(1)

      // 測試非活躍會員
      const inactiveMemberships = membershipLogic.getInactiveMemberships()
      expect(inactiveMemberships).toHaveLength(1)

      // 測試狀態文字顯示
      expect(statusLogic.getStatusText('ACTIVE')).toBe('已加入')
      expect(statusLogic.getStatusText('INVITED')).toBe('邀請中')
      expect(statusLogic.getStatusText('INACTIVE')).toBe('已離開')
    })
  })
})
