<template>
  <div class="fixed inset-0 z-[200] pointer-events-none">
    <div class="absolute inset-0 bg-black/60 backdrop-blur-sm pointer-events-auto" @click="sidebarStore.close()"></div>
    <div class="absolute left-0 top-0 bottom-0 w-72 bg-slate-900 border-r border-white/10 shadow-2xl pointer-events-auto animate-slide-in">
      <div class="p-4 border-b border-white/10">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-slate-100">選單</h2>
          <button @click="sidebarStore.close()" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
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
          class="flex items-center gap-3 p-3 rounded-lg hover:bg-white/5 transition-colors"
          :class="{ 'bg-primary-500/20 text-primary-500': route.path === '/teacher/dashboard' }"
          @click="sidebarStore.close()"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span>我的課表</span>
        </NuxtLink>

        <NuxtLink
          to="/teacher/profile"
          class="flex items-center gap-3 p-3 rounded-lg hover:bg-white/5 transition-colors"
          :class="{ 'bg-primary-500/20 text-primary-500': route.path === '/teacher/profile' }"
          @click="sidebarStore.close()"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          <span>個人檔案</span>
        </NuxtLink>

        <NuxtLink
          to="/teacher/exceptions"
          class="flex items-center gap-3 p-3 rounded-lg hover:bg-white/5 transition-colors"
          :class="{ 'bg-primary-500/20 text-primary-500': route.path === '/teacher/exceptions' }"
          @click="sidebarStore.close()"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <span>例外申請</span>
        </NuxtLink>

        <NuxtLink
          to="/teacher/export"
          class="flex items-center gap-3 p-3 rounded-lg hover:bg-white/5 transition-colors"
          :class="{ 'bg-primary-500/20 text-primary-500': route.path === '/teacher/export' }"
          @click="sidebarStore.close()"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          <span>匯出課表</span>
        </NuxtLink>

        <button
          @click="handleLogout"
          class="w-full flex items-center gap-3 p-3 rounded-lg hover:bg-critical-500/10 text-critical-500 transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          <span>登出</span>
        </button>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const sidebarStore = useSidebar()

const handleLogout = () => {
  if (confirm('確定要登出嗎？')) {
    sidebarStore.close()
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
