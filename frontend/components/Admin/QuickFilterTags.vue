<template>
  <div class="quick-filter-tags">
    <div class="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-hide">
      <span class="text-sm text-slate-400 shrink-0">快速篩選：</span>
      
      <!-- 熱門技能標籤 -->
      <button
        v-for="skill in popularSkills"
        :key="skill"
        @click="$emit('filterBySkill', skill)"
        :class="[
          'px-3 py-1.5 rounded-full text-sm whitespace-nowrap transition-all',
          activeSkills.includes(skill)
            ? 'bg-indigo-500 text-white shadow-lg shadow-indigo-500/25'
            : 'bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white'
        ]"
      >
        <svg class="w-3.5 h-3.5 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
        </svg>
        {{ skill }}
      </button>

      <!-- 分隔線 -->
      <div class="w-px h-6 bg-white/10 shrink-0" />

      <!-- 熱門標籤 -->
      <button
        v-for="tag in popularTags"
        :key="tag"
        @click="$emit('filterByTag', tag)"
        :class="[
          'px-3 py-1.5 rounded-full text-sm whitespace-nowrap transition-all',
          activeTags.includes(tag)
            ? 'bg-primary-500 text-white shadow-lg shadow-primary-500/25'
            : 'bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white'
        ]"
      >
        <span class="text-primary-400">#</span>{{ tag }}
      </button>

      <!-- 評分標籤 -->
      <div class="flex items-center gap-1 ml-2">
        <button
          v-for="rating in ratingOptions"
          :key="rating"
          @click="$emit('filterByRating', rating)"
          :class="[
            'px-2 py-1 rounded-full text-xs whitespace-nowrap transition-all flex items-center gap-1',
            activeRating === rating
              ? 'bg-yellow-500 text-white shadow-lg shadow-yellow-500/25'
              : 'bg-white/5 text-slate-400 hover:bg-white/10'
          ]"
        >
          <svg class="w-3 h-3 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
          </svg>
          {{ rating }}+
        </button>
      </div>

      <!-- 清除篩選 -->
      <button
        v-if="hasActiveFilters"
        @click="$emit('clearAll')"
        class="px-3 py-1.5 rounded-full text-sm whitespace-nowrap bg-red-500/20 text-red-400 hover:bg-red-500/30 transition-all ml-auto shrink-0"
      >
        <svg class="w-3.5 h-3.5 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
        清除篩選
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  popularSkills: string[]
  popularTags: string[]
  activeSkills: string[]
  activeTags: string[]
  activeRating: number
}>()

defineEmits<{
  filterBySkill: [skill: string]
  filterByTag: [tag: string]
  filterByRating: [rating: number]
  clearAll: []
}>()

const ratingOptions = [3, 4, 5]

// 計算屬性：是否有啟用的篩選
const hasActiveFilters = computed(() => {
  return props.activeSkills.length > 0 || 
         props.activeTags.length > 0 || 
         props.activeRating > 0
})
</script>

<style scoped>
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
</style>
