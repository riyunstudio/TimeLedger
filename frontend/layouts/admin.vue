<template>
  <div class="min-h-screen bg-slate-900">
    <AdminHeader @toggle-sidebar="sidebarOpen = !sidebarOpen" />

    <div class="flex flex-col md:flex-row min-h-[calc(100vh-64px)]">
      <main class="flex-1 p-4 md:p-6 overflow-auto">
        <slot />
      </main>

      <AdminSidebar v-if="sidebarOpen" />
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth-admin'],
})

const sidebarOpen = ref(false)

onMounted(() => {
  // Default to open on desktop, closed on mobile
  sidebarOpen.value = window.innerWidth >= 768
})
</script>
