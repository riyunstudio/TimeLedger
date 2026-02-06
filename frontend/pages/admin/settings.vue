<template>
  <div class="p-4 md:p-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
        帳號設定
      </h1>
      <p class="text-slate-400 text-sm md:text-base">
        管理您的帳號資訊和密碼
      </p>
    </div>

    <!-- 個人資料卡片 -->
    <BaseGlassCard class="mb-6">
      <div class="p-6">
        <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          個人資料
        </h2>

        <div class="grid md:grid-cols-2 gap-4">
          <div class="bg-white/5 rounded-xl p-4">
            <p class="text-slate-400 text-sm mb-1">姓名</p>
            <p class="text-white font-medium">{{ profile.name || '-' }}</p>
          </div>

          <div class="bg-white/5 rounded-xl p-4">
            <p class="text-slate-400 text-sm mb-1">Email</p>
            <p class="text-white font-medium">{{ profile.email || '-' }}</p>
          </div>

          <div class="bg-white/5 rounded-xl p-4">
            <p class="text-slate-400 text-sm mb-1">角色</p>
            <p class="text-white font-medium">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="{
                  'bg-purple-500/20 text-purple-400': profile.role === 'OWNER',
                  'bg-blue-500/20 text-blue-400': profile.role === 'ADMIN',
                  'bg-green-500/20 text-green-400': profile.role === 'STAFF'
                }"
              >
                {{ roleText(profile.role) }}
              </span>
            </p>
          </div>

          <div class="bg-white/5 rounded-xl p-4">
            <p class="text-slate-400 text-sm mb-1">所屬中心</p>
            <p class="text-white font-medium">{{ profile.center_name || '-' }}</p>
          </div>
        </div>
      </div>
    </BaseGlassCard>

    <!-- 修改密碼卡片 -->
    <BaseGlassCard class="mb-6">
      <div class="p-6">
        <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
          修改密碼
        </h2>

        <form @submit.prevent="changePassword" class="space-y-4">
          <div>
            <label class="block text-sm text-slate-400 mb-2">舊密碼</label>
            <BaseInput
              v-model="passwordForm.oldPassword"
              type="password"
              placeholder="請輸入舊密碼"
              class="w-full"
            />
          </div>

          <div>
            <label class="block text-sm text-slate-400 mb-2">新密碼</label>
            <BaseInput
              v-model="passwordForm.newPassword"
              type="password"
              placeholder="請輸入新密碼（至少 6 個字元）"
              class="w-full"
            />
          </div>

          <div>
            <label class="block text-sm text-slate-400 mb-2">確認新密碼</label>
            <BaseInput
              v-model="passwordForm.confirmPassword"
              type="password"
              placeholder="請再次輸入新密碼"
              class="w-full"
              :error="passwordForm.confirmPassword && passwordForm.newPassword !== passwordForm.confirmPassword ? '密碼不一致' : ''"
            />
          </div>

          <div class="flex items-center gap-4 pt-2">
            <button
              type="submit"
              :disabled="passwordLoading || !isPasswordFormValid"
              class="px-6 py-2.5 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="passwordLoading" class="flex items-center gap-2">
                <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                處理中...
              </span>
              <span v-else>更新密碼</span>
            </button>

            <span v-if="passwordMessage" :class="passwordSuccess ? 'text-green-400' : 'text-red-400'" class="text-sm">
              {{ passwordMessage }}
            </span>
          </div>
        </form>
      </div>
    </BaseGlassCard>

    <!-- LINE 綁定區塊 -->
    <BaseGlassCard class="mb-6">
      <div class="p-6">
        <h2 class="text-lg font-semibold text-white mb-6 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          LINE 通知設定
        </h2>

        <AdminLineBindingSettings />
      </div>
    </BaseGlassCard>

    <!-- 登出按鈕 -->
    <BaseGlassCard>
      <div class="p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-slate-300 mb-1">登出帳號</p>
            <p class="text-slate-500 text-sm">安全的登出您的帳號</p>
          </div>

          <button
            @click="logout"
            class="px-6 py-2.5 bg-critical-500/20 border border-critical-500/50 text-critical-400 rounded-xl hover:bg-critical-500/30 transition-colors"
          >
            登出
          </button>
        </div>
      </div>
    </BaseGlassCard>
  </div>
</template>

<script setup lang="ts">
import BaseGlassCard from '~/components/base/BaseGlassCard.vue'
import BaseInput from '~/components/base/BaseInput.vue'
import AdminLineBindingSettings from '~/components/Admin/AdminLineBindingSettings.vue'
import { alertError, alertSuccess, alertConfirm } from '~/composables/useAlert'

definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

const router = useRouter()
const config = useRuntimeConfig()

// API 基礎 URL
const API_BASE = config.public.apiBase

// 個人資料
const profile = ref({
  name: '',
  email: '',
  role: '',
  center_name: '',
})

// 密碼表單
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const passwordLoading = ref(false)
const passwordMessage = ref('')
const passwordSuccess = ref(false)

// 取得管理員資料
const fetchProfile = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/me/profile`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (response.ok) {
      const data = await response.json()
      profile.value = data.datas
    }
  } catch (err) {
    console.error('取得個人資料失敗:', err)
  }
}

// 角色文字
const roleText = (role: string) => {
  const roles: Record<string, string> = {
    OWNER: '擁有者',
    ADMIN: '管理員',
    STAFF: '員工',
  }
  return roles[role] || role
}

// 密碼表單驗證
const isPasswordFormValid = computed(() => {
  return (
    passwordForm.value.oldPassword.length >= 6 &&
    passwordForm.value.newPassword.length >= 6 &&
    passwordForm.value.newPassword === passwordForm.value.confirmPassword
  )
})

// 修改密碼
const changePassword = async () => {
  if (!isPasswordFormValid.value) return

  passwordLoading.value = true
  passwordMessage.value = ''
  passwordSuccess.value = false

  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/me/change-password`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        old_password: passwordForm.value.oldPassword,
        new_password: passwordForm.value.newPassword,
      }),
    })

    const data = await response.json()

    if (response.ok) {
      passwordSuccess.value = true
      passwordMessage.value = '密碼已成功修改'
      passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
      await alertSuccess('密碼已成功修改')
    } else {
      passwordSuccess.value = false
      passwordMessage.value = data.message || '修改失敗'
      await alertError(data.message || '舊密碼錯誤')
    }
  } catch (err) {
    passwordSuccess.value = false
    passwordMessage.value = '修改失敗，請稍後再試'
    await alertError('修改失敗，請稍後再試')
  } finally {
    passwordLoading.value = false
  }
}

// 登出
const logout = async () => {
  const confirmed = await alertConfirm('確定要登出嗎？')
  if (!confirmed) return

  localStorage.removeItem('admin_token')
  localStorage.removeItem('admin_user')
  router.push('/admin/login')
}

// 頁面載入時取得資料
onMounted(async () => {
  await Promise.all([
    fetchProfile(),
  ])
})
</script>
