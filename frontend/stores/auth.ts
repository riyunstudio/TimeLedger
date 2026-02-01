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
    const userData = user.value as any
    return !!user.value && !userData?.user_type && !userData?.role
  })

  const isAdmin = computed(() => {
    if (!user.value) {
      console.log('[DEBUG isAdmin] user.value is null, returning false')
      return false
    }
    const userData = user.value as any
    const role = userData.user_type || userData.role
    console.log('[DEBUG isAdmin] userData:', JSON.stringify(userData))
    console.log('[DEBUG isAdmin] role:', role)
    const result = role === 'ADMIN' || role === 'OWNER' || role === 'STAFF'
    console.log('[DEBUG isAdmin] result:', result)
    return result
  })

  // Actions
  function login(authData: AuthResponse) {
    console.log('[DEBUG login] authData:', JSON.stringify(authData))
    console.log('[DEBUG login] authData.token:', authData.token)
    console.log('[DEBUG login] authData.user:', JSON.stringify(authData.user))

    // 保存 token 到 localStorage（用於 API 請求認證）
    const userType = authData.user?.user_type || authData.teacher?.user_type || 'ADMIN'
    const storageKey = userType === 'ADMIN' || userType === 'OWNER' || userType === 'STAFF' 
      ? 'admin_token' 
      : 'teacher_token'
    localStorage.setItem(storageKey, authData.token)
    console.log('[DEBUG login] saved token to localStorage:', storageKey)

    token.value = authData.token
    refreshToken.value = authData.refresh_token

    if (authData.teacher) {
      user.value = authData.teacher
    } else if (authData.user) {
      user.value = authData.user as any
      console.log('[DEBUG login] user.value after assignment:', JSON.stringify(user.value))
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
