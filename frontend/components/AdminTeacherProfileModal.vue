<template>
  <div class="fixed inset-0 z-[200] flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-md max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
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

        <!-- 證照數量 -->
        <div class="flex items-center justify-between p-3 rounded-lg bg-white/5">
          <div class="flex items-center gap-2">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
            </svg>
            <span class="text-slate-300">持有證照</span>
          </div>
          <span class="text-lg font-semibold text-slate-100">{{ teacher?.certificates?.length || 0 }} 張</span>
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
  level: string
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
  certificates?: any[]
}

const props = defineProps<{
  teacher: TeacherProfile | null
}>()

const emit = defineEmits<{
  close: []
}>()
</script>
