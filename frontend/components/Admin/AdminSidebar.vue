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
              <Icon icon="calendar" size="lg" />
              <span>排課管理</span>
            </div>
            <Icon
              icon="chevron-down"
              size="sm"
              :class="{ 'rotate-180': expandedMenus.scheduling }"
            />
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
          <Icon icon="check-circle" size="lg" />
          <span>審核管理</span>
        </NuxtLink>

        <!-- 人才庫 -->
        <div>
          <button
            @click="toggleSubmenu('talent')"
            class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
          >
            <div class="flex items-center gap-3">
              <Icon icon="people" size="lg" />
              <span>人才庫</span>
            </div>
            <Icon
              icon="chevron-down"
              size="sm"
              :class="{ 'rotate-180': expandedMenus.talent }"
            />
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
          <Icon icon="template" size="lg" />
          <span>範本管理</span>
        </NuxtLink>

        <!-- 資源管理 -->
        <NuxtLink
          to="/admin/resources"
          class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
          active-class="bg-primary-500/20 !text-primary-400"
        >
          <Icon icon="resource" size="lg" />
          <span>資源管理</span>
        </NuxtLink>

        <!-- 系統設定 -->
        <div>
          <button
            @click="toggleSubmenu('settings')"
            class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
          >
            <div class="flex items-center gap-3">
              <Icon icon="settings" size="lg" />
              <span>系統設定</span>
            </div>
            <Icon
              icon="chevron-down"
              size="sm"
              :class="{ 'rotate-180': expandedMenus.settings }"
            />
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
            <Icon icon="close" size="lg" />
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
                <Icon icon="calendar" size="lg" />
                <span>排課管理</span>
              </div>
              <Icon
                icon="chevron-down"
                size="sm"
                :class="{ 'rotate-180': expandedMenus.scheduling }"
              />
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
            <Icon icon="check-circle" size="lg" />
            <span>審核管理</span>
          </NuxtLink>

          <!-- 人才庫 -->
          <div>
            <button
              @click="toggleSubmenu('talent')"
              class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
            >
              <div class="flex items-center gap-3">
                <Icon icon="people" size="lg" />
                <span>人才庫</span>
              </div>
              <Icon
                icon="chevron-down"
                size="sm"
                :class="{ 'rotate-180': expandedMenus.talent }"
              />
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
            <Icon icon="template" size="lg" />
            <span>範本管理</span>
          </NuxtLink>

          <!-- 資源管理 -->
          <NuxtLink
            to="/admin/resources"
            class="flex items-center gap-3 px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
            active-class="bg-primary-500/20 !text-primary-400"
            @click="$emit('close')"
          >
            <Icon icon="resource" size="lg" />
            <span>資源管理</span>
          </NuxtLink>

          <!-- 系統設定 -->
          <div>
            <button
              @click="toggleSubmenu('settings')"
              class="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
            >
              <div class="flex items-center gap-3">
                <Icon icon="settings" size="lg" />
                <span>系統設定</span>
              </div>
              <Icon
                icon="chevron-down"
                size="sm"
                :class="{ 'rotate-180': expandedMenus.settings }"
              />
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
