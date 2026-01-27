<template>
  <div class="talent-stats-panel">
    <!-- 統計卡片網格 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <!-- 總人數 -->
      <div class="stat-card bg-white/5 rounded-xl p-4 hover:bg-white/10 transition-colors">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 rounded-lg bg-primary-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <div class="text-xs px-2 py-1 rounded-full" :class="stats.monthlyChange > 0 ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'">
            {{ stats.monthlyChange > 0 ? '+' : '' }}{{ stats.monthlyChange }}%
          </div>
        </div>
        <div class="text-2xl font-bold text-white mb-1">{{ formatNumber(stats.totalCount) }}</div>
        <div class="text-xs text-slate-400">人才庫總人數</div>
        <!-- 趨勢圖 -->
        <div class="mt-3 h-8 flex items-end gap-1">
          <div
            v-for="(value, index) in stats.monthlyTrend"
            :key="index"
            class="flex-1 bg-gradient-to-t from-primary-500/50 to-primary-500/20 rounded-t"
            :style="{ height: `${(value / maxTrendValue) * 100}%` }"
          />
        </div>
      </div>

      <!-- 開放徵才 -->
      <div class="stat-card bg-white/5 rounded-xl p-4 hover:bg-white/10 transition-colors">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 rounded-lg bg-green-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="text-xs px-2 py-1 rounded-full bg-green-500/20 text-green-400">
            {{ ((stats.openHiringCount / stats.totalCount) * 100).toFixed(0) }}%
          </div>
        </div>
        <div class="text-2xl font-bold text-white mb-1">{{ formatNumber(stats.openHiringCount) }}</div>
        <div class="text-xs text-slate-400">開放徵才</div>
        <div class="mt-2 text-xs text-slate-500">
          <span class="text-green-400">{{ stats.pendingInvites }}</span> 待回覆邀請
        </div>
      </div>

      <!-- 中心成員 -->
      <div class="stat-card bg-white/5 rounded-xl p-4 hover:bg-white/10 transition-colors">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 rounded-lg bg-blue-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          </div>
        </div>
        <div class="text-2xl font-bold text-white mb-1">{{ formatNumber(stats.memberCount) }}</div>
        <div class="text-xs text-slate-400">中心成員</div>
        <div class="mt-2 flex items-center gap-1">
          <div class="flex-1 h-1.5 bg-white/10 rounded-full overflow-hidden">
            <div class="h-full bg-blue-500 rounded-full" :style="{ width: `${(stats.memberCount / stats.totalCount) * 100}%` }" />
          </div>
        </div>
      </div>

      <!-- 平均評分 -->
      <div class="stat-card bg-white/5 rounded-xl p-4 hover:bg-white/10 transition-colors">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 rounded-lg bg-yellow-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
          </div>
          <div class="text-xs px-2 py-1 rounded-full bg-yellow-500/20 text-yellow-400">
            {{ stats.averageRating.toFixed(1) }}
          </div>
        </div>
        <div class="flex items-end gap-1 mb-1">
          <span class="text-2xl font-bold text-white">{{ stats.averageRating.toFixed(1) }}</span>
          <span class="text-sm text-slate-400 mb-1">/ 5.0</span>
        </div>
        <div class="text-xs text-slate-400">平均評分</div>
        <div class="flex items-center gap-0.5 mt-2">
          <svg v-for="i in 5" :key="i" class="w-3 h-3" :class="i <= Math.round(stats.averageRating) ? 'text-yellow-400' : 'text-slate-600'" fill="currentColor" viewBox="0 0 20 20">
            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.364 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
          </svg>
        </div>
      </div>
    </div>

    <!-- 第二行：地區分布與熱門技能 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
      <!-- 地區分布 -->
      <div class="bg-white/5 rounded-xl p-4">
        <h4 class="text-sm font-medium text-white mb-3 flex items-center gap-2">
          <svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          地區分布
        </h4>
        <div class="space-y-2">
          <div
            v-for="city in cityDistribution"
            :key="city.name"
            class="flex items-center gap-2"
          >
            <span class="text-xs text-slate-400 w-12">{{ city.name }}</span>
            <div class="flex-1 h-2 bg-white/10 rounded-full overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-blue-500 to-blue-400 rounded-full transition-all duration-500"
                :style="{ width: `${(city.count / maxCityCount) * 100}%` }"
              />
            </div>
            <span class="text-xs text-slate-500 w-8 text-right">{{ city.count }}</span>
          </div>
        </div>
      </div>

      <!-- 熱門技能 Top 5 -->
      <div class="bg-white/5 rounded-xl p-4">
        <h4 class="text-sm font-medium text-white mb-3 flex items-center gap-2">
          <svg class="w-4 h-4 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
          </svg>
          熱門技能
        </h4>
        <div class="flex flex-wrap gap-2">
          <div
            v-for="(skill, index) in topSkills"
            :key="skill.name"
            class="flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-white/5"
          >
            <span class="text-xs font-medium text-white">{{ skill.name }}</span>
            <span class="text-xs text-slate-500">{{ skill.count }}</span>
          </div>
        </div>
      </div>

      <!-- 招募狀態 -->
      <div class="bg-white/5 rounded-xl p-4">
        <h4 class="text-sm font-medium text-white mb-3 flex items-center gap-2">
          <svg class="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          招募狀態
        </h4>
        <div class="space-y-2">
          <div class="flex items-center justify-between text-sm">
            <span class="text-slate-400">已接受邀請</span>
            <span class="text-green-400 font-medium">{{ stats.acceptedInvites }}</span>
          </div>
          <div class="flex items-center justify-between text-sm">
            <span class="text-slate-400">待回覆</span>
            <span class="text-yellow-400 font-medium">{{ stats.pendingInvites }}</span>
          </div>
          <div class="flex items-center justify-between text-sm">
            <span class="text-slate-400">已拒絕</span>
            <span class="text-red-400 font-medium">{{ stats.declinedInvites }}</span>
          </div>
          <div class="pt-2 mt-2 border-t border-white/10">
            <div class="flex items-center justify-between text-sm">
              <span class="text-slate-400">回覆率</span>
              <span class="text-primary-400 font-medium">{{ replyRate }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Stats {
  totalCount: number
  openHiringCount: number
  memberCount: number
  averageRating: number
  monthlyChange: number
  monthlyTrend: number[]
  pendingInvites: number
  acceptedInvites: number
  declinedInvites: number
}

interface CityDistribution {
  name: string
  count: number
}

interface SkillCount {
  name: string
  count: number
}

const props = defineProps<{
  stats: Stats
  cityDistribution: CityDistribution[]
  topSkills: SkillCount[]
}>()

// 計算屬性：最大趨勢值
const maxTrendValue = computed(() => {
  return Math.max(...props.stats.monthlyTrend, 1)
})

// 計算屬性：最大城市數量
const maxCityCount = computed(() => {
  return Math.max(...props.cityDistribution.map(c => c.count), 1)
})

// 計算屬性：回覆率
const replyRate = computed(() => {
  const total = props.stats.acceptedInvites + props.stats.pendingInvites + props.stats.declinedInvites
  if (total === 0) return 0
  return Math.round((props.stats.acceptedInvites / (props.stats.acceptedInvites + props.stats.declinedInvites)) * 100)
})

// 格式化數字
const formatNumber = (num: number): string => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}
</script>
