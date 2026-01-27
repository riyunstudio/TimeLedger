<template>
  <div class="skills-distribution-chart bg-white/5 rounded-xl p-4">
    <!-- 標題與檢視切換 -->
    <div class="flex items-center justify-between mb-4">
      <h4 class="text-sm font-medium text-white flex items-center gap-2">
        <svg class="w-4 h-4 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        技能分布
      </h4>
      <div class="flex items-center gap-1 bg-white/5 rounded-lg p-1">
        <button
          :class="['p-1 rounded', viewMode === 'bar' ? 'bg-indigo-500 text-white' : 'text-slate-400 hover:text-white']"
          @click="viewMode = 'bar'"
          title="長條圖"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 8v8m-4-5v5m-4-2v2m-2 4h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </button>
        <button
          :class="['p-1 rounded', viewMode === 'pie' ? 'bg-indigo-500 text-white' : 'text-slate-400 hover:text-white']"
          @click="viewMode = 'pie'"
          title="圓餅圖"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 長條圖檢視 -->
    <div v-if="viewMode === 'bar'" class="space-y-3">
      <div
        v-for="skill in skills"
        :key="skill.name"
        class="skill-bar group"
      >
        <div class="flex items-center justify-between mb-1">
          <div class="flex items-center gap-2">
            <span class="text-sm text-slate-300 group-hover:text-white transition-colors">{{ skill.name }}</span>
            <BaseBadge v-if="isTopSkill(skill.name)" variant="warning" size="xs">
              Hot
            </BaseBadge>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xs text-slate-500">{{ skill.count }} 人</span>
            <span class="text-xs text-slate-400">{{ getPercentage(skill.count) }}%</span>
          </div>
        </div>
        <div class="h-3 bg-white/10 rounded-full overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-700 ease-out"
            :class="getBarColor(skill.name)"
            :style="{ width: `${(skill.count / maxCount) * 100}%` }"
          />
        </div>
      </div>
    </div>

    <!-- 圓餅圖檢視 -->
    <div v-else class="flex items-center gap-6">
      <!-- SVG 圓餅圖 -->
      <div class="relative w-40 h-40 shrink-0">
        <svg viewBox="0 0 100 100" class="w-full h-full -rotate-90">
          <circle
            v-for="(skill, index) in skills"
            :key="skill.name"
            cx="50"
            cy="50"
            r="40"
            fill="transparent"
            :stroke="getPieColor(index)"
            stroke-width="20"
            :stroke-dasharray="getStrokeDasharray(skill.count, index)"
            :stroke-dashoffset="getStrokeDashoffset(index)"
            class="transition-all duration-700 ease-out hover:opacity-80 cursor-pointer"
            @click="$emit('selectSkill', skill.name)"
          />
        </svg>
        <!-- 中心文字 -->
        <div class="absolute inset-0 flex flex-col items-center justify-center">
          <span class="text-2xl font-bold text-white">{{ totalCount }}</span>
          <span class="text-xs text-slate-400">總計</span>
        </div>
      </div>

      <!-- 圖例 -->
      <div class="flex-1 grid grid-cols-2 gap-2">
        <div
          v-for="(skill, index) in skills"
          :key="skill.name"
          class="flex items-center gap-2 p-2 rounded-lg hover:bg-white/5 cursor-pointer transition-colors"
          @click="$emit('selectSkill', skill.name)"
        >
          <div
            class="w-3 h-3 rounded-full shrink-0"
            :style="{ backgroundColor: getPieColor(index) }"
          />
          <div class="min-w-0">
            <div class="text-sm text-slate-300 truncate">{{ skill.name }}</div>
            <div class="text-xs text-slate-500">{{ getPercentage(skill.count) }}%</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 技能詳細統計 -->
    <div class="mt-4 pt-4 border-t border-white/10">
      <div class="flex items-center justify-between text-xs">
        <span class="text-slate-500">涵蓋 {{ skills.length }} 種技能類型</span>
        <span class="text-slate-400">平均每人 {{ averageSkillsPerPerson.toFixed(1) }} 種技能</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { SKILL_CATEGORIES } from '~/types'

interface SkillCount {
  name: string
  count: number
}

const props = defineProps<{
  skills: SkillCount[]
  totalTeachers: number
}>()

const emit = defineEmits<{
  selectSkill: [skillName: string]
}>()

const viewMode = ref<'bar' | 'pie'>('bar')

// 計算屬性：最大數量
const maxCount = computed(() => {
  return Math.max(...props.skills.map(s => s.count), 1)
})

// 計算屬性：總數
const totalCount = computed(() => {
  return props.skills.reduce((sum, s) => sum + s.count, 0)
})

// 計算屬性：每人平均技能數
const averageSkillsPerPerson = computed(() => {
  if (props.totalTeachers === 0) return 0
  return totalCount.value / props.totalTeachers
})

// 取得百分比
const getPercentage = (count: number): string => {
  return ((count / totalCount.value) * 100).toFixed(1)
}

// 判斷是否為熱門技能
const isTopSkill = (name: string): boolean => {
  const topSkills = ['瑜珈', '鋼琴', '舞蹈', '美術', '英語']
  return topSkills.includes(name)
}

// 取得長條圖顏色
const getBarColor = (name: string): string => {
  const topSkills = ['瑜珈', '鋼琴', '舞蹈']
  if (topSkills.includes(name)) {
    return 'bg-gradient-to-r from-yellow-500 to-orange-500'
  }
  return 'bg-gradient-to-r from-primary-500 to-secondary-500'
}

// 取得圓餅圖顏色
const getPieColor = (index: number): string => {
  const colors = [
    '#6366f1', // indigo
    '#a855f7', // purple
    '#f59e0b', // yellow
    '#10b981', // green
    '#3b82f6', // blue
    '#ef4444', // red
    '#ec4899', // pink
    '#14b8a6'  // teal
  ]
  return colors[index % colors.length]
}

// 計算 stroke-dasharray
const getStrokeDasharray = (count: number, index: number): string => {
  const circumference = 2 * Math.PI * 40 // r=40
  const percentage = count / totalCount.value
  const dashLength = percentage * circumference
  return `${dashLength} ${circumference}`
}

// 計算 stroke-dashoffset
const getStrokeDashoffset = (index: number): string => {
  const circumference = 2 * Math.PI * 40
  let offset = 0
  
  for (let i = 0; i < index; i++) {
    const percentage = props.skills[i].count / totalCount.value
    offset -= percentage * circumference
  }
  
  return String(offset)
}
</script>
