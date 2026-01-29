import { defineStore } from 'pinia'
import type { Notification } from '~/types'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])
  const unreadCount = ref<number>(0)
  const lastFetchTime = ref<number>(0)

  const fetchNotifications = async (forceRefresh = false) => {
    // 如果最近 30 秒內已經 fetch 過，且不是強制刷新，則跳過
    const now = Date.now()
    if (!forceRefresh && now - lastFetchTime.value < 30000 && notifications.value.length > 0) {
      return
    }

    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: { notifications: Notification[]; unread_count: number } }>(
        '/notifications?limit=50'
      )
      notifications.value = response.datas?.notifications || []
      unreadCount.value = response.datas?.unread_count || 0
      lastFetchTime.value = now
    } catch (error) {
      console.error('Failed to fetch notifications:', error)
    }
  }

  const markAsRead = async (id: number) => {
    try {
      const api = useApi()
      await api.post(`/notifications/${id}/read`, {})

      const notification = notifications.value.find((n) => n.id === id)
      if (notification) {
        notification.is_read = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch (error) {
      console.error('Failed to mark as read:', error)
    }
  }

  const markAllAsRead = async () => {
    try {
      const api = useApi()
      await api.post('/notifications/read-all', {})

      notifications.value.forEach((n) => (n.is_read = true))
      unreadCount.value = 0
    } catch (error) {
      console.error('Failed to mark all as read:', error)
    }
  }

  const addNotification = (notification: Notification) => {
    notifications.value.unshift(notification)
    unreadCount.value++
  }

  return {
    notifications,
    unreadCount,
    fetchNotifications,
    markAsRead,
    markAllAsRead,
    addNotification,
  }
})
