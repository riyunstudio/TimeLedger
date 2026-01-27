<template>
  <div class="talent-filter-panel">
    <!-- 展開/收合按鈕 -->
    <button
      @click="expanded = !expanded"
      class="flex items-center gap-2 px-4 py-2 rounded-lg bg-white/5 hover:bg-white/10 text-slate-300 transition-colors text-sm"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
      </svg>
      進階篩選
      <svg 
        class="w-4 h-4 transition-transform" 
        :class="{ 'rotate-180': expanded }"
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
      <!-- 已套用篩選數量 -->
      <BaseBadge v-if="activeFilterCount > 0" variant="info" size="xs">
        {{ activeFilterCount }}
      </BaseBadge>
    </button>

    <!-- 篩選面板內容 -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-show="expanded"
        class="mt-4 p-4 bg-white/5 rounded-xl border border-white/10"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <!-- 城市選擇 -->
          <div>
            <label class="block text-slate-300 mb-2">城市</label>
            <select
              v-model="filters.city"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-indigo-500"
            >
              <option value="">全部城市</option>
              <option v-for="city in cities" :key="city" :value="city">{{ city }}</option>
            </select>
          </div>

          <!-- 區域選擇 -->
          <div>
            <label class="block text-slate-300 mb-2">區域</label>
            <select
              v-model="filters.district"
              :disabled="!filters.city"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-indigo-500 disabled:opacity-50"
            >
              <option value="">全部區域</option>
              <option v-for="district in districts" :key="district" :value="district">{{ district }}</option>
            </select>
          </div>

          <!-- 技能類別 -->
          <div>
            <label class="block text-slate-300 mb-2">技能類別</label>
            <select
              v-model="filters.skillCategory"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-indigo-500"
            >
              <option value="">全部類別</option>
              <option v-for="cat in skillCategories" :key="cat.key" :value="cat.key">
                {{ cat.icon }} {{ cat.name }}
              </option>
            </select>
          </div>

          <!-- 徵才狀態 -->
          <div>
            <label class="block text-slate-300 mb-2">徵才狀態</label>
            <div class="flex gap-1">
              <button
                v-for="status in hiringStatuses"
                :key="status.value"
                @click="filters.hiringStatus = status.value"
                :class="[
                  'flex-1 px-2 py-1.5 rounded-lg text-xs transition-colors',
                  filters.hiringStatus === status.value
                    ? 'bg-indigo-500 text-white'
                    : 'bg-white/5 text-slate-400 hover:bg-white/10'
                ]"
              >
                {{ status.label }}
              </button>
            </div>
          </div>
        </div>

        <!-- 第二行篩選 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4 pt-4 border-t border-white/10">
          <!-- 評分篩選 -->
          <div>
            <label class="block text-slate-300 mb-2">
              最低評分 
              <span class="text-indigo-400">{{ filters.minRating }}+</span>
            </label>
            <div class="flex items-center gap-3">
              <input
                type="range"
                v-model.number="filters.minRating"
                min="0"
                max="5"
                step="0.5"
                class="flex-1 h-2 bg-white/10 rounded-full appearance-none cursor-pointer"
              />
              <button
                v-if="filters.minRating > 0"
                @click="filters.minRating = 0"
                class="text-xs text-slate-400 hover:text-slate-300"
              >
                清除
              </button>
            </div>
          </div>

          <!-- 技能數量 -->
          <div>
            <label class="block text-slate-300 mb-2">
              最低技能數
              <span class="text-indigo-400">{{ filters.minSkills }}+</span>
            </label>
            <div class="flex items-center gap-3">
              <input
                type="range"
                v-model.number="filters.minSkills"
                min="0"
                max="10"
                step="1"
                class="flex-1 h-2 bg-white/10 rounded-full appearance-none cursor-pointer"
              />
              <button
                v-if="filters.minSkills > 0"
                @click="filters.minSkills = 0"
                class="text-xs text-slate-400 hover:text-slate-300"
              >
                清除
              </button>
            </div>
          </div>

          <!-- 是否為中心成員 -->
          <div>
            <label class="block text-slate-300 mb-2">會籍狀態</label>
            <div class="flex gap-2">
              <button
                @click="filters.membership = filters.membership === 'member' ? '' : 'member'"
                :class="[
                  'flex-1 px-3 py-1.5 rounded-lg text-sm transition-colors',
                  filters.membership === 'member'
                    ? 'bg-green-500/20 text-green-400 border border-green-500/30'
                    : 'bg-white/5 text-slate-400 hover:bg-white/10'
                ]"
              >
                中心成員
              </button>
              <button
                @click="filters.membership = filters.membership === 'external' ? '' : 'external'"
                :class="[
                  'flex-1 px-3 py-1.5 rounded-lg text-sm transition-colors',
                  filters.membership === 'external'
                    ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30'
                    : 'bg-white/5 text-slate-400 hover:bg-white/10'
                ]"
              >
                外部人才
              </button>
            </div>
          </div>
        </div>

        <!-- 篩選操作 -->
        <div class="flex items-center gap-3 mt-4 pt-4 border-t border-white/10">
          <button
            @click="applyFilters"
            class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors"
          >
            套用篩選
          </button>
          <button
            @click="clearFilters"
            class="px-4 py-2 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 transition-colors"
          >
            清除全部
          </button>
          <button
            @click="saveAsPreset"
            class="px-4 py-2 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 transition-colors flex items-center gap-1 ml-auto"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
            </svg>
            儲存為常用
          </button>
        </div>
      </div>
    </Transition>

    <!-- 已套用篩選標籤 -->
    <div v-if="activeFilters.length > 0" class="flex flex-wrap gap-2 mt-3">
      <BaseBadge
        v-for="filter in activeFilters"
        :key="`${filter.key}-${filter.value}`"
        variant="primary"
        size="sm"
        removable
        @remove="removeFilter(filter.key, filter.value)"
      >
        {{ filter.label }}
      </BaseBadge>
    </div>
  </div>
</template>

<script setup lang="ts">
import { SKILL_CATEGORIES } from '~/types'

interface FilterState {
  city: string
  district: string
  skillCategory: string
  hiringStatus: string
  minRating: number
  minSkills: number
  membership: string
}

interface FilterPreset {
  id: string
  name: string
  icon: string
  filters: FilterState
}

const emit = defineEmits<{
  apply: [filters: FilterState]
  clear: []
}>()

const expanded = ref(false)

// 篩選狀態
const filters = ref<FilterState>({
  city: '',
  district: '',
  skillCategory: '',
  hiringStatus: '',
  minRating: 0,
  minSkills: 0,
  membership: ''
})

// 選項資料
const cities = ref(['台北市', '新北市', '桃園市', '台中市', '台南市', '高雄市'])
const districts = ref(['中山區', '大安區', '信義區', '松山區', '萬華區', '中正區'])

const skillCategories = computed(() => {
  return Object.entries(SKILL_CATEGORIES).map(([key, cat]) => ({
    key,
    name: cat.name || key,
    icon: cat.icon || '✨'
  }))
})

const hiringStatuses = ref([
  { value: 'open', label: '開放中' },
  { value: 'closed', label: '已關閉' },
  { value: '', label: '全部' }
])

// 計算屬性：已套用的篩選數量
const activeFilterCount = computed(() => {
  let count = 0
  if (filters.value.city) count++
  if (filters.value.district) count++
  if (filters.value.skillCategory) count++
  if (filters.value.hiringStatus) count++
  if (filters.value.minRating > 0) count++
  if (filters.value.minSkills > 0) count++
  if (filters.value.membership) count++
  return count
})

// 計算屬性：已套用篩選的標籤
const activeFilters = computed(() => {
  const labels: Array<{ key: string; value: string; label: string }> = []
  
  if (filters.value.city) {
    labels.push({ key: 'city', value: filters.value.city, label: `城市: ${filters.value.city}` })
  }
  if (filters.value.district) {
    labels.push({ key: 'district', value: filters.value.district, label: `區域: ${filters.value.district}` })
  }
  if (filters.value.skillCategory) {
    const cat = skillCategories.value.find(c => c.key === filters.value.skillCategory)
    labels.push({ key: 'skillCategory', value: filters.value.skillCategory, label: `類別: ${cat?.icon} ${cat?.name}` })
  }
  if (filters.value.hiringStatus) {
    const status = hiringStatuses.value.find(s => s.value === filters.value.hiringStatus)
    labels.push({ key: 'hiringStatus', value: filters.value.hiringStatus, label: status?.label || '' })
  }
  if (filters.value.minRating > 0) {
    labels.push({ key: 'minRating', value: String(filters.value.minRating), label: `${filters.value.minRating}+ 評分` })
  }
  if (filters.value.minSkills > 0) {
    labels.push({ key: 'minSkills', value: String(filters.value.minSkills), label: `${filters.value.minSkills}+ 技能` })
  }
  if (filters.value.membership) {
    labels.push({ key: 'membership', value: filters.value.membership, label: filters.value.membership === 'member' ? '中心成員' : '外部人才' })
  }
  
  return labels
})

// 套用篩選
const applyFilters = () => {
  emit('apply', { ...filters.value })
}

// 清除篩選
const clearFilters = () => {
  filters.value = {
    city: '',
    district: '',
    skillCategory: '',
    hiringStatus: '',
    minRating: 0,
    minSkills: 0,
    membership: ''
  }
  emit('clear')
}

// 移除單一篩選
const removeFilter = (key: string, value: string) => {
  switch (key) {
    case 'city': filters.value.city = ''; break
    case 'district': filters.value.district = ''; break
    case 'skillCategory': filters.value.skillCategory = ''; break
    case 'hiringStatus': filters.value.hiringStatus = ''; break
    case 'minRating': filters.value.minRating = 0; break
    case 'minSkills': filters.value.minSkills = 0; break
    case 'membership': filters.value.membership = ''; break
  }
  emit('apply', { ...filters.value })
}

// 儲存為常用篩選
const saveAsPreset = () => {
  // TODO: 實現儲存預設功能
  console.log('Saving filter preset:', filters.value)
}

// 暴露方法
defineExpose({
  clearFilters,
  expanded
})
</script>

<style scoped>
input[type="range"] {
  -webkit-appearance: none;
  background: transparent;
}

input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 16px;
  height: 16px;
  background: linear-gradient(135deg, #6366f1, #a855f7);
  border-radius: 50%;
  cursor: pointer;
}

input[type="range"]::-webkit-slider-runnable-track {
  height: 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}
</style>
