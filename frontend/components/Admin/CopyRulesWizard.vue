<template>
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" @click.self="$emit('close')">
    <div class="bg-slate-800 rounded-2xl w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col">
      <!-- 標題列 -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-white/10">
        <h2 class="text-xl font-bold text-slate-100">複製規則精靈</h2>
        <button
          @click="$emit('close')"
          class="text-slate-400 hover:text-white transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- 步驟指示器 -->
      <div class="flex items-center justify-center px-6 py-4 border-b border-white/10">
        <div class="flex items-center gap-4">
          <div
            v-for="step in 4"
            :key="step"
            class="flex items-center"
          >
            <div
              class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-colors"
              :class="currentStep >= step
                ? 'bg-primary-500 text-white'
                : 'bg-white/10 text-slate-400'"
            >
              <svg v-if="currentStep > step" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span v-else>{{ step }}</span>
            </div>
            <div
              v-if="step < 4"
              class="w-16 h-0.5 mx-2"
              :class="currentStep > step ? 'bg-primary-500' : 'bg-white/10'"
            />
          </div>
        </div>
      </div>

      <!-- 步驟內容 -->
      <div class="flex-1 overflow-y-auto p-6">
        <!-- 步驟 1：選擇來源學期 -->
        <div v-if="currentStep === 1" class="space-y-4">
          <h3 class="text-lg font-medium text-slate-100 mb-4">選擇來源學期</h3>
          <p class="text-slate-400 mb-4">選擇要複製規則的來源學期</p>

          <div class="space-y-3">
            <label
              v-for="term in terms"
              :key="term.id"
              class="flex items-center p-4 rounded-xl border cursor-pointer transition-all"
              :class="sourceTermId === term.id
                ? 'bg-primary-500/20 border-primary-500'
                : 'bg-white/5 border-white/10 hover:bg-white/10'"
            >
              <input
                type="radio"
                :value="term.id"
                v-model="sourceTermId"
                class="sr-only"
              />
              <div class="flex-1">
                <div class="font-medium text-slate-100">{{ term.name }}</div>
                <div class="text-sm text-slate-400 mt-1">
                  {{ formatDateRange(term.start_date, term.end_date) }}
                </div>
              </div>
              <div
                v-if="sourceTermId === term.id"
                class="w-6 h-6 rounded-full bg-primary-500 flex items-center justify-center"
              >
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </div>
            </label>
          </div>

          <div v-if="terms.length === 0" class="text-center py-8 text-slate-400">
            尚未建立任何學期
          </div>
        </div>

        <!-- 步驟 2：選擇目標學期 -->
        <div v-else-if="currentStep === 2" class="space-y-4">
          <h3 class="text-lg font-medium text-slate-100 mb-4">選擇目標學期</h3>
          <p class="text-slate-400 mb-4">選擇要複製規則的目標學期</p>

          <div class="space-y-3">
            <label
              v-for="term in availableTargetTerms"
              :key="term.id"
              class="flex items-center p-4 rounded-xl border cursor-pointer transition-all"
              :class="targetTermId === term.id
                ? 'bg-primary-500/20 border-primary-500'
                : 'bg-white/5 border-white/10 hover:bg-white/10'"
            >
              <input
                type="radio"
                :value="term.id"
                v-model="targetTermId"
                class="sr-only"
              />
              <div class="flex-1">
                <div class="font-medium text-slate-100">{{ term.name }}</div>
                <div class="text-sm text-slate-400 mt-1">
                  {{ formatDateRange(term.start_date, term.end_date) }}
                </div>
              </div>
              <div
                v-if="targetTermId === term.id"
                class="w-6 h-6 rounded-full bg-primary-500 flex items-center justify-center"
              >
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </div>
            </label>
          </div>

          <div v-if="availableTargetTerms.length === 0" class="text-center py-8 text-slate-400">
            沒有可用的目標學期
          </div>
        </div>

        <!-- 步驟 3：選擇規則 -->
        <div v-else-if="currentStep === 3" class="space-y-4">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-medium text-slate-100">選擇要複製的規則</h3>
            <div class="flex items-center gap-4 text-sm">
              <button
                @click="selectAll"
                class="text-primary-400 hover:text-primary-300 transition-colors"
              >
                全選
              </button>
              <button
                @click="deselectAll"
                class="text-slate-400 hover:text-white transition-colors"
              >
                取消全選
              </button>
            </div>
          </div>

          <!-- 規則列表 -->
          <div v-if="loadingRules" class="flex items-center justify-center py-12">
            <div class="animate-spin w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full"></div>
          </div>

          <div v-else-if="rules.length === 0" class="text-center py-8 text-slate-400">
            來源學期沒有規則
          </div>

          <div v-else class="space-y-2 max-h-96 overflow-y-auto">
            <label
              v-for="rule in rules"
              :key="rule.id"
              class="flex items-center p-3 rounded-xl border cursor-pointer transition-all"
              :class="selectedRuleIds.has(rule.id)
                ? 'bg-primary-500/20 border-primary-500'
                : 'bg-white/5 border-white/10 hover:bg-white/10'"
            >
              <input
                type="checkbox"
                :value="rule.id"
                :checked="selectedRuleIds.has(rule.id)"
                @change="toggleRule(rule.id)"
                class="sr-only"
              />
              <div class="flex-1">
                <div class="font-medium text-slate-100">
                  {{ rule.offering_name || rule.course_name }}
                </div>
                <div class="text-sm text-slate-400 mt-1">
                  {{ getWeekdayText(rule.weekday) }} {{ rule.start_time }} - {{ rule.end_time }}
                  <span v-if="rule.teacher_name">｜{{ rule.teacher_name }}</span>
                  <span v-if="rule.room_name">｜{{ rule.room_name }}</span>
                </div>
              </div>
              <div
                v-if="selectedRuleIds.has(rule.id)"
                class="w-5 h-5 rounded bg-primary-500 flex items-center justify-center"
              >
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </div>
            </label>
          </div>

          <div class="text-sm text-slate-400">
            已選擇 {{ selectedRuleIds.size }} 個規則
          </div>
        </div>

        <!-- 步驟 4：確認複製 -->
        <div v-else-if="currentStep === 4" class="space-y-6">
          <h3 class="text-lg font-medium text-slate-100 mb-4">確認複製</h3>

          <div class="glass-card p-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <div class="text-sm text-slate-400">來源學期</div>
                <div class="font-medium text-slate-100 mt-1">
                  {{ getTermName(sourceTermId) }}
                </div>
              </div>
              <div>
                <div class="text-sm text-slate-400">目標學期</div>
                <div class="font-medium text-slate-100 mt-1">
                  {{ getTermName(targetTermId) }}
                </div>
              </div>
              <div>
                <div class="text-sm text-slate-400">將複製規則數</div>
                <div class="font-medium text-primary-400 mt-1">
                  {{ selectedRuleIds.size }} 個
                </div>
              </div>
              <div>
                <div class="text-sm text-slate-400">日期調整</div>
                <div class="font-medium text-slate-100 mt-1">
                  自動對齊目標學期
                </div>
              </div>
            </div>
          </div>

          <div class="bg-yellow-500/10 border border-yellow-500/30 rounded-xl p-4">
            <div class="flex items-start gap-3">
              <svg class="w-5 h-5 text-yellow-400 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <div class="text-sm text-yellow-200">
                <div class="font-medium">注意事項</div>
                <ul class="mt-2 space-y-1 text-slate-300">
                  <li>• 複製後的規則日期會自動調整為目標學期的日期範圍</li>
                  <li>• 如果目標學期已有相同規則，可能會產生衝突</li>
                  <li>• 此操作不可復原，請確認後再執行</li>
                </ul>
              </div>
            </div>
          </div>

          <label class="flex items-center gap-3 cursor-pointer">
            <input
              type="checkbox"
              v-model="confirmed"
              class="w-5 h-5 rounded border-white/20 bg-white/5 text-primary-500 focus:ring-primary-500"
            />
            <span class="text-sm text-slate-300">
              我已確認了解上述注意事項，執行複製
            </span>
          </label>
        </div>
      </div>

      <!-- 底部按鈕 -->
      <div class="flex items-center justify-between px-6 py-4 border-t border-white/10">
        <button
          v-if="currentStep > 1"
          @click="currentStep--"
          class="px-4 py-2 bg-white/5 border border-white/10 text-slate-300 rounded-xl hover:bg-white/10 transition-colors"
        >
          上一步
        </button>
        <div v-else />

        <div class="flex items-center gap-3">
          <button
            @click="$emit('close')"
            class="px-4 py-2 bg-white/5 border border-white/10 text-slate-300 rounded-xl hover:bg-white/10 transition-colors"
          >
            取消
          </button>
          <button
            v-if="currentStep < 4"
            @click="nextStep"
            :disabled="!canProceed"
            class="px-6 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            下一步
          </button>
          <button
            v-else
            @click="executeCopy"
            :disabled="!confirmed || copying"
            class="px-6 py-2 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <svg v-if="copying" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ copying ? '複製中...' : '執行複製' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Term, OccupancyRule } from '~/types/scheduling'

// ========== Emits ==========

const emit = defineEmits<{
  close: []
  copied: []
}>()

// ========== Props ==========

const props = defineProps<{
  terms: Term[]
}>()

// ========== State ==========

const api = useApi()
const { success, error: toastError } = useToast()

const currentStep = ref(1)
const sourceTermId = ref<number | null>(null)
const targetTermId = ref<number | null>(null)
const loadingRules = ref(false)
const rules = ref<OccupancyRule[]>([])
const selectedRuleIds = ref<Set<number>>(new Set())
const copying = ref(false)
const confirmed = ref(false)

// ========== Computed ==========

/**
 * 可用的目標學期（排除來源學期）
 */
const availableTargetTerms = computed(() => {
  return props.terms.filter(term => term.id !== sourceTermId.value)
})

/**
 * 是否可以進行下一步
 */
const canProceed = computed(() => {
  switch (currentStep.value) {
    case 1:
      return sourceTermId.value !== null
    case 2:
      return targetTermId.value !== null
    case 3:
      return selectedRuleIds.value.size > 0
    case 4:
      return confirmed.value
    default:
      return false
  }
})

// ========== Methods ==========

/**
 * 格式化日期區間
 */
const formatDateRange = (startDate: string, endDate: string) => {
  const start = new Date(startDate)
  const end = new Date(endDate)
  return `${start.getFullYear()}/${start.getMonth() + 1}/${start.getDate()} - ${end.getFullYear()}/${end.getMonth() + 1}/${end.getDate()}`
}

/**
 * 取得星期文字
 */
const getWeekdayText = (weekday: number): string => {
  const days = ['週日', '週一', '週二', '週三', '週四', '週五', '週六']
  return days[weekday === 7 ? 0 : weekday]
}

/**
 * 取得學期名稱
 */
const getTermName = (termId: number | null): string => {
  if (!termId) return '-'
  const term = props.terms.find(t => t.id === termId)
  return term?.name || '-'
}

/**
 * 下一步
 */
const nextStep = async () => {
  if (currentStep.value < 4) {
    currentStep.value++

    // 當進入步驟 3 時，取得規則列表
    if (currentStep.value === 3) {
      await fetchRules()
    }
  }
}

// 取得規則列表
const fetchRules = async () => {
  if (!sourceTermId.value) return

  loadingRules.value = true

  try {
    // 使用 /admin/rules API 獲取所有規則
    const response = await api.get<any[]>('/admin/rules')

    if (response) {
      // 規則為循環規則，無需按日期過濾，直接使用所有規則
      rules.value = response
    } else {
      rules.value = []
    }
  } catch (err) {
    console.error('Failed to fetch rules:', err)
    toastError('載入規則失敗')
  } finally {
    loadingRules.value = false
  }
}

/**
 * 切換規則選擇
 */
const toggleRule = (ruleId: number) => {
  if (selectedRuleIds.value.has(ruleId)) {
    selectedRuleIds.value.delete(ruleId)
  } else {
    selectedRuleIds.value.add(ruleId)
  }
}

/**
 * 全選
 */
const selectAll = () => {
  rules.value.forEach(rule => {
    selectedRuleIds.value.add(rule.id)
  })
}

/**
 * 取消全選
 */
const deselectAll = () => {
  selectedRuleIds.value.clear()
}

/**
 * 執行複製
 */
const executeCopy = async () => {
  if (!sourceTermId.value || !targetTermId.value || selectedRuleIds.value.size === 0) return

  copying.value = true

  try {
    await api.post('/admin/terms/copy-rules', {
      source_term_id: sourceTermId.value,
      target_term_id: targetTermId.value,
      rule_ids: Array.from(selectedRuleIds.value),
    })

    success('規則複製成功')
    emit('copied')
  } catch (err) {
    console.error('Failed to copy rules:', err)
    toastError('複製失敗，請稍後再試')
  } finally {
    copying.value = false
  }
}

// ========== Watch ==========

/**
 * 當來源學期改變時，重置後續狀態
 */
watch(sourceTermId, () => {
  targetTermId.value = null
  selectedRuleIds.value.clear()
  rules.value = []
})
</script>
