<template>
  <Teleport to="body">
    <div
      v-if="show"
      class="fixed inset-0 z-[2000] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
      @click.self="$emit('close')"
    >
      <div class="glass-card w-full max-w-md animate-spring">
        <div class="flex items-center justify-between p-4 border-b border-white/10">
          <h3 class="text-lg font-semibold text-white">{{ isSuspendMode ? '選擇停課範圍' : '選擇更新範圍' }}</h3>
          <button @click="$emit('close')" class="p-2 rounded-lg hover:bg-white/10">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-4 space-y-3">
          <p class="text-slate-400 text-sm mb-4">{{ isSuspendMode ? '請選擇要停課的範圍：' : '請選擇要更新哪些場次：' }}</p>

          <button
            v-for="option in visibleOptions"
            :key="option.value"
            @click="selectOption(option.value)"
            class="w-full p-4 rounded-xl border text-left transition-all"
            :class="selectedMode === option.value
              ? 'bg-primary-500/20 border-primary-500 text-white'
              : 'bg-slate-800/50 border-white/10 text-slate-300 hover:bg-slate-700/50'"
          >
            <div class="font-medium mb-1">{{ option.label }}</div>
            <div class="text-xs text-slate-400">{{ option.description }}</div>
          </button>
        </div>

        <div class="flex gap-3 p-4 border-t border-white/10">
          <button
            @click="$emit('close')"
            class="flex-1 py-2.5 rounded-xl font-medium text-slate-300 hover:bg-white/10 transition-colors"
          >
            取消
          </button>
          <button
            @click="confirm"
            :disabled="!selectedMode"
            class="flex-1 py-2.5 rounded-xl font-medium bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isSuspendMode ? '確認停課' : '確認修改' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
interface Props {
  show: boolean
  ruleName?: string
  ruleDate?: string
  isSuspendMode?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  show: false,
  ruleName: '',
  ruleDate: '',
  isSuspendMode: false,
})

const emit = defineEmits<{
  close: []
  confirm: [mode: string]
}>()

const selectedMode = ref<string>('')

const isSuspendMode = computed(() => props.isSuspendMode || false)

const updateOptions = [
  {
    value: 'SINGLE',
    label: '只修改這一天',
    description: `僅修改 ${props.ruleDate || '當前'} 這個日期的課程`,
  },
  {
    value: 'FUTURE',
    label: '修改這天及之後',
    description: `修改 ${props.ruleDate || '當前'} 起至課程結束的所有場次`,
  },
  {
    value: 'ALL',
    label: '修改全部',
    description: '修改這個循環規則的所有場次（包括過去）',
  },
]

const suspendOptions = [
  {
    value: 'SINGLE',
    label: '僅停課這一天',
    description: `僅停課 ${props.ruleDate || '當前'} 這個日期的課程`,
  },
  {
    value: 'FUTURE',
    label: '從此以後停課',
    description: `從 ${props.ruleDate || '當前'} 起至課程結束的所有場次都停課`,
  },
]

const visibleOptions = computed(() => {
  return isSuspendMode.value ? suspendOptions : updateOptions
})

const selectOption = (mode: string) => {
  selectedMode.value = mode
}

const confirm = () => {
  if (selectedMode.value) {
    emit('confirm', selectedMode.value)
  }
}

// Reset selection when modal opens
watch(() => props.show, (show) => {
  if (show) {
    selectedMode.value = ''
  }
})
</script>
