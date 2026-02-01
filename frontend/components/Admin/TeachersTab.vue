<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-slate-100">老師列表</h2>
      <button
        @click="showInviteModal = true"
        class="btn-primary px-4 py-2 text-sm font-medium"
      >
        + 邀請老師
      </button>
    </div>

    <div class="mb-4">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜尋老師..."
        class="input-field"
        :disabled="loading"
      />
    </div>

    <!-- 骨架屏載入狀態 -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 gap-4">
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
    <div v-else-if="filteredTeachers.length === 0" class="text-center py-12 text-slate-500 glass-card">
      <svg class="w-16 h-16 mx-auto mb-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20H2v-2a3 3 0 015.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
      <p>{{ searchQuery ? '找不到符合的老師' : '尚未添加老師' }}</p>
    </div>

    <!-- 老師列表 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="teacher in filteredTeachers"
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
            @click="removeTeacher(teacher)"
            class="w-full glass-btn px-4 py-2 rounded-xl text-sm text-red-400 hover:bg-red-500/10"
          >
            移除老師
          </button>
        </div>
      </div>
    </div>

    <TeacherInviteModal
      v-if="showInviteModal"
      @close="showInviteModal = false"
      @invited="fetchTeachers"
    />

    <AdminTeacherProfileModal
      v-if="selectedTeacher"
      :teacher="selectedTeacher"
      @close="selectedTeacher = null"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// Alert composable
const { error: alertError, confirm: alertConfirm } = useAlert()

// 資源快取
const { invalidate } = useResourceCache()

const showInviteModal = ref(false)
const searchQuery = ref('')
const showMenu = ref<Record<number, boolean>>({})
const loading = ref(false)
const selectedTeacher = ref<any>(null)

const teachers = ref<any[]>([])

const filteredTeachers = computed(() => {
  if (!searchQuery.value) return teachers.value

  const query = searchQuery.value.toLowerCase()
  return teachers.value.filter(t =>
    t.name?.toLowerCase().includes(query) ||
    t.email?.toLowerCase().includes(query)
  )
})

const fetchTeachers = async () => {
  loading.value = true
  try {
    const api = useApi()
    const result = await api.get<any[]>('/admin/teachers')
    teachers.value = result

    // 背景抓取每位老師的評分資料以填充 UI 中的星星
    await Promise.all(
      teachers.value.map(async (teacher) => {
        try {
          const noteResponse = await api.get<any>(
            `/admin/teachers/${teacher.id}/note`
          )
          if (noteResponse) {
            teacher.note = noteResponse
          }
        } catch {
          // 無評分資料為正常情況
        }
      })
    )
  } catch (error) {
    console.error('Failed to fetch teachers:', error)
    teachers.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchTeachers()
})

const viewProfile = (teacher: any) => {
  selectedTeacher.value = teacher
  showMenu.value[teacher.id] = false
}

const sendMessage = (teacher: any) => {
}

const removeTeacher = async (teacher: any) => {
  if (!await alertConfirm(`確定要移除老師「${teacher.name}」嗎？此操作將刪除該老師的帳號。`)) {
    return
  }

  try {
    const api = useApi()
    await api.delete(`/teachers/${teacher.id}`)
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
