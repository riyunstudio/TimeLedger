<template>
  <div class="fixed inset-0 z-[200] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-md sm:max-w-lg max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          {{ course ? '編輯課程' : '新增課程' }}
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="p-4 space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">課程代號</label>
            <input
              v-model="form.code"
              type="text"
              placeholder="例：Piano-101"
              class="input-field text-sm sm:text-base"
            />
          </div>
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">課程名稱</label>
            <input
              v-model="form.name"
              type="text"
              placeholder="例：鋼琴基礎"
              class="input-field text-sm sm:text-base"
              required
            />
          </div>
        </div>

        <!-- 課程類別 -->
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">課程類別</label>
          <select
            v-model="form.category"
            class="input-field text-sm sm:text-base"
          >
            <option value="">無（不指定類別）</option>
            <option v-for="cat in categories" :key="cat" :value="cat">
              {{ cat }}
            </option>
          </select>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">課程時長</label>
            <div class="flex items-center gap-2">
              <input
                v-model.number="form.duration"
                type="number"
                placeholder="60"
                min="1"
                class="input-field text-sm sm:text-base"
                required
              />
              <span class="text-sm text-slate-400">分鐘</span>
            </div>
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">課程顏色</label>
            <div class="flex items-center gap-2">
              <input
                v-model="form.color_hex"
                type="color"
                class="w-10 h-10 rounded-lg cursor-pointer border border-white/20"
                required
              />
              <input
                v-model="form.color_hex"
                type="text"
                placeholder="#3B82F6"
                class="input-field text-sm sm:text-base flex-1"
                required
              />
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">老師緩衝時間</label>
            <div class="flex items-center gap-2">
              <input
                v-model.number="form.teacher_buffer_min"
                type="number"
                min="0"
                placeholder="10"
                class="input-field text-sm sm:text-base"
                required
              />
              <span class="text-sm text-slate-400">分鐘</span>
            </div>
          </div>

          <div>
            <label class="block text-slate-300 mb-2 font-medium text-sm sm:text-base">教室緩衝時間</label>
            <div class="flex items-center gap-2">
              <input
                v-model.number="form.room_buffer_min"
                type="number"
                min="0"
                placeholder="5"
                class="input-field text-sm sm:text-base"
                required
              />
              <span class="text-sm text-slate-400">分鐘</span>
            </div>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-slate-700/30 border border-white/10">
          <p class="text-sm text-slate-300">
            <span class="font-medium text-slate-400">💡 提示：</span>
            緩衝時間是連續課程之間的間隔時間
          </p>
        </div>

        <div class="flex gap-3 pt-2">
          <button
            type="button"
            @click="emit('close')"
            class="flex-1 glass-btn py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="flex-1 btn-primary py-2.5 sm:py-3 rounded-xl font-medium text-sm sm:text-base"
          >
            {{ loading ? '儲存中...' : '儲存' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertError, alertSuccess } from '~/composables/useAlert'

const props = defineProps<{
  course: any | null
  defaultDuration?: number
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const loading = ref(false)
const { getCenterId } = useCenterId()
const form = ref({
  code: props.course?.code || '',
  name: props.course?.name || '',
  category: props.course?.category || '',
  duration: props.course?.default_duration || props.defaultDuration || 60,
  color_hex: props.course?.color_hex || '#3B82F6',
  teacher_buffer_min: props.course?.teacher_buffer_min || 10,
  room_buffer_min: props.course?.room_buffer_min || 5,
})

// 從設定取得課程類別
const categories = ref<string[]>([])
const fetchCategories = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${window.location.origin}/api/v1/admin/course-categories`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    if (response.ok) {
      const data = await response.json()
      // 轉換為名稱陣列
      categories.value = (data.datas || []).map((c: any) => c.name)
    }
  } catch (err) {
    console.error('取得課程類別失敗:', err)
  }
}

onMounted(() => {
  fetchCategories()
})

watch(() => props.course, (newCourse) => {
  if (newCourse) {
    form.value = {
      code: newCourse.code || '',
      name: newCourse.name,
      category: newCourse.category || '',
      duration: newCourse.default_duration || props.defaultDuration || 60,
      color_hex: newCourse.color_hex || '#3B82F6',
      teacher_buffer_min: newCourse.teacher_buffer_min,
      room_buffer_min: newCourse.room_buffer_min,
    }
  } else {
    form.value = {
      code: '',
      name: '',
      category: '',
      duration: props.defaultDuration || 60,
      color_hex: '#3B82F6',
      teacher_buffer_min: 10,
      room_buffer_min: 5,
    }
  }
}, { immediate: true })

watch(() => props.defaultDuration, (newDuration) => {
  if (!props.course && newDuration) {
    form.value.duration = newDuration
  }
})

const handleSubmit = async () => {
  loading.value = true

  try {
    const api = useApi()
    const courseData = {
      code: form.value.code.trim() || null,
      name: form.value.name,
      category: form.value.category || null,
      duration: form.value.duration,
      color_hex: form.value.color_hex,
      teacher_buffer_min: form.value.teacher_buffer_min,
      room_buffer_min: form.value.room_buffer_min,
    }

    if (props.course && props.course.id) {
      await api.put(`/admin/courses/${props.course.id}`, courseData)
      await alertSuccess('課程更新成功')
    } else {
      await api.post(`/admin/courses`, courseData)
      await alertSuccess('課程建立成功')
    }

    emit('saved')
    emit('close')
  } catch (error) {
    console.error('Failed to save course:', error)
    await alertError('儲存失敗，請稍後再試')
  } finally {
    loading.value = false
  }
}
</script>
