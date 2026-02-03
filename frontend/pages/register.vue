<template>
  <div class="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
    <!-- 背景裝飾 -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-primary-500/20 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-indigo-500/20 rounded-full blur-3xl"></div>
    </div>

    <div class="relative w-full max-w-md">
      <!-- 載入中 -->
      <div v-if="loading" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="inline-block w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full animate-spin mb-4"></div>
        <p class="text-slate-300">正在處理中...</p>
      </div>

      <!-- 錯誤訊息 -->
      <div v-else-if="error" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-critical-500/20 flex items-center justify-center">
          <svg class="w-8 h-8 text-critical-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        <h2 class="text-xl font-bold text-white mb-2">發生錯誤</h2>
        <p class="text-slate-400 mb-6">{{ error }}</p>
        <button
          @click="goHome"
          class="px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors"
        >
          返回首頁
        </button>
      </div>

      <!-- LINE 登入按鈕 -->
      <div v-else-if="!hasLineUserId" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20">
        <div class="text-center mb-6">
          <h1 class="text-2xl font-bold text-white mb-1">加入 TimeLedger</h1>
          <p class="text-primary-400 font-medium">開始管理您的教學課表</p>
        </div>

        <button
          @click="lineLogin"
          class="w-full py-4 bg-[#06C755] hover:bg-[#05B54A] text-white font-medium rounded-xl transition-colors flex items-center justify-center gap-3"
        >
          <svg class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
          </svg>
          透過 LINE 快速註冊
        </button>
      </div>

      <!-- 註冊表單 -->
      <div v-else class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20">
        <div class="text-center mb-6">
          <h1 class="text-2xl font-bold text-white mb-1">填寫基本資料</h1>
          <p class="text-primary-400 font-medium">即將完成您的帳號建立</p>
        </div>

        <form @submit.prevent="handleRegister" class="space-y-4">
          <div>
            <label class="block text-slate-300 text-sm mb-2">姓名</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white focus:border-primary-500 transition-colors"
            />
          </div>
          <div>
            <label class="block text-slate-300 text-sm mb-2">Email</label>
            <input
              v-model="form.email"
              type="email"
              required
              class="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white focus:border-primary-500 transition-colors"
            />
          </div>
          <button
            type="submit"
            :disabled="submitting"
            class="w-full py-4 bg-primary-500 hover:bg-primary-600 text-white font-medium rounded-xl transition-colors disabled:opacity-50"
          >
            {{ submitting ? '註冊中...' : '提交註冊' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref('')
const hasLineUserId = ref(false)
const submitting = ref(false)

const form = ref({
  line_user_id: '',
  name: '',
  email: '',
})

const init = () => {
  const urlParams = new URLSearchParams(window.location.search)
  const lineUserId = urlParams.get('line_user_id')
  if (lineUserId) {
    form.value.line_user_id = lineUserId
    hasLineUserId.value = true
    window.history.replaceState({}, '', window.location.pathname)
  }
}

const lineLogin = () => {
  const liffUrl = `https://liff.line.me/${config.public.liffId}/teacher/login?redirect=${encodeURIComponent(window.location.href)}`
  window.location.href = liffUrl
}

const handleRegister = async () => {
  submitting.value = true
  try {
    const response = await fetch(`${config.public.apiBase}/teacher/public/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value),
    })
    const data = await response.json()
    if (!response.ok) throw new Error(data.message)
    
    if (data.datas?.token) {
      authStore.login({ token: data.datas.token, teacher: data.datas.teacher })
      router.push('/teacher/dashboard')
    }
  } catch (err: any) {
    error.value = err.message
  } finally {
    submitting.value = false
  }
}

const goHome = () => router.push('/')

onMounted(init)
</script>
