<template>
  <div class="teacher-timeline bg-white/5 rounded-xl p-4">
    <!-- 標題列 -->
    <div class="flex items-center gap-3 mb-4">
      <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0">
        <span class="text-white font-medium">{{ teacher.teacher_name?.charAt(0) || '?' }}</span>
      </div>
      <div class="flex-1">
        <h4 class="text-white font-medium">{{ teacher.teacher_name }}</h4>
        <p class="text-xs text-slate-400">課表衝突檢視</p>
      </div>
      <BaseBadge
        :variant="availabilityVariant"
        size="sm"
      >
        {{ availabilityText }}
      </BaseBadge>
    </div>

    <!-- 時間軸容器 -->
    <div class="timeline-container relative pl-10">
      <!-- 時間刻度 -->
      <div class="time-axis absolute left-0 top-0 bottom-0 w-8 flex flex-col text-xs text-slate-500">
        <div v-for="hour in hours" :key="hour" class="h-12 flex items-center justify-end pr-1">
          {{ hour }}:00
        </div>
      </div>

      <!-- 軌道區域 -->
      <div class="tracks space-y-1">
        <!-- 目標時段軌道 -->
        <div class="track relative h-10 bg-indigo-500/20 rounded-lg border border-indigo-500/30">
          <div
            class="absolute h-full bg-indigo-500 rounded-lg flex items-center justify-center text-white text-xs font-medium"
            :style="targetSlotStyle"
          >
            目標時段
          </div>
        </div>

        <!-- 現有課程軌道 -->
        <div class="track relative h-10 bg-white/5 rounded-lg">
          <div
            v-for="session in existingSessions"
            :key="session.id"
            class="absolute h-8 top-1 bg-slate-600/50 rounded border-l-4 border-slate-500 overflow-hidden"
            :style="getSessionStyle(session)"
          >
            <span class="text-xs px-2 text-white truncate block">{{ session.course_name }}</span>
          </div>
        </div>

        <!-- 衝突標註 -->
        <div
          v-if="hasConflict"
          class="conflict-indicator flex items-center gap-2 mt-2 p-2 bg-red-500/10 border border-red-500/20 rounded-lg"
        >
          <svg class="w-4 h-4 text-red-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <span class="text-sm text-red-400">{{ conflictMessage }}</span>
          <button
            v-if="availability === 'BUFFER_CONFLICT'"
            @click="$emit('override')"
            class="ml-auto px-3 py-1 rounded-lg bg-yellow-500/20 text-yellow-400 hover:bg-yellow-500/30 text-xs transition-colors"
          >
            仍要安排
          </button>
        </div>

        <!-- 無衝突提示 -->
        <div
          v-else
          class="flex items-center gap-2 mt-2 p-2 bg-green-500/10 border border-green-500/20 rounded-lg"
        >
          <svg class="w-4 h-4 text-green-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="text-sm text-green-400">此時段完全可用，無課程衝突</span>
        </div>
      </div>
    </div>

    <!-- 衝突細節說明 -->
    <div v-if="conflictDetails.length > 0" class="mt-4 pt-4 border-t border-white/10">
      <p class="text-xs text-slate-400 mb-2">衝突課程：</p>
      <div class="space-y-2">
        <div
          v-for="(detail, index) in conflictDetails"
          :key="index"
          class="flex items-center gap-2 text-xs"
        >
          <span
            :class="[
              'w-2 h-2 rounded-full shrink-0',
              detail.type === 'OVERLAP' ? 'bg-red-500' : 'bg-yellow-500'
            ]"
          />
          <span class="text-slate-300">{{ detail.course_name }}</span>
          <span class="text-slate-500">{{ detail.start }} - {{ detail.end }}</span>
          <span class="text-slate-400">{{ detail.room }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Session {
  id: number
  course_name: string
  start_time: string
  end_time: string
  room_name?: string
}

interface Props {
  teacher: {
    teacher_id: number
    teacher_name: string
    availability: string
  }
  targetStart: string
  targetEnd: string
  existingSessions: Session[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  override: []
}>()

// 產生小時刻度（6:00 - 22:00）
const hours = computed(() => {
  const result = []
  for (let h = 6; h <= 22; h++) {
    result.push(h)
  }
  return result
})

// 可用性文字
const availabilityText = computed(() => {
  switch (props.teacher.availability) {
    case 'AVAILABLE':
      return '完全可用'
    case 'BUFFER_CONFLICT':
      return '緩衝衝突'
    case 'OVERLAP':
      return '時間重疊'
    default:
      return props.teacher.availability
  }
})

// 可用性標籤樣式
const availabilityVariant = computed(() => {
  switch (props.teacher.availability) {
    case 'AVAILABLE':
      return 'success'
    case 'BUFFER_CONFLICT':
      return 'warning'
    case 'OVERLAP':
      return 'error'
    default:
      return 'secondary'
  }
})

// 是否有衝突
const hasConflict = computed(() => {
  return props.teacher.availability !== 'AVAILABLE'
})

// 衝突訊息
const conflictMessage = computed(() => {
  if (props.teacher.availability === 'BUFFER_CONFLICT') {
    return '與現有課程緩衝時間衝突'
  }
  return '與現有課程時間重疊'
})

// 衝突細節
const conflictDetails = computed(() => {
  const details: Array<{
    type: string
    course_name: string
    start: string
    end: string
    room: string
  }> = []

  const targetStart = new Date(props.targetStart)
  const targetEnd = new Date(props.targetEnd)

  props.existingSessions.forEach(session => {
    const sessionStart = new Date(session.start_time)
    const sessionEnd = new Date(session.end_time)

    // 檢查是否重疊
    if (targetStart < sessionEnd && targetEnd > sessionStart) {
      // 判斷衝突類型
      const hasBufferConflict = targetStart < sessionEnd && targetEnd > sessionStart

      details.push({
        type: hasBufferConflict ? 'BUFFER_CONFLICT' : 'OVERLAP',
        course_name: session.course_name,
        start: formatTime(sessionStart),
        end: formatTime(sessionEnd),
        room: session.room_name || '未指定教室'
      })
    }
  })

  return details
})

// 格式化時間
const formatTime = (date: Date): string => {
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

// 計算時間位置百分比
const getTimePosition = (time: string): number => {
  const date = new Date(time)
  const hours = date.getHours() + date.getMinutes() / 60
  // 6:00 = 0%, 22:00 = 100%
  return ((hours - 6) / 16) * 100
}

// 目標時段樣式
const targetSlotStyle = computed(() => {
  const startPos = getTimePosition(props.targetStart)
  const endPos = getTimePosition(props.targetEnd)
  const width = endPos - startPos

  return {
    left: `${startPos}%`,
    width: `${width}%`
  }
})

// 課程區段樣式
const getSessionStyle = (session: Session) => {
  const startPos = getTimePosition(session.start_time)
  const endPos = getTimePosition(session.end_time)
  const width = endPos - startPos

  return {
    left: `${startPos}%`,
    width: `${width}%`
  }
}
</script>
