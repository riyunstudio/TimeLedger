<template>
   <div class="min-h-screen bg-slate-900">
 
     <main class="p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <h1 class="text-2xl font-bold text-white">課程時段管理</h1>
        <button
          @click="showModal = true"
          class="px-4 py-2 rounded-lg bg-primary-500 text-white hover:bg-primary-600 transition-colors"
        >
          新增時段
        </button>
      </div>

      <div class="glass-card p-6">
        <div v-if="loading" class="text-center py-8 text-slate-400">
          載入中...
        </div>

        <div v-else-if="rules.length === 0" class="text-center py-8 text-slate-400">
          尚未建立課程時段
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="text-left text-slate-400 text-sm border-b border-white/10">
                <th class="pb-3 pl-2">課程</th>
                <th class="pb-3">星期</th>
                <th class="pb-3">時間</th>
                <th class="pb-3">教室</th>
                <th class="pb-3">老師</th>
                <th class="pb-3">狀態</th>
                <th class="pb-3 pr-2 text-right">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="rule in rules"
                :key="rule.id"
                class="border-b border-white/5 hover:bg-white/5"
              >
                <td class="py-3 pl-2 text-white">{{ rule.offering?.name || '-' }}</td>
                <td class="py-3 text-slate-300">{{ getWeekdayText(rule.weekday) }}</td>
                <td class="py-3 text-slate-300">{{ rule.start_time }} - {{ rule.end_time }}</td>
                <td class="py-3 text-slate-300">{{ rule.room?.name || '-' }}</td>
                <td class="py-3 text-slate-300">{{ rule.teacher?.name || '-' }}</td>
                <td class="py-3">
                  <span
                    class="px-2 py-1 rounded-full text-xs"
                    :class="getStatusClass(rule)"
                  >
                    {{ getStatusText(rule) }}
                  </span>
                </td>
                <td class="py-3 pr-2 text-right">
                  <button
                    @click="deleteRule(rule.id)"
                    class="text-critical-500 hover:text-critical-400"
                  >
                    刪除
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>

    <ScheduleRuleModal
      v-if="showModal"
      @close="showModal = false"
      @saved="fetchRules"
    />

    <NotificationDropdown
      v-if="notificationUI.show.value"
      @close="notificationUI.close()"
    />
  </div>
</template>

<script setup lang="ts">
 definePageMeta({
   middleware: 'auth-admin',
   layout: 'admin',
 })

 const notificationUI = useNotification()
const showModal = ref(false)
const loading = ref(true)
const rules = ref<any[]>([])

const fetchRules = async () => {
  loading.value = true
  try {
    const api = useApi()
    const response = await api.get<{ code: number; datas: any[] }>(`/admin/scheduling/rules`)
    rules.value = response.datas || []
  } catch (error) {
    console.error('Failed to fetch rules:', error)
  } finally {
    loading.value = false
  }
}

const deleteRule = async (id: number) => {
  if (!confirm('確定要刪除此課程時段？')) return

  try {
    const api = useApi()
    await api.delete(`/admin/scheduling/rules/${id}`)
    await fetchRules()
  } catch (error) {
    console.error('Failed to delete rule:', error)
    alert('刪除失敗')
  }
}

const getWeekdayText = (weekday: number): string => {
  const days = ['日', '一', '二', '三', '四', '五', '六']
  return days[weekday] || '-'
}

const getStatusClass = (rule: any): string => {
  const now = new Date()
  const startDate = new Date(rule.effective_range?.start_date)
  const endDate = rule.effective_range?.end_date ? new Date(rule.effective_range.end_date) : null

  if (endDate && now > endDate) return 'bg-slate-500/20 text-slate-400'
  if (now < startDate) return 'bg-primary-500/20 text-primary-500'
  return 'bg-success-500/20 text-success-500'
}

const getStatusText = (rule: any): string => {
  const now = new Date()
  const startDate = new Date(rule.effective_range?.start_date)
  const endDate = rule.effective_range?.end_date ? new Date(rule.effective_range.end_date) : null

  if (endDate && now > endDate) return '已結束'
  if (now < startDate) return '尚未開始'
  return '進行中'
}

onMounted(() => {
  fetchRules()
})
</script>
