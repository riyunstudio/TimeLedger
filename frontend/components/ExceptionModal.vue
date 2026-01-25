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
                {{ rule.title }} - {{ formatDate(rule.original_date) }} {{ rule.start_time }}-{{ rule.end_time }}
              </option>
            </select>
            <p v-if="loadingRules" class="text-xs text-slate-500 mt-1">載入中...</p>
            <p v-else-if="form.center_id && displayScheduleRules.length === 0" class="text-xs text-slate-500 mt-1">該中心暫無您的課程</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">申請類型</label>
            <div class="flex gap-3">
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" v-model="form.type" value="CANCEL" class="accent-primary-500" />
                <span class="text-white">停課</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" v-model="form.type" value="RESCHEDULE" class="accent-primary-500" />
                <span class="text-white">改期</span>
              </label>
            </div>
          </div>

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

interface Props {
  exception?: ScheduleException
  centers: Array<{ center_id: number; center_name: string }>
  scheduleRules?: Array<{ id: number; title: string; original_date: string; start_time: string; end_time: string }>
}

const props = defineProps<Props>()
const emit = defineEmits(['close', 'submit'])

const teacherStore = useTeacherStore()
const loading = ref(false)
const localScheduleRules = ref<Array<{ id: number; title: string; original_date: string; start_time: string; end_time: string }>>([])
const loadingRules = ref(false)

const isEdit = computed(() => !!props.exception)

const form = reactive({
  center_id: props.exception?.center_id || 0,
  rule_id: props.exception?.rule_id || 0,
  original_date: props.exception?.original_date || '',
  type: props.exception?.type || 'CANCEL' as 'CANCEL' | 'RESCHEDULE',
  new_start_at: props.exception?.new_start_at || '',
  new_end_at: props.exception?.new_end_at || '',
  reason: props.exception?.reason || '',
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
    const response = await api.get<{ code: number; datas: any[] }>(`/teacher/me/centers/${centerId}/schedule-rules`)
    if (response.code === 0 && response.datas) {
      // 轉換資料結構
      localScheduleRules.value = response.datas.map((rule: any) => ({
        id: rule.id,
        title: rule.offering?.name || '課程',
        original_date: rule.original_date,
        start_time: rule.start_time,
        end_time: rule.end_time,
      }))
    }
  } catch (error) {
    console.error('Failed to fetch schedule rules:', error)
  } finally {
    loadingRules.value = false
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW')
}

const handleSubmit = async () => {
  loading.value = true
  try {
    await teacherStore.createException({
      center_id: form.center_id,
      rule_id: form.rule_id,
      original_date: form.original_date,
      type: form.type,
      new_start_at: form.new_start_at || undefined,
      new_end_at: form.new_end_at || undefined,
      reason: form.reason,
    })
    emit('submit')
    emit('close')
  } catch (error) {
    console.error('Failed to create exception:', error)
  } finally {
    loading.value = false
  }
}
</script>
