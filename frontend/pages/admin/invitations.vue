<template>
  <div class="p-4 md:p-6 max-w-7xl mx-auto">
    <div class="mb-6 md:mb-8">
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
            邀請紀錄
          </h1>
          <p class="text-slate-400 text-sm md:text-base">
            查看邀請老師的歷史記錄與處理結果
          </p>
        </div>
        <button
          @click="showGenerateModal = true"
          class="px-4 py-2 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
          </svg>
          產生邀請連結
        </button>
      </div>
    </div>

    <!-- 統計卡片 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <div class="bg-white/5 rounded-xl p-4 border border-white/10">
        <p class="text-slate-400 text-sm">待處理</p>
        <p class="text-2xl font-bold text-warning-500">{{ stats.pending }}</p>
      </div>
      <div class="bg-white/5 rounded-xl p-4 border border-white/10">
        <p class="text-slate-400 text-sm">已接受</p>
        <p class="text-2xl font-bold text-success-500">{{ stats.accepted }}</p>
      </div>
      <div class="bg-white/5 rounded-xl p-4 border border-white/10">
        <p class="text-slate-400 text-sm">已婉拒</p>
        <p class="text-2xl font-bold text-critical-500">{{ stats.declined }}</p>
      </div>
      <div class="bg-white/5 rounded-xl p-4 border border-white/10">
        <p class="text-slate-400 text-sm">已過期</p>
        <p class="text-2xl font-bold text-slate-400">{{ stats.expired }}</p>
      </div>
    </div>

    <!-- 標籤頁 -->
    <div class="mb-6">
      <div class="flex gap-2 border-b border-white/10 pb-2">
        <button
          @click="activeTab = 'invitations'"
          class="px-4 py-2 rounded-t-lg transition-colors"
          :class="activeTab === 'invitations' ? 'bg-white/10 text-white' : 'text-slate-400 hover:text-white'"
        >
          邀請記錄
        </button>
        <button
          @click="activeTab = 'links'"
          class="px-4 py-2 rounded-t-lg transition-colors"
          :class="activeTab === 'links' ? 'bg-white/10 text-white' : 'text-slate-400 hover:text-white'"
        >
          邀請連結
          <span v-if="links.length > 0" class="ml-1 px-1.5 py-0.5 text-xs bg-primary-500/20 text-primary-400 rounded">
            {{ links.length }}
          </span>
        </button>
        <button
          @click="activeTab = 'general'"
          class="px-4 py-2 rounded-t-lg transition-colors"
          :class="activeTab === 'general' ? 'bg-white/10 text-white' : 'text-slate-400 hover:text-white'"
        >
          通用連結
          <span v-if="generalLink" class="ml-1 px-1.5 py-0.5 text-xs bg-success-500/20 text-success-400 rounded">
            已啟用
          </span>
        </button>
      </div>
    </div>

    <!-- 邀請記錄標籤 -->
    <div v-if="activeTab === 'invitations'">
      <!-- 篩選器 -->
      <div class="flex flex-wrap gap-3 mb-6">
        <select
          v-model="filters.inviteType"
          class="px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-slate-300 focus:outline-none focus:border-primary-500"
        >
          <option value="">全部類型</option>
          <option value="GENERAL">通用邀請</option>
          <option value="TEACHER">指定邀請</option>
        </select>

        <select
          v-model="filters.status"
          class="px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-slate-300 focus:outline-none focus:border-primary-500"
        >
          <option value="">全部狀態</option>
          <option value="PENDING">待處理</option>
          <option value="ACCEPTED">已接受</option>
          <option value="DECLINED">已婉拒</option>
          <option value="EXPIRED">已過期</option>
        </select>

        <input
          v-model="filters.search"
          type="text"
          placeholder="搜尋 Email..."
          class="px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-slate-300 focus:outline-none focus:border-primary-500 flex-1 min-w-[200px]"
        />

        <button
          @click="refreshData"
          class="px-4 py-2 bg-primary-500/20 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/30 transition-colors"
        >
          重新整理
        </button>
      </div>

      <!-- 邀請列表 -->
      <div class="bg-white/5 rounded-xl border border-white/10 overflow-hidden">
        <div v-if="loading" class="p-8 text-center">
          <div class="inline-block w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin"></div>
          <p class="text-slate-400 mt-2">載入中...</p>
        </div>

        <div v-else-if="filteredInvitations.length === 0" class="p-8 text-center text-slate-400">
          暫無邀請記錄
        </div>

        <table v-else class="w-full">
          <thead class="bg-white/5">
            <tr>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">邀請類型</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">Email</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">狀態</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">邀請時間</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">回應時間</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="invitation in filteredInvitations" :key="invitation.id" class="hover:bg-white/5">
              <td class="px-4 py-3">
                <span
                  class="px-2 py-1 rounded-full text-xs font-medium"
                  :class="{
                    'bg-primary-500/20 text-primary-400': invitation.invite_type === 'GENERAL',
                    'bg-slate-500/20 text-slate-400': invitation.invite_type !== 'GENERAL'
                  }"
                >
                  {{ inviteTypeText(invitation.invite_type) }}
                </span>
              </td>
              <td class="px-4 py-3 text-slate-300">{{ invitation.email || '-' }}</td>
              <td class="px-4 py-3">
                <span
                  class="px-2 py-1 rounded-full text-xs font-medium"
                  :class="{
                    'bg-warning-500/20 text-warning-500': invitation.status === 'PENDING',
                    'bg-success-500/20 text-success-500': invitation.status === 'ACCEPTED',
                    'bg-critical-500/20 text-critical-500': invitation.status === 'DECLINED',
                    'bg-slate-500/20 text-slate-400': invitation.status === 'EXPIRED',
                  }"
                >
                  {{ statusText(invitation.status) }}
                </span>
              </td>
              <td class="px-4 py-3 text-slate-400 text-sm">
                {{ formatDate(invitation.created_at) }}
              </td>
              <td class="px-4 py-3 text-slate-400 text-sm">
                {{ invitation.responded_at ? formatDate(invitation.responded_at) : '-' }}
              </td>
              <td class="px-4 py-3">
                <span class="text-slate-500 text-sm">-</span>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- 分頁 -->
        <div v-if="pagination.totalPages > 1" class="px-4 py-3 bg-white/5 border-t border-white/10 flex items-center justify-between">
          <p class="text-slate-400 text-sm">
            第 {{ pagination.page }} 頁，共 {{ pagination.totalPages }} 頁
          </p>
          <div class="flex gap-2">
            <button
              @click="changePage(pagination.page - 1)"
              :disabled="pagination.page === 1"
              class="px-3 py-1 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              上一頁
            </button>
            <button
              @click="changePage(pagination.page + 1)"
              :disabled="pagination.page === pagination.totalPages"
              class="px-3 py-1 rounded-lg bg-white/5 text-slate-300 hover:bg-white/10 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              下一頁
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 邀請連結標籤 -->
    <div v-if="activeTab === 'links'">
      <div class="bg-white/5 rounded-xl border border-white/10 overflow-hidden">
        <div v-if="loadingLinks" class="p-8 text-center">
          <div class="inline-block w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin"></div>
          <p class="text-slate-400 mt-2">載入中...</p>
        </div>

        <div v-else-if="links.length === 0" class="p-8 text-center text-slate-400">
          暫無有效邀請連結
        </div>

        <table v-else class="w-full">
          <thead class="bg-white/5">
            <tr>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">Email</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">職位</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">有效期限</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-slate-400">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="link in links" :key="link.id" class="hover:bg-white/5">
              <td class="px-4 py-3 text-slate-300">{{ link.email }}</td>
              <td class="px-4 py-3">
                <span class="px-2 py-1 rounded-full text-xs font-medium bg-primary-500/20 text-primary-400">
                  {{ roleText(link.role) }}
                </span>
              </td>
              <td class="px-4 py-3 text-slate-400 text-sm">
                {{ link.expires_at ? formatDate(link.expires_at) : '無期限' }}
              </td>
              <td class="px-4 py-3">
                <div class="flex gap-2">
                  <button
                    @click="copyLink(link)"
                    class="text-primary-500 hover:text-primary-400 text-sm flex items-center gap-1"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-28 5av-1M2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
                    </svg>
                    複製
                  </button>
                  <button
                    @click="revokeLink(link)"
                    class="text-critical-500 hover:text-critical-400 text-sm flex items-center gap-1"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    撤回
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 通用邀請連結標籤 -->
    <div v-if="activeTab === 'general'">
      <div class="bg-white/5 rounded-xl border border-white/10 overflow-hidden">
        <div class="p-6">
          <div class="flex items-start justify-between mb-6">
            <div>
              <h3 class="text-lg font-bold text-white mb-2">通用邀請連結</h3>
              <p class="text-slate-400 text-sm">
                產生一個不綁定 Email 的通用邀請連結，可分享給多位老師重複使用
              </p>
            </div>
          </div>

          <!-- 通用連結狀態卡片 -->
          <div v-if="loadingGeneral" class="p-8 text-center">
            <div class="inline-block w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin"></div>
            <p class="text-slate-400 mt-2">載入中...</p>
          </div>

          <div v-else-if="generalLink" class="space-y-4">
            <!-- 連結資訊 -->
            <div class="p-4 bg-success-500/10 border border-success-500/30 rounded-xl">
              <div class="flex items-center gap-2 mb-3">
                <span class="px-2 py-1 rounded-full text-xs font-medium bg-success-500/20 text-success-400">
                  已啟用
                </span>
                <span class="text-slate-400 text-sm">
                  職位：{{ roleText(generalLink.role) }}
                </span>
              </div>

              <div class="mb-4">
                <label class="block text-slate-400 text-sm mb-2">邀請連結</label>
                <div class="flex gap-2">
                  <input
                    :value="generalLink.invite_link"
                    readonly
                    class="flex-1 px-3 py-2 bg-white/5 border border-white/10 rounded-lg text-slate-300 text-sm"
                  />
                  <button
                    @click="copyGeneralLink"
                    class="px-3 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600 transition-colors flex items-center gap-1"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
                    </svg>
                    複製
                  </button>
                </div>
              </div>

              <div class="flex items-center justify-between text-sm">
                <span class="text-slate-400">
                  有效期限：無期限
                </span>
              </div>
            </div>

            <!-- 操作按鈕 -->
            <div class="flex gap-3">
              <button
                @click="toggleGeneralLink"
                class="flex-1 py-3 rounded-xl transition-colors flex items-center justify-center gap-2"
                :class="generalLink.status === 'PENDING' ? 'bg-warning-500/20 text-warning-400 hover:bg-warning-500/30 border border-warning-500/30' : 'bg-success-500/20 text-success-400 hover:bg-success-500/30 border border-success-500/30'"
              >
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                </svg>
                {{ generalLink.status === 'PENDING' ? '停用連結' : '重新啟用' }}
              </button>
              <button
                @click="showGeneralModal = true"
                class="flex-1 py-3 bg-primary-500/20 text-primary-400 border border-primary-500/30 rounded-xl hover:bg-primary-500/30 transition-colors flex items-center justify-center gap-2"
              >
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                重新產生
              </button>
            </div>
          </div>

          <div v-else class="p-8 text-center">
            <div class="mb-4">
              <svg class="w-12 h-12 mx-auto text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
              </svg>
            </div>
            <p class="text-slate-400 mb-4">尚未產生通用邀請連結</p>
            <button
              @click="showGeneralModal = true"
              class="px-6 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors flex items-center gap-2 mx-auto"
            >
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              產生通用連結
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 產生邀請連結 Modal -->
    <div v-if="showGenerateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div class="bg-slate-800 rounded-2xl p-6 w-full max-w-md border border-white/10">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-bold text-white">產生邀請連結</h3>
          <button @click="closeGenerateModal" class="text-slate-400 hover:text-white">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="generateLink">
          <div class="mb-4">
            <label class="block text-slate-400 text-sm mb-2">Email 地址</label>
            <input
              v-model="generateForm.email"
              type="email"
              required
              placeholder="輸入要邀請的 Email"
              class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary-500"
            />
          </div>

          <div class="mb-4">
            <label class="block text-slate-400 text-sm mb-2">職位</label>
            <select
              v-model="generateForm.role"
              required
              class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary-500"
            >
              <option value="TEACHER">正式老師</option>
              <option value="SUBSTITUTE">代課老師</option>
            </select>
          </div>

          <div class="mb-6">
            <label class="block text-slate-400 text-sm mb-2">邀請訊息（選填）</label>
            <textarea
              v-model="generateForm.message"
              rows="3"
              placeholder="輸入邀請訊息..."
              class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary-500 resize-none"
            ></textarea>
          </div>

          <div v-if="generatedLink" class="mb-6 p-4 bg-success-500/10 border border-success-500/30 rounded-xl">
            <p class="text-success-500 text-sm font-medium mb-2">邀請連結已產生！</p>
            <div class="flex gap-2">
              <input
                :value="generatedLink.invite_link"
                readonly
                class="flex-1 px-3 py-2 bg-white/5 border border-white/10 rounded-lg text-slate-300 text-sm"
              />
              <button
                @click="copyGeneratedLink"
                type="button"
                class="px-3 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600 transition-colors"
              >
                複製
              </button>
            </div>
          </div>

          <div class="flex gap-3">
            <button
              type="button"
              @click="closeGenerateModal"
              class="flex-1 py-3 bg-white/10 text-white rounded-xl hover:bg-white/20 transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="generating"
              class="flex-1 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <span v-if="generating" class="inline-block w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              {{ generating ? '產生中...' : '產生連結' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 產生通用邀請連結 Modal -->
    <div v-if="showGeneralModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div class="bg-slate-800 rounded-2xl p-6 w-full max-w-md border border-white/10">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-bold text-white">產生通用邀請連結</h3>
          <button @click="closeGeneralModal" class="text-slate-400 hover:text-white">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="generateGeneralLink">
          <div class="mb-6">
            <label class="block text-slate-400 text-sm mb-2">邀請訊息（選填）</label>
            <textarea
              v-model="generalForm.message"
              rows="3"
              placeholder="輸入邀請訊息..."
              class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary-500 resize-none"
            ></textarea>
          </div>

          <div class="flex gap-3">
            <button
              type="button"
              @click="closeGeneralModal"
              class="flex-1 py-3 bg-white/10 text-white rounded-xl hover:bg-white/20 transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="generating"
              class="flex-1 py-3 bg-primary-500 text-white rounded-xl hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <span v-if="generating" class="inline-block w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              {{ generating ? '產生中...' : '產生連結' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Toast 通知 -->
    <div
      v-if="toast.show"
      class="fixed bottom-6 right-6 px-6 py-3 rounded-xl text-white shadow-lg z-50"
      :class="toast.type === 'success' ? 'bg-success-500' : 'bg-critical-500'"
    >
      {{ toast.message }}
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  auth: 'ADMIN',
  layout: 'admin',
})

const config = useRuntimeConfig()
const authStore = useAuthStore()

// 取得當前登入管理員的中心 ID
const centerId = computed(() => {
  return authStore.user?.center_id || 1
})

// 標籤頁
const activeTab = ref('invitations')

// 介面
interface Invitation {
  id: number
  email: string
  status: string
  invite_type: string
  created_at: string
  responded_at?: string
  expires_at: string | null
}

interface InvitationLink {
  id: number
  email: string
  role: string
  invite_link: string
  expires_at: string | null
  created_at: string
}

// 狀態
const loading = ref(true)
const loadingLinks = ref(false)
const invitations = ref<Invitation[]>([])
const links = ref<InvitationLink[]>([])
const stats = ref({
  pending: 0,
  accepted: 0,
  declined: 0,
  expired: 0,
})

const filters = ref({
  status: '',
  search: '',
  inviteType: '',
})

const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  totalPages: 0,
})

// 產生連結表單
const showGenerateModal = ref(false)
const generating = ref(false)
const generatedLink = ref<InvitationLink | null>(null)
const generateForm = ref({
  email: '',
  role: 'TEACHER',
  message: '',
})

// 通用連結相關狀態
const showGeneralModal = ref(false)
const loadingGeneral = ref(false)
const generalLink = ref<InvitationLink | null>(null)
const generalForm = ref({
  message: '',
})

// Toast
const toast = ref({
  show: false,
  message: '',
  type: 'success',
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

// 狀態文字
const statusText = (status: string) => {
  const texts: Record<string, string> = {
    PENDING: '待處理',
    ACCEPTED: '已接受',
    DECLINED: '已婉拒',
    EXPIRED: '已過期',
  }
  return texts[status] || status
}

// 邀請類型文字
const inviteTypeText = (inviteType: string) => {
  const texts: Record<string, string> = {
    GENERAL: '通用邀請',
    TEACHER: '指定邀請',
    TALENT_POOL: '人才庫',
    MEMBER: '會員',
  }
  return texts[inviteType] || inviteType
}

// 角色文字
const roleText = (role: string) => {
  const texts: Record<string, string> = {
    TEACHER: '正式老師',
    SUBSTITUTE: '代課老師',
  }
  return texts[role] || role
}

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// 取得統計
const fetchStats = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId.value}/invitations/stats`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    if (response.ok) {
      const data = await response.json()
      stats.value = {
        pending: data.datas?.pending || 0,
        accepted: data.datas?.accepted || 0,
        declined: data.datas?.declined || 0,
        expired: data.datas?.expired || 0,
      }
    }
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

// 取得邀請列表
const fetchInvitations = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    if (filters.value.status) {
      params.append('status', filters.value.status)
    }

    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId.value}/invitations?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    if (response.ok) {
      const data = await response.json()
      invitations.value = data.datas?.data || []
      pagination.value.total = data.datas?.total || 0
      pagination.value.totalPages = Math.ceil(pagination.value.total / pagination.value.limit)
    }
  } catch (error) {
    console.error('Failed to fetch invitations:', error)
  } finally {
    loading.value = false
  }
}

// 取得邀請連結列表
const fetchLinks = async () => {
  loadingLinks.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId.value}/invitations/links`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    if (response.ok) {
      const data = await response.json()
      links.value = data.datas || []
    }
  } catch (error) {
    console.error('Failed to fetch links:', error)
  } finally {
    loadingLinks.value = false
  }
}

// 重新整理資料
const refreshData = () => {
  pagination.value.page = 1
  fetchStats()
  fetchInvitations()
  if (activeTab.value === 'links') {
    fetchLinks()
  }
}

// 切換頁碼
const changePage = (page: number) => {
  pagination.value.page = page
  fetchInvitations()
}

// 篩選後的邀請
const filteredInvitations = computed(() => {
  let result = invitations.value

  // 篩選邀請類型
  if (filters.value.inviteType) {
    result = result.filter(inv => inv.invite_type === filters.value.inviteType)
  }

  // 搜尋 Email
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    result = result.filter(inv =>
      (inv.email && inv.email.toLowerCase().includes(search))
    )
  }

  return result
})

// 產生邀請連結
const generateLink = async () => {
  if (!generateForm.value.email || !generateForm.value.role) return

  generating.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId.value}/invitations/generate-link`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(generateForm.value),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || '產生邀請連結失敗')
    }

    generatedLink.value = data.datas
    fetchLinks()
    showToast('邀請連結產生成功')
  } catch (error: any) {
    showToast(error.message || '產生邀請連結失敗', 'error')
  } finally {
    generating.value = false
  }
}

// 複製連結
const copyLink = async (link: InvitationLink) => {
  try {
    await navigator.clipboard.writeText(link.invite_link)
    showToast('連結已複製到剪貼簿')
  } catch (error) {
    showToast('複製失敗，請手動複製', 'error')
  }
}

// 複製產生的連結
const copyGeneratedLink = async () => {
  if (!generatedLink.value) return
  try {
    await navigator.clipboard.writeText(generatedLink.value.invite_link)
    showToast('連結已複製到剪貼簿')
  } catch (error) {
    showToast('複製失敗，請手動複製', 'error')
  }
}

// 撤回連結
const revokeLink = async (link: InvitationLink) => {
  if (!confirm(`確定要撤回「${link.email}」的邀請連結嗎？`)) return

  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/invitations/links/${link.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || '撤回邀請連結失敗')
    }

    fetchLinks()
    showToast('邀請連結已撤回')
  } catch (error: any) {
    showToast(error.message || '撤回邀請連結失敗', 'error')
  }
}

// 關閉產生 Modal
const closeGenerateModal = () => {
  showGenerateModal.value = false
  generatedLink.value = null
  generateForm.value = {
    email: '',
    role: 'TEACHER',
    message: '',
  }
}

// 取得通用邀請連結
const fetchGeneralLink = async () => {
  loadingGeneral.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId.value}/invitations/general-link`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    if (response.ok) {
      const data = await response.json()
      generalLink.value = data.datas || null
    } else {
      generalLink.value = null
    }
  } catch (error) {
    console.error('Failed to fetch general link:', error)
    generalLink.value = null
  } finally {
    loadingGeneral.value = false
  }
}

// 產生通用邀請連結
const generateGeneralLink = async () => {
  generating.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId.value}/invitations/general-link`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        role: 'TEACHER',
        message: generalForm.value.message,
      }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || '產生通用邀請連結失敗')
    }

    generalLink.value = data.datas
    showGeneralModal.value = false
    generalForm.value = {
      role: 'TEACHER',
      message: '',
    }
    showToast('通用邀請連結產生成功')
  } catch (error: any) {
    showToast(error.message || '產生通用邀請連結失敗', 'error')
  } finally {
    generating.value = false
  }
}

// 切換通用邀請連結狀態
const toggleGeneralLink = async () => {
  if (!generalLink.value) return

  const action = generalLink.value.status === 'PENDING' ? '停用' : '重新啟用'
  if (!confirm(`確定要${action}通用邀請連結嗎？`)) return

  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.public.apiBase}/admin/centers/${centerId.value}/invitations/general-link/toggle`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || `${action}失敗`)
    }

    await fetchGeneralLink()
    showToast(`通用邀請連結已${action}`)
  } catch (error: any) {
    showToast(error.message || `${action}失敗`, 'error')
  }
}

// 複製通用連結
const copyGeneralLink = async () => {
  if (!generalLink.value) return
  try {
    await navigator.clipboard.writeText(generalLink.value.invite_link)
    showToast('連結已複製到剪貼簿')
  } catch (error) {
    showToast('複製失敗，請手動複製', 'error')
  }
}

// 關閉通用連結 Modal
const closeGeneralModal = () => {
  showGeneralModal.value = false
  generalForm.value = {
    message: '',
  }
}

// 監聽標籤頁切換
watch(activeTab, (newTab) => {
  if (newTab === 'links' && links.value.length === 0) {
    fetchLinks()
  }
  if (newTab === 'general') {
    fetchGeneralLink()
  }
})

// 監聽篩選器變化
watch(filters, () => {
  pagination.value.page = 1
  fetchInvitations()
}, { deep: true })

onMounted(() => {
  refreshData()
})
</script>
