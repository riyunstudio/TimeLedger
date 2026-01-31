<template>
  <Teleport to="body">
    <TransitionGroup
      name="toast"
      tag="div"
      class="fixed bottom-4 right-4 z-[9999] space-y-3"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="glass-card px-4 py-3 rounded-xl shadow-xl flex items-center gap-3 min-w-[280px] max-w-md"
        :class="getToastClass(toast.type)"
      >
        <div class="shrink-0">
          <svg v-if="toast.type === 'success'" class="w-5 h-5 text-success-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <svg v-else-if="toast.type === 'error'" class="w-5 h-5 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <svg v-else-if="toast.type === 'warning'" class="w-5 h-5 text-warning-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <svg v-else class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="flex-1">
          <p v-if="toast.title" class="font-medium text-slate-100 text-sm">{{ toast.title }}</p>
          <p class="text-slate-300 text-sm">{{ toast.message }}</p>
        </div>
        <button
          @click="removeToast(toast.id)"
          class="shrink-0 p-1 rounded-lg hover:bg-white/10 transition-colors"
        >
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </TransitionGroup>
  </Teleport>
</template>

<script setup lang="ts">
interface Toast {
  id: number
  type: 'success' | 'error' | 'warning' | 'info'
  title?: string
  message: string
}

const toasts = ref<Toast[]>([])
let toastId = 0

const getToastClass = (type: string) => {
  switch (type) {
    case 'success':
      return 'border-l-4 border-l-success-500'
    case 'error':
      return 'border-l-4 border-l-critical-500'
    case 'warning':
      return 'border-l-4 border-l-warning-500'
    default:
      return 'border-l-4 border-l-primary-500'
  }
}

const removeToast = (id: number) => {
  toasts.value = toasts.value.filter(t => t.id !== id)
}

const addToast = (type: Toast['type'], message: string, title?: string) => {
  const id = ++toastId
  toasts.value.push({ id, type, message, title })
  setTimeout(() => removeToast(id), 5000)
}

// 暴露方法給外部使用
defineExpose({
  success: (message: string, title?: string) => addToast('success', message, title),
  error: (message: string, title?: string) => addToast('error', message, title),
  warning: (message: string, title?: string) => addToast('warning', message, title),
  info: (message: string, title?: string) => addToast('info', message, title),
})
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100px);
}

.toast-move {
  transition: transform 0.3s ease;
}
</style>
