export default defineNuxtPlugin(() => {
  const authStore = useAuthStore()

  // 認證狀態已由 pinia-plugin-persistedstate 自動恢復
  // 不需要額外的初始化邏輯
})
