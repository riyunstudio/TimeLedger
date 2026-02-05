<template>
  <div class="p-4 md:p-6 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
            老師管理
          </h1>
          <p class="text-slate-400 text-sm md:text-base">
            管理中心的老師帳號與綁定狀態
          </p>
        </div>
        <!-- 手動新增老師按鈕 -->
        <button
          @click="showCreateModal = true"
          class="px-4 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          手動新增老師
        </button>
      </div>
    </div>

    <!-- 篩選器 -->
    <div class="mb-6 flex flex-wrap items-center gap-4">
      <!-- 搜尋框 -->
      <div class="flex-1 min-w-[200px]">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜尋老師名稱或 Email..."
          class="w-full px-4 py-2 rounded-lg bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500"
        />
      </div>

      <!-- 佔位篩選 -->
      <select
        v-model="filterPlaceholder"
        class="px-3 py-2 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 min-w-[140px]"
      >
        <option value="">全部狀態</option>
        <option value="placeholder">佔位老師</option>
        <option value="bound">已綁定老師</option>
      </select>

      <!-- 清除篩選 -->
      <button
        v-if="searchQuery || filterPlaceholder"
        @click="clearFilters"
        class="text-sm text-primary-400 hover:text-primary-300 transition-colors"
      >
        清除篩選
      </button>
    </div>

    <!-- 老師列表 -->
    <BaseGlassCard>
      <div class="p-6">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="text-left text-slate-400 text-sm border-b border-white/10">
                <th class="pb-4 pr-4 font-medium">姓名</th>
                <th class="pb-4 pr-4 font-medium">Email</th>
                <th class="pb-4 pr-4 font-medium">狀態</th>
                <th class="pb-4 pr-4 font-medium">LINE 綁定</th>
                <th class="pb-4 font-medium">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="teacher in filteredTeachers"
                :key="teacher.id"
                class="border-b border-white/5 hover:bg-white/5 transition-colors"
              >
                <td class="py-4 pr-4">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-primary-500/20 flex items-center justify-center">
                      <span class="text-primary-400 font-medium text-sm">
                        {{ teacher.name.charAt(0) }}
                      </span>
                    </div>
                    <div>
                      <span class="text-white font-medium">{{ teacher.name }}</span>
                      <!-- 佔位老師標籤 -->
                      <span
                        v-if="teacher.is_placeholder"
                        class="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-warning-500/20 text-warning-400"
                      >
                        佔位老師
                      </span>
                    </div>
                  </div>
                </td>
                <td class="py-4 pr-4 text-slate-300">{{ teacher.email || '-' }}</td>
                <td class="py-4 pr-4">
                  <span
                    class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="teacher.is_placeholder
                      ? 'bg-warning-500/20 text-warning-400'
                      : 'bg-success-500/20 text-success-400'"
                  >
                    <span
                      class="w-1.5 h-1.5 rounded-full"
                      :class="teacher.is_placeholder ? 'bg-warning-400' : 'bg-success-400'"
                    />
                    {{ teacher.is_placeholder ? '待綁定' : '已激活' }}
                  </span>
                </td>
                <td class="py-4 pr-4">
                  <div v-if="teacher.line_user_id" class="flex items-center gap-1.5 text-success-400">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    <span class="text-xs">已綁定</span>
                  </div>
                  <span v-else class="text-slate-500 text-sm">未綁定</span>
                </td>
                <td class="py-4">
                  <div class="flex items-center gap-2">
                    <!-- 合併帳號按鈕（僅佔位老師顯示） -->
                    <button
                      v-if="teacher.is_placeholder"
                      @click="openMergeModal(teacher)"
                      class="px-3 py-1.5 rounded-lg bg-primary-500/20 text-primary-400 hover:bg-primary-500/30 transition-colors text-sm flex items-center gap-1.5"
                      title="將此佔位老師合併至已綁定老師帳號"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                      </svg>
                      合併帳號
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- 空狀態 -->
          <div v-if="filteredTeachers.length === 0 && !loading" class="text-center py-12">
            <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-slate-700/50 flex items-center justify-center">
              <svg class="w-8 h-8 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </div>
            <p class="text-slate-400">暫無符合條件的老師</p>
          </div>

          <!-- 載入中 -->
          <div v-if="loading" class="text-center py-12">
            <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
            <p class="text-slate-400 mt-4">載入中...</p>
          </div>
        </div>
      </div>
    </BaseGlassCard>

    <!-- 新增佔位老師 Modal -->
    <Teleport to="body">
      <div
        v-if="showCreateModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="showCreateModal = false"
      >
        <div class="glass-card w-full max-w-md">
          <div class="p-6 border-b border-white/10">
            <h3 class="text-lg font-semibold text-white">手動新增老師</h3>
            <p class="text-sm text-slate-400 mt-1">
              建立佔位老師，稍後可合併至正式 LINE 帳號
            </p>
          </div>

          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm text-slate-400 mb-2">
                姓名 <span class="text-critical-500">*</span>
              </label>
              <BaseInput
                v-model="createForm.name"
                placeholder="請輸入老師姓名"
                class="w-full"
              />
            </div>

            <div>
              <label class="block text-sm text-slate-400 mb-2">
                Email <span class="text-slate-500">(選填)</span>
              </label>
              <BaseInput
                v-model="createForm.email"
                type="email"
                placeholder="請輸入 Email"
                class="w-full"
              />
            </div>

            <!-- 提示資訊 -->
            <div class="p-3 rounded-lg bg-warning-500/10 border border-warning-500/30">
              <p class="text-sm text-warning-400">
                建立後，老師可透過 LINE 登入並自動綁定此帳號。
              </p>
            </div>
          </div>

          <div class="p-6 border-t border-white/10 flex items-center gap-4">
            <button
              @click="showCreateModal = false"
              class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              取消
            </button>
            <button
              @click="createPlaceholderTeacher"
              :disabled="!isCreateValid || creating"
              class="flex-1 px-4 py-2 rounded-lg bg-primary-500/30 border border-primary-500 text-primary-400 hover:bg-primary-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="creating" class="flex items-center justify-center gap-2">
                <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                處理中...
              </span>
              <span v-else>新增老師</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 合併帳號 Modal -->
    <Teleport to="body">
      <div
        v-if="showMergeModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="closeMergeModal"
      >
        <div class="glass-card w-full max-w-md">
          <div class="p-6 border-b border-white/10">
            <h3 class="text-lg font-semibold text-white">合併帳號</h3>
            <p class="text-sm text-slate-400 mt-1">
              將佔位老師的資料合併至已綁定老師帳號
            </p>
          </div>

          <div v-if="mergingTeacher" class="p-6 space-y-4">
            <!-- 來源老師（佔位） -->
            <div class="p-4 rounded-lg bg-warning-500/10 border border-warning-500/30">
              <p class="text-sm text-warning-400 mb-2">來源（將被合併）</p>
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-warning-500/20 flex items-center justify-center">
                  <span class="text-warning-400 font-medium">
                    {{ mergingTeacher.name.charAt(0) }}
                  </span>
                </div>
                <div>
                  <p class="text-white font-medium">{{ mergingTeacher.name }}</p>
                  <p class="text-sm text-slate-400">{{ mergingTeacher.email || '無 Email' }}</p>
                </div>
              </div>
              <div class="mt-2 flex items-center gap-1.5 text-warning-400 text-xs">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <span>佔位老師，合併後將被刪除</span>
              </div>
            </div>

            <!-- 目標老師選擇 -->
            <div>
              <label class="block text-sm text-slate-400 mb-2">
                選擇目標老師 <span class="text-critical-500">*</span>
              </label>
              <select
                v-model="mergeForm.targetTeacherId"
                class="w-full px-3 py-2 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500"
              >
                <option value="">請選擇目標老師</option>
                <option
                  v-for="teacher in boundTeachers"
                  :key="teacher.id"
                  :value="teacher.id"
                  :disabled="teacher.id === mergingTeacher?.id"
                >
                  {{ teacher.name }} ({{ teacher.email || '無 Email' }})
                </option>
              </select>
              <p class="mt-1 text-xs text-slate-500">
                選擇要接收資料的已綁定老師帳號
              </p>
            </div>

            <!-- 合併說明 -->
            <div class="p-3 rounded-lg bg-primary-500/10 border border-primary-500/30">
              <p class="text-sm text-primary-400">
                合併將會：
              </p>
              <ul class="mt-2 text-xs text-slate-400 space-y-1">
                <li>- 遷移所有課表、例外記錄、課程筆記</li>
                <li>- 遷移私人行程和技能證照</li>
                <li>- 刪除來源佔位老師帳號</li>
                <li>- 此操作無法復原</li>
              </ul>
            </div>
          </div>

          <div class="p-6 border-t border-white/10 flex items-center gap-4">
            <button
              @click="closeMergeModal"
              class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              取消
            </button>
            <button
              @click="executeMerge"
              :disabled="!isMergeValid || merging"
              class="flex-1 px-4 py-2 rounded-lg bg-warning-500/30 border border-warning-500 text-warning-400 hover:bg-warning-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="merging" class="flex items-center justify-center gap-2">
                <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                處理中...
              </span>
              <span v-else>確認合併</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 通知組件 -->
    <NotificationDropdown
      v-if="notificationUI.show.value"
      @close="notificationUI.close()"
    />
  </div>
</template>

<script setup lang="ts">
import BaseGlassCard from '~/components/base/BaseGlassCard.vue'
import BaseInput from '~/components/base/BaseInput.vue'
import NotificationDropdown from '~/components/Navigation/NotificationDropdown.vue'
import type { Teacher } from '~/types/teacher'

definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

const notificationUI = useNotification()
const { confirm: alertConfirm } = useAlert()

// API 實例
const api = useApi()

// 老師列表
const teachers = ref<Teacher[]>([])
const loading = ref(true)

// 篩選條件
const searchQuery = ref('')
const filterPlaceholder = ref('')

// 篩選後的老師列表
const filteredTeachers = computed(() => {
  let result = teachers.value

  // 搜尋過濾
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(t =>
      t.name.toLowerCase().includes(query) ||
      (t.email && t.email.toLowerCase().includes(query))
    )
  }

  // 佔位狀態過濾
  if (filterPlaceholder.value) {
    if (filterPlaceholder.value === 'placeholder') {
      result = result.filter(t => t.is_placeholder)
    } else if (filterPlaceholder.value === 'bound') {
      result = result.filter(t => !t.is_placeholder)
    }
  }

  return result
})

// 已綁定的老師列表（用於合併目標選擇）
const boundTeachers = computed(() => {
  return teachers.value.filter(t => !t.is_placeholder)
})

// 清除篩選
const clearFilters = () => {
  searchQuery.value = ''
  filterPlaceholder.value = ''
}

// 取得老師列表
const fetchTeachers = async () => {
  loading.value = true
  try {
    const response = await api.get<Teacher[]>('/teachers')
    teachers.value = response || []
  } catch (error) {
    console.error('取得老師列表失敗:', error)
    notificationUI.error('載入老師列表失敗')
    teachers.value = []
  } finally {
    loading.value = false
  }
}

// 新增佔位老師 Modal
const showCreateModal = ref(false)
const createForm = ref({
  name: '',
  email: '',
})
const creating = ref(false)

// 驗證新增表單
const isCreateValid = computed(() => {
  return createForm.value.name.trim().length > 0
})

// 新增佔位老師
const createPlaceholderTeacher = async () => {
  if (!isCreateValid.value) return

  creating.value = true
  try {
    await api.post('/admin/teachers/placeholder', {
      name: createForm.value.name.trim(),
      email: createForm.value.email.trim() || undefined,
    })

    notificationUI.success('佔位老師已建立')
    showCreateModal.value = false
    createForm.value = { name: '', email: '' }
    await fetchTeachers()
  } catch (error) {
    console.error('建立佔位老師失敗:', error)
    notificationUI.error('建立失敗，請稍後再試')
  } finally {
    creating.value = false
  }
}

// 合併帳號 Modal
const showMergeModal = ref(false)
const mergingTeacher = ref<Teacher | null>(null)
const mergeForm = ref({
  targetTeacherId: '',
})
const merging = ref(false)

// 開啟合併 Modal
const openMergeModal = (teacher: Teacher) => {
  mergingTeacher.value = teacher
  mergeForm.value.targetTeacherId = ''
  showMergeModal.value = true
}

// 關閉合併 Modal
const closeMergeModal = () => {
  showMergeModal.value = false
  mergingTeacher.value = null
  mergeForm.value.targetTeacherId = ''
}

// 驗證合併表單
const isMergeValid = computed(() => {
  return mergingTeacher.value !== null && mergeForm.value.targetTeacherId !== ''
})

// 執行合併
const executeMerge = async () => {
  if (!isMergeValid.value || !mergingTeacher.value) return

  const targetId = parseInt(mergeForm.value.targetTeacherId)

  // 確認對話
  if (!await alertConfirm(
    `確定要將「${mergingTeacher.value.name}」合併至目標老師嗎？\n\n此操作將遷移所有資料並刪除佔位帳號，無法復原。`
  )) {
    return
  }

  merging.value = true
  try {
    await api.post('/admin/teachers/merge', {
      source_teacher_id: mergingTeacher.value.id,
      target_teacher_id: targetId,
    })

    notificationUI.success('帳號合併成功')
    closeMergeModal()
    await fetchTeachers()
  } catch (error) {
    console.error('合併帳號失敗:', error)
    notificationUI.error('合併失敗，請稍後再試')
  } finally {
    merging.value = false
  }
}

// 頁面載入時取得資料
onMounted(() => {
  fetchTeachers()
})
</script>
