// Global Alert Composable
import { inject } from 'vue'

export type AlertType = 'info' | 'warning' | 'error' | 'success'

export interface AlertButton {
  text: string
  style?: 'primary' | 'warning' | 'critical' | 'secondary'
  action?: () => void
}

export interface AlertOptions {
  title?: string
  message: string
  type?: AlertType
  confirmText?: string
  cancelText?: string
  buttons?: AlertButton[]
  onConfirm?: () => void
  onCancel?: () => void
  autoCloseDelay?: number
  width?: 'sm' | 'md' | 'lg' | 'xl'
}

interface AlertComposable {
  showAlert: (options: AlertOptions) => Promise<boolean>
  info: (message: string, title?: string) => Promise<boolean>
  warning: (message: string, title?: string) => Promise<boolean>
  error: (message: string, title?: string) => Promise<boolean>
  success: (message: string, title?: string) => Promise<boolean>
  confirm: (message: string, title?: string) => Promise<boolean>
  dialog: (message: string, buttons: AlertButton[], title?: string, type?: AlertType) => Promise<string>
}

// Vue composable
export const useAlert = () => {
  return inject<AlertComposable>('useAlert', {
    showAlert: async () => {
      console.warn('[useAlert] Fallback showAlert called')
      const result = await (window as any).$confirm?.('提示')
      return result
    },
    info: async () => false,
    warning: async () => false,
    error: async () => false,
    success: async () => false,
    confirm: async (message: string) => {
      console.warn('[useAlert] Fallback confirm called:', message)
      const result = await (window as any).$confirm?.(message)
      return result === true
    },
    dialog: async () => '',
  })
}

// ============ Standalone Functions ============

const getGlobalAlert = () => {
  if (typeof window !== 'undefined') {
    return (window as any).$alert
  }
  return null
}

const getGlobalConfirm = () => {
  if (typeof window !== 'undefined') {
    return (window as any).$confirm
  }
  return null
}

export const showAlert = (options: AlertOptions): Promise<boolean> => {
  const alertFn = getGlobalAlert()
  if (alertFn) {
    return alertFn(options)
  }
  return Promise.resolve(false)
}

export const alertInfo = (message: string, title?: string): Promise<boolean> => {
  const alertFn = getGlobalAlert()
  if (alertFn) {
    return alertFn({ message, title, type: 'info' })
  }
  return Promise.resolve(false)
}

export const alertWarning = (message: string, title?: string): Promise<boolean> => {
  const alertFn = getGlobalAlert()
  if (alertFn) {
    return alertFn({ message, title, type: 'warning' })
  }
  return Promise.resolve(false)
}

export const alertError = (message: string, title?: string): Promise<boolean> => {
  const alertFn = getGlobalAlert()
  if (alertFn) {
    return alertFn({ message, title, type: 'error' })
  }
  return Promise.resolve(false)
}

export const alertSuccess = (message: string, title?: string): Promise<boolean> => {
  const alertFn = getGlobalAlert()
  if (alertFn) {
    return alertFn({ message, title, type: 'success' })
  }
  return Promise.resolve(false)
}

export const alertConfirm = (message: string, title?: string): Promise<boolean> => {
  const confirmFn = getGlobalConfirm()
  if (confirmFn) {
    return confirmFn(message, title)
  }
  return Promise.resolve(false)
}

export const alertDialog = (
  message: string,
  buttons: AlertButton[],
  title?: string,
  type?: AlertType
): Promise<string> => {
  const alertFn = getGlobalAlert()
  if (alertFn) {
    return alertFn({
      message,
      title: title || (type ? { info: '提示', warning: '提醒', error: '操作失敗', success: '操作成功' }[type] : '提示'),
      type,
      buttons,
      width: buttons.length > 2 ? 'lg' : 'md',
    }) as Promise<string>
  }
  return Promise.resolve('')
}

export const nativeAlert = (message: string, title?: string) => {
  alert(message + (title ? `\n${title}` : ''))
}

export const $alert = showAlert
export const $confirm = alertConfirm
export const $info = alertInfo
export const $warning = alertWarning
export const $error = alertError
export const $success = alertSuccess
export const $dialog = alertDialog
