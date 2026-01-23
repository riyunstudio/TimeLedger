<template>
  <div class="h-full flex flex-col lg:flex-row gap-6">
    <ScheduleGrid class="flex-1 min-w-0" />
    <ScheduleResourcePanel class="lg:w-80 shrink-0" />
  </div>

  <ScheduleDetailPanel />

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
 definePageMeta({
   middleware: 'auth-admin',
   layout: 'admin',
 })

 const notificationStore = useNotificationStore()
const notificationUI = useNotification()

onMounted(() => {
  notificationStore.fetchNotifications()
})
</script>
