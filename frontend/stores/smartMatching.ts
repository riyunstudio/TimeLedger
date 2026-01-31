/**
 * 智慧媒合 Store
 *
 * 為智慧媒合功能建立統一的狀態管理
 * 包含人才庫搜尋、智慧媒合、替代時段建議等功能
 */

import { defineStore } from 'pinia'
import type {
  SmartMatchingRequest,
  SmartMatchingResult,
  SmartMatchingResponse,
  TalentPoolStats,
  TalentPoolStatsResponse,
  TalentSearchRequest,
  TalentSearchResponse,
  TalentCard,
  InviteTalentRequest,
  InviteTalentResponse,
  AlternativeSlotsRequest,
  AlternativeSlotsResponse,
  AlternativeSlot,
  TeacherScheduleQueryRequest,
  TeacherScheduleResponse,
  TeacherScheduleItem,
} from '~/types/matching'
import { withLoading } from '~/utils/loadingHelper'

// ==================== 搜尋相關類型 ====================

/**
 * 智慧媒合搜尋參數
 */
export interface MatchingParams {
  /** 房間 ID 清單 */
  room_ids?: number[]
  /** 開始時間 (YYYY-MM-DDTHH:mm:ss+08:00) */
  start_time: string
  /** 結束時間 (YYYY-MM-DDTHH:mm:ss+08:00) */
  end_time: string
  /** 需要的技能 */
  required_skills?: string[]
}

/**
 * 搜尋結果狀態
 */
export interface SearchState {
  /** 媒合結果清單 */
  results: SmartMatchingResult[]
  /** 是否已搜尋 */
  hasSearched: boolean
  /** 搜尋進度 */
  progress: number
  /** 搜尋步驟 */
  steps: SearchStep[]
}

/**
 * 搜尋步驟
 */
export interface SearchStep {
  /** 步驟名稱 */
  name: string
  /** 是否完成 */
  completed: boolean
  /** 是否進行中 */
  active: boolean
}

// ==================== 人才庫相關類型 ====================

/**
 * 人才搜尋參數
 */
export interface TalentSearchParams {
  /** 搜尋關鍵字 */
  keyword?: string
  /** 技能分類 */
  category?: string
  /** 技能名稱 */
  skill_name?: string
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 是否僅顯示開放應徵 */
  open_to_hiring_only?: boolean
  /** 最低評分 */
  min_rating?: number
  /** 排序欄位 */
  sort_by?: 'rating' | 'name' | 'recent'
  /** 排序方向 */
  sort_order?: 'ASC' | 'DESC'
  /** 頁碼 */
  page?: number
  /** 每頁數量 */
  limit?: number
}

/**
 * 人才庫搜尋狀態
 */
export interface TalentSearchState {
  /** 人才清單 */
  results: TalentCard[]
  /** 總數量 */
  totalItems: number
  /** 總頁數 */
  totalPages: number
  /** 目前頁碼 */
  currentPage: number
  /** 是否已搜尋 */
  hasSearched: boolean
  /** 統計資料 */
  stats: TalentPoolStats | null
  /** 縣市分布 */
  cityDistribution: Array<{ name: string; count: number }>
  /** 熱門技能 */
  topSkills: Array<{ name: string; count: number }>
}

// ==================== 教師課表相關類型 ====================

/**
 * 教師課表查詢參數
 */
export interface TeacherScheduleParams {
  /** 教師 ID */
  teacher_id: number
  /** 開始日期 */
  start_date: string
  /** 結束日期 */
  end_date: string
}

/**
 * 教師課表狀態
 */
export interface TeacherScheduleState {
  /** 選中的教師 */
  selectedTeacher: SmartMatchingResult | null
  /** 課表項目 */
  sessions: TeacherScheduleItem[]
  /** 替代時段 */
  alternativeSlots: AlternativeSlot[]
}

// ==================== 比較模式相關類型 ====================

/**
 * 比較模式狀態
 */
export interface CompareState {
  /** 選中比較清單 */
  selectedIds: Set<number>
  /** 檢視模式 */
  viewMode: 'card' | 'list' | 'compare'
}

/**
 * 邀請狀態
 */
export interface InvitationStatus {
  /** 是否已發送 */
  sent: boolean
  /** 變體 */
  variant: 'success' | 'warning' | 'error' | 'secondary'
  /** 文字 */
  text: string
}

// ==================== Smart Matching Store ====================

export const useSmartMatchingStore = defineStore('smartMatching', () => {
  // ==================== API 實例 ====================
  const api = useApi()

  // ==================== 搜尋相關狀態 ====================
  const searchResults = ref<SmartMatchingResult[]>([])
  const hasSearched = ref(false)
  const searchProgress = ref(0)
  const searchSteps = ref<SearchStep[]>([
    { name: '取得老師清單', completed: false, active: false },
    { name: '分析可用時間', completed: false, active: false },
    { name: '計算技能匹配度', completed: false, active: false },
    { name: '評估緩衝時間', completed: false, active: false },
    { name: '產生媒合結果', completed: false, active: false },
  ])

  // ==================== 人才庫相關狀態 ====================
  const talentResults = ref<TalentCard[]>([])
  const talentTotalItems = ref(0)
  const talentTotalPages = ref(0)
  const talentCurrentPage = ref(1)
  const hasSearchedTalent = ref(false)
  const talentStats = ref<TalentPoolStats | null>(null)
  const cityDistribution = ref<Array<{ name: string; count: number }>>([])
  const topSkills = ref<Array<{ name: string; count: number }>>([])

  // ==================== 教師課表相關狀態 ====================
  const selectedTeacher = ref<SmartMatchingResult | null>(null)
  const teacherSessions = ref<TeacherScheduleItem[]>([])
  const alternativeSlots = ref<AlternativeSlot[]>([])

  // ==================== 比較相關狀態 ====================
  const selectedForCompare = ref<Set<number>>(new Set())
  const viewMode = ref<'card' | 'list' | 'compare'>('card')

  // ==================== 邀請相關狀態 ====================
  const invitationStatuses = ref<Map<number, InvitationStatus>>(new Map())
  const inviteLoadingIds = ref<Set<number>>(new Set())
  const bulkLoading = ref(false)
  const bulkProgress = ref(0)

  // ==================== Loading 狀態 ====================
  const isSearching = ref(false)
  const isSearchingTalent = ref(false)
  const isFetchingStats = ref(false)
  const isFetchingSchedule = ref(false)
  const isInviting = ref(false)
  const isBulkInviting = ref(false)

  // ==================== 排序狀態 ====================
  const sortBy = ref<'score' | 'availability' | 'rating' | 'skill-match'>('score')
  const sortOrder = ref<'asc' | 'desc'>('desc')
  const talentSortBy = ref<'name' | 'skills' | 'rating' | 'city'>('name')
  const talentSortOrder = ref<'asc' | 'desc'>('asc')

  // ==================== Computed 屬性 ====================

  /**
   * 排序後的媒合結果
   */
  const sortedResults = computed(() => {
    const result = [...searchResults.value]

    result.sort((a, b) => {
      let comparison = 0

      switch (sortBy.value) {
        case 'score':
          comparison = a.match_score - b.match_score
          break
        case 'availability':
          const availOrder = { AVAILABLE: 0, BUSY: 1, UNAVAILABLE: 2 }
          comparison = (availOrder[a.availability] || 0) -
                       (availOrder[b.availability] || 0)
          break
        case 'rating':
          comparison = (a.rating || 0) - (b.rating || 0)
          break
        case 'skill-match':
          comparison = (a.score_detail?.match_score || 0) - (b.score_detail?.match_score || 0)
          break
      }

      return sortOrder.value === 'desc' ? -comparison : comparison
    })

    return result
  })

  /**
   * 選中比較的教師數量
   */
  const selectedCount = computed(() => selectedForCompare.value.size)

  // ==================== 搜尋相關方法 ====================

  /**
   * 重置搜尋進度
   */
  const resetSearchProgress = () => {
    searchProgress.value = 0
    searchSteps.value = [
      { name: '取得老師清單', completed: false, active: true },
      { name: '分析可用時間', completed: false, active: false },
      { name: '計算技能匹配度', completed: false, active: false },
      { name: '評估緩衝時間', completed: false, active: false },
      { name: '產生媒合結果', completed: false, active: false },
    ]
  }

  /**
   * 完成當前步驟並開始下一個
   */
  const completeSearchStep = (stepIndex: number) => {
    if (stepIndex < searchSteps.value.length) {
      searchSteps.value[stepIndex].completed = true
      searchSteps.value[stepIndex].active = false
      if (stepIndex + 1 < searchSteps.value.length) {
        searchSteps.value[stepIndex + 1].active = true
      }

      // 更新進度
      const completedCount = searchSteps.value.filter(s => s.completed).length
      searchProgress.value = Math.floor((completedCount / searchSteps.value.length) * 100)
    }
  }

  /**
   * 模擬搜尋進度
   */
  const simulateProgress = async () => {
    await new Promise(resolve => setTimeout(resolve, 200))
    completeSearchStep(0)

    await new Promise(resolve => setTimeout(resolve, 150))
    completeSearchStep(1)

    await new Promise(resolve => setTimeout(resolve, 150))
    completeSearchStep(2)

    await new Promise(resolve => setTimeout(resolve, 100))
    completeSearchStep(3)

    await new Promise(resolve => setTimeout(resolve, 100))
    completeSearchStep(4)
  }

  /**
   * 智慧媒合搜尋
   */
  const searchMatches = async (params: MatchingParams): Promise<SmartMatchingResult[]> => {
    return withLoading(isSearching, async () => {
      try {
        // 重置狀態
        hasSearched.value = true
        searchResults.value = []
        selectedForCompare.value.clear()
        viewMode.value = 'card'
        selectedTeacher.value = null
        resetSearchProgress()

        // 並行執行：API 呼叫 + 進度模擬
        const [response] = await Promise.all([
          api.post<SmartMatchingResponse>('/admin/smart-matching/matches', params),
          simulateProgress(),
        ])

        if (response.code === 0) {
          searchResults.value = response.datas || []
          return searchResults.value
        } else {
          console.error('搜尋失敗:', response.message)
          return []
        }
      } catch (error) {
        console.error('智慧媒合搜尋失敗:', error)
        return []
      }
    })
  }

  // ==================== 人才庫相關方法 ====================

  /**
   * 取得人才庫統計
   */
  const fetchTalentStats = async (): Promise<void> => {
    return withLoading(isFetchingStats, async () => {
      try {
        const response = await api.get<TalentPoolStatsResponse>(
          '/admin/smart-matching/talent/stats'
        )

        if (response.code === 0 && response.datas) {
          const stats = response.datas
          talentStats.value = stats
          cityDistribution.value = stats.city_distribution || []
          topSkills.value = stats.top_skills || []
        }
      } catch (error) {
        console.error('取得人才庫統計失敗:', error)
      }
    })
  }

  /**
   * 人才庫搜尋
   */
  const searchTalent = async (
    params: TalentSearchParams,
    page: number = 1
  ): Promise<{ results: TalentCard[]; totalItems: number; totalPages: number }> => {
    return withLoading(isSearchingTalent, async () => {
      try {
        hasSearchedTalent.value = true

        // 建立查詢參數
        const queryParams = new URLSearchParams()
        if (params.keyword) queryParams.append('keyword', params.keyword)
        if (params.city) queryParams.append('city', params.city)
        if (params.skill_name) queryParams.append('skill_name', params.skill_name)
        if (params.district) queryParams.append('district', params.district)
        if (params.open_to_hiring_only) queryParams.append('open_to_hiring_only', 'true')
        if (params.min_rating) queryParams.append('min_rating', params.min_rating.toString())
        if (params.sort_by) queryParams.append('sort_by', params.sort_by)
        if (params.sort_order) queryParams.append('sort_order', params.sort_order)
        queryParams.append('page', page.toString())
        queryParams.append('limit', (params.limit || 20).toString())

        const response = await api.get<TalentSearchResponse>(
          `/admin/smart-matching/talent/search?${queryParams.toString()}`
        )

        if (response.code === 0 && response.datas) {
          const data = response.datas
          talentResults.value = data.data || []
          talentTotalItems.value = data.pagination?.total || 0
          talentTotalPages.value = data.pagination?.total_pages || 0
          talentCurrentPage.value = page

          return {
            results: talentResults.value,
            totalItems: talentTotalItems.value,
            totalPages: talentTotalPages.value,
          }
        } else {
          talentResults.value = []
          talentTotalItems.value = 0
          talentTotalPages.value = 0
          return { results: [], totalItems: 0, totalPages: 0 }
        }
      } catch (error) {
        console.error('人才庫搜尋失敗:', error)
        talentResults.value = []
        talentTotalItems.value = 0
        talentTotalPages.value = 0
        return { results: [], totalItems: 0, totalPages: 0 }
      }
    })
  }

  /**
   * 邀請人才
   */
  const inviteTalent = async (
    teacherId: number,
    message?: string
  ): Promise<boolean> => {
    return withLoading(isInviting, async () => {
      try {
        inviteLoadingIds.value.add(teacherId)

        const response = await api.post<InviteTalentResponse>(
          '/admin/smart-matching/talent/invite',
          {
            teacher_ids: [teacherId],
            message: message || '',
          } as InviteTalentRequest
        )

        if (response.code === 0) {
          // 更新邀請狀態
          invitationStatuses.value.set(teacherId, {
            sent: true,
            variant: 'success',
            text: '已邀請',
          })
          return true
        } else {
          invitationStatuses.value.set(teacherId, {
            sent: false,
            variant: 'error',
            text: '邀請失敗',
          })
          return false
        }
      } catch (error) {
        console.error('邀請人才失敗:', error)
        invitationStatuses.value.set(teacherId, {
          sent: false,
          variant: 'error',
          text: '邀請失敗',
        })
        return false
      } finally {
        inviteLoadingIds.value.delete(teacherId)
      }
    })
  }

  /**
   * 批量邀請人才
   */
  const bulkInviteTalents = async (
    teacherIds: number[],
    message?: string
  ): Promise<{ success: number; failed: number }> => {
    return withLoading(isBulkInviting, async () => {
      try {
        bulkLoading.value = true
        bulkProgress.value = 0

        const response = await api.post<InviteTalentResponse>(
          '/admin/smart-matching/talent/invite',
          {
            teacher_ids: teacherIds,
            message: message || '',
          } as InviteTalentRequest
        )

        if (response.code === 0) {
          const result = response.datas
          const successCount = result.success_count || 0
          const failedCount = result.failed_count || 0

          // 更新邀請狀態
          teacherIds.forEach(teacherId => {
            if (result.failed_ids?.includes(teacherId)) {
              invitationStatuses.value.set(teacherId, {
                sent: false,
                variant: 'error',
                text: '邀請失敗',
              })
            } else {
              invitationStatuses.value.set(teacherId, {
                sent: true,
                variant: 'success',
                text: '已邀請',
              })
            }
          })

          bulkProgress.value = 100
          return { success: successCount, failed: failedCount }
        } else {
          // 標記所有選中為失敗
          teacherIds.forEach(teacherId => {
            invitationStatuses.value.set(teacherId, {
              sent: false,
              variant: 'error',
              text: '邀請失敗',
            })
          })
          return { success: 0, failed: teacherIds.length }
        }
      } catch (error) {
        console.error('批量邀請失敗:', error)
        teacherIds.forEach(teacherId => {
          invitationStatuses.value.set(teacherId, {
            sent: false,
            variant: 'error',
            text: '邀請失敗',
          })
        })
        return { success: 0, failed: teacherIds.length }
      } finally {
        bulkLoading.value = false
        bulkProgress.value = 0
      }
    })
  }

  // ==================== 教師課表相關方法 ====================

  /**
   * 取得教師課表
   */
  const fetchTeacherSchedule = async (
    params: TeacherScheduleParams
  ): Promise<TeacherScheduleItem[]> => {
    return withLoading(isFetchingSchedule, async () => {
      try {
        const queryParams = new URLSearchParams()
        queryParams.append('start_date', params.start_date)
        queryParams.append('end_date', params.end_date)

        const response = await api.get<TeacherScheduleResponse>(
          `/admin/teachers/${params.teacher_id}/sessions?${queryParams.toString()}`
        )

        if (response.code === 0 && response.datas) {
          teacherSessions.value = response.datas
          return teacherSessions.value
        } else {
          teacherSessions.value = []
          return []
        }
      } catch (error) {
        console.error('取得教師課表失敗:', error)
        teacherSessions.value = []
        return []
      }
    })
  }

  /**
   * 取得替代時段建議
   */
  const fetchAlternativeSlots = async (
    teacherId: number,
    originalDate: string,
    originalStart: string,
    originalEnd: string,
    duration: number = 90
  ): Promise<AlternativeSlot[]> => {
    try {
      const response = await api.post<AlternativeSlotsResponse>(
        '/admin/smart-matching/alternatives',
        {
          teacher_id: teacherId,
          original_date: originalDate,
          original_start_time: originalStart,
          original_end_time: originalEnd,
          duration,
        } as AlternativeSlotsRequest
      )

      if (response.code === 0 && response.datas) {
        alternativeSlots.value = response.datas
        return alternativeSlots.value
      } else {
        alternativeSlots.value = []
        return []
      }
    } catch (error) {
      console.error('取得替代時段失敗:', error)
      alternativeSlots.value = []
      return []
    }
  }

  // ==================== 比較相關方法 ====================

  /**
   * 切換比較選取
   */
  const toggleCompare = (teacherId: number, selected: boolean): boolean => {
    if (selected) {
      if (selectedForCompare.value.size >= 3) {
        return false
      }
      selectedForCompare.value.add(teacherId)
    } else {
      selectedForCompare.value.delete(teacherId)
    }
    return true
  }

  /**
   * 從比較中移除
   */
  const removeFromCompare = (teacherId: number) => {
    selectedForCompare.value.delete(teacherId)
  }

  /**
   * 清除所有比較選取
   */
  const clearCompare = () => {
    selectedForCompare.value.clear()
  }

  /**
   * 退出比較模式
   */
  const exitCompareMode = () => {
    viewMode.value = 'card'
  }

  // ==================== 選取相關方法 ====================

  /**
   * 選取教師
   */
  const selectTeacher = (teacher: SmartMatchingResult) => {
    selectedTeacher.value = teacher
  }

  /**
   * 清除選取
   */
  const clearSelection = () => {
    selectedTeacher.value = null
    teacherSessions.value = []
    alternativeSlots.value = []
  }

  // ==================== 重置相關方法 ====================

  /**
   * 重置搜尋結果
   */
  const resetSearchResults = () => {
    searchResults.value = []
    hasSearched.value = false
    selectedForCompare.value.clear()
    viewMode.value = 'card'
    selectedTeacher.value = null
    teacherSessions.value = []
    alternativeSlots.value = []
    resetSearchProgress()
  }

  /**
   * 重置人才庫搜尋
   */
  const resetTalentSearch = () => {
    talentResults.value = []
    talentTotalItems.value = 0
    talentTotalPages.value = 0
    talentCurrentPage.value = 1
    hasSearchedTalent.value = false
    invitationStatuses.value.clear()
  }

  /**
   * 重置所有狀態
   */
  const reset = () => {
    resetSearchResults()
    resetTalentSearch()
    talentStats.value = null
    cityDistribution.value = []
    topSkills.value = []
  }

  // ==================== 回傳 ====================

  return {
    // ==================== 搜尋相關狀態 ====================
    searchResults,
    hasSearched,
    searchProgress,
    searchSteps,
    sortedResults,

    // ==================== 人才庫相關狀態 ====================
    talentResults,
    talentTotalItems,
    talentTotalPages,
    talentCurrentPage,
    hasSearchedTalent,
    talentStats,
    cityDistribution,
    topSkills,

    // ==================== 教師課表相關狀態 ====================
    selectedTeacher,
    teacherSessions,
    alternativeSlots,

    // ==================== 比較相關狀態 ====================
    selectedForCompare,
    viewMode,
    selectedCount,

    // ==================== 邀請相關狀態 ====================
    invitationStatuses,
    inviteLoadingIds,
    bulkLoading,
    bulkProgress,

    // ==================== Loading 狀態 ====================
    isSearching,
    isSearchingTalent,
    isFetchingStats,
    isFetchingSchedule,
    isInviting,
    isBulkInviting,

    // ==================== 排序狀態 ====================
    sortBy,
    sortOrder,
    talentSortBy,
    talentSortOrder,

    // ==================== 搜尋相關方法 ====================
    searchMatches,
    resetSearchProgress,
    completeSearchStep,

    // ==================== 人才庫相關方法 ====================
    fetchTalentStats,
    searchTalent,
    inviteTalent,
    bulkInviteTalents,

    // ==================== 教師課表相關方法 ====================
    fetchTeacherSchedule,
    fetchAlternativeSlots,

    // ==================== 比較相關方法 ====================
    toggleCompare,
    removeFromCompare,
    clearCompare,
    exitCompareMode,

    // ==================== 選取相關方法 ====================
    selectTeacher,
    clearSelection,

    // ==================== 重置相關方法 ====================
    resetSearchResults,
    resetTalentSearch,
    reset,
  }
})
