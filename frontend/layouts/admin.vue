<template>
  <div class="min-h-screen bg-slate-900 flex flex-col md:flex-row">
    <AdminHeader />

    <main class="flex-1 p-4 md:p-6 overflow-auto">
      <slot />
    </main>

    <AdminSidebar />
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth-admin'],
})
</script>
