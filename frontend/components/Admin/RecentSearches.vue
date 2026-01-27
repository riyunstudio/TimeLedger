<template>
  <div v-if="recentSearches.length > 0" class="recent-searches mb-4">
    <div class="flex items-center justify-between mb-2">
      <p class="text-sm text-slate-400">近期搜尋</p>
      <button
        @click="clearRecentSearches"
        class="text-xs text-slate-500 hover:text-slate-300 transition-colors"
      >
        清除紀錄
      </button>
    </div>
    <div class="flex flex-wrap gap-2">
      <button
        v-for="search in recentSearches"
        :key="search.id"
        @click="loadSearch(search)"
        class="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-white/5 hover:bg-white/10 border border-white/10 hover:border-white/20 transition-all text-sm text-slate-300"
      >
        <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ formatSearchSummary(search) }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface RecentSearch {
  id: string
  start_time: string
  end_time: string
  room_ids: number[]
  skills: string[]
  created_at: number
}

const emit = defineEmits<{
  loadSearch: [search: RecentSearch]
  clearAll: []
}>()

// 從 localStorage 載入近期搜尋
const recentSearches = ref<RecentSearch[]>([])

const STORAGE_KEY = 'timeledger_recent_searches'
const MAX_STORES = 5

const loadFromStorage = () => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      recentSearches.value = JSON.parse(stored)
    }
  } catch (e) {
    console.error('Failed to load recent searches:', e)
  }
}

const saveToStorage = () => {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(recentSearches.value))
  } catch (e) {
    console.error('Failed to save recent searches:', e)
  }
}

// 格式化搜尋摘要
const formatSearchSummary = (search: RecentSearch): string => {
  const date = new Date(search.start_time)
  const dateStr = `${date.getMonth() + 1}/${date.getDate()}`
  const timeStr = `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
  
  const skillsText = search.skills.length > 0 
    ? search.skills.slice(0, 2).join(', ') + (search.skills.length > 2 ? '...' : '')
    : '無技能條件'
  
  return `${dateStr} ${timeStr} · ${skillsText}`
}

// 載入搜尋條件
const loadSearch = (search: RecentSearch) => {
  emit('loadSearch', search)
}

// 清除所有近期搜尋
const clearRecentSearches = () => {
  recentSearches.value = []
  saveToStorage()
  emit('clearAll')
}

// 新增搜尋到歷史記錄
const addSearch = (search: Omit<RecentSearch, 'id' | 'created_at'>) => {
  // 移除相同的搜尋條件
  recentSearches.value = recentSearches.value.filter(
    s => !(s.start_time === search.start_time && s.end_time === search.end_time)
  )
  
  // 新增到最前面
  recentSearches.value.unshift({
    ...search,
    id: Date.now().toString(),
    created_at: Date.now()
  })
  
  // 保留最多 MAX_STORES 筆
  if (recentSearches.value.length > MAX_STORES) {
    recentSearches.value = recentSearches.value.slice(0, MAX_STORES)
  }
  
  saveToStorage()
}

// 暴露方法給父組件
defineExpose({
  addSearch,
  loadFromStorage
})

onMounted(() => {
  loadFromStorage()
})
</script>
