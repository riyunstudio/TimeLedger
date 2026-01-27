<template>
  <div class="p-4 md:p-6">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-white">系統監控</h1>
        <p class="text-sm text-slate-400 mt-1">監控通知佇列與系統狀態</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="refreshStats"
          :disabled="loading"
          class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors flex items-center gap-2 disabled:opacity-50"
        >
          <svg class="w-4 h-4" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          重新整理
        </button>
      </div>
    </div>

    <!-- 通知佇列狀態 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-6">
      <!-- 待處理 -->
      <div class="glass-card p-4">
        <div class="flex items-center justify-between mb-2">
          <div class="w-10 h-10 rounded-lg bg-primary-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <span class="text-xs px-2 py-1 rounded-full bg-primary-500/20 text-primary-400">待處理</span>
        </div>
        <p class="text-3xl font-bold text-white">{{ queueStats.pending }}</p>
        <p class="text-xs text-slate-400 mt-1">筆等待發送</p>
      </div>

      <!-- 重試中 -->
      <div class="glass-card p-4">
        <div class="flex items-center justify-between mb-2">
          <div class="w-10 h-10 rounded-lg bg-yellow-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-yellow-500 animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </div>
          <span class="text-xs px-2 py-1 rounded-full bg-yellow-500/20 text-yellow-400">重試中</span>
        </div>
        <p class="text-3xl font-bold text-white">{{ queueStats.retry }}</p>
        <p class="text-xs text-slate-400 mt-1">筆稍後重試</p>
      </div>

      <!-- 已完成 -->
      <div class="glass-card p-4">
        <div class="flex items-center justify-between mb-2">
          <div class="w-10 h-10 rounded-lg bg-success-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-success-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <span class="text-xs px-2 py-1 rounded-full bg-success-500/20 text-success-400">已完成</span>
        </div>
        <p class="text-3xl font-bold text-white">{{ queueStats.total }}</p>
        <p class="text-xs text-slate-400 mt-1">筆已發送</p>
      </div>

      <!-- 已重試 -->
      <div class="glass-card p-4">
        <div class="flex items-center justify-between mb-2">
          <div class="w-10 h-10 rounded-lg bg-blue-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </div>
          <span class="text-xs px-2 py-1 rounded-full bg-blue-500/20 text-blue-400">已重試</span>
        </div>
        <p class="text-3xl font-bold text-white">{{ queueStats.retried }}</p>
        <p class="text-xs text-slate-400 mt-1">筆重新處理</p>
      </div>

      <!-- 失敗 -->
      <div class="glass-card p-4">
        <div class="flex items-center justify-between mb-2">
          <div class="w-10 h-10 rounded-lg bg-red-500/20 flex items-center justify-center">
            <svg class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <span class="text-xs px-2 py-1 rounded-full bg-red-500/20 text-red-400">失敗</span>
        </div>
        <p class="text-3xl font-bold text-white">{{ queueStats.failed }}</p>
        <p class="text-xs text-slate-400 mt-1">筆發送失敗</p>
      </div>
    </div>

    <!-- 失敗率警示 -->
    <div v-if="failureRate > 10" class="mb-6 p-4 bg-red-500/10 border border-red-500/30 rounded-xl">
      <div class="flex items-center gap-3">
        <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div>
          <p class="text-red-400 font-medium">警告：失敗率過高</p>
          <p class="text-sm text-slate-400">当前失败率为 {{ failureRate.toFixed(1) }}%，請檢查 LINE API 連線或通知設定</p>
        </div>
      </div>
    </div>

    <!-- 系統狀態與人才庫統計 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Redis 連線狀態 -->
      <div class="glass-card p-6">
        <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
          </svg>
          Redis 連線狀態
        </h3>

        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-white/5 rounded-lg">
            <span class="text-slate-400">連線狀態</span>
            <div class="flex items-center gap-2">
              <span
                :class="[
                  'w-2 h-2 rounded-full',
                  redisStatus.connected ? 'bg-success-500' : 'bg-red-500'
                ]"
              />
              <span :class="redisStatus.connected ? 'text-success-500' : 'text-red-500'">
                {{ redisStatus.connected ? '已連線' : '未連線' }}
              </span>
            </div>
          </div>

          <div class="flex items-center justify-between p-3 bg-white/5 rounded-lg">
            <span class="text-slate-400">回應時間</span>
            <span class="text-white font-mono">{{ redisStatus.latency }}ms</span>
          </div>

          <div class="flex items-center justify-between p-3 bg-white/5 rounded-lg">
            <span class="text-slate-400">記憶體使用</span>
            <span class="text-white font-mono">{{ redisStatus.memoryUsage }}</span>
          </div>
        </div>
      </div>

      <!-- 人才庫邀請統計 -->
      <div class="glass-card p-6">
        <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          人才庫邀請統計
        </h3>

        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-white/5 rounded-lg">
            <span class="text-slate-400">待回覆邀請</span>
            <span class="text-yellow-400 font-medium">{{ invitationStats.pending }}</span>
          </div>

          <div class="flex items-center justify-between p-3 bg-white/5 rounded-lg">
            <span class="text-slate-400">已接受</span>
            <span class="text-success-500 font-medium">{{ invitationStats.accepted }}</span>
          </div>

          <div class="flex items-center justify-between p-3 bg-white/5 rounded-lg">
            <span class="text-slate-400">已拒絕</span>
            <span class="text-red-400 font-medium">{{ invitationStats.declined }}</span>
          </div>

          <div class="pt-2 border-t border-white/10">
            <div class="flex items-center justify-between">
              <span class="text-slate-400">回覆率</span>
              <span class="text-primary-400 font-bold">{{ replyRate }}%</span>
            </div>
            <div class="h-2 bg-white/10 rounded-full mt-2 overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-primary-500 to-success-500 transition-all duration-500"
                :style="{ width: `${replyRate}%` }"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 最近的失敗記錄 -->
    <div class="mt-6 glass-card p-6">
      <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
        <svg class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        最近的失敗記錄
      </h3>

      <div v-if="failedRecords.length === 0" class="text-center py-8 text-slate-500">
        <svg class="w-12 h-12 mx-auto mb-3 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p>目前沒有失敗記錄</p>
      </div>

      <div v-else class="space-y-2">
        <div
          v-for="(record, index) in failedRecords"
          :key="index"
          class="flex items-center justify-between p-3 bg-white/5 rounded-lg hover:bg-white/10 transition-colors"
        >
          <div class="flex items-center gap-3">
            <span class="text-red-500">✕</span>
            <div>
              <p class="text-white text-sm">{{ record.type }}</p>
              <p class="text-xs text-slate-500">{{ record.recipient }}</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-xs text-slate-400">{{ record.time }}</p>
            <p class="text-xs text-red-400">{{ record.error }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth-admin',
  layout: 'admin'
})

const { warning: alertWarning } = useAlert()
const api = useApi()

const loading = ref(false)

// 通知佇列統計
const queueStats = ref({
  pending: '0',
  retry: '0',
  total: '0',
  retried: '0',
  failed: '0'
})

// Redis 狀態
const redisStatus = ref({
  connected: false,
  latency: 0,
  memoryUsage: '0 MB'
})

// 人才庫邀請統計
const invitationStats = ref({
  pending: 0,
  accepted: 0,
  declined: 0
})

// 失敗記錄
const failedRecords = ref<Array<{
  type: string
  recipient: string
  time: string
  error: string
}>>([])

// 計算失敗率
const failureRate = computed(() => {
  const total = parseInt(queueStats.value.total) || 0
  const failed = parseInt(queueStats.value.failed) || 0
  if (total === 0) return 0
  return (failed / total) * 100
})

// 計算回覆率
const replyRate = computed(() => {
  const total = invitationStats.value.accepted + invitationStats.value.declined
  if (total === 0) return 0
  return Math.round((invitationStats.value.accepted / total) * 100)
})

// 取得通知佇列統計
const fetchQueueStats = async () => {
  try {
    const response = await api.get<{ code: number; datas: any }>(
      '/admin/notifications/queue-stats'
    )
    if (response.code === 0 && response.datas) {
      queueStats.value = {
        pending: response.datas.pending || '0',
        retry: response.datas.retry || '0',
        total: response.datas.total || '0',
        retried: response.datas.retried || '0',
        failed: response.datas.failed || '0'
      }
    }
  } catch (error) {
    console.error('Failed to fetch queue stats:', error)
  }
}

// 取得 Redis 狀態
const fetchRedisStatus = async () => {
  try {
    // 嘗試取得 Redis 健康檢查
    const startTime = Date.now()
    const response = await api.get<{ code: number; datas: any }>(
      '/admin/health/redis'
    )
    const latency = Date.now() - startTime

    redisStatus.value = {
      connected: response.code === 0,
      latency: latency,
      memoryUsage: response.datas?.memory_usage || 'N/A'
    }
  } catch (error) {
    redisStatus.value = {
      connected: false,
      latency: 0,
      memoryUsage: 'N/A'
    }
  }
}

// 取得人才庫邀請統計
const fetchInvitationStats = async () => {
  try {
    const { getCenterId } = useCenterId()
    const centerId = getCenterId()

    const response = await api.get<{ code: number; datas: any }>(
      '/admin/smart-matching/talent/stats'
    )
    if (response.code === 0 && response.datas) {
      invitationStats.value = {
        pending: response.datas.pending_invites || 0,
        accepted: response.datas.accepted_invites || 0,
        declined: response.datas.declined_invites || 0
      }
    }
  } catch (error) {
    console.error('Failed to fetch invitation stats:', error)
  }
}

// 取得失敗記錄
const fetchFailedRecords = async () => {
  // 這裡應該呼叫 API 取得失敗記錄，目前使用模擬資料
  failedRecords.value = [
    {
      type: '例外通知',
      recipient: '管理員 #3',
      time: '2 分鐘前',
      error: 'LINE API 逾時'
    },
    {
      type: '邀請通知',
      recipient: '老師 #12',
      time: '15 分鐘前',
      error: '無效的 LineUserID'
    }
  ]
}

// 重新整理所有數據
const refreshStats = async () => {
  loading.value = true
  try {
    await Promise.all([
      fetchQueueStats(),
      fetchRedisStatus(),
      fetchInvitationStats(),
      fetchFailedRecords()
    ])
  } finally {
    loading.value = false
  }
}

// 自動重新整理（每 30 秒）
let autoRefreshInterval: NodeJS.Timeout | null = null

onMounted(() => {
  refreshStats()

  // 啟動自動重新整理
  autoRefreshInterval = setInterval(() => {
    fetchQueueStats()
    fetchRedisStatus()
  }, 30000)
})

onUnmounted(() => {
  if (autoRefreshInterval) {
    clearInterval(autoRefreshInterval)
  }
})
</script>
