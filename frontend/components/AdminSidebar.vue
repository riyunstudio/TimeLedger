<template>
  <div class="contents">
    <!-- Desktop sidebar (always visible) -->
    <div class="hidden md:block w-64 bg-slate-800/50 border-r border-slate-700 h-screen flex flex-col shrink-0">
      <div class="p-6 border-b border-slate-700">
        <h2 class="text-white font-bold text-lg">Admin Sidebar</h2>
      </div>

      <nav class="flex-1 flex-col overflow-y-auto p-4 space-y-2">
        <NuxtLink
          to="/admin/dashboard"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          Dashboard
        </NuxtLink>
        <NuxtLink
          to="/admin/schedules"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          Schedules
        </NuxtLink>
        <NuxtLink
          to="/admin/templates"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          Templates
        </NuxtLink>
        <NuxtLink
          to="/admin/resources"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          Resources
        </NuxtLink>
        <NuxtLink
          to="/admin/approval"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          Approvals
        </NuxtLink>
        <NuxtLink
          to="/admin/matching"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          Matching
        </NuxtLink>
        <NuxtLink
          to="/admin/holidays"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          Holidays
        </NuxtLink>
        <NuxtLink
          to="/admin/teacher-ratings"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          Teacher Ratings
        </NuxtLink>
        <NuxtLink
          to="/admin/line-bind"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          LINE 通知
        </NuxtLink>
        <NuxtLink
          to="/admin/queue-monitor"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          系統監控
        </NuxtLink>
        <NuxtLink
          to="/admin/admin-list"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          管理員
        </NuxtLink>
        <NuxtLink
          to="/admin/settings"
          class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
          active-class="bg-primary-500/20 text-primary-400"
        >
          Settings
        </NuxtLink>
      </nav>
    </div>

    <!-- Mobile sidebar overlay -->
    <Teleport to="body">
      <div
        v-if="isOpen"
        class="fixed inset-y-0 left-0 z-50 w-64 bg-slate-800/95 border-r border-slate-700 h-full flex flex-col transform transition-transform duration-300 ease-in-out"
      >
        <!-- Mobile header with close button -->
        <div class="flex items-center justify-between p-4 border-b border-slate-700">
          <h2 class="text-white font-bold text-lg">Admin Sidebar</h2>
          <button @click="$emit('close')" class="p-2 rounded-lg hover:bg-white/10 text-slate-300">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <nav class="flex-1 flex-col overflow-y-auto p-4 space-y-2">
          <NuxtLink
            to="/admin/dashboard"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            Dashboard
          </NuxtLink>
          <NuxtLink
            to="/admin/schedules"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            Schedules
          </NuxtLink>
          <NuxtLink
            to="/admin/templates"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            Templates
          </NuxtLink>
          <NuxtLink
            to="/admin/resources"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            Resources
          </NuxtLink>
          <NuxtLink
            to="/admin/approval"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            Approvals
          </NuxtLink>
          <NuxtLink
            to="/admin/matching"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            Matching
          </NuxtLink>
          <NuxtLink
            to="/admin/holidays"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            Holidays
          </NuxtLink>
          <NuxtLink
            to="/admin/teacher-ratings"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            Teacher Ratings
          </NuxtLink>
          <NuxtLink
            to="/admin/line-bind"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            LINE 通知
          </NuxtLink>
          <NuxtLink
            to="/admin/queue-monitor"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            系統監控
          </NuxtLink>
          <NuxtLink
            to="/admin/admin-list"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            管理員
          </NuxtLink>
          <NuxtLink
            to="/admin/settings"
            class="block px-4 py-2 rounded-lg bg-slate-700/50 hover:bg-slate-600 text-slate-300 hover:text-white transition-colors"
            active-class="bg-primary-500/20 text-primary-400"
            @click="$emit('close')"
          >
            Settings
          </NuxtLink>
        </nav>
      </div>

      <!-- Mobile overlay backdrop -->
      <div
        v-if="isOpen"
        class="fixed inset-0 bg-black/50 z-40"
        @click="$emit('close')"
      />
    </Teleport>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  isOpen: boolean
}>()

defineEmits<{
  close: []
}>()
</script>
