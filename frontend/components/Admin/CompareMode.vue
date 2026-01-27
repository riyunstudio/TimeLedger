<template>
  <div class="compare-mode">
    <!-- 比較標題列 -->
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-white">老師比較</h3>
      <div class="flex items-center gap-2">
        <span class="text-sm text-slate-400">
          最多可選取 3 位老師进行比较
        </span>
        <button
          @click="$emit('exit')"
          class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-white transition-colors"
          title="退出比較"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 比較表格 -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr>
            <th class="w-32 p-3 text-left text-sm font-medium text-slate-400 bg-white/5 rounded-tl-lg">
              項目
            </th>
            <th
              v-for="teacher in selectedTeachers"
              :key="teacher.teacher_id"
              class="p-3 text-center bg-white/5 min-w-[200px]"
              :class="{
                'rounded-tr-lg': teacher === selectedTeachers[selectedTeachers.length - 1]
              }"
            >
              <div class="flex flex-col items-center">
                <div class="w-12 h-12 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center mb-2">
                  <span class="text-white font-medium">{{ teacher.teacher_name?.charAt(0) || '?' }}</span>
                </div>
                <span class="text-white font-medium">{{ teacher.teacher_name }}</span>
                <button
                  @click="$emit('remove', teacher)"
                  class="mt-2 text-xs text-slate-400 hover:text-red-400 transition-colors"
                >
                  移除
                </button>
              </div>
            </th>
            <!-- 填充空白欄位 -->
            <th
              v-for="n in Math.max(0, 3 - selectedTeachers.length)"
              :key="n"
              class="p-3 bg-white/5 border border-dashed border-white/10"
            >
              <div class="flex items-center justify-center h-full text-slate-500 text-sm">
                選擇老師加入比較
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in compareRows" :key="row.key" class="border-t border-white/5">
            <td class="p-3 text-sm text-slate-400 bg-white/5">
              {{ row.label }}
            </td>
            <td
              v-for="teacher in selectedTeachers"
              :key="teacher.teacher_id"
              :class="[
                'p-3 text-center',
                row.highlight && row.getValue(teacher) === row.maxValue ? 'bg-green-500/10' : ''
              ]"
            >
              <template v-if="row.key === 'score'">
                <div class="flex items-center justify-center gap-2">
                  <div class="w-20 h-2 bg-white/10 rounded-full overflow-hidden">
                    <div
                      class="h-full bg-gradient-to-r from-indigo-500 to-purple-500 rounded-full"
                      :style="{ width: `${row.getValue(teacher)}%` }"
                    />
                  </div>
                  <span class="text-white font-medium">{{ row.getValue(teacher) }}%</span>
                </div>
              </template>
              <template v-else-if="row.key === 'is_member'">
                <span :class="row.getValue(teacher) ? 'text-green-400' : 'text-slate-400'">
                  {{ row.getValue(teacher) ? '✓ 是' : '✗ 否' }}
                </span>
              </template>
              <template v-else-if="row.key === 'rating'">
                <div class="flex items-center justify-center gap-1">
                  <span class="text-yellow-400">{{ '★'.repeat(row.getValue(teacher)) }}</span>
                  <span class="text-slate-400 ml-1">{{ row.getValue(teacher) }}.0</span>
                </div>
              </template>
              <template v-else-if="row.key === 'availability'">
                <span
                  :class="[
                    'px-2 py-1 rounded text-xs font-medium',
                    row.getValue(teacher) === 'AVAILABLE' ? 'bg-green-500/20 text-green-400' :
                    row.getValue(teacher) === 'BUFFER_CONFLICT' ? 'bg-yellow-500/20 text-yellow-400' :
                    'bg-red-500/20 text-red-400'
                  ]"
                >
                  {{ availabilityText(row.getValue(teacher)) }}
                </span>
              </template>
              <template v-else>
                <span class="text-white">{{ row.getValue(teacher) }}</span>
              </template>
            </td>
            <!-- 填充空白 -->
            <td
              v-for="n in Math.max(0, 3 - selectedTeachers.length)"
              :key="n"
              class="p-3 border border-dashed border-white/10"
            ></td>
          </tr>
          
          <!-- 操作列 -->
          <tr class="border-t border-white/5">
            <td class="p-3 bg-white/5"></td>
            <td
              v-for="teacher in selectedTeachers"
              :key="teacher.teacher_id"
              class="p-3"
            >
              <button
                @click="$emit('select', teacher)"
                class="w-full px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors text-sm"
              >
                選擇這位老師
              </button>
            </td>
            <td
              v-for="n in Math.max(0, 3 - selectedTeachers.length)"
              :key="n"
              class="p-3 border border-dashed border-white/10"
            ></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 圖例說明 -->
    <div class="mt-4 flex flex-wrap items-center gap-4 text-sm">
      <div class="flex items-center gap-2">
        <div class="w-3 h-3 rounded-full bg-green-500/20 border border-green-500/30"></div>
        <span class="text-slate-400">該項目表現最佳</span>
      </div>
      <div class="flex items-center gap-2">
        <div class="flex items-center gap-1 text-xs text-slate-500">
          <span class="px-1.5 py-0.5 rounded bg-green-500/20 text-green-400">✓ 是</span>
          <span>中心成員</span>
        </div>
      </div>
      <div class="flex items-center gap-1 text-xs text-slate-500">
        <span class="px-1.5 py-0.5 rounded bg-yellow-500/20 text-yellow-400">緩衝衝突</span>
        <span>可_override</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface MatchResult {
  teacher_id: number
  teacher_name: string
  match_score: number
  availability: string
  availability_score: number
  internal_score: number
  skill_match: number
  skill_score: number
  rating?: number
  is_member: boolean
  notes?: string
}

const props = defineProps<{
  selectedTeachers: MatchResult[]
}>()

const emit = defineEmits<{
  remove: [teacher: MatchResult]
  select: [teacher: MatchResult]
  exit: []
}>()

// 可用性文字轉換
const availabilityText = (availability: string): string => {
  switch (availability) {
    case 'AVAILABLE':
      return '完全可用'
    case 'BUFFER_CONFLICT':
      return '緩衝衝突'
    case 'OVERLAP':
      return '時間重疊'
    default:
      return availability
  }
}

// 比較欄位定義
const compareRows = computed(() => {
  if (props.selectedTeachers.length === 0) return []
  
  const teachers = props.selectedTeachers
  
  const rows = [
    {
      key: 'score',
      label: '媒合分數',
      highlight: true,
      getValue: (t: MatchResult) => t.match_score,
      maxValue: Math.max(...teachers.map(t => t.match_score))
    },
    {
      key: 'availability',
      label: '可用性',
      highlight: false,
      getValue: (t: MatchResult) => t.availability,
      maxValue: ''
    },
    {
      key: 'internal_score',
      label: '內部評分',
      highlight: true,
      getValue: (t: MatchResult) => t.internal_score,
      maxValue: Math.max(...teachers.map(t => t.internal_score))
    },
    {
      key: 'skill_match',
      label: '技能匹配',
      highlight: true,
      getValue: (t: MatchResult) => t.skill_match,
      maxValue: Math.max(...teachers.map(t => t.skill_match))
    },
    {
      key: 'rating',
      label: '星等評分',
      highlight: true,
      getValue: (t: MatchResult) => t.rating || 0,
      maxValue: Math.max(...teachers.map(t => t.rating || 0))
    },
    {
      key: 'is_member',
      label: '中心成員',
      highlight: false,
      getValue: (t: MatchResult) => t.is_member,
      maxValue: ''
    }
  ]
  
  return rows
})
</script>
