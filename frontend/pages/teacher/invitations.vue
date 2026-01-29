<template>
  <div class="p-4 md:p-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
        🎁 邀請通知
      </h1>
      <p class="text-slate-400 text-sm md:text-base">
        管理您收到的中心邀請
      </p>
    </div>

    <!-- 篩選標籤 -->
    <div class="flex gap-2 mb-6 overflow-x-auto pb-2">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        @click="activeTab = tab.value"
        class="px-4 py-2 rounded-lg whitespace-nowrap transition-colors"
        :class="activeTab === tab.value 
          ? 'bg-primary-500/30 text-primary-400 border border-primary-500/50' 
          : 'bg-white/5 text-slate-400 hover:bg-white/10 border border-transparent'"
      >
        {{ tab.label }}
        <span v-if="tab.count > 0" class="ml-2 px-2 py-0.5 rounded-full text-xs bg-primary-500/20 text-primary-400">
          {{ tab.count }}
        </span>
      </button>
    </div>

    <!-- 邀請列表 -->
    <div class="space-y-4">
      <div v-if="loading" class="text-center py-12">
        <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full mx-auto mb-4"></div>
        <p class="text-slate-400">載入中...</p>
      </div>

      <div v-else-if="filteredInvitations.length === 0" class="text-center py-12">
        <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-slate-700/50 flex items-center justify-center">
          <svg class="w-10 h-10 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold text-slate-300 mb-2">沒有邀請</h3>
        <p class="text-slate-500">
          {{ activeTab === 'PENDING' ? '目前沒有待處理的邀請' : '沒有相關的邀請記錄' }}
        </p>
      </div>

      <div
        v-else
        v-for="invitation in filteredInvitations"
        :key="invitation.id"
        class="glass-card p-6"
      >
        <div class="flex items-start justify-between mb-4">
          <div class="flex items-center gap-4">
            <div 
              class="w-12 h-12 rounded-xl flex items-center justify-center"
              :class="getInviteTypeBgClass(invitation.invite_type)"
            >
              <span class="text-2xl">{{ getInviteTypeIcon(invitation.invite_type) }}</span>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-slate-100">{{ invitation.center_name }}</h3>
              <p class="text-sm text-slate-400">
                {{ getInviteTypeName(invitation.invite_type) }}
              </p>
            </div>
          </div>
          <span 
            class="px-3 py-1 rounded-full text-sm font-medium"
            :class="getStatusClass(invitation.status)"
          >
            {{ getStatusText(invitation.status) }}
          </span>
        </div>

        <!-- 邀請訊息 -->
        <div v-if="invitation.message" class="bg-white/5 rounded-lg p-4 mb-4">
          <p class="text-slate-300 text-sm">{{ invitation.message }}</p>
        </div>

        <!-- 時間資訊 -->
        <div class="flex items-center gap-4 text-sm text-slate-400 mb-4">
          <div class="flex items-center gap-1">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>收到於 {{ formatDate(invitation.created_at) }}</span>
          </div>
          <div v-if="invitation.status === 'PENDING'" class="flex items-center gap-1">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>過期於 {{ formatDate(invitation.expires_at) }}</span>
          </div>
          <div v-if="invitation.responded_at" class="flex items-center gap-1">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span>回應於 {{ formatDate(invitation.responded_at) }}</span>
          </div>
        </div>

        <!-- 操作按鈕 -->
        <div v-if="invitation.status === 'PENDING'" class="flex gap-3">
          <button
            @click="respondToInvitation(invitation.id, 'ACCEPT')"
            :disabled="respondingId === invitation.id"
            class="flex-1 px-4 py-3 bg-green-500/20 border border-green-500/50 text-green-400 rounded-xl hover:bg-green-500/30 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
          >
            <svg v-if="respondingId === invitation.id" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span v-else>✅</span>
            <span>接受邀請</span>
          </button>
          <button
            @click="respondToInvitation(invitation.id, 'DECLINE')"
            :disabled="respondingId === invitation.id"
            class="flex-1 px-4 py-3 bg-red-500/20 border border-red-500/50 text-red-400 rounded-xl hover:bg-red-500/30 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
          >
            <svg v-if="respondingId === invitation.id" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span v-else>❌</span>
            <span>婉拒</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 確認對話框 -->
    <GlobalAlert
      v-if="showConfirmDialog"
      :type="confirmType"
      :title="confirmTitle"
      :message="confirmMessage"
      :confirmText="confirmButtonText"
      cancelText="取消"
      @confirm="handleConfirm"
      @cancel="showConfirmDialog = false"
    />
  </div>
</template>

<script setup lang="ts">
import BaseGlassCard from '~/components/base/BaseGlassCard.vue'
import GlobalAlert from '~/components/GlobalAlert.vue'
import { alertError, alertSuccess } from '~/composables/useAlert'

definePageMeta({
  middleware: 'auth-teacher',
  layout: 'default',
})

const config = useRuntimeConfig()
const { success, error } = useToast()

const API_BASE = config.public.apiBase

// 狀態
const loading = ref(true)
const invitations = ref<any[]>([])
const activeTab = ref('PENDING')
const respondingId = ref<number | null>(null)

// 確認對話框
const showConfirmDialog = ref(false)
const confirmType = ref<'warning' | 'info'>('warning')
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmButtonText = ref('')
const pendingResponse = ref<{id: number, response: string} | null>(null)

// 標籤
const tabs = computed(() => [
  { label: '待處理', value: 'PENDING', count: invitations.value.filter(i => i.status === 'PENDING').length },
  { label: '已接受', value: 'ACCEPTED', count: invitations.value.filter(i => i.status === 'ACCEPTED').length },
  { label: '已婉拒', value: 'DECLINED', count: invitations.value.filter(i => i.status === 'DECLINED').length },
  { label: '已過期', value: 'EXPIRED', count: invitations.value.filter(i => i.status === 'EXPIRED').length },
  { label: '全部', value: 'ALL', count: invitations.value.length },
])

const filteredInvitations = computed(() => {
  if (activeTab.value === 'ALL') {
    return invitations.value
  }
  return invitations.value.filter(i => i.status === activeTab.value)
})

// 取得邀請列表
const fetchInvitations = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('teacher_token')
    const response = await fetch(`${API_BASE}/teacher/me/invitations`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (response.ok) {
      const data = await response.json()
      invitations.value = data.datas || []
    }
  } catch (err) {
    console.error('取得邀請列表失敗:', err)
    alertError('載入邀請失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}

// 回應邀請
const respondToInvitation = (id: number, response: 'ACCEPT' | 'DECLINE') => {
  pendingResponse.value = { id, response }
  
  if (response === 'ACCEPT') {
    confirmType.value = 'info'
    confirmTitle.value = '接受邀請'
    confirmMessage.value = '確定要接受這個邀請嗎？接受後將成為該中心的合作老師。'
    confirmButtonText.value = '確定接受'
  } else {
    confirmType.value = 'warning'
    confirmTitle.value = '婉拒邀請'
    confirmMessage.value = '確定要婉拒這個邀請嗎？'
    confirmButtonText.value = '確定婉拒'
  }
  
  showConfirmDialog.value = true
}

const handleConfirm = async () => {
  if (!pendingResponse.value) return
  
  const { id, response } = pendingResponse.value
  showConfirmDialog.value = false
  
  respondingId.value = id
  try {
    const token = localStorage.getItem('teacher_token')
    const res = await fetch(`${API_BASE}/teacher/me/invitations/respond`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        invitation_id: id,
        response: response,
      }),
    })

    if (res.ok) {
      const data = await res.json()
      success(data.message || '操作成功')
      await fetchInvitations()
    } else {
      const data = await res.json()
      alertError(data.message || '操作失敗')
    }
  } catch (err) {
    console.error('回應邀請失敗:', err)
    alertError('操作失敗，請稍後再試')
  } finally {
    respondingId.value = null
    pendingResponse.value = null
  }
}

// 樣式輔助函數
const getStatusClass = (status: string) => {
  switch (status) {
    case 'PENDING':
      return 'bg-warning-500/20 text-warning-500'
    case 'ACCEPTED':
      return 'bg-success-500/20 text-success-500'
    case 'DECLINED':
      return 'bg-critical-500/20 text-critical-500'
    case 'EXPIRED':
      return 'bg-slate-500/20 text-slate-400'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'PENDING':
      return '待處理'
    case 'ACCEPTED':
      return '已接受'
    case 'DECLINED':
      return '已婉拒'
    case 'EXPIRED':
      return '已過期'
    default:
      return status
  }
}

const getInviteTypeIcon = (type: string) => {
  switch (type) {
    case 'TALENT_POOL':
      return '⭐'
    case 'TEACHER':
      return '👨‍🏫'
    case 'MEMBER':
      return '🎫'
    default:
      return '📬'
  }
}

const getInviteTypeName = (type: string) => {
  switch (type) {
    case 'TALENT_POOL':
      return '人才庫邀請'
    case 'TEACHER':
      return '老師邀請'
    case 'MEMBER':
      return '會員邀請'
    default:
      return '邀請'
  }
}

const getInviteTypeBgClass = (type: string) => {
  switch (type) {
    case 'TALENT_POOL':
      return 'bg-yellow-500/20'
    case 'TEACHER':
      return 'bg-blue-500/20'
    case 'MEMBER':
      return 'bg-purple-500/20'
    default:
      return 'bg-slate-500/20'
  }
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

// 頁面載入時取得邀請列表
onMounted(() => {
  fetchInvitations()
})
</script>
