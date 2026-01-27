<template>
  <div
    :class="[
      'talent-card relative p-4 rounded-xl border transition-all',
      selected
        ? 'bg-indigo-500/10 border-indigo-500/30'
        : 'bg-white/5 border-white/10 hover:bg-white/10 hover:border-white/20'
    ]"
  >
    <!-- 選取核取方塊 -->
    <div
      v-if="showCheckbox && teacher.is_open_to_hiring"
      class="absolute top-3 right-3 z-10"
    >
      <input
        type="checkbox"
        :checked="selected"
        @change="$emit('update:selected', !selected)"
        class="w-5 h-5 rounded border-white/20 bg-white/10 text-indigo-500 focus:ring-indigo-500 focus:ring-offset-0"
      />
    </div>

    <!-- 老師基本資訊 -->
    <div class="flex items-start gap-3 mb-3">
      <div class="w-12 h-12 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0">
        <span class="text-white font-medium">{{ teacher.name?.charAt(0) || '?' }}</span>
      </div>
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 mb-1">
          <h4 class="text-white font-medium truncate">{{ teacher.name }}</h4>
          <BaseBadge
            v-if="teacher.is_member"
            variant="success"
            size="xs"
          >
            中心成員
          </BaseBadge>
          <BaseBadge
            v-if="!teacher.is_open_to_hiring"
            variant="secondary"
            size="xs"
          >
            已關閉徵才
          </BaseBadge>
        </div>
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
    <p v-if="teacher.bio" class="text-sm text-slate-400 mb-3 line-clamp-2">
      {{ teacher.bio }}
    </p>

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
    <div v-if="teacher.personal_hashtags?.length" class="flex flex-wrap gap-1 mb-4">
      <span
        v-for="(tag, index) in teacher.personal_hashtags.slice(0, 5)"
        :key="index"
        class="px-2 py-0.5 rounded-full text-xs bg-primary-500/20 text-primary-400"
      >
        {{ tag }}
      </span>
    </div>

    <!-- 邀請狀態 -->
    <div v-if="invitationStatus" class="mb-3">
      <BaseBadge :variant="invitationStatus.variant" size="sm">
        {{ invitationStatus.text }}
      </BaseBadge>
    </div>

    <!-- 操作按鈕 -->
    <div class="flex gap-2">
      <button
        v-if="teacher.is_open_to_hiring"
        @click="$emit('invite', teacher)"
        :disabled="inviteLoading || invitationStatus?.sent"
        class="flex-1 px-3 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors text-sm flex items-center justify-center gap-1 disabled:opacity-50"
      >
        <svg v-if="inviteLoading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>
        {{ invitationStatus?.sent ? '已邀請' : '邀請合作' }}
      </button>
      <button
        v-else
        disabled
        class="flex-1 px-3 py-2 rounded-lg bg-white/5 text-slate-500 text-sm flex items-center justify-center gap-1 cursor-not-allowed"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
        </svg>
        已關閉徵才
      </button>
      <button
        @click="$emit('view', teacher)"
        class="px-3 py-2 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white transition-colors"
        title="查看詳細資訊"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
      </button>
      <button
        @click="$emit('compare', teacher)"
        class="px-3 py-2 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white transition-colors"
        title="加入比較"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      </button>
    </div>

    <!-- 聯絡資訊（邀請成功後顯示） -->
    <div
      v-if="showContactInfo && teacher.public_contact_info"
      class="mt-3 pt-3 border-t border-white/10"
    >
      <p class="text-xs text-slate-400">
        <span class="font-medium text-slate-300">聯絡方式：</span>
        {{ teacher.public_contact_info }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { SKILL_CATEGORIES } from '~/types'

interface Teacher {
  id: number
  name: string
  bio?: string
  city?: string
  district?: string
  skills?: Array<{ name: string; category: string }>
  personal_hashtags?: string[]
  is_open_to_hiring: boolean
  is_member: boolean
  public_contact_info?: string
}

interface InvitationStatus {
  sent: boolean
  variant: 'success' | 'warning' | 'error' | 'secondary'
  text: string
}

const props = defineProps<{
  teacher: Teacher
  selected?: boolean
  showCheckbox?: boolean
  inviteLoading?: boolean
  invitationStatus?: InvitationStatus | null
  showContactInfo?: boolean
}>()

defineEmits<{
  invite: [teacher: Teacher]
  view: [teacher: Teacher]
  compare: [teacher: Teacher]
  'update:selected': [selected: boolean]
}>()

// 技能類別相關函數
const getSkillCategoryIcon = (category: string): string => {
  return SKILL_CATEGORIES[category as keyof typeof SKILL_CATEGORIES]?.icon || '✨'
}

const getSkillCategoryStyle = (category: string): string => {
  return SKILL_CATEGORIES[category as keyof typeof SKILL_CATEGORIES]?.color || 'bg-slate-500/20 text-slate-400 border-slate-500/30'
}
</script>
