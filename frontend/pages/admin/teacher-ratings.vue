<template>
  <div class="p-4 md:p-6 max-w-7xl mx-auto">
    <div class="mb-6 md:mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
        老師評分
      </h1>
      <p class="text-slate-400 text-sm md:text-base">
        為中心的老師進行評分與備註，供智慧媒合參考
      </p>
    </div>

    <!-- 搜尋與篩選 -->
    <div class="mb-6 glass-card p-4">
      <div class="flex flex-col sm:flex-row gap-4">
        <div class="flex-1">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜尋老師名稱..."
            class="w-full px-4 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500"
          />
        </div>
        <select
          v-model="filterRating"
          class="px-4 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500"
        >
          <option value="">全部評分</option>
          <option value="5">5 星</option>
          <option value="4">4 星以上</option>
          <option value="3">3 星以上</option>
          <option value="2">2 星以上</option>
          <option value="1">1 星以上</option>
          <option value="0">未評分</option>
        </select>
      </div>
    </div>

    <!-- 老師列表 -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
    </div>

    <div v-else-if="filteredTeachers.length === 0" class="text-center py-16">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-white/5 flex items-center justify-center">
        <svg class="w-8 h-8 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </div>
      <p class="text-slate-500 mb-4">目前沒有已加入的老師</p>
      <p class="text-sm text-slate-600">請先邀請老師加入中心</p>
    </div>

    <div v-else class="grid gap-4">
      <div
        v-for="teacher in filteredTeachers"
        :key="teacher.id"
        class="glass-card p-4"
      >
        <div class="flex flex-col sm:flex-row sm:items-center gap-4">
          <div class="flex items-center gap-4 flex-1">
            <div class="w-12 h-12 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0">
              <span class="text-white font-medium">{{ teacher.name?.charAt(0) || '?' }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="text-white font-medium truncate">{{ teacher.name }}</h3>
              <p class="text-sm text-slate-400 truncate">{{ teacher.city }} {{ teacher.district }}</p>
            </div>
          </div>

          <div class="flex items-center gap-4">
            <!-- 評分顯示 -->
            <div class="flex items-center gap-2">
              <div class="flex">
                <template v-for="star in 5" :key="star">
                  <svg
                    class="w-5 h-5"
                    :class="star <= (teacher.note?.rating || 0) ? 'text-warning-500' : 'text-slate-600'"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                  </svg>
                </template>
              </div>
              <span class="text-sm text-slate-400">
                {{ teacher.note?.rating || 0 }} / 5
              </span>
            </div>

            <!-- 編輯按鈕 -->
            <button
              @click="openEditModal(teacher)"
              class="px-4 py-2 rounded-lg bg-primary-500/20 text-primary-500 hover:bg-primary-500/30 transition-colors text-sm"
            >
              編輯評分
            </button>
          </div>
        </div>

        <!-- 備註預覽 -->
        <div v-if="teacher.note?.internal_note" class="mt-3 pt-3 border-t border-white/10">
          <p class="text-sm text-slate-400">
            <span class="font-medium text-slate-500">備註：</span>
            {{ teacher.note.internal_note }}
          </p>
        </div>
      </div>
    </div>

    <!-- 統計資訊 -->
    <div v-if="teachers.length > 0" class="mt-6 grid grid-cols-2 sm:grid-cols-4 gap-4">
      <div class="glass-card p-4 text-center">
        <div class="text-2xl font-bold text-white">{{ teachers.length }}</div>
        <div class="text-sm text-slate-400">總老師數</div>
      </div>
      <div class="glass-card p-4 text-center">
        <div class="text-2xl font-bold text-warning-500">{{ ratedCount }}</div>
        <div class="text-sm text-slate-400">已評分</div>
      </div>
      <div class="glass-card p-4 text-center">
        <div class="text-2xl font-bold text-slate-400">{{ unratedCount }}</div>
        <div class="text-sm text-slate-400">未評分</div>
      </div>
      <div class="glass-card p-4 text-center">
        <div class="text-2xl font-bold text-primary-500">{{ averageRating }}</div>
        <div class="text-sm text-slate-400">平均評分</div>
      </div>
    </div>
  </div>

  <!-- 編輯評分 Modal -->
  <Teleport to="body">
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" @click.self="closeEditModal">
      <div class="glass-card w-full max-w-lg">
        <div class="flex items-center justify-between p-4 border-b border-white/10">
          <h3 class="text-lg font-semibold text-white">編輯老師評分</h3>
          <button @click="closeEditModal" class="p-2 rounded-lg hover:bg-white/10">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div v-if="editingTeacher" class="p-4 space-y-4">
          <!-- 老師資訊 -->
          <div class="flex items-center gap-4 p-4 rounded-xl bg-white/5">
            <div class="w-12 h-12 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
              <span class="text-white font-medium">{{ editingTeacher.name?.charAt(0) || '?' }}</span>
            </div>
            <div>
              <h4 class="text-white font-medium">{{ editingTeacher.name }}</h4>
              <p class="text-sm text-slate-400">{{ editingTeacher.city }} {{ editingTeacher.district }}</p>
            </div>
          </div>

          <!-- 評分 -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">評分 (0-5 星)</label>
            <div class="flex items-center gap-1">
              <button
                v-for="star in 5"
                :key="star"
                @click="editForm.rating = star"
                class="p-1 rounded-lg transition-colors"
                :class="editForm.rating >= star ? 'bg-primary-500/20' : 'hover:bg-white/5'"
              >
                <svg
                  class="w-8 h-8"
                  :class="editForm.rating >= star ? 'text-warning-500' : 'text-slate-600'"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              </button>
              <button
                @click="editForm.rating = 0"
                class="ml-2 px-3 py-1 rounded-lg text-sm transition-colors"
                :class="editForm.rating === 0 ? 'bg-critical-500/20 text-critical-500' : 'bg-white/5 text-slate-400 hover:bg-white/10'"
              >
                清除
              </button>
            </div>
            <p class="mt-2 text-sm" :class="editForm.rating > 0 ? 'text-warning-400' : 'text-slate-500'">
              {{ ratingLabels[editForm.rating] }}
            </p>
          </div>

          <!-- 內部備註 -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-1">內部備註</label>
            <textarea
              v-model="editForm.internal_note"
              rows="4"
              placeholder="記錄老師的表現特點、專長領域、合作經驗等..."
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-primary-500 resize-none"
            ></textarea>
            <p class="mt-1 text-xs text-slate-500">
              此備註僅供內部管理使用，影響智慧媒合的評分權重
            </p>
          </div>

          <!-- 關鍵字提示 -->
          <div class="p-3 rounded-xl bg-primary-500/10 border border-primary-500/30">
            <p class="text-sm text-primary-400 mb-2">可用關鍵字（可提升媒合分數）</p>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="keyword in keywords"
                :key="keyword"
                @click="addKeyword(keyword)"
                class="px-2 py-1 rounded-full text-xs bg-primary-500/20 text-primary-500 cursor-pointer hover:bg-primary-500/30 transition-colors"
              >
                {{ keyword }}
              </span>
            </div>
          </div>

          <div class="flex gap-3 pt-4">
            <button
              v-if="editingTeacher.note?.id"
              @click="deleteNote"
              :disabled="saving"
              class="px-4 py-2 rounded-lg bg-critical-500/20 text-critical-500 hover:bg-critical-500/30 transition-colors"
            >
              刪除評分
            </button>
            <button @click="closeEditModal" class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors">
              取消
            </button>
            <button @click="saveNote" :disabled="saving" class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50">
              {{ saving ? '儲存中...' : '儲存' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth-admin',
  layout: 'admin',
})

const notificationUI = useNotification()
const { confirm: alertConfirm } = useAlert()
const { getCenterId } = useCenterId()
const api = useApi()

interface TeacherNote {
  id?: number
  rating: number
  internal_note: string
}

interface Teacher {
  id: number
  name: string
  city: string
  district: string
  note?: TeacherNote
}

const loading = ref(false)
const saving = ref(false)
const showEditModal = ref(false)
const teachers = ref<Teacher[]>([])
const searchQuery = ref('')
const filterRating = ref('')
const editingTeacher = ref<Teacher | null>(null)

const editForm = reactive({
  rating: 0,
  internal_note: ''
})

const ratingLabels: Record<number, string> = {
  0: '未評分',
  1: '需改進',
  2: '一般',
  3: '良好',
  4: '優良',
  5: '優秀'
}

const keywords = ['推薦', '優秀', '穩定', '經驗豐富', '專業', '認真', '負責', '熱心']

const filteredTeachers = computed(() => {
  let result = teachers.value

  // 搜尋過濾
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(t => t.name.toLowerCase().includes(query))
  }

  // 評分過濾
  if (filterRating.value) {
    const minRating = parseInt(filterRating.value)
    if (minRating === 0) {
      result = result.filter(t => !t.note || t.note.rating === 0)
    } else {
      result = result.filter(t => t.note && t.note.rating >= minRating)
    }
  }

  return result
})

const ratedCount = computed(() => {
  return teachers.value.filter(t => t.note && t.note.rating > 0).length
})

const unratedCount = computed(() => {
  return teachers.value.length - ratedCount.value
})

const averageRating = computed(() => {
  const rated = teachers.value.filter(t => t.note && t.note.rating > 0)
  if (rated.length === 0) return '0.0'
  const sum = rated.reduce((acc, t) => acc + (t.note?.rating || 0), 0)
  return (sum / rated.length).toFixed(1)
})

const fetchTeachers = async () => {
  loading.value = true
  try {
    const response = await api.get<{ code: number; datas: Teacher[] }>(
      '/teachers'
    )

    teachers.value = response.datas || []

    // 並行為每位老師取得評分資料
    await Promise.all(
      teachers.value.map(async (teacher) => {
        try {
          const noteResponse = await api.get<{ code: number; datas: TeacherNote }>(
            `/admin/teachers/${teacher.id}/note`
          )
          if (noteResponse.datas) {
            teacher.note = noteResponse.datas
          }
        } catch {
          // 沒有評分資料是正常的
        }
      })
    )
  } catch (error) {
    console.error('Failed to fetch teachers:', error)
    notificationUI.error('載入老師列表失敗')
    teachers.value = []
  } finally {
    loading.value = false
  }
}

const openEditModal = async (teacher: Teacher) => {
  editingTeacher.value = teacher

  // 確保載入最新評分資料
  if (!teacher.note) {
    try {
      const noteResponse = await api.get<{ code: number; datas: TeacherNote }>(
        `/admin/teachers/${teacher.id}/note`
      )
      if (noteResponse.datas) {
        teacher.note = noteResponse.datas
      }
    } catch {
      // 沒有評分資料是正常的
    }
  }

  editForm.rating = teacher.note?.rating || 0
  editForm.internal_note = teacher.note?.internal_note || ''
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editingTeacher.value = null
}

const addKeyword = (keyword: string) => {
  const textarea = document.querySelector('textarea[placeholder*="記錄老師"]') as HTMLTextAreaElement
  if (textarea) {
    const current = editForm.internal_note
    editForm.internal_note = current ? `${current} ${keyword}` : keyword
    textarea.focus()
  }
}

const saveNote = async () => {
  if (!editingTeacher.value) return

  saving.value = true
  try {
    await api.put(`/admin/teachers/${editingTeacher.value.id}/note`, {
      rating: editForm.rating,
      internal_note: editForm.internal_note
    })

    notificationUI.success('評分已儲存')
    closeEditModal()
    await fetchTeachers()
  } catch (error) {
    console.error('Failed to save note:', error)
    notificationUI.error('儲存評分失敗')
  } finally {
    saving.value = false
  }
}

const deleteNote = async () => {
  if (!editingTeacher.value || !await alertConfirm('確定要刪除這位老師的評分嗎？')) return

  saving.value = true
  try {
    await api.delete(`/admin/teachers/${editingTeacher.value.id}/note`)
    notificationUI.success('評分已刪除')
    closeEditModal()
    await fetchTeachers()
  } catch (error) {
    console.error('Failed to delete note:', error)
    notificationUI.error('刪除評分失敗')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchTeachers()
})
</script>
