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

  // Loading 狀態
  const isFetching = ref(false)
  const isMarkingNotificationRead = ref(false)

  // 通知相關方法
  const fetchNotifications = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const response = await api.get<{ code: number; message: string; datas: Notification[] }>('/notifications')
        notifications.value = response.datas || []
      } catch (error) {
        console.error('Failed to fetch notifications:', error)
        throw error
      }
    })
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
      } catch (error) {
        console.error('Failed to mark notification as read:', error)
        throw error
      }
    })
  }

  return {
    // 資料狀態
    notifications,

    // Loading 狀態
    isFetching,
    isMarkingNotificationRead,

    // 方法
    fetchNotifications,
    markNotificationRead,
  }
})
