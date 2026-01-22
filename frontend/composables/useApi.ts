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

  const checkResponse = (response: Response) => {
    if (response.status === 401) {
      handleUnauthorized()
      throw new Error('Unauthorized')
    }
    if (!response.ok) throw new Error(`HTTP ${response.status}`)
  }

  const get = async <T>(endpoint: string): Promise<T> => {
    const headers: Record<string, string> = {
      ...getAuthHeader(),
      'Content-Type': 'application/json',
    }
    const response = await fetch(`${apiBase}${endpoint}`, { headers })
    checkResponse(response)
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
    checkResponse(response)
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
    checkResponse(response)
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
    checkResponse(response)
    return response.json()
  }

  return { get, post, put, delete: del }
}
