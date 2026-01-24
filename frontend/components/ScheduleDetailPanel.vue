<template>
  <Teleport to="body">
    <div
      v-if="showPanel"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm isolate"
      @click.self="$emit('close')"
    >
      <div class="glass-card w-full max-w-md max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
        <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
          <h3 class="text-lg font-semibold text-slate-100">
            {{ schedule ? schedule.offering_name : '選擇時段' }}
          </h3>
          <button
            @click="$emit('close')"
            class="p-2 rounded-lg hover:bg-white/10 transition-colors"
          >
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

    <div v-if="schedule" class="space-y-4">
      <div class="glass p-3 rounded-xl">
        <h4 class="text-sm font-medium text-slate-300 mb-2">課程資訊</h4>
        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-slate-400">課程</span>
            <span class="text-slate-100">{{ schedule.offering_name }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">老師</span>
            <span class="text-slate-100">{{ schedule.teacher_name }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">教室</span>
            <span class="text-slate-100">{{ schedule.room_name }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">時間</span>
            <span class="text-slate-100">{{ formatTime(props.time ?? 0) }} - {{ formatTime((props.time ?? 0) + 1) }}</span>
          </div>
        </div>
      </div>

      <div
        v-if="validation && !validation.valid"
        class="p-3 rounded-xl border"
        :class="validation.conflicts?.some((c: any) => c.type === 'TEACHER_OVERLAP')
          ? 'border-critical-500/50 bg-critical-500/10'
          : 'border-warning-500/50 bg-warning-500/10'"
      >
        <div class="flex items-start gap-2">
          <svg class="w-5 h-5 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 2.502-3.25V6.718c0-1.583 1.962-3.25 3.25H5.082c-1.54 0-2.502 1.667-2.502 3.25v8.016c0 1.583-1.962 3.25-3.25H12zM9 12h.01" />
          </svg>
          <div>
            <h4 class="font-medium text-slate-100 mb-1">
              {{ validation.conflicts?.some((c: any) => c.type === 'TEACHER_OVERLAP')
                ? '老師衝突'
                : '緩衝警告' }}
            </h4>
            <ul class="space-y-1 text-sm text-slate-400">
              <li
                v-for="(conflict, index) in validation.conflicts"
                :key="index"
                class="flex items-start gap-2"
              >
                <span class="text-primary-500">•</span>
                <span>{{ conflict.message }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <div
        v-if="validation && validation.valid"
        class="p-3 rounded-xl bg-success-500/10 border border-success-500/30"
      >
        <div class="flex items-center gap-2">
          <svg class="w-5 h-5 text-success-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <span class="text-sm font-medium text-success-500">可以排入此時段</span>
        </div>
      </div>

      <div class="space-y-3">
        <div v-if="schedule" class="flex items-center gap-2">
          <button
            @click="handleEdit"
            class="flex-1 glass-btn py-3 rounded-xl font-medium flex items-center justify-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 0L21.828 3.172a2 2 0 010-2.828l-7-7a2 2 0 00-2.828 0L2.172 20.828a2 2 0 010 2.828l7 7a2 2 0 0012.828 0l7.172-7.172z" />
            </svg>
            編輯
          </button>
          <HelpTooltip
            title="編輯排課"
            description="修改已存在的排課規則資訊。"
            :usage="['點擊開啟編輯模式', '可修改老師、教室、時間', '儲存後自動更新課表']"
          />
        </div>

        <div v-if="schedule" class="flex items-center gap-2">
          <button
            @click="handleDelete"
            class="flex-1 btn-critical py-3 rounded-xl font-medium flex items-center justify-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            刪除
          </button>
          <HelpTooltip
            title="刪除排課"
            description="移除已建立的排課規則。"
            :usage="['點擊刪除按鈕', '系統會跳出確認訊息', '確認後將永久刪除']"
          />
        </div>

        <div v-if="!schedule && validation && validation.valid" class="flex items-center gap-2">
          <button
            @click="handleCreate"
            class="flex-1 btn-primary py-3 rounded-xl font-medium flex items-center justify-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            建立課程
          </button>
          <HelpTooltip
            title="建立課程"
            description="在選定的時段建立新的排課。"
            :usage="['確認時段無衝突後點擊', '選擇課程、老師、教室', '完成後課表自動更新']"
          />
        </div>

        <div v-if="!schedule && validation && validation.valid" class="flex items-center gap-2">
          <button
            @click="handleFindSubstitute"
            class="flex-1 btn-secondary py-3 rounded-xl font-medium flex items-center justify-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            找代課老師
          </button>
          <HelpTooltip
            title="找代課老師"
            description="智慧媒合可替代的老師人選。"
            :usage="['系統會根據技能、時間評估', '顯示符合條件的老師列表', '可快速指派代課']"
          />
        </div>
      </div>
    </div>

    <div
      v-else
      class="text-center py-12 text-slate-500"
    >
      點擊左側網格查看詳情
    </div>
    </div>
  </div>
  </Teleport>
</template>

<script setup lang="ts">
import { alertConfirm } from '~/composables/useAlert'

const props = defineProps<{
  time?: number
  weekday?: number
  schedule?: any
  validation?: any
}>()

const emit = defineEmits<{
  close: []
  edit: []
  delete: []
  create: []
  findSubstitute: []
}>()

const showPanel = computed(() => !!props.schedule || (!!props.validation && props.validation.valid))

const formatTime = (hour: number): string => {
  return `${hour}:00`
}

const handleEdit = () => {
  emit('edit')
}

const handleDelete = async () => {
  if (await alertConfirm('確定要刪除此排課？')) {
    emit('delete')
  }
}

const handleCreate = () => {
  emit('create')
}

const handleFindSubstitute = () => {
  emit('findSubstitute')
}
</script>
