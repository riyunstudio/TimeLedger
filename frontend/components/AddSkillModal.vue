<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-md sm:max-w-lg max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          {{ isEditing ? '編輯技能' : '新增技能' }}
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
            <option v-for="cat in categories" :key="cat.value" :value="cat.value">
              {{ cat.icon }} {{ cat.label }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">技能名稱</label>
          <input
            v-model="form.skill_name"
            @input="checkDuplicate"
            type="text"
            placeholder="例：鋼琴、小提琴、吉他..."
            class="input-field text-sm sm:text-base"
            required
          />
          <!-- 重複警告 -->
          <div
            v-if="isDuplicate"
            class="mt-2 p-2 rounded-lg bg-warning-500/20 border border-warning-500/30 flex items-center gap-2"
          >
            <svg class="w-4 h-4 text-warning-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <span class="text-sm text-warning-500">此技能已存在</span>
          </div>
        </div>

        <!-- 技能標籤 -->
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">
            技能標籤
            <span class="text-xs text-slate-500 font-normal ml-2">(選填，最多 3 個)</span>
          </label>
          <div class="flex flex-wrap gap-2 mb-2">
            <span
              v-for="(tag, index) in form.hashtags"
              :key="index"
              class="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs bg-primary-500/20 text-primary-400"
            >
              #{{ tag }}
              <button
                type="button"
                @click="removeHashtag(index)"
                class="hover:text-primary-300"
              >
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
              :disabled="form.hashtags.length >= 3"
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
                <span>#{{ suggestion.name }}</span>
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
            :disabled="loading || isDuplicate || !form.category || !form.skill_name"
            class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            {{ loading ? (isEditing ? '儲存中...' : '新增中...') : (isEditing ? '儲存' : '新增') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertError } from '~/composables/useAlert'
import { SKILL_CATEGORIES } from '~/types'

const props = defineProps<{
  skill?: any
  existingSkills?: any[]
}>()

const emit = defineEmits<{
  close: []
  added: []
  updated: []
}>()

const teacherStore = useTeacherStore()
const loading = ref(false)
const hashtagInput = ref('')
const hashtagSuggestions = ref<{ id: number; name: string; usage_count: number }[]>([])
let hashtagSearchTimeout: ReturnType<typeof setTimeout> | null = null

const isEditing = computed(() => !!props.skill)

const categories = Object.entries(SKILL_CATEGORIES).map(([value, data]) => ({
  value,
  ...data,
}))

const form = ref({
  category: props.skill?.category || '' as 'MUSIC' | 'ART' | 'DANCE' | 'LANGUAGE' | 'SPORTS' | 'OTHER' | '',
  skill_name: props.skill?.skill_name || '',
  hashtags: props.skill?.hashtags?.map((h: any) => h.hashtag?.name || h.name || h).slice(0, 3) || [] as string[],
})

const isDuplicate = computed(() => {
  if (!form.value.skill_name || !props.existingSkills) return false
  const normalizedInput = form.value.skill_name.toLowerCase().trim()
  return props.existingSkills.some((s: any) => {
    if (isEditing.value && s.id === props.skill.id) return false
    return s.skill_name.toLowerCase().trim() === normalizedInput
  })
})

const checkDuplicate = () => {
  // Trigger reactivity
}

const addHashtag = () => {
  let tag = hashtagInput.value.trim()
  if (!tag) return

  tag = tag.replace(/^#+/, '')

  if (tag.length < 2) return
  if (form.value.hashtags.includes(tag)) {
    hashtagInput.value = ''
    hashtagSuggestions.value = []
    return
  }
  if (form.value.hashtags.length < 3) {
    form.value.hashtags.push(tag)
  }

  hashtagInput.value = ''
  hashtagSuggestions.value = []
}

const removeHashtag = (index: number) => {
  form.value.hashtags.splice(index, 1)
}

const handleHashtagBackspace = () => {
  if (hashtagInput.value === '' && form.value.hashtags.length > 0) {
    form.value.hashtags.pop()
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
  if (!form.value.hashtags.includes(tag) && form.value.hashtags.length < 3) {
    form.value.hashtags.push(tag)
  }
  hashtagInput.value = ''
  hashtagSuggestions.value = []
}

const handleSubmit = async () => {
  if (isDuplicate.value || !form.value.category || !form.value.skill_name) {
    return
  }

  loading.value = true

  try {
    if (isEditing.value) {
      await teacherStore.updateSkill(props.skill.id, {
        category: form.value.category,
        skill_name: form.value.skill_name,
        hashtags: form.value.hashtags,
      })
      emit('updated')
    } else {
      // 轉換標籤名稱為 hashtag_ids
      const hashtagIds: number[] = []
      for (const tagName of form.value.hashtags) {
        // 搜尋現有標籤或創建新標籤
        try {
          const api = useApi()
          const response = await api.get<{ code: number; datas: any[] }>('/hashtags/search', { q: tagName })
          const existing = response.datas?.find((h: any) => h.name === tagName || h.name === '#' + tagName)
          if (existing) {
            hashtagIds.push(existing.id)
          } else {
            // 創建新標籤
            const createResponse = await api.post<{ code: number; datas: any }>('/hashtags', { name: '#' + tagName })
            if (createResponse.datas?.id) {
              hashtagIds.push(createResponse.datas.id)
            }
          }
        } catch (e) {
          console.error('Failed to process hashtag:', tagName, e)
        }
      }

      await teacherStore.createSkill({
        category: form.value.category,
        skill_name: form.value.skill_name,
        hashtag_ids: hashtagIds,
      })
      emit('added')
    }
    emit('close')
  } catch (error: any) {
    console.error('Failed to save skill:', error)
    await alertError(error?.message || '儲存失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}
</script>
