<template>
  <Teleport to="body">
    <div
      v-if="show"
      class="fixed inset-0 z-[250] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div class="glass-card w-full max-w-md animate-spring">
        <div class="p-6 border-b border-white/10">
          <h3 class="text-lg font-semibold text-white">合併帳號</h3>
          <p class="text-sm text-slate-400 mt-1">
            將佔位老師的資料合併至已綁定老師帳號
          </p>
        </div>

        <div v-if="teacher" class="p-6 space-y-4">
          <!-- 來源老師（佔位） -->
          <div class="p-4 rounded-lg bg-warning-500/10 border border-warning-500/30">
            <p class="text-sm text-warning-400 mb-2">來源（將被合併）</p>
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-warning-500/20 flex items-center justify-center">
                <span class="text-warning-400 font-medium">
                  {{ teacher.name.charAt(0) }}
                </span>
              </div>
              <div>
                <p class="text-white font-medium">{{ teacher.name }}</p>
                <p class="text-sm text-slate-400">{{ teacher.email || '無 Email' }}</p>
              </div>
            </div>
            <div class="mt-2 flex items-center gap-1.5 text-warning-400 text-xs">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <span>佔位老師，合併後將被刪除</span>
            </div>
          </div>

          <!-- 目標老師選擇 -->
          <div>
            <label class="block text-sm text-slate-400 mb-2">
              選擇目標老師 <span class="text-critical-500">*</span>
            </label>
            <select
              v-model="targetTeacherId"
              class="w-full px-3 py-2 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500"
            >
              <option value="">請選擇目標老師</option>
              <option
                v-for="bt in boundTeachers"
                :key="bt.id"
                :value="bt.id"
                :disabled="bt.id === teacher.id"
              >
                {{ bt.name }} ({{ bt.email || '無 Email' }})
              </option>
            </select>
            <p class="mt-1 text-xs text-slate-500">
              選擇要接收資料的已綁定老師帳號
            </p>
          </div>

          <!-- 合併說明 -->
          <div class="p-3 rounded-lg bg-primary-500/10 border border-primary-500/30">
            <p class="text-sm text-primary-400">
              合併將會：
            </p>
            <ul class="mt-2 text-xs text-slate-400 space-y-1">
              <li>- 遷移所有課表、例外記錄、課程筆記</li>
              <li>- 遷移私人行程和技能證照</li>
              <li>- 刪除來源佔位老師帳號</li>
              <li>- 此操作無法復原</li>
            </ul>
          </div>
        </div>

        <div class="p-6 border-t border-white/10 flex items-center gap-4">
          <button
            @click="emit('close')"
            class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
          >
            取消
          </button>
          <button
            @click="handleMerge"
            :disabled="!targetTeacherId || merging"
            class="flex-1 px-4 py-2 rounded-lg bg-warning-500/30 border border-warning-500 text-warning-400 hover:bg-warning-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="merging" class="flex items-center justify-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              處理中...
            </span>
            <span v-else>確認合併</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
const props = defineProps<{
  show: boolean
  teacher: any | null
  boundTeachers: any[]
}>()

const emit = defineEmits<{
  close: []
  merged: []
}>()

const { confirm: alertConfirm, error: alertError } = useAlert()
const api = useApi()

const targetTeacherId = ref('')
const merging = ref(false)

const handleMerge = async () => {
  if (!targetTeacherId.value || !props.teacher) return

  const confirmed = await alertConfirm(
    `確定要將「${props.teacher.name}」合併至目標老師嗎？\n\n此操作將遷移所有資料並刪除佔位帳號，無法復原。`
  )

  if (!confirmed) return

  merging.value = true
  try {
    await api.post('/admin/teachers/merge', {
      source_teacher_id: props.teacher.id,
      target_teacher_id: parseInt(targetTeacherId.value),
    })

    emit('merged')
    emit('close')
  } catch (err) {
    console.error('Failed to merge teacher:', err)
    await alertError('合併失敗，請稍後再試')
  } finally {
    merging.value = false
  }
}

watch(() => props.show, (newVal) => {
  if (newVal) {
    targetTeacherId.value = ''
  }
})
</script>
