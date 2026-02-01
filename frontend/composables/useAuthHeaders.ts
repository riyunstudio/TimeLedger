/**
 * useAuthHeaders Composable
 * Centralizes Authorization header construction for API requests.
 *
 * This composable provides a unified way to build authentication headers
 * for all API requests, handling both admin and teacher tokens.
 */

export function useAuthHeaders() {
  const config = useRuntimeConfig()

  /**
   * Retrieves the authorization header for API requests.
   * Supports multiple token storage keys for different user types.
   *
   * @returns Record containing Authorization header if token exists, empty object otherwise
   */
  function getAuthHeader(): Record<string, string> {
    if (process.server) {
      // Server-side: Check authorization header from incoming request
      return {}
    }

    // Client-side: Check multiple possible token storage locations
    const token = localStorage.getItem('admin_token') ||
      localStorage.getItem('teacher_token') ||
      localStorage.getItem('token')

    if (!token) {
      return {}
    }

    return {
      Authorization: `Bearer ${token}`,
    }
  }

  /**
   * Gets the content type header for JSON requests.
   */
  function getContentTypeHeader(): Record<string, string> {
    return {
      'Content-Type': 'application/json',
    }
  }

  /**
   * Builds combined headers for standard API requests.
   * Includes both authorization, content type, and cache-control headers.
   */
  function buildStandardHeaders(): Record<string, string> {
    return {
      ...getAuthHeader(),
      ...getContentTypeHeader(),
      // 防止瀏覽器快取 GET 請求
      'Cache-Control': 'no-cache, no-store, must-revalidate',
      'Pragma': 'no-cache',
      'Expires': '0',
    }
  }

  /**
   * Checks if a valid token exists (any user type).
   */
  function hasValidToken(): boolean {
    if (process.server) {
      return false
    }
    return !!(
      localStorage.getItem('admin_token') ||
      localStorage.getItem('teacher_token') ||
      localStorage.getItem('token')
    )
  }

  /**
   * Gets the current user type based on stored tokens.
   * Returns 'admin', 'teacher', or null if no token found.
   */
  function getCurrentUserType(): 'admin' | 'teacher' | null {
    if (process.server) {
      return null
    }

    if (localStorage.getItem('admin_token')) {
      return 'admin'
    }
    if (localStorage.getItem('teacher_token')) {
      return 'teacher'
    }
    if (localStorage.getItem('token')) {
      // Legacy token format, treat as teacher
      return 'teacher'
    }
    return null
  }

  return {
    getAuthHeader,
    getContentTypeHeader,
    buildStandardHeaders,
    hasValidToken,
    getCurrentUserType,
  }
}
