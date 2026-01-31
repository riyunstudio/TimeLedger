<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
      @click.self="$emit('close')"
    >
      <div class="glass-card w-full max-w-sm">
        <div class="p-4 border-b border-white/10">
          <h3 class="text-lg font-semibold text-white">選擇操作</h3>
          <p v-if="item" class="text-sm text-slate-400 mt-1">
            {{ item.offering_name }}
          </p>
        </div>
        <div class="p-4 space-y-3">
          <!-- 個人行程選項 -->
          <template v-if="item?.is_personal_event">
            <button
              @click="$emit('action', 'edit')"
              class="w-full p-4 rounded-lg bg-success-500/20 border border-success-500/30 hover:bg-success-500/30 transition-colors text-left"
            >
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-success-500/30 flex items-center justify-center">
                  <svg class="w-5 h-5 text-success-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-white">編輯行程</div>
                  <div class="text-xs text-slate-400">修改行程時間或內容</div>
                </div>
              </div>
            </button>
            <button
              @click="$emit('action', 'note')"
              class="w-full p-4 rounded-lg bg-primary-500/20 border border-primary-500/30 hover:bg-primary-500/30 transition-colors text-left"
            >
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-primary-500/30 flex items-center justify-center">
                  <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-white">行程備註</div>
                  <div class="text-xs text-slate-400">為行程添加備註資訊</div>
                </div>
              </div>
            </button>
            <button
              @click="$emit('action', 'delete')"
              class="w-full p-4 rounded-lg bg-critical-500/20 border border-critical-500/30 hover:bg-critical-500/30 transition-colors text-left"
            >
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-critical-500/30 flex items-center justify-center">
                  <svg class="w-5 h-5 text-critical-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-white">刪除行程</div>
                  <div class="text-xs text-slate-400">移除此個人行程</div>
                </div>
              </div>
            </button>
          </template>

          <!-- 中心課程選項 -->
          <template v-else>
            <button
              @click="$emit('action', 'exception')"
              class="w-full p-4 rounded-lg bg-warning-500/20 border border-warning-500/30 hover:bg-warning-500/30 transition-colors text-left"
            >
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-warning-500/30 flex items-center justify-center">
                  <svg class="w-5 h-5 text-warning-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-white">課程例外申請</div>
                  <div class="text-xs text-slate-400">申請調課、請假或找代課</div>
                </div>
              </div>
            </button>
            <button
              @click="$emit('action', 'note')"
              class="w-full p-4 rounded-lg bg-primary-500/20 border border-primary-500/30 hover:bg-primary-500/30 transition-colors text-left"
            >
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-primary-500/30 flex items-center justify-center">
                  <svg class="w-5 h-5 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-white">課堂備註</div>
                  <div class="text-xs text-slate-400">撰寫或查看課程筆記</div>
                </div>
              </div>
            </button>
          </template>
        </div>
        <div class="p-4 border-t border-white/10">
          <button
            @click="$emit('close')"
            class="w-full px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
          >
            取消
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
// ============================================
// Props 定義
// ============================================

defineProps<{
  // 是否顯示
  visible: boolean
  // 選中的項目
  item: any
}>()

// ============================================
// Emits 定義
// ============================================

defineEmits<{
  close: []
  action: [action: 'exception' | 'note' | 'edit' | 'delete']
}>()
</script>
