<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-md animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10">
        <h3 class="text-lg font-semibold text-slate-100">邀請老師</h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="p-4 space-y-4">
        <div>
          <label class="block text-slate-300 mb-2 font-medium">老師 Email</label>
          <input
            v-model="form.email"
            type="email"
            placeholder="teacher@example.com"
            class="input-field"
            required
          />
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium">角色</label>
          <select v-model="form.role" class="input-field">
            <option value="TEACHER">老師</option>
            <option value="SUBSTITUTE">代課老師</option>
          </select>
        </div>

        <div v-if="inviteLink" class="p-3 rounded-xl bg-success-500/10 border border-success-500/20">
          <p class="text-sm text-success-500 font-medium mb-1">邀請連結已建立！</p>
          <input
            :value="inviteLink"
            readonly
            class="input-field text-xs bg-slate-800"
            @click="$event.target?.select()"
          />
          <p class="text-xs text-slate-400 mt-1">連結將在 72 小時後過期</p>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium">邀請訊息</label>
          <textarea
            v-model="form.message"
            placeholder="歡迎加入我們的中心..."
            rows="3"
            class="input-field resize-none"
          />
        </div>

        <div class="flex gap-3">
          <button
            @click="emit('close')"
            class="flex-1 glass-btn py-3 rounded-xl font-medium"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="flex-1 btn-primary"
          >
            {{ loading ? '發送中...' : '發送邀請' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  close: []
  invited: []
}>()

const loading = ref(false)
const inviteLink = ref('')
const { getCenterId } = useCenterId()
const form = ref({
  email: '',
  role: 'TEACHER' as 'TEACHER' | 'SUBSTITUTE',
  message: '',
})

const handleSubmit = async () => {
  loading.value = true

  try {
    const api = useApi()
    const centerId = getCenterId()
    const response = await api.post<any>(`/admin/invitations`, {
      email: form.value.email,
      role: form.value.role,
      message: form.value.message,
    })

    if (response && response.token) {
      inviteLink.value = `${window.location.origin}/register?token=${response.token}&center=${centerId}`
    } else {
      emit('invited')
      emit('close')
      alert('邀請已發送！')
    }
  } catch (error) {
    console.error('Failed to send invitation:', error)
    alert('發送失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}
</script>
