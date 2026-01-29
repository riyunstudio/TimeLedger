/**
 * useTokenManager Composable
 * Centralizes token management operations for authentication.
 *
 * This composable handles:
 * - Token storage and retrieval
 * - Token clearing on logout
 * - User type detection
 * - Redirect logic after logout
 */

export function useTokenManager() {
  const router = useRouter()

  // Token storage keys
  const TOKEN_KEYS = {
    admin: 'admin_token',
    teacher: 'teacher_token',
    legacy: 'token',
    userType: 'current_user_type',
  } as const

  /**
   * Stores a token with its associated user type.
   *
   * @param token - The JWT token to store
   * @param userType - The type of user ('admin' or 'teacher')
   */
  function setToken(token: string, userType: 'admin' | 'teacher'): void {
    if (process.server) {
      return
    }

    const storageKey = userType === 'admin' ? TOKEN_KEYS.admin : TOKEN_KEYS.teacher
    localStorage.setItem(storageKey, token)
    localStorage.setItem(TOKEN_KEYS.userType, userType)
  }

  /**
   * Clears all authentication tokens from storage.
   * This should be called on logout or token expiration.
   */
  function clearAllTokens(): void {
    if (process.server) {
      return
    }

    localStorage.removeItem(TOKEN_KEYS.admin)
    localStorage.removeItem(TOKEN_KEYS.teacher)
    localStorage.removeItem(TOKEN_KEYS.legacy)
    localStorage.removeItem(TOKEN_KEYS.userType)
  }

  /**
   * Determines if the user was an admin before logout.
   * Used to decide which login page to redirect to.
   */
  function wasAdminUser(): boolean {
    if (process.server) {
      return false
    }

    const userType = localStorage.getItem(TOKEN_KEYS.userType)
    return userType === 'admin' || !!localStorage.getItem(TOKEN_KEYS.admin)
  }

  /**
   * Determines if the user was a teacher before logout.
   */
  function wasTeacherUser(): boolean {
    if (process.server) {
      return false
    }

    const userType = localStorage.getItem(TOKEN_KEYS.userType)
    return userType === 'teacher' || !!localStorage.getItem(TOKEN_KEYS.teacher)
  }

  /**
   * Handles unauthorized response by clearing tokens and redirecting.
   * Automatically determines the correct login page based on user type.
   *
   * @param currentPath - The current path to redirect back to after login
   */
  function handleUnauthorized(currentPath?: string): void {
    if (process.server) {
      return
    }

    // Determine which login page to redirect to
    let redirectPath = '/'

    if (currentPath) {
      const encodedRedirect = encodeURIComponent(currentPath)

      if (wasTeacherUser() && !currentPath.includes('/teacher/login')) {
        redirectPath = `/teacher/login?redirect=${encodedRedirect}`
      } else if (wasAdminUser() && !currentPath.includes('/admin/login')) {
        redirectPath = `/admin/login?redirect=${encodedRedirect}`
      } else if (!currentPath.includes('/login')) {
        // Default to home if unknown
        redirectPath = '/'
      }
    }

    // Clear all tokens first
    clearAllTokens()

    // Navigate to the appropriate page
    router.push(redirectPath)
  }

  /**
   * Gets the stored token for a specific user type.
   *
   * @param userType - The user type to get token for
   * @returns The token string or null if not found
   */
  function getToken(userType: 'admin' | 'teacher'): string | null {
    if (process.server) {
      return null
    }

    const key = userType === 'admin' ? TOKEN_KEYS.admin : TOKEN_KEYS.teacher
    return localStorage.getItem(key)
  }

  /**
   * Checks if a token exists for any user type.
   */
  function hasAnyToken(): boolean {
    if (process.server) {
      return false
    }

    return !!(
      localStorage.getItem(TOKEN_KEYS.admin) ||
      localStorage.getItem(TOKEN_KEYS.teacher) ||
      localStorage.getItem(TOKEN_KEYS.legacy)
    )
  }

  return {
    TOKEN_KEYS,
    setToken,
    clearAllTokens,
    wasAdminUser,
    wasTeacherUser,
    handleUnauthorized,
    getToken,
    hasAnyToken,
  }
}
