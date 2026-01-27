<template>
  <div class="space-y-4">
    <!-- 標題區域 -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <h2 class="text-xl font-semibold text-slate-100">課程列表</h2>
        <span class="text-sm text-slate-500">({{ courses.length }})</span>
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
          @click="showCreateModal = true"
          class="btn-primary px-4 py-2 text-sm font-medium"
        >
          + 新增課程
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
            <div class="flex items-center gap-2">
              <div class="w-5 h-5 bg-white/10 rounded"></div>
              <div class="h-4 w-20 bg-white/10 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空狀態 -->
    <div v-else-if="courses.length === 0" class="text-center py-12 text-slate-500 glass-card">
      <svg class="w-16 h-16 mx-auto text-slate-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
      </svg>
      <p class="text-lg mb-2">尚未添加課程</p>
      <p class="text-sm text-slate-500">點擊上方按鈕新增第一個課程</p>
    </div>

    <!-- 課程列表 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="course in courses"
          :key="course.id"
          class="glass-card p-5 hover:bg-white/5 transition-all"
        >
          <div class="flex items-start justify-between mb-3">
            <div class="flex-1">
              <h3 class="text-lg font-semibold text-slate-100">{{ course.name }}</h3>
              <p class="text-sm text-slate-400 mt-1">Course ID: {{ course.id }}</p>
            </div>
            <div class="flex items-center gap-1">
              <button
                @click="editCourse(course)"
                class="p-2 rounded-lg hover:bg-white/10 transition-colors"
                title="編輯"
              >
                <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 0L21.828 5.172a2 2 0 010-2.828l-7.414-7.414a2 2 0 00-2.828 0L2.172 15.828a2 2 0 010 2.828l7.414 7.414a2 2 0 002.828 0z" />
                </svg>
              </button>
              <button
                @click="deleteCourse(course)"
                class="p-2 rounded-lg hover:bg-red-500/20 transition-colors"
                title="刪除"
              >
                <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>

        <div class="space-y-2">
          <div class="flex items-center gap-2">
            <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0zm-9 9h.01" />
            </svg>
            <span class="text-sm text-slate-300">老師緩衝</span>
            <span class="text-sm font-medium text-slate-100">{{ course.teacher_buffer_min }} 分鐘</span>
          </div>
          <div class="flex items-center gap-2">
            <svg class="w-5 h-5 text-secondary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H9" />
            </svg>
            <span class="text-sm text-slate-300">教室緩衝</span>
            <span class="text-sm font-medium text-slate-100">{{ course.room_buffer_min }} 分鐘</span>
          </div>
        </div>
      </div>
    </div>

    <CourseModal
      v-if="showCreateModal"
      :course="null"
      @close="showCreateModal = false"
      @saved="fetchCourses"
    />

    <CourseModal
      v-if="showEditModal"
      :course="editingCourse"
      @close="showEditModal = false"
      @saved="fetchCourses"
    />
  </div>
</template>

<script setup lang="ts">
import { alertError, alertConfirm, alertSuccess } from '~/composables/useAlert'

const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingCourse = ref<any>(null)
const loading = ref(false)
const { getCenterId } = useCenterId()

const courses = ref<any[]>([])

const fetchCourses = async () => {
  loading.value = true
  try {
    const api = useApi()
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/courses`)
    courses.value = response.datas || []
  } catch (error: any) {
    console.error('Failed to fetch courses:', error)
    courses.value = []
    const errorMessage = error?.response?.data?.message || error?.message || '無法載入課程列表，請檢查網路連線'
    await alertError(errorMessage)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCourses()
})

const editCourse = (course: any) => {
  editingCourse.value = { ...course }
  showEditModal.value = true
}

const deleteCourse = async (course: any) => {
  if (!await alertConfirm(`確定要刪除課程「${course.name}」嗎？`)) {
    return
  }

  try {
    const api = useApi()
    const centerId = getCenterId()
    await api.delete(`/admin/courses/${course.id}`)
    await fetchCourses()
  } catch (err) {
    console.error('Failed to delete course:', err)
    await alertError('刪除失敗，請稍後再試')
  }
}
</script>
