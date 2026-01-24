<template>
  <div class="h-full flex flex-col lg:flex-row gap-6">
    <!-- 矩陣視圖 -->
    <ScheduleMatrixView
      class="flex-1 min-w-0"
      v-model:resource-type="resourceType"
      @select-cell="handleSelectResource"
    />
    <ScheduleResourcePanel
      class="lg:w-80 shrink-0"
      :view-mode="resourcePanelViewMode"
      @select-resource="handleSelectResource"
    />
  </div>

  <ScheduleDetailPanel
    v-if="selectedCell"
    :time="selectedCell.time"
    :weekday="selectedCell.weekday"
    :schedule="selectedSchedule"
    @close="selectedCell = null"
  />

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

// 資源類型：'teacher' | 'room'
const resourceType = ref<'teacher' | 'room'>('teacher')

// 選中的資源（用於詳情顯示）
const selectedCell = ref<{ time: number; weekday: number; resource: any } | null>(null)
const selectedSchedule = ref<any>(null)

// 資源面板的視角模式
const resourcePanelViewMode = computed(() => {
  return resourceType.value
})

const handleSelectResource = (data: { resource: any; time: number; weekday: number } | null) => {
  if (data) {
    selectedCell.value = data
    // 從 scheduleMap 中獲取完整的 schedule 資料
    // 這部分由 ScheduleMatrixView 處理
  }
}

onMounted(() => {
  notificationStore.fetchNotifications()
})
</script>
