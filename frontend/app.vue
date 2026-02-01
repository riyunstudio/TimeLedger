<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <ClientOnly>
    <GlobalAlert />
  </ClientOnly>
</template>

<script setup lang="ts">
import GlobalAlert from '~/components/base/GlobalAlert.vue'

// Pinia 持久化外掛會自動從 storage 恢復狀態
// 不需要額外的初始化
</script>
