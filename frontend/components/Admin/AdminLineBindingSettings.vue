<template>
  <div class="space-y-6">
    <!-- 已綁定狀態 -->
    <div v-if="bindingStatus.isBound" class="text-center">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-green-500/20 flex items-center justify-center">
        <svg class="w-8 h-8 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <h3 class="text-lg font-bold text-slate-100 mb-2">已成功綁定 LINE</h3>
      <p class="text-slate-400 text-sm mb-6">
        您可以收到老師提交的例外申請通知
      </p>

      <!-- 通知開關 -->
      <div class="bg-white/5 rounded-xl p-4 mb-6 text-left">
        <h4 class="text-sm font-medium text-slate-300 mb-4">功能開關</h4>

        <div class="space-y-4">
          <label class="flex items-center justify-between cursor-pointer">
            <div>
              <p class="text-slate-200 text-sm">接收新例外申請通知</p>
              <p class="text-xs text-slate-500">當有老師提交代課或調課申請時通知我</p>
            </div>
            <div
              class="w-10 h-5 rounded-full transition-colors relative"
              :class="notifySettings.newException ? 'bg-green-500' : 'bg-slate-600'"
              @click="toggleNotifySetting('newException')"
            >
              <div
                class="absolute top-0.5 w-4 h-4 bg-white rounded-full transition-all"
                :class="notifySettings.newException ? 'left-5' : 'left-0.5'"
              ></div>
            </div>
          </label>

          <label class="flex items-center justify-between cursor-pointer">
            <div>
              <p class="text-slate-200 text-sm">接收審核結果通知</p>
              <p class="text-xs text-slate-500">當其他管理員處理申請時通知我</p>
            </div>
            <div
              class="w-10 h-5 rounded-full transition-colors relative"
              :class="notifySettings.reviewResult ? 'bg-green-500' : 'bg-slate-600'"
              @click="toggleNotifySetting('reviewResult')"
            >
              <div
                class="absolute top-0.5 w-4 h-4 bg-white rounded-full transition-all"
                :class="notifySettings.reviewResult ? 'left-5' : 'left-0.5'"
              ></div>
            </div>
          </label>
        </div>
      </div>

      <!-- 解除綁定按鈕 -->
      <button
        @click="showUnbindConfirm = true"
        class="text-sm text-red-400 hover:text-red-300 transition-colors"
      >
        解除 LINE 帳號綁定
      </button>
    </div>

    <!-- 未綁定狀態 -->
    <div v-else-if="!showQRCode" class="text-center">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-slate-700/50 flex items-center justify-center">
        <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
      </div>
      <h3 class="text-lg font-bold text-slate-100 mb-2">尚未綁定 LINE</h3>
      <p class="text-slate-400 text-sm mb-6 max-w-sm mx-auto">
        綁定 LINE 後，當老師提交例外申請時，您會立即收到通知。
      </p>

      <button
        @click="initBinding"
        :disabled="loading"
        class="px-8 py-2.5 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors disabled:opacity-50 flex items-center gap-2 mx-auto"
      >
        <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span>開始綁定</span>
      </button>
    </div>

    <!-- QR Code 綁定區域 -->
    <div v-else class="text-center space-y-6">
      <h3 class="text-lg font-bold text-slate-100">掃描 QR Code 或傳送驗證碼</h3>

      <div class="flex flex-col items-center">
        <div class="bg-white p-4 rounded-xl mb-4">
          <div class="w-40 h-40 flex items-center justify-center bg-slate-100 rounded-lg">
            <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="QR Code" class="w-full h-full" />
            <div v-else class="animate-spin w-6 h-6 border-2 border-primary-500 border-t-transparent rounded-full"></div>
          </div>
        </div>

        <div v-if="bindingCode">
          <p class="text-slate-400 text-sm mb-2">官方帳號驗證碼：</p>
          <div class="text-2xl font-mono font-bold text-primary-400 tracking-widest bg-white/5 px-6 py-2 rounded-lg">
            {{ bindingCode }}
          </div>
          <p class="text-slate-500 text-xs mt-3">
            驗證碼將在 {{ expiresIn }} 後過期
          </p>
        </div>
      </div>

      <div class="bg-white/5 rounded-xl p-4 text-left text-sm">
        <p class="text-slate-300 font-medium mb-2">如何綁定？</p>
        <ul class="space-y-1 text-slate-400">
          <li>1. 掃描上方 QR Code 加入官方帳號</li>
          <li>2. 傳送驗證碼 <strong class="text-primary-400">{{ bindingCode }}</strong> 給帳號</li>
          <li>3. 收到「綁定成功」即完成設定</li>
        </ul>
      </div>

      <button
        @click="cancelBinding"
        class="text-sm text-slate-500 hover:text-slate-400 transition-colors"
      >
        取消並返回設定
      </button>
    </div>

    <!-- 解除綁定確認 -->
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
import GlobalAlert from '~/components/base/GlobalAlert.vue'

const emit = defineEmits(['status-updated'])

const config = useRuntimeConfig()
const api = useApi()
const { success, error } = useToast()
const { confirm: alertConfirm, error: alertError } = useAlert()

// 狀態
const loading = ref(false)
const showQRCode = ref(false)
const qrCodeUrl = ref('')
const bindingCode = ref('')
const bindingExpiresAt = ref<Date | null>(null)
const expiresIn = ref('')
const timer = ref<number | null>(null)
const showUnbindConfirm = ref(false)

const bindingStatus = ref({
  isBound: false,
  lineUserID: '',
  notifyEnabled: true,
})

const notifySettings = ref({
  newException: true,
  reviewResult: true,
})

// 取得綁定狀態
const fetchStatus = async () => {
  try {
    const response = await api.get<any>('/admin/me/line-binding')
    if (response) {
      // API 返回 snake_case，需轉換為前端使用的 camelCase
      bindingStatus.value = {
        isBound: response.is_bound ?? false,
        lineUserID: response.line_user_id ?? '',
        notifyEnabled: response.notify_enabled ?? true,
      }
      notifySettings.value = {
        newException: response.notify_enabled ?? true,
        reviewResult: response.notify_enabled ?? true,
      }
      emit('status-updated', response)
    }
  } catch (err) {
    console.error('Failed to fetch LINE binding status:', err)
  }
}

// 初始化綁定
const initBinding = async () => {
  loading.value = true
  try {
    const response = await api.post<any>('/admin/me/line/bind', {})
    if (response) {
      bindingCode.value = response.code
      bindingExpiresAt.value = new Date(response.expires_at)
      showQRCode.value = true

      // 取得 QR Code
      await fetchQRCode(response.code)
      startCountdown()
      success('已產生驗證碼')
    }
  } catch (err: any) {
    alertError(err.message || '初始化綁定失敗')
  } finally {
    loading.value = false
  }
}

const fetchQRCode = async (code: string) => {
  try {
    console.log('Fetching QR Code from:', `/admin/me/line/qrcode-with-code?code=${code}`)

    // 取得 base64 格式的 QR Code
    const response = await api.get<{ image: string }>(`/admin/me/line/qrcode-with-code?code=${code}`)

    console.log('Response received:', response)

    if (response && response.image) {
      // 直接使用 data URL（已經是完整的 data:image/png;base64,... 格式）
      qrCodeUrl.value = response.image
      console.log('QR Code loaded successfully')
    } else {
      console.error('Invalid response format:', response)
      error('QR Code 資料格式錯誤')
    }
  } catch (err: any) {
    console.error('Failed to fetch QR Code:', err)
    error(err.message || '取得 QR Code 失敗，請稍後再試')
  }
}

const startCountdown = () => {
  const update = () => {
    if (!bindingExpiresAt.value) return
    const diff = bindingExpiresAt.value.getTime() - Date.now()
    if (diff <= 0) {
      expiresIn.value = '已過期'
      if (timer.value) clearInterval(timer.value)
      return
    }
    const m = Math.floor(diff / 60000)
    const s = Math.floor((diff % 60000) / 1000)
    expiresIn.value = `${m}分${s}秒`
  }
  update()
  timer.value = window.setInterval(update, 1000)
}

const cancelBinding = () => {
  showQRCode.value = false
  if (timer.value) clearInterval(timer.value)
}

const unbindLINE = async () => {
  try {
    await api.delete('/admin/me/line/unbind')
    success('已解除綁定')
    showUnbindConfirm.value = false
    await fetchStatus()
  } catch (err: any) {
    alertError(err.message || '解除失敗')
  }
}

const toggleNotifySetting = async (setting: 'newException' | 'reviewResult') => {
  const newValue = !notifySettings.value[setting]
  notifySettings.value[setting] = newValue

  try {
    await api.patch('/admin/me/line/notify-settings', {
      enabled: newValue,
    })
    success('設定已更新')
  } catch (err) {
    notifySettings.value[setting] = !newValue // 回滾
    error('更新失敗')
  }
}

onMounted(fetchStatus)
onUnmounted(() => {
  if (timer.value) clearInterval(timer.value)
})
</script>
