<template>
  <div class="space-y-6">
    <!-- 標題與新增按鈕 -->
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-bold text-slate-100">學期管理</h2>
      <div class="flex items-center gap-3">
        <button
          @click="openCopyModal"
          :disabled="terms.length < 2"
          class="px-4 py-2 bg-amber-500/20 border border-amber-500/50 text-amber-400 rounded-xl hover:bg-amber-500/30 transition-colors flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <span>複製規則</span>
        </button>
        <button
          @click="openAddModal"
          class="px-4 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span>新增學期</span>
        </button>
      </div>
    </div>

    <!-- 學期列表 -->
    <div class="bg-white/5 rounded-xl overflow-hidden">
      <div v-if="loading && terms.length === 0" class="flex items-center justify-center py-12">
        <div class="animate-spin w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full"></div>
      </div>

      <div v-else-if="terms.length === 0" class="text-center py-12">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-slate-700/50 flex items-center justify-center">
          <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <p class="text-slate-400 mb-4">尚未建立任何學期</p>
        <button
          @click="openAddModal"
          class="text-primary-400 hover:text-primary-300 transition-colors"
        >
        </button>
      </div>

      <table v-else class="w-full">
        <thead class="bg-white/5">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
              學期名稱
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
              開始日期
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
              結束日期
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-slate-400 uppercase tracking-wider">
              操作
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="term in terms" :key="term.id" class="hover:bg-white/5 transition-colors">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-slate-200">{{ term.name }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-slate-400">{{ formatDate(term.start_date) }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-slate-400">{{ formatDate(term.end_date) }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button
                @click="openEditModal(term as Term)"
                class="text-primary-400 hover:text-primary-300 transition-colors mr-4"
              >
                編輯
              </button>
              <button
                @click="openDeleteConfirm(term as Term)"
                class="text-red-400 hover:text-red-300 transition-colors"
              >
                刪除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 新增/編輯學期 Modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      @click.self="closeModal"
    >
      <div class="bg-slate-800 rounded-2xl w-full max-w-md p-6">
        <h3 class="text-lg font-bold text-slate-100 mb-6">
          {{ isEditing ? '編輯學期' : '新增學期' }}
        </h3>

        <form @submit.prevent="submitForm" class="space-y-4">
          <!-- 學期名稱 -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              學期名稱 <span class="text-red-400">*</span>
            </label>
            <input
              v-model="form.name"
              type="text"
              placeholder="例如：114學年度上學期"
              class="w-full px-4 py-2.5 bg-white/5 border border-white/10 rounded-xl text-slate-100 placeholder-slate-500 focus:outline-none focus:border-primary-500 transition-colors"
              required
            />
          </div>

          <!-- 開始日期 -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              開始日期 <span class="text-red-400">*</span>
            </label>
            <input
              v-model="form.start_date"
              type="date"
              class="w-full px-4 py-2.5 bg-white/5 border border-white/10 rounded-xl text-slate-100 focus:outline-none focus:border-primary-500 transition-colors"
              required
            />
          </div>

          <!-- 結束日期 -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              結束日期 <span class="text-red-400">*</span>
            </label>
            <input
              v-model="form.end_date"
              type="date"
              class="w-full px-4 py-2.5 bg-white/5 border border-white/10 rounded-xl text-slate-100 focus:outline-none focus:border-primary-500 transition-colors"
              required
            />
          </div>

          <!-- 錯誤訊息 -->
          <div v-if="errorMessage" class="text-sm text-red-400 bg-red-400/10 px-4 py-2 rounded-lg">
            {{ errorMessage }}
          </div>

          <!-- 按鈕 -->
          <div class="flex justify-end gap-3 mt-6">
            <button
              type="button"
              @click="closeModal"
              class="px-4 py-2 bg-white/5 border border-white/10 text-slate-300 rounded-xl hover:bg-white/10 transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="submitting"
              class="px-4 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors disabled:opacity-50 flex items-center gap-2"
            >
              <svg v-if="submitting" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>{{ isEditing ? '儲存變更' : '新增' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 刪除確認 -->
    <GlobalAlert
      v-if="showDeleteConfirm"
      type="warning"
      title="刪除學期"
      message="確定要刪除此學期嗎？刪除後將無法恢復。"
      confirmText="確定刪除"
      cancelText="取消"
      @confirm="deleteTerm"
      @cancel="showDeleteConfirm = false"
    />

    <!-- 複製規則精靈 -->
    <TermCopyModal
      :visible="showCopyModal"
      @close="showCopyModal = false"
      @success="onCopySuccess"
    />
  </div>
</template>

<script setup lang="ts">
import type { Term } from '~/types/scheduling'
import GlobalAlert from '~/components/base/GlobalAlert.vue'
import TermCopyModal from './TermCopyModal.vue'

// 表單類型
interface TermForm {
  name: string
  start_date: string
  end_date: string
}

// 狀態
const loading = ref(true)
const submitting = ref(false)
const terms = ref<Term[]>([])
const showModal = ref(false)
const showDeleteConfirm = ref(false)
const showCopyModal = ref(false)
const isEditing = ref(false)
const selectedTerm = ref<Term | null>(null)
const errorMessage = ref('')

const api = useApi()
const { success, error: toastError } = useToast()
const { confirm: alertConfirm, error: alertError } = useAlert()

const form = ref<TermForm>({
  name: '',
  start_date: '',
  end_date: ''
})

// 取得學期列表
const fetchTerms = async () => {
  loading.value = true
  try {
    const response = await api.get<Term[]>('/admin/terms')
    if (response) {
      terms.value = response
    }
  } catch (err: any) {
    console.error('Failed to fetch terms:', err)
    toastError('載入學期失敗')
  } finally {
    loading.value = false
  }
}

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 打開新增 Modal
const openAddModal = () => {
  isEditing.value = false
  form.value = {
    name: '',
    start_date: '',
    end_date: ''
  }
  errorMessage.value = ''
  showModal.value = true
}

// 打開編輯 Modal
const openEditModal = (term: Term) => {
  isEditing.value = true
  selectedTerm.value = term
  form.value = {
    name: term.name,
    start_date: term.start_date,
    end_date: term.end_date
  }
  errorMessage.value = ''
  showModal.value = true
}

// 關閉 Modal
const closeModal = () => {
  showModal.value = false
  selectedTerm.value = null
  form.value = {
    name: '',
    start_date: '',
    end_date: ''
  }
  errorMessage.value = ''
}

// 打開刪除確認
const openDeleteConfirm = (term: Term) => {
  selectedTerm.value = term
  showDeleteConfirm.value = true
}

// 打開複製規則精靈
const openCopyModal = () => {
  showCopyModal.value = true
}

// 複製成功回調
const onCopySuccess = () => {
  // 可選：複製成功後執行額外操作
  // 例如：重新整理規則列表或顯示提示
}

// 提交表單
const submitForm = async () => {
  errorMessage.value = ''

  // 驗證日期
  if (new Date(form.value.end_date) < new Date(form.value.start_date)) {
    errorMessage.value = '結束日期必須晚於開始日期'
    return
  }

  submitting.value = true
  try {
    if (isEditing.value && selectedTerm.value) {
      // 更新學期
      await api.put(`/admin/terms/${selectedTerm.value.id}`, form.value)
      success('學期更新成功')
    } else {
      // 新增學期
      await api.post('/admin/terms', form.value)
      success('學期建立成功')
    }
    closeModal()
    await fetchTerms()
  } catch (err: any) {
    errorMessage.value = err.message || (isEditing.value ? '更新學期失敗' : '建立學期失敗')
  } finally {
    submitting.value = false
  }
}

// 刪除學期
const deleteTerm = async () => {
  if (!selectedTerm.value) return

  try {
    await api.delete(`/admin/terms/${selectedTerm.value.id}`)
    success('學期已刪除')
    showDeleteConfirm.value = false
    selectedTerm.value = null
    await fetchTerms()
  } catch (err: any) {
    alertError(err.message || '刪除學期失敗')
  }
}

// 生命週期
onMounted(fetchTerms)
</script>
