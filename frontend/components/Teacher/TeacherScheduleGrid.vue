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

// 調試：監控 schedules prop 的變化
watch(() => props.schedules, (newSchedules) => {
  console.log('[TeacherScheduleGrid] schedules prop changed:', {
    count: newSchedules?.length || 0,
    firstItem: newSchedules?.[0] || null
  })
}, { immediate: true })

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
  console.log('[TeacherScheduleGrid] handleWeekStartChange:', date)
  emit('update:weekStart', date)
}

// ============================================
// 監聽週起始日期變化（同步到內部狀態）
// ============================================

watch(() => props.weekStart, (newVal) => {
  if (newVal) {
    console.log('[TeacherScheduleGrid] weekStart changed:', newVal)
  }
}, { immediate: true })

// 調試：組件掛載時記錄狀態
onMounted(() => {
  console.log('[TeacherScheduleGrid] mounted with schedules:', {
    count: props.schedules?.length || 0
  })
})
</script>
