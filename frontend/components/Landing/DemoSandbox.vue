<template>
  <div class="py-20 px-4">
    <!-- Section Header -->
    <div class="max-w-3xl mx-auto text-center mb-10">
      <h2 class="text-3xl lg:text-4xl font-bold text-slate-100 mb-3">
        {{ $t('landing.demo.title') }}
      </h2>
      <p class="text-lg text-slate-400">
        {{ $t('landing.demo.subtitle') }}
      </p>
    </div>

    <!-- Mini Sandbox -->
    <div class="max-w-4xl mx-auto glass-card p-3 sm:p-4 lg:p-6 overflow-x-auto">
      <!-- Mobile Navigation -->
      <div v-if="isMobile" class="flex items-center justify-between mb-3 sm:mb-4">
        <button
          @click="changeMobileDays(-1)"
          class="glass-btn p-2.5 rounded-lg shrink-0 transition-all duration-300 hover:bg-white/10 hover:scale-105 active:scale-95"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <span class="text-[clamp(12px,2.5vw,16px)] font-medium text-slate-200">{{ mobileDayLabel }}</span>
        <button
          @click="changeMobileDays(1)"
          class="glass-btn p-2.5 rounded-lg shrink-0 transition-all duration-300 hover:bg-white/10 hover:scale-105 active:scale-95"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>

      <div class="min-w-[350px] sm:min-w-[500px]">
        <div class="grid" :class="displayGridCols" gap-0.5 sm:gap-1>
          <!-- Header Row -->
          <div class="p-1 sm:p-2 text-center bg-white/5 rounded-tl-lg">
            <span class="text-[clamp(9px,1.5vw,11px)] text-slate-400"></span>
          </div>
          <div v-for="day in displayWeekDays" :key="`header-${day.value}`" class="p-1 sm:p-2 text-center bg-white/5">
            <span class="text-[clamp(10px,2vw,14px)] text-slate-300 font-medium">{{ day.name }}</span>
          </div>

          <!-- Time Slots -->
          <template v-for="(hour, hourIndex) in timeSlots" :key="hour">
            <!-- Time Label -->
            <div class="p-1 flex items-center justify-center border-t border-white/5">
              <span class="text-[clamp(9px,1.8vw,12px)] text-slate-500">{{ formatTime(hour) }}</span>
            </div>

            <!-- Day Cells -->
            <div
              v-for="day in displayWeekDays"
              :key="`${hour}-${day.value}`"
              class="p-0.5 min-h-[45px] sm:min-h-[60px] border-t border-l border-white/5 relative transition-all duration-300"
              :class="{
                'bg-success-500/5': !getScheduleAt(hour, day.value) && !hasConflict,
                'ring-2 ring-primary-500/50': selectedCell && isTargetCell(hour, day.value),
                'bg-white/5': hourIndex % 2 === 0
              }"
              @dragenter.prevent="handleDragEnter(hour, day.value)"
              @dragleave.prevent="handleDragLeave"
              @drop.prevent="handleDrop"
              @dragover.prevent
              @click="selectCell(hour, day.value)"
            >
              <div
                v-if="getScheduleAt(hour, day.value)"
                draggable="true"
                @dragstart="handleScheduleDragStart(hour, day.value, $event)"
                @dragend="handleDragEnd"
                @click="handleScheduleClick(hour, day.value, $event)"
                class="schedule-card p-1.5 rounded-lg cursor-grab transition-all duration-300 hover:scale-[1.02] h-full relative overflow-hidden"
                :class="{
                  'bg-critical-500/20 border-critical-500/50': hasConflict,
                  'ring-2 ring-primary-500': selectedCell?.time === hour && selectedCell?.day === day.value
                }"
              >
                <!-- Gradient overlay -->
                <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/20 via-primary-500/10 to-purple-500/20 opacity-0 hover:opacity-100 transition-opacity duration-300"></div>
                
                <!-- Content -->
                <div class="relative z-10">
                  <div class="font-medium text-[clamp(11px,2.5vw,15px)] truncate text-white leading-tight">
                    {{ getScheduleAt(hour, day.value)?.offering_name }}
                  </div>
                  <div class="text-[clamp(9px,2vw,12px)] text-slate-400 truncate leading-tight">
                    {{ getScheduleAt(hour, day.value)?.teacher_name }}
                  </div>
                </div>
              </div>

              <!-- Drop zone highlight -->
              <div
                v-if="!getScheduleAt(hour, day.value) && isDragging && isTargetCell(hour, day.value)"
                class="absolute inset-0 border-2 border-dashed border-primary-500/50 bg-primary-500/10 flex items-center justify-center rounded-lg"
              >
                <span class="text-[clamp(9px,1.8vw,11px)] text-primary-400">{{ $t('landing.demo.drop_here') }}</span>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- Status Bar -->
    <div class="max-w-4xl mx-auto text-center mt-6">
      <div class="inline-flex items-center gap-2 mb-3">
        <!-- Conflict Status -->
        <div
          v-if="hasConflict"
          class="flex items-center gap-2 p-3 rounded-xl bg-critical-500/20 border border-critical-500/50 transition-all duration-300"
        >
          <svg class="w-5 h-5 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 2.502-3.25V6.718c0-1.583 1.667-2.502-3.25-2.502-1.657 0-2.502 0.845-2.502 2.502 0 1.657 0.845 2.502 2.502 2.502" />
          </svg>
          <span class="text-sm text-critical-500 font-medium">{{ $t('landing.demo.conflict') }}</span>
        </div>
        
        <!-- Ready to Drop -->
        <div
          v-else-if="isDragging"
          class="flex items-center gap-2 p-3 rounded-xl bg-success-500/20 border border-success-500/50 transition-all duration-300"
        >
          <svg class="w-5 h-5 text-success-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <span class="text-sm text-success-500 font-medium">{{ $t('landing.demo.drop_ready') }}</span>
        </div>
        
        <!-- Cell Selected -->
        <div
          v-else-if="selectedCell"
          class="flex items-center gap-2 p-3 rounded-xl bg-primary-500/20 border border-primary-500/50 transition-all duration-300"
        >
          <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="text-sm text-primary-500 font-medium">{{ $t('landing.demo.cell_selected') }}</span>
        </div>
        
        <!-- Default hint -->
        <div v-else class="text-slate-400 text-sm">
          {{ $t('landing.demo.drag_hint') }}
        </div>
      </div>

      <!-- Reset Button -->
      <button
        @click="resetDemo"
        class="glass-btn group relative inline-flex items-center justify-center gap-2 px-8 py-3.5 rounded-xl font-medium text-slate-300 transition-all duration-300 hover:text-white hover:bg-white/10 hover:shadow-indigo-500/20 hover:scale-105 active:scale-95"
      >
        <!-- Icon -->
        <svg class="w-5 h-5 transition-transform duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        <span>{{ $t('landing.demo.reset') }}</span>

        <!-- Hover glow -->
        <div class="absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300" style="box-shadow: 0 0 20px rgba(99, 102, 241, 0.2);"></div>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Schedule {
  id: number
  offering_name: string
  teacher_name: string
  room_id: number
}

interface DragTarget {
  time: number
  day: number
}

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

const isMobile = ref(false)
const mobileDayOffset = ref(0)

const schedules = ref<Record<string, Schedule>>({
  '14-1': { id: 1, offering_name: '鋼琴基礎', teacher_name: 'Alice', room_id: 1 },
  '14-2': { id: 2, offering_name: '鋼琴基礎', teacher_name: 'Alice', room_id: 1 },
  '14-3': { id: 3, offering_name: '小提琴入門', teacher_name: 'Bob', room_id: 2 },
  '15-1': { id: 4, offering_name: '鋼琴進階', teacher_name: 'Alice', room_id: 1 },
})

const isDragging = ref(false)
const dragTarget = ref<DragTarget | null>(null)
const hasConflict = ref(false)
const selectedCell = ref<DragTarget | null>(null)

const displayWeekDays = computed(() => {
  if (!isMobile.value) return weekDays
  const start = mobileDayOffset.value
  const days = []
  for (let i = 0; i < 3; i++) {
    const index = (start + i) % 7
    days.push(weekDays[index])
  }
  return days
})

const displayGridCols = computed(() => {
  return isMobile.value ? 'grid-cols-[50px_repeat(3,1fr)]' : 'grid-cols-[50px_repeat(7,1fr)]'
})

const mobileDayLabel = computed(() => {
  const start = displayWeekDays.value[0]
  const end = displayWeekDays.value[2]
  return `${start.name} - ${end.name}`
})

const changeMobileDays = (delta: number) => {
  mobileDayOffset.value = (mobileDayOffset.value + delta + 7) % 7
}

const getScheduleAt = (hour: number, weekday: number): Schedule | undefined => {
  return schedules.value[`${hour}-${weekday}`]
}

const isTargetCell = (hour: number, weekday: number): boolean => {
  return dragTarget.value?.time === hour && dragTarget.value?.day === weekday
}

const handleScheduleDragStart = (hour: number, weekday: number, event: DragEvent) => {
  isDragging.value = true
  const schedule = getScheduleAt(hour, weekday)
  event.dataTransfer?.setData('application/json', JSON.stringify({
    type: 'schedule',
    time: hour,
    day: weekday,
    data: schedule,
  }))
}

const handleDragEnd = () => {
  isDragging.value = false
  dragTarget.value = null
}

const handleDragEnter = (hour: number, weekday: number) => {
  dragTarget.value = { time: hour, day: weekday }
  checkConflict(hour, weekday)
}

const handleDragLeave = () => {
  // Don't clear dragTarget here, only clear on drop
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()

  if (!isDragging.value || !dragTarget.value) return

  const data = event.dataTransfer?.getData('application/json')
  if (!data) return

  const parsed = JSON.parse(data)
  if (parsed.type === 'schedule') {
    const targetKey = `${dragTarget.value.time}-${dragTarget.value.day}`
    const sourceKey = `${parsed.time}-${parsed.day}`

    schedules.value[targetKey] = parsed.data
    delete schedules.value[sourceKey]

    hasConflict.value = false
    isDragging.value = false
    dragTarget.value = null
  }
}

const selectCell = (hour: number, weekday: number) => {
  if (selectedCell.value) {
    const sourceKey = `${selectedCell.value.time}-${selectedCell.value.day}`
    const targetKey = `${hour}-${weekday}`
    const schedule = schedules.value[sourceKey]

    if (schedule) {
      schedules.value[targetKey] = schedule
      delete schedules.value[sourceKey]
      selectedCell.value = null
      hasConflict.value = false
    }
  }
}

const handleScheduleClick = (hour: number, weekday: number, event: MouseEvent) => {
  event.stopPropagation()

  if (selectedCell.value?.time === hour && selectedCell.value?.day === weekday) {
    selectedCell.value = null
  } else {
    selectedCell.value = { time: hour, day: weekday }
  }
}

const checkConflict = (hour: number, weekday: number) => {
  const key = `${hour}-${weekday}`
  const current = getScheduleAt(hour, weekday)
  const newSchedule = getScheduleAt(hour, weekday)

  if (current && newSchedule && current.id !== newSchedule.id) {
    hasConflict.value = true
  } else {
    hasConflict.value = false
  }
}

const resetDemo = () => {
  Object.keys(schedules.value).forEach(key => {
    if (key !== '14-1') {
      delete schedules.value[key]
    }
  })
  hasConflict.value = false
  selectedCell.value = null
}

const formatTime = (hour: number): string => {
  return `${hour}:00`
}

onMounted(() => {
  const checkMobile = () => {
    isMobile.value = window.innerWidth < 640
  }
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', () => {})
})
</script>

<style scoped>
/* Premium Glass Card Styling */
.glass-card {
  @apply bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl shadow-2xl shadow-black/20;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.glass-card:hover {
  @apply border-white/20;
}

/* Glass Button Styling */
.glass-btn {
  @apply bg-white/5 border border-white/10 transition-all duration-300;
}

.glass-btn:hover {
  @apply border-white/20 bg-white/10;
}

/* Enhanced Schedule Card Styling */
.schedule-card {
  @apply bg-gradient-to-br from-indigo-500/20 via-primary-500/10 to-purple-500/20 border border-white/10;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.schedule-card:hover {
  @apply shadow-indigo-500/20;
  transform: translateY(-2px) scale(1.02);
  box-shadow: 0 10px 40px -15px rgba(99, 102, 241, 0.3);
}

.schedule-card:active {
  @apply scale-[0.98];
}

/* Drop zone highlight */
.border-dashed {
  border-style: dashed;
}

/* Hover effects */
.glass-btn {
  @apply hover:shadow-indigo-500/20;
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .glass-card {
    @apply p-3 rounded-xl;
  }
}
</style>
