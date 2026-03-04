<template>
  <div class="p-4 md:p-6">
    <!-- 操作指南區域 -->
    <div class="mb-6">
      <button
        @click="showGuide = !showGuide"
        class="flex items-center gap-2 text-slate-300 hover:text-white transition-colors"
      >
        <svg class="w-5 h-5" :class="showGuide ? 'rotate-90' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
        <span class="font-medium">操作指南</span>
      </button>

      <div
        v-if="showGuide"
        class="mt-4 p-4 bg-slate-800/50 rounded-lg border border-white/10"
      >
        <h4 class="text-white font-medium mb-3">課表模板使用流程</h4>
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <!-- 步驟 1 -->
          <div class="flex gap-3">
            <div class="flex-shrink-0 w-8 h-8 rounded-full bg-primary-500/20 flex items-center justify-center text-primary-500 font-bold">1</div>
            <div>
              <h5 class="text-white font-medium text-sm">建立模板</h5>
              <p class="text-slate-400 text-xs mt-1">選擇模板名稱和視角類型（教室/老師）</p>
            </div>
          </div>
          <!-- 步驟 2 -->
          <div class="flex gap-3">
            <div class="flex-shrink-0 w-8 h-8 rounded-full bg-primary-500/20 flex items-center justify-center text-primary-500 font-bold">2</div>
            <div>
              <h5 class="text-white font-medium text-sm">新增格子</h5>
              <p class="text-slate-400 text-xs mt-1">設定時間段並綁定教室或老師</p>
            </div>
          </div>
          <!-- 步驟 3 -->
          <div class="flex gap-3">
            <div class="flex-shrink-0 w-8 h-8 rounded-full bg-primary-500/20 flex items-center justify-center text-primary-500 font-bold">3</div>
            <div>
              <h5 class="text-white font-medium text-sm">選擇課程</h5>
              <p class="text-slate-400 text-xs mt-1">為模板選擇要套用的課程班別</p>
            </div>
          </div>
          <!-- 步驟 4 -->
          <div class="flex gap-3">
            <div class="flex-shrink-0 w-8 h-8 rounded-full bg-primary-500/20 flex items-center justify-center text-primary-500 font-bold">4</div>
            <div>
              <h5 class="text-white font-medium text-sm">套用模板</h5>
              <p class="text-slate-400 text-xs mt-1">選擇日期範圍和星期，自動生成課表</p>
            </div>
          </div>
        </div>

        <!-- 視角說明 -->
        <div class="mt-4 p-3 bg-blue-500/10 rounded-lg">
          <h5 class="text-blue-400 text-sm font-medium mb-2">視角類型說明</h5>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs text-slate-300">
            <div class="flex items-start gap-2">
              <span class="text-primary-500">●</span>
              <span><strong class="text-white">教室視角：</strong>按教室分配時間格，適合固定教室的課程排班</span>
            </div>
            <div class="flex items-start gap-2">
              <span class="text-secondary-500">●</span>
              <span><strong class="text-white">老師視角：</strong>按老師分配時間格，適合指名老師的課程排班</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold text-white">課表模板</h1>
      <button
        @click="showModal = true"
        class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新增模板
      </button>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="template in templates"
        :key="template.id"
        class="glass-card p-4"
      >
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 mb-3">
          <h3 class="text-lg font-medium text-white">{{ template.name }}</h3>
          <span
            class="px-2 py-1 rounded-full text-xs w-fit"
            :class="template.row_type === 'ROOM' ? 'bg-primary-500/20 text-primary-500' : 'bg-secondary-500/20 text-secondary-500'"
          >
            {{ template.row_type === 'ROOM' ? '教室視角' : '老師視角' }}
          </span>
        </div>

        <div class="flex items-center justify-between text-sm text-slate-400 mb-4">
          <span>建立於 {{ formatDate(template.created_at) }}</span>
          <span>{{ template.is_active !== false ? '啟用' : '停用' }}</span>
        </div>

        <div class="flex gap-2">
          <button
            @click="viewTemplate(template)"
            class="flex-1 px-3 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors text-sm flex items-center justify-center gap-1"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            查看格子
          </button>
          <button
            @click="openApplyModal(template)"
            class="px-3 py-2 rounded-lg bg-primary-500/20 text-primary-500 hover:bg-primary-500/30 transition-colors text-sm flex items-center justify-center gap-1"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            套用
          </button>
          <button
            @click="deleteTemplate(template.id)"
            class="px-3 py-2 rounded-lg bg-critical-500/20 text-critical-500 hover:bg-critical-500/30 transition-colors text-sm"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>
      </div>

      <div
        v-if="templates.length === 0"
        class="col-span-full text-center py-16"
      >
        <div class="mb-4">
          <svg class="w-16 h-16 mx-auto text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-white mb-2">尚未建立課表模板</h3>
        <p class="text-slate-400 mb-4">建立模板可以快速重複套用相同的課表結構</p>
        <button
          @click="showModal = true"
          class="px-6 py-3 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors inline-flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          建立第一個模板
        </button>
      </div>
    </div>
  </div>

  <div
    v-if="showModal"
    class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    @click.self="showModal = false"
  >
    <div class="glass-card w-full max-w-md p-6">
      <h3 class="text-lg font-semibold text-white mb-2">新增模板</h3>
      <p class="text-sm text-slate-400 mb-4">建立課表模板，快速重複套用相同的排課結構</p>

      <div class="mb-4">
        <label for="template-name" class="block text-slate-300 mb-2">
          模板名稱
          <span class="text-slate-500 text-xs ml-2">必填</span>
        </label>
        <input
          id="template-name"
          v-model="form.name"
          type="text"
          placeholder="例如：週一上午班表、瑜珈教室模板"
          class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white placeholder-slate-500"
          required
        />
      </div>
      <div class="mb-4">
        <label for="template-type" class="block text-slate-300 mb-2">
          視角類型
          <span class="text-slate-500 text-xs ml-2">必填</span>
        </label>
        <div class="grid grid-cols-2 gap-3 mb-2">
          <label
            class="flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-colors"
            :class="form.row_type === 'ROOM' ? 'bg-primary-500/20 border-primary-500' : 'bg-white/5 border-white/10 hover:bg-white/10'"
          >
            <input
              type="radio"
              v-model="form.row_type"
              value="ROOM"
              class="sr-only"
            />
            <div class="flex-1">
              <div class="text-white font-medium text-sm">教室視角</div>
              <div class="text-slate-400 text-xs mt-0.5">按教室分配時間</div>
            </div>
            <svg v-if="form.row_type === 'ROOM'" class="w-5 h-5 text-primary-500" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
          </label>
          <label
            class="flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-colors"
            :class="form.row_type === 'TEACHER' ? 'bg-secondary-500/20 border-secondary-500' : 'bg-white/5 border-white/10 hover:bg-white/10'"
          >
            <input
              type="radio"
              v-model="form.row_type"
              value="TEACHER"
              class="sr-only"
            />
            <div class="flex-1">
              <div class="text-white font-medium text-sm">老師視角</div>
              <div class="text-slate-400 text-xs mt-0.5">按老師分配時間</div>
            </div>
            <svg v-if="form.row_type === 'TEACHER'" class="w-5 h-5 text-secondary-500" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
          </label>
        </div>
        <p class="text-xs text-slate-500">
          {{ form.row_type === 'ROOM' ? '建立後需為每個格子選擇教室' : '建立後需為每個格子選擇老師' }}
        </p>
      </div>
      <div class="flex gap-3">
        <button
          type="button"
          @click="showModal = false"
          class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
        >
          取消
        </button>
        <button
          type="submit"
          @click="createTemplate"
          :disabled="creating || !form.name"
          class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ creating ? '建立中...' : '建立模板' }}
        </button>
      </div>
    </div>
  </div>

  <div
    v-if="selectedTemplate"
    class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    @click.self="selectedTemplate = null"
  >
    <div class="glass-card w-full max-w-2xl p-6 max-h-[80vh] overflow-y-auto">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h3 class="text-lg font-semibold text-white">{{ selectedTemplate.name }}</h3>
          <p class="text-sm text-slate-400 mt-1">
            {{ selectedTemplate.row_type === 'ROOM' ? '教室視角模板' : '老師視角模板' }}
            · {{ cells.length }} 個時間格
          </p>
        </div>
        <button @click="selectedTemplate = null" class="text-slate-400 hover:text-white">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- 視角提示 -->
      <div class="mb-4 p-3 rounded-lg" :class="selectedTemplate.row_type === 'ROOM' ? 'bg-primary-500/10 border border-primary-500/30' : 'bg-secondary-500/10 border border-secondary-500/30'">
        <div class="flex items-start gap-2">
          <svg class="w-5 h-5 mt-0.5 flex-shrink-0" :class="selectedTemplate.row_type === 'ROOM' ? 'text-primary-500' : 'text-secondary-500'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <p class="text-sm font-medium" :class="selectedTemplate.row_type === 'ROOM' ? 'text-primary-400' : 'text-secondary-400'">
              {{ selectedTemplate.row_type === 'ROOM' ? '教室視角' : '老師視角' }}操作說明
            </p>
            <p class="text-xs text-slate-400 mt-1">
              {{ selectedTemplate.row_type === 'ROOM'
                ? '為每個時間格選擇教室。系統會檢查教室時間衝突與緩衝時間。'
                : '為每個時間格選擇老師。系統會檢查老師時間衝突、緩衝時間與私人行程。' }}
            </p>
          </div>
        </div>
      </div>

      <div v-if="cells.length === 0" class="text-center py-12">
        <div class="mb-4">
          <svg class="w-16 h-16 mx-auto text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
        </div>
        <h4 class="text-white font-medium mb-2">尚未建立時間格</h4>
        <p class="text-slate-400 text-sm mb-4">新增時間格來定義課表結構</p>
        <div class="flex items-center justify-center gap-4 text-xs text-slate-500">
          <span class="flex items-center gap-1">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
            列：時段（如 9:00-10:00）
          </span>
          <span class="flex items-center gap-1">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2" />
            </svg>
            行：星期（週一～週日）
          </span>
        </div>
      </div>

      <div v-else>
        <!-- 拖曳排序提示 -->
        <div class="mb-3 flex items-center gap-2 text-xs text-slate-500">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16" />
          </svg>
          <span>拖曳格子可調整順序（先後順序影響套用時的預設排序）</span>
        </div>

        <!-- 格子列表（支援拖曳） -->
        <div class="space-y-2">
          <div
            v-for="(cell, index) in cells"
            :key="cell.id"
            class="p-3 rounded-lg transition-all cursor-move"
            :class="[
              draggedCell?.id === cell.id ? 'opacity-50 scale-102 bg-primary-500/20 border-2 border-primary-500' :
              dragOverCell?.id === cell.id ? 'bg-primary-500/10 border-2 border-primary-500 transform scale-[1.02]' :
              'bg-white/5 hover:bg-white/10 border-2 border-transparent'
            ]"
            draggable="true"
            @dragstart="handleDragStart($event, cell)"
            @dragend="handleDragEnd"
            @dragover.prevent="handleDragOver($event, cell)"
            @dragleave="handleDragLeave"
            @drop="handleDrop($event, cell, index)"
          >
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
              <div class="flex items-center gap-3">
                <!-- 拖曳把手圖示 -->
                <div class="flex-shrink-0 w-8 h-8 flex items-center justify-center text-slate-500">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h4m-4 8h4m-4 8h4m-4 8h4M6 4h12v4H6z" />
                  </svg>
                </div>
                <span class="inline-flex items-center justify-center w-10 h-10 text-xs bg-white/10 rounded-lg text-slate-300 font-mono">
                  {{ cell.row_no }}-{{ cell.col_no }}
                </span>
                <div>
                  <div class="text-white font-medium">{{ cell.start_time }} - {{ cell.end_time }}</div>
                  <div class="text-xs text-slate-500">
                    {{ selectedTemplate.row_type === 'ROOM' ? '教室' : '老師' }}：
                    <span :class="getCellResourceClass(cell)">
                      {{ getCellResourceName(cell) }}
                    </span>
                  </div>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-xs text-slate-600 px-2 py-1 bg-white/5 rounded">
                  #{{ index + 1 }}
                </span>
                <button
                  @click="deleteCell(cell.id)"
                  class="p-2 rounded-lg bg-critical-500/10 text-critical-500 hover:bg-critical-500/20 transition-colors"
                  title="刪除格子"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 新增格子表單 -->
      <div class="mt-4 pt-4 border-t border-white/10">
        <h4 class="text-sm font-medium text-white mb-3">新增時間格</h4>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-3">
          <div>
            <label class="block text-slate-400 text-xs mb-1">列（時段）</label>
            <input
              v-model.number="newCell.row_no"
              type="number"
              min="1"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white text-sm"
            />
          </div>
          <div>
            <label class="block text-slate-400 text-xs mb-1">行（星期）</label>
            <input
              v-model.number="newCell.col_no"
              type="number"
              min="1"
              max="7"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white text-sm"
            />
          </div>
          <div>
            <label class="block text-slate-400 text-xs mb-1">開始時間</label>
            <input
              v-model="newCell.start_time"
              type="time"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white text-sm"
            />
          </div>
          <div>
            <label class="block text-slate-400 text-xs mb-1">結束時間</label>
            <input
              v-model="newCell.end_time"
              type="time"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white text-sm"
            />
          </div>
        </div>

        <!-- 教室選擇（教室視角） -->
        <div v-if="selectedTemplate.row_type === 'ROOM'" class="mb-3">
          <label class="block text-slate-400 text-xs mb-1">
            選擇教室
            <span class="text-red-400 ml-1">*必選</span>
          </label>
          <select
            v-model="newCell.room_id"
            class="w-full px-3 py-2 rounded-lg bg-slate-800 border border-white/10 text-white text-sm cursor-pointer appearance-none"
          >
            <option value="">請選擇教室</option>
            <option v-for="room in rooms" :key="room.id" :value="room.id">
              {{ room.name }}
            </option>
          </select>
          <p v-if="rooms.length === 0" class="text-xs text-yellow-500 mt-1">尚無教室資料，請先至「資源管理」建立教室</p>
        </div>

        <!-- 老師選擇（老師視角） -->
        <div v-else class="mb-3">
          <SearchableSelect
            v-model="newCell.teacher_id"
            :options="teacherOptions"
            label="選擇老師"
            placeholder="請選擇老師"
            required
          />
          <p v-if="teachers.length === 0" class="text-xs text-yellow-500 mt-1">尚無老師資料，請先邀請老師加入中心</p>
        </div>

        <div class="flex gap-2">
          <button
            @click="addCells"
            :disabled="addingCell || !canAddCell"
            class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            {{ addingCell ? '新增中...' : '新增時間格' }}
          </button>
        </div>
      </div>
    </div>
  </div>

  <div
    v-if="showApplyModal"
    class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    @click.self="closeApplyModal"
  >
    <div class="glass-card w-full max-w-lg p-6 max-h-[80vh] overflow-y-auto">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h3 class="text-lg font-semibold text-white">套用模板</h3>
          <p class="text-sm text-slate-400 mt-1">{{ applyForm.templateName }}</p>
        </div>
        <button @click="closeApplyModal" class="text-slate-400 hover:text-white">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- 套用步驟指示 -->
      <div class="mb-4 flex items-center gap-2 text-xs">
        <span class="flex items-center gap-1 text-slate-400">
          <span class="w-5 h-5 rounded-full bg-primary-500 flex items-center justify-center text-white">1</span>
          選擇課程
        </span>
        <svg class="w-4 h-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
        <span class="flex items-center gap-1 text-slate-400">
          <span class="w-5 h-5 rounded-full bg-primary-500 flex items-center justify-center text-white">2</span>
          設定日期
        </span>
        <svg class="w-4 h-4 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
        <span class="flex items-center gap-1 text-white font-medium">
          <span class="w-5 h-5 rounded-full bg-primary-500 flex items-center justify-center text-white">3</span>
          確認套用
        </span>
      </div>

      <!-- 衝突警告區域 -->
      <div
        v-if="applyConflicts.length > 0 || validationWarnings.length > 0"
        class="mb-4 p-4 bg-yellow-500/10 border border-yellow-500/30 rounded-lg"
      >
        <div class="flex items-center gap-2 mb-2">
          <svg class="w-5 h-5 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <span class="text-yellow-500 font-medium">
            檢測到 {{ applyConflicts.length }} 個衝突、{{ validationWarnings.length }} 個警告
          </span>
        </div>

        <!-- 衝突列表 -->
        <div v-if="applyConflicts.length > 0" class="space-y-2 mb-3 max-h-48 overflow-y-auto">
          <div
            v-for="(conflict, index) in applyConflicts"
            :key="'conflict-' + index"
            class="text-sm p-2 rounded"
            :class="conflict.can_override ? 'bg-yellow-500/10' : 'bg-red-500/10'"
          >
            <div class="font-medium" :class="conflict.can_override ? 'text-yellow-400' : 'text-red-400'">
              {{ getConflictTypeLabel(conflict.conflict_type || conflict.type) }}
            </div>
            <div class="text-slate-400 mt-1 text-xs">{{ conflict.message }}</div>
          </div>
        </div>

        <!-- 警告列表 -->
        <div v-if="validationWarnings.length > 0" class="space-y-2 max-h-48 overflow-y-auto">
          <div
            v-for="(warning, index) in validationWarnings"
            :key="'warning-' + index"
            class="text-sm p-2 rounded bg-blue-500/10"
          >
            <div class="font-medium text-blue-400">
              {{ getConflictTypeLabel(warning.warning_type || warning.type) }}
            </div>
            <div class="text-slate-400 mt-1 text-xs">{{ warning.message }}</div>
          </div>
        </div>

        <div v-if="hasOverrideableConflicts" class="mt-3 p-2 bg-yellow-500/10 rounded text-xs text-yellow-400">
          💡 提示：可勾選「允許覆蓋 Buffer 衝突」來強制套用
        </div>
      </div>

      <form @submit.prevent="applyTemplate">
        <div class="mb-4">
          <label for="offering-select" class="block text-slate-300 mb-2 font-medium">
            選擇課程
            <span class="text-red-400 ml-1">*必選</span>
          </label>
          <div class="relative">
            <select
              id="offering-select"
              v-model="applyForm.offeringId"
              class="input-field"
              :disabled="offeringsLoading"
              required
            >
              <option value="">請選擇課程 ({{ offerings.length }} 筆資料)</option>
              <option v-for="offering in offerings" :key="offering.id" :value="offering.id">
                {{ offering.name }}
              </option>
            </select>
            <!-- 載入指示器 -->
            <div v-if="offeringsLoading" class="absolute right-10 top-1/2 -translate-y-1/2">
              <svg class="w-5 h-5 text-primary-500 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
            </div>
          </div>
          <p v-if="offeringsLoading" class="text-xs text-slate-500 mt-1">載入課程資料中...</p>
          <p v-else-if="offerings.length === 0" class="text-xs text-yellow-500 mt-1">尚無課程資料，請先至「課程管理」建立課程</p>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-4">
          <div>
            <label for="start-date" class="block text-slate-300 mb-2">
              開始日期
              <span class="text-red-400 ml-1">*必填</span>
            </label>
            <input
              id="start-date"
              v-model="applyForm.startDate"
              type="date"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
              required
            />
          </div>
          <div>
            <label for="end-date" class="block text-slate-300 mb-2">
              結束日期
              <span class="text-red-400 ml-1">*必填</span>
            </label>
            <input
              id="end-date"
              v-model="applyForm.endDate"
              type="date"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
              required
            />
          </div>
        </div>

        <div class="mb-4">
          <span class="block text-slate-300 mb-2">
            選擇星期
            <span class="text-red-400 ml-1">*必選</span>
          </span>
          <div class="flex flex-wrap gap-2" role="group" aria-label="選擇星期">
            <label
              v-for="day in weekdays"
              :key="day.value"
              class="flex items-center gap-2 px-3 py-2 rounded-lg bg-white/5 cursor-pointer hover:bg-white/10 transition-colors"
              :class="{ 'bg-primary-500/20 border border-primary-500': applyForm.weekdays.includes(day.value) }"
            >
              <input
                type="checkbox"
                :value="day.value"
                v-model="applyForm.weekdays"
                class="w-4 h-4 rounded border-white/20 bg-white/10 text-primary-500"
              />
              <span class="text-white text-sm">{{ day.label }}</span>
            </label>
          </div>
          <p class="text-xs text-slate-500 mt-1">系統會在選定的星期自動套用模板中的時間格</p>
        </div>

        <div class="mb-4">
          <label for="duration-input" class="block text-slate-300 mb-2">
            每堂課時長（分鐘）
            <span class="text-slate-500 ml-1">預設 60 分鐘</span>
          </label>
          <input
            id="duration-input"
            v-model.number="applyForm.duration"
            type="number"
            min="1"
            class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white"
          />
        </div>

        <!-- Buffer Override 選項 -->
        <div class="mb-4 p-3 bg-white/5 rounded-lg border border-white/10">
          <label class="flex items-start gap-3 cursor-pointer">
            <input
              type="checkbox"
              v-model="applyForm.override_buffer"
              class="w-5 h-5 mt-0.5 rounded border-white/20 bg-white/10 text-primary-500"
            />
            <div>
              <span class="text-white font-medium">允許覆蓋 Buffer 衝突</span>
              <p class="text-xs text-slate-400 mt-1">
                勾選後，即使時間間隔不足緩衝時間也會強制套用模板。<br />
                適用於緊急排課或特殊情況。
              </p>
            </div>
          </label>
        </div>

        <div class="flex gap-3">
          <button
            type="button"
            @click="closeApplyModal"
            class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="applying || isValidating || !canApply"
            class="flex-1 px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg v-if="isValidating" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ isValidating ? '檢查中...' : applying ? '套用中...' : '確認套用' }}
          </button>
        </div>
      </form>
    </div>
  </div>

  <NotificationDropdown
    v-if="notificationUI.show.value"
    @close="notificationUI.close()"
  />
</template>

<script setup lang="ts">
import NotificationDropdown from '~/components/Navigation/NotificationDropdown.vue'
import SearchableSelect, { type SelectOption } from '~/components/Common/SearchableSelect.vue'
definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

 const notificationUI = useNotification()
const showModal = ref(false)
const showApplyModal = ref(false)
const showGuide = ref(false) // 是否顯示操作指南
const selectedTemplate = ref<any>(null)
const templates = ref<any[]>([])
const cells = ref<any[]>([])
const offerings = ref<any[]>([])
const offeringsLoading = ref(true)
const rooms = ref<any[]>([])
const teachers = ref<any[]>([])
const applyConflicts = ref<any[]>([])
const validationWarnings = ref<any[]>([])
const isValidating = ref(false)
const creating = ref(false)
const applying = ref(false)
const addingCell = ref(false)
const { confirm: alertConfirm, error: alertError, warning: alertWarning, success: alertSuccess } = useAlert()

// 拖曳排序相關狀態
const draggedCell = ref<any>(null)
const dragOverCell = ref<any>(null)
const isReordering = ref(false)

// 開始拖曳
const handleDragStart = (event: DragEvent, cell: any) => {
  draggedCell.value = cell
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.dropEffect = 'move'
  }
}

// 結束拖曳
const handleDragEnd = () => {
  draggedCell.value = null
  dragOverCell.value = null
}

// 拖曳經過
const handleDragOver = (event: DragEvent, cell: any) => {
  if (draggedCell.value && draggedCell.value.id !== cell.id) {
    dragOverCell.value = cell
  }
}

// 拖曳離開
const handleDragLeave = () => {
  dragOverCell.value = null
}

// 放置
const handleDrop = async (event: DragEvent, targetCell: any, targetIndex: number) => {
  event.preventDefault()

  if (!draggedCell.value || draggedCell.value.id === targetCell.id) {
    handleDragEnd()
    return
  }

  // 重新排序
  const sourceIndex = cells.value.findIndex(c => c.id === draggedCell.value.id)
  if (sourceIndex === -1) {
    handleDragEnd()
    return
  }

  // 更新順序
  const newCells = [...cells.value]
  const [removed] = newCells.splice(sourceIndex, 1)
  newCells.splice(targetIndex, 0, removed)

  // 更新 cells 並重新計算排序
  cells.value = newCells.map((cell, index) => ({
    ...cell,
    sort_order: index + 1
  }))

  handleDragEnd()

  // 儲存新的排序到後端
  await saveCellOrder()
}

const saveCellOrder = async () => {
  if (!selectedTemplate.value || cells.value.length === 0) return

  isReordering.value = true
  try {
    const api = useApi()
    const orderData = cells.value.map((cell, index) => ({
      id: cell.id,
      sort_order: index + 1
    }))
    await api.put(`/admin/templates/${selectedTemplate.value.id}/cells/reorder`, { cells: orderData })
  } catch (error) {
    console.error('Failed to save cell order:', error)
    await alertError('排序儲存失敗，請重新整理頁面')
  } finally {
    isReordering.value = false
  }
}

const newCell = ref({
  row_no: 1,
  col_no: 1,
  start_time: '09:00',
  end_time: '10:00',
  room_id: null as number | null,
  teacher_id: null as number | null
})

const applyForm = ref({
  templateId: 0,
  templateName: '',
  offeringId: '',
  startDate: '',
  endDate: '',
  weekdays: [] as number[],
  duration: 60,
  override_buffer: false
})

// 是否可以新增格子
const canAddCell = computed(() => {
  if (!selectedTemplate.value) return false
  if (selectedTemplate.value.row_type === 'ROOM') {
    return newCell.value.room_id !== null
  } else {
    return newCell.value.teacher_id !== null
  }
})

// 是否可以套用模板
const canApply = computed(() => {
  return applyForm.value.offeringId !== '' &&
    applyForm.value.startDate !== '' &&
    applyForm.value.endDate !== '' &&
    applyForm.value.weekdays.length > 0
})

// 是否有可覆蓋的衝突
const hasOverrideableConflicts = computed(() => {
  return applyConflicts.value.some(c => c.can_override)
})

const weekdays = [
  { value: 1, label: '週一' },
  { value: 2, label: '週二' },
  { value: 3, label: '週三' },
  { value: 4, label: '週四' },
  { value: 5, label: '週五' },
  { value: 6, label: '週六' },
  { value: 7, label: '週日' }
]

const form = ref({
  name: '',
  row_type: 'ROOM'
})

// 取得教室名稱
const getRoomName = (roomId: number): string => {
  const room = rooms.value.find(r => r.id === roomId)
  return room ? room.name : `教室 ${roomId}`
}

// 取得老師名稱
const getTeacherName = (teacherId: number): string => {
  const teacher = teachers.value.find(t => t.id === teacherId)
  return teacher ? teacher.name : `老師 ${teacherId}`
}

// 取得格子資源名稱
const getCellResourceName = (cell: any): string => {
  if (cell.room_id) {
    return getRoomName(cell.room_id)
  } else if (cell.teacher_id) {
    return getTeacherName(cell.teacher_id)
  }
  return '未設定'
}

// 取得格子資源樣式類別
const getCellResourceClass = (cell: any): string => {
  if (cell.room_id || cell.teacher_id) {
    return 'text-primary-400'
  }
  return 'text-yellow-500'
}

// 衝突警告類型標籤
const getConflictTypeLabel = (conflictType: string): string => {
  const labels: Record<string, string> = {
    'ROOM_OVERLAP': '教室時間衝突',
    'TEACHER_OVERLAP': '老師時間衝突',
    'PERSONAL_EVENT': '老師私人行程衝突',
    'TEACHER_BUFFER': '老師緩衝時間不足',
    'ROOM_BUFFER': '教室緩衝時間不足',
    'WARNING': '警告',
    'SCHEDULE_WARNING': '課表警告',
    'CAPACITY_WARNING': '容量警告'
  }
  return labels[conflictType] || conflictType
}

// 驗證套用模板（預檢查）
const validateApplyTemplate = async (): Promise<{ hasConflicts: boolean; conflicts: any[]; warnings: any[] }> => {
  const conflicts: any[] = []
  const warnings: any[] = []

  try {
    const api = useApi()
    const response = await api.post<any>(`/admin/templates/${applyForm.value.templateId}/validate-apply`, {
      offering_id: Number(applyForm.value.offeringId),
      start_date: applyForm.value.startDate,
      end_date: applyForm.value.endDate,
      weekdays: applyForm.value.weekdays,
      duration: applyForm.value.duration,
      override_buffer: applyForm.value.override_buffer
    })

    // 解析驗證回應
    if (response.datas) {
      if (Array.isArray(response.datas.conflicts)) {
        conflicts.push(...response.datas.conflicts)
      }
      if (Array.isArray(response.datas.warnings)) {
        warnings.push(...response.datas.warnings)
      }
    }

    return { hasConflicts: conflicts.length > 0, conflicts, warnings }
  } catch (error: any) {
    console.error('Failed to validate template apply:', error)

    // 嘗試從錯誤回應中解析衝突資訊
    try {
      const errorData = (error as any).data || {}
      if (errorData.datas?.conflicts) {
        conflicts.push(...errorData.datas.conflicts)
      }
      if (errorData.datas?.warnings) {
        warnings.push(...errorData.datas.warnings)
      }
    } catch {
      // 忽略解析錯誤
    }

    return { hasConflicts: conflicts.length > 0, conflicts, warnings }
  }
}

// 取得教室列表
const fetchRooms = async () => {
  try {
    const api = useApi()
    const response = await api.get<any>('/admin/rooms')
    if (response.datas?.rooms) {
      rooms.value = response.datas.rooms
    } else if (response.datas) {
      rooms.value = Array.isArray(response.datas) ? response.datas : []
    } else {
      rooms.value = []
    }
  } catch (error) {
    console.error('Failed to fetch rooms:', error)
    rooms.value = []
  }
}

// 取得老師列表
const fetchTeachers = async () => {
  try {
    const api = useApi()
    const response = await api.get<any>('/admin/teachers?limit=1000')
    if (response.datas?.teachers) {
      teachers.value = response.datas.teachers
    } else if (response.datas) {
      teachers.value = Array.isArray(response.datas) ? response.datas : []
    } else {
      teachers.value = []
    }
  } catch (error) {
    console.error('Failed to fetch teachers:', error)
    teachers.value = []
  }
}

// 將老師轉換為 SearchableSelect 選項格式
const teacherOptions = computed<SelectOption[]>(() =>
  teachers.value.map(t => ({
    id: t.id,
    name: t.name
  }))
)

const fetchTemplates = async () => {
  try {
    const api = useApi()
    // parseResponse 已經提取了 datas 欄位，所以 response 就是模板陣列本身
    const response = await api.get<any[]>('/admin/templates')
    templates.value = response || []
  } catch (error) {
    console.error('Failed to fetch templates:', error)
  }
}

const createTemplate = async () => {
  creating.value = true
  try {
    const api = useApi()
    await api.post('/admin/templates', form.value)
    showModal.value = false
    form.value = { name: '', row_type: 'ROOM' }
    await fetchTemplates()
  } catch (error) {
    console.error('Failed to create template:', error)
    await alertError('建立失敗')
  } finally {
    creating.value = false
  }
}

const deleteTemplate = async (id: number) => {
  if (!await alertConfirm('確定要刪除此模板？')) return

  try {
    const api = useApi()
    await api.delete(`/admin/templates/${id}`)
    await fetchTemplates()
  } catch (error) {
    console.error('Failed to delete template:', error)
    await alertError('刪除失敗')
  }
}

const viewTemplate = async (template: any) => {
  selectedTemplate.value = template
  // 重置 newCell
  newCell.value = {
    row_no: 1,
    col_no: 1,
    start_time: '09:00',
    end_time: '10:00',
    room_id: null,
    teacher_id: null
  }
  try {
    const api = useApi()
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/templates/${template.id}/cells`)
    cells.value = response.datas || []
  } catch (error) {
    console.error('Failed to fetch cells:', error)
    cells.value = []
  }
}

const addCells = async () => {
  if (!selectedTemplate.value) return

  // 驗證資源選擇
  if (selectedTemplate.value.row_type === 'ROOM' && !newCell.value.room_id) {
    await alertWarning('請選擇教室')
    return
  }
  if (selectedTemplate.value.row_type === 'TEACHER' && !newCell.value.teacher_id) {
    await alertWarning('請選擇老師')
    return
  }

  addingCell.value = true
  try {
    const api = useApi()
    const cellData = {
      row_no: newCell.value.row_no,
      col_no: newCell.value.col_no,
      start_time: newCell.value.start_time,
      end_time: newCell.value.end_time,
      room_id: selectedTemplate.value.row_type === 'ROOM' ? newCell.value.room_id : null,
      teacher_id: selectedTemplate.value.row_type === 'TEACHER' ? newCell.value.teacher_id : null
    }
    await api.post(`/admin/templates/${selectedTemplate.value.id}/cells`, [cellData])
    await viewTemplate(selectedTemplate.value)

    // 自動遞增列和行號
    newCell.value.col_no++
    if (newCell.value.col_no > 4) {
      newCell.value.col_no = 1
      newCell.value.row_no++
    }

    await alertSuccess('格子新增成功')
  } catch (error) {
    console.error('Failed to add cells:', error)
    await alertError('新增格子失敗：' + (error as Error).message)
  } finally {
    addingCell.value = false
  }
}

const deleteCell = async (cellId: number) => {
  if (!await alertConfirm('確定要刪除此格子？')) return

  try {
    const api = useApi()
    await api.delete(`/admin/templates/cells/${cellId}`)
    await viewTemplate(selectedTemplate.value)
    await alertSuccess('格子已刪除')
  } catch (error) {
    console.error('Failed to delete cell:', error)
    await alertError('刪除格子失敗')
  }
}

const fetchOfferings = async () => {
  offeringsLoading.value = true
  try {
    const api = useApi()
    const response = await api.get<any>('/admin/offerings')

    // useApi 的 parseResponse 已經提取了 datas 欄位
    // API 回應格式: { Offerings: [...], Pagination: {...} }
    // 注意：使用駝峰式命名 (Offerings 不是 offerings)
    if (response && typeof response === 'object') {
      if (Array.isArray(response)) {
        // 直接是陣列格式
        offerings.value = response
      } else if (Array.isArray(response.Offerings)) {
        // { Offerings: [...] } 格式
        offerings.value = response.Offerings
      } else if (Array.isArray(response.offerings)) {
        // { offerings: [...] } 格式（小寫）
        offerings.value = response.offerings
      } else {
        offerings.value = []
      }
    } else {
      offerings.value = []
    }
  } catch (error) {
    offerings.value = []
  } finally {
    offeringsLoading.value = false
  }
}

const openApplyModal = async (template: any) => {
  selectedTemplate.value = null
  applyForm.value = {
    templateId: template.id,
    templateName: template.name,
    offeringId: '',
    startDate: '',
    endDate: '',
    weekdays: [],
    duration: 60,
    override_buffer: false
  }
  applyConflicts.value = []
  validationWarnings.value = []
  showApplyModal.value = true

  // 確保課程資料已載入
  if (offerings.value.length === 0 && !offeringsLoading.value) {
    await fetchOfferings()
  }
}

const closeApplyModal = () => {
  showApplyModal.value = false
  applyConflicts.value = []
  validationWarnings.value = []
}

const applyTemplate = async () => {
  if (!applyForm.value.offeringId || applyForm.value.weekdays.length === 0) {
    await alertWarning('請填寫完整資訊')
    return
  }

  // 先進行驗證
  isValidating.value = true
  applyConflicts.value = []
  validationWarnings.value = []

  try {
    const validationResult = await validateApplyTemplate()
    isValidating.value = false

    // 如果有衝突或警告，顯示在對話框中並讓用戶確認
    if (validationResult.hasConflicts || validationResult.warnings.length > 0) {
      applyConflicts.value = validationResult.conflicts
      validationWarnings.value = validationResult.warnings

      // 構建確認訊息
      let confirmMessage = '偵測到以下問題：\n\n'
      if (validationResult.conflicts.length > 0) {
        confirmMessage += `⚠️ ${validationResult.conflicts.length} 個衝突\n`
      }
      if (validationResult.warnings.length > 0) {
        confirmMessage += `💡 ${validationResult.warnings.length} 個警告\n`
      }
      confirmMessage += '\n是否仍要套用模板？'

      // 讓用戶確認是否繼續
      if (!await alertConfirm(confirmMessage)) {
        return
      }
    }

    // 用戶確認或無問題，執行套用
    applying.value = true

    const api = useApi()
    const response = await api.post<any>(`/admin/templates/${applyForm.value.templateId}/apply`, {
      offering_id: Number(applyForm.value.offeringId),
      start_date: applyForm.value.startDate,
      end_date: applyForm.value.endDate,
      weekdays: applyForm.value.weekdays,
      duration: applyForm.value.duration,
      override_buffer: applyForm.value.override_buffer
    })

    // 檢查是否為衝突警告回應 (40003)
    if (response.code === 40003 && response.datas?.conflicts) {
      applyConflicts.value = response.datas.conflicts
      if (!applyForm.value.override_buffer) {
        await alertWarning('偵測到 Buffer 衝突，請勾選「允許覆蓋 Buffer 衝突」後再試')
      }
      applying.value = false
      return
    }

    // 如果有不可覆蓋的衝突 (40002)
    if (response.code === 40002 && response.datas?.conflicts) {
      applyConflicts.value = response.datas.conflicts
      await alertError('無法套用模板：存在不可覆蓋的時間衝突，請查看詳細資訊')
      applying.value = false
      return
    }

    // 檢查其他錯誤碼
    if (response.code !== 0 && response.code !== 200) {
      // 顯示後端回傳的錯誤訊息
      const errorMsg = response.message || `錯誤碼: ${response.code}`
      await alertError(errorMsg)
      applying.value = false
      return
    }

    closeApplyModal()
    notificationUI.showSuccess('模板套用成功')
  } catch (error: any) {
    console.error('Failed to apply template:', error)
    applying.value = false
    isValidating.value = false

    // 嘗試解析錯誤回應
    try {
      // 解析後端可能回傳的衝突資訊
      const errorData = (error as any).data || {}
      const errorMessage = (error as any).message || ''
      const errorCode = (error as any).code || (errorData.code || 0)

      // 如果有衝突資訊，顯示衝突詳情
      if (errorData.datas?.conflicts) {
        applyConflicts.value = errorData.datas.conflicts
      }

      // 根據錯誤碼顯示對應訊息
      if (errorCode === 40002) {
        await alertError('無法套用模板：存在不可覆蓋的時間衝突，請查看詳細資訊')
      } else if (errorCode === 40003) {
        await alertWarning('偵測到 Buffer 衝突，請勾選「允許覆蓋 Buffer 衝突」後再試')
      } else if (errorMessage) {
        await alertError(errorMessage)
      } else {
        await alertError('套用失敗，請稍後再試')
      }
    } catch {
      await alertError('套用失敗，請稍後再試')
    }
  } finally {
    applying.value = false
    isValidating.value = false
  }
}

const formatDate = (dateStr: string): string => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-TW')
}

onMounted(async () => {
  await Promise.all([fetchTemplates(), fetchOfferings(), fetchRooms(), fetchTeachers()])
})
</script>

<style scoped>
.input-field {
  @apply w-full px-3 py-2 rounded-lg bg-slate-800 border border-white/10 text-white cursor-pointer appearance-none;
}

.input-field:focus {
  @apply outline-none border-primary-500;
}

select.input-field {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 0.5rem center;
  background-repeat: no-repeat;
  background-size: 1.5em 1.5em;
  padding-right: 2.5rem;
}

/* 展開/收合動畫 */
.rotate-90 {
  transform: rotate(90deg);
}

/* 拖曳相關樣式 */
.cursor-move {
  cursor: grab;
}

.cursor-move:active {
  cursor: grabbing;
}

/* 拖曳時的變形效果 */
.scale-102 {
  transform: scale(1.02);
}

/* 拖曳區域動畫 */
.transition-all {
  transition-property: all;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}

/* 拖曳處理圖示 */
.drag-handle {
  opacity: 0.5;
  transition: opacity 0.2s;
}

.draggable-cell:hover .drag-handle {
  opacity: 1;
}
</style>
