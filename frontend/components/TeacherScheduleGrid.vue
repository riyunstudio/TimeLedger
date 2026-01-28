<template>
  <ScheduleGrid
    mode="teacher"
    :api-endpoint="apiEndpoint"
    card-info-type="center"
    :show-matrix-view="false"
    :show-view-mode-selector="false"
    :show-create-button="false"
    :show-personal-event-button="true"
    :show-exception-button="true"
    :show-export-button="true"
    :show-help-tooltip="effectiveShowHelpTooltip"
    :view-mode="viewMode"
    :selected-resource-id="null"
    :schedules="props.schedules"
    @update:week-start="handleWeekStartChange"
    @select-schedule="$emit('select-schedule', $event)"
    @add-personal-event="$emit('add-personal-event')"
    @add-exception="$emit('add-exception')"
    @export="$emit('export')"
    @edit-personal-event="$emit('edit-personal-event', $event)"
    @delete-personal-event="$emit('delete-personal-event', $event)"
    @personal-event-note="$emit('personal-event-note', $event)"
  />
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  schedules: any[]
  weekStart?: Date
  showHelpTooltip?: boolean
}>()

// ============================================
// Emits 定義
// ============================================

const emit = defineEmits<{
  'update:weekStart': [value: Date]
  'select-schedule': [schedule: any]
  'add-personal-event': []
  'add-exception': []
  'export': []
  'edit-personal-event': [event: any]
  'delete-personal-event': [event: any]
  'personal-event-note': [event: any]
}>()

// ============================================
// 常量
// ============================================

const apiEndpoint = '/teacher/schedules'

// ============================================
// 狀態
// ============================================

const viewMode = ref<'calendar'>('calendar')

// 安全的 showHelpTooltip 計算屬性
const effectiveShowHelpTooltip = computed(() => props.showHelpTooltip ?? true)

// ============================================
// 事件處理
// ============================================

const handleWeekStartChange = (date: Date) => {
  emit('update:weekStart', date)
}

// ============================================
// 監聽週起始日期變化（同步到內部狀態）
// ============================================

watch(() => props.weekStart, (newVal) => {
  // 當外部傳入 weekStart 變化時，可以進行相應處理
  if (newVal) {
    // 可以在這裡重新獲取資料或進行其他同步操作
  }
}, { immediate: true })
</script>
