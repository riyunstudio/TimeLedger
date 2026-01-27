<template>
  <div
    ref="containerRef"
    class="virtual-scroll-container"
    :style="containerStyle"
    role="listbox"
    :aria-label="ariaLabel || '虛擬滾動列表'"
    @scroll="onScroll"
  >
    <div
      class="virtual-scroll-spacer"
      :style="{ height: `${totalHeight}px` }"
      role="presentation"
    ></div>
    <div
      class="virtual-scroll-content"
      :style="{ transform: `translateY(${offsetY}px)` }"
      role="presentation"
    >
      <div
        v-for="item in visibleItems"
        :key="item.key"
        class="virtual-scroll-item"
        :style="{ height: `${itemHeight}px` }"
        role="option"
        :aria-selected="selectedIndex === item.index"
      >
        <slot
          name="item"
          :item="item.data"
          :index="item.index"
          :style="getItemStyle(item.data, item.index)"
        >
          <!-- 預設插槽內容 -->
          <div :style="getItemStyle(item.data, item.index)">
            {{ getItemContent(item.data, item.index) }}
          </div>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, type PropType } from 'vue'

interface VirtualScrollItem {
  key: string | number
  data: any
  index: number
}

interface Props {
  items: any[]
  itemHeight: number
  itemKey?: string
  overscan?: number
  buffer?: number
  ariaLabel?: string
  selectedIndex?: number
}

const props = withDefaults(defineProps<Props>(), {
  itemKey: 'id',
  overscan: 3,
  buffer: 5,
  selectedIndex: -1,
})

const emit = defineEmits<{
  scroll: [event: Event]
  reachTop: []
  reachBottom: []
}>()

const containerRef = ref<HTMLElement | null>(null)
const scrollTop = ref(0)
const containerHeight = ref(0)

// 計算總高度
const totalHeight = computed(() => props.items.length * props.itemHeight)

// 計算可見範圍
const visibleRange = computed(() => {
  const start = Math.max(0, Math.floor(scrollTop.value / props.itemHeight) - props.overscan)
  const end = Math.min(
    props.items.length,
    Math.ceil((scrollTop.value + containerHeight.value) / props.itemHeight) + props.overscan
  )
  return { start, end }
})

// 計算偏移量
const offsetY = computed(() => visibleRange.value.start * props.itemHeight)

// 可見項目
const visibleItems = computed(() => {
  const items: VirtualScrollItem[] = []
  const { start, end } = visibleRange.value

  for (let i = start; i < end; i++) {
    const item = props.items[i]
    if (item !== undefined) {
      items.push({
        key: typeof item === 'object' ? item[props.itemKey] : item,
        data: item,
        index: i,
      })
    }
  }

  return items
})

// 容器樣式
const containerStyle = computed(() => ({
  height: `${containerHeight.value}px`,
  overflowY: 'auto' as const,
}))

// 滾動事件處理
const onScroll = (event: Event) => {
  const target = event.target as HTMLElement
  scrollTop.value = target.scrollTop

  // 檢查是否到達頂部或底部
  if (scrollTop.value <= 0) {
    emit('reachTop')
  }
  if (scrollTop.value + containerHeight.value >= totalHeight.value - props.buffer * props.itemHeight) {
    emit('reachBottom')
  }

  emit('scroll', event)
}

// 監聽容器尺寸變化
let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  if (containerRef.value) {
    containerHeight.value = containerRef.value.clientHeight

    resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        containerHeight.value = entry.contentRect.height
      }
    })
    resizeObserver.observe(containerRef.value)
  }
})

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
})

// 可覆寫的函數
const getItemStyle = (_data: any, _index: number) => ({})

const getItemContent = (data: any, index: number) => {
  if (typeof data === 'object' && data !== null) {
    return JSON.stringify(data)
  }
  return data
}

// 暴露方法給父組件
defineExpose({
  scrollToIndex: (index: number) => {
    if (containerRef.value) {
      containerRef.value.scrollTop = index * props.itemHeight
    }
  },
  scrollToTop: () => {
    if (containerRef.value) {
      containerRef.value.scrollTop = 0
    }
  },
  scrollToBottom: () => {
    if (containerRef.value) {
      containerRef.value.scrollTop = totalHeight.value - containerHeight.value
    }
  },
  refresh: () => {
    // 強制重新計算
    scrollTop.value = containerRef.value?.scrollTop || 0
  },
})
</script>

<style scoped>
.virtual-scroll-container {
  position: relative;
  width: 100%;
  contain: strict;
}

.virtual-scroll-spacer {
  width: 1px;
  will-change: transform;
}

.virtual-scroll-content {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  will-change: transform;
}

.virtual-scroll-item {
  contain: content;
  will-change: transform, height;
}
</style>
