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
            <span class="text-white">{{ invitation.expires_at ? formatDate(invitation.expires_at) : '無期限' }}</span>
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
            LINE 登入並接受邀請
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
            {{ accepting ? $t('common.submitting') : '接受邀請並加入' }}
          </button>
          <button
            @click="logout"
            class="w-full py-3 text-slate-400 hover:text-white transition-colors"
          >
            {{ $t('auth.invite.logoutOther') }}
          </button>
        </div>

        <!-- 已接受成功 -->
        <div v-if="accepted" class="text-center">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-success-500/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-success-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h2 class="text-xl font-bold text-white mb-2">
            {{ isAlreadyMember ? $t('auth.invite.welcomeBackMember') : $t('auth.invite.joinSuccess') }}
          </h2>
          <p class="text-slate-400 mb-6">
            {{ isAlreadyMember ? '' : $t('auth.invite.welcomeTo', { center: invitation.center_name }) }}
          </p>
          <button
            @click="goToTeacherDashboard"
            class="px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors"
          >
            {{ $t('auth.invite.goToDashboard') }}
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
const { $liff } = useNuxtApp()

// 取得 token
const token = computed(() => route.params.token as string)

// 狀態
const loading = ref(true)
const error = ref('')
const invitation = ref<any>(null)
const accepting = ref(false)
const accepted = ref(false)
const isAlreadyMember = ref(false)

// 檢查是否已登入
const isLoggedIn = computed(() => authStore.isAuthenticated && authStore.isTeacher)

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

// 取得用戶資訊（包含 email）
const getLineUserInfo = async () => {
  // 防禦性檢查：確保 $liff 已初始化
  if (!$liff || typeof $liff.getProfile !== 'function') {
    throw new Error('LIFF SDK 尚未完全載入')
  }

  // 嘗試從 ID Token 獲取 email（需要 openid 權限）
  let email = ''
  if (typeof $liff.getDecodedIDToken === 'function') {
    try {
      const decodedToken = await $liff.getDecodedIDToken()
      email = decodedToken?.email || ''
    } catch (e) {
      console.warn('無法獲取 ID Token email:', e)
    }
  }

  // 從 Profile 獲取基本資訊
  const profile = await $liff.getProfile()
  const accessToken = $liff.getAccessToken()

  if (!accessToken || !profile) {
    throw new Error('無法取得 LINE 資訊，請重新登入')
  }

  return {
    userId: profile.userId,
    displayName: profile.displayName,
    pictureUrl: profile.pictureUrl,
    email,
    accessToken
  }
}

// LINE 登入
const lineLogin = async () => {
  // 防禦性檢查：確保 $liff 已初始化
  if (!$liff || typeof $liff.isLoggedIn !== 'function') {
    error.value = 'LIFF SDK 尚未完全載入，請重新整理頁面'
    return
  }

  try {
    // 檢查是否已登入 LINE
    const isLoggedInLine = $liff.isLoggedIn()

    if (!isLoggedInLine) {
      // 未登入，使用 LIFF SDK 登入
      $liff.login()
      return
    }

    // 已登入 LINE，取得用戶資訊
    const userInfo = await getLineUserInfo()

    // 儲存邀請 token 和用戶資訊到 localStorage
    localStorage.setItem('invitation_token', token.value)
    localStorage.setItem('invitation_line_user_id', userInfo.userId)
    localStorage.setItem('invitation_access_token', userInfo.accessToken)

    // 如果有 email，也儲存起來
    if (userInfo.email) {
      localStorage.setItem('invitation_email', userInfo.email)
    }

    // 執行登入並接受邀請
    await performLoginAndAccept(userInfo.userId, userInfo.accessToken)
  } catch (err: any) {
    error.value = err.message || 'LINE 登入失敗，請稍後再試'
  }
}

// 登出
const logout = () => {
  localStorage.removeItem('teacher_token')
  localStorage.removeItem('invitation_token')
  localStorage.removeItem('invitation_line_user_id')
  localStorage.removeItem('invitation_access_token')
  authStore.logout()
  // 重新整理頁面
  window.location.reload()
}

// 執行接受邀請並登入
const performLoginAndAccept = async (lineUserId: string, accessToken: string) => {
  accepting.value = true
  error.value = ''

  try {
    // 1. 先接受邀請（重要：後端會自動為新老師建立帳號）
    await acceptInvitationWithToken(lineUserId, accessToken)

    // 2. 接受成功後，再執行 LINE 登入取得 JWT Token
    const response = await fetch(`${config.public.apiBase}/auth/teacher/line/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ line_user_id: lineUserId, access_token: accessToken }),
    })

    const data = await response.json()

    if (!response.ok || !data.datas?.token) {
      throw new Error(data.message || '登入失敗')
    }

    // 3. 登入成功，儲存 token 並更新狀態
    localStorage.setItem('teacher_token', data.datas.token)
    authStore.login({ token: data.datas.token, teacher: data.datas.teacher })
  } catch (err: any) {
    error.value = err.message || '處理邀請時發生錯誤，請稍後再試'
  } finally {
    accepting.value = false
  }
}

// 使用指定 token 接受邀請
const acceptInvitationWithToken = async (lineUserId: string, accessToken: string) => {
  // 嘗試獲取已儲存的 email
  const savedEmail = localStorage.getItem('invitation_email') || ''

  try {
    const response = await fetch(`${config.public.apiBase}/invitations/${token.value}/accept`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        line_user_id: lineUserId,
        access_token: accessToken,
        email: savedEmail, // 傳遞 LINE ID Token 中的 email
      }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || '接受邀請失敗')
    }

    // 檢查是否早就是成員 (由後端回傳 status: 'ALREADY_MEMBER')
    if (data.datas?.status === 'ALREADY_MEMBER') {
      isAlreadyMember.value = true
    }

    accepted.value = true

    // 清除邀請相關的 localStorage
    localStorage.removeItem('invitation_token')
    localStorage.removeItem('invitation_line_user_id')
    localStorage.removeItem('invitation_access_token')
    localStorage.removeItem('invitation_email')
  } catch (err: any) {
    error.value = err.message || '接受邀請失敗，請稍後再試'
  } finally {
    accepting.value = false
  }
}

// 接受邀請
const acceptInvitation = async () => {
  if (!authStore.user?.line_user_id) {
    error.value = '無法取得登入資訊，請重新登入'
    return
  }

  // 取得 Access Token（使用可選鏈結防禦性檢查）
  const accessToken = $liff?.getAccessToken?.() || null
  if (!accessToken) {
    error.value = '無法取得 LINE Access Token，請重新登入'
    return
  }

  accepting.value = true
  try {
    const response = await fetch(`${config.public.apiBase}/invitations/${token.value}/accept`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ line_user_id: authStore.user.line_user_id, access_token: accessToken }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || '接受邀請失敗')
    }

    accepted.value = true

    // 檢查是否早就是成員 (由後端回傳 status: 'ALREADY_MEMBER')
    if (data.datas?.status === 'ALREADY_MEMBER') {
      isAlreadyMember.value = true
    }

    // === 新增：使用後端返回的 Token 自動登入 ===
    if (data.datas?.token && data.datas?.teacher) {
      const authData = {
        token: data.datas.token,
        refresh_token: '',
        teacher: {
          id: data.datas.teacher.id,
          name: data.datas.teacher.name,
          email: data.datas.teacher.email,
          line_user_id: data.datas.teacher.line_user_id,
          avatar_url: data.datas.teacher.avatar_url,
        }
      }

      // 保存 token 到 localStorage
      localStorage.setItem('teacher_token', data.datas.token)

      // 更新 authStore
      authStore.login(authData)
    }
    // ==========================================

    // 清除邀請相關的 localStorage
    localStorage.removeItem('invitation_token')
    localStorage.removeItem('invitation_line_user_id')
    localStorage.removeItem('invitation_access_token')
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
  const savedLineUserId = localStorage.getItem('invitation_line_user_id')
  const savedAccessToken = localStorage.getItem('invitation_access_token')

  if (savedToken && savedToken === token.value && savedLineUserId && savedAccessToken) {
    // 使用者已經透過邀請頁面導向登入，回來後自動執行登入並接受邀請
    setTimeout(() => {
      performLoginAndAccept(savedLineUserId, savedAccessToken)
    }, 500)
  }
}

onMounted(async () => {
  await fetchInvitation()
  checkInvitationLogin()
})
</script>
