<template>
  <div class="room-card-select">
    <label class="block text-slate-300 mb-3">教室（可多選）</label>
    
    <!-- 卡片式教室選擇 -->
    <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 max-h-64 overflow-y-auto p-1">
      <div
        v-for="room in rooms"
        :key="room.id"
        :class="[
          'room-card p-3 rounded-lg border-2 cursor-pointer transition-all',
          selectedIds.includes(room.id)
            ? 'border-indigo-500 bg-indigo-500/10'
            : 'border-white/10 bg-white/5 hover:border-white/20'
        ]"
        @click="toggleRoom(room.id)"
      >
        <div class="flex items-center justify-between mb-2">
          <span class="font-medium text-white text-sm">{{ room.name }}</span>
          <svg
            v-if="selectedIds.includes(room.id)"
            class="w-5 h-5 text-indigo-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        
        <!-- 容量資訊 -->
        <div class="flex items-center gap-1 text-xs text-slate-400 mb-2">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          <span>{{ room.capacity || 20 }} 人</span>
        </div>
        
        <!-- 設備標籤 -->
        <div v-if="room.equipment && room.equipment.length > 0" class="flex flex-wrap gap-1">
          <span
            v-for="(equip, index) in room.equipment.slice(0, 3)"
            :key="index"
            :class="[
              'px-1.5 py-0.5 rounded text-xs',
              selectedIds.includes(room.id)
                ? 'bg-indigo-500/20 text-indigo-300'
                : 'bg-white/10 text-slate-400'
            ]"
          >
            {{ equip }}
          </span>
          <span
            v-if="room.equipment.length > 3"
            :class="[
              'px-1.5 py-0.5 rounded text-xs',
              selectedIds.includes(room.id)
                ? 'bg-indigo-500/20 text-indigo-300'
                : 'bg-white/10 text-slate-400'
            ]"
          >
            +{{ room.equipment.length - 3 }}
          </span>
        </div>
        
        <!-- 無設備時顯示提示 -->
        <div v-else class="text-xs text-slate-500">
          基本設施
        </div>
      </div>
    </div>
    
    <!-- 選取狀態提示 -->
    <div class="mt-3 flex items-center gap-2">
      <p v-if="selectedIds.length === 0" class="text-xs text-slate-500">
        <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        未選擇教室，將搜尋所有可用教室
      </p>
      <p v-else class="text-xs text-indigo-400">
        <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        已選擇 {{ selectedIds.length }} 間教室
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Room {
  id: number
  name: string
  capacity?: number
  equipment?: string[]
}

const props = defineProps<{
  rooms: Room[]
  modelValue: number[]
}>()

const emit = defineEmits<{
  'update:modelValue': [ids: number[]]
}>()

// 計算屬性：雙向綁定
const selectedIds = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 切換教室選取狀態
const toggleRoom = (roomId: number) => {
  const index = selectedIds.value.indexOf(roomId)
  if (index === -1) {
    selectedIds.value = [...selectedIds.value, roomId]
  } else {
    selectedIds.value = selectedIds.value.filter(id => id !== roomId)
  }
}
</script>
