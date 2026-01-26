<template>
  <Teleport to="body">
    <div
      v-if="showCreateModal || showEditModal"
      class="fixed inset-0 z-[1000] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm isolate"
      @click.self="handleClose"
    >
      <div class="glass-card w-full max-w-lg sm:max-w-xl max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
        <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
          <h3 class="text-lg font-semibold text-slate-100">
            {{ showEditModal ? '編輯排課規則' : '新增排課規則' }}
          </h3>
          <button @click="handleClose" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

      <!-- 載入中 -->
      <div v-if="dataLoading" class="p-8 text-center">
        <div class="inline-flex items-center gap-2 text-slate-400">
          <svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span>載入中...</span>
        </div>
      </div>

      <!-- 衝突提示 -->
      <div v-if="conflictError" class="p-4">
        <div class="bg-critical-500/10 border border-critical-500/30 rounded-xl p-4 mb-4">
          <div class="flex items-start gap-3">
            <svg class="w-6 h-6 text-critical-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <div>
              <h4 class="font-medium text-critical-500 mb-2">排課時間衝突</h4>
              <p class="text-sm text-slate-400 mb-3">以下時間已有排課，請選擇其他時間或教室：</p>
              <ul class="space-y-1">
                <li v-for="(conflict, index) in conflictErrors" :key="index" class="text-sm text-slate-300 flex items-center gap-2">
                  <span class="w-1.5 h-1.5 rounded-full bg-critical-500"></span>
                  {{ conflict }}
                </li>
              </ul>
            </div>
          </div>
        </div>
        <button @click="conflictError = null; conflictErrors = []" class="btn-secondary w-full py-3 rounded-xl font-medium">
          我知道了
        </button>
      </div>

      <!-- 自定義 Alert 彈窗 -->
      <div v-else-if="showAlert" class="p-4">
        <div class="glass-card p-6">
          <div class="flex items-start gap-4">
            <!-- 警告圖示 -->
            <div v-if="alertType === 'warning'" class="w-10 h-10 rounded-full bg-warning-500/20 flex items-center justify-center flex-shrink-0">
              <svg class="w-6 h-6 text-warning-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <!-- 錯誤圖示 -->
            <div v-else-if="alertType === 'error'" class="w-10 h-10 rounded-full bg-critical-500/20 flex items-center justify-center flex-shrink-0">
              <svg class="w-6 h-6 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
            <!-- 資訊圖示 -->
            <div v-else class="w-10 h-10 rounded-full bg-primary-500/20 flex items-center justify-center flex-shrink-0">
              <svg class="w-6 h-6 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            
            <div class="flex-1">
              <h4 v-if="alertTitle" class="font-medium text-white mb-2">{{ alertTitle }}</h4>
              <p class="text-sm text-slate-400">{{ alertMessage }}</p>
            </div>
          </div>
          
          <div class="flex gap-3 mt-6">
            <button 
              @click="showAlert = false" 
              class="flex-1 py-2.5 rounded-xl font-medium transition-all"
              :class="alertType === 'error' ? 'btn-critical' : 'btn-primary'"
            >
              確定
            </button>
          </div>
        </div>
      </div>

      <!-- 錯誤訊息 -->
      <div v-else-if="error" class="p-8 text-center">
        <div class="text-critical-500 mb-2">
          <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <p class="text-slate-300 mb-4">{{ error }}</p>
        <button @click="fetchData" class="btn-primary px-4 py-2 text-sm">
          重試
        </button>
      </div>

      <!-- 表單內容 -->
      <form v-show="!dataLoading && !error" @submit.prevent="handleSubmit" class="p-4 space-y-4">
        <!-- 空資料提示 -->
        <div v-if="offerings.length === 0 || rooms.length === 0 || teachers.length === 0" class="mb-4 p-4 rounded-lg bg-warning-500/10 border border-warning-500/30">
          <p class="text-warning-500 text-sm">
            <span v-if="offerings.length === 0">尚未建立課程班別，請先至「資源管理」建立</span>
            <span v-if="rooms.length === 0">尚未建立教室</span>
            <span v-if="teachers.length === 0">尚未建立老師</span>
          </p>
        </div>

        <!-- 規則名稱 -->
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">規則名稱</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="例：週一上午鋼琴課"
            class="input-field text-sm sm:text-base"
            required
          />
        </div>

        <!-- 課程和老師 -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">課程</label>
            <select v-model="form.offering_id" class="input-field text-sm sm:text-base" required>
              <option value="">請選擇課程</option>
              <option v-for="offering in offerings" :key="offering.id" :value="offering.id">
                {{ offering.name || `班別 #${offering.id}` }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">老師</label>
            <select v-model="form.teacher_id" class="input-field text-sm sm:text-base">
              <option value="">請選擇老師</option>
              <option v-for="teacher in teachers" :key="teacher.id" :value="teacher.id">
                {{ teacher.name }}
              </option>
            </select>
          </div>
        </div>

        <!-- 教室和時間 -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">教室</label>
            <select v-model="form.room_id" class="input-field text-sm sm:text-base">
              <option value="">請選擇教室</option>
              <option v-for="room in rooms" :key="room.id" :value="room.id">
                {{ room.name }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">開始時間</label>
            <input
              v-model="form.start_time"
              type="time"
              class="input-field text-sm sm:text-base"
              required
            />
          </div>
        </div>

        <!-- 結束時間和時長 -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">結束時間</label>
            <input
              v-model="form.end_time"
              type="time"
              class="input-field text-sm sm:text-base"
              required
            />
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">課程時長（分鐘）</label>
            <input
              v-model.number="form.duration"
              type="number"
              min="30"
              step="30"
              class="input-field text-sm sm:text-base"
              required
            />
          </div>
        </div>

        <!-- 重複星期 -->
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">重複星期</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="day in weekDays"
              :key="day.value"
              type="button"
              @click="toggleDay(day.value)"
              class="px-3 py-2 rounded-lg text-sm font-medium transition-all"
              :class="form.weekdays.includes(day.value) ? 'bg-primary-500 text-white' : 'bg-slate-700/50 text-slate-400 hover:bg-slate-700'"
            >
              {{ day.name }}
            </button>
          </div>
          <HelpTooltip
            class="mt-2"
            title="重複星期"
            description="選擇此排課規則適用的星期幾。"
            :usage="['可選擇多個星期', '形成每週重複的排課']"
          />
        </div>

        <!-- 開始和結束日期 -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">開始日期</label>
            <input
              v-model="form.start_date"
              type="date"
              class="input-field text-sm sm:text-base"
              required
            />
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">結束日期</label>
            <input
              v-model="form.end_date"
              type="date"
              class="input-field text-sm sm:text-base"
            />
          </div>
        </div>

        <div class="flex gap-3 pt-2">
          <button
            type="button"
            @click="handleClose"
            class="flex-1 glass-btn py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            {{ loading ? '儲存中...' : (showEditModal ? '儲存修改' : '新增') }}
          </button>
        </div>
      </form>
    </div>
  </div>
  </Teleport>
</template>

<script setup lang="ts">
// Props for create mode
const props = defineProps<{
  editingRule?: any | null
  updateMode?: string
}>()

const emit = defineEmits<{
  close: []
  submit: [formData: any, updateMode: string]
}>()

const loading = ref(false)
const dataLoading = ref(true)
const error = ref<string | null>(null)
const conflictError = ref<string | null>(null)
const conflictErrors = ref<string[]>([])
const showAlert = ref(false)
const alertTitle = ref('')
const alertMessage = ref('')
const alertType = ref<'info' | 'warning' | 'error'>('info')
const showCreateModal = computed(() => !props.editingRule)
const showEditModal = computed(() => !!props.editingRule)

const weekDays = [
  { value: 1, name: '週一' },
  { value: 2, name: '週二' },
  { value: 3, name: '週三' },
  { value: 4, name: '週四' },
  { value: 5, name: '週五' },
  { value: 6, name: '週六' },
  { value: 7, name: '週日' },
]

const form = ref({
  name: '',
  offering_id: '',
  teacher_id: null as number | null,
  room_id: null as number | null,
  start_time: '09:00',
  end_time: '10:00',
  duration: 60,
  weekdays: [1] as number[],
  start_date: new Date().toISOString().split('T')[0],
  end_date: '',
})

// 使用共享的資源緩存
const { resourceCache, fetchAllResources } = useResourceCache()

// 從共享緩存取得資料
const offerings = computed(() => resourceCache.value.offerings)
const teachers = computed(() => Array.from(resourceCache.value.teachers.values()))
const rooms = computed(() => Array.from(resourceCache.value.rooms.values()))

const { getCenterId } = useCenterId()

// 載入編輯資料
watch(() => props.editingRule, (rule) => {
  if (rule) {
    form.value = {
      name: rule.name || '',
      offering_id: String(rule.offering_id || ''),
      teacher_id: rule.teacher_id || null,
      room_id: rule.room_id || null,
      start_time: rule.start_time || '09:00',
      end_time: rule.end_time || '10:00',
      duration: rule.duration || 60,
      weekdays: [rule.weekday] || [1],
      start_date: rule.effective_range?.start_date?.split('T')[0] || new Date().toISOString().split('T')[0],
      end_date: rule.effective_range?.end_date?.split('T')[0] || '',
    }
  }
}, { immediate: true })

// 監聽課程選擇，自動帶入預設老師和教室
watch(() => form.value.offering_id, (newOfferingId) => {
  // 編輯模式不自動帶入預設值
  if (showEditModal.value) return
  if (!newOfferingId) return

  const selectedOffering = offerings.value.find(o => o.id === parseInt(newOfferingId))
  if (selectedOffering) {
    // 自動帶入預設老師（如果還沒有選老師）
    if (selectedOffering.default_teacher_id && !form.value.teacher_id) {
      form.value.teacher_id = selectedOffering.default_teacher_id
    }
    // 自動帶入預設教室（如果還沒有選教室）
    if (selectedOffering.default_room_id && !form.value.room_id) {
      form.value.room_id = selectedOffering.default_room_id
    }
    // 自動帶入名稱（如果還沒有填名稱）
    if (!form.value.name) {
      form.value.name = selectedOffering.name
    }
  }
})

const fetchData = async () => {
  dataLoading.value = true
  error.value = null

  try {
    await fetchAllResources()
  } catch (err: any) {
    console.error('Failed to fetch data:', err)
    error.value = err.message || '載入資料失敗'
  } finally {
    dataLoading.value = false
  }
}

const toggleDay = (day: number) => {
  const index = form.value.weekdays.indexOf(day)
  if (index === -1) {
    form.value.weekdays.push(day)
  } else {
    form.value.weekdays.splice(index, 1)
  }
}

const handleClose = () => {
  emit('close')
}

// 自定義 Alert 函數
const customAlert = (message: string, title?: string, type: 'info' | 'warning' | 'error' = 'info') => {
  alertTitle.value = title || (type === 'error' ? '操作失敗' : type === 'warning' ? '提醒' : '提示')
  alertMessage.value = message
  alertType.value = type
  showAlert.value = true
}

const handleSubmit = async () => {
  if (form.value.weekdays.length === 0) {
    customAlert('請至少選擇一個星期', '驗證提醒', 'warning')
    return
  }

  if (!form.value.offering_id) {
    customAlert('請選擇課程', '驗證提醒', 'warning')
    return
  }

  loading.value = true

  try {
    const data: any = {
      name: form.value.name,
      offering_id: parseInt(form.value.offering_id),
      start_time: form.value.start_time,
      end_time: form.value.end_time,
      duration: form.value.duration,
      weekdays: form.value.weekdays,
      start_date: form.value.start_date,
      end_date: form.value.end_date || null,
    }

    // 只有當有選擇老師時才傳送
    if (form.value.teacher_id) {
      data.teacher_id = form.value.teacher_id
    }

    // 只有當有選擇教室時才傳送
    if (form.value.room_id) {
      data.room_id = form.value.room_id
    }

    if (showEditModal.value) {
      // 編輯模式：發射表單資料給父元件處理
      emit('submit', data, props.updateMode || 'ALL')
      handleClose()
    } else {
      // 新增模式：直接呼叫 API
      const api = useApi()
      await api.post('/admin/rules', data)
      handleClose()
      // 通知父元件刷新列表
      refreshParent()
    }
  } catch (error: any) {
    console.error('Failed to save schedule rule:', error)

    // 處理衝突錯誤
    if (error.response?.data?.code === 40002 || error.response?.data?.code === 'OVERLAP' || error.response?.data?.code === 20002) {
      conflictError.value = error.response?.data?.message || '排課時間與現有規則衝突'
      conflictErrors.value = error.response?.data?.datas?.conflicts || []
    } else {
      customAlert(showEditModal.value ? '更新失敗，請稍後再試' : '新增失敗，請稍後再試', '操作失敗', 'error')
    }
  } finally {
    loading.value = false
  }
}

// 刷新父元件列表
const refreshParent = async () => {
  try {
    const api = useApi()
    // 觸發重新整理事件，這裡透過重新導向來刷新
    await api.get('/admin/rules')
  } catch (e) {
    console.error('Failed to refresh:', e)
  }
}

onMounted(() => {
  fetchData()
})
</script>
