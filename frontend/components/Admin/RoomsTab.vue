<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-slate-100">教室列表</h2>
      <button
        @click="showCreateModal = true"
        class="btn-primary px-4 py-2 text-sm font-medium flex items-center gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新增教室
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
              <div class="h-4 w-16 bg-white/10 rounded"></div>
            </div>
            <div class="flex gap-2">
              <div class="h-6 w-16 bg-white/10 rounded-full"></div>
              <div class="h-6 w-16 bg-white/10 rounded-full"></div>
            </div>
          </div>
          <div class="mt-3 pt-3 border-t border-white/10">
            <div class="h-3 w-24 bg-white/10 rounded"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空狀態 -->
    <div v-else-if="rooms.length === 0" class="text-center py-12 text-slate-500 glass-card">
      <svg class="w-16 h-16 mx-auto mb-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
      </svg>
      <p class="mb-4">尚未添加教室</p>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors"
      >
        新增第一間教室
      </button>
    </div>

    <!-- 教室列表 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="room in rooms"
        :key="room.id"
        class="glass-card p-5 hover:bg-white/5 transition-all"
      >
        <div class="flex items-start justify-between mb-3">
          <div>
            <h3 class="text-lg font-semibold text-slate-100 mb-1">{{ room.name }}</h3>
            <p class="text-sm text-slate-400">Room {{ room.id }}</p>
          </div>
          <div class="flex gap-2">
            <button
              @click="toggleRoom(room)"
              class="p-2 rounded-lg transition-colors"
              :class="room.status === 'ACTIVE' ? 'hover:bg-success-500/20 text-success-500' : 'hover:bg-slate-500/20 text-slate-400'"
              :title="room.status === 'ACTIVE' ? '停用' : '啟用'"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="room.status === 'ACTIVE'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                <path v-if="room.status === 'ACTIVE'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </button>
            <button
              @click="editRoom(room)"
              class="p-2 rounded-lg hover:bg-white/10 transition-colors"
              title="編輯"
            >
              <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 0L21.828 3.172a2 2 0 010-2.828l-7-7a2 2 0 00-2.828 0L2.172 20.828a2 2 0 010 2.828l7 7a2 2 0 0012.828 0l7.172-7.172z" />
              </svg>
            </button>
            <button
              @click="deleteRoom(room.id)"
              class="p-2 rounded-lg hover:bg-critical-500/20 text-critical-500 transition-colors"
              title="刪除"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>

        <div class="space-y-2">
          <div class="flex items-center gap-2">
            <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20H2v-2a3 3 0 015.356-1.857m0 0a5 5 0 0111 0 5 5 0 0111 0z" />
            </svg>
            <span class="text-sm text-slate-400">容量</span>
            <span class="text-sm font-medium text-slate-100">{{ room.capacity }} 人</span>
          </div>

          <div v-if="room.equipment" class="flex flex-wrap gap-2">
            <span
              v-for="(eq, index) in room.equipment"
              :key="index"
              class="px-2 py-1 rounded-full text-xs font-medium bg-secondary-500/20 text-secondary-500"
            >
              {{ eq }}
            </span>
          </div>

          <div class="flex items-center gap-2">
            <svg class="w-5 h-5 text-success-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span
              class="px-2 py-1 rounded-full text-xs font-medium"
              :class="room.status === 'ACTIVE' ? 'bg-success-500/20 text-success-500' : room.status === 'BUSY' ? 'bg-warning-500/20 text-warning-500' : 'bg-slate-500/20 text-slate-400'"
            >
              {{ room.status === 'ACTIVE' ? '可使用' : room.status === 'BUSY' ? '使用中' : '已停用' }}
            </span>
          </div>
        </div>

        <div class="text-xs text-slate-500 mt-3 pt-3 border-t border-white/10">
          創建於：{{ formatDate(room.created_at) }}
        </div>
      </div>
    </div>

    <RoomModal
      v-if="showCreateModal"
      :room="null"
      @close="showCreateModal = false"
      @saved="fetchRooms"
    />

    <RoomModal
      v-if="showEditModal"
      :room="editingRoom"
      @close="showEditModal = false"
      @saved="fetchRooms"
    />
  </div>
</template>

<script setup lang="ts">
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingRoom = ref<any>(null)
const loading = ref(false)
const { getCenterId } = useCenterId()

// Alert composable
const { error: alertError, confirm: alertConfirm } = useAlert()

const rooms = ref<any[]>([])

const fetchRooms = async () => {
  loading.value = true
  try {
    const api = useApi()
    const centerId = getCenterId()
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    rooms.value = response.datas || []
  } catch (error) {
    console.error('Failed to fetch rooms:', error)
    rooms.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchRooms()
})

const editRoom = (room: any) => {
  editingRoom.value = { ...room }
  showEditModal.value = true
}

const deleteRoom = async (id: number) => {
  if (await alertConfirm('確定要刪除此教室？')) {
    try {
      const api = useApi()
      const centerId = getCenterId()
      await api.delete(`/admin/rooms/${id}`)
      rooms.value = rooms.value.filter(r => r.id !== id)
    } catch (err) {
      console.error('Failed to delete room:', err)
      await alertError('刪除失敗，請稍後再試')
    }
  }
}

const toggleRoom = async (room: any) => {
  const newStatus = room.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE'
  const actionText = room.status === 'ACTIVE' ? '停用' : '啟用'

  if (!await alertConfirm(`確定要${actionText}教室「${room.name}」嗎？`)) {
    return
  }

  try {
    const api = useApi()
    await api.patch(`/admin/rooms/${room.id}/toggle-active`, { is_active: room.status !== 'ACTIVE' })
    room.status = newStatus
  } catch (err) {
    console.error('Failed to toggle room:', err)
    await alertError(`${actionText}失敗，請稍後再試`)
  }
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}
</script>
