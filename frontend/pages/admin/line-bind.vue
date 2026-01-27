<template>
  <div class="p-4 md:p-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
        💬 LINE 通知設定
      </h1>
      <p class="text-slate-400 text-sm md:text-base">
        綁定 LINE 帳號以接收即時例外通知
      </p>
    </div>

    <!-- 綁定狀態卡片 -->
    <BaseGlassCard class="mb-6">
      <div class="p-6">
        <!-- 已綁定狀態 -->
        <div v-if="bindingStatus.isBound" class="text-center">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-green-500/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h2 class="text-xl font-bold text-slate-100 mb-2">已成功綁定 LINE</h2>
          <p class="text-slate-400 mb-6">
            您可以收到老師提交的例外申請通知
          </p>

          <!-- 通知開關 -->
          <div class="bg-white/5 rounded-xl p-4 mb-6">
            <h3 class="text-sm font-medium text-slate-300 mb-4">通知設定</h3>

            <div class="space-y-3">
              <label class="flex items-center justify-between cursor-pointer">
                <span class="text-slate-300">接收新例外申請通知</span>
                <div
                  class="w-12 h-6 rounded-full transition-colors relative"
                  :class="notifySettings.newException ? 'bg-green-500' : 'bg-slate-600'"
                  @click="toggleNotifySetting('newException')"
                >
                  <div
                    class="absolute top-1 w-4 h-4 bg-white rounded-full transition-all"
                    :class="notifySettings.newException ? 'left-7' : 'left-1'"
                  ></div>
                </div>
              </label>

              <label class="flex items-center justify-between cursor-pointer">
                <span class="text-slate-300">接收審核結果通知</span>
                <div
                  class="w-12 h-6 rounded-full transition-colors relative"
                  :class="notifySettings.reviewResult ? 'bg-green-500' : 'bg-slate-600'"
                  @click="toggleNotifySetting('reviewResult')"
                >
                  <div
                    class="absolute top-1 w-4 h-4 bg-white rounded-full transition-all"
                    :class="notifySettings.reviewResult ? 'left-7' : 'left-1'"
                  ></div>
                </div>
              </label>
            </div>
          </div>

          <!-- 解除綁定按鈕 -->
          <button
            @click="showUnbindConfirm = true"
            class="px-6 py-3 bg-red-500/20 border border-red-500/50 text-red-400 rounded-xl hover:bg-red-500/30 transition-colors"
          >
            解除綁定
          </button>
        </div>

        <!-- 未綁定狀態 -->
        <div v-else class="text-center">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-slate-700/50 flex items-center justify-center">
            <svg class="w-10 h-10 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
          </div>
          <h2 class="text-xl font-bold text-slate-100 mb-2">尚未綁定 LINE</h2>
          <p class="text-slate-400 mb-6 max-w-md mx-auto">
            綁定 LINE 後，當老師提交例外申請時，您會立即收到通知，不再錯過任何重要申請。
          </p>

          <!-- 開始綁定按鈕 -->
          <button
            @click="initBinding"
            :disabled="loading"
            class="px-8 py-3 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors disabled:opacity-50"
          >
            <span v-if="loading" class="flex items-center gap-2">
              <svg class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              處理中...
            </span>
            <span v-else>開始綁定</span>
          </button>
        </div>
      </div>
    </BaseGlassCard>

    <!-- QR Code 綁定區域 -->
    <BaseGlassCard v-if="showQRCode" class="mb-6">
      <div class="p-6">
        <h3 class="text-lg font-bold text-slate-100 mb-4 text-center">掃描 QR Code 綁定</h3>

        <!-- QR Code 顯示 -->
        <div class="flex flex-col items-center mb-6">
          <div class="bg-white p-4 rounded-xl mb-4">
            <div class="w-48 h-48 flex items-center justify-center bg-slate-100 rounded-lg">
              <!-- 顯示真實的 QR Code -->
              <img
                v-if="qrCodeUrl"
                :src="qrCodeUrl"
                alt="LINE 綁定 QR Code"
                class="w-full h-full object-contain"
              />
              <!-- 載入中的顯示 -->
              <div v-else class="text-center">
                <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full mx-auto mb-2"></div>
                <p class="text-sm text-slate-500">產生中...</p>
              </div>
            </div>
          </div>

          <!-- 驗證碼顯示 -->
          <div v-if="bindingCode" class="text-center">
            <p class="text-slate-400 mb-2">或傳送驗證碼給 LINE 官方帳號：</p>
            <div class="inline-flex items-center gap-3 bg-white/10 px-6 py-3 rounded-xl">
              <span class="text-3xl font-mono font-bold text-primary-400 tracking-widest">{{ bindingCode }}</span>
              <button
                @click="copyCode"
                class="p-2 hover:bg-white/10 rounded-lg transition-colors"
                title="複製驗證碼"
              >
                <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              </button>
            </div>
            <p class="text-slate-500 text-sm mt-3">
              驗證碼將在 {{ expiresIn }} 後過期
            </p>
          </div>
        </div>

        <!-- 綁定說明 -->
        <div class="bg-white/5 rounded-xl p-4 text-left">
          <h4 class="text-sm font-medium text-slate-300 mb-3">綁定步驟：</h4>
          <ol class="space-y-2 text-slate-400 text-sm">
            <li class="flex items-start gap-2">
              <span class="w-5 h-5 bg-primary-500/30 text-primary-400 rounded-full flex items-center justify-center text-xs flex-shrink-0 mt-0.5">1</span>
              <span>開啟 LINE，搜尋官方帳號「TimeLedger」</span>
            </li>
            <li class="flex items-start gap-2">
              <span class="w-5 h-5 bg-primary-500/30 text-primary-400 rounded-full flex items-center justify-center text-xs flex-shrink-0 mt-0.5">2</span>
              <span>傳送驗證碼 <strong class="text-primary-400">{{ bindingCode }}</strong> 給官方帳號</span>
            </li>
            <li class="flex items-start gap-2">
              <span class="w-5 h-5 bg-primary-500/30 text-primary-400 rounded-full flex items-center justify-center text-xs flex-shrink-0 mt-0.5">3</span>
              <span>官方帳號回覆「綁定成功」即完成</span>
            </li>
          </ol>
        </div>

        <!-- 取消按鈕 -->
        <div class="mt-6 text-center">
          <button
            @click="cancelBinding"
            class="text-slate-400 hover:text-slate-300 transition-colors"
          >
            取消綁定
          </button>
        </div>
      </div>
    </BaseGlassCard>

    <!-- 功能說明 -->
    <BaseGlassCard>
      <div class="p-6">
        <h3 class="text-lg font-bold text-slate-100 mb-4">LINE 通知特色</h3>

        <div class="grid md:grid-cols-2 gap-4">
          <div class="flex items-start gap-3">
            <div class="w-10 h-10 bg-primary-500/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div>
              <h4 class="text-slate-200 font-medium mb-1">即時通知</h4>
              <p class="text-slate-400 text-sm">老師提交例外申請後，馬上收到 LINE 通知</p>
            </div>
          </div>

          <div class="flex items-start gap-3">
            <div class="w-10 h-10 bg-primary-500/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <h4 class="text-slate-200 font-medium mb-1">一鍵處理</h4>
              <p class="text-slate-400 text-sm">點擊通知即可開啟後台，快速處理例外申請</p>
            </div>
          </div>

          <div class="flex items-start gap-3">
            <div class="w-10 h-10 bg-primary-500/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <h4 class="text-slate-200 font-medium mb-1">彈性設定</h4>
              <p class="text-slate-400 text-sm">可選擇性開關不同類型的通知</p>
            </div>
          </div>

          <div class="flex items-start gap-3">
            <div class="w-10 h-10 bg-primary-500/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </div>
            <div>
              <h4 class="text-slate-200 font-medium mb-1">全員通知</h4>
              <p class="text-slate-400 text-sm">中心所有管理員都會收到通知，隨時可處理</p>
            </div>
          </div>
        </div>
      </div>
    </BaseGlassCard>

    <!-- 解除綁定確認對話框 -->
    <GlobalAlert
      v-if="showUnbindConfirm"
      type="warning"
      title="解除 LINE 綁定"
      message="確定要解除 LINE 綁定嗎？解除後將無法收到即時例外通知。"
      confirmText="確定解除"
      cancelText="取消"
      @confirm="unbindLINE"
      @cancel="showUnbindConfirm = false"
    />
  </div>
</template>

<script setup lang="ts">
import BaseGlassCard from '~/components/base/BaseGlassCard.vue'
import GlobalAlert from '~/components/GlobalAlert.vue'
import { alertError, alertSuccess, alertWarning, alertConfirm } from '~/composables/useAlert'
import { useToast } from '~/composables/useToast'

definePageMeta({
  middleware: 'auth-admin',
  layout: 'admin',
})

const config = useRuntimeConfig()
const { success, error } = useToast()

// API 基礎 URL
const API_BASE = config.public.apiBase

// LINE 官方帳號 ID（從環境變數取得）
const lineOfficialAccountId = config.public.lineOfficialAccountId || '@timeledger'

// 狀態
const loading = ref(false)
const showQRCode = ref(false)
const qrCodeUrl = ref('')
const bindingCode = ref('')
const bindingExpiresAt = ref<Date | null>(null)
const expiresIn = ref('')
const timer = ref<number | null>(null)
const showUnbindConfirm = ref(false)

// 綁定狀態
const bindingStatus = ref({
  isBound: false,
  lineUserID: '',
  boundAt: null as Date | null,
  notifyEnabled: true,
  welcomeSent: false,
})

// 通知設定
const notifySettings = ref({
  newException: true,
  reviewResult: true,
})

// 取得綁定狀態
const fetchBindingStatus = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/api/v1/admin/me/line-binding`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (response.ok) {
      const data = await response.json()
      bindingStatus.value = data.datas

      // 初始化通知設定
      notifySettings.value = {
        newException: data.datas.notify_enabled,
        reviewResult: data.datas.notify_enabled,
      }
    }
  } catch (err) {
    console.error('取得綁定狀態失敗:', err)
  }
}

// 取得 QR Code
const fetchQRCode = async (code: string) => {
  try {
    const token = localStorage.getItem('admin_token')
    // 使用含驗證碼的 QR Code API
    const response = await fetch(`${API_BASE}/api/v1/admin/me/line/qrcode-with-code?code=${code}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (response.ok) {
      const blob = await response.blob()
      qrCodeUrl.value = URL.createObjectURL(blob)
    } else {
      console.error('取得 QR Code 失敗')
    }
  } catch (err) {
    console.error('取得 QR Code 失敗:', err)
  }
}

// 初始化綁定
const initBinding = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/api/v1/admin/me/line/bind`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    if (response.ok) {
      const data = await response.json()
      bindingCode.value = data.datas.code
      bindingExpiresAt.value = new Date(data.datas.expires_at)
      showQRCode.value = true

      // 取得 QR Code
      await fetchQRCode(data.datas.code)

      // 啟動倒數計時
      startCountdown()

      success('已產生驗證碼，請使用 LINE 掃描或傳送驗證碼')
    } else {
      const data = await response.json()
      alertError(data.message || '初始化綁定失敗')
    }
  } catch (err) {
    console.error('初始化綁定失敗:', err)
    alertError('初始化綁定失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}

// 倒數計時
const startCountdown = () => {
  const updateExpiresIn = () => {
    if (!bindingExpiresAt.value) return

    const now = new Date()
    const diff = bindingExpiresAt.value.getTime() - now.getTime()

    if (diff <= 0) {
      expiresIn.value = '已過期'
      if (timer.value) {
        clearInterval(timer.value)
      }
      return
    }

    const minutes = Math.floor(diff / 60000)
    const seconds = Math.floor((diff % 60000) / 1000)
    expiresIn.value = `${minutes} 分 ${seconds} 秒`
  }

  updateExpiresIn()
  timer.value = window.setInterval(updateExpiresIn, 1000)
}

// 複製驗證碼
const copyCode = async () => {
  try {
    await navigator.clipboard.writeText(bindingCode.value)
    success('已複製驗證碼')
  } catch (err) {
    error('複製失敗，請手動複製')
  }
}

// 取消綁定
const cancelBinding = () => {
  showQRCode.value = false
  bindingCode.value = ''
  qrCodeUrl.value = ''
  if (timer.value) {
    clearInterval(timer.value)
  }
}

// 解除綁定
const unbindLINE = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/api/v1/admin/me/line/unbind`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (response.ok) {
      success('已解除 LINE 綁定')
      showUnbindConfirm.value = false
      await fetchBindingStatus()
    } else {
      const data = await response.json()
      alertError(data.message || '解除綁定失敗')
    }
  } catch (err) {
    console.error('解除綁定失敗:', err)
    alertError('解除綁定失敗，請稍後再試')
  }
}

// 切換通知設定
const toggleNotifySetting = async (setting: 'newException' | 'reviewResult') => {
  notifySettings.value[setting] = !notifySettings.value[setting]

  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/api/v1/admin/me/line/notify-settings`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        enabled: notifySettings.value[setting],
      }),
    })

    if (!response.ok) {
      const data = await response.json()
      alertError(data.message || '更新通知設定失敗')
      // 回滾狀態
      notifySettings.value[setting] = !notifySettings.value[setting]
    }
  } catch (err) {
    console.error('更新通知設定失敗:', err)
    alertError('更新通知設定失敗，請稍後再試')
    // 回滾狀態
    notifySettings.value[setting] = !notifySettings.value[setting]
  }
}

// 頁面載入時取得綁定狀態
onMounted(() => {
  fetchBindingStatus()
})

// 頁面卸載時清除計時器
onUnmounted(() => {
  if (timer.value) {
    clearInterval(timer.value)
  }
})
</script>
