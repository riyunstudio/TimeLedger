<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          求職設定
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-4 space-y-6">
        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="text-slate-300 font-medium">開放中心搜尋我的檔案</label>
            <button
              @click="form.is_open_to_hiring = !form.is_open_to_hiring"
              class="relative w-14 h-7 rounded-full transition-colors duration-300"
              :class="form.is_open_to_hiring ? 'bg-success-500' : 'bg-slate-700'"
            >
              <span
                class="absolute top-1 left-1 w-5 h-5 rounded-full bg-white transition-transform duration-300"
                :class="{ 'translate-x-7': form.is_open_to_hiring }"
              />
            </button>
          </div>
          <p class="text-sm text-slate-400">
            開啟後，中心可以透過智慧媒合系統找到您
          </p>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium">
            公開聯繫方式
          </label>
          <textarea
            v-model="form.public_contact_info"
            placeholder="例：LINE ID: @yourline 或 電話: 0912345678"
            rows="3"
            :disabled="!form.is_open_to_hiring"
            class="input-field resize-none"
            :class="{ 'opacity-50': !form.is_open_to_hiring }"
          />
          <p class="text-xs text-slate-500 mt-2">
            此資訊僅在您開啟求職時，對媒合您的中心可見
          </p>
        </div>

        <div class="flex items-start gap-3 p-4 rounded-xl bg-primary-500/10 border border-primary-500/30">
          <svg class="w-5 h-5 text-primary-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9 9h.01" />
          </svg>
          <p class="text-sm text-slate-400">
            開啟求職後，您的技能與證照將自動對媒合中心可見（唯讀）。您可以隨時關閉此設定。
          </p>
        </div>

        <div class="flex gap-3">
          <button
            @click="emit('close')"
            class="flex-1 glass-btn py-3 rounded-xl font-medium"
          >
            取消
          </button>
          <button
            @click="handleSubmit"
            :disabled="loading"
            class="flex-1 btn-primary"
          >
            {{ loading ? '儲存中...' : '儲存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertError } from '~/composables/useAlert'

const emit = defineEmits<{
  close: []
}>()

const authStore = useAuthStore()
const loading = ref(false)
const form = ref({
  is_open_to_hiring: false,
  public_contact_info: '',
})

// 載入最新資料
const loadData = async () => {
  try {
    const api = useApi()
    const response = await api.get<{ code: number; data: any }>('/teacher/me/profile')
    if (response.code === 0 && response.data) {
      form.value.is_open_to_hiring = response.data.is_open_to_hiring || false
      // public_contact_info 可能是 JSON 字串或物件
      const contactInfo = response.data.public_contact_info
      if (typeof contactInfo === 'string') {
        form.value.public_contact_info = contactInfo
      } else if (typeof contactInfo === 'object' && contactInfo !== null) {
        form.value.public_contact_info = contactInfo.content || JSON.stringify(contactInfo)
      } else {
        form.value.public_contact_info = ''
      }
    }
  } catch (error) {
    console.error('Failed to load hiring settings:', error)
  }
}

const handleSubmit = async () => {
  loading.value = true

  try {
    const api = useApi()
    await api.put('/teacher/me/profile', form.value)

    authStore.user = { ...authStore.user, ...form.value } as any
    localStorage.setItem('teacher_user', JSON.stringify(authStore.user))

    emit('close')
  } catch (error) {
    console.error('Failed to update hiring settings:', error)
    await alertError('儲存失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>
