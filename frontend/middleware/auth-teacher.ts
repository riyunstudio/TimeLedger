export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()

  if (to.path === '/' || to.path.startsWith('/teacher/login')) {
    return
  }

  authStore.initFromStorage()

  if (!authStore.isAuthenticated) {
    return navigateTo('/teacher/login')
  }

  if (!authStore.isTeacher) {
    return navigateTo('/admin/dashboard')
  }
})
