<template>
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
      @click.self="handleClose"
    >
      <div class="glass-card w-full max-w-lg animate-spring" @click.stop>
        <!-- Header -->
        <div class="flex items-center justify-between p-4 border-b border-white/10">
          <div>
            <h3 class="text-lg font-semibold text-slate-100">合併老師帳號</h3>
            <p class="text-sm text-slate-400 mt-0.5">
              將佔位老師的資料合併至已綁定的真實老師
            </p>
          </div>
          <button
            @click="handleClose"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Content -->
        <div class="p-4 space-y-4">
          <!-- Source Teacher (Placeholder) -->
          <div v-if="sourceTeacher" class="p-4 rounded-xl bg-warning-500/10 border border-warning-500/30">
            <div class="flex items-center justify-between mb-2">
              <p class="text-sm text-warning-400 font-medium">來源（將被合併）</p>
              <span class="px-2 py-0.5 text-xs rounded bg-warning-500/20 text-warning-400">
                佔位老師
              </span>
            </div>
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-warning-500/20 flex items-center justify-center">
                <span class="text-warning-400 font-medium text-lg">
                  {{ sourceTeacher.name.charAt(0) }}
                </span>
              </div>
              <div>
                <p class="text-white font-medium">{{ sourceTeacher.name }}</p>
                <p class="text-sm text-slate-400">
                  {{ sourceTeacher.email || '無 Email' }}
                </p>
              </div>
            </div>
          </div>

          <!-- Search Target Teacher -->
          <div>
            <label class="block text-sm text-slate-300 mb-2 font-medium">
              選擇目標老師 <span class="text-critical-500">*</span>
            </label>
            <div class="relative">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="搜尋老師名稱或 Email..."
                class="input-field pr-10"
                @input="handleSearch"
              />
              <svg
                class="absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>

            <!-- Teacher List -->
            <div
              v-if="filteredTeachers.length > 0"
              class="mt-2 max-h-48 overflow-y-auto space-y-1 border border-white/10 rounded-lg bg-slate-900/50"
            >
              <button
                v-for="teacher in filteredTeachers"
                :key="teacher.id"
                :disabled="teacher.id === sourceTeacher?.id"
                :class="[
                  'w-full px-3 py-2 text-left flex items-center gap-3 transition-colors',
                  teacher.id === sourceTeacher?.id
                    ? 'opacity-50 cursor-not-allowed'
                    : 'hover:bg-white/5',
                  selectedTargetTeacher?.id === teacher.id
                    ? 'bg-primary-500/20 border border-primary-500/30'
                    : ''
                ]"
                @click="selectTargetTeacher(teacher)"
              >
                <div class="w-8 h-8 rounded-full bg-primary-500/20 flex items-center justify-center flex-shrink-0">
                  <span class="text-primary-400 text-sm font-medium">
                    {{ teacher.name.charAt(0) }}
                  </span>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm text-white truncate">{{ teacher.name }}</p>
                  <p class="text-xs text-slate-400 truncate">
                    {{ teacher.email || '無 Email' }}
                  </p>
                </div>
                <svg
                  v-if="selectedTargetTeacher?.id === teacher.id"
                  class="w-5 h-5 text-primary-400 flex-shrink-0"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fill-rule="evenodd"
                    d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                    clip-rule="evenodd"
                  />
                </svg>
              </button>
            </div>
            <p v-else-if="searchQuery && !loading" class="mt-2 text-sm text-slate-500">
              沒有找到符合的老師
            </p>
          </div>

          <!-- Selected Target Teacher -->
          <div
            v-if="selectedTargetTeacher"
            class="p-4 rounded-xl bg-primary-500/10 border border-primary-500/30"
          >
            <p class="text-sm text-primary-400 font-medium mb-2">目標（接收資料）</p>
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full bg-primary-500/20 flex items-center justify-center">
                <span class="text-primary-400 font-medium text-lg">
                  {{ selectedTargetTeacher.name.charAt(0) }}
                </span>
              </div>
              <div>
                <p class="text-white font-medium">{{ selectedTargetTeacher.name }}</p>
                <p class="text-sm text-slate-400">
                  {{ selectedTargetTeacher.email || '無 Email' }}
                </p>
              </div>
            </div>
          </div>

          <!-- Migration Warning -->
          <div
            v-if="selectedTargetTeacher"
            class="p-4 rounded-xl bg-critical-500/10 border border-critical-500/30"
          >
            <div class="flex items-center gap-2 mb-3">
              <svg class="w-5 h-5 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                />
              </svg>
              <p class="text-sm text-critical-400 font-medium">資料遷移警告</p>
            </div>

            <p class="text-xs text-slate-400 mb-3">
              合併後，以下來源的資料將遷移至目標老師，來源老師將被軟刪除：
            </p>

            <!-- Migration Stats -->
            <div class="grid grid-cols-2 gap-2">
              <div class="p-2 rounded-lg bg-slate-800/50">
                <p class="text-xs text-slate-500">課表規則</p>
                <p class="text-sm text-white font-medium">{{ migrationStats.scheduleRules }}</p>
              </div>
              <div class="p-2 rounded-lg bg-slate-800/50">
                <p class="text-xs text-slate-500">例外記錄</p>
                <p class="text-sm text-white font-medium">{{ migrationStats.exceptions }}</p>
              </div>
              <div class="p-2 rounded-lg bg-slate-800/50">
                <p class="text-xs text-slate-500">課程筆記</p>
                <p class="text-sm text-white font-medium">{{ migrationStats.sessionNotes }}</p>
              </div>
              <div class="p-2 rounded-lg bg-slate-800/50">
                <p class="text-xs text-slate-500">私人行程</p>
                <p class="text-sm text-white font-medium">{{ migrationStats.personalEvents }}</p>
              </div>
              <div class="p-2 rounded-lg bg-slate-800/50">
                <p class="text-xs text-slate-500">課程項目</p>
                <p class="text-sm text-white font-medium">{{ migrationStats.offerings }}</p>
              </div>
              <div class="p-2 rounded-lg bg-slate-800/50">
                <p class="text-xs text-slate-500">技能</p>
                <p class="text-sm text-white font-medium">{{ migrationStats.skills }}</p>
              </div>
              <div class="p-2 rounded-lg bg-slate-800/50">
                <p class="text-xs text-slate-500">證照</p>
                <p class="text-sm text-white font-medium">{{ migrationStats.certificates }}</p>
              </div>
              <div class="p-2 rounded-lg bg-slate-800/50">
                <p class="text-xs text-slate-500">個人標籤</p>
                <p class="text-sm text-white font-medium">{{ migrationStats.hashtags }}</p>
              </div>
            </div>

            <div class="mt-3 p-2 rounded-lg bg-critical-500/10 border border-critical-500/20">
              <p class="text-xs text-critical-400">
                ⚠️ 此操作無法復原，來源老師將被刪除
              </p>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="p-4 border-t border-white/10 flex items-center gap-3">
          <button
            @click="handleClose"
            class="flex-1 py-3 rounded-xl font-medium bg-white/5 text-white hover:bg-white/10 transition-colors"
          >
            取消
          </button>
          <button
            @click="handleMerge"
            :disabled="!canMerge || loading"
            class="flex-1 py-3 rounded-xl font-medium bg-critical-500/30 text-critical-400 border border-critical-500/50 hover:bg-critical-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading" class="flex items-center justify-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                />
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                />
              </svg>
              處理中...
            </span>
            <span v-else>確認合併</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { alertError, alertConfirm } from '~/composables/useAlert'
import type { Teacher } from '~/types/teacher'

interface Props {
  isOpen: boolean
  sourceTeacher: Teacher | null
}

interface Emits {
  close: []
  merged: []
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const { getCenterId } = useCenterId()

// State
const loading = ref(false)
const searchQuery = ref('')
const allTeachers = ref<Teacher[]>([])
const selectedTargetTeacher = ref<Teacher | null>(null)
const migrationStats = ref({
  scheduleRules: 0,
  exceptions: 0,
  sessionNotes: 0,
  personalEvents: 0,
  offerings: 0,
  timetableCells: 0,
  skills: 0,
  certificates: 0,
  hashtags: 0
})

// Filtered teachers (exclude source teacher and placeholders)
const filteredTeachers = computed(() => {
  let result = allTeachers.value.filter(t => !t.is_placeholder)

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(
      t =>
        t.name.toLowerCase().includes(query) ||
        (t.email && t.email.toLowerCase().includes(query))
    )
  }

  // Exclude source teacher
  if (props.sourceTeacher) {
    result = result.filter(t => t.id !== props.sourceTeacher?.id)
  }

  return result
})

// Can merge validation
const canMerge = computed(() => {
  return (
    props.sourceTeacher !== null &&
    selectedTargetTeacher.value !== null &&
    selectedTargetTeacher.value.id !== props.sourceTeacher.id
  )
})

// Fetch bound teachers
const fetchBoundTeachers = async () => {
  if (!props.isOpen) return

  try {
    const api = useApi()
    const response = await api.get<Teacher[]>('/teachers')
    allTeachers.value = response || []

    // Filter only bound teachers (non-placeholder)
    allTeachers.value = allTeachers.value.filter(t => !t.is_placeholder)
  } catch (error) {
    console.error('取得老師列表失敗:', error)
    alertError('載入老師列表失敗')
  }
}

// Search handler (debounced)
let searchTimeout: ReturnType<typeof setTimeout> | null = null
const handleSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  // Debounce is minimal since we're filtering locally
}

// Select target teacher
const selectTargetTeacher = async (teacher: Teacher) => {
  selectedTargetTeacher.value = teacher
  await fetchMigrationStats(teacher.id)
}

// Fetch migration statistics
const fetchMigrationStats = async (targetId: number) => {
  if (!props.sourceTeacher) return

  try {
    const api = useApi()
    const centerId = getCenterId()

    // 嘗試從 API 獲取統計數據，如果沒有則使用預估數值
    // 這裡假設後端可能需要新增一個 endpoint 來獲取預覽數據
    // 目前使用預估數值
    migrationStats.value = {
      scheduleRules: 0, // 需要後端 API 支援
      exceptions: 0,
      sessionNotes: 0,
      personalEvents: 0,
      offerings: 0,
      timetableCells: 0,
      skills: props.sourceTeacher.skills?.length || 0,
      certificates: props.sourceTeacher.certificates?.length || 0,
      hashtags: props.sourceTeacher.personal_hashtags?.length || 0
    }
  } catch (error) {
    console.error('取得遷移統計失敗:', error)
  }
}

// Handle merge
const handleMerge = async () => {
  if (!canMerge.value || !props.sourceTeacher || !selectedTargetTeacher.value) return

  const confirmMessage = `確定要將「${props.sourceTeacher.name}」合併至「${selectedTargetTeacher.value.name}」嗎？

此操作將遷移所有資料並軟刪除來源老師，無法復原。`

  if (!(await alertConfirm(confirmMessage))) {
    return
  }

  loading.value = true

  try {
    const api = useApi()
    const centerId = getCenterId()

    await api.post('/admin/teachers/merge', {
      source_teacher_id: props.sourceTeacher.id,
      target_teacher_id: selectedTargetTeacher.value.id,
      center_id: centerId
    })

    emit('merged')
    emit('close')
    await alertSuccess('帳號合併成功')
  } catch (error) {
    console.error('合併失敗:', error)
    await alertError('合併失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}

// Handle close
const handleClose = () => {
  searchQuery.value = ''
  selectedTargetTeacher.value = null
  migrationStats.value = {
    scheduleRules: 0,
    exceptions: 0,
    sessionNotes: 0,
    personalEvents: 0,
    offerings: 0,
    timetableCells: 0,
    skills: 0,
    certificates: 0,
    hashtags: 0
  }
  emit('close')
}

// Watch for modal open to fetch teachers
watch(
  () => props.isOpen,
  async isOpen => {
    if (isOpen) {
      await fetchBoundTeachers()
    }
  }
)
</script>
