import { defineStore } from 'pinia'
import type { Teacher, AuthResponse } from '~/types'
import { withLoading } from '~/utils/loadingHelper'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<Teacher | null>(null)
  const token = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)
  const isRefreshing = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!token.value)

  const isTeacher = computed(() => {
    if (!user.value) return false
    const userData = user.value as any
    const role = userData.user_type || userData.role
    return role !== 'ADMIN' && role !== 'OWNER' && role !== 'STAFF'
  })

  const isAdmin = computed(() => {
    if (!user.value) {
      return false
    }
    const userData = user.value as any
    const role = userData.user_type || userData.role
    return role === 'ADMIN' || role === 'OWNER' || role === 'STAFF'
  })

  // Actions
  function login(authData: AuthResponse) {
    // 保存 token 到 localStorage（用於 API 請求認證）
    const userType = authData.user?.user_type || (authData.teacher ? 'TEACHER' : 'ADMIN')
    const storageKey = userType === 'ADMIN' || userType === 'OWNER' || userType === 'STAFF' 
      ? 'admin_token' 
      : 'teacher_token'
    localStorage.setItem(storageKey, authData.token)

    token.value = authData.token
    refreshToken.value = authData.refresh_token

    if (authData.teacher) {
      user.value = authData.teacher
    } else if (authData.user) {
      user.value = authData.user as any
    }
  }

  function logout() {
    // 清除 localStorage 中的 token
    const userType = user.value?.user_type
    const storageKey = userType === 'ADMIN' || userType === 'OWNER' || userType === 'STAFF'
      ? 'admin_token'
      : 'teacher_token'
    localStorage.removeItem(storageKey)
    localStorage.removeItem('token') // 清除可能的其他 token key

    user.value = null
    token.value = null
    refreshToken.value = null
  }

  async function refreshAccessToken(): Promise<boolean> {
    if (!refreshToken.value) return false

    return withLoading(isRefreshing.value, async () => {
      try {
        const api = useApi()
        const response = await api.post<{ code: number; message: string; datas: { token: string; refresh_token: string } }>('/auth/refresh', {
          refresh_token: refreshToken.value,
        })

        token.value = response.datas?.token || ''
        refreshToken.value = response.datas?.refresh_token || ''

        return true
      } catch (error) {
        logout()
        return false
      }
    })
  }

  return {
    user,
    token,
    refreshToken,
    isLoading,
    isRefreshing,
    isAuthenticated,
    isTeacher,
    isAdmin,
    login,
    logout,
    refreshAccessToken,
  }
}, {
  persist: {
    paths: ['user', 'token', 'refreshToken'],
  },
})
