<template>
  <div class="h-full flex flex-col glass-card overflow-hidden">
    <div class="p-4 border-b border-white/10 shrink-0">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex items-center gap-4">
          <!-- 週導航區域 -->
          <div class="flex items-center gap-2">
            <button
              @click="changeWeek(-1)"
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
              @click="changeWeek(1)"
              class="p-2 rounded-lg hover:bg-white/10 transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>

            <HelpTooltip
              v-if="effectiveShowHelpTooltip"
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
                v-model="selectedTeacherId"
                class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 min-w-[120px]"
              >
                <option :value="null">所有老師</option>
                <option v-for="teacher in teacherList" :key="'teacher-' + teacher.id" :value="teacher.id">
                  {{ teacher.name }}
                </option>
              </select>
              <HelpTooltip
                v-if="effectiveShowHelpTooltip"
                title="老師篩選"
                description="選擇特定老師，過濾顯示該老師的課程。"
                :usage="['選擇「所有老師」顯示所有課程', '選擇老師僅顯示該老師的課程']"
              />
            </div>

            <!-- 教室篩選 -->
            <div class="flex items-center gap-2">
              <select
                v-model="selectedRoomId"
                class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 min-w-[120px]"
              >
                <option :value="null">所有教室</option>
                <option v-for="room in roomList" :key="'room-' + room.id" :value="room.id">
                  {{ room.name }}
                </option>
              </select>
              <HelpTooltip
                v-if="effectiveShowHelpTooltip"
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
          <template v-if="effectiveShowCreateButton">
            <button
              @click="showCreateModal = true"
              class="btn-primary px-4 py-2 text-sm font-medium"
            >
              + 新增排課規則
            </button>
            <HelpTooltip
              v-if="effectiveShowHelpTooltip"
              title="新增排課規則"
              description="建立新的課程排課規則，設定課程、老師、教室、時間等資訊。"
              :usage="['點擊按鈕開啟新增表單', '選擇課程、老師、教室', '設定每週固定上課日與時段', '設定有效期限後儲存']"
              shortcut="Ctrl + N"
            />
          </template>

          <!-- 個人行程按鈕（老師端） -->
          <template v-if="effectiveShowPersonalEventButton">
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
          <template v-if="effectiveShowExceptionButton">
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
          <template v-if="effectiveShowExportButton">
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
          <button @click="clearTeacherFilter" class="hover:text-primary-300">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </span>
        <span v-if="selectedRoomId !== null" class="inline-flex items-center gap-1 px-2 py-1 bg-primary-500/20 rounded text-sm">
          <span class="text-white">{{ selectedRoomName }}</span>
          <button @click="clearRoomFilter" class="hover:text-primary-300">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </span>
        <button
          v-if="selectedTeacherId !== null || selectedRoomId !== null"
          @click="clearAllFilters"
          class="ml-auto text-xs text-slate-400 hover:text-white transition-colors"
        >
          清除全部
        </button>
      </div>
    </div>

    <div
      class="flex-1 overflow-auto p-4"
      @dragover.prevent="handleDragOver"
      @drop="handleDrop"
    >
      <!-- 週曆視圖 -->
      <div class="min-w-[600px] relative" ref="calendarContainerRef">
        <!-- 表頭 -->
        <div class="grid sticky top-0 z-10 bg-slate-900/95 backdrop-blur-sm" style="grid-template-columns: 80px repeat(7, 1fr);">
          <div class="p-2 border-b border-white/10 text-center">
            <span class="text-xs text-slate-400">時段</span>
          </div>
          <div
            v-for="day in weekDays"
            :key="day.value"
            class="p-2 border-b border-white/10 text-center"
          >
            <span class="text-sm font-medium text-slate-100">{{ day.name }}</span>
          </div>
        </div>

        <!-- 時間列和網格區域 -->
        <div class="relative">
          <!-- 時間格子 - 先渲染，在底部 -->
          <div
            v-for="time in timeSlots"
            :key="time"
            class="grid relative z-0"
            style="grid-template-columns: 80px repeat(7, 1fr);"
          >
            <!-- 時間標籤 -->
            <div class="p-2 border-r border-b border-white/5 text-right text-xs text-slate-400">
              {{ formatTime(time) }}
            </div>

            <!-- 每日網格 -->
            <div
              v-for="day in weekDays"
              :key="`${time}-${day.value}`"
              class="p-0 min-h-[60px] border-b border-white/5 border-r relative z-0"
              :class="getCellClass(time, day.value)"
              @dragenter="handleDragEnter(time, day.value)"
              @dragleave="handleDragLeave"
              @dragover.prevent
            />
          </div>

          <!-- 課程卡片層 - 虛擬滾動優化，z-index 較高確保在網格之上 -->
          <div class="absolute top-0 left-0 right-0 bottom-0 pointer-events-none z-10">
            <DynamicScroller
              :items="virtualizedSchedules"
              :min-item-size="60"
              class="h-full"
              key-field="key"
              v-if="filteredSchedules.length > 0"
            >
              <template #default="{ item, index, active }">
                <DynamicScrollerItem
                  :item="item"
                  :active="active"
                  :size-dependencies="[
                    item.offering_name,
                    item.start_time,
                    item.end_time,
                    item.teacher_name,
                    item.has_exception
                  ]"
                  :data-index="index"
                >
                  <template v-if="item.is_personal_event">
                    <!-- 個人行程保持原樣顯示 -->
                    <div
                      class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-90 transition-opacity pointer-events-auto"
                      :class="getScheduleCardClass(item)"
                      :style="[getScheduleStyle(item), { backgroundColor: item.color_hex + '40', borderColor: item.color_hex + '60' }]"
                      @click="selectSchedule(item)"
                    >
                      <div class="font-medium truncate text-white">
                        {{ item.offering_name }}
                      </div>
                      <div class="text-slate-400 truncate text-[10px]">
                        {{ item.start_time }} - {{ item.end_time }}
                      </div>
                    </div>
                  </template>
                  <template v-else>
                    <!-- 中心課程：檢查是否有重疊 -->
                    <!-- 僅顯示第一個課程（帶重疊指示器） -->
                    <div
                      v-if="getOverlapCount(item) === 1"
                      class="absolute rounded-lg p-2 text-xs cursor-pointer hover:opacity-90 transition-opacity pointer-events-auto"
                      :class="getScheduleCardClass(item)"
                      :style="getScheduleStyle(item)"
                      @click="selectSchedule(item)"
                    >
                      <div class="font-medium truncate">
                        {{ item.offering_name }}
                      </div>
                      <div v-if="effectiveCardInfoType === 'teacher'" class="text-slate-400 truncate">
                        {{ item.teacher_name }}
                      </div>
                      <div v-else class="text-slate-400 truncate">
                        {{ item.center_name }}
                      </div>
                      <div class="text-slate-500 text-[10px] mt-0.5">
                        {{ item.start_time }} - {{ item.end_time }}
                      </div>
                    </div>
                    <!-- 重疊指示器（顯示數量） -->
                    <div
                      v-else-if="getOverlapCount(item) > 1 && isFirstInOverlap(item)"
                      class="absolute rounded-lg bg-warning-500/20 border border-warning-500/50 p-2 text-xs cursor-pointer hover:bg-warning-500/30 transition-opacity pointer-events-auto"
                      :style="getScheduleStyle(item)"
                      @click="handleOverlapClick(item)"
                    >
                      <div class="flex items-center justify-center h-full">
                        <span class="text-warning-400 font-bold text-lg">
                          {{ getOverlapCount(item) }}
                        </span>
                        <span class="text-warning-300 ml-1 text-xs">堂課程</span>
                      </div>
                    </div>
                  </template>
                </DynamicScrollerItem>
              </template>
            </DynamicScroller>
          </div>
        </div>

      <!-- 矩陣視圖（已停用） -->
      <!--
        矩陣視圖功能已停用，如有需要可重新啟用
      -->
      </div>
    </div>

    <!-- 管理員專屬彈窗 -->
    <Teleport to="body">
      <ScheduleDetailPanel
        v-if="selectedSchedule && mode === 'admin'"
        :time="selectedSchedule.start_hour"
        :weekday="selectedSchedule.weekday"
        :schedule="selectedSchedule"
        @close="closeSchedulePanel"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </Teleport>

    <Teleport to="body">
      <UpdateModeModal
        v-if="showUpdateModeModal && mode === 'admin'"
        :show="showUpdateModeModal"
        :rule-name="editingRule?.offering_name"
        :rule-date="editingRule?.date ? new Date(editingRule.date).toLocaleDateString('zh-TW', { year: 'numeric', month: 'long', day: 'numeric' }) : ''"
        @close="handleUpdateModeClose"
        @confirm="handleUpdateModeConfirm"
      />
    </Teleport>

    <Teleport to="body">
      <ScheduleRuleModal
        v-if="showCreateModal && mode === 'admin'"
        @close="showCreateModal = false"
        @created="handleRuleCreated"
      />
      <ScheduleRuleModal
        v-if="showEditModal && mode === 'admin'"
        :editing-rule="editingRule"
        :update-mode="pendingUpdateMode"
        @close="handleEditModalClose"
        @submit="handleRuleUpdated"
      />
    </Teleport>

    <!-- 老師端專屬彈窗 - 動作選擇對話框 -->
    <Teleport to="body">
      <div
        v-if="showActionDialog && mode === 'teacher'"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="closeActionDialog"
      >
        <div class="glass-card w-full max-w-sm">
          <div class="p-4 border-b border-white/10">
            <h3 class="text-lg font-semibold text-white">選擇操作</h3>
            <p v-if="actionDialogItem" class="text-sm text-slate-400 mt-1">
              {{ actionDialogItem.offering_name }}
            </p>
          </div>
          <div class="p-4 space-y-3">
            <!-- 個人行程選項 -->
            <template v-if="actionDialogItem?.is_personal_event">
              <button
                @click="handleActionSelect('edit')"
                class="w-full p-4 rounded-lg bg-success-500/20 border border-success-500/30 hover:bg-success-500/30 transition-colors text-left"
              >
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-success-500/30 flex items-center justify-center">
                    <svg class="w-5 h-5 text-success-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </div>
                  <div>
                    <div class="font-medium text-white">編輯行程</div>
                    <div class="text-xs text-slate-400">修改行程時間或內容</div>
                  </div>
                </div>
              </button>
              <button
                @click="handleActionSelect('note')"
                class="w-full p-4 rounded-lg bg-primary-500/20 border border-primary-500/30 hover:bg-primary-500/30 transition-colors text-left"
              >
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-primary-500/30 flex items-center justify-center">
                    <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                  </div>
                  <div>
                    <div class="font-medium text-white">行程備註</div>
                    <div class="text-xs text-slate-400">為行程添加備註資訊</div>
                  </div>
                </div>
              </button>
              <button
                @click="handleActionSelect('delete')"
                class="w-full p-4 rounded-lg bg-critical-500/20 border border-critical-500/30 hover:bg-critical-500/30 transition-colors text-left"
              >
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-critical-500/30 flex items-center justify-center">
                    <svg class="w-5 h-5 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </div>
                  <div>
                    <div class="font-medium text-white">刪除行程</div>
                    <div class="text-xs text-slate-400">移除此個人行程</div>
                  </div>
                </div>
              </button>
            </template>

            <!-- 中心課程選項 -->
            <template v-else>
              <button
                @click="handleActionSelect('exception')"
                class="w-full p-4 rounded-lg bg-warning-500/20 border border-warning-500/30 hover:bg-warning-500/30 transition-colors text-left"
              >
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-warning-500/30 flex items-center justify-center">
                    <svg class="w-5 h-5 text-warning-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                  </div>
                  <div>
                    <div class="font-medium text-white">課程例外申請</div>
                    <div class="text-xs text-slate-400">申請調課、請假或找代課</div>
                  </div>
                </div>
              </button>
              <button
                @click="handleActionSelect('note')"
                class="w-full p-4 rounded-lg bg-primary-500/20 border border-primary-500/30 hover:bg-primary-500/30 transition-colors text-left"
              >
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-primary-500/30 flex items-center justify-center">
                    <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </div>
                  <div>
                    <div class="font-medium text-white">課堂備註</div>
                    <div class="text-xs text-slate-400">撰寫或查看課程筆記</div>
                  </div>
                </div>
              </button>
            </template>
          </div>
          <div class="p-4 border-t border-white/10">
            <button
              @click="closeActionDialog"
              class="w-full px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              取消
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 管理員端專屬彈窗 - 重疊課程選擇對話框 -->
    <Teleport to="body">
      <div
        v-if="showOverlapDialog && mode === 'admin'"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="closeOverlapDialog"
      >
        <div class="glass-card w-full max-w-sm">
          <div class="p-4 border-b border-white/10">
            <h3 class="text-lg font-semibold text-white">選擇課程</h3>
            <p v-if="overlapTimeSlot" class="text-sm text-slate-400 mt-1">
              {{ getWeekdayText(overlapTimeSlot.weekday) }}
              {{ overlapTimeSlot.start_hour.toString().padStart(2, '0') }}:{{ overlapTimeSlot.start_minute.toString().padStart(2, '0') }}
            </p>
          </div>
          <div class="max-h-96 overflow-y-auto">
            <div
              v-for="schedule in overlapSchedules"
              :key="schedule.id"
              class="p-4 border-b border-white/5 hover:bg-white/5 cursor-pointer transition-colors"
              @click="selectFromOverlap(schedule)"
            >
              <div class="flex items-center justify-between">
                <div>
                  <div class="font-medium text-white">{{ schedule.offering_name }}</div>
                  <div class="text-sm text-slate-400 mt-1">
                    {{ schedule.start_time }} - {{ schedule.end_time }}
                  </div>
                  <div v-if="effectiveCardInfoType === 'teacher'" class="text-xs text-slate-500 mt-1">
                    {{ schedule.teacher_name }}
                  </div>
                  <div v-else class="text-xs text-slate-500 mt-1">
                    {{ schedule.center_name }}
                  </div>
                </div>
                <div v-if="schedule.room_name" class="text-xs text-slate-500">
                  {{ schedule.room_name }}
                </div>
              </div>
            </div>
          </div>
          <div class="p-4 border-t border-white/10">
            <button
              @click="closeOverlapDialog"
              class="w-full px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              取消
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { formatDateToString } from '~/composables/useTaiwanTime'
import { nextTick, ref, computed, watch, onMounted, onUnmounted } from 'vue'

// ============================================
// Props 定義
// ============================================

const props = defineProps<{
  // 模式：'admin' 或 'teacher'
  mode: 'admin' | 'teacher'
  // 排課資料（可選，如果提供則使用此資料，否則自動獲取）
  schedules?: any[]
  // API 端點
  apiEndpoint: string
  // 卡片顯示類型：'teacher' 顯示老師名稱，'center' 顯示中心名稱
  cardInfoType?: 'teacher' | 'center'
  // 是否顯示矩陣視圖
  showMatrixView?: boolean
  // 是否顯示視圖切換器
  showViewModeSelector?: boolean
  // 是否顯示新增排課按鈕
  showCreateButton?: boolean
  // 是否顯示個人行程按鈕（老師端）
  showPersonalEventButton?: boolean
  // 是否顯示請假/調課按鈕（老師端）
  showExceptionButton?: boolean
  // 是否顯示匯出按鈕（老師端）
  showExportButton?: boolean
  // 是否顯示說明提示
  showHelpTooltip?: boolean
}>()

// 為可選屬性提供計算後的預設值
const effectiveApiEndpoint = computed(() => props.apiEndpoint || '/admin/expand-rules')
const effectiveCardInfoType = computed(() => props.cardInfoType || 'teacher')
const effectiveShowMatrixView = computed(() => props.showMatrixView ?? true)
const effectiveShowViewModeSelector = computed(() => props.showViewModeSelector ?? true)
const effectiveShowCreateButton = computed(() => props.showCreateButton ?? true)
const effectiveShowPersonalEventButton = computed(() => props.showPersonalEventButton ?? false)
const effectiveShowExceptionButton = computed(() => props.showExceptionButton ?? false)
const effectiveShowExportButton = computed(() => props.showExportButton ?? false)
const effectiveShowHelpTooltip = computed(() => props.showHelpTooltip ?? true)

// ============================================
// Emits 定義
// ============================================

const emit = defineEmits<{
  selectCell: { time: number, weekday: number }
  'update:weekStart': [value: Date]
  'select-schedule': [schedule: any]
  'add-personal-event': []
  'add-exception': []
  'export': []
  'edit-personal-event': [event: any]
  'delete-personal-event': [event: any]
  'personal-event-note': [event: any]
}>()

// ============================================
// Composables
// ============================================

const { confirm: confirmDialog, error: alertError } = useAlert()
const { resourceCache, fetchAllResources } = useResourceCache()
const { getCenterId } = useCenterId()

// ============================================
// DOM 引用
// ============================================

const calendarContainerRef = ref<HTMLElement | null>(null)
const slotWidth = ref(100)

// ============================================
// 常量定義
// ============================================

const TIME_SLOT_HEIGHT = 60 // 每個時段格子的高度 (px)
const TIME_COLUMN_WIDTH = 80 // 時間列寬度 (px)

// ============================================
// 狀態管理
// ============================================

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showUpdateModeModal = ref(false)
const editingRule = ref<any>(null)
const pendingUpdateMode = ref<string>('')
const selectedCell = ref<{ time: number, day: number } | null>(null)
const selectedSchedule = ref<any>(null)
const dragTarget = ref<{ time: number, day: number } | null>(null)
const validationResults = ref<Record<string, any>>({})

// 老師端動作選擇對話框狀態
const showActionDialog = ref(false)
const actionDialogItem = ref<any>(null)

// 管理員端重疊課程選擇對話框狀態
const showOverlapDialog = ref(false)
const overlapSchedules = ref<any[]>([])
const overlapTimeSlot = ref<{ weekday: number, start_hour: number, start_minute: number } | null>(null)

// 週起始日期
const getWeekStart = (date: Date): Date => {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  return new Date(d.setDate(diff))
}

const weekStart = ref(getWeekStart(new Date()))
const weekEnd = computed(() => {
  const end = new Date(weekStart.value)
  end.setDate(end.getDate() + 6)
  return end
})

const weekLabel = computed(() => {
  const start = weekStart.value.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric' })
  const end = weekEnd.value.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric', year: 'numeric' })
  return `${start} - ${end}`
})

// ============================================
// 視圖模式管理
// ============================================

// 內部篩選狀態
const selectedTeacherId = ref<number | null>(null)
const selectedRoomId = ref<number | null>(null)

// 老師列表
const teacherList = computed(() => {
  return Array.from(resourceCache.value.teachers.values())
})

// 教室列表
const roomList = computed(() => {
  return Array.from(resourceCache.value.rooms.values())
})

// 資源列表（向後相容）
const resourceList = computed(() => {
  return teacherList.value
})

// 矩陣視圖的資源列表（向後相容）
const matrixResources = computed(() => {
  return teacherList.value
})

// 選中老師名稱
const selectedTeacherName = computed(() => {
  const teacher = resourceCache.value.teachers.get(selectedTeacherId.value)
  if (teacher) {
    return teacher.name
  }
  return ''
})

// 選中教室名稱
const selectedRoomName = computed(() => {
  const room = resourceCache.value.rooms.get(selectedRoomId.value)
  if (room) {
    return room.name
  }
  return ''
})

const clearTeacherFilter = () => {
  selectedTeacherId.value = null
}

const clearRoomFilter = () => {
  selectedRoomId.value = null
}

const clearAllFilters = () => {
  selectedTeacherId.value = null
  selectedRoomId.value = null
}

// ============================================
// 時間和日期配置
// ============================================

// 時間段 - 包含 00:00-03:00 和 22:00-23:00
const timeSlots = [0, 1, 2, 3, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23]

const weekDays = [
  { value: 1, name: '週一' },
  { value: 2, name: '週二' },
  { value: 3, name: '週三' },
  { value: 4, name: '週四' },
  { value: 5, name: '週五' },
  { value: 6, name: '週六' },
  { value: 7, name: '週日' },
]

// ============================================
// 排課資料
// ============================================

const schedules = ref<any[]>([])

// 週切換
const changeWeek = (delta: number) => {
  weekStart.value = getWeekStart(new Date(weekStart.value.getTime() + delta * 7 * 24 * 60 * 60 * 1000))
  emit('update:weekStart', weekStart.value)
}

// 監聽週變化（僅在沒有傳入 schedules prop 時才獲取資料）
watch(weekStart, async () => {
  if (!props.schedules || props.schedules.length === 0) {
    await fetchSchedules()
  }
})

// 取得排課資料
const fetchSchedules = async () => {
  try {
    const api = useApi()

    // 取得當前週的日期範圍
    const startDate = formatDateToString(weekStart.value)
    const endDate = formatDateToString(weekEnd.value)

    // 根據模式選擇 API
    let response
    if (props.mode === 'teacher') {
      // 老師端使用 /teacher/schedules
      response = await api.get<{ code: number; datas: any[] }>('/teacher/schedules', {
        start_date: startDate,
        end_date: endDate,
      })
    } else {
      // 管理員使用 expand-rules API
      response = await api.post<{ code: number; datas: any[] }>(effectiveApiEndpoint.value, {
        rule_ids: [],
        start_date: startDate,
        end_date: endDate,
      })
    }

    const expandedSchedules = response.datas || []

    // 將展開後的排課轉換為前端格式
    const scheduleList = expandedSchedules.map((schedule: any) => {
      const date = new Date(schedule.date)
      const weekday = date.getDay() === 0 ? 7 : date.getDay()
      const startTime = schedule.start_time || '09:00'
      const endTime = schedule.end_time || '10:00'
      const [startHour, startMinute] = startTime.split(':').map(Number)
      const durationMinutes = calculateDurationMinutes(startTime, endTime)

      return {
        id: schedule.rule_id || schedule.id,
        key: `${schedule.rule_id || schedule.id}-${weekday}-${startTime}-${schedule.date}`,
        offering_name: schedule.offering_name || '-',
        teacher_name: schedule.teacher_name || '-',
        teacher_id: schedule.teacher_id,
        center_name: schedule.center_name || '-',
        center_id: schedule.center_id,
        room_id: schedule.room_id,
        room_name: schedule.room_name || '-',
        weekday: weekday,
        start_time: startTime,
        end_time: endTime,
        start_hour: startHour,
        start_minute: startMinute,
        duration_minutes: durationMinutes,
        date: schedule.date,
        has_exception: schedule.has_exception || false,
        exception_type: schedule.exception_type || null,
        exception_info: schedule.exception_info || null,
        rule: schedule.rule || null,
        offering_id: schedule.offering_id,
        effective_range: schedule.effective_range || null,
      }
    })

    schedules.value = scheduleList

    // 等待 DOM 更新後計算槽寬度
    await nextTick()
    calculateSlotWidth()
  } catch (error) {
    console.error('Failed to fetch schedules:', error)
    schedules.value = []
  }
}

// 去重後的排課（優先使用 prop 資料，否則使用內部資料）
const displaySchedules = computed(() => {
  const sourceSchedules = props.schedules && props.schedules.length > 0
    ? props.schedules
    : schedules.value

  const seen = new Set<string>()
  const result: any[] = []

  for (const schedule of sourceSchedules) {
    const key = `${schedule.id}-${schedule.weekday}-${schedule.start_time}`
    if (!seen.has(key)) {
      seen.add(key)
      result.push(schedule)
    }
  }

  return result
})

// 計算每個時段的重疊數量（基於過濾後的課程）
const overlapCountMap = computed(() => {
  const countMap: Record<string, number> = {}

  for (const schedule of filteredSchedules.value) {
    const key = `${schedule.weekday}-${schedule.start_hour}-${schedule.start_minute}`
    countMap[key] = (countMap[key] || 0) + 1
  }

  return countMap
})

// 取得某個課程的重疊數量
const getOverlapCount = (schedule: any) => {
  const key = `${schedule.weekday}-${schedule.start_hour}-${schedule.start_minute}`
  return overlapCountMap.value[key] || 1
}

// 判斷是否為重疊群組中的第一個（用於顯示指示器）
const isFirstInOverlap = (schedule: any) => {
  const key = `${schedule.weekday}-${schedule.start_hour}-${schedule.start_minute}`
  const allAtSameTime = filteredSchedules.value.filter(s =>
    `${s.weekday}-${s.start_hour}-${s.start_minute}` === key
  )
  // 按 id 排序，返回第一個
  allAtSameTime.sort((a, b) => a.id - b.id)
  return allAtSameTime[0]?.id === schedule.id
}

// 取得同一時段的所有課程
const getSchedulesAtSameTime = (schedule: any) => {
  return filteredSchedules.value.filter(s =>
    s.weekday === schedule.weekday &&
    s.start_hour === schedule.start_hour &&
    s.start_minute === schedule.start_minute
  )
}

// 週曆視圖顯示的課程（根據選中的老師和教室過濾）
const filteredSchedules = computed(() => {
  // 如果都沒有選中，返回全部
  if (!selectedTeacherId.value && !selectedRoomId.value) {
    return displaySchedules.value
  }

  // 同時篩選老師和教室（AND 邏輯）
  return displaySchedules.value.filter(schedule => {
    // 檢查老師（如果選中了老師）
    const teacherMatch = !selectedTeacherId.value || schedule.teacher_id === selectedTeacherId.value
    // 檢查教室（如果選中了教室）
    const roomMatch = !selectedRoomId.value || schedule.room_id === selectedRoomId.value
    // 兩者都符合才返回
    return teacherMatch && roomMatch
  })
})

// 虛擬滾動用的課程列表
// 當課程數量超過閾值時啟用虛擬滾動
const virtualizedSchedules = computed(() => {
  const schedules = filteredSchedules.value
  // 如果課程數量少於 50 筆，不需要虛擬滾動，直接返回原陣列
  if (schedules.length < 50) {
    return schedules
  }
  // 否則返回原陣列，DynamicScroller 會自動處理虛擬化
  return schedules
})

// ============================================
// 卡片定位計算
// ============================================

// 計算格子寬度
const calculateSlotWidth = () => {
  if (calendarContainerRef.value) {
    const containerWidth = calendarContainerRef.value.offsetWidth
    slotWidth.value = Math.max(80, (containerWidth - TIME_COLUMN_WIDTH) / 7)
  }
}

// 計算課程持續分鐘數
const calculateDurationMinutes = (startTime: string, endTime: string): number => {
  const [startHour, startMinute] = startTime.split(':').map(Number)
  const [endHour, endMinute] = endTime.split(':').map(Number)

  const startMinutes = startHour * 60 + startMinute
  let endMinutes = endHour * 60 + endMinute

  // 跨日處理
  if (endMinutes <= startMinutes) {
    endMinutes += 24 * 60
  }

  return endMinutes - startMinutes
}

const formatTime = (hour: number): string => {
  return `${hour.toString().padStart(2, '0')}:00`
}

// 取得星期文字
const getWeekdayText = (weekday: number): string => {
  const days = ['週日', '週一', '週二', '週三', '週四', '週五', '週六']
  return days[weekday - 1] || ''
}

// 週曆視圖：計算課程卡片樣式
const getScheduleStyle = (schedule: any) => {
  const { weekday, start_hour, start_minute, duration_minutes } = schedule

  // 計算水平位置 - 對齊到星期網格
  const dayIndex = weekday - 1 // 0-6
  const left = TIME_COLUMN_WIDTH + (dayIndex * slotWidth.value)

  // 計算垂直位置
  let topSlotIndex = 0
  for (let t = 0; t < start_hour; t++) {
    if (t >= 0 && t <= 3) {
      topSlotIndex++ // 0-3 時段每個都算
    } else if (t >= 9) {
      topSlotIndex++ // 9 以後的時段每個都算
    }
  }

  const slotHeight = TIME_SLOT_HEIGHT
  const baseTop = topSlotIndex * slotHeight
  const minuteOffset = (start_minute / 60) * slotHeight
  const top = baseTop + minuteOffset

  // 計算高度（持續分鐘數轉像素）
  const height = (duration_minutes / 60) * slotHeight

  // 計算寬度（略小於格子寬度以留邊距）
  const width = slotWidth.value - 4

  return {
    left: `${left}px`,
    top: `${top}px`,
    width: `${width}px`,
    height: `${height}px`,
  }
}

// ============================================
// 樣式相關
// ============================================

const getCellClass = (time: number, weekday: number): string => {
  const key = `${time}-${weekday}`
  const validation = validationResults.value[key]

  if (validation?.valid === false) {
    return 'bg-critical-500/10 border-critical-500/50'
  } else if (validation?.warning) {
    return 'bg-warning-500/10 border-warning-500/50'
  } else if (validation?.valid === true) {
    return 'bg-success-500/10 border-success-500/50'
  }

  return 'hover:bg-white/5'
}

const getScheduleCardClass = (schedule: any): string => {
  if (!schedule) return ''

  // 個人行程使用不同的樣式
  if (schedule.is_personal_event) {
    const baseColor = schedule.color_hex || '#a855f7' // 紫色預設
    return 'border border-white/20'
  }

  if (schedule.has_exception) {
    switch (schedule.exception_type) {
      case 'CANCEL':
        return 'bg-critical-500/30 border border-critical-500/50 line-through'
      case 'RESCHEDULE':
        return 'bg-warning-500/30 border border-warning-500/50'
      case 'SWAP':
        return 'bg-primary-500/30 border border-primary-500/50'
      default:
        return 'bg-slate-700/80 border border-white/10'
    }
  }

  return 'bg-slate-700/80 border border-white/10'
}

// ============================================
// 互動處理
// ============================================

const selectSchedule = (schedule: any) => {
  if (props.mode === 'teacher') {
    // 老師端顯示動作選擇對話框
    actionDialogItem.value = schedule
    showActionDialog.value = true
    emit('select-schedule', schedule)
  } else {
    // 管理員端顯示詳情面板
    selectedSchedule.value = schedule
    emit('selectCell', { time: schedule.start_hour, weekday: schedule.weekday })
  }
}

// 關閉動作選擇對話框
const closeActionDialog = () => {
  showActionDialog.value = false
  actionDialogItem.value = null
}

// 關閉管理員課程詳情面板
const closeSchedulePanel = () => {
  selectedSchedule.value = null
  selectedCell.value = null
}

// 處理重疊指示器點擊（管理員端）
const handleOverlapClick = (schedule: any) => {
  const schedulesAtSameTime = getSchedulesAtSameTime(schedule)

  if (schedulesAtSameTime.length === 1) {
    // 只有一堂課，直接選擇
    selectSchedule(schedule)
  } else {
    // 多堂課，顯示選擇對話框
    overlapSchedules.value = schedulesAtSameTime
    overlapTimeSlot.value = {
      weekday: schedule.weekday,
      start_hour: schedule.start_hour,
      start_minute: schedule.start_minute
    }
    showOverlapDialog.value = true
  }
}

// 關閉重疊選擇對話框
const closeOverlapDialog = () => {
  showOverlapDialog.value = false
  overlapSchedules.value = []
  overlapTimeSlot.value = null
}

// 選擇重疊課程中的某一堂
const selectFromOverlap = (schedule: any) => {
  selectSchedule(schedule)
  closeOverlapDialog()
}

// 處理動作選擇
const handleActionSelect = (action: 'exception' | 'note' | 'edit' | 'delete') => {
  const item = actionDialogItem.value
  if (!item) return

  if (item.is_personal_event) {
    // 個人行程處理
    if (action === 'edit') {
      emit('edit-personal-event', item)
    } else if (action === 'delete') {
      emit('delete-personal-event', item)
    } else if (action === 'note') {
      // 個人行程備註 - 延遲 emit 確保對話框已關閉
      setTimeout(() => {
        emit('personal-event-note', { ...item, action: 'note' })
      }, 100)
    }
  } else {
    // 中心課程處理
    if (action === 'exception') {
      emit('add-exception', item)
    } else if (action === 'note') {
      // 課堂備註功能 - 延遲 emit 確保對話框已關閉
      setTimeout(() => {
        emit('select-schedule', { ...item, action: 'note' })
      }, 100)
    }
  }

  closeActionDialog()
}

// 管理員端：編輯
const handleEdit = () => {
  if (selectedSchedule.value) {
    editingRule.value = selectedSchedule.value
    showUpdateModeModal.value = true
  }
}

// 管理員端：更新模式確認
const handleUpdateModeClose = () => {
  showUpdateModeModal.value = false
  editingRule.value = null
}

const handleUpdateModeConfirm = (mode: string) => {
  pendingUpdateMode.value = mode
  showUpdateModeModal.value = false
  showEditModal.value = true
}

// 管理員端：編輯 modal 關閉
const handleEditModalClose = () => {
  showEditModal.value = false
  editingRule.value = null
  pendingUpdateMode.value = ''
}

// 管理員端：刪除
const handleDelete = async () => {
  const confirmed = await confirmDialog('確定要刪除此排課規則？')
  if (!confirmed || !selectedSchedule.value) return

  try {
    const api = useApi()
    await api.delete(`/admin/rules/${selectedSchedule.value.id}`)
    selectedCell.value = null
    selectedSchedule.value = null
    await fetchSchedules()
  } catch (err) {
    console.error('Failed to delete rule:', err)
    await alertError('刪除失敗，請稍後再試')
  }
}

// 管理員端：更新排課規則
const handleRuleUpdated = async (formData: any, updateMode: string) => {
  try {
    const api = useApi()
    await api.put(`/admin/rules/${editingRule.value.id}`, {
      ...formData,
      update_mode: updateMode,
    })
    await fetchSchedules()
    selectedCell.value = null
    selectedSchedule.value = null
    editingRule.value = null
    pendingUpdateMode.value = ''
  } catch (err) {
    console.error('Failed to update rule:', err)
    await alertError('更新失敗，請稍後再試')
  }
}

// 新增排課規則後
const handleRuleCreated = () => {
  fetchSchedules()
}

// ============================================
// 拖曳處理（僅管理員）
// ============================================

const handleDragOver = (event: DragEvent) => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.time}-${dragTarget.value.day}`
    validationResults.value[key] = { valid: true }
  }
}

const handleDragEnter = (time: number, day: number) => {
  dragTarget.value = { time, day }
}

const handleDragLeave = () => {
  if (dragTarget.value) {
    const key = `${dragTarget.value.time}-${dragTarget.value.day}`
    delete validationResults.value[key]
  }
  dragTarget.value = null
}

const handleDrop = async (event: DragEvent) => {
  event.preventDefault()

  if (!dragTarget.value) return

  const data = event.dataTransfer?.getData('application/json')
  if (!data) return

  const parsed = JSON.parse(data)
  const { type, item } = parsed

  const key = `${dragTarget.value.time}-${dragTarget.value.day}`
  validationResults.value[key] = { valid: 'checking' }

  try {
    const api = useApi()
    const teacherId = type === 'teacher' ? item.id : (item.teacher_id || null)
    const roomId = type === 'room' ? item.id : (item.room_id || null)

    const response = await api.post<any>('/admin/scheduling/check-overlap', {
      teacher_id: teacherId,
      room_id: roomId,
      start_time: `${formatDateToString(weekStart.value)}T${formatTime(dragTarget.value.time)}:00`,
      end_time: `${formatDateToString(weekStart.value)}T${formatTime(dragTarget.value.time + 1)}:00`,
    })

    if (response.data.valid) {
      validationResults.value[key] = { valid: true }
    } else {
      validationResults.value[key] = { valid: false, conflicts: response.data.conflicts }
    }
  } catch (error) {
    console.error('Validation failed:', error)
    validationResults.value[key] = { valid: false, error: true }
  }

  dragTarget.value = null
}

// ============================================
// 生命週期
// ============================================

// ResizeObserver 引用（用於清理）
const resizeObserver = ref<ResizeObserver | null>(null)

onMounted(async () => {
  // 只有當沒有傳入 schedules prop 時才獲取資料
  if (!props.schedules || props.schedules.length === 0) {
    await fetchSchedules()
  }

  // 管理員模式需要載入資源
  if (props.mode === 'admin') {
    fetchAllResources()
  }

  // 等待 DOM 更新後計算槽寬度
  await nextTick()
  calculateSlotWidth()

  // 監控容器大小變化
  if (calendarContainerRef.value) {
    resizeObserver.value = new ResizeObserver(() => {
      calculateSlotWidth()
    })
    resizeObserver.value.observe(calendarContainerRef.value)
  }
})

// 確保在組件卸載時正確清理
onUnmounted(() => {
  if (resizeObserver.value) {
    resizeObserver.value.disconnect()
  }
})
</script>
