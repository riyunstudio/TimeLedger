<template>
  <div class="conflict-legend flex flex-wrap items-center gap-4 p-3 bg-white/5 rounded-lg">
    <span class="text-xs text-slate-400 font-medium shrink-0">圖例說明</span>
    
    <!-- 完全可用 -->
    <div class="flex items-center gap-2">
      <div class="w-6 h-6 rounded bg-green-500/20 border border-green-500/30 flex items-center justify-center">
        <svg class="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <div class="text-xs">
        <span class="text-white block">完全可用</span>
        <span class="text-slate-500">可立即安排</span>
      </div>
    </div>

    <!-- 分隔線 -->
    <div class="w-px h-10 bg-white/10 shrink-0" />

    <!-- 緩衝衝突 -->
    <div class="flex items-center gap-2">
      <div class="w-6 h-6 rounded bg-yellow-500/20 border border-yellow-500/30 flex items-center justify-center">
        <svg class="w-4 h-4 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <div class="text-xs">
        <span class="text-white block">緩衝衝突</span>
        <span class="text-slate-500">可選擇_override</span>
      </div>
    </div>

    <!-- 分隔線 -->
    <div class="w-px h-10 bg-white/10 shrink-0" />

    <!-- 時間重疊 -->
    <div class="flex items-center gap-2">
      <div class="w-6 h-6 rounded bg-red-500/20 border border-red-500/30 flex items-center justify-center">
        <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </div>
      <div class="text-xs">
        <span class="text-white block">時間重疊</span>
        <span class="text-slate-500">不可排課</span>
      </div>
    </div>

    <!-- 分隔線 -->
    <div class="w-px h-10 bg-white/10 shrink-0" />

    <!-- 目標時段 -->
    <div class="flex items-center gap-2">
      <div class="w-6 h-6 rounded bg-indigo-500/20 border border-indigo-500/30 flex items-center justify-center">
        <div class="w-3 h-3 bg-indigo-500 rounded-sm" />
      </div>
      <div class="text-xs">
        <span class="text-white block">目標時段</span>
        <span class="text-slate-500">您想安排的時間</span>
      </div>
    </div>

    <!-- 分隔線 -->
    <div class="w-px h-10 bg-white/10 shrink-0" />

    <!-- 現有課程 -->
    <div class="flex items-center gap-2">
      <div class="w-6 h-6 rounded bg-slate-600/50 border border-slate-500 flex items-center justify-center">
        <div class="w-4 h-2 bg-slate-400 rounded-sm" />
      </div>
      <div class="text-xs">
        <span class="text-white block">現有課程</span>
        <span class="text-slate-500">老師已排課表</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 此組件為純展示性質，無需處理邏輯
</script>
