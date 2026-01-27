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
        <div class="relative">
          <label for="weekday-filter" class="sr-only">篩選星期</label>
          <select
            id="weekday-filter"
            v-model="filterWeekday"
            aria-label="篩選星期"
            class="px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white focus:outline-none focus:border-primary-500/50 appearance-none"
          >
          <option value="">全部星期</option>
          <option v-for="(day, index) in ['週日', '週一', '週二', '週三', '週四', '週五', '週六']" :key="index" :value="index === 0 ? 7 : index">
            {{ day }}
          </option>
        </select>
        <!-- 狀態篩選 -->
        <div class="relative">
          <label for="status-filter" class="sr-only">篩選狀態</label>
          <select
            id="status-filter"
            v-model="filterStatus"
            aria-label="篩選課程狀態"
            class="px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white focus:outline-none focus:border-primary-500/50 appearance-none"
          >
            <option value="">全部狀態</option>
            <option value="upcoming">尚未開始</option>
            <option value="ongoing">進行中</option>
            <option value="ended">已結束</option>
          </select>
        </div>
        <!-- 清除篩選 -->
        <button
          v-if="searchQuery || filterWeekday || filterStatus"
          @click="clearFilters"
          aria-label="清除所有篩選條件"
          class="px-4 py-2 text-slate-400 hover:text-white transition-colors"
        >
          清除篩選
        </button>
      </div>
      <!-- 篩選結果計數 -->
      <div v-if="rules.length > 0" class="mt-2 text-sm text-slate-500" role="status" aria-live="polite">
        顯示 {{ filteredRules.length }} / {{ rules.length }} 筆資料
      </div>
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

      <div v-else class="overflow-x-auto -mx-6">
        <table class="w-full min-w-[600px]" role="table" aria-label="課程時段列表">
          <thead class="sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
            <tr class="text-left text-slate-400 text-sm border-b border-white/10">
              <th class="pb-3 pl-4" scope="col">課程</th>
              <th class="pb-3" scope="col">星期</th>
              <th class="pb-3" scope="col">時間</th>
              <th class="pb-3" scope="col">教室</th>
              <th class="pb-3" scope="col">老師</th>
              <th class="pb-3" scope="col">狀態</th>
              <th class="pb-3 pr-4 text-right" scope="col">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="rule in filteredRules"
              :key="rule.id"
              class="border-b border-white/5 hover:bg-white/5 transition-colors"
            >
              <td class="py-3 pl-4 text-white">{{ rule.offering?.name || '-' }}</td>
              <td class="py-3 text-slate-300">{{ getWeekdayText(rule.weekday) }}</td>
              <td class="py-3 text-slate-300">{{ rule.start_time }} - {{ rule.end_time }}</td>
              <td class="py-3 text-slate-300">{{ rule.room?.name || '-' }}</td>
              <td class="py-3 text-slate-300">{{ rule.teacher?.name || '-' }}</td>
              <td class="py-3">
                <span
                  class="px-2 py-1 rounded-full text-xs"
                  :class="getStatusClass(rule)"
                >
                  {{ getStatusText(rule) }}
                </span>
              </td>
              <td class="py-3 pr-4 text-right">
                <button
                  @click="editRule(rule)"
                  aria-label="編輯課程時段"
                  class="text-primary-500 hover:text-primary-400 mr-3"
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
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <ScheduleRuleModal
    v-if="showModal"
    :editing-rule="editingRule"
    @close="handleModalClose"
    @submit="handleModalSubmit"
  />

  <UpdateModeModal
    v-if="showUpdateModeModal"
    :show="showUpdateModeModal"
    :rule-date="editingRule ? new Date(editingRule.effective_range?.start_date).toLocaleDateString('zh-TW', { year: 'numeric', month: 'long', day: 'numeric' }) : ''"
    @close="showUpdateModeModal = false; showModal = true; pendingEditData = null"
    @confirm="handleUpdateModeConfirm"
  />

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
  </div>
</template>

<script setup lang="ts">
 definePageMeta({
   middleware: 'auth-admin',
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

// Alert composable
const { error: alertError, confirm: alertConfirm } = useAlert()

// 清除所有篩選條件
const clearFilters = () => {
  searchQuery.value = ''
  filterWeekday.value = ''
  filterStatus.value = ''
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
      const startDate = new Date(rule.effective_range?.start_date)
      const endDate = rule.effective_range?.end_date ? new Date(rule.effective_range.end_date) : null

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
    const response = await api.get<{ code: number; datas: any[] }>('/admin/rules')
    rules.value = response.datas || []
  } catch (error) {
    console.error('Failed to fetch rules:', error)
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

const handleModalSubmit = (formData: any) => {
  // 如果編輯模式下有修改日期相關內容，需要詢問更新模式
  if (editingRule.value && formData.start_date) {
    const originalStartDate = editingRule.value.effective_range?.start_date?.split('T')[0]
    if (originalStartDate && originalStartDate !== formData.start_date) {
      // 日期有變更，顯示更新模式選擇
      const ruleDate = new Date(editingRule.value.effective_range?.start_date).toLocaleDateString('zh-TW', {
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
    await api.put(`/admin/rules/${editingRule.value.id}`, formData)
    await fetchRules()
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
  const startDate = new Date(rule.effective_range?.start_date)
  const endDate = rule.effective_range?.end_date ? new Date(rule.effective_range.end_date) : null

  if (endDate && now > endDate) return 'bg-slate-500/20 text-slate-400'
  if (now < startDate) return 'bg-primary-500/20 text-primary-500'
  return 'bg-success-500/20 text-success-500'
}

const getStatusText = (rule: any): string => {
  const now = new Date()
  const startDate = new Date(rule.effective_range?.start_date)
  const endDate = rule.effective_range?.end_date ? new Date(rule.effective_range.end_date) : null

  if (endDate && now > endDate) return '已結束'
  if (now < startDate) return '尚未開始'
  return '進行中'
}

onMounted(() => {
  fetchRules()
})
</script>
