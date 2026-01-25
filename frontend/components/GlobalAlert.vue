<script setup lang="ts">
import { ref, provide, computed } from 'vue'

type AlertType = 'info' | 'warning' | 'error' | 'success'

interface AlertOptions {
  title?: string
  message: string
  type?: AlertType
  confirmText?: string
  cancelText?: string
  onConfirm?: () => void
  onCancel?: () => void
}

const visible = ref(false)
const alertConfig = ref<AlertOptions | null>(null)
const alertResolve = ref<((value: boolean) => void) | null>(null)

const defaultTitles: Record<AlertType, string> = {
  info: '提示',
  warning: '提醒',
  error: '操作失敗',
  success: '操作成功',
}

const bgColors: Record<AlertType, string> = {
  info: 'bg-blue-500',
  warning: 'bg-yellow-500',
  error: 'bg-red-500',
  success: 'bg-green-500',
}

const buttonColors: Record<AlertType, string> = {
  info: 'bg-blue-500 hover:bg-blue-600',
  warning: 'bg-yellow-500 hover:bg-yellow-600',
  error: 'bg-red-500 hover:bg-red-600',
  success: 'bg-green-500 hover:bg-green-600',
}

const showAlert = (options: AlertOptions): Promise<boolean> => {
  alertConfig.value = {
    title: options.title || defaultTitles[options.type || 'info'],
    message: options.message,
    type: options.type || 'info',
    confirmText: options.confirmText || '確定',
    cancelText: options.cancelText || '取消',
    onConfirm: options.onConfirm,
    onCancel: options.onCancel,
  }
  visible.value = true
  
  return new Promise((resolve) => {
    alertResolve.value = resolve
  })
}

const confirm = (message: string, title?: string): Promise<boolean> => {
  return showAlert({
    message,
    title: title || '確認操作',
    type: 'warning',
    confirmText: '確認',
    cancelText: '取消',
  })
}

const info = (message: string, title?: string) => showAlert({ message, title, type: 'info' })
const warning = (message: string, title?: string) => showAlert({ message, title, type: 'warning' })
const error = (message: string, title?: string) => showAlert({ message, title, type: 'error' })
const success = (message: string, title?: string) => showAlert({ message, title, type: 'success' })

const handleConfirm = () => {
  alertConfig.value?.onConfirm?.()
  if (alertResolve.value) {
    alertResolve.value(true)
  }
  visible.value = false
  alertConfig.value = null
  alertResolve.value = null
}

const handleCancel = () => {
  alertConfig.value?.onCancel?.()
  if (alertResolve.value) {
    alertResolve.value(false)
  }
  visible.value = false
  alertConfig.value = null
  alertResolve.value = null
}

// Provide to child components
provide('useAlert', { showAlert, info, warning, error, success, confirm })

// Window globals
if (typeof window !== 'undefined') {
  ;(window as any).$alert = showAlert
  ;(window as any).$confirm = confirm
  ;(window as any).$info = info
  ;(window as any).$warning = warning
  ;(window as any).$error = error
  ;(window as any).$success = success
}
</script>

<template>
  <!-- 使用 transition包裹 -->
  <Transition name="alert-fade">
    <div 
      v-if="visible" 
      class="fixed inset-0 z-[9999] flex items-center justify-center p-4 bg-black/50"
      @click="handleCancel"
    >
      <!-- Dialog -->
      <div 
        class="glass-card w-96 p-6 relative z-10"
        @click.stop
      >
        <!-- Icon -->
        <div 
          class="w-12 h-12 rounded-full flex items-center justify-center mb-4"
          :class="alertConfig ? bgColors[alertConfig.type!] : 'bg-blue-500'"
        >
          <svg v-if="alertConfig?.type === 'warning'" class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <svg v-else-if="alertConfig?.type === 'error'" class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
          <svg v-else-if="alertConfig?.type === 'success'" class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <svg v-else class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        
        <!-- Title -->
        <h3 class="text-lg font-bold text-white mb-2">
          {{ alertConfig?.title || '提示' }}
        </h3>
        
        <!-- Message -->
        <p class="text-slate-300 mb-6">
          {{ alertConfig?.message || '' }}
        </p>
        
        <!-- Buttons -->
        <div class="flex gap-3">
          <button 
            v-if="alertConfig?.type === 'warning' || alertConfig?.type === 'error'"
            @click="handleCancel"
            class="flex-1 py-2.5 px-4 rounded-lg font-medium bg-slate-700 text-white hover:bg-slate-600 transition-colors"
          >
            {{ alertConfig?.cancelText || '取消' }}
          </button>
          <button 
            @click="handleConfirm"
            class="flex-1 py-2.5 px-4 rounded-lg font-medium text-white transition-colors"
            :class="alertConfig ? buttonColors[alertConfig.type!] : 'bg-blue-500 hover:bg-blue-600'"
          >
            {{ alertConfig?.confirmText || '確定' }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.alert-fade-enter-active,
.alert-fade-leave-active {
  transition: opacity 0.2s ease;
}

.alert-fade-enter-from,
.alert-fade-leave-to {
  opacity: 0;
}

.alert-fade-enter-active > div,
.alert-fade-leave-active > div {
  transition: transform 0.3s ease;
}

.alert-fade-enter-from > div,
.alert-fade-leave-to > div {
  transform: scale(0.95);
}
</style>
