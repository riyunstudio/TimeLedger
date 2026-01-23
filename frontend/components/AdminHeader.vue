<template>
  <header class="glass sticky top-0 z-40 border-b border-white/10">
    <div class="flex items-center justify-between px-4 py-3 max-w-7xl mx-auto">
      <div class="flex items-center gap-3">
        <h1 class="text-xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">TimeLedger</h1>
      </div>

      <nav class="flex items-center gap-6">
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
          課課模板
        </NuxtLink>
        <NuxtLink
          to="/admin/resources"
          class="text-slate-300 hover:text-primary-400 transition-colors-btn py-2"
          active-class="text-primary-500"
        >
          資源管理
        </NuxtLink>
        <NuxtLink
          to="/admin/approval"
          class="text-slate-300 hover:text-primary-400 transition-colors-btn py-2"
          active-class="text-primary-500"
        >
          審核中心
        </NuxtLink>
        <NuxtLink
          to="/admin/matching"
          class="text-slate-300 hover:text-primary-400 transition-colors-btn py-2"
          active-class="text-primary-500"
        >
          智慧媒合
        </NuxtLink>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
// TODO: Implement notificationUI composable if needed
// const notificationUI = useNotification()
</script>
