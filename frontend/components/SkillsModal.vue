<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg sm:max-w-xl max-h-[90vh] overflow-visible flex flex-col animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          技能與證照
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

        <div class="p-4 flex-1 overflow-y-auto space-y-6">
          <div>
            <div class="flex items-center justify-between mb-3">
              <h4 class="font-medium text-slate-100 text-base sm:text-lg">技能</h4>
            <button
              @click="showAddSkill = true"
              class="text-sm text-primary-500 hover:text-primary-400 font-medium"
            >
              + 新增技能
            </button>
          </div>

          <div
            v-if="skills.length === 0"
            class="text-center py-6 text-slate-500"
          >
            尚未添加技能
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="skill in skills"
              :key="skill.id"
              class="glass p-3 rounded-xl"
            >
              <div class="flex items-center justify-between">
                <div>
                  <h5 class="font-medium text-slate-100 text-sm sm:text-base">{{ skill.skill_name }}</h5>
                  <span
                    class="px-2 py-1 rounded-full text-xs font-medium mt-1 inline-block"
                    :class="getLevelClass(skill.level)"
                  >
                    {{ getLevelText(skill.level) }}
                  </span>
                </div>
                <button
                  @click="deleteSkill(skill.id)"
                  class="p-2 rounded-lg hover:bg-critical-500/20 text-critical-500 transition-colors"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

          <div>
            <div class="flex items-center justify-between mb-3">
              <h4 class="font-medium text-slate-100 text-base sm:text-lg">證照</h4>
            <button
              @click="showAddCertificate = true"
              class="text-sm text-primary-500 hover:text-primary-400 font-medium"
            >
              + 新增證照
            </button>
          </div>

          <div
            v-if="certificates.length === 0"
            class="text-center py-6 text-slate-500"
          >
            尚未添加證照
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="cert in certificates"
              :key="cert.id"
              class="glass p-3 rounded-xl"
            >
              <div class="flex items-center justify-between">
                <div>
                  <h5 class="font-medium text-slate-100 text-sm sm:text-base">{{ cert.certificate_name }}</h5>
                  <p v-if="cert.issued_by" class="text-sm text-slate-400">
                    {{ cert.issued_by }}
                  </p>
                  <p v-if="cert.issued_date" class="text-xs sm:text-sm text-slate-500">
                    發證日期: {{ formatDate(cert.issued_date) }}
                  </p>
                </div>
                <button
                  @click="deleteCertificate(cert.id)"
                  class="p-2 rounded-lg hover:bg-critical-500/20 text-critical-500 transition-colors"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <AddSkillModal
        v-if="showAddSkill"
        @close="showAddSkill = false"
        @added="fetchData"
      />

      <AddCertificateModal
        v-if="showAddCertificate"
        @close="showAddCertificate = false"
        @added="fetchData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertConfirm, alertError } from '~/composables/useAlert'

const emit = defineEmits<{
  close: []
}>()

const authStore = useAuthStore()
const loading = ref(false)
const skills = ref<any[]>([])
const certificates = ref<any[]>([])
const showAddSkill = ref(false)
const showAddCertificate = ref(false)

const fetchData = async () => {
  try {
    const api = useApi()
    const [skillsRes, certsRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>('/teacher/me/skills'),
      api.get<{ code: number; datas: any[] }>('/teacher/me/certificates')
    ])
    skills.value = skillsRes.datas || []
    certificates.value = certsRes.datas || []
  } catch (error) {
    console.error('Failed to fetch skills and certificates:', error)
  }
}

const deleteSkill = async (id: number) => {
  if (!await alertConfirm('確定要刪除此技能？')) return

  try {
    const api = useApi()
    await api.delete(`/teacher/me/skills/${id}`)
    await fetchData()
  } catch (error) {
    console.error('Failed to delete skill:', error)
    await alertError('刪除失敗')
  }
}

const deleteCertificate = async (id: number) => {
  if (!await alertConfirm('確定要刪除此證照？')) return

  try {
    const api = useApi()
    await api.delete(`/teacher/me/certificates/${id}`)
    await fetchData()
  } catch (error) {
    console.error('Failed to delete certificate:', error)
    await alertError('刪除失敗')
  }
}

const getLevelClass = (level: string): string => {
  switch (level) {
    case 'Beginner':
      return 'bg-slate-500/20 text-slate-400'
    case 'Intermediate':
      return 'bg-primary-500/20 text-primary-500'
    case 'Advanced':
      return 'bg-secondary-500/20 text-secondary-500'
    case 'Expert':
      return 'bg-warning-500/20 text-warning-500'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}

const getLevelText = (level: string): string => {
  switch (level) {
    case 'Beginner':
      return '初級'
    case 'Intermediate':
      return '中級'
    case 'Advanced':
      return '高級'
    case 'Expert':
      return '專家'
    default:
      return level
  }
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

onMounted(() => {
  fetchData()
})
</script>
