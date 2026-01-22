export default defineNuxtRouteMiddleware((to, from) => {
  if (to.path.startsWith('/admin')) {
    return
  }
})
