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
        <p class="text-slate-300">正在載入邀請資訊...</p>
      </div>

      <!-- 錯誤訊息 -->
      <div v-else-if="error" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20 text-center">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-critical-500/20 flex items-center justify-center">
          <svg class="w-8 h-8 text-critical-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        <h2 class="text-xl font-bold text-white mb-2">無法載入邀請</h2>
        <p class="text-slate-400 mb-6">{{ error }}</p>
        <button
          @click="goHome"
          class="px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors"
        >
          返回首頁
        </button>
      </div>

      <!-- 邀請資訊 -->
      <div v-else-if="invitation" class="bg-white/10 backdrop-blur-xl rounded-2xl p-8 border border-white/20">
        <!-- 中心資訊 -->
        <div class="text-center mb-6">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-primary-500/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-white mb-1">{{ invitation.center_name || 'XX 中心' }}</h1>
          <p class="text-primary-400 font-medium">邀請您加入</p>
        </div>

        <!-- 邀請詳情 -->
        <div class="bg-white/5 rounded-xl p-4 mb-6 space-y-3">
          <div class="flex justify-between items-center">
            <span class="text-slate-400">職位</span>
            <span class="text-white font-medium">{{ roleText }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-slate-400">邀請時間</span>
            <span class="text-white">{{ formatDate(invitation.created_at) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-slate-400">有效期限</span>
            <span class="text-white">{{ formatDate(invitation.expires_at) }}</span>
          </div>
        </div>

        <!-- 邀請訊息 -->
        <div v-if="invitation.message" class="bg-white/5 rounded-xl p-4 mb-6">
          <p class="text-slate-400 text-sm mb-2">邀請訊息</p>
          <p class="text-slate-200">{{ invitation.message }}</p>
        </div>

        <!-- 登入按鈕 -->
        <div v-if="!isLoggedIn" class="space-y-3">
          <button
            @click="lineLogin"
            class="w-full py-4 bg-[#06C755] hover:bg-[#05B54A] text-white font-medium rounded-xl transition-colors flex items-center justify-center gap-3"
          >
            <svg class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
            </svg>
            LINE 快速登入並接受邀請
          </button>
          <p class="text-center text-slate-500 text-sm">
            點擊上方按鈕透過 LINE 登入並接受邀請
          </p>
        </div>

        <!-- 已登入但非受邀請者 -->
        <div v-else-if="isLoggedIn && !isInvitedUser" class="space-y-3">
          <div class="bg-critical-500/20 border border-critical-500/30 rounded-xl p-4 text-center">
            <p class="text-critical-500 font-medium mb-1">此邀請不屬於您</p>
            <p class="text-slate-400 text-sm">您目前登入的帳號與邀請不符</p>
          </div>
          <button
            @click="logout"
            class="w-full py-3 bg-white/10 text-white rounded-xl hover:bg-white/20 transition-colors"
          >
            登出並使用其他帳號
          </button>
        </div>

        <!-- 已登入且是受邀請者 -->
        <div v-else-if="isLoggedIn && isInvitedUser" class="space-y-3">
          <div class="bg-success-500/20 border border-success-500/30 rounded-xl p-4 text-center">
            <p class="text-success-500 font-medium mb-1">歡迎回來！</p>
            <p class="text-slate-400 text-sm">您已登入，可以接受邀請</p>
          </div>
          <button
            @click="acceptInvitation"
            :disabled="accepting"
            class="w-full py-4 bg-primary-500 hover:bg-primary-600 text-white font-medium rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <span v-if="accepting" class="inline-block w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            {{ accepting ? '處理中...' : '接受邀請並加入' }}
          </button>
          <button
            @click="logout"
            class="w-full py-3 text-slate-400 hover:text-white transition-colors"
          >
            使用其他帳號
          </button>
        </div>

        <!-- 已接受成功 -->
        <div v-if="accepted" class="text-center">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-success-500/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-success-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h2 class="text-xl font-bold text-white mb-2">加入成功！</h2>
          <p class="text-slate-400 mb-6">歡迎加入 {{ invitation.center_name }}</p>
          <button
            @click="goToTeacherDashboard"
            class="px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors"
          >
            前往老師後台
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

const config = useRuntimeConfig()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// 取得 token
const token = computed(() => route.params.token as string)

// 狀態
const loading = ref(true)
const error = ref('')
const invitation = ref<any>(null)
const accepting = ref(false)
const accepted = ref(false)

// 檢查是否已登入
const isLoggedIn = computed(() => authStore.isAuthenticated && authStore.userType === 'teacher')

// 檢查是否為受邀請者
const isInvitedUser = computed(() => {
  if (!isLoggedIn.value || !invitation.value) return false
  return authStore.user?.email === invitation.value.email
})

// 角色文字
const roleText = computed(() => {
  const roleMap: Record<string, string> = {
    TEACHER: '正式老師',
    SUBSTITUTE: '代課老師',
  }
  return roleMap[invitation.value?.role] || invitation.value?.role || '老師'
})

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

// 取得邀請資訊
const fetchInvitation = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await fetch(`${config.public.apiBase}/invitations/${token.value}`)
    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || '邀請不存在或已過期')
    }

    invitation.value = data.datas
  } catch (err: any) {
    error.value = err.message || '無法載入邀請資訊'
  } finally {
    loading.value = false
  }
}

// LINE 登入
const lineLogin = () => {
  // 儲存邀請 token 到 localStorage，接受邀請時使用
  localStorage.setItem('invitation_token', token.value)

  // 導向 LINE LIFF 登入
  const liffUrl = `${config.public.liffUrl}/teacher/login?redirect=${encodeURIComponent(window.location.href)}`
  window.location.href = liffUrl
}

// 登出
const logout = () => {
  localStorage.removeItem('teacher_token')
  localStorage.removeItem('invitation_token')
  authStore.logout()
  // 重新整理頁面
  window.location.reload()
}

// 接受邀請
const acceptInvitation = async () => {
  if (!authStore.user?.id_token) {
    error.value = '無法取得登入資訊，請重新登入'
    return
  }

  accepting.value = true
  try {
    const response = await fetch(`${config.public.apiBase}/invitations/${token.value}/accept`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id_token: authStore.user.id_token,
      }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || '接受邀請失敗')
    }

    accepted.value = true

    // 清除邀請 token
    localStorage.removeItem('invitation_token')
  } catch (err: any) {
    error.value = err.message || '接受邀請失敗，請稍後再試'
  } finally {
    accepting.value = false
  }
}

// 返回首頁
const goHome = () => {
  router.push('/')
}

// 前往老師後台
const goToTeacherDashboard = () => {
  router.push('/teacher/dashboard')
}

// 檢查是否有待處理的邀請登入
const checkInvitationLogin = () => {
  const savedToken = localStorage.getItem('invitation_token')
  if (savedToken && savedToken === token.value) {
    // 使用者已經透過邀請頁面導向登入，回來後自動接受邀請
    setTimeout(() => {
      if (isLoggedIn.value && isInvitedUser.value) {
        acceptInvitation()
      }
    }, 500)
  }
}

onMounted(async () => {
  await fetchInvitation()
  checkInvitationLogin()
})
</script>
