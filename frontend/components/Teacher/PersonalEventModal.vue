<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg sm:max-w-xl max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
        <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
          <h3 class="text-lg font-semibold text-slate-100">
            {{ isEditing ? '編輯個人行程' : '新增個人行程' }}
          </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="p-4 space-y-4">
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">標題</label>
          <input
            v-model="form.title"
            type="text"
            placeholder="例如：休息時間"
            class="input-field text-sm sm:text-base"
            required
          />
        </div>

        <div class="flex gap-4">
          <div class="flex-1">
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">開始日期</label>
            <input
              v-model="form.start_date"
              type="date"
              class="input-field text-sm sm:text-base w-full"
              required
            />
          </div>

          <div class="flex-1">
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">結束日期</label>
            <input
              v-model="form.end_date"
              type="date"
              class="input-field text-sm sm:text-base w-full"
              required
            />
          </div>
        </div>

        <div class="flex gap-4">
          <div class="flex-1">
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">開始時間</label>
            <input
              v-model="form.start_time"
              type="time"
              class="input-field text-sm sm:text-base w-full"
              required
            />
          </div>

          <div class="flex-1">
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">結束時間</label>
            <input
              v-model="form.end_time"
              type="time"
              class="input-field text-sm sm:text-base w-full"
              required
            />
          </div>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">重複設定</label>
          <select v-model="form.recurrence" class="input-field text-sm sm:text-base">
            <option value="NONE">不重複</option>
            <option value="DAILY">每天</option>
            <option value="WEEKLY">每週</option>
            <option value="BIWEEKLY">每兩週</option>
            <option value="MONTHLY">每月</option>
          </select>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">顏色標籤</label>
          <div class="flex gap-2 flex-wrap">
            <button
              v-for="color in colors"
              :key="color"
              type="button"
              @click="form.color_hex = color"
              class="w-10 h-10 rounded-xl transition-transform hover:scale-110"
              :class="form.color_hex === color ? 'ring-2 ring-white' : ''"
              :style="{ backgroundColor: color }"
            />
          </div>
        </div>

        <div class="flex gap-3 pt-2">
          <button
            type="button"
            @click="emit('close')"
            class="flex-1 glass-btn py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            {{ loading ? '儲存中...' : '儲存' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { toRef, toRaw, nextTick } from 'vue'
import { alertError, alertSuccess } from '~/composables/useAlert'
import { getTodayString, formatDateToString } from '~/composables/useTaiwanTime'

const props = defineProps<{
  editingEvent?: any
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

// 使用 toRef 保留響應式，避免解構破壞響應式追蹤
const editingEventRef = toRef(props, 'editingEvent')

const scheduleStore = useScheduleStore()
const loading = ref(false)

const isEditing = computed(() => !!editingEventRef.value)

// 追蹤 Modal 是否已經初始化過
const isModalInitialized = ref(false)

// 使用台灣時區的今天日期
const today = getTodayString()

const form = ref({
  title: '',
  start_date: today,
  start_time: '09:00',
  end_date: today,
  end_time: '10:00',
  recurrence: 'NONE' as 'NONE' | 'DAILY' | 'WEEKLY' | 'BIWEEKLY' | 'MONTHLY',
  color_hex: '#6366F1',
})

// Initialize form based on editing mode
const initializeForm = () => {
  // 使用 toRaw 確保取得原始物件，避免響應式追蹤問題
  const event = toRaw(editingEventRef.value)

  // 優先從頂層獲取資料，其次從 data 獲取
  // 支援兩種格式：頂層有資料或在 data 中
  const eventData = event?.data || event

  // 確保 event 存在且有必要的屬性
  if (eventData && eventData.start_at) {
    try {
      const startDate = new Date(eventData.start_at)
      const endDate = new Date(eventData.end_at)

      // 確保日期解析正確
      const startDateStr = !isNaN(startDate.getTime())
        ? formatDateToString(startDate)
        : today
      const endDateStr = !isNaN(endDate.getTime())
        ? formatDateToString(endDate)
        : today

      // 解析時間，確保格式正確
      const startTimeStr = !isNaN(startDate.getTime())
        ? startDate.toTimeString().slice(0, 5)
        : '09:00'
      const endTimeStr = !isNaN(endDate.getTime())
        ? endDate.toTimeString().slice(0, 5)
        : '10:00'

      // 優先使用 title，其次使用 offering_name（個人行程的標題可能在這兩個欄位）
      const title = eventData.title || eventData.offering_name || event?.offering_name || event?.title || ''

      // 優先使用 color_hex，其次使用 color
      const colorHex = eventData.color_hex || eventData.color || event?.color_hex || event?.color || '#6366F1'

      form.value = {
        title: title,
        start_date: startDateStr,
        start_time: startTimeStr,
        end_date: endDateStr,
        end_time: endTimeStr,
        recurrence: eventData.recurrence_rule?.type ||
                    eventData.recurrence_rule?.recurrence_type ||
                    eventData.recurrence_type ||
                    event?.recurrence_rule?.type ||
                    event?.recurrence_rule?.recurrence_type ||
                    event?.recurrence_type ||
                    'NONE',
        color_hex: colorHex,
      }
    } catch (error) {
      console.error('[PersonalEventModal] Error initializing form:', error)
      resetForm()
    }
  } else {
    resetForm()
  }
}

// 重置表單到預設值
const resetForm = () => {
  form.value = {
    title: '',
    start_date: today,
    start_time: '09:00',
    end_date: today,
    end_time: '10:00',
    recurrence: 'NONE' as 'NONE' | 'DAILY' | 'WEEKLY' | 'BIWEEKLY' | 'MONTHLY',
    color_hex: '#6366F1',
  }
}

// 監聽 editingEvent 變化，當有編輯資料時初始化表單
watch(editingEventRef, (newVal) => {
  if (newVal) {
    // 確保在下一個 tick 初始化表單
    nextTick(() => {
      initializeForm()
      isModalInitialized.value = true
    })
  }
}, { immediate: true })

// 監聽 Modal 開啟狀態，確保開啟時表單已初始化
watch(isEditing, (wasEditing) => {
  // 當編輯狀態為 true 且表單尚未初始化時，強制初始化
  if (wasEditing && !isModalInitialized.value) {
    nextTick(() => {
      initializeForm()
      isModalInitialized.value = true
    })
  }
}, { immediate: true })

const colors = [
  '#6366F1',
  '#A855F7',
  '#10B981',
  '#F43F5E',
  '#F59E0B',
  '#3B82F6',
  '#EC4899',
]

const formatDateTimeForApi = (date: string, time: string): string => {
  if (!date || !time) return ''
  // 使用台灣時區格式化 datetime 字串
  const [year, month, day] = date.split('-').map(Number)
  const [hours, minutes] = time.split(':').map(Number)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${year}-${pad(month)}-${pad(day)}T${pad(hours)}:${pad(minutes)}:00+08:00`
}

const handleSubmit = async () => {
  loading.value = true

  try {
    const data = {
      title: form.value.title,
      start_at: formatDateTimeForApi(form.value.start_date, form.value.start_time),
      end_at: formatDateTimeForApi(form.value.end_date, form.value.end_time),
      color_hex: form.value.color_hex,
    }

    if (form.value.recurrence !== 'NONE') {
      (data as any).recurrence_rule = {
        type: form.value.recurrence,
        interval: 1,
      }
    }

    if (isEditing.value) {
      // Update existing event
      await scheduleStore.updatePersonalEvent(props.editingEvent.id, data)
    } else {
      // Create new event
      await scheduleStore.createPersonalEvent(data)
    }
    await scheduleStore.fetchPersonalEvents()
    await scheduleStore.fetchSchedule()
    await alertSuccess('儲存成功')
    emit('saved')
    emit('close')
  } catch (err: any) {
    console.error('Failed to save personal event:', err)
    // 顯示後端返回的錯誤訊息
    const errorMessage = err.message || '儲存失敗，請稍後再試'
    await alertError(errorMessage)
  } finally {
    loading.value = false
  }
}
</script>
