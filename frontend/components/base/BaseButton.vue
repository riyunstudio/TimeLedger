<template>
  <button
    :class="[
      variantClasses,
      sizeClasses,
      disabled ? 'opacity-50 cursor-not-allowed' : 'hover:scale-105 active:scale-95',
      'transition-all duration-300 shadow-lg hover:shadow-xl rounded-xl font-semibold animate-spring',
    ]"
    :disabled="disabled"
    @click="$emit('click', $event)"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'success' | 'critical' | 'warning'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false,
})

defineEmits<{
  click: [event: Event]
}>()

const variantClasses = computed(() => {
  const variants = {
    primary: 'bg-gradient-to-r from-primary-500 to-primary-600 text-white',
    secondary: 'bg-gradient-to-r from-secondary-500 to-secondary-600 text-white',
    success: 'bg-gradient-to-r from-success-500 to-success-600 text-white',
    critical: 'bg-gradient-to-r from-critical-500 to-critical-600 text-white',
    warning: 'bg-gradient-to-r from-warning-500 to-warning-600 text-white',
  }
  return variants[props.variant]
})

const sizeClasses = computed(() => {
  const sizes = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-6 py-3 text-base',
    lg: 'px-8 py-4 text-lg',
  }
  return sizes[props.size]
})
</script>
