<template>
  <div class="text-center mb-8">
    <div
      class="w-24 h-24 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center mx-auto mb-4"
    >
      <span class="text-4xl font-bold text-white">
        {{ authStore.user?.name?.charAt(0) || 'T' }}
      </span>
    </div>
    <h1 class="text-2xl font-bold text-slate-100 mb-1">
      {{ authStore.user?.name }}
    </h1>
    <p v-if="authStore.user?.bio" class="text-slate-400">
      {{ authStore.user?.bio }}
    </p>

    <!-- 檔案完整度指標 -->
    <div class="mt-4 max-w-xs mx-auto">
      <div class="flex items-center justify-between text-xs text-slate-400 mb-1">
        <span>檔案完整度</span>
        <span class="text-white font-medium">{{ profileCompleteness }}%</span>
      </div>
      <div class="h-2 bg-white/10 rounded-full overflow-hidden">
        <div
          class="h-full rounded-full transition-all duration-500"
          :class="getCompletenessColor()"
          :style="{ width: `${profileCompleteness}%` }"
        ></div>
      </div>
      <p class="text-xs text-slate-500 mt-1">{{ getCompletenessHint() }}</p>
    </div>
  </div>

  <div class="space-y-4">
    <button
      @click="showProfileModal = true"
      class="glass-card w-full p-4 flex items-center gap-4 hover:bg-white/5 transition-colors"
    >
      <div class="p-3 rounded-xl bg-primary-500/20">
        <svg class="w-6 h-6 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
        </svg>
      </div>
      <div class="flex-1 text-left">
        <h3 class="font-medium text-slate-100">個人檔案</h3>
        <p class="text-sm text-slate-400">編輯基本資料</p>
      </div>
      <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
    </button>

    <button
      @click="showHiringModal = true"
      class="glass-card w-full p-4 flex items-center gap-4 hover:bg-white/5 transition-colors"
    >
      <div class="p-3 rounded-xl bg-secondary-500/20">
        <svg class="w-6 h-6 text-secondary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
      </div>
      <div class="flex-1 text-left">
        <h3 class="font-medium text-slate-100">求職設定</h3>
        <p class="text-sm text-slate-400">開放媒合狀態</p>
      </div>
      <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
    </button>

    <button
      @click="showSkillsModal = true"
      class="glass-card w-full p-4 flex items-center gap-4 hover:bg-white/5 transition-colors"
    >
      <div class="p-3 rounded-xl bg-success-500/20">
        <svg class="w-6 h-6 text-success-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
        </svg>
      </div>
      <div class="flex-1 text-left">
        <h3 class="font-medium text-slate-100">技能與證照</h3>
        <p class="text-sm text-slate-400">管理教學專長</p>
      </div>
      <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
    </button>

    <button
      @click="showExportModal = true"
      class="glass-card w-full p-4 flex items-center gap-4 hover:bg-white/5 transition-colors"
    >
      <div class="p-3 rounded-xl bg-warning-500/20">
        <svg class="w-6 h-6 text-warning-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
      </div>
      <div class="flex-1 text-left">
        <h3 class="font-medium text-slate-100">匯出課表</h3>
        <p class="text-sm text-slate-400">精美圖片分享</p>
      </div>
      <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
    </button>
  </div>

  <div class="mt-8">
    <h2 class="text-lg font-semibold text-slate-100 mb-4">我的中心</h2>
    <div
      v-if="teacherStore.centers.length > 0"
      class="space-y-3"
    >
      <div
        v-for="membership in teacherStore.centers"
        :key="membership.id"
        class="glass-card p-4"
      >
        <div class="flex items-center justify-between">
          <div>
            <h3 class="font-medium text-slate-100">
              {{ membership.center_name || '未知中心' }}
            </h3>
          </div>
          <span
            class="px-3 py-1 rounded-full text-sm font-medium"
            :class="getMembershipStatusClass(membership.status)"
          >
            {{ getMembershipStatusText(membership.status) }}
          </span>
        </div>
      </div>
    </div>

    <div
      v-else
      class="text-center py-8 text-slate-500"
    >
      尚未加入任何中心
    </div>
  </div>

  <ProfileModal
    v-if="showProfileModal"
    @close="showProfileModal = false"
  />
  <HiringModal
    v-if="showHiringModal"
    @close="showHiringModal = false"
  />
  <SkillsModal
    v-if="showSkillsModal"
    @close="showSkillsModal = false"
  />
  <ExportModal
    v-if="showExportModal"
    @close="showExportModal = false"
  />

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />

  <TeacherSidebar
    v-if="sidebarStore.isOpen.value"
    @close="sidebarStore.close()"
  />
</template>

<script setup lang="ts">
 definePageMeta({
   middleware: 'auth-teacher',
   layout: 'default',
 })

 const authStore = useAuthStore()
 const teacherStore = useTeacherStore()
 const sidebarStore = useSidebar()
 const notificationUI = useNotification()

 const showProfileModal = ref(false)
 const showHiringModal = ref(false)
 const showSkillsModal = ref(false)
 const showExportModal = ref(false)

 // 計算檔案完整度
 const profileCompleteness = computed(() => {
   let score = 0
   const user = authStore.user
   const teacher = teacherStore.centers.length > 0 ? teacherStore.centers[0] : null

   // 基本資料（30分）
   if (user?.name) score += 10
   if (user?.email) score += 10
   if (user?.phone) score += 10

   // 簡介（20分）
   if (user?.bio && user.bio.length >= 10) score += 20

   // 技能（30分）
   if (teacher?.skills && teacher.skills.length > 0) {
     score += Math.min(teacher.skills.length * 10, 30)
   }

   // 證照（20分）
   if (teacher?.certificates && teacher.certificates.length > 0) {
     score += Math.min(teacher.certificates.length * 5, 20)
   }

   return Math.min(score, 100)
 })

 const getCompletenessColor = () => {
   const score = profileCompleteness.value
   if (score >= 80) return 'bg-success-500'
   if (score >= 50) return 'bg-warning-500'
   return 'bg-critical-500'
 }

 const getCompletenessHint = () => {
   const score = profileCompleteness.value
   if (score >= 100) return '太棒了！您的檔案已經完整'
   if (score >= 80) return '很不錯！再完善一些就能獲得更好的曝光'
   if (score >= 50) return '請補完基本資料和技能，讓更多人認識您'
   return '建議您先填寫基本資料和簡介'
 }

 const getMembershipStatusClass = (status: string): string => {
   switch (status) {
     case 'ACTIVE':
       return 'bg-success-500/20 text-success-500'
     case 'INVITED':
       return 'bg-warning-500/20 text-warning-500'
     case 'INACTIVE':
       return 'bg-slate-500/20 text-slate-400'
     default:
       return 'bg-slate-500/20 text-slate-400'
   }
 }

 const getMembershipStatusText = (status: string): string => {
   switch (status) {
     case 'ACTIVE':
       return '已加入'
     case 'INVITED':
       return '邀請中'
     case 'INACTIVE':
       return '已離開'
     default:
       return status
   }
 }

onMounted(() => {
  teacherStore.fetchCenters()
})
</script>
