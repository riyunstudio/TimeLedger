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
      localStorage.removeItem('token')
      localStorage.removeItem('admin_token')
      localStorage.removeItem('teacher_token')
      localStorage.removeItem('current_user_type')
      
      const currentPath = window.location.pathname
      if (!currentPath.includes('/login')) {
        router.push(`/admin/login?redirect=${encodeURIComponent(currentPath)}`)
      } else {
        router.push('/admin/login')
      }
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

  return { get, post, put, patch, delete: del }
}
