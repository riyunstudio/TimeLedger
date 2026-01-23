<template>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
    <ScheduleResourcePanel class="col-span-1" />
    <ScheduleGrid class="col-span-2" />
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
