const toastComponent = ref<any>(null)

export function useToast() {
  const showToast = (type: 'success' | 'error' | 'warning' | 'info', message: string, title?: string) => {
    if (toastComponent.value) {
      toastComponent.value[type](message, title)
    } else {
      // Fallback to native alert if toast component not available
      alert(message)
    }
  }

  return {
    toast: toastComponent,
    success: (message: string, title?: string) => showToast('success', message, title),
    error: (message: string, title?: string) => showToast('error', message, title),
    warning: (message: string, title?: string) => showToast('warning', message, title),
    info: (message: string, title?: string) => showToast('info', message, title),
  }
}

// 用於在全局註冊 toast component
export function registerToast(component: any) {
  toastComponent.value = component
}
