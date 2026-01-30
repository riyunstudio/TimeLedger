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
          <svg class="w-5 h-5 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2 2 0 0118 22v-4.317l-1.405 1.405A2 2 0 0115 17h5l-1.405-1.405A2 2 0 0118 12v6a2 2 0 01-2 2h-2m-6 3h2m8 0h2M3 8h18M3 8v10a2 2 0 002 2h14a2 2 0 002-2V8z" />
          </svg>
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
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
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
    authStore.logout()
    router.push('/admin/login')
  }
}

// Fetch notifications on mount
onMounted(() => {
  notificationStore.fetchNotifications()
})
</script>
