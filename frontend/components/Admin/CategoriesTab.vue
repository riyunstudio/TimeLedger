<template>
  <div class="space-y-6">
    <!-- 標題與新增按鈕 -->
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-bold text-slate-100">課程類別管理</h2>
      <button
        @click="openAddModal"
        class="px-4 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        <span>新增類別</span>
      </button>
    </div>

    <!-- 說明文字 -->
    <div class="bg-blue-500/10 border border-blue-500/20 rounded-xl p-4">
      <p class="text-sm text-blue-300">
        課程類別可用於對課程進行分類管理，建立課程時可選擇所屬類別。於排課規則中也可依類別篩選。
      </p>
    </div>

    <!-- 類別列表 -->
    <div class="bg-white/5 rounded-xl overflow-hidden">
      <div v-if="loading && categories.length === 0" class="flex items-center justify-center py-12">
        <div class="animate-spin w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full"></div>
      </div>

      <div v-else-if="categories.length === 0" class="text-center py-12">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-slate-700/50 flex items-center justify-center">
          <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
          </svg>
        </div>
        <p class="text-slate-400 mb-4">尚未建立任何課程類別</p>
        <button
          @click="openAddModal"
          class="text-primary-400 hover:text-primary-300 transition-colors"
        >
          點擊新增第一個類別
        </button>
      </div>

      <table v-else class="w-full">
        <thead class="bg-white/5">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
              類別名稱
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-slate-400 uppercase tracking-wider">
              操作
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="category in categories" :key="category.id" class="hover:bg-white/5 transition-colors">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center gap-2">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-500/20 text-primary-400">
                  {{ category.name }}
                </span>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button
                @click="openEditModal(category)"
                class="text-primary-400 hover:text-primary-300 transition-colors mr-4"
              >
                編輯
              </button>
              <button
                @click="openDeleteConfirm(category)"
                class="text-red-400 hover:text-red-300 transition-colors"
              >
                刪除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 新增/編輯類別 Modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      @click.self="closeModal"
    >
      <div class="bg-slate-800 rounded-2xl w-full max-w-md p-6 border border-white/10">
        <h3 class="text-lg font-semibold text-slate-100 mb-4">
          {{ modalMode === 'add' ? '新增類別' : '編輯類別' }}
        </h3>

        <form @submit.prevent="handleSubmit">
          <div class="space-y-4">
            <div>
              <label class="block text-sm text-slate-300 mb-2">類別名稱</label>
              <input
                v-model="formData.name"
                type="text"
                placeholder="例如：音樂、語言、運動"
                class="w-full px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-primary-500"
                required
              />
            </div>
          </div>

          <div class="flex justify-end gap-3 mt-6">
            <button
              type="button"
              @click="closeModal"
              class="px-4 py-2 bg-white/5 text-slate-300 rounded-xl hover:bg-white/10 transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="saving"
              class="px-4 py-2 bg-primary-500/30 border border-primary-500 text-primary-400 rounded-xl hover:bg-primary-500/40 transition-colors disabled:opacity-50"
            >
              {{ saving ? '儲存中...' : '儲存' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 刪除確認 Modal -->
    <div
      v-if="showDeleteModal"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      @click.self="showDeleteModal = false"
    >
      <div class="bg-slate-800 rounded-2xl w-full max-w-md p-6 border border-white/10">
        <h3 class="text-lg font-semibold text-slate-100 mb-4">
          確認刪除
        </h3>
        <p class="text-slate-300 mb-6">
          確定要刪除類別「<span class="text-white font-medium">{{ deleteTarget?.name }}</span>」嗎？
        </p>
        <div class="flex justify-end gap-3">
          <button
            @click="showDeleteModal = false"
            class="px-4 py-2 bg-white/5 text-slate-300 rounded-xl hover:bg-white/10 transition-colors"
          >
            取消
          </button>
          <button
            @click="handleDelete"
            :disabled="saving"
            class="px-4 py-2 bg-red-500/30 border border-red-500 text-red-400 rounded-xl hover:bg-red-500/40 transition-colors disabled:opacity-50"
          >
            {{ saving ? '刪除中...' : '確認刪除' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertError, alertSuccess } from '~/composables/useAlert'

const { getCenterId } = useCenterId()

const config = useRuntimeConfig()
const API_BASE = config.public.apiBase

const loading = ref(false)
const saving = ref(false)

// Category 類型定義
interface Category {
  id: number
  center_id: number
  name: string
}

const categories = ref<Category[]>([])

// Modal 狀態
const showModal = ref(false)
const modalMode = ref<'add' | 'edit'>('add')
const editingCategory = ref<Category | null>(null)

const formData = ref({
  name: '',
})

// 刪除 Modal 狀態
const showDeleteModal = ref(false)
const deleteTarget = ref<Category | null>(null)

// 取得類別列表
const fetchCategories = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${API_BASE}/admin/course-categories`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (response.ok) {
      const data = await response.json()
      categories.value = data.datas || []
    }
  } catch (err) {
    console.error('取得課程類別失敗:', err)
    alertError('取得課程類別失敗')
  } finally {
    loading.value = false
  }
}

// 打開新增 Modal
const openAddModal = () => {
  modalMode.value = 'add'
  formData.value.name = ''
  editingCategory.value = null
  showModal.value = true
}

// 打開編輯 Modal
const openEditModal = (category: Category) => {
  modalMode.value = 'edit'
  formData.value.name = category.name
  editingCategory.value = category
  showModal.value = true
}

// 關閉 Modal
const closeModal = () => {
  showModal.value = false
  formData.value.name = ''
  editingCategory.value = null
}

// 儲存類別
const handleSubmit = async () => {
  const name = formData.value.name.trim()
  if (!name) {
    alertError('請輸入類別名稱')
    return
  }

  // 檢查是否重複
  const existing = categories.value.find(c => c.name === name)
  if (existing && (modalMode.value === 'add' || existing.id !== editingCategory.value?.id)) {
    alertError('此類別名稱已存在')
    return
  }

  saving.value = true
  try {
    const token = localStorage.getItem('admin_token')

    let response
    if (modalMode.value === 'add') {
      // 新增
      response = await fetch(`${API_BASE}/admin/course-categories`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name }),
      })
    } else {
      // 編輯
      response = await fetch(`${API_BASE}/admin/course-categories/${editingCategory.value?.id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name }),
      })
    }

    if (response.ok) {
      await fetchCategories()
      closeModal()
      alertSuccess(modalMode.value === 'add' ? '類別新增成功' : '類別更新成功')
    } else {
      const data = await response.json()
      alertError(data.message || '儲存失敗')
    }
  } catch (err) {
    console.error('儲存類別失敗:', err)
    alertError('儲存失敗，請稍後再試')
  } finally {
    saving.value = false
  }
}

// 打開刪除確認
const openDeleteConfirm = (category: Category) => {
  deleteTarget.value = category
  showDeleteModal.value = true
}

// 刪除類別
const handleDelete = async () => {
  if (!deleteTarget.value) return

  saving.value = true
  try {
    const token = localStorage.getItem('admin_token')

    const response = await fetch(`${API_BASE}/admin/course-categories/${deleteTarget.value.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (response.ok) {
      await fetchCategories()
      showDeleteModal.value = false
      deleteTarget.value = null
      alertSuccess('類別刪除成功')
    } else {
      const data = await response.json()
      alertError(data.message || '刪除失敗')
    }
  } catch (err) {
    console.error('刪除類別失敗:', err)
    alertError('刪除失敗，請稍後再試')
  } finally {
    saving.value = false
  }
}

// 初始載入
onMounted(() => {
  fetchCategories()
})
</script>
