<template>
  <div class="min-h-screen bg-slate-900 flex flex-col">
    <AdminHeader />

    <main class="flex-1 p-6">
      <div class="grid grid-cols-3 gap-6 h-[calc(100vh-80px)]">
        <ScheduleResourcePanel class="col-span-1" />
        <ScheduleGrid class="col-span-2" />
      </div>
    </main>

    <ScheduleDetailPanel />

    <NotificationDropdown
      v-if="notificationUI.show.value"
      @close="notificationUI.close()"
    />
  </div>
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
