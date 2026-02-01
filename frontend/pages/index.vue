<template>
  <div class="min-h-screen bg-gradient-mesh flex flex-col">
    <LandingHeroSection :loading="loading" @line-login="handleLineLogin" />

    <LandingDemoSandbox />

    <LandingFeatureShowcase :loading="loading" @line-login="handleLineLogin" />
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
})

const router = useRouter()

const loading = ref(false)
const error = ref('')

const handleLineLogin = async () => {
  loading.value = true

  try {
    // 跳轉到教師登入頁面，帶上預設測試老師的 LINE User ID 和 Access Token
    const demoLineUserId = 'U1234567890abcdef1234567890abcd'
    const demoAccessToken = 'test_access_token_demo001'
    router.push(`/teacher/login?line_user_id=${demoLineUserId}&access_token=${demoAccessToken}`)
  } catch (err) {
    console.error('Login failed:', err)
    error.value = '登入失敗，請稍後再試'
  } finally {
    loading.value = false
  }
}
</script>
