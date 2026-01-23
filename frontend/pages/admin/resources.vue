<template>
  <div class="p-4 md:p-6 max-w-7xl mx-auto">
    <div class="mb-6 md:mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
        資源管理
      </h1>
      <p class="text-slate-400 text-sm md:text-base">
        管理教室、課程、待排課程、老師
      </p>
    </div>

    <div class="mb-6 flex flex-col sm:flex-row gap-3 overflow-x-auto pb-2">
      <button
        @click="activeTab = 'rooms'"
        class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
        :class="activeTab === 'rooms' ? 'bg-primary-500/30 border-primary-500' : ''"
      >
        教室
      </button>
      <button
        @click="activeTab = 'courses'"
        class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
        :class="activeTab === 'courses' ? 'bg-primary-500/30 border-primary-500' : ''"
      >
        課程
      </button>
      <button
        @click="activeTab = 'offerings'"
        class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
        :class="activeTab === 'offerings' ? 'bg-primary-500/30 border-primary-500' : ''"
      >
        待排課程
      </button>
      <button
        @click="activeTab = 'teachers'"
        class="glass-btn px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap"
        :class="activeTab === 'teachers' ? 'bg-primary-500/30 border-primary-500' : ''"
      >
        老師
      </button>
    </div>

    <RoomsTab v-if="activeTab === 'rooms'" />
    <CoursesTab v-if="activeTab === 'courses'" />
    <OfferingsTab v-if="activeTab === 'offerings'" />
    <TeachersTab v-if="activeTab === 'teachers'" />
  </div>

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
 definePageMeta({
   middleware: 'auth-admin',
   layout: 'admin',
 })

 const activeTab = ref('rooms')
const notificationUI = useNotification()
</script>
