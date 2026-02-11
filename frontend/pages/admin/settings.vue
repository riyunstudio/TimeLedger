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
      <!-- 營業時間卡片 -->
      <BaseGlassCard class="mb-6">
        <div class="p-6">
          <div class="flex items-center justify-between mb-2">
            <h2 class="text-lg font-semibold text-white flex items-center gap-2">
              <BaseIcon icon="clock" class="text-primary-400" />
              營業時間
            </h2>
            <!-- 未儲存變更提示 -->
            <span
              v-if="operatingHoursDirty && !operatingHoursLoading"
              class="flex items-center gap-1.5 px-2.5 py-1 bg-amber-500/20 text-amber-400 rounded-full text-xs animate-pulse"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              未儲存變更
            </span>
          </div>
          <p class="text-slate-400 text-sm mb-6">設定中心的營業時間，課表將根據此範圍顯示</p>

          <form @submit.prevent="saveOperatingHours">
            <div class="grid md:grid-cols-2 gap-6">
              <!-- 開始時間 -->
              <div>
                <label class="block text-sm text-slate-400 mb-2">開始時間</label>
                <div class="relative">
                  <select
                    v-model="courseSettings.operating_start_time"
                    class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-primary-500 transition-colors appearance-none cursor-pointer"
                  >
                    <option v-for="hour in 24" :key="hour - 1" :value="String(hour - 1).padStart(2, '0') + ':00'">
                      {{ String(hour - 1).padStart(2, '0') }}:00
                    </option>
                  </select>
                  <svg class="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500 pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </div>

              <!-- 結束時間 -->
              <div>
                <label class="block text-sm text-slate-400 mb-2">結束時間</label>
                <div class="relative">
                  <select
                    v-model="courseSettings.operating_end_time"
                    class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-primary-500 transition-colors appearance-none cursor-pointer"
                  >
                    <option v-for="hour in 24" :key="hour - 1" :value="String(hour - 1).padStart(2, '0') + ':00'">
                      {{ String(hour - 1).padStart(2, '0') }}:00
                    </option>
                  </select>
                  <svg class="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500 pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </div>
              </div>
            </div>

            <!-- 驗證錯誤 -->
            <div v-if="operatingHoursError" class="flex items-center gap-2 mt-3 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <svg class="w-5 h-5 text-red-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p class="text-red-400 text-sm">{{ operatingHoursError }}</p>
            </div>

            <!-- 操作按鈕 -->
            <div class="flex items-center gap-3 mt-6">
              <button
                type="submit"
                :disabled="!!operatingHoursError || operatingHoursLoading || !operatingHoursDirty"
                class="relative px-6 py-2.5 rounded-xl transition-all duration-200 disabled:cursor-not-allowed"
                :class="[
                  operatingHoursLoading
                    ? 'bg-primary-500/20 border-primary-500/30 text-primary-400'
                    : operatingHoursDirty
                      ? 'bg-primary-500/30 border border-primary-500 text-primary-400 hover:bg-primary-500/40 hover:border-primary-400'
                      : 'bg-slate-800/50 border border-slate-700 text-slate-500'
                ]"
              >
                <span v-if="operatingHoursLoading" class="flex items-center gap-2">
                  <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                  儲存中...
                </span>
                <span v-else class="flex items-center gap-2">
                  <svg v-if="operatingHoursDirty" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
                  </svg>
                  <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  {{ operatingHoursDirty ? '儲存變更' : '已儲存' }}
                </span>
              </button>

              <button
                type="button"
                @click="resetOperatingHours"
                :disabled="operatingHoursLoading"
                class="px-4 py-2.5 bg-slate-700/50 border border-slate-600 text-slate-300 rounded-xl hover:bg-slate-700 hover:border-slate-500 transition-colors text-sm disabled:opacity-50 disabled:cursor-not-allowed"
              >
                重設為預設值
              </button>
            </div>
          </form>
        </div>
      </BaseGlassCard>

      <!-- 課程預設值卡片 -->
      <BaseGlassCard>
        <div class="p-6">
          <div class="flex items-center justify-between mb-2">
            <h2 class="text-lg font-semibold text-white flex items-center gap-2">
              <BaseIcon icon="calendar" class="text-primary-400" />
              課程預設值
            </h2>
            <!-- 未儲存變更提示 -->
            <span
              v-if="courseSettingsDirty && !courseSettingsLoading"
              class="flex items-center gap-1.5 px-2.5 py-1 bg-amber-500/20 text-amber-400 rounded-full text-xs animate-pulse"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              未儲存變更
            </span>
          </div>
          <p class="text-slate-400 text-sm mb-6">設定新課程的預設時長，建立課程時將自動套用此設定</p>

          <form @submit.prevent="saveCourseSettings">
            <div class="grid md:grid-cols-2 gap-6">
              <!-- 預設課程時長 -->
              <div>
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
                <!-- 驗證錯誤提示 -->
                <div v-if="courseDurationError" class="flex items-center gap-1.5 mt-1.5 text-red-400 text-xs">
                  <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  {{ courseDurationError }}
                </div>
              </div>

              <!-- 異動申請截止天數 -->
              <div>
                <label class="block text-sm text-slate-400 mb-2">異動申請截止天數</label>
                <BaseInput
                  v-model="courseSettings.exception_lead_days"
                  type="number"
                  placeholder="例如：14"
                  class="w-full"
                  :min="0"
                  :max="365"
                  :error="exceptionLeadDaysError"
                />
                <div v-if="exceptionLeadDaysError" class="flex items-center gap-1.5 mt-1.5 text-red-400 text-xs">
                  <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  {{ exceptionLeadDaysError }}
                </div>
                <p class="text-slate-500 text-xs mt-2">設定老師申請停課或調課需提前的天數（預設為 14 天，設定為 0 表示不限制）</p>
              </div>
            </div>

            <!-- 儲存按鈕 -->
            <div class="flex items-center gap-3 mt-6">
              <button
                type="submit"
                :disabled="courseSettingsLoading || !!courseDurationError || !!exceptionLeadDaysError || !courseSettingsDirty"
                class="relative px-6 py-2.5 rounded-xl transition-all duration-200 disabled:cursor-not-allowed"
                :class="[
                  courseSettingsLoading
                    ? 'bg-primary-500/20 border-primary-500/30 text-primary-400'
                    : courseSettingsDirty
                      ? 'bg-primary-500/30 border border-primary-500 text-primary-400 hover:bg-primary-500/40 hover:border-primary-400'
                      : 'bg-slate-800/50 border border-slate-700 text-slate-500'
                ]"
              >
                <span v-if="courseSettingsLoading" class="flex items-center gap-2">
                  <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                  儲存中...
                </span>
                <span v-else class="flex items-center gap-2">
                  <svg v-if="courseSettingsDirty" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
                  </svg>
                  <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  {{ courseSettingsDirty ? '儲存變更' : '已儲存' }}
                </span>
              </button>
            </div>
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
  exception_lead_days: 14,
  operating_start_time: '00:00',
  operating_end_time: '23:00',
})

// 儲存原始值用於 Dirty State 偵測
const originalSettings = ref({
  default_course_duration: 60,
  exception_lead_days: 14,
  operating_start_time: '00:00',
  operating_end_time: '23:00',
})

const courseSettingsLoading = ref(false)
const operatingHoursLoading = ref(false)

// Toast 通知
const { success: showSuccess, error: showError } = useToast()

// Dirty State 偵測
const operatingHoursDirty = computed(() => {
  return courseSettings.value.operating_start_time !== originalSettings.value.operating_start_time ||
         courseSettings.value.operating_end_time !== originalSettings.value.operating_end_time
})

const courseSettingsDirty = computed(() => {
  return courseSettings.value.default_course_duration !== originalSettings.value.default_course_duration ||
         courseSettings.value.exception_lead_days !== originalSettings.value.exception_lead_days
})

// 課程時長驗證錯誤
const courseDurationError = computed(() => {
  const duration = courseSettings.value.default_course_duration
  if (duration && (duration < 15 || duration > 480)) {
    return '時長需介於 15 到 480 分鐘之間'
  }
  return ''
})

// 例外申請截止天數驗證錯誤
const exceptionLeadDaysError = computed(() => {
  const days = courseSettings.value.exception_lead_days
  if (days !== undefined && days !== null && (days < 0 || days > 365)) {
    return '天數需介於 0 到 365 天之間（0 表示不限制）'
  }
  return ''
})

// 營業時間驗證錯誤
const operatingHoursError = computed(() => {
  const start = courseSettings.value.operating_start_time
  const end = courseSettings.value.operating_end_time

  if (start && end) {
    const startHour = parseInt(start.split(':')[0], 10)
    const endHour = parseInt(end.split(':')[0], 10)

    if (endHour <= startHour) {
      return '結束時間必須晚於開始時間'
    }
  }
  return ''
})

// 重設營業時間為預設值
const resetOperatingHours = async () => {
  const confirmed = await alertConfirm('確定要將營業時間重設為預設值 (00:00 - 23:00) 嗎？')
  if (!confirmed) return

  courseSettings.value.operating_start_time = '00:00'
  courseSettings.value.operating_end_time = '23:00'
}

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
    if (data.datas) {
      const newSettings = {
        default_course_duration: data.datas.default_course_duration ?? 60,
        exception_lead_days: data.datas.exception_lead_days ?? 14,
        operating_start_time: data.datas.operating_start_time || '00:00',
        operating_end_time: data.datas.operating_end_time || '23:00',
      }

      courseSettings.value = { ...newSettings }
      originalSettings.value = { ...newSettings }
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
  if (courseDurationError.value || exceptionLeadDaysError.value) return
  if (!courseSettingsDirty.value) return // 沒有變更則不送出

  courseSettingsLoading.value = true

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
        exception_lead_days: Number(courseSettings.value.exception_lead_days),
      }),
    })

    const data = await response.json()

    if (response.ok) {
      // 更新原始值
      originalSettings.value.default_course_duration = courseSettings.value.default_course_duration
      originalSettings.value.exception_lead_days = courseSettings.value.exception_lead_days
      showSuccess('課程預設值已成功儲存', '儲存成功')
    } else {
      showError(data.message || '儲存失敗', '錯誤')
    }
  } catch (err) {
    showError('儲存失敗，請稍後再試', '錯誤')
  } finally {
    courseSettingsLoading.value = false
  }
}

// 儲存營業時間設定
const saveOperatingHours = async () => {
  if (operatingHoursError.value) return
  if (!operatingHoursDirty.value) return // 沒有變更則不送出

  operatingHoursLoading.value = true

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
        operating_start_time: courseSettings.value.operating_start_time,
        operating_end_time: courseSettings.value.operating_end_time,
      }),
    })

    const data = await response.json()

    if (response.ok) {
      // 更新原始值
      originalSettings.value.operating_start_time = courseSettings.value.operating_start_time
      originalSettings.value.operating_end_time = courseSettings.value.operating_end_time
      showSuccess('營業時間已成功儲存', '儲存成功')
    } else {
      showError(data.message || '儲存失敗', '錯誤')
    }
  } catch (err) {
    showError('儲存失敗，請稍後再試', '錯誤')
  } finally {
    operatingHoursLoading.value = false
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
