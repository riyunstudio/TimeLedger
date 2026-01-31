<template>
  <div class="p-4 border-b border-white/10 shrink-0">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div class="flex items-center gap-4">
        <!-- 週導航區域 -->
        <div class="flex items-center gap-2">
          <button
            @click="$emit('change-week', -1)"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>

          <h2 class="text-lg font-semibold text-slate-100">
            {{ weekLabel }}
          </h2>

          <button
            @click="$emit('change-week', 1)"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

          <HelpTooltip
            v-if="showHelpTooltip"
            placement="bottom"
            title="週期導航"
            description="查看不同週期的排課狀況，預設顯示本週。"
            :usage="['點擊左右箭頭切換上週/下週', '可跨月、跨年查看', '所有視角共用同一週期']"
          />
        </div>

        <!-- 資源篩選器（僅管理員模式） -->
        <div v-if="mode === 'admin'" class="flex items-center gap-3">
          <!-- 老師篩選 -->
          <div class="flex items-center gap-2">
            <select
              :value="selectedTeacherId ?? null"
              @change="$emit('update:selected-teacher-id', Number($event.target) || null)"
              class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 min-w-[120px]"
            >
              <option :value="null">所有老師</option>
              <option v-for="teacher in teacherList" :key="'teacher-' + teacher.id" :value="teacher.id">
                {{ teacher.name }}
              </option>
            </select>
            <HelpTooltip
              v-if="showHelpTooltip"
              title="老師篩選"
              description="選擇特定老師，過濾顯示該老師的課程。"
              :usage="['選擇「所有老師」顯示所有課程', '選擇老師僅顯示該老師的課程']"
            />
          </div>

          <!-- 教室篩選 -->
          <div class="flex items-center gap-2">
            <select
              :value="selectedRoomId ?? null"
              @change="$emit('update:selected-room-id', Number($event.target) || null)"
              class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 min-w-[120px]"
            >
              <option :value="null">所有教室</option>
              <option v-for="room in roomList" :key="'room-' + room.id" :value="room.id">
                {{ room.name }}
              </option>
            </select>
            <HelpTooltip
              v-if="showHelpTooltip"
              title="教室篩選"
              description="選擇特定教室，過濾顯示該教室的課程。"
              :usage="['選擇「所有教室」顯示所有課程', '選擇教室僅顯示該教室的課程']"
            />
          </div>
        </div>
      </div>

      <!-- 右側按鈕區域 -->
      <div class="flex items-center gap-2 ml-auto">
        <!-- 新增排課按鈕（管理員） -->
        <template v-if="showCreateButton">
          <button
            @click="$emit('create-schedule')"
            class="btn-primary px-4 py-2 text-sm font-medium"
          >
            + 新增排課規則
          </button>
          <HelpTooltip
            v-if="showHelpTooltip"
            title="新增排課規則"
            description="建立新的課程排課規則，設定課程、老師、教室、時間等資訊。"
            :usage="['點擊按鈕開啟新增表單', '選擇課程、老師、教室', '設定每週固定上課日與時段', '設定有效期限後儲存']"
            shortcut="Ctrl + N"
          />
        </template>

        <!-- 個人行程按鈕（老師端） -->
        <template v-if="showPersonalEventButton">
          <button
            @click="$emit('add-personal-event')"
            class="px-4 py-2 rounded-lg bg-purple-500/20 text-purple-400 hover:bg-purple-500/30 transition-colors flex items-center gap-2 text-sm"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            個人行程
          </button>
        </template>

        <!-- 請假/調課按鈕（老師端） -->
        <template v-if="showExceptionButton">
          <button
            @click="$emit('add-exception')"
            class="px-4 py-2 rounded-lg bg-warning-500/20 text-warning-400 hover:bg-warning-500/30 transition-colors flex items-center gap-2 text-sm"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01" />
            </svg>
            請假/調課
          </button>
        </template>

        <!-- 匯出按鈕（老師端） -->
        <template v-if="showExportButton">
          <button
            @click="$emit('export')"
            class="px-4 py-2 rounded-lg bg-secondary-500/20 text-secondary-500 hover:bg-secondary-500/30 transition-colors flex items-center gap-2 text-sm"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            匯出
          </button>
        </template>
      </div>
    </div>

    <!-- 篩選提示（僅管理員模式） -->
    <div
      v-if="mode === 'admin' && (selectedTeacherId !== null || selectedRoomId !== null)"
      class="mt-3 flex items-center gap-2 px-3 py-2 bg-primary-500/10 border border-primary-500/30 rounded-lg"
    >
      <span class="text-sm text-primary-400">已篩選：</span>
      <span v-if="selectedTeacherId !== null" class="inline-flex items-center gap-1 px-2 py-1 bg-primary-500/20 rounded text-sm">
        <span class="text-white">{{ selectedTeacherName }}</span>
        <button @click="$emit('clear-teacher-filter')" class="hover:text-primary-300">
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </span>
      <span v-if="selectedRoomId !== null" class="inline-flex items-center gap-1 px-2 py-1 bg-primary-500/20 rounded text-sm">
        <span class="text-white">{{ selectedRoomName }}</span>
        <button @click="$emit('clear-room-filter')" class="hover:text-primary-300">
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </span>
      <button
        v-if="selectedTeacherId !== null || selectedRoomId !== null"
        @click="$emit('clear-all-filters')"
        class="ml-auto text-xs text-slate-400 hover:text-white transition-colors"
      >
        清除全部
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  // 模式：'admin' 或 'teacher'
  mode: 'admin' | 'teacher'
  // 週標籤
  weekLabel: string
  // 老師列表
  teacherList: Array<{ id: number; name: string }>
  // 教室列表
  roomList: Array<{ id: number; name: string }>
  // 選中的老師 ID
  selectedTeacherId: number | null
  // 選中的教室 ID
  selectedRoomId: number | null
  // 選中老師名稱
  selectedTeacherName: string
  // 選中教室名稱
  selectedRoomName: string
  // 是否顯示新增排課按鈕
  showCreateButton: boolean
  // 是否顯示個人行程按鈕（老師端）
  showPersonalEventButton: boolean
  // 是否顯示請假/調課按鈕（老師端）
  showExceptionButton: boolean
  // 是否顯示匯出按鈕（老師端）
  showExportButton: boolean
  // 是否顯示說明提示
  showHelpTooltip: boolean
}>()

// ============================================
// Emits 定義
// ============================================

defineEmits<{
  'change-week': [delta: number]
  'create-schedule': []
  'add-personal-event': []
  'add-exception': []
  'export': []
  'update:selected-teacher-id': [value: number | null]
  'update:selected-room-id': [value: number | null]
  'clear-teacher-filter': []
  'clear-room-filter': []
  'clear-all-filters': []
}>()
</script>
