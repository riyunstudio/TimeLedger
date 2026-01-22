<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-md animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10">
        <h3 class="text-lg font-semibold text-slate-100">
          {{ room ? '編輯教室' : '新增教室' }}
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="p-4 space-y-4">
        <div>
          <label class="block text-slate-300 mb-2 font-medium">教室名稱</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="例：Room A"
            class="input-field"
            required
          />
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium">容量</label>
          <input
            v-model.number="form.capacity"
            type="number"
            min="1"
            placeholder="1"
            class="input-field"
            required
          />
        </div>

        <div class="flex gap-3">
          <button
            @click="emit('close')"
            class="flex-1 glass-btn py-3 rounded-xl font-medium"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="flex-1 btn-primary"
          >
            {{ loading ? '儲存中...' : '儲存' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  room: any | null
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const loading = ref(false)
const { getCenterId } = useCenterId()
const form = ref({
  name: props.room?.name || '',
  capacity: props.room?.capacity || 1,
})

watch(() => props.room, (newRoom) => {
  if (newRoom) {
    form.value = {
      name: newRoom.name,
      capacity: newRoom.capacity,
    }
  } else {
    form.value = {
      name: '',
      capacity: 1,
    }
  }
})

const handleSubmit = async () => {
  loading.value = true

  try {
    const api = useApi()
    const centerId = getCenterId()

    const roomData = {
      name: form.value.name,
      capacity: parseInt(form.value.capacity.toString()),
    }

    if (props.room && props.room.id) {
      await api.put(`/admin/rooms/${props.room.id}`, roomData)
    } else {
      await api.post(`/admin/rooms`, roomData)
    }

    emit('saved')
    emit('close')
  } catch (error) {
    console.error('Failed to save room:', error)
    alert('儲存失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}
</script>
