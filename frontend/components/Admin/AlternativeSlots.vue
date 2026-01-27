<template>
  <div class="alternative-slots">
    <div class="flex items-center justify-between mb-4">
      <h4 class="font-medium text-white">替代時段建議</h4>
      <span class="text-xs text-slate-400">
        基於老師 {{ teacherName }} 的課表
      </span>
    </div>

    <!-- 日期選擇標籤 -->
    <div class="flex gap-2 mb-4 overflow-x-auto pb-2">
      <button
        v-for="date in availableDates"
        :key="date.value"
        :class="[
          'px-3 py-1.5 rounded-lg text-sm whitespace-nowrap transition-colors',
          selectedDate === date.value
            ? 'bg-indigo-500 text-white'
            : 'bg-white/5 text-slate-400 hover:bg-white/10'
        ]"
        @click="selectedDate = date.value"
      >
        {{ date.label }}
      </button>
    </div>

    <!-- 時段卡片網格 -->
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-2">
      <button
        v-for="slot in filteredSlots"
        :key="`${slot.date}-${slot.start}`"
        :class="[
          'slot-button p-3 rounded-lg border-2 text-left transition-all',
          slot.available
            ? 'border-green-500/30 bg-green-500/5 hover:bg-green-500/10 hover:border-green-500'
            : 'border-white/10 bg-white/5 opacity-50 cursor-not-allowed'
        ]"
        :disabled="!slot.available"
        @click="selectSlot(slot)"
      >
        <div class="flex items-center justify-between mb-1">
          <span class="text-sm font-medium text-white">{{ slot.dateLabel }}</span>
          <span
            v-if="slot.available"
            class="text-xs text-green-400"
          >
            可安排
          </span>
          <span
            v-else
            class="text-xs text-red-400"
          >
            衝突
          </span>
        </div>
        <div class="text-xs text-slate-400">
          {{ slot.start }} - {{ slot.end }}
        </div>
        
        <!-- 可用教室標籤 -->
        <div v-if="slot.available && slot.availableRooms.length > 0" class="mt-2 flex flex-wrap gap-1">
          <BaseBadge
            v-for="room in slot.availableRooms.slice(0, 2)"
            :key="room.id"
            variant="success"
            size="xs"
          >
            {{ room.name }}
          </BaseBadge>
          <BaseBadge
            v-if="slot.availableRooms.length > 2"
            variant="secondary"
            size="xs"
          >
            +{{ slot.availableRooms.length - 2 }}
          </BaseBadge>
        </div>

        <!-- 衝突原因 -->
        <div v-if="!slot.available && slot.conflictReason" class="mt-2">
          <span class="text-xs text-red-400">{{ slot.conflictReason }}</span>
        </div>
      </button>

      <!-- 空狀態 -->
      <div
        v-if="filteredSlots.length === 0"
        class="col-span-full py-8 text-center text-slate-500"
      >
        <svg class="w-12 h-12 mx-auto mb-2 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p>暫無替代時段建議</p>
        <p class="text-xs text-slate-400 mt-1">嘗試選擇其他老師或調整時段</p>
      </div>
    </div>

    <!-- 自訂時段按鈕 -->
    <div class="mt-4 pt-4 border-t border-white/10">
      <button
        @click="$emit('custom')"
        class="w-full px-4 py-2 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white transition-colors text-sm flex items-center justify-center gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
        自訂其他時段
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Room {
  id: number
  name: string
}

interface AlternativeSlot {
  date: string
  dateLabel: string
  start: string
  end: string
  available: boolean
  availableRooms: Room[]
  conflictReason?: string
}

interface Props {
  teacherName: string
  slots: AlternativeSlot[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  select: [slot: AlternativeSlot]
  custom: []
}>()

// 取得可用日期（去重並排序）
const availableDates = computed(() => {
  const dateMap = new Map<string, string>()
  
  props.slots.forEach(slot => {
    if (!dateMap.has(slot.date)) {
      dateMap.set(slot.date, slot.dateLabel)
    }
  })
  
  return Array.from(dateMap.entries()).map(([value, label]) => ({
    value,
    label
  }))
})

// 篩選選定日期的時段
const filteredSlots = computed(() => {
  const firstDate = availableDates.value[0]?.value
  const selected = selectedDate.value || firstDate
  
  return props.slots.filter(slot => slot.date === selected)
})

// 選定的日期
const selectedDate = ref('')

// 選擇時段
const selectSlot = (slot: AlternativeSlot) => {
  if (!slot.available) return
  emit('select', slot)
}
</script>
