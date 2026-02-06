<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" @click.self="close">
      <div class="glass-card w-full max-w-lg">
        <div class="flex items-center justify-between p-4 border-b border-white/10">
          <h3 class="text-lg font-semibold text-white">編輯老師評分</h3>
          <button @click="close" class="p-2 rounded-lg hover:bg-white/10">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div v-if="teacher" class="p-4 space-y-4">
          <!-- 老師資訊 -->
          <div class="flex items-center gap-4 p-4 rounded-xl bg-white/5">
            <div class="w-12 h-12 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
              <span class="text-white font-medium">{{ teacher.name?.charAt(0) || '?' }}</span>
            </div>
            <div>
              <h4 class="text-white font-medium">{{ teacher.name }}</h4>
              <p class="text-sm text-slate-400">{{ teacher.email }}</p>
            </div>
          </div>

          <!-- 評分 -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">評分 (0-5 星)</label>
            <div class="flex items-center gap-1">
              <button
                v-for="star in 5"
                :key="star"
                @click="form.rating = star"
                class="p-1 rounded-lg transition-colors"
                :class="form.rating >= star ? 'bg-primary-500/20' : 'hover:bg-white/5'"
              >
                <svg
                  class="w-8 h-8"
                  :class="form.rating >= star ? 'text-warning-500' : 'text-slate-600'"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              </button>
              <button
                @click="form.rating = 0"
                class="ml-2 px-3 py-1 rounded-lg text-sm transition-colors"
                :class="form.rating === 0 ? 'bg-critical-500/20 text-critical-500' : 'bg-white/5 text-slate-400 hover:bg-white/10'"
              >
                清除
              </button>
            </div>
            <p class="mt-2 text-sm" :class="form.rating > 0 ? 'text-warning-400' : 'text-slate-500'">
              {{ ratingLabels[form.rating] }}
            </p>
          </div>

          <!-- 內部備註 -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">內部備註</label>
            <textarea
              v-model="form.internal_note"
              rows="4"
              placeholder="記錄老師的表現特點、專長領域、合作經驗等..."
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500 resize-none"
            ></textarea>
            <p class="mt-1 text-xs text-slate-500">
              此備註僅供內部管理使用，影響智慧媒合的評分權重
            </p>
          </div>

          <div class="flex gap-3 pt-4">
            <button
              v-if="teacher.note?.id"
              @click="handleDelete"
              :disabled="saving"
              class="px-4 py-2 rounded-lg bg-critical-500/20 text-critical-500 hover:bg-critical-500/30 transition-colors"
            >
              刪除評分
            </button>
            <button @click="close" class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors">
              取消
            </button>
            <button @click="handleSave" :disabled="saving" class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50">
              {{ saving ? '儲存中...' : '儲存' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
const props = defineProps({
  show: Boolean,
  teacher: Object as PropType<any>
})

const emit = defineEmits(['close', 'saved', 'deleted'])

const api = useApi()
const notificationUI = useNotification()
const { confirm: alertConfirm } = useAlert()
const { invalidate } = useResourceCache()

const saving = ref(false)
const form = reactive({
  rating: 0,
  internal_note: ''
})

const ratingLabels: Record<number, string> = {
  0: '未評分',
  1: '需改進',
  2: '一般',
  3: '良好',
  4: '優良',
  5: '優秀'
}

watch(() => props.teacher, (newTeacher) => {
  if (newTeacher) {
    form.rating = newTeacher.note?.rating || 0
    form.internal_note = newTeacher.note?.internal_note || ''
  }
}, { immediate: true })

const close = () => {
  emit('close')
}

const handleSave = async () => {
  if (!props.teacher) return

  saving.value = true
  try {
    await api.put(`/admin/teachers/${props.teacher.id}/note`, {
      rating: form.rating,
      internal_note: form.internal_note
    })

    notificationUI.success('評分已儲存')
    invalidate('teachers')
    emit('saved')
    close()
  } catch (error) {
    console.error('Failed to save note:', error)
    notificationUI.error('儲存評分失敗')
  } finally {
    saving.value = false
  }
}

const handleDelete = async () => {
  if (!props.teacher || !await alertConfirm('確定要刪除這位老師的評分嗎？')) return

  saving.value = true
  try {
    await api.delete(`/admin/teachers/${props.teacher.id}/note`)
    notificationUI.success('評分已刪除')
    invalidate('teachers')
    emit('deleted')
    close()
  } catch (error) {
    console.error('Failed to delete note:', error)
    notificationUI.error('刪除評分失敗')
  } finally {
    saving.value = false
  }
}
</script>
