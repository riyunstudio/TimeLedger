<template>
  <div class="min-h-screen bg-slate-900">
    <AdminHeader />

    <main class="h-[calc(100vh-64px)] overflow-auto">
      <div class="p-4 md:p-6">
        <slot />
      </div>
    </main>

    <ToastNotification ref="toastRef" />
    <NotificationDropdown
      v-if="notificationUI.show.value"
      @close="notificationUI.close()"
    />
  </div>
</template>

<script setup lang="ts">
import { registerToast } from '~/composables/useToast'

const toastRef = ref<any>(null)
const notificationUI = useNotification()

onMounted(() => {
  if (toastRef.value) {
    registerToast(toastRef.value)
  }
})
</script>
