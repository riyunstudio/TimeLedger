<template>
  <div class="min-h-screen bg-slate-900">
    <TeacherHeader ref="headerRef" />

    <main class="p-4 pb-24 max-w-4xl mx-auto">
      <div class="flex items-center justify-between mb-6">
        <button
          @click="changeWeek(-1)"
          class="glass-btn p-2 rounded-lg shrink-0"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>

        <h2 class="text-lg sm:text-xl font-semibold text-slate-100 truncate px-2">
          {{ teacherStore.weekLabel }}
        </h2>

        <button
          @click="changeWeek(1)"
          class="glass-btn p-2 rounded-lg shrink-0"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>

      <div class="flex items-center justify-end mb-4">
        <div class="glass rounded-lg p-1 flex">
          <button
            @click="viewMode = 'grid'"
            class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
            :class="viewMode === 'grid' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-slate-200'"
          >
            <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
            </svg>
            網格
          </button>
          <button
            @click="viewMode = 'list'"
            class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
            :class="viewMode === 'list' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-slate-200'"
          >
            <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
            </svg>
            列表
          </button>
        </div>
      </div>

      <div
        v-if="teacherStore.schedule"
        class="space-y-4"
      >
        <!-- Grid View -->
        <div
          v-if="viewMode === 'grid'"
          class="glass-card p-3 sm:p-4 overflow-x-auto"
        >
          <!-- Mobile Navigation for Grid -->
          <div v-if="isMobile" class="flex items-center justify-between mb-2">
            <button
              @click="changeGridDay(-1)"
              class="glass-btn p-2 rounded-lg shrink-0"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <span class="text-sm font-medium text-slate-200">{{ gridDayLabel }}</span>
            <button
              @click="changeGridDay(1)"
              class="glass-btn p-2 rounded-lg shrink-0"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>

          <div class="min-w-[350px] sm:min-w-[500px]">
            <div class="grid" :class="gridColsClass" gap-0.5 sm:gap-1>
              <!-- Header Row -->
              <div class="p-1 sm:p-2 text-center bg-white/5 rounded-tl-lg">
                <span class="text-[10px] text-slate-400"></span>
              </div>
              <div
                v-for="day in displayWeekDays"
                :key="day.date"
                class="p-1 sm:p-2 text-center bg-white/5"
              >
                <div class="text-[10px] text-slate-400">{{ day.weekday }}</div>
                <div class="text-[clamp(10px,2vw,14px)] font-medium text-slate-100">{{ day.day }}</div>
              </div>

              <!-- Time Slots -->
              <template v-for="hour in timeSlots" :key="hour">
                <div class="p-1 sm:p-2 flex items-center justify-center border-t border-white/5">
                  <span class="text-[clamp(9px,1.8vw,12px)] text-slate-500">{{ hour }}:00</span>
                </div>
                
                <div
                  v-for="day in displayWeekDays"
                  :key="`${hour}-${day.date}`"
                  class="p-0.5 min-h-[45px] sm:min-h-[50px] border-t border-l border-white/5 relative"
                  :class="getGridCellClass(day.date, hour)"
                >
                  <div
                    v-for="item in getScheduleItemsAt(day.date, hour)"
                    :key="item.id"
                    class="rounded p-1 text-xs cursor-pointer hover:opacity-80 transition-opacity"
                    :class="getItemBgClass(item)"
                    @click="openItemDetail(item)"
                  >
                    <div class="font-medium text-[clamp(10px,2vw,14px)] truncate text-white leading-tight">{{ item.title }}</div>
                    <div class="text-[9px] sm:text-[10px] text-slate-300 truncate">{{ item.start_time }}</div>
                  </div>
                </div>
              </template>
            </div>
          </div>
        </div>

        <!-- List View - Daily View -->
        <div v-else class="space-y-4">
          <div class="flex items-center justify-between mb-4">
            <button
              @click="changeListDay(-1)"
              class="glass-btn p-2 rounded-lg"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>

            <h3 class="text-lg font-semibold text-slate-100">
              {{ formatDate(listCurrentDate) }}
            </h3>

            <button
              @click="changeListDay(1)"
              class="glass-btn p-2 rounded-lg"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>

          <div class="glass-card p-4">
            <div
              v-if="currentDayItems.length === 0"
              class="text-center py-12 text-slate-500"
            >
              今日無行程
            </div>

            <div
              v-else
              class="space-y-3"
            >
              <div
                v-for="item in currentDayItems"
                :key="item.id"
                class="border rounded-xl p-4 cursor-pointer hover:bg-white/5 transition-all"
                :class="getItemBorderClass(item)"
                @click="openItemDetail(item)"
              >
                <div class="flex items-start gap-4">
                  <div class="flex-shrink-0 w-16 text-center">
                    <div class="text-xs text-slate-500 mb-1">{{ item.start_time }}</div>
                    <div class="text-xs text-slate-600">-</div>
                    <div class="text-xs text-slate-500 mt-1">{{ item.end_time }}</div>
                  </div>
                  
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-1">
                      <span
                        class="w-2 h-2 rounded-full"
                        :style="{ backgroundColor: item.color || '#10B981' }"
                      ></span>
                      <h4 class="font-medium text-slate-100 truncate">
                        {{ item.title }}
                        <span v-if="(item.data as any)?.center_name" class="text-primary-400 font-normal">@{{ (item.data as any).center_name }}</span>
                      </h4>
                    </div>
                    <p v-if="item.type === 'SCHEDULE_RULE'" class="text-sm text-slate-400">
                      課程時段
                    </p>
                    <p v-else class="text-sm text-slate-400">
                      個人行程
                    </p>
                  </div>

                  <div
                    v-if="item.status"
                    class="flex-shrink-0 px-2 py-1 rounded-full text-xs font-medium"
                    :class="getStatusClass(item.status)"
                  >
                    {{ getStatusText(item.status) }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div
        v-else
        class="text-center py-12 text-slate-500"
      >
        載入中...
      </div>
    </main>

    <button
      @click="showPersonalEventModal = true"
      class="fixed bottom-6 right-6 w-14 h-14 rounded-full bg-gradient-to-r from-primary-500 to-secondary-500 flex items-center justify-center shadow-xl hover:scale-110 transition-transform duration-300 z-50"
    >
      <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
    </button>

    <PersonalEventModal
      v-if="showPersonalEventModal"
      @close="showPersonalEventModal = false"
    />

    <NotificationDropdown
      v-if="notificationUI.show.value"
      @close="notificationUI.close()"
    />

    <TeacherSidebar
      v-if="sidebarStore.isOpen.value"
      @close="sidebarStore.close()"
    />
    <SessionNoteModal
      :is-open="showSessionNoteModal"
      :schedule-item="selectedScheduleItem"
      @close="handleNoteModalClose"
      @saved="handleNoteModalSaved"
    />
  </div>
</template>

<script setup lang="ts">
import type { ScheduleItem, WeekSchedule } from '~/types'

definePageMeta({
  middleware: 'auth-teacher',
})

const teacherStore = useTeacherStore()
const sidebarStore = useSidebar()
const notificationUI = useNotification()
const showPersonalEventModal = ref(false)
const viewMode = ref('grid')
const listCurrentDate = ref('')
const isMobile = ref(false)
const gridDayOffset = ref(0)

const timeSlots = [9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]

const allWeekDays = computed(() => {
  if (!teacherStore.schedule) return []
  
  return teacherStore.schedule.days.map(day => {
    const date = new Date(day.date)
    return {
      date: day.date,
      weekday: date.toLocaleDateString('zh-TW', { weekday: 'short' }),
      day: date.getDate(),
    }
  })
})

const displayWeekDays = computed(() => {
  if (!isMobile.value) return allWeekDays.value
  
  const start = gridDayOffset.value
  const days = []
  for (let i = 0; i < 3; i++) {
    const index = (start + i) % 7
    if (allWeekDays.value[index]) {
      days.push(allWeekDays.value[index])
    }
  }
  return days
})

const gridColsClass = computed(() => {
  return isMobile.value ? 'grid-cols-[50px_repeat(3,1fr)]' : 'grid-cols-[80px_repeat(7,1fr)]'
})

const gridDayLabel = computed(() => {
  if (!displayWeekDays.value.length) return ''
  const start = displayWeekDays.value[0]
  const end = displayWeekDays.value[displayWeekDays.value.length - 1]
  return `${start.weekday} - ${end.weekday}`
})

const changeGridDay = (delta: number) => {
  gridDayOffset.value = (gridDayOffset.value + delta + 7) % 7
}

const weekDayHeaders = computed(() => allWeekDays.value)

const currentDayItems = computed(() => {
  if (!teacherStore.schedule || !listCurrentDate.value) return []
  
  const day = teacherStore.schedule.days.find(d => d.date === listCurrentDate.value)
  return day?.items || []
})

const changeListDay = (delta: number) => {
  if (!teacherStore.schedule) return
  
  const currentIndex = teacherStore.schedule.days.findIndex(d => d.date === listCurrentDate.value)
  const newIndex = currentIndex + delta
  
  if (newIndex >= 0 && newIndex < teacherStore.schedule.days.length) {
    listCurrentDate.value = teacherStore.schedule.days[newIndex].date
  }
}

const getScheduleItemsAt = (date: string, hour: number): ScheduleItem[] => {
  const day = teacherStore.schedule?.days.find(d => d.date === date)
  if (!day) return []
  
  const hourNum = hour
  return day.items.filter(item => {
    const startHour = parseInt(item.start_time.split(':')[0])
    const endHour = parseInt(item.end_time.split(':')[0])
    return hourNum >= startHour && hourNum < endHour
  })
}

const getGridCellClass = (date: string, hour: number): string => {
  const items = getScheduleItemsAt(date, hour)
  if (items.length > 0) return ''
  
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const cellDate = new Date(date)
  const isPast = cellDate < today
  
  if (isPast) return 'bg-slate-800/50'
  return 'hover:bg-white/5'
}

const getItemBgClass = (item: ScheduleItem): string => {
  if (item.type === 'PERSONAL_EVENT') {
    return 'bg-purple-500/30 border border-purple-500/50'
  }
  
  // CENTER_SESSION 或 SCHEDULE_RULE 都視為課程
  const data = item.data as any
  if (data?.has_exception) {
    return 'bg-warning-500/30 border border-warning-500/50'
  }
  
  return 'bg-success-500/20 border border-success-500/30'
}

const changeWeek = (delta: number) => {
  teacherStore.changeWeek(delta)
  gridDayOffset.value = 0
  teacherStore.fetchSchedule().then(() => {
    if (teacherStore.schedule?.days.length) {
      listCurrentDate.value = teacherStore.schedule.days[0].date
    }
  })
}

const openItemDetail = (item: ScheduleItem) => {
  // SCHEDULE_RULE 或 CENTER_SESSION 都視為課程
  if ((item.type === 'SCHEDULE_RULE' || item.type === 'CENTER_SESSION') && item.data?.id) {
    selectedScheduleItem.value = item
    showSessionNoteModal.value = true
  } else {
    console.log('Open item detail:', item)
  }
}

const showSessionNoteModal = ref(false)
const selectedScheduleItem = ref<ScheduleItem | null>(null)

const handleNoteModalClose = () => {
  showSessionNoteModal.value = false
  selectedScheduleItem.value = null
}

const handleNoteModalSaved = () => {
  // Optionally refresh or show toast
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const diffDays = Math.floor((date.getTime() - today.getTime()) / (1000 * 60 * 60 * 24))

  if (diffDays === 0) return '今天'
  if (diffDays === 1) return '明天'
  if (diffDays === -1) return '昨天'

  return date.toLocaleDateString('zh-TW', {
    month: 'long',
    day: 'numeric',
    weekday: 'short',
  })
}

const getItemBorderClass = (item: ScheduleItem): string => {
  if (item.type === 'PERSONAL_EVENT') {
    return 'border-purple-500/50 bg-purple-500/10'
  }

  // CENTER_SESSION 或 SCHEDULE_RULE 都視為課程
  const data = item.data as any
  if (data?.has_exception) {
    return 'border-warning-500/50 bg-warning-500/10'
  }

  return 'border-success-500/50 bg-success-500/10'
}

const getStatusClass = (status: string): string => {
  switch (status) {
    case 'PENDING':
      return 'bg-warning-500/20 text-warning-500'
    case 'APPROVED':
      return 'bg-success-500/20 text-success-500'
    case 'REJECTED':
      return 'bg-critical-500/20 text-critical-500'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const getStatusText = (status: string): string => {
  switch (status) {
    case 'PENDING':
      return '待審核'
    case 'APPROVED':
      return '已核准'
    case 'REJECTED':
      return '已拒絕'
    default:
      return status
  }
}

onMounted(() => {
  const checkMobile = () => {
    isMobile.value = window.innerWidth < 640
  }
  checkMobile()
  window.addEventListener('resize', checkMobile)
  
  teacherStore.fetchCenters()
  teacherStore.fetchSchedule()
  
  if (teacherStore.schedule?.days.length) {
    listCurrentDate.value = teacherStore.schedule.days[0].date
  }
})
</script>
