<template>
  <div class="space-y-4">
    <!-- 重複類型選擇 -->
    <div v-if="showTypeSelector">
      <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
        重複模式
      </label>
      <div class="grid grid-cols-3 sm:grid-cols-5 gap-2">
        <button
          v-for="type in recurrenceTypes"
          :key="type.value"
          type="button"
          @click="selectType(type.value)"
          class="px-3 py-2 rounded-lg text-sm font-medium transition-all text-center"
          :class="selectedType === type.value
            ? 'bg-primary-500 text-white'
            : 'bg-slate-700/50 text-slate-400 hover:bg-slate-700'"
        >
          {{ type.label }}
        </button>
      </div>
      <p v-if="typeDescription" class="text-xs text-slate-400 mt-2">
        {{ typeDescription }}
      </p>
    </div>

    <!-- 星期選擇（僅當選擇週循環時顯示） -->
    <div v-if="showWeekdaySelector">
      <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
        {{ weekdayLabel }}
      </label>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="day in weekDays"
          :key="day.value"
          type="button"
          @click="toggleWeekday(day.value)"
          class="px-3 py-2 rounded-lg text-sm font-medium transition-all"
          :class="selectedWeekdays.includes(day.value)
            ? 'bg-primary-500 text-white'
            : 'bg-slate-700/50 text-slate-400 hover:bg-slate-700'"
        >
          {{ day.name }}
        </button>
      </div>
      <HelpTooltip
        v-if="weekdayHelpText"
        class="mt-2"
        :title="weekdayLabel"
        :description="weekdayHelpText"
        :usage="weekdayUsageTips"
      />
      <span v-if="weekdayError" class="text-critical-500 text-xs mt-1 block">
        {{ weekdayError }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import HelpTooltip from '~/components/base/HelpTooltip.vue'
interface Props {
  /** 選中的星期陣列（1-7，代表週一到週日） */
  modelValue: number[]
  /** 重複類型：NONE, DAILY, WEEKLY, MONTHLY, CUSTOM */
  recurrenceType?: 'NONE' | 'DAILY' | 'WEEKLY' | 'MONTHLY' | 'CUSTOM'
  /** 是否顯示重複類型選擇器 */
  showTypeSelector?: boolean
  /** 是否顯示星期選擇器 */
  showWeekdaySelector?: boolean
  /** 星期標籤文字 */
  weekdayLabel?: string
  /** 星期說明文字 */
  weekdayHelpText?: string
  /** 星期使用提示 */
  weekdayUsageTips?: string[]
  /** 星期驗證錯誤訊息 */
  weekdayError?: string
  /** 隱藏的星期值（用於 DAILY 等類型） */
  hiddenWeekdays?: number[]
}

const props = withDefaults(defineProps<Props>(), {
  recurrenceType: 'WEEKLY',
  showTypeSelector: false,
  showWeekdaySelector: true,
  weekdayLabel: '重複星期',
  weekdayHelpText: '選擇此排課規則適用的星期幾。',
  weekdayUsageTips: () => ['可選擇多個星期', '形成每週重複的排課'],
  modelValue: () => [],
  hiddenWeekdays: () => [],
})

const emit = defineEmits<{
  'update:modelValue': [weekdays: number[]]
  'update:recurrenceType': [type: 'NONE' | 'DAILY' | 'WEEKLY' | 'MONTHLY' | 'CUSTOM']
}>()

/** 重複類型選項 */
const recurrenceTypes = [
  { value: 'NONE' as const, label: '不重複', description: '僅單次排課' },
  { value: 'DAILY' as const, label: '每日', description: '每天重複' },
  { value: 'WEEKLY' as const, label: '每週', description: '每週特定日期重複' },
  { value: 'MONTHLY' as const, label: '每月', description: '每月相同日期重複' },
  { value: 'CUSTOM' as const, label: '自訂', description: '自訂循環規則' },
]

/** 星期選項 */
const weekDays = [
  { value: 1, name: '週一' },
  { value: 2, name: '週二' },
  { value: 3, name: '週三' },
  { value: 4, name: '週四' },
  { value: 5, name: '週五' },
  { value: 6, name: '週六' },
  { value: 7, name: '週日' },
]

/** 選中的重複類型 */
const selectedType = ref<'NONE' | 'DAILY' | 'WEEKLY' | 'MONTHLY' | 'CUSTOM'>(props.recurrenceType)

/** 選中的星期 */
const selectedWeekdays = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

/** 是否顯示星期選擇器 */
const showWeekdaySelector = computed(() => {
  if (!props.showWeekdaySelector) return false
  // 週循環類型需要選擇星期
  return ['WEEKLY', 'CUSTOM'].includes(selectedType.value)
})

/** 是否顯示類型選擇器 */
const showTypeSelector = computed(() => props.showTypeSelector)

/** 當前類型的說明 */
const typeDescription = computed(() => {
  const type = recurrenceTypes.find(t => t.value === selectedType.value)
  return type?.description
})

/** 選擇重複類型 */
const selectType = (type: 'NONE' | 'DAILY' | 'WEEKLY' | 'MONTHLY' | 'CUSTOM') => {
  selectedType.value = type
  emit('update:recurrenceType', type)

  // 根據類型自動設定星期
  if (type === 'DAILY') {
    // 每日：選擇所有星期
    emit('update:modelValue', [1, 2, 3, 4, 5, 6, 7])
  } else if (type === 'NONE') {
    // 不重複：只選擇今天
    const today = new Date().getDay() || 7
    emit('update:modelValue', [today])
  }
  // WEEKLY 和 CUSTOM 需要手動選擇星期，不自動設定
}

/** 切換星期選擇 */
const toggleWeekday = (day: number) => {
  const current = selectedWeekdays.value
  const index = current.indexOf(day)
  if (index === -1) {
    emit('update:modelValue', [...current, day].sort((a, b) => a - b))
  } else {
    emit('update:modelValue', current.filter((d) => d !== day))
  }
}

// 監聽外部 recurrenceType 變化
watch(() => props.recurrenceType, (newType) => {
  if (newType) {
    selectedType.value = newType
  }
})
</script>
