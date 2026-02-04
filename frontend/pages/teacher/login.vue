<template>
  <div class="min-h-screen bg-gray-900 flex items-center justify-center p-4">
    <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-full max-w-md border border-gray-700">
      <h1 class="text-2xl font-bold mb-6 text-center text-gray-100">老師登入</h1>

      <div class="mb-6 p-4 bg-blue-900/30 border border-blue-800 rounded-lg">
        <p class="text-sm text-blue-300">
          請使用 LINE 帳號登入。輸入您的 LINE User ID。
        </p>
      </div>

      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block text-gray-300 mb-2">LINE User ID</label>
          <input
            v-model="lineUserId"
            type="text"
            class="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="輸入 LINE User ID"
            required
          />
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:bg-gray-600 transition-colors"
        >
          {{ loading ? '登入中...' : '登入' }}
        </button>
      </form>

      <div v-if="error" class="mt-4 p-3 bg-red-900/50 text-red-300 rounded-lg border border-red-800">
        {{ error }}
      </div>

      <div v-if="success" class="mt-4 p-3 bg-green-900/50 text-green-300 rounded-lg border border-green-800">
        登入成功，正在跳轉...
      </div>

      <div class="mt-6 text-center">
        <NuxtLink to="/" class="text-blue-400 hover:text-blue-300 transition-colors">
          回首頁
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const lineUserId = ref(route.query.line_user_id as string || '')
const loading = ref(false)
const error = ref('')
const success = ref(false)

async function handleLogin() {
  if (!lineUserId.value) {
    error.value = '請輸入 LINE User ID'
    return
  }

  loading.value = true
  error.value = ''
  success.value = false

  try {
    const response = await $fetch('/api/v1/auth/teacher/line/login', {
      method: 'POST',
      body: {
        line_user_id: lineUserId.value
      }
    })

    // 檢查 response.data 或 response.datas
    const responseData = (response as any).data || (response as any).datas
    const responseCode = (response as any).code

    if (responseCode === 0 && responseData) {
      const token = responseData.token
      const user = responseData.user

      // 設置 authStore 和 localStorage
      authStore.login({
        token,
        refresh_token: '',
        teacher: user,
      })

      success.value = true

      setTimeout(() => {
        router.push('/teacher/dashboard')
      }, 1000)
    } else {
      error.value = (response as any)?.message || '登入失敗'
    }
  } catch (err: any) {
    console.error('Login error:', err)
    error.value = err.data?.message || err.message || '登入失敗，請稍後再試'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  const token = localStorage.getItem('teacher_token')
  if (token) {
    router.push('/teacher/dashboard')
  }
})
</script>
