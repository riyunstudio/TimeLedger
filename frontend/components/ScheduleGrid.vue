<template>
  <div class="h-full flex flex-col glass-card overflow-hidden">
    <div class="p-4 border-b border-white/10 shrink-0">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex items-center gap-4">
          <button
            @click="changeWeek(-1)"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>

          <h2 class="text-lg font-semibold text-slate-100">
            {{ weekLabel }}
          </h2>

          <button
            @click="changeWeek(1)"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>

        <div class="flex items-center gap-4">
          <!-- 視角切換器 -->
          <div class="flex items-center gap-1 bg-slate-800/80 rounded-lg p-1">
            <button
              @click="viewModeModel = 'calendar'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="viewMode === 'calendar' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              週曆
            </button>
            <button
              @click="viewModeModel = 'teacher_matrix'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="viewMode === 'teacher_matrix' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              老師矩陣
            </button>
            <button
              @click="viewModeModel = 'room_matrix'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="viewMode === 'room_matrix' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              教室矩陣
            </button>
          </div>

          <!-- 矩陣視角選擇器 -->
          <select
            v-if="viewMode !== 'calendar'"
            v-model="selectedResourceIdModel"
            class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500"
          >
            <option :value="null">選擇{{ viewMode === 'teacher_matrix' ? '老師' : '教室' }}...</option>
            <option v-for="resource in resourceList" :key="resource.id" :value="resource.id">
              {{ resource.name }}
            </option>
          </select>

          <button
            @click="showCreateModal = true"
            class="btn-primary px-4 py-2 text-sm font-medium"
          >
            + 新增排課規則
          </button>
        </div>
      </div>

      <!-- 選中資源提示 -->
      <div
        v-if="viewMode !== 'calendar' && selectedResourceName"
        class="mt-3 flex items-center gap-2 px-3 py-2 bg-primary-500/10 border border-primary-500/30 rounded-lg"
      >
        <span class="text-sm text-primary-400">
          {{ viewMode === 'teacher_matrix' ? '老師' : '教室' }}矩陣：
        </span>
        <span class="text-sm font-medium text-white">{{ selectedResourceName }}</span>
        <button
          @click="clearViewMode"
          class="ml-auto p-1 hover:bg-white/10 rounded transition-colors"
        >
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <div
      class="flex-1 overflow-auto p-4"
      @dragover.prevent="handleDragOver"
      @drop="handleDrop"
    >
      <div class="min-w-[600px]">
        <div class="grid grid-cols-[80px_repeat(7)] sticky top-0 z-10 bg-slate-800/90 backdrop-blur-sm">
          <div class="p-2 border-b border-white/10 text-center">
            <span class="text-xs text-slate-400">時段</span>
          </div>
          <div
            v-for="day in weekDays"
            :key="day.value"
            class="p-2 border-b border-white/10 text-center"
          >
            <span class="text-sm font-medium text-slate-100">{{ day.name }}</span>
          </div>
        </div>

        <div
          v-for="time in timeSlots"
          :key="time"
          class="grid grid-cols-[80px_repeat(7)]"
        >
          <div class="p-2 border-r border-b border-white/5 text-right text-xs text-slate-400">
            {{ formatTime(time) }}
          </div>

          <div
            v-for="day in weekDays"
            :key="`${time}-${day.value}`"
            class="p-1 min-h-[60px] border-b border-white/5 border-r relative"
            :class="getCellClass(time, day.value)"
            @dragenter="handleDragEnter(time, day.value)"
            @dragleave="handleDragLeave"
            @dragover.prevent
          >
            <div
              v-if="getScheduleAt(time, day.value)"
              class="rounded-lg p-1 text-xs cursor-pointer hover:opacity-80 transition-opacity group relative"
              :class="getScheduleCardClass(getScheduleAt(time, day.value))"
              @click="selectSchedule(time, day.value)"
            >
              <!-- 簡短資訊 -->
              <div class="font-medium truncate">
                {{ getScheduleAt(time, day.value)?.offering_name }}
              </div>
              <div class="text-slate-400 truncate">
                {{ getScheduleAt(time, day.value)?.teacher_name }}
              </div>

              <!-- 懸停 Tooltip -->
              <div class="absolute z-50 left-0 bottom-full mb-2 hidden group-hover:block w-64">
                <div class="glass p-3 rounded-lg shadow-xl border border-white/10">
                  <div class="font-medium text-white mb-2">{{ getScheduleAt(time, day.value)?.offering_name }}</div>
                  <div class="space-y-1 text-xs">
                    <div class="flex justify-between text-slate-400">
                      <span>老師：</span>
                      <span class="text-slate-200">{{ getScheduleAt(time, day.value)?.teacher_name }}</span>
                    </div>
                    <div class="flex justify-between text-slate-400">
                      <span>教室：</span>
                      <span class="text-slate-200">{{ getScheduleAt(time, day.value)?.room_name }}</span>
                    </div>
                    <div class="flex justify-between text-slate-400">
                      <span>時間：</span>
                      <span class="text-slate-200">{{ formatTime(time) }} - {{ formatTime(time + 1) }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ScheduleDetailPanel
      v-if="selectedCell"
      :time="selectedCell.time"
      :weekday="selectedCell.day"
      :schedule="selectedSchedule"
      @close="selectedCell = null"
    />

    <ScheduleRuleModal
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @created="handleRuleCreated"
    />
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  selectCell: { time: number, weekday: number }
  'update:viewMode': [value: 'calendar' | 'teacher_matrix' | 'room_matrix']
  'update:selectedResourceId': [value: number | null]
}>()

// Props
const props = defineProps<{
  viewMode: 'calendar' | 'teacher_matrix' | 'room_matrix'
  selectedResourceId: number | null
}>()

// Computed with setter for v-model support
const viewModeModel = computed({
  get: () => props.viewMode,
  set: (value) => emit('update:viewMode', value)
})

const selectedResourceIdModel = computed({
  get: () => props.selectedResourceId,
  set: (value) => emit('update:selectedResourceId', value)
})

const showCreateModal = ref(false)
const selectedCell = ref<{ time: number, day: number } | null>(null)
const selectedSchedule = ref<any>(null)
const dragTarget = ref<{ time: number, day: number } | null>(null)
const validationResults = ref<Record<string, any>>({})

// 資源快取（用於顯示選中的資源名稱）
const resourceCache = ref<{ teachers: Map<number, any>, rooms: Map<number, any> }>({
  teachers: new Map(),
  rooms: new Map(),
})

// 資源列表（根據視角模式動態取得）
const resourceList = computed(() => {
  if (props.viewMode === 'teacher_matrix') {
    return Array.from(resourceCache.value.teachers.values())
  } else if (props.viewMode === 'room_matrix') {
    return Array.from(resourceCache.value.rooms.values())
  }
  return []
})

const getWeekStart = (date: Date): Date => {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  return new Date(d.setDate(diff))
}

const weekStart = ref(getWeekStart(new Date()))
const weekEnd = computed(() => {
  const end = new Date(weekStart.value)
  end.setDate(end.getDate() + 6)
  return end
})

const weekLabel = computed(() => {
  const start = weekStart.value.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric' })
  const end = weekEnd.value.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric', year: 'numeric' })
  return `${start} - ${end}`
})

const timeSlots = [9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]

const weekDays = [
  { value: 1, name: '週一' },
  { value: 2, name: '週二' },
  { value: 3, name: '週三' },
  { value: 4, name: '週四' },
  { value: 5, name: '週五' },
  { value: 6, name: '週六' },
  { value: 7, name: '週日' },
]

const schedules = ref<Record<string, any>>({})
const { getCenterId } = useCenterId()

const selectedResourceName = computed(() => {
  if (props.viewMode === 'teacher_matrix') {
    return resourceCache.value.teachers.get(props.selectedResourceId)?.name || '未知老師'
  } else if (props.viewMode === 'room_matrix') {
    return resourceCache.value.rooms.get(props.selectedResourceId)?.name || '未知教室'
  }
  return ''
})

const clearViewMode = () => {
  viewModeModel.value = 'calendar'
  selectedResourceIdModel.value = null
}

const fetchSchedules = async () => {
  try {
    const api = useApi()
    const centerId = getCenterId()
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/centers/${centerId}/rules`)
    const rules = response.datas || []

    // 將規則轉換為 schedule map
    const scheduleMap: Record<string, any> = {}
    rules.forEach((rule: any) => {
      rule.weekdays?.forEach((day: number) => {
        const key = `${rule.start_time.split(':')[0]}-${day}`
        scheduleMap[key] = {
          id: rule.id,
          offering_name: rule.offering?.name || '-',
          teacher_name: rule.teacher?.name || '-',
          teacher_id: rule.teacher_id,
          room_id: rule.room_id,
          room_name: rule.room?.name || '-',
          ...rule,
        }
      })
    })
    schedules.value = scheduleMap
  } catch (error) {
    console.error('Failed to fetch schedules:', error)
    schedules.value = {}
  }
}

// 根據視角模式過濾排課
const filteredSchedules = computed(() => {
  // 週曆視圖顯示全部
  if (props.viewMode === 'calendar' || !props.selectedResourceId) {
    return schedules.value
  }

  // 矩陣視圖過濾特定資源
  const filtered: Record<string, any> = {}
  Object.entries(schedules.value).forEach(([key, schedule]) => {
    if (props.viewMode === 'teacher_matrix') {
      if (schedule.teacher_id === props.selectedResourceId) {
        filtered[key] = schedule
      }
    } else if (props.viewMode === 'room_matrix') {
      if (schedule.room_id === props.selectedResourceId) {
        filtered[key] = schedule
      }
    }
  })
  return filtered
})

const changeWeek = (delta: number) => {
  weekStart.value = getWeekStart(new Date(weekStart.value.getTime() + delta * 7 * 24 * 60 * 60 * 1000))
}

const formatTime = (hour: number): string => {
  return `${hour}:00`
}

const getScheduleAt = (time: number, weekday: number) => {
  return filteredSchedules.value[`${time}-${weekday}`]
}

const getCellClass = (time: number, weekday: number): string => {
  const key = `${time}-${weekday}`
  const validation = validationResults.value[key]

  if (validation?.valid === false) {
    return 'bg-critical-500/10 border-critical-500/50'
  } else if (validation?.warning) {
    return 'bg-warning-500/10 border-warning-500/50'
  } else if (validation?.valid === true) {
    return 'bg-success-500/10 border-success-500/50'
  } else if (getScheduleAt(time, weekday)) {
    return 'bg-primary-500/10'
  }

  return 'hover:bg-white/5'
}

const getScheduleCardClass = (schedule: any): string => {
  if (!schedule) return ''

  return 'bg-slate-700/80 border border-white/10'
}

const selectSchedule = (time: number, weekday: number) => {
  selectedCell.value = { time, day: weekday }
  selectedSchedule.value = getScheduleAt(time, weekday)
  emit('selectCell', { time, weekday: weekday })
}

const handleDragOver = (event: DragEvent) => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.time}-${dragTarget.value.day}`
    validationResults.value[key] = { valid: true }
  }
}

const handleDragEnter = (time: number, day: number) => {
  dragTarget.value = { time, day }
}

const handleDragLeave = () => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.time}-${dragTarget.value.day}`
    delete validationResults.value[key]
  }
  dragTarget.value = null
}

const handleDrop = async (event: DragEvent) => {
  event.preventDefault()

  if (!dragTarget.value) return

  const data = event.dataTransfer?.getData('application/json')
  if (!data) return

  const parsed = JSON.parse(data)
  const { type, item } = parsed

  const key = `${dragTarget.value.time}-${dragTarget.value.day}`

  validationResults.value[key] = { valid: 'checking' }

  try {
    const api = useApi()
    const centerId = getCenterId()
    const teacherId = type === 'teacher' ? item.id : (item.teacher_id || null)
    const roomId = type === 'room' ? item.id : (item.room_id || null)
    
    const response = await api.post<any>('/admin/scheduling/check-overlap', {
      center_id: parseInt(centerId),
      teacher_id: teacherId,
      room_id: roomId,
      start_time: `${formatDate(weekStart.value)}T${formatTime(dragTarget.value.time)}:00`,
      end_time: `${formatDate(weekStart.value)}T${formatTime(dragTarget.value.time + 1)}:00`,
    })

    if (response.data.valid) {
      validationResults.value[key] = { valid: true }
      schedules.value[key] = {
        id: Date.now(),
        offering_name: item.name || item.title,
        teacher_name: item.name || item.title,
        room_id: item.id,
      }
    } else {
      validationResults.value[key] = { valid: false, conflicts: response.data.conflicts }
    }
  } catch (error) {
    console.error('Validation failed:', error)
    validationResults.value[key] = { valid: false, error: true }
  }

  dragTarget.value = null
}

const formatDate = (date: Date): string => {
  return date.toISOString().split('T')[0]
}

const handleRuleCreated = () => {
  fetchSchedules()
}

const fetchResourceCache = async () => {
  try {
    const api = useApi()
    const [teachersRes, roomsRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>('/teachers'),
      api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    ])

    teachersRes.datas?.forEach((t: any) => {
      resourceCache.value.teachers.set(t.id, t)
    })
    roomsRes.datas?.forEach((r: any) => {
      resourceCache.value.rooms.set(r.id, r)
    })
  } catch (error) {
    console.error('Failed to fetch resource cache:', error)
  }
}

onMounted(() => {
  fetchSchedules()
  fetchResourceCache()
})
</script>
