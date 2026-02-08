<template>
  <div class="min-h-screen bg-slate-900 flex">
    <!-- Desktop Sidebar -->
    <aside class="hidden lg:flex w-64 bg-slate-800/50 border-r border-slate-700 flex-col shrink-0 relative z-20">
      <!-- Logo 區域 -->
      <div class="p-5 border-b border-slate-700">
        <NuxtLink to="/admin/dashboard" class="flex items-center gap-3 group">
          <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center group-hover:scale-105 transition-transform duration-200">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <h1 class="text-lg font-bold text-white group-hover:text-primary-300 transition-colors">TimeLedger</h1>
            <p class="text-xs text-slate-400">排課管理平台</p>
          </div>
        </NuxtLink>
      </div>

      <!-- 導航選單 -->
      <nav class="flex-1 overflow-y-auto p-3 space-y-1">
        <!-- 排課管理 -->
        <div>
          <button
            @click="toggleSubmenu('scheduling')"
            class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200 group"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg bg-slate-700/50 flex items-center justify-center group-hover:bg-primary-500/20 transition-colors">
                <svg class="w-4 h-4 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <span class="font-medium">排課管理</span>
            </div>
            <svg
              class="w-4 h-4 text-slate-500 transition-transform duration-200"
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
            class="ml-4 mt-1 space-y-0.5"
          >
            <NuxtLink
              to="/admin/dashboard"
              class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
              active-class="!text-white bg-primary-500/10"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
              首頁
            </NuxtLink>
            <NuxtLink
              to="/admin/schedules"
              class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
              active-class="!text-white bg-primary-500/10"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
              排課管理
            </NuxtLink>
            <NuxtLink
              to="/admin/resource-occupancy"
              class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
              active-class="!text-white bg-primary-500/10"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
              資源佔用表
            </NuxtLink>
          </div>
        </div>

        <!-- 審核管理 -->
        <NuxtLink
          to="/admin/approval"
          class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200 group"
          active-class="bg-primary-500/20 !text-primary-400"
        >
          <div class="w-8 h-8 rounded-lg bg-slate-700/50 flex items-center justify-center group-hover:bg-amber-500/20 transition-colors">
            <svg class="w-4 h-4 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <span class="font-medium">審核管理</span>
        </NuxtLink>

        <!-- 一鍵公告 -->
        <NuxtLink
          to="/admin/broadcast"
          class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200 group"
          active-class="bg-primary-500/20 !text-primary-400"
        >
          <div class="w-8 h-8 rounded-lg bg-slate-700/50 flex items-center justify-center group-hover:bg-orange-500/20 transition-colors">
            <svg class="w-4 h-4 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />
            </svg>
          </div>
          <span class="font-medium">一鍵公告</span>
        </NuxtLink>

        <!-- 人才庫 -->
        <div>
          <button
            @click="toggleSubmenu('talent')"
            class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200 group"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg bg-slate-700/50 flex items-center justify-center group-hover:bg-violet-500/20 transition-colors">
                <svg class="w-4 h-4 text-violet-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <span class="font-medium">人才庫</span>
            </div>
            <svg
              class="w-4 h-4 text-slate-500 transition-transform duration-200"
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
            class="ml-4 mt-1 space-y-0.5"
          >
            <NuxtLink
              to="/admin/matching"
              class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
              active-class="!text-white bg-violet-500/10"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
              智慧媒合
            </NuxtLink>
            <NuxtLink
              to="/admin/invitations"
              class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
              active-class="!text-white bg-violet-500/10"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
              邀請紀錄
            </NuxtLink>
          </div>
        </div>

        <!-- 範本管理 -->
        <NuxtLink
          to="/admin/templates"
          class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200 group"
          active-class="bg-primary-500/20 !text-primary-400"
        >
          <div class="w-8 h-8 rounded-lg bg-slate-700/50 flex items-center justify-center group-hover:bg-emerald-500/20 transition-colors">
            <svg class="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
            </svg>
          </div>
          <span class="font-medium">範本管理</span>
        </NuxtLink>

        <!-- 資源管理 -->
        <NuxtLink
          to="/admin/resources"
          class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200 group"
          active-class="bg-primary-500/20 !text-primary-400"
        >
          <div class="w-8 h-8 rounded-lg bg-slate-700/50 flex items-center justify-center group-hover:bg-cyan-500/20 transition-colors">
            <svg class="w-4 h-4 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
            </svg>
          </div>
          <span class="font-medium">資源管理</span>
        </NuxtLink>

        <!-- 系統設定 -->
        <div>
          <button
            @click="toggleSubmenu('settings')"
            class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200 group"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg bg-slate-700/50 flex items-center justify-center group-hover:bg-slate-600 transition-colors">
                <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </div>
              <span class="font-medium">系統設定</span>
            </div>
            <svg
              class="w-4 h-4 text-slate-500 transition-transform duration-200"
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
            class="ml-4 mt-1 space-y-0.5"
          >
            <NuxtLink
              to="/admin/settings"
              class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
              active-class="!text-white bg-slate-700/50"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
              系統設定
            </NuxtLink>
            <NuxtLink
              to="/admin/admin-list"
              class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
              active-class="!text-white bg-slate-700/50"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
              管理員
            </NuxtLink>
            <NuxtLink
              to="/admin/queue-monitor"
              class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
              active-class="!text-white bg-slate-700/50"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
              系統監控
            </NuxtLink>
            <NuxtLink
              to="/admin/holidays"
              class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
              active-class="!text-white bg-slate-700/50"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-slate-500" />
              假日設定
            </NuxtLink>
          </div>
        </div>
      </nav>

      <!-- 用戶資訊 -->
      <div class="p-4 border-t border-slate-700">
        <div class="flex items-center gap-3 px-3 py-2.5 rounded-lg bg-slate-700/30">
          <div class="w-9 h-9 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center text-white font-medium text-sm">
            管
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-white truncate">超級管理員</p>
            <p class="text-xs text-slate-400 truncate">admin@timeledger.com</p>
          </div>
        </div>
      </div>
    </aside>

    <!-- Mobile Sidebar Overlay -->
    <Teleport to="body">
      <Transition name="slide">
        <div
          v-if="mobileMenuOpen"
          class="fixed inset-y-0 left-0 z-50 w-72 bg-slate-800 border-r border-slate-700 flex flex-col lg:hidden"
        >
          <!-- Mobile Header -->
          <div class="p-5 border-b border-slate-700 flex items-center justify-between">
            <NuxtLink to="/admin/dashboard" class="flex items-center gap-3 group" @click="mobileMenuOpen = false">
              <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center group-hover:scale-105 transition-transform duration-200">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <h1 class="text-lg font-bold text-white group-hover:text-primary-300 transition-colors">TimeLedger</h1>
                <p class="text-xs text-slate-400">排課管理平台</p>
              </div>
            </NuxtLink>
            <button
              @click="mobileMenuOpen = false"
              class="p-2 rounded-lg hover:bg-slate-700 transition-colors"
            >
              <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Mobile Navigation -->
          <nav class="flex-1 overflow-y-auto p-3 space-y-1">
            <!-- 排課管理 -->
            <div>
              <button
                @click="toggleSubmenu('scheduling')"
                class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200"
              >
                <div class="flex items-center gap-3">
                  <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  <span class="font-medium">排課管理</span>
                </div>
                <svg
                  class="w-4 h-4 text-slate-500 transition-transform duration-200"
                  :class="{ 'rotate-180': expandedMenus.scheduling }"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div v-show="expandedMenus.scheduling" class="ml-4 mt-1 space-y-0.5">
                <NuxtLink
                  to="/admin/dashboard"
                  class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
                  active-class="!text-white bg-primary-500/10"
                  @click="mobileMenuOpen = false"
                >
                  首頁
                </NuxtLink>
                <NuxtLink
                  to="/admin/schedules"
                  class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
                  active-class="!text-white bg-primary-500/10"
                  @click="mobileMenuOpen = false"
                >
                  排課管理
                </NuxtLink>
                <NuxtLink
                  to="/admin/resource-occupancy"
                  class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
                  active-class="!text-white bg-primary-500/10"
                  @click="mobileMenuOpen = false"
                >
                  資源佔用表
                </NuxtLink>
              </div>
            </div>

            <!-- 審核管理 -->
            <NuxtLink
              to="/admin/approval"
              class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200"
              active-class="bg-primary-500/20 !text-primary-400"
              @click="mobileMenuOpen = false"
            >
              <svg class="w-5 h-5 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="font-medium">審核管理</span>
            </NuxtLink>

            <!-- 一鍵公告 -->
            <NuxtLink
              to="/admin/broadcast"
              class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200"
              active-class="bg-primary-500/20 !text-primary-400"
              @click="mobileMenuOpen = false"
            >
              <svg class="w-5 h-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />
              </svg>
              <span class="font-medium">一鍵公告</span>
            </NuxtLink>

            <!-- 人才庫 -->
            <div>
              <button
                @click="toggleSubmenu('talent')"
                class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200"
              >
                <div class="flex items-center gap-3">
                  <svg class="w-5 h-5 text-violet-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  <span class="font-medium">人才庫</span>
                </div>
                <svg
                  class="w-4 h-4 text-slate-500 transition-transform duration-200"
                  :class="{ 'rotate-180': expandedMenus.talent }"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div v-show="expandedMenus.talent" class="ml-4 mt-1 space-y-0.5">
                <NuxtLink
                  to="/admin/matching"
                  class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
                  active-class="!text-white bg-violet-500/10"
                  @click="mobileMenuOpen = false"
                >
                  智慧媒合
                </NuxtLink>
                <NuxtLink
                  to="/admin/invitations"
                  class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
                  active-class="!text-white bg-violet-500/10"
                  @click="mobileMenuOpen = false"
                >
                  邀請紀錄
                </NuxtLink>
              </div>
            </div>

            <!-- 範本管理 -->
            <NuxtLink
              to="/admin/templates"
              class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200"
              active-class="bg-primary-500/20 !text-primary-400"
              @click="mobileMenuOpen = false"
            >
              <svg class="w-5 h-5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5z" />
              </svg>
              <span class="font-medium">範本管理</span>
            </NuxtLink>

            <!-- 資源管理 -->
            <NuxtLink
              to="/admin/resources"
              class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200"
              active-class="bg-primary-500/20 !text-primary-400"
              @click="mobileMenuOpen = false"
            >
              <svg class="w-5 h-5 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
              </svg>
              <span class="font-medium">資源管理</span>
            </NuxtLink>

            <!-- 系統設定 -->
            <div>
              <button
                @click="toggleSubmenu('settings')"
                class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-all duration-200"
              >
                <div class="flex items-center gap-3">
                  <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  </svg>
                  <span class="font-medium">系統設定</span>
                </div>
                <svg
                  class="w-4 h-4 text-slate-500 transition-transform duration-200"
                  :class="{ 'rotate-180': expandedMenus.settings }"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div v-show="expandedMenus.settings" class="ml-4 mt-1 space-y-0.5">
                <NuxtLink
                  to="/admin/settings"
                  class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
                  active-class="!text-white bg-slate-700/50"
                  @click="mobileMenuOpen = false"
                >
                  系統設定
                </NuxtLink>
                <NuxtLink
                  to="/admin/admin-list"
                  class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
                  active-class="!text-white bg-slate-700/50"
                  @click="mobileMenuOpen = false"
                >
                  管理員
                </NuxtLink>
                <NuxtLink
                  to="/admin/queue-monitor"
                  class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
                  active-class="!text-white bg-slate-700/50"
                  @click="mobileMenuOpen = false"
                >
                  系統監控
                </NuxtLink>
                <NuxtLink
                  to="/admin/holidays"
                  class="flex items-center gap-3 px-4 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/30 transition-all duration-200 text-sm"
                  active-class="!text-white bg-slate-700/50"
                  @click="mobileMenuOpen = false"
                >
                  假日設定
                </NuxtLink>
              </div>
            </div>
          </nav>

          <!-- Mobile User Info -->
          <div class="p-4 border-t border-slate-700">
            <div class="flex items-center gap-3 px-3 py-2.5 rounded-lg bg-slate-700/30">
              <div class="w-9 h-9 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center text-white font-medium text-sm">
                管
              </div>
              <div class="flex-1">
                <p class="text-sm font-medium text-white">超級管理員</p>
                <p class="text-xs text-slate-400">admin@timeledger.com</p>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Mobile Backdrop -->
      <Transition name="fade">
        <div
          v-if="mobileMenuOpen"
          class="fixed inset-0 bg-black/60 z-40 lg:hidden"
          @click="mobileMenuOpen = false"
        />
      </Transition>
    </Teleport>

    <!-- 主內容區域 -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- 頂部 Header -->
      <header class="h-16 bg-slate-800/50 border-b border-slate-700 flex items-center justify-between px-4 lg:px-6">
        <!-- 左側：選單按鈕 + 標題 -->
        <div class="flex items-center gap-3">
          <!-- Mobile Menu Button -->
          <button
            @click="mobileMenuOpen = true"
            class="lg:hidden p-2 rounded-lg hover:bg-slate-700/50 transition-colors"
          >
            <svg class="w-6 h-6 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <h2 class="text-lg font-semibold text-white">{{ pageTitle }}</h2>
        </div>

        <!-- 右側功能 -->
        <div class="flex items-center gap-2">
          <!-- 通知按鈕 -->
          <button
            @click="notificationUI.open()"
            class="relative p-2 rounded-lg hover:bg-slate-700/50 transition-colors"
          >
            <svg class="w-5 h-5 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2 2 0 0118 22v-4.317l-1.405 1.405A2 2 0 0115 17h5l-1.405-1.405A2 2 0 0118 12v6a2 2 0 01-2 2h-2m-6 3h2m8 0h2M3 8h18M3 8v10a2 2 0 002 2h14a2 2 0 002-2V8z" />
            </svg>
            <span
              v-if="notificationStore.unreadCount > 0"
              class="absolute -top-0.5 -right-0.5 w-5 h-5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center font-medium"
            >
              {{ notificationStore.unreadCount > 9 ? '9+' : notificationStore.unreadCount }}
            </span>
          </button>

          <!-- 登出按鈕 -->
          <button
            @click="handleLogout"
            class="flex items-center gap-2 px-3 py-2 rounded-lg text-slate-300 hover:text-red-400 hover:bg-red-500/10 transition-all duration-200"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            <span class="text-sm font-medium hidden sm:inline">登出</span>
          </button>
        </div>
      </header>

      <!-- 主要內容 -->
      <main class="flex-1 overflow-auto p-4 lg:p-6">
        <slot />
      </main>
    </div>

    <!-- 通知彈窗 -->
    <ToastNotification ref="toastRef" />
    <NotificationDropdown
      v-if="notificationUI.show.value"
      @close="notificationUI.close()"
    />
  </div>
</template>

<script setup lang="ts">
import ToastNotification from '~/components/base/ToastNotification.vue'
import NotificationDropdown from '~/components/Navigation/NotificationDropdown.vue'
import { registerToast } from '~/composables/useToast'
import { alertConfirm } from '~/composables/useAlert'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const notificationUI = useNotification()
const notificationStore = useNotificationStore()

const toastRef = ref<any>(null)

// Mobile menu state
const mobileMenuOpen = ref(false)

// 子選單展開狀態
const expandedMenus = reactive({
  scheduling: true,
  talent: false,
  settings: false
})

// 頁面標題
const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/admin/dashboard': '首頁',
    '/admin/schedules': '排課管理',
    '/admin/resource-occupancy': '資源佔用表',
    '/admin/approval': '審核管理',
    '/admin/broadcast': '一鍵公告',
    '/admin/matching': '智慧媒合',
    '/admin/templates': '範本管理',
    '/admin/resources': '資源管理',
    '/admin/settings': '系統設定',
    '/admin/admin-list': '管理員帳號',
    '/admin/queue-monitor': '系統監控',
    '/admin/holidays': '假日設定',
    '/admin/invitations': '邀請紀錄'
  }
  return titles[route.path] || 'TimeLedger'
})

// 切換子選單
function toggleSubmenu(menu: keyof typeof expandedMenus) {
  expandedMenus[menu] = !expandedMenus[menu]
}

// 自動展開選單邏輯
watch(() => route.path, (path) => {
  if (path.startsWith('/admin/settings') || path.startsWith('/admin/admin-list') || path.startsWith('/admin/queue-monitor') || path.startsWith('/admin/holidays')) {
    expandedMenus.settings = true
  }
  if (path.startsWith('/admin/dashboard') || path.startsWith('/admin/schedules') || path.startsWith('/admin/resource-occupancy')) {
    expandedMenus.scheduling = true
  }
  if (path.startsWith('/admin/matching') || path.startsWith('/admin/invitations')) {
    expandedMenus.talent = true
  }
}, { immediate: true })

// 登出
const handleLogout = async () => {
  if (await alertConfirm('確定要登出嗎？')) {
    authStore.logout()
    router.push('/admin/login')
  }
}

// 初始化通知
onMounted(() => {
  if (toastRef.value) {
    registerToast(toastRef.value)
  }
  notificationStore.fetchNotifications()
})
</script>

<style scoped>
/* Slide animation for mobile sidebar */
.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(-100%);
}

/* Fade animation for backdrop */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
