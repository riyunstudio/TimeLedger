<template>
  <div class="fixed inset-0 z-[200] pointer-events-none">
    <div class="absolute inset-0 bg-black/60 backdrop-blur-sm pointer-events-auto" @click="sidebarStore.close()"></div>
    <div class="absolute left-0 top-0 bottom-0 w-72 bg-slate-900 border-r border-white/10 shadow-2xl pointer-events-auto animate-slide-in">
      <div class="p-4 border-b border-white/10">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-slate-100">選單</h2>
          <button @click="sidebarStore.close()" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
            <BaseIcon icon="close" size="lg" />
          </button>
        </div>
        
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
            <span class="text-white font-semibold">{{ authStore.user?.name?.charAt(0) || 'T' }}</span>
          </div>
          <div>
            <p class="font-medium text-slate-100">{{ authStore.user?.name }}</p>
            <p class="text-xs text-slate-400">老師</p>
          </div>
        </div>
      </div>

      <nav class="p-4 space-y-2">
        <NuxtLink
          to="/teacher/dashboard"
          class="flex items-center gap-3 p-3 rounded-lg hover:bg-white/5 transition-colors relative"
          active-class="bg-primary-500/20 !text-primary-500"
          @click="sidebarStore.close()"
        >
          <BaseIcon icon="calendar" class="w-5 h-5" />
          <span>我的課表</span>
        </NuxtLink>

        <NuxtLink
          to="/teacher/profile"
          class="flex items-center gap-3 p-3 rounded-lg hover:bg-white/5 transition-colors"
          active-class="bg-primary-500/20 !text-primary-500"
          @click="sidebarStore.close()"
        >
          <BaseIcon icon="user" class="w-5 h-5" />
          <span>個人檔案</span>
        </NuxtLink>

        <NuxtLink
          to="/teacher/exceptions"
          class="flex items-center gap-3 p-3 rounded-lg hover:bg-white/5 transition-colors relative"
          active-class="bg-primary-500/20 !text-primary-500"
          @click="sidebarStore.close()"
        >
          <BaseIcon icon="warning" class="w-5 h-5" />
          <span>例外申請</span>
          <span
            v-if="pendingExceptions > 0"
            class="absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 rounded-full bg-warning-500 text-white text-xs flex items-center justify-center"
          >
            {{ pendingExceptions > 9 ? '9+' : pendingExceptions }}
          </span>
        </NuxtLink>

        <NuxtLink
          to="/teacher/export"
          class="flex items-center gap-3 p-3 rounded-lg hover:bg-white/5 transition-colors"
          active-class="bg-primary-500/20 !text-primary-500"
          @click="sidebarStore.close()"
        >
          <BaseIcon icon="download" class="w-5 h-5" />
          <span>匯出課表</span>
        </NuxtLink>

        <NuxtLink
          to="/teacher/invitations"
          class="flex items-center gap-3 p-3 rounded-lg hover:bg-white/5 transition-colors relative"
          active-class="bg-primary-500/20 !text-primary-500"
          @click="sidebarStore.close()"
        >
          <BaseIcon icon="chat" class="w-5 h-5" />
          <span>邀請通知</span>
          <span
            v-if="pendingInvitations > 0"
            class="absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 rounded-full bg-primary-500 text-white text-xs flex items-center justify-center"
          >
            {{ pendingInvitations > 9 ? '9+' : pendingInvitations }}
          </span>
        </NuxtLink>

        <button
          @click="handleLogout"
          class="w-full flex items-center gap-3 p-3 rounded-lg hover:bg-critical-500/10 text-critical-500 transition-colors"
        >
          <BaseIcon icon="lock" class="w-5 h-5" />
          <span>登出</span>
        </button>
      </nav>

      <!-- 快速統計 -->
      <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-white/10 bg-slate-900/95">
        <div class="text-center">
          <p class="text-2xl font-bold text-white">{{ pendingExceptions }}</p>
          <p class="text-xs text-slate-400">待審核例外申請</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertConfirm } from '~/composables/useAlert'

const config = useRuntimeConfig()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const scheduleStore = useScheduleStore()
const sidebarStore = useSidebar()
const pendingInvitations = ref(0)

// 待審核申請數
const pendingExceptions = computed(() => {
  return scheduleStore.exceptions.filter(e => e.status === 'PENDING').length
})

// 取得待處理邀請數量
const fetchPendingInvitations = async () => {
  try {
    const token = localStorage.getItem('teacher_token')
    const response = await fetch(`${config.public.apiBase}/teacher/me/invitations/pending-count`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    if (response.ok) {
      const data = await response.json()
      pendingInvitations.value = data.datas?.count || 0
    }
  } catch (err) {
    console.error('取得待處理邀請數量失敗:', err)
  }
}

// 頁面載入時取得資料
onMounted(() => {
  fetchPendingInvitations()
})

const handleLogout = async () => {
  if (await alertConfirm('確定要登出嗎？')) {
    sidebarStore.close()
    
    try {
      // 呼叫後端登出 API（將 Token 加入黑名單）
      const api = useApi()
      await api.post('/auth/logout', {})
    } catch (error) {
      console.error('Logout API failed:', error)
    }
    
    authStore.logout()
    router.push('/')
  }
}
</script>

<style scoped>
@keyframes slide-in {
  from {
    transform: translateX(-100%);
  }
  to {
    transform: translateX(0);
  }
}

.animate-slide-in {
  animation: slide-in 0.3s ease-out;
}
</style>
