<template>
  <div class="p-4 md:p-6">
    <!-- 頁面標題 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-white">廣播管理</h1>
      <p class="text-slate-400 mt-1">發送 LINE 廣播訊息給中心的老師</p>
    </div>

    <!-- 衝突警告 -->
    <div
      v-if="conflictWarning"
      class="mb-6 p-4 bg-yellow-500/10 border border-yellow-500/30 rounded-lg"
    >
      <div class="flex items-start gap-3">
        <svg class="w-5 h-5 text-yellow-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div>
          <h3 class="text-yellow-500 font-medium">⚠️ 訊息衝突檢測</h3>
          <p class="text-slate-300 text-sm mt-1">{{ conflictWarning }}</p>
        </div>
      </div>
    </div>

    <!-- 成功訊息 -->
    <div
      v-if="successMessage"
      class="mb-6 p-4 bg-green-500/10 border border-green-500/30 rounded-lg"
    >
      <div class="flex items-start gap-3">
        <svg class="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div>
          <h3 class="text-green-500 font-medium">廣播發送完成</h3>
          <p class="text-slate-300 text-sm mt-1">{{ successMessage }}</p>
        </div>
      </div>
    </div>

    <!-- 主要內容區域 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 左側：公告內容輸入區 -->
      <div class="glass-card p-6">
        <h2 class="text-lg font-semibold text-white mb-4">公告內容</h2>

        <form @submit.prevent="handleBroadcast">
          <!-- 訊息類型選擇 -->
          <div class="mb-4">
            <label class="block text-slate-300 mb-2">
              訊息類型
              <span class="text-slate-500 text-xs ml-2">必填</span>
            </label>
            <div class="grid grid-cols-2 gap-3">
              <label
                class="flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-colors"
                :class="form.type === 'GENERAL' ? 'bg-primary-500/20 border-primary-500' : 'bg-white/5 border-white/10 hover:bg-white/10'"
              >
                <input
                  type="radio"
                  v-model="form.type"
                  value="GENERAL"
                  class="sr-only"
                />
                <div class="flex-1">
                  <div class="text-white font-medium text-sm">一般公告</div>
                  <div class="text-slate-400 text-xs mt-0.5">一般通知訊息</div>
                </div>
                <svg v-if="form.type === 'GENERAL'" class="w-5 h-5 text-primary-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
              </label>
              <label
                class="flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-colors"
                :class="form.type === 'URGENT' ? 'bg-red-500/20 border-red-500' : 'bg-white/5 border-white/10 hover:bg-white/10'"
              >
                <input
                  type="radio"
                  v-model="form.type"
                  value="URGENT"
                  class="sr-only"
                />
                <div class="flex-1">
                  <div class="text-white font-medium text-sm">緊急通知</div>
                  <div class="text-slate-400 text-xs mt-0.5">需要立即注意的訊息</div>
                </div>
                <svg v-if="form.type === 'URGENT'" class="w-5 h-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
              </label>
            </div>
          </div>

          <!-- 訊息標題 -->
          <div class="mb-4">
            <label for="title" class="block text-slate-300 mb-2">
              標題
              <span class="text-slate-500 text-xs ml-2">必填</span>
            </label>
            <input
              id="title"
              v-model="form.title"
              type="text"
              placeholder="輸入訊息標題"
              maxlength="50"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white placeholder-slate-500 focus:outline-none focus:border-primary-500"
              required
            />
            <p class="text-xs text-slate-500 mt-1 text-right">
              {{ form.title.length }}/50
            </p>
          </div>

          <!-- 訊息內容 -->
          <div class="mb-4">
            <label for="message" class="block text-slate-300 mb-2">
              訊息內容
              <span class="text-slate-500 text-xs ml-2">必填（最多 2000 字）</span>
            </label>
            <textarea
              id="message"
              v-model="form.message"
              placeholder="輸入公告訊息內容..."
              rows="8"
              maxlength="2000"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white placeholder-slate-500 focus:outline-none focus:border-primary-500 resize-none"
              required
            ></textarea>
            <p class="text-xs text-slate-500 mt-1 text-right">
              {{ form.message.length }}/2000
            </p>
          </div>

          <!-- 警告訊息（可選） -->
          <div class="mb-4">
            <label for="warning" class="block text-slate-300 mb-2">
              警告提示
              <span class="text-slate-500 text-xs ml-2">可選，用於提醒注意事項</span>
            </label>
            <textarea
              id="warning"
              v-model="form.warning"
              placeholder="輸入警告提示資訊..."
              rows="2"
              maxlength="200"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white placeholder-slate-500 focus:outline-none focus:border-yellow-500 resize-none"
            ></textarea>
            <p class="text-xs text-slate-500 mt-1 text-right">
              {{ form.warning.length }}/200
            </p>
          </div>

          <!-- 動作按鈕（可選） -->
          <div class="mb-4">
            <label for="actionLabel" class="block text-slate-300 mb-2">
              按鈕文字
              <span class="text-slate-500 text-xs ml-2">可選，顯示動作按鈕時必填</span>
            </label>
            <input
              id="actionLabel"
              v-model="form.actionLabel"
              type="text"
              placeholder="例如：前往查看"
              maxlength="20"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white placeholder-slate-500 focus:outline-none focus:border-primary-500"
            />
          </div>

          <!-- 動作連結（可選） -->
          <div v-if="form.actionLabel" class="mb-6">
            <label for="actionUrl" class="block text-slate-300 mb-2">
              連結網址
              <span class="text-slate-500 text-xs ml-2">必填</span>
            </label>
            <input
              id="actionUrl"
              v-model="form.actionUrl"
              type="url"
              placeholder="https://timeledger.app/..."
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white placeholder-slate-500 focus:outline-none focus:border-primary-500"
            />
          </div>

          <!-- 預估人數 -->
          <div class="mb-6 p-3 bg-white/5 rounded-lg">
            <div class="flex items-center justify-between">
              <span class="text-slate-400 text-sm">預估接收人數</span>
              <span class="text-white font-medium">{{ estimatedRecipients }} 位老師</span>
            </div>
            <p class="text-xs text-slate-500 mt-1">
              {{ form.teacherIds.length > 0 ? '已選擇特定老師' : '已綁定 LINE 且非佔位符的老師' }}
            </p>
          </div>

          <!-- 發送按鈕 -->
          <div class="flex gap-3">
            <button
              type="button"
              @click="handlePreview"
              :disabled="!canPreview || previewing"
              class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <svg v-if="previewing" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              {{ previewing ? '預覽中...' : '重新整理預覽' }}
            </button>
            <button
              type="submit"
              :disabled="!canSubmit || sending"
              class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <svg v-if="sending" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
              </svg>
              {{ sending ? '發送中...' : '發送廣播' }}
            </button>
          </div>
        </form>
      </div>

      <!-- 右側：LINE 訊息預覽 -->
      <div class="glass-card p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-white">LINE 預覽</h2>
          <span class="text-xs text-slate-500">模擬畫面</span>
        </div>

        <LineFlexPreview
          :title="previewTitle"
          :content="previewContent"
          :warning="form.warning || undefined"
          :action-label="form.actionLabel || undefined"
          :action-url="form.actionUrl || undefined"
          :disabled="false"
          :show-date="true"
        />

        <div class="mt-4 p-3 bg-white/5 rounded-lg">
          <h4 class="text-slate-400 text-xs font-medium mb-2">預覽說明</h4>
          <ul class="text-slate-500 text-xs space-y-1">
            <li>• 預覽僅供參考，實際顯示可能略有差異</li>
            <li>• 標題前綴會根據訊息類型自動添加</li>
            <li>• 動作按鈕僅在填寫按鈕文字後才會顯示</li>
          </ul>
        </div>
      </div>
    </div>
  </div>

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
import LineFlexPreview from '~/components/Notification/LineFlexPreview.vue'

definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

const notificationUI = useNotification()
const api = useApi()
const { confirm: alertConfirm, error: alertError, success: alertSuccess } = useAlert()

// 表單資料
const form = ref({
  type: 'GENERAL' as 'GENERAL' | 'URGENT',
  title: '',
  message: '',
  warning: '',
  actionLabel: '',
  actionUrl: '',
  teacherIds: [] as number[]
})

// 狀態
const sending = ref(false)
const previewing = ref(false)
const conflictWarning = ref('')
const successMessage = ref('')
const teacherCount = ref(0)

// 計算預估接收人數
const estimatedRecipients = computed(() => {
  if (form.value.teacherIds.length > 0) {
    return form.value.teacherIds.length
  }
  return teacherCount.value
})

// 預覽標題（帶前綴）
const previewTitle = computed(() => {
  const prefix = form.value.type === 'URGENT' ? '🚨 ' : '🔔 '
  return prefix + (form.value.title || '新的通知')
})

// 預覽內容
const previewContent = computed(() => {
  return form.value.message || '輸入訊息內容...'
})

// 是否可以預覽
const canPreview = computed(() => {
  return form.value.title.length > 0 && form.value.message.length > 0
})

// 是否可以發送
const canSubmit = computed(() => {
  return form.value.title.length > 0 &&
    form.value.message.length > 0 &&
    form.value.message.length <= 2000 &&
    (!form.value.actionLabel || form.value.actionUrl.length > 0)
})

// 取得老師數量統計（計算可接收 LINE 廣播的老師數量）
const fetchTeacherCount = async () => {
  try {
    const response = await api.get<any>('/admin/teachers')
    // API 回應結構：response.data 是老師陣列（useApi 已提取 datas 層）
    // response = { data: [...], total: 4, page: 1, ... }
    const teachers = response.data
    if (Array.isArray(teachers)) {
      // 篩選條件：
      // 1. line_user_id 存在（已綁定 LINE 才能收到廣播）
      // 2. is_placeholder 為 false（真實老師，不是佔位符）
      teacherCount.value = teachers.filter((t: any) =>
        t.line_user_id && t.line_user_id.length > 0 && !t.is_placeholder
      ).length
    } else if (Array.isArray(response)) {
      // 備用：直接使用 response（如果沒有分頁結構）
      teacherCount.value = response.filter((t: any) =>
        t.line_user_id && t.line_user_id.length > 0 && !t.is_placeholder
      ).length
    }
  } catch (error) {
    console.error('Failed to fetch teacher count:', error)
    teacherCount.value = 0
  }
}

// 處理預覽
const handlePreview = async () => {
  if (!canPreview.value) return

  previewing.value = true
  conflictWarning.value = ''

  try {
    // 模擬預覽檢查
    if (form.value.message.length > 1800) {
      conflictWarning.value = '訊息內容接近上限（2000 字），實際顯示時可能會被截斷。'
    }
  } finally {
    previewing.value = false
  }
}

// 處理發送廣播
const handleBroadcast = async () => {
  if (!canSubmit.value) return

  // 二次確認
  const confirmed = await alertConfirm(
    `確定要發送廣播訊息嗎？\n\n` +
    `• 標題：${form.value.title}\n` +
    `• 內容：${form.value.message.length} 字\n` +
    `• 預估人數：${estimatedRecipients.value} 位老師\n\n` +
    `此操作將透過 LINE 發送訊息給所有已綁定的老師。`
  )

  if (!confirmed) return

  sending.value = true
  successMessage.value = ''
  conflictWarning.value = ''

  try {
    const requestBody = {
      message: form.value.message,
      teacher_ids: form.value.teacherIds.length > 0 ? form.value.teacherIds : undefined
    }

    const response = await api.post<any>('/admin/notifications/broadcast', requestBody)

    if (response.code === 0 || response.code === 200) {
      const data = response.datas || response
      successMessage.value = `成功發送給 ${data.success_count || 0} 位老師，失敗 ${data.failed_count || 0} 位。`

      // 重置表單
      form.value = {
        type: 'GENERAL',
        title: '',
        message: '',
        warning: '',
        actionLabel: '',
        actionUrl: '',
        teacherIds: []
      }

      // 重新取得老師數量
      await fetchTeacherCount()
    } else {
      await alertError(response.message || '發送失敗，請稍後再試')
    }
  } catch (error: any) {
    console.error('Broadcast failed:', error)

    // 處理衝突警告
    if (error.data?.datas?.conflicts) {
      const conflicts = error.data.datas.conflicts
      if (conflicts.some((c: any) => !c.can_override)) {
        await alertError('偵測到衝突，無法發送廣播訊息')
      } else {
        conflictWarning.value = '偵測到部分衝突，是否仍要發送？'
        // 如果有可覆蓋的衝突，這裡可以讓用戶選擇是否繼續
      }
    } else {
      await alertError(error.message || '發送失敗，請稍後再試')
    }
  } finally {
    sending.value = false
  }
}

// 初始化
onMounted(async () => {
  await fetchTeacherCount()
})
</script>

<style scoped>
.glass-card {
  @apply bg-slate-800/50 backdrop-blur-sm border border-white/10 rounded-xl;
}

textarea:focus,
input:focus {
  outline: none;
}
</style>
