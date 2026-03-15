<template>
  <div class="p-4 md:p-6">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold text-white">課程時段管理</h1>
      <button
        @click="showModal = true"
        class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors"
      >
        新增時段
      </button>
    </div>

    <!-- 搜尋與篩選區域 -->
    <div class="glass-card p-4 mb-6">
      <div class="flex flex-col md:flex-row gap-4">
        <!-- 搜尋框 -->
        <div class="flex-1">
          <label for="search-input" class="sr-only">搜尋課程、老師或教室名稱</label>
          <input
            id="search-input"
            v-model="searchQuery"
            type="text"
            placeholder="搜尋課程、老師或教室名稱..."
            aria-label="搜尋課程、老師或教室名稱"
            class="w-full px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-primary-500/50"
          />
        </div>
        <!-- 星期篩選 -->
        <div>
          <label for="weekday-filter" class="sr-only">篩選星期</label>
          <select
            id="weekday-filter"
            v-model="filterWeekday"
            aria-label="篩選星期"
            class="px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white focus:outline-none focus:border-primary-500/50 appearance-none w-full md:w-auto"
          >
          <option value="">全部星期</option>
          <option v-for="(day, index) in ['週日', '週一', '週二', '週三', '週四', '週五', '週六']" :key="index" :value="index === 0 ? 7 : index">
            {{ day }}
          </option>
        </select>
        </div>
        <!-- 狀態篩選 -->
        <div>
          <label for="status-filter" class="sr-only">篩選狀態</label>
          <select
            id="status-filter"
            v-model="filterStatus"
            aria-label="篩選課程狀態"
            class="px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white focus:outline-none focus:border-primary-500/50 appearance-none w-full md:w-auto"
          >
            <option value="">全部狀態</option>
            <option value="upcoming">尚未開始</option>
            <option value="ongoing">進行中</option>
            <option value="ended">已結束</option>
          </select>
        </div>
        <!-- 類別篩選 -->
        <div>
          <label for="category-filter" class="sr-only">篩選類別</label>
          <select
            id="category-filter"
            v-model="filterCategory"
            aria-label="篩選課程類別"
            class="px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white focus:outline-none focus:border-primary-500/50 appearance-none w-full md:w-auto"
          >
            <option value="">全部類別</option>
            <option v-for="cat in categories" :key="cat" :value="cat">
              {{ cat }}
            </option>
          </select>
        </div>
        <!-- 清除篩選 -->
        <button
          v-if="searchQuery || filterWeekday || filterStatus || filterCategory"
          @click="clearFilters"
          aria-label="清除所有篩選條件"
          class="px-4 py-2 text-slate-400 hover:text-white transition-colors"
        >
          清除篩選
        </button>
      </div>
    </div>

    <!-- 篩選結果計數 -->
    <div v-if="rules.length > 0" class="mb-4 text-sm text-slate-500 text-right" role="status" aria-live="polite">
      共 {{ totalCount }} 筆資料
    </div>

    <div class="glass-card p-6" role="region" aria-label="課程時段列表">
      <div v-if="loading" class="text-center py-8 text-slate-400" role="status" aria-live="polite">
        載入中...
      </div>

      <div v-else-if="rules.length === 0" class="text-center py-8 text-slate-400" role="status">
        尚未建立課程時段
      </div>

      <div v-else-if="filteredRules.length === 0" class="text-center py-8 text-slate-400" role="status">
        沒有符合搜尋條件的課程時段
      </div>

      <div v-else class="overflow-x-auto -mx-4 sm:-mx-6 px-4 sm:px-6">
        <table class="w-full min-w-[800px]" role="table" aria-label="課程時段列表">
          <thead class="bg-white/5">
            <tr class="text-slate-400 text-sm border-b border-white/10">
              <th class="p-3 text-center w-20" scope="col">課程代號</th>
              <th class="p-3 text-center w-28" scope="col">課程</th>
              <th class="p-3 text-center w-16" scope="col">星期</th>
              <th class="p-3 text-center w-36" scope="col">課程期間</th>
              <th class="p-3 text-center w-28" scope="col">課程時間</th>
              <th class="p-3 text-center w-24" scope="col">教室</th>
              <th class="p-3 text-center w-24" scope="col">老師</th>
              <th class="p-3 text-center w-20" scope="col">狀態</th>
              <th class="p-3 text-center w-24" scope="col">課程類別</th>
              <th class="p-3 text-center w-20" scope="col">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="rule in filteredRules"
              :key="rule.id"
              class="border-b border-white/5 hover:bg-white/5 transition-colors"
            >
              <td class="p-3 text-center text-slate-300 font-mono text-sm">{{ rule.code || '-' }}</td>
              <td class="p-3 text-center text-slate-200">{{ rule.offering?.name || '-' }}</td>
              <td class="p-3 text-center text-slate-300">{{ getWeekdayText(rule.weekday) }}</td>
              <td class="p-3 text-center text-slate-300">{{ formatDateRange(rule.effective_range) }}</td>
              <td class="p-3 text-center text-slate-300">{{ rule.start_time }} - {{ rule.end_time }}</td>
              <td class="p-3 text-center text-slate-300">{{ rule.room?.name || '-' }}</td>
              <td class="p-3 text-center text-slate-300">{{ rule.teacher?.name || '-' }}</td>
              <td class="p-3 text-center">
                <span
                  class="px-2 py-1 rounded-full text-xs"
                  :class="getStatusClass(rule)"
                >
                  {{ getStatusText(rule) }}
                </span>
              </td>
              <td class="p-3 text-center text-slate-300">{{ rule.offering?.course?.category || '-' }}</td>
              <td class="p-3">
                <div class="flex items-center justify-center gap-3">
                  <button
                    @click="editRule(rule)"
                    aria-label="編輯課程時段"
                    class="text-primary-500 hover:text-primary-400"
                  >
                    編輯
                  </button>
                  <button
                    @click="deleteRule(rule.id)"
                    :aria-label="'刪除課程時段 ' + (rule.offering?.name || '')"
                    class="text-critical-500 hover:text-critical-400"
                  >
                    刪除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- 分頁控制 -->
        <div v-if="totalPages >= 1" class="flex items-center justify-between mt-4 pt-4 border-t border-white/10">
          <div class="text-sm text-slate-400">
            每頁 {{ pageSize }} 筆
          </div>
          <div class="flex items-center gap-1">
            <button
              @click="changePage(currentPage - 1)"
              :disabled="currentPage === 1"
              class="px-3 py-1.5 rounded-lg text-sm bg-white/5 text-slate-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              上一頁
            </button>

            <template v-for="page in visiblePages" :key="page">
              <span v-if="page === '...'" class="px-2 text-slate-500">...</span>
              <button
                v-else
                @click="changePage(page as number)"
                class="px-3 py-1.5 rounded-lg text-sm"
                :class="page === currentPage
                  ? 'bg-primary-500 text-white'
                  : 'bg-white/5 text-slate-300 hover:bg-white/10'"
              >
                {{ page }}
              </button>
            </template>

            <button
              @click="changePage(currentPage + 1)"
              :disabled="currentPage === totalPages"
              class="px-3 py-1.5 rounded-lg text-sm bg-white/5 text-slate-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              下一頁
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <ScheduleRuleModal
    v-if="showModal"
    :editing-rule="editingRule"
    @close="handleModalClose"
    @submit="handleModalSubmit"
    @created="handleRuleCreated"
  />

  <UpdateModeModal
    v-if="showUpdateModeModal"
    :show="showUpdateModeModal"
    :rule-date="editingRule ? new Date(editingRule.effective_from || editingRule.effective_range?.start_date).toLocaleDateString('zh-TW', { year: 'numeric', month: 'long', day: 'numeric' }) : ''"
    @close="showUpdateModeModal = false; showModal = true; pendingEditData = null"
    @confirm="handleUpdateModeConfirm"
  />

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
// 明確導入組件（確保 Nuxt 可以解析）
import ScheduleRuleModal from '~/components/Scheduling/ScheduleRuleModal.vue'
import UpdateModeModal from '~/components/Scheduling/UpdateModeModal.vue'
import NotificationDropdown from '~/components/Navigation/NotificationDropdown.vue'
import { watch } from 'vue'

definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

 const notificationUI = useNotification()
const showModal = ref(false)
const loading = ref(true)
const rules = ref<any[]>([])
const editingRule = ref<any | null>(null)
const showUpdateModeModal = ref(false)
const pendingEditData = ref<any>(null)
const { getCenterId } = useCenterId()

// 搜尋與篩選
const searchQuery = ref('')
const filterWeekday = ref('')
const filterStatus = ref('')
const filterCategory = ref('')

// 課程類別
const categories = ref<string[]>([])
const fetchCategories = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${window.location.origin}/api/v1/admin/course-categories`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    if (response.ok) {
      const data = await response.json()
      // 轉換為名稱陣列
      categories.value = (data.datas || []).map((c: any) => c.name)
    }
  } catch (err) {
    console.error('取得課程類別失敗:', err)
  }
}

// 分頁狀態
const currentPage = ref(1)
const totalPages = ref(0)
const totalCount = ref(0)
const pageSize = ref(20)

// 改變頁碼
const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchRules()
  }
}

// 計算顯示的頁碼陣列
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value
  const pages: (number | string)[] = []

  if (total <= 7) {
    // 總頁數 <= 7，顯示全部
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    // 總頁數 > 7，需要顯示省略號
    if (current <= 4) {
      // 靠近開頭：1, 2, 3, 4, 5, ..., last
      for (let i = 1; i <= 5; i++) pages.push(i)
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      // 靠近結尾：1, ..., last-4, last-3, last-2, last-1, last
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) pages.push(i)
    } else {
      // 中間：1, ..., current-1, current, current+1, ..., last
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) pages.push(i)
      pages.push('...')
      pages.push(total)
    }
  }

  return pages
})

// 監聽篩選變化，重置到第一頁
watch([searchQuery, filterWeekday, filterStatus, filterCategory], () => {
  currentPage.value = 1
  fetchRules()
})

// Alert composable
const { error: alertError, confirm: alertConfirm } = useAlert()

// 清除所有篩選條件
const clearFilters = () => {
  searchQuery.value = ''
  filterWeekday.value = ''
  filterStatus.value = ''
  filterCategory.value = ''
}

// 篩選後的規則列表
const filteredRules = computed(() => {
  let result = [...rules.value]

  // 搜尋過濾
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(rule =>
      rule.offering?.name?.toLowerCase().includes(query) ||
      rule.teacher?.name?.toLowerCase().includes(query) ||
      rule.room?.name?.toLowerCase().includes(query)
    )
  }

  // 星期過濾
  if (filterWeekday.value) {
    const weekdayValue = parseInt(filterWeekday.value)
    result = result.filter(rule => rule.weekday === weekdayValue)
  }

  // 狀態過濾
  if (filterStatus.value) {
    const now = new Date()
    result = result.filter(rule => {
      const startDate = new Date(rule.effective_from || rule.effective_range?.start_date)
      const endDate = (rule.effective_to || rule.effective_range?.end_date) ? new Date(rule.effective_to || rule.effective_range?.end_date) : null

      switch (filterStatus.value) {
        case 'upcoming':
          return now < startDate
        case 'ongoing':
          return now >= startDate && (!endDate || now <= endDate)
        case 'ended':
          return endDate && now > endDate
        default:
          return true
      }
    })
  }

  return result
})

const fetchRules = async () => {
  loading.value = true
  try {
    const api = useApi()
    const params: any = {
      page: currentPage.value,
      limit: pageSize.value
    }

    // 加入類別篩選參數
    if (filterCategory.value) {
      params.category = filterCategory.value
    }

    const response = await api.get<any>('/admin/rules', params)

    // 分頁格式：{ data: [...], total: X, page: X, total_pages: X }
    // useApi 已自動提取 data.data 或 data.datas
    if (response && response.data) {
      rules.value = response.data || []
      totalCount.value = response.total || 0
      totalPages.value = response.total_pages || 1
    } else {
      rules.value = []
      totalCount.value = 0
      totalPages.value = 0
    }
  } catch (error) {
    console.error('Failed to fetch rules:', error)
    rules.value = []
    totalCount.value = 0
    totalPages.value = 0
  } finally {
    loading.value = false
  }
}

const deleteRule = async (id: number) => {
  if (!await alertConfirm('確定要刪除此課程時段？')) return

  try {
    const api = useApi()
    await api.delete(`/admin/rules/${id}`)
    await fetchRules()
  } catch (err) {
    console.error('Failed to delete rule:', err)
    await alertError('刪除失敗，請稍後再試')
  }
}

const editRule = (rule: any) => {
  editingRule.value = rule
  showModal.value = true
}

const handleUpdateModeConfirm = async (updateMode: string) => {
  if (!pendingEditData.value || !updateMode) return

  try {
    const api = useApi()
    await api.put(`/admin/rules/${pendingEditData.value.id}`, {
      ...pendingEditData.value.formData,
      update_mode: updateMode,
      // 排除自己，避免與自己衝突
      exclude_rule_id: pendingEditData.value.id,
    })
    await fetchRules()
    showUpdateModeModal.value = false
    pendingEditData.value = null
    editingRule.value = null
    showModal.value = false
  } catch (err) {
    console.error('Failed to update rule:', err)
    await alertError('更新失敗，請稍後再試')
  }
}

const handleModalClose = () => {
  showModal.value = false
  editingRule.value = null
}

// 處理新規則建立事件
const handleRuleCreated = async () => {
  await fetchRules()
  // 清除資源快取，確保下次開啟 Modal 時載入最新資料
  const { invalidate } = useResourceCache()
  invalidate()
}

const handleModalSubmit = (formData: any) => {
  // 如果編輯模式下有修改日期相關內容，需要詢問更新模式
  if (editingRule.value && formData.start_date) {
    const originalStartDate = (editingRule.value.effective_from || editingRule.value.effective_range?.start_date)?.split('T')[0]
    if (originalStartDate && originalStartDate !== formData.start_date) {
      // 日期有變更，顯示更新模式選擇
      const ruleDate = new Date(editingRule.value.effective_from || editingRule.value.effective_range?.start_date).toLocaleDateString('zh-TW', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
      pendingEditData.value = {
        id: editingRule.value.id,
        formData: formData,
      }
      showModal.value = false
      showUpdateModeModal.value = true
      return
    }
  }

  // 無需詢問更新模式，直接提交
  submitDirectly(formData)
}

const submitDirectly = async (formData: any) => {
  try {
    const api = useApi()
    await api.put(`/admin/rules/${editingRule.value.id}`, {
      ...formData,
      // 排除自己，避免與自己衝突
      exclude_rule_id: editingRule.value.id,
    })
    await fetchRules()
    // 清除資源快取，確保下次開啟 Modal 時載入最新資料
    const { invalidate } = useResourceCache()
    invalidate()
    showModal.value = false
    editingRule.value = null
  } catch (err) {
    console.error('Failed to update rule:', err)
    await alertError('更新失敗，請稍後再試')
  }
}

const getWeekdayText = (weekday: number): string => {
  const days = ['日', '一', '二', '三', '四', '五', '六']
  // 我們的系統使用 7 表示週日，但 JavaScript 的 Date.getDay() 返回 0
  const dayIndex = weekday === 7 ? 0 : weekday
  return days[dayIndex] || '-'
}

const getStatusClass = (rule: any): string => {
  const now = new Date()
  const startDate = new Date(rule.effective_from || rule.effective_range?.start_date)
  const endDate = (rule.effective_to || rule.effective_range?.end_date) ? new Date(rule.effective_to || rule.effective_range?.end_date) : null

  if (endDate && now > endDate) return 'bg-slate-500/20 text-slate-400'
  if (now < startDate) return 'bg-primary-500/20 text-primary-500'
  return 'bg-success-500/20 text-success-500'
}

const getStatusText = (rule: any): string => {
  const now = new Date()
  const startDate = new Date(rule.effective_from || rule.effective_range?.start_date)
  const endDate = (rule.effective_to || rule.effective_range?.end_date) ? new Date(rule.effective_to || rule.effective_range?.end_date) : null

  if (endDate && now > endDate) return '已結束'
  if (now < startDate) return '尚未開始'
  return '進行中'
}

const formatDateRange = (effectiveRange: any): string => {
  // 支援新格式 (effective_from/effective_to) 和舊格式 (effective_range.start_date/end_date)
  const startDate = effectiveRange?.start_date || effectiveRange?.effective_from
  if (!startDate) return '-'

  const start = new Date(startDate)
  const startStr = start.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })

  const endDate = effectiveRange?.end_date || effectiveRange?.effective_to
  if (!endDate) {
    return `${startStr} 起`
  }

  const endDate = new Date(effectiveRange.end_date)
  const endStr = endDate.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })

  return `${startStr} ~ ${endStr}`
}

onMounted(() => {
  fetchRules()
  fetchCategories()
})
</script>
