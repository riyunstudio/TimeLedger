<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
const authStore = useAuthStore()

onMounted(() => {
  if (process.client) {
    authStore.initFromStorage()
  }
})
</script>
