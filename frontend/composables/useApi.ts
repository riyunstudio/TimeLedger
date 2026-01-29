/**
 * useApi Composable
 * Centralized API client for making HTTP requests to the backend.
 *
 * This composable provides a unified interface for all API calls,
 * handling authentication headers, error responses, and token management.
 *
 * Composables used:
 * - useAuthHeaders: For building authorization headers
 * - useTokenManager: For token storage and unauthorized handling
 */

export const useApi = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  // Import composable functions
  const { getAuthHeader, buildStandardHeaders } = useAuthHeaders()
  const { handleUnauthorized: clearAndRedirect } = useTokenManager()

  /**
   * Checks the response and handles common error cases.
   * Automatically handles 401 Unauthorized by clearing tokens and redirecting.
   *
   * @param response - The fetch Response object
   * @throws Error with descriptive message for non-2xx responses
   */
  async function checkResponse(response: Response): Promise<void> {
    if (response.status === 401) {
      const currentPath = process.client ? window.location.pathname : '/'
      clearAndRedirect(currentPath)
      throw new Error('Unauthorized')
    }

    if (!response.ok) {
      // Attempt to parse error message from backend response
      const errorText = await response.text()
      let errorMessage = `HTTP ${response.status}`

      try {
        const errorData = JSON.parse(errorText)
        if (errorData.message) {
          errorMessage = errorData.message
        }
      } catch {
        // Use raw text if JSON parsing fails
        if (errorText) {
          errorMessage = errorText
        }
      }

      throw new Error(errorMessage)
    }
  }

  /**
   * Builds query string from parameters object.
   * Filters out undefined, null, and empty string values.
   *
   * @param params - Object containing query parameters
   * @returns Query string with leading '?' or empty string if no params
   */
  function buildQueryString(params?: Record<string, any>): string {
    if (!params) {
      return ''
    }

    const searchParams = new URLSearchParams()

    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        searchParams.append(key, String(value))
      }
    })

    const queryString = searchParams.toString()
    return queryString ? `?${queryString}` : ''
  }

  /**
   * Makes a GET request to the specified endpoint.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint (will be appended to apiBase)
   * @param params - Optional query parameters
   * @returns Promise resolving to the response data
   */
  async function get<T>(endpoint: string, params?: Record<string, any>): Promise<T> {
    const url = `${apiBase}${endpoint}${buildQueryString(params)}`
    const headers = buildStandardHeaders()

    const response = await fetch(url, { headers })
    await checkResponse(response)

    return response.json()
  }

  /**
   * Makes a POST request to the specified endpoint.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint
   * @param data - Request body data (will be JSON serialized)
   * @returns Promise resolving to the response data
   */
  async function post<T>(endpoint: string, data: any): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify(data),
    })
    await checkResponse(response)

    return response.json()
  }

  /**
   * Makes a PUT request to the specified endpoint.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint
   * @param data - Request body data
   * @returns Promise resolving to the response data
   */
  async function put<T>(endpoint: string, data: any): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    const response = await fetch(url, {
      method: 'PUT',
      headers,
      body: JSON.stringify(data),
    })
    await checkResponse(response)

    return response.json()
  }

  /**
   * Makes a PATCH request to the specified endpoint.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint
   * @param data - Request body data
   * @returns Promise resolving to the response data
   */
  async function patch<T>(endpoint: string, data: any): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    const response = await fetch(url, {
      method: 'PATCH',
      headers,
      body: JSON.stringify(data),
    })
    await checkResponse(response)

    return response.json()
  }

  /**
   * Makes a DELETE request to the specified endpoint.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint
   * @returns Promise resolving to the response data
   */
  async function del<T>(endpoint: string): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    const response = await fetch(url, {
      method: 'DELETE',
      headers,
    })
    await checkResponse(response)

    return response.json()
  }

  /**
   * Uploads a file to the specified endpoint using multipart/form-data.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint for file upload
   * @param file - The File object to upload
   * @param fieldName - The form field name for the file (default: 'file')
   * @returns Promise resolving to the upload response
   */
  async function upload<T>(endpoint: string, file: File, fieldName: string = 'file'): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = getAuthHeader() // No Content-Type for FormData

    const formData = new FormData()
    formData.append(fieldName, file)

    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: formData,
    })
    await checkResponse(response)

    return response.json()
  }

  return {
    get,
    post,
    put,
    patch,
    delete: del,
    upload,
  }
}
