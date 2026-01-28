// 全域通知狀態
const notificationState = {
  show: ref(false),
  // Toast 通知狀態
  toast: ref<{
    type: 'success' | 'error' | 'warning' | 'info'
    message: string
    visible: boolean
  } | null>(null),
}

let toastTimeout: ReturnType<typeof setTimeout> | null = null

// 顯示 Toast 並自動消失
const showToast = (type: 'success' | 'error' | 'warning' | 'info', message: string, duration: number = 3000) => {
  // 清除之前的計時器
  if (toastTimeout) {
    clearTimeout(toastTimeout)
    toastTimeout = null
  }

  // 設定新的 toast
  notificationState.toast.value = {
    type,
    message,
    visible: true,
  }

  // 自動隱藏
  toastTimeout = setTimeout(() => {
    notificationState.toast.value = null
    toastTimeout = null
  }, duration)
}

export const useNotification = () => {
  return {
    show: notificationState.show,
    toast: notificationState.toast,
    toggle: () => {
      notificationState.show.value = !notificationState.show.value
    },
    open: () => {
      notificationState.show.value = true
    },
    close: () => {
      notificationState.show.value = false
    },
    success: (message: string) => {
      showToast('success', message)
    },
    error: (message: string) => {
      showToast('error', message)
    },
    warning: (message: string) => {
      showToast('warning', message)
    },
    info: (message: string) => {
      showToast('info', message)
    },
    showSuccess: (message: string) => {
      showToast('success', message)
    },
    showError: (message: string) => {
      showToast('error', message)
    },
  }
}

// Toast 組件需要的 CSS 樣式
export const toastStyles = `
.toast-container {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.toast {
  padding: 12px 20px;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  animation: toastSlideUp 0.3s ease-out;
  display: flex;
  align-items: center;
  gap: 8px;
}

.toast.success {
  background: linear-gradient(135deg, #10b981, #059669);
}

.toast.error {
  background: linear-gradient(135deg, #ef4444, #dc2626);
}

.toast.warning {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.toast.info {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
}

@keyframes toastSlideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
`
