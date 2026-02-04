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
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-primary-500/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-white mb-1">老師登入</h1>
          <p class="text-primary-400 font-medium">歡迎回來</p>
        </div>

        <div class="bg-white/5 rounded-xl p-4 mb-6">
          <p class="text-slate-300 text-sm">
            請使用 LINE 帳號登入，開始管理您的課表。
          </p>
        </div>

        <!-- LINE 登入按鈕 -->
        <button
          @click="lineLogin"
          class="w-full py-4 bg-[#06C755] hover:bg-[#05B54A] text-white font-medium rounded-xl transition-colors flex items-center justify-center gap-3"
        >
          <svg class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
          </svg>
          LINE 登入
        </button>
        <p class="text-center text-slate-500 text-sm mt-3">
          點擊上方按鈕透過 LINE 登入
        </p>
      </div>

      <!-- 登入中 -->
      <div v-else-if="loggingIn" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="inline-block w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full animate-spin mb-4"></div>
        <p class="text-slate-300">正在登入...</p>
      </div>

      <!-- 登入成功 -->
      <div v-else-if="loginSuccess" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-success-500/20 flex items-center justify-center">
          <svg class="w-8 h-8 text-success-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 class="text-xl font-bold text-white mb-2">登入成功！</h2>
        <p class="text-slate-400">正在跳轉到後台...</p>
      </div>

      <!-- 登入失敗 -->
      <div v-else-if="loginError" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-critical-500/20 flex items-center justify-center">
          <svg class="w-8 h-8 text-critical-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        <div class="text-center mb-6">
          <h2 class="text-xl font-bold text-white mb-2">登入失敗</h2>
          <p class="text-slate-400">{{ loginError }}</p>
        </div>
        <button
          @click="retryLogin"
          class="w-full py-3 bg-white/10 text-white rounded-xl hover:bg-white/20 transition-colors"
        >
          重新登入
        </button>
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

// 狀態
const loading = ref(true)
const error = ref('')
const hasLineUserId = ref(false)
const loggingIn = ref(false)
const loginSuccess = ref(false)
const loginError = ref('')

// LINE User ID
const lineUserId = ref('')

// 初始化
const initLiff = async () => {
  loading.value = true
  error.value = ''

  try {
    // 嘗試從 URL 參數取得 line_user_id（從 LIFF 返回時）
    const urlParams = new URLSearchParams(window.location.search)
    const idFromUrl = urlParams.get('line_user_id')

    if (idFromUrl) {
      // 從 LIFF 返回，已取得 line_user_id
      lineUserId.value = idFromUrl
      hasLineUserId.value = true

      // 清除 URL 參數
      window.history.replaceState({}, '', window.location.pathname)

      // 執行登入
      await performLogin()
    } else {
      // 檢查是否有 stored line_user_id
      const storedLineUserId = localStorage.getItem('login_line_user_id')
      if (storedLineUserId) {
        lineUserId.value = storedLineUserId
        hasLineUserId.value = true

        // 執行登入
        await performLogin()
      } else {
        // 需要先登入 LINE
        hasLineUserId.value = false
      }
    }
  } catch (err: any) {
    error.value = err.message || '初始化失敗，請重新整理頁面'
  } finally {
    loading.value = false
  }
}

// LINE 登入
const lineLogin = () => {
  // 儲存登入狀態
  localStorage.setItem('login_state', 'in_progress')

  // 導向 LINE LIFF 登入頁面
  const liffUrl = `https://liff.line.me/${config.public.liffId}/teacher/login?redirect=${encodeURIComponent(window.location.href)}`
  window.location.href = liffUrl
}

// 執行登入
const performLogin = async () => {
  loggingIn.value = true
  loginError.value = ''

  try {
    // 儲存 line_user_id
    localStorage.setItem('login_line_user_id', lineUserId.value)

    const response = await $fetch('/api/v1/auth/teacher/line/login', {
      method: 'POST',
      body: {
        line_user_id: lineUserId.value
      }
    })

    // 檢查 response.data 或 response.datas
    const responseData = (response as any).data || (response as any).datas
    const responseCode = (response as any).code

    if (responseCode === 0 && responseData) {
      const token = responseData.token
      const user = responseData.user

      // 設置 authStore 和 localStorage
      authStore.login({
        token,
        refresh_token: '',
        teacher: user,
      })

      loginSuccess.value = true

      setTimeout(() => {
        router.push('/teacher/dashboard')
      }, 1500)
    } else {
      loginError.value = (response as any)?.message || '登入失敗，請稍後再試'
    }
  } catch (err: any) {
    console.error('Login error:', err)
    loginError.value = err.data?.message || err.message || '登入失敗，請稍後再試'
  } finally {
    loggingIn.value = false
  }
}

// 重新嘗試初始化
const retryInit = () => {
  window.location.reload()
}

// 重新登入
const retryLogin = () => {
  hasLineUserId.value = false
  lineUserId.value = ''
  loginError.value = ''
  localStorage.removeItem('login_line_user_id')
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
