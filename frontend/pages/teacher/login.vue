<template>
  <div class="min-h-screen bg-gray-900 flex items-center justify-center">
    <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-full max-w-md border border-gray-700">
      <h1 class="text-2xl font-bold mb-6 text-center text-gray-100">老師登入</h1>
      
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
        
        <div class="mb-6">
          <label class="block text-gray-300 mb-2">Access Token</label>
          <input
            v-model="accessToken"
            type="text"
            class="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="輸入 Access Token"
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
      
      <div class="mt-4 p-4 bg-gray-700/50 rounded-lg text-sm border border-gray-600">
        <p class="font-bold text-gray-200 mb-2">測試帳號:</p>
        <p class="text-gray-300">LINE User ID: <code class="bg-gray-600 px-2 py-0.5 rounded text-blue-300">LINE_TEACHER_001</code></p>
        <p class="text-gray-300 mt-1">Access Token: <code class="bg-gray-600 px-2 py-0.5 rounded text-blue-300">mock_token</code></p>
      </div>
    </div>
  </div>
</template>

<script setup>
 definePageMeta({
   layout: false,
 })

 const router = useRouter()

const lineUserId = ref('LINE_TEACHER_001')
const accessToken = ref('mock_token')
const loading = ref(false)
const error = ref('')
const success = ref(false)

async function handleLogin() {
  loading.value = true
  error.value = ''
  success.value = false
  
  try {
    const response = await $fetch('/api/v1/auth/teacher/line/login', {
      method: 'POST',
      body: {
        line_user_id: lineUserId.value,
        access_token: accessToken.value
      }
    })
    
    console.log('Login response:', response)
    
    if (response && response.code === 0 && response.datas) {
      const token = response.datas.token
      const user = response.datas.user
      
      localStorage.setItem('teacher_token', token)
      localStorage.setItem('teacher_user', JSON.stringify(user))
      localStorage.setItem('current_user_type', 'teacher')
      
      success.value = true
      
      setTimeout(() => {
        router.push('/teacher/dashboard')
      }, 1000)
    } else {
      error.value = response?.message || '登入失敗'
    }
  } catch (err) {
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
