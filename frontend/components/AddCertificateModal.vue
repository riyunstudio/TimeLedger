<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-md sm:max-w-lg max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          新增證照
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="p-4 space-y-4">
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">證照名稱</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="例：ABRSM Grade 8"
            class="input-field text-sm sm:text-base"
            required
          />
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">發證日期</label>
          <input
            v-model="form.issued_at"
            type="datetime-local"
            class="input-field text-sm sm:text-base"
          />
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">上傳證照檔案</label>
          <div
            class="border-2 border-dashed border-slate-700 rounded-xl p-6 sm:p-8 text-center hover:border-primary-500 transition-colors cursor-pointer"
            @click="triggerFileInput"
          >
            <input
              ref="fileInput"
              type="file"
              accept=".pdf,.jpg,.jpeg,.png"
              class="hidden"
              @change="handleFileChange"
            />
            <svg v-if="!fileName" class="w-10 h-10 sm:w-12 sm:h-12 mx-auto mb-3 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            <p v-if="!fileName" class="text-slate-400 text-sm sm:text-base">點擊或拖曳檔案至此</p>
            <p v-else class="text-slate-100 text-sm sm:text-base">{{ fileName }}</p>
          </div>
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
            :disabled="loading || uploading"
            class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            {{ uploading ? '上傳中...' : loading ? '儲存中...' : '新增' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  close: []
  added: []
}>()

const teacherStore = useTeacherStore()
const loading = ref(false)
const uploading = ref(false)
const fileInput = ref<HTMLInputElement>()
const fileName = ref('')

const form = ref({
  name: '',
  issued_at: '',
})

const formatDateTimeForApi = (datetimeLocal: string): string => {
  if (!datetimeLocal) return ''
  return new Date(datetimeLocal).toISOString()
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    fileName.value = target.files[0].name
  }
}

const handleSubmit = async () => {
  loading.value = true

  try {
    const fileUrl = fileInput.value?.files && fileInput.value.files.length > 0
      ? `uploads/${fileName.value}`
      : undefined

    await teacherStore.createCertificate({
      name: form.value.name,
      file_url: fileUrl,
      issued_at: formatDateTimeForApi(form.value.issued_at),
    })

    await teacherStore.fetchCertificates()
    emit('added')
    emit('close')
  } catch (error) {
    console.error('Failed to add certificate:', error)
    alert('新增失敗，請稍後再試')
  } finally {
    loading.value = false
    uploading.value = false
  }
}
</script>
