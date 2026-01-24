const notificationState = {
  show: ref(false),
}

export const useNotification = () => {
  return {
    show: notificationState.show,
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
      console.log('Success:', message)
      // TODO: 實作成功通知 UI
    },
    error: (message: string) => {
      console.error('Error:', message)
      // TODO: 實作錯誤通知 UI
    },
    showSuccess: (message: string) => {
      console.log('Success:', message)
      // TODO: 實作成功通知 UI
    },
  }
}
