<template>
  <Teleport to="body">
    <!-- 遮罩層 -->
    <div
      v-if="visible"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      @click.self="handleClose"
    >
      <div class="bg-slate-800 rounded-2xl w-full max-w-3xl overflow-hidden shadow-2xl">
        <!-- 標題欄 -->
        <div class="px-6 py-4 border-b border-white/10 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-primary-500/20 flex items-center justify-center">
              <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-bold text-slate-100">複製課表規則</h3>
              <p class="text-sm text-slate-400">將來源學期的規則複製到目標學期</p>
            </div>
          </div>
          <button
            @click="handleClose"
            class="text-slate-400 hover:text-slate-200 transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- 步驟指示器 -->
        <div class="px-6 py-4 border-b border-white/10">
          <div class="flex items-center justify-between">
            <div
              v-for="(step, index) in steps"
              :key="step.key"
              class="flex items-center"
              :class="{ 'flex-1': index < steps.length - 1 }"
            >
              <div class="flex items-center gap-2">
                <div
                  class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-all duration-300"
                  :class="[
                    currentStep > index + 1 ? 'bg-green-500/20 text-green-400' :
                    currentStep === index + 1 ? 'bg-primary-500/20 text-primary-400 ring-2 ring-primary-500/50' :
                    'bg-white/5 text-slate-500'
                  ]"
                >
                  <svg v-if="currentStep > index + 1" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <span v-else>{{ index + 1 }}</span>
                </div>
                <span
                  class="text-sm font-medium transition-colors"
                  :class="[
                    currentStep >= index + 1 ? 'text-slate-200' : 'text-slate-500'
                  ]"
                >
                  {{ step.label }}
                </span>
              </div>
              <div
                v-if="index < steps.length - 1"
                class="flex-1 h-0.5 mx-4 transition-colors"
                :class="currentStep > index + 1 ? 'bg-green-500/50' : 'bg-white/10'"
              ></div>
            </div>
          </div>
        </div>

        <!-- 內容區域 -->
        <div class="p-6 min-h-[400px]">
          <!-- Step 1: 選擇來源學期 -->
          <div v-if="currentStep === 1" class="space-y-4">
            <h4 class="text-lg font-semibold text-slate-100 mb-4">選擇來源學期</h4>
            <p class="text-slate-400 mb-4">選擇要複製規則的來源學期</p>

            <div v-if="loadingTerms" class="flex items-center justify-center py-12">
              <div class="animate-spin w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full"></div>
            </div>

            <div v-else-if="terms.length === 0" class="text-center py-12">
              <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-slate-700/50 flex items-center justify-center">
                <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <p class="text-slate-400">尚未建立任何學期</p>
            </div>

            <div v-else class="space-y-2">
              <label
                v-for="term in terms"
                :key="term.id"
                class="flex items-center p-4 rounded-xl border cursor-pointer transition-all"
                :class="[
                  form.sourceTermId === term.id
                    ? 'border-primary-500 bg-primary-500/10'
                    : 'border-white/10 bg-white/5 hover:bg-white/10'
                ]"
                @click="selectSourceTerm(term)"
              >
                <input
                  type="radio"
                  :value="term.id"
                  v-model="form.sourceTermId"
                  class="sr-only"
                />
                <div
                  class="w-5 h-5 rounded-full border-2 mr-3 flex items-center justify-center transition-colors"
                  :class="form.sourceTermId === term.id ? 'border-primary-500' : 'border-slate-500'"
                >
                  <div
                    v-if="form.sourceTermId === term.id"
                    class="w-2.5 h-2.5 rounded-full bg-primary-500"
                  ></div>
                </div>
                <div class="flex-1">
                  <div class="font-medium text-slate-200">{{ term.name }}</div>
                  <div class="text-sm text-slate-400">
                    {{ formatDate(term.start_date) }} ~ {{ formatDate(term.end_date) }}
                  </div>
                </div>
              </label>
            </div>
          </div>

          <!-- Step 2: 選擇要複製的規則 -->
          <div v-if="currentStep === 2" class="space-y-4">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h4 class="text-lg font-semibold text-slate-100">選擇要複製的規則</h4>
                <p class="text-sm text-slate-400 mt-1">
                  來源學期：<span class="text-primary-400">{{ getSourceTermName() }}</span>
                </p>
              </div>
              <div class="flex items-center gap-2">
                <button
                  @click="toggleSelectAll"
                  class="text-sm text-primary-400 hover:text-primary-300 transition-colors"
                >
                  {{ allRulesSelected ? '取消全選' : '全選' }}
                </button>
              </div>
            </div>

            <div v-if="loadingRules" class="flex items-center justify-center py-12">
              <div class="animate-spin w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full"></div>
            </div>

            <div v-else-if="sourceRules.length === 0" class="text-center py-12">
              <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-slate-700/50 flex items-center justify-center">
                <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                </svg>
              </div>
              <p class="text-slate-400">該學期沒有課表規則</p>
            </div>

            <div v-else class="space-y-2 max-h-[350px] overflow-y-auto pr-2">
              <label
                v-for="rule in sourceRules"
                :key="rule.id"
                class="flex items-center p-4 rounded-xl border cursor-pointer transition-all"
                :class="[
                  form.ruleIds.includes(rule.id)
                    ? 'border-primary-500 bg-primary-500/10'
                    : 'border-white/10 bg-white/5 hover:bg-white/10'
                ]"
              >
                <input
                  type="checkbox"
                  :value="rule.id"
                  v-model="form.ruleIds"
                  class="sr-only"
                />
                <div
                  class="w-5 h-5 rounded border-2 mr-3 flex items-center justify-center transition-colors"
                  :class="form.ruleIds.includes(rule.id) ? 'border-primary-500 bg-primary-500' : 'border-slate-500'"
                >
                  <svg v-if="form.ruleIds.includes(rule.id)" class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div class="flex-1">
                  <div class="font-medium text-slate-200">
                    {{ rule.offering_name || rule.course_name || '未命名規則' }}
                  </div>
                  <div class="text-sm text-slate-400 mt-1">
                    <span v-if="rule.teacher_name" class="mr-3">👤 {{ rule.teacher_name }}</span>
                    <span v-if="rule.room_name">📍 {{ rule.room_name }}</span>
                  </div>
                  <div class="text-sm text-slate-500 mt-1">
                    {{ getWeekdayName(rule.weekday) }} {{ rule.start_time }} ~ {{ rule.end_time }}
                  </div>
                </div>
              </label>
            </div>

            <div v-if="sourceRules.length > 0" class="text-sm text-slate-400 mt-4">
              已選擇 <span class="text-primary-400 font-medium">{{ form.ruleIds.length }}</span> 條規則
            </div>
          </div>

          <!-- Step 3: 選擇目標學期 -->
          <div v-if="currentStep === 3" class="space-y-4">
            <h4 class="text-lg font-semibold text-slate-100 mb-4">選擇目標學期</h4>
            <p class="text-slate-400 mb-4">選擇要接收複製規則的目標學期</p>

            <div v-if="loadingTerms" class="flex items-center justify-center py-12">
              <div class="animate-spin w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full"></div>
            </div>

            <div v-else class="space-y-2">
              <label
                v-for="term in availableTargetTerms"
                :key="term.id"
                class="flex items-center p-4 rounded-xl border cursor-pointer transition-all"
                :class="[
                  form.targetTermId === term.id
                    ? 'border-primary-500 bg-primary-500/10'
                    : 'border-white/10 bg-white/5 hover:bg-white/10'
                ]"
                @click="form.targetTermId = term.id"
              >
                <input
                  type="radio"
                  :value="term.id"
                  v-model="form.targetTermId"
                  class="sr-only"
                />
                <div
                  class="w-5 h-5 rounded-full border-2 mr-3 flex items-center justify-center transition-colors"
                  :class="form.targetTermId === term.id ? 'border-primary-500' : 'border-slate-500'"
                >
                  <div
                    v-if="form.targetTermId === term.id"
                    class="w-2.5 h-2.5 rounded-full bg-primary-500"
                  ></div>
                </div>
                <div class="flex-1">
                  <div class="font-medium text-slate-200">{{ term.name }}</div>
                  <div class="text-sm text-slate-400">
                    {{ formatDate(term.start_date) }} ~ {{ formatDate(term.end_date) }}
                  </div>
                </div>
                <div
                  v-if="term.id === form.sourceTermId"
                  class="text-amber-400 text-sm"
                >
                  (來源學期)
                </div>
              </label>
            </div>

            <div v-if="availableTargetTerms.length === 0" class="text-center py-8">
              <p class="text-slate-400">沒有可用的目標學期</p>
            </div>
          </div>

          <!-- Step 4: 確認摘要 -->
          <div v-if="currentStep === 4" class="space-y-6">
            <h4 class="text-lg font-semibold text-slate-100">確認複製</h4>

            <!-- 摘要卡片 -->
            <div class="bg-white/5 rounded-xl p-6 space-y-4">
              <!-- 來源學期 -->
              <div class="flex items-start gap-4">
                <div class="w-10 h-10 rounded-lg bg-primary-500/20 flex items-center justify-center flex-shrink-0">
                  <svg class="w-5 h-5 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                  </svg>
                </div>
                <div class="flex-1">
                  <div class="text-sm text-slate-400">來源學期</div>
                  <div class="font-medium text-slate-200">{{ getSourceTermName() }}</div>
                  <div class="text-sm text-slate-400">
                    {{ formatDate(getSourceTerm()?.start_date) }} ~ {{ formatDate(getSourceTerm()?.end_date) }}
                  </div>
                </div>
              </div>

              <!-- 目標學期 -->
              <div class="flex items-start gap-4">
                <div class="w-10 h-10 rounded-lg bg-green-500/20 flex items-center justify-center flex-shrink-0">
                  <svg class="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div class="flex-1">
                  <div class="text-sm text-slate-400">目標學期</div>
                  <div class="font-medium text-slate-200">{{ getTargetTermName() }}</div>
                  <div class="text-sm text-slate-400">
                    {{ formatDate(getTargetTerm()?.start_date) }} ~ {{ formatDate(getTargetTerm()?.end_date) }}
                  </div>
                </div>
              </div>

              <!-- 規則數量 -->
              <div class="pt-4 border-t border-white/10">
                <div class="flex items-center justify-between">
                  <div class="text-slate-400">將複製的規則數量</div>
                  <div class="text-2xl font-bold text-primary-400">{{ form.ruleIds.length }}</div>
                </div>
              </div>
            </div>

            <!-- 日期調整說明 -->
            <div class="bg-amber-500/10 border border-amber-500/30 rounded-xl p-4">
              <div class="flex items-start gap-3">
                <svg class="w-5 h-5 text-amber-400 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <div class="text-sm text-amber-200">
                  <div class="font-medium mb-1">日期調整說明</div>
                  <p>複製時，規則的有效日期範圍將自動調整為目標學期的開始和結束日期。</p>
                  <p class="mt-1">原日期：{{ formatDate(getSourceTerm()?.start_date) }} ~ {{ formatDate(getSourceTerm()?.end_date) }}</p>
                  <p>新日期：{{ formatDate(getTargetTerm()?.start_date) }} ~ {{ formatDate(getTargetTerm()?.end_date) }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- 處理中狀態 -->
          <div v-if="submitting" class="flex flex-col items-center justify-center py-12">
            <div class="animate-spin w-12 h-12 border-4 border-primary-500 border-t-transparent rounded-full mb-4"></div>
            <p class="text-slate-300">正在複製規則...</p>
            <p class="text-sm text-slate-500 mt-2">這可能需要幾秒鐘</p>
          </div>

          <!-- 成功狀態 -->
          <div v-if="submitSuccess" class="flex flex-col items-center justify-center py-8">
            <div class="w-16 h-16 rounded-full bg-green-500/20 flex items-center justify-center mb-4">
              <svg class="w-8 h-8 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h4 class="text-xl font-bold text-slate-100 mb-2">複製成功！</h4>
            <p class="text-slate-400 text-center">
              已成功複製 <span class="text-green-400 font-medium">{{ copiedCount }}</span> 條規則到
              <span class="text-primary-400 font-medium">{{ getTargetTermName() }}</span>
            </p>
          </div>

          <!-- 錯誤訊息 -->
          <div v-if="errorMessage" class="mt-4 p-4 bg-red-500/10 border border-red-500/30 rounded-xl">
            <div class="flex items-center gap-2 text-red-400">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{{ errorMessage }}</span>
            </div>
          </div>
        </div>

        <!-- 底部按鈕 -->
        <div v-if="!submitting && !submitSuccess" class="px-6 py-4 border-t border-white/10 flex justify-between">
          <button
            v-if="currentStep > 1"
            @click="prevStep"
            class="px-4 py-2 bg-white/5 border border-white/10 text-slate-300 rounded-xl hover:bg-white/10 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            上一步
          </button>
          <div v-else></div>

          <div class="flex items-center gap-3">
            <button
              @click="handleClose"
              class="px-4 py-2 bg-white/5 border border-white/10 text-slate-300 rounded-xl hover:bg-white/10 transition-colors"
            >
              取消
            </button>
            <button
              v-if="currentStep < 4"
              @click="nextStep"
              :disabled="!canProceed"
              class="px-4 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              下一步
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
            <button
              v-if="currentStep === 4"
              @click="submitCopy"
              :disabled="submitting || form.ruleIds.length === 0"
              class="px-4 py-2 bg-green-500/30 border border-green-500 text-green-400 rounded-xl hover:bg-green-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              確認複製
            </button>
          </div>
        </div>

        <!-- 成功狀態的底部按鈕 -->
        <div v-if="submitSuccess" class="px-6 py-4 border-t border-white/10 flex justify-end">
          <button
            @click="handleClose"
            class="px-4 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors flex items-center gap-2"
          >
            完成
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Term } from '~/types/scheduling'

// 規則類型（簡化版）
interface RuleInfo {
  id: number
  offering_id: number
  offering_name?: string
  course_id?: number
  course_name?: string
  teacher_id?: number
  teacher_name?: string
  room_id: number
  room_name?: string
  weekday: number
  start_time: string
  end_time: string
}

// 複製規則請求
interface CopyRulesRequest {
  source_term_id: number
  target_term_id: number
  rule_ids: number[]
}

// 複製規則回應
interface CopyRulesResponse {
  copied_count: number
  rules: Array<{
    original_rule_id: number
    new_rule_id: number
    offering_name: string
    weekday: number
    start_time: string
    end_time: string
  }>
}

// Props
const props = defineProps<{
  visible: boolean
}>()

// Emits
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

// 步驟定義
const steps = [
  { key: 'source', label: '選擇來源' },
  { key: 'rules', label: '選擇規則' },
  { key: 'target', label: '選擇目標' },
  { key: 'confirm', label: '確認複製' }
]

// 當前步驟
const currentStep = ref(1)

// 表單資料
const form = ref({
  sourceTermId: 0 as number,
  targetTermId: 0 as number,
  ruleIds: [] as number[]
})

// 學期列表
const terms = ref<Term[]>([])
const loadingTerms = ref(false)

// 來源規則列表
const sourceRules = ref<RuleInfo[]>([])
const loadingRules = ref(false)

// 提交狀態
const submitting = ref(false)
const submitSuccess = ref(false)
const copiedCount = ref(0)
const errorMessage = ref('')

// API
const api = useApi()
const { success: toastSuccess, error: toastError } = useToast()

// 計算屬性
const allRulesSelected = computed(() =>
  sourceRules.value.length > 0 &&
  form.value.ruleIds.length === sourceRules.value.length
)

const availableTargetTerms = computed(() =>
  terms.value.filter(term => term.id !== form.value.sourceTermId)
)

const canProceed = computed(() => {
  switch (currentStep.value) {
    case 1:
      return form.value.sourceTermId > 0
    case 2:
      return form.value.ruleIds.length > 0
    case 3:
      return form.value.targetTermId > 0
    case 4:
      return true
    default:
      return false
  }
})

// 監聽 visible 變化
watch(() => props.visible, (val) => {
  if (val) {
    resetForm()
    fetchTerms()
  }
})

// 獲取學期列表
const fetchTerms = async () => {
  loadingTerms.value = true
  try {
    const response = await api.get<Term[]>('/admin/terms')
    if (response) {
      terms.value = response
    }
  } catch (err: any) {
    console.error('Failed to fetch terms:', err)
    toastError('載入學期失敗')
  } finally {
    loadingTerms.value = false
  }
}

// 獲取來源學期的規則
const fetchSourceRules = async () => {
  if (!form.value.sourceTermId) return

  loadingRules.value = true
  try {
    const sourceTerm = terms.value.find(t => t.id === form.value.sourceTermId)
    if (!sourceTerm) return

    // 使用 /admin/rules API 獲取所有規則
    const response = await api.get<any[]>('/admin/rules')

    if (response) {
      // 規則為循環規則，無需按日期過濾，直接使用所有規則
      const allRules: RuleInfo[] = response.map((rule: any) => ({
        id: rule.id,
        offering_id: rule.offering_id,
        offering_name: rule.offering_name,
        course_id: rule.course_id,
        course_name: rule.course_name,
        teacher_id: rule.teacher_id,
        teacher_name: rule.teacher_name,
        room_id: rule.room_id,
        room_name: rule.room_name,
        weekday: rule.weekday,
        start_time: rule.start_time,
        end_time: rule.end_time
      }))
      sourceRules.value = allRules
    } else {
      sourceRules.value = []
    }
  } catch (err: any) {
    console.error('Failed to fetch rules:', err)
    toastError('載入規則失敗')
  } finally {
    loadingRules.value = false
  }
}

// 選擇來源學期
const selectSourceTerm = (term: Term) => {
  form.value.sourceTermId = term.id
}

// 切換全選
const toggleSelectAll = () => {
  if (allRulesSelected.value) {
    form.value.ruleIds = []
  } else {
    form.value.ruleIds = sourceRules.value.map(r => r.id)
  }
}

// 下一步
const nextStep = () => {
  if (currentStep.value === 1) {
    // 從步驟 1 到步驟 2，需要載入規則
    fetchSourceRules()
  }
  if (currentStep.value < 4) {
    currentStep.value++
  }
}

// 上一步
const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

// 提交複製
const submitCopy = async () => {
  if (!form.value.sourceTermId || !form.value.targetTermId || form.value.ruleIds.length === 0) {
    return
  }

  submitting.value = true
  errorMessage.value = ''

  try {
    const request: CopyRulesRequest = {
      source_term_id: form.value.sourceTermId,
      target_term_id: form.value.targetTermId,
      rule_ids: form.value.ruleIds
    }

    const response = await api.post<CopyRulesResponse>('/admin/terms/copy-rules', request)

    if (response) {
      submitSuccess.value = true
      copiedCount.value = response.copied_count || 0
      emit('success')
    } else {
      errorMessage.value = '複製失敗，請稍後再試'
    }
  } catch (err: any) {
    console.error('Failed to copy rules:', err)
    errorMessage.value = err.message || '複製失敗，請稍後再試'
  } finally {
    submitting.value = false
  }
}

// 重置表單
const resetForm = () => {
  currentStep.value = 1
  form.value = {
    sourceTermId: 0,
    targetTermId: 0,
    ruleIds: []
  }
  sourceRules.value = []
  submitting.value = false
  submitSuccess.value = false
  copiedCount.value = 0
  errorMessage.value = ''
}

// 關閉
const handleClose = () => {
  resetForm()
  emit('close')
}

// 輔助函數
const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const getSourceTerm = () => {
  return terms.value.find(t => t.id === form.value.sourceTermId)
}

const getTargetTerm = () => {
  return terms.value.find(t => t.id === form.value.targetTermId)
}

const getSourceTermName = () => {
  const term = getSourceTerm()
  return term?.name || '-'
}

const getTargetTermName = () => {
  const term = getTargetTerm()
  return term?.name || '-'
}

const getWeekdayName = (weekday: number) => {
  const names = ['週日', '週一', '週二', '週三', '週四', '週五', '週六']
  return names[weekday] || '-'
}
</script>
