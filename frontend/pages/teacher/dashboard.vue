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

  <!-- 快捷操作按鈕 -->
  <div class="mb-6 flex flex-wrap gap-3" role="toolbar" aria-label="快捷操作">
    <button
      @click="showPersonalEventModal = true"
      aria-label="新增個人行程"
      class="px-4 py-2 rounded-lg bg-primary-500/20 text-primary-500 hover:bg-primary-500/30 transition-colors flex items-center gap-2 text-sm"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      新增個人行程
    </button>
    <button
      @click="goToExceptions"
      aria-label="申請請假或調課"
      class="px-4 py-2 rounded-lg bg-warning-500/20 text-warning-500 hover:bg-warning-500/30 transition-colors flex items-center gap-2 text-sm"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      請假/調課
    </button>
    <button
      @click="goToExport"
      aria-label="匯出課表"
      class="px-4 py-2 rounded-lg bg-secondary-500/20 text-secondary-500 hover:bg-secondary-500/30 transition-colors flex items-center gap-2 text-sm"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
      </svg>
      匯出課表
    </button>
  </div>

  <!-- 視圖切換和課表顯示 -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-4">
    <div class="flex items-center gap-4">
      <!-- 週次導航 -->
      <div class="flex items-center gap-1">
        <button
          @click="changeWeek(-1)"
          class="glass-btn p-2 rounded-lg hover:bg-white/10 transition-colors"
          title="上一週"
        >
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <span class="text-sm font-medium text-slate-200 min-w-[140px] text-center">
          {{ weekRangeLabel }}
        </span>
        <button
          @click="changeWeek(1)"
          class="glass-btn p-2 rounded-lg hover:bg-white/10 transition-colors"
          title="下一週"
        >
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>

      <div class="glass rounded-lg p-1 flex" role="tablist" aria-label="課表檢視模式">
        <button
          @click="viewMode = 'grid'"
          role="tab"
          aria-selected="true"
          aria-label="網格檢視"
          class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
          :class="viewMode === 'grid' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-slate-200'"
        >
          <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
          網格
        </button>
        <button
          @click="viewMode = 'list'"
          role="tab"
          aria-selected="false"
          aria-label="列表檢視"
          class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
          :class="viewMode === 'list' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-slate-200'"
        >
          <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
          </svg>
          列表
        </button>
      </div>
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
              @dragenter.prevent="handleDragEnter(hour, day.date)"
              @dragleave.prevent="handleDragLeave"
              @drop.prevent="handleDrop(hour, day.date)"
              @dragover.prevent
            >
              <div
                v-for="item in getScheduleItemsAt(day.date, hour)"
                :key="item.id"
                class="rounded p-1 text-xs cursor-grab hover:opacity-80 transition-opacity"
                :class="getItemBgClass(item)"
                draggable="true"
                @dragstart="handleDragStart(item, hour, day.date, $event)"
                @dragend="handleDragEnd"
                @click="openItemDetail(item)"
              >
                <div class="font-medium text-[clamp(10px,2vw,14px)] truncate text-white leading-tight">
                  {{ item.title }}
                </div>
                <div class="text-[9px] sm:text-[10px] text-slate-300 truncate">
                  {{ item.start_time }}
                  <span v-if="item.center_name" class="text-primary-400 ml-1">@{{ item.center_name }}</span>
                </div>
              </div>

              <div
                v-if="isDragging && isTargetCell(hour, day.date)"
                class="absolute inset-0 border-2 border-dashed border-primary-500/50 bg-primary-500/10 flex items-center justify-center rounded"
              >
                <span class="text-xs text-primary-400">放置</span>
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
</template>

<script setup lang="ts">
 import type { ScheduleItem, WeekSchedule } from '~/types'
 import { alertError } from '~/composables/useAlert'

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
 const viewMode = ref('grid')
 const listCurrentDate = ref('')
 const isMobile = ref(false)
 const gridDayOffset = ref(0)

 const timeSlots = [9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]

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
   today.setHours(0, 0, 0, 0)
   const todayStr = today.toISOString().split('T')[0]

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

// 週次範圍標籤（使用 store 中的 weekLabel）
const weekRangeLabel = computed(() => {
  return teacherStore.weekLabel || '本週'
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
  let items: ScheduleItem[] = day?.items || []

  // Add personal events
  const personalEventsAtHour = teacherStore.personalEvents
    .filter(event => {
      // Skip events with invalid dates
      if (!event.start_at || !event.end_at) return false
      const eventDateObj = new Date(event.start_at)
      if (isNaN(eventDateObj.getTime())) return false
      const eventDate = eventDateObj.toISOString().split('T')[0]
      if (eventDate !== date) return false
      const localStartHour = eventDateObj.getHours()
      const localEndHour = new Date(event.end_at).getHours()
      return hour >= localStartHour && hour < localEndHour
    })
    .map(event => {
      const startDateObj = new Date(event.start_at)
      const endDateObj = new Date(event.end_at)
      return {
        id: event.id,
        type: 'PERSONAL_EVENT' as const,
        title: event.title,
        date: date,
        start_time: `${startDateObj.getHours().toString().padStart(2, '0')}:${startDateObj.getMinutes().toString().padStart(2, '0')}`,
        end_time: `${endDateObj.getHours().toString().padStart(2, '0')}:${endDateObj.getMinutes().toString().padStart(2, '0')}`,
        room_id: 0,
        teacher_id: event.teacher_id,
        center_id: 0,
        center_name: '',
        status: 'APPROVED' as const,
        color: event.color,
        data: event,
      }
    })

  items = items.concat(personalEventsAtHour)

  const hourNum = hour
  return items.filter(item => {
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
    // 重新計算今日統計
    calculateTodayStats()
  })
  teacherStore.fetchPersonalEvents()
}

const openItemDetail = (item: ScheduleItem) => {
  // SCHEDULE_RULE 或 CENTER_SESSION 都視為課程
  if ((item.type === 'SCHEDULE_RULE' || item.type === 'CENTER_SESSION') && item.data?.id) {
    selectedScheduleItem.value = item
    showSessionNoteModal.value = true
  } else if (item.type === 'PERSONAL_EVENT') {
    // 編輯個人行程
    editingEvent.value = item.data as any
    showPersonalEventModal.value = true
  }
}

const showSessionNoteModal = ref(false)
const selectedScheduleItem = ref<ScheduleItem | null>(null)

const isDragging = ref(false)
const dragTarget = ref<{ time: number, date: string } | null>(null)
const draggedItem = ref<ScheduleItem | null>(null)
const sourceDate = ref<string>('')
const sourceHour = ref<number>(0)

const handleDragStart = (item: ScheduleItem, hour: number, date: string, event: DragEvent) => {
  isDragging.value = true
  draggedItem.value = item
  sourceHour.value = hour
  sourceDate.value = date
  event.dataTransfer?.setData('application/json', JSON.stringify({
    type: 'schedule',
    item,
    sourceHour: hour,
    sourceDate: date,
  }))
}

const handleDragEnd = () => {
  isDragging.value = false
  dragTarget.value = null
  draggedItem.value = null
}

const handleDragEnter = (hour: number, date: string) => {
  dragTarget.value = { time: hour, date }
}

const handleDragLeave = () => {
}

const handleDrop = async (hour: number, date: string) => {
  if (!isDragging.value || !draggedItem.value) return

  const sourceKey = `${sourceHour.value}-${sourceDate.value}`
  const targetKey = `${hour}-${date}`

  if (sourceKey !== targetKey && teacherStore.schedule) {
    const day = teacherStore.schedule.days.find(d => d.date === date)
    if (day && draggedItem.value) {
      const itemData = draggedItem.value.data
      let itemId = itemData?.id || draggedItem.value.id

      // Handle string IDs like "center_1023_rule_344_20260122"
      if (typeof itemId === 'string') {
        const numericMatch = itemId.match(/_(\d+)$/)
        if (numericMatch) {
          itemId = parseInt(numericMatch[1])
        } else {
          itemId = parseInt(itemId.replace(/\D/g, '')) || 0
        }
      }

      if (itemId) {
        const duration = parseInt(draggedItem.value.end_time.split(':')[0]) - parseInt(draggedItem.value.start_time.split(':')[0])
        const newEndHour = hour + duration
        const newStartTime = `${hour.toString().padStart(2, '0')}:00`
        const newEndTime = `${newEndHour.toString().padStart(2, '0')}:00`

        try {
          await teacherStore.moveScheduleItem({
            item_id: itemId,
            item_type: draggedItem.value.type as 'SCHEDULE_RULE' | 'PERSONAL_EVENT' | 'CENTER_SESSION',
            center_id: draggedItem.value.center_id || 1,
            new_date: date,
            new_start_time: newStartTime,
            new_end_time: newEndTime,
          })

          await teacherStore.fetchSchedule()
        } catch (error) {
          console.error('Failed to move schedule:', error)
          await alertError('更新失敗，請稍後再試')
        }
      }
    }
  }

  isDragging.value = false
  dragTarget.value = null
  draggedItem.value = null
}

const isTargetCell = (hour: number, date: string): boolean => {
  return dragTarget.value?.time === hour && dragTarget.value?.date === date
}

const handleNoteModalClose = () => {
  showSessionNoteModal.value = false
  selectedScheduleItem.value = null
}

const handleNoteModalSaved = () => {
  // Optionally refresh or show toast
}

const handlePersonalEventSaved = () => {
  // Refresh personal events and schedule
  teacherStore.fetchPersonalEvents()
  teacherStore.fetchSchedule().then(() => {
    // 重新計算統計
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
  teacherStore.fetchPersonalEvents()
  teacherStore.fetchExceptions()

  // 等待資料載入完成後計算統計
  const checkAndCalculateStats = () => {
    if (teacherStore.schedule) {
      calculateTodayStats()
    }
  }

  // 使用 watch 監聽 schedule 變化
  watch(() => teacherStore.schedule, () => {
    calculateTodayStats()
    if (teacherStore.schedule?.days.length) {
      listCurrentDate.value = teacherStore.schedule.days[0].date
    }
  }, { immediate: true })

  // 監聽例外申請變化
  watch(() => teacherStore.exceptions, () => {
    calculatePendingExceptions()
  }, { immediate: true })
})
</script>
