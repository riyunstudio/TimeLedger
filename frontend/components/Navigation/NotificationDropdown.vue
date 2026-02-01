<template>
  <div class="fixed inset-0 z-[150] pointer-events-none">
    <div class="absolute inset-0 bg-black/40 backdrop-blur-sm pointer-events-auto" @click="notificationUI.close()"></div>
    <div class="absolute top-14 right-4 w-80 glass-card shadow-2xl pointer-events-auto animate-fade-in max-h-[70vh] overflow-hidden flex flex-col">
      <div class="p-4 border-b border-white/10 shrink-0">
        <div class="flex items-center justify-between">
          <h3 class="font-semibold text-slate-100">通知</h3>
          <button
            @click="markAllAsRead"
            :disabled="notificationDataStore.unreadCount === 0 || notificationDataStore.isMarkingAllRead"
            class="text-sm transition-colors"
            :class="notificationDataStore.unreadCount === 0 ? 'text-slate-600 cursor-not-allowed' : 'text-primary-500 hover:text-primary-400'"
          >
            {{ notificationDataStore.isMarkingAllRead ? '處理中...' : '全部標記為已讀' }}
          </button>
        </div>
      </div>

      <div class="flex-1 overflow-y-auto">
        <div
          v-if="notificationDataStore.notifications.length === 0"
          class="p-8 text-center text-slate-500"
        >
          無通知
        </div>

        <div
          v-else
          class="divide-y divide-white/10"
        >
          <div
            v-for="notification in notificationDataStore.notifications"
            :key="notification.id"
            class="p-4 hover:bg-white/5 transition-colors cursor-pointer"
            :class="{ 'bg-primary-500/10': !notification.is_read }"
            @click="handleNotificationClick(notification)"
          >
            <div class="flex items-start gap-3">
              <div class="flex-1">
                <h4 class="font-medium text-slate-100 text-sm mb-1">
                  {{ notification.title }}
                </h4>
                <p class="text-sm text-slate-400">
                  {{ notification.message }}
                </p>
                <p class="text-xs text-slate-500 mt-2">
                  {{ formatTime(notification.created_at) }}
                </p>
              </div>

              <div
                v-if="!notification.is_read"
                class="w-2 h-2 rounded-full bg-primary-500 shrink-0 mt-1"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const notificationDataStore = useNotificationStore()
const notificationUI = useNotification()
const router = useRouter()

// 除錯：監聽 show 狀態變化
console.log('NotificationDropdown mounted, show:', notificationUI.show.value)

// 監聽通知彈窗顯示狀態，開啟時重新整理
watch(
  () => notificationUI.show.value,
  (isShown) => {
    console.log('Show state changed:', isShown)
    if (isShown) {
      notificationDataStore.fetchNotifications(true)
    }
  },
  { immediate: true }
)

const handleNotificationClick = async (notification: any) => {
  if (!notification.is_read) {
    await notificationDataStore.markNotificationRead(notification.id)
  }

  // 根據通知類型跳轉到相應頁面
  const userType = notification.user_type || notificationUI.getUserType?.() || 'TEACHER'
  
  if (notification.title === '新例外申請通知' || notification.type === 'EXCEPTION') {
    // 管理員跳轉到審核頁面
    if (userType === 'ADMIN') {
      router.push('/admin/approval')
    }
  } else if (notification.type === 'REVIEW_RESULT' || notification.title === '例外單審核結果') {
    // 老師跳轉到自己的例外申請頁面
    if (userType === 'TEACHER') {
      router.push('/teacher/exceptions')
    }
  } else if (notification.type === 'APPROVAL') {
    // 管理員跳轉到審核頁面查看結果
    if (userType === 'ADMIN') {
      router.push('/admin/approval')
    }
  }

  notificationUI.close()
}

const markAllAsRead = async () => {
  // 檢查是否有未讀通知
  if (notificationDataStore.unreadCount === 0) {
    return
  }

  // 檢查是否正在處理中
  if (notificationDataStore.isMarkingAllRead) {
    return
  }

  try {
    await notificationDataStore.markAllAsRead()
  } catch (error) {
    console.error('Failed to mark all as read:', error)
  }
}

const formatTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / (1000 * 60))
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffMins < 1) return '剛剛'
  if (diffMins < 60) return `${diffMins} 分鐘前`
  if (diffHours < 24) return `${diffHours} 小時前`
  if (diffDays < 7) return `${diffDays} 天前`

  return date.toLocaleDateString('zh-TW', {
    month: 'short',
    day: 'numeric',
  })
}
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.2s ease-out;
}
</style>
