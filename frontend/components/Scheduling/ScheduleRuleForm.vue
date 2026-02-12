<template>
  <form @submit.prevent="onFormSubmit" class="space-y-4">
    <!-- 空資料提示 -->
    <div
      v-if="offerings.length === 0 || rooms.length === 0 || teachers.length === 0"
      class="mb-4 p-4 rounded-lg bg-warning-500/10 border border-warning-500/30"
    >
      <p class="text-warning-500 text-sm">
        <span v-if="offerings.length === 0">尚未建立課程班別，請先至「資源管理」建立</span>
        <span v-if="rooms.length === 0">尚未建立教室</span>
        <span v-if="teachers.length === 0">尚未建立老師</span>
      </p>
    </div>

    <!-- 規則名稱 -->
    <div>
      <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
        規則名稱
      </label>
      <input
        :value="values.name"
        @input="(e) => setFieldValue('name', (e.target as HTMLInputElement).value)"
        type="text"
        placeholder="例：週一上午鋼琴課"
        class="input-field text-sm sm:text-base"
      />
      <span v-if="errors.name" class="text-critical-500 text-xs mt-1">
        {{ errors.name }}
      </span>
    </div>

    <!-- 課程狀態 (PLANNED/CONFIRMED) -->
    <div>
      <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
        課程狀態
      </label>
      <div class="relative">
        <select
          :value="values.status"
          @change="(e) => setFieldValue('status', (e.target as HTMLSelectElement).value)"
          class="input-field text-sm sm:text-base appearance-none cursor-pointer pr-10"
        >
          <option value="CONFIRMED">正式課 (已開成)</option>
          <option value="PLANNED">預計課 (佔位用)</option>
        </select>
        <!-- 自定義下拉箭頭 -->
        <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </div>
      </div>
      <p class="mt-2 text-xs text-slate-400">
        正式課：顯示實色背景。預計課：顯示斜紋虛線背景，標題加上「預計」字樣。兩者都會佔用時段進行衝突檢查。
      </p>
    </div>

    <!-- 課程和老師 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <SearchableSelect
          :model-value="values.offering_id"
          @update:model-value="(val) => setFieldValue('offering_id', val)"
          :options="offeringOptions"
          label="課程"
          placeholder="請選擇課程"
          required
          :error="errors.offering_id"
        />
      </div>

      <div>
        <SearchableSelect
          :model-value="values.teacher_id"
          @update:model-value="(val) => setFieldValue('teacher_id', val)"
          :options="teacherOptions"
          label="老師"
          placeholder="請選擇老師"
        />
      </div>
    </div>

    <!-- 教室和時間 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <SearchableSelect
          :model-value="values.room_id"
          @update:model-value="(val) => setFieldValue('room_id', val)"
          :options="roomOptions"
          label="教室"
          placeholder="請選擇教室"
        />
      </div>

      <div>
        <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
          開始時間
        </label>
        <input
          :value="values.start_time"
          @input="(e) => setFieldValue('start_time', (e.target as HTMLInputElement).value)"
          type="time"
          class="input-field text-sm sm:text-base"
        />
        <span v-if="errors.start_time" class="text-critical-500 text-xs mt-1">
          {{ errors.start_time }}
        </span>
      </div>
    </div>

    <!-- 結束時間和時長 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
          結束時間
        </label>
        <input
          :value="values.end_time"
          @input="(e) => setFieldValue('end_time', (e.target as HTMLInputElement).value)"
          type="time"
          class="input-field text-sm sm:text-base"
        />
        <span v-if="errors.end_time" class="text-critical-500 text-xs mt-1">
          {{ errors.end_time }}
        </span>
      </div>

      <div>
        <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
          課程時長（分鐘）
        </label>
        <input
          :value="values.duration"
          @input="(e) => setFieldValue('duration', Number((e.target as HTMLInputElement).value))"
          type="number"
          min="1"
          class="input-field text-sm sm:text-base"
        />
        <span v-if="errors.duration" class="text-critical-500 text-xs mt-1">
          {{ errors.duration }}
        </span>
      </div>
    </div>

    <!-- 重複星期 -->
    <RecurrencePicker
      v-model="weekdaysValue"
      :weekday-error="errors.weekdays"
      :weekday-label="'重複星期'"
      :weekday-help-text="'選擇此排課規則適用的星期幾。'"
      :weekday-usage-tips="['可選擇多個星期', '形成每週重複的排課']"
    />

    <!-- 例假日停課開關 -->
    <div class="mt-4">
      <label class="flex items-center cursor-pointer">
        <div class="relative inline-block w-12 h-7 align-middle select-none transition duration-200 ease-in-out">
          <input
            :checked="values.skip_holiday"
            @change="setFieldValue('skip_holiday', !values.skip_holiday)"
            type="checkbox"
            class="toggle-checkbox absolute block w-5 h-5 rounded-full bg-white border-4 appearance-none cursor-pointer transition-all duration-300 ease-in-out"
            :class="[
              values.skip_holiday ? 'left-6 border-primary-500' : 'left-0 border-slate-500'
            ]"
            style="top: 1px;"
          />
          <span
            class="toggle-label block overflow-hidden h-7 rounded-full transition-colors duration-300 ease-in-out"
            :class="values.skip_holiday ? 'bg-primary-500/30' : 'bg-slate-700'"
          ></span>
        </div>
        <span class="ml-3 text-sm sm:text-base text-slate-300 font-medium">
          例假日是否停課
        </span>
      </label>
      <p class="mt-1.5 ml-15 text-xs text-slate-400">
        開啟後，若遇一般例假日將自動停課
      </p>
    </div>

    <!-- 停課日期管理 -->
    <div class="mt-4">
      <button
        type="button"
        @click="openSuspendedDatesModal"
        class="w-full glass-btn py-2.5 rounded-xl font-medium text-sm sm:text-base flex items-center justify-center gap-2"
        :class="suspendedDatesCount > 0 ? 'border-warning-500/50 text-warning-500' : 'text-slate-300'"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
        </svg>
        <span>停課日期管理</span>
        <span v-if="suspendedDatesCount > 0" class="bg-warning-500 text-white text-xs px-2 py-0.5 rounded-full">
          {{ suspendedDatesCount }}
        </span>
      </button>
      <p class="mt-1.5 text-xs text-slate-400">
        設定特定日期停課，例如國定假日、補課日等
      </p>
    </div>

    <!-- 編輯模式的日期欄位 -->
    <div v-if="isEditMode" class="mb-4 p-3 rounded-lg bg-slate-800/50 border border-slate-700/50">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
            開始日期
          </label>
          <input
            :value="values.start_date"
            @input="(e) => setFieldValue('start_date', (e.target as HTMLInputElement).value)"
            type="date"
            class="input-field text-sm sm:text-base"
          />
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
            結束日期
          </label>
          <input
            :value="values.end_date"
            @input="(e) => setFieldValue('end_date', (e.target as HTMLInputElement).value)"
            type="date"
            class="input-field text-sm sm:text-base"
          />
        </div>
      </div>
      <p class="text-xs text-slate-400 mt-2">
        <span class="text-warning-500">💡 提示：</span>如只修改課程內容（老師、教室、時間），日期可留空以保留現有日期範圍。
      </p>
    </div>

    <!-- 新增模式的必填日期欄位 -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
          開始日期
        </label>
        <input
          :value="values.start_date"
          @input="(e) => setFieldValue('start_date', (e.target as HTMLInputElement).value)"
          type="date"
          class="input-field text-sm sm:text-base"
        />
        <span v-if="errors.start_date" class="text-critical-500 text-xs mt-1">
          {{ errors.start_date }}
        </span>
      </div>

      <div>
        <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
          結束日期
        </label>
        <input
          :value="values.end_date"
          @input="(e) => setFieldValue('end_date', (e.target as HTMLInputElement).value)"
          type="date"
          class="input-field text-sm sm:text-base"
        />
      </div>
    </div>

    <!-- 提交按鈕 -->
    <div class="flex gap-3 pt-2">
      <button
        type="button"
        @click="$emit('cancel')"
        class="flex-1 glass-btn py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
      >
        取消
      </button>
      <button
        type="submit"
        :disabled="isSubmitting || validationLoading"
        class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
      >
        {{ validationLoading ? '驗證中...' : (isSubmitting ? '儲存中...' : (isEditMode ? '儲存修改' : '新增')) }}
      </button>
    </div>
  </form>

  <!-- 停課日期管理 Modal -->
  <Teleport to="body">
    <div
      v-if="showSuspendedDatesModal"
      class="fixed inset-0 z-[1100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
      @click.self="closeSuspendedDatesModal"
    >
      <div class="bg-slate-800 rounded-2xl w-full max-w-3xl max-h-[80vh] overflow-hidden shadow-2xl border border-slate-700">
        <!-- Modal Header -->
        <div class="p-4 border-b border-slate-700 flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-white">停課日期管理</h3>
            <p class="text-sm text-slate-400 mt-1">
              勾選要在 {{ values.start_date }} ~ {{ values.end_date || '無限期' }} 期間停課的日期
            </p>
          </div>
          <button
            @click="closeSuspendedDatesModal"
            class="p-2 hover:bg-slate-700 rounded-lg transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>

        <!-- Modal Content -->
        <div class="p-4 overflow-y-auto max-h-[60vh]">
          <!-- 篩選器 -->
          <div class="mb-4 flex flex-wrap gap-2 items-center">
            <button
              @click="filterMode = 'all'"
              class="px-3 py-1.5 text-sm rounded-lg transition-colors"
              :class="filterMode === 'all' ? 'bg-primary-500 text-white' : 'bg-slate-700 text-slate-300 hover:bg-slate-600'"
            >
              全部 ({{ allDates.length }})
            </button>
            <button
              @click="filterMode = 'suspended'"
              class="px-3 py-1.5 text-sm rounded-lg transition-colors"
              :class="filterMode === 'suspended' ? 'bg-warning-500 text-white' : 'bg-slate-700 text-slate-300 hover:bg-slate-600'"
            >
              已選停課 ({{ suspendedDates.length }})
            </button>
            <button
              @click="filterMode = 'available'"
              class="px-3 py-1.5 text-sm rounded-lg transition-colors"
              :class="filterMode === 'available' ? 'bg-green-500 text-white' : 'bg-slate-700 text-slate-300 hover:bg-slate-600'"
            >
              可選日期 ({{ allDates.length - suspendedDates.length }})
            </button>

            <!-- 隱藏例假日開關 -->
            <label
              v-if="centerHolidays.length > 0"
              class="ml-auto flex items-center gap-2 cursor-pointer"
            >
              <span class="text-sm text-slate-400">隱藏例假日</span>
              <div class="relative inline-block w-10 h-6 align-middle select-none transition duration-200 ease-in-out">
                <input
                  v-model="hideHolidays"
                  type="checkbox"
                  class="toggle-checkbox absolute block w-4 h-4 rounded-full bg-white border-2 appearance-none cursor-pointer transition-all duration-200 ease-in-out"
                  :class="hideHolidays ? 'left-5 border-primary-500' : 'left-0 border-slate-500'"
                  style="top: 2px;"
                />
                <span
                  class="toggle-label block overflow-hidden h-5 rounded-full transition-colors duration-200 ease-in-out"
                  :class="hideHolidays ? 'bg-primary-500/30' : 'bg-slate-700'"
                ></span>
              </div>
            </label>
          </div>

          <!-- 群組操作按鈕 -->
          <div class="mb-4 flex gap-2">
            <button
              v-if="filterMode !== 'suspended'"
              @click="selectAllVisible"
              class="px-3 py-1.5 text-sm bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg transition-colors"
            >
              全選可見日期
            </button>
            <button
              v-if="filterMode !== 'available'"
              @click="deselectAllVisible"
              class="px-3 py-1.5 text-sm bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg transition-colors"
            >
              取消全選可見日期
            </button>
          </div>

          <!-- 日期列表 -->
          <div v-if="filteredDates.length > 0" class="space-y-4">
            <div
              v-for="(dates, monthKey) in groupedDates"
              :key="monthKey"
              class="bg-slate-700/50 rounded-lg p-3"
            >
              <h4 class="text-sm font-medium text-slate-300 mb-2">{{ monthKey }}</h4>
              <div class="grid grid-cols-4 sm:grid-cols-6 gap-2">
                <label
                  v-for="date in dates"
                  :key="date.value"
                  class="flex items-center gap-2 p-2 rounded-lg cursor-pointer transition-all relative"
                  :class="isDateSuspended(date.value)
                    ? 'bg-warning-500/20 border border-warning-500/50'
                    : date.isHoliday
                      ? 'bg-blue-500/10 border border-blue-500/30 hover:bg-blue-500/20'
                      : 'bg-slate-700 hover:bg-slate-600 border border-transparent'"
                >
                  <input
                    type="checkbox"
                    :checked="isDateSuspended(date.value)"
                    @change="toggleSuspendedDate(date.value)"
                    class="w-4 h-4 rounded border-slate-500 text-warning-500 focus:ring-warning-500 z-10"
                  />
                  <div class="flex flex-col min-w-0 z-10">
                    <span class="text-xs text-slate-300">{{ date.weekday }}</span>
                    <span class="text-sm font-medium" :class="isDateSuspended(date.value) ? 'text-warning-500' : 'text-white'">
                      {{ date.day }}
                    </span>
                  </div>

                  <!-- 假日標記 -->
                  <div
                    v-if="date.isHoliday"
                    class="absolute top-1 right-1 flex items-center gap-0.5"
                    :title="date.holidayName || '例假日'"
                  >
                    <span class="text-[10px] font-bold text-blue-400">H</span>
                  </div>

                  <!-- 假日名稱顯示 -->
                  <div
                    v-if="date.isHoliday"
                    class="absolute -bottom-1 left-2 right-2 flex justify-center"
                  >
                    <span class="text-[9px] text-blue-400 bg-slate-800/90 px-1 rounded truncate max-w-full" :title="date.holidayName || ''">
                      {{ date.holidayName || '' }}
                    </span>
                  </div>
                </label>
              </div>
            </div>
          </div>

          <!-- 無日期提示 -->
          <div v-else class="text-center py-8 text-slate-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto mb-3 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <p v-if="!values.start_date || !values.end_date">
              請先設定開始日期和結束日期
            </p>
            <p v-else-if="allDates.length === 0">
              在指定的日期範圍內沒有符合的重複星期
            </p>
            <p v-else-if="hideHolidays && allDates.filter(d => !d.isHoliday).length === 0">
              所有日期皆為例假日，請取消隱藏例假日
            </p>
            <p v-else>
              沒有符合篩選條件的日期
            </p>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-4 border-t border-slate-700 flex items-center justify-between">
          <div class="text-sm text-slate-400">
            <span class="text-warning-500 font-medium">{{ suspendedDatesCount }}</span> 個停課日期
            <span v-if="hideHolidays" class="ml-2 text-blue-400">
              (已隱藏 {{ allDates.filter(d => d.isHoliday).length }} 個例假日)
            </span>
          </div>
          <div class="flex gap-2">
            <button
              @click="clearAllSuspendedDates"
              class="px-4 py-2 text-sm text-slate-300 hover:text-white transition-colors"
              :disabled="suspendedDatesCount === 0"
            >
              清除全部
            </button>
            <button
              @click="closeSuspendedDatesModal"
              class="px-4 py-2 text-sm bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
            >
              確定
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.toggle-checkbox:checked {
  @apply border-primary-500;
}

.toggle-checkbox:not(:checked) {
  @apply border-slate-500;
}

.toggle-checkbox:checked + .toggle-label {
  @apply bg-primary-500/30;
}

.toggle-checkbox:not(:checked) + .toggle-label {
  @apply bg-slate-700;
}
</style>

<script setup lang="ts">
import { z } from 'zod'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { formatDateToString } from '~/composables/useTaiwanTime'
import { alertWarning } from '~/composables/useAlert'
import { useCenterId } from '~/composables/useCenterId'
import { useApi } from '~/composables/useApi'
import RecurrencePicker from './RecurrencePicker.vue'
import SearchableSelect, { type SelectOption } from '~/components/Common/SearchableSelect.vue'

// 全域預設課程時長（分鐘）
const DEFAULT_COURSE_DURATION = 60

// Props
interface Props {
  editingRule?: any | null
  updateMode?: string
  validationLoading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  editingRule: null,
  updateMode: 'ALL',
  validationLoading: false,
})

const emit = defineEmits<{
  cancel: []
  submit: [formData: Record<string, unknown>, updateMode: string]
  validate: [formData: Record<string, unknown>] // 用於父元件執行驗證
}>()

// 計算屬性
const isEditMode = computed(() => !!props.editingRule)

// 停課日期管理相關
const showSuspendedDatesModal = ref(false)
const suspendedDates = ref<string[]>([])
const filterMode = ref<'all' | 'suspended' | 'available'>('all')

// 假日管理相關
const centerHolidays = ref<Array<{
  id: number
  date: string
  name: string
  force_cancel: boolean
}>>([])
const hideHolidays = ref(false)
const holidaysLoading = ref(false)

// 星期對照表
const weekdayNames = ['日', '一', '二', '三', '四', '五', '六']

// 取得中心假日資料
async function fetchHolidays() {
  const { get: apiGet } = useApi()
  const { getCenterId } = useCenterId()
  const centerId = getCenterId()

  if (!centerId) return

  holidaysLoading.value = true
  try {
    const response = await apiGet<{ data: Array<{ id: number; date: string; name: string; force_cancel: boolean }> }>(
      `/admin/centers/${centerId}/holidays`,
      {
        start_date: values.start_date,
        end_date: values.end_date || values.start_date,
      }
    )
    centerHolidays.value = response.data || []
  } catch (error) {
    console.error('取得假日資料失敗:', error)
    centerHolidays.value = []
  } finally {
    holidaysLoading.value = false
  }
}

// 計算所有可能的上课日期
const allDates = computed(() => {
  const startDate = values.start_date
  const endDate = values.end_date
  const weekdays = values.weekdays || []

  if (!startDate || !endDate || weekdays.length === 0) {
    return []
  }

  const dates: Array<{
    value: string
    day: string
    weekday: string
    monthKey: string
    isHoliday: boolean
    holidayName: string | null
  }> = []

  const start = new Date(startDate)
  const end = new Date(endDate)
  const current = new Date(start)

  // 設定為當天開始
  current.setHours(0, 0, 0, 0)
  end.setHours(23, 59, 59, 999)

  // 建立假日快速查找 Map
  const holidayMap = new Map<string, { name: string; force_cancel: boolean }>()
  centerHolidays.value.forEach((holiday) => {
    holidayMap.set(holiday.date, {
      name: holiday.name,
      force_cancel: holiday.force_cancel,
    })
  })

  while (current <= end) {
    const dayOfWeek = current.getDay()
    if (weekdays.includes(dayOfWeek)) {
      const dateStr = formatDateToString(current)
      const monthKey = `${current.getFullYear()}年${current.getMonth() + 1}月`
      const holidayInfo = holidayMap.get(dateStr)

      dates.push({
        value: dateStr,
        day: `${current.getDate()}`,
        weekday: weekdayNames[dayOfWeek],
        monthKey,
        isHoliday: !!holidayInfo,
        holidayName: holidayInfo?.name || null,
      })
    }
    current.setDate(current.getDate() + 1)
  }

  return dates
})

// 過濾後的日期
const filteredDates = computed(() => {
  // 先根據篩選模式過濾
  let dates = allDates.value
  switch (filterMode.value) {
    case 'suspended':
      dates = dates.filter((d) => suspendedDates.value.includes(d.value))
      break
    case 'available':
      dates = dates.filter((d) => !suspendedDates.value.includes(d.value))
      break
  }

  // 如果開啟隱藏假日，則排除假日
  if (hideHolidays.value) {
    dates = dates.filter((d) => !d.isHoliday)
  }

  return dates
})

// 按月份分組
const groupedDates = computed(() => {
  const groups: Record<string, typeof filteredDates.value> = {}
  for (const date of filteredDates.value) {
    if (!groups[date.monthKey]) {
      groups[date.monthKey] = []
    }
    groups[date.monthKey].push(date)
  }
  return groups
})

// 已選擇的停課日期數量
const suspendedDatesCount = computed(() => suspendedDates.value.length)

// 檢查日期是否已選擇停課
function isDateSuspended(date: string): boolean {
  return suspendedDates.value.includes(date)
}

// 切換停課日期
function toggleSuspendedDate(date: string) {
  const index = suspendedDates.value.indexOf(date)
  if (index > -1) {
    suspendedDates.value.splice(index, 1)
  } else {
    suspendedDates.value.push(date)
  }
  // 同步到表單值
  setFieldValue('suspended_dates', [...suspendedDates.value])
}

// 全選可見日期
function selectAllVisible() {
  for (const date of filteredDates.value) {
    if (!suspendedDates.value.includes(date.value)) {
      suspendedDates.value.push(date.value)
    }
  }
  setFieldValue('suspended_dates', [...suspendedDates.value])
}

// 取消全選可見日期
function deselectAllVisible() {
  for (const date of filteredDates.value) {
    const index = suspendedDates.value.indexOf(date.value)
    if (index > -1) {
      suspendedDates.value.splice(index, 1)
    }
  }
  setFieldValue('suspended_dates', [...suspendedDates.value])
}

// 清除全部停課日期
function clearAllSuspendedDates() {
  suspendedDates.value = []
  setFieldValue('suspended_dates', [])
}

// 開啟停課日期 Modal
async function openSuspendedDatesModal() {
  if (!values.start_date || !values.end_date) {
    alertWarning('請先設定開始日期和結束日期')
    return
  }
  if (!values.weekdays || values.weekdays.length === 0) {
    alertWarning('請先選擇重複星期')
    return
  }
  // 取得假日資料
  await fetchHolidays()
  showSuspendedDatesModal.value = true
}

// 關閉停課日期 Modal
function closeSuspendedDatesModal() {
  showSuspendedDatesModal.value = false
}

// 從共享緩存取得資料
const { resourceCache } = useResourceCache()
const offerings = computed(() => resourceCache.value.offerings)
const teachers = computed(() => Array.from(resourceCache.value.teachers.values()))
const rooms = computed(() => Array.from(resourceCache.value.rooms.values()))

// 轉換為 SearchableSelect 選項格式
const offeringOptions = computed<SelectOption[]>(() =>
  offerings.value.map(o => ({
    id: o.id,
    name: o.name || `班別 #${o.id}`
  }))
)

const teacherOptions = computed<SelectOption[]>(() =>
  teachers.value.map(t => ({
    id: t.id,
    name: t.name
  }))
)

const roomOptions = computed<SelectOption[]>(() =>
  rooms.value.map(r => ({
    id: r.id,
    name: r.name
  }))
)

// Zod 驗證 Schema
const createValidationSchema = () => {
  const baseSchema = {
    name: z.string().optional(),
    offering_id: z.union([z.string(), z.number(), z.null()]).refine((val) => {
      // 確保值不為空（字串長度 > 0 或數字 > 0）
      if (val === null) return false
      if (typeof val === 'string') return val.length > 0
      return val > 0
    }, { message: '請選擇課程' }),
    teacher_id: z.union([z.string(), z.number(), z.null()]).optional(),
    room_id: z.union([z.string(), z.number(), z.null()]).optional(),
    start_time: z.string().min(1, '請選擇開始時間'),
    end_time: z.string().min(1, '請選擇結束時間'),
    duration: z.number().positive().min(1, '課程時長必須為正數'),
    weekdays: z.array(z.number()).min(1, '請至少選擇一個星期'),
    start_date: z.string().min(1, '請選擇開始日期'),
    end_date: z.string().optional(),
    status: z.enum(['PLANNED', 'CONFIRMED']).default('CONFIRMED'),
    skip_holiday: z.boolean().default(true),
    suspended_dates: z.array(z.string()).default([]),
  }

  return z.object(baseSchema)
}

// 初始化表單值
const getInitialValues = () => {
  if (props.editingRule) {
    // 解析 suspended_dates（可能是 JSON 字串或陣列）
    let suspendedDatesData: string[] = []
    if (props.editingRule.suspended_dates) {
      if (Array.isArray(props.editingRule.suspended_dates)) {
        suspendedDatesData = props.editingRule.suspended_dates
      } else if (typeof props.editingRule.suspended_dates === 'string') {
        try {
          suspendedDatesData = JSON.parse(props.editingRule.suspended_dates)
        } catch {
          suspendedDatesData = []
        }
      }
    }
    // 同步到組件狀態
    suspendedDates.value = suspendedDatesData

    return {
      name: props.editingRule.name || '',
      offering_id: props.editingRule.offering_id || null,
      teacher_id: props.editingRule.teacher_id || null,
      room_id: props.editingRule.room_id || null,
      start_time: props.editingRule.start_time || '09:00',
      end_time: props.editingRule.end_time || '10:00',
      duration: props.editingRule.duration || DEFAULT_COURSE_DURATION,
      weekdays: props.editingRule.weekdays || (props.editingRule.weekday !== undefined ? [props.editingRule.weekday] : [1]),
      start_date: props.editingRule.effective_range?.start_date?.split(/[T ]/)[0] || formatDateToString(new Date()),
      end_date: props.editingRule.effective_range?.end_date?.split(/[T ]/)[0] || '',
      status: props.editingRule.status || 'CONFIRMED',
      skip_holiday: props.editingRule.skip_holiday ?? true,
      suspended_dates: suspendedDatesData,
    }
  }

  return {
    name: '',
    offering_id: null,
    teacher_id: null,
    room_id: null,
    start_time: '09:00',
    end_time: '10:00',
    duration: DEFAULT_COURSE_DURATION,
    weekdays: [1] as number[],
    start_date: formatDateToString(new Date()),
    end_date: '',
    status: 'CONFIRMED',
    skip_holiday: true,
    suspended_dates: [] as string[],
  }
}

// 使用 vee-validate 的 useForm
const { handleSubmit, isSubmitting, errors, values, setFieldValue, setValues } = useForm({
  validationSchema: toTypedSchema(createValidationSchema()),
  initialValues: getInitialValues(),
}) as any

// 監聽編輯資料變化，使用 setValues 正確更新所有表單值
watch(
  () => props.editingRule,
  (rule) => {
    if (rule) {
      // 解析 suspended_dates
      let suspendedDatesData: string[] = []
      if (rule.suspended_dates) {
        if (Array.isArray(rule.suspended_dates)) {
          suspendedDatesData = rule.suspended_dates
        } else if (typeof rule.suspended_dates === 'string') {
          try {
            suspendedDatesData = JSON.parse(rule.suspended_dates)
          } catch {
            suspendedDatesData = []
          }
        }
      }
      suspendedDates.value = suspendedDatesData

      // 使用 setValues 一次更新所有欄位
      setValues({
        name: rule.name || '',
        offering_id: rule.offering_id || null,
        teacher_id: rule.teacher_id || null,
        room_id: rule.room_id || null,
        start_time: rule.start_time || '09:00',
        end_time: rule.end_time || '10:00',
        duration: rule.duration || DEFAULT_COURSE_DURATION,
        weekdays: rule.weekdays || (rule.weekday !== undefined ? [rule.weekday] : [1]),
        start_date: rule.effective_range?.start_date?.split(/[T ]/)[0] || formatDateToString(new Date()),
        end_date: rule.effective_range?.end_date?.split(/[T ]/)[0] || null,
        status: rule.status || 'CONFIRMED',
        skip_holiday: rule.skip_holiday ?? true,
        suspended_dates: suspendedDatesData,
      })
    }
  },
  { immediate: true }
)

// 建立欄位屬性物件（用於 v-bind）
const fieldAttrs = computed(() => {
  return {
    name: {
      value: values.name,
      onChange: (val: string) => setFieldValue('name', val),
      error: errors.name,
    },
    offering_id: {
      value: values.offering_id,
      onChange: (val: any) => setFieldValue('offering_id', val),
      error: errors.offering_id,
    },
    teacher_id: {
      value: values.teacher_id,
      onChange: (val: any) => setFieldValue('teacher_id', val),
      error: errors.teacher_id,
    },
    room_id: {
      value: values.room_id,
      onChange: (val: any) => setFieldValue('room_id', val),
      error: errors.room_id,
    },
    start_time: {
      value: values.start_time,
      onChange: (val: string) => setFieldValue('start_time', val),
      error: errors.start_time,
    },
    end_time: {
      value: values.end_time,
      onChange: (val: string) => setFieldValue('end_time', val),
      error: errors.end_time,
    },
    duration: {
      value: values.duration,
      onChange: (val: number) => setFieldValue('duration', val),
      error: errors.duration,
    },
    weekdays: {
      value: values.weekdays,
      onChange: (val: number[]) => setFieldValue('weekdays', val),
      error: errors.weekdays,
    },
    start_date: {
      value: values.start_date,
      onChange: (val: string) => setFieldValue('start_date', val),
      error: errors.start_date,
    },
    end_date: {
      value: values.end_date,
      onChange: (val: string) => setFieldValue('end_date', val),
      error: errors.end_date,
    },
  }
})

// weekdays 的值（用於 UI 顯示）
const weekdaysValue = computed({
  get: () => values.weekdays as number[],
  set: (val) => setFieldValue('weekdays', val),
})

// 追蹤是否為初始載入（用於區分初始化和用戶操作）
const isInitialLoad = ref(true)

// 監聽課程選擇，自動帶入預設老師、教室、開始結束時間和時長
watch(
  () => values.offering_id,
  (newOfferingId, oldOfferingId) => {
    // 如果沒有選擇課程，不處理
    if (!newOfferingId) return

    const selectedOffering = offerings.value.find((o) => String(o.id) === String(newOfferingId))
    if (!selectedOffering) return

    // 判斷是否為用戶手動更換課程（非初始載入）
    const isUserChangedCourse = oldOfferingId !== undefined && oldOfferingId !== newOfferingId

    // 編輯模式下：
    // - 初始載入時跳過（保持原有的編輯資料）
    // - 用戶手動更換課程時，帶入新課程的預設時間和時長
    if (isEditMode.value && !isUserChangedCourse) {
      // 這是初始載入，跳過自動帶入
      return
    }

    // 自動帶入預設開始時間
    if (selectedOffering.default_start_time && selectedOffering.default_start_time !== '') {
      setFieldValue('start_time', selectedOffering.default_start_time)
    }
    // 自動帶入預設結束時間
    if (selectedOffering.default_end_time && selectedOffering.default_end_time !== '') {
      setFieldValue('end_time', selectedOffering.default_end_time)
    }
    // 自動帶入課程時長（從關聯課程取得，如果有的話）
    if (selectedOffering.course_duration && selectedOffering.course_duration > 0) {
      setFieldValue('duration', selectedOffering.course_duration)
    }
    // 新增模式：自動帶入預設老師（如果還沒有選老師）
    if (!isEditMode.value && selectedOffering.default_teacher_id && !values.teacher_id) {
      setFieldValue('teacher_id', selectedOffering.default_teacher_id)
    }
    // 新增模式：自動帶入預設教室（如果還沒有選教室）
    if (!isEditMode.value && selectedOffering.default_room_id && !values.room_id) {
      setFieldValue('room_id', selectedOffering.default_room_id)
    }
    // 新增模式：自動帶入名稱（如果還沒有填名稱）
    if (!isEditMode.value && !values.name) {
      setFieldValue('name', selectedOffering.name)
    }
  }
)

// 監聽初始載入完成
watch(
  () => props.editingRule,
  () => {
    // 當 editingRule 載入完成後，標記為非初始狀態
    if (props.editingRule) {
      nextTick(() => {
        isInitialLoad.value = false
      })
    }
  },
  { immediate: true }
)

// 提交處理
const onFormSubmit = handleSubmit(async (formValues) => {
  const data: Record<string, unknown> = {
    name: formValues.name,
    offering_id: parseInt(formValues.offering_id as string),
    start_time: formValues.start_time,
    end_time: formValues.end_time,
    duration: formValues.duration,
    weekdays: formValues.weekdays,
    start_date: formValues.start_date,
    end_date: formValues.end_date || null,
    status: formValues.status,
    skip_holiday: formValues.skip_holiday,
    suspended_dates: formValues.suspended_dates || [],
  }

  // 只有當有選擇老師時才傳送
  if (formValues.teacher_id) {
    data.teacher_id = formValues.teacher_id
  }

  // 只有當有選擇教室時才傳送
  if (formValues.room_id) {
    data.room_id = formValues.room_id
  }

  // 編輯模式：處理日期欄位
  if (isEditMode.value) {
    // 如果日期為空，從 data 中移除，讓後端保留現有值
    if (!data.start_date) {
      delete data.start_date
    }
    if (!data.end_date) {
      delete data.end_date
    }
    emit('submit', data, props.updateMode || 'ALL')
  } else {
    // 新增模式
    emit('submit', data, 'ALL')
  }
})

</script>
