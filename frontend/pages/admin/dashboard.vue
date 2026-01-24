<template>
  <div class="h-full flex flex-col lg:flex-row gap-6">
    <ScheduleGrid
      class="flex-1 min-w-0"
      :view-mode="viewMode"
      :selected-resource-id="selectedResourceId"
      @select-resource="handleSelectResource"
    />
    <ScheduleResourcePanel
      class="lg:w-80 shrink-0"
      :view-mode="viewMode"
      @select-resource="handleSelectResource"
    />
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

// 視角模式：'all' | 'teacher' | 'room'
const viewMode = ref<'all' | 'teacher' | 'room'>('all')
// 選中的資源 ID（老師或教室）
const selectedResourceId = ref<number | null>(null)

const handleSelectResource = (resource: { type: 'teacher' | 'room', id: number } | null) => {
  if (!resource) {
    viewMode.value = 'all'
    selectedResourceId.value = null
  } else {
    viewMode.value = resource.type
    selectedResourceId.value = resource.id
  }
}

onMounted(() => {
  notificationStore.fetchNotifications()
})
</script>
