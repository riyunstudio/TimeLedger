<script setup lang="ts">
import { ref } from 'vue'
import { useQuickError, useFormErrors, useAsyncOperation } from '~/composables/useQuickError'

// ==================== 使用簡化版錯誤處理 ====================

const {
  showError,
  showSuccess,
  showWarning,
  showInfo,
  confirm,
  request,
  handle,
  handlePermission,
  handleValidation,
  handleConflicts,
} = useQuickError()

// ==================== 表單錯誤處理 ====================

const {
  errors: formErrors,
  hasErrors,
  setFieldError,
  setErrors,
  clearFieldError,
  clearErrors,
  getFieldError,
} = useFormErrors()

// ==================== 異步操作處理 ====================

const { loading, execute } = useAsyncOperation()

// ==================== 範例資料 ====================

const formData = ref({
  name: '',
  email: '',
  phone: '',
})

const userData = ref<any>(null)
const apiLoading = ref(false)

// ==================== 錯誤處理範例 ====================

/**
 * 範例 1：顯示成功訊息
 */
const handleSaveSuccess = async () => {
  // 模擬儲存成功
  await new Promise(resolve => setTimeout(resolve, 500))

  await showSuccess('資料已成功儲存', '儲存成功')
}

/**
 * 範例 2：顯示錯誤訊息
 */
const handleShowError = async () => {
  await showError('發生錯誤，請稍後再試', '操作失敗')
}

/**
 * 範例 3：顯示警告訊息
 */
const handleShowWarning = async () => {
  await showWarning('此操作將無法復原，是否繼續？', '警告')
}

/**
 * 範例 4：顯示確認對話框
 */
const handleDelete = async () => {
  const shouldDelete = await confirm('確定要刪除此項目嗎？此操作無法復原。', '確認刪除')

  if (shouldDelete) {
    await showSuccess('刪除成功')
  }
}

/**
 * 範例 5：處理 API 回應（自動顯示錯誤）
 */
const handleApiRequest = async () => {
  apiLoading.value = true

  try {
    // 模擬 API 請求
    const response = await mockApiCall('success')

    const result = await request(() => Promise.resolve(response))

    if (result) {
      await showSuccess('取得資料成功')
    }
  } finally {
    apiLoading.value = false
  }
}

/**
 * 範例 6：處理 API 回應（手動控制）
 */
const handleManualApi = async () => {
  const response = await mockApiCall('error')

  const result = await handle(response)

  if (!result.success) {
    // 這裡可以自定義錯誤處理邏輯
    console.log('錯誤訊息：', result.error)
    return
  }

  // 成功處理
  console.log('資料：', result.data)
}

/**
 * 範例 7：處理權限錯誤
 */
const handlePermissionError = async () => {
  await handlePermission('FORBIDDEN', '您沒有權限執行此操作')
}

/**
 * 範例 8：處理驗證錯誤
 */
const handleValidationError = () => {
  handleValidation({
    name: ['姓名必填'],
    email: ['電子郵件格式錯誤'],
    phone: ['電話號碼格式不正確'],
  })
}

/**
 * 範例 9：處理排課衝突
 */
const handleScheduleConflicts = async () => {
  await handleConflicts([
    {
      type: 'TEACHER',
      message: '老師張三在 14:00-15:00 已有課程',
    },
    {
      type: 'ROOM',
      message: '教室 A 在 14:00-15:00 已被預約',
    },
  ])
}

/**
 * 範例 10：表單驗證
 */
const handleFormSubmit = async () => {
  // 清除舊錯誤
  clearErrors()

  // 驗證表單
  const errors: Record<string, string> = {}

  if (!formData.value.name) {
    errors.name = '姓名必填'
  }

  if (!formData.value.email) {
    errors.email = '電子郵件必填'
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.value.email)) {
    errors.email = '電子郵件格式錯誤'
  }

  if (Object.keys(errors).length > 0) {
    setErrors(errors)
    await showWarning('請修正表單中的錯誤', '驗證失敗')
    return
  }

  await showSuccess('表單驗證通過')
}

/**
 * 範例 11：異步操作（自動錯誤處理）
 */
const handleAsyncOperation = async () => {
  const result = await execute(async () => {
    await new Promise(resolve => setTimeout(resolve, 1000))
    throw new Error('操作失敗')
  })

  if (result.success) {
    console.log('操作成功：', result.data)
  } else {
    console.log('操作失敗：', result.error)
  }
}

/**
 * 範例 12：使用 async/await 捕獲錯誤
 */
const handleWithTryCatch = async () => {
  try {
    const response = await mockApiCall('success')
    await showSuccess('操作成功')
    return response.data
  } catch (error) {
    await showError(String(error))
    return null
  }
}

// ==================== 模擬 API 函數 ====================

type MockApiType = 'success' | 'error' | 'validation' | 'permission'

async function mockApiCall(type: MockApiType) {
  await new Promise(resolve => setTimeout(resolve, 500))

  switch (type) {
    case 'success':
      return { code: 'SUCCESS', message: '操作成功', data: { id: 1, name: '測試' } }

    case 'error':
      return { code: 'SQL_ERROR', message: '資料庫操作錯誤', data: null }

    case 'validation':
      return { code: 'VALIDATION_ERROR', message: '輸入資料驗證失敗', data: null }

    case 'permission':
      return { code: 'FORBIDDEN', message: '沒有權限', data: null }

    default:
      return { code: 'UNKNOWN_ERROR', message: '未知錯誤', data: null }
  }
}
</script>

<template>
  <div class="p-6 max-w-4xl mx-auto">
    <h1 class="text-2xl font-bold text-white mb-6">錯誤處理系統範例</h1>

    <!-- 成功訊息範例 -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">1. 成功訊息</h2>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
          @click="handleSaveSuccess"
        >
          顯示成功訊息
        </button>
      </div>
    </section>

    <!-- 錯誤訊息範例 -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">2. 錯誤訊息</h2>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
          @click="handleShowError"
        >
          顯示錯誤訊息
        </button>
        <button
          class="px-4 py-2 bg-yellow-500 text-white rounded-lg hover:bg-yellow-600 transition-colors"
          @click="handleShowWarning"
        >
          顯示警告訊息
        </button>
        <button
          class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
          @click="handleShowWarning"
        >
          顯示資訊訊息
        </button>
      </div>
    </section>

    <!-- 確認對話框範例 -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">3. 確認對話框</h2>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
          @click="handleDelete"
        >
          刪除（需確認）
        </button>
      </div>
    </section>

    <!-- API 請求範例 -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">4. API 請求處理</h2>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
          :disabled="apiLoading"
          @click="handleApiRequest"
        >
          {{ apiLoading ? '載入中...' : '自動處理錯誤' }}
        </button>
        <button
          class="px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors"
          @click="handleManualApi"
        >
          手動控制
        </button>
      </div>
    </section>

    <!-- 權限錯誤範例 -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">5. 特殊錯誤處理</h2>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 transition-colors"
          @click="handlePermissionError"
        >
          權限錯誤
        </button>
        <button
          class="px-4 py-2 bg-teal-500 text-white rounded-lg hover:bg-teal-600 transition-colors"
          @click="handleValidationError"
        >
          驗證錯誤
        </button>
        <button
          class="px-4 py-2 bg-pink-500 text-white rounded-lg hover:bg-pink-600 transition-colors"
          @click="handleScheduleConflicts"
        >
          排課衝突
        </button>
      </div>
    </section>

    <!-- 表單驗證範例 -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">6. 表單驗證</h2>

      <div class="space-y-4 max-w-md">
        <!-- 姓名欄位 -->
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-1">姓名</label>
          <input
            v-model="formData.name"
            type="text"
            class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
            :class="{ 'border-red-500': getFieldError('name') }"
          />
          <p v-if="getFieldError('name')" class="mt-1 text-sm text-red-400">
            {{ getFieldError('name') }}
          </p>
        </div>

        <!-- 電子郵件欄位 -->
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-1">電子郵件</label>
          <input
            v-model="formData.email"
            type="email"
            class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
            :class="{ 'border-red-500': getFieldError('email') }"
          />
          <p v-if="getFieldError('email')" class="mt-1 text-sm text-red-400">
            {{ getFieldError('email') }}
          </p>
        </div>

        <!-- 電話欄位 -->
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-1">電話</label>
          <input
            v-model="formData.phone"
            type="tel"
            class="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
          />
        </div>

        <button
          class="w-full px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
          @click="handleFormSubmit"
        >
          提交表單
        </button>
      </div>
    </section>

    <!-- 異步操作範例 -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">7. 異步操作</h2>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 bg-indigo-500 text-white rounded-lg hover:bg-indigo-600 transition-colors disabled:opacity-50"
          :disabled="loading"
          @click="handleAsyncOperation"
        >
          {{ loading ? '執行中...' : '執行異步操作' }}
        </button>
      </div>
    </section>

    <!-- Try/Catch 範例 -->
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">8. Try/Catch 模式</h2>
      <div class="flex gap-3">
        <button
          class="px-4 py-2 bg-cyan-500 text-white rounded-lg hover:bg-cyan-600 transition-colors"
          @click="handleWithTryCatch"
        >
          使用 Try/Catch
        </button>
      </div>
    </section>
  </div>
</template>
