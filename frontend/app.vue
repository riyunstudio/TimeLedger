<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <ClientOnly>
    <GlobalAlert />
  </ClientOnly>
</template>

<script setup lang="ts">
import GlobalAlert from '~/components/GlobalAlert.vue'

const authStore = useAuthStore()

onMounted(() => {
  console.log('[App] mounted')
  if (process.client) {
    authStore.initFromStorage()
  }
})
</script>
