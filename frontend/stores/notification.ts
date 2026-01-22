import { defineStore } from 'pinia'
import type { Notification } from '~/types'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])
  const unreadCount = ref<number>(0)
  const isMock = ref(false)

  const loadMockNotifications = () => {
    notifications.value = [
      {
        id: 1,
        user_id: 1,
        user_type: 'TEACHER' as const,
        title: '課程已核准',
        message: '您的「鋼琴基礎」課程已通過審核',
        is_read: false,
        created_at: new Date().toISOString(),
        type: 'APPROVAL' as const,
      },
      {
        id: 2,
        user_id: 1,
        user_type: 'TEACHER' as const,
        title: '新課表提醒',
        message: '本週課表已更新，請查看',
        is_read: false,
        created_at: new Date(Date.now() - 3600000).toISOString(),
        type: 'SCHEDULE' as const,
      },
      {
        id: 3,
        user_id: 1,
        user_type: 'TEACHER' as const,
        title: '中心邀請',
        message: '「藝術中心」邀請您加入',
        is_read: true,
        created_at: new Date(Date.now() - 86400000).toISOString(),
        type: 'CENTER_INVITE' as const,
      },
    ]
    unreadCount.value = 2
    isMock.value = true
  }

  const fetchNotifications = async () => {
    if (isMock.value) return

    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: { notifications: Notification[]; unread_count: number } }>(
        '/notifications?limit=50'
      )
      notifications.value = response.datas?.notifications || []
      unreadCount.value = response.datas?.unread_count || 0
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
    isMock,
    loadMockNotifications,
    fetchNotifications,
    markAsRead,
    markAllAsRead,
    addNotification,
  }
})
