<template>
  <!-- 今日課表摘要區塊 -->
  <div class="mb-6" role="region" aria-label="今日課表摘要">
    <h2 class="text-lg font-semibold text-white mb-4">今日課表摘要</h2>
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4" role="list">
      <!-- 今日課程數 -->
      <div class="glass-card p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-400">今日課程</p>
            <p class="text-2xl font-bold text-white mt-1">{{ todayStats.totalSessions }}</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-primary-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </div>
        </div>
        <div class="mt-2 flex items-center gap-2 text-xs">
          <span class="text-success-500">{{ todayStats.completedSessions }} 已完成</span>
          <span class="text-slate-600">|</span>
          <span class="text-primary-500">{{ todayStats.upcomingSessions }} 待上課</span>
        </div>
      </div>

      <!-- 進行中課程 -->
      <div class="glass-card p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-400">進行中</p>
            <p class="text-2xl font-bold text-white mt-1">{{ todayStats.inProgressSessions }}</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-yellow-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-yellow-500 animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
        <div class="mt-2 text-xs text-slate-400 truncate">
          {{ todayStats.inProgressTeacherNames.length > 0 ? todayStats.inProgressTeacherNames.join('、') : '無進行中課程' }}
        </div>
      </div>

      <!-- 即將開始 -->
      <div class="glass-card p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-400">即將開始</p>
            <p class="text-2xl font-bold text-white mt-1">{{ upcomingSessions.length }}</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-secondary-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-secondary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
        <div class="mt-2 text-xs text-slate-400">
          <span v-if="upcomingSessions.length > 0">下一堂 {{ upcomingSessions[0]?.minutesUntil }} 分鐘後</span>
          <span v-else>今日無課程</span>
        </div>
      </div>

      <!-- 待審核申請 -->
      <div class="glass-card p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-400">待審核</p>
            <p class="text-2xl font-bold text-white mt-1">{{ pendingExceptions }}</p>
          </div>
          <div class="w-10 h-10 rounded-lg bg-warning-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-warning-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
        </div>
        <div class="mt-2 text-xs text-slate-400">
          <NuxtLink to="/teacher/exceptions" class="text-primary-500 hover:text-primary-400">
            查看詳情 →
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- 即將開始的課程列表 -->
    <div v-if="upcomingSessions.length > 0" class="mt-4 glass-card p-4" role="region" aria-label="即將開始的課程">
      <h3 class="text-sm font-medium text-white mb-3">即將開始</h3>
      <div class="space-y-2" role="list">
        <div
          v-for="session in upcomingSessions.slice(0, 3)"
          :key="session.id"
          role="listitem"
          class="flex items-center justify-between p-2 rounded-lg bg-white/5 hover:bg-white/10 transition-colors cursor-pointer"
          @click="openItemDetail(session)"
        >
          <div class="flex items-center gap-3">
            <span class="text-primary-500 font-mono text-sm">{{ session.time }}</span>
            <div>
              <p class="text-white text-sm">{{ session.title }}</p>
              <p class="text-slate-500 text-xs">{{ session.centerName }}</p>
            </div>
          </div>
          <span class="text-xs text-slate-400">{{ session.minutesUntil }} 分鐘後</span>
        </div>
      </div>
      <button
        v-if="upcomingSessions.length > 3"
        class="mt-3 text-sm text-primary-500 hover:text-primary-400 flex items-center gap-1"
      >
        查看全部 {{ upcomingSessions.length }} 堂課
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  </div>

  <!-- 教師課表週曆組件（使用內建導航列） -->
  <div class="mb-4">
    <TeacherScheduleGrid
      v-if="teacherStore.schedule"
      :schedules="transformedSchedules"
      :week-start="teacherStore.weekStart"
      @update:weekStart="handleWeekStartChange"
      @select-schedule="handleScheduleNoteAction"
      @add-personal-event="showPersonalEventModal = true"
      @add-exception="goToExceptions"
      @export="goToExport"
      @edit-personal-event="handleEditPersonalEvent"
      @delete-personal-event="handleDeletePersonalEvent"
      @personal-event-note="handlePersonalEventNote"
    />
  </div>

  <!-- 載入狀態 -->
  <div
    v-if="!teacherStore.schedule && !teacherStore.loading"
    class="text-center py-12 text-slate-500"
  >
    載入中...
  </div>

  <!-- 骨架屏載入狀態 -->
  <div v-else-if="teacherStore.loading" class="space-y-4">
    <!-- 今日摘要骨架屏 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div v-for="i in 4" :key="i" class="glass-card p-4 animate-pulse">
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <div class="h-3 w-16 bg-white/10 rounded mb-2"></div>
            <div class="h-6 w-12 bg-white/10 rounded"></div>
          </div>
          <div class="w-10 h-10 bg-white/10 rounded-lg"></div>
        </div>
      </div>
    </div>

    <!-- 快捷操作骨架屏 -->
    <div class="flex gap-3">
      <div v-for="i in 3" :key="i" class="h-10 w-32 bg-white/10 rounded-lg animate-pulse"></div>
    </div>

    <!-- 課表骨架屏 -->
    <div class="glass-card p-4 animate-pulse">
      <div class="h-10 w-48 bg-white/10 rounded mb-4"></div>
      <div class="space-y-2">
        <div v-for="i in 5" :key="i" class="h-12 bg-white/10 rounded"></div>
      </div>
    </div>
  </div>

  <button
    @click="showPersonalEventModal = true"
    class="fixed bottom-24 md:bottom-6 right-6 w-14 h-14 rounded-full bg-gradient-to-r from-primary-500 to-secondary-500 flex items-center justify-center shadow-xl hover:scale-110 transition-transform duration-300 z-50"
  >
    <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
    </svg>
  </button>

  <PersonalEventModal
    v-if="showPersonalEventModal"
    :editing-event="editingEvent"
    @close="showPersonalEventModal = false; editingEvent = null"
    @saved="handlePersonalEventSaved"
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

  <PersonalEventNoteModal
    :is-open="showPersonalEventNoteModal"
    :event="selectedPersonalEvent"
    @close="showPersonalEventNoteModal = false; selectedPersonalEvent = null"
    @saved="handlePersonalEventNoteSaved"
  />
</template>

<script setup lang="ts">
import type { ScheduleItem, WeekSchedule } from '~/types'
import { alertError } from '~/composables/useAlert'
import { formatDateToString } from '~/composables/useTaiwanTime'

definePageMeta({
  middleware: 'auth-teacher',
  layout: 'default',
})

const teacherStore = useTeacherStore()
const sidebarStore = useSidebar()
const notificationUI = useNotification()
const router = useRouter()
const showPersonalEventModal = ref(false)
const editingEvent = ref<any>(null)

// 課表資料轉換（包含中心課程和個人行程）
const transformedSchedules = computed(() => {
  if (!teacherStore.schedule) return []

  const result: any[] = []

  teacherStore.schedule.days.forEach(day => {
    const date = new Date(day.date)
    const weekday = date.getDay() === 0 ? 7 : date.getDay()

    // 處理中心課程
    day.items.forEach(item => {
      const [startHour, startMinute] = item.start_time.split(':').map(Number)
      const [endHour, endMinute] = item.end_time.split(':').map(Number)
      const durationMinutes = (endHour * 60 + endMinute) - (startHour * 60 + startMinute)

      result.push({
        id: item.id,
        key: `${item.id}-${weekday}-${item.start_time}`,
        offering_name: item.title,
        center_name: (item.data as any)?.center_name || '',
        teacher_name: '',
        weekday: weekday,
        start_time: item.start_time,
        end_time: item.end_time,
        start_hour: startHour,
        start_minute: startMinute,
        duration_minutes: durationMinutes,
        date: day.date,
        has_exception: (item.data as any)?.has_exception || false,
        exception_type: (item.data as any)?.exception_type || null,
        data: item.data,
        is_personal_event: false,
        type: item.type,
        rule_id: item.rule_id, // 確保 rule_id 可用於課堂備註
      })
    })

    // 處理個人行程
    const dayEvents = teacherStore.personalEvents.filter(e => {
      const eventDate = new Date(e.start_at).toISOString().split('T')[0]
      return eventDate === day.date
    })

    dayEvents.forEach(event => {
      const startDate = new Date(event.start_at)
      const endDate = new Date(event.end_at)
      const [startHour, startMinute] = [startDate.getHours(), startDate.getMinutes()]
      const [endHour, endMinute] = [endDate.getHours(), endDate.getMinutes()]
      const durationMinutes = (endDate.getTime() - startDate.getTime()) / (1000 * 60)

      result.push({
        id: event.id,
        key: `personal_${event.id}-${weekday}-${startHour.toString().padStart(2, '0')}:${startMinute.toString().padStart(2, '0')}`,
        offering_name: event.title,
        center_name: '',
        teacher_name: '',
        weekday: weekday,
        start_time: `${startHour.toString().padStart(2, '0')}:${startMinute.toString().padStart(2, '0')}`,
        end_time: `${endHour.toString().padStart(2, '0')}:${endMinute.toString().padStart(2, '0')}`,
        start_hour: startHour,
        start_minute: startMinute,
        duration_minutes: durationMinutes,
        date: day.date,
        has_exception: false,
        exception_type: null,
        data: event,
        is_personal_event: true,
        type: 'PERSONAL_EVENT',
        color_hex: event.color_hex,
      })
    })
  })

  return result
})

// 週次範圍標籤
const weekRangeLabel = computed(() => {
  return teacherStore.weekLabel || '本週'
})

const changeWeek = (delta: number) => {
  teacherStore.changeWeek(delta)
  teacherStore.fetchSchedule().then(() => {
    calculateTodayStats()
  })
  teacherStore.fetchPersonalEvents()
}

const handleWeekStartChange = (newStart: Date) => {
  teacherStore.weekStart = newStart
}

// 處理課表備註動作（來自 ScheduleGrid）
const handleScheduleNoteAction = (schedule: any) => {
  // 如果是課堂備註動作
  if (schedule.action === 'note') {
    const itemData = schedule.data || {}
    selectedScheduleItem.value = {
      id: schedule.id,
      type: 'SCHEDULE_RULE' as const,
      title: schedule.offering_name,
      start_time: schedule.start_time,
      end_time: schedule.end_time,
      date: schedule.date,
      room_id: itemData?.room_id || 0,
      center_id: itemData?.center_id || 0,
      center_name: itemData?.center_name || '',
      status: 'APPROVED' as const,
      data: itemData,
    }
    showSessionNoteModal.value = true
  }
}

// 處理編輯個人行程
const handleEditPersonalEvent = (event: any) => {
  editingEvent.value = event
  showPersonalEventModal.value = true
}

// 處理刪除個人行程
const handleDeletePersonalEvent = async (event: any) => {
  const confirmed = await alertConfirm(`確定要刪除行程「${event.title}」嗎？`)
  if (confirmed) {
    try {
      await teacherStore.deletePersonalEvent(event.id)
      await teacherStore.fetchSchedule()
    } catch (error) {
      console.error('Failed to delete personal event:', error)
      await alertError('刪除失敗，請稍後再試')
    }
  }
}

// 今日課表摘要相關
const todayStats = ref({
  totalSessions: 0,
  completedSessions: 0,
  upcomingSessions: 0,
  inProgressSessions: 0,
  inProgressTeacherNames: [] as string[]
})

const upcomingSessions = ref<Array<{
  id: number | string
  time: string
  title: string
  centerName: string
  minutesUntil: number
  data?: any
}>>([])

const pendingExceptions = ref(0)

// 計算今日課表統計
const calculateTodayStats = () => {
  if (!teacherStore.schedule) return

  const today = new Date()
  const todayStr = formatDateToString(today)

  const todayDay = teacherStore.schedule.days.find(d => d.date === todayStr)
  if (!todayDay) {
    todayStats.value = {
      totalSessions: 0,
      completedSessions: 0,
      upcomingSessions: 0,
      inProgressSessions: 0,
      inProgressTeacherNames: []
    }
    upcomingSessions.value = []
    return
  }

  const now = new Date()
  const currentHour = now.getHours()
  const currentMinute = now.getMinutes()

  let completed = 0
  let inProgress = 0
  let upcoming = 0
  const inProgressTeachers: string[] = []
  const upcomingList: typeof upcomingSessions.value = []

  todayDay.items.forEach((item: ScheduleItem) => {
    const [startHour, startMinute] = item.start_time.split(':').map(Number)
    const [endHour, endMinute] = item.end_time.split(':').map(Number)
    const startTime = startHour * 60 + startMinute
    const endTime = endHour * 60 + endMinute
    const currentTime = currentHour * 60 + currentMinute

    if (currentTime >= endTime) {
      completed++
    } else if (currentTime >= startTime && currentTime < endTime) {
      inProgress++
      const centerName = (item.data as any)?.center_name || '未知中心'
      if (!inProgressTeachers.includes(centerName)) {
        inProgressTeachers.push(centerName)
      }
    } else {
      upcoming++
      const minutesUntil = startTime - currentTime
      upcomingList.push({
        id: item.id,
        time: item.start_time,
        title: item.title,
        centerName: (item.data as any)?.center_name || '',
        minutesUntil,
        data: item.data
      })
    }
  })

  upcomingList.sort((a, b) => a.minutesUntil - b.minutesUntil)

  todayStats.value = {
    totalSessions: todayDay.items.length,
    completedSessions: completed,
    upcomingSessions: upcoming,
    inProgressSessions: inProgress,
    inProgressTeacherNames: inProgressTeachers
  }

  upcomingSessions.value = upcomingList
}

// 計算待審核申請數
const calculatePendingExceptions = () => {
  pendingExceptions.value = teacherStore.exceptions.filter(e => e.status === 'PENDING').length
}

// 跳轉到例外申請頁面
const goToExceptions = () => {
  router.push('/teacher/exceptions')
}

// 跳轉到匯出頁面
const goToExport = () => {
  router.push('/teacher/export')
}

const showSessionNoteModal = ref(false)
const selectedScheduleItem = ref<ScheduleItem | null>(null)
const showPersonalEventNoteModal = ref(false)
const selectedPersonalEvent = ref<any>(null)

// 處理個人行程備註
const handlePersonalEventNote = (event: any) => {
  selectedPersonalEvent.value = event
  showPersonalEventNoteModal.value = true
}

const handlePersonalEventNoteSaved = () => {
  // 可以選擇性重新獲取資料或顯示提示
}

const handleNoteModalClose = () => {
  showSessionNoteModal.value = false
  selectedScheduleItem.value = null
}

const handleNoteModalSaved = () => {
  // Optionally refresh or show toast
}

const handlePersonalEventSaved = () => {
  teacherStore.fetchPersonalEvents()
  teacherStore.fetchSchedule().then(() => {
    calculateTodayStats()
  })
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
  teacherStore.fetchCenters()
  teacherStore.fetchSchedule()
  teacherStore.fetchPersonalEvents()
  teacherStore.fetchExceptions()

  // 使用 watch 監聽 schedule 變化
  watch(() => teacherStore.schedule, () => {
    calculateTodayStats()
  }, { immediate: true })

  // 監聽例外申請變化
  watch(() => teacherStore.exceptions, () => {
    calculatePendingExceptions()
  }, { immediate: true })
})
</script>
