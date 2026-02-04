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
        <p class="text-slate-300">正在初始化...</p>
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
          @click="retryInit"
          class="px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors"
        >
          重新整理
        </button>
      </div>

      <!-- 未登入 LINE -->
      <div v-else-if="!hasLineUserId" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20">
        <div class="text-center mb-6">
          <h1 class="text-2xl font-bold text-white mb-1">加入 TimeLedger</h1>
          <p class="text-primary-400 font-medium">開始管理您的教學課表</p>
        </div>

        <div class="bg-white/5 rounded-xl p-4 mb-6">
          <p class="text-slate-300 text-sm">
            請使用 LINE 帳號註冊，開始管理您的教學課表。
          </p>
        </div>

        <button
          @click="lineLogin"
          class="w-full py-4 bg-[#06C755] hover:bg-[#05B54A] text-white font-medium rounded-xl transition-colors flex items-center justify-center gap-3"
        >
          <svg class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
          </svg>
          LINE 註冊
        </button>
        <p class="text-center text-slate-500 text-sm mt-3">
          點擊上方按鈕透過 LINE 註冊
        </p>
      </div>

      <!-- 註冊中 -->
      <div v-else-if="submitting" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="inline-block w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full animate-spin mb-4"></div>
        <p class="text-slate-300">正在註冊...</p>
      </div>

      <!-- 註冊失敗 -->
      <div v-else-if="registerError" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-critical-500/20 flex items-center justify-center">
          <svg class="w-8 h-8 text-critical-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        <h2 class="text-xl font-bold text-white mb-2">註冊失敗</h2>
        <p class="text-slate-400 mb-6">{{ registerError }}</p>
        <button
          @click="retryRegister"
          class="px-6 py-3 bg-white/10 text-white rounded-xl hover:bg-white/20 transition-colors"
        >
          重新註冊
        </button>
      </div>

      <!-- 註冊成功 -->
      <div v-else-if="registerSuccess" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-success-500/20 flex items-center justify-center">
          <svg class="w-8 h-8 text-success-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 class="text-xl font-bold text-white mb-2">註冊成功！</h2>
        <p class="text-slate-400">正在跳轉到後台...</p>
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
            class="w-full py-4 bg-primary-500 hover:bg-primary-600 text-white font-medium rounded-xl transition-colors"
          >
            提交註冊
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

const config = useRuntimeConfig()
const router = useRouter()
const authStore = useAuthStore()
const { $liff } = useNuxtApp()

// 狀態
const loading = ref(true)
const error = ref('')
const hasLineUserId = ref(false)
const submitting = ref(false)
const registerError = ref('')
const registerSuccess = ref(false)

// LINE User ID
const lineUserId = ref('')

// 表單
const form = ref({
  name: '',
  email: '',
})

// 初始化 LIFF
const initLiff = async () => {
  loading.value = true
  error.value = ''

  try {
    // 檢查 LIFF 是否已初始化
    if (!$liff) {
      throw new Error('LIFF 尚未初始化，請重新整理頁面')
    }

    // 檢查是否已登入 LINE
    const isLoggedIn = $liff.isLoggedIn()

    if (isLoggedIn) {
      // 已登入 LINE，取得用戶資訊
      const profile = await $liff.getProfile()
      lineUserId.value = profile.userId
      hasLineUserId.value = true
    } else {
      // 未登入 LINE，需要先登入
      hasLineUserId.value = false
    }
  } catch (err: any) {
    error.value = err.message || '初始化失敗，請重新整理頁面'
  } finally {
    loading.value = false
  }
}

// LINE 登入
const lineLogin = async () => {
  if (!$liff) {
    error.value = 'LIFF 尚未初始化，請重新整理頁面'
    return
  }

  try {
    // 使用 LIFF SDK 登入
    $liff.login()
  } catch (err: any) {
    error.value = err.message || 'LINE 登入失敗，請稍後再試'
  }
}

// 處理註冊
const handleRegister = async () => {
  submitting.value = true
  registerError.value = ''

  try {
    // 取得 LINE Access Token
    const accessToken = $liff.getAccessToken()
    if (!accessToken) {
      throw new Error('無法取得 LINE Access Token，請重新登入')
    }

    const response = await $fetch(`${config.public.apiBase}/teacher/public/register`, {
      method: 'POST',
      body: {
        line_user_id: lineUserId.value,
        access_token: accessToken,
        name: form.value.name,
        email: form.value.email,
      },
    })

    const responseData = (response as any).data || (response as any).datas
    const responseCode = (response as any).code

    if (responseCode === 0 && responseData?.token) {
      // 註冊成功
      authStore.login({
        token: responseData.token,
        refresh_token: responseData.refresh_token || '',
        teacher: responseData.teacher,
      })

      registerSuccess.value = true

      setTimeout(() => {
        router.push('/teacher/dashboard')
      }, 1500)
    } else {
      registerError.value = (response as any)?.message || '註冊失敗，請稍後再試'
    }
  } catch (err: any) {
    console.error('Register error:', err)
    registerError.value = err.data?.message || err.message || '註冊失敗，請稍後再試'
  } finally {
    submitting.value = false
  }
}

// 重新嘗試初始化
const retryInit = () => {
  window.location.reload()
}

// 重新註冊
const retryRegister = () => {
  registerError.value = ''
  submitting.value = false
  hasLineUserId.value = false
  lineUserId.value = ''

  // 重新登入 LINE
  lineLogin()
}

onMounted(() => {
  // 檢查是否已登入
  const token = localStorage.getItem('teacher_token')
  if (token) {
    // 已登入，直接跳轉
    router.push('/teacher/dashboard')
    return
  }

  initLiff()
})
</script>
