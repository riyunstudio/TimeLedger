<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeOnBackdrop ? $emit('update:modelValue', false) : null" />
        <div
          class="relative glass-card max-h-[90vh] overflow-y-auto"
          :class="[
            sizeClasses,
            mobilePositionClasses
          ]"
        >
          <div class="flex items-center justify-between mb-4">
            <h2 v-if="title" class="text-xl font-bold text-white">{{ title }}</h2>
            <button
              class="text-slate-400 hover:text-white transition-colors"
              @click="$emit('update:modelValue', false)"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12M6 6v12" />
              </svg>
            </button>
          </div>
          <slot />
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
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  closeOnBackdrop: true,
  mobilePosition: 'bottom',
})

defineEmits<{
  'update:modelValue': [value: boolean]
}>()

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
