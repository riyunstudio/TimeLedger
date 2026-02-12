<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="close"
      >
        <div class="glass-card w-full max-w-2xl max-h-[90vh] overflow-hidden flex flex-col">
          <!-- 標題 -->
          <div class="flex items-center justify-between p-4 border-b border-white/10">
            <h3 class="text-lg font-semibold text-white">人才詳細資訊</h3>
            <button
              @click="close"
              class="p-2 rounded-lg hover:bg-white/10 text-slate-400 hover:text-white transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- 內容 -->
          <div class="flex-1 overflow-y-auto p-4">
            <!-- 基本資訊 -->
            <div class="flex items-start gap-4 mb-6">
              <div class="w-16 h-16 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0">
                <span class="text-white text-xl font-medium">{{ teacher?.name?.charAt(0) || '?' }}</span>
              </div>
              <div class="flex-1">
                <h4 class="text-xl font-semibold text-white mb-1">{{ teacher?.name }}</h4>
                <div class="flex flex-wrap items-center gap-3 text-sm text-slate-400">
                  <span v-if="teacher?.city" class="flex items-center gap-1">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    {{ teacher?.city }}{{ teacher?.district ? ` ${teacher.district}` : '' }}
                  </span>
                  <span v-if="teacher?.rating" class="flex items-center gap-1">
                    <svg class="w-4 h-4 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                    {{ teacher.rating.toFixed(1) }}
                  </span>
                  <span v-if="teacher?.is_member" class="px-2 py-0.5 rounded bg-green-500/20 text-green-400 text-xs">
                    已加入中心
                  </span>
                  <span v-else class="px-2 py-0.5 rounded bg-slate-500/20 text-slate-400 text-xs">
                    未加入中心
                  </span>
                </div>
              </div>
            </div>

            <!-- 個人簡介 -->
            <div v-if="teacher?.bio" class="mb-6">
              <h5 class="text-sm font-medium text-slate-300 mb-2">個人簡介</h5>
              <p class="text-slate-400 text-sm leading-relaxed">{{ teacher.bio }}</p>
            </div>

            <!-- 技能列表 -->
            <div v-if="teacher?.skills?.length" class="mb-6">
              <h5 class="text-sm font-medium text-slate-300 mb-2">專業技能</h5>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="(skill, index) in teacher.skills"
                  :key="index"
                  class="px-3 py-1.5 rounded-full bg-primary-500/20 text-primary-400 text-sm"
                >
                  {{ typeof skill === 'string' ? skill : skill.name }}
                </span>
              </div>
            </div>

            <!-- 證照列表 -->
            <div v-if="teacher?.certificates?.length" class="mb-6">
              <h5 class="text-sm font-medium text-slate-300 mb-2">專業證照</h5>
              <div class="space-y-2">
                <div
                  v-for="cert in teacher.certificates"
                  :key="cert.id"
                  class="flex items-center gap-3 p-3 rounded-lg bg-white/5 border border-white/10"
                >
                  <div class="w-10 h-10 rounded-lg bg-yellow-500/20 flex items-center justify-center shrink-0">
                    <svg class="w-5 h-5 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
                    </svg>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-white text-sm font-medium truncate">{{ cert.name }}</p>
                    <p v-if="cert.issuer" class="text-slate-500 text-xs mt-0.5">{{ cert.issuer }}</p>
                  </div>
                  <div v-if="cert.expiry_date" class="shrink-0">
                    <span class="px-2 py-0.5 rounded bg-slate-500/20 text-slate-400 text-xs">
                      到期日：{{ cert.expiry_date }}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 個人標籤 -->
            <div v-if="teacher?.personal_hashtags?.length" class="mb-6">
              <h5 class="text-sm font-medium text-slate-300 mb-2">個人標籤</h5>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="(tag, index) in teacher.personal_hashtags"
                  :key="index"
                  class="px-3 py-1 rounded-full bg-white/10 text-slate-300 text-sm"
                >
                  #{{ typeof tag === 'string' ? tag : tag.name }}
                </span>
              </div>
            </div>

            <!-- 公開聯絡資訊 -->
            <div v-if="teacher?.public_contact_info" class="mb-6">
              <h5 class="text-sm font-medium text-slate-300 mb-2">聯絡方式</h5>
              <div class="p-3 rounded-lg bg-white/5 border border-white/10">
                <p class="text-slate-400 text-sm whitespace-pre-wrap">{{ teacher.public_contact_info }}</p>
              </div>
            </div>

            <!-- 徵才狀態 -->
            <div class="p-4 rounded-lg" :class="[
              teacher?.is_open_to_hiring
                ? 'bg-green-500/10 border border-green-500/20'
                : 'bg-slate-500/10 border border-slate-500/20'
            ]">
              <div class="flex items-center gap-2">
                <svg v-if="teacher?.is_open_to_hiring" class="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <svg v-else class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                </svg>
                <span :class="teacher?.is_open_to_hiring ? 'text-green-400' : 'text-slate-400'">
                  {{ teacher?.is_open_to_hiring ? '開放應徵中' : '目前未開放應徵' }}
                </span>
              </div>
              <p v-if="teacher?.is_member" class="mt-2 text-sm text-slate-400">
                此人才已是中心成員，無法再次邀請
              </p>
            </div>
          </div>

          <!-- 底部按鈕 -->
          <div class="flex items-center justify-end gap-3 p-4 border-t border-white/10">
            <button
              @click="close"
              class="px-4 py-2 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white transition-colors"
            >
              關閉
            </button>
            <button
              v-if="canInvite"
              @click="onInvite"
              class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors"
            >
              邀請合作
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
interface TeacherCertificate {
  id: number
  name: string
  issuer?: string
  obtained_at?: string
  expiry_date?: string
}

interface TeacherDetail {
  id: number
  name: string
  bio?: string
  city?: string
  district?: string
  rating?: number
  skills?: Array<{ name: string; category: string }>
  personal_hashtags?: Array<{ name: string } | string>
  is_open_to_hiring: boolean
  is_member: boolean
  public_contact_info?: string
  certificates?: TeacherCertificate[]
}

const props = defineProps<{
  teacher: TeacherDetail | null
  visible: boolean
}>()

const emit = defineEmits<{
  close: []
  invite: [teacher: TeacherDetail]
}>()

const close = () => {
  emit('close')
}

const onInvite = () => {
  if (props.teacher) {
    emit('invite', props.teacher)
    close()
  }
}

// 檢查是否可以邀請（開放應徵 且 尚未加入中心）
const canInvite = computed(() => {
  return props.teacher?.is_open_to_hiring && !props.teacher?.is_member
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active > div:last-child,
.modal-leave-active > div:last-child {
  transition: transform 0.3s ease;
}

.modal-enter-from > div:last-child,
.modal-leave-to > div:last-child {
  transform: scale(0.95);
}
</style>
