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

    <!-- 課程和老師 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <SearchableSelect
          v-model="values.offering_id"
          :options="offeringOptions"
          label="課程"
          placeholder="請選擇課程"
          required
          :error="errors.offering_id"
        />
      </div>

      <div>
        <SearchableSelect
          v-model="values.teacher_id"
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
          v-model="values.room_id"
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
import RecurrencePicker from './RecurrencePicker.vue'
import SearchableSelect, { type SelectOption } from '~/components/Common/SearchableSelect.vue'

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
    offering_id: z.string().min(1, '請選擇課程'),
    teacher_id: z.union([z.string(), z.number(), z.null()]).optional(),
    room_id: z.union([z.string(), z.number(), z.null()]).optional(),
    start_time: z.string().min(1, '請選擇開始時間'),
    end_time: z.string().min(1, '請選擇結束時間'),
    duration: z.number().positive().min(1, '課程時長必須為正數'),
    weekdays: z.array(z.number()).min(1, '請至少選擇一個星期'),
    start_date: z.string().min(1, '請選擇開始日期'),
    end_date: z.string().optional(),
    skip_holiday: z.boolean().default(true),
  }

  return z.object(baseSchema)
}

// 初始化表單值
const getInitialValues = () => {
  if (props.editingRule) {
    return {
      name: props.editingRule.name || '',
      offering_id: String(props.editingRule.offering_id || ''),
      teacher_id: props.editingRule.teacher_id || null,
      room_id: props.editingRule.room_id || null,
      start_time: props.editingRule.start_time || '09:00',
      end_time: props.editingRule.end_time || '10:00',
      duration: props.editingRule.duration || 60,
      weekdays: [props.editingRule.weekday] || [1],
      start_date: props.editingRule.effective_range?.start_date?.split('T')[0] || formatDateToString(new Date()),
      end_date: props.editingRule.effective_range?.end_date?.split('T')[0] || '',
      skip_holiday: props.editingRule.skip_holiday ?? true,
    }
  }

  return {
    name: '',
    offering_id: '',
    teacher_id: null,
    room_id: null,
    start_time: '09:00',
    end_time: '10:00',
    duration: 60,
    weekdays: [1] as number[],
    start_date: formatDateToString(new Date()),
    end_date: '',
    skip_holiday: true,
  }
}

// 使用 vee-validate 的 useForm
const { handleSubmit, isSubmitting, errors, values, setFieldValue } = useForm({
  validationSchema: toTypedSchema(createValidationSchema()),
  initialValues: getInitialValues(),
})

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
      onChange: (val: string) => setFieldValue('offering_id', val),
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

// 監聽課程選擇，自動帶入預設老師和教室
watch(
  () => values.offering_id,
  (newOfferingId) => {
    // 編輯模式不自動帶入預設值
    if (isEditMode.value) return
    if (!newOfferingId) return

    const selectedOffering = offerings.value.find((o) => o.id === parseInt(newOfferingId))
    if (selectedOffering) {
      // 自動帶入預設老師（如果還沒有選老師）
      if (selectedOffering.default_teacher_id && !values.teacher_id) {
        setFieldValue('teacher_id', selectedOffering.default_teacher_id)
      }
      // 自動帶入預設教室（如果還沒有選教室）
      if (selectedOffering.default_room_id && !values.room_id) {
        setFieldValue('room_id', selectedOffering.default_room_id)
      }
      // 自動帶入名稱（如果還沒有填名稱）
      if (!values.name) {
        setFieldValue('name', selectedOffering.name)
      }
    }
  }
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
    skip_holiday: formValues.skip_holiday,
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

// 監聽編輯資料變化，更新表單值
watch(
  () => props.editingRule,
  (rule) => {
    if (rule) {
      setFieldValue('name', rule.name || '')
      setFieldValue('offering_id', String(rule.offering_id || ''))
      setFieldValue('teacher_id', rule.teacher_id || null)
      setFieldValue('room_id', rule.room_id || null)
      setFieldValue('start_time', rule.start_time || '09:00')
      setFieldValue('end_time', rule.end_time || '10:00')
      setFieldValue('duration', rule.duration || 60)
      setFieldValue('weekdays', [rule.weekday] || [1])
      setFieldValue(
        'start_date',
        rule.effective_range?.start_date?.split('T')[0] || formatDateToString(new Date())
      )
      setFieldValue('end_date', rule.effective_range?.end_date?.split('T')[0] || '')
      setFieldValue('skip_holiday', rule.skip_holiday ?? true)
    }
  },
  { immediate: true }
)
</script>
