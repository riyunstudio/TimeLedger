<template>
  <div class="glass-card">
    <slot />
  </div>
</template>

<script setup lang="ts">
interface Props {
  padding?: string
  margin?: string
}

const props = withDefaults(defineProps<Props>(), {
  padding: 'p-6',
  margin: 'm-0',
})
</script>
