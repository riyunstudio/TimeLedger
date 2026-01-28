<template>
  <div class="p-4 md:p-6 max-w-6xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl md:text-3xl font-bold text-slate-100 mb-2">
            管理員管理
          </h1>
          <p class="text-slate-400 text-sm md:text-base">
            管理中心的管理員帳號（僅擁有者可執行）
          </p>
        </div>
        <!-- 新增管理員按鈕 -->
        <button
          v-if="canAddAdmin"
          @click="showCreateModal = true"
          class="px-4 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          新增管理員
        </button>
      </div>
    </div>

    <!-- 篩選器 -->
    <div class="mb-6 flex flex-wrap items-center gap-4">
      <!-- 角色篩選 -->
      <div class="flex items-center gap-2">
        <label class="text-sm text-slate-400">角色篩選：</label>
        <select
          v-model="selectedRole"
          class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 min-w-[120px]"
        >
          <option value="">全部角色</option>
          <option value="OWNER">擁有者</option>
          <option value="ADMIN">管理員</option>
          <option value="STAFF">員工</option>
        </select>
      </div>

      <!-- 狀態篩選 -->
      <div class="flex items-center gap-2">
        <label class="text-sm text-slate-400">狀態篩選：</label>
        <select
          v-model="selectedStatus"
          class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 min-w-[120px]"
        >
          <option value="">全部狀態</option>
          <option value="ACTIVE">啟用</option>
          <option value="INACTIVE">停用</option>
        </select>
      </div>

      <!-- 清除篩選 -->
      <button
        v-if="selectedRole || selectedStatus"
        @click="clearFilters"
        class="text-sm text-primary-400 hover:text-primary-300 transition-colors"
      >
        清除篩選
      </button>
    </div>

    <!-- 管理員列表 -->
    <BaseGlassCard>
      <div class="p-6">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="text-left text-slate-400 text-sm border-b border-white/10">
                <th class="pb-4 pr-4 font-medium">姓名</th>
                <th class="pb-4 pr-4 font-medium">Email</th>
                <th class="pb-4 pr-4 font-medium">角色</th>
                <th class="pb-4 pr-4 font-medium">狀態</th>
                <th class="pb-4 pr-4 font-medium">LINE</th>
                <th class="pb-4 font-medium">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="admin in filteredAdminList"
                :key="admin.id"
                class="border-b border-white/5 hover:bg-white/5 transition-colors"
              >
                <td class="py-4 pr-4">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-primary-500/20 flex items-center justify-center">
                      <span class="text-primary-400 font-medium text-sm">
                        {{ admin.name.charAt(0) }}
                      </span>
                    </div>
                    <span class="text-white font-medium">{{ admin.name }}</span>
                  </div>
                </td>
                <td class="py-4 pr-4 text-slate-300">{{ admin.email }}</td>
                <td class="py-4 pr-4">
                  <!-- 角色顯示/選擇 -->
                  <div v-if="canManageAdmin(admin)" class="relative">
                    <select
                      v-model="admin.role"
                      @change="changeRole(admin)"
                      class="px-3 py-1.5 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500 cursor-pointer"
                      :class="{
                        'border-purple-500/50': admin.role === 'OWNER',
                        'border-blue-500/50': admin.role === 'ADMIN',
                        'border-green-500/50': admin.role === 'STAFF'
                      }"
                    >
                      <option value="ADMIN">管理員</option>
                      <option value="STAFF">員工</option>
                      <option value="OWNER">擁有者</option>
                    </select>
                  </div>
                  <span
                    v-else
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="{
                      'bg-purple-500/20 text-purple-400': admin.role === 'OWNER',
                      'bg-blue-500/20 text-blue-400': admin.role === 'ADMIN',
                      'bg-green-500/20 text-green-400': admin.role === 'STAFF'
                    }"
                  >
                    {{ roleText(admin.role) }}
                  </span>
                </td>
                <td class="py-4 pr-4">
                  <span
                    class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="{
                      'bg-success-500/20 text-success-400': admin.status === 'ACTIVE',
                      'bg-critical-500/20 text-critical-400': admin.status === 'INACTIVE'
                    }"
                  >
                    <span
                      class="w-1.5 h-1.5 rounded-full"
                      :class="{
                        'bg-success-400': admin.status === 'ACTIVE',
                        'bg-critical-400': admin.status === 'INACTIVE'
                      }"
                    />
                    {{ admin.status === 'ACTIVE' ? '啟用' : '停用' }}
                  </span>
                </td>
                <td class="py-4 pr-4">
                  <div v-if="admin.line_user_id" class="flex items-center gap-1.5 text-success-400">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    <span class="text-xs">已綁定</span>
                  </div>
                  <span v-else class="text-slate-500 text-sm">未綁定</span>
                </td>
                <td class="py-4">
                  <div class="flex items-center gap-2">
                    <!-- 重設密碼按鈕 -->
                    <button
                      v-if="canManageAdmin(admin)"
                      @click="openResetPasswordModal(admin)"
                      class="p-2 rounded-lg hover:bg-white/10 text-slate-400 hover:text-white transition-colors"
                      title="重設密碼"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                      </svg>
                    </button>

                    <!-- 停用/啟用按鈕 -->
                    <button
                      v-if="canManageAdmin(admin)"
                      @click="toggleAdminStatus(admin)"
                      class="p-2 rounded-lg transition-colors"
                      :class="admin.status === 'ACTIVE'
                        ? 'hover:bg-critical-500/20 text-slate-400 hover:text-critical-400'
                        : 'hover:bg-success-500/20 text-slate-400 hover:text-success-400'"
                      :title="admin.status === 'ACTIVE' ? '停用' : '啟用'"
                    >
                      <svg v-if="admin.status === 'ACTIVE'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                      </svg>
                      <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                      </svg>
                    </button>

                    <!-- 自己的帳號標記 -->
                    <span v-if="admin.id === currentAdminId" class="text-xs text-slate-500 ml-2">
                      (本人)
                    </span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- 空狀態 -->
          <div v-if="filteredAdminList.length === 0 && !loading" class="text-center py-12">
            <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-slate-700/50 flex items-center justify-center">
              <svg class="w-8 h-8 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </div>
            <p class="text-slate-400">暫無符合條件的管理員</p>
          </div>

          <!-- 載入中 -->
          <div v-if="loading" class="text-center py-12">
            <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
            <p class="text-slate-400 mt-4">載入中...</p>
          </div>
        </div>
      </div>
    </BaseGlassCard>

    <!-- 重設密碼對話框 -->
    <Teleport to="body">
      <div
        v-if="showResetPasswordModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="closeResetPasswordModal"
      >
        <div class="glass-card w-full max-w-md">
          <div class="p-6 border-b border-white/10">
            <h3 class="text-lg font-semibold text-white">重設密碼</h3>
            <p v-if="selectedAdmin" class="text-sm text-slate-400 mt-1">
              將重設 {{ selectedAdmin.name }} 的密碼
            </p>
          </div>

          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm text-slate-400 mb-2">新密碼</label>
              <BaseInput
                v-model="resetPasswordForm.newPassword"
                type="password"
                placeholder="請輸入新密碼（至少 6 個字元）"
                class="w-full"
              />
            </div>

            <div>
              <label class="block text-sm text-slate-400 mb-2">確認新密碼</label>
              <BaseInput
                v-model="resetPasswordForm.confirmPassword"
                type="password"
                placeholder="請再次輸入新密碼"
                class="w-full"
                :error="resetPasswordForm.confirmPassword && resetPasswordForm.newPassword !== resetPasswordForm.confirmPassword ? '密碼不一致' : ''"
              />
            </div>
          </div>

          <div class="p-6 border-t border-white/10 flex items-center gap-4">
            <button
              @click="closeResetPasswordModal"
              class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              取消
            </button>
            <button
              @click="resetPassword"
              :disabled="!isResetPasswordValid || resetting"
              class="flex-1 px-4 py-2 rounded-lg bg-primary-500/30 border border-primary-500 text-primary-400 hover:bg-primary-500/40 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="resetting" class="flex items-center justify-center gap-2">
                <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                處理中...
              </span>
              <span v-else>確認重設</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 停用/啟用確認對話框 -->
    <Teleport to="body">
      <div
        v-if="showToggleStatusModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="closeToggleStatusModal"
      >
        <div class="glass-card w-full max-w-sm">
          <div class="p-6 border-b border-white/10">
            <h3 class="text-lg font-semibold text-white">
              {{ toggleAction === 'disable' ? '停用管理員' : '啟用管理員' }}
            </h3>
            <p v-if="selectedAdmin" class="text-sm text-slate-400 mt-1">
              確定要{{ toggleAction === 'disable' ? '停用' : '啟用' }} {{ selectedAdmin.name }} 的帳號嗎？
            </p>
          </div>

          <div class="p-6 border-t border-white/10 flex items-center gap-4">
            <button
              @click="closeToggleStatusModal"
              class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              取消
            </button>
            <button
              @click="confirmToggleStatus"
              :disabled="toggling"
              class="flex-1 px-4 py-2 rounded-lg transition-colors"
              :class="toggleAction === 'disable'
                ? 'bg-critical-500/30 border border-critical-500 text-critical-400 hover:bg-critical-500/40'
                : 'bg-success-500/30 border border-success-500 text-success-400 hover:bg-success-500/40'"
            >
              <span v-if="toggling" class="flex items-center justify-center gap-2">
                <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </span>
              <span v-else>確認{{ toggleAction === 'disable' ? '停用' : '啟用' }}</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 角色變更確認對話框 -->
    <Teleport to="body">
      <div
        v-if="showRoleChangeModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="closeRoleChangeModal"
      >
        <div class="glass-card w-full max-w-sm">
          <div class="p-6 border-b border-white/10">
            <h3 class="text-lg font-semibold text-white">變更角色</h3>
            <p v-if="selectedAdmin" class="text-sm text-slate-400 mt-1">
              確定要將 {{ selectedAdmin.name }} 的角色變更為 {{ newRoleText }} 嗎？
            </p>
          </div>

          <div class="p-6 border-t border-white/10 flex items-center gap-4">
            <button
              @click="closeRoleChangeModal"
              class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              取消
            </button>
            <button
              @click="confirmRoleChange"
              :disabled="changingRole"
              class="flex-1 px-4 py-2 rounded-lg bg-primary-500/30 border border-primary-500 text-primary-400 hover:bg-primary-500/40 transition-colors"
            >
              <span v-if="changingRole" class="flex items-center justify-center gap-2">
                <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </span>
              <span v-else>確認變更</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 新增管理員對話框 -->
    <Teleport to="body">
      <div
        v-if="showCreateModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="showCreateModal = false"
      >
        <div class="glass-card w-full max-w-md">
          <div class="p-6 border-b border-white/10">
            <h3 class="text-lg font-semibold text-white">新增管理員</h3>
            <p class="text-sm text-slate-400 mt-1">
              直接建立管理員帳號
            </p>
          </div>

          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm text-slate-400 mb-2">姓名 <span class="text-critical-500">*</span></label>
              <BaseInput
                v-model="createForm.name"
                placeholder="請輸入姓名"
                class="w-full"
              />
            </div>

            <div>
              <label class="block text-sm text-slate-400 mb-2">Email <span class="text-critical-500">*</span></label>
              <BaseInput
                v-model="createForm.email"
                type="email"
                placeholder="請輸入 Email"
                class="w-full"
              />
            </div>

            <div>
              <label class="block text-sm text-slate-400 mb-2">密碼 <span class="text-critical-500">*</span></label>
              <BaseInput
                v-model="createForm.password"
                type="password"
                placeholder="請輸入密碼（至少6個字元）"
                class="w-full"
              />
            </div>

            <div>
              <label class="block text-sm text-slate-400 mb-2">角色</label>
              <select
                v-model="createForm.role"
                class="w-full px-3 py-2 rounded-lg text-sm bg-slate-800/80 border border-white/10 text-slate-300 focus:outline-none focus:border-primary-500"
              >
                <option value="ADMIN">管理員</option>
                <option value="STAFF">員工</option>
              </select>
            </div>
          </div>

          <div class="p-6 border-t border-white/10 flex items-center gap-4">
            <button
              @click="showCreateModal = false"
              class="flex-1 px-4 py-2 rounded-lg bg-white/5 text-white hover:bg-white/10 transition-colors"
            >
              取消
            </button>
            <button
              @click="createAdmin"
              :disabled="creating"
              class="flex-1 px-4 py-2 rounded-lg bg-primary-500/30 border border-primary-500 text-primary-400 hover:bg-primary-500/40 transition-colors"
            >
              <span v-if="creating" class="flex items-center justify-center gap-2">
                <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </span>
              <span v-else>新增管理員</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import BaseGlassCard from '~/components/base/BaseGlassCard.vue'
import BaseInput from '~/components/base/BaseInput.vue'
import { alertError, alertSuccess } from '~/composables/useAlert'

definePageMeta({
  middleware: 'auth-admin',
  layout: 'admin',
})

const config = useRuntimeConfig()
const API_BASE = config.public.apiBase

// 管理員列表
const adminList = ref<any[]>([])
const loading = ref(true)

// 目前登入的管理員資訊
const currentAdminId = ref<number | null>(null)
const currentAdminRole = ref<string>('')

// 篩選條件
const selectedRole = ref('')
const selectedStatus = ref('')

// 篩選後的列表
const filteredAdminList = computed(() => {
  let result = adminList.value

  if (selectedRole.value) {
    result = result.filter(admin => admin.role === selectedRole.value)
  }

  if (selectedStatus.value) {
    result = result.filter(admin => admin.status === selectedStatus.value)
  }

  return result
})

// 重設密碼對話框
const showResetPasswordModal = ref(false)
const selectedAdmin = ref<any>(null)
const resetPasswordForm = ref({
  newPassword: '',
  confirmPassword: '',
})
const resetting = ref(false)

// 停用/啟用對話框
const showToggleStatusModal = ref(false)
const toggleAction = ref<'disable' | 'enable'>('disable')
const toggling = ref(false)

// 角色變更對話框
const showRoleChangeModal = ref(false)
const pendingRoleChange = ref<{ admin: any; newRole: string } | null>(null)
const changingRole = ref(false)

// 角色文字
const roleText = (role: string) => {
  const roles: Record<string, string> = {
    OWNER: '擁有者',
    ADMIN: '管理員',
    STAFF: '員工',
  }
  return roles[role] || role
}

// 新角色文字
const newRoleText = computed(() => {
  if (pendingRoleChange.value) {
    return roleText(pendingRoleChange.value.newRole)
  }
  return ''
})

// 清除篩選
const clearFilters = () => {
  selectedRole.value = ''
  selectedStatus.value = ''
}

// 檢查是否可以管理該管理員
const canManageAdmin = (admin: any) => {
  // 只能管理非 OWNER 的帳號，且不能管理自己
  if (admin.role === 'OWNER') return false
  if (admin.id === currentAdminId.value) return false
  // 只有 OWNER 可以管理
  if (currentAdminRole.value !== 'OWNER') return false
  return true
}

// 取得管理員列表
const fetchAdminList = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/admins`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (response.ok) {
      const data = await response.json()
      adminList.value = data.datas || []
    }
  } catch (err) {
    console.error('取得管理員列表失敗:', err)
    await alertError('載入管理員列表失敗')
  } finally {
    loading.value = false
  }
}

// 取得目前管理員資料
const fetchCurrentAdmin = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/me/profile`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (response.ok) {
      const data = await response.json()
      currentAdminId.value = data.datas.id
      currentAdminRole.value = data.datas.role
    }
  } catch (err) {
    console.error('取得目前管理員資料失敗:', err)
  }
}

// 開啟重設密碼對話框
const openResetPasswordModal = (admin: any) => {
  selectedAdmin.value = admin
  resetPasswordForm.value = { newPassword: '', confirmPassword: '' }
  showResetPasswordModal.value = true
}

// 關閉重設密碼對話框
const closeResetPasswordModal = () => {
  showResetPasswordModal.value = false
  selectedAdmin.value = null
}

// 重設密碼驗證
const isResetPasswordValid = computed(() => {
  return (
    resetPasswordForm.value.newPassword.length >= 6 &&
    resetPasswordForm.value.newPassword === resetPasswordForm.value.confirmPassword
  )
})

// 執行重設密碼
const resetPassword = async () => {
  if (!isResetPasswordValid.value || !selectedAdmin.value) return

  resetting.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/admins/reset-password`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        target_admin_id: selectedAdmin.value.id,
        new_password: resetPasswordForm.value.newPassword,
      }),
    })

    const data = await response.json()

    if (response.ok) {
      await alertSuccess('密碼已重設')
      closeResetPasswordModal()
    } else {
      await alertError(data.message || '重設密碼失敗')
    }
  } catch (err) {
    await alertError('重設密碼失敗，請稍後再試')
  } finally {
    resetting.value = false
  }
}

// 開啟停用/啟用對話框
const toggleAdminStatus = (admin: any) => {
  selectedAdmin.value = admin
  toggleAction.value = admin.status === 'ACTIVE' ? 'disable' : 'enable'
  showToggleStatusModal.value = true
}

// 關閉停用/啟用對話框
const closeToggleStatusModal = () => {
  showToggleStatusModal.value = false
  selectedAdmin.value = null
}

// 確認停用/啟用
const confirmToggleStatus = async () => {
  if (!selectedAdmin.value) return

  toggling.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const newStatus = toggleAction.value === 'disable' ? 'INACTIVE' : 'ACTIVE'

    const response = await fetch(`${API_BASE}/admin/admins/toggle-status`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        target_admin_id: selectedAdmin.value.id,
        new_status: newStatus,
      }),
    })

    const data = await response.json()

    if (response.ok) {
      await alertSuccess(newStatus === 'ACTIVE' ? '已啟用管理員' : '已停用管理員')
      closeToggleStatusModal()
      await fetchAdminList() // 重新載入列表
    } else {
      await alertError(data.message || '操作失敗')
    }
  } catch (err) {
    await alertError('操作失敗，請稍後再試')
  } finally {
    toggling.value = false
  }
}

// 變更角色
const changeRole = (admin: any) => {
  pendingRoleChange.value = { admin, newRole: admin.role }
  showRoleChangeModal.value = true
}

// 關閉角色變更對話框
const closeRoleChangeModal = () => {
  showRoleChangeModal.value = false
  pendingRoleChange.value = null
  // 重新載入列表以還原顯示
  fetchAdminList()
}

// 確認角色變更
const confirmRoleChange = async () => {
  if (!pendingRoleChange.value) return

  changingRole.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/admins/change-role`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        target_admin_id: pendingRoleChange.value.admin.id,
        new_role: pendingRoleChange.value.newRole,
      }),
    })

    const data = await response.json()

    if (response.ok) {
      await alertSuccess(`已將 ${pendingRoleChange.value.admin.name} 變更為 ${roleText(pendingRoleChange.value.newRole)}`)
      closeRoleChangeModal()
      await fetchAdminList() // 重新載入列表
    } else {
      await alertError(data.message || '角色變更失敗')
    }
  } catch (err) {
    await alertError('角色變更失敗，請稍後再試')
  } finally {
    changingRole.value = false
  }
}

// 頁面載入時取得資料
onMounted(async () => {
  await Promise.all([fetchAdminList(), fetchCurrentAdmin()])
})

// 新增管理員相關
const showCreateModal = ref(false)
const createForm = ref({
  email: '',
  name: '',
  role: 'ADMIN',
  password: '',
})
const creating = ref(false)

// 檢查是否可以新增管理員（僅 OWNER）
const canAddAdmin = computed(() => {
  return currentAdminRole.value === 'OWNER'
})

// 新增管理員
const createAdmin = async () => {
  if (!createForm.value.email || !createForm.value.name || !createForm.value.password) {
    await alertError('請填寫完整資訊')
    return
  }

  if (createForm.value.password.length < 6) {
    await alertError('密碼至少需要 6 個字元')
    return
  }

  creating.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/admins`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: createForm.value.email,
        name: createForm.value.name,
        role: createForm.value.role,
        password: createForm.value.password,
      }),
    })

    const data = await response.json()

    if (response.ok) {
      await alertSuccess(`管理員 ${createForm.value.name} 已成功建立`)
      showCreateModal.value = false
      createForm.value = { email: '', name: '', role: 'ADMIN', password: '' }
      await fetchAdminList()
    } else {
      await alertError(data.message || '新增失敗')
    }
  } catch (err) {
    await alertError('新增失敗，請稍後再試')
  } finally {
    creating.value = false
  }
}
</script>
