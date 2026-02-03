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
 *
 * Schema Validation (Optional):
 * - Pass a Zod schema as the 4th parameter to validate responses
 * - Validation only runs in development mode
 * - Validation failures are logged to console.warn but don't throw
 */

import type { ApiResponse } from '~/types/api'
import type { ZodSchema, ZodError } from 'zod'
import { isSuccessCode } from '~/constants/errorCodes'
import { ErrorHandler, type ApiError, type NetworkError } from '~/utils/errorHandler'

export const useApi = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  // Import composable functions
  const { getAuthHeader, buildStandardHeaders } = useAuthHeaders()
  const { handleUnauthorized: clearAndRedirect } = useTokenManager()

  /**
   * 驗證開發環境
   */
  const isDev = () => process.dev

  /**
   * 取得驗證失敗的視覺化樣式 (開發環境專用)
   */
  const getValidationWarningStyle = () =>
    `
    background: #fff3cd;
    border: 2px solid #ffc107;
    border-radius: 8px;
    padding: 12px;
    margin: 8px 0;
    font-family: 'Consolas', 'Monaco', monospace;
    font-size: 13px;
    line-height: 1.6;
    color: #856404;
    `.trim()

  /**
   * 取得驗證問題的詳細格式化 (開發環境專用)
   */
  const formatValidationIssues = (issues: Array<{
    path: Array<string | number>
    message: string
    expected?: string
    received?: string
  }>): string => {
    return issues.map((issue) => {
      const path = issue.path.length > 0 ? `.${issue.path.join('.')}` : '(root)'
      let detail = `  🔸 Path: ${path}`
      detail += `\n     Message: ${issue.message}`
      if (issue.expected && issue.received) {
        detail += `\n     Expected: ${issue.expected}`
        detail += `\n     Received: ${issue.received}`
      }
      return detail
    }).join('\n')
  }

  /**
   * Schema 驗證函數 (開發環境專用)
   *
   * @param data - 要驗證的資料
   * @param schema - Zod schema
   * @param endpoint - API 端點名稱 (用於錯誤訊息)
   * @returns 驗證結果 (不拋錯，僅記錄警告)
   */
  function validateWithSchema<T>(
    data: unknown,
    schema: ZodSchema<T>,
    endpoint: string
  ): { success: true; data: T } | { success: false; issues: Array<{ path: string[]; message: string }> } {
    // 僅在開發環境執行驗證
    if (!isDev()) {
      return { success: true, data: data as T }
    }

    const result = schema.safeParse(data)

    if (result.success) {
      return { success: true, data: result.data }
    }

    // 開發環境：顯示顯眼的 console.warn
    const issues = result.error.issues.map((issue) => ({
      path: issue.path.map((p) => String(p)),
      message: issue.message,
      expected: issue.expected,
      received: issue.received,
    }))

    const formattedIssues = formatValidationIssues(issues)

    // eslint-disable-next-line no-console
    console.warn(
      `%c⚠️ Schema 驗證失敗 (開發環境僅警告) %c[${endpoint}]`,
      'font-weight: bold; color: #ffc107; font-size: 14px;',
      'color: #6c757d; font-size: 12px;'
    )
    // eslint-disable-next-line no-console
    console.warn(getValidationWarningStyle())
    // eslint-disable-next-line no-console
    console.warn(`驗證問題 (共 ${issues.length} 項):\n${formattedIssues}`)
    // eslint-disable-next-line no-console
    console.warn('─'.repeat(60))

    return { success: false, issues }
  }

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

    // Handle empty response
    if (!text || text.trim() === '') {
      // For 2xx status codes with empty body, return empty object or null
      if (response.ok) {
        return {} as T
      }

      // For error status codes with empty body
      const networkError: NetworkError = {
        status: response.status,
        statusText: response.statusText,
        message: `HTTP ${response.status}: ${response.statusText}`,
        originalError: undefined,
      }
      throw networkError
    }

    try {
      const data = JSON.parse(text) as ApiResponse<T>

      // Check for API-level error (use isSuccessCode to handle '0', 'SUCCESS', etc.)
      if (!isSuccessCode(data.code)) {
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
   * @param schema - Optional Zod schema for response validation (dev mode only)
   * @returns Promise resolving to the response data
   */
  async function get<T>(
    endpoint: string,
    params?: Record<string, any>,
    options?: { showAlert?: boolean; context?: Record<string, unknown> },
    schema?: ZodSchema<T>
  ): Promise<T> {
    const url = `${apiBase}${endpoint}${buildQueryString(params)}`
    const headers = buildStandardHeaders()

    try {
      const response = await fetch(url, { headers })
      await checkResponse(response)
      const data = await parseResponse<T>(response)

      // Schema 驗證 (開發環境專用)
      if (schema) {
        const validation = validateWithSchema(data, schema, `GET ${endpoint}`)
        // 驗證失敗時仍回傳原始資料，避免影響功能
        return validation.success ? validation.data : data
      }

      return data
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
   * @param schema - Optional Zod schema for response validation (dev mode only)
   * @returns Promise resolving to the response data
   */
  async function post<T>(
    endpoint: string,
    data: any,
    options?: { showAlert?: boolean; context?: Record<string, unknown> },
    schema?: ZodSchema<T>
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
      const result = await parseResponse<T>(response)

      // Schema 驗證 (開發環境專用)
      if (schema) {
        const validation = validateWithSchema(result, schema, `POST ${endpoint}`)
        return validation.success ? validation.data : result
      }

      return result
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
   * @param schema - Optional Zod schema for response validation (dev mode only)
   * @returns Promise resolving to the response data
   */
  async function put<T>(
    endpoint: string,
    data: any,
    options?: { showAlert?: boolean; context?: Record<string, unknown> },
    schema?: ZodSchema<T>
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
      const result = await parseResponse<T>(response)

      // Schema 驗證 (開發環境專用)
      if (schema) {
        const validation = validateWithSchema(result, schema, `PUT ${endpoint}`)
        return validation.success ? validation.data : result
      }

      return result
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
   * @param schema - Optional Zod schema for response validation (dev mode only)
   * @returns Promise resolving to the response data
   */
  async function patch<T>(
    endpoint: string,
    data: any,
    options?: { showAlert?: boolean; context?: Record<string, unknown> },
    schema?: ZodSchema<T>
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
      const result = await parseResponse<T>(response)

      // Schema 驗證 (開發環境專用)
      if (schema) {
        const validation = validateWithSchema(result, schema, `PATCH ${endpoint}`)
        return validation.success ? validation.data : result
      }

      return result
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
   * @param schema - Optional Zod schema for response validation (dev mode only)
   * @returns Promise resolving to the response data
   */
  async function del<T>(
    endpoint: string,
    options?: { showAlert?: boolean; context?: Record<string, unknown> },
    schema?: ZodSchema<T>
  ): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    try {
      const response = await fetch(url, {
        method: 'DELETE',
        headers,
      })
      await checkResponse(response)
      const result = await parseResponse<T>(response)

      // Schema 驗證 (開發環境專用)
      if (schema) {
        const validation = validateWithSchema(result, schema, `DELETE ${endpoint}`)
        return validation.success ? validation.data : result
      }

      return result
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
   * @param schema - Optional Zod schema for response validation (dev mode only)
   * @returns Promise resolving to the upload response
   */
  async function upload<T>(
    endpoint: string,
    file: File,
    fieldName: string = 'file',
    options?: { showAlert?: boolean; context?: Record<string, unknown> },
    schema?: ZodSchema<T>
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
      const result = await parseResponse<T>(response)

      // Schema 驗證 (開發環境專用)
      if (schema) {
        const validation = validateWithSchema(result, schema, `UPLOAD ${endpoint}`)
        return validation.success ? validation.data : result
      }

      return result
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
   * @param schema - Optional Zod schema for response validation (dev mode only)
   * @returns Promise resolving to the raw response
   */
  async function raw<T>(
    endpoint: string,
    options?: { showAlert?: boolean; context?: Record<string, unknown> },
    schema?: ZodSchema<T>
  ): Promise<T> {
    const url = `${apiBase}${endpoint}`
    const headers = buildStandardHeaders()

    try {
      const response = await fetch(url, { headers })
      await checkResponse(response)
      const result = await parseResponse<T>(response, true)

      // Schema 驗證 (開發環境專用)
      if (schema) {
        const validation = validateWithSchema(result, schema, `RAW ${endpoint}`)
        return validation.success ? validation.data : result
      }

      return result
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
