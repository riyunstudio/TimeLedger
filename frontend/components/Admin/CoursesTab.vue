<template>
  <div class="space-y-4">
    <!-- 標題區域 -->
    <div class="flex flex-col gap-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <h2 class="text-xl font-semibold text-slate-100">課程列表</h2>
          <span class="text-sm text-slate-500">({{ pagination.total }})</span>
        </div>
        <div class="flex items-center gap-2">
          <button
            @click="fetchCourses"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
            title="重新整理"
          >
            <svg class="w-5 h-5 text-slate-400" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
          <button
            @click="openCreateModal"
            class="btn-primary px-4 py-2 text-sm font-medium"
          >
            + 新增課程
          </button>
        </div>
      </div>

      <!-- 搜尋框 -->
      <div class="relative">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜尋課程名稱..."
          class="w-full pl-10 pr-4 py-2 bg-white/5 border border-white/10 rounded-lg text-slate-100 placeholder-slate-500 focus:outline-none focus:border-primary-500 focus:ring-1 focus:ring-primary-500 transition-colors"
        />
      </div>
    </div>

    <!-- 分頁控制 -->
    <div v-if="pagination.totalPages > 1" class="flex items-center justify-between px-2">
      <span class="text-sm text-slate-500">
        第 {{ pagination.currentPage }} 頁 / 共 {{ pagination.totalPages }} 頁
      </span>
      <div class="flex items-center gap-2">
        <button
          @click="goToPage(pagination.currentPage - 1)"
          :disabled="pagination.currentPage <= 1"
          class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors"
          :class="pagination.currentPage <= 1
            ? 'bg-white/5 text-slate-600 cursor-not-allowed'
            : 'bg-white/10 text-slate-300 hover:bg-white/20'"
        >
          上一頁
        </button>
        <button
          @click="goToPage(pagination.currentPage + 1)"
          :disabled="pagination.currentPage >= pagination.totalPages"
          class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors"
          :class="pagination.currentPage >= pagination.totalPages
            ? 'bg-white/5 text-slate-600 cursor-not-allowed'
            : 'bg-white/10 text-slate-300 hover:bg-white/20'"
        >
          下一頁
        </button>
      </div>
    </div>

    <!-- 骨架屏載入狀態 -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="i in 6"
        :key="i"
        class="glass-card p-5"
      >
        <div class="animate-pulse">
          <div class="flex items-start justify-between mb-3">
            <div class="flex-1">
              <div class="h-5 w-32 bg-white/10 rounded mb-2"></div>
              <div class="h-4 w-20 bg-white/10 rounded"></div>
            </div>
            <div class="flex gap-2">
              <div class="w-8 h-8 bg-white/10 rounded-lg"></div>
              <div class="w-8 h-8 bg-white/10 rounded-lg"></div>
            </div>
          </div>
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <div class="w-5 h-5 bg-white/10 rounded"></div>
              <div class="h-4 w-20 bg-white/10 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空狀態 -->
    <div v-else-if="courses.length === 0" class="glass-card p-12 text-center">
      <svg class="w-16 h-16 mx-auto mb-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
      </svg>
      <h3 class="text-lg font-medium text-slate-300 mb-2">
        {{ searchQuery ? '沒有符合的課程' : '尚無課程' }}
      </h3>
      <p class="text-slate-500 mb-4">
        {{ searchQuery ? '嘗試其他搜尋關鍵字' : '建立第一個課程來開始使用' }}
      </p>
      <button
        v-if="!searchQuery"
        @click="showCreateModal = true"
        class="btn-primary px-4 py-2"
      >
        + 新增課程
      </button>
    </div>

    <!-- 課程列表 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="course in courses"
        :key="course.id"
        class="glass-card p-5 hover:bg-white/5 transition-colors"
      >
        <div class="flex items-start justify-between mb-3">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span
                v-if="course.code"
                class="px-2 py-0.5 text-xs rounded bg-primary-500/20 text-primary-400 font-mono"
              >
                {{ course.code }}
              </span>
              <h3 class="text-lg font-medium text-slate-100">{{ course.name }}</h3>
            </div>
          </div>
          <div class="flex gap-1">
            <button
              @click="toggleCourseActive(course)"
              class="p-1.5 rounded-lg transition-colors"
              :class="course.is_active
                ? 'hover:bg-green-500/20 text-green-400'
                : 'hover:bg-slate-500/20 text-slate-500'"
              :title="course.is_active ? '停用課程' : '啟用課程'"
            >
              <svg class="w-4 h-4" :class="{ 'opacity-50': togglingId === course.id }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="course.is_active" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
            <button
              @click="editCourse(course)"
              class="p-1.5 rounded-lg hover:bg-white/10 transition-colors"
              title="編輯"
            >
              <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </button>
            <button
              @click="deleteCourse(course)"
              class="p-1.5 rounded-lg hover:bg-red-500/20 transition-colors"
              title="刪除"
            >
              <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>

        <div class="space-y-2 text-sm">
          <div class="flex items-center gap-2 text-slate-400">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{{ course.default_duration || 60 }} 分鐘</span>
          </div>
        </div>

        <div class="mt-4 pt-3 border-t border-white/10 flex items-center justify-between">
          <span
            class="px-2 py-1 text-xs rounded"
            :class="course.is_active ? 'bg-green-500/20 text-green-400' : 'bg-slate-500/20 text-slate-400'"
          >
            {{ course.is_active ? '啟用中' : '已停用' }}
          </span>
          <span class="text-xs text-slate-500">
            ID: {{ course.id }}
          </span>
        </div>
      </div>
    </div>

    <!-- 分頁控制 -->
    <div v-if="pagination.totalPages > 1" class="flex items-center justify-between px-2">
      <span class="text-sm text-slate-500">
        第 {{ pagination.currentPage }} 頁 / 共 {{ pagination.totalPages }} 頁
      </span>
      <div class="flex items-center gap-2">
        <button
          @click="goToPage(pagination.currentPage - 1)"
          :disabled="pagination.currentPage <= 1"
          class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors"
          :class="pagination.currentPage <= 1
            ? 'bg-white/5 text-slate-600 cursor-not-allowed'
            : 'bg-white/10 text-slate-300 hover:bg-white/20'"
        >
          上一頁
        </button>
        <button
          @click="goToPage(pagination.currentPage + 1)"
          :disabled="pagination.currentPage >= pagination.totalPages"
          class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors"
          :class="pagination.currentPage >= pagination.totalPages
            ? 'bg-white/5 text-slate-600 cursor-not-allowed'
            : 'bg-white/10 text-slate-300 hover:bg-white/20'"
        >
          下一頁
        </button>
      </div>
    </div>

    <!-- 新增/編輯課程 Modal -->
    <div
      v-if="showCreateModal || editingCourse"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      @click.self="closeModal"
    >
      <div class="glass-card w-full max-w-md p-6">
        <h3 class="text-xl font-semibold text-slate-100 mb-6">
          {{ editingCourse ? '編輯課程' : '新增課程' }}
        </h3>

        <form @submit.prevent="saveCourse" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-2">
                課程代號
              </label>
              <input
                v-model="form.code"
                type="text"
                class="input-field"
                placeholder="例：Piano-101"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-2">
                課程名稱 <span class="text-red-400">*</span>
              </label>
              <input
                v-model="form.name"
                type="text"
                class="input-field"
                placeholder="輸入課程名稱"
                required
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              預設時長（分鐘）
            </label>
            <input
              v-model.number="form.default_duration"
              type="number"
              class="input-field"
              placeholder="60"
              min="1"
            />
          </div>

          <div class="flex items-center gap-2">
            <input
              type="checkbox"
              id="is_active"
              v-model="form.is_active"
              class="w-4 h-4 rounded border-slate-600 bg-slate-700 text-primary-500 focus:ring-primary-500"
            />
            <label for="is_active" class="text-sm text-slate-300">
              啟用此課程
            </label>
          </div>

          <div class="flex justify-end gap-3 mt-6">
            <button
              type="button"
              @click="closeModal"
              class="px-4 py-2 text-slate-300 hover:text-white transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              class="btn-primary px-6"
              :disabled="saving"
            >
              {{ saving ? '儲存中...' : '儲存' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertError, alertSuccess } from '~/composables/useAlert'

interface Course {
  id: number
  code: string
  name: string
  default_duration: number
  is_active: boolean
}

interface PaginationState {
  currentPage: number
  totalPages: number
  total: number
  limit: number
}

interface CenterSettings {
  default_course_duration: number
}

const api = useApi()

const courses = ref<Course[]>([])
const loading = ref(false)
const saving = ref(false)
const togglingId = ref<number | null>(null)
const showCreateModal = ref(false)
const editingCourse = ref<Course | null>(null)
const centerSettings = ref<CenterSettings | null>(null)

// 搜尋與分頁狀態
const searchQuery = ref('')
const pagination = ref<PaginationState>({
  currentPage: 1,
  totalPages: 1,
  total: 0,
  limit: 12,
})
const debounceTimer = ref<NodeJS.Timeout | null>(null)

const form = reactive({
  code: '',
  name: '',
  default_duration: 60,
  is_active: true,
})

// 取得中心設定
async function fetchCenterSettings() {
  try {
    const { getCenterId } = useCenterId()
    const centerId = getCenterId()
    if (!centerId) return

    const response = await api.get<CenterSettings>(`/admin/centers/${centerId}/settings`)
    centerSettings.value = response
    // 如果是新增課程，使用中心預設時長
    if (!editingCourse.value) {
      form.default_duration = response.default_course_duration || 60
    }
  } catch (error) {
    console.error('取得中心設定失敗:', error)
    // 使用預設值
    form.default_duration = 60
  }
}

// 打開新增課程 Modal
function openCreateModal() {
  editingCourse.value = null
  resetForm()
  fetchCenterSettings()
  showCreateModal.value = true
}

// Debounce 搜尋
function updateSearch() {
  if (debounceTimer.value) {
    clearTimeout(debounceTimer.value)
  }
  debounceTimer.value = setTimeout(() => {
    pagination.value.currentPage = 1 // 搜尋時重置到第一頁
    fetchCourses()
  }, 300) // 300ms debounce
}

// 監聽搜尋文字變化
watch(searchQuery, () => {
  updateSearch()
})

async function fetchCourses() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (searchQuery.value.trim()) {
      params.set('query', searchQuery.value.trim())
    }
    params.set('page', String(pagination.value.currentPage))
    params.set('limit', String(pagination.value.limit))

    const queryString = params.toString()
    const url = `/admin/courses${queryString ? `?${queryString}` : ''}`

    const response = await api.get<{ data: Course[]; total: number; total_pages: number }>(url)
    courses.value = response.data || []

    // 更新分頁資訊
    pagination.value.total = response.total || 0
    pagination.value.totalPages = response.total_pages || 1
  } catch (error) {
    console.error('取得課程列表失敗:', error)
    courses.value = []
  } finally {
    loading.value = false
  }
}

function goToPage(page: number) {
  if (page >= 1 && page <= pagination.value.totalPages) {
    pagination.value.currentPage = page
    fetchCourses()
  }
}

function editCourse(course: Course) {
  editingCourse.value = course
  form.code = course.code || ''
  form.name = course.name
  form.default_duration = course.default_duration || 60
  form.is_active = course.is_active
}

async function deleteCourse(course: Course) {
  if (!await alertConfirm(`確定要刪除課程「${course.name}」嗎？`)) {
    return
  }

  try {
    await api.delete(`/admin/courses/${course.id}`)
    await alertSuccess('課程已刪除')
    courses.value = []
    await fetchCourses()
  } catch (error) {
    console.error('刪除課程失敗:', error)
    await alertError('刪除課程失敗')
  }
}

async function toggleCourseActive(course: Course) {
  if (togglingId.value === course.id) return

  togglingId.value = course.id
  try {
    const newStatus = !course.is_active
    await api.patch(`/admin/courses/${course.id}/toggle-active`, {
      is_active: newStatus,
    })
    course.is_active = newStatus
    await alertSuccess(newStatus ? '課程已啟用' : '課程已停用')
  } catch (error) {
    console.error('切換課程狀態失敗:', error)
    await alertError('操作失敗，請稍後再試')
  } finally {
    togglingId.value = null
  }
}

function closeModal() {
  showCreateModal.value = false
  editingCourse.value = null
  resetForm()
}

function resetForm() {
  form.code = ''
  form.name = ''
  form.default_duration = centerSettings.value?.default_course_duration || 60
  form.is_active = true
}

async function saveCourse() {
  if (!form.name.trim()) {
    await alertError('請輸入課程名稱')
    return
  }

  saving.value = true
  try {
    const data = {
      code: form.code.trim() || null,
      name: form.name,
      duration: form.default_duration,
      is_active: form.is_active,
      color_hex: '#3b82f6',
    }

    if (editingCourse.value) {
      await api.put(`/admin/courses/${editingCourse.value.id}`, data)
      await alertSuccess('課程已更新')
    } else {
      await api.post('/admin/courses', data)
      await alertSuccess('課程已建立')
    }

    courses.value = []
    closeModal()
    await fetchCourses()
  } catch (error) {
    console.error('儲存課程失敗:', error)
    await alertError('儲存課程失敗')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchCourses()
})

// 元件卸載時清除計時器
onUnmounted(() => {
  if (debounceTimer.value) {
    clearTimeout(debounceTimer.value)
  }
})
</script>
