<template>
  <div class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
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

      <form @submit.prevent="handleSubmit" class="p-4 space-y-4">
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

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">課程</label>
            <select v-model="form.offering_id" class="input-field text-sm sm:text-base" required>
              <option value="">請選擇課程</option>
              <option v-for="offering in offerings" :key="offering.id" :value="offering.id">
                {{ offering.name }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">老師</label>
            <select v-model="form.teacher_id" class="input-field text-sm sm:text-base" required>
              <option value="">請選擇老師</option>
              <option v-for="teacher in teachers" :key="teacher.id" :value="teacher.id">
                {{ teacher.name }}
              </option>
            </select>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">教室</label>
            <select v-model="form.room_id" class="input-field text-sm sm:text-base" required>
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
        </div>

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
  teacher_id: '',
  room_id: '',
  start_time: '09:00',
  end_time: '10:00',
  duration: 60,
  weekdays: [1] as number[],
  start_date: new Date().toISOString().split('T')[0],
  end_date: '',
})

const offerings = ref<any[]>([])
const teachers = ref<any[]>([])
const rooms = ref<any[]>([])

const { getCenterId } = useCenterId()

const fetchData = async () => {
  try {
    const api = useApi()
    const centerId = getCenterId()
    
    const [coursesRes, roomsRes, teachersRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>(`/admin/courses`),
      api.get<{ code: number; datas: any[] }>(`/admin/rooms`),
      api.get<{ code: number; datas: any[] }>('/teachers')
    ])
    
    offerings.value = coursesRes.datas || []
    rooms.value = roomsRes.datas || []
    teachers.value = teachersRes.datas || []
    
    console.log('載入資料:', {
      centerId,
      offerings: offerings.value.length,
      rooms: rooms.value.length,
      teachers: teachers.value.length
    })
  } catch (error) {
    console.error('Failed to fetch data:', error)
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

  loading.value = true

  try {
    const api = useApi()
    await api.post(`/admin/scheduling/rules`, {
      name: form.value.name,
      offering_id: parseInt(form.value.offering_id),
      teacher_id: parseInt(form.value.teacher_id),
      room_id: parseInt(form.value.room_id),
      start_time: form.value.start_time,
      end_time: form.value.end_time,
      duration: form.value.duration,
      weekdays: form.value.weekdays,
      start_date: form.value.start_date,
      end_date: form.value.end_date || null,
    })

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
