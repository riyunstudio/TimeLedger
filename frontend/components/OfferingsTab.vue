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

    <div v-if="offerings.length === 0" class="text-center py-12 text-slate-500 glass-card">
      尚未添加待排課程
    </div>

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

        <p v-if="(offering as any).default_teacher_id" class="text-sm text-slate-300 mb-1">
          <span class="text-slate-400">預設老師：</span>
          <span class="text-slate-100">Teacher {{ (offering as any).default_teacher_id }}</span>
        </p>
        <p v-if="offering.default_room_id" class="text-sm text-slate-300">
          <span class="text-slate-400">預設教室：</span>
          <span class="text-slate-100">Room {{ offering.default_room_id }}</span>
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
const showModal = ref(false)
const editingOffering = ref<any>(null)
const offerings = ref<any[]>([])
const { getCenterId } = useCenterId()

const handleDragStart = (offering: any, event: DragEvent) => {
  event.dataTransfer?.setData('application/json', JSON.stringify({
    type: 'offering',
    item: offering,
  }))
  event.dataTransfer.effectAllowed = 'copy'
}

const fetchOfferings = async () => {
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
  if (!confirm(`確定要刪除待排課程「${offering.name}」嗎？`)) {
    return
  }

  try {
    const api = useApi()
    const centerId = getCenterId()
    await api.delete(`/admin/offerings/${offering.id}`)
    await fetchOfferings()
  } catch (error) {
    console.error('Failed to delete offering:', error)
    alert('刪除失敗，請稍後再試')
  }
}

onMounted(() => {
  fetchOfferings()
})
</script>
