<template>
  <div class="p-4 md:p-6 max-w-7xl mx-auto">
    <div class="mb-6 md:mb-8 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
          假日管理
        </h1>
        <p class="text-slate-400 text-sm md:text-base">
          管理中心的國定假日與特殊假日
        </p>
      </div>
      <div class="flex gap-3">
        <button
          @click="showBulkModal = true"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
        >
          <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
          </svg>
          批次匯入
        </button>
        <button
          @click="showAddModal = true"
          class="px-4 py-2 rounded-xl bg-primary-500 text-white text-sm font-medium hover:bg-primary-600 transition-colors whitespace-nowrap"
        >
          <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          新增假日
        </button>
      </div>
    </div>

    <!-- 月份選擇 -->
    <div class="mb-6 flex flex-col sm:flex-row gap-4 items-center justify-between">
      <div class="flex items-center gap-3">
        <button
          @click="prevMonth"
          class="p-2 rounded-lg hover:bg-white/10 transition-colors"
        >
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h2 class="text-lg font-semibold text-white min-w-[140px] text-center">
          {{ currentYear }} 年 {{ currentMonth + 1 }} 月
        </h2>
        <button
          @click="nextMonth"
          class="p-2 rounded-lg hover:bg-white/10 transition-colors"
        >
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>
      <div class="flex gap-2">
        <button
          @click="goToToday"
          class="px-3 py-1.5 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-white/10 transition-colors"
        >
          今天
        </button>
        <select
          v-model="selectedYear"
          @change="fetchHolidays"
          class="px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-white text-sm"
        >
          <option v-for="year in yearOptions" :key="year" :value="year">{{ year }}年</option>
        </select>
      </div>
    </div>

    <!-- 假日列表 -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
    </div>

    <div v-else-if="currentMonthHolidays.length === 0" class="text-center py-16">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-white/5 flex items-center justify-center">
        <svg class="w-8 h-8 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      </div>
      <p class="text-slate-500 mb-4">目前沒有設定假日</p>
      <button
        @click="showAddModal = true"
        class="px-4 py-2 rounded-lg bg-primary-500 text-white text-sm hover:bg-primary-600 transition-colors"
      >
        新增假日
      </button>
    </div>

    <div v-else class="grid gap-3">
      <div
        v-for="holiday in currentMonthHolidays"
        :key="holiday.id"
        class="glass-card p-4 flex items-center justify-between"
      >
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-warning-500/20 flex items-center justify-center">
            <span class="text-warning-500 font-bold">{{ new Date(holiday.date).getDate() }}</span>
          </div>
          <div>
            <h3 class="text-white font-medium flex items-center gap-2">
              {{ holiday.name }}
              <span v-if="holiday.force_cancel" class="px-2 py-0.5 rounded-full text-xs font-medium bg-critical-500/20 text-critical-500">
                強制停課
              </span>
            </h3>
            <p class="text-sm text-slate-400">{{ formatFullDate(holiday.date) }}</p>
          </div>
        </div>
        <button
          @click="deleteHoliday(holiday.id)"
          class="p-2 rounded-lg hover:bg-critical-500/20 text-slate-400 hover:text-critical-500 transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 統計資訊 -->
    <div v-if="currentMonthHolidays.length > 0" class="mt-6 p-4 rounded-xl bg-white/5 border border-white/10">
      <div class="flex items-center justify-between text-sm">
        <span class="text-slate-400">本月已設定 {{ currentMonthHolidays.length }} 天假日</span>
        <span class="text-slate-500">提醒：假日期間的課程將自動停課</span>
      </div>
    </div>
  </div>

  <!-- 新增假日 Modal -->
  <Teleport to="body">
    <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" @click.self="showAddModal = false">
      <div class="glass-card w-full max-w-md">
        <div class="flex items-center justify-between p-4 border-b border-white/10">
          <h3 class="text-lg font-semibold text-white">新增假日</h3>
          <button @click="showAddModal = false" class="p-2 rounded-lg hover:bg-white/10">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <form @submit.prevent="handleAddHoliday" class="p-4 space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">假日名稱</label>
            <input
              v-model="addForm.name"
              type="text"
              placeholder="例如：春節、農曆新年"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500"
              required
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">日期</label>
            <input
              v-model="addForm.date"
              type="date"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500"
              required
            />
          </div>
          <div class="flex items-center gap-3">
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                v-model="addForm.force_cancel"
                type="checkbox"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-white/10 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-primary-500/30 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-critical-500"></div>
              <span class="ms-3 text-sm font-medium text-slate-300">強制停課</span>
            </label>
            <span class="text-xs text-slate-500">啟用後，假日期間的所有課程將自動取消</span>
          </div>
          <div class="flex gap-3 pt-4">
            <button type="button" @click="showAddModal = false" class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors">
              取消
            </button>
            <button type="submit" :disabled="saving" class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50">
              {{ saving ? '儲存中...' : '儲存' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>

  <!-- 批次匯入 Modal -->
  <Teleport to="body">
    <div v-if="showBulkModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" @click.self="showBulkModal = false">
      <div class="glass-card w-full max-w-lg">
        <div class="flex items-center justify-between p-4 border-b border-white/10">
          <h3 class="text-lg font-semibold text-white">批次匯入假日</h3>
          <button @click="showBulkModal = false" class="p-2 rounded-lg hover:bg-white/10">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-4 space-y-4">
          <div class="p-4 rounded-xl bg-blue-500/10 border border-blue-500/30">
            <p class="text-sm text-blue-400 mb-2">
              支援 JSON 格式，每筆需包含 <code class="px-1 bg-blue-500/20 rounded">date</code> 與 <code class="px-1 bg-blue-500/20 rounded">name</code> 欄位，可選 <code class="px-1 bg-blue-500/20 rounded">force_cancel</code> 強制停課
            </p>
            <p class="text-xs text-blue-300">
              💡 提示：請使用標準的半形雙引號 (")，不要使用中文引號「」或全形引號
            </p>
            <button
              @click="bulkForm.jsonData = defaultHolidayJSON"
              class="mt-2 text-xs text-blue-400 hover:text-blue-300 underline"
            >
              填入範例資料
            </button>
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">假日列表 (JSON)</label>
            <textarea
              v-model="bulkForm.jsonData"
              rows="10"
              placeholder='[
  {"date": "2026-01-01", "name": "元旦"},
  {"date": "2026-02-11", "name": "春節"},
  {"date": "2026-10-10", "name": "國慶日", "force_cancel": true}
]'
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white font-mono text-sm focus:outline-none focus:border-primary-500 resize-none"
            ></textarea>
          </div>
          <div class="flex gap-3 pt-4">
            <button @click="showBulkModal = false" class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors">
              取消
            </button>
            <button @click="handleBulkImport" :disabled="saving" class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50">
              {{ saving ? '匯入中...' : '匯入' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
import NotificationDropdown from '~/components/Navigation/NotificationDropdown.vue'
definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

const notificationUI = useNotification()
const { confirm: alertConfirm } = useAlert()
const { getCenterId } = useCenterId()
const api = useApi()

const loading = ref(false)
const saving = ref(false)
const showAddModal = ref(false)
const showBulkModal = ref(false)
const currentMonth = ref(new Date().getMonth())
const selectedYear = ref(new Date().getFullYear())
const holidays = ref<any[]>([])

const addForm = reactive({
  name: '',
  date: '',
  force_cancel: false
})

const bulkForm = reactive({
  jsonData: ''
})

// 預設的假日 JSON 範例
const defaultHolidayJSON = JSON.stringify([
  { "date": "2026-01-01", "name": "元旦" },
  { "date": "2026-02-11", "name": "春節" },
  { "date": "2026-02-12", "name": "春節初二" },
  { "date": "2026-02-13", "name": "春節初三" },
  { "date": "2026-02-14", "name": "春節初四" },
  { "date": "2026-02-28", "name": "和平紀念日" },
  { "date": "2026-04-03", "name": "清明節" },
  { "date": "2026-04-04", "name": "清明節連假" },
  { "date": "2026-05-01", "name": "勞動節" },
  { "date": "2026-05-05", "name": "端午節" },
  { "date": "2026-06-20", "name": "中秋節" },
  { "date": "2026-10-10", "name": "國慶日", "force_cancel": true }
], null, 2)

const currentYear = computed(() => selectedYear.value)

// 修復：新增 computed property 來過濾當月假日
const currentMonthHolidays = computed(() => {
  return holidays.value.filter(holiday => {
    const holidayDate = new Date(holiday.date)
    return holidayDate.getMonth() === currentMonth.value
  })
})

const yearOptions = computed(() => {
  const years = []
  const current = new Date().getFullYear()
  for (let i = current - 1; i <= current + 2; i++) {
    years.push(i)
  }
  return years
})

const fetchHolidays = async () => {
  loading.value = true
  try {
    const centerId = getCenterId()
    const startDate = `${selectedYear.value}-01-01`
    const endDate = `${selectedYear.value}-12-31`
    // parseResponse 已經提取了 datas 欄位，所以 response 就是假日陣列本身
    const response = await api.get<any[]>(
      `/admin/centers/${centerId}/holidays?start_date=${startDate}&end_date=${endDate}`
    )
    holidays.value = response || []
  } catch (error) {
    console.error('Failed to fetch holidays:', error)
    notificationUI.error('載入假日失敗')
    holidays.value = []
  } finally {
    loading.value = false
  }
}

const handleAddHoliday = async () => {
  saving.value = true
  try {
    const centerId = getCenterId()
    await api.post(`/admin/centers/${centerId}/holidays`, {
      date: addForm.date,
      name: addForm.name,
      force_cancel: addForm.force_cancel
    })
    notificationUI.success('假日已新增')
    showAddModal.value = false
    addForm.name = ''
    addForm.date = ''
    addForm.force_cancel = false
    await fetchHolidays()
  } catch (error) {
    console.error('Failed to add holiday:', error)
    notificationUI.error('新增假日失敗')
  } finally {
    saving.value = false
  }
}

const handleBulkImport = async () => {
  saving.value = true
  try {
    let holidaysData
    try {
      holidaysData = JSON.parse(bulkForm.jsonData)
    } catch (e: any) {
      console.error('JSON parse error:', e)
      // 提供更詳細的錯誤訊息
      let errorMsg = 'JSON 格式錯誤'
      if (e.message) {
        errorMsg += `：${e.message}`
      }
      // 檢查常見問題
      const text = bulkForm.jsonData
      if (text.includes('「') || text.includes('」') || text.includes('"')) {
        errorMsg += '。注意：請使用標準的雙引號 (")，不要使用中文引號「」或全形引號'
      }
      notificationUI.error(errorMsg)
      saving.value = false
      return
    }

    if (!Array.isArray(holidaysData)) {
      notificationUI.error('資料必須是陣列格式')
      saving.value = false
      return
    }

    // 驗證每筆資料
    for (let i = 0; i < holidaysData.length; i++) {
      const item = holidaysData[i]
      if (!item.date || !item.name) {
        notificationUI.error(`第 ${i + 1} 筆資料缺少必要欄位 (date, name)`)
        saving.value = false
        return
      }
    }

    const centerId = getCenterId()

    const response = await api.post(`/admin/centers/${centerId}/holidays/bulk`, { holidays: holidaysData })
    notificationUI.success(`已成功匯入 ${holidaysData.length} 天假日`)
    showBulkModal.value = false
    bulkForm.jsonData = ''
    await fetchHolidays()
  } catch (error: any) {
    notificationUI.error(`批次匯入失敗：${error?.message || '未知錯誤'}`)
    notificationUI.error(error?.message?.includes('404') ? 'API 端點不存在，請聯絡管理員' : '匯入失敗，請稍後再試')
  } finally {
    saving.value = false
  }
}

const deleteHoliday = async (id: number) => {
  if (!await alertConfirm('確定要刪除這個假日嗎？')) return

  try {
    const centerId = getCenterId()
    await api.delete(`/admin/centers/${centerId}/holidays/${id}`)
    notificationUI.success('假日已刪除')
    await fetchHolidays()
  } catch (error) {
    console.error('Failed to delete holiday:', error)
    notificationUI.error('刪除假日失敗')
  }
}

const prevMonth = async () => {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    selectedYear.value--
  } else {
    currentMonth.value--
  }
  await fetchHolidays()
}

const nextMonth = async () => {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    selectedYear.value++
  } else {
    currentMonth.value++
  }
  await fetchHolidays()
}

const goToToday = async () => {
  const now = new Date()
  currentMonth.value = now.getMonth()
  selectedYear.value = now.getFullYear()
  await fetchHolidays()
}

const formatFullDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'short'
  })
}

onMounted(() => {
  fetchHolidays()
})
</script>
