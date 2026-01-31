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
 * - errorHandler: For unified error handling
 */

import type { ApiResponse } from '~/types/api'
import { ErrorHandler, type ApiError, type NetworkError } from '~/utils/errorHandler'

export const useApi = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  // Import composable functions
  const { getAuthHeader, buildStandardHeaders } = useAuthHeaders()
  const { handleUnauthorized: clearAndRedirect } = useTokenManager()

  /**
   * Parses the API response and handles errors.
   *
   * @param response - The fetch Response object
   * @param rawResponse - Whether to return raw response without parsing
   * @returns The parsed response data
   * @throws ApiError for API-level errors (non-SUCCESS code)
   * @throws NetworkError for HTTP errors
   */
  async function parseResponse<T>(response: Response, rawResponse: boolean = false): Promise<T> {
    // Handle raw response (for file downloads, etc.)
    if (rawResponse) {
      return response.blob() as unknown as T
    }

    // Parse JSON response
    const text = await response.text()

    try {
      const data = JSON.parse(text) as ApiResponse<T>

      // Check for API-level error
      if (data.code !== 'SUCCESS') {
        const apiError: ApiError = {
          code: data.code,
          message: data.message || '操作失敗',
          data: data.data,
        }
        throw apiError
      }

      // Return data or datas based on response structure
      return (data.data ?? data.datas) as T
    } catch (error) {
      // Check if it's already an ApiError
      if (ErrorHandler.isApiError(error)) {
        throw error
      }

      // Try to parse as JSON error
      try {
        const errorData = JSON.parse(text)
        const networkError: NetworkError = {
          status: response.status,
          statusText: response.statusText,
          message: errorData.message || `HTTP ${response.status}`,
          originalError: error instanceof Error ? error : undefined,
        }
        throw networkError
      } catch {
        // Fallback to simple error
        const networkError: NetworkError = {
          status: response.status,
          statusText: response.statusText,
          message: text || `HTTP ${response.status}: ${response.statusText}`,
          originalError: error instanceof Error ? error : undefined,
        }
        throw networkError
      }
    }
  }

  /**
   * Checks the response and handles common error cases.
   * Automatically handles 401 Unauthorized by clearing tokens and redirecting.
   *
   * @param response - The fetch Response object
   * @throws ApiError for API-level errors
   * @throws NetworkError for HTTP errors
   */
  async function checkResponse(response: Response): Promise<void> {
    // Handle 401 Unauthorized
    if (response.status === 401) {
      const currentPath = process.client ? window.location.pathname : '/'
      clearAndRedirect(currentPath)

      const unauthorizedError: NetworkError = {
        status: 401,
        statusText: 'Unauthorized',
        message: '登入已過期，請重新登入',
      }
      throw unauthorizedError
    }

    // Handle other HTTP errors
    if (!response.ok) {
      const networkError: NetworkError = {
        status: response.status,
        statusText: response.statusText,
        message: `HTTP ${response.status}`,
      }
      throw networkError
    }
  }

  /**
   * Handles errors consistently across all API methods.
   *
   * @param error - The error that occurred
   * @param options - Error handling options
   * @throws Re-throws the error after handling
   */
  function handleError(error: unknown, options?: { showAlert?: boolean; context?: Record<string, unknown> }): never {
    if (ErrorHandler.isApiError(error)) {
      ErrorHandler.handleApiError(error as ApiError, {
        showAlert: options?.showAlert,
        context: options?.context,
      })
    } else if (ErrorHandler.isNetworkError(error)) {
      ErrorHandler.handleNetworkError(error as NetworkError, {
        showAlert: options?.showAlert,
        context: options?.context,
      })
    } else {
      ErrorHandler.handleUnknownError(error, {
        showAlert: options?.showAlert,
        context: options?.context,
      })
    }
    throw error
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
   * @param options - Request options (showAlert, context)
   * @returns Promise resolving to the response data
   */
  async function get<T>(
    endpoint: string,
    params?: Record<string, any>,
    options?: { showAlert?: boolean; context?: Record<string, unknown> }
  ): Promise<T> {
    const url = `${apiBase}${endpoint}${buildQueryString(params)}`
    const headers = buildStandardHeaders()

    try {
      const response = await fetch(url, { headers })
      await checkResponse(response)
      return await parseResponse<T>(response)
    } catch (error) {
      handleError(error, options)
    }
  }

  /**
   * Makes a POST request to the specified endpoint.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint
   * @param data - Request body data (will be JSON serialized)
   * @param options - Request options (showAlert, context)
   * @returns Promise resolving to the response data
   */
  async function post<T>(
    endpoint: string,
    data: any,
    options?: { showAlert?: boolean; context?: Record<string, unknown> }
  ): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers,
        body: JSON.stringify(data),
      })
      await checkResponse(response)
      return await parseResponse<T>(response)
    } catch (error) {
      handleError(error, options)
    }
  }

  /**
   * Makes a PUT request to the specified endpoint.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint
   * @param data - Request body data
   * @param options - Request options (showAlert, context)
   * @returns Promise resolving to the response data
   */
  async function put<T>(
    endpoint: string,
    data: any,
    options?: { showAlert?: boolean; context?: Record<string, unknown> }
  ): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    try {
      const response = await fetch(url, {
        method: 'PUT',
        headers,
        body: JSON.stringify(data),
      })
      await checkResponse(response)
      return await parseResponse<T>(response)
    } catch (error) {
      handleError(error, options)
    }
  }

  /**
   * Makes a PATCH request to the specified endpoint.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint
   * @param data - Request body data
   * @param options - Request options (showAlert, context)
   * @returns Promise resolving to the response data
   */
  async function patch<T>(
    endpoint: string,
    data: any,
    options?: { showAlert?: boolean; context?: Record<string, unknown> }
  ): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    try {
      const response = await fetch(url, {
        method: 'PATCH',
        headers,
        body: JSON.stringify(data),
      })
      await checkResponse(response)
      return await parseResponse<T>(response)
    } catch (error) {
      handleError(error, options)
    }
  }

  /**
   * Makes a DELETE request to the specified endpoint.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint
   * @param options - Request options (showAlert, context)
   * @returns Promise resolving to the response data
   */
  async function del<T>(
    endpoint: string,
    options?: { showAlert?: boolean; context?: Record<string, unknown> }
  ): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    try {
      const response = await fetch(url, {
        method: 'DELETE',
        headers,
      })
      await checkResponse(response)
      return await parseResponse<T>(response)
    } catch (error) {
      handleError(error, options)
    }
  }

  /**
   * Uploads a file to the specified endpoint using multipart/form-data.
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint for file upload
   * @param file - The File object to upload
   * @param fieldName - The form field name for the file (default: 'file')
   * @param options - Request options (showAlert, context)
   * @returns Promise resolving to the upload response
   */
  async function upload<T>(
    endpoint: string,
    file: File,
    fieldName: string = 'file',
    options?: { showAlert?: boolean; context?: Record<string, unknown> }
  ): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = getAuthHeader() // No Content-Type for FormData

    const formData = new FormData()
    formData.append(fieldName, file)

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers,
        body: formData,
      })
      await checkResponse(response)
      return await parseResponse<T>(response)
    } catch (error) {
      handleError(error, options)
    }
  }

  /**
   * Makes a request that returns raw binary data (e.g., file download).
   *
   * @typeParam T - The expected response data type
   * @param endpoint - API endpoint
   * @param options - Request options (showAlert, context)
   * @returns Promise resolving to the raw response
   */
  async function raw<T>(
    endpoint: string,
    options?: { showAlert?: boolean; context?: Record<string, unknown> }
  ): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    try {
      const response = await fetch(url, { headers })
      await checkResponse(response)
      return await parseResponse<T>(response, true)
    } catch (error) {
      handleError(error, options)
    }
  }

  return {
    get,
    post,
    put,
    patch,
    delete: del,
    upload,
    raw,
  }
}
