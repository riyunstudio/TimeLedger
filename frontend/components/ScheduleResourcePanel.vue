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
        draggable="true"
        @dragstart="handleDragStart(item, $event)"
      >
        <h4 class="font-medium text-slate-100 text-sm mb-1">
          {{ item.name || item.title }}
        </h4>
        <p v-if="item.subtitle" class="text-xs text-slate-400">
          {{ item.subtitle }}
        </p>
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const activeTab = ref('offering')
const searchQuery = ref('')
const { getCenterId } = useCenterId()

const offerings = ref<any[]>([])
const teachers = ref<any[]>([])
const rooms = ref<any[]>([])

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
      item.subtitle?.toLowerCase().includes(query)
    )
  }

  return items
})

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
