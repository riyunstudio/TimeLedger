<template>
  <div class="contents">
    <!-- Desktop sidebar (always visible) -->
    <div class="hidden md:block w-64 bg-slate-800/50 border-r border-slate-700 h-screen flex flex-col shrink-0">
      <div class="p-6 border-b border-slate-700">
        <h2 class="text-white font-bold text-lg">TimeLedger</h2>
      </div>

      <nav class="flex-1 flex-col overflow-y-auto p-4 space-y-1">
        <!-- 排課管理 -->
        <div>
          <button
            @click="toggleSubmenu('scheduling')"
            class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
          >
            <div class="flex items-center gap-3">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span>排課管理</span>
            </div>
            <svg
              class="w-4 h-4 transition-transform duration-200"
              :class="{ 'rotate-180': expandedMenus.scheduling }"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          <div
            v-show="expandedMenus.scheduling"
            class="ml-4 mt-1 space-y-1 border-l border-slate-700 pl-3"
          >
            <NuxtLink
              to="/admin/dashboard"
              class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
              active-class="!text-primary-400 bg-primary-500/10"
            >
              首頁
            </NuxtLink>
            <NuxtLink
              to="/admin/schedules"
              class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
              active-class="!text-primary-400 bg-primary-500/10"
            >
              排課管理
            </NuxtLink>
          </div>
        </div>

        <!-- 審核管理 -->
        <NuxtLink
          to="/admin/approval"
          class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
          active-class="bg-primary-500/20 !text-primary-400"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>審核管理</span>
        </NuxtLink>

        <!-- 人才庫 -->
        <div>
          <button
            @click="toggleSubmenu('talent')"
            class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
          >
            <div class="flex items-center gap-3">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
              <span>人才庫</span>
            </div>
            <svg
              class="w-4 h-4 transition-transform duration-200"
              :class="{ 'rotate-180': expandedMenus.talent }"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          <div
            v-show="expandedMenus.talent"
            class="ml-4 mt-1 space-y-1 border-l border-slate-700 pl-3"
          >
            <NuxtLink
              to="/admin/matching"
              class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
              active-class="!text-primary-400 bg-primary-500/10"
            >
              智慧媒合
            </NuxtLink>
            <NuxtLink
              to="/admin/teacher-ratings"
              class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
              active-class="!text-primary-400 bg-primary-500/10"
            >
              老師評價
            </NuxtLink>
          </div>
        </div>

        <!-- 範本管理 -->
        <NuxtLink
          to="/admin/templates"
          class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
          active-class="bg-primary-500/20 !text-primary-400"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
          </svg>
          <span>範本管理</span>
        </NuxtLink>

        <!-- 資源管理 -->
        <NuxtLink
          to="/admin/resources"
          class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
          active-class="bg-primary-500/20 !text-primary-400"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
          </svg>
          <span>資源管理</span>
        </NuxtLink>

        <!-- 系統設定 -->
        <div>
          <button
            @click="toggleSubmenu('settings')"
            class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
          >
            <div class="flex items-center gap-3">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span>系統設定</span>
            </div>
            <svg
              class="w-4 h-4 transition-transform duration-200"
              :class="{ 'rotate-180': expandedMenus.settings }"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          <div
            v-show="expandedMenus.settings"
            class="ml-4 mt-1 space-y-1 border-l border-slate-700 pl-3"
          >
            <NuxtLink
              to="/admin/settings"
              class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
              active-class="!text-primary-400 bg-primary-500/10"
            >
              系統設定
            </NuxtLink>
            <NuxtLink
              to="/admin/admin-list"
              class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
              active-class="!text-primary-400 bg-primary-500/10"
            >
              管理員
            </NuxtLink>
            <NuxtLink
              to="/admin/line-bind"
              class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
              active-class="!text-primary-400 bg-primary-500/10"
            >
              LINE 綁定
            </NuxtLink>
            <NuxtLink
              to="/admin/queue-monitor"
              class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
              active-class="!text-primary-400 bg-primary-500/10"
            >
              系統監控
            </NuxtLink>
            <NuxtLink
              to="/admin/holidays"
              class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
              active-class="!text-primary-400 bg-primary-500/10"
            >
              假日設定
            </NuxtLink>
          </div>
        </div>
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
          <h2 class="text-white font-bold text-lg">TimeLedger</h2>
          <button @click="$emit('close')" class="p-2 rounded-lg hover:bg-white/10 text-slate-300">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <nav class="flex-1 flex-col overflow-y-auto p-4 space-y-1">
          <!-- 排課管理 -->
          <div>
            <button
              @click="toggleSubmenu('scheduling')"
              class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
            >
              <div class="flex items-center gap-3">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <span>排課管理</span>
              </div>
              <svg
                class="w-4 h-4 transition-transform duration-200"
                :class="{ 'rotate-180': expandedMenus.scheduling }"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            <div
              v-show="expandedMenus.scheduling"
              class="ml-4 mt-1 space-y-1 border-l border-slate-700 pl-3"
            >
              <NuxtLink
                to="/admin/dashboard"
                class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
                active-class="!text-primary-400 bg-primary-500/10"
                @click="$emit('close')"
              >
                首頁
              </NuxtLink>
              <NuxtLink
                to="/admin/schedules"
                class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
                active-class="!text-primary-400 bg-primary-500/10"
                @click="$emit('close')"
              >
                排課管理
              </NuxtLink>
            </div>
          </div>

          <!-- 審核管理 -->
          <NuxtLink
            to="/admin/approval"
            class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
            active-class="bg-primary-500/20 !text-primary-400"
            @click="$emit('close')"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>審核管理</span>
          </NuxtLink>

          <!-- 人才庫 -->
          <div>
            <button
              @click="toggleSubmenu('talent')"
              class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
            >
              <div class="flex items-center gap-3">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
                <span>人才庫</span>
              </div>
              <svg
                class="w-4 h-4 transition-transform duration-200"
                :class="{ 'rotate-180': expandedMenus.talent }"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            <div
              v-show="expandedMenus.talent"
              class="ml-4 mt-1 space-y-1 border-l border-slate-700 pl-3"
            >
              <NuxtLink
                to="/admin/matching"
                class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
                active-class="!text-primary-400 bg-primary-500/10"
                @click="$emit('close')"
              >
                智慧媒合
              </NuxtLink>
              <NuxtLink
                to="/admin/teacher-ratings"
                class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
                active-class="!text-primary-400 bg-primary-500/10"
                @click="$emit('close')"
              >
                老師評價
              </NuxtLink>
            </div>
          </div>

          <!-- 範本管理 -->
          <NuxtLink
            to="/admin/templates"
            class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
            active-class="bg-primary-500/20 !text-primary-400"
            @click="$emit('close')"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
            </svg>
            <span>範本管理</span>
          </NuxtLink>

          <!-- 資源管理 -->
          <NuxtLink
            to="/admin/resources"
            class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
            active-class="bg-primary-500/20 !text-primary-400"
            @click="$emit('close')"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
            </svg>
            <span>資源管理</span>
          </NuxtLink>

          <!-- 系統設定 -->
          <div>
            <button
              @click="toggleSubmenu('settings')"
              class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
            >
              <div class="flex items-center gap-3">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                <span>系統設定</span>
              </div>
              <svg
                class="w-4 h-4 transition-transform duration-200"
                :class="{ 'rotate-180': expandedMenus.settings }"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            <div
              v-show="expandedMenus.settings"
              class="ml-4 mt-1 space-y-1 border-l border-slate-700 pl-3"
            >
              <NuxtLink
                to="/admin/settings"
                class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
                active-class="!text-primary-400 bg-primary-500/10"
                @click="$emit('close')"
              >
                系統設定
              </NuxtLink>
              <NuxtLink
                to="/admin/admin-list"
                class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
                active-class="!text-primary-400 bg-primary-500/10"
                @click="$emit('close')"
              >
                管理員
              </NuxtLink>
              <NuxtLink
                to="/admin/line-bind"
                class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
                active-class="!text-primary-400 bg-primary-500/10"
                @click="$emit('close')"
              >
                LINE 綁定
              </NuxtLink>
              <NuxtLink
                to="/admin/queue-monitor"
                class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
                active-class="!text-primary-400 bg-primary-500/10"
                @click="$emit('close')"
              >
                系統監控
              </NuxtLink>
              <NuxtLink
                to="/admin/holidays"
                class="block px-4 py-2 rounded-lg text-sm text-slate-400 hover:text-white hover:bg-slate-700/30 transition-colors"
                active-class="!text-primary-400 bg-primary-500/10"
                @click="$emit('close')"
              >
                假日設定
              </NuxtLink>
            </div>
          </div>
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

// 子選單展開狀態
const expandedMenus = reactive({
  scheduling: true,  // 預設展開排課管理
  talent: false,
  settings: false
})

// 切換子選單展開/收合
function toggleSubmenu(menu: keyof typeof expandedMenus) {
  expandedMenus[menu] = !expandedMenus[menu]
}
</script>
