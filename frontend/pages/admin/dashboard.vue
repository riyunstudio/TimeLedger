<template>
  <div class="h-full flex flex-col lg:flex-row gap-6">
    <!-- 時間軸週曆 -->
    <ScheduleTimelineView
      class="flex-1 min-w-0"
      v-model:view-mode="viewMode"
      v-model:selected-resource-id="selectedResourceId"
    />
    <ScheduleResourcePanel
      class="lg:w-80 shrink-0"
      :view-mode="resourcePanelViewMode"
      @select-resource="handleSelectResource"
    />
  </div>

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

// 視圖模式：'all' | 'teacher' | 'room'
const viewMode = ref<'all' | 'teacher' | 'room'>('all')
// 選中的資源 ID
const selectedResourceId = ref<number | null>(null)

// 資源面板的視角模式
const resourcePanelViewMode = computed(() => {
  if (viewMode.value === 'teacher') return 'teacher'
  if (viewMode.value === 'room') return 'room'
  return 'offering'
})

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
