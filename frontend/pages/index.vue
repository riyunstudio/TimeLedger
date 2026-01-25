<template>
  <div class="min-h-screen bg-gradient-mesh flex flex-col">
    <!-- Hero Section -->
    <header class="relative overflow-hidden">
      <div class="absolute inset-0 z-0">
        <!-- Animated Mesh Gradient Background -->
        <div class="absolute inset-0 animate-gradient">
          <div class="gradient-orb-1"></div>
          <div class="gradient-orb-2"></div>
          <div class="gradient-orb-3"></div>
          <div class="gradient-orb-4"></div>
          <div class="gradient-orb-5"></div>
          <div class="gradient-orb-6"></div>
        </div>
      </div>

      <!-- Content Overlay -->
      <div class="relative z-10">
        <div class="max-w-6xl mx-auto px-4 py-20 text-center">
          <!-- Logo -->
          <div class="mb-8">
            <div class="w-24 h-24 mx-auto mb-4 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-full flex items-center justify-center shadow-2xl">
              <svg class="w-12 h-12 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-6-4 4m6-6h6a2 2 0 012 0V8a2 2 0 01-2-2H6a2 2 0 00-2 2v12" />
              </svg>
            </div>
          </div>

          <!-- Title -->
          <h1 class="text-5xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-show-in-view-fade mb-4">
            TimeLedger
          </h1>
          <p class="text-2xl text-slate-300 font-medium">
            您的時間資產管家
          </p>

          <!-- Slogan -->
          <div class="mb-12">
            <p class="text-xl text-white font-medium">
              多中心排課
              <span class="mx-4 text-slate-400">× </span>
              老師人才庫
              <span class="text-white/80">× </span>
              專業課表分享
            </p>
          </div>

          <!-- CTA Buttons -->
          <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <button
              @click="handleLineLogin"
              :disabled="loading"
              class="w-full sm:w-auto btn-primary flex items-center justify-center gap-3 py-4 px-8 text-lg shadow-2xl"
            >
              <svg v-if="loading" class="animate-spin w-6 h-6" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v16a8 8 0 00-8 8h16a8 8 0 018-8zm-9 8a8 8 0 11-14 0v16a8 8 0 00-16 8H4a8 8 0 01-8 8v-16zm10 0a6 6 0 00-6 6v20a6 6 0 016 6H4a6 6 0 0016 6v20zM2 2a8 8 0 01-8 8v24a8 8 0 018-8 8H4a8 8 0 0016 8V6z" />
              </svg>
              <svg v-else class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2C6.477 2 2 6.477 2 12c0 1.891.5 3.657 1.363 5.203L2 22l4.797-1.363A9.956 9.956 0 0 0 12 22c5.523 0 10-4.477 10-10S17.523 2 12 2zm0 18c-1.455 0-2.822-.395-4-1.088V16c0-1.103.897-2 2-2h4c1.103 0 2 .897 2 2v.912C14.822 19.605 13.455 20 12 20zm6-5h-1v-1c0-2.206-1.794-4-4-4s-4 1.794-4 4v1H8v-1c0-2.762 2.238-5 5-5s5 2.238 5 5v1z"/>
              </svg>
              <span class="hidden sm:inline">使用 LINE 登入</span>
              <span class="sm:hidden">LINE 登入</span>
            </button>
            <NuxtLink
              to="/admin/login"
              class="w-full sm:w-auto glass-btn px-8 py-4 text-lg font-medium hover:bg-white/10 transition-all"
            >
              中心管理員登入
            </NuxtLink>
          </div>

          <!-- Admin Entry Link -->
          <p class="mt-4">
            <NuxtLink to="/admin/login" class="text-sm text-slate-400 underline-offset-2 hover:text-primary-500 transition-colors">
              已是中心管理員？點擊此處登入
            </NuxtLink>
          </p>
        </div>
      </div>

      <!-- Demo Section -->
      <div class="relative z-10 mt-16">
        <div class="text-center mb-8">
          <h2 class="text-3xl font-bold text-slate-100 mb-2">
            🎯 立手體驗
          </h2>
          <p class="text-lg text-slate-400">
            拖曳課程卡片到時間格中，感受流暢度
          </p>
        </div>

        <!-- Mini Sandbox -->
        <div class="max-w-3xl mx-auto glass-card p-3 sm:p-4 lg:p-6 overflow-x-auto">
          <!-- Mobile Navigation -->
          <div v-if="isMobile" class="flex items-center justify-between mb-2 sm:mb-3">
            <button
              @click="changeMobileDays(-1)"
              class="glass-btn p-2 rounded-lg shrink-0"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <span class="text-[clamp(12px,2.5vw,16px)] font-medium text-slate-200">{{ mobileDayLabel }}</span>
            <button
              @click="changeMobileDays(1)"
              class="glass-btn p-2 rounded-lg shrink-0"
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
                  class="p-0.5 min-h-[45px] sm:min-h-[60px] border-t border-l border-white/5 relative transition-all"
                  :class="{
                    'bg-success-500/10': !getScheduleAt(hour, day.value) && !hasConflict,
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
                    class="schedule-card p-1 rounded cursor-grab transition-all hover:scale-105 h-full"
                    :class="{
                      'bg-critical-500/20 border-critical-500/50': hasConflict,
                      'ring-2 ring-primary-500': selectedCell?.time === hour && selectedCell?.day === day.value
                    }"
                  >
                    <div class="font-medium text-[clamp(11px,2.5vw,15px)] truncate text-white leading-tight">
                      {{ getScheduleAt(hour, day.value)?.offering_name }}
                    </div>
                    <div class="text-[clamp(9px,2vw,12px)] text-slate-400 truncate leading-tight">
                      {{ getScheduleAt(hour, day.value)?.teacher_name }}
                    </div>
                  </div>

                  <div
                    v-if="!getScheduleAt(hour, day.value) && isDragging && isTargetCell(hour, day.value)"
                    class="absolute inset-0 border-2 border-dashed border-primary-500/50 bg-primary-500/10 flex items-center justify-center rounded"
                  >
                    <span class="text-[clamp(9px,1.8vw,11px)] text-primary-400">放置</span>
                  </div>
                </div>
              </template>
            </div>
          </div>
        </div>

        <div class="text-center mt-6">
          <div class="inline-flex items-center gap-2 mb-2">
             <div
               v-if="hasConflict"
               class="flex items-center gap-2 p-3 rounded-xl bg-critical-500/20 border border-critical-500/50"
             >
               <svg class="w-5 h-5 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                 <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 2.502-3.25V6.718c0-1.583 1.667-2.502-3.25H5.082c0-1.54 0-2.502-1.657-2.502-1.657a2 2 0 002.502 1.657a3.25z" />
               </svg>
               <span class="text-sm text-critical-500 font-medium">時段衝突！</span>
             </div>
             <div
               v-else-if="isDragging"
               class="flex items-center gap-2 p-3 rounded-xl bg-success-500/20 border border-success-500/50"
             >
               <svg class="w-5 h-5 text-success-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                 <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
               </svg>
               <span class="text-sm text-success-500 font-medium">可以放置</span>
             </div>
             <div
               v-else-if="selectedCell"
               class="flex items-center gap-2 p-3 rounded-xl bg-primary-500/20 border border-primary-500/50"
             >
               <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                 <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
               </svg>
               <span class="text-sm text-primary-500 font-medium">已選中課程，點擊目標格子移動</span>
             </div>
             <div v-else class="text-slate-400 text-sm">
               電腦：拖曳移動 / 手機：點擊選取再點擊放置
             </div>
          </div>

          <button
            @click="resetDemo"
            class="glass-btn px-6 py-3 rounded-xl font-medium"
          >
            重置
          </button>
        </div>
      </div>

      <!-- Feature Showcase -->
      <div class="max-w-5xl mx-auto mt-20 px-4">
        <h2 class="text-3xl font-bold text-slate-100 text-center mb-8">
          炫麗特色
        </h2>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <!-- Feature 1: Beautiful Export -->
          <div class="glass-card p-6 rounded-2xl">
            <div class="aspect-[9/16] rounded-xl overflow-hidden mb-4 relative">
              <div class="absolute inset-0 bg-gradient-mesh flex items-center justify-center">
                <div class="text-center">
                  <h3 class="text-2xl font-bold text-white mb-2">我的課表</h3>
                  <div class="flex items-center justify-center gap-2 mb-4">
                    <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
                      <span class="text-xl font-bold text-white">A</span>
                    </div>
                    <div class="space-y-2">
                      <div class="bg-white/20 rounded-lg p-2">
                        <p class="text-xs text-slate-800">14:00 鋼琴基礎</p>
                        <p class="text-xs text-slate-800">15:00 小提琴入門</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <h3 class="text-lg font-semibold text-slate-100 text-center">精美分享圖</h3>
            <p class="text-sm text-slate-400 text-center mt-2">
              多種漸層主題、個人品牌資訊、社交分享
            </p>
          </div>

          <!-- Feature 2: Smart Matching -->
          <div class="glass-card p-6 rounded-2xl">
            <div class="p-4 mb-4 rounded-xl bg-secondary-500/20 border border-secondary-500/50 flex items-center gap-4">
              <svg class="w-10 h-10 text-secondary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path d="M9 12l2 2 4-4m6 2v10a2 2 0 012 0V5a2 2 0 00-2-2H7a2 2 0 00-2 2v7a2 2 0 002-2 2h3.583c1.428-1.15 3.897a5 5 0 0 0 .0 2-3.5 0 00-2-3.5 0 00-3.5a6 6 0 0 0 3.5a6 6 0 0 0 0 0 0 0 4 0 0z" />
              </svg>
              <div>
                <h4 class="font-semibold text-slate-100">智慧媒合</h4>
                <p class="text-sm text-slate-400 mt-2">
                  技能、證照、評分綜合評比
                </p>
              </div>
            </div>
            <div class="mt-4">
              <div class="bg-slate-800/50 rounded-lg p-3">
                <p class="text-xs text-white mb-1">🎯 Alice</p>
                <p class="text-xs text-slate-200">鋼琴 · 高級 · 5星</p>
              </div>
            </div>
            <div class="flex items-center gap-4">
              <div class="bg-primary-500/20 rounded-lg px-3 py-2">
                <p class="text-xs font-white font-medium">技能</p>
                <span class="px-2 py-0.5 rounded-full bg-white/20 text-xs text-primary-500">鋼琴</span>
                <span class="px-2 py-0.5 rounded-full bg-white/20 text-xs text-secondary-500">古典</span>
              </div>
              <div class="bg-warning-500/20 rounded-lg px-3 py-2">
                <p class="text-xs text-white font-medium">證照</p>
                <span class="px-2 py-0.5 rounded-full bg-success-500/20 text-xs text-white">ABRSM</span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <p class="text-xs text-slate-200">96 分匹配</p>
            </div>
          </div>

          <!-- Feature 3: Real-time Notifications -->
          <div class="glass-card p-6 rounded-2xl">
            <div class="p-4 mb-4 rounded-xl bg-critical-500/20 border border-critical-500/50 flex items-center gap-4">
              <svg class="w-10 h-10 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5v-2a3 3 0 00-3 3H6a3 3 0 00-3 3v11a3 3 0 013 3 0 011.593a1.5 3 0 003.167a1.5 3 0 01.593a1.5 3 0 011.593a1.5.3 0 011.593a1.5.3 0 013.3 0 00-3.167a1.5 3 0 011.593a1.5.3 0 003.167a1.5.3 0 011.593a1.5.3 0 00.3.167a1.5.3 0 013.593a1.5 3 0 011.593a1.5 3 0 011.593a1.5 3 0 011.593a1.5 3 0 013.3 0 011.593a1.5 3 0 011.593a1.5 3 0 011.593a1.5.3 0 011.593a1.5 3 0 013.3 0 011.593a1.5 3 0 011.593a1.5 3 0 013.3 0 011.593a1.5 3 0 011.593a1.5.3 0 011.593a1.5 3 0 011.593a1.5.3 0 013.3 0 011.593a1.5 3 0 011.593a1.5 3 0 011.593a1.5 3 0 011.593a1.5.3 0 011.593a1.5 3 0 013.3 0 011.593a1.5.3 0 0 011.593a1.5 3 0 011.593a1.5 3 0 013.3 0 011.593a1.5 3 0 011.593a1.1 1 1 0 0.6 16.5v5a1.1 0 015.716 3.953a2 2.0 0 01.034-8.808-4.053a2.2 0 018.953a3.953a3.953 0 016 8 0 018.953a5.953a4.808 3 4.053a2 2 0 018.8.96-1.062.744.0 00.905 1.5 0 006 0z" />
              </svg>
              <div>
                <h4 class="font-semibold text-slate-100">即時通知</h4>
                <p class="text-sm text-slate-400 mt-2">
                  LINE Notify 推送課程變更、審核結果
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- CTA -->
        <div class="text-center mt-12">
          <button
            @click="handleLineLogin"
            :disabled="loading"
            class="btn-primary px-8 py-4 text-lg shadow-2xl animate-bounce"
          >
            {{ loading ? '連接中...' : '開始使用 TimeLedger' }}
          </button>
        </div>

        <!-- Footer -->
        <div class="text-center mt-20 text-slate-500 text-sm">
          <p>© 2026 TimeLedger Team. All rights reserved.</p>
        </div>
      </div>
    </header>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

const config = useRuntimeConfig()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref('')

let liff: any = null

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

const getScheduleAt = (hour: number, weekday: number): any => {
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

const handleLineLogin = async () => {
  loading.value = true

  try {
    // 跳轉到教師登入頁面，帶上預設測試老師的 LINE User ID 和 Access Token
    const demoLineUserId = 'U1234567890abcdef1234567890abcd'
    const demoAccessToken = 'test_access_token_demo001'
    router.push(`/teacher/login?line_user_id=${demoLineUserId}&access_token=${demoAccessToken}`)
  } catch (err) {
    console.error('Login failed:', err)
    error.value = '登入失敗，請稍後再試'
  } finally {
    loading.value = false
  }
}

const schedules = ref<Record<string, any>>({
  '14-1': { id: 1, offering_name: '鋼琴基礎', teacher_name: 'Alice', room_id: 1 },
  '14-2': { id: 2, offering_name: '鋼琴基礎', teacher_name: 'Alice', room_id: 1 },
  '14-3': { id: 3, offering_name: '小提琴入門', teacher_name: 'Bob', room_id: 2 },
  '15-1': { id: 4, offering_name: '鋼琴進階', teacher_name: 'Alice', room_id: 1 },
})

const isDragging = ref(false)
const dragTarget = ref<{ time: number, day: number } | null>(null)
const hasConflict = ref(false)
const selectedCell = ref<{ time: number, day: number } | null>(null)

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
</script>

<style scoped>
.animate-gradient {
  animation: gradient-shift 20s linear infinite;
}

@keyframes gradient-shift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

.gradient-orb-1,
.gradient-orb-2,
.gradient-orb-3,
.gradient-orb-4,
.gradient-orb-5,
.gradient-orb-6 {
  position: absolute;
  border-radius: 50%;
  animation: gradient-shift 20s linear infinite;
}

.gradient-orb-1 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #6366F1 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-2 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #A855F7 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-3 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #1e3a8a 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-4 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #1e3a8a 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-5 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #1e3a8a 0%, transparent 50%);
  background-size: 200% 200%;
}

.gradient-orb-6 {
  top: -50%;
  left: -50%;
  background: radial-gradient(circle at 50% 50%, #6366F1 0%, transparent 50%);
  background-size: 200% 200%;
}

.schedule-card {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.schedule-card:hover {
  transform: scale(1.02);
}

.animate-bounce {
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-20px);
  }
 60% {
    transform: translateY(-10px);
  }
}

.bg-show-in-view-fade {
  animation: show-in-fade 1s ease-out forwards;
}

@keyframes show-in-fade {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
