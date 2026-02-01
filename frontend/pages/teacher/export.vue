<template>
  <div class="mb-6">
    <div class="mb-4">
      <h2 class="text-2xl font-bold text-white mb-1">匯出課表</h2>
      <p class="text-slate-400 text-sm">預覽並下載您的詳細課表</p>
    </div>
    <!-- 按鈕區域 - 手機版自動換行 -->
    <div class="flex flex-wrap gap-2 sm:gap-3">
      <button
        @click="router.push('/teacher/dashboard')"
        class="px-3 py-2 sm:px-4 sm:py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors text-sm sm:text-base flex items-center gap-1 sm:gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        返回
      </button>
      <button
        @click="handleDownloadICal"
        class="px-3 py-2 sm:px-4 sm:py-2 rounded-lg bg-secondary-500 text-white hover:bg-secondary-600 transition-colors text-sm sm:text-base flex items-center gap-1 sm:gap-2"
        title="匯出到日曆 App"
      >
        <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <span class="hidden xs:inline">iCal</span>
        <span class="xs:hidden">日曆</span>
      </button>
      <button
        @click="handleShareLINE"
        class="px-3 py-2 sm:px-4 sm:py-2 rounded-lg bg-[#06C755] text-white hover:bg-[#05b546] transition-colors text-sm sm:text-base flex items-center gap-1 sm:gap-2"
        title="分享到 LINE"
      >
        <svg class="w-4 h-4 sm:w-5 sm:h-5" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12.6 9.5c-.17 1. 2zm41-1.3 2.08-2.5 2.55-.35.14-.6.22-.86.22-.32 0-.64-.1-.86-.35-.22-.25-.33-.55-.33-.9 0-.6.45-1.1 1-1.35.1-.05.18-.08.28-.08.55 0 1.05.45 1.1 1 .02.1.02.2.02.32zm-3.1 2.7c-.1.05-.18.08-.28.08-.55 0-1.05-.45-1.1-1-.02-.1-.02-.2-.02-.32 0-.6.45-1.1 1-1.35.1-.05.18-.08.28-.08.32 0 .64.1.86.35.22.25.33.55.33.9 0 .6-.45 1.1-1 1.1zm-3.1-2.7c.17-1.1 1.3-2.08 2.5-2.55.35-.14.6-.22.86-.22.32 0 .64.1.86.35.22.25.33.55.33.9 0 .6-.45 1.1-1 1.35-.1.05-.18.08-.28.08-.55 0-1.05-.45-1.1-1-.02-.1-.02-.2-.02-.32 0 .6.45 1.1 1 1.35.1.05.18.08.28.08.32 0 .64-.1.86-.35.22-.25.33-.55.33-.9 0-.6-.45-1.1-1-1.35-.1-.05-.18-.08-.28-.08-.55 0-1.05.45-1.1 1-.02.1-.02.2-.02.32 0-.6.45-1.1 1-1.35z"/>
        </svg>
        <span class="hidden xs:inline">LINE</span>
        <span class="xs:hidden">分享</span>
      </button>
      <button
        @click="handleDownloadPDF"
        class="px-3 py-2 sm:px-4 sm:py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors text-sm sm:text-base flex items-center gap-1 sm:gap-2"
      >
        <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        <span class="hidden xs:inline">下載 PDF</span>
        <span class="xs:hidden">PDF</span>
      </button>
      <button
        @click="handleExportAsImage"
        class="px-3 py-2 sm:px-4 sm:py-2 rounded-lg bg-secondary-500 text-white hover:bg-secondary-600 transition-colors text-sm sm:text-base flex items-center gap-1 sm:gap-2"
      >
        <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <span class="hidden xs:inline">下載圖片</span>
        <span class="xs:hidden">圖片</span>
      </button>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
    <div class="lg:col-span-3">
      <div ref="scheduleRef" class="glass-card p-8" :class="currentTheme.cardGradient">
        <div class="flex items-center justify-between mb-8 pb-6" :class="currentTheme.borderClass">
          <div class="flex items-center gap-4">
            <div class="w-16 h-16 rounded-full flex items-center justify-center" :class="currentTheme.avatarClass">
              <span class="text-2xl font-bold text-white">{{ authStore.user?.name?.charAt(0) || 'T' }}</span>
            </div>
            <div>
              <h1 class="text-2xl font-bold text-center" :class="currentTheme.titleClass">{{ authStore.user?.name }}</h1>
              <p class="text-sm text-center" :class="currentTheme.subtitleClass">個人課表</p>
              <div v-if="options.showPersonalInfo && teacherHashtags.length > 0" class="flex flex-wrap gap-2 mt-2 justify-center">
                <span
                  v-for="(tag, index) in teacherHashtags.slice(0, 5)"
                  :key="tag.id"
                  class="px-2 py-0.5 rounded-full text-xs"
                  :class="getTagClass(index)"
                >
                  {{ tag.name.replace('#', '') }}
                </span>
              </div>
            </div>
          </div>
          <div class="text-right">
            <p class="text-sm" :class="currentTheme.subtitleClass">匯出日期</p>
            <p class="font-medium" :class="currentTheme.titleClass">{{ formatDate(new Date()) }}</p>
          </div>
        </div>

        <div class="mb-6">
          <div class="flex items-center gap-2 mb-4">
            <div class="flex gap-1">
              <button
                @click="changeWeek(-1)"
                class="p-1 rounded transition-colors"
                :class="currentTheme.buttonClass"
              >
                <svg class="w-4 h-4" :class="currentTheme.iconClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
              </button>
            </div>
            <span class="font-medium" :class="currentTheme.titleClass">{{ weekLabel }}</span>
            <div class="flex gap-1">
              <button
                @click="changeWeek(1)"
                class="p-1 rounded transition-colors"
                :class="currentTheme.buttonClass"
              >
                <svg class="w-4 h-4" :class="currentTheme.iconClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </button>
            </div>
          </div>
          <!-- 視圖切換標籤 -->
          <div class="flex gap-2 mt-4">
            <button
              v-for="viewOption in viewOptions"
              :key="viewOption.id"
              @click="selectedViewOption = viewOption.id"
              class="px-3 py-1.5 rounded-lg text-sm transition-colors"
              :class="selectedViewOption === viewOption.id 
                ? 'bg-primary-500/30 text-primary-400 border border-primary-500/50' 
                : 'text-slate-400 hover:text-white hover:bg-white/5'"
            >
              {{ viewOption.name }}
            </button>
          </div>
        </div>

        <!-- 列表視圖 -->
        <div v-if="selectedViewOption === 'list'" class="space-y-6">
          <div
            v-for="day in scheduleDays"
            :key="day.date"
            class="rounded-xl overflow-hidden"
          >
            <div class="flex items-center justify-between px-4 py-3" :class="currentTheme.dayHeaderClass">
              <div class="flex items-center gap-3">
                <span class="font-semibold" :class="currentTheme.dayTextClass">{{ formatWeekday(day.date) }}</span>
                <span :class="currentTheme.subtitleClass">{{ formatMonthDay(day.date) }}</span>
              </div>
              <span :class="currentTheme.subtitleClass">{{ day.items.length }} 課程</span>
            </div>

            <div v-if="day.items.length > 0" :class="currentTheme.divideClass">
              <div
                v-for="item in day.items"
                :key="item.id"
                class="flex items-start px-4 py-3 relative"
                :class="[currentTheme.itemClass, currentTheme.itemGradient]"
              >
                <div class="absolute left-0 top-0 bottom-0 w-1" :class="currentTheme.itemAccentClass"></div>
                <!-- 時間區塊 -->
                <div class="w-20 flex-shrink-0 time-block">
                  <div class="font-semibold text-sm whitespace-nowrap" :class="currentTheme.timeClass">{{ item.start_time }}</div>
                  <div class="text-xs" :class="currentTheme.subtitleClass">{{ item.end_time }}</div>
                </div>
                <!-- 課程資訊區塊 -->
                <div class="flex-1 min-w-0 ml-3 flex flex-col justify-center">
                  <div class="flex items-center justify-center gap-2 mb-1 flex-wrap">
                    <span
                      class="w-2 h-2 rounded-full flex-shrink-0"
                      :style="{ backgroundColor: item.color || '#10B981' }"
                    ></span>
                    <h4 class="font-medium list-item-title text-center" :class="currentTheme.itemTitleClass">
                      {{ item.title }}
                    </h4>
                    <span
                      v-if="item.status && item.status !== 'APPROVED' && item.status !== 'NORMAL'"
                      class="px-2 py-0.5 rounded-full text-xs flex-shrink-0"
                      :class="getStatusClass(item.status)"
                    >
                      {{ getStatusText(item.status) }}
                    </span>
                  </div>
                  <div v-if="getValidCenterName(item.center_name)" class="text-xs truncate text-center" :class="currentTheme.centerClass">
                    {{ getValidCenterName(item.center_name) }}
                  </div>
                </div>
                <!-- 時長區塊 -->
                <div class="flex-shrink-0 text-right ml-2 time-block">
                  <p class="text-sm font-medium whitespace-nowrap" :class="currentTheme.timeClass">
                    {{ getDuration(item.start_time, item.end_time) }}
                  </p>
                  <p class="text-xs" :class="currentTheme.subtitleClass">分鐘</p>
                </div>
              </div>
            </div>

            <div v-else class="px-4 py-6 text-center" :class="currentTheme.emptyClass">
              休息日
            </div>
          </div>
        </div>

        <!-- 網格視圖（週課表格式） -->
        <div v-else-if="selectedViewOption === 'grid'" class="overflow-x-auto">
          <div class="min-w-[700px]">
            <!-- 網格標題列 -->
            <div class="grid grid-cols-8 gap-1 mb-1">
              <div class="text-center py-2 text-xs font-medium" :class="currentTheme.subtitleClass">時間</div>
              <div
                v-for="(day, index) in weekDays"
                :key="index"
                class="text-center py-2 rounded-lg"
                :class="currentTheme.dayHeaderClass"
              >
                <div class="text-xs font-medium" :class="currentTheme.dayTextClass">{{ day.weekday }}</div>
                <div class="text-xs" :class="currentTheme.subtitleClass">{{ day.monthDay }}</div>
              </div>
            </div>
            <!-- 網格內容 - 使用 flex 佈局確保每列高度一致 -->
            <div class="grid grid-cols-8 gap-1">
              <!-- 時間列 -->
              <div class="flex flex-col">
                <div
                  v-for="hour in timeSlots"
                  :key="hour"
                  class="flex items-center justify-center text-xs border-r border-white/5"
                  :style="{ height: `${GRID_HOUR_HEIGHT}px` }"
                  :class="currentTheme.subtitleClass"
                >
                  {{ hour }}
                </div>
              </div>
              <!-- 每天的課程列 -->
              <div
                v-for="(day, dayIndex) in weekDays"
                :key="dayIndex"
                class="flex flex-col relative"
              >
                <!-- 背景網格線 -->
                <div
                  v-for="hour in timeSlots"
                  :key="hour"
                  class="border-b border-white/5"
                  :style="{ height: `${GRID_HOUR_HEIGHT}px` }"
                  :class="currentTheme.itemGradient"
                ></div>
                <!-- 課程區塊 - 使用 flex 縱向排列，根據時間計算偏移 -->
                <div class="absolute inset-0 flex flex-col p-1">
                  <div
                    v-for="item in getDayScheduleItems(day.date)"
                    :key="item.id"
                    class="absolute left-1 right-1 rounded p-1.5 text-xs z-10"
                    :style="{
                      backgroundColor: `${item.color || '#10B981'}30`,
                      borderLeft: `3px solid ${item.color || '#10B981'}`,
                      top: getGridItemTopOffset(item),
                      height: `${Math.max(getGridItemHeight(item), 65)}px`
                    }"
                  >
                    <div class="flex items-center gap-1 mb-0.5 leading-none">
                      <span class="font-medium truncate" :class="currentTheme.itemTitleClass">
                        {{ item.start_time }}
                      </span>
                      <span v-if="item.status && item.status !== 'APPROVED' && item.status !== 'NORMAL'" 
                        class="px-1 py-0.5 rounded text-[10px] flex-shrink-0"
                        :class="getStatusClass(item.status)"
                      >
                        {{ getStatusText(item.status) }}
                      </span>
                    </div>
                    <div class="font-medium leading-tight mb-0.5 text-center" :class="[currentTheme.itemTitleClass, getGridTitleClass(item)]">
                      {{ item.title }}
                    </div>
                    <div v-if="getValidCenterName(item.center_name)" class="text-[10px] leading-none truncate opacity-80 text-center" :class="currentTheme.centerClass">
                      {{ getValidCenterName(item.center_name) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-8 pt-6 flex items-center justify-between" :class="currentTheme.borderClass">
          <div class="flex gap-4">
            <div class="flex items-center gap-2">
              <span class="w-3 h-3 rounded-full" :class="currentTheme.statusDotGradient"></span>
              <span class="text-sm" :class="currentTheme.subtitleClass">已確認</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-3 h-3 rounded-full" :class="currentTheme.pendingDotGradient"></span>
              <span class="text-sm" :class="currentTheme.subtitleClass">待審核</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-3 h-3 rounded-full" :class="currentTheme.rejectedDotGradient"></span>
              <span class="text-sm" :class="currentTheme.subtitleClass">已拒絕</span>
            </div>
          </div>
          <p class="text-sm" :class="currentTheme.subtitleClass">TimeLedger 課表管理系統</p>
        </div>
      </div>
    </div>

    <div class="lg:col-span-1 space-y-4">
      <div class="glass-card p-4">
        <h3 class="text-white font-semibold mb-4">選擇風格</h3>
        <div class="grid grid-cols-2 gap-2">
          <button
            v-for="theme in themes"
            :key="theme.id"
            @click="selectedTheme = theme.id"
            class="p-3 rounded-lg text-left transition-all"
            :class="selectedTheme === theme.id ? 'bg-primary-500/20 border-2 border-primary-500' : 'bg-white/5 border-2 border-transparent hover:bg-white/10'"
          >
            <div
              class="w-full h-8 rounded mb-2"
              :style="{ background: theme.preview }"
            ></div>
            <p class="text-xs font-medium" :class="selectedTheme === theme.id ? 'text-primary-400' : 'text-slate-300'">
              {{ theme.name }}
            </p>
          </button>
        </div>
      </div>

      <div class="glass-card p-4">
        <h3 class="text-white font-semibold mb-4">匯出選項</h3>
        <div class="space-y-4">
          <label class="flex items-center gap-3 cursor-pointer">
            <input type="checkbox" v-model="options.showPersonalInfo" class="accent-primary-500 w-4 h-4" />
            <span class="text-slate-300 text-sm">顯示個人資訊</span>
          </label>
          <label class="flex items-center gap-3 cursor-pointer">
            <input type="checkbox" v-model="options.showStats" class="accent-primary-500 w-4 h-4" />
            <span class="text-slate-300 text-sm">顯示統計資料</span>
          </label>
          <label class="flex items-center gap-3 cursor-pointer">
            <input type="checkbox" v-model="options.includeNotes" class="accent-primary-500 w-4 h-4" />
            <span class="text-slate-300 text-sm">包含備註</span>
          </label>
        </div>
      </div>

      <div v-if="options.showStats" class="glass-card p-4">
        <h3 class="text-white font-semibold mb-4">本週統計</h3>
        <div class="space-y-3">
          <div class="flex justify-between">
            <span class="text-slate-400 text-sm">總課程數</span>
            <span class="text-white font-medium">{{ totalLessons }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400 text-sm">總時數</span>
            <span class="text-white font-medium">{{ totalHours }} 小時</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400 text-sm">教學天數</span>
            <span class="text-white font-medium">{{ teachingDays }} 天</span>
          </div>
        </div>
      </div>

      <!-- iCal 訂閱區塊 -->
      <div class="glass-card p-4">
        <h3 class="text-white font-semibold mb-4">iCal 訂閱</h3>

        <!-- 尚未建立訂閱 -->
        <div v-if="!subscriptionUrl" class="space-y-3">
          <p class="text-slate-400 text-sm">建立課表訂閱連結，同步到您的日曆 App</p>
          <button
            @click="handleCreateSubscription"
            class="w-full px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors text-sm flex items-center justify-center gap-2"
            :disabled="scheduleStore.isCreatingSubscription"
          >
            <svg v-if="!scheduleStore.isCreatingSubscription" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            <BaseLoading v-else :loading="true" size="sm" />
            建立訂閱連結
          </button>
        </div>

        <!-- 已建立訂閱 -->
        <div v-else class="space-y-3">
          <div class="flex items-center gap-2 p-2 rounded-lg bg-white/5">
            <div class="flex-1 min-w-0">
              <p class="text-xs text-slate-400 mb-1">訂閱連結</p>
              <p class="text-xs text-slate-300 truncate font-mono">{{ subscriptionUrl }}</p>
            </div>
          </div>

          <div class="flex gap-2">
            <button
              @click="handleCopySubscriptionUrl"
              class="flex-1 px-3 py-2 rounded-lg bg-white/10 text-white hover:bg-white/20 transition-colors text-sm flex items-center justify-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
              </svg>
              複製連結
            </button>
            <button
              @click="handleDeleteSubscription"
              class="px-3 py-2 rounded-lg bg-critical-500/20 text-critical-500 hover:bg-critical-500/30 transition-colors text-sm flex items-center justify-center"
              :disabled="scheduleStore.isDeletingSubscription"
            >
              <svg v-if="!scheduleStore.isDeletingSubscription" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              <BaseLoading v-else :loading="true" size="sm" />
            </button>
          </div>
        </div>
      </div>

      <div class="glass-card p-4">
        <h3 class="text-white font-semibold mb-4">快速操作</h3>
        <div class="space-y-2">
          <button
            @click="router.push('/teacher/dashboard')"
            class="w-full px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors text-left text-sm"
          >
            查看完整課表
          </button>
          <button
            @click="router.push('/teacher/exceptions')"
            class="w-full px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors text-left text-sm"
          >
            例外申請
          </button>
        </div>
      </div>
    </div>
  </div>

  <NavigationNotificationDropdown v-if="notificationUI.show.value" @close="notificationUI.close()" />
  <TeacherSidebar v-if="sidebarStore.isOpen.value" @close="sidebarStore.close()" />
</template>

<script setup lang="ts">
definePageMeta({
  auth: 'TEACHER',
  layout: 'default',
})

import { formatDateToString } from '~/composables/useTaiwanTime'

interface Theme {
  id: string
  name: string
  preview: string
  cardClass: string
  cardGradient: string
  avatarClass: string
  titleClass: string
  subtitleClass: string
  tagClass: string
  tagClass2: string
  tagClass3: string
  dayHeaderClass: string
  dayTextClass: string
  itemClass: string
  itemGradient: string
  itemAccentClass: string
  itemTitleClass: string
  centerClass: string
  personalClass: string
  timeClass: string
  timeGradient: string
  borderClass: string
  divideClass: string
  emptyClass: string
  buttonClass: string
  iconClass: string
  statusDotClass: string
  statusDotGradient: string
  pendingDotClass: string
  pendingDotGradient: string
  rejectedDotClass: string
  rejectedDotGradient: string
}

const themes: Theme[] = [
  {
    id: 'dustyRose',
    name: '玫瑰灰',
    preview: 'linear-gradient(135deg, #c9b1b1 0%, #dcc8c8 50%, #ebdada 100%)',
    cardClass: 'bg-[#f8f5f5]',
    cardGradient: 'bg-gradient-to-br from-[#faf5f5] via-[#f8f2f2] to-[#f5eeee]',
    avatarClass: 'bg-gradient-to-br from-[#c9b1b1] to-[#dfc4c4]',
    titleClass: 'text-[#6b5555]',
    subtitleClass: 'text-[#9a8888]',
    tagClass: 'bg-gradient-to-r from-[#ebe2e2] to-[#e5d8d8] text-[#6b5555]',
    tagClass2: 'bg-gradient-to-r from-[#e5dada] to-[#ded0d0] text-[#6b5555]',
    tagClass3: 'bg-gradient-to-r from-[#dfd2d2] to-[#d8c8c8] text-[#6b5555]',
    dayHeaderClass: 'bg-gradient-to-r from-[#f5f0ee] to-[#f0ebe9] border-b border-[#e5dada]',
    dayTextClass: 'text-[#5a4a4a]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#f5f0ee] hover:to-[#f0ebe9]',
    itemGradient: 'bg-gradient-to-r from-[#faf7f7] to-[#f8f4f4]',
    itemAccentClass: 'bg-gradient-to-b from-[#c9b1b1] to-[#b5a0a0]',
    itemTitleClass: 'text-[#5a4a4a]',
    centerClass: 'text-[#9a8888]',
    personalClass: 'text-[#9a8888]',
    timeClass: 'text-[#6b5555]',
    timeGradient: 'bg-gradient-to-b from-[#6b5555] to-[#5a4545]',
    borderClass: 'border-t border-[#e5dada]',
    divideClass: 'divide-y divide-[#f0ebe9]',
    emptyClass: 'text-[#b5a5a5]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#f0ebe9] hover:to-[#eae4e0]',
    iconClass: 'text-[#9a8888]',
    statusDotClass: 'bg-[#c9b1b1]',
    statusDotGradient: 'bg-gradient-to-br from-[#c9b1b1] to-[#dfc4c4]',
    pendingDotClass: 'bg-[#d4c4c4]',
    pendingDotGradient: 'bg-gradient-to-br from-[#d4c4c4] to-[#e0d4d4]',
    rejectedDotClass: 'bg-[#dfd5d5]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#dfd5d5] to-[#ebe0e0]',
  },
  {
    id: 'sageGreen',
    name: '鼠尾草綠',
    preview: 'linear-gradient(135deg, #a8b598 0%, #bdd4b8 50%, #d2e8d8 100%)',
    cardClass: 'bg-[#f5f7f5]',
    cardGradient: 'bg-gradient-to-br from-[#f7faf7] via-[#f5f7f2] to-[#f2f5ee]',
    avatarClass: 'bg-gradient-to-br from-[#a8b598] to-[#c8d8b8]',
    titleClass: 'text-[#5a6b55]',
    subtitleClass: 'text-[#8a9e88]',
    tagClass: 'bg-gradient-to-r from-[#e5ede5] to-[#dce5dc] text-[#5a6b55]',
    tagClass2: 'bg-gradient-to-r from-[#dce5dc] to-[#d3ddd3] text-[#5a6b55]',
    tagClass3: 'bg-gradient-to-r from-[#d3ddd3] to-[#cad5ca] text-[#5a6b55]',
    dayHeaderClass: 'bg-gradient-to-r from-[#ebf0eb] to-[#e5ede5] border-b border-[#dce5dc]',
    dayTextClass: 'text-[#4a5a45]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#ebf0eb] hover:to-[#e5ede5]',
    itemGradient: 'bg-gradient-to-r from-[#f7faf7] to-[#f4f7f4]',
    itemAccentClass: 'bg-gradient-to-b from-[#a8b598] to-[#95a585]',
    itemTitleClass: 'text-[#4a5a45]',
    centerClass: 'text-[#8a9e88]',
    personalClass: 'text-[#8a9e88]',
    timeClass: 'text-[#5a6b55]',
    timeGradient: 'bg-gradient-to-b from-[#5a6b55] to-[#4a5a45]',
    borderClass: 'border-t border-[#dce5dc]',
    divideClass: 'divide-y divide-[#ebf0eb]',
    emptyClass: 'text-[#95a585]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#ebf0eb] hover:to-[#e5ede5]',
    iconClass: 'text-[#8a9e88]',
    statusDotClass: 'bg-[#a8b598]',
    statusDotGradient: 'bg-gradient-to-br from-[#a8b598] to-[#c8d8b8]',
    pendingDotClass: 'bg-[#b8c4a8]',
    pendingDotGradient: 'bg-gradient-to-br from-[#b8c4a8] to-[#d0dcc0]',
    rejectedDotClass: 'bg-[#c8d4b8]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#c8d4b8] to-[#e0e8d0]',
  },
  {
    id: 'mutedBlue',
    name: '霧霾藍',
    preview: 'linear-gradient(135deg, #9bafba 0%, #b5c8d4 50%, #cfdce8 100%)',
    cardClass: 'bg-[#f5f7f9]',
    cardGradient: 'bg-gradient-to-br from-[#f8fafb] via-[#f5f7f9] to-[#f2f5f7]',
    avatarClass: 'bg-gradient-to-br from-[#9bafba] to-[#bbc8d8]',
    titleClass: 'text-[#565d6b]',
    subtitleClass: 'text-[#7a8a99]',
    tagClass: 'bg-gradient-to-r from-[#e5eaed] to-[#dce2e8] text-[#565d6b]',
    tagClass2: 'bg-gradient-to-r from-[#dce2e8] to-[#d3dbe2] text-[#565d6b]',
    tagClass3: 'bg-gradient-to-r from-[#d3dbe2] to-[#cad4dc] text-[#565d6b]',
    dayHeaderClass: 'bg-gradient-to-r from-[#e9eff2] to-[#e5eaed] border-b border-[#dce2e8]',
    dayTextClass: 'text-[#464d5a]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#e9eff2] hover:to-[#e5eaed]',
    itemGradient: 'bg-gradient-to-r from-[#f8fafb] to-[#f5f8fa]',
    itemAccentClass: 'bg-gradient-to-b from-[#9bafba] to-[#8899a8]',
    itemTitleClass: 'text-[#464d5a]',
    centerClass: 'text-[#7a8a99]',
    personalClass: 'text-[#7a8a99]',
    timeClass: 'text-[#565d6b]',
    timeGradient: 'bg-gradient-to-b from-[#565d6b] to-[#464d5a]',
    borderClass: 'border-t border-[#dce2e8]',
    divideClass: 'divide-y divide-[#e9eff2]',
    emptyClass: 'text-[#8899a8]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#e9eff2] hover:to-[#e5eaed]',
    iconClass: 'text-[#7a8a99]',
    statusDotClass: 'bg-[#9bafba]',
    statusDotGradient: 'bg-gradient-to-br from-[#9bafba] to-[#bbc8d8]',
    pendingDotClass: 'bg-[#a8bcc8]',
    pendingDotGradient: 'bg-gradient-to-br from-[#a8bcc8] to-[#c0d0e0]',
    rejectedDotClass: 'bg-[#b8c9d8]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#b8c9d8] to-[#d0dce8]',
  },
  {
    id: 'warmBeige',
    name: '暖米色',
    preview: 'linear-gradient(135deg, #c4b7a6 0%, #d8c8b0 50%, #ecdac0 100%)',
    cardClass: 'bg-[#faf8f5]',
    cardGradient: 'bg-gradient-to-br from-[#fdfbf7] via-[#faf8f5] to-[#f7f5f0]',
    avatarClass: 'bg-gradient-to-br from-[#c4b7a6] to-[#e0d0b8]',
    titleClass: 'text-[#5c5650]',
    subtitleClass: 'text-[#8a847a]',
    tagClass: 'bg-gradient-to-r from-[#f0ede8] to-[#e9e4dc] text-[#5c5650]',
    tagClass2: 'bg-gradient-to-r from-[#e9e4dc] to-[#e2dbd0] text-[#5c5650]',
    tagClass3: 'bg-gradient-to-r from-[#e2dbd0] to-[#dbd2c4] text-[#5c5650]',
    dayHeaderClass: 'bg-gradient-to-r from-[#f5f1eb] to-[#f0ede8] border-b border-[#e9e4dc]',
    dayTextClass: 'text-[#4a4540]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#f5f1eb] hover:to-[#f0ede8]',
    itemGradient: 'bg-gradient-to-r from-[#fdfbf8] to-[#faf8f5]',
    itemAccentClass: 'bg-gradient-to-b from-[#c4b7a6] to-[#b0a08a]',
    itemTitleClass: 'text-[#4a4540]',
    centerClass: 'text-[#8a847a]',
    personalClass: 'text-[#8a847a]',
    timeClass: 'text-[#5c5650]',
    timeGradient: 'bg-gradient-to-b from-[#5c5650] to-[#4a4540]',
    borderClass: 'border-t border-[#e9e4dc]',
    divideClass: 'divide-y divide-[#f5f1eb]',
    emptyClass: 'text-[#b5aa9a]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#f5f1eb] hover:to-[#f0ede8]',
    iconClass: 'text-[#8a847a]',
    statusDotClass: 'bg-[#c4b7a6]',
    statusDotGradient: 'bg-gradient-to-br from-[#c4b7a6] to-[#e0d0b8]',
    pendingDotClass: 'bg-[#d0c4b0]',
    pendingDotGradient: 'bg-gradient-to-br from-[#d0c4b0] to-[#e0d8c8]',
    rejectedDotClass: 'bg-[#dcd4c0]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#dcd4c0] to-[#eee8d8]',
  },
  {
    id: 'lavender',
    name: '薰衣草灰',
    preview: 'linear-gradient(135deg, #b5aab5 0%, #c8bdc8 50%, #dbd0db 100%)',
    cardClass: 'bg-[#f8f7f8]',
    cardGradient: 'bg-gradient-to-br from-[#faf9fa] via-[#f8f7f8] to-[#f5f5f5]',
    avatarClass: 'bg-gradient-to-br from-[#b5aab5] to-[#d5c8d5]',
    titleClass: 'text-[#5d595f]',
    subtitleClass: 'text-[#8a8088]',
    tagClass: 'bg-gradient-to-r from-[#f0eef0] to-[#e8e5e8] text-[#5d595f]',
    tagClass2: 'bg-gradient-to-r from-[#e8e5e8] to-[#e0dce0] text-[#5d595f]',
    tagClass3: 'bg-gradient-to-r from-[#e0dce0] to-[#d8d3d8] text-[#5d595f]',
    dayHeaderClass: 'bg-gradient-to-r from-[#f2eef2] to-[#f0eef0] border-b border-[#e8e5e8]',
    dayTextClass: 'text-[#4d494f]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#f2eef2] hover:to-[#f0eef0]',
    itemGradient: 'bg-gradient-to-r from-[#faf9fa] to-[#f7f7f7]',
    itemAccentClass: 'bg-gradient-to-b from-[#b5aab5] to-[#9a8a9a]',
    itemTitleClass: 'text-[#4d494f]',
    centerClass: 'text-[#8a8088]',
    personalClass: 'text-[#8a8088]',
    timeClass: 'text-[#5d595f]',
    timeGradient: 'bg-gradient-to-b from-[#5d595f] to-[#4d494f]',
    borderClass: 'border-t border-[#e8e5e8]',
    divideClass: 'divide-y divide-[#f2eef2]',
    emptyClass: 'text-[#a99fa8]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#f2eef2] hover:to-[#f0eef0]',
    iconClass: 'text-[#8a8088]',
    statusDotClass: 'bg-[#b5aab5]',
    statusDotGradient: 'bg-gradient-to-br from-[#b5aab5] to-[#d5c8d5]',
    pendingDotClass: 'bg-[#c4bac4]',
    pendingDotGradient: 'bg-gradient-to-br from-[#c4bac4] to-[#d8d0d8]',
    rejectedDotClass: 'bg-[#d4cad4]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#d4cad4] to-[#e8e0e8]',
  },
  {
    id: 'warmGrey',
    name: '溫柔灰',
    preview: 'linear-gradient(135deg, #b8b8b8 0%, #d0d0d0 50%, #e8e8e8 100%)',
    cardClass: 'bg-[#f8f8f8]',
    cardGradient: 'bg-gradient-to-br from-[#fafafa] via-[#f8f8f8] to-[#f5f5f5]',
    avatarClass: 'bg-gradient-to-br from-[#b8b8b8] to-[#d8d8d8]',
    titleClass: 'text-[#555555]',
    subtitleClass: 'text-[#888888]',
    tagClass: 'bg-gradient-to-r from-[#f0f0f0] to-[#e8e8e8] text-[#555555]',
    tagClass2: 'bg-gradient-to-r from-[#e8e8e8] to-[#e0e0e0] text-[#555555]',
    tagClass3: 'bg-gradient-to-r from-[#e0e0e0] to-[#d8d8d8] text-[#555555]',
    dayHeaderClass: 'bg-gradient-to-r from-[#f2f2f2] to-[#f0f0f0] border-b border-[#e8e8e8]',
    dayTextClass: 'text-[#454545]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#f2f2f2] hover:to-[#f0f0f0]',
    itemGradient: 'bg-gradient-to-r from-[#fafafa] to-[#f7f7f7]',
    itemAccentClass: 'bg-gradient-to-b from-[#b8b8b8] to-[#a0a0a0]',
    itemTitleClass: 'text-[#454545]',
    centerClass: 'text-[#888888]',
    personalClass: 'text-[#888888]',
    timeClass: 'text-[#555555]',
    timeGradient: 'bg-gradient-to-b from-[#555555] to-[#454545]',
    borderClass: 'border-t border-[#e8e8e8]',
    divideClass: 'divide-y divide-[#f2f2f2]',
    emptyClass: 'text-[#a8a8a8]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#f2f2f2] hover:to-[#f0f0f0]',
    iconClass: 'text-[#888888]',
    statusDotClass: 'bg-[#b8b8b8]',
    statusDotGradient: 'bg-gradient-to-br from-[#b8b8b8] to-[#d8d8d8]',
    pendingDotClass: 'bg-[#c8c8c8]',
    pendingDotGradient: 'bg-gradient-to-br from-[#c8c8c8] to-[#e0e0e0]',
    rejectedDotClass: 'bg-[#d8d8d8]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#d8d8d8] to-[#f0f0f0]',
  },
  {
    id: 'domeMouse',
    name: '多梅鼠',
    preview: 'linear-gradient(135deg, #b5a9a0 0%, #c8bbb2 50%, #dbcdc4 100%)',
    cardClass: 'bg-[#f8f6f5]',
    cardGradient: 'bg-gradient-to-br from-[#faf8f6] via-[#f8f6f4] to-[#f5f2ef]',
    avatarClass: 'bg-gradient-to-br from-[#b5a9a0] to-[#d4c8bc]',
    titleClass: 'text-[#5a524d]',
    subtitleClass: 'text-[#8a7e75]',
    tagClass: 'bg-gradient-to-r from-[#ebe6e2] to-[#e3ddd8] text-[#5a524d]',
    tagClass2: 'bg-gradient-to-r from-[#e3ddd8] to-[#dbd4ce] text-[#5a524d]',
    tagClass3: 'bg-gradient-to-r from-[#dbd4ce] to-[#d3cbc4] text-[#5a524d]',
    dayHeaderClass: 'bg-gradient-to-r from-[#f5f1ed] to-[#f0ece8] border-b border-[#e3ddd8]',
    dayTextClass: 'text-[#4a4440]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#f5f1ed] hover:to-[#f0ece8]',
    itemGradient: 'bg-gradient-to-r from-[#faf8f6] to-[#f7f4f1]',
    itemAccentClass: 'bg-gradient-to-b from-[#b5a9a0] to-[#9a8e82]',
    itemTitleClass: 'text-[#4a4440]',
    centerClass: 'text-[#8a7e75]',
    personalClass: 'text-[#8a7e75]',
    timeClass: 'text-[#5a524d]',
    timeGradient: 'bg-gradient-to-b from-[#5a524d] to-[#4a4440]',
    borderClass: 'border-t border-[#e3ddd8]',
    divideClass: 'divide-y divide-[#f5f1ed]',
    emptyClass: 'text-[#a89e94]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#f5f1ed] hover:to-[#f0ece8]',
    iconClass: 'text-[#8a7e75]',
    statusDotClass: 'bg-[#b5a9a0]',
    statusDotGradient: 'bg-gradient-to-br from-[#b5a9a0] to-[#d4c8bc]',
    pendingDotClass: 'bg-[#c4b8b0]',
    pendingDotGradient: 'bg-gradient-to-br from-[#c4b8b0] to-[#d8d0c8]',
    rejectedDotClass: 'bg-[#d4ccc4]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#d4ccc4] to-[#e8e0d8]',
  },
  {
    id: 'deepDomeMouse',
    name: '深梅鼠',
    preview: 'linear-gradient(135deg, #9a8a7a 0%, #b5a598 50%, #d0c0b0 100%)',
    cardClass: 'bg-[#f5f2ef]',
    cardGradient: 'bg-gradient-to-br from-[#f7f4f1] via-[#f5f2ed] to-[#f2efe9]',
    avatarClass: 'bg-gradient-to-br from-[#9a8a7a] to-[#c4b4a4]',
    titleClass: 'text-[#4a4038]',
    subtitleClass: 'text-[#7a6a5a]',
    tagClass: 'bg-gradient-to-r from-[#e5ded8] to-[#dcd4cc] text-[#4a4038]',
    tagClass2: 'bg-gradient-to-r from-[#dcd4cc] to-[#d3cac0] text-[#4a4038]',
    tagClass3: 'bg-gradient-to-r from-[#d3cac0] to-[#cac0b4] text-[#4a4038]',
    dayHeaderClass: 'bg-gradient-to-r from-[#f0ede8] to-[#ebe6e0] border-b border-[#dcd4cc]',
    dayTextClass: 'text-[#3a342a]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#f0ede8] hover:to-[#ebe6e0]',
    itemGradient: 'bg-gradient-to-r from-[#f7f4f1] to-[#f4f0ec]',
    itemAccentClass: 'bg-gradient-to-b from-[#9a8a7a] to-[#7a6a5a]',
    itemTitleClass: 'text-[#3a342a]',
    centerClass: 'text-[#7a6a5a]',
    personalClass: 'text-[#7a6a5a]',
    timeClass: 'text-[#4a4038]',
    timeGradient: 'bg-gradient-to-b from-[#4a4038] to-[#3a342a]',
    borderClass: 'border-t border-[#dcd4cc]',
    divideClass: 'divide-y divide-[#f0ede8]',
    emptyClass: 'text-[#8a7a6a]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#f0ede8] hover:to-[#ebe6e0]',
    iconClass: 'text-[#7a6a5a]',
    statusDotClass: 'bg-[#9a8a7a]',
    statusDotGradient: 'bg-gradient-to-br from-[#9a8a7a] to-[#c4b4a4]',
    pendingDotClass: 'bg-[#a99a8a]',
    pendingDotGradient: 'bg-gradient-to-br from-[#a99a8a] to-[#c8bcb0]',
    rejectedDotClass: 'bg-[#c0b4a4]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#c0b4a4] to-[#e0d8cc]',
  },
  {
    id: 'lightDomeMouse',
    name: '淡梅鼠',
    preview: 'linear-gradient(135deg, #d4ccc4 0%, #e0d8d0 50%, #ece8e0 100%)',
    cardClass: 'bg-[#faf9f7]',
    cardGradient: 'bg-gradient-to-br from-[#fdfcfb] via-[#faf9f7] to-[#f7f5f2]',
    avatarClass: 'bg-gradient-to-br from-[#d4ccc4] to-[#ede8e0]',
    titleClass: 'text-[#6a6058]',
    subtitleClass: 'text-[#9a8e82]',
    tagClass: 'bg-gradient-to-r from-[#f5f2ee] to-[#f0ede8] text-[#6a6058]',
    tagClass2: 'bg-gradient-to-r from-[#f0ede8] to-[#ebe6df] text-[#6a6058]',
    tagClass3: 'bg-gradient-to-r from-[#ebe6df] to-[#e6dfd6] text-[#6a6058]',
    dayHeaderClass: 'bg-gradient-to-r from-[#f7f4f0] to-[#f4f0ec] border-b border-[#f0ede8]',
    dayTextClass: 'text-[#5a5248]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#f7f4f0] hover:to-[#f4f0ec]',
    itemGradient: 'bg-gradient-to-r from-[#fdfcfb] to-[#faf9f7]',
    itemAccentClass: 'bg-gradient-to-b from-[#d4ccc4] to-[#b4aca4]',
    itemTitleClass: 'text-[#5a5248]',
    centerClass: 'text-[#9a8e82]',
    personalClass: 'text-[#9a8e82]',
    timeClass: 'text-[#6a6058]',
    timeGradient: 'bg-gradient-to-b from-[#6a6058] to-[#5a5248]',
    borderClass: 'border-t border-[#f0ede8]',
    divideClass: 'divide-y divide-[#f7f4f0]',
    emptyClass: 'text-[#c4bcb0]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#f7f4f0] hover:to-[#f4f0ec]',
    iconClass: 'text-[#9a8e82]',
    statusDotClass: 'bg-[#d4ccc4]',
    statusDotGradient: 'bg-gradient-to-br from-[#d4ccc4] to-[#ede8e0]',
    pendingDotClass: 'bg-[#e0dcd4]',
    pendingDotGradient: 'bg-gradient-to-br from-[#e0dcd4] to-[#f0ede8]',
    rejectedDotClass: 'bg-[#ece8e0]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#ece8e0] to-[#fdfbf9]',
  },
  {
    id: 'coralRose',
    name: '珊瑚玫瑰',
    preview: 'linear-gradient(135deg, #9e7a7a 0%, #b8958f 50%, #d2b0a8 100%)',
    cardClass: 'bg-[#f8f5f4]',
    cardGradient: 'bg-gradient-to-br from-[#faf7f6] via-[#f8f5f3] to-[#f5f1ef]',
    avatarClass: 'bg-gradient-to-br from-[#9e7a7a] to-[#c4a8a0]',
    titleClass: 'text-[#5a4040]',
    subtitleClass: 'text-[#8a6868]',
    tagClass: 'bg-gradient-to-r from-[#ede6e6] to-[#e5d8d8] text-[#5a4040]',
    tagClass2: 'bg-gradient-to-r from-[#e5d8d8] to-[#ddcaca] text-[#5a4040]',
    tagClass3: 'bg-gradient-to-r from-[#ddcaca] to-[#d5bcbc] text-[#5a4040]',
    dayHeaderClass: 'bg-gradient-to-r from-[#f5f0ef] to-[#f0eae8] border-b border-[#e5d8d8]',
    dayTextClass: 'text-[#4a3535]',
    itemClass: 'hover:bg-gradient-to-r hover:from-[#f5f0ef] hover:to-[#f0eae8]',
    itemGradient: 'bg-gradient-to-r from-[#faf7f6] to-[#f7f3f1]',
    itemAccentClass: 'bg-gradient-to-b from-[#9e7a7a] to-[#7e5a5a]',
    itemTitleClass: 'text-[#4a3535]',
    centerClass: 'text-[#8a6868]',
    personalClass: 'text-[#8a6868]',
    timeClass: 'text-[#5a4040]',
    timeGradient: 'bg-gradient-to-b from-[#5a4040] to-[#4a3535]',
    borderClass: 'border-t border-[#e5d8d8]',
    divideClass: 'divide-y divide-[#f5f0ef]',
    emptyClass: 'text-[#9e7a7a]',
    buttonClass: 'hover:bg-gradient-to-r hover:from-[#f5f0ef] hover:to-[#f0eae8]',
    iconClass: 'text-[#8a6868]',
    statusDotClass: 'bg-[#9e7a7a]',
    statusDotGradient: 'bg-gradient-to-br from-[#9e7a7a] to-[#c4a8a0]',
    pendingDotClass: 'bg-[#b08a8a]',
    pendingDotGradient: 'bg-gradient-to-br from-[#b08a8a] to-[#d0b8b0]',
    rejectedDotClass: 'bg-[#d0a8a0]',
    rejectedDotGradient: 'bg-gradient-to-br from-[#d0a8a0] to-[#e8d0c8]',
  },
]

const router = useRouter()
const authStore = useAuthStore()
const scheduleStore = useScheduleStore()
const sidebarStore = useSidebar()
const notificationUI = useNotification()
const scheduleRef = ref<HTMLElement>()

const selectedTheme = ref('domeMouse')
const selectedViewOption = ref('list') // 'list' | 'grid'
const options = reactive({
  showPersonalInfo: true,
  showStats: true,
  includeNotes: false,
})

// 視圖選項
const viewOptions = [
  { id: 'list', name: '列表' },
  { id: 'grid', name: '網格' },
]

// 網格視圖所需的計算屬性
const weekDays = computed(() => {
  const days = []
  const weekdays = ['週日', '週一', '週二', '週三', '週四', '週五', '週六']
  scheduleDays.value.forEach(day => {
    const date = new Date(day.date)
    days.push({
      date: day.date,
      weekday: weekdays[date.getDay()],
      monthDay: date.toLocaleDateString('zh-TW', { month: 'numeric', day: 'numeric' })
    })
  })
  return days
})

// 時間槽（每小時）- 從 00:00 到 23:00，涵蓋全天 24 小時
const timeSlots = computed(() => {
  const slots = []
  for (let i = 0; i < 24; i++) {
    slots.push(`${i.toString().padStart(2, '0')}:00`)
  }
  return slots
})

// 網格每小時的高度（像素）- 增加高度以提供更多垂直空間
const GRID_HOUR_HEIGHT = 80

// 取得指定日期和時間的課程項目
const getScheduleItemsForHour = (date: string, hour: string) => {
  const day = scheduleDays.value.find(d => d.date === date)
  if (!day) return []
  
  return day.items.filter(item => {
    return item.start_time.startsWith(hour)
  })
}

// 取得指定日期的所有課程（用於網格視圖）
const getDayScheduleItems = (date: string) => {
  const day = scheduleDays.value.find(d => d.date === date)
  if (!day) return []
  return day.items
}

// 計算網格視圖中項目的頂部偏移（使用可配置的高度）
const getGridItemTopOffset = (item: any) => {
  const [hours, minutes] = item.start_time.split(':').map(Number)
  const startHour = 0 // 開始時間改為 00:00
  const offsetHours = hours - startHour
  const offsetMinutes = minutes
  return `${(offsetHours * GRID_HOUR_HEIGHT + offsetMinutes * GRID_HOUR_HEIGHT / 60)}px`
}

// 計算網格視圖中項目的高度
const getGridItemHeight = (item: any) => {
  const [startH, startM] = item.start_time.split(':').map(Number)
  const [endH, endM] = item.end_time.split(':').map(Number)
  const durationMinutes = (endH * 60 + endM) - (startH * 60 + startM)
  return durationMinutes * GRID_HOUR_HEIGHT / 60 // 使用可配置的高度
}

// 計算項目在格子中的頂部偏移（以小時為單位，1小時 = 60px）
const getItemTopOffset = (item: any) => {
  const [hours, minutes] = item.start_time.split(':').map(Number)
  const startHour = 8 // 開始時間
  const offsetHours = hours - startHour
  const offsetMinutes = minutes
  return `${(offsetHours * 60 + offsetMinutes) * 1}px` // 1px per minute
}

// 計算項目高度
const getItemHeight = (item: any) => {
  const [startH, startM] = item.start_time.split(':').map(Number)
  const [endH, endM] = item.end_time.split(':').map(Number)
  const durationMinutes = (endH * 60 + endM) - (startH * 60 + startM)
  return `${durationMinutes * 1}px` // 1px per minute
}

const currentTheme = computed(() => {
  return themes.find(t => t.id === selectedTheme.value) || themes[0]
})

const weekLabel = computed(() => scheduleStore.weekLabel)

const scheduleDays = computed(() => {
  return scheduleStore.schedule?.days || []
})

const totalLessons = computed(() => {
  return scheduleDays.value.reduce((sum, day) => sum + day.items.length, 0)
})

const totalHours = computed(() => {
  let total = 0
  scheduleDays.value.forEach(day => {
    day.items.forEach(item => {
      total += getDuration(item.start_time, item.end_time)
    })
  })
  return Math.round(total / 60 * 10) / 10
})

const teachingDays = computed(() => {
  return scheduleDays.value.filter(day => day.items.length > 0).length
})

const changeWeek = (delta: number) => {
  scheduleStore.changeWeek(delta)
  scheduleStore.fetchSchedule()
}

const formatDate = (date: Date): string => {
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const formatWeekday = (dateStr: string): string => {
  const date = new Date(dateStr)
  const weekdays = ['週日', '週一', '週二', '週三', '週四', '週五', '週六']
  return weekdays[date.getDay()]
}

const formatMonthDay = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    month: 'long',
    day: 'numeric',
  })
}

const getDuration = (start: string, end: string): number => {
  const [startH, startM] = start.split(':').map(Number)
  const [endH, endM] = end.split(':').map(Number)
  return (endH * 60 + endM) - (startH * 60 + startM)
}

// 檢查中心名稱是否有效（過濾掉 NORMAL、空值等無意義的內容）
const isValidCenterName = (name?: string): boolean => {
  if (!name || typeof name !== 'string') {
    return false
  }
  const trimmed = name.trim()
  // 過濾掉 "NORMAL"、空字串、僅數字等無意義值
  if (!trimmed) return false
  // 檢查是否為 NORMAL（忽略大小寫和空白）
  const normalized = trimmed.toUpperCase().replace(/\s+/g, '')
  if (normalized === 'NORMAL') {
    return false
  }
  // 檢查是否只包含數字
  if (/^\d+$/.test(trimmed)) return false
  return true
}

// 取得有效的中心名稱
const getValidCenterName = (name?: string): string => {
  return isValidCenterName(name) ? name!.trim() : ''
}

const getStatusClass = (status: string) => {
  switch (status) {
    case 'PENDING':
      return 'bg-warning-500/20 text-warning-500'
    case 'APPROVED':
      return 'bg-success-500/20 text-success-500'
    case 'REJECTED':
      return 'bg-critical-500/20 text-critical-500'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'PENDING':
      return '待審核'
    case 'APPROVED':
      return '已核准'
    case 'REJECTED':
      return '已拒絕'
    default:
      return status
  }
}

// 取得老師的個人標籤
const teacherHashtags = computed(() => {
  const user = authStore.user as any
  if (user?.personal_hashtags && Array.isArray(user.personal_hashtags)) {
    return user.personal_hashtags
  }
  return []
})

// 根據索引取得標籤樣式
const getTagClass = (index: number) => {
  const classes = [
    currentTheme.value.tagClass,
    currentTheme.value.tagClass2,
    currentTheme.value.tagClass3,
  ]
  return classes[index % classes.length]
}

// 根據標題長度動態調整字體大小，避免遮蓋
const getGridTitleClass = (item: any) => {
  if (item.title?.length > 10) return 'text-[10px]'
  return 'text-xs'
}

const getBackgroundColor = () => {
  const theme = currentTheme.value
  if (theme.id === 'dustyRose') return '#f8f5f5'
  if (theme.id === 'sageGreen') return '#f5f7f5'
  if (theme.id === 'mutedBlue') return '#f5f7f9'
  if (theme.id === 'warmBeige') return '#faf8f5'
  if (theme.id === 'lavender') return '#f8f7f8'
  if (theme.id === 'warmGrey') return '#f8f8f8'
  if (theme.id === 'domeMouse') return '#f8f6f5'
  if (theme.id === 'deepDomeMouse') return '#f5f2ef'
  if (theme.id === 'lightDomeMouse') return '#faf9f7'
  if (theme.id === 'coralRose') return '#f8f5f4'
  return '#f8f8f8'
}

// 創建乾淨的 DOM 元素用於匯出
const createCleanExportElement = (view: 'list' | 'grid' = 'list'): HTMLElement | null => {
  // 使用簡單的白色背景
  const bgColor = '#ffffff'

  // 創建容器
  const container = document.createElement('div')
  container.id = 'export-container'
  container.style.cssText = `background-color: ${bgColor}; padding: 24px; font-family: Arial, sans-serif; width: ${view === 'grid' ? '900px' : '650px'};`

  // 添加標題區域
  const header = document.createElement('div')
  header.style.cssText = 'display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; padding-bottom: 15px; border-bottom: 2px solid #eee;'
  header.innerHTML = `
    <div style="display: flex; align-items: center; gap: 12px;">
      <div style="width: 50px; height: 50px; border-radius: 50%; background: #6366F1; display: flex; align-items: center; justify-content: center; color: white; font-size: 20px; font-weight: bold;">${authStore.user?.name?.charAt(0) || 'T'}</div>
      <div>
        <h1 style="margin: 0; font-size: 20px; color: #333;">${authStore.user?.name}</h1>
        <p style="margin: 2px 0 0; color: #888; font-size: 12px;">個人課表</p>
      </div>
    </div>
    <div style="text-align: right;">
      <p style="margin: 0; color: #888; font-size: 12px;">匯出日期</p>
      <p style="margin: 2px 0 0; font-size: 14px; color: #333;">${formatDate(new Date())}</p>
    </div>
  `
  container.appendChild(header)

  // 添加週標籤
  const weekLabelDiv = document.createElement('div')
  weekLabelDiv.style.cssText = 'font-size: 16px; font-weight: bold; margin-bottom: 16px; color: #333; padding: 8px 12px; background: #f5f5f5; border-radius: 6px; display: inline-block;'
  weekLabelDiv.textContent = weekLabel.value
  container.appendChild(weekLabelDiv)

  if (view === 'grid') {
    // ========== 網格視圖（週曆格式）- 使用 60px/小時 ==========
    const gridContainer = document.createElement('div')
    gridContainer.style.cssText = 'display: grid; grid-template-columns: repeat(8, 1fr); gap: 2px;'

    // 左上角空白格
    const emptyCell = document.createElement('div')
    emptyCell.style.cssText = 'border: 1px solid #ddd; background: #f0f0f0; height: 40px;'
    gridContainer.appendChild(emptyCell)

    // 星期標題列
    const weekdays = ['日', '一', '二', '三', '四', '五', '六']
    const weekDaysList = scheduleDays.value.slice(0, 7)

    weekDaysList.forEach((day) => {
      const date = new Date(day.date)
      const monthDay = date.toLocaleDateString('zh-TW', { month: 'numeric', day: 'numeric' })

      const dayHeader = document.createElement('div')
      dayHeader.style.cssText = 'padding: 4px; text-align: center; border: 1px solid #ddd; background: #f8f8f8; height: 40px; display: flex; flex-direction: column; justify-content: center;'
      dayHeader.innerHTML = `<span style="font-weight: bold; font-size: 11px; color: #333;">${weekdays[date.getDay()]}</span><span style="font-size: 9px; color: #888;">${monthDay}</span>`
      gridContainer.appendChild(dayHeader)
    })

    // 時間列（左側）- 60px/小時
    const timeColumn = document.createElement('div')
    timeColumn.style.cssText = 'display: flex; flex-direction: column;'
    for (let i = 0; i < 24; i++) {
      const hour = `${i.toString().padStart(2, '0')}:00`
      const timeCell = document.createElement('div')
      timeCell.style.cssText = `height: 60px; display: flex; align-items: flex-start; justify-content: center; font-size: 10px; color: #888; border: 1px solid #ddd; background: #fafafa; box-sizing: border-box; padding-top: 2px;`
      timeCell.textContent = hour
      timeColumn.appendChild(timeCell)
    }
    gridContainer.appendChild(timeColumn)

    // 每天的課程列 - 總高度 1440px (24 * 60)
    weekDaysList.forEach((day) => {
      const dayColumn = document.createElement('div')
      dayColumn.style.cssText = 'display: flex; flex-direction: column; position: relative; height: 1440px; border: 1px solid #ddd; background: white; margin: 0; overflow: hidden;'

      // 背景網格線 - 在最下層
      const gridBg = document.createElement('div')
      gridBg.style.cssText = 'position: absolute; top: 0; left: 0; right: 0; bottom: 0; z-index: 1;'
      for (let i = 0; i < 24; i++) {
        const line = document.createElement('div')
        line.style.cssText = 'height: 60px; border-bottom: 1px solid #eee; box-sizing: border-box;'
        gridBg.appendChild(line)
      }
      dayColumn.appendChild(gridBg)

      // 課程區塊 - 在上層，z-index 更高
      const itemsContainer = document.createElement('div')
      itemsContainer.style.cssText = 'position: absolute; top: 0; left: 0; right: 0; bottom: 0; padding: 2px; z-index: 10;'

      day.items.forEach(item => {
        const [startH, startM] = item.start_time.split(':').map(Number)
        const [endH, endM] = item.end_time.split(':').map(Number)
        const duration = (endH * 60 + endM) - (startH * 60 + startM)
        const startMinutes = startH * 60 + startM
        const top = startMinutes // 每分鐘 1px
        const height = Math.max(duration, 30) // 最小 30px

        // 取得有效的中心名稱
        const centerName = getValidCenterName(item.center_name)
        const displayText = centerName || item.title

        const itemDiv = document.createElement('div')
        itemDiv.style.cssText = `
          position: absolute;
          left: 2px;
          right: 2px;
          top: ${top}px;
          height: ${height}px;
          background: ${item.color || '#10B981'}60;
          border-left: 4px solid ${item.color || '#10B981'};
          border-radius: 4px;
          display: flex;
          align-items: center;
          justify-content: center;
          box-sizing: border-box;
          padding: 0 6px;
          z-index: 20;
          overflow: visible;
        `
        const titleSpan = document.createElement('span')
        titleSpan.style.cssText = `
          font-size: 11px;
          color: #8a7e75;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          max-width: 100%;
          text-align: center;
          line-height: 1.2;
          font-weight: 500;
        `
        titleSpan.textContent = displayText
        itemDiv.appendChild(titleSpan)
        itemsContainer.appendChild(itemDiv)
      })

      dayColumn.appendChild(itemsContainer)
      gridContainer.appendChild(dayColumn)
    })

    container.appendChild(gridContainer)
  } else {
    // ========== 列表視圖 ==========
    scheduleDays.value.forEach(day => {
      const dayDiv = document.createElement('div')
      dayDiv.style.cssText = 'margin-bottom: 16px; border: 1px solid #ddd; border-radius: 8px; overflow: hidden; background: white;'

      const date = new Date(day.date)
      const weekday = ['週日', '週一', '週二', '週三', '週四', '週五', '週六'][date.getDay()]
      const monthDay = date.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric' })

      // 標題
      const dayHeader = document.createElement('div')
      dayHeader.style.cssText = 'display: flex; justify-content: space-between; align-items: center; padding: 10px 14px; background: #f8f8f8; border-bottom: 1px solid #eee;'
      dayHeader.innerHTML = `<span style="font-weight: bold; color: #333; font-size: 14px;">${weekday} ${monthDay}</span><span style="color: #888; font-size: 12px;">${day.items.length} 課程</span>`
      dayDiv.appendChild(dayHeader)

      // 課程列表
      if (day.items.length > 0) {
        const itemsDiv = document.createElement('div')
        itemsDiv.style.cssText = 'padding: 6px;'

        day.items.forEach(item => {
          const itemDiv = document.createElement('div')
          itemDiv.style.cssText = 'display: flex; align-items: center; padding: 10px 12px; margin-bottom: 4px; background: #fafafa; border-radius: 6px; border-left: 3px solid ' + (item.color || '#10B981') + ';'
          itemDiv.innerHTML = `
            <div style="width: 65px; flex-shrink: 0; text-align: center;">
              <div style="font-weight: bold; font-size: 13px; color: #333;">${item.start_time}</div>
              <div style="font-size: 11px; color: #999;">${item.end_time}</div>
            </div>
            <div style="flex: 1; min-width: 0; margin-left: 8px;">
              <div style="font-weight: bold; font-size: 14px; color: #333;">${item.title}</div>
              ${getValidCenterName(item.center_name) ? `<div style="font-size: 12px; color: #888; margin-top: 2px;">${getValidCenterName(item.center_name)}</div>` : ''}
            </div>
            <div style="width: 55px; flex-shrink: 0; text-align: right;">
              <div style="font-weight: bold; font-size: 13px; color: #333;">${getDuration(item.start_time, item.end_time)}</div>
              <div style="font-size: 11px; color: #999;">分鐘</div>
            </div>
          `
          itemsDiv.appendChild(itemDiv)
        })
        dayDiv.appendChild(itemsDiv)
      } else {
        const emptyDiv = document.createElement('div')
        emptyDiv.style.cssText = 'padding: 20px; text-align: center; color: #aaa; font-size: 13px; background: #fafafa;'
        emptyDiv.textContent = '休息日'
        dayDiv.appendChild(emptyDiv)
      }

      container.appendChild(dayDiv)
    })
  }

  // 添加頁腳
  const footer = document.createElement('div')
  footer.style.cssText = 'margin-top: 16px; padding-top: 12px; border-top: 1px solid #eee; text-align: center; color: #999; font-size: 11px;'
  footer.textContent = 'TimeLedger 課表管理系統'
  container.appendChild(footer)

  return container
}

// 清理匯出元素
const cleanupExportElement = () => {
  const element = document.getElementById('export-container')
  if (element && element.parentNode) {
    element.parentNode.removeChild(element)
  }
}

const handleExportAsImage = async () => {
  // 直接截圖週曆區域
  const exportElement = scheduleRef.value
  if (!exportElement) {
    notificationUI.showError('無法找到課表區域')
    return
  }

  // 暫時移除 overflow 限制以確保完整截圖
  const originalOverflowX = exportElement.style.overflowX
  const originalOverflowY = exportElement.style.overflowY
  const originalOverflow = exportElement.style.overflow
  const originalMaxHeight = exportElement.style.maxHeight
  const originalPosition = exportElement.style.position

  // 確保捕獲完整內容
  exportElement.style.overflow = 'visible'
  exportElement.style.maxHeight = 'none'
  exportElement.style.position = 'relative'

  // 添加匯出用的 CSS 類別來隱藏時間和中心名稱
  exportElement.classList.add('export-mode')

  // 等待瀏覽器重新渲染
  await new Promise(resolve => setTimeout(resolve, 100))

  try {
    const { default: html2canvas } = await import('html2canvas')

    // 獲取完整的滾動區域尺寸
    const scrollWidth = exportElement.scrollWidth
    const scrollHeight = exportElement.scrollHeight

    const canvas = await html2canvas(exportElement, {
      backgroundColor: null,
      scale: 2,
      useCORS: true,
      logging: false,
      allowTaint: true,
      width: scrollWidth,
      height: scrollHeight,
      scrollX: 0,
      scrollY: 0,
      windowWidth: scrollWidth,
      windowHeight: scrollHeight,
    })

    // 產生下載
    const link = document.createElement('a')
    link.download = `課表-${weekLabel.value}-${selectedViewOption.value === 'grid' ? '網格' : '列表'}-${Date.now()}.png`
    link.href = canvas.toDataURL('image/png')
    link.click()
  } catch (error) {
    console.error('Image export failed:', error)
    notificationUI.showError('圖片匯出失敗，請稍後再試')
  } finally {
    // 移除匯出用的 CSS 類別
    exportElement.classList.remove('export-mode')
    // 恢復原始樣式
    exportElement.style.overflowX = originalOverflowX
    exportElement.style.overflowY = originalOverflowY
    exportElement.style.overflow = originalOverflow
    exportElement.style.maxHeight = originalMaxHeight
    exportElement.style.position = originalPosition
  }
}

const handleDownloadPDF = async () => {
  notificationUI.showLoading('正在生成 PDF...')

  // 直接截圖週曆區域
  const exportElement = scheduleRef.value
  if (!exportElement) {
    notificationUI.hideLoading()
    notificationUI.showError('無法找到課表區域')
    return
  }

  // 暫時移除 overflow 限制以確保完整截圖
  const originalOverflow = exportElement.style.overflow
  const originalMaxHeight = exportElement.style.maxHeight
  const originalPosition = exportElement.style.position

  exportElement.style.overflow = 'visible'
  exportElement.style.maxHeight = 'none'
  exportElement.style.position = 'relative'

  // 添加匯出用的 CSS 類別來隱藏時間和中心名稱
  exportElement.classList.add('export-mode')

  // 等待瀏覽器重新渲染
  await new Promise(resolve => setTimeout(resolve, 100))

  try {
    const { default: html2canvas } = await import('html2canvas')
    const { default: jsPDF } = await import('jspdf')

    // 獲取完整的滾動區域尺寸
    const scrollWidth = exportElement.scrollWidth
    const scrollHeight = exportElement.scrollHeight

    const canvas = await html2canvas(exportElement, {
      backgroundColor: null,
      scale: 2,
      useCORS: true,
      logging: false,
      allowTaint: true,
      width: scrollWidth,
      height: scrollHeight,
      scrollX: 0,
      scrollY: 0,
      windowWidth: scrollWidth,
      windowHeight: scrollHeight,
    })

    const imgData = canvas.toDataURL('image/png', 1.0)
    const pdf = new jsPDF({
      orientation: scrollWidth > scrollHeight ? 'landscape' : 'portrait',
      unit: 'mm',
      format: 'a4',
    })

    const pdfWidth = pdf.internal.pageSize.getWidth()
    const pdfHeight = pdf.internal.pageSize.getHeight()
    const imgWidth = canvas.width
    const imgHeight = canvas.height
    const ratio = Math.min(pdfWidth / imgWidth, pdfHeight / imgHeight)

    const imgX = (pdfWidth - imgWidth * ratio) / 2
    const imgY = 10

    // 如果內容超過一頁，分頁處理
    const contentHeight = imgHeight * ratio
    if (contentHeight > pdfHeight - 20) {
      const pageCount = Math.ceil(contentHeight / (pdfHeight - 20))
      const pageHeight = (pdfHeight - 20) / ratio

      for (let i = 0; i < pageCount; i++) {
        if (i > 0) {
          pdf.addPage()
        }

        const sourceY = i * pageHeight
        const sourceHeight = Math.min(pageHeight, imgHeight - sourceY)

        const tempCanvas = document.createElement('canvas')
        tempCanvas.width = imgWidth
        tempCanvas.height = sourceHeight
        const tempCtx = tempCanvas.getContext('2d')

        if (tempCtx) {
          tempCtx.fillStyle = '#ffffff'
          tempCtx.fillRect(0, 0, imgWidth, sourceHeight)
          tempCtx.drawImage(canvas, 0, sourceY, imgWidth, sourceHeight, 0, 0, imgWidth, sourceHeight)

          const tempImgData = tempCanvas.toDataURL('image/png', 1.0)
          const yPos = i === 0 ? imgY : 10
          pdf.addImage(tempImgData, 'PNG', imgX, yPos, imgWidth * ratio, sourceHeight * ratio)
        }
      }
    } else {
      pdf.addImage(imgData, 'PNG', imgX, imgY, imgWidth * ratio, imgHeight * ratio)
    }

    pdf.save(`課表-${weekLabel.value}-${selectedViewOption.value === 'grid' ? '網格' : '列表'}-${Date.now()}.pdf`)
    notificationUI.hideLoading()
    notificationUI.showSuccess('PDF 已下載')
  } catch (error) {
    console.error('PDF generation failed:', error)
    notificationUI.hideLoading()
    notificationUI.showError('PDF 生成失敗，請稍後再試')
  } finally {
    // 移除匯出用的 CSS 類別
    exportElement.classList.remove('export-mode')
    // 恢復原始樣式
    exportElement.style.overflow = originalOverflow
    exportElement.style.maxHeight = originalMaxHeight
    exportElement.style.position = originalPosition
  }
}

// 產生 iCal 格式的課表資料
const generateICalData = (): string => {
  const now = new Date()
  const events: string[] = []

  // iCal 檔案頭
  events.push('BEGIN:VCALENDAR')
  events.push('VERSION:2.0')
  events.push('PRODID:-//TimeLedger//課表管理系統//ZH')
  events.push('CALSCALE:GREGORIAN')
  events.push('METHOD:PUBLISH')
  events.push('X-WR-CALNAME:TimeLedger 課表')
  events.push(`DTSTART:${now.toISOString().replace(/[-:]/g, '').split('.')[0]}Z`)
  events.push(`DTEND:${now.toISOString().replace(/[-:]/g, '').split('.')[0]}Z`)

  // 產生每個課程事件
  scheduleDays.value.forEach(day => {
    day.items.forEach(item => {
      const [startHour, startMinute] = item.start_time.split(':').map(Number)
      const [endHour, endMinute] = item.end_time.split(':').map(Number)

      const dayDate = new Date(day.date)
      const startDate = new Date(dayDate)
      startDate.setHours(startHour, startMinute, 0, 0)
      const endDate = new Date(dayDate)
      endDate.setHours(endHour, endMinute, 0, 0)

      const formatICalDate = (date: Date) => {
        return date.toISOString().replace(/[-:]/g, '').split('.')[0] + 'Z'
      }

      const uid = `${item.id}-${day.date}@timeledger`

      events.push('BEGIN:VEVENT')
      events.push(`UID:${uid}`)
      events.push(`DTSTAMP:${formatICalDate(now)}`)
      events.push(`DTSTART:${formatICalDate(startDate)}`)
      events.push(`DTEND:${formatICalDate(endDate)}`)
      events.push(`SUMMARY:${item.title}`)
      if (item.center_name) {
        events.push(`LOCATION:${item.center_name}`)
      }
      if (item.status && item.status !== 'APPROVED') {
        events.push(`STATUS:CONFIRMED`)
      }
      events.push('END:VEVENT')
    })
  })

  // iCal 檔案尾
  events.push('END:VCALENDAR')

  return events.join('\r\n')
}

const weekStart = computed(() => scheduleStore.weekStart)
const weekEnd = computed(() => scheduleStore.weekEnd)

// 使用 API 下載 iCal 檔案
const handleDownloadICal = async () => {
  try {
    const api = useApi()
    // 從 scheduleStore 取得週開始和結束日期，並格式化為 YYYY-MM-DD
    const startDate = weekStart.value ? formatDateToString(weekStart.value) : ''
    const endDate = weekEnd.value ? formatDateToString(weekEnd.value) : ''

    if (!startDate || !endDate) {
      console.error('Invalid date range')
      notificationUI.showError('無法取得課表日期範圍')
      return
    }

    const response = await api.raw<Blob>(`/api/v1/teacher/me/schedule.ics?start_date=${startDate}&end_date=${endDate}`)

    // 建立下載連結
    const url = URL.createObjectURL(response)
    const link = document.createElement('a')
    link.href = url
    link.download = `課表-${weekLabel.value}.ics`
    link.click()

    // 清理
    URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Failed to download iCal:', error)
    // 如果 API 失敗，回退到本地生成
    const icalData = generateICalData()
    const blob = new Blob([icalData], { type: 'text/calendar;charset=utf-8' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `課表-${weekLabel.value}.ics`
    link.click()
    URL.revokeObjectURL(link.href)
  }
}

// 分享到 LINE（產生圖片後開啟 LINE）
const handleShareLINE = async () => {
  // 直接截圖週曆區域
  const exportElement = scheduleRef.value
  if (!exportElement) {
    notificationUI.showError('無法找到課表區域')
    return
  }

  // 暫時移除 overflow 限制以確保完整截圖
  const originalOverflow = exportElement.style.overflow
  const originalMaxHeight = exportElement.style.maxHeight
  const originalPosition = exportElement.style.position

  exportElement.style.overflow = 'visible'
  exportElement.style.maxHeight = 'none'
  exportElement.style.position = 'relative'

  // 添加匯出用的 CSS 類別來隱藏時間和中心名稱
  exportElement.classList.add('export-mode')

  // 等待瀏覽器重新渲染
  await new Promise(resolve => setTimeout(resolve, 100))

  try {
    const { default: html2canvas } = await import('html2canvas')

    // 獲取完整的滾動區域尺寸
    const scrollWidth = exportElement.scrollWidth
    const scrollHeight = exportElement.scrollHeight

    const canvas = await html2canvas(exportElement, {
      backgroundColor: null,
      scale: 2,
      useCORS: true,
      logging: false,
      allowTaint: true,
      width: scrollWidth,
      height: scrollHeight,
      scrollX: 0,
      scrollY: 0,
      windowWidth: scrollWidth,
      windowHeight: scrollHeight,
    })

    // 將圖片轉為 base64
    const imageData = canvas.toDataURL('image/png')

    // 嘗試開啟 LINE
    const lineUrl = `https://line.me/R/msg/text/?${encodeURIComponent(`${authStore.user?.name} 的課表 - ${weekLabel.value}`)}`
    window.open(lineUrl, '_blank')

    // 提示使用者手動分享
    notificationUI.showSuccess('圖片已產生，請在 LINE 中貼上')

    // 下載圖片讓使用者手動分享
    const link = document.createElement('a')
    link.href = imageData
    link.download = `課表-${weekLabel.value}-${selectedViewOption.value === 'grid' ? '網格' : '列表'}-${Date.now()}.png`
    link.click()
  } catch (error) {
    console.error('Failed to share to LINE:', error)
    notificationUI.showError('分享失敗，請嘗試下載圖片')
  } finally {
    // 移除匯出用的 CSS 類別
    exportElement.classList.remove('export-mode')
    // 恢復原始樣式
    exportElement.style.overflow = originalOverflow
    exportElement.style.maxHeight = originalMaxHeight
    exportElement.style.position = originalPosition
  }
}

// 取得訂閱連結
const subscriptionUrl = computed(() => scheduleStore.subscriptionUrl)

// 建立訂閱連結
const handleCreateSubscription = async () => {
  try {
    await scheduleStore.createSubscription()
    notificationUI.showSuccess('已建立訂閱連結')
  } catch (error) {
    console.error('Failed to create subscription:', error)
    notificationUI.showError('建立訂閱連結失敗，請稍後再試')
  }
}

// 複製訂閱連結
const handleCopySubscriptionUrl = async () => {
  if (!subscriptionUrl.value) return

  try {
    await navigator.clipboard.writeText(subscriptionUrl.value)
    notificationUI.showSuccess('已複製訂閱連結')
  } catch (error) {
    console.error('Failed to copy:', error)
    notificationUI.showError('複製失敗，請手動複製')
  }
}

// 刪除訂閱連結
const handleDeleteSubscription = async () => {
  const confirmed = await alertConfirm('確定要取消訂閱嗎？取消後將無法自動同步課表。')
  if (!confirmed) return

  try {
    await scheduleStore.deleteSubscription()
    notificationUI.showSuccess('已取消訂閱')
  } catch (error) {
    console.error('Failed to delete subscription:', error)
    notificationUI.showError('取消訂閱失敗，請稍後再試')
  }
}

// 使用 API 下載圖片
const handleDownloadImage = async () => {
  try {
    // 從 scheduleStore 取得週開始和結束日期，並格式化為 YYYY-MM-DD
    const startDate = weekStart.value ? formatDateToString(weekStart.value) : ''
    const endDate = weekEnd.value ? formatDateToString(weekEnd.value) : ''

    if (!startDate || !endDate) {
      console.error('Invalid date range')
      notificationUI.showError('無法取得課表日期範圍')
      return
    }

    await scheduleStore.downloadImage(startDate, endDate)
  } catch (error) {
    console.error('Failed to download image:', error)
    notificationUI.showError('下載失敗，請稍後再試')
  }
}

onMounted(() => {
  scheduleStore.fetchSchedule()
})
</script>

<style scoped>
/* 匯出模式：隱藏時間和中心名稱 */

/* 列表視圖 */
:deep(.export-mode) .time-block {
  display: none !important;
}

:deep(.export-mode) .time-block + div {
  margin-left: 0 !important;
}

:deep(.export-mode) .flex-shrink-0.text-right {
  display: none !important;
}

:deep(.export-mode) .flex-1.min-w-0 > .text-xs.truncate {
  display: none !important;
}

/* 網格視圖 */
:deep(.export-mode) .grid .grid-cols-8 > div:first-child {
  display: none !important;
}

:deep(.export-mode) .grid .grid-cols-8 > div:nth-child(2) {
  border-top-left-radius: 0.5rem !important;
}

:deep(.export-mode) .grid .grid-cols-8 > div:nth-child(8) {
  border-top-right-radius: 0.5rem !important;
}

/* 網格課程卡片內的時間和中心名稱 */
:deep(.export-mode) .grid .absolute.inset-0 .flex.items-center.gap-1 {
  display: none !important;
}

:deep(.export-mode) .grid .absolute.inset-0 .text-\[10px\].leading-none.truncate.opacity-80 {
  display: none !important;
}

:deep(.export-mode) .grid .absolute.inset-0 .mb-0\.5 {
  margin-bottom: 0 !important;
}

/* 匯出模式：強制桌面板版 */
:deep(.export-mode) {
  /* 強制使用桌面板網格佈局 */
  display: block !important;
  width: 100% !important;
  max-width: none !important;
  overflow: visible !important;
}

:deep(.export-mode) .grid {
  display: grid !important;
  grid-template-columns: repeat(8, 1fr) !important;
  min-width: 900px !important;
}

:deep(.export-mode) .grid > div {
  flex: none !important;
  width: auto !important;
}

:deep(.export-mode) .lg\:col-span-3 {
  width: 100% !important;
  max-width: none !important;
  flex: none !important;
}

:deep(.export-mode) .glass-card {
  padding: 32px !important;
  width: 100% !important;
}

:deep(.export-mode) .min-w-\[700px\] {
  min-width: 900px !important;
}

:deep(.export-mode) .overflow-x-auto {
  overflow: visible !important;
}

/* 匯出模式：課程卡片內容置中 */
:deep(.export-mode) .grid .absolute.inset-0 {
  display: flex !important;
  flex-direction: column !important;
  align-items: center !important;
  justify-content: center !important;
}

:deep(.export-mode) .grid .absolute.inset-0 .flex.items-center.gap-1 {
  justify-content: center !important;
}
</style>
