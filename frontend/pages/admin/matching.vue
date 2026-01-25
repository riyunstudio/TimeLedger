<template>
  <div class="p-4 md:p-6">
    <h1 class="text-2xl font-bold text-white mb-6">智慧媒合</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="glass-card p-4 md:p-6">
        <h2 class="text-lg font-semibold text-white mb-4">搜尋條件</h2>
        <form @submit.prevent="findMatches" class="space-y-4">
          <div>
            <label class="block text-slate-300 mb-2">課程時段</label>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-slate-400 text-sm mb-1">開始時間</label>
                <input
                  v-model="form.start_time"
                  type="datetime-local"
                  class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
                />
              </div>
              <div>
                <label class="block text-slate-400 text-sm mb-1">結束時間</label>
                <input
                  v-model="form.end_time"
                  type="datetime-local"
                  class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
                />
              </div>
            </div>
          </div>

          <div>
            <label class="block text-slate-300 mb-2">教室（可多選）</label>
            <div class="grid grid-cols-2 gap-2 max-h-48 overflow-y-auto p-2 rounded-lg bg-white/5 border border-white/10">
              <label
                v-for="room in rooms"
                :key="room.id"
                class="flex items-center gap-2 p-2 rounded-lg cursor-pointer hover:bg-white/10 transition-colors"
              >
                <input
                  type="checkbox"
                  :value="room.id"
                  v-model="form.room_ids"
                  class="w-4 h-4 rounded border-white/20 bg-white/10 text-primary-500 focus:ring-primary-500"
                />
                <span class="text-sm text-slate-300">{{ room.name }}</span>
              </label>
            </div>
            <p v-if="form.room_ids.length === 0" class="text-xs text-slate-500 mt-1">
              未選擇教室，將搜尋可用教室
            </p>
            <p v-else class="text-xs text-slate-500 mt-1">
              已選擇 {{ form.room_ids.length }} 間教室
            </p>
          </div>

          <div>
            <label class="block text-slate-300 mb-2">所需技能</label>
            <input
              v-model="form.skills"
              type="text"
              placeholder="例如：鋼琴、小提琴（用逗號分隔）"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
            />
          </div>

          <div class="flex flex-col sm:flex-row gap-3 pt-4">
            <button
              type="button"
              @click="clearForm"
              class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              清除
            </button>
            <button
              type="submit"
              :disabled="searching"
              class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50"
            >
              {{ searching ? '搜尋中...' : '開始媒合' }}
            </button>
          </div>
        </form>
      </div>

      <div class="glass-card p-4 md:p-6">
        <h2 class="text-lg font-semibold text-white mb-4">媒合結果</h2>

        <div v-if="!hasSearched" class="text-center py-12 text-slate-500">
          請設定搜尋條件並點擊「開始媒合」
        </div>

        <div v-else-if="matches.length === 0" class="text-center py-12 text-slate-500">
          沒有找到符合條件的老師
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="match in matches"
            :key="match.teacher_id"
            class="p-4 rounded-lg bg-white/5"
          >
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-3">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0">
                  <span class="text-white font-medium">{{ match.teacher_name?.charAt(0) || '?' }}</span>
                </div>
                <div>
                  <h3 class="text-white font-medium">{{ match.teacher_name }}</h3>
                  <p class="text-sm text-slate-400">匹配度 {{ match.match_score }}%</p>
                </div>
              </div>
              <div class="text-right">
                <div class="text-2xl font-bold text-primary-500">{{ match.match_score }}%</div>
              </div>
            </div>

            <div class="flex flex-wrap items-center gap-3 text-sm text-slate-400 mb-3">
              <span class="flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                技能匹配: {{ match.skill_match }}%
              </span>
              <span class="flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                </svg>
                評分: {{ match.rating?.toFixed(1) || '-' }}
              </span>
            </div>

            <!-- 可用教室標註 -->
            <div v-if="match.available_rooms?.length" class="mb-3">
              <p class="text-xs text-slate-500 mb-2">可授課教室：</p>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="room in match.available_rooms"
                  :key="room.id"
                  class="px-2 py-1 rounded-md text-xs bg-success-500/20 text-success-500"
                >
                  {{ room.name }}
                </span>
              </div>
            </div>

            <p v-if="match.notes" class="text-sm text-slate-400 mb-3">
              {{ match.notes }}
            </p>

            <button
              @click="selectTeacher(match)"
              class="w-full px-4 py-2 rounded-lg bg-primary-500/20 text-primary-500 hover:bg-primary-500/30 transition-colors text-sm"
            >
              選擇這位老師
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-6 glass-card p-4 md:p-6">
      <h2 class="text-lg font-semibold text-white mb-4">人才庫搜尋</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div>
          <label class="block text-slate-300 mb-2">城市</label>
          <input
            v-model="talentSearch.city"
            type="text"
            placeholder="例如：台北市"
            class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
          />
        </div>
        <div>
          <label class="block text-slate-300 mb-2">技能關鍵字</label>
          <input
            v-model="talentSearch.skills"
            type="text"
            placeholder="例如：鋼琴"
            class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
          />
        </div>
        <div>
          <label class="block text-slate-300 mb-2">標籤</label>
          <input
            v-model="talentSearch.hashtags"
            type="text"
            placeholder="例如：古典 兒童"
            class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
          />
        </div>
        <div class="flex items-end">
          <button
            @click="searchTalent"
            :disabled="talentSearching"
            class="w-full px-4 py-2 rounded-lg bg-secondary-500 text-white hover:bg-secondary-600 transition-colors disabled:opacity-50"
          >
            {{ talentSearching ? '搜尋中...' : '搜尋人才' }}
          </button>
        </div>
      </div>

      <div v-if="talentResults.length > 0" class="mt-4 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="teacher in talentResults"
          :key="teacher.id"
          class="p-4 rounded-lg bg-white/5"
        >
          <div class="flex items-center gap-3 mb-3">
            <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0">
              <span class="text-white text-sm font-medium">{{ teacher.name?.charAt(0) || '?' }}</span>
            </div>
            <div>
              <h4 class="text-white font-medium">{{ teacher.name }}</h4>
              <p class="text-xs text-slate-400 flex items-center gap-1">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                {{ teacher.city }}{{ teacher.district }}
              </p>
            </div>
          </div>

          <!-- 個人簡介 -->
          <p v-if="teacher.bio" class="text-sm text-slate-400 mb-3 line-clamp-2">{{ teacher.bio }}</p>

          <!-- 技能標籤 -->
          <div v-if="teacher.skills?.length" class="mb-3">
            <p class="text-xs text-slate-500 mb-2">專長技能</p>
            <div class="flex flex-wrap gap-1">
              <span
                v-for="(skill, index) in teacher.skills.slice(0, 4)"
                :key="index"
                class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs"
                :class="getSkillCategoryStyle(skill.category)"
              >
                <span>{{ getSkillCategoryIcon(skill.category) }}</span>
                <span>{{ skill.name }}</span>
              </span>
              <span
                v-if="teacher.skills.length > 4"
                class="px-2 py-1 rounded-md text-xs bg-slate-500/20 text-slate-400"
              >
                +{{ teacher.skills.length - 4 }}
              </span>
            </div>
          </div>

          <!-- 個人品牌標籤 -->
          <div v-if="teacher.personal_hashtags?.length" class="flex flex-wrap gap-1">
            <span
              v-for="(tag, index) in teacher.personal_hashtags.slice(0, 5)"
              :key="index"
              class="px-2 py-0.5 rounded-full text-xs bg-primary-500/20 text-primary-400"
            >
              {{ tag }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
import { SKILL_CATEGORIES } from '~/types'

 definePageMeta({
   middleware: 'auth-admin',
   layout: 'admin',
 })

 const notificationUI = useNotification()
const { warning: alertWarning, success: alertSuccess } = useAlert()
const searching = ref(false)
const talentSearching = ref(false)
const hasSearched = ref(false)
const matches = ref<any[]>([])
const talentResults = ref<any[]>([])
const rooms = ref<any[]>([])
const { getCenterId } = useCenterId()

const form = ref({
  start_time: '',
  end_time: '',
  room_ids: [] as number[],
  skills: ''
})

const talentSearch = ref({
  city: '',
  skills: '',
  hashtags: ''
})

// 技能類別相關函數
const getSkillCategoryIcon = (category: string): string => {
  return SKILL_CATEGORIES[category as keyof typeof SKILL_CATEGORIES]?.icon || '✨'
}

const getSkillCategoryStyle = (category: string): string => {
  return SKILL_CATEGORIES[category as keyof typeof SKILL_CATEGORIES]?.color || 'bg-slate-500/20 text-slate-400 border-slate-500/30'
}

const findMatches = async () => {
  if (!form.value.start_time || !form.value.end_time) {
    await alertWarning('請填寫開始時間和結束時間')
    return
  }

  searching.value = true
  hasSearched.value = true

  try {
    const api = useApi()
    const response = await api.post<{ code: number; datas: any[] }>(`/admin/smart-matching/matches`, {
      room_ids: form.value.room_ids.length > 0 ? form.value.room_ids : undefined,
      start_time: new Date(form.value.start_time).toISOString(),
      end_time: new Date(form.value.end_time).toISOString(),
      required_skills: form.value.skills.split(',').map(s => s.trim()).filter(Boolean)
    })
    matches.value = response.datas || []
  } catch (error) {
    console.error('Failed to find matches:', error)
    matches.value = []
  } finally {
    searching.value = false
  }
}

const searchTalent = async () => {
  talentSearching.value = true

  try {
    const api = useApi()
    const params = new URLSearchParams()
    if (talentSearch.value.city) params.append('city', talentSearch.value.city)
    if (talentSearch.value.skills) params.append('skills', talentSearch.value.skills)
    if (talentSearch.value.hashtags) params.append('hashtags', talentSearch.value.hashtags)
    
    const response = await api.get<{ code: number; datas: any[] }>(
      `/admin/smart-matching/talent/search?${params.toString()}`
    )
    talentResults.value = response.datas || []
  } catch (error) {
    console.error('Failed to search talent:', error)
    talentResults.value = []
  } finally {
    talentSearching.value = false
  }
}

const selectTeacher = (match: any) => {
  alertSuccess(`已選擇 ${match.teacher_name}`)
}

const clearForm = () => {
  form.value = {
    start_time: '',
    end_time: '',
    room_ids: [],
    skills: ''
  }
  hasSearched.value = false
  matches.value = []
}

const fetchRooms = async () => {
  try {
    const api = useApi()
    const centerId = getCenterId()
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    rooms.value = response.datas || []
  } catch (error) {
    console.error('Failed to fetch rooms:', error)
  }
}

onMounted(() => {
  fetchRooms()
})
</script>
