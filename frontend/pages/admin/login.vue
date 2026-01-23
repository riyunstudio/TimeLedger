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

       <div class="mt-4 pt-4 border-t border-white/10">
         <button
           @click="handleMockLogin"
           :disabled="loading"
           class="w-full glass-btn flex items-center justify-center gap-2 text-sm"
         >
           <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
             <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
           </svg>
           Mock 登入 (無需後端)
         </button>
       </div>

      <div class="mt-6 text-center">
        <NuxtLink
          to="/"
          class="text-slate-400 hover:text-primary-500 transition-colors duration-300"
        >
          老師登入請點此
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

const authStore = useAuthStore()
const router = useRouter()

const email = ref('admin@timeledger.com')
const password = ref('admin123')
const loading = ref(false)

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
    alert('登入失敗，請檢查 Email 和密碼')
  } finally {
    loading.value = false
  }
}

const handleMockLogin = () => {
  loading.value = true

  setTimeout(() => {
    authStore.login({
      user: {
        id: 1,
        name: 'Mock Admin',
        email: 'admin@example.com',
        role: 'admin',
      },
      token: 'mock-token-' + Date.now(),
    })
    router.push('/admin/dashboard')
    loading.value = false
  }, 500)
}
</script>
