<template>
  <div class="min-h-screen bg-slate-900">
    <AdminHeader />

    <main class="p-6 max-w-7xl mx-auto">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-slate-100 mb-2">
          資源管理
        </h1>
        <p class="text-slate-400">
          管理教室、課程、待排課程、老師
        </p>
      </div>

      <div class="mb-6 flex gap-3">
        <button
          @click="activeTab = 'rooms'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium"
          :class="activeTab === 'rooms' ? 'bg-primary-500/30 border-primary-500' : ''"
        >
          教室
        </button>
        <button
          @click="activeTab = 'courses'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium"
          :class="activeTab === 'courses' ? 'bg-primary-500/30 border-primary-500' : ''"
        >
          課程
        </button>
        <button
          @click="activeTab = 'offerings'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium"
          :class="activeTab === 'offerings' ? 'bg-primary-500/30 border-primary-500' : ''"
        >
          待排課程
        </button>
        <button
          @click="activeTab = 'teachers'"
          class="glass-btn px-4 py-2 rounded-xl text-sm font-medium"
          :class="activeTab === 'teachers' ? 'bg-primary-500/30 border-primary-500' : ''"
        >
          老師
        </button>
      </div>

      <RoomsTab v-if="activeTab === 'rooms'" />
      <CoursesTab v-if="activeTab === 'courses'" />
      <OfferingsTab v-if="activeTab === 'offerings'" />
      <TeachersTab v-if="activeTab === 'teachers'" />
    </main>

    <NotificationDropdown
      v-if="notificationUI.show.value"
      @close="notificationUI.close()"
    />
  </div>
</template>

<script setup lang="ts">
 definePageMeta({
   middleware: 'auth-admin',
   layout: 'admin',
 })

 const activeTab = ref('rooms')
const notificationUI = useNotification()
</script>
