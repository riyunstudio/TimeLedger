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
        class="px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-colors"
        :class="activeTab === 'rooms' ? 'bg-primary-500/30 border border-primary-500 text-primary-400' : 'bg-white/5 text-slate-300 hover:bg-white/10'"
      >
        教室
      </button>
      <button
        @click="activeTab = 'courses'"
        class="px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-colors"
        :class="activeTab === 'courses' ? 'bg-primary-500/30 border border-primary-500 text-primary-400' : 'bg-white/5 text-slate-300 hover:bg-white/10'"
      >
        課程
      </button>
      <button
        @click="activeTab = 'offerings'"
        class="px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-colors"
        :class="activeTab === 'offerings' ? 'bg-primary-500/30 border border-primary-500 text-primary-400' : 'bg-white/5 text-slate-300 hover:bg-white/10'"
      >
        待排課程
      </button>
      <button
        @click="activeTab = 'teachers'"
        class="px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-colors"
        :class="activeTab === 'teachers' ? 'bg-primary-500/30 border border-primary-500 text-primary-400' : 'bg-white/5 text-slate-300 hover:bg-white/10'"
      >
        老師
      </button>
      <button
        @click="activeTab = 'terms'"
        class="px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-colors"
        :class="activeTab === 'terms' ? 'bg-primary-500/30 border border-primary-500 text-primary-400' : 'bg-white/5 text-slate-300 hover:bg-white/10'"
      >
        學期期間
      </button>
    </div>

    <!-- Tab Content -->
    <RoomsTab v-if="activeTab === 'rooms'" />
    <CoursesTab v-else-if="activeTab === 'courses'" />
    <OfferingsTab v-else-if="activeTab === 'offerings'" />
    <TeachersTab v-else-if="activeTab === 'teachers'" />
    <TermsTab v-else-if="activeTab === 'terms'" />
  </div>
</template>

<script setup lang="ts">
import RoomsTab from '~/components/Admin/RoomsTab.vue'
import CoursesTab from '~/components/Admin/CoursesTab.vue'
import OfferingsTab from '~/components/Admin/OfferingsTab.vue'
import TeachersTab from '~/components/Admin/TeachersTab.vue'
import TermsTab from '~/components/Admin/TermsTab.vue'

definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

const activeTab = ref('rooms')
</script>
