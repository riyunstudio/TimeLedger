<template>
  <div class="min-h-screen bg-slate-900 pb-20 md:pb-0">
    <TeacherHeader v-if="showHeader" />

    <main class="p-4 md:p-6 w-full">
      <slot />
    </main>

    <nav v-if="showNav" class="fixed bottom-0 left-0 right-0 bg-slate-800/95 backdrop-blur-md border-t border-white/10 md:hidden">
      <div class="flex justify-around items-center h-16">
        <NuxtLink
          to="/teacher/dashboard"
          class="flex flex-col items-center justify-center flex-1 h-full h-full text-slate-400 hover:text-primary-400 transition-colors"
          active-class="text-primary-500"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4l-7-7 7-7z" />
          </svg>
          <span class="text-xs mt-1">課表</span>
        </NuxtLink>

        <NuxtLink
          to="/teacher/profile"
          class="flex flex-col items-center justify-center flex-1 h-full h-full text-slate-400 hover:text-primary-400 transition-colors"
          active-class="text-primary-500"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 011-4 1 1 0V6M6 18L4 6l2 2m6 6v12a2 2 0 012 2-2 2 0 012 2l4 4m6 6v12a2 2 0 012 2 2 2 0012 2-4 4m6 6v12a2 2 0 012 2 2 2 0012 2-4 4m6 6v12a2 2 0 012 2 2 2 0012 2" />
          </svg>
          <span class="text-xs mt-1">個人</span>
        </NuxtLink>

        <button
          class="flex flex-col items-center justify-center flex-1 h-full h-full text-slate-400 hover:text-primary-400 transition-colors"
          @click="notificationUI.open()"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2 2 0 0118 22v-4.317l-1.405 1.405A2 2 0 0115 17h5l-1.405-1.405A2 2 0 0118 12v6a2 2 0 01-2 2h-2m-6 3h2m8 0h2M3 8h18M3 8v10a2 2 0 002 2h14a2 2 0 002-2V8z" />
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
const notificationUI = useNotification()
</script>
