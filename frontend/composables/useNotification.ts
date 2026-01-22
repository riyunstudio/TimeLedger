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
  }
}
