<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg sm:max-w-xl max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          新增個人行程
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

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">開始日期</label>
            <input
              v-model="form.start_date"
              type="date"
              class="input-field text-sm sm:text-base"
              required
            />
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">開始時間</label>
            <input
              v-model="form.start_time"
              type="time"
              class="input-field text-sm sm:text-base"
              required
            />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">結束日期</label>
            <input
              v-model="form.end_date"
              type="date"
              class="input-field text-sm sm:text-base"
              required
            />
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">結束時間</label>
            <input
              v-model="form.end_time"
              type="time"
              class="input-field text-sm sm:text-base"
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
const emit = defineEmits<{
  close: []
}>()

const teacherStore = useTeacherStore()
const loading = ref(false)

const now = new Date()

const form = ref({
  title: '',
  start_date: now.toISOString().split('T')[0],
  start_time: now.toTimeString().slice(0, 5),
  end_date: now.toISOString().split('T')[0],
  end_time: new Date(now.getTime() + 60 * 60 * 1000).toTimeString().slice(0, 5),
  recurrence: 'NONE' as 'NONE' | 'DAILY' | 'WEEKLY' | 'BIWEEKLY' | 'MONTHLY',
  color_hex: '#6366F1',
})

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
  return new Date(`${date}T${time}:00`).toISOString()
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

    await teacherStore.createPersonalEvent(data)
    await teacherStore.fetchPersonalEvents()
    emit('close')
  } catch (error) {
    console.error('Failed to create personal event:', error)
    alert('儲存失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}
</script>
