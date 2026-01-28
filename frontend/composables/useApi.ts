export const useApi = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase
  const router = useRouter()

  const getAuthHeader = (): Record<string, string> => {
    if (process.client) {
      const token = localStorage.getItem('admin_token') || localStorage.getItem('teacher_token') || localStorage.getItem('token')
      return token ? { Authorization: `Bearer ${token}` } : {}
    }
    return {}
  }

  const handleUnauthorized = () => {
    if (process.client) {
      // 先保存當前用戶類型（因為馬上要清除）
      const currentUserType = localStorage.getItem('current_user_type')
      const wasTeacher = localStorage.getItem('teacher_token') || currentUserType === 'teacher'
      const wasAdmin = localStorage.getItem('admin_token') || currentUserType === 'admin'

      // 清除所有認證資料
      localStorage.removeItem('token')
      localStorage.removeItem('admin_token')
      localStorage.removeItem('teacher_token')
      localStorage.removeItem('current_user_type')

      const currentPath = window.location.pathname

      // 根據之前的用戶類型決定轉向哪個登入頁
      let redirectPath = '/'
      if (wasTeacher && !currentPath.includes('/teacher/login')) {
        redirectPath = `/teacher/login?redirect=${encodeURIComponent(currentPath)}`
      } else if (wasAdmin && !currentPath.includes('/admin/login')) {
        redirectPath = `/admin/login?redirect=${encodeURIComponent(currentPath)}`
      } else if (!currentPath.includes('/login')) {
        // 無法判斷用戶類型，轉到首頁
        redirectPath = '/'
      } else {
        // 已經在登入頁，轉到首頁
        redirectPath = '/'
      }

      router.push(redirectPath)
    }
  }

  const checkResponse = async (response: Response) => {
    if (response.status === 401) {
      handleUnauthorized()
      throw new Error('Unauthorized')
    }
    if (!response.ok) {
      // 嘗試讀取後端返回的錯誤訊息
      const errorText = await response.text()
      let errorMessage = `HTTP ${response.status}`
      try {
        const errorData = JSON.parse(errorText)
        if (errorData.message) {
          errorMessage = errorData.message
        }
      } catch {
        // 如果解析失敗，使用原始錯誤文字
        if (errorText) {
          errorMessage = errorText
        }
      }
      throw new Error(errorMessage)
    }
  }

  const get = async <T>(endpoint: string, params?: Record<string, any>): Promise<T> => {
    let url = `${apiBase}${endpoint}`
    if (params) {
      const searchParams = new URLSearchParams()
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          searchParams.append(key, String(value))
        }
      })
      const queryString = searchParams.toString()
      if (queryString) {
        url += `?${queryString}`
      }
    }
    const headers: Record<string, string> = {
      ...getAuthHeader(),
      'Content-Type': 'application/json',
    }
    const response = await fetch(url, { headers })
    await checkResponse(response)
    return response.json()
  }

  const post = async <T>(endpoint: string, data: any): Promise<T> => {
    const headers: Record<string, string> = {
      ...getAuthHeader(),
      'Content-Type': 'application/json',
    }
    const response = await fetch(`${apiBase}${endpoint}`, {
      method: 'POST',
      headers,
      body: JSON.stringify(data),
    })
    await checkResponse(response)
    return response.json()
  }

  const put = async <T>(endpoint: string, data: any): Promise<T> => {
    const headers: Record<string, string> = {
      ...getAuthHeader(),
      'Content-Type': 'application/json',
    }
    const response = await fetch(`${apiBase}${endpoint}`, {
      method: 'PUT',
      headers,
      body: JSON.stringify(data),
    })
    await checkResponse(response)
    return response.json()
  }

  const patch = async <T>(endpoint: string, data: any): Promise<T> => {
    const headers: Record<string, string> = {
      ...getAuthHeader(),
      'Content-Type': 'application/json',
    }
    const response = await fetch(`${apiBase}${endpoint}`, {
      method: 'PATCH',
      headers,
      body: JSON.stringify(data),
    })
    await checkResponse(response)
    return response.json()
  }

  const del = async <T>(endpoint: string): Promise<T> => {
    const headers: Record<string, string> = {
      ...getAuthHeader(),
      'Content-Type': 'application/json',
    }
    const response = await fetch(`${apiBase}${endpoint}`, {
      method: 'DELETE',
      headers,
    })
    await checkResponse(response)
    return response.json()
  }

  const upload = async <T>(endpoint: string, file: File, fieldName: string = 'file'): Promise<T> => {
    const headers: Record<string, string> = {
      ...getAuthHeader(),
    }
    const formData = new FormData()
    formData.append(fieldName, file)

    const response = await fetch(`${apiBase}${endpoint}`, {
      method: 'POST',
      headers,
      body: formData,
    })
    await checkResponse(response)
    return response.json()
  }

  return { get, post, put, patch, delete: del, upload }
}
