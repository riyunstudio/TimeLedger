export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()

  if (to.path === '/admin/login') {
    return
  }

  authStore.initFromStorage()

  if (!authStore.isAuthenticated) {
    return navigateTo('/admin/login')
  }

  if (!authStore.isAdmin) {
    return navigateTo('/teacher/dashboard')
  }
})
