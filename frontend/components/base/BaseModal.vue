<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div
          class="absolute inset-0 bg-black/60 backdrop-blur-sm transition-opacity"
          :class="{ 'cursor-pointer': closeOnBackdrop }"
          @click="closeOnBackdrop ? handleClose() : null"
        />
        <div
          class="relative glass-card w-full overflow-y-auto"
          :class="[
            sizeClasses,
            mobilePositionClasses
          ]"
          @click.stop
        >
          <!-- Header -->
          <div v-if="title || $slots.header" class="flex items-center justify-between mb-4">
            <slot name="header">
              <h2 v-if="title" class="text-xl font-bold text-white">{{ title }}</h2>
            </slot>
            <button
              v-if="showCloseButton"
              class="text-slate-400 hover:text-white transition-colors p-1 rounded-lg hover:bg-white/10"
              @click="handleClose"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <!-- Body -->
          <slot />
          <!-- Footer -->
          <div v-if="$slots.footer" class="mt-6 pt-4 border-t border-white/10">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
interface Props {
  modelValue: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  closeOnBackdrop?: boolean
  mobilePosition?: 'center' | 'bottom'
  showCloseButton?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  closeOnBackdrop: true,
  mobilePosition: 'bottom',
  showCloseButton: true,
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'close': []
}>()

const handleClose = () => {
  emit('update:modelValue', false)
  emit('close')
}

const sizeClasses = computed(() => {
  const sizes = {
    sm: 'w-full max-w-sm p-4',
    md: 'w-full max-w-md p-6',
    lg: 'w-full max-w-lg p-6',
    xl: 'w-full max-w-2xl p-6',
  }
  return sizes[props.size]
})

const mobilePositionClasses = computed(() => {
  if (typeof window === 'undefined') return ''
  const isMobile = window.innerWidth < 768

  if (!isMobile) return ''

  return props.mobilePosition === 'bottom' ? 'md:rounded-2xl rounded-t-3xl' : 'rounded-2xl'
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .glass-card,
.modal-leave-to .glass-card {
  transform: translateY(20px);
}
</style>
