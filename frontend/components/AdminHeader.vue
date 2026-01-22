<template>
   <header class="glass sticky top-0 z-40 border-b border-white/10 bg-slate-900 relative">
     <div class="flex items-center justify-between px-6 py-3">
       <div class="flex items-center gap-6">
         <h1 class="text-xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">
           TimeLedger
         </h1>
          <nav class="flex items-center gap-4">
            <NuxtLink
              to="/admin/dashboard"
              class="text-slate-300 hover:text-primary-500 transition-colors font-medium"
              :class="{ 'text-primary-500': isActive('/admin/dashboard') }"
            >
              排課表
            </NuxtLink>
            <NuxtLink
              to="/admin/schedules"
              class="text-slate-300 hover:text-primary-500 transition-colors font-medium"
              :class="{ 'text-primary-500': isActive('/admin/schedules') }"
            >
              課程時段
            </NuxtLink>
            <NuxtLink
              to="/admin/templates"
              class="text-slate-300 hover:text-primary-500 transition-colors font-medium"
              :class="{ 'text-primary-500': isActive('/admin/templates') }"
            >
              課表模板
            </NuxtLink>
            <NuxtLink
              to="/admin/matching"
              class="text-slate-300 hover:text-primary-500 transition-colors font-medium"
              :class="{ 'text-primary-500': isActive('/admin/matching') }"
            >
              智慧媒合
            </NuxtLink>
            <NuxtLink
              to="/admin/approval"
              class="text-slate-300 hover:text-primary-500 transition-colors font-medium"
              :class="{ 'text-primary-500': isActive('/admin/approval') }"
            >
              審核中心
            </NuxtLink>
            <NuxtLink
              to="/admin/resources"
              class="text-slate-300 hover:text-primary-500 transition-colors font-medium"
              :class="{ 'text-primary-500': isActive('/admin/resources') }"
            >
              資源管理
            </NuxtLink>
          </nav>
       </div>

       <div class="flex items-center gap-4">
         <button
           @click="notificationUI.toggle()"
           class="relative p-2 rounded-lg hover:bg-white/10 transition-colors"
         >
           <svg class="w-6 h-6 text-slate-100" fill="none" stroke="currentColor" viewBox="0 0 24 24">
             <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
           </svg>
           <span
             v-if="notificationStore.unreadCount > 0"
             class="absolute -top-1 -right-1 w-5 h-5 rounded-full bg-critical-500 text-white text-xs flex items-center justify-center"
           >
             {{ notificationStore.unreadCount > 9 ? '9+' : notificationStore.unreadCount }}
           </span>
         </button>

         <div class="flex items-center gap-3">
           <div class="text-right">
             <p class="text-sm font-medium text-slate-100">
               {{ authStore.user?.name }}
             </p>
             <p class="text-xs text-slate-400">
               {{ (authStore.user as any)?.role }}
             </p>
           </div>
           <button
             @click="logout"
             class="p-2 rounded-lg hover:bg-critical-500/20 transition-colors"
           >
             <svg class="w-5 h-5 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
               <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H3m4 4h9a2 2 0 002-2V6a2 2 0 00-2-2H8a2 2 0 00-2 2v9a2 2 0 002 2h9" />
             </svg>
           </button>
         </div>
       </div>
     </div>
   </header>
 </template>

<script setup lang="ts">
const route = useRoute()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const notificationUI = useNotification()
const router = useRouter()

const isActive = (path: string): boolean => {
  return route.path === path
}

const logout = () => {
  if (confirm('確定要登出嗎？')) {
    authStore.logout()
    router.push('/')
  }
}

onMounted(() => {
  notificationStore.fetchNotifications()
})
</script>
