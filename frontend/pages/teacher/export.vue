<template>
  <div class="flex items-center justify-between mb-6">
    <div>
      <h2 class="text-2xl font-bold text-white mb-1">匯出課表</h2>
      <p class="text-slate-400 text-sm">預覽並下載您的詳細課表</p>
    </div>
    <div class="flex gap-3">
      <button
        @click="router.push('/teacher/dashboard')"
        class="px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
      >
        返回
      </button>
      <button
        @click="handleDownloadPDF"
        class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        下載 PDF
      </button>
      <button
        @click="handleDownloadImage"
        class="px-4 py-2 rounded-lg bg-secondary-500 text-white hover:bg-secondary-600 transition-colors flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        下載圖片
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
              <h1 class="text-2xl font-bold" :class="currentTheme.titleClass">{{ authStore.user?.name }}</h1>
              <p class="text-sm" :class="currentTheme.subtitleClass">個人課表</p>
              <div v-if="options.showPersonalInfo" class="flex gap-2 mt-2">
                <span class="px-2 py-0.5 rounded-full text-xs" :class="currentTheme.tagClass">鋼琴</span>
                <span class="px-2 py-0.5 rounded-full text-xs" :class="currentTheme.tagClass2">古典</span>
                <span class="px-2 py-0.5 rounded-full text-xs" :class="currentTheme.tagClass3">樂理</span>
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
            <span class="font-medium" :class="currentTheme.titleClass">{{ weekLabel }}</span>
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
        </div>

        <div class="space-y-6">
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
                class="flex items-center px-4 py-4 relative overflow-hidden"
                :class="[currentTheme.itemClass, currentTheme.itemGradient]"
              >
                <div class="absolute left-0 top-0 bottom-0 w-1" :class="currentTheme.itemAccentClass"></div>
                <div class="w-20 flex-shrink-0">
                  <div class="font-medium bg-clip-text text-transparent" :class="currentTheme.timeGradient">{{ item.start_time }}</div>
                  <div class="text-sm" :class="currentTheme.subtitleClass">{{ item.end_time }}</div>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-1 flex-wrap">
                    <span
                      class="w-2 h-2 rounded-full flex-shrink-0"
                      :style="{ backgroundColor: item.color || '#10B981' }"
                    ></span>
                    <h4 class="font-medium truncate" :class="currentTheme.itemTitleClass">
                      {{ item.title }}
                      <span v-if="item.center_name" class="font-normal" :class="currentTheme.centerClass">{{ item.center_name }}</span>
                    </h4>
                    <span
                      v-if="item.status && item.status !== 'APPROVED'"
                      class="px-2 py-0.5 rounded-full text-xs flex-shrink-0"
                      :class="getStatusClass(item.status)"
                    >
                      {{ getStatusText(item.status) }}
                    </span>
                  </div>
                  <div class="flex items-center gap-2 text-sm flex-wrap">
                    <span v-if="item.type === 'PERSONAL_EVENT'" :class="currentTheme.personalClass">個人行程</span>
                    <template v-else>
                      <span :class="currentTheme.subtitleClass">課程時段</span>
                    </template>
                  </div>
                </div>
                <div class="flex-shrink-0 text-right ml-4">
                  <p class="text-sm bg-clip-text text-transparent font-medium" :class="currentTheme.timeGradient">
                    {{ getDuration(item.start_time, item.end_time) }}
                  </p>
                  <p class="text-xs" :class="currentTheme.subtitleClass">分鐘</p>
                </div>
              </div>
            </div>

            <div v-else class="px-4 py-8 text-center" :class="currentTheme.emptyClass">
              休息日
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

  <NotificationDropdown v-if="notificationUI.show.value" @close="notificationUI.close()" />
  <TeacherSidebar v-if="sidebarStore.isOpen.value" @close="sidebarStore.close()" />
</template>

<script setup lang="ts">
 definePageMeta({
   middleware: 'auth-teacher',
   layout: 'default',
 })

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
const teacherStore = useTeacherStore()
const sidebarStore = useSidebar()
const notificationUI = useNotification()
const scheduleRef = ref<HTMLElement>()

const selectedTheme = ref('domeMouse')
const options = reactive({
  showPersonalInfo: true,
  showStats: true,
  includeNotes: false,
})

const currentTheme = computed(() => {
  return themes.find(t => t.id === selectedTheme.value) || themes[0]
})

const weekLabel = computed(() => teacherStore.weekLabel)

const scheduleDays = computed(() => {
  return teacherStore.schedule?.days || []
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
  teacherStore.changeWeek(delta)
  teacherStore.fetchSchedule()
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

const handleDownloadImage = () => {
  if (scheduleRef.value) {
    import('html2canvas').then(({ default: html2canvas }) => {
      html2canvas(scheduleRef.value!, {
        backgroundColor: getBackgroundColor(),
        scale: 2,
      }).then(canvas => {
        const link = document.createElement('a')
        link.download = `課表-${weekLabel.value}-${Date.now()}.png`
        link.href = canvas.toDataURL('image/png')
        link.click()
      })
    })
  }
}

const handleDownloadPDF = () => {
  if (scheduleRef.value) {
    import('html2canvas').then(({ default: html2canvas }) => {
      import('jspdf').then(({ default: jsPDF }) => {
        html2canvas(scheduleRef.value!, {
          backgroundColor: getBackgroundColor(),
          scale: 2,
        }).then(canvas => {
          const imgData = canvas.toDataURL('image/png')
          const pdf = new jsPDF({
            orientation: 'portrait',
            unit: 'mm',
            format: 'a4',
          })
          const imgWidth = 210
          const imgHeight = (canvas.height * imgWidth) / canvas.width
          pdf.addImage(imgData, 'PNG', 0, 0, imgWidth, imgHeight)
          pdf.save(`課表-${weekLabel.value}-${Date.now()}.pdf`)
        })
      })
    })
  }
}

onMounted(() => {
  teacherStore.fetchSchedule()
})
</script>
