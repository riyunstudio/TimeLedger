<template>
  <div class="flex items-center justify-between mb-6">
    <div>
      <h2 class="text-xl font-semibold text-white">例外申請</h2>
      <p class="text-sm text-slate-400 mt-1">管理您的請假和調課申請</p>
    </div>
    <button
      @click="showModal = true"
      class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors flex items-center gap-2"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      新增申請
    </button>
  </div>

  <!-- 申請統計摘要 -->
  <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
    <div class="glass-card p-4">
      <p class="text-sm text-slate-400">全部申請</p>
      <p class="text-2xl font-bold text-white mt-1">{{ teacherStore.exceptions.length }}</p>
    </div>
    <div class="glass-card p-4">
      <p class="text-sm text-slate-400">待審核</p>
      <p class="text-2xl font-bold text-warning-500 mt-1">{{ statusCounts.PENDING }}</p>
    </div>
    <div class="glass-card p-4">
      <p class="text-sm text-slate-400">已核准</p>
      <p class="text-2xl font-bold text-success-500 mt-1">{{ statusCounts.APPROVED }}</p>
    </div>
    <div class="glass-card p-4">
      <p class="text-sm text-slate-400">已拒絕</p>
      <p class="text-2xl font-bold text-critical-500 mt-1">{{ statusCounts.REJECTED }}</p>
    </div>
  </div>

  <div class="flex gap-2 mb-4 overflow-x-auto pb-2">
    <button
      v-for="status in statusFilters"
      :key="status.value"
      @click="currentFilter = status.value"
      class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all whitespace-nowrap"
      :class="currentFilter === status.value ? 'bg-primary-500 text-white' : 'bg-white/5 text-slate-400 hover:text-white'"
    >
      {{ status.label }}
      <span
        v-if="status.count > 0"
        class="ml-1 px-1.5 py-0.5 rounded-full text-xs"
        :class="currentFilter === status.value ? 'bg-white/20' : 'bg-white/10'"
      >
        {{ status.count }}
      </span>
    </button>
  </div>

  <div class="space-y-3">
    <div
      v-for="exception in filteredExceptions"
      :key="exception.id"
      class="glass-card p-4"
    >
      <div class="flex items-start justify-between">
        <div class="flex-1">
          <div class="flex items-center gap-2 mb-2">
            <span
              class="px-2 py-0.5 rounded-full text-xs font-medium"
              :class="getStatusClass(exception.status)"
            >
              {{ getStatusText(exception.status) }}
            </span>
            <span
              class="px-2 py-0.5 rounded-full text-xs font-medium"
              :class="exception.type === 'CANCEL' ? 'bg-critical-500/20 text-critical-500' : 'bg-secondary-500/20 text-secondary-500'"
            >
              {{ exception.type === 'CANCEL' ? '停課' : '改期' }}
            </span>
            <span class="text-xs text-slate-500">
              {{ formatDateTime(exception.created_at) }}
            </span>
          </div>
          <div class="text-white font-medium mb-1">
            {{ formatDate(exception.original_date) }}
          </div>
          <p class="text-sm text-slate-400 line-clamp-2">{{ exception.reason }}</p>

          <!-- 詳細資訊展開區塊 -->
          <div v-if="expandedException === exception.id" class="mt-4 p-3 bg-white/5 rounded-lg">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div>
                <p class="text-slate-400">申請時間</p>
                <p class="text-white">{{ formatDateTime(exception.created_at) }}</p>
              </div>
              <div v-if="exception.type === 'RESCHEDULE'">
                <p class="text-slate-400">新時間</p>
                <p class="text-white">{{ formatDate(exception.new_start_at || '') }} - {{ formatDate(exception.new_end_at || '') }}</p>
              </div>
              <div v-if="exception.new_teacher_name">
                <p class="text-slate-400">代課老師</p>
                <p class="text-white">{{ exception.new_teacher_name }}</p>
              </div>
              <div v-if="exception.reviewed_by">
                <p class="text-slate-400">審核人</p>
                <p class="text-white">{{ exception.reviewed_by }}</p>
              </div>
              <div v-if="exception.reviewed_at">
                <p class="text-slate-400">審核時間</p>
                <p class="text-white">{{ formatDateTime(exception.reviewed_at) }}</p>
              </div>
              <div v-if="exception.review_note">
                <p class="text-slate-400">審核回覆</p>
                <p class="text-white">{{ exception.review_note }}</p>
              </div>
            </div>
          </div>
        </div>
        <div class="flex gap-2">
          <button
            @click="toggleExpand(exception.id)"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
            title="查看詳情"
          >
            <svg
              class="w-4 h-4 text-slate-400 transition-transform"
              :class="expandedException === exception.id ? 'rotate-180' : ''"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          <button
            v-if="exception.status === 'PENDING'"
            @click="handleRevoke(exception.id)"
            class="px-3 py-1 rounded-lg bg-critical-500/20 text-critical-500 text-sm hover:bg-critical-500/30 transition-colors"
          >
            撤回
          </button>
        </div>
      </div>
    </div>

    <div v-if="teacherStore.loading" class="space-y-3">
    <!-- 統計摘要骨架屏 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div v-for="i in 4" :key="i" class="glass-card p-4 animate-pulse">
        <div class="h-3 w-16 bg-white/10 rounded mb-2"></div>
        <div class="h-7 w-12 bg-white/10 rounded"></div>
      </div>
    </div>

    <!-- 篩選按鈕骨架屏 -->
    <div class="flex gap-2 overflow-x-auto pb-2">
      <div v-for="i in 5" :key="i" class="h-9 w-20 bg-white/10 rounded-lg flex-shrink-0 animate-pulse"></div>
    </div>

    <!-- 申請列表骨架屏 -->
    <div v-for="i in 3" :key="i" class="glass-card p-4 animate-pulse">
      <div class="flex items-start justify-between">
        <div class="flex-1">
          <div class="flex items-center gap-2 mb-2">
            <div class="w-16 h-5 bg-white/10 rounded-full"></div>
            <div class="w-24 h-3 bg-white/10 rounded"></div>
          </div>
          <div class="h-4 w-32 bg-white/10 rounded mb-1"></div>
          <div class="h-3 w-48 bg-white/10 rounded"></div>
        </div>
        <div class="w-12 h-8 bg-white/10 rounded"></div>
      </div>
    </div>
  </div>

  <div v-else-if="filteredExceptions.length === 0" class="text-center py-12 text-slate-500 glass-card">
      <svg class="w-16 h-16 mx-auto mb-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
      </svg>
      <p class="text-lg mb-2">暫無例外申請紀錄</p>
      <p class="text-sm text-slate-500 mb-4">點擊上方按鈕新增申請</p>
      <button
        @click="showModal = true"
        class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors"
      >
        新增申請
      </button>
    </div>
  </div>

  <ExceptionModal
    v-if="showModal"
    :centers="centers"
    :schedule-rules="scheduleRules"
    :prefill="prefillData"
    @close="showModal = false; prefillData = null"
    @submit="fetchExceptions"
  />

  <NotificationDropdown v-if="notificationUI.show.value" @close="notificationUI.close()" />
  <TeacherSidebar v-if="sidebarStore.isOpen.value" @close="sidebarStore.close()" />
</template>

<script setup lang="ts">
import { formatDateTime } from '~/composables/useTaiwanTime'
import type { ScheduleException } from '~/types'

 definePageMeta({
   middleware: 'auth-teacher',
   layout: 'default',
 })

 const route = useRoute()
 const router = useRouter()
 const teacherStore = useTeacherStore()
 const sidebarStore = useSidebar()
 const notificationUI = useNotification()
 const { confirm: alertConfirm } = useAlert()
 const showModal = ref(false)
 const currentFilter = ref('')
 const expandedException = ref<number | null>(null)

 // 預填資料（從課表拖曳過來）
 const prefillData = ref<{
   rule_id?: number
   center_id?: number
   course_name?: string
   original_date?: string
   original_time?: string
 } | null>(null)

 const statusFilters = computed(() => {
   const counts = {
     '': teacherStore.exceptions.length,
     PENDING: teacherStore.exceptions.filter(e => e.status === 'PENDING').length,
     APPROVED: teacherStore.exceptions.filter(e => e.status === 'APPROVED' || e.status === 'APPROVE').length,
     REJECTED: teacherStore.exceptions.filter(e => e.status === 'REJECTED' || e.status === 'REJECT').length,
     REVOKED: teacherStore.exceptions.filter(e => e.status === 'REVOKED').length,
   }
   return [
     { value: '', label: '全部', count: counts[''] },
     { value: 'PENDING', label: '待審核', count: counts.PENDING },
     { value: 'APPROVED', label: '已核准', count: counts.APPROVED },
     { value: 'REJECTED', label: '已拒絕', count: counts.REJECTED },
     { value: 'REVOKED', label: '已撤回', count: counts.REVOKED },
   ]
 })

 const statusCounts = computed(() => ({
   PENDING: teacherStore.exceptions.filter(e => e.status === 'PENDING').length,
   APPROVED: teacherStore.exceptions.filter(e => e.status === 'APPROVED' || e.status === 'APPROVE').length,
   REJECTED: teacherStore.exceptions.filter(e => e.status === 'REJECTED' || e.status === 'REJECT').length,
   REVOKED: teacherStore.exceptions.filter(e => e.status === 'REVOKED').length,
 }))

 const filteredExceptions = computed(() => {
  if (!currentFilter.value) return teacherStore.exceptions
  return teacherStore.exceptions.filter(e => e.status === currentFilter.value || e.status === currentFilter.value + 'D')
 })

 // 中心列表
 const centers = computed(() => teacherStore.centers.map(c => ({
   center_id: c.center_id,
   center_name: c.center_name
 })))

 // 課程規則列表（用於預填）
 const scheduleRules = ref<any[]>([])

 const toggleExpand = (id: number) => {
   expandedException.value = expandedException.value === id ? null : id
 }

 const fetchExceptions = async () => {
   await teacherStore.fetchExceptions(currentFilter.value || undefined)
 }

// 监听筛选器变化，重新获取数据
watch(currentFilter, async () => {
  await fetchExceptions()
})

const getStatusClass = (status: string) => {
  switch (status.toUpperCase()) {
    case 'PENDING':
      return 'bg-warning-500/20 text-warning-500'
    case 'APPROVED':
    case 'APPROVE': // 向后兼容旧数据
      return 'bg-success-500/20 text-success-500'
    case 'REJECTED':
    case 'REJECT': // 向后兼容旧数据
      return 'bg-critical-500/20 text-critical-500'
    case 'REVOKED':
      return 'bg-slate-500/20 text-slate-400'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const getStatusText = (status: string) => {
  switch (status.toUpperCase()) {
    case 'PENDING':
      return '待審核'
    case 'APPROVED':
    case 'APPROVE': // 向后兼容旧数据
      return '已核准'
    case 'REJECTED':
    case 'REJECT': // 向后兼容旧数据
      return '已拒絕'
    case 'REVOKED':
      return '已撤回'
    default:
      return status
  }
}

const handleRevoke = async (id: number) => {
  if (await alertConfirm('確定要撤回此申請嗎？')) {
    await teacherStore.revokeException(id)
  }
}

onMounted(async () => {
  await Promise.all([
    teacherStore.fetchCenters(),
    fetchExceptions(),
  ])

  // 檢查是否有 query 參數（從課表拖曳過來）
  if (route.query.action === 'create') {
    prefillData.value = {
      rule_id: route.query.rule_id ? Number(route.query.rule_id) : undefined,
      center_id: route.query.center_id ? Number(route.query.center_id) : undefined,
      course_name: route.query.course_name as string | undefined,
      original_date: route.query.original_date as string | undefined,
      original_time: route.query.original_time as string | undefined,
    }
    showModal.value = true

    // 清除 query 參數，避免重新整理時再次彈出
    router.replace({ path: route.path, query: {} })
  }
})
</script>
