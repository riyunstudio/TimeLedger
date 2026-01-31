<script setup lang="ts">
interface Props {
  type: 'text' | 'avatar' | 'card' | 'table' | 'title'
  lines?: number
  avatarSize?: 'sm' | 'md' | 'lg'
  avatarShape?: 'circle' | 'square'
  cardHeight?: string
  tableColumns?: number
  tableRows?: number
}

const props = withDefaults(defineProps<Props>(), {
  lines: 3,
  avatarSize: 'md',
  avatarShape: 'circle',
  cardHeight: '120px',
  tableColumns: 4,
  tableRows: 5
})

const avatarSizeClass = {
  sm: 'w-8 h-8',
  md: 'w-12 h-12',
  lg: 'w-16 h-16'
}

const baseClass = 'animate-pulse bg-gray-200 dark:bg-gray-700 rounded'
</script>

<template>
  <!-- Text Skeleton -->
  <div
    v-if="type === 'text'"
    class="space-y-2"
  >
    <div
      v-for="i in lines"
      :key="i"
      :class="[baseClass, i === lines ? 'w-3/4' : 'w-full']"
      :style="{ height: '1rem' }"
    />
  </div>

  <!-- Avatar Skeleton -->
  <div
    v-else-if="type === 'avatar'"
    :class="[baseClass, avatarSizeClass[avatarSize], avatarShape === 'circle' ? 'rounded-full' : 'rounded-lg']"
  />

  <!-- Card Skeleton -->
  <div
    v-else-if="type === 'card'"
    class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3"
  >
    <!-- Header with avatar and title -->
    <div class="flex items-center space-x-3">
      <div
        :class="[baseClass, avatarSizeClass.avatarSize]"
        :style="{ width: '2.5rem', height: '2.5rem' }"
      />
      <div class="flex-1 space-y-2">
        <div
          :class="[baseClass, 'w-1/2']"
          style="height: 0.875rem"
        />
        <div
          :class="[baseClass, 'w-1/3']"
          style="height: 0.75rem"
        />
      </div>
    </div>
    <!-- Content lines -->
    <div
      v-for="i in lines"
      :key="i"
      :class="[baseClass, i === lines ? 'w-3/4' : 'w-full']"
      style="height: 1rem"
    />
    <!-- Action button -->
    <div
      :class="[baseClass, 'w-24 mt-2']"
      style="height: 2rem"
    />
  </div>

  <!-- Table Skeleton -->
  <div
    v-else-if="type === 'table'"
    class="w-full"
  >
    <!-- Header -->
    <div class="flex border-b border-gray-200 dark:border-gray-700">
      <div
        v-for="i in tableColumns"
        :key="`header-${i}`"
        :class="[baseClass, 'flex-1 py-3']"
      />
    </div>
    <!-- Rows -->
    <div
      v-for="row in tableRows"
      :key="`row-${row}`"
      class="flex border-b border-gray-100 dark:border-gray-800"
    >
      <div
        v-for="col in tableColumns"
        :key="`cell-${row}-${col}`"
        class="flex-1 py-3 px-1"
      >
        <div
          :class="[baseClass, 'w-full']"
          style="height: 0.875rem"
        />
      </div>
    </div>
  </div>

  <!-- Title Skeleton -->
  <div
    v-else-if="type === 'title'"
    class="space-y-2"
  >
    <div
      :class="[baseClass, 'w-1/2']"
      style="height: 1.5rem"
    />
    <div
      :class="[baseClass, 'w-1/3']"
      style="height: 1rem"
    />
  </div>
</template>
