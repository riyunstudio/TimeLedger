<template>
  <div class="input-wrapper" :class="wrapperClasses">
    <label v-if="label" :class="labelClasses" class="block mb-2 font-semibold text-sm">
      {{ label }}
      <span v-if="required" class="text-critical-500 ml-1">*</span>
    </label>
    <input
      :value="modelValue"
      :type="type"
      :placeholder="placeholder"
      :disabled="disabled"
      :class="inputClasses"
      @input="handleInput"
      @blur="$emit('blur', $event)"
      @focus="$emit('focus', $event)"
    />
    <p v-if="error" class="text-critical-500 text-sm mt-1">{{ error }}</p>
    <p v-else-if="helper" class="text-slate-400 text-sm mt-1">{{ helper }}</p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue: string | number
  type?: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url'
  label?: string
  placeholder?: string
  disabled?: boolean
  required?: boolean
  error?: string
  helper?: string
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  placeholder: '',
  disabled: false,
  required: false,
  size: 'md',
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
  blur: [event: Event]
  focus: [event: Event]
}>()

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}

const wrapperClasses = computed(() => {
  return props.disabled ? 'opacity-50 cursor-not-allowed' : ''
})

const labelClasses = computed(() => {
  return {
    'text-slate-300': !document.documentElement.classList.contains('light'),
    'text-slate-700': document.documentElement.classList.contains('light'),
  }
})

const inputClasses = computed(() => {
  const sizes = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-3 text-base',
    lg: 'px-5 py-4 text-lg',
  }

  const baseClasses = 'w-full rounded-xl border backdrop-blur-glass bg-glass-dark border-white/10 text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-300'
  const lightClasses = 'light:bg-glass-light border-slate-300/30 text-slate-900 placeholder-slate-500'

  return `${baseClasses} ${sizes[props.size]} ${props.error ? 'border-critical-500 ring-1 ring-critical-500' : ''}`
})
</script>
