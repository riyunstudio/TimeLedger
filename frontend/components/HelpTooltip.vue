<template>
  <div ref="wrapperRef" class="help-tooltip-wrapper">
    <!-- 問號按鈕 -->
    <button
      type="button"
      class="help-tooltip-btn inline-flex items-center justify-center w-5 h-5 rounded-full bg-slate-700/50 text-slate-400 hover:text-white hover:bg-slate-600 transition-colors cursor-help shrink-0"
      :class="buttonClass"
      @click="toggle"
      @mouseenter="show"
      @mouseleave="hide"
    >
      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    </button>

    <!-- 提示框 - 使用 Teleport 移到 body -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-150"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition ease-in duration-100"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div
          v-if="isVisible"
          class="help-tooltip-overlay"
          @click="hide"
        >
          <div
            ref="tooltipRef"
            class="help-tooltip-content"
            :style="tooltipStyle"
            @click.stop
          >
            <!-- 標題 -->
            <div v-if="title" class="font-medium text-white mb-2 flex items-center gap-2">
              <svg class="w-4 h-4 text-primary-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{{ title }}</span>
            </div>

            <!-- 說明內容 -->
            <div class="text-sm text-slate-300 space-y-2">
              <p v-if="description" class="leading-relaxed">{{ description }}</p>

              <!-- 使用方式 -->
              <div v-if="usage && usage.length > 0" class="mt-3 pt-3 border-t border-white/10">
                <div class="text-xs text-slate-500 mb-2 uppercase tracking-wide font-medium">使用方式</div>
                <ul class="space-y-1.5">
                  <li v-for="(item, index) in usage" :key="index" class="flex items-start gap-2">
                    <svg class="w-4 h-4 text-primary-400 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span class="leading-relaxed">{{ item }}</span>
                  </li>
                </ul>
              </div>

              <!-- 快捷鍵 -->
              <div v-if="shortcut" class="mt-3 pt-3 border-t border-white/10">
                <div class="text-xs text-slate-500 mb-2 uppercase tracking-wide font-medium">快捷鍵</div>
                <div class="flex items-center gap-2">
                  <kbd class="px-2 py-1 bg-slate-700/70 rounded text-xs text-slate-200 font-mono">{{ shortcut }}</kbd>
                </div>
              </div>
            </div>

            <!-- 箭頭 -->
            <div class="help-tooltip-arrow" :style="arrowStyle"></div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title?: string
  description?: string
  usage?: string[]
  shortcut?: string
  placement?: 'top' | 'bottom' | 'left' | 'right'
  buttonClass?: string
}

const props = withDefaults(defineProps<Props>(), {
  placement: 'top',
})

const isVisible = ref(false)
let hideTimer: ReturnType<typeof setTimeout> | null = null
const wrapperRef = ref<HTMLElement | null>(null)

const show = () => {
  if (hideTimer) {
    clearTimeout(hideTimer)
    hideTimer = null
  }
  isVisible.value = true
}

const hide = () => {
  hideTimer = setTimeout(() => {
    isVisible.value = false
  }, 200)
}

const toggle = () => {
  isVisible.value = !isVisible.value
}

// 檢查某個位置是否適合顯示提示框
const checkFits = (rect: DOMRect, placement: string, tooltipWidth: number, tooltipHeight: number): boolean => {
  const padding = 16
  const arrowSize = 12
  const gap = 8

  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight

  switch (placement) {
    case 'top':
      return rect.top - gap - arrowSize - tooltipHeight >= padding
    case 'bottom':
      return rect.bottom + gap + arrowSize + tooltipHeight <= viewportHeight - padding
    case 'left':
      return rect.left - gap - arrowSize - tooltipWidth >= padding
    case 'right':
      return rect.right + gap + arrowSize + tooltipWidth <= viewportWidth - padding
    default:
      return false
  }
}

// 計算最佳位置
const getBestPlacement = (rect: DOMRect, tooltipWidth: number, tooltipHeight: number): string => {
  const placements = ['top', 'bottom', 'left', 'right']
  const preferredOrder = [props.placement, ...placements.filter(p => p !== props.placement)]

  for (const placement of preferredOrder) {
    if (checkFits(rect, placement, tooltipWidth, tooltipHeight)) {
      return placement
    }
  }

  // 如果首選方向都不行，嘗試找任何一個可以顯示的方向
  for (const placement of placements) {
    if (checkFits(rect, placement, tooltipWidth, tooltipHeight)) {
      return placement
    }
  }

  // 如果都放不下，就用首選方向（會稍微超出邊界）
  return props.placement
}

const tooltipStyle = computed(() => {
  if (!wrapperRef.value) return {}

  const rect = wrapperRef.value.getBoundingClientRect()
  const tooltipWidth = 288 // min-w-72
  const tooltipHeight = 200 // 估計高度
  const arrowSize = 12
  const gap = 8

  // 自動選擇最佳位置
  const bestPlacement = getBestPlacement(rect, tooltipWidth, tooltipHeight)

  let top = 0
  let left = 0
  let transform = 'translateX(-50%)'

  switch (bestPlacement) {
    case 'top':
      top = rect.top - gap - arrowSize
      left = rect.left + rect.width / 2
      transform = 'translateX(-50%) translateY(-100%)'
      break
    case 'bottom':
      top = rect.bottom + gap + arrowSize
      left = rect.left + rect.width / 2
      transform = 'translateX(-50%)'
      break
    case 'left':
      top = rect.top + rect.height / 2
      left = rect.left - gap - arrowSize
      transform = 'translateX(-100%) translateY(-50%)'
      break
    case 'right':
      top = rect.top + rect.height / 2
      left = rect.right + gap + arrowSize
      transform = 'translateY(-50%)'
      break
  }

  // 邊界微調（當所有方向都放不下時使用）
  const padding = 16
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight

  // 確保不超出水平邊界
  if (left - tooltipWidth / 2 < padding) {
    left = tooltipWidth / 2 + padding
    transform = transform.replace('translateX(-50%)', '')
  } else if (left + tooltipWidth / 2 > viewportWidth - padding) {
    left = viewportWidth - tooltipWidth / 2 - padding
    transform = transform.replace('translateX(-50%)', '')
  }

  // 確保不超出垂直邊界
  if (top < padding) {
    top = padding
  } else if (top + tooltipHeight > viewportHeight - padding) {
    top = viewportHeight - tooltipHeight - padding
  }

  return {
    top: `${top}px`,
    left: `${left}px`,
    transform,
  }
})

const arrowStyle = computed(() => {
  // 根據當前位置計算箭頭位置（基於 tooltipStyle 的邏輯）
  const rect = wrapperRef.value?.getBoundingClientRect()
  if (!rect) return {}

  const tooltipWidth = 288
  const tooltipHeight = 200
  const arrowSize = 12
  const gap = 8
  const padding = 16

  const bestPlacement = getBestPlacement(rect, tooltipWidth, tooltipHeight)
  const arrowOffset = 8

  switch (bestPlacement) {
    case 'top':
      return {
        bottom: `${-arrowOffset}px`,
        left: '50%',
        transform: 'translateX(-50%) rotate(45deg)',
      }
    case 'bottom':
      return {
        top: `${-arrowOffset}px`,
        left: '50%',
        transform: 'translateX(-50%) rotate(45deg)',
      }
    case 'left':
      return {
        top: '50%',
        right: `${-arrowOffset}px`,
        transform: 'translateY(-50%) rotate(45deg)',
      }
    case 'right':
      return {
        top: '50%',
        left: `${-arrowOffset}px`,
        transform: 'translateY(-50%) rotate(45deg)',
      }
    default:
      return {
        bottom: `${-arrowOffset}px`,
        left: '50%',
        transform: 'translateX(-50%) rotate(45deg)',
      }
  }
})
</script>

<style scoped>
.help-tooltip-wrapper {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.help-tooltip-btn {
  flex-shrink: 0;
}

.help-tooltip-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.help-tooltip-content {
  position: fixed;
  min-width: 280px;
  max-width: 320px;
  padding: 16px;
  border-radius: 12px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.1);
  z-index: 51;
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.98) 0%, rgba(30, 41, 59, 0.98) 100%);
  backdrop-filter: blur(16px);
  pointer-events: auto;
}

.help-tooltip-arrow {
  position: absolute;
  width: 10px;
  height: 10px;
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.98) 0%, rgba(30, 41, 59, 0.98) 100%);
  border: 1px solid rgba(255, 255, 255, 0.1);
  z-index: 52;
}
</style>
