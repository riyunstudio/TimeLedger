<template>
  <div class="animate-pulse">
    <div class="space-y-4">
      <div
        v-for="i in rows"
        :key="i"
        class="flex items-center space-x-4"
      >
        <div :class="avatarClass" />
        <div class="flex-1 space-y-2">
          <div :style="{ width: titleWidth }" class="h-4 rounded bg-slate-700/50"></div>
          <div :style="{ width: subtitleWidth }" class="h-3 rounded bg-slate-700/30"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  rows?: number
  avatarClass?: string
  titleWidth?: string
  subtitleWidth?: string
}

const props = withDefaults(defineProps<Props>(), {
  rows: 5,
  avatarClass: 'h-12 w-12 flex-shrink-0',
  titleWidth: '40%',
  subtitleWidth: '70%',
})
</script>
