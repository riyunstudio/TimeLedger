<script setup lang="ts">
interface Props {
  loading: boolean
  size?: 'sm' | 'md' | 'lg'
  text?: string
  fullScreen?: boolean
  transparent?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  fullScreen: false,
  transparent: false
})

const sizeClasses = {
  sm: 'w-4 h-4',
  md: 'w-8 h-8',
  lg: 'w-12 h-12'
}

const containerClass = computed(() => {
  if (props.fullScreen) {
    return 'fixed inset-0 flex flex-col items-center justify-center z-50'
  }
  return 'flex flex-col items-center justify-center p-4'
})

const overlayClass = computed(() => {
  if (props.transparent) {
    return 'bg-transparent'
  }
  return 'bg-white/80 dark:bg-gray-900/80'
})
</script>

<template>
  <div
    v-if="loading"
    :class="[containerClass, overlayClass]"
  >
    <!-- Loading Spinner -->
    <div
      class="relative"
      :class="sizeClasses[size]"
    >
      <div
        class="absolute inset-0 rounded-full border-2 border-gray-200 dark:border-gray-700"
      />
      <div
        class="absolute inset-0 rounded-full border-2 border-indigo-600 border-t-transparent animate-spin"
      />
    </div>

    <!-- Loading Text -->
    <p
      v-if="text"
      class="mt-3 text-sm text-gray-600 dark:text-gray-400 font-medium"
    >
      {{ text }}
    </p>

    <!-- Slot for custom content -->
    <slot />
  </div>

  <!-- Default slot when not loading -->
  <slot v-else />
</template>
