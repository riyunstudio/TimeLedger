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
        <div v-if="teacher?.certificates?.length > 0">
          <h5 class="text-sm font-medium text-slate-400 mb-2">持有證照 ({{ teacher.certificates.length }} 張)</h5>
          <div class="space-y-2">
            <div
              v-for="cert in teacher.certificates"
              :key="cert.id"
              class="flex items-center gap-3 p-3 rounded-lg bg-white/5 hover:bg-white/10 transition-colors"
            >
              <!-- 證照圖示 -->
              <div class="flex-shrink-0 w-10 h-10 rounded-lg bg-amber-500/20 flex items-center justify-center">
                <svg v-if="getCertIcon(cert.name) === 'pdf'" class="w-5 h-5 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
                <svg v-else class="w-5 h-5 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              
              <!-- 證照資訊 -->
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-slate-200 truncate">{{ cert.name }}</p>
                <p class="text-xs text-slate-500">
                  發照日期：{{ formatDate(cert.issued_at) }}
                </p>
              </div>
              
              <!-- 查看按鈕 -->
              <a
                v-if="cert.file_url"
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
}

const props = defineProps<{
  teacher: TeacherProfile | null
}>()

const emit = defineEmits<{
  close: []
}>()

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
</script>
