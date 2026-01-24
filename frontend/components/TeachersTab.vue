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
      />
    </div>

    <div v-if="filteredTeachers.length === 0" class="text-center py-12 text-slate-500 glass-card">
      {{ searchQuery ? '找不到符合的老師' : '尚未添加老師' }}
    </div>

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

          <div v-if="teacher.certificates" class="text-sm text-slate-400">
            {{ Array.isArray(teacher.certificates) ? teacher.certificates.length : teacher.certificates }} 個證照
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
          <button
            @click="sendMessage(teacher)"
            class="w-full glass-btn px-4 py-2 rounded-xl text-sm"
          >
            發送訊息
          </button>
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
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// Alert composable
const { error: alertError, confirm: alertConfirm } = useAlert()

const showInviteModal = ref(false)
const searchQuery = ref('')
const showMenu = ref<Record<number, boolean>>({})
const loading = ref(false)

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
    const response = await api.get<{ code: number; datas: any[] }>('/teachers')
    teachers.value = response.datas || []
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
  console.log('View profile:', teacher)
}

const sendMessage = (teacher: any) => {
  console.log('Send message to:', teacher)
}

const removeTeacher = async (teacher: any) => {
  if (!await alertConfirm(`確定要移除老師「${teacher.name}」嗎？此操作將刪除該老師的帳號。`)) {
    return
  }

  try {
    const api = useApi()
    await api.delete(`/teachers/${teacher.id}`)
    await fetchTeachers()
    showMenu.value[teacher.id] = false
  } catch (err) {
    console.error('Failed to remove teacher:', err)
    await alertError('移除失敗，請稍後再試')
  }
}
</script>
