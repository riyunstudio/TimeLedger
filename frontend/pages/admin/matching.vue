<template>
  <div class="min-h-screen bg-slate-900">
    <main class="p-6">
      <h1 class="text-2xl font-bold text-white mb-6">智慧媒合</h1>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="glass-card p-6">
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
              <label class="block text-slate-300 mb-2">教室</label>
              <select
                v-model="form.room_id"
                class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
              >
                <option value="">選擇教室</option>
                <option v-for="room in rooms" :key="room.id" :value="room.id">
                  {{ room.name }}
                </option>
              </select>
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

            <div class="flex gap-3 pt-4">
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

        <div class="glass-card p-6">
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
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
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

              <div class="flex items-center gap-4 text-sm text-slate-400">
                <span>技能匹配: {{ match.skill_match }}%</span>
                <span>評分: {{ match.rating || '-' }}</span>
              </div>

              <p v-if="match.notes" class="mt-2 text-sm text-slate-400">
                {{ match.notes }}
              </p>

              <button
                @click="selectTeacher(match)"
                class="mt-3 w-full px-4 py-2 rounded-lg bg-primary-500/20 text-primary-500 hover:bg-primary-500/30 transition-colors text-sm"
              >
                選擇這位老師
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-6 glass-card p-6">
        <h2 class="text-lg font-semibold text-white mb-4">人才庫搜尋</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
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

        <div v-if="talentResults.length > 0" class="mt-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="teacher in talentResults"
            :key="teacher.id"
            class="p-4 rounded-lg bg-white/5"
          >
            <div class="flex items-center gap-3 mb-2">
              <div class="w-8 h-8 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
                <span class="text-white text-sm">{{ teacher.name?.charAt(0) || '?' }}</span>
              </div>
              <div>
                <h4 class="text-white font-medium">{{ teacher.name }}</h4>
                <p class="text-xs text-slate-400">{{ teacher.city }}</p>
              </div>
            </div>
            <div v-if="teacher.skills?.length" class="flex flex-wrap gap-1 mt-2">
              <span
                v-for="skill in teacher.skills.slice(0, 3)"
                :key="skill"
                class="px-2 py-0.5 rounded-full text-xs bg-primary-500/20 text-primary-500"
              >
                {{ skill }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </main>

    <NotificationDropdown
      v-if="notificationUI.show.value"
      @close="notificationUI.close()"
    />
  </div>
</template>

<script setup lang="ts">
 definePageMeta({
   middleware: 'auth-admin',
   layout: 'admin',
 })

 const notificationUI = useNotification()
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
  room_id: '',
  skills: ''
})

const talentSearch = ref({
  city: '',
  skills: ''
})

const findMatches = async () => {
  if (!form.value.start_time || !form.value.end_time || !form.value.room_id) {
    alert('請填寫開始時間、結束時間和教室')
    return
  }

  searching.value = true
  hasSearched.value = true

  try {
    const api = useApi()
    const centerId = getCenterId()
    const response = await api.post<{ code: number; datas: any[] }>(`/admin/matching/teachers`, {
      center_id: parseInt(centerId),
      room_id: parseInt(form.value.room_id),
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
    const centerId = getCenterId()
    const response = await api.get<{ code: number; datas: any[] }>(
      `/admin/matching/teachers/search?city=${talentSearch.value.city}&skills=${talentSearch.value.skills}`
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
  alert(`已選擇 ${match.teacher_name}`)
}

const clearForm = () => {
  form.value = {
    start_time: '',
    end_time: '',
    room_id: '',
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
