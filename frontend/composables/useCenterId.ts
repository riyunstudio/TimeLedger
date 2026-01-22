export const useCenterId = () => {
  const authStore = useAuthStore()

  const getCenterId = (): string => {
    if (process.client) {
      const userData = authStore.user as any
      if (userData?.center_id) {
        return userData.center_id.toString()
      }
      return localStorage.getItem('center_id') || '1'
    }
    return '1'
  }

  return { getCenterId }
}
