<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg sm:max-w-xl max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          匯出課表圖片
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-4 space-y-6">
        <div>
          <label class="block text-slate-300 mb-3 font-medium">主題選擇</label>
          <div class="grid grid-cols-3 gap-3">
            <button
              v-for="theme in themes"
              :key="theme.id"
              type="button"
              @click="selectedTheme = theme.id"
              class="aspect-square rounded-xl p-3 transition-all duration-300 hover:scale-105"
              :class="selectedTheme === theme.id ? 'ring-2 ring-white' : ''"
              :style="{ background: theme.preview }"
            >
              <span class="text-sm font-medium" :style="{ color: theme.textColor }">
                {{ theme.name }}
              </span>
            </button>
          </div>
        </div>

        <div>
          <label class="block text-slate-300 mb-3 font-medium">顯示模式</label>
          <div class="grid grid-cols-2 gap-3">
            <button
              type="button"
              @click="selectedView = 'grid'"
              class="p-4 rounded-xl glass-btn transition-all duration-300"
              :class="selectedView === 'grid' ? 'bg-primary-500/30 border-primary-500' : ''"
            >
              <div class="text-center">
                <div class="grid grid-cols-3 gap-1 w-12 mx-auto mb-2">
                  <div v-for="i in 6" :key="i" class="aspect-square bg-slate-600/50 rounded"></div>
                </div>
                <span class="text-sm font-medium text-slate-100">網格</span>
              </div>
            </button>

            <button
              type="button"
              @click="selectedView = 'list'"
              class="p-4 rounded-xl glass-btn transition-all duration-300"
              :class="selectedView === 'list' ? 'bg-primary-500/30 border-primary-500' : ''"
            >
              <div class="text-center">
                <div class="space-y-1 w-12 mx-auto mb-2">
                  <div v-for="i in 3" :key="i" class="h-3 bg-slate-600/50 rounded"></div>
                </div>
                <span class="text-sm font-medium text-slate-100">列表</span>
              </div>
            </button>
          </div>
        </div>

        <div>
          <label class="block text-slate-300 mb-3 font-medium">內容選項</label>
          <div class="space-y-3">
            <label class="flex items-center gap-3 cursor-pointer">
              <input
                v-model="options.showPersonalInfo"
                type="checkbox"
                class="w-5 h-5 rounded accent-primary-500"
              />
              <span class="text-slate-100">顯示個人品牌資訊（頭像、姓名、標籤）</span>
            </label>

            <label class="flex items-center gap-3 cursor-pointer">
              <input
                v-model="options.compactMode"
                type="checkbox"
                class="w-5 h-5 rounded accent-primary-500"
              />
              <span class="text-slate-100">精簡模式（縮小字體與間距，放更多課程）</span>
            </label>

            <label class="flex items-center gap-3 cursor-pointer">
              <input
                v-model="options.showTime"
                type="checkbox"
                class="w-5 h-5 rounded accent-primary-500"
              />
              <span class="text-slate-100">顯示課程時間</span>
            </label>

            <label class="flex items-center gap-3 cursor-pointer">
              <input
                v-model="options.privacyMode"
                type="checkbox"
                class="w-5 h-5 rounded accent-primary-500"
              />
              <span class="text-slate-100">隱私模式（個人行程僅顯示「已保留」）</span>
            </label>

            <label class="flex items-center gap-3 cursor-pointer">
              <input
                v-model="options.showFullWeek"
                type="checkbox"
                class="w-5 h-5 rounded accent-primary-500"
              />
              <span class="text-slate-100">顯示完整週（預設：3 天）</span>
            </label>
          </div>
        </div>

        <div>
          <label class="block text-slate-300 mb-3 font-medium">格式</label>
          <div class="grid grid-cols-2 gap-3">
            <button
              type="button"
              @click="selectedFormat = 'story'"
              class="p-4 rounded-xl glass-btn transition-all duration-300"
              :class="selectedFormat === 'story' ? 'bg-primary-500/30 border-primary-500' : ''"
            >
              <div class="text-center">
                <div class="aspect-[9/16] w-16 mx-auto mb-2 rounded-lg bg-slate-700 flex items-center justify-center">
                  <svg class="w-6 h-6 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                  </svg>
                </div>
                <span class="text-sm font-medium text-slate-100">IG Story</span>
                <p class="text-xs text-slate-400">9:16</p>
              </div>
            </button>

            <button
              type="button"
              @click="selectedFormat = 'post'"
              class="p-4 rounded-xl glass-btn transition-all duration-300"
              :class="selectedFormat === 'post' ? 'bg-primary-500/30 border-primary-500' : ''"
            >
              <div class="text-center">
                <div class="aspect-[4/3] w-16 mx-auto mb-2 rounded-lg bg-slate-700 flex items-center justify-center">
                  <svg class="w-6 h-6 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                  </svg>
                </div>
                <span class="text-sm font-medium text-slate-100">貼文</span>
                <p class="text-xs text-slate-400">4:3</p>
              </div>
            </button>
          </div>
        </div>

        <div class="flex flex-col sm:flex-row gap-3">
          <button
            @click="handlePreview"
            :disabled="generating"
            class="flex-1 glass-btn py-3 rounded-xl font-medium order-2 sm:order-1"
          >
            {{ generating ? '生成中...' : '預覽' }}
          </button>
          <button
            @click="handleExport"
            :disabled="generating"
            class="flex-1 btn-primary py-3 rounded-xl font-medium order-1 sm:order-2"
          >
            {{ generating ? '生成中...' : '下載並分享' }}
          </button>
        </div>
      </div>
    </div>

    <TeacherExportPreviewModal
      v-if="showPreview"
      :theme="getThemeById(selectedTheme)"
      :options="options"
      :format="selectedFormat"
      :view="selectedView"
      @close="showPreview = false"
    />
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  close: []
}>()

const generating = ref(false)
const selectedTheme = ref('midnight')
const selectedFormat = ref('story')
const selectedView = ref('grid')
const showPreview = ref(false)

const options = ref({
  showPersonalInfo: true,
  compactMode: false,
  showTime: true,
  privacyMode: false,
  showFullWeek: false,
})

const themes = [
  {
    id: 'midnight',
    name: 'Midnight Glow',
    preview: 'linear-gradient(135deg, #1e3a8a 0%, #6366F1 25%, #A855F7 50%, #1e3a8a 75%, #0f172a 100%)',
    textColor: '#ffffff',
  },
  {
    id: 'emerald',
    name: 'Emerald Mist',
    preview: 'linear-gradient(135deg, #064e3b 0%, #10B981 25%, #34D399 50%, #064e3b 75%, #022c22 100%)',
    textColor: '#ffffff',
  },
  {
    id: 'sunset',
    name: 'Sunset Quartz',
    preview: 'linear-gradient(135deg, #78350f 0%, #F59E0B 25%, #FBBF24 50%, #78350f 75%, #451a03 100%)',
    textColor: '#ffffff',
  },
]

const getThemeById = (id: string) => {
  return themes.find(t => t.id === id)
}

const handlePreview = () => {
  showPreview.value = true
}

const handleExport = async () => {
  generating.value = true

  try {
    showPreview.value = true

    await new Promise(resolve => setTimeout(resolve, 2000))

    await alert('圖片已生成！請點擊分享按鈕分享到 LINE 或其他平台。')
  } catch (error) {
    await alert('匯出失敗，請稍後再試')
  } finally {
    generating.value = false
  }
}
</script>
