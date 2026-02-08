<template>
  <div class="p-4 md:p-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
        設定
      </h1>
      <p class="text-slate-400 text-sm md:text-base">
        管理您的帳號與中心設定
      </p>
    </div>

    <!-- Tabs -->
    <div class="flex border-b border-white/10 mb-6">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="activeTab = tab.id"
        class="px-6 py-3 text-sm font-medium transition-colors border-b-2 -mb-px"
        :class="[
          activeTab === tab.id
            ? 'text-primary-400 border-primary-400'
            : 'text-slate-400 border-transparent hover:text-slate-300'
        ]"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Tab Content: 帳號設定 -->
    <div v-show="activeTab === 'account'">
      <!-- 個人資料卡片 -->
      <BaseGlassCard class="mb-6">
        <div class="p-6">
          <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <BaseIcon icon="user" class="text-primary-400" />
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
            <BaseIcon icon="lock" class="text-primary-400" />
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
                  <BaseIcon icon="spinner" class="animate-spin w-4 h-4" />
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
            <BaseIcon icon="chat" class="text-primary-400" />
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

    <!-- Tab Content: 中心設定 -->
    <div v-show="activeTab === 'center'">
      <!-- 課程預設值卡片 -->
      <BaseGlassCard>
        <div class="p-6">
          <h2 class="text-lg font-semibold text-white mb-2 flex items-center gap-2">
            <BaseIcon icon="clock" class="text-primary-400" />
            課程預設值
          </h2>
          <p class="text-slate-400 text-sm mb-6">設定新課程的預設時長，建立課程時將自動套用此設定</p>

          <form @submit.prevent="saveCourseSettings">
            <div class="flex items-end gap-4">
              <div class="flex-1 max-w-xs">
                <label class="block text-sm text-slate-400 mb-2">預設課程時長（分鐘）</label>
                <BaseInput
                  v-model="courseSettings.default_course_duration"
                  type="number"
                  placeholder="例如：60"
                  class="w-full"
                  :min="15"
                  :max="480"
                  :error="courseDurationError"
                />
              </div>

              <button
                type="submit"
                :disabled="courseSettingsLoading || !!courseDurationError"
                class="px-6 py-2.5 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span v-if="courseSettingsLoading" class="flex items-center gap-2">
                  <BaseIcon icon="spinner" class="animate-spin w-4 h-4" />
                  儲存中...
                </span>
                <span v-else>儲存設定</span>
              </button>
            </div>

            <p v-if="courseSettingsMessage" :class="courseSettingsSuccess ? 'text-green-400' : 'text-red-400'" class="text-sm mt-3">
              {{ courseSettingsMessage }}
            </p>
          </form>
        </div>
      </BaseGlassCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseGlassCard from '~/components/base/BaseGlassCard.vue'
import BaseInput from '~/components/base/BaseInput.vue'
import BaseIcon from '~/components/base/Icon.vue'
import AdminLineBindingSettings from '~/components/Admin/AdminLineBindingSettings.vue'
import { alertError, alertSuccess, alertConfirm } from '~/composables/useAlert'

interface TabItem {
  id: string
  label: string
}

definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

const router = useRouter()
const config = useRuntimeConfig()

// API 基礎 URL
const API_BASE = config.public.apiBase

// Center ID
const { getCenterId } = useCenterId()

// Tabs
const tabs: TabItem[] = [
  { id: 'account', label: '帳號設定' },
  { id: 'center', label: '中心設定' },
]
const activeTab = ref('account')

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

// 課程設定
const courseSettings = ref({
  default_course_duration: 60,
})

const courseSettingsLoading = ref(false)
const courseSettingsMessage = ref('')
const courseSettingsSuccess = ref(false)

// 課程時長驗證錯誤
const courseDurationError = computed(() => {
  const duration = courseSettings.value.default_course_duration
  if (duration && (duration < 15 || duration > 480)) {
    return '時長需介於 15 到 480 分鐘之間'
  }
  return ''
})

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

// 取得中心設定
const fetchCenterSettings = async () => {
  try {
    const centerId = getCenterId()
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/centers/${centerId}/settings`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      console.error('取得中心設定失敗: HTTP', response.status)
      return
    }

    const text = await response.text()
    if (!text.trim()) {
      console.error('取得中心設定失敗: 空响应')
      return
    }

    const data = JSON.parse(text)
    if (data.datas && data.datas.default_course_duration) {
      courseSettings.value.default_course_duration = data.datas.default_course_duration
    }
  } catch (err) {
    console.error('取得中心設定失敗:', err)
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

// 儲存課程設定
const saveCourseSettings = async () => {
  if (courseDurationError.value) return

  courseSettingsLoading.value = true
  courseSettingsMessage.value = ''
  courseSettingsSuccess.value = false

  try {
    const centerId = getCenterId()
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/centers/${centerId}/settings`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        default_course_duration: Number(courseSettings.value.default_course_duration),
      }),
    })

    const data = await response.json()

    if (response.ok) {
      courseSettingsSuccess.value = true
      courseSettingsMessage.value = '設定已成功儲存'
      await alertSuccess('設定已成功儲存')
    } else {
      courseSettingsSuccess.value = false
      courseSettingsMessage.value = data.message || '儲存失敗'
      await alertError(data.message || '儲存失敗')
    }
  } catch (err) {
    courseSettingsSuccess.value = false
    courseSettingsMessage.value = '儲存失敗，請稍後再試'
    await alertError('儲存失敗，請稍後再試')
  } finally {
    courseSettingsLoading.value = false
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
    fetchCenterSettings(),
  ])
})
</script>
