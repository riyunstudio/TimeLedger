<template>
  <div class="fixed inset-0 z-[200] flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-lg max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10">
        <h3 class="text-lg font-semibold text-slate-100">
          老師檔案
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-4 space-y-4">
        <!-- 頭像和姓名 -->
        <div class="flex items-center gap-4">
          <div class="w-16 h-16 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
            <span class="text-2xl font-bold text-white">{{ teacher?.name?.charAt(0) || '?' }}</span>
          </div>
          <div>
            <h4 class="text-xl font-semibold text-slate-100">{{ teacher?.name }}</h4>
            <span 
              class="px-2 py-1 rounded-full text-xs font-medium"
              :class="teacher?.is_active ? 'bg-success-500/20 text-success-500' : 'bg-warning-500/20 text-warning-500'"
            >
              {{ teacher?.is_active ? '活躍中' : '未活躍' }}
            </span>
          </div>
        </div>

        <!-- 聯繫資訊 -->
        <div class="space-y-3">
          <div class="flex items-center gap-3 p-3 rounded-lg bg-white/5">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            <span class="text-slate-300">{{ teacher?.email || '未設定' }}</span>
          </div>

          <div class="flex items-center gap-3 p-3 rounded-lg bg-white/5">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
            </svg>
            <span class="text-slate-300">{{ teacher?.phone || '未設定' }}</span>
          </div>

          <div class="flex items-center gap-3 p-3 rounded-lg bg-white/5">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span class="text-slate-300">{{ teacher?.city || '未設定' }} {{ teacher?.district || '' }}</span>
          </div>
        </div>

        <!-- 技能標籤 -->
        <div v-if="teacher?.skills?.length > 0">
          <h5 class="text-sm font-medium text-slate-400 mb-2">技能專長</h5>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="skill in teacher.skills"
              :key="skill.id"
              class="px-3 py-1.5 rounded-full text-sm font-medium bg-primary-500/20 text-primary-500"
            >
              {{ skill.skill_name }}
            </span>
          </div>
        </div>

        <!-- 證照列表 -->
        <div v-if="visibleCertificates.length > 0">
          <h5 class="text-sm font-medium text-slate-400 mb-2">持有證照 ({{ visibleCertificates.length }} 張)</h5>
          <div class="space-y-2">
            <div
              v-for="cert in visibleCertificates"
              :key="cert.id"
              class="flex items-center gap-3 p-3 rounded-lg bg-white/5 hover:bg-white/10 transition-colors"
            >
              <!-- 證照圖示 -->
              <div
                class="flex-shrink-0 w-10 h-10 rounded-lg flex items-center justify-center"
                :class="isRestrictedCert(cert) ? 'bg-slate-500/20' : 'bg-amber-500/20'"
              >
                <svg v-if="getCertIcon(cert.name) === 'pdf'" class="w-5 h-5" :class="isRestrictedCert(cert) ? 'text-slate-400' : 'text-amber-500'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
                <svg v-else class="w-5 h-5" :class="isRestrictedCert(cert) ? 'text-slate-400' : 'text-amber-500'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>

              <!-- 證照資訊 -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2">
                  <p class="text-sm font-medium text-slate-200 truncate">{{ cert.name }}</p>
                  <span v-if="isRestrictedCert(cert)" class="flex-shrink-0 flex items-center gap-1 px-1.5 py-0.5 rounded text-xs bg-slate-500/20 text-slate-400">
                    <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                    僅限名稱
                  </span>
                </div>
                <p class="text-xs text-slate-500">
                  發照日期：{{ formatDate(cert.issued_at) }}
                </p>
              </div>

              <!-- 查看按鈕 -->
              <a
                v-if="canViewCert(cert)"
                :href="cert.file_url"
                target="_blank"
                class="flex-shrink-0 p-2 rounded-lg bg-primary-500/20 text-primary-500 hover:bg-primary-500/30 transition-colors"
                title="查看證照"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
              </a>
            </div>
          </div>
        </div>

        <!-- 無證照提示 -->
        <div v-else class="flex items-center gap-3 p-3 rounded-lg bg-white/5">
          <svg class="w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
          </svg>
          <span class="text-slate-400">尚未上傳任何證照</span>
        </div>

        <!-- 評分與備註區塊 -->
        <div v-if="teacherNote" class="border-t border-white/10 pt-4 mt-4">
          <h5 class="text-sm font-medium text-slate-400 mb-3">評分與備註</h5>

          <!-- 評分顯示 -->
          <div class="flex items-center gap-3 mb-3">
            <div class="flex items-center gap-1">
              <template v-for="star in 5" :key="star">
                <svg
                  class="w-5 h-5"
                  :class="star <= (teacherNote.rating || 0) ? 'text-warning-500' : 'text-slate-600'"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              </template>
            </div>
            <span class="text-sm text-slate-300">
              {{ teacherNote.rating || 0 }} / 5
              <span v-if="teacherNote.rating > 0" class="text-warning-400 ml-1">
                ({{ ratingLabels[teacherNote.rating] }})
              </span>
            </span>
          </div>

          <!-- 內部備註 -->
          <div v-if="teacherNote.internal_note" class="p-3 rounded-lg bg-primary-500/10 border border-primary-500/20">
            <p class="text-sm text-slate-300 whitespace-pre-wrap">{{ teacherNote.internal_note }}</p>
          </div>

          <!-- 無備註提示 -->
          <div v-else class="text-sm text-slate-500 italic">
            尚未設定內部備註
          </div>
        </div>

        <!-- 載入中狀態 -->
        <div v-else-if="loadingNote" class="border-t border-white/10 pt-4 mt-4">
          <div class="animate-pulse flex items-center gap-3">
            <div class="flex gap-1">
              <div v-for="i in 5" :key="i" class="w-5 h-5 bg-white/10 rounded"></div>
            </div>
            <div class="h-4 w-20 bg-white/10 rounded"></div>
          </div>
        </div>
      </div>

      <div class="p-4 border-t border-white/10">
        <button
          @click="emit('close')"
          class="w-full glass-btn py-3 rounded-xl font-medium"
        >
          關閉
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface TeacherSkill {
  id: number
  skill_name: string
  category?: string
  level?: string
}

interface Certificate {
  id: number
  name: string
  file_url?: string
  issued_at: string
  created_at: string
  visibility?: number
}

// 證照隱私設定常數
const CERT_VISIBILITY = {
  NAME_ONLY: 1, // 僅限名稱
}

interface TeacherNote {
  id?: number
  rating: number
  internal_note: string
}

interface TeacherProfile {
  id: number
  name: string
  email: string
  phone?: string
  city?: string
  district?: string
  is_active: boolean
  skills?: TeacherSkill[]
  certificates?: Certificate[]
  note?: TeacherNote
}

const props = defineProps<{
  teacher: TeacherProfile | null
}>()

const emit = defineEmits<{
  close: []
}>()

const teacherNote = ref<TeacherNote | null>(null)
const loadingNote = ref(false)

// 取得證照圖示類型
const getCertIcon = (name: string): string => {
  const lowerName = name.toLowerCase()
  if (lowerName.includes('pdf')) return 'pdf'
  return 'image'
}

// 格式化日期
const formatDate = (dateStr: string): string => {
  if (!dateStr) return '未設定'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 獲取老師評分資料
const fetchTeacherNote = async () => {
  if (!props.teacher?.id) return

  loadingNote.value = true
  try {
    const api = useApi()
    const response = await api.get<TeacherNote>(
      `/admin/teachers/${props.teacher.id}/note`
    )
    if (response) {
      teacherNote.value = response
    }
  } catch (err) {
    console.error('Failed to fetch teacher note:', err)
  } finally {
    loadingNote.value = false
  }
}

// 監聽 teacher 變化，獲取評分資料
watch(() => props.teacher, async (newTeacher) => {
  if (newTeacher?.id) {
    await fetchTeacherNote()
  }
}, { immediate: true })

// 評分標籤文字
const ratingLabels: Record<number, string> = {
  0: '未評分',
  1: '需改進',
  2: '一般',
  3: '良好',
  4: '優良',
  5: '優秀'
}

// 判斷證照是否為受限狀態（僅限名稱或無檔案）
const isRestrictedCert = (cert: Certificate): boolean => {
  return cert.visibility === CERT_VISIBILITY.NAME_ONLY || !cert.file_url
}

// 判斷是否可以查看證照圖片
const canViewCert = (cert: Certificate): boolean => {
  return !!cert.file_url && cert.visibility !== CERT_VISIBILITY.NAME_ONLY
}

// 取得所有可見的證照列表（後端已過濾掉完全私密的）
const visibleCertificates = computed(() => {
  return props.teacher?.certificates || []
})
</script>
