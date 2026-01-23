<template>
  <div class="min-h-screen bg-slate-900">
    <AdminHeader />

    <main class="h-[calc(100vh-64px)] overflow-auto">
      <div class="p-4 md:p-6">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth-admin'],
})
</script>
