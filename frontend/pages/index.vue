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

const handleLineLogin = async () => {
  loading.value = true

  try {
    // 導向教師登入頁面，由 SDK 統一處理登入邏輯
    router.push('/teacher/login')
  } catch (err) {
    console.error('Login redirect failed:', err)
    loading.value = false
  }
}
</script>
