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

const email = ref('')
const password = ref('')
const loading = ref(false)

const handleLogin = async () => {
  loading.value = true

  try {
    const api = useApi()
    // api.post 返回的是 datas 字段的值，即 LoginResponse（包含 token 和 user）
    const response = await api.post<{ token: string; user: any }>('/auth/admin/login', {
      email: email.value,
      password: password.value,
    })

    authStore.login(response)
    router.push('/admin/dashboard')
  } catch (error) {
    console.error('Login failed:', error)
    await alertError('登入失敗，請檢查 Email 和密碼')
  } finally {
    loading.value = false
  }
}
</script>
