<template>
  <div class="p-4 md:p-6">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold text-white">課表模板</h1>
      <button
        @click="showModal = true"
        class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors"
      >
        新增模板
      </button>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="template in templates"
        :key="template.id"
        class="glass-card p-4"
      >
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 mb-3">
          <h3 class="text-lg font-medium text-white">{{ template.name }}</h3>
          <span
            class="px-2 py-1 rounded-full text-xs w-fit"
            :class="template.row_type === 'ROOM' ? 'bg-primary-500/20 text-primary-500' : 'bg-secondary-500/20 text-secondary-500'"
          >
            {{ template.row_type === 'ROOM' ? '教室視角' : '老師視角' }}
          </span>
        </div>

        <div class="flex items-center justify-between text-sm text-slate-400 mb-4">
          <span>建立於 {{ formatDate(template.created_at) }}</span>
          <span>{{ template.is_active !== false ? '啟用' : '停用' }}</span>
        </div>

        <div class="flex gap-2">
          <button
            @click="viewTemplate(template)"
            class="flex-1 px-3 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors text-sm"
          >
            查看格子
          </button>
          <button
            @click="openApplyModal(template)"
            class="px-3 py-2 rounded-lg bg-primary-500/20 text-primary-500 hover:bg-primary-500/30 transition-colors text-sm"
          >
            套用模板
          </button>
          <button
            @click="deleteTemplate(template.id)"
            class="px-3 py-2 rounded-lg bg-critical-500/20 text-critical-500 hover:bg-critical-500/30 transition-colors text-sm"
          >
            刪除
          </button>
        </div>
      </div>

      <div
        v-if="templates.length === 0"
        class="col-span-full text-center py-12 text-slate-500"
      >
        尚未建立課表模板
      </div>
    </div>
  </div>

  <div
    v-if="showModal"
    class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    @click.self="showModal = false"
  >
    <div class="glass-card w-full max-w-md p-6">
      <h3 class="text-lg font-semibold text-white mb-4">新增模板</h3>
      <form @submit.prevent="createTemplate">
        <div class="mb-4">
          <label for="template-name" class="block text-slate-300 mb-2">模板名稱</label>
          <input
            id="template-name"
            v-model="form.name"
            type="text"
            class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
            required
          />
        </div>
        <div class="mb-4">
          <label for="template-type" class="block text-slate-300 mb-2">視角類型</label>
          <select
            id="template-type"
            v-model="form.row_type"
            class="w-full px-3 py-2 rounded-lg bg-slate-800 border border-white/10 text-white cursor-pointer appearance-none"
          >
            <option value="ROOM">教室視角</option>
            <option value="TEACHER">老師視角</option>
          </select>
        </div>
        <div class="flex gap-3">
          <button
            type="button"
            @click="showModal = false"
            class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="creating"
            class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50"
          >
            {{ creating ? '建立中...' : '建立' }}
          </button>
        </div>
      </form>
    </div>
  </div>

  <div
    v-if="selectedTemplate"
    class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    @click.self="selectedTemplate = null"
  >
    <div class="glass-card w-full max-w-2xl p-6 max-h-[80vh] overflow-y-auto">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-white">{{ selectedTemplate.name }} - 格子</h3>
        <button @click="selectedTemplate = null" class="text-slate-400 hover:text-white">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div v-if="cells.length === 0" class="text-center py-8 text-slate-500">
        此模板尚未建立格子
      </div>

      <div v-else class="space-y-2">
        <div
          v-for="cell in cells"
          :key="cell.id"
          class="p-3 rounded-lg bg-white/5 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2"
        >
          <div class="text-white">
            {{ cell.start_time }} - {{ cell.end_time }}
          </div>
          <div class="text-sm text-slate-400">
            <span v-if="cell.room_id">教室 {{ cell.room_id }}</span>
            <span v-else-if="cell.teacher_id">老師 {{ cell.teacher_id }}</span>
            <span v-else>-</span>
          </div>
        </div>
      </div>

      <div class="mt-4 pt-4 border-t border-white/10">
        <button
          @click="addCells"
          class="w-full px-4 py-2 rounded-lg bg-primary-500/20 text-primary-500 hover:bg-primary-500/30 transition-colors"
        >
          新增格子
        </button>
      </div>
    </div>
  </div>

  <div
    v-if="showApplyModal"
    class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    @click.self="showApplyModal = false"
  >
    <div class="glass-card w-full max-w-lg p-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-white">套用模板 - {{ applyForm.templateName }}</h3>
        <button @click="showApplyModal = false" class="text-slate-400 hover:text-white">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="applyTemplate">
        <div class="mb-4">
          <label for="offering-select" class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">選擇課程</label>
          <select
            id="offering-select"
            v-model="applyForm.offeringId"
            class="input-field text-sm sm:text-base"
            required
          >
            <option value="">請選擇課程</option>
            <option v-for="offering in offerings" :key="offering.id" :value="offering.id">
              {{ offering.name }}
            </option>
          </select>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-4">
          <div>
            <label for="start-date" class="block text-slate-300 mb-2">開始日期</label>
            <input
              id="start-date"
              v-model="applyForm.startDate"
              type="date"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
              required
            />
          </div>
          <div>
            <label for="end-date" class="block text-slate-300 mb-2">結束日期</label>
            <input
              id="end-date"
              v-model="applyForm.endDate"
              type="date"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
              required
            />
          </div>
        </div>

        <div class="mb-4">
          <span class="block text-slate-300 mb-2">選擇星期</span>
          <div class="flex flex-wrap gap-2" role="group" aria-label="選擇星期">
            <label
              v-for="day in weekdays"
              :key="day.value"
              class="flex items-center gap-2 px-3 py-2 rounded-lg bg-white/5 cursor-pointer hover:bg-white/10 transition-colors"
            >
              <input
                type="checkbox"
                :value="day.value"
                v-model="applyForm.weekdays"
                class="w-4 h-4 rounded border-white/20 bg-white/10 text-primary-500"
              />
              <span class="text-white text-sm">{{ day.label }}</span>
            </label>
          </div>
        </div>

        <div class="mb-4">
          <label for="duration-input" class="block text-slate-300 mb-2">每堂課時長（分鐘）</label>
          <input
            id="duration-input"
            v-model.number="applyForm.duration"
            type="number"
            min="30"
            step="15"
            class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
            required
          />
        </div>

        <div class="flex gap-3">
          <button
            type="button"
            @click="showApplyModal = false"
            class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="applying"
            class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50"
          >
            {{ applying ? '套用中...' : '確認套用' }}
          </button>
        </div>
      </form>
    </div>
  </div>

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
 definePageMeta({
   middleware: 'auth-admin',
   layout: 'admin',
 })

 const notificationUI = useNotification()
const showModal = ref(false)
const showApplyModal = ref(false)
const selectedTemplate = ref<any>(null)
const templates = ref<any[]>([])
const cells = ref<any[]>([])
const offerings = ref<any[]>([])
const creating = ref(false)
const applying = ref(false)
const { getCenterId } = useCenterId()
const { confirm: alertConfirm, error: alertError, warning: alertWarning } = useAlert()

const applyForm = ref({
  templateId: 0,
  templateName: '',
  offeringId: '',
  startDate: '',
  endDate: '',
  weekdays: [] as number[],
  duration: 60
})

const weekdays = [
  { value: 1, label: '週一' },
  { value: 2, label: '週二' },
  { value: 3, label: '週三' },
  { value: 4, label: '週四' },
  { value: 5, label: '週五' },
  { value: 6, label: '週六' },
  { value: 7, label: '週日' }
]

const form = ref({
  name: '',
  row_type: 'ROOM'
})

const fetchTemplates = async () => {
  try {
    const api = useApi()
    const centerId = getCenterId()
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/centers/${centerId}/templates`)
    templates.value = response.datas || []
  } catch (error) {
    console.error('Failed to fetch templates:', error)
  }
}

const createTemplate = async () => {
  creating.value = true
  try {
    const api = useApi()
    const centerId = getCenterId()
    await api.post(`/admin/centers/${centerId}/templates`, form.value)
    showModal.value = false
    form.value = { name: '', row_type: 'ROOM' }
    await fetchTemplates()
  } catch (error) {
    console.error('Failed to create template:', error)
    await alertError('建立失敗')
  } finally {
    creating.value = false
  }
}

const deleteTemplate = async (id: number) => {
  if (!await alertConfirm('確定要刪除此模板？')) return

  try {
    const api = useApi()
    const centerId = getCenterId()
    await api.delete(`/admin/centers/${centerId}/templates/${id}`)
    await fetchTemplates()
  } catch (error) {
    console.error('Failed to delete template:', error)
    await alertError('刪除失敗')
  }
}

const viewTemplate = async (template: any) => {
  selectedTemplate.value = template
  try {
    const api = useApi()
    const centerId = getCenterId()
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/centers/${centerId}/templates/${template.id}/cells`)
    cells.value = response.datas || []
  } catch (error) {
    console.error('Failed to fetch cells:', error)
    cells.value = []
  }
}

const addCells = async () => {
  if (!selectedTemplate.value) return

  try {
    const api = useApi()
    const centerId = getCenterId()
    await api.post(`/admin/centers/${centerId}/templates/${selectedTemplate.value.id}/cells`, [
      { row_no: 1, col_no: 1, start_time: '09:00', end_time: '10:00' }
    ])
    await viewTemplate(selectedTemplate.value)
  } catch (error) {
    console.error('Failed to add cells:', error)
  }
}

const fetchOfferings = async () => {
  try {
    const api = useApi()
    const response = await api.get<any>('/admin/offerings')
    console.log('Offerings API response:', response)
    
    // 處理不同格式的 API 回應
    if (response.datas?.offerings) {
      offerings.value = response.datas.offerings
    } else if (response.datas) {
      offerings.value = response.datas
    } else if (Array.isArray(response)) {
      offerings.value = response
    } else {
      offerings.value = []
    }
    
    console.log('Loaded offerings:', offerings.value)
  } catch (error) {
    console.error('Failed to fetch offerings:', error)
    offerings.value = []
  }
}

const openApplyModal = (template: any) => {
  selectedTemplate.value = null
  applyForm.value = {
    templateId: template.id,
    templateName: template.name,
    offeringId: '',
    startDate: '',
    endDate: '',
    weekdays: [],
    duration: 60
  }
  showApplyModal.value = true
}

const applyTemplate = async () => {
  if (!applyForm.value.offeringId || applyForm.value.weekdays.length === 0) {
    await alertWarning('請填寫完整資訊')
    return
  }

  applying.value = true
  try {
    const api = useApi()
    const centerId = getCenterId()
    await api.post(`/admin/centers/${centerId}/templates/${applyForm.value.templateId}/apply`, {
      offering_id: Number(applyForm.value.offeringId),
      start_date: applyForm.value.startDate,
      end_date: applyForm.value.endDate,
      weekdays: applyForm.value.weekdays,
      duration: applyForm.value.duration
    })
    showApplyModal.value = false
    notificationUI.showSuccess('模板套用成功')
  } catch (error) {
    console.error('Failed to apply template:', error)
    await alertError('套用失敗，請稍後再試')
  } finally {
    applying.value = false
  }
}

const formatDate = (dateStr: string): string => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-TW')
}

onMounted(async () => {
  await Promise.all([fetchTemplates(), fetchOfferings()])
})
</script>
