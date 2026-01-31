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
              :disabled="smartMatchingStore.isSearching"
              aria-label="開始媒合搜尋"
              :aria-busy="smartMatchingStore.isSearching"
              class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50"
            >
              <span v-if="smartMatchingStore.isSearching" role="status">
                <span class="animate-pulse">搜尋中...</span>
              </span>
              <span v-else>開始媒合</span>
            </button>
          </div>
        </form>
      </div>

      <div class="glass-card p-4 md:p-6" role="region" aria-label="媒合結果">
        <h2 class="text-lg font-semibold text-white mb-4">媒合結果</h2>

        <div v-if="!smartMatchingStore.hasSearched" class="text-center py-12 text-slate-500" role="status">
          請設定搜尋條件並點擊「開始媒合」
        </div>

        <div v-else-if="smartMatchingStore.searchResults.length === 0" class="text-center py-12 text-slate-500" role="status">
          沒有找到符合條件的老師
        </div>

        <!-- 比較模式 -->
        <div v-else-if="smartMatchingStore.viewMode === 'compare'">
          <AdminCompareMode
            :selected-teachers="smartMatchingStore.sortedResults.filter(m => smartMatchingStore.selectedForCompare.has(m.teacher_id))"
            @remove="smartMatchingStore.removeFromCompare"
            @select="selectTeacher"
            @exit="smartMatchingStore.exitCompareMode"
          />
        </div>

        <!-- 搜尋中載入進度 -->
        <div v-else-if="smartMatchingStore.isSearching && smartMatchingStore.searchResults.length === 0" class="py-12">
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
              <span>{{ smartMatchingStore.searchProgress }}%</span>
            </div>
            <div class="h-2 bg-white/10 rounded-full overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-primary-500 to-secondary-500 transition-all duration-300"
                :style="{ width: `${smartMatchingStore.searchProgress}%` }"
              ></div>
            </div>
          </div>

          <!-- 搜尋步驟 -->
          <div class="mt-6 space-y-2 max-w-sm mx-auto">
            <div
              v-for="(step, index) in smartMatchingStore.searchSteps"
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
            v-if="smartMatchingStore.searchResults.length > 0"
            v-model:sort-by="smartMatchingStore.sortBy"
            v-model:sort-order="smartMatchingStore.sortOrder"
            v-model:view-mode="smartMatchingStore.viewMode"
            :selected-count="smartMatchingStore.selectedCount"
          />

          <!-- 結果列表 -->
          <div v-if="smartMatchingStore.viewMode === 'card'" class="grid grid-cols-1 gap-3">
            <AdminEnhancedMatchCard
              v-for="match in smartMatchingStore.sortedResults"
              :key="match.teacher_id"
              :match="match"
              :selected="smartMatchingStore.selectedForCompare.has(match.teacher_id)"
              :show-checkbox="smartMatchingStore.sortedResults.length > 1"
              @update:selected="smartMatchingStore.toggleCompare(match.teacher_id, $event)"
              @click="viewTeacherTimeline(match)"
            />
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="match in smartMatchingStore.sortedResults"
              :key="match.teacher_id"
              :class="[
                'p-4 rounded-lg transition-all cursor-pointer',
                smartMatchingStore.selectedForCompare.has(match.teacher_id)
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
                    <div class="flex items-center gap-2">
                      <p class="text-sm text-slate-400">匹配度 {{ match.match_score }}%</p>
                      <ScoreBreakdownTooltip />
                    </div>
                  </div>
                </div>
                <div class="flex items-center gap-3">
                  <div class="text-right">
                    <div class="flex items-center gap-1">
                      <div class="text-2xl font-bold text-primary-500">{{ match.match_score }}%</div>
                      <ScoreBreakdownTooltip />
                    </div>
                  </div>
                  <input
                    v-if="smartMatchingStore.sortedResults.length > 1"
                    type="checkbox"
                    :checked="smartMatchingStore.selectedForCompare.has(match.teacher_id)"
                    @change="smartMatchingStore.toggleCompare(match.teacher_id, !smartMatchingStore.selectedForCompare.has(match.teacher_id))"
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
                  技能匹配: {{ match.score_detail?.match_score || 0 }}%
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
              <div v-if="(match as any).available_rooms?.length" class="mb-3">
                <p class="text-xs text-slate-500 mb-2">可授課教室：</p>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="room in (match as any).available_rooms"
                    :key="room.id"
                    class="px-2 py-1 rounded-md text-xs bg-success-500/20 text-success-500"
                  >
                    {{ room.name }}
                  </span>
                </div>
              </div>

              <p v-if="match?.notes" class="text-sm text-slate-400 mb-3">
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
          <div v-if="smartMatchingStore.selectedTeacher" class="mt-6">
            <!-- 衝突圖例 -->
            <AdminConflictLegend class="mb-4" />

            <!-- 時間軸 -->
            <AdminTeacherTimeline
              :teacher="smartMatchingStore.selectedTeacher"
              :target-start="form.start_time"
              :target-end="form.end_time"
              :existing-sessions="smartMatchingStore.teacherSessions"
              @override="confirmOverride"
            />

            <!-- 替代時段建議 -->
            <div v-if="smartMatchingStore.alternativeSlots.length > 0" class="mt-4">
              <AdminAlternativeSlots
                :teacher-name="smartMatchingStore.selectedTeacher.teacher_name"
                :slots="smartMatchingStore.alternativeSlots"
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
          v-if="smartMatchingStore.talentResults.length > 0"
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
        <div v-if="showStats && smartMatchingStore.talentResults.length > 0" class="mb-6">
          <!-- 統計卡片 -->
          <AdminTalentStatsPanel
            v-if="smartMatchingStore.talentStats"
            :stats="smartMatchingStore.talentStats"
            :city-distribution="smartMatchingStore.cityDistribution"
            :top-skills="smartMatchingStore.topSkills"
          />

          <!-- 技能分布圖 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
            <AdminSkillsDistributionChart
              :skills="skillDistribution"
              :total-teachers="smartMatchingStore.talentResults.length"
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
        v-if="smartMatchingStore.talentResults.length > 0"
        :popular-skills="popularSearchSkills"
        :popular-tags="['古典', '兒童', '成人', '入門', '進階']"
        :active-skills="activeSkillFilters"
        :active-tags="activeTagFilters"
        @filter-by-skill="onSkillFilter"
        @filter-by-tag="onTagFilter"
        @clear-all="clearQuickFilters"
      />

      <!-- 搜尋結果數量與分頁 -->
      <div v-if="smartMatchingStore.talentResults.length > 0" class="mt-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <span class="text-sm text-slate-400">
          找到 <span class="text-white font-medium">{{ smartMatchingStore.talentTotalItems }}</span> 位人才
          <span class="text-slate-500 ml-2">(第 {{ smartMatchingStore.talentCurrentPage }}/{{ smartMatchingStore.talentTotalPages }} 頁)</span>
        </span>

        <div class="flex items-center gap-2">
          <!-- 排序選單 -->
          <select
            v-model="smartMatchingStore.talentSortBy"
            @change="searchTalent(undefined, 1)"
            class="px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-white text-sm"
          >
            <option value="name">姓名</option>
            <option value="skills">技能數</option>
            <option value="rating">評分</option>
            <option value="city">地區</option>
          </select>

          <button
            @click="smartMatchingStore.talentSortOrder = smartMatchingStore.talentSortOrder === 'asc' ? 'desc' : 'asc'; searchTalent(undefined, 1)"
            class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-white transition-colors"
            title="切換排序順序"
          >
            <svg v-if="smartMatchingStore.talentSortOrder === 'asc'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h9m5-4v12m0 0l-4-4m4 4l4-4" />
            </svg>
          </button>
        </div>
      </div>

      <!-- 分頁控制項 -->
      <div v-if="smartMatchingStore.talentTotalPages > 1" class="mt-4 flex items-center justify-center gap-2">
        <button
          @click="goToPage(smartMatchingStore.talentCurrentPage - 1)"
          :disabled="smartMatchingStore.talentCurrentPage <= 1"
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
            :class="page === smartMatchingStore.talentCurrentPage
              ? 'bg-primary-500 text-white'
              : 'bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white'"
          >
            {{ page }}
          </button>
        </div>

        <button
          @click="goToPage(smartMatchingStore.talentCurrentPage + 1)"
          :disabled="smartMatchingStore.talentCurrentPage >= smartMatchingStore.talentTotalPages"
          class="px-3 py-1.5 rounded-lg bg-white/5 text-slate-400 hover:bg-white/10 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors text-sm"
        >
          下一頁
        </button>

        <span class="text-xs text-slate-500 ml-2">
          跳至第
          <input
            type="number"
            min="1"
            :max="smartMatchingStore.talentTotalPages"
            v-model.number="smartMatchingStore.talentCurrentPage"
            @keyup.enter="goToPage(smartMatchingStore.talentCurrentPage)"
            class="w-12 px-2 py-1 mx-1 rounded bg-white/5 border border-white/10 text-white text-center text-sm"
          />
          頁
        </span>
      </div>

      <!-- 搜尋結果 -->
      <div v-if="smartMatchingStore.talentResults.length > 0" class="mt-4">
        <!-- 批量操作工具列 -->
        <AdminBulkActions
          :selected-count="selectedTalents.size"
          :bulk-loading="smartMatchingStore.bulkLoading"
          :bulk-progress="smartMatchingStore.bulkProgress"
          @clear="clearTalentSelection"
          @compare="smartMatchingStore.viewMode = 'compare'"
          @export="exportContactInfo"
          @bulk-invite="bulkInviteTalents"
        />

        <!-- 人才卡片網格 -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-20">
          <AdminTalentCard
            v-for="teacher in smartMatchingStore.talentResults"
            :key="teacher.id"
            :teacher="teacher"
            :selected="selectedTalents.has(teacher.id)"
            :show-checkbox="teacher.is_open_to_hiring"
            :invite-loading="smartMatchingStore.inviteLoadingIds.has(teacher.id)"
            :invitation-status="smartMatchingStore.invitationStatuses.get(teacher.id) || null"
            @update:selected="toggleTalentSelection(teacher.id, $event)"
            @invite="inviteSingleTalent"
            @view="viewTalentDetail"
            @compare="addTalentToCompare"
          />
        </div>
      </div>

      <!-- 空狀態 -->
      <div v-else-if="smartMatchingStore.hasSearchedTalent" class="text-center py-12 text-slate-500">
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
import ScoreBreakdownTooltip from '~/components/Scheduling/ScoreBreakdownTooltip.vue'
import type { SmartMatchingResult, TalentCard } from '~/types/matching'

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

definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

 const notificationUI = useNotification()
 const { warning: alertWarning, success: alertSuccess, error: alertError, confirm: alertConfirm } = useAlert()

 // 使用 Smart Matching Store
 const smartMatchingStore = useSmartMatchingStore()

 // 搜尋表單
 const form = ref({
   start_time: '',
   end_time: '',
   room_ids: [] as number[],
   skills: '' as string
 })

 // 人才搜尋表單
 const talentSearch = ref({
   query: '',
   city: '',
   skills: '',
   hashtags: ''
 })

 // 選中的技能列表
 const selectedSkills = ref<SelectedSkill[]>([])

 // 房間列表
 const rooms = ref<any[]>([])
 const { getCenterId } = useCenterId()

 // 近期搜尋組件引用
 const recentSearchesRef = ref<InstanceType<typeof AdminRecentSearches> | null>(null)

 // 人才庫管理狀態
 const selectedTalents = ref<Set<number>>(new Set())
 const showStats = ref(false)
 const showCustomTime = ref(false)

 // 篩選狀態
 const activeSkillFilters = ref<string[]>([])
 const activeTagFilters = ref<string[]>([])

 // 人才篩選面板引用
 const talentFilterRef = ref<InstanceType<typeof AdminTalentFilterPanel> | null>(null)

 // 熱門搜尋技能
 const popularSearchSkills = computed(() => {
   return smartMatchingStore.topSkills.slice(0, 6).map((s: any) => s.name)
 })

 // 技能分布
 const skillDistribution = computed(() => {
   return smartMatchingStore.topSkills.slice(0, 8)
 })

 // 計算可見的頁碼
 const visiblePages = computed(() => {
   const pages: number[] = []
   const maxVisible = 5
   const totalPages = smartMatchingStore.talentTotalPages
   const currentPage = smartMatchingStore.talentCurrentPage

   if (totalPages <= maxVisible) {
     for (let i = 1; i <= totalPages; i++) {
       pages.push(i)
     }
   } else {
     let start = Math.max(1, currentPage - Math.floor(maxVisible / 2))
     let end = start + maxVisible - 1

     if (end > totalPages) {
       end = totalPages
       start = Math.max(1, end - maxVisible + 1)
     }

     for (let i = start; i <= end; i++) {
       pages.push(i)
     }
   }

   return pages
 })

 // ==================== 搜尋相關方法 ====================

 /**
  * 智慧媒合搜尋
  */
 const findMatches = async () => {
   if (!form.value.start_time || !form.value.end_time) {
     await alertWarning('請填寫開始時間和結束時間')
     return
   }

   const skills = selectedSkills.value.length > 0
     ? selectedSkills.value.map(s => s.name)
     : form.value.skills.split(',').map(s => s.trim()).filter(Boolean)

   await smartMatchingStore.searchMatches({
     room_ids: form.value.room_ids.length > 0 ? form.value.room_ids : undefined,
     start_time: formatDateToString(new Date(form.value.start_time)),
     end_time: formatDateToString(new Date(form.value.end_time)),
     required_skills: skills.length > 0 ? skills : undefined,
   })

   if (smartMatchingStore.searchResults.length === 0) {
     await alertWarning('沒有找到符合條件的老師')
   } else {
     // 儲存到近期搜尋
     if (recentSearchesRef.value) {
       recentSearchesRef.value.addSearch({
         start_time: formatDateToString(new Date(form.value.start_time)),
         end_time: formatDateToString(new Date(form.value.end_time)),
         room_ids: [...form.value.room_ids],
         skills: [...selectedSkills.value]
       })
     }
   }
 }

 /**
  * 清除搜尋表單
  */
 const clearForm = () => {
   form.value = {
     start_time: '',
     end_time: '',
     room_ids: [],
     skills: ''
   }
   selectedSkills.value = []
   smartMatchingStore.resetSearchResults()
 }

 /**
  * 從近期搜尋載入
  */
 const onLoadRecentSearch = (search: RecentSearch) => {
   form.value.start_time = search.start_time.slice(0, 16)
   form.value.end_time = search.end_time.slice(0, 16)
   form.value.room_ids = [...search.room_ids]
   selectedSkills.value = [...search.skills]
 }

 const onClearRecentSearches = () => {
   // 清除時不執行額外動作
 }

 /**
  * 查看教師課表時間軸
  */
 const viewTeacherTimeline = async (teacher: SmartMatchingResult) => {
   smartMatchingStore.selectTeacher(teacher)

   if (!form.value.start_time || !form.value.end_time) {
     return
   }

   const startDate = formatDateToString(new Date(form.value.start_time))
   const endDate = formatDateToString(new Date(form.value.end_time))

   await smartMatchingStore.fetchTeacherSchedule({
     teacher_id: teacher.teacher_id,
     start_date: startDate,
     end_date: endDate,
   })

   // 取得替代時段建議
   await smartMatchingStore.fetchAlternativeSlots(
     teacher.teacher_id,
     formatDateToString(new Date(form.value.start_time)),
     formatDateToString(new Date(form.value.start_time)),
     formatDateToString(new Date(form.value.end_time)),
     90
   )

   await alertSuccess(`正在查看 ${teacher.teacher_name} 的課表`)
 }

 /**
  * 選擇替代時段
  */
 const selectAlternativeSlot = (slot: any) => {
   form.value.start_time = `${slot.date}T${slot.start}`
   form.value.end_time = `${slot.date}T${slot.end}`

   alertSuccess(`已選擇替代時段 ${slot.dateLabel} ${slot.start}-${slot.end}`)

   // 重新搜尋
   findMatches()
 }

 /**
  * 確認 Override
  */
 const confirmOverride = async () => {
   if (await alertConfirm('確定要在緩衝衝突的情況下安排這位老師嗎？')) {
     if (smartMatchingStore.selectedTeacher) {
       await alertSuccess(`已選擇 ${smartMatchingStore.selectedTeacher.teacher_name}`)
     }
   }
 }

 /**
  * 選取教師
  */
 const selectTeacher = async (teacher: SmartMatchingResult) => {
   await alertSuccess(`已選擇 ${teacher.teacher_name}`)
 }

 // ==================== 人才庫相關方法 ====================

 /**
  * 人才庫搜尋
  */
 const searchTalent = async (query?: string, page: number = 1) => {
   await smartMatchingStore.searchTalent({
     keyword: query || talentSearch.value.query || undefined,
     city: talentSearch.value.city || undefined,
     skill_name: talentSearch.value.skills || undefined,
     district: undefined,
     open_to_hiring_only: undefined,
     min_rating: undefined,
     sort_by: smartMatchingStore.talentSortBy as any,
     sort_order: smartMatchingStore.talentSortOrder,
     page,
     limit: 20,
   })

   showStats.value = true
 }

 /**
  * 切換分頁
  */
 const goToPage = (page: number) => {
   if (page >= 1 && page <= smartMatchingStore.talentTotalPages) {
     searchTalent(undefined, page)
     window.scrollTo({ top: 0, behavior: 'smooth' })
   }
 }

 /**
  * 搜尋建議選取
  */
 const onSuggestionSelected = (suggestion: { type: string; value: string }) => {
   talentSearch.value.query = suggestion.value
   if (suggestion.type === 'skill') {
     talentSearch.value.skills = suggestion.value
   } else if (suggestion.type === 'tag') {
     talentSearch.value.hashtags = suggestion.value
   }
   searchTalent()
 }

 /**
  * 技能選取
  */
 const onSkillSelected = (skillName: string) => {
   talentSearch.value.skills = skillName
   activeSkillFilters.value = [skillName]
   searchTalent()
 }

 /**
  * 技能篩選
  */
 const onSkillFilter = (skill: string) => {
   if (activeSkillFilters.value.includes(skill)) {
     activeSkillFilters.value = activeSkillFilters.value.filter(s => s !== skill)
   } else {
     activeSkillFilters.value = [...activeSkillFilters.value, skill]
   }
 }

 /**
  * 標籤篩選
  */
 const onTagFilter = (tag: string) => {
   if (activeTagFilters.value.includes(tag)) {
     activeTagFilters.value = activeTagFilters.value.filter(t => t !== tag)
   } else {
     activeTagFilters.value = [...activeTagFilters.value, tag]
   }
 }

 /**
  * 清除快速篩選
  */
 const clearQuickFilters = () => {
   activeSkillFilters.value = []
   activeTagFilters.value = []
 }

 /**
  * 人才篩選面板 - 套用篩選
  */
 const applyTalentFilters = (filters: any) => {
   if (filters.city) talentSearch.value.city = filters.city
   if (filters.skills) talentSearch.value.skills = filters.skills
   if (filters.hashtags) talentSearch.value.hashtags = filters.hashtags
   searchTalent()
 }

 /**
  * 人才篩選面板 - 清除篩選
  */
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

 /**
  * 清除人才搜尋
  */
 const clearTalentSearch = () => {
   talentSearch.value = {
     query: '',
     city: '',
     skills: '',
     hashtags: ''
   }
   smartMatchingStore.resetTalentSearch()
   clearQuickFilters()
   showStats.value = false
 }

 // ==================== 人才邀請相關方法 ====================

 /**
  * 邀請單一人才
  */
 const inviteSingleTalent = async (teacher: TalentCard) => {
   const success = await smartMatchingStore.inviteTalent(teacher.id)
   if (success) {
     await alertSuccess(`已向 ${teacher.name} 發送合作邀請`)
   } else {
     await alertWarning('邀請失敗，請稍後再試')
   }
 }

 /**
  * 批量邀請人才
  */
 const bulkInviteTalents = async () => {
   if (selectedTalents.value.size === 0) {
     await alertWarning('請先選取要邀請的人才')
     return
   }

   if (!(await alertConfirm(`確定要邀請選取的 ${selectedTalents.value.size} 位人才嗎？`))) {
     return
   }

   const result = await smartMatchingStore.bulkInviteTalents(
     Array.from(selectedTalents.value)
   )

   await alertSuccess(`已成功邀請 ${result.success} 位人才，失敗 ${result.failed} 位`)
   selectedTalents.value.clear()
 }

 /**
  * 切換人才選取
  */
 const toggleTalentSelection = (teacherId: number, selected: boolean) => {
   if (selected) {
     selectedTalents.value.add(teacherId)
   } else {
     selectedTalents.value.delete(teacherId)
   }
 }

 /**
  * 清除人才選取
  */
 const clearTalentSelection = () => {
   selectedTalents.value.clear()
 }

 /**
  * 匯出聯絡資訊
  */
 const exportContactInfo = () => {
   const selectedData = smartMatchingStore.talentResults.filter(t => selectedTalents.value.has(t.id))

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
     (t as any).public_contact_info || '',
     t.skills?.map((s: any) => s.name).join(', ') || '',
     (t as any).personal_hashtags?.join(', ') || ''
   ])

   const csv = [headers.join(','), ...rows.map(r => r.map((cell: string) => `"${cell}"`).join(','))].join('\n')

   // 下載檔案
   const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
   const link = document.createElement('a')
   link.href = URL.createObjectURL(blob)
   link.download = `人才庫匯出_${getTodayString()}.csv`
   link.click()
   URL.revokeObjectURL(link.href)

   alertSuccess(`已匯出 ${selectedData.length} 位人才的聯絡資訊`)
 }

 /**
  * 查看人才詳細資訊
  */
 const viewTalentDetail = (teacher: TalentCard) => {
   alertSuccess(`正在查看 ${teacher.name} 的詳細資訊`)
 }

 /**
  * 加入比較
  */
 const addTalentToCompare = (teacher: TalentCard) => {
   if (smartMatchingStore.selectedForCompare.size >= 3) {
     alertWarning('最多只能選取 3 位老師進行比較')
     return
   }
   smartMatchingStore.selectedForCompare.add(teacher.id)
   alertSuccess(`已將 ${teacher.name} 加入比較`)
 }

 // ==================== 取得房間列表 ====================

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

 // ==================== 生命週期 ====================

 onMounted(async () => {
   fetchRooms()
   if (recentSearchesRef.value) {
     recentSearchesRef.value.loadFromStorage()
   }

   // 載入人才庫統計
   await smartMatchingStore.fetchTalentStats()
 })
</script>
