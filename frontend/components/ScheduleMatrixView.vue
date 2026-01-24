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
          <!-- 切換回週曆 -->
          <button
            @click="$emit('switchToCalendar')"
            class="px-3 py-1.5 rounded-lg text-sm bg-slate-700/50 text-slate-300 hover:text-white hover:bg-slate-600 transition-colors"
          >
            返回週曆
          </button>

          <!-- 資源類型切換 -->
          <div class="flex items-center gap-1 bg-slate-800/80 rounded-lg p-1">
            <button
              @click="resourceType = 'teacher'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="resourceType === 'teacher' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              老師
            </button>
            <button
              @click="resourceType = 'room'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="resourceType === 'room' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              教室
            </button>
          </div>

          <button
            @click="showCreateModal = true"
            class="btn-primary px-4 py-2 text-sm font-medium"
          >
            + 新增排課規則
          </button>
        </div>
      </div>
    </div>

    <div
      class="flex-1 overflow-auto"
      @dragover.prevent="handleDragOver"
      @drop="handleDrop"
    >
      <div class="min-w-[800px]">
        <!-- 表頭：時段 -->
        <div class="grid" :style="{ gridTemplateColumns: `120px repeat(${timeSlots.length}, 1fr)` }">
          <div class="p-3 border-b border-white/10 text-center bg-slate-800/50">
            <span class="text-sm font-medium text-slate-300">資源 / 時段</span>
          </div>
          <div
            v-for="time in timeSlots"
            :key="time"
            class="p-3 border-b border-white/10 text-center bg-slate-800/50"
          >
            <span class="text-sm font-medium text-slate-300">{{ formatTime(time) }}</span>
          </div>
        </div>

        <!-- 資源列 -->
        <div
          v-for="resource in resourceList"
          :key="resource.id"
          class="grid hover:bg-white/5 transition-colors"
          :style="{ gridTemplateColumns: `120px repeat(${timeSlots.length}, 1fr)` }"
        >
          <!-- 資源名稱 -->
          <div class="p-3 border-r border-b border-white/10 flex items-center">
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0 mr-2">
              <span class="text-white text-sm font-medium">{{ resource.name?.charAt(0) || '?' }}</span>
            </div>
            <span class="text-sm text-slate-300 truncate">{{ resource.name }}</span>
          </div>

          <!-- 時段儲存格 -->
          <div
            v-for="time in timeSlots"
            :key="`${resource.id}-${time}`"
            class="p-1 min-h-[60px] border-b border-white/5 border-r relative"
            :class="getCellClass(resource.id, time)"
            @dragenter="handleDragEnter(resource.id, time)"
            @dragleave="handleDragLeave"
            @dragover.prevent
          >
            <div
              v-if="getScheduleAt(resource.id, time)"
              class="h-full rounded-lg p-2 text-xs cursor-pointer hover:opacity-80 transition-opacity group relative"
              :class="getScheduleCardClass(getScheduleAt(resource.id, time))"
              @click="selectSchedule(resource, time)"
            >
              <!-- 簡短資訊 -->
              <div class="font-medium truncate text-slate-100">
                {{ getScheduleAt(resource.id, time)?.offering_name }}
              </div>
              <div class="text-slate-400 truncate">
                {{ getScheduleAt(resource.id, time)?.teacher_name || getScheduleAt(resource.id, time)?.room_name }}
              </div>

              <!-- 懸停 Tooltip -->
              <div class="absolute z-50 left-0 bottom-full mb-2 hidden group-hover:block w-64">
                <div class="glass p-3 rounded-lg shadow-xl border border-white/10">
                  <div class="font-medium text-white mb-2">{{ getScheduleAt(resource.id, time)?.offering_name }}</div>
                  <div class="space-y-1 text-xs">
                    <div class="flex justify-between text-slate-400">
                      <span>老師：</span>
                      <span class="text-slate-200">{{ getScheduleAt(resource.id, time)?.teacher_name }}</span>
                    </div>
                    <div class="flex justify-between text-slate-400">
                      <span>教室：</span>
                      <span class="text-slate-200">{{ getScheduleAt(resource.id, time)?.room_name }}</span>
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

        <!-- 無資源時顯示 -->
        <div v-if="resourceList.length === 0" class="text-center py-12 text-slate-500">
          沒有{{ resourceType === 'teacher' ? '老師' : '教室' }}資料
        </div>
      </div>
    </div>

    <ScheduleDetailPanel
      v-if="selectedCell"
      :time="selectedCell.time"
      :weekday="selectedCell.weekday"
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
  selectCell: { resource: any, time: number, weekday: number }
  switchToCalendar: []
}>()

const showCreateModal = ref(false)
const selectedCell = ref<{ resource: any, time: number, weekday: number } | null>(null)
const selectedSchedule = ref<any>(null)
const dragTarget = ref<{ resourceId: number, time: number } | null>(null)
const validationResults = ref<Record<string, any>>({})

const resourceType = ref<'teacher' | 'room'>('teacher')
const { getCenterId } = useCenterId()

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

const schedules = ref<{ teachers: Record<string, any>, rooms: Record<string, any> }>({
  teachers: {},
  rooms: {}
})
const teachers = ref<any[]>([])
const rooms = ref<any[]>([])

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

const resourceList = computed(() => {
  return resourceType.value === 'teacher' ? teachers.value : rooms.value
})

const formatTime = (hour: number): string => {
  return `${hour}:00`
}

const fetchData = async () => {
  try {
    const api = useApi()
    const centerId = getCenterId()

    const [rulesRes, teachersRes, roomsRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>(`/admin/centers/${centerId}/rules`),
      api.get<{ code: number; datas: any[] }>('/teachers'),
      api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    ])

    // 處理老師資料
    teachers.value = teachersRes.datas || []

    // 處理教室資料
    rooms.value = roomsRes.datas || []

    // 將規則轉換為 schedule map（分別存儲老師和教室）
    const teacherScheduleMap: Record<string, any> = {}
    const roomScheduleMap: Record<string, any> = {}
    const rules = rulesRes.datas || []
    rules.forEach((rule: any) => {
      rule.weekdays?.forEach((day: number) => {
        const time = rule.start_time.split(':')[0]
        const scheduleData = {
          id: rule.id,
          offering_name: rule.offering?.name || '-',
          teacher_name: rule.teacher?.name || '-',
          teacher_id: rule.teacher_id,
          room_id: rule.room_id,
          room_name: rule.room?.name || '-',
          ...rule,
        }

        // 根據資源類型分別存儲
        if (rule.teacher_id) {
          const key = `${rule.teacher_id}-${time}-${day}`
          teacherScheduleMap[key] = scheduleData
        }
        if (rule.room_id) {
          const key = `${rule.room_id}-${time}-${day}`
          roomScheduleMap[key] = scheduleData
        }
      })
    })

    // 存儲為包含兩個視圖的對象
    schedules.value = {
      teachers: teacherScheduleMap,
      rooms: roomScheduleMap
    }
  } catch (error) {
    console.error('Failed to fetch data:', error)
  }
}

const getScheduleAt = (resourceId: number, time: number) => {
  // 確保 schedules 結構存在
  if (!schedules.value || !schedules.value.teachers || !schedules.value.rooms) {
    return null
  }

  // 根據資源類型取得對應的排課 map
  const scheduleMap = resourceType.value === 'teacher' ? schedules.value.teachers : schedules.value.rooms

  // 搜尋所有星期
  for (const day of weekDays) {
    const key = `${resourceId}-${time}-${day.value}`
    if (scheduleMap[key]) {
      return scheduleMap[key]
    }
  }
  return null
}

const getCellClass = (resourceId: number, time: number): string => {
  const scheduleMap = resourceType.value === 'teacher' ? schedules.value.teachers : schedules.value.rooms
  const key = `${resourceId}-${time}`
  const validation = validationResults.value[key]

  if (validation?.valid === false) {
    return 'bg-critical-500/10 border-critical-500/50'
  } else if (validation?.warning) {
    return 'bg-warning-500/10 border-warning-500/50'
  } else if (validation?.valid === true) {
    return 'bg-success-500/10 border-success-500/50'
  } else if (getScheduleAt(resourceId, time)) {
    return 'bg-primary-500/10'
  }

  return ''
}

const getScheduleCardClass = (schedule: any): string => {
  if (!schedule) return ''
  return 'bg-slate-700/80 border border-white/10'
}

const selectSchedule = (resource: any, time: number) => {
  const scheduleMap = resourceType.value === 'teacher' ? schedules.value.teachers : schedules.value.rooms

  // 找到對應的星期
  for (const day of weekDays) {
    const key = `${resource.id}-${time}-${day.value}`
    if (scheduleMap[key]) {
      selectedCell.value = { resource, time, weekday: day.value }
      selectedSchedule.value = scheduleMap[key]
      emit('selectCell', { resource, time, weekday: day.value })
      return
    }
  }
}

const changeWeek = (delta: number) => {
  weekStart.value = getWeekStart(new Date(weekStart.value.getTime() + delta * 7 * 24 * 60 * 60 * 1000))
}

const handleDragEnter = (resourceId: number, time: number) => {
  dragTarget.value = { resourceId, time }
}

const handleDragLeave = () => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.resourceId}-${dragTarget.value.time}`
    delete validationResults.value[key]
  }
  dragTarget.value = null
}

const handleDragOver = (event: DragEvent) => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.resourceId}-${dragTarget.value.time}`
    validationResults.value[key] = { valid: true }
  }
}

const handleDrop = async (event: DragEvent) => {
  event.preventDefault()

  if (!dragTarget.value) return

  const data = event.dataTransfer?.getData('application/json')
  if (!data) return

  const parsed = JSON.parse(data)
  const { type, item } = parsed

  const key = `${dragTarget.value.resourceId}-${dragTarget.value.time}`
  validationResults.value[key] = { valid: 'checking' }

  try {
    const api = useApi()
    const teacherId = type === 'teacher' ? item.id : (item.teacher_id || null)
    const roomId = type === 'room' ? item.id : (item.room_id || null)

    const response = await api.post<any>('/admin/scheduling/check-overlap', {
      center_id: 1,
      teacher_id: teacherId,
      room_id: roomId,
      start_time: `${formatDate(weekStart.value)}T${formatTime(dragTarget.value.time)}:00`,
      end_time: `${formatDate(weekStart.value)}T${formatTime(dragTarget.value.time + 1)}:00`,
    })

    if (response.data.valid) {
      validationResults.value[key] = { valid: true }
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
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>
