<template>
  <div class="relative">
    <label v-if="label" class="block mb-2 font-semibold text-sm text-slate-300">
      {{ label }}
      <span v-if="required" class="text-critical-500 ml-1">*</span>
    </label>

    <Combobox v-model="selectedId" :disabled="disabled" @update:modelValue="handleSelect">
      <div class="relative">
        <div
          :class="[
            'flex items-center rounded-xl border backdrop-blur-glass bg-glass-dark border-white/10 text-slate-100 placeholder-slate-400 transition-all duration-300',
            disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer',
            isOpen ? 'ring-2 ring-primary-500 border-transparent' : '',
            error ? 'border-critical-500 ring-1 ring-critical-500' : ''
          ]"
        >
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

        <TransitionRoot
          leave="transition ease-in duration-100"
          leave-from="opacity-100"
          leave-to="opacity-0"
          @after-leave="query = ''"
        >
          <ComboboxOptions
            class="absolute z-50 mt-2 w-full rounded-xl border backdrop-blur-glass bg-glass-dark border-white/10 shadow-xl overflow-hidden focus:outline-none max-h-60 scroll-py-1"
            static
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

    <p v-if="error" class="text-critical-500 text-sm mt-1">{{ error }}</p>
    <p v-else-if="helper" class="text-slate-400 text-sm mt-1">{{ helper }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
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
  modelValue: number | string | null
  options: SelectOption[]
  placeholder?: string
  label?: string
  loading?: boolean
  disabled?: boolean
  required?: boolean
  error?: string
  helper?: string
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '請選擇',
  loading: false,
  disabled: false,
  required: false
})

const emit = defineEmits<{
  'update:modelValue': [value: number | string | null]
  change: [option: SelectOption | null]
}>()

const query = ref('')
const selectedId = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const isOpen = ref(false)

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
  return selected?.name || ''
}

const handleSelect = (value: number | string) => {
  const selectedOption = props.options.find((opt) => opt.id === value) || null
  emit('change', selectedOption)
}
</script>
