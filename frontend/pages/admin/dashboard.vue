<template>
  <div class="h-full flex flex-col lg:flex-row gap-6">
    <!-- 排課網格（支援週曆/矩陣視圖） -->
    <ScheduleGrid
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

// 視圖模式：'calendar' | 'teacher_matrix' | 'room_matrix'
const viewMode = ref<'calendar' | 'teacher_matrix' | 'room_matrix'>('calendar')
// 選中的資源 ID
const selectedResourceId = ref<number | null>(null)

// 資源面板的視角模式
const resourcePanelViewMode = computed(() => {
  if (viewMode.value === 'teacher_matrix') return 'teacher'
  if (viewMode.value === 'room_matrix') return 'room'
  return 'offering'
})

const handleSelectResource = (resource: { type: 'teacher' | 'room', id: number } | null) => {
  if (!resource) {
    viewMode.value = 'calendar'
    selectedResourceId.value = null
  } else {
    if (resource.type === 'teacher') {
      viewMode.value = 'teacher_matrix'
    } else {
      viewMode.value = 'room_matrix'
    }
    selectedResourceId.value = resource.id
  }
}

onMounted(() => {
  notificationStore.fetchNotifications()
})
</script>
