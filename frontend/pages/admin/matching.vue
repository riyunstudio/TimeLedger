<template>
  <div class="p-4 md:p-6" role="main" aria-label="智慧媒合系統">
    <h1 class="text-2xl font-bold text-white mb-6">智慧媒合</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="glass-card p-4 md:p-6" role="region" aria-label="搜尋條件">
        <h2 class="text-lg font-semibold text-white mb-4">搜尋條件</h2>

        <!-- 近期搜尋快速載入 -->
        <AdminRecentSearches
          ref="recentSearchesRef"
          @load-search="onLoadRecentSearch"
          @clear-all="onClearRecentSearches"
        />

        <form @submit.prevent="findMatches" class="space-y-4" aria-label="媒合搜尋表單">
          <div>
            <label id="datetime-range-label" class="block text-slate-300 mb-2">課程時段</label>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4" role="group" aria-labelledby="datetime-range-label">
              <div>
                <label for="start-time" class="block text-slate-400 text-sm mb-1">開始時間</label>
                <input
                  id="start-time"
                  v-model="form.start_time"
                  type="datetime-local"
                  aria-label="搜尋開始時間"
                  class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
                />
              </div>
              <div>
                <label for="end-time" class="block text-slate-400 text-sm mb-1">結束時間</label>
                <input
                  id="end-time"
                  v-model="form.end_time"
                  type="datetime-local"
                  aria-label="搜尋結束時間"
                  class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
                />
              </div>
            </div>
          </div>

          <!-- 教室卡片式選擇 -->
          <AdminRoomCardSelect
            :rooms="rooms"
            v-model="form.room_ids"
            role="group"
            aria-label="選擇教室"
          />

          <!-- 技能 autocomplete 選擇器 -->
          <AdminSkillSelector v-model="selectedSkills" aria-label="選擇技能" />

          <div class="flex flex-col sm:flex-row gap-3 pt-4">
            <button
              type="button"
              @click="clearForm"
              aria-label="清除搜尋條件"
              class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              清除
            </button>
            <button
              type="submit"
              :disabled="searching"
              aria-label="開始媒合搜尋"
              :aria-busy="searching"
              class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50"
            >
              <span v-if="searching" role="status">
                <span class="animate-pulse">搜尋中...</span>
              </span>
              <span v-else>開始媒合</span>
            </button>
          </div>
        </form>
      </div>

      <div class="glass-card p-4 md:p-6" role="region" aria-label="媒合結果">
        <h2 class="text-lg font-semibold text-white mb-4">媒合結果</h2>

        <div v-if="!hasSearched" class="text-center py-12 text-slate-500" role="status">
          請設定搜尋條件並點擊「開始媒合」
        </div>

        <div v-else-if="matches.length === 0" class="text-center py-12 text-slate-500" role="status">
          沒有找到符合條件的老師
        </div>

        <!-- 比較模式 -->
        <div v-else-if="viewMode === 'compare'">
          <AdminCompareMode
            :selected-teachers="sortedMatches.filter(m => selectedForCompare.has(m.teacher_id))"
            @remove="removeFromCompare"
            @select="selectTeacher"
            @exit="exitCompareMode"
          />
        </div>

        <!-- 搜尋中載入進度 -->
        <div v-else-if="searching && matches.length === 0" class="py-12">
          <div class="text-center mb-4">
            <div class="inline-flex items-center gap-3">
              <div class="flex items-center gap-1">
                <span class="w-3 h-3 rounded-full bg-primary-500 animate-bounce" style="animation-delay: 0ms;"></span>
                <span class="w-3 h-3 rounded-full bg-primary-500 animate-bounce" style="animation-delay: 150ms;"></span>
                <span class="w-3 h-3 rounded-full bg-primary-500 animate-bounce" style="animation-delay: 300ms;"></span>
              </div>
              <span class="text-slate-400">正在搜尋符合條件的老師...</span>
            </div>
          </div>

          <!-- 進度條 -->
          <div class="max-w-md mx-auto">
            <div class="flex items-center justify-between text-xs text-slate-500 mb-1">
              <span>正在分析老師可用性</span>
              <span>{{ searchProgress }}%</span>
            </div>
            <div class="h-2 bg-white/10 rounded-full overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-primary-500 to-secondary-500 transition-all duration-300"
                :style="{ width: `${searchProgress}%` }"
              ></div>
            </div>
          </div>

          <!-- 搜尋步驟 -->
          <div class="mt-6 space-y-2 max-w-sm mx-auto">
            <div
              v-for="(step, index) in searchSteps"
              :key="step.name"
              class="flex items-center gap-3 text-sm"
              :class="step.completed ? 'text-success-500' : step.active ? 'text-white' : 'text-slate-600'"
            >
              <div class="w-5 h-5 rounded-full flex items-center justify-center">
                <svg v-if="step.completed" class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <span v-else-if="step.active" class="w-2 h-2 rounded-full bg-primary-500 animate-pulse"></span>
                <span v-else class="w-2 h-2 rounded-full bg-slate-600"></span>
              </div>
              <span>{{ step.name }}</span>
            </div>
          </div>
        </div>

        <!-- 卡片/列表檢視 -->
        <template v-else>
          <!-- 排序控制 -->
          <AdminSortControls
            v-if="matches.length > 0"
            v-model:sort-by="sortBy"
            v-model:sort-order="sortOrder"
            v-model:view-mode="viewMode"
            :selected-count="selectedForCompare.size"
          />

          <!-- 結果列表 -->
          <div v-if="viewMode === 'card'" class="grid grid-cols-1 gap-3">
            <AdminEnhancedMatchCard
              v-for="match in sortedMatches"
              :key="match.teacher_id"
              :match="match"
              :selected="selectedForCompare.has(match.teacher_id)"
              :show-checkbox="sortedMatches.length > 1"
              @update:selected="toggleCompare(match.teacher_id, $event)"
              @click="viewTeacherTimeline(match)"
            />
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="match in sortedMatches"
              :key="match.teacher_id"
              :class="[
                'p-4 rounded-lg transition-all cursor-pointer',
                selectedForCompare.has(match.teacher_id)
                  ? 'bg-indigo-500/10 border border-indigo-500/30'
                  : 'bg-white/5 border border-white/10 hover:bg-white/10'
              ]"
              @click="viewTeacherTimeline(match)"
            >
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-3">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center shrink-0">
                    <span class="text-white font-medium">{{ match.teacher_name?.charAt(0) || '?' }}</span>
                  </div>
                  <div>
                    <h3 class="text-white font-medium">{{ match.teacher_name }}</h3>
                    <p class="text-sm text-slate-400">匹配度 {{ match.match_score }}%</p>
                  </div>
                </div>
                <div class="flex items-center gap-3">
                  <div class="text-right">
                    <div class="text-2xl font-bold text-primary-500">{{ match.match_score }}%</div>
                  </div>
                  <input
                    v-if="sortedMatches.length > 1"
                    type="checkbox"
                    :checked="selectedForCompare.has(match.teacher_id)"
                    @change="toggleCompare(match.teacher_id, !selectedForCompare.has(match.teacher_id))"
                    @click.stop
                    class="w-5 h-5 rounded border-white/20 bg-white/10 text-indigo-500 focus:ring-indigo-500"
                  />
                </div>
              </div>

              <div class="flex flex-wrap items-center gap-3 text-sm text-slate-400 mb-3">
                <span class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  技能匹配: {{ match.skill_match }}%
                </span>
                <span class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                  </svg>
                  評分: {{ match.rating?.toFixed(1) || '-' }}
                </span>
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-xs',
                    match.availability === 'AVAILABLE' ? 'bg-green-500/20 text-green-400' :
                    match.availability === 'BUFFER_CONFLICT' ? 'bg-yellow-500/20 text-yellow-400' :
                    'bg-red-500/20 text-red-400'
                  ]"
                >
                  {{ match.availability === 'AVAILABLE' ? '完全可用' :
                     match.availability === 'BUFFER_CONFLICT' ? '緩衝衝突' : '時間重疊' }}
                </span>
              </div>

              <!-- 可用教室標註 -->
              <div v-if="match.available_rooms?.length" class="mb-3">
                <p class="text-xs text-slate-500 mb-2">可授課教室：</p>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="room in match.available_rooms"
                    :key="room.id"
                    class="px-2 py-1 rounded-md text-xs bg-success-500/20 text-success-500"
                  >
                    {{ room.name }}
                  </span>
                </div>
              </div>

              <p v-if="match.notes" class="text-sm text-slate-400 mb-3">
                {{ match.notes }}
              </p>

              <!-- 查看課表按鈕 -->
              <button
                @click.stop="viewTeacherTimeline(match)"
                class="mt-2 px-3 py-1.5 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white text-sm transition-colors flex items-center gap-1"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                查看課表衝突
              </button>
            </div>
          </div>

          <!-- 教師課表時間軸（展開時顯示） -->
          <div v-if="selectedTeacher" class="mt-6">
            <!-- 衝突圖例 -->
            <AdminConflictLegend class="mb-4" />

            <!-- 時間軸 -->
            <AdminTeacherTimeline
              :teacher="selectedTeacher"
              :target-start="form.start_time"
              :target-end="form.end_time"
              :existing-sessions="teacherSessions"
              @override="confirmOverride"
            />

            <!-- 替代時段建議 -->
            <div v-if="alternativeSlots.length > 0" class="mt-4">
              <AdminAlternativeSlots
                :teacher-name="selectedTeacher.teacher_name"
                :slots="alternativeSlots"
                @select="selectAlternativeSlot"
                @custom="showCustomTime = true"
              />
            </div>
          </div>
        </template>
      </div>
    </div>

    <div class="mt-6 glass-card p-4 md:p-6">
      <!-- 標題與展開收合 -->
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold text-white">人才庫搜尋</h2>
        <button
          v-if="talentResults.length > 0"
          @click="showStats = !showStats"
          class="text-sm text-slate-400 hover:text-white transition-colors flex items-center gap-1"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          {{ showStats ? '隱藏統計' : '顯示統計' }}
        </button>
      </div>

      <!-- 統計面板（可展開） -->
      <Transition
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 -translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition-all duration-200 ease-in"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 -translate-y-2"
      >
        <div v-if="showStats && talentResults.length > 0" class="mb-6">
          <!-- 統計卡片 -->
          <AdminTalentStatsPanel
            :stats="talentStats"
            :city-distribution="cityDistribution"
            :top-skills="topSkills"
          />
          
          <!-- 技能分布圖 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
            <AdminSkillsDistributionChart
              :skills="skillDistribution"
              :total-teachers="talentResults.length"
              @select-skill="onSkillSelected"
            />
            
            <!-- 搜尋建議熱門技能 -->
            <div class="bg-white/5 rounded-xl p-4">
              <h4 class="text-sm font-medium text-white mb-3 flex items-center gap-2">
                <svg class="w-4 h-4 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                </svg>
                熱門搜尋技能
              </h4>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="skill in popularSearchSkills"
                  :key="skill"
                  @click="onSkillSelected(skill)"
                  class="px-3 py-1.5 rounded-lg bg-white/5 hover:bg-white/10 text-slate-300 hover:text-white text-sm transition-colors"
                >
                  {{ skill }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- 搜尋輸入框 -->
      <AdminSearchSuggestions
        v-model="talentSearch.query"
        @search="searchTalent"
        @select-suggestion="onSuggestionSelected"
      />

      <!-- 進階篩選面板 -->
      <div class="mt-4">
        <AdminTalentFilterPanel
          ref="talentFilterRef"
          @apply="applyTalentFilters"
          @clear="clearTalentFilters"
        />
      </div>

      <!-- 快速篩選標籤 -->
      <AdminQuickFilterTags
        v-if="talentResults.length > 0"
        :popular-skills="popularSearchSkills"
        :popular-tags="['古典', '兒童', '成人', '入門', '進階']"
        :active-skills="activeSkillFilters"
        :active-tags="activeTagFilters"
        @filter-by-skill="onSkillFilter"
        @filter-by-tag="onTagFilter"
        @clear-all="clearQuickFilters"
      />

      <!-- 搜尋結果數量與分頁 -->
      <div v-if="talentResults.length > 0" class="mt-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <span class="text-sm text-slate-400">
          找到 <span class="text-white font-medium">{{ totalItems }}</span> 位人才
          <span class="text-slate-500 ml-2">(第 {{ currentPage }}/{{ totalPages }} 頁)</span>
        </span>

        <div class="flex items-center gap-2">
          <!-- 排序選單 -->
          <select
            v-model="talentSortBy"
            @change="searchTalent(undefined, 1)"
            class="px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-white text-sm"
          >
            <option value="name">姓名</option>
            <option value="skills">技能數</option>
            <option value="rating">評分</option>
            <option value="city">地區</option>
          </select>

          <button
            @click="talentSortOrder = talentSortOrder === 'asc' ? 'desc' : 'asc'; searchTalent(undefined, 1)"
            class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-white transition-colors"
            title="切換排序順序"
          >
            <svg v-if="talentSortOrder === 'asc'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h9m5-4v12m0 0l-4-4m4 4l4-4" />
            </svg>
          </button>
        </div>
      </div>

      <!-- 分頁控制項 -->
      <div v-if="totalPages > 1" class="mt-4 flex items-center justify-center gap-2">
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage <= 1"
          class="px-3 py-1.5 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors text-sm"
        >
          上一頁
        </button>

        <div class="flex items-center gap-1">
          <button
            v-for="page in visiblePages"
            :key="page"
            @click="goToPage(page)"
            class="w-8 h-8 rounded-lg text-sm transition-colors"
            :class="page === currentPage
              ? 'bg-primary-500 text-white'
              : 'bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white'"
          >
            {{ page }}
          </button>
        </div>

        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage >= totalPages"
          class="px-3 py-1.5 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors text-sm"
        >
          下一頁
        </button>

        <span class="text-xs text-slate-500 ml-2">
          跳至第
          <input
            type="number"
            min="1"
            :max="totalPages"
            v-model.number="currentPage"
            @keyup.enter="goToPage(currentPage)"
            class="w-12 px-2 py-1 mx-1 rounded bg-white/5 border border-white/10 text-white text-center text-sm"
          />
          頁
        </span>
      </div>

      <!-- 搜尋結果 -->
      <div v-if="talentResults.length > 0" class="mt-4">
        <!-- 批量操作工具列 -->
        <AdminBulkActions
          :selected-count="selectedTalents.size"
          :bulk-loading="bulkLoading"
          :bulk-progress="bulkProgress"
          @clear="clearTalentSelection"
          @compare="viewMode = 'compare'"
          @export="exportContactInfo"
          @bulk-invite="bulkInviteTalents"
        />

        <!-- 人才卡片網格 -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-20">
          <AdminTalentCard
            v-for="teacher in filteredResults"
            :key="teacher.id"
            :teacher="teacher"
            :selected="selectedTalents.has(teacher.id)"
            :show-checkbox="teacher.is_open_to_hiring"
            :invite-loading="inviteLoadingIds.has(teacher.id)"
            :invitation-status="invitationStatuses.get(teacher.id) || null"
            @update:selected="toggleTalentSelection(teacher.id, $event)"
            @invite="inviteTalent"
            @view="viewTalentDetail"
            @compare="addTalentToCompare"
          />
        </div>
      </div>

      <!-- 空狀態 -->
      <div v-else-if="hasSearchedTalent" class="text-center py-12 text-slate-500">
        <svg class="w-16 h-16 mx-auto mb-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p>沒有找到符合條件的人才</p>
        <button
          @click="clearTalentSearch"
          class="mt-4 px-4 py-2 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 transition-colors"
        >
          清除搜尋條件
        </button>
      </div>
    </div>
  </div>

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
import { SKILL_CATEGORIES } from '~/types'
import { formatDateToString, getTodayString } from '~/composables/useTaiwanTime'
import AdminRecentSearches from '~/components/Admin/RecentSearches.vue'
import AdminRoomCardSelect from '~/components/Admin/RoomCardSelect.vue'
import AdminSkillSelector from '~/components/Admin/SkillSelector.vue'
import AdminSortControls from '~/components/Admin/SortControls.vue'
import AdminCompareMode from '~/components/Admin/CompareMode.vue'
import AdminEnhancedMatchCard from '~/components/Admin/EnhancedMatchCard.vue'
import AdminTeacherTimeline from '~/components/Admin/TeacherTimeline.vue'
import AdminConflictLegend from '~/components/Admin/ConflictLegend.vue'
import AdminAlternativeSlots from '~/components/Admin/AlternativeSlots.vue'
import AdminTalentCard from '~/components/Admin/TalentCard.vue'
import AdminBulkActions from '~/components/Admin/BulkActions.vue'
import AdminTalentFilterPanel from '~/components/Admin/TalentFilterPanel.vue'
import AdminQuickFilterTags from '~/components/Admin/QuickFilterTags.vue'
import AdminSearchSuggestions from '~/components/Admin/SearchSuggestions.vue'
import AdminTalentStatsPanel from '~/components/Admin/TalentStatsPanel.vue'
import AdminSkillsDistributionChart from '~/components/Admin/SkillsDistributionChart.vue'

interface SelectedSkill {
  id: number
  name: string
  category: string
  teacherCount?: number
}

interface RecentSearch {
  id: string
  start_time: string
  end_time: string
  room_ids: number[]
  skills: SelectedSkill[]
  created_at: number
}

interface Session {
  id: number
  course_name: string
  start_time: string
  end_time: string
  room_name?: string
}

interface AlternativeSlot {
  date: string
  dateLabel: string
  start: string
  end: string
  available: boolean
  availableRooms: { id: number; name: string }[]
  conflictReason?: string
}

interface MatchResult {
  teacher_id: number
  teacher_name: string
  match_score: number
  availability: string
  availability_score: number
  internal_score: number
  skill_match: number
  skill_score: number
  region_score?: number
  rating?: number
  is_member: boolean
  notes?: string
  available_rooms?: { id: number; name: string }[]
  skills?: { name: string; category: string }[]
}

interface TalentResult {
  id: number
  name: string
  bio?: string
  city?: string
  district?: string
  skills?: Array<{ name: string; category: string }>
  personal_hashtags?: string[]
  is_open_to_hiring: boolean
  is_member: boolean
  internal_rating: number
  public_contact_info?: string
}

interface InvitationStatus {
  sent: boolean
  variant: 'success' | 'warning' | 'error' | 'secondary'
  text: string
}

type SortBy = 'score' | 'availability' | 'rating' | 'skill-match'
type SortOrder = 'asc' | 'desc'
type ViewMode = 'card' | 'list' | 'compare'

 definePageMeta({
   middleware: 'auth-admin',
   layout: 'admin',
 })

 const notificationUI = useNotification()
const { warning: alertWarning, success: alertSuccess } = useAlert()
const searching = ref(false)
const talentSearching = ref(false)
const hasSearched = ref(false)
const matches = ref<MatchResult[]>([])
const rooms = ref<any[]>([])
const { getCenterId } = useCenterId()

// 搜尋進度相關
const searchProgress = ref(0)
const searchSteps = ref([
  { name: '取得老師清單', completed: false, active: false },
  { name: '分析可用時間', completed: false, active: false },
  { name: '計算技能匹配度', completed: false, active: false },
  { name: '評估緩衝時間', completed: false, active: false },
  { name: '產生媒合結果', completed: false, active: false }
])

// 更新搜尋進度
const updateSearchProgress = () => {
  const steps = searchSteps.value
  let progress = 0
  let activeIndex = -1

  // 找到當前進度
  for (let i = 0; i < steps.length; i++) {
    if (steps[i].completed) {
      progress = Math.floor(((i + 1) / steps.length) * 100)
      activeIndex = i
    }
  }

  // 計算實際進度
  progress = Math.min(progress, 100)
  searchProgress.value = progress

  // 更新步驟狀態
  steps.forEach((step, index) => {
    step.active = index === activeIndex + 1 && !step.completed
  })
}

// 重置搜尋進度
const resetSearchProgress = () => {
  searchProgress.value = 0
  searchSteps.value = [
    { name: '取得老師清單', completed: false, active: true },
    { name: '分析可用時間', completed: false, active: false },
    { name: '計算技能匹配度', completed: false, active: false },
    { name: '評估緩衝時間', completed: false, active: false },
    { name: '產生媒合結果', completed: false, active: false }
  ]
}

// 完成當前步驟並開始下一個
const completeSearchStep = (stepIndex: number) => {
  if (stepIndex < searchSteps.value.length) {
    searchSteps.value[stepIndex].completed = true
    searchSteps.value[stepIndex].active = false
    if (stepIndex + 1 < searchSteps.value.length) {
      searchSteps.value[stepIndex + 1].active = true
    }
    updateSearchProgress()
  }
}

// 近期搜尋組件引用
const recentSearchesRef = ref<InstanceType<typeof AdminRecentSearches> | null>(null)

// 選取的技能列表
const selectedSkills = ref<SelectedSkill[]>([])

// 排序和檢視模式
const sortBy = ref<SortBy>('score')
const sortOrder = ref<SortOrder>('desc')
const viewMode = ref<ViewMode>('card')
const selectedForCompare = ref<Set<number>>(new Set())

// 選中的教師時間軸
const selectedTeacher = ref<MatchResult | null>(null)
const teacherSessions = ref<Session[]>([])
const alternativeSlots = ref<AlternativeSlot[]>([])
const showCustomTime = ref(false)

const form = ref({
  start_time: '',
  end_time: '',
  room_ids: [] as number[],
  skills: '' as string
})

const talentSearch = ref({
  query: '',
  city: '',
  skills: '',
  hashtags: ''
})

// 人才庫管理狀態
const talentResults = ref<TalentResult[]>([])
const selectedTalents = ref<Set<number>>(new Set())
const invitationStatuses = ref<Map<number, InvitationStatus>>(new Map())
const inviteLoadingIds = ref<Set<number>>(new Set())
const bulkLoading = ref(false)
const bulkProgress = ref(0)
const hasSearchedTalent = ref(false)
const showStats = ref(false)

// 分頁與排序狀態
const talentSortBy = ref('name')
const talentSortOrder = ref<'asc' | 'desc'>('asc')
const activeSkillFilters = ref<string[]>([])
const activeTagFilters = ref<string[]>([])

// 分頁狀態
const currentPage = ref(1)
const itemsPerPage = ref(20)
const totalItems = ref(0)
const totalPages = ref(0)

// 統計資料 - 從 API 載入
const talentStats = ref({
  totalCount: 0,
  openHiringCount: 0,
  memberCount: 0,
  averageRating: 0,
  monthlyChange: 0,
  monthlyTrend: [],
  pendingInvites: 0,
  acceptedInvites: 0,
  declinedInvites: 0
})

const cityDistribution = ref<Array<{ name: string; count: number }>>([])
const topSkills = ref<Array<{ name: string; count: number }>>([])
const skillDistribution = ref<Array<{ name: string; count: number }>>([])
const popularSearchSkills = ref<string[]>([])

// 取得人才庫統計
const fetchTalentStats = async () => {
  try {
    const api = useApi()
    const response = await api.get<{ code: number; datas: any }>(
      '/admin/smart-matching/talent/stats'
    )

    if (response.code === 0 && response.datas) {
      const stats = response.datas
      talentStats.value = {
        totalCount: stats.total_count || 0,
        openHiringCount: stats.open_hiring_count || 0,
        memberCount: stats.member_count || 0,
        averageRating: stats.average_rating || 0,
        monthlyChange: stats.monthly_change || 0,
        monthlyTrend: stats.monthly_trend || [],
        pendingInvites: stats.pending_invites || 0,
        acceptedInvites: stats.accepted_invites || 0,
        declinedInvites: stats.declined_invites || 0
      }

      // 更新城市分布
      if (stats.city_distribution) {
        cityDistribution.value = stats.city_distribution
      }

      // 更新熱門技能
      if (stats.top_skills) {
        topSkills.value = stats.top_skills
        skillDistribution.value = stats.top_skills.slice(0, 8)
        popularSearchSkills.value = stats.top_skills.slice(0, 6).map((s: any) => s.name)
      }
    }
  } catch (error) {
    console.error('Failed to fetch talent stats:', error)
    // 發生錯誤時清空統計
    talentStats.value = {
      totalCount: 0,
      openHiringCount: 0,
      memberCount: 0,
      averageRating: 0,
      monthlyChange: 0,
      monthlyTrend: [],
      pendingInvites: 0,
      acceptedInvites: 0,
      declinedInvites: 0
    }
    cityDistribution.value = []
    topSkills.value = []
    skillDistribution.value = []
    popularSearchSkills.value = []
  }
}

// 篩選後的結果（用於顯示，實際分頁由後端處理）
const filteredResults = computed(() => {
  return talentResults.value
})

// 計算可見的頁碼
const visiblePages = computed(() => {
  const pages: number[] = []
  const maxVisible = 5 // 最多顯示 5 個頁碼

  if (totalPages.value <= maxVisible) {
    // 如果總頁數少於最大顯示數，顯示所有頁碼
    for (let i = 1; i <= totalPages.value; i++) {
      pages.push(i)
    }
  } else {
    // 計算顯示範圍
    let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
    let end = start + maxVisible - 1

    if (end > totalPages.value) {
      end = totalPages.value
      start = Math.max(1, end - maxVisible + 1)
    }

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }
  }

  return pages
})

// 技能類別相關函數
const getSkillCategoryIcon = (category: string): string => {
  return SKILL_CATEGORIES[category as keyof typeof SKILL_CATEGORIES]?.icon || '✨'
}

const getSkillCategoryStyle = (category: string): string => {
  return SKILL_CATEGORIES[category as keyof typeof SKILL_CATEGORIES]?.color || 'bg-slate-500/20 text-slate-400 border-slate-500/30'
}

// 排序後的結果
const sortedMatches = computed(() => {
  const result = [...matches.value]
  
  result.sort((a, b) => {
    let comparison = 0
    
    switch (sortBy.value) {
      case 'score':
        comparison = a.match_score - b.match_score
        break
      case 'availability':
        const availOrder = { 'AVAILABLE': 0, 'BUFFER_CONFLICT': 1, 'OVERLAP': 2 }
        comparison = (availOrder[a.availability as keyof typeof availOrder] || 0) - 
                     (availOrder[b.availability as keyof typeof availOrder] || 0)
        break
      case 'rating':
        comparison = (a.rating || 0) - (b.rating || 0)
        break
      case 'skill-match':
        comparison = a.skill_match - b.skill_match
        break
    }
    
    return sortOrder.value === 'desc' ? -comparison : comparison
  })
  
  return result
})

// 從近期搜尋載入
const onLoadRecentSearch = (search: RecentSearch) => {
  form.value.start_time = search.start_time.slice(0, 16)
  form.value.end_time = search.end_time.slice(0, 16)
  form.value.room_ids = [...search.room_ids]
  selectedSkills.value = [...search.skills]
}

const onClearRecentSearches = () => {
  // 清除時不執行額外動作
}

// 比較功能
const toggleCompare = (teacherId: number, selected: boolean) => {
  if (selected) {
    if (selectedForCompare.value.size >= 3) {
      alertWarning('最多只能選取 3 位老師進行比較')
      return
    }
    selectedForCompare.value.add(teacherId)
  } else {
    selectedForCompare.value.delete(teacherId)
  }
}

const removeFromCompare = (teacher: MatchResult) => {
  selectedForCompare.value.delete(teacher.teacher_id)
}

const exitCompareMode = () => {
  viewMode.value = 'card'
}

// 人才庫管理功能
// 切換人才選取
const toggleTalentSelection = (teacherId: number, selected: boolean) => {
  if (selected) {
    selectedTalents.value.add(teacherId)
  } else {
    selectedTalents.value.delete(teacherId)
  }
}

// 清除人才選取
const clearTalentSelection = () => {
  selectedTalents.value.clear()
}

// 邀請合作
const inviteTalent = async (teacher: TalentResult) => {
  inviteLoadingIds.value.add(teacher.id)

  try {
    const api = useApi()
    const response = await api.post<{ code: number; datas: any }>(
      '/admin/smart-matching/talent/invite',
      {
        teacher_ids: [teacher.id],
        message: ''
      }
    )

    if (response.code === 0) {
      // 更新邀請狀態
      invitationStatuses.value.set(teacher.id, {
        sent: true,
        variant: 'success',
        text: '已邀請'
      })

      alertSuccess(`已向 ${teacher.name} 發送合作邀請`)
    } else {
      throw new Error(response.message || '邀請失敗')
    }
  } catch (error) {
    console.error('Failed to send invitation:', error)
    invitationStatuses.value.set(teacher.id, {
      sent: false,
      variant: 'error',
      text: '邀請失敗'
    })
    alertWarning('邀請失敗，請稍後再試')
  } finally {
    inviteLoadingIds.value.delete(teacher.id)
  }
}

// 批量邀請
const bulkInviteTalents = async () => {
  if (selectedTalents.value.size === 0) {
    await alertWarning('請先選取要邀請的人才')
    return
  }

  if (!(await alertConfirm(`確定要邀請選取的 ${selectedTalents.value.size} 位人才嗎？`))) {
    return
  }

  bulkLoading.value = true
  bulkProgress.value = 0

  const selectedIds = Array.from(selectedTalents.value)

  try {
    const api = useApi()
    const response = await api.post<{ code: number; datas: any }>(
      '/admin/smart-matching/talent/invite',
      {
        teacher_ids: selectedIds,
        message: ''
      }
    )

    if (response.code === 0) {
      const result = response.datas

      // 更新邀請狀態
      selectedIds.forEach(teacherId => {
        if (result.failed_ids && result.failed_ids.includes(teacherId)) {
          invitationStatuses.value.set(teacherId, {
            sent: false,
            variant: 'error',
            text: '邀請失敗'
          })
        } else {
          invitationStatuses.value.set(teacherId, {
            sent: true,
            variant: 'success',
            text: '已邀請'
          })
        }
      })

      bulkProgress.value = 100
      alertSuccess(`已成功邀請 ${result.invited_count} 位人才，失敗 ${result.failed_count} 位`)
    } else {
      throw new Error(response.message || '邀請失敗')
    }
  } catch (error) {
    console.error('Failed to bulk invite:', error)
    alertWarning('批量邀請失敗，請稍後再試')

    // 標記所有選中為失敗
    selectedIds.forEach(teacherId => {
      invitationStatuses.value.set(teacherId, {
        sent: false,
        variant: 'error',
        text: '邀請失敗'
      })
    })
  } finally {
    bulkLoading.value = false
    bulkProgress.value = 0
    selectedTalents.value.clear()
  }
}

// 匯出聯絡資訊
const exportContactInfo = () => {
  const selectedData = talentResults.value.filter(t => selectedTalents.value.has(t.id))
  
  if (selectedData.length === 0) {
    alertWarning('請先選取要匯出的人才')
    return
  }

  // 生成 CSV 格式
  const headers = ['姓名', '城市', '區域', '聯絡方式', '技能', '標籤']
  const rows = selectedData.map(t => [
    t.name,
    t.city || '',
    t.district || '',
    t.public_contact_info || '',
    t.skills?.map(s => s.name).join(', ') || '',
    t.personal_hashtags?.join(', ') || ''
  ])
  
  const csv = [headers.join(','), ...rows.map(r => r.map(cell => `"${cell}"`).join(','))].join('\n')
  
  // 下載檔案
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `人才庫匯出_${getTodayString()}.csv`
  link.click()
  URL.revokeObjectURL(link.href)
  
  alertSuccess(`已匯出 ${selectedData.length} 位人才的聯絡資訊`)
}

// 人才庫搜尋（整合新功能）
const searchTalent = async (query?: string, page: number = 1) => {
  talentSearching.value = true
  hasSearchedTalent.value = true

  try {
    const api = useApi()
    const params = new URLSearchParams()

    if (query || talentSearch.value.query) {
      params.append('keyword', query || talentSearch.value.query)
    }
    if (talentSearch.value.city) params.append('city', talentSearch.value.city)
    if (talentSearch.value.skills) params.append('skills', talentSearch.value.skills)
    if (talentSearch.value.hashtags) params.append('hashtags', talentSearch.value.hashtags)

    // 分頁參數
    params.append('page', page.toString())
    params.append('limit', itemsPerPage.value.toString())

    // 排序參數
    params.append('sort_by', talentSortBy.value)
    params.append('sort_order', talentSortOrder.value)

    const response = await api.get<{ code: number; datas: any }>(
      `/admin/smart-matching/talent/search?${params.toString()}`
    )

    if (response.code === 0 && response.datas) {
      const data = response.datas
      talentResults.value = data.talents || []

      // 更新分頁資訊
      if (data.pagination) {
        currentPage.value = data.pagination.page
        itemsPerPage.value = data.pagination.limit
        totalItems.value = data.pagination.total
        totalPages.value = data.pagination.total_pages
      }
    } else {
      talentResults.value = []
      totalItems.value = 0
      totalPages.value = 0
    }

    // 顯示統計面板
    showStats.value = true
  } catch (error) {
    console.error('Failed to search talent:', error)
    talentResults.value = []
    totalItems.value = 0
    totalPages.value = 0
    showStats.value = true
    await alertError('搜尋人才失敗，請稍後再試')
  } finally {
    talentSearching.value = false
  }
}

// 切換分頁
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    searchTalent(undefined, page)
    // 捲動到頂部
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

// 搜尋建議選取
const onSuggestionSelected = (suggestion: { type: string; value: string }) => {
  talentSearch.value.query = suggestion.value
  if (suggestion.type === 'skill') {
    talentSearch.value.skills = suggestion.value
  } else if (suggestion.type === 'tag') {
    talentSearch.value.hashtags = suggestion.value
  }
  searchTalent()
}

// 技能選取
const onSkillSelected = (skillName: string) => {
  talentSearch.value.skills = skillName
  activeSkillFilters.value = [skillName]
  searchTalent()
}

// 技能篩選
const onSkillFilter = (skill: string) => {
  if (activeSkillFilters.value.includes(skill)) {
    activeSkillFilters.value = activeSkillFilters.value.filter(s => s !== skill)
  } else {
    activeSkillFilters.value = [...activeSkillFilters.value, skill]
  }
}

// 標籤篩選
const onTagFilter = (tag: string) => {
  if (activeTagFilters.value.includes(tag)) {
    activeTagFilters.value = activeTagFilters.value.filter(t => t !== tag)
  } else {
    activeTagFilters.value = [...activeTagFilters.value, tag]
  }
}

// 清除快速篩選
const clearQuickFilters = () => {
  activeSkillFilters.value = []
  activeTagFilters.value = []
}

// 人才篩選面板 - 套用篩選
const applyTalentFilters = (filters: any) => {
  if (filters.city) talentSearch.value.city = filters.city
  if (filters.skills) talentSearch.value.skills = filters.skills
  if (filters.hashtags) talentSearch.value.hashtags = filters.hashtags
  searchTalent()
}

// 人才篩選面板 - 清除篩選
const clearTalentFilters = () => {
  talentSearch.value = {
    query: '',
    city: '',
    skills: '',
    hashtags: ''
  }
  clearQuickFilters()
  searchTalent()
}

// 清除人才搜尋
const clearTalentSearch = () => {
  talentSearch.value = {
    query: '',
    city: '',
    skills: '',
    hashtags: ''
  }
  hasSearchedTalent.value = false
  talentResults.value = []
  clearQuickFilters()
  showStats.value = false
}

// 查看人才詳細資訊
const viewTalentDetail = (teacher: TalentResult) => {
  alertSuccess(`正在查看 ${teacher.name} 的詳細資訊`)
}

// 加入比較
const addTalentToCompare = (teacher: TalentResult) => {
  if (selectedForCompare.value.size >= 3) {
    alertWarning('最多只能選取 3 位老師進行比較')
    return
  }
  selectedForCompare.value.add(teacher.id)
  alertSuccess(`已將 ${teacher.name} 加入比較`)
}

// 查看教師課表時間軸
const viewTeacherTimeline = async (teacher: MatchResult) => {
  selectedTeacher.value = teacher

  try {
    const api = useApi()
    const startDate = form.value.start_time ? formatDateToString(new Date(form.value.start_time)) : getTodayString()
    const endDate = form.value.end_time ? formatDateToString(new Date(form.value.end_time)) : formatDateToString(new Date(Date.now() + 7 * 24 * 60 * 60 * 1000))

    const response = await api.get<{ code: number; datas: any }>(
      `/admin/teachers/${teacher.teacher_id}/sessions?start_date=${startDate}&end_date=${endDate}`
    )

    if (response.code === 0 && response.datas) {
      const data = response.datas
      teacherSessions.value = data.sessions.map((s: any) => ({
        id: s.id,
        course_name: s.course_name,
        start_time: s.start_time,
        end_time: s.end_time,
        room_name: s.room_name || '未指定教室'
      }))
    } else {
      // 沒有課表資料
      teacherSessions.value = []
    }
  } catch (error) {
    console.error('Failed to fetch teacher sessions:', error)
    teacherSessions.value = []
    await alertError('取得老師課表失敗')
  }

  // 取得替代時段建議
  await fetchAlternativeSlots(teacher)

  alertSuccess(`正在查看 ${teacher.teacher_name} 的課表`)
}

// 取得替代時段建議
const fetchAlternativeSlots = async (teacher: MatchResult) => {
  if (!form.value.start_time || !form.value.end_time) {
    // 沒有指定時間時，顯示空替代時段
    alternativeSlots.value = []
    return
  }

  try {
    const api = useApi()
    const response = await api.post<{ code: number; datas: any[] }>(
      '/admin/smart-matching/alternatives',
      {
        teacher_id: teacher.teacher_id,
        original_start: formatDateToString(new Date(form.value.start_time)),
        original_end: formatDateToString(new Date(form.value.end_time)),
        duration: 90
      }
    )

    if (response.code === 0 && response.datas) {
      alternativeSlots.value = response.datas.map((slot: any) => ({
        date: slot.date,
        dateLabel: slot.date_label,
        start: slot.start,
        end: slot.end,
        available: slot.available,
        availableRooms: slot.available_rooms || [],
        conflictReason: slot.conflict_reason
      }))
    } else {
      alternativeSlots.value = []
    }
  } catch (error) {
    console.error('Failed to fetch alternative slots:', error)
    alternativeSlots.value = []
  }
}

// 選擇替代時段
const selectAlternativeSlot = (slot: AlternativeSlot) => {
  form.value.start_time = `${slot.date}T${slot.start}`
  form.value.end_time = `${slot.date}T${slot.end}`
  
  alertSuccess(`已選擇替代時段 ${slot.dateLabel} ${slot.start}-${slot.end}`)
  
  // 重新搜尋
  findMatches()
}

// 確認 Override
const confirmOverride = async () => {
  if (await alertConfirm('確定要在緩衝衝突的情況下安排這位老師嗎？')) {
    selectTeacher(selectedTeacher.value!)
  }
}

const findMatches = async () => {
  if (!form.value.start_time || !form.value.end_time) {
    await alertWarning('請填寫開始時間和結束時間')
    return
  }

  searching.value = true
  hasSearched.value = true
  matches.value = []

  // 重置比較狀態
  selectedForCompare.value.clear()
  viewMode.value = 'card'
  selectedTeacher.value = null

  // 初始化搜尋進度
  resetSearchProgress()

  // 模擬進度（Go 處理很快，這只是 UI 效果）
  const simulateProgress = async () => {
    // 步驟 1: 取得老師清單
    await delay(200)
    completeSearchStep(0)

    // 步驟 2: 分析可用時間
    await delay(150)
    completeSearchStep(1)

    // 步驟 3: 計算技能匹配度
    await delay(150)
    completeSearchStep(2)

    // 步驟 4: 評估緩衝時間
    await delay(100)
    completeSearchStep(3)

    // 步驟 5: 產生媒合結果
    await delay(100)
    completeSearchStep(4)
  }

  try {
    const api = useApi()

    // 並行執行：API 呼叫 + 進度模擬
    const [response] = await Promise.all([
      api.post<{ code: number; datas: MatchResult[] }>(
        `/admin/smart-matching/matches`,
        {
          room_ids: form.value.room_ids.length > 0 ? form.value.room_ids : undefined,
          start_time: formatDateToString(new Date(form.value.start_time)),
          end_time: formatDateToString(new Date(form.value.end_time)),
          required_skills: selectedSkills.value.length > 0
            ? selectedSkills.value.map(s => s.name)
            : form.value.skills.split(',').map(s => s.trim()).filter(Boolean)
        }
      ),
      simulateProgress()
    ])

    matches.value = (response as any).datas || []

    // 儲存到近期搜尋
    if (recentSearchesRef.value) {
      recentSearchesRef.value.addSearch({
        start_time: formatDateToString(new Date(form.value.start_time)),
        end_time: formatDateToString(new Date(form.value.end_time)),
        room_ids: [...form.value.room_ids],
        skills: [...selectedSkills.value]
      })
    }

    if (matches.value.length === 0) {
      alertWarning('沒有找到符合條件的老師')
    }
  } catch (error) {
    console.error('Failed to find matches:', error)
    matches.value = []
    await alertWarning('搜尋失敗，請稍後再試')
  } finally {
    searching.value = false
  }
}

// 延遲輔助函數
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

const selectTeacher = (match: MatchResult) => {
  alertSuccess(`已選擇 ${match.teacher_name}`)
}

const clearForm = () => {
  form.value = {
    start_time: '',
    end_time: '',
    room_ids: [],
    skills: ''
  }
  selectedSkills.value = []
  hasSearched.value = false
  matches.value = []
  selectedForCompare.value.clear()
  selectedTeacher.value = null
  teacherSessions.value = []
  alternativeSlots.value = []
}

const fetchRooms = async () => {
  try {
    const api = useApi()
    const centerId = getCenterId()
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    rooms.value = response.datas || []
  } catch (error) {
    console.error('Failed to fetch rooms:', error)
  }
}

onMounted(async () => {
  fetchRooms()
  if (recentSearchesRef.value) {
    recentSearchesRef.value.loadFromStorage()
  }

  // 載入人才庫統計
  await fetchTalentStats()
})
</script>
