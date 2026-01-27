<template>
  <div
    :class="[
      'match-card relative p-4 rounded-xl border transition-all cursor-pointer group',
      selected
        ? 'bg-indigo-500/10 border-indigo-500/30'
        : 'bg-white/5 border-white/10 hover:bg-white/10 hover:border-white/20'
    ]"
    @click="$emit('click')"
  >
    <!-- 選取核取方塊 -->
    <div
      v-if="showCheckbox"
      class="absolute top-3 right-3 z-10"
      @click.stop
    >
      <input
        type="checkbox"
        :checked="selected"
        @change="$emit('update:selected', !selected)"
        class="w-5 h-5 rounded border-white/20 bg-white/10 text-indigo-500 focus:ring-indigo-500 focus:ring-offset-0"
      />
    </div>

    <!-- 老師基本資訊 -->
    <div class="flex items-start gap-3 mb-4">
      <div class="w-12 h-12 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0">
        <span class="text-white font-medium">{{ match.teacher_name?.charAt(0) || '?' }}</span>
      </div>
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 mb-1">
          <h4 class="text-white font-medium truncate">{{ match.teacher_name }}</h4>
          <BaseBadge
            v-if="match.is_member"
            variant="success"
            size="xs"
          >
            中心成員
          </BaseBadge>
        </div>
        <!-- 可用性標籤 -->
        <div class="flex items-center gap-2">
          <span
            :class="[
              'px-2 py-0.5 rounded text-xs font-medium',
              match.availability === 'AVAILABLE' ? 'bg-green-500/20 text-green-400' :
              match.availability === 'BUFFER_CONFLICT' ? 'bg-yellow-500/20 text-yellow-400' :
              'bg-red-500/20 text-red-400'
            ]"
          >
            {{ availabilityText }}
          </span>
        </div>
      </div>
    </div>

    <!-- 總分進度條 -->
    <div class="mb-4">
      <div class="flex items-center justify-between mb-1">
        <span class="text-xs text-slate-400">媒合分數</span>
        <span class="text-lg font-bold text-primary-500">{{ match.match_score }}%</span>
      </div>
      <div class="h-2 bg-white/10 rounded-full overflow-hidden">
        <div
          class="h-full bg-gradient-to-r from-indigo-500 to-purple-500 rounded-full transition-all duration-500"
          :style="{ width: `${match.match_score}%` }"
        />
      </div>
    </div>

    <!-- 分數 breakdown -->
    <div class="grid grid-cols-2 gap-2 mb-4">
      <div class="flex items-center gap-2 text-xs">
        <div class="w-6 h-6 rounded bg-green-500/20 flex items-center justify-center">
          <svg class="w-3.5 h-3.5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div>
          <span class="text-slate-400 block">可用性</span>
          <span class="text-white font-medium">{{ match.availability_score }}/40</span>
        </div>
      </div>
      <div class="flex items-center gap-2 text-xs">
        <div class="w-6 h-6 rounded bg-yellow-500/20 flex items-center justify-center">
          <svg class="w-3.5 h-3.5 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
          </svg>
        </div>
        <div>
          <span class="text-slate-400 block">內部評分</span>
          <span class="text-white font-medium">{{ match.internal_score }}/40</span>
        </div>
      </div>
      <div class="flex items-center gap-2 text-xs">
        <div class="w-6 h-6 rounded bg-blue-500/20 flex items-center justify-center">
          <svg class="w-3.5 h-3.5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
          </svg>
        </div>
        <div>
          <span class="text-slate-400 block">技能匹配</span>
          <span class="text-white font-medium">{{ match.skill_score }}/10</span>
        </div>
      </div>
      <div class="flex items-center gap-2 text-xs">
        <div class="w-6 h-6 rounded bg-purple-500/20 flex items-center justify-center">
          <svg class="w-3.5 h-3.5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>
        <div>
          <span class="text-slate-400 block">地區匹配</span>
          <span class="text-white font-medium">{{ match.region_score || 0 }}/10</span>
        </div>
      </div>
    </div>

    <!-- 懸停顯示的詳細資訊 -->
    <div
      class="absolute inset-x-0 bottom-0 p-4 bg-slate-800/95 backdrop-blur-sm rounded-b-xl border-t border-white/10 transition-all duration-200"
      :class="hovered ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2 pointer-events-none'"
    >
      <h5 class="text-sm font-medium text-white mb-2">專長技能</h5>
      <div v-if="match.skills?.length" class="flex flex-wrap gap-1 mb-3">
        <BaseBadge
          v-for="skill in match.skills.slice(0, 5)"
          :key="skill.name"
          variant="primary"
          size="xs"
        >
          {{ skill.name }}
        </BaseBadge>
      </div>
      <p v-else class="text-xs text-slate-400 mb-3">無技能資料</p>
      
      <h5 v-if="match.notes" class="text-sm font-medium text-white mb-1">內部備註</h5>
      <p v-if="match.notes" class="text-xs text-slate-400">{{ match.notes }}</p>
      
      <div class="flex items-center gap-2 mt-3 pt-3 border-t border-white/10">
        <span v-if="match.rating" class="text-xs text-yellow-400">
          {{ '★'.repeat(match.rating) }}{{ '☆'.repeat(5 - match.rating) }}
          <span class="text-slate-400 ml-1">{{ match.rating }}.0</span>
        </span>
        <span class="text-xs text-slate-400">
          技能匹配 {{ match.skill_match }}%
        </span>
      </div>
    </div>

    <!-- 懸停提示 -->
    <div
      class="absolute inset-x-0 bottom-0 h-8 rounded-b-xl transition-colors"
      :class="hovered ? 'bg-transparent' : 'bg-gradient-to-t from-white/5 to-transparent'"
    />
  </div>
</template>

<script setup lang="ts">
interface Skill {
  name: string
  category: string
}

interface MatchResult {
  teacher_id: number
  teacher_name: string
  match_score: number
  availability: string
  availability_score: number
  internal_score: number
  skill_match: number
  skill_score: number
  region_score?: number
  rating?: number
  is_member: boolean
  notes?: string
  skills?: Skill[]
}

const props = defineProps<{
  match: MatchResult
  selected?: boolean
  showCheckbox?: boolean
}>()

const emit = defineEmits<{
  click: []
  'update:selected': [selected: boolean]
}>()

const hovered = ref(false)

// 可用性文字轉換
const availabilityText = computed(() => {
  switch (props.match.availability) {
    case 'AVAILABLE':
      return '完全可用'
    case 'BUFFER_CONFLICT':
      return '緩衝衝突'
    case 'OVERLAP':
      return '時間重疊'
    default:
      return props.match.availability
  }
})
</script>
