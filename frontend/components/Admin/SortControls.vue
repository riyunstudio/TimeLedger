<template>
  <div class="sort-controls flex flex-wrap items-center gap-4 mb-4">
    <span class="text-sm text-slate-400">排序方式：</span>
    
    <div class="flex items-center gap-2">
      <select
        v-model="localSortBy"
        class="px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-white text-sm focus:outline-none focus:border-indigo-500"
      >
        <option value="score">媒合分數</option>
        <option value="availability">可用性</option>
        <option value="rating">內部評分</option>
        <option value="skill-match">技能匹配度</option>
      </select>
      
      <button
        @click="toggleSortOrder"
        :class="[
          'p-1.5 rounded-lg border transition-colors',
          sortOrder === 'desc'
            ? 'bg-indigo-500/20 border-indigo-500/30 text-indigo-400'
            : 'bg-white/5 border-white/10 text-slate-400 hover:bg-white/10'
        ]"
        :title="sortOrder === 'desc' ? '高到低' : '低到高'"
      >
        <svg v-if="sortOrder === 'desc'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
        </svg>
        <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h9m5-4v12m0 0l-4-4m4 4l4-4" />
        </svg>
      </button>
    </div>
    
    <!-- 檢視模式切換 -->
    <div class="flex items-center gap-1 ml-auto border border-white/10 rounded-lg overflow-hidden">
      <button
        :class="[
          'p-1.5 transition-colors',
          viewMode === 'card' ? 'bg-indigo-500 text-white' : 'bg-white/5 text-slate-400 hover:bg-white/10'
        ]"
        @click="$emit('update:viewMode', 'card')"
        title="卡片檢視"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
        </svg>
      </button>
      <button
        :class="[
          'p-1.5 transition-colors',
          viewMode === 'list' ? 'bg-indigo-500 text-white' : 'bg-white/5 text-slate-400 hover:bg-white/10'
        ]"
        @click="$emit('update:viewMode', 'list')"
        title="列表檢視"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
        </svg>
      </button>
      <button
        :class="[
          'p-1.5 transition-colors',
          viewMode === 'compare' ? 'bg-indigo-500 text-white' : 'bg-white/5 text-slate-400 hover:bg-white/10'
        ]"
        @click="$emit('update:viewMode', 'compare')"
        title="比較檢視"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      </button>
    </div>
    
    <!-- 比較模式入口 -->
    <div v-if="viewMode !== 'compare' && selectedCount > 0" class="flex items-center gap-2">
      <span class="text-xs text-slate-400">
        已選 {{ selectedCount }} 位
      </span>
      <button
        @click="$emit('update:viewMode', 'compare')"
        class="px-3 py-1.5 rounded-lg bg-indigo-500/20 border border-indigo-500/30 text-indigo-400 hover:bg-indigo-500/30 transition-colors text-sm"
      >
        開始比較
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
type SortBy = 'score' | 'availability' | 'rating' | 'skill-match'
type SortOrder = 'asc' | 'desc'
type ViewMode = 'card' | 'list' | 'compare'

const props = defineProps<{
  sortBy: SortBy
  sortOrder: SortOrder
  viewMode: ViewMode
  selectedCount: number
}>()

const emit = defineEmits<{
  'update:sortBy': [value: SortBy]
  'update:sortOrder': [value: SortOrder]
  'update:viewMode': [value: ViewMode]
}>()

// 本地狀態（避免直接修改 prop）
const localSortBy = computed({
  get: () => props.sortBy,
  set: (value) => emit('update:sortBy', value)
})

// 切換排序方向
const toggleSortOrder = () => {
  emit('update:sortOrder', props.sortOrder === 'desc' ? 'asc' : 'desc')
}
</script>
