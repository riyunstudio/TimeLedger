import { defineStore } from 'pinia'
import type { Teacher, AuthResponse } from '~/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<Teacher | null>(null)
  const token = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isAuthenticated = computed(() => !!token.value)
  const isTeacher = computed(() => !!user.value && !isAdmin.value)
  const isAdmin = computed(() => {
    if (!user.value) {
      console.log('[DEBUG isAdmin] user.value is null, returning false')
      return false
    }
    const userData = user.value as any
    // 同時支援 user_type（來自登入 API）和 role（來自 /admin/me/profile API）
    const role = userData.user_type || userData.role
    console.log('[DEBUG isAdmin] userData:', JSON.stringify(userData))
    console.log('[DEBUG isAdmin] role:', role)
    const result = role === 'ADMIN' || role === 'OWNER' || role === 'STAFF'
    console.log('[DEBUG isAdmin] result:', result)
    return result
  })

  const login = (authData: AuthResponse) => {
    console.log('[DEBUG login] authData:', JSON.stringify(authData))
    console.log('[DEBUG login] authData.token:', authData.token)
    console.log('[DEBUG login] authData.user:', JSON.stringify(authData.user))

    token.value = authData.token
    refreshToken.value = authData.refresh_token

    if (authData.teacher) {
      user.value = authData.teacher
      localStorage.setItem('teacher_user', JSON.stringify(authData.teacher))
      localStorage.setItem('teacher_token', authData.token)
      localStorage.setItem('teacher_refresh_token', authData.refresh_token || '')
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_refresh_token')
      localStorage.removeItem('admin_user')
    } else if (authData.user) {
      user.value = authData.user as any
      console.log('[DEBUG login] user.value after assignment:', JSON.stringify(user.value))
      localStorage.setItem('admin_user', JSON.stringify(authData.user))
      localStorage.setItem('admin_token', authData.token)
      console.log('[DEBUG login] saved admin_token:', authData.token)
      localStorage.setItem('admin_refresh_token', authData.refresh_token || '')
      localStorage.removeItem('teacher_token')
      localStorage.removeItem('teacher_refresh_token')
      localStorage.removeItem('teacher_user')
    }

    localStorage.setItem('current_user_type', authData.teacher ? 'teacher' : 'admin')
  }

  const logout = () => {
    user.value = null
    token.value = null
    refreshToken.value = null
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_refresh_token')
    localStorage.removeItem('admin_user')
    localStorage.removeItem('teacher_token')
    localStorage.removeItem('teacher_refresh_token')
    localStorage.removeItem('teacher_user')
    localStorage.removeItem('current_user_type')
  }

  const initFromStorage = () => {
    const userType = localStorage.getItem('current_user_type')

    if (userType === 'admin') {
      const storedToken = localStorage.getItem('admin_token')
      const storedUser = localStorage.getItem('admin_user')
      const storedRefresh = localStorage.getItem('admin_refresh_token')

      if (storedToken) {
        token.value = storedToken
      }
      if (storedRefresh) {
        refreshToken.value = storedRefresh
      }
      if (storedUser) {
        try {
          user.value = JSON.parse(storedUser)
        } catch (e) {
          localStorage.removeItem('admin_user')
        }
      }
    } else if (userType === 'teacher') {
      const storedToken = localStorage.getItem('teacher_token')
      const storedUser = localStorage.getItem('teacher_user')
      const storedRefresh = localStorage.getItem('teacher_refresh_token')

      if (storedToken) {
        token.value = storedToken
      }
      if (storedRefresh) {
        refreshToken.value = storedRefresh
      }
      if (storedUser) {
        try {
          user.value = JSON.parse(storedUser)
        } catch (e) {
          localStorage.removeItem('teacher_user')
        }
      }
    }
  }

  const refreshAccessToken = async () => {
    if (!refreshToken.value) return false

    try {
      const api = useApi()
      const response = await api.post<{ code: number; message: string; datas: { token: string; refresh_token: string } }>('/auth/refresh', {
        refresh_token: refreshToken.value,
      })

      token.value = response.datas?.token || ''
      refreshToken.value = response.datas?.refresh_token || ''

      const userType = localStorage.getItem('current_user_type')
      if (userType === 'admin') {
        localStorage.setItem('admin_token', response.datas?.token || '')
        localStorage.setItem('admin_refresh_token', response.datas?.refresh_token || '')
      } else if (userType === 'teacher') {
        localStorage.setItem('teacher_token', response.datas?.token || '')
        localStorage.setItem('teacher_refresh_token', response.datas?.refresh_token || '')
      }

      return true
    } catch (error) {
      logout()
      return false
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    isTeacher,
    isAdmin,
    login,
    logout,
    refreshAccessToken,
    initFromStorage,
  }
})
