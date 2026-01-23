<template>
  <div class="min-h-screen bg-gradient-mesh">
    <NuxtPage />
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()

onMounted(() => {
  if (process.client) {
    authStore.initFromStorage()
  }
})
</script>
