<template>
  <div class="h-full flex flex-col glass-card overflow-hidden">
    <div class="p-4 border-b border-white/10">
      <div class="flex gap-2 mb-3">
        <button
          @click="activeTab = 'offering'"
          class="flex-1 glass-btn px-4 py-2 rounded-xl text-sm font-medium transition-all"
          :class="activeTab === 'offering' ? 'bg-primary-500/30 border-primary-500' : ''"
        >
          待排課程
        </button>
        <button
          @click="activeTab = 'teacher'"
          class="flex-1 glass-btn px-4 py-2 rounded-xl text-sm font-medium transition-all"
          :class="activeTab === 'teacher' ? 'bg-primary-500/30 border-primary-500' : ''"
        >
          老師列表
        </button>
        <button
          @click="activeTab = 'room'"
          class="flex-1 glass-btn px-4 py-2 rounded-xl text-sm font-medium transition-all"
          :class="activeTab === 'room' ? 'bg-primary-500/30 border-primary-500' : ''"
        >
          教室列表
        </button>
      </div>

      <div class="relative">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜尋..."
          class="input-field pl-10 text-sm"
        />
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto p-4 space-y-3">
      <div
        v-if="filteredItems.length === 0"
        class="text-center py-8 text-slate-500"
      >
        無{{ getTabName() }}資料
      </div>

      <div
        v-for="item in filteredItems"
        :key="item.id"
        class="glass p-3 rounded-xl cursor-pointer hover:bg-white/10 transition-all select-none"
        :class="isSelected(item) ? 'ring-2 ring-primary-500 bg-primary-500/10' : ''"
        @click="handleItemClick(item)"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <h4 class="font-medium text-slate-100 text-sm mb-1">
              {{ item.name || item.title }}
            </h4>
            <p v-if="item.subtitle" class="text-xs text-slate-400">
              {{ item.subtitle }}
            </p>
            <p v-if="item.skill_name && activeTab === 'teacher'" class="text-xs text-slate-500 mt-1">
              {{ item.skill_name }}
            </p>
          </div>
          <div v-if="isSelected(item)" class="flex-shrink-0">
            <svg class="w-5 h-5 text-primary-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
          </div>
        </div>
        <div v-if="item.tag" class="mt-2">
          <span
            class="px-2 py-1 rounded-full text-xs font-medium"
            :class="getTagClass(item.tag)"
          >
            {{ item.tag }}
          </span>
        </div>
      </div>
    </div>

    <!-- 清除篩選 -->
    <div
      v-if="selectedId"
      class="p-3 border-t border-white/10"
    >
      <button
        @click="clearSelection"
        class="w-full px-4 py-2 rounded-lg bg-slate-700/50 text-slate-300 hover:text-white hover:bg-slate-600 transition-colors text-sm"
      >
        顯示全部課表
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

// Props
const props = defineProps<{
  viewMode: 'all' | 'teacher' | 'room'
}>()

// Emit
const emit = defineEmits<{
  selectResource: { type: 'teacher' | 'room', id: number } | null
}>()

const activeTab = ref('offering')
const searchQuery = ref('')
const selectedId = ref<number | null>(null)
const { getCenterId } = useCenterId()

const offerings = ref<any[]>([])
const teachers = ref<any[]>([])
const rooms = ref<any[]>([])

// 監聽視角模式變化，自動切換到對應的 tab
watch(() => props.viewMode, (newMode) => {
  if (newMode === 'teacher') {
    activeTab.value = 'teacher'
  } else if (newMode === 'room') {
    activeTab.value = 'room'
  } else {
    activeTab.value = 'offering'
    clearSelection()
  }
})

const fetchData = async () => {
  try {
    const api = useApi()
    const centerId = getCenterId()

    const [offeringsRes, teachersRes, roomsRes] = await Promise.all([
      api.get<{ code: number; datas: any }>(`/admin/offerings`),
      api.get<{ code: number; datas: any[] }>('/teachers'),
      api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    ])

    offerings.value = offeringsRes.datas?.offerings || []
    teachers.value = teachersRes.datas || []
    rooms.value = roomsRes.datas || []
  } catch (error) {
    console.error('Failed to fetch data:', error)
    offerings.value = []
    teachers.value = []
    rooms.value = []
  }
}

const filteredItems = computed(() => {
  let items: any[] = []

  switch (activeTab.value) {
    case 'offering':
      items = offerings.value
      break
    case 'teacher':
      items = teachers.value
      break
    case 'room':
      items = rooms.value
      break
  }

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    items = items.filter(item =>
      (item.name || item.title)?.toLowerCase().includes(query) ||
      item.subtitle?.toLowerCase().includes(query) ||
      item.skill_name?.toLowerCase().includes(query)
    )
  }

  return items
})

// 檢查項目是否被選中
const isSelected = (item: any) => {
  return selectedId.value === item.id &&
    ((activeTab.value === 'teacher' && props.viewMode === 'teacher') ||
     (activeTab.value === 'room' && props.viewMode === 'room'))
}

// 處理項目點擊
const handleItemClick = (item: any) => {
  if (activeTab.value === 'offering') {
    // 課程不支援篩選視角，只處理拖放
    return
  }

  if (selectedId.value === item.id) {
    // 取消選中
    clearSelection()
  } else {
    // 選中該資源
    selectedId.value = item.id
    emit('selectResource', {
      type: activeTab.value as 'teacher' | 'room',
      id: item.id
    })
  }
}

const clearSelection = () => {
  selectedId.value = null
  emit('selectResource', null)
}

onMounted(() => {
  fetchData()
})

const getTabName = (): string => {
  switch (activeTab.value) {
    case 'offering':
      return '課程'
    case 'teacher':
      return '老師'
    case 'room':
      return '教室'
    default:
      return ''
  }
}

const getTagClass = (tag: string): string => {
  switch (tag) {
    case 'Piano':
    case 'Violin':
    case 'Theory':
      return 'bg-secondary-500/20 text-secondary-500'
    case 'ACTIVE':
      return 'bg-success-500/20 text-success-500'
    case 'INACTIVE':
      return 'bg-slate-500/20 text-slate-400'
    case 'AVAILABLE':
      return 'bg-success-500/20 text-success-500'
    case 'BUSY':
      return 'bg-warning-500/20 text-warning-500'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const handleDragStart = (item: any, event: DragEvent) => {
  const typeMap: Record<string, string> = {
    offering: 'offering',
    teacher: 'teacher',
    room: 'room',
  }
  event.dataTransfer?.setData('application/json', JSON.stringify({
    type: typeMap[activeTab.value] || activeTab.value,
    item: item,
  }))
  event.dataTransfer.effectAllowed = 'copy'
}
</script>
