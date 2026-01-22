<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          編輯個人檔案
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-4 space-y-4">
        <div class="text-center pt-2">
          <div
            class="w-20 h-20 sm:w-24 sm:h-24 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center mx-auto mb-3"
          >
            <span class="text-3xl sm:text-4xl font-bold text-white">
              {{ form.name?.charAt(0) || 'T' }}
            </span>
          </div>
          <p class="text-xs sm:text-sm text-slate-400">點擊頭像更換照片</p>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">姓名</label>
          <input
            v-model="form.name"
            type="text"
            class="input-field text-sm sm:text-base"
            required
          />
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">Email</label>
          <input
            v-model="form.email"
            type="email"
            class="input-field text-sm sm:text-base"
            required
          />
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">Bio</label>
          <textarea
            v-model="form.bio"
            placeholder="簡短介紹自己..."
            rows="3"
            class="input-field resize-none text-sm sm:text-base"
          />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">縣市</label>
            <select v-model="form.city" class="input-field text-sm sm:text-base">
              <option value="">請選擇</option>
              <option v-for="city in cities" :key="city" :value="city">
                {{ city }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">區域</label>
            <select v-model="form.district" class="input-field text-sm sm:text-base">
              <option value="">請選擇</option>
              <option v-for="district in districts" :key="district" :value="district">
                {{ district }}
              </option>
            </select>
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
            :disabled="loading"
            class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            {{ loading ? '儲存中...' : '儲存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  close: []
}>()

const authStore = useAuthStore()

const loading = ref(false)
const form = ref({
  name: authStore.user?.name || '',
  email: authStore.user?.email || '',
  bio: authStore.user?.bio || '',
  city: authStore.user?.city || '',
  district: authStore.user?.district || '',
})

const cities = [
  '臺北市',
  '新北市',
  '桃園市',
  '臺中市',
  '臺南市',
  '高雄市',
  '基隆市',
  '新竹市',
  '新竹縣',
  '苗栗縣',
  '彰化縣',
  '南投縣',
  '雲林縣',
  '嘉義市',
  '嘉義縣',
  '屏東縣',
  '宜蘭縣',
  '花蓮縣',
  '臺東縣',
  '澎湖縣',
  '金門縣',
  '連江縣',
]

const districts = [
  '中正區',
  '大同區',
  '中山區',
  '松山區',
  '大安區',
  '萬華區',
  '文山區',
  '南港區',
  '內湖區',
  '士林區',
  '北投區',
  '信義區',
]

const handleSubmit = async () => {
  loading.value = true

  try {
    const api = useApi()
    await api.put('/teacher/me/profile', form.value)

    authStore.user = { ...authStore.user, ...form.value } as any
    localStorage.setItem('user', JSON.stringify(authStore.user))

    emit('close')
  } catch (error) {
    console.error('Failed to update profile:', error)
    alert('儲存失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}
</script>
