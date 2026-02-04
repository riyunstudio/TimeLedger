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
        <p class="text-slate-300">{{ $t('auth.initializing') }}</p>
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
          {{ $t('auth.refresh') }}
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
          <h1 class="text-2xl font-bold text-white mb-1">{{ $t('auth.loginTitle') }}</h1>
          <p class="text-primary-400 font-medium">{{ $t('auth.welcomeBack') }}</p>
        </div>

        <div class="bg-white/5 rounded-xl p-4 mb-6">
          <p class="text-slate-300 text-sm">
            {{ $t('auth.loginHint') }}
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
          {{ $t('auth.lineLogin') }}
        </button>
        <p class="text-center text-slate-500 text-sm mt-3">
          {{ $t('auth.lineLoginFull') }}
        </p>
      </div>

      <!-- 登入中 -->
      <div v-else-if="loggingIn" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="inline-block w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full animate-spin mb-4"></div>
        <p class="text-slate-300">{{ $t('auth.loggingIn') }}</p>
      </div>

      <!-- 登入成功 -->
      <div v-else-if="loginSuccess" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-success-500/20 flex items-center justify-center">
          <svg class="w-8 h-8 text-success-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 class="text-xl font-bold text-white mb-2">{{ $t('auth.loginSuccess') }}</h2>
        <p class="text-slate-400">{{ $t('auth.redirecting') }}</p>
      </div>

      <!-- 登入失敗 -->
      <div v-else-if="loginError" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-critical-500/20 flex items-center justify-center">
          <svg class="w-8 h-8 text-critical-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        <div class="text-center mb-6">
          <h2 class="text-xl font-bold text-white mb-2">{{ $t('auth.loginFailed') }}</h2>
          <p class="text-slate-400">{{ loginError }}</p>
        </div>
        <button
          @click="retryLogin"
          class="w-full py-3 bg-white/10 text-white rounded-xl hover:bg-white/20 transition-colors"
        >
          {{ $t('auth.retryLogin') }}
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
const route = useRoute()
const authStore = useAuthStore()
const { $liff } = useNuxtApp()

// 狀態
const loading = ref(true)
const error = ref('')
const hasLineUserId = ref(false)
const loggingIn = ref(false)
const loginSuccess = ref(false)
const loginError = ref('')

// LINE User ID
const lineUserId = ref('')

// 等待 LIFF SDK 初始化完成
const waitForLiffInit = (maxWait = 5000): Promise<boolean> => {
  return new Promise((resolve) => {
    const startTime = Date.now()
    const checkInterval = 100

    const check = () => {
      if ($liff) {
        // SDK 已初始化，檢查是否可以使用
        try {
          // 嘗試存取 isLoggedIn 方法
          if (typeof $liff.isLoggedIn === 'function') {
            resolve(true)
            return
          }
        } catch (e) {
          // SDK 尚未完全就緒，繼續等待
        }
      }

      if (Date.now() - startTime > maxWait) {
        console.warn('LIFF SDK 初始化超時')
        resolve(false)
        return
      }

      setTimeout(check, checkInterval)
    }

    check()
  })
}

// 處理 OAuth 重導回來的參數
const handleOAuthCallback = async () => {
  const code = route.query.code as string
  const state = route.query.state as string

  console.log('[OAuth] 處理 callback，code:', !!code, 'state:', !!state)

  // 【優化】首先等待 LIFF SDK 初始化完成
  // 這是關鍵步驟，確保 SDK 就緒後才能進行後續操作
  const initialized = await waitForLiffInit()

  if (!initialized) {
    console.error('[OAuth] LIFF SDK 初始化超時')
    throw new Error('LIFF SDK 初始化超時，請重新整理頁面')
  }

  console.log('[OAuth] SDK 初始化完成')

  // 【優化】檢查是否有 OAuth 回調參數 (code 和 state)
  if (code && state) {
    console.log('[OAuth] 檢測到 OAuth 回調參數，嘗試獲取登入狀態...')

    // 嘗試獲取登入狀態
    const isLoggedIn = $liff.isLoggedIn()
    console.log('[OAuth] LINE 登入狀態:', isLoggedIn)

    if (isLoggedIn) {
      // 登入成功，取得用戶資訊
      try {
        const profile = await $liff.getProfile()
        lineUserId.value = profile.userId
        hasLineUserId.value = true
        console.log('[OAuth] 已登入，userId:', profile.userId)
        await performLogin()
      } catch (err) {
        console.error('[OAuth] 獲取用戶資訊失敗:', err)
        throw new Error('無法獲取 LINE 用戶資訊，請重新登入')
      }
    } else {
      // 【優化】如果 SDK 認為未登入但有 OAuth 參數，嘗試檢查 Access Token
      // 有些情況下 SDK 可能在 callback URL 解析後需要額外處理
      try {
        const accessToken = $liff.getAccessToken()
        if (accessToken) {
          console.log('[OAuth] 找到 Access Token，嘗試獲取用戶資訊...')
          const profile = await $liff.getProfile()
          lineUserId.value = profile.userId
          hasLineUserId.value = true
          console.log('[OAuth] 已登入 (透過 token)，userId:', profile.userId)
          await performLogin()
        } else {
          console.warn('[OAuth] SDK 未自動完成登入，無 Access Token')
          hasLineUserId.value = false
        }
      } catch (err) {
        console.error('[OAuth] 處理 OAuth 回調時發生錯誤:', err)
        hasLineUserId.value = false
      }
    }
  } else {
    // 無 OAuth 回調參數，清除狀態
    console.log('[OAuth] 無 OAuth 回調參數，清除 URL 參數')
    hasLineUserId.value = false

    // 【優化】清除 URL 中的殘留參數（避免用戶刷新頁面時重複處理）
    if (route.query.code || route.query.state) {
      await router.replace({ path: route.path, query: {} })
    }
  }
}

// 初始化
const initLiff = async () => {
  loading.value = true
  error.value = ''

  console.log('[Login] 開始初始化，URL:', window.location.href)
  console.log('[Login] Route query:', route.query)

  try {
    // 檢查 LIFF 是否已初始化
    if (!$liff) {
      throw new Error('LIFF 尚未初始化，請重新整理頁面')
    }

    console.log('[Login] LIFF 已初始化')

    // 檢查 URL 中是否有 OAuth callback 參數
    if (route.query.code) {
      console.log('[Login] 檢測到 OAuth callback，開始處理...')
      await handleOAuthCallback()
      return
    }

    console.log('[Login] 無 OAuth callback，檢查登入狀態...')

    // 檢查是否已登入 LINE
    const isLoggedIn = $liff.isLoggedIn()
    console.log('[Login] LINE 登入狀態:', isLoggedIn)

    if (isLoggedIn) {
      // 已登入 LINE，取得用戶資訊
      const profile = await $liff.getProfile()
      lineUserId.value = profile.userId
      hasLineUserId.value = true
      console.log('[Login] 已登入，userId:', profile.userId)

      // 執行登入
      await performLogin()
    } else {
      // 未登入 LINE，需要先登入
      hasLineUserId.value = false
      console.log('[Login] 未登入，等待用戶點擊登入按鈕')
    }
  } catch (err: any) {
    error.value = err.message || '初始化失敗，請重新整理頁面'
    console.error('[Login] 初始化錯誤:', err)
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

// 執行登入
const performLogin = async () => {
  loggingIn.value = true
  loginError.value = ''

  try {
    // 取得 LINE Access Token
    const accessToken = $liff.getAccessToken()
    if (!accessToken) {
      throw new Error('無法取得 LINE Access Token，請重新登入')
    }

    const response = await $fetch('/api/v1/auth/teacher/line/login', {
      method: 'POST',
      body: {
        line_user_id: lineUserId.value,
        access_token: accessToken
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

      // 檢查是否有邀請 token，若有則導向邀請頁面
      setTimeout(() => {
        const savedInviteToken = localStorage.getItem('invitation_token')
        if (savedInviteToken) {
          router.push(`/invite/${savedInviteToken}`)
        } else {
          router.push('/teacher/dashboard')
        }
      }, 1500)
    } else if (responseCode === 40010) {
      // TEACHER_NOT_REGISTERED - 老師尚未註冊
      // 重導到註冊頁面，並攜帶 LINE 用戶資訊

      // 儲存 LINE 用戶資訊到 localStorage，供註冊頁面使用
      localStorage.setItem('register_line_user_id', lineUserId.value)

      // 重導到註冊頁面
      router.push('/teacher/register')
      return
    } else {
      loginError.value = (response as any)?.message || '登入失敗，請稍後再試'
    }
  } catch (err: any) {
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
