<template>
  <BaseModal
    v-model="isOpen"
    title="新增證照"
    size="md"
    @close="handleClose"
  >
    <form @submit.prevent="handleSubmit" class="space-y-4">
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
          <Icon v-if="!fileName" icon="upload" size="3xl" class="mx-auto mb-3 text-slate-400" />
          <p v-if="!fileName" class="text-slate-400 text-sm sm:text-base">點擊或拖曳檔案至此</p>
          <p v-else class="text-slate-100 text-sm sm:text-base">{{ fileName }}</p>
        </div>
      </div>
    </form>

    <template #footer>
      <div class="flex gap-3 pt-2">
        <button
          type="button"
          @click="handleClose"
          class="flex-1 glass-btn py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
        >
          取消
        </button>
        <button
          type="button"
          :disabled="loading || uploading"
          class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          @click="handleSubmit"
        >
          {{ uploading ? '上傳中...' : loading ? '儲存中...' : '新增' }}
        </button>
      </div>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { alertError, alertSuccess } from '~/composables/useAlert'
import { formatDateToString } from '~/composables/useTaiwanTime'
import Icon from '~/components/base/Icon.vue'

const emit = defineEmits<{
  close: []
  added: []
}>()

const isOpen = ref(true)
const profileStore = useProfileStore()
const api = useApi()
const loading = ref(false)
const uploading = ref(false)
const fileInput = ref<HTMLInputElement>()
const fileName = ref('')
const selectedFile = ref<File | null>(null)

const form = ref({
  name: '',
  issued_at: '',
})

const formatDateTimeForApi = (datetimeLocal: string): string => {
  if (!datetimeLocal) return ''
  // 使用台灣時區格式化
  const date = new Date(datetimeLocal)
  return formatDateToString(date)
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    selectedFile.value = target.files[0]
    fileName.value = target.files[0].name
  }
}

const handleClose = () => {
  emit('close')
}

const handleSubmit = async () => {
  if (!form.value.name) {
    await alertError('請輸入證照名稱')
    return
  }

  loading.value = true
  uploading.value = true

  try {
    let fileUrl: string | undefined

    // 如果有選擇檔案，先上傳
    if (selectedFile.value) {
      const uploadResponse = await api.upload<{ code: number; message: string; datas: { file_url: string; file_name: string; file_size: number } }>(
        '/teacher/me/certificates/upload',
        selectedFile.value
      )

      if (uploadResponse.code === 0) {
        fileUrl = uploadResponse.datas.file_url
      } else {
        throw new Error(uploadResponse.message || '上傳失敗')
      }
    }

    // 建立證照記錄
    await profileStore.createCertificate({
      name: form.value.name,
      file_url: fileUrl,
      issued_at: formatDateTimeForApi(form.value.issued_at),
    })

    await profileStore.fetchCertificates()
    await alertSuccess('證照新增成功')
    emit('added')
    handleClose()
  } catch (error: any) {
    console.error('Failed to add certificate:', error)
    await alertError(error.message || '新增失敗，請稍後再試')
  } finally {
    loading.value = false
    uploading.value = false
  }
}
</script>
