/**
 * 統一認證中介層
 *
 * 讀取頁面的 definePageMeta({ auth: '...' }) 設定，
 * 根據 authStore.isAdmin 或 authStore.isTeacher 判斷是否有權存取。
 *
 * 使用方式：
 * - 在頁面中設定 definePageMeta({ auth: 'ADMIN' }) 或 definePageMeta({ auth: 'TEACHER' })
 * - 不需要 auth 設定的頁面會被視為 PUBLIC
 */

export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()

  // 從頁面 meta 取得 auth 設定，預設為 PUBLIC
  const requiredAuth = (to.meta.auth as string) || 'PUBLIC'

  // 【新增】LINE OAuth 回調檢測：如果 URL 包含 code 和 state 參數，強制重導向到登入頁
  // 這樣可以確保無論 LINE 跳回哪一頁，都會被強制送往登入處理頁
  const hasCode = to.query.code !== undefined && to.query.code !== null && to.query.code !== ''
  const hasState = to.query.state !== undefined && to.query.state !== null && to.query.state !== ''

  if (hasCode && hasState && to.path !== '/teacher/login') {
    // 這是 LINE OAuth 回調，強制重導向到老師登入頁處理
    return navigateTo({
      path: '/teacher/login',
      query: {
        code: to.query.code,
        state: to.query.state
      }
    })
  }

  // 登入頁面允許所有人存取
  if (to.path === '/admin/login' || to.path === '/teacher/login') {
    return
  }

  // 判斷是否有登入
  if (!authStore.isAuthenticated) {
    // 未登入時，根據頁面類型導向對應的登入頁
    if (requiredAuth === 'ADMIN') {
      return navigateTo('/admin/login')
    } else if (requiredAuth === 'TEACHER') {
      return navigateTo('/teacher/login')
    } else {
      // PUBLIC 頁面不需要登入
      return
    }
  }

  // 已登入時，檢查權限是否足夠
  if (requiredAuth === 'ADMIN') {
    if (!authStore.isAdmin) {
      // 已是登入狀態但權限不足，導向其所屬角色的儀表板
      if (authStore.isTeacher) {
        return navigateTo('/teacher/dashboard')
      }
      return navigateTo('/admin/login')
    }
  } else if (requiredAuth === 'TEACHER') {
    if (!authStore.isTeacher) {
      // 已是登入狀態但權限不足，導向其所屬角色的儀表板
      if (authStore.isAdmin) {
        return navigateTo('/admin/dashboard')
      }
      return navigateTo('/teacher/login')
    }
  }

  // 權限檢查通過
  return
})
