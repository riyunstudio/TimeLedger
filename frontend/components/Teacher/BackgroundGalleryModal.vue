<template>
  <BaseModal
    :model-value="isOpen"
    :title="title"
    size="lg"
    @update:modelValue="$emit('update:isOpen', $event)"
    @close="$emit('close')"
  >
    <template #default>
      <div class="space-y-6">
        <!-- 上傳區域 -->
        <div
          class="border-2 border-dashed rounded-xl p-6 text-center cursor-pointer transition-all"
          :class="[
            isDragging ? 'border-primary-500 bg-primary-500/10' : 'border-white/20 hover:border-white/40 hover:bg-white/5',
            isUploading ? 'opacity-50 cursor-not-allowed' : ''
          ]"
          @dragover.prevent="isDragging = true"
          @dragleave.prevent="isDragging = false"
          @drop.prevent="handleDrop"
          @click="triggerFileInput"
        >
          <input
            ref="fileInput"
            type="file"
            accept="image/jpeg,image/png"
            class="hidden"
            @change="handleFileSelect"
          >

          <div v-if="isUploading" class="py-4">
            <div class="animate-spin w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full mx-auto mb-3"></div>
            <p class="text-slate-400">正在上傳圖片...</p>
          </div>

          <template v-else>
            <svg class="w-12 h-12 mx-auto mb-3 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <p class="text-slate-300 font-medium mb-1">點擊或拖曳上傳圖片</p>
            <p class="text-slate-500 text-sm">支援 JPEG、PNG 格式，最大 5MB</p>
          </template>
        </div>

        <!-- 錯誤提示 -->
        <div
          v-if="errorMessage"
          class="bg-red-500/10 border border-red-500/30 rounded-lg p-3 text-red-400 text-sm"
        >
          {{ errorMessage }}
        </div>

        <!-- 背景圖片列表 -->
        <div v-if="backgrounds.length > 0">
          <h3 class="text-sm font-medium text-slate-400 mb-3">我的背景圖片</h3>
          <div class="grid grid-cols-3 gap-4">
            <div
              v-for="bg in backgrounds"
              :key="bg.id"
              class="group relative aspect-video rounded-lg overflow-hidden bg-white/5 cursor-pointer ring-2 transition-all"
              :class="[
                selectedPath === bg.url ? 'ring-primary-500' : 'ring-transparent hover:ring-white/20'
              ]"
              @click="$emit('select', bg.url)"
            >
              <img
                :src="bg.url"
                :alt="bg.file_url"
                class="w-full h-full object-cover"
              >

              <!-- 懸停時的遮罩 -->
              <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                <button
                  v-if="selectable"
                  class="px-3 py-1 bg-primary-500 text-white text-sm rounded-lg hover:bg-primary-600 transition-colors"
                  @click.stop="$emit('select', bg.url)"
                >
                  選擇
                </button>
                <button
                  class="p-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
                  @click.stop="confirmDelete(bg)"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>

              <!-- 已選擇標記 -->
              <div
                v-if="selectedPath === bg.url"
                class="absolute top-2 right-2 w-6 h-6 bg-primary-500 rounded-full flex items-center justify-center"
              >
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </div>
            </div>
          </div>
        </div>

        <!-- 空狀態 -->
        <div
          v-else-if="!isLoading"
          class="text-center py-12"
        >
          <svg class="w-16 h-16 mx-auto mb-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <p class="text-slate-400">尚未上傳任何背景圖片</p>
        </div>

        <!-- 載入中 -->
        <div v-if="isLoading" class="flex justify-center py-12">
          <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full"></div>
        </div>
      </div>
    </template>

    <!-- Modal Footer -->
    <template #footer>
      <div class="flex justify-end gap-3">
        <button
          class="px-4 py-2 text-slate-400 hover:text-white transition-colors"
          @click="$emit('close')"
        >
          取消
        </button>
      </div>
    </template>
  </BaseModal>

  <!-- 刪除確認對話框 -->
  <GlobalAlert
    :is-open="showDeleteConfirm"
    type="warning"
    title="刪除背景圖片"
    :message="`確定要刪除「${deletingBg?.file_url}」嗎？此操作無法復原。`"
    confirm-text="刪除"
    confirm-class="bg-red-500 hover:bg-red-600"
    @confirm="handleDelete"
    @cancel="showDeleteConfirm = false"
  />
</template>

<script setup lang="ts">
import type { BackgroundImage } from '~/stores/useProfileStore'
import { alertError, alertSuccess } from '~/composables/useAlert'
import GlobalAlert from '~/components/base/GlobalAlert.vue'
import BaseModal from '~/components/base/BaseModal.vue'

const props = withDefaults(defineProps<{
  isOpen: boolean
  title?: string
  selectable?: boolean
  selectedPath?: string | null
}>(), {
  title: '背景圖片管理',
  selectable: false,
  selectedPath: null,
})

const emit = defineEmits<{
  'update:isOpen': [value: boolean]
  close: []
  select: [path: string]
}>()

const profileStore = useProfileStore()

const fileInput = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
const isUploading = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')
const showDeleteConfirm = ref(false)
const deletingBg = ref<BackgroundImage | null>(null)

// 計算屬性
const backgrounds = computed(() => profileStore.backgrounds)

// 方法
const triggerFileInput = () => {
  if (!isUploading.value) {
    fileInput.value?.click()
  }
}

const validateFile = (file: File): boolean => {
  const maxSize = 5 * 1024 * 1024 // 5MB
  const allowedTypes = ['image/jpeg', 'image/png']

  if (!allowedTypes.includes(file.type)) {
    errorMessage.value = '只支援 JPEG 和 PNG 格式'
    return false
  }

  if (file.size > maxSize) {
    errorMessage.value = '圖片大小不能超過 5MB'
    return false
  }

  errorMessage.value = ''
  return true
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (file && validateFile(file)) {
    uploadFile(file)
  }

  // 清空 input 以便重複選擇同一檔案
  target.value = ''
}

const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  const file = event.dataTransfer?.files[0]

  if (file && validateFile(file)) {
    uploadFile(file)
  }
}

const uploadFile = async (file: File) => {
  isUploading.value = true
  errorMessage.value = ''

  try {
    await profileStore.uploadBackground(file)
    await alertSuccess('上傳成功')
  } catch (error) {
    console.error('上傳失敗:', error)
    await alertError('上傳失敗，請稍後再試')
  } finally {
    isUploading.value = false
  }
}

const confirmDelete = (bg: BackgroundImage) => {
  deletingBg.value = bg
  showDeleteConfirm.value = true
}

const handleDelete = async () => {
  if (!deletingBg.value) return

  try {
    await profileStore.deleteBackground(deletingBg.value.id)
    showDeleteConfirm.value = false
    deletingBg.value = null
    await alertSuccess('刪除成功')
  } catch (error) {
    console.error('刪除失敗:', error)
    await alertError('刪除失敗，請稍後再試')
  }
}

// 生命週期
onMounted(async () => {
  if (props.isOpen) {
    await loadBackgrounds()
  }
})

watch(() => props.isOpen, async (newVal) => {
  if (newVal) {
    await loadBackgrounds()
  }
})

const loadBackgrounds = async () => {
  isLoading.value = true
  try {
    await profileStore.fetchBackgrounds()
  } catch (error) {
    console.error('載入背景圖片失敗:', error)
  } finally {
    isLoading.value = false
  }
}
</script>
