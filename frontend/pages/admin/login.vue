<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-mesh p-4">
    <div class="glass-card p-8 max-w-md w-full">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-slate-100 mb-2">
          管理員登入
        </h1>
        <p class="text-slate-400">TimeLedger 中心後台</p>
      </div>

       <form @submit.prevent="handleLogin" class="space-y-6">
         <div>
           <label class="block text-slate-300 mb-2 font-medium">
             Email
           </label>
           <input
             v-model="email"
             type="email"
             placeholder="admin@example.com"
             class="input-field"
             required
           />
         </div>

         <div>
           <label class="block text-slate-300 mb-2 font-medium">
             密碼
           </label>
           <input
             v-model="password"
             type="password"
             placeholder="••••••••"
             class="input-field"
             required
           />
         </div>

         <button
           type="submit"
           :disabled="loading"
           class="w-full btn-primary flex items-center justify-center gap-2"
         >
           <svg v-if="loading" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
             <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
             <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
           </svg>
           {{ loading ? '登入中...' : '登入' }}
         </button>
       </form>

      <!-- 快速登入測試區域 -->
      <div class="mt-8 p-4 bg-slate-800/50 rounded-lg border border-slate-700">
        <p class="text-sm text-slate-400 mb-3 text-center">🧪 測試快速登入</p>
        <div class="grid grid-cols-3 gap-2">
          <button
            @click="quickLogin('admin@timeledger.com', 'admin123')"
            class="px-3 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm rounded transition-colors"
          >
            擁有者
          </button>
          <button
            @click="quickLogin('manager@timeledger.com', 'admin123')"
            class="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded transition-colors"
          >
            管理員
          </button>
          <button
            @click="quickLogin('staff@timeledger.com', 'admin123')"
            class="px-3 py-2 bg-purple-600 hover:bg-purple-700 text-white text-sm rounded transition-colors"
          >
            工作人員
          </button>
        </div>
      </div>

      <div class="mt-6 text-center">
        <NuxtLink
          to="/teacher/login"
          class="text-slate-400 hover:text-primary-500 transition-colors duration-300"
        >
          老師登入請點此
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertError } from '~/composables/useAlert'

definePageMeta({
  layout: false,
})

const authStore = useAuthStore()
const router = useRouter()

const email = ref('admin@timeledger.com')
const password = ref('admin123')
const loading = ref(false)

const quickLogin = (testEmail: string, testPassword: string) => {
  email.value = testEmail
  password.value = testPassword
  handleLogin()
}

const handleLogin = async () => {
  loading.value = true

  try {
    const api = useApi()
    const response = await api.post<{ code: number; message: string; datas: any }>('/auth/admin/login', {
      email: email.value,
      password: password.value,
    })

    authStore.login(response.datas)
    router.push('/admin/dashboard')
  } catch (error) {
    console.error('Login failed:', error)
    await alertError('登入失敗，請檢查 Email 和密碼')
  } finally {
    loading.value = false
  }
}
</script>
