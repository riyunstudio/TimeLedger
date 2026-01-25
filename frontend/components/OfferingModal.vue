<template>
  <div class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="emit('close')">
    <div class="glass-card w-full max-w-md max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
      <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
        <h3 class="text-lg font-semibold text-slate-100">
          {{ offering ? '編輯待排課程' : '新增待排課程' }}
        </h3>
        <button @click="emit('close')" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
          <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- 載入中 -->
      <div v-if="loading" class="p-8 text-center">
        <div class="inline-flex items-center gap-2 text-slate-400">
          <svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span>載入中...</span>
        </div>
      </div>

      <!-- 表單 -->
      <form v-else @submit.prevent="handleSubmit" class="p-4 space-y-4">
        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm">課程</label>
          <select v-model="form.course_id" class="input-field text-sm" required>
            <option value="">請選擇課程</option>
            <option v-for="course in courses" :key="course.id" :value="course.id">
              {{ course.name }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm">名稱（可選）</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="例：週一鋼琴班"
            class="input-field text-sm"
          />
          <p class="text-xs text-slate-500 mt-1">留空將自動使用課程名稱</p>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm">預設老師（可選）</label>
          <select v-model="form.default_teacher_id" class="input-field text-sm">
            <option :value="null">未指定</option>
            <option v-for="teacher in teachers" :key="teacher.id" :value="teacher.id">
              {{ teacher.name }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-slate-300 mb-2 font-medium text-sm">預設教室（可選）</label>
          <select v-model="form.default_room_id" class="input-field text-sm">
            <option :value="null">未指定</option>
            <option v-for="room in rooms" :key="room.id" :value="room.id">
              {{ room.name }}
            </option>
          </select>
        </div>

        <div class="flex items-center gap-2">
          <input
            type="checkbox"
            id="allow_buffer_override"
            v-model="form.allow_buffer_override"
            class="w-4 h-4 rounded bg-slate-800 border-slate-600 text-primary-500 focus:ring-primary-500"
          />
          <label for="allow_buffer_override" class="text-sm text-slate-300">
            允許覆蓋緩衝時間
          </label>
        </div>

        <div class="flex gap-3 pt-2">
          <button
            type="button"
            @click="emit('close')"
            class="flex-1 glass-btn py-2.5 rounded-xl font-medium text-sm"
          >
            取消
          </button>
          <button
            type="submit"
            :disabled="submitting"
            class="flex-1 btn-primary py-2.5 rounded-xl font-medium text-sm"
          >
            {{ submitting ? '儲存中...' : '儲存' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { alertError } from '~/composables/useAlert'

const props = defineProps<{
  offering?: any
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const loading = ref(true)
const submitting = ref(false)

const courses = ref<any[]>([])
const teachers = ref<any[]>([])
const rooms = ref<any[]>([])

const form = ref({
  course_id: '',
  name: '',
  default_teacher_id: null as number | null,
  default_room_id: null as number | null,
  allow_buffer_override: false,
})

const { getCenterId } = useCenterId()

const fetchData = async () => {
  try {
    const api = useApi()
    const centerId = getCenterId()

    const [coursesRes, teachersRes, roomsRes] = await Promise.all([
      api.get<{ code: number; datas: any[] }>(`/admin/courses`),
      api.get<{ code: number; datas: any[] }>('/teachers'),
      api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
    ])

    courses.value = coursesRes.datas || []
    teachers.value = teachersRes.datas || []
    rooms.value = roomsRes.datas || []

    // 如果是編輯模式，載入現有資料
    if (props.offering) {
      form.value = {
        course_id: props.offering.course_id || '',
        name: props.offering.name || '',
        default_teacher_id: props.offering.default_teacher_id || null,
        default_room_id: props.offering.default_room_id || null,
        allow_buffer_override: props.offering.allow_buffer_override || false,
      }
    }
  } catch (error) {
    console.error('Failed to fetch data:', error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true

  try {
    const api = useApi()
    const centerId = getCenterId()

    const data = {
      course_id: parseInt(form.value.course_id),
      name: form.value.name || null,
      default_teacher_id: form.value.default_teacher_id,
      default_room_id: form.value.default_room_id,
      allow_buffer_override: form.value.allow_buffer_override,
    }

    if (props.offering) {
      await api.put(`/admin/offerings/${props.offering.id}`, data)
    } else {
      await api.post(`/admin/offerings`, data)
    }

    emit('saved')
    emit('close')
  } catch (error) {
    console.error('Failed to save offering:', error)
    await alertError('儲存失敗，請稍後再試')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
