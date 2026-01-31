<template>
  <Teleport to="body">
    <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50" @click="handleClose"></div>
      <div class="relative glass-card w-full max-w-lg p-6" :class="theme.cardGradient">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold" :class="theme.titleClass">課堂筆記</h3>
          <button @click="handleClose" class="p-1 rounded-lg hover:bg-white/10 transition-colors">
            <svg class="w-5 h-5" :class="theme.subtitleClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="mb-4 p-3 rounded-lg" :class="theme.dayHeaderClass">
          <p class="font-medium" :class="theme.itemTitleClass">{{ scheduleItem?.title }}</p>
          <p class="text-sm" :class="theme.subtitleClass">
            {{ formatDate(scheduleItem?.date) }} {{ scheduleItem?.start_time }} - {{ scheduleItem?.end_time }}
          </p>
          <p v-if="scheduleItem?.center_name" class="text-sm" :class="theme.centerClass">
            {{ scheduleItem.center_name }}
          </p>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2" :class="theme.subtitleClass">教學筆記</label>
            <textarea
              v-model="form.content"
              rows="4"
              class="w-full p-3 rounded-lg border resize-none focus:outline-none focus:ring-2 focus:ring-primary-500/50"
              :class="[theme.cardClass, theme.titleClass]"
              placeholder="記錄本次上課內容、進度、學生表現..."
            ></textarea>
          </div>

          <div>
            <label class="block text-sm font-medium mb-2" :class="theme.subtitleClass">備課筆記</label>
            <textarea
              v-model="form.prepNote"
              rows="3"
              class="w-full p-3 rounded-lg border resize-none focus:outline-none focus:ring-2 focus:ring-primary-500/50"
              :class="[theme.cardClass, theme.titleClass]"
              placeholder="下次上課準備事項、教材..."
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
import type { ScheduleItem, SessionNote } from '~/types'

const { isOpen, scheduleItem } = defineProps<{
  isOpen: boolean
  scheduleItem: ScheduleItem | null
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const scheduleStore = useScheduleStore()
const isSaving = ref(false)

const form = reactive({
  content: '',
  prepNote: '',
})

const theme = computed(() => {
  const exportThemes = [
    { id: 'domeMouse', cardGradient: 'bg-gradient-to-br from-[#fdfbf7] via-[#faf8f5] to-[#f7f5f0]', titleClass: 'text-[#5a524d]', subtitleClass: 'text-[#8a7e75]', dayHeaderClass: 'bg-[#f5f1ed] border-b border-[#e9e4dc]', itemTitleClass: 'text-[#4a4540]', centerClass: 'text-[#8a7e75]', cardClass: 'bg-[#f8f6f5]', buttonClass: 'hover:bg-[#f0ede8] text-[#5a524d]' },
    { id: 'dustyRose', cardGradient: 'bg-gradient-to-br from-[#faf5f5] via-[#f8f2f2] to-[#f5eeee]', titleClass: 'text-[#6b5555]', subtitleClass: 'text-[#9a8888]', dayHeaderClass: 'bg-[#f5f0ee] border-b border-[#e5dada]', itemTitleClass: 'text-[#5a4a4a]', centerClass: 'text-[#9a8888]', cardClass: 'bg-[#f8f5f5]', buttonClass: 'hover:bg-[#f0ebe9] text-[#6b5555]' },
    { id: 'sageGreen', cardGradient: 'bg-gradient-to-br from-[#f7faf7] via-[#f5f7f2] to-[#f2f5ee]', titleClass: 'text-[#5a6b55]', subtitleClass: 'text-[#8a9e88]', dayHeaderClass: 'bg-[#ebf0eb] border-b border-[#dce5dc]', itemTitleClass: 'text-[#4a5a45]', centerClass: 'text-[#8a9e88]', cardClass: 'bg-[#f5f7f5]', buttonClass: 'hover:bg-[#ebf0eb] text-[#5a6b55]' },
    { id: 'mutedBlue', cardGradient: 'bg-gradient-to-br from-[#f8fafb] via-[#f5f7f9] to-[#f2f5f7]', titleClass: 'text-[#565d6b]', subtitleClass: 'text-[#7a8a99]', dayHeaderClass: 'bg-[#e9eff2] border-b border-[#dce2e8]', itemTitleClass: 'text-[#464d5a]', centerClass: 'text-[#7a8a99]', cardClass: 'bg-[#f5f7f9]', buttonClass: 'hover:bg-[#e9eff2] text-[#565d6b]' },
    { id: 'warmBeige', cardGradient: 'bg-gradient-to-br from-[#fdfbf7] via-[#faf8f5] to-[#f7f5f0]', titleClass: 'text-[#5c5650]', subtitleClass: 'text-[#8a847a]', dayHeaderClass: 'bg-[#f5f1eb] border-b border-[#e9e4dc]', itemTitleClass: 'text-[#4a4540]', centerClass: 'text-[#8a847a]', cardClass: 'bg-[#faf8f5]', buttonClass: 'hover:bg-[#f5f1eb] text-[#5c5650]' },
    { id: 'lavender', cardGradient: 'bg-gradient-to-br from-[#faf9fa] via-[#f8f7f8] to-[#f5f5f5]', titleClass: 'text-[#5d595f]', subtitleClass: 'text-[#8a8088]', dayHeaderClass: 'bg-[#f2eef2] border-b border-[#e8e5e8]', itemTitleClass: 'text-[#4d494f]', centerClass: 'text-[#8a8088]', cardClass: 'bg-[#f8f7f8]', buttonClass: 'hover:bg-[#f2eef2] text-[#5d595f]' },
    { id: 'warmGrey', cardGradient: 'bg-gradient-to-br from-[#fafafa] via-[#f8f8f8] to-[#f5f5f5]', titleClass: 'text-[#555555]', subtitleClass: 'text-[#888888]', dayHeaderClass: 'bg-[#f2f2f2] border-b border-[#e8e8e8]', itemTitleClass: 'text-[#454545]', centerClass: 'text-[#888888]', cardClass: 'bg-[#f8f8f8]', buttonClass: 'hover:bg-[#f2f2f2] text-[#555555]' },
  ]
  return exportThemes[0]
})

const loadNote = async () => {
  // 優先使用 rule_id，其次從 data 取得
  const ruleId = scheduleItem?.rule_id || scheduleItem?.data?.id
  if (!ruleId || !scheduleItem?.date) return

  const note = await scheduleStore.fetchSessionNote(ruleId, scheduleItem.date)
  if (note) {
    form.content = note.content || ''
    form.prepNote = note.prep_note || ''
  } else {
    form.content = ''
    form.prepNote = ''
  }
}

watch(() => isOpen, (isOpen) => {
  if (isOpen && scheduleItem) {
    loadNote()
  }
})

const handleClose = () => {
  emit('close')
}

const handleSave = async () => {
  // 優先使用 rule_id，其次從 data 取得
  const ruleId = scheduleItem?.rule_id || scheduleItem?.data?.id
  if (!ruleId || !scheduleItem?.date) return

  isSaving.value = true
  try {
    await scheduleStore.saveSessionNote(
      ruleId,
      scheduleItem.date,
      form.content,
      form.prepNote
    )
    emit('saved')
    handleClose()
  } catch (error) {
    console.error('Failed to save note:', error)
  } finally {
    isSaving.value = false
  }
}
</script>
