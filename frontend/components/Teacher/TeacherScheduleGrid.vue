<template>
  <SchedulingScheduleGrid
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
    :week-start="props.weekStart"
    @update:week-start="handleWeekStartChange"
    @select-date="handleDateSelect"
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
import { ref, watch, computed, onMounted } from 'vue'

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
  'select-date': [date: Date]
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

// 處理日期選擇 - 直接跳轉到指定日期
const handleDateSelect = (date: Date) => {
  emit('select-date', date)
}

// ============================================
// 監聽週起始日期變化（同步到內部狀態）
// ============================================

watch(() => props.weekStart, (newVal) => {
  if (newVal) {
    // 週起始日期變化時，同步到內部狀態
  }
}, { immediate: true })
</script>
