<template>
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
    >
      <div class="absolute inset-0 bg-black/50" @click="handleClose"></div>
      <div class="relative glass-card w-full max-w-lg p-6" :class="theme.cardGradient">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold" :class="theme.titleClass">行程備註</h3>
          <button @click="handleClose" class="p-1 rounded-lg hover:bg-white/10 transition-colors">
            <svg class="w-5 h-5" :class="theme.subtitleClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="mb-4 p-3 rounded-lg" :class="theme.dayHeaderClass">
          <p class="font-medium" :class="theme.itemTitleClass">{{ eventData?.title }}</p>
          <p class="text-sm" :class="theme.subtitleClass">
            {{ formatDate(eventData?.start_at) }}
          </p>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2" :class="theme.subtitleClass">備註內容</label>
            <textarea
              v-model="form.content"
              rows="4"
              class="w-full p-3 rounded-lg border resize-none focus:outline-none focus:ring-2 focus:ring-primary-500/50"
              :class="[theme.cardClass, theme.titleClass]"
              placeholder="記錄此行程的相關資訊..."
            ></textarea>
          </div>
        </div>

        <div class="flex justify-end gap-3 mt-6">
          <button
            @click="handleClose"
            class="px-4 py-2 rounded-lg transition-colors"
            :class="theme.buttonClass"
          >
            取消
          </button>
          <button
            @click="handleSave"
            :disabled="isSaving"
            class="px-4 py-2 rounded-lg transition-colors flex items-center gap-2"
            :class="isSaving ? 'opacity-50 cursor-not-allowed' : 'bg-primary-500 text-white hover:bg-primary-600'"
          >
            <svg v-if="isSaving" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
            {{ isSaving ? '儲存中...' : '儲存' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { formatDate } from '~/composables/useTaiwanTime'
import type { PersonalEvent } from '~/types'

interface Props {
  isOpen: boolean
  event: PersonalEvent | null
}

const { isOpen, event: eventProp } = defineProps<Props>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const isSaving = ref(false)

const eventData = computed(() => eventProp)

const form = reactive({
  content: '',
})

const theme = computed(() => {
  return {
    cardGradient: 'bg-gradient-to-br from-[#fdfbf7] via-[#faf8f5] to-[#f7f5f0]',
    titleClass: 'text-[#5a524d]',
    subtitleClass: 'text-[#8a7e75]',
    dayHeaderClass: 'bg-[#f5f1ed] border-b border-[#e9e4dc]',
    itemTitleClass: 'text-[#4a4540]',
    centerClass: 'text-[#8a7e75]',
    cardClass: 'bg-[#f8f6f5]',
    buttonClass: 'hover:bg-[#f0ede8] text-[#5a524d]',
  }
})

const loadNote = async () => {
  if (!eventData.value) return

  // 個人行程備註使用不同的 API 端點
  try {
    const api = useApi()
    const originalId = typeof eventData.value.id === 'string'
      ? parseInt(eventData.value.id.split('_')[0])
      : eventData.value.id

    const response = await api.get<{ code: number; message: string; datas: { content: string } }>(
      `/teacher/me/personal-events/${originalId}/note`
    )
    form.content = response.datas?.content || ''
  } catch (error) {
    // 如果沒有備註，返回空
    form.content = ''
  }
}

watch(() => isOpen, (isOpen) => {
  if (isOpen && eventData.value) {
    loadNote()
  }
})

const handleClose = () => {
  emit('close')
}

const handleSave = async () => {
  if (!eventData.value) return

  isSaving.value = true
  try {
    const api = useApi()
    const originalId = typeof eventData.value.id === 'string'
      ? parseInt(eventData.value.id.split('_')[0])
      : eventData.value.id

    await api.put(`/teacher/me/personal-events/${originalId}/note`, {
      content: form.content,
    })
    emit('saved')
    handleClose()
  } catch (error) {
    console.error('Failed to save note:', error)
  } finally {
    isSaving.value = false
  }
}
</script>
