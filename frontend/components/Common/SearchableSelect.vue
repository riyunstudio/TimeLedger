<template>
  <div ref="dropdownRef" class="relative">
    <label v-if="label" class="block mb-2 font-semibold text-sm text-slate-300">
      {{ label }}
      <span v-if="required" class="text-critical-500 ml-1">*</span>
    </label>

    <!-- 單選模式 -->
    <Combobox v-if="!multiple" v-model="selectedId" :disabled="disabled" @update:modelValue="handleSelect">
      <div class="relative" @click="isOpen = true">
        <div
          :class="[
            'flex items-center rounded-xl border backdrop-blur-glass bg-glass-dark border-white/10 text-slate-100 placeholder-slate-400 transition-all duration-300',
            disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer',
            error ? 'border-critical-500 ring-1 ring-critical-500' : ''
          ]"
        >
          <!-- 當有選中值但選項還沒載入時，顯示載入中的值 -->
          <ComboboxInput
            class="w-full px-4 py-3 bg-transparent border-none outline-none text-slate-100 placeholder-slate-400 focus:outline-none"
            :display-value="displaySelectedName"
            :placeholder="placeholder"
            @change="query = $event.target.value"
          />
          <div class="pr-3 flex items-center gap-2">
            <svg
              v-if="loading"
              class="animate-spin h-5 w-5 text-primary-500"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
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
            <ComboboxButton class="p-1 hover:text-primary-500 transition-colors">
              <svg
                class="w-5 h-5 text-slate-400 transition-transform duration-200"
                :class="{ 'rotate-180': isOpen }"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 9l-7 7-7-7"
                />
              </svg>
            </ComboboxButton>
          </div>
        </div>

        <!-- 使用 TransitionRoot 並添加 :show 屬性 -->
        <TransitionRoot
          :show="isOpen"
          enter="transition ease-out duration-100"
          enter-from="opacity-0"
          enter-to="opacity-100"
          leave="transition ease-in duration-100"
          leave-from="opacity-100"
          leave-to="opacity-0"
        >
          <ComboboxOptions
            class="absolute z-[1000] mt-2 w-full rounded-xl border backdrop-blur-glass bg-glass-dark border-white/10 shadow-xl overflow-hidden focus:outline-none max-h-60 scroll-py-1"
          >
            <div
              v-if="filteredOptions.length === 0 && query !== ''"
              class="px-4 py-3 text-slate-400 text-sm"
            >
              找不到符合的選項
            </div>

            <ComboboxOption
              v-for="option in filteredOptions"
              :key="option.id"
              :value="option.id"
              v-slot="{ selected, active }"
              as="template"
            >
              <li
                :class="[
                  'px-4 py-3 text-sm cursor-pointer transition-colors duration-150',
                  active ? 'bg-primary-500/20 text-primary-500' : 'text-slate-200',
                  selected ? 'bg-primary-500/10 text-primary-500 font-medium' : ''
                ]"
              >
                <div class="flex items-center justify-between">
                  <span class="truncate">{{ option.name }}</span>
                  <svg
                    v-if="selected"
                    class="w-4 h-4 text-primary-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                </div>
              </li>
            </ComboboxOption>

            <div
              v-if="filteredOptions.length === 0 && query === ''"
              class="px-4 py-3 text-slate-400 text-sm"
            >
              請輸入搜尋關鍵字
            </div>
          </ComboboxOptions>
        </TransitionRoot>
      </div>
    </Combobox>

    <!-- 多選模式 -->
    <Combobox v-if="multiple" :model-value="selectedIds" :disabled="disabled" @update:modelValue="handleMultiSelect">
      <div class="relative" @click="isOpen = true">
        <div
          :class="[
            'flex flex-wrap gap-2 p-2 rounded-xl border backdrop-blur-glass bg-glass-dark border-white/10 text-slate-100 placeholder-slate-400 transition-all duration-300 min-h-[52px]',
            disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer',
            error ? 'border-critical-500 ring-1 ring-critical-500' : ''
          ]"
        >
          <!-- 已選中的標籤 -->
          <div
            v-for="option in selectedOptions"
            :key="option.id"
            class="flex items-center gap-1 px-2 py-1 bg-primary-500/30 border border-primary-500/50 rounded-lg text-sm"
          >
            <span class="truncate">{{ option.name }}</span>
            <button
              v-if="!disabled"
              type="button"
              class="p-0.5 hover:bg-primary-500/50 rounded transition-colors"
              @click.stop="removeOption(option.id)"
            >
              <svg class="w-3 h-3 text-primary-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- 輸入框 -->
          <ComboboxInput
            class="flex-1 min-w-[100px] px-1 py-1 bg-transparent border-none outline-none text-slate-100 placeholder-slate-400 focus:outline-none text-sm"
            :placeholder="selectedOptions.length === 0 ? placeholder : ''"
            @change="query = $event.target.value"
          />
        </div>

        <!-- 多選下拉選單 -->
        <TransitionRoot
          :show="isOpen"
          enter="transition ease-out duration-100"
          enter-from="opacity-0"
          enter-to="opacity-100"
          leave="transition ease-in duration-100"
          leave-from="opacity-100"
          leave-to="opacity-0"
        >
          <ComboboxOptions
            v-if="filteredOptions.length > 0"
            class="absolute z-[1000] mt-2 w-full rounded-xl border backdrop-blur-glass bg-glass-dark border-white/10 shadow-xl overflow-hidden focus:outline-none max-h-60 scroll-py-1"
          >
            <ComboboxOption
              v-for="option in filteredOptions"
              :key="option.id"
              :value="option.id"
              v-slot="{ selected, active }"
              as="template"
            >
              <li
                :class="[
                  'px-4 py-3 text-sm cursor-pointer transition-colors duration-150',
                  active ? 'bg-primary-500/20 text-primary-500' : 'text-slate-200',
                  selected ? 'bg-primary-500/10 text-primary-500 font-medium' : ''
                ]"
              >
                <div class="flex items-center justify-between">
                  <span class="truncate">{{ option.name }}</span>
                  <svg
                    v-if="selected"
                    class="w-4 h-4 text-primary-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </div>
              </li>
            </ComboboxOption>
          </ComboboxOptions>
        </TransitionRoot>
      </div>
    </Combobox>

    <p v-if="multiple && error" class="text-critical-500 text-sm mt-1">{{ error }}</p>
    <p v-else-if="multiple && helper" class="text-slate-400 text-sm mt-1">{{ helper }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  Combobox,
  ComboboxInput,
  ComboboxButton,
  ComboboxOptions,
  ComboboxOption,
  TransitionRoot
} from '@headlessui/vue'

export interface SelectOption {
  id: number | string
  name: string
}

interface Props {
  modelValue: number | string | null | (number | string)[]
  options: SelectOption[]
  placeholder?: string
  label?: string
  loading?: boolean
  disabled?: boolean
  required?: boolean
  error?: string
  helper?: string
  multiple?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '請選擇',
  loading: false,
  disabled: false,
  required: false,
  multiple: false
})

const emit = defineEmits<{
  'update:modelValue': [value: number | string | null | (number | string)[]]
  change: [option: SelectOption | SelectOption[] | null]
}>()

const query = ref('')
const selectedId = computed({
  get: () => {
    // 多選模式下，返回選中的第一個 ID
    if (props.multiple) {
      const value = props.modelValue
      if (Array.isArray(value) && value.length > 0) {
        return value[0]
      }
      return null
    }
    return props.modelValue as number | string | null
  },
  set: (value) => emit('update:modelValue', value)
})

const isOpen = ref(false)

// 點擊外部關閉下拉選單
const dropdownRef = ref<HTMLElement | null>(null)

const handleClickOutside = (event: MouseEvent) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside)
})

const filteredOptions = computed(() => {
  if (query.value === '') {
    return props.options
  }

  const searchTerm = query.value.toLowerCase().trim()

  return props.options.filter((option) => {
    return option.name.toLowerCase().includes(searchTerm)
  })
})

const displaySelectedName = (id: number | string | null) => {
  if (id === null || id === undefined) return ''

  const selected = props.options.find((opt) => opt.id === id)
  if (selected) return selected.name

  // 當選項還沒載入但有選中值時，嘗試從現有選項中查找
  // 這裡不返回任何內容，讓 ComboboxInput 保持顯示當前值
  return ''
}

// 多選模式相關
const selectedIds = computed(() => {
  if (!props.multiple) return []
  const value = props.modelValue
  if (Array.isArray(value)) return value
  if (value === null || value === undefined) return []
  return [value]
})

const selectedOptions = computed(() => {
  return props.options.filter(opt => selectedIds.value.includes(opt.id))
})

const isOptionSelected = (id: number | string) => {
  return selectedIds.value.includes(id)
}

const removeOption = (id: number | string) => {
  const currentIds = [...selectedIds.value]
  const index = currentIds.indexOf(id)
  if (index > -1) {
    currentIds.splice(index, 1)
    emit('update:modelValue', currentIds.length > 0 ? currentIds : null)
    emit('change', selectedOptions.value)
  }
}

const handleMultiSelect = (id: number | string) => {
  const currentIds = [...selectedIds.value]
  if (!currentIds.includes(id)) {
    currentIds.push(id)
    emit('update:modelValue', currentIds)
    emit('change', [...selectedOptions.value, props.options.find(opt => opt.id === id)!])
  }
  query.value = ''
  isOpen.value = false
}

const handleSelect = (value: number | string) => {
  const selectedOption = props.options.find((opt) => opt.id === value) || null
  emit('change', selectedOption)
  isOpen.value = false
}
</script>
