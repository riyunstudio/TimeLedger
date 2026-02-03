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
          @click="goHome"
          class="px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors"
        >
          返回首頁
        </button>
      </div>

      <!-- LINE 登入狀態 -->
      <div v-else-if="!hasLineUserId" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20">
        <div class="text-center mb-6">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-primary-500/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-white mb-1">教師註冊</h1>
          <p class="text-primary-400 font-medium">加入 TimeLedger 人才庫</p>
        </div>

        <div class="bg-white/5 rounded-xl p-4 mb-6">
          <h3 class="text-slate-300 font-medium mb-3">加入 TimeLedger 人才庫，您可以：</h3>
          <ul class="text-slate-400 text-sm space-y-2">
            <li class="flex items-start gap-2">
              <svg class="w-5 h-5 text-primary-400 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span>接案賺取額外收入</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="w-5 h-5 text-primary-400 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span>自由管理您的課表</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="w-5 h-5 text-primary-400 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span>獲得多元教學機會</span>
            </li>
          </ul>
        </div>

        <!-- LINE 登入按鈕 -->
        <button
          @click="lineLogin"
          class="w-full py-4 bg-[#06C755] hover:bg-[#05B54A] text-white font-medium rounded-xl transition-colors flex items-center justify-center gap-3"
        >
          <svg class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
          </svg>
          LINE 快速登入
        </button>
        <p class="text-center text-slate-500 text-sm mt-3">
          點擊上方按鈕透過 LINE 登入繼續
        </p>
      </div>

      <!-- 註冊表單 -->
      <div v-else class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20">
        <div class="text-center mb-6">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-primary-500/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-white mb-1">完成註冊</h1>
          <p class="text-primary-400 font-medium">請填寫您的個人資料</p>
        </div>

        <!-- 成功提示 -->
        <div v-if="registered" class="text-center">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-success-500/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-success-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h2 class="text-xl font-bold text-white mb-2">註冊成功！</h2>
          <p class="text-slate-400 mb-6">歡迎加入 TimeLedger</p>
          <button
            @click="goToDashboard"
            class="px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors"
          >
            前往老師後台
          </button>
        </div>

        <!-- 註冊表單 -->
        <form v-else @submit.prevent="handleRegister" class="space-y-4">
          <!-- LINE ID 隱藏欄位 -->
          <input type="hidden" v-model="form.line_user_id" />

          <!-- 姓名 -->
          <div>
            <label class="block text-slate-300 text-sm mb-2">姓名</label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="請輸入您的姓名"
              class="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-primary-500 transition-colors"
            />
          </div>

          <!-- Email -->
          <div>
            <label class="block text-slate-300 text-sm mb-2">Email</label>
            <input
              v-model="form.email"
              type="email"
              required
              placeholder="請輸入您的 Email"
              class="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-primary-500 transition-colors"
            />
          </div>

          <!-- 錯誤訊息 -->
          <div v-if="formError" class="bg-critical-500/20 border border-critical-500/30 rounded-xl p-3">
            <p class="text-critical-500 text-sm">{{ formError }}</p>
          </div>

          <!-- 提交按鈕 -->
          <button
            type="submit"
            :disabled="submitting"
            class="w-full py-4 bg-primary-500 hover:bg-primary-600 text-white font-medium rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <span v-if="submitting" class="inline-block w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            {{ submitting ? '處理中...' : '完成註冊' }}
          </button>
        </form>

        <!-- 取消按鈕 -->
        <button
          @click="logout"
          class="w-full py-3 text-slate-400 hover:text-white transition-colors mt-2"
        >
          使用其他帳號
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
const registered = ref(false)
const submitting = ref(false)
const formError = ref('')

// 表單資料
const form = ref({
  line_user_id: '',
  name: '',
  email: '',
})

// 初始化 LIFF
const initLiff = async () => {
  loading.value = true
  error.value = ''

  try {
    // 嘗試從 URL 參數取得 line_user_id（從 LIFF 返回時）
    const urlParams = new URLSearchParams(window.location.search)
    const lineUserId = urlParams.get('line_user_id')

    if (lineUserId) {
      // 從 LIFF 返回，已取得 line_user_id
      form.value.line_user_id = lineUserId
      hasLineUserId.value = true

      // 清除 URL 參數
      window.history.replaceState({}, '', window.location.pathname)
    } else {
      // 檢查是否有 stored line_user_id
      const storedLineUserId = localStorage.getItem('register_line_user_id')
      if (storedLineUserId) {
        form.value.line_user_id = storedLineUserId
        hasLineUserId.value = true
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
  // 儲存註冊狀態
  localStorage.setItem('register_state', 'in_progress')

  // 導向 LINE LIFF 登入頁面
  // 使用 liffId 構造 LIFF URL
  const liffUrl = `https://liff.line.me/${config.public.liffId}/teacher/login?redirect=${encodeURIComponent(window.location.href)}`
  window.location.href = liffUrl
}

// 登出
const logout = () => {
  localStorage.removeItem('register_line_user_id')
  localStorage.removeItem('register_state')
  authStore.logout()
  // 重新整理頁面
  window.location.reload()
}

// 提交註冊
const handleRegister = async () => {
  if (!form.value.line_user_id) {
    formError.value = 'LINE 登入資訊遺失，請重新登入'
    return
  }

  if (!form.value.name.trim()) {
    formError.value = '請輸入姓名'
    return
  }

  if (!form.value.email.trim()) {
    formError.value = '請輸入 Email'
    return
  }

  submitting.value = true
  formError.value = ''

  try {
    const response = await fetch(`${config.public.apiBase}/teacher/public/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        line_user_id: form.value.line_user_id,
        name: form.value.name,
        email: form.value.email,
      }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || '註冊失敗，請稍後再試')
    }

    // 註冊成功，儲存 token
    if (data.datas?.token) {
      localStorage.setItem('teacher_token', data.datas.token)
      authStore.login({
        token: data.datas.token,
        teacher: data.datas.teacher,
      })
    }

    // 清除註冊狀態
    localStorage.removeItem('register_line_user_id')
    localStorage.removeItem('register_state')

    registered.value = true
  } catch (err: any) {
    formError.value = err.message || '註冊失敗，請稍後再試'
  } finally {
    submitting.value = false
  }
}

// 返回首頁
const goHome = () => {
  router.push('/')
}

// 前往老師後台
const goToDashboard = () => {
  router.push('/teacher/dashboard')
}

onMounted(() => {
  initLiff()
})
</script>
