<template>
  <div class="fixed inset-0 z-[1000] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg sm:max-w-xl max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          新增排課規則
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
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
            @click="emit('close')"
            class="flex-1 glass-btn py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            {{ loading ? '儲存中...' : '新增' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  close: []
  created: []
}>()

const loading = ref(false)
const dataLoading = ref(true)
const error = ref<string | null>(null)

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

// 監聽課程選擇，自動帶入預設老師和教室
watch(() => form.value.offering_id, (newOfferingId) => {
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
    console.log('等待載入資源資料...')
    await fetchAllResources()
    console.log('載入資料完成:', {
      offerings: offerings.value.length,
      rooms: rooms.value.length,
      teachers: teachers.value.length
    })
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

const handleSubmit = async () => {
  if (form.value.weekdays.length === 0) {
    alert('請至少選擇一個星期')
    return
  }

  if (!form.value.offering_id) {
    alert('請選擇課程')
    return
  }

  loading.value = true

  try {
    const api = useApi()
    const centerId = getCenterId()

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

    console.log('提交排課規則資料:', JSON.stringify(data, null, 2))

    const response = await api.post('/admin/rules', data)

    console.log('排課規則建立成功:', response)

    emit('created')
    emit('close')
  } catch (error) {
    console.error('Failed to create schedule rule:', error)
    alert('新增失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
