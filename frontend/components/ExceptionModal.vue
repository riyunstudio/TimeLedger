<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" @click.self="$emit('close')">
      <div class="glass-card w-full max-w-lg max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between p-4 border-b border-white/10">
          <h3 class="text-lg font-semibold text-white">{{ isEdit ? '編輯例外申請' : '新增例外申請' }}</h3>
          <button @click="$emit('close')" class="p-2 rounded-lg hover:bg-white/10">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="p-4 space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">中心</label>
            <select v-model="form.center_id" class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500" required>
              <option value="">選擇中心</option>
              <option v-for="center in centers" :key="center.center_id" :value="center.center_id">
                {{ center.center_name }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">課程時段</label>
            <select v-model="form.rule_id" class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500" required :disabled="loadingRules">
              <option value="">選擇課程時段</option>
              <option v-for="rule in displayScheduleRules" :key="rule.id" :value="rule.id">
                {{ formatRuleDisplay(rule) }}
              </option>
            </select>
            <p v-if="loadingRules" class="text-xs text-slate-500 mt-1">載入中...</p>
            <p v-else-if="form.center_id && displayScheduleRules.length === 0" class="text-xs text-slate-500 mt-1">該中心暫無您的課程</p>
          </div>

          <!-- 選擇具體日期 -->
          <div v-if="form.rule_id">
            <label class="block text-sm font-medium text-slate-300 mb-1">調整日期</label>
            <input
              type="date"
              v-model="form.original_date"
              :min="today"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500"
              required
            />
            <p class="text-xs text-slate-500 mt-1">選擇要調整的具體日期</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">申請類型</label>
            <div class="flex gap-3 flex-wrap">
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" v-model="form.type" value="CANCEL" class="accent-primary-500" />
                <span class="text-white">停課</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" v-model="form.type" value="RESCHEDULE" class="accent-primary-500" />
                <span class="text-white">改期</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" v-model="form.type" value="REPLACE_TEACHER" class="accent-primary-500" />
                <span class="text-white">找代課</span>
              </label>
            </div>
          </div>

          <!-- 改期時間選擇 -->
          <div v-if="form.type === 'RESCHEDULE'">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-slate-300 mb-1">新開始時間</label>
                <input type="datetime-local" v-model="form.new_start_at" class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500" />
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-300 mb-1">新結束時間</label>
                <input type="datetime-local" v-model="form.new_end_at" class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500" />
              </div>
            </div>
          </div>

          <!-- 代課老師選擇 -->
          <div v-if="form.type === 'REPLACE_TEACHER'" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-2">代課方式</label>
              <div class="space-y-2">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="replaceTeacherMode" value="center" class="accent-primary-500" />
                  <span class="text-white">由中心幫我找代課老師</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="replaceTeacherMode" value="manual" class="accent-primary-500" />
                  <span class="text-white">我已找到代課老師</span>
                </label>
              </div>
            </div>

            <!-- 手動輸入代課老師姓名 -->
            <div v-if="replaceTeacherMode === 'manual'">
              <label class="block text-sm font-medium text-slate-300 mb-1">代課老師姓名</label>
              <input
                v-model="form.new_teacher_name"
                type="text"
                placeholder="請輸入代課老師姓名"
                class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white placeholder-slate-500 focus:outline-none focus:border-primary-500"
              />
              <p class="text-xs text-slate-500 mt-1">請輸入您已聯繫好的代課老師姓名</p>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">原因</label>
            <textarea v-model="form.reason" rows="3" class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500 resize-none" required></textarea>
          </div>

          <div class="flex gap-3 pt-4">
            <button type="button" @click="$emit('close')" class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors">
              取消
            </button>
            <button type="submit" :disabled="loading" class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50">
              {{ loading ? '提交中...' : '提交申請' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import type { ScheduleException } from '~/types'
import { alertError } from '~/composables/useAlert'
import { getTodayString, formatDateToString } from '~/composables/useTaiwanTime'

interface ScheduleRuleData {
  id: number
  title: string
  weekday: number
  weekday_text: string
  start_time: string
  end_time: string
  effective_start_date: string
  effective_end_date: string
}

interface Props {
  exception?: ScheduleException
  centers: Array<{ center_id: number; center_name: string }>
  scheduleRules?: ScheduleRuleData[]
  prefill?: {
    rule_id?: number
    center_id?: number
    course_name?: string
    original_date?: string
    original_time?: string
  }
}

const props = defineProps<Props>()
const emit = defineEmits(['close', 'submit'])

const teacherStore = useTeacherStore()
const loading = ref(false)
const localScheduleRules = ref<ScheduleRuleData[]>([])
const loadingRules = ref(false)

const isEdit = computed(() => !!props.exception)

const form = reactive({
  center_id: props.exception?.center_id || 0,
  rule_id: props.exception?.rule_id || 0,
  original_date: props.exception?.original_date || '',
  type: props.exception?.type || 'CANCEL' as 'CANCEL' | 'RESCHEDULE' | 'REPLACE_TEACHER',
  new_start_at: props.exception?.new_start_at || '',
  new_end_at: props.exception?.new_end_at || '',
  new_teacher_id: 0, // 不再使用
  new_teacher_name: props.exception?.new_teacher_name || '',
  reason: props.exception?.reason || '',
})

// 代課模式：'center' = 由中心找, 'manual' = 手動輸入
const replaceTeacherMode = ref<'center' | 'manual'>('center')

// 處理預填資料
onMounted(async () => {
  if (props.prefill) {
    const { center_id, rule_id, original_date, original_time } = props.prefill

    if (center_id) {
      form.center_id = center_id
    }

    if (original_date) {
      form.original_date = original_date
    }

    // 如果有 center_id，載入課程列表
    if (form.center_id && form.center_id > 0) {
      await fetchScheduleRules(form.center_id)

      // 如果有 rule_id，選中對應課程
      if (rule_id) {
        form.rule_id = rule_id
      }
    }
  }
})

// 今天日期（用於日期選擇的最小值）
const today = computed(() => {
  return getTodayString()
})

// 監聽中心選擇變化，載入該中心的課程
watch(() => form.center_id, async (newCenterId) => {
  // 清空課程選擇
  form.rule_id = 0
  form.original_date = ''
  localScheduleRules.value = []

  if (newCenterId && newCenterId > 0) {
    await fetchScheduleRules(newCenterId)
  }
})

// 監聽類型變化，清空相關欄位
watch(() => form.type, (newType) => {
  if (newType === 'CANCEL') {
    form.new_start_at = ''
    form.new_end_at = ''
    form.new_teacher_name = ''
    replaceTeacherMode.value = 'center'
  } else if (newType === 'RESCHEDULE') {
    form.new_teacher_name = ''
    replaceTeacherMode.value = 'center'
  } else if (newType === 'REPLACE_TEACHER') {
    form.new_start_at = ''
    form.new_end_at = ''
  }
})

// 如果有傳入 scheduleRules，使用傳入的；否則使用本地的
const displayScheduleRules = computed(() => {
  if (props.scheduleRules && props.scheduleRules.length > 0) {
    return props.scheduleRules
  }
  return localScheduleRules.value
})

const fetchScheduleRules = async (centerId: number) => {
  try {
    loadingRules.value = true
    const api = useApi()
    const response = await api.get<{ code: number; datas: ScheduleRuleData[] }>(`/teacher/me/centers/${centerId}/schedule-rules`)
    if (response.code === 0 && response.datas) {
      localScheduleRules.value = response.datas
    }
  } catch (error) {
    console.error('Failed to fetch schedule rules:', error)
  } finally {
    loadingRules.value = false
  }
}

// 生成課程顯示文字
const formatRuleDisplay = (rule: ScheduleRuleData): string => {
  const dateRange = rule.effective_start_date || rule.effective_end_date
    ? `${rule.effective_start_date || '-'} ~ ${rule.effective_end_date || '-'}`
    : ''

  return `${rule.title} (${rule.weekday_text} ${rule.start_time}-${rule.end_time})${dateRange ? ' ' + dateRange : ''}`
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW')
}

const handleSubmit = async () => {
  loading.value = true
  try {
    // 確保日期格式正確 (ISO 8601 台灣時區)
    const originalDate = form.original_date
      ? formatDateToString(new Date(form.original_date))
      : ''

    const submitData: any = {
      center_id: form.center_id,
      rule_id: form.rule_id,
      original_date: originalDate,
      type: form.type,
      reason: form.reason,
    }

    // 根據類型添加相應欄位
    if (form.type === 'RESCHEDULE') {
      if (form.new_start_at) {
        submitData.new_start_at = formatDateToString(new Date(form.new_start_at))
      }
      if (form.new_end_at) {
        submitData.new_end_at = formatDateToString(new Date(form.new_end_at))
      }
    } else if (form.type === 'REPLACE_TEACHER') {
      // 如果選擇手動輸入代課老師，才傳送名字
      if (replaceTeacherMode.value === 'manual' && form.new_teacher_name && form.new_teacher_name.trim()) {
        submitData.new_teacher_name = form.new_teacher_name.trim()
      }
      // 如果選擇由中心找，不傳送 new_teacher_name
    }

    await teacherStore.createException(submitData)
    emit('submit')
    emit('close')
  } catch (error: any) {
    console.error('Failed to create exception:', error)
    // 顯示錯誤彈窗，從錯誤響應中提取訊息
    const errorMessage = error?.response?._data?.message || error?.message || '操作失敗，請稍後再試'
    await alertError(errorMessage)
  } finally {
    loading.value = false
  }
}
</script>
