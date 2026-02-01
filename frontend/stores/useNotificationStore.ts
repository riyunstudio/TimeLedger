import { defineStore } from 'pinia'
import type { Notification } from '~/types'
import { withLoading } from '~/utils/loadingHelper'

// 台灣時區 datetime 格式化（用於 API 傳送）
const formatDateTimeForApi = (date: Date): string => {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}+08:00`
}

export const useNotificationStore = defineStore('notification', () => {
  // 資料狀態
  const notifications = ref<Notification[]>([])
  const unreadCount = ref(0)

  // Loading 狀態
  const isFetching = ref(false)
  const isMarkingNotificationRead = ref(false)
  const isMarkingAllRead = ref(false)

  // 計算屬性
  const hasUnread = computed(() => unreadCount.value > 0)

  // 通知相關方法
  const fetchNotifications = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const response = await api.get<{ code: number; message: string; datas: Notification[] }>('/notifications')
        notifications.value = response.datas || []
        // 取得通知後，同時更新未讀數量
        await fetchUnreadCount()
      } catch (error) {
        console.error('Failed to fetch notifications:', error)
        throw error
      }
    })
  }

  // 取得未讀通知數量
  const fetchUnreadCount = async () => {
    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; data: { count: number } }>('/notifications/unread-count')
      unreadCount.value = response.data?.count || 0
    } catch (error) {
      console.error('Failed to fetch unread count:', error)
      // 失敗時從本地 notifications 計算
      unreadCount.value = notifications.value.filter(n => !n.is_read).length
    }
  }

  const markNotificationRead = async (notificationId: number) => {
    return withLoading(isMarkingNotificationRead, async () => {
      try {
        const api = useApi()
        await api.post(`/notifications/${notificationId}/read`, {})
        const notification = notifications.value.find(n => n.id === notificationId)
        if (notification) {
          notification.is_read = true
          notification.read_at = formatDateTimeForApi(new Date())
        }
        // 更新未讀數量
        await fetchUnreadCount()
      } catch (error) {
        console.error('Failed to mark notification as read:', error)
        throw error
      }
    })
  }

  // 標記全部通知為已讀
  const markAllAsRead = async () => {
    return withLoading(isMarkingAllRead, async () => {
      try {
        const api = useApi()
        await api.post('/notifications/read-all', {})
        // 將所有通知標記為已讀
        notifications.value = notifications.value.map(n => ({
          ...n,
          is_read: true,
          read_at: n.read_at || formatDateTimeForApi(new Date())
        }))
        // 重置未讀數量
        unreadCount.value = 0
      } catch (error) {
        console.error('Failed to mark all notifications as read:', error)
        throw error
      }
    })
  }

  return {
    // 資料狀態
    notifications,
    unreadCount,

    // 計算屬性
    hasUnread,

    // Loading 狀態
    isFetching,
    isMarkingNotificationRead,
    isMarkingAllRead,

    // 方法
    fetchNotifications,
    fetchUnreadCount,
    markNotificationRead,
    markAllAsRead,
  }
}, {
  persist: {
    key: 'timeledger-notification',
    paths: ['notifications'],
  },
})
