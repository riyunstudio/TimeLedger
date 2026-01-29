<template>
  <div class="p-4 md:p-6 max-w-7xl mx-auto">
    <div class="mb-6 md:mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
        邀請紀錄
      </h1>
      <p class="text-slate-400 text-sm md:text-base">
        查看邀請老師的歷史記錄與處理結果
      </p>
    </div>

    <!-- 統計卡片 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <div class="bg-white/5 rounded-xl p-4 border border-white/10">
        <p class="text-slate-400 text-sm">待處理</p>
        <p class="text-2xl font-bold text-warning-500">{{ stats.pending }}</p>
      </div>
      <div class="bg-white/5 rounded-xl p-4 border border-white/10">
        <p class="text-slate-400 text-sm">已接受</p>
        <p class="text-2xl font-bold text-success-500">{{ stats.accepted }}</p>
      </div>
      <div class="bg-white/5 rounded-xl p-4 border border-white/10">
        <p class="text-slate-400 text-sm">已婉拒</p>
        <p class="text-2xl font-bold text-critical-500">{{ stats.declined }}</p>
      </div>
      <div class="bg-white/5 rounded-xl p-4 border border-white/10">
        <p class="text-slate-400 text-sm">已過期</p>
        <p class="text-2xl font-bold text-slate-400">{{ stats.expired }}</p>
      </div>
    </div>

    <!-- 篩選器 -->
    <div class="flex flex-wrap gap-3 mb-6">
      <select
        v-model="filters.status"
        class="px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-slate-300 focus:outline-none focus:border-primary-500"
      >
        <option value="">全部狀態</option>
        <option value="PENDING">待處理</option>
        <option value="ACCEPTED">已接受</option>
        <option value="DECLINED">已婉拒</option>
        <option value="EXPIRED">已過期</option>
      </select>

      <input
        v-model="filters.search"
        type="text"
        placeholder="搜尋 Email 或姓名..."
        class="px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-slate-300 focus:outline-none focus:border-primary-500 flex-1 min-w-[200px]"
      />

      <button
        @click="refreshData"
        class="px-4 py-2 bg-primary-500/20 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/30 transition-colors"
      >
        重新整理
      </button>
    </div>

    <!-- 邀請列表 -->
    <div class="bg-white/5 rounded-xl border border-white/10 overflow-hidden">
      <div v-if="loading" class="p-8 text-center">
        <div class="inline-block w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-slate-400 mt-2">載入中...</p>
      </div>

      <div v-else-if="filteredInvitations.length === 0" class="p-8 text-center text-slate-400">
        暫無邀請記錄
      </div>

      <table v-else class="w-full">
        <thead class="bg-white/5">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">Email</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">狀態</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">邀請時間</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">回應時間</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="invitation in filteredInvitations" :key="invitation.id" class="hover:bg-white/5">
            <td class="px-4 py-3 text-slate-300">{{ invitation.email }}</td>
            <td class="px-4 py-3">
              <span
                class="px-2 py-1 rounded-full text-xs font-medium"
                :class="{
                  'bg-warning-500/20 text-warning-500': invitation.status === 'PENDING',
                  'bg-success-500/20 text-success-500': invitation.status === 'ACCEPTED',
                  'bg-critical-500/20 text-critical-500': invitation.status === 'DECLINED',
                  'bg-slate-500/20 text-slate-400': invitation.status === 'EXPIRED',
                }"
              >
                {{ statusText(invitation.status) }}
              </span>
            </td>
            <td class="px-4 py-3 text-slate-400 text-sm">
              {{ formatDate(invitation.created_at) }}
            </td>
            <td class="px-4 py-3 text-slate-400 text-sm">
              {{ invitation.responded_at ? formatDate(invitation.responded_at) : '-' }}
            </td>
            <td class="px-4 py-3">
              <button
                v-if="invitation.status === 'PENDING'"
                @click="resendInvitation(invitation)"
                class="text-primary-500 hover:text-primary-400 text-sm"
              >
                重新傳送
              </button>
              <span v-else class="text-slate-500 text-sm">-</span>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- 分頁 -->
      <div v-if="pagination.totalPages > 1" class="px-4 py-3 bg-white/5 border-t border-white/10 flex items-center justify-between">
        <p class="text-slate-400 text-sm">
          第 {{ pagination.page }} 頁，共 {{ pagination.totalPages }} 頁
        </p>
        <div class="flex gap-2">
          <button
            @click="changePage(pagination.page - 1)"
            :disabled="pagination.page === 1"
            class="px-3 py-1 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            上一頁
          </button>
          <button
            @click="changePage(pagination.page + 1)"
            :disabled="pagination.page === pagination.totalPages"
            class="px-3 py-1 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            下一頁
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth-admin',
  layout: 'admin',
})

const config = useRuntimeConfig()
const authStore = useAuthStore()

// 取得當前登入管理員的中心 ID
const centerId = computed(() => {
  return authStore.user?.center_id || 1
})

interface Invitation {
  id: number
  email: string
  status: string
  created_at: string
  responded_at?: string
  expires_at: string
}

const loading = ref(true)
const invitations = ref<Invitation[]>([])
const stats = ref({
  pending: 0,
  accepted: 0,
  declined: 0,
  expired: 0,
})

const filters = ref({
  status: '',
  search: '',
})

const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  totalPages: 0,
})

const statusText = (status: string) => {
  const texts: Record<string, string> = {
    PENDING: '待處理',
    ACCEPTED: '已接受',
    DECLINED: '已婉拒',
    EXPIRED: '已過期',
  }
  return texts[status] || status
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const fetchStats = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId.value}/invitations/stats`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    if (response.ok) {
      const data = await response.json()
      stats.value = {
        pending: data.datas?.pending || 0,
        accepted: data.datas?.accepted || 0,
        declined: data.datas?.declined || 0,
        expired: data.datas?.expired || 0,
      }
    }
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

const fetchInvitations = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    if (filters.value.status) {
      params.append('status', filters.value.status)
    }

    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId.value}/invitations?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    if (response.ok) {
      const data = await response.json()
      invitations.value = data.datas?.data || []
      pagination.value.total = data.datas?.total || 0
      pagination.value.totalPages = Math.ceil(pagination.value.total / pagination.value.limit)
    }
  } catch (error) {
    console.error('Failed to fetch invitations:', error)
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  pagination.value.page = 1
  fetchStats()
  fetchInvitations()
}

const changePage = (page: number) => {
  pagination.value.page = page
  fetchInvitations()
}

const resendInvitation = async (invitation: Invitation) => {
  // TODO: 實作重新傳送邀請功能
  alert('重新傳送功能待實作')
}

const filteredInvitations = computed(() => {
  if (!filters.value.search) return invitations.value
  const search = filters.value.search.toLowerCase()
  return invitations.value.filter(inv =>
    inv.email.toLowerCase().includes(search)
  )
})

// 監聽篩選器變化
watch(filters, () => {
  pagination.value.page = 1
  fetchInvitations()
}, { deep: true })

onMounted(() => {
  refreshData()
})
</script>
