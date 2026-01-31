<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-slate-100">待排課程</h2>
      <button
        @click="createNewOffering"
        class="btn-primary px-4 py-2 text-sm font-medium flex items-center gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新增待排課程
      </button>
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
              <div class="h-5 w-40 bg-white/10 rounded mb-2"></div>
              <div class="h-4 w-24 bg-white/10 rounded"></div>
            </div>
            <div class="flex gap-2">
              <div class="w-8 h-8 bg-white/10 rounded-lg"></div>
              <div class="w-8 h-8 bg-white/10 rounded-lg"></div>
            </div>
          </div>
          <div class="space-y-2">
            <div class="h-4 w-32 bg-white/10 rounded"></div>
            <div class="h-4 w-32 bg-white/10 rounded"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空狀態 -->
    <div v-else-if="offerings.length === 0" class="text-center py-12 text-slate-500 glass-card">
      <svg class="w-16 h-16 mx-auto mb-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
      </svg>
      <p class="mb-4">尚未添加待排課程</p>
      <button
        @click="createNewOffering"
        class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors"
      >
        新增第一個待排課程
      </button>
    </div>

    <!-- 待排課程列表 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="offering in offerings"
        :key="offering.id"
        class="glass-card p-5 hover:bg-white/5 transition-all cursor-pointer"
        draggable="true"
        @dragstart="(e: DragEvent) => handleDragStart(offering, e)"
      >
        <div class="flex items-start justify-between mb-3">
          <div>
            <h3 class="text-lg font-semibold text-slate-100 mb-1">{{ offering.name }}</h3>
            <p class="text-sm text-slate-400">Course ID: {{ offering.course_id }}</p>
          </div>
          <div class="flex gap-1">
            <button
              @click.stop="toggleOffering(offering)"
              class="p-2 rounded-lg transition-colors"
              :class="offering.status === 'ACTIVE' ? 'hover:bg-success-500/20 text-success-500' : 'hover:bg-slate-500/20 text-slate-400'"
              :title="offering.status === 'ACTIVE' ? '停用' : '啟用'"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="offering.status === 'ACTIVE'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
              </svg>
            </button>
            <button
              @click.stop="copyOffering(offering)"
              class="p-2 rounded-lg hover:bg-primary-500/20 text-primary-500 transition-colors"
              title="複製"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
            <button
              @click.stop="editOffering(offering)"
              class="p-2 rounded-lg hover:bg-white/10 transition-colors"
              title="編輯"
            >
              <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 0L21.828 3.172a2 2 0 010-2.828l-7-7a2 2 0 00-2.828 0L2.172 20.828a2 2 0 010 2.828l7 7a2 2 0 0012.828 0l7.172-7.172z" />
              </svg>
            </button>
            <button
              @click.stop="deleteOffering(offering)"
              class="p-2 rounded-lg hover:bg-red-500/20 transition-colors"
              title="刪除"
            >
              <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>

        <p v-if="offering.default_teacher_id" class="text-sm text-slate-300 mb-1">
          <span class="text-slate-400">預設老師：</span>
          <span class="text-slate-100">{{ getTeacherName(offering.default_teacher_id) }}</span>
        </p>
        <p v-else class="text-sm text-slate-500 mb-1">
          <span class="text-slate-400">預設老師：</span>
          <span>未指定</span>
        </p>

        <p v-if="offering.default_room_id" class="text-sm text-slate-300">
          <span class="text-slate-400">預設教室：</span>
          <span class="text-slate-100">{{ getRoomName(offering.default_room_id) }}</span>
        </p>
        <p v-else class="text-sm text-slate-500">
          <span class="text-slate-400">預設教室：</span>
          <span>未指定</span>
        </p>

        <p v-if="!offering.allow_buffer_override" class="text-xs text-warning-500 mt-1">
          ⚠️ 不允許覆蓋緩衝時間
        </p>
      </div>
    </div>

    <OfferingModal
      v-if="showModal"
      :offering="editingOffering"
      @close="showModal = false; editingOffering = null"
      @saved="fetchOfferings"
    />
  </div>
</template>

<script setup lang="ts">
import { alertConfirm, alertError } from '~/composables/useAlert'

const showModal = ref(false)
const editingOffering = ref<any>(null)
const offerings = ref<any[]>([])
const loading = ref(false)

// 資源快取
const teachersCache = ref<Map<number, any>>(new Map())
const roomsCache = ref<Map<number, any>>(new Map())

const { getCenterId } = useCenterId()

const handleDragStart = (offering: any, event: DragEvent) => {
  event.dataTransfer?.setData('application/json', JSON.stringify({
    type: 'offering',
    item: offering,
  }))
  event.dataTransfer.effectAllowed = 'copy'
}

// 取得老師名稱
const getTeacherName = (teacherId: number): string => {
  const teacher = teachersCache.value.get(teacherId)
  return teacher?.name || `老師 ${teacherId}`
}

// 取得教室名稱
const getRoomName = (roomId: number): string => {
  const room = roomsCache.value.get(roomId)
  return room?.name || `教室 ${roomId}`
}

const fetchResources = async () => {
  try {
    const api = useApi()
    const [teachersRes, roomsRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>('/teachers'),
      api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    ])

    teachersRes.datas?.forEach((t: any) => {
      teachersCache.value.set(t.id, t)
    })
    roomsRes.datas?.forEach((r: any) => {
      roomsCache.value.set(r.id, r)
    })
  } catch (error) {
    console.error('Failed to fetch resources:', error)
  }
}

const fetchOfferings = async () => {
  loading.value = true
  try {
    const api = useApi()
    const centerId = getCenterId()
    const response = await api.get<{ code: number; datas: any }>(`/admin/offerings`)
    if (response.datas?.offerings) {
      offerings.value = response.datas.offerings
    } else {
      offerings.value = []
    }
  } catch (error) {
    console.error('Failed to fetch offerings:', error)
    offerings.value = []
  } finally {
    loading.value = false
  }
}

const editOffering = (offering: any) => {
  editingOffering.value = offering
  showModal.value = true
}

const createNewOffering = () => {
  editingOffering.value = null
  showModal.value = true
}

const deleteOffering = async (offering: any) => {
  if (!await alertConfirm(`確定要刪除待排課程「${offering.name}」嗎？`)) {
    return
  }

  try {
    const api = useApi()
    const centerId = getCenterId()
    await api.delete(`/admin/offerings/${offering.id}`)
    await fetchOfferings()
  } catch (error) {
    console.error('Failed to delete offering:', error)
    await alertError('刪除失敗，請稍後再試')
  }
}

const toggleOffering = async (offering: any) => {
  const newStatus = offering.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE'
  const actionText = offering.status === 'ACTIVE' ? '停用' : '啟用'

  if (!await alertConfirm(`確定要${actionText}待排課程「${offering.name}」嗎？`)) {
    return
  }

  try {
    const api = useApi()
    await api.patch(`/admin/offerings/${offering.id}/toggle-active`, { is_active: offering.status !== 'ACTIVE' })
    offering.status = newStatus
  } catch (error) {
    console.error('Failed to toggle offering:', error)
    await alertError(`${actionText}失敗，請稍後再試`)
  }
}

const copyOffering = async (offering: any) => {
  const newName = prompt(`請輸入複製後的名稱：`, `${offering.name} (複製)`)
  if (!newName || newName.trim() === '') {
    return
  }

  try {
    const api = useApi()
    const centerId = getCenterId()
    await api.post(`/admin/centers/${centerId}/offerings/${offering.id}/copy`, {
      new_name: newName.trim(),
      copy_teacher: false
    })
    await fetchOfferings()
    await alertSuccess('複製成功')
  } catch (error) {
    console.error('Failed to copy offering:', error)
    await alertError('複製失敗，請稍後再試')
  }
}

onMounted(() => {
  fetchResources()
  fetchOfferings()
})
</script>
