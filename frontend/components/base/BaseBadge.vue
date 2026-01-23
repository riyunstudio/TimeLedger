<template>
  <span
    :class="variantClasses"
    class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium"
  >
    <slot />
  </span>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'success' | 'critical' | 'warning' | 'info'
  size?: 'sm' | 'md'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
})

const variantClasses = computed(() => {
  const variants = {
    primary: 'bg-primary-500/20 text-primary-300 border border-primary-500/30',
    secondary: 'bg-secondary-500/20 text-secondary-300 border border-secondary-500/30',
    success: 'bg-success-500/20 text-success-300 border border-success-500/30',
    critical: 'bg-critical-500/20 text-critical-300 border border-critical-500/30',
    warning: 'bg-warning-500/20 text-warning-300 border border-warning-500/30',
    info: 'bg-blue-500/20 text-blue-300 border border-blue-500/30',
  }
  return variants[props.variant]
})
</script>
