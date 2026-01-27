<template>
  <div class="p-4 md:p-6">
    <!-- 今日課表摘要 -->
    <div class="mb-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold text-white">今日課表摘要</h2>
        <span class="text-sm text-slate-400">{{ formatDate(today) }}</span>
      </div>

      <!-- 摘要統計卡片 -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
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
          <div class="mt-2 text-xs text-slate-400">
            {{ todayStats.inProgressTeacherNames.length > 0 ? todayStats.inProgressTeacherNames.join('、') : '無進行中課程' }}
          </div>
        </div>

        <!-- 待審核申請 -->
        <div class="glass-card p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-slate-400">待審核</p>
              <p class="text-2xl font-bold text-white mt-1">{{ todayStats.pendingExceptions }}</p>
            </div>
            <div class="w-10 h-10 rounded-lg bg-warning-500/20 flex items-center justify-center">
              <svg class="w-5 h-5 text-warning-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
          </div>
          <div class="mt-2 text-xs text-slate-400">
            {{ todayStats.pendingExceptions > 0 ? '點擊查看詳情' : '無待審核項目' }}
          </div>
        </div>

        <!-- 異動提醒 -->
        <div class="glass-card p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-slate-400">異動提醒</p>
              <p class="text-2xl font-bold text-white mt-1">{{ todayStats.changesCount }}</p>
            </div>
            <div class="w-10 h-10 rounded-lg bg-critical-500/20 flex items-center justify-center">
              <svg class="w-5 h-5 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
              </svg>
            </div>
          </div>
          <div class="mt-2 text-xs text-slate-400">
            {{ todayStats.hasScheduleChanges ? '課表有調整' : '無異動' }}
          </div>
        </div>
      </div>

      <!-- 即將開始的課程 -->
      <div v-if="upcomingSessions.length > 0" class="glass-card p-4">
        <h3 class="text-sm font-medium text-white mb-3">即將開始</h3>
        <div class="space-y-2">
          <div
            v-for="session in upcomingSessions.slice(0, 3)"
            :key="session.id"
            class="flex items-center justify-between p-2 rounded-lg bg-white/5"
          >
            <div class="flex items-center gap-3">
              <span class="text-primary-500 font-mono text-sm">{{ session.time }}</span>
              <div>
                <p class="text-white text-sm">{{ session.courseName }}</p>
                <p class="text-slate-500 text-xs">{{ session.teacherName }} · {{ session.roomName }}</p>
              </div>
            </div>
            <span class="text-xs text-slate-400">{{ session.minutesUntil }} 分鐘後</span>
          </div>
        </div>
        <button
          v-if="upcomingSessions.length > 3"
          @click="viewAllUpcoming"
          class="mt-3 text-sm text-primary-500 hover:text-primary-400 flex items-center gap-1"
        >
          查看全部 {{ upcomingSessions.length }} 堂課
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 排課網格（支援週曆/矩陣視圖） -->
    <div class="h-full flex flex-col lg:flex-row gap-6">
      <ScheduleGrid
        class="flex-1 min-w-0"
        v-model:view-mode="viewMode"
        v-model:selected-resource-id="selectedResourceId"
      />
      <ScheduleResourcePanel
        class="lg:w-80 shrink-0"
        :view-mode="resourcePanelViewMode"
        @select-resource="handleSelectResource"
      />
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

const notificationStore = useNotificationStore()
const notificationUI = useNotification()
const router = useRouter()

// 視圖模式：'calendar' | 'teacher_matrix' | 'room_matrix'
const viewMode = ref<'calendar' | 'teacher_matrix' | 'room_matrix'>('calendar')
// 選中的資源 ID
const selectedResourceId = ref<number | null>(null)

// 今日日期
const today = computed(() => new Date())

// 格式化日期
const formatDate = (date: Date): string => {
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
}

// 今日課表統計
const todayStats = ref({
  totalSessions: 0,
  completedSessions: 0,
  upcomingSessions: 0,
  inProgressSessions: 0,
  inProgressTeacherNames: [] as string[],
  pendingExceptions: 0,
  changesCount: 0,
  hasScheduleChanges: false
})

// 即將開始的課程
const upcomingSessions = ref<Array<{
  id: number
  time: string
  courseName: string
  teacherName: string
  roomName: string
  minutesUntil: number
}>>([])

// 模擬今日課表數據（實際應從 API 獲取）
const fetchTodayStats = async () => {
  try {
    const api = useApi()
    // 嘗試獲取今日課表數據
    const response = await api.get<any>('/admin/dashboard/today-summary')

    if (response.datas) {
      const sessions = response.datas.sessions || []
      const now = new Date()

      // 計算統計
      let completed = 0
      let inProgress = 0
      let upcoming = 0
      const inProgressTeachers: string[] = []
      const upcomingList: typeof upcomingSessions.value = []

      sessions.forEach((session: any) => {
        const startTime = new Date(session.start_time)
        const endTime = new Date(session.end_time)

        if (now >= endTime) {
          completed++
        } else if (now >= startTime && now < endTime) {
          inProgress++
          inProgressTeachers.push(session.teacher?.name || '未知老師')
        } else {
          upcoming++
          const minutesUntil = Math.floor((startTime.getTime() - now.getTime()) / 60000)
          upcomingList.push({
            id: session.id,
            time: startTime.toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit' }),
            courseName: session.offering?.name || '未知課程',
            teacherName: session.teacher?.name || '未知老師',
            roomName: session.room?.name || '未知教室',
            minutesUntil
          })
        }
      })

      upcomingList.sort((a, b) => a.minutesUntil - b.minutesUntil)

      todayStats.value = {
        totalSessions: response.datas.totalSessions || sessions.length,
        completedSessions: completed,
        upcomingSessions: upcoming,
        inProgressSessions: inProgress,
        inProgressTeacherNames: inProgressTeachers,
        pendingExceptions: response.datas.pendingExceptions || 0,
        changesCount: response.datas.changesCount || 0,
        hasScheduleChanges: response.datas.hasScheduleChanges || false
      }

      upcomingSessions.value = upcomingList
    }
  } catch (error) {
    console.log('今日課表 API 尚未實作或無數據，使用模擬數據')
    // 使用模擬數據
    loadMockTodayStats()
  }
}

// 模擬今日課表數據（展示用）
const loadMockTodayStats = () => {
  const now = new Date()
  const currentHour = now.getHours()

  todayStats.value = {
    totalSessions: 8,
    completedSessions: currentHour > 12 ? 5 : 2,
    upcomingSessions: currentHour > 12 ? 3 : 6,
    inProgressSessions: currentHour >= 10 && currentHour < 12 ? 1 : 0,
    inProgressTeacherNames: currentHour >= 10 && currentHour < 12 ? ['林老師'] : [],
    pendingExceptions: 3,
    changesCount: 2,
    hasScheduleChanges: true
  }

  // 模擬即將開始的課程
  upcomingSessions.value = [
    { id: 1, time: '14:00', courseName: '瑜珈基礎', teacherName: '林老師', roomName: 'A教室', minutesUntil: 30 },
    { id: 2, time: '15:30', courseName: '鋼琴入門', teacherName: '陳老師', roomName: 'B教室', minutesUntil: 90 },
    { id: 3, time: '16:00', courseName: '舞蹈課程', teacherName: '王老師', roomName: '多功能廳', minutesUntil: 120 }
  ]
}

// 查看全部即將開始的課程
const viewAllUpcoming = () => {
  // 跳轉到課表頁面，並顯示今日視圖
  router.push('/admin/schedules')
}

// 資源面板的視角模式
const resourcePanelViewMode = computed(() => {
  if (viewMode.value === 'teacher_matrix') return 'teacher'
  if (viewMode.value === 'room_matrix') return 'room'
  return 'offering'
})

const handleSelectResource = (resource: { type: 'teacher' | 'room', id: number } | null) => {
  if (!resource) {
    viewMode.value = 'calendar'
    selectedResourceId.value = null
  } else {
    if (resource.type === 'teacher') {
      viewMode.value = 'teacher_matrix'
    } else {
      viewMode.value = 'room_matrix'
    }
    selectedResourceId.value = resource.id
  }
}

onMounted(async () => {
  notificationStore.fetchNotifications()
  await fetchTodayStats()
})
</script>
