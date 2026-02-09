<template>
  <div class="h-full flex flex-col glass-card overflow-hidden">
    <!-- 頂部工具列 -->
    <div class="p-4 border-b border-white/10 shrink-0">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <!-- 週導航 -->
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

          <button
            @click="goToCurrentWeek"
            class="px-3 py-1.5 rounded-lg text-sm bg-slate-700/50 text-slate-300 hover:text-white hover:bg-slate-600 transition-colors"
          >
            回到本週
          </button>
        </div>

        <!-- 視圖控制 -->
        <div class="flex items-center gap-4">
          <!-- 資源類型切換 -->
          <div class="flex items-center gap-1 bg-slate-800/80 rounded-lg p-1">
            <button
              @click="viewModeModel = 'all'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="viewModeModel === 'all' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              全部
            </button>
            <button
              @click="viewModeModel = 'teacher'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="viewModeModel === 'teacher' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              老師
            </button>
            <button
              @click="viewModeModel = 'room'"
              class="px-3 py-1.5 rounded-md text-sm font-medium transition-all"
              :class="viewModeModel === 'room' ? 'bg-primary-500 text-white' : 'text-slate-400 hover:text-white'"
            >
              教室
            </button>
          </div>

          <!-- 資源選擇下拉 -->
          <select
            v-if="viewModeModel !== 'all'"
            v-model="selectedResourceIdModel"
            class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500"
          >
            <option :value="null">選擇{{ viewModeModel === 'teacher' ? '老師' : '教室' }}...</option>
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

      <!-- 已選資源提示 -->
      <div
        v-if="viewModeModel !== 'all' && selectedResourceName"
        class="mt-3 flex items-center gap-2 px-3 py-2 bg-primary-500/10 border border-primary-500/30 rounded-lg"
      >
        <span class="text-sm text-primary-400">
          {{ viewModeModel === 'teacher' ? '老師' : '教室' }}：
        </span>
        <span class="text-sm font-medium text-white">{{ selectedResourceName }}</span>
        <button
          @click="clearFilter"
          class="ml-auto p-1 hover:bg-white/10 rounded transition-colors"
        >
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 時間軸週曆 -->
    <div class="flex-1 overflow-auto">
      <div class="min-w-[900px]">
        <!-- 表頭：星期 -->
        <div class="grid" :style="{ gridTemplateColumns: `60px repeat(7, 1fr)` }">
          <div class="p-3 border-b border-white/10 bg-slate-800/50"></div>
          <div
            v-for="day in weekDays"
            :key="day.value"
            class="p-3 border-b border-white/10 text-center bg-slate-800/50"
          >
            <div class="text-sm font-medium text-slate-300">{{ day.name }}</div>
            <div class="text-xs text-slate-500">{{ getDayDate(day.value) }}</div>
          </div>
        </div>

        <!-- 時間軸 -->
        <div class="relative" :style="{ height: `${timeSlots.length * 60}px` }">
          <!-- 時間標記 -->
          <div
            v-for="hour in timeSlots"
            :key="hour"
            class="grid absolute w-full"
            :style="{
              gridTemplateColumns: `60px repeat(7, 1fr)`,
              top: `${hour * 60}px`,
              height: '60px'
            }"
          >
            <div class="p-2 text-right text-xs text-slate-500 border-r border-white/5">
              {{ formatTime(hour) }}
            </div>
            <div
              v-for="day in weekDays"
              :key="`${hour}-${day.value}`"
              class="relative border-r border-white/5 border-b border-dashed border-white/10"
            >
              <!-- 小時線 -->
              <div class="absolute inset-x-0 top-1/2 h-px bg-white/5"></div>
            </div>
          </div>

          <!-- 課程卡片容器 -->
          <div class="absolute inset-0" :style="{ top: '0' }">
            <!-- 課程卡片 -->
            <div
              v-for="session in filteredSessions"
              :key="`${session.id}-${session.date}`"
              class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-80 transition-opacity group"
              :style="getSessionStyle(session)"
              @click="selectSession(session)"
            >
              <div class="font-medium truncate text-white">
                {{ session.offering_name }}
              </div>
              <div class="text-slate-400 truncate">
                {{ session.time_range }}
              </div>
              <div v-if="viewModeModel === 'all'" class="text-slate-500 truncate">
                {{ session.teacher_name }}
              </div>

              <!-- 懸停 Tooltip -->
              <div
                v-if="!selectedSession"
                class="absolute left-0 bottom-full mb-2 hidden group-hover:block w-64 z-[150] pointer-events-none"
              >
                <div class="glass p-3 rounded-lg shadow-xl border border-white/10">
                  <div class="font-medium text-white mb-2">{{ session.offering_name }}</div>
                  <div class="space-y-1 text-xs">
                    <div class="flex justify-between text-slate-400">
                      <span>日期：</span>
                      <span class="text-slate-200">{{ session.formatted_date }}</span>
                    </div>
                    <div class="flex justify-between text-slate-400">
                      <span>時間：</span>
                      <span class="text-slate-200">{{ session.time_range }}</span>
                    </div>
                    <div class="flex justify-between text-slate-400">
                      <span>老師：</span>
                      <span class="text-slate-200">{{ session.teacher_name }}</span>
                    </div>
                    <div class="flex justify-between text-slate-400">
                      <span>教室：</span>
                      <span class="text-slate-200">{{ session.room_name }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 課程詳情 Modal -->
    <Teleport to="body">
      <ScheduleDetailPanel
        v-if="selectedSession"
        :schedule="selectedSession"
        @close="selectedSession = null"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </Teleport>

    <!-- 新增/編輯排課 Modal -->
    <Teleport to="body">
      <ScheduleRuleModal
        v-if="showCreateModal"
        @close="showCreateModal = false"
        @created="handleRuleCreated"
      />
      <ScheduleRuleModal
        v-if="showEditModal"
        :editing-rule="editingRule"
        @close="showEditModal = false; editingRule = null"
        @updated="handleRuleUpdated"
      />
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { formatDateToString } from '~/composables/useTaiwanTime'

const emit = defineEmits<{
  'update:viewMode': [value: 'all' | 'teacher' | 'room']
  'update:selectedResourceId': [value: number | null]
}>()

// Alert composable
const { confirm: confirmDialog, error: alertError } = useAlert()

// Props
const props = defineProps<{
  viewMode: 'all' | 'teacher' | 'room'
  selectedResourceId: number | null
}>()

// Computed with setter
const viewModeModel = computed({
  get: () => props.viewMode,
  set: (value) => emit('update:viewMode', value)
})

const selectedResourceIdModel = computed({
  get: () => props.selectedResourceId,
  set: (value) => emit('update:selectedResourceId', value)
})

const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingRule = ref<any>(null)
const selectedSession = ref<any>(null)
const { getCenterId } = useCenterId()

const handleEdit = () => {
  if (selectedSession.value) {
    editingRule.value = selectedSession.value
    showEditModal.value = true
  }
}

const handleDelete = async () => {
  if (!selectedSession.value || !(await confirmDialog('確定要刪除此排課規則？'))) return

  try {
    const api = useApi()
    await api.delete(`/admin/rules/${selectedSession.value.id}`)
    selectedSession.value = null
    await fetchSessions()
  } catch (err) {
    console.error('Failed to delete rule:', err)
    await alertError('刪除失敗，請稍後再試')
  }
}

const handleRuleUpdated = async () => {
  await fetchSessions()
  selectedSession.value = null
}

const timeSlots = [0, 1, 2, 3, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23]

const weekDays = [
  { value: 1, name: '週一' },
  { value: 2, name: '週二' },
  { value: 3, name: '週三' },
  { value: 4, name: '週四' },
  { value: 5, name: '週五' },
  { value: 6, name: '週六' },
  { value: 7, name: '週日' },
]

const sessions = ref<any[]>([])
const teachers = ref<any[]>([])
const rooms = ref<any[]>([])
const resourceCache = ref<{ teachers: Map<number, any>, rooms: Map<number, any> }>({
  teachers: new Map(),
  rooms: new Map(),
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

const resourceList = computed(() => {
  return viewModeModel.value === 'teacher'
    ? Array.from(resourceCache.value.teachers.values())
    : viewModeModel.value === 'room'
      ? Array.from(resourceCache.value.rooms.values())
      : []
})

const selectedResourceName = computed(() => {
  if (viewModeModel.value === 'teacher') {
    return resourceCache.value.teachers.get(selectedResourceIdModel.value)?.name || ''
  } else if (viewModeModel.value === 'room') {
    return resourceCache.value.rooms.get(selectedResourceIdModel.value)?.name || ''
  }
  return ''
})

const filteredSessions = computed(() => {
  let result = sessions.value

  // 過濾
  if (viewModeModel.value === 'teacher' && selectedResourceIdModel.value) {
    result = result.filter(session => session.teacher_id === selectedResourceIdModel.value)
  } else if (viewModeModel.value === 'room' && selectedResourceIdModel.value) {
    result = result.filter(session => session.room_id === selectedResourceIdModel.value)
  }

  // 計算衝突位置（用於全部視圖）
  if (viewModeModel.value === 'all') {
    const sessionGroups: Record<string, any[]> = {}
    result.forEach(session => {
      const key = `${session.weekday}-${session.start_time}`
      if (!sessionGroups[key]) sessionGroups[key] = []
      sessionGroups[key].push(session)
    })

    // 為每個衝突群組分配位置
    Object.values(sessionGroups).forEach(group => {
      if (group.length > 1) {
        // 有衝突，分配水平位置
        group.forEach((session, index) => {
          session.conflictIndex = index
          session.conflictCount = group.length
        })
      } else if (group.length === 1) {
        // 單一課程，佔滿寬度
        group[0].conflictIndex = 0
        group[0].conflictCount = 1
      }
    })
  } else {
    // 非全部視圖，不需要衝突處理
    result.forEach(session => {
      session.conflictIndex = 0
      session.conflictCount = 1
    })
  }

  return result
})

const formatTime = (hour: number): string => {
  return `${hour}:00`
}

const getDayDate = (weekday: number): string => {
  const date = new Date(weekStart.value)
  const diff = weekday - 1
  date.setDate(date.getDate() + diff)
  return date.getDate().toString()
}

const getSessionStyle = (session: any) => {
  // 計算位置和高度
  const [startHour, startMin] = session.start_time.split(':').map(Number)
  const [endHour, endMin] = session.end_time.split(':').map(Number)

  // 計算持續時間（分鐘）
  const duration = (endHour - startHour) * 60 + (endMin - startMin)

  // 計算 top 位置（從 0:00 開始）
  const startOffset = startHour * 60 + startMin

  // 計算星期位置（0-6）
  const weekdayIndex = session.weekday - 1

  // 每天寬度
  const dayWidth = `(100% - 60px) / 7`

  // 衝突處理：計算水平位置
  const conflictWidth = session.conflictCount > 1
    ? `calc((${dayWidth} - 8px) / ${session.conflictCount})`
    : `calc(${dayWidth} - 8px)`

  const conflictLeft = session.conflictCount > 1
    ? `calc(60px + ${weekdayIndex} * ${dayWidth} + ${session.conflictIndex} * (${dayWidth} - 8px) / ${session.conflictCount} + 4px)`
    : `calc(60px + ${weekdayIndex} * ${dayWidth} + 4px)`

  return {
    left: conflictLeft,
    width: conflictWidth,
    top: `${startOffset}px`,
    height: `${duration}px`,
    backgroundColor: getCourseColor(session.offering_name)
  }
}

const getCourseColor = (courseName: string): string => {
  const colors: Record<string, string> = {
    '瑜珈': 'bg-emerald-500/80',
    '瑜伽': 'bg-emerald-500/80',
    '有氧': 'bg-blue-500/80',
    '舞蹈': 'bg-purple-500/80',
    '鋼琴': 'bg-amber-500/80',
    '小提琴': 'bg-rose-500/80',
    '游泳': 'bg-cyan-500/80',
    '健身': 'bg-orange-500/80',
    '拳擊': 'bg-red-500/80',
    '芭蕾': 'bg-pink-500/80',
  }

  for (const [key, value] of Object.entries(colors)) {
    if (courseName.includes(key)) return value
  }
  return 'bg-primary-500/80'
}

const changeWeek = (delta: number) => {
  // 直接加減 7 天，不再強制對齊週一
  const currentDate = new Date(weekStart.value)
  currentDate.setDate(currentDate.getDate() + delta * 7)
  weekStart.value = currentDate
  fetchSessions()
}

const goToCurrentWeek = () => {
  // 回到今天，但不強制對齊週一
  weekStart.value = new Date()
  fetchSessions()
}

const clearFilter = () => {
  viewModeModel.value = 'all'
  selectedResourceIdModel.value = null
}

const selectSession = (session: any) => {
  selectedSession.value = session
}

const fetchSessions = async () => {
  try {
    const api = useApi()

    // 取得排課規則
    const rulesRes = await api.get<{ code: number; datas: any[] }>(
      '/admin/rules'
    )

    // 取得老師和教室
    const [teachersRes, roomsRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>('/teachers'),
      api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    ])

    teachers.value = teachersRes.datas || []
    rooms.value = roomsRes.datas || []

    // 更新快取
    teachersRes.datas?.forEach((t: any) => resourceCache.value.teachers.set(t.id, t))
    roomsRes.datas?.forEach((r: any) => resourceCache.value.rooms.set(r.id, r))

    // 展開規則為 sessions
    const sessionList: any[] = []
    const rules = rulesRes.datas || []

    rules.forEach((rule: any) => {
      // 後端返回的是 weekday（單一值），不是 weekdays 陣列
      const day = rule.weekday
      if (!day) return

      const startParts = rule.start_time.split(':')
      const endParts = rule.end_time.split(':')

      const startDate = new Date(weekStart.value)
      startDate.setDate(startDate.getDate() + (day - 1))

      sessionList.push({
        id: rule.id,
        rule_id: rule.id,
        date: formatDateToString(startDate),
        formatted_date: startDate.toLocaleDateString('zh-TW', {
          month: 'long',
          day: 'numeric',
          weekday: 'short'
        }),
        weekday: day,
        start_time: rule.start_time,
        end_time: rule.end_time,
        time_range: `${rule.start_time} - ${rule.end_time}`,
        offering_name: rule.offering?.name || '-',
        offering_id: rule.offering_id,
        teacher_id: rule.teacher_id,
        teacher_name: rule.teacher?.name || '-',
        room_id: rule.room_id,
        room_name: rule.room?.name || '-',
        color: getCourseColor(rule.offering?.name || '')
      })
    })

    sessions.value = sessionList
  } catch (error) {
    console.error('Failed to fetch sessions:', error)
    sessions.value = []
  }
}

const handleRuleCreated = () => {
  fetchSessions()
}

onMounted(() => {
  fetchSessions()
})
</script>
