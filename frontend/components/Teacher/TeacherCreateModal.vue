<template>
  <Teleport to="body">
    <div
      v-if="show"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div class="glass-card w-full max-w-md animate-spring">
        <div class="p-6 border-b border-white/10 flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-white">手動新增老師</h3>
            <p class="text-sm text-slate-400 mt-1">
              建立佔位老師，稍後可合併至正式 LINE 帳號
            </p>
          </div>
          <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm text-slate-400 mb-2">
              姓名 <span class="text-critical-500">*</span>
            </label>
            <input
              v-model="form.name"
              type="text"
              placeholder="請輸入老師姓名"
              class="input-field w-full"
              required
            />
          </div>

          <div>
            <label class="block text-sm text-slate-400 mb-2">
              Email <span class="text-slate-500">(選填)</span>
            </label>
            <input
              v-model="form.email"
              type="email"
              placeholder="teacher@example.com"
              class="input-field w-full"
            />
          </div>

          <div class="p-3 rounded-lg bg-warning-500/10 border border-warning-500/30">
            <p class="text-sm text-warning-400">
              建立後，老師可透過 LINE 登入並自動綁定此帳號。
            </p>
          </div>
        </div>

        <div class="p-6 border-t border-white/10 flex items-center gap-4">
          <button
            @click="emit('close')"
            class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors text-sm font-medium"
          >
            取消
          </button>
          <button
            @click="handleSubmit"
            :disabled="!isValid || loading"
            class="flex-1 btn-primary py-2 text-sm font-medium disabled:opacity-50"
          >
            <span v-if="loading" class="flex items-center justify-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              處理中...
            </span>
            <span v-else>新增老師</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  close: []
  created: []
}>()

const api = useApi()
const loading = ref(false)
const form = ref({
  name: '',
  email: '',
})

const isValid = computed(() => form.value.name.trim().length > 0)

const handleSubmit = async () => {
  if (!isValid.value) return

  loading.value = true
  try {
    await api.post('/admin/teachers/placeholder', {
      name: form.value.name.trim(),
      email: form.value.email.trim() || undefined,
    })

    emit('created')
    emit('close')
    form.value = { name: '', email: '' }
  } catch (error) {
    console.error('Failed to create placeholder teacher:', error)
    // 這裡通常會由 api client 處理錯誤提示，或由父組件處理
  } finally {
    loading.value = false
  }
}
</script>
