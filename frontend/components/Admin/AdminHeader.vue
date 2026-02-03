<template>
  <header class="glass sticky top-0 z-40 border-b border-white/10">
    <div class="flex items-center justify-between px-4 py-3 max-w-7xl mx-auto">
      <div class="flex items-center gap-3">
        <h1 class="text-xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">TimeLedger</h1>
      </div>

      <div class="flex items-center gap-2">
        <!-- Notification bell -->
        <button
          @click="notificationUI.open()"
          class="relative p-2 rounded-lg hover:bg-white/10 transition-colors"
        >
          <Icon icon="bell" size="lg" />
          <span
            v-if="notificationStore.unreadCount > 0"
            class="absolute -top-1 -right-1 w-5 h-5 bg-critical-500 text-white text-xs rounded-full flex items-center justify-center"
          >
            {{ notificationStore.unreadCount > 9 ? '9+' : notificationStore.unreadCount }}
          </span>
        </button>

        <button
          @click="handleLogout"
          class="flex items-center gap-2 px-3 py-2 rounded-lg text-slate-300 hover:text-white hover:bg-white/10 transition-colors"
        >
          <Icon icon="logout" size="lg" />
          <span class="hidden sm:inline">登出</span>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { alertConfirm } from '~/composables/useAlert'

const router = useRouter()
const authStore = useAuthStore()
const notificationUI = useNotification()
const notificationStore = useNotificationStore()

const handleLogout = async () => {
  if (await alertConfirm('確定要登出嗎？')) {
    try {
      // 呼叫後端登出 API（將 Token 加入黑名單）
      const api = useApi()
      await api.post('/auth/logout', {})
    } catch (error) {
      // 即使後端呼叫失敗，仍然清除本地狀態
      console.error('Logout API failed:', error)
    }
    
    authStore.logout()
    router.push('/admin/login')
  }
}

// Fetch notifications on mount
onMounted(() => {
  notificationStore.fetchNotifications()
  // 同時獲取未讀數量以確保徽章顯示正確
  notificationStore.fetchUnreadCount()
})
</script>
