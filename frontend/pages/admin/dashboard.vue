<template>
  <div class="h-full flex flex-col lg:flex-row gap-6">
    <!-- 週曆視圖 -->
    <template v-if="viewMode === 'calendar'">
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
    </template>

    <!-- 矩陣視圖 -->
    <template v-else>
      <ScheduleMatrixView class="flex-1 min-w-0" />
      <ScheduleResourcePanel
        class="lg:w-80 shrink-0"
        :view-mode="resourcePanelViewMode"
        @select-resource="handleSelectResource"
      />
    </template>
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

// 視角模式：'calendar' | 'teacher_matrix' | 'room_matrix'
const viewMode = ref<'calendar' | 'teacher_matrix' | 'room_matrix'>('calendar')
// 選中的資源 ID（老師或教室）
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
    viewMode.value = resource.type === 'teacher' ? 'teacher_matrix' : 'room_matrix'
    selectedResourceId.value = resource.id
  }
}

onMounted(() => {
  notificationStore.fetchNotifications()
})
</script>
