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

      <form @submit.prevent="handleSubmit" class="p-4 space-y-4">
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
            <select v-model="form.city" class="input-field text-sm sm:text-base" :disabled="loadingCities">
              <option value="">請選擇</option>
              <option v-for="city in cities" :key="city" :value="city">
                {{ city }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">區域</label>
            <select v-model="form.district" class="input-field text-sm sm:text-base" :disabled="loadingCities || !form.city">
              <option value="">請選擇</option>
              <option v-for="district in districts" :key="district" :value="district">
                {{ district }}
              </option>
            </select>
          </div>
        </div>

        <!-- 專業連結 -->
        <div>
          <label class="block text-slate-300 mb-3 font-medium text-sm sm:text-base">專業連結</label>
          <div class="space-y-3">
            <!-- Instagram -->
            <div class="flex items-center gap-3">
              <div class="flex-shrink-0 w-10 h-10 rounded-xl bg-pink-500/20 flex items-center justify-center">
                <svg class="w-5 h-5 text-pink-500" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"/>
                </svg>
              </div>
              <input
                v-model="form.public_contact_info.instagram"
                type="text"
                placeholder="Instagram 帳號"
                class="input-field text-sm sm:text-base flex-1"
              />
            </div>

            <!-- YouTube -->
            <div class="flex items-center gap-3">
              <div class="flex-shrink-0 w-10 h-10 rounded-xl bg-red-500/20 flex items-center justify-center">
                <svg class="w-5 h-5 text-red-500" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/>
                </svg>
              </div>
              <input
                v-model="form.public_contact_info.youtube"
                type="text"
                placeholder="YouTube 頻道"
                class="input-field text-sm sm:text-base flex-1"
              />
            </div>

            <!-- Website -->
            <div class="flex items-center gap-3">
              <div class="flex-shrink-0 w-10 h-10 rounded-xl bg-primary-500/20 flex items-center justify-center">
                <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                </svg>
              </div>
              <input
                v-model="form.public_contact_info.website"
                type="text"
                placeholder="個人網站"
                class="input-field text-sm sm:text-base flex-1"
              />
            </div>

            <!-- Other -->
            <div class="flex items-center gap-3">
              <div class="flex-shrink-0 w-10 h-10 rounded-xl bg-slate-500/20 flex items-center justify-center">
                <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                </svg>
              </div>
              <input
                v-model="form.public_contact_info.other"
                type="text"
                placeholder="其他連結"
                class="input-field text-sm sm:text-base flex-1"
              />
            </div>
          </div>
        </div>

        <!-- 個人品牌標籤 -->
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
            個人品牌標籤
            <span class="text-xs text-slate-500 font-normal ml-2">(3-5 個標籤)</span>
          </label>
          <div class="flex flex-wrap gap-2 mb-2">
            <span
              v-for="(tag, index) in form.personal_hashtags"
              :key="index"
              class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm bg-primary-500/20 text-primary-400"
            >
              {{ tag.startsWith('#') ? tag : '#' + tag }}
              <button
                type="button"
                @click="removeHashtag(index)"
                class="hover:text-primary-300"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </span>
          </div>
          <div class="relative">
            <input
              v-model="hashtagInput"
              @keydown.enter.prevent="addHashtag"
              @keydown.comma.prevent="addHashtag"
              @keydown.backspace="handleHashtagBackspace"
              @input="onHashtagInput"
              type="text"
              placeholder="輸入標籤後按 Enter..."
              class="input-field text-sm sm:text-base w-full"
              :disabled="form.personal_hashtags.length >= 5"
            />
            <div
              v-if="hashtagSuggestions.length > 0"
              class="absolute z-10 w-full mt-1 bg-slate-800 border border-white/10 rounded-lg shadow-lg max-h-40 overflow-y-auto"
            >
              <button
                v-for="suggestion in hashtagSuggestions"
                :key="suggestion.id"
                type="button"
                @click="selectHashtagSuggestion(suggestion)"
                class="w-full px-3 py-2 text-left text-sm text-slate-300 hover:bg-white/5 flex items-center justify-between"
              >
                <span>{{ suggestion.name.startsWith('#') ? suggestion.name : '#' + suggestion.name }}</span>
                <span class="text-xs text-slate-500">{{ suggestion.usage_count }} 次使用</span>
              </button>
            </div>
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
            :disabled="loading || !isFormValid"
            class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            {{ loading ? '儲存中...' : '儲存' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertError } from '~/composables/useAlert'
import { watch } from 'vue'

const emit = defineEmits<{
  close: []
}>()

const authStore = useAuthStore()

interface PublicContactInfo {
  instagram?: string
  youtube?: string
  website?: string
  other?: string
}

const loading = ref(false)
const hashtagInput = ref('')
const hashtagSuggestions = ref<{ id: number; name: string; usage_count: number }[]>([])
let hashtagSearchTimeout: ReturnType<typeof setTimeout> | null = null

const parseContactInfo = (info: string | PublicContactInfo | undefined): PublicContactInfo => {
  if (!info) return { instagram: '', youtube: '', website: '', other: '' }
  if (typeof info === 'string') {
    try {
      return JSON.parse(info)
    } catch {
      return { instagram: '', youtube: '', website: '', other: '' }
    }
  }
  return {
    instagram: info.instagram || '',
    youtube: info.youtube || '',
    website: info.website || '',
    other: info.other || '',
  }
}

// 輔助函數：提取標籤名稱（確保是字串）
const extractTagName = (tag: any): string => {
  if (!tag) return ''
  if (typeof tag === 'string') return tag
  if (typeof tag === 'object') {
    // 處理 { hashtag: { name: "xxx" } } 結構
    if (tag.hashtag && typeof tag.hashtag === 'object') {
      return tag.hashtag.name || ''
    }
    // 處理 { name: "xxx" } 結構
    if (tag.name) return tag.name
  }
  return ''
}

const form = ref({
  name: authStore.user?.name || '',
  email: authStore.user?.email || '',
  bio: authStore.user?.bio || '',
  city: authStore.user?.city || '',
  district: authStore.user?.district || '',
  public_contact_info: parseContactInfo(authStore.user?.public_contact_info),
  personal_hashtags: ((authStore.user as any)?.personal_hashtags || []).map((h: any) => extractTagName(h)).filter((t: string) => t),
})

// 城市區域資料（從 API 載入）
interface GeoCity {
  id: number
  name: string
  districts: GeoDistrict[]
}

interface GeoDistrict {
  id: number
  city_id: number
  name: string
}

const citiesData = ref<GeoCity[]>([])
const loadingCities = ref(false)

const cities = computed(() => {
  return citiesData.value.map(city => city.name)
})

const districts = computed(() => {
  const selectedCity = citiesData.value.find(city => city.name === form.value.city)
  return selectedCity?.districts.map(d => d.name) || []
})

const fetchCities = async () => {
  try {
    loadingCities.value = true
    const api = useApi()
    const response = await api.get<{ code: number; data: GeoCity[] }>('/geo/cities')
    if (response.code === 0 && response.data) {
      citiesData.value = response.data
    }
  } catch (error) {
    console.error('Failed to fetch cities:', error)
    // 如果 API 失敗，使用備用資料
    citiesData.value = []
  } finally {
    loadingCities.value = false
  }
}

const isFormValid = computed(() => {
  return form.value.name.trim() !== '' &&
    form.value.email.trim() !== '' &&
    form.value.personal_hashtags.length <= 5
})

// 監聽城市變化，清空區域選擇
watch(() => form.value.city, () => {
  form.value.district = ''
})

const addHashtag = () => {
  let tag = hashtagInput.value.trim()
  if (!tag) return

  // 移除 # 符號
  tag = tag.replace(/^#+/, '')

  // 至少 2 個字元
  if (tag.length < 2) {
    return
  }

  // 不重複
  if (form.value.personal_hashtags.includes(tag)) {
    hashtagInput.value = ''
    hashtagSuggestions.value = []
    return
  }

  // 未滿 5 個
  if (form.value.personal_hashtags.length < 5) {
    form.value.personal_hashtags.push(tag)
  }

  hashtagInput.value = ''
  hashtagSuggestions.value = []
}

const removeHashtag = (index: number) => {
  form.value.personal_hashtags.splice(index, 1)
}

const handleHashtagBackspace = () => {
  if (hashtagInput.value === '' && form.value.personal_hashtags.length > 0) {
    form.value.personal_hashtags.pop()
  }
}

const onHashtagInput = () => {
  if (hashtagSearchTimeout) {
    clearTimeout(hashtagSearchTimeout)
  }

  const query = hashtagInput.value.replace(/^#+/, '').trim()
  if (query.length < 1) {
    hashtagSuggestions.value = []
    return
  }

  // 延遲 300ms 發送搜尋請求
  hashtagSearchTimeout = setTimeout(async () => {
    try {
      const api = useApi()
      const response = await api.get<{ code: number; datas: any[] }>('/hashtags/search', { q: query })
      hashtagSuggestions.value = response.datas || []
    } catch {
      hashtagSuggestions.value = []
    }
  }, 300)
}

const selectHashtagSuggestion = (suggestion: { id: number; name: string; usage_count: number }) => {
  const tag = suggestion.name
  if (!form.value.personal_hashtags.includes(tag) && form.value.personal_hashtags.length < 5) {
    form.value.personal_hashtags.push(tag)
  }
  hashtagInput.value = ''
  hashtagSuggestions.value = []
}

const handleSubmit = async () => {
  if (!isFormValid.value) {
    await alertError('請填寫完整資訊')
    return
  }

  loading.value = true

  try {
    const api = useApi()
    await api.put('/teacher/me/profile', {
      name: form.value.name,
      email: form.value.email,
      bio: form.value.bio,
      city: form.value.city,
      district: form.value.district,
      public_contact_info: JSON.stringify(form.value.public_contact_info),
      personal_hashtags: form.value.personal_hashtags,
    })

    const updatedUser = {
      ...authStore.user,
      name: form.value.name,
      email: form.value.email,
      bio: form.value.bio,
      city: form.value.city,
      district: form.value.district,
      public_contact_info: JSON.stringify(form.value.public_contact_info),
      personal_hashtags: form.value.personal_hashtags,
    }
    authStore.user = updatedUser as any
    localStorage.setItem('teacher_user', JSON.stringify(updatedUser))

    emit('close')
  } catch (error) {
    console.error('Failed to update profile:', error)
    await alertError('儲存失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCities()
})
</script>
