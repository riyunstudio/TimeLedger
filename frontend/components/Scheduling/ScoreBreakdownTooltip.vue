<template>
  <HelpTooltip
    title="媒合分數說明"
    :description="scoreDescription"
    :usage="scoreBreakdownItems"
    placement="top"
  />
</template>

<script setup lang="ts">
/**
 * 媒合分數說明 Tooltip 組件
 *
 * 說明智慧媒合的分數組成：
 * - 可用性分數 (0-40分)：完全空閒 40 分，緩衝衝突 15 分，時間重疊 0 分
 * - 內部評價分數 (0-40分)：星等評分正規化 + 內部備註關鍵字額外加分
 * - 技能/地區匹配分數 (0-20分)：技能命中 10 分，標籤命中 8 分，地區命中 10 分
 */

const scoreDescription = '媒合分數由三個維度組成，總分 100 分。分數越高代表越適合安排課程。'

const scoreBreakdownItems = computed(() => [
  {
    label: '可用性 (40分)',
    value: '完全空閒 40 分，緩衝衝突 15 分，時間重疊 0 分',
  },
  {
    label: '內部評價 (40分)',
    value: '星等評分正規化 0-30 分 + 備註關鍵字額外 10 分',
  },
  {
    label: '技能匹配 (20分)',
    value: '技能命中 10 分，標籤命中 8 分，地區命中 10 分',
  },
])
</script>
