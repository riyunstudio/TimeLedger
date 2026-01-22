<template>
  <div class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-2xl max-h-[90vh] overflow-hidden animate-spring flex flex-col" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          課表預覽
        </h3>
        <div class="flex gap-2">
          <button
            @click="router.push('/teacher/export')"
            class="px-4 py-2 rounded-lg glass-btn flex items-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2" />
            </svg>
            全版面
          </button>
          <button
            @click="handleShare"
            class="px-4 py-2 rounded-lg glass-btn flex items-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 12v8a2 2 0 002 2h12a2 2 0 002-2v-8" />
              <polyline points="16 6 12 2 8 6" />
              <line x1="12" y1="2" x2="12" y2="15" />
            </svg>
            分享
          </button>
          <button
            @click="handleDownload"
            class="px-4 py-2 rounded-lg glass-btn flex items-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            下載
          </button>
          <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>

      <div class="flex-1 overflow-auto p-6 flex justify-center">
        <div
          ref="scheduleRef"
          class="rounded-3xl shadow-2xl"
          :style="[scheduleStyle, containerStyle]"
        >
          <div
            v-if="options.showPersonalInfo"
            class="text-center pb-4 border-b border-white/20"
            :class="options.compactMode ? 'mb-3' : 'mb-6'"
          >
            <div
              class="rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center mx-auto mb-2"
              :class="options.compactMode ? 'w-12 h-12' : 'w-20 h-20'"
            >
              <span :class="options.compactMode ? 'text-lg font-bold' : 'text-2xl font-bold'" class="text-white">
                {{ authStore.user?.name?.charAt(0) || 'T' }}
              </span>
            </div>
            <h2 :class="options.compactMode ? 'text-base font-bold' : 'text-xl font-bold'" class="text-white mb-2">
              {{ authStore.user?.name }}
            </h2>
            <div class="flex flex-wrap gap-1 justify-center">
              <span
                v-for="i in 3"
                :key="i"
                class="px-2 py-0.5 rounded-full text-xs bg-white/20 text-white"
              >
                #{{ ['鋼琴', '古典', '樂理'][i - 1] }}
              </span>
            </div>
          </div>

          <div v-if="view === 'grid'" class="grid gap-2" :class="gridColsClass">
            <div
              v-for="day in previewDays"
              :key="day.date"
              class="rounded-xl p-2 bg-white/10"
            >
              <div class="text-center mb-1.5">
                <h4 :class="options.compactMode ? 'text-xs' : 'text-sm'" class="font-semibold text-white">
                  {{ formatDateShort(day.date) }}
                </h4>
                <span class="text-white/60 text-xs">
                  {{ day.items.length }} 課
                </span>
              </div>
              <div class="space-y-1">
                <div
                  v-for="item in getVisibleItems(day.items)"
                  :key="item.id"
                  class="rounded-lg p-1.5 bg-white/20"
                  :class="options.privacyMode && item.type === 'PERSONAL_EVENT' ? 'bg-white/10' : 'bg-white/20'"
                >
                  <div v-if="options.showTime" class="text-white/90 text-xs mb-0.5">
                    {{ item.start_time }}
                  </div>
                  <h5 :class="options.compactMode ? 'text-xs' : 'text-sm'" class="font-medium text-white truncate">
                    {{ item.type === 'PERSONAL_EVENT' ? (options.privacyMode ? '已保留' : item.title) : item.title }}
                  </h5>
                </div>
                <div
                  v-if="day.items.length > maxVisibleItems"
                  class="text-center text-white/50 text-xs py-1"
                >
                  +{{ day.items.length - maxVisibleItems }} 更多
                </div>
                <div
                  v-if="day.items.length === 0"
                  class="text-center py-2 text-white/40 text-xs"
                >
                  休息
                </div>
              </div>
            </div>
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="day in previewDays"
              :key="day.date"
            >
              <div class="flex items-center justify-between mb-1">
                <h4 :class="options.compactMode ? 'text-sm' : 'text-base'" class="font-semibold text-white">
                  {{ formatDate(day.date) }}
                </h4>
                <span class="text-white/70 text-xs">
                  {{ day.items.length }} 個課程
                </span>
              </div>

              <div class="space-y-1.5">
                <div
                  v-for="item in day.items"
                  :key="item.id"
                  class="rounded-lg p-2"
                  :class="options.privacyMode && item.type === 'PERSONAL_EVENT' ? 'bg-white/10' : 'bg-white/20'"
                >
                  <div v-if="options.showTime" class="text-white/90 text-xs mb-0.5">
                    {{ item.start_time }} - {{ item.end_time }}
                  </div>
                  <h5 :class="options.compactMode ? 'text-xs' : 'text-sm'" class="font-medium text-white">
                    {{ item.type === 'PERSONAL_EVENT' ? (options.privacyMode ? '已保留' : item.title) : item.title }}
                  </h5>
                </div>
              </div>

              <div
                v-if="day.items.length === 0"
                class="text-center py-2 text-white/50 text-xs"
              >
                無行程
              </div>
            </div>
          </div>

          <div class="text-center pt-3 border-t border-white/20">
            <p class="text-white/60 text-xs">TimeLedger</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const router = useRouter()
const props = defineProps<{
  theme: any
  options: any
  format: string
  view?: 'grid' | 'list'
}>()

const emit = defineEmits<{
  close: []
}>()

const authStore = useAuthStore()
const scheduleRef = ref<HTMLElement>()
const teacherStore = useTeacherStore()

const maxVisibleItems = computed(() => {
  if (props.options.compactMode) {
    return props.view === 'grid' ? 4 : 8
  }
  return props.view === 'grid' ? 3 : 5
})

const gridColsClass = computed(() => {
  const days = props.options.showFullWeek ? 7 : 3
  if (days === 7) return 'grid-cols-7'
  return 'grid-cols-3'
})

const containerStyle = computed(() => {
  const padding = props.options.compactMode ? 'p-4' : 'p-6'
  return {
    padding,
  }
})

const scheduleStyle = computed(() => ({
  background: props.theme?.preview || 'linear-gradient(135deg, #1e3a8a 0%, #6366F1 25%, #A855F7 50%, #1e3a8a 75%, #0f172a 100%)',
}))

const getVisibleItems = (items: any[]) => {
  const max = maxVisibleItems.value
  if (items.length <= max) return items
  return items.slice(0, max)
}

const previewDays = computed(() => {
  const days = props.options.showFullWeek ? 7 : 3
  return teacherStore.schedule?.days?.slice(0, days) || []
})

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const diffDays = Math.floor((date.getTime() - today.getTime()) / (1000 * 60 * 60 * 24))

  if (diffDays === 0) return '今天'
  if (diffDays === 1) return '明天'

  return date.toLocaleDateString('zh-TW', {
    month: 'long',
    day: 'numeric',
    weekday: 'short',
  })
}

const formatDateShort = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    month: 'numeric',
    day: 'numeric',
    weekday: 'short',
  })
}

const handleDownload = () => {
  if (scheduleRef.value) {
    import('html2canvas').then(({ default: html2canvas }) => {
      html2canvas(scheduleRef.value!, {
        backgroundColor: null,
        scale: 2,
      }).then(canvas => {
        const link = document.createElement('a')
        link.download = `timeledger-schedule-${Date.now()}.png`
        link.href = canvas.toDataURL('image/png')
        link.click()
      })
    })
  }
}

const handleShare = async () => {
  if (scheduleRef.value) {
    await handleDownload()

    if (navigator.share) {
      try {
        await navigator.share({
          title: '我的課表 - TimeLedger',
          text: '查看我的課表！',
        })
      } catch (err) {
        console.log('Share failed:', err)
      }
    }
  }
}
</script>
