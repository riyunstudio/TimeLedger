<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-md sm:max-w-lg max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          新增技能
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="p-4 space-y-4">
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">技能類別</label>
          <select v-model="form.category" class="input-field text-sm sm:text-base" required>
            <option value="">請選擇類別</option>
            <option value="MUSIC">音樂 (MUSIC)</option>
            <option value="ART">美術 (ART)</option>
            <option value="DANCE">舞蹈 (DANCE)</option>
            <option value="LANGUAGE">語言 (LANGUAGE)</option>
            <option value="SPORTS">運動 (SPORTS)</option>
            <option value="OTHER">其他 (OTHER)</option>
          </select>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">技能名稱</label>
          <input
            v-model="form.skill_name"
            type="text"
            placeholder="例：鋼琴、小提琴、吉他..."
            class="input-field text-sm sm:text-base"
            required
          />
        </div>

        <div class="flex gap-3 pt-2">
          <button
            type="button"
            @click="emit('close')"
            class="flex-1 glass-btn py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            {{ loading ? '新增中...' : '新增' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertError } from '~/composables/useAlert'

const emit = defineEmits<{
  close: []
  added: []
}>()

const teacherStore = useTeacherStore()
const loading = ref(false)
const form = ref({
  category: '' as 'MUSIC' | 'ART' | 'DANCE' | 'LANGUAGE' | 'SPORTS' | 'OTHER' | '',
  skill_name: '',
})

const handleSubmit = async () => {
  loading.value = true

  try {
    await teacherStore.createSkill(form.value)
    await teacherStore.fetchSkills()
    emit('added')
    emit('close')
  } catch (error) {
    console.error('Failed to add skill:', error)
    await alertError('新增失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}
</script>
