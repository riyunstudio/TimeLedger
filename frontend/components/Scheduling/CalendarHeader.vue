<template>
  <div class="p-4 border-b border-white/10 shrink-0">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div class="flex items-center gap-4">
        <!-- 週導航區域（僅在大螢幕顯示，小螢幕由 MobileWeekView 顯示） -->
        <div class="hidden md:flex items-center gap-2">
          <button
            @click="$emit('change-week', -1)"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <Icon icon="chevron-left" size="lg" />
          </button>

          <!-- 日期選擇區域 - 可點擊彈出 DatePicker -->
          <div class="relative">
            <button
              @click="showDatePicker = true"
              class="flex items-center gap-2 px-3 py-1.5 rounded-lg hover:bg-white/10 transition-colors"
            >
              <h2 class="text-lg font-semibold text-slate-100 whitespace-nowrap">
                {{ weekLabel }}
              </h2>
              <Icon icon="calendar" size="sm" class="text-slate-400" />
            </button>

            <!-- DatePicker 彈出層 -->
            <Teleport to="body">
              <div
                v-if="showDatePicker"
                class="fixed inset-0 z-50 flex items-center justify-center"
                @click.self="showDatePicker = false"
              >
                <!-- 遮罩層 -->
                <div class="absolute inset-0 bg-black/50" @click="showDatePicker = false"></div>

                <!-- DatePicker 容器 -->
                <div class="relative glass-card p-4 rounded-xl shadow-2xl w-auto">
                  <div class="flex items-center justify-between mb-4">
                    <h3 class="text-white font-medium">選擇日期</h3>
                    <button
                      @click="showDatePicker = false"
                      class="p-1 rounded-lg hover:bg-white/10 transition-colors"
                    >
                      <Icon icon="close" size="sm" />
                    </button>
                  </div>

                  <!-- 原生 Date Picker -->
                  <input
                    type="date"
                    v-model="selectedDate"
                    :min="minDate"
                    :max="maxDate"
                    class="w-full px-4 py-3 bg-slate-800 border border-white/10 rounded-lg text-white text-lg focus:outline-none focus:border-primary-500"
                    @change="handleDateSelect"
                  />

                  <!-- 快捷按鈕 -->
                  <div class="mt-4 flex gap-2">
                    <button
                      @click="quickSelectToday"
                      class="flex-1 px-3 py-2 rounded-lg bg-primary-500/20 text-primary-400 hover:bg-primary-500/30 transition-colors text-sm"
                    >
                      回到今天
                    </button>
                    <button
                      @click="showDatePicker = false"
                      class="flex-1 px-3 py-2 rounded-lg bg-slate-700 text-white hover:bg-slate-600 transition-colors text-sm"
                    >
                      取消
                    </button>
                    <button
                      @click="confirmDateSelection"
                      class="flex-1 px-3 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors text-sm"
                    >
                      確定
                    </button>
                  </div>
                </div>
              </div>
            </Teleport>
          </div>

          <button
            @click="$emit('change-week', 1)"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <Icon icon="chevron-right" size="lg" />
          </button>

          <HelpTooltip
            v-if="showHelpTooltip"
            placement="bottom"
            :title="$t('help.weekNavigation')"
            :description="$t('help.weekNavigationHelp')"
            :usage="['點擊左右箭頭切換上週/下週', '可跨月、跨年查看', '所有視角共用同一週期']"
          />
        </div>

        <!-- 資源篩選器（僅管理員模式） -->
        <div v-if="mode === 'admin'" class="flex flex-wrap justify-center items-center gap-2 md:gap-3 w-full">
          <!-- 老師篩選 -->
          <div class="flex items-center gap-1 md:gap-2 w-full sm:w-auto justify-center">
            <select
              :value="selectedTeacherId ?? -1"
              @change="$emit('update:selected-teacher-id', Number($event.target.value) === -1 ? null : Number($event.target.value))"
              class="w-full sm:w-auto px-2 py-1.5 md:px-3 md:py-1.5 rounded-lg text-xs md:text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 min-w-[100px] md:min-w-[120px] max-w-[150px] md:max-w-none overflow-hidden text-ellipsis text-center"
            >
              <option :value="-1">{{ $t('schedule.allTeachers') }}</option>
              <option v-for="teacher in teacherList" :key="'teacher-' + teacher.id" :value="teacher.id">
                {{ teacher.name }}
              </option>
            </select>
            <HelpTooltip
              v-if="showHelpTooltip"
              :title="$t('filter.teacherFilter')"
              :description="$t('filter.filterByTeacher')"
              :usage="[$t('filter.showAllCourses'), $t('filter.showTeacherCourses')]"
            />
          </div>

          <!-- 教室篩選 -->
          <div class="flex items-center gap-1 md:gap-2 w-full sm:w-auto justify-center">
            <select
              :value="selectedRoomId ?? -1"
              @change="$emit('update:selected-room-id', Number($event.target.value) === -1 ? null : Number($event.target.value))"
              class="w-full sm:w-auto px-2 py-1.5 md:px-3 md:py-1.5 rounded-lg text-xs md:text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 min-w-[100px] md:min-w-[120px] max-w-[150px] md:max-w-none overflow-hidden text-ellipsis text-center"
            >
              <option :value="-1">{{ $t('schedule.allRooms') }}</option>
              <option v-for="room in roomList" :key="'room-' + room.id" :value="room.id">
                {{ room.name }}
              </option>
            </select>
            <HelpTooltip
              v-if="showHelpTooltip"
              :title="$t('filter.roomFilter')"
              :description="$t('filter.filterByRoom')"
              :usage="[$t('filter.showAllRooms'), $t('filter.showRoomCourses')]"
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
            + {{ $t('schedule.addScheduleRule') }}
          </button>
          <HelpTooltip
            v-if="showHelpTooltip"
            :title="$t('schedule.addScheduleRule')"
            :description="$t('help.addScheduleRuleHelp')"
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
            <Icon icon="plus" size="md" />
            {{ $t('schedule.personalEvent') }}
          </button>
        </template>

        <!-- 請假/調課按鈕（老師端） -->
        <template v-if="showExceptionButton">
          <button
            @click="$emit('add-exception')"
            class="px-4 py-2 rounded-lg bg-warning-500/20 text-warning-400 hover:bg-warning-500/30 transition-colors flex items-center gap-2 text-sm"
          >
            <Icon icon="warning" size="md" />
            {{ $t('schedule.leaveSwap') }}
          </button>
        </template>

        <!-- 匯出按鈕（老師端） -->
        <template v-if="showExportButton">
          <button
            @click="$emit('export')"
            class="px-4 py-2 rounded-lg bg-secondary-500/20 text-secondary-500 hover:bg-secondary-500/30 transition-colors flex items-center gap-2 text-sm"
          >
            <Icon icon="download" size="md" />
            {{ $t('schedule.export') }}
          </button>
        </template>
      </div>
    </div>

    <!-- 篩選提示（僅管理員模式） -->
    <div
      v-if="mode === 'admin' && (selectedTeacherId > 0 || selectedRoomId > 0)"
      class="mt-3 flex items-center gap-2 px-3 py-2 bg-primary-500/10 border border-primary-500/30 rounded-lg"
    >
      <span class="text-sm text-primary-400">{{ $t('filter.filtered') }}</span>
      <span v-if="selectedTeacherId > 0" class="inline-flex items-center gap-1 px-2 py-1 bg-primary-500/20 rounded text-sm">
        <span class="text-white">{{ selectedTeacherName }}</span>
        <button @click="$emit('clear-teacher-filter')" class="hover:text-primary-300">
          <Icon icon="close" size="xs" />
        </button>
      </span>
      <span v-if="selectedRoomId > 0" class="inline-flex items-center gap-1 px-2 py-1 bg-primary-500/20 rounded text-sm">
        <span class="text-white">{{ selectedRoomName }}</span>
        <button @click="$emit('clear-room-filter')" class="hover:text-primary-300">
          <Icon icon="close" size="xs" />
        </button>
      </span>
      <button
        v-if="selectedTeacherId > 0 || selectedRoomId > 0"
        @click="$emit('clear-all-filters')"
        class="ml-auto text-xs text-slate-400 hover:text-white transition-colors"
      >
        {{ $t('common.clearAll') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

// 引入圖示和提示元件
import Icon from '~/components/base/Icon.vue'
import HelpTooltip from '~/components/base/HelpTooltip.vue'

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
// DatePicker 相關狀態
// ============================================

const showDatePicker = ref(false)
const selectedDate = ref('')
const today = new Date()
const minDate = computed(() => {
  const d = new Date(today)
  d.setFullYear(d.getFullYear() - 2)
  return d.toISOString().split('T')[0]
})
const maxDate = computed(() => {
  const d = new Date(today)
  d.setFullYear(d.getFullYear() + 2)
  return d.toISOString().split('T')[0]
})

// 初始化今天的日期
const initializeToday = () => {
  selectedDate.value = today.toISOString().split('T')[0]
}

// 快速選擇今天
const quickSelectToday = () => {
  initializeToday()
  emit('select-date', new Date())
  showDatePicker.value = false
}

// 當日期改變時自動更新
const handleDateSelect = () => {
  // 選擇日期後自動確認並發送事件
  if (selectedDate.value) {
    emit('select-date', new Date(selectedDate.value))
    showDatePicker.value = false
  }
}

// 確認選擇
const confirmDateSelection = () => {
  if (selectedDate.value) {
    emit('select-date', new Date(selectedDate.value))
  }
  showDatePicker.value = false
}

// ============================================
// Emits 定義
// ============================================

const emit = defineEmits<{
  'change-week': [delta: number]
  'select-date': [date: Date]
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

// 初始化
initializeToday()
</script>
