<template>
  <div class="skill-selector">
    <label class="block text-slate-300 mb-2">所需技能</label>
    
    <!-- 輸入框 -->
    <div class="relative">
      <input
        :value="inputValue"
        type="text"
        placeholder="輸入技能名稱（如：瑜珈、鋼琴）"
        class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 transition-colors"
        @input="onInput"
        @keydown.enter.prevent="addCurrentInput"
        @keydown.backspace="onBackspace"
        @focus="showSuggestions = true"
      />
      
      <!-- 建議技能下拉選單 -->
      <div
        v-if="showSuggestions && filteredSuggestions.length > 0"
        class="suggestions-dropdown absolute z-50 w-full mt-1 bg-slate-800 border border-white/10 rounded-lg shadow-xl overflow-hidden"
      >
        <button
          v-for="(skill, index) in filteredSuggestions"
          :key="skill.id"
          :class="[
            'w-full px-3 py-2 text-left flex items-center justify-between hover:bg-white/5 transition-colors',
            focusedIndex === index ? 'bg-white/10' : ''
          ]"
          @click="selectSkill(skill)"
          @mouseenter="focusedIndex = index"
        >
          <div class="flex items-center gap-2">
            <span class="font-medium text-white">{{ skill.name }}</span>
            <BaseBadge variant="secondary" size="xs">
              {{ skill.category }}
            </BaseBadge>
          </div>
          <span class="text-xs text-slate-400">
            {{ skill.teacherCount }} 位老師
          </span>
        </button>
      </div>
    </div>
    
    <!-- 已選技能標籤 -->
    <div v-if="selectedSkills.length > 0" class="flex flex-wrap gap-2 mt-3">
      <div
        v-for="(skill, index) in selectedSkills"
        :key="skill.id"
        class="group flex items-center gap-1 px-2 py-1 rounded-lg bg-indigo-500/20 border border-indigo-500/30 text-indigo-300 text-sm"
      >
        <span>{{ skill.name }}</span>
        <button
          @click="removeSkill(index)"
          class="ml-1 p-0.5 rounded hover:bg-indigo-500/30 transition-colors"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
    
    <!-- 熱門技能快捷選取 -->
    <div v-if="!showSuggestions && popularSkills.length > 0" class="mt-3">
      <p class="text-xs text-slate-500 mb-2">熱門技能</p>
      <div class="flex flex-wrap gap-1.5">
        <button
          v-for="skill in popularSkills"
          :key="skill.id"
          @click="addSkill(skill)"
          :disabled="isSkillSelected(skill.id)"
          :class="[
            'px-2 py-1 rounded-md text-xs transition-all',
            isSkillSelected(skill.id)
              ? 'bg-indigo-500/20 text-indigo-400 cursor-not-allowed'
              : 'bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white'
          ]"
        >
          + {{ skill.name }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { SKILL_CATEGORIES } from '~/types'

interface Skill {
  id: number
  name: string
  category: string
  teacherCount?: number
}

const props = defineProps<{
  modelValue: Skill[]
}>()

const emit = defineEmits<{
  'update:modelValue': [skills: Skill[]]
}>()

const inputValue = ref('')
const showSuggestions = ref(false)
const focusedIndex = ref(-1)
const suggestions = ref<Skill[]>([])
const popularSkills = ref<Skill[]>([])

// 計算屬性：雙向綁定
const selectedSkills = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 熱門技能列表（從類別定義取得）
const initPopularSkills = () => {
  const popular: Skill[] = []
  Object.entries(SKILL_CATEGORIES).forEach(([key, category]) => {
    if (category.skills) {
      category.skills.slice(0, 3).forEach((skillName: string) => {
        popular.push({
          id: popular.length + 1,
          name: skillName,
          category: key,
          teacherCount: Math.floor(Math.random() * 50) + 5 // 模擬數量
        })
      })
    }
  })
  popularSkills.value = popular.slice(0, 8)
}

// 篩選建議技能
const filteredSuggestions = computed(() => {
  if (!inputValue.value.trim()) return []
  
  const query = inputValue.value.toLowerCase().trim()
  
  // 先從熱門技能中搜尋
  let results = popularSkills.value.filter(
    skill => skill.name.toLowerCase().includes(query) && !isSkillSelected(skill.id)
  )
  
  // 如果結果不夠，從類別定義中補充
  if (results.length < 3) {
    Object.entries(SKILL_CATEGORIES).forEach(([key, category]) => {
      if (category.skills) {
        category.skills.forEach((skillName: string) => {
          if (skillName.toLowerCase().includes(query)) {
            const existing = results.find(s => s.name === skillName)
            if (!existing && !isSkillSelectedByName(skillName)) {
              results.push({
                id: results.length + 100,
                name: skillName,
                category: key,
                teacherCount: Math.floor(Math.random() * 30) + 1
              })
            }
          }
        })
      }
    })
  }
  
  return results.slice(0, 8)
})

// 檢查技能是否已選取
const isSkillSelected = (skillId: number): boolean => {
  return selectedSkills.value.some(s => s.id === skillId)
}

const isSkillSelectedByName = (name: string): boolean => {
  return selectedSkills.value.some(s => s.name === name)
}

// 輸入事件
const onInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  inputValue.value = target.value
  showSuggestions.value = true
  focusedIndex.value = -1
  
  // 模擬搜尋建議（實際應該呼叫 API）
  if (inputValue.value.trim()) {
    // 這裡可以呼叫後端 API 取得技能建議
  }
}

// 選擇技能
const selectSkill = (skill: Skill) => {
  addSkill(skill)
  showSuggestions.value = false
  inputValue.value = ''
  focusedIndex.value = -1
}

// 新增技能
const addSkill = (skill: Skill) => {
  if (isSkillSelected(skill.id)) return
  
  selectedSkills.value = [...selectedSkills.value, skill]
  inputValue.value = ''
  showSuggestions.value = false
}

// 新增目前輸入框中的文字
const addCurrentInput = () => {
  if (focusedIndex.value >= 0 && filteredSuggestions.value[focusedIndex.value]) {
    selectSkill(filteredSuggestions.value[focusedIndex.value])
    return
  }
  
  const skillName = inputValue.value.trim()
  if (skillName && !isSkillSelectedByName(skillName)) {
    selectedSkills.value = [...selectedSkills.value, {
      id: Date.now(),
      name: skillName,
      category: '其他',
      teacherCount: 0
    }]
    inputValue.value = ''
  }
  showSuggestions.value = false
}

// 移除技能
const removeSkill = (index: number) => {
  selectedSkills.value = selectedSkills.value.filter((_, i) => i !== index)
}

// 按下 backspace 且輸入框為空時，移除最後一個選取的技能
const onBackspace = () => {
  if (inputValue.value === '' && selectedSkills.value.length > 0) {
    selectedSkills.value = selectedSkills.value.slice(0, -1)
  }
}

// 點擊外部關閉建議
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.skill-selector')) {
    showSuggestions.value = false
  }
}

onMounted(() => {
  initPopularSkills()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
