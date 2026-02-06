<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-slate-100">老師列表</h2>
      <div class="flex items-center gap-2">
        <button
          @click="showCreateModal = true"
          class="glass-btn px-4 py-2 text-sm font-medium"
        >
          手動新增老師
        </button>
        <button
          @click="showInviteModal = true"
          class="btn-primary px-4 py-2 text-sm font-medium"
        >
          + 邀請老師
        </button>
      </div>
    </div>

    <div class="mb-4">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜尋老師..."
        class="input-field"
        :disabled="loading"
        @input="debouncedSearch"
      />
    </div>

    <!-- 骨架屏載入狀態 -->
    <div v-if="loading && teachers.length === 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="i in 6"
        :key="i"
        class="glass-card p-5"
      >
        <div class="animate-pulse">
          <div class="flex items-start justify-between mb-3">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-white/10"></div>
              <div>
                <div class="h-5 w-32 bg-white/10 rounded mb-2"></div>
                <div class="h-4 w-48 bg-white/10 rounded"></div>
              </div>
            </div>
            <div class="w-8 h-8 bg-white/10 rounded-lg"></div>
          </div>
          <div class="space-y-2">
            <div class="flex gap-2">
              <div class="h-6 w-16 bg-white/10 rounded-full"></div>
              <div class="h-6 w-16 bg-white/10 rounded-full"></div>
              <div class="h-6 w-16 bg-white/10 rounded-full"></div>
            </div>
            <div class="flex items-center gap-2">
              <div class="h-5 w-20 bg-white/10 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空狀態 -->
    <div v-else-if="teachers.length === 0" class="text-center py-12 text-slate-500 glass-card">
      <svg class="w-16 h-16 mx-auto mb-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20H2v-2a3 3 0 015.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
      <p>{{ searchQuery ? '找不到符合的老師' : '尚未添加老師' }}</p>
    </div>

    <!-- 老師列表 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="teacher in teachers"
        :key="teacher.id"
        class="glass-card p-5 hover:bg-white/5 transition-all"
      >
        <div class="flex items-start justify-between mb-3">
          <div class="flex items-center gap-3">
            <div
              class="w-12 h-12 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center"
            >
              <span class="text-lg font-bold text-white">{{ teacher.name?.charAt(0) }}</span>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-slate-100">{{ teacher.name }}</h3>
              <p class="text-sm text-slate-400">{{ teacher.email }}</p>
            </div>
          </div>
          <button
            @click="showMenu[teacher.id] = !showMenu[teacher.id]"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 19v.01M7 5h10M7 19h10" />
            </svg>
          </button>
        </div>

        <div class="space-y-2">
          <div v-if="teacher.skills" class="flex flex-wrap gap-1">
            <span
              v-for="(skill, index) in teacher.skills.slice(0, 3)"
              :key="index"
              class="px-2 py-1 rounded-full text-xs font-medium bg-secondary-500/20 text-secondary-500"
            >
              {{ skill.skill_name }}
            </span>
            <span v-if="teacher.skills.length > 3" class="text-xs text-slate-400">
              +{{ teacher.skills.length - 3 }} 更多
            </span>
          </div>

          <div v-if="teacher.certificates && teacher.certificates.length > 0" class="flex items-center gap-2 text-sm text-slate-400">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
            </svg>
            {{ teacher.certificates.length }} 張證照
          </div>

          <!-- 評分顯示 -->
          <div class="flex items-center gap-2">
            <div class="flex items-center gap-0.5">
              <template v-for="star in 5" :key="star">
                <svg
                  class="w-4 h-4"
                  :class="star <= (teacher.note?.rating || 0) ? 'text-warning-500' : 'text-slate-600'"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              </template>
            </div>
            <span class="text-xs text-slate-500">
              {{ teacher.note?.rating || 0 }}
            </span>
          </div>

          <div class="flex items-center gap-2">
            <span class="px-2 py-1 rounded-full text-xs font-medium"
              :class="teacher.is_active ? 'bg-success-500/20 text-success-500' : 'bg-warning-500/20 text-warning-500'"
            >
              {{ teacher.is_active ? '活躍中' : '未活躍' }}
            </span>
          </div>
        </div>

        <div v-if="showMenu[teacher.id]" class="mt-3 pt-3 border-t border-white/10 space-y-2">
          <button
            @click="viewProfile(teacher)"
            class="w-full glass-btn px-4 py-2 rounded-xl text-sm"
          >
            查看個人檔案
          </button>
          <!-- 發送訊息功能暫時隱藏
          <button
            @click="sendMessage(teacher)"
            class="w-full glass-btn px-4 py-2 rounded-xl text-sm"
          >
            發送訊息
          </button>
          -->
          <button
            v-if="teacher.is_placeholder"
            @click="openMergeModal(teacher)"
            class="w-full glass-btn px-4 py-2 rounded-xl text-sm text-primary-400 hover:bg-primary-500/10"
          >
            合併帳號
          </button>
          <button
            @click="removeTeacher(teacher)"
            class="w-full glass-btn px-4 py-2 rounded-xl text-sm text-red-400 hover:bg-red-500/10"
          >
            移除老師
          </button>
          <button
            @click="openRatingModal(teacher)"
            class="w-full glass-btn px-4 py-2 rounded-xl text-sm text-warning-400 hover:bg-warning-500/10"
          >
            編輯評分
          </button>
        </div>
      </div>
    </div>

    <!-- 分頁控制 -->
    <div v-if="totalPages > 1" class="flex items-center justify-between pt-4">
      <div class="text-sm text-slate-400">
        共 {{ totalCount }} 位老師
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage === 1"
          class="glass-btn px-3 py-1.5 text-sm disabled:opacity-50 disabled:cursor-not-allowed"
        >
          上一頁
        </button>
        <span class="text-sm text-slate-300">
          第 {{ currentPage }} / {{ totalPages }} 頁
        </span>
        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="glass-btn px-3 py-1.5 text-sm disabled:opacity-50 disabled:cursor-not-allowed"
        >
          下一頁
        </button>
      </div>
    </div>

    <TeacherInviteModal
      v-if="showInviteModal"
      @close="showInviteModal = false"
      @invited="fetchTeachers"
    />

    <TeacherCreateModal
      :show="showCreateModal"
      @close="showCreateModal = false"
      @created="fetchTeachers"
    />

    <AdminTeacherProfileModal
      v-if="selectedTeacher"
      :teacher="selectedTeacher"
      @close="selectedTeacher = null"
    />

    <TeacherMergeModal
      :show="showMergeModal"
      :teacher="mergingTeacher"
      :bound-teachers="allTeachers.filter(t => !t.is_placeholder)"
      @close="closeMergeModal"
      @merged="onMerged"
    />

    <TeacherRatingModal
      :show="showRatingModal"
      :teacher="ratingTeacher"
      @close="closeRatingModal"
      @saved="onRatingSaved"
      @deleted="onRatingSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useDebounceFn } from '@vueuse/core'

// Alert composable
const { error: alertError, confirm: alertConfirm } = useAlert()

// 資源快取
const { invalidate } = useResourceCache()

import TeacherCreateModal from '~/components/Teacher/TeacherCreateModal.vue'
import TeacherMergeModal from '~/components/Teacher/TeacherMergeModal.vue'
import TeacherRatingModal from '~/components/Teacher/TeacherRatingModal.vue'

// 分頁常數
const PAGE_LIMIT = 20

const showInviteModal = ref(false)
const showCreateModal = ref(false)
const searchQuery = ref('')
const showMenu = ref<Record<number, boolean>>({})
const loading = ref(false)
const selectedTeacher = ref<any>(null)

// 分頁狀態
const currentPage = ref(1)
const totalPages = ref(0)
const totalCount = ref(0)

// 合併狀態
const showMergeModal = ref(false)
const mergingTeacher = ref<any>(null)
const allTeachers = ref<any[]>([]) // 用於選擇目標老師 (通常需要全部已綁定的老師)

// 評分狀態
const showRatingModal = ref(false)
const ratingTeacher = ref<any>(null)

const teachers = ref<any[]>([])

// 去抖動搜尋函數
const debouncedSearch = useDebounceFn(() => {
  currentPage.value = 1 // 搜尋時回到第一頁
  fetchTeachers()
}, 300)

// 帶搜尋和分頁的 API 請求
const fetchTeachers = async () => {
  loading.value = true
  try {
    const api = useApi()
    const params = new URLSearchParams()

    // 伺服器端搜尋
    if (searchQuery.value.trim()) {
      params.append('q', searchQuery.value.trim())
    }

    // 分頁參數
    params.append('page', currentPage.value.toString())
    params.append('limit', PAGE_LIMIT.toString())

    const response = await api.get<any>(`/admin/teachers?${params.toString()}`)

    // API 現在直接回傳分頁格式
    if (response && response.data) {
      teachers.value = response.data
    }

    // 更新分頁資訊
    if (response) {
      totalCount.value = response.total || 0
      totalPages.value = response.total_pages || 1
    }
  } catch (error) {
    console.error('Failed to fetch teachers:', error)
    teachers.value = []
    totalCount.value = 0
    totalPages.value = 0
  } finally {
    loading.value = false
  }
}

// 獲取所有老師 (不分頁，用於合併時選擇目標)
const fetchAllTeachers = async () => {
  try {
    const api = useApi()
    const response = await api.get<any>('/admin/teachers?limit=1000') // 假設不超過 1000 位
    if (response && response.data) {
      allTeachers.value = response.data
    } else {
      allTeachers.value = Array.isArray(response) ? response : []
    }
  } catch (err) {
    console.error('Failed to fetch all teachers:', err)
  }
}

// 跳轉到指定頁面
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchTeachers()
  }
}

// 監聽搜尋變化
watch(searchQuery, () => {
  // 搜尋時會由 debouncedSearch 處理
  // 這裡不需要額外邏輯
})

onMounted(() => {
  fetchTeachers()
  fetchAllTeachers()
})

const viewProfile = (teacher: any) => {
  selectedTeacher.value = teacher
  showMenu.value[teacher.id] = false
}

const sendMessage = (teacher: any) => {
}

const openMergeModal = (teacher: any) => {
  mergingTeacher.value = teacher
  showMergeModal.value = true
  showMenu.value[teacher.id] = false
}

const closeMergeModal = () => {
  showMergeModal.value = false
  mergingTeacher.value = null
}

const onMerged = async () => {
  await fetchTeachers()
  await fetchAllTeachers()
  invalidate('teachers')
}

const openRatingModal = (teacher: any) => {
  ratingTeacher.value = teacher
  showRatingModal.value = true
  showMenu.value[teacher.id] = false
}

const closeRatingModal = () => {
  showRatingModal.value = false
  ratingTeacher.value = null
}

const onRatingSaved = async () => {
  await fetchTeachers()
}

const removeTeacher = async (teacher: any) => {
  const authStore = useAuthStore()
  const centerId = (authStore.user as any)?.center_id

  if (!centerId) {
    await alertError('未找到中心資訊，請重新登入')
    return
  }

  if (!await alertConfirm(`確定要將老師「${teacher.name}」從此中心移除嗎？`)) {
    return
  }

  try {
    const api = useApi()
    await api.delete(`/admin/centers/${centerId}/teachers/${teacher.id}`)

    // 移除成功後，如果目前頁面沒有資料了，回到前一頁
    if (teachers.value.length === 1 && currentPage.value > 1) {
      currentPage.value--
    }

    await fetchTeachers()
    // 清除老師快取，下次存取會自動重新載入
    invalidate('teachers')
    showMenu.value[teacher.id] = false
  } catch (err) {
    console.error('Failed to remove teacher:', err)
    await alertError('移除失敗，請稍後再試')
  }
}
</script>
