<template>
  <div class="h-full flex flex-col glass-card overflow-hidden">
    <div class="p-4 border-b border-white/10">
      <div class="flex items-center justify-between">
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

        <div class="flex items-center gap-2">
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
              class="rounded-lg p-1 text-xs cursor-pointer hover:opacity-80 transition-opacity"
              :class="getScheduleCardClass(getScheduleAt(time, day.value))"
              @click="selectSchedule(time, day.value)"
            >
              <div class="font-medium truncate">
                {{ getScheduleAt(time, day.value)?.offering_name }}
              </div>
              <div class="text-slate-400 truncate">
                {{ getScheduleAt(time, day.value)?.teacher_name }}
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
}>()

const showCreateModal = ref(false)
const selectedCell = ref<{ time: number, day: number } | null>(null)
const selectedSchedule = ref<any>(null)
const dragTarget = ref<{ time: number, day: number } | null>(null)
const validationResults = ref<Record<string, any>>({})

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

const schedules = ref<Record<string, any>>({
  '9-1': { id: 1, offering_name: '鋼琴基礎', teacher_name: 'Alice', room_id: 1 },
  '10-1': { id: 2, offering_name: '鋼琴基礎', teacher_name: 'Alice', room_id: 1 },
  '14-1': { id: 3, offering_name: '鋼琴進階', teacher_name: 'Alice', room_id: 1 },
  '15-3': { id: 4, offering_name: '小提琴', teacher_name: 'Bob', room_id: 2 },
  '16-3': { id: 5, offering_name: '小提琴', teacher_name: 'Bob', room_id: 2 },
})

const changeWeek = (delta: number) => {
  weekStart.value = getWeekStart(new Date(weekStart.value.getTime() + delta * 7 * 24 * 60 * 60 * 1000))
}

const formatTime = (hour: number): string => {
  return `${hour}:00`
}

const getScheduleAt = (time: number, weekday: number) => {
  return schedules.value[`${time}-${weekday}`]
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
  emit('selectCell', { time, weekday: time })
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
  console.log('Schedule rule created')
}
</script>
