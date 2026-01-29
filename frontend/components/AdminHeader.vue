<template>
  <header class="glass sticky top-0 z-40 border-b border-white/10">
    <div class="flex items-center justify-between px-4 py-3 max-w-7xl mx-auto relative">
      <div class="flex items-center gap-3">
        <h1 class="text-xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">TimeLedger</h1>
      </div>

      <!-- Desktop nav -->
      <nav class="hidden md:flex items-center gap-6">
        <NuxtLink
          to="/admin/dashboard"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          課表
        </NuxtLink>
        <NuxtLink
          to="/admin/schedules"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          課程時段
        </NuxtLink>
        <NuxtLink
          to="/admin/templates"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          課表模板
        </NuxtLink>
        <NuxtLink
          to="/admin/resources"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          資源管理
        </NuxtLink>
        <NuxtLink
          to="/admin/approval"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          審核中心
        </NuxtLink>
        <NuxtLink
          to="/admin/matching"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          智慧媒合
        </NuxtLink>
        <NuxtLink
          to="/admin/holidays"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          例假日
        </NuxtLink>
        <NuxtLink
          to="/admin/invitations"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          邀請紀錄
        </NuxtLink>
        <NuxtLink
          to="/admin/teacher-ratings"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          老師評分
        </NuxtLink>
        <NuxtLink
          to="/admin/admin-list"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          管理員
        </NuxtLink>
        <NuxtLink
          to="/admin/settings"
          class="text-slate-300 hover:text-primary-400 transition-colors font-medium"
          active-class="text-primary-500"
        >
          設定
        </NuxtLink>
      </nav>

      <!-- Mobile nav -->
      <div v-if="mobileMenuOpen" class="md:hidden absolute top-full left-0 right-0 bg-slate-800/95 backdrop-blur-sm border-b border-white/10 z-50">
        <nav class="flex flex-col p-4 space-y-2">
          <NuxtLink
            to="/admin/dashboard"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            課表
          </NuxtLink>
          <NuxtLink
            to="/admin/schedules"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            課程時段
          </NuxtLink>
          <NuxtLink
            to="/admin/templates"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            課表模板
          </NuxtLink>
          <NuxtLink
            to="/admin/resources"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            資源管理
          </NuxtLink>
          <NuxtLink
            to="/admin/approval"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            審核中心
          </NuxtLink>
          <NuxtLink
            to="/admin/matching"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            智慧媒合
          </NuxtLink>
          <NuxtLink
            to="/admin/holidays"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            例假日
          </NuxtLink>
          <NuxtLink
            to="/admin/invitations"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            邀請紀錄
          </NuxtLink>
          <NuxtLink
            to="/admin/teacher-ratings"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            老師評分
          </NuxtLink>
          <NuxtLink
            to="/admin/admin-list"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            管理員
          </NuxtLink>
          <NuxtLink
            to="/admin/settings"
            class="block px-4 py-2 rounded-lg text-slate-300 hover:text-primary-400 hover:bg-white/10 transition-colors font-medium"
            active-class="text-primary-500 bg-primary-500/10"
            @click="mobileMenuOpen = false"
          >
            設定
          </NuxtLink>
        </nav>
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

        <!-- Mobile menu button -->
        <button
          @click="mobileMenuOpen = !mobileMenuOpen"
          class="md:hidden p-2 rounded-lg hover:bg-white/10 transition-colors"
        >
          <svg class="w-5 h-5 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
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

const mobileMenuOpen = ref(false)

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
