<template>
  <div class="min-h-screen bg-slate-900 pb-20 md:pb-0">
    <TeacherHeader v-if="showHeader" />

    <main class="p-4 md:p-6">
      <slot />
    </main>

    <nav v-if="showNav" class="fixed bottom-0 left-0 right-0 bg-slate-800/95 backdrop-blur-md border-t border-white/10 md:hidden">
      <div class="flex justify-around items-center h-16">
        <NuxtLink
          to="/teacher/dashboard"
          class="flex flex-col items-center justify-center flex-1 h-full text-slate-400 hover:text-primary-400 transition-colors"
          active-class="text-primary-500"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span class="text-xs mt-1">課表</span>
        </NuxtLink>

        <NuxtLink
          to="/teacher/profile"
          class="flex flex-col items-center justify-center flex-1 h-full text-slate-400 hover:text-primary-400 transition-colors"
          active-class="text-primary-500"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          <span class="text-xs mt-1">個人</span>
        </NuxtLink>

        <button
          class="flex flex-col items-center justify-center flex-1 h-full text-slate-400 hover:text-primary-400 transition-colors"
          @click="showNotifications = true"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
          <span class="text-xs mt-1">通知</span>
        </button>
      </div>
    </nav>
  </div>
</template>

<script setup lang="ts">
const showHeader = ref(true)
const showNav = ref(true)
const showNotifications = ref(false)
</script>
