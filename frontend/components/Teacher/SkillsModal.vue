<template>
  <div class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg sm:max-w-xl max-h-[90vh] flex flex-col animate-spring" @click.stop>
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
        <!-- 技能區塊 -->
        <div>
          <div class="flex items-center justify-between mb-3">
            <h4 class="font-medium text-slate-100 text-base sm:text-lg">技能</h4>
            <button
              @click="openAddSkill"
              class="text-sm text-primary-500 hover:text-primary-400 font-medium flex items-center gap-1"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              新增技能
            </button>
          </div>

          <div
            v-if="skills.length === 0"
            class="text-center py-8 text-slate-500"
          >
            <svg class="w-12 h-12 mx-auto mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
            </svg>
            尚未添加技能
          </div>

          <div v-else class="grid gap-2">
            <div
              v-for="skill in skills"
              :key="skill.id"
              class="glass p-3 rounded-xl hover:bg-white/5 transition-colors cursor-pointer group"
              @click="openEditSkill(skill)"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <!-- 類別標籤 -->
                  <span
                    class="px-2 py-1 rounded-lg text-xs font-medium border"
                    :class="getCategoryStyle(skill.category)"
                  >
                    {{ getCategoryIcon(skill.category) }} {{ getCategoryLabel(skill.category) }}
                  </span>
                  <div>
                    <h5 class="font-medium text-slate-100 text-sm sm:text-base">{{ skill.skill_name }}</h5>
                    <p v-if="skill.hashtags && skill.hashtags.length > 0" class="text-xs text-slate-500 mt-1">
                      {{ skill.hashtags.filter((t: any) => t && typeof t === 'string').map((t: string) => t.startsWith('#') ? t : '#' + t).join(' ') }}
                    </p>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <!-- 編輯提示 -->
                  <span class="text-xs text-slate-600 group-hover:text-slate-400 transition-colors">
                    點擊編輯
                  </span>
                  <button
                    @click.stop="deleteSkill(skill.id)"
                    class="p-2 rounded-lg hover:bg-critical-500/20 text-critical-500 transition-colors opacity-60 hover:opacity-100"
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

        <!-- 證照區塊 -->
        <div>
          <div class="flex items-center justify-between mb-3">
            <h4 class="font-medium text-slate-100 text-base sm:text-lg">證照</h4>
            <button
              @click="showAddCertificate = true"
              class="text-sm text-primary-500 hover:text-primary-400 font-medium flex items-center gap-1"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              新增證照
            </button>
          </div>

          <div
            v-if="certificates.length === 0"
            class="text-center py-8 text-slate-500"
          >
            <svg class="w-12 h-12 mx-auto mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
            </svg>
            尚未添加證照
          </div>

          <div v-else class="grid grid-cols-2 gap-3">
            <div
              v-for="cert in certificates"
              :key="cert.id"
              class="glass p-3 rounded-xl"
            >
              <!-- 證照圖片預覽 -->
              <div
                v-if="cert.file_url"
                class="aspect-video rounded-lg bg-slate-800 mb-2 overflow-hidden cursor-pointer hover:ring-2 ring-primary-500 transition-all"
                :class="{ 'opacity-75 cursor-not-allowed': !cert.file_url }"
                @click="cert.file_url ? openCertificatePreview(cert) : null"
              >
                <img
                  v-if="cert.file_url"
                  :src="cert.file_url"
                  :alt="cert.name"
                  class="w-full h-full object-cover"
                  @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
              </div>
              <!-- 無圖片時的佔位符 -->
              <div
                v-else
                class="aspect-video rounded-lg bg-slate-800 mb-2 flex items-center justify-center"
              >
                <svg class="w-8 h-8 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <div class="flex items-start justify-between">
                <div class="min-w-0 flex-1">
                  <h5 class="font-medium text-slate-100 text-sm truncate">{{ cert.name }}</h5>
                  <p v-if="cert.issued_at" class="text-xs text-slate-500 mt-1">
                    {{ formatDate(cert.issued_at) }}
                  </p>
                  <!-- 隱私狀態顯示 -->
                  <div class="flex items-center gap-1 mt-1">
                    <span
                      class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-xs font-medium"
                      :class="getVisibilityClass(cert.visibility)"
                    >
                      <svg v-if="cert.visibility === 2" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <svg v-else-if="cert.visibility === 1" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                      </svg>
                      <svg v-else class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                      </svg>
                      {{ getVisibilityLabel(cert.visibility) }}
                    </span>
                  </div>
                </div>
                <button
                  @click="deleteCertificate(cert.id)"
                  class="p-1.5 rounded-lg hover:bg-critical-500/20 text-critical-500 transition-colors ml-2 flex-shrink-0"
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
    </div>

    <!-- 新增/編輯技能 Modal - AddSkillModal 自己管理 isOpen 狀態，不需要 Teleport -->
    <AddSkillModal
      v-if="showAddSkill || editingSkill"
      :skill="editingSkill"
      :existing-skills="skills"
      @close="closeSkillModal"
      @added="fetchData"
      @updated="fetchData"
    />

    <!-- 新增證照 Modal -->
    <AddCertificateModal
      v-if="showAddCertificate"
      @close="showAddCertificate = false"
      @added="fetchData"
    />

    <!-- 證照圖片預覽 Modal -->
    <div
      v-if="previewingCertificate"
      class="fixed inset-0 z-[200] flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm"
      @click="previewingCertificate = null"
    >
      <div class="relative max-w-3xl w-full" @click.stop>
        <button
          @click="previewingCertificate = null"
          class="absolute -top-10 right-0 text-white hover:text-primary-400 transition-colors"
        >
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <img
          v-if="previewingCertificate.file_url"
          :src="previewingCertificate.file_url"
          :alt="previewingCertificate.name"
          class="w-full h-auto rounded-lg shadow-2xl"
        />
        <div class="mt-4 text-center">
          <h4 class="text-lg font-semibold text-white">{{ previewingCertificate.name }}</h4>
          <p v-if="previewingCertificate.issued_at" class="text-sm text-slate-500">
            發證日期：{{ formatDate(previewingCertificate.issued_at) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertConfirm, alertError, alertWarning } from '~/composables/useAlert'
import { SKILL_CATEGORIES } from '~/types'
import { CERTIFICATE_VISIBILITY_LABELS, type CertificateVisibility } from '~/types/teacher'
import AddSkillModal from './AddSkillModal.vue'
import AddCertificateModal from './AddCertificateModal.vue'

const emit = defineEmits<{
  close: []
}>()

const authStore = useAuthStore()
const loading = ref(false)
const skills = ref<any[]>([])
const certificates = ref<any[]>([])
const showAddSkill = ref(false)
const showAddCertificate = ref(false)
const editingSkill = ref<any>(null)
const previewingCertificate = ref<any>(null)

const getCategoryLabel = (category: string): string => {
  return SKILL_CATEGORIES[category as keyof typeof SKILL_CATEGORIES]?.label || category
}

const getCategoryIcon = (category: string): string => {
  return SKILL_CATEGORIES[category as keyof typeof SKILL_CATEGORIES]?.icon || '✨'
}

const getCategoryStyle = (category: string): string => {
  return SKILL_CATEGORIES[category as keyof typeof SKILL_CATEGORIES]?.color || 'bg-slate-500/20 text-slate-400 border-slate-500/30'
}

/**
 * 取得隱私狀態顯示文字
 */
const getVisibilityLabel = (visibility?: number): string => {
  return CERTIFICATE_VISIBILITY_LABELS[visibility as keyof typeof CERTIFICATE_VISIBILITY_LABELS] || '私密'
}

/**
 * 取得隱私狀態樣式
 */
const getVisibilityClass = (visibility?: number): string => {
  switch (visibility) {
    case 2: // 公開
      return 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30'
    case 1: // 僅限名稱
      return 'bg-amber-500/20 text-amber-400 border-amber-500/30'
    case 0: // 私密
    default:
      return 'bg-slate-500/20 text-slate-400 border-slate-500/30'
  }
}

const fetchData = async () => {
  try {
    const api = useApi()
    const [skillsRes, certsRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>('/teacher/me/skills'),
      api.get<{ code: number; datas: any[] }>('/teacher/me/certificates')
    ])
    // 轉換技能資料，確保標籤是字串
    skills.value = (skillsRes.datas || []).map((skill: any) => ({
      ...skill,
      hashtags: (skill.hashtags || []).map((h: any) => {
        // 從各種可能的結構中提取標籤名稱
        if (typeof h === 'string') return h
        // 處理 { hashtag: { name: "xxx" } } 結構
        if (h.hashtag && typeof h.hashtag === 'object') {
          return h.hashtag.name || ''
        }
        // 處理 { name: "xxx" } 結構
        if (h.name) return h.name
        // 如果都沒有，回傳空字串
        return ''
      }).filter((t: string) => typeof t === 'string' && t.length > 0) // 只保留非空字串
    }))
    // 轉換證照資料，確保欄位與 API 回應一致
    certificates.value = (certsRes.datas || []).map((cert: any) => ({
      ...cert,
      // 確保使用 API 回應的欄位名稱
      name: cert.name,
      issued_at: cert.issued_at,
      file_url: cert.file_url,
    }))
  } catch (error) {
    console.error('Failed to fetch skills and certificates:', error)
  }
}

const openAddSkill = () => {
  editingSkill.value = null
  showAddSkill.value = true
}

const openEditSkill = (skill: any) => {
  editingSkill.value = skill
  showAddSkill.value = true
}

const closeSkillModal = () => {
  showAddSkill.value = false
  editingSkill.value = null
}

const openCertificatePreview = (cert: any) => {
  previewingCertificate.value = cert
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
