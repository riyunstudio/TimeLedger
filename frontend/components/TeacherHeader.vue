<template>
  <header class="glass sticky top-0 z-40 border-b border-white/10">
    <div class="flex items-center justify-between px-4 py-3 w-full">
      <button
        @click="sidebarStore.toggle()"
        class="p-2 rounded-lg hover:bg-white/10 transition-colors"
      >
        <svg class="w-6 h-6 text-slate-100" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>

      <div class="flex items-center gap-3">
        <NuxtLink
          to="/teacher/export"
          class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          title="匯出課表"
        >
          <svg class="w-6 h-6 text-slate-100" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
        </NuxtLink>

        <button
          @click="notificationUI.toggle()"
          class="relative p-2 rounded-lg hover:bg-white/10 transition-colors"
        >
          <svg class="w-6 h-6 text-slate-100" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
          <span
            v-if="notificationDataStore.unreadCount > 0"
            class="absolute -top-1 -right-1 w-5 h-5 rounded-full bg-critical-500 text-white text-xs flex items-center justify-center"
          >
            {{ notificationDataStore.unreadCount > 9 ? '9+' : notificationDataStore.unreadCount }}
          </span>
        </button>

        <NuxtLink
          to="/teacher/profile"
          class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0"
        >
          <span class="text-white font-semibold">
            {{ authStore.user?.name?.charAt(0) || 'T' }}
          </span>
        </NuxtLink>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const notificationDataStore = useNotificationStore()
const notificationUI = useNotification()
const sidebarStore = useSidebar()

onMounted(() => {
  notificationDataStore.fetchNotifications()
})
</script>
