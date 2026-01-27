<template>
  <div class="search-suggestions relative">
    <!-- 搜尋輸入框 -->
    <div class="relative">
      <input
        :value="query"
        @input="onInput"
        @focus="showSuggestions = true"
        @keydown.enter.prevent="executeSearch"
        @keydown.escape="showSuggestions = false"
        type="text"
        placeholder="搜尋人才姓名、技能、標籤..."
        class="w-full px-4 py-3 pl-10 rounded-xl bg-white/5 border border-white/10 text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
      />
      
      <!-- 搜尋圖標 -->
      <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      
      <!-- 清除按鈕 -->
      <button
        v-if="query"
        @click="clearQuery"
        class="absolute right-3 top-1/2 -translate-y-1/2 p-1 rounded-lg hover:bg-white/10 text-slate-400 transition-colors"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- 建議下拉選單 -->
    <Transition
      enter-active-class="transition-all duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="showSuggestions && (suggestions.length > 0 || recentSearches.length > 0 || trendingSearches.length > 0)"
        class="absolute z-50 w-full mt-2 bg-slate-800 border border-white/10 rounded-xl shadow-2xl overflow-hidden"
      >
        <!-- 熱門搜尋 -->
        <div v-if="trendingSearches.length > 0 && !query" class="p-3 border-b border-white/10">
          <div class="flex items-center gap-2 mb-2">
            <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
            </svg>
            <span class="text-xs text-slate-400">熱門搜尋</span>
          </div>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="(trend, index) in trendingSearches"
              :key="index"
              @click="selectSuggestion({ type: 'trending', value: trend })"
              class="px-3 py-1 rounded-lg bg-white/5 hover:bg-white/10 text-slate-300 text-sm transition-colors"
            >
              {{ trend }}
            </button>
          </div>
        </div>

        <!-- 搜尋建議 -->
        <div v-if="suggestions.length > 0" class="p-3">
          <p v-if="query" class="text-xs text-slate-500 mb-2">搜尋建議</p>
          <button
            v-for="suggestion in suggestions"
            :key="`${suggestion.type}-${suggestion.value}`"
            @click="selectSuggestion(suggestion)"
            class="w-full px-3 py-2 text-left flex items-center gap-3 hover:bg-white/5 rounded-lg transition-colors group"
          >
            <div class="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center group-hover:bg-indigo-500/20 transition-colors">
              <svg class="w-4 h-4 text-slate-400 group-hover:text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <div class="flex-1">
              <span class="text-white">{{ suggestion.value }}</span>
              <span
                :class="[
                  'ml-2 px-2 py-0.5 rounded text-xs',
                  suggestion.type === 'skill' ? 'bg-yellow-500/20 text-yellow-400' :
                  suggestion.type === 'tag' ? 'bg-primary-500/20 text-primary-400' :
                  'bg-blue-500/20 text-blue-400'
                ]"
              >
                {{ getTypeLabel(suggestion.type) }}
              </span>
            </div>
            <svg class="w-4 h-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>

        <!-- 近期搜尋 -->
        <div v-if="recentSearches.length > 0 && !query" class="p-3 border-t border-white/10">
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="text-xs text-slate-400">近期搜尋</span>
            </div>
            <button
              @click="clearRecentSearches"
              class="text-xs text-slate-500 hover:text-slate-300"
            >
              清除
            </button>
          </div>
          <button
            v-for="search in recentSearches"
            :key="search.id"
            @click="selectRecentSearch(search)"
            class="w-full px-3 py-2 text-left flex items-center gap-3 hover:bg-white/5 rounded-lg transition-colors"
          >
            <span class="text-slate-300">{{ search.query }}</span>
            <span class="text-xs text-slate-500 ml-auto">{{ formatTime(search.timestamp) }}</span>
          </button>
        </div>

        <!-- 快捷鍵提示 -->
        <div class="px-3 py-2 border-t border-white/10 bg-white/5">
          <div class="flex items-center gap-4 text-xs text-slate-500">
            <span class="flex items-center gap-1">
              <kbd class="px-1.5 py-0.5 rounded bg-white/10 text-slate-400">Enter</kbd>
              搜尋
            </span>
            <span class="flex items-center gap-1">
              <kbd class="px-1.5 py-0.5 rounded bg-white/10 text-slate-400">Esc</kbd>
              關閉
            </span>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 點擊外部關閉 -->
    <div
      v-if="showSuggestions"
      class="fixed inset-0 z-40"
      @click="showSuggestions = false"
    />
  </div>
</template>

<script setup lang="ts">
interface Suggestion {
  type: 'skill' | 'tag' | 'name'
  value: string
}

interface RecentSearch {
  id: string
  query: string
  timestamp: number
}

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  search: [query: string]
  selectSuggestion: [suggestion: Suggestion]
}>()

const query = ref(props.modelValue)
const showSuggestions = ref(false)
const suggestions = ref<Suggestion[]>([])
const recentSearches = ref<RecentSearch[]>([])
const trendingSearches = ref<string[]>(['瑜珈', '鋼琴', '舞蹈', '美術', '英語'])

// 監聽 modelValue 變化
watch(() => props.modelValue, (newVal) => {
  query.value = newVal
})

// 輸入事件
const onInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  query.value = target.value
  emit('update:modelValue', query.value)
  
  if (query.value.length > 0) {
    // 生成搜尋建議
    generateSuggestions()
  } else {
    suggestions.value = []
  }
}

// 生成搜尋建議
const generateSuggestions = () => {
  const q = query.value.toLowerCase()
  
  // 模擬建議（實際應該呼叫 API）
  const skillSuggestions: Suggestion[] = [
    { type: 'skill', value: '瑜珈' },
    { type: 'skill', value: '皮拉提斯' },
    { type: 'skill', value: '有氧舞蹈' },
    { type: 'skill', value: '鋼琴教學' },
    { type: 'skill', value: '小提琴' }
  ].filter(s => s.value.toLowerCase().includes(q))
  
  const tagSuggestions: Suggestion[] = [
    { type: 'tag', value: '古典' },
    { type: 'tag', value: '兒童' },
    { type: 'tag', value: '成人' },
    { type: 'tag', value: '進階' },
    { type: 'tag', value: '入門' }
  ].filter(s => s.value.toLowerCase().includes(q))
  
  suggestions.value = [...skillSuggestions, ...tagSuggestions].slice(0, 6)
}

// 選擇建議
const selectSuggestion = (suggestion: Suggestion) => {
  query.value = suggestion.value
  emit('update:modelValue', suggestion.value)
  emit('selectSuggestion', suggestion)
  showSuggestions.value = false
  
  // 加入近期搜尋
  addToRecentSearches(suggestion.value)
}

// 執行搜尋
const executeSearch = () => {
  if (query.value.trim()) {
    emit('search', query.value)
    showSuggestions.value = false
    addToRecentSearches(query.value)
  }
}

// 選擇近期搜尋
const selectRecentSearch = (search: RecentSearch) => {
  query.value = search.query
  emit('update:modelValue', search.query)
  showSuggestions.value = false
  emit('search', search.query)
}

// 加入近期搜尋
const addToRecentSearches = (searchQuery: string) => {
  const newSearch: RecentSearch = {
    id: Date.now().toString(),
    query: searchQuery,
    timestamp: Date.now()
  }
  
  // 移除相同的搜尋
  recentSearches.value = recentSearches.value.filter(s => s.query !== searchQuery)
  
  // 新增到最前面
  recentSearches.value.unshift(newSearch)
  
  // 保留最多 10 筆
  if (recentSearches.value.length > 10) {
    recentSearches.value = recentSearches.value.slice(0, 10)
  }
  
  // 儲存到 localStorage
  saveRecentSearches()
}

// 清除近期搜尋
const clearRecentSearches = () => {
  recentSearches.value = []
  localStorage.removeItem('talent_recent_searches')
}

// 清除查詢
const clearQuery = () => {
  query.value = ''
  emit('update:modelValue', '')
  suggestions.value = []
}

// 取得類型標籤
const getTypeLabel = (type: string): string => {
  switch (type) {
    case 'skill': return '技能'
    case 'tag': return '標籤'
    case 'name': return '姓名'
    default: return ''
  }
}

// 格式化時間
const formatTime = (timestamp: number): string => {
  const diff = Date.now() - timestamp
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (minutes < 1) return '剛剛'
  if (minutes < 60) return `${minutes} 分鐘前`
  if (hours < 24) return `${hours} 小時前`
  return `${days} 天前`
}

// 儲存近期搜尋到 localStorage
const saveRecentSearches = () => {
  try {
    localStorage.setItem('talent_recent_searches', JSON.stringify(recentSearches.value))
  } catch (e) {
    console.error('Failed to save recent searches:', e)
  }
}

// 從 localStorage 載入近期搜尋
const loadRecentSearches = () => {
  try {
    const stored = localStorage.getItem('talent_recent_searches')
    if (stored) {
      recentSearches.value = JSON.parse(stored)
    }
  } catch (e) {
    console.error('Failed to load recent searches:', e)
  }
}

onMounted(() => {
  loadRecentSearches()
})
</script>
