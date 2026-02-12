<template>
  <Teleport to="body">
    <div
      v-if="showCreateModal || showEditModal"
      class="fixed inset-0 z-[1000] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm isolate"
      @click.self="handleClose"
    >
      <div class="glass-card w-full max-w-lg sm:max-w-xl max-h-[90vh] overflow-y-auto animate-spring" @click.stop>
        <div class="flex items-center justify-between p-4 border-b border-white/10 sticky top-0 bg-slate-900/95 backdrop-blur-sm z-10">
          <h3 class="text-lg font-semibold text-slate-100">
            {{ showEditModal ? $t('schedule.editScheduleRule') : $t('schedule.addScheduleRule') }}
          </h3>
          <button @click="handleClose" class="p-2 rounded-lg hover:bg-white/10 transition-colors">
            <Icon icon="close" size="lg" />
          </button>
        </div>

        <!-- 載入中 -->
        <div v-if="dataLoading" class="p-8 text-center">
          <div class="inline-flex items-center gap-2 text-slate-400">
            <svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ $t('common.loading') }}</span>
          </div>
        </div>

        <!-- 衝突提示 -->
        <ConflictAlert
          v-if="conflictError"
          :conflict-messages="conflictErrors"
          :conflict-type="conflictType"
          @dismiss="clearConflictError"
        />

        <!-- 自定義 Alert 彈窗 -->
        <div v-else-if="showAlert" class="p-4">
          <div class="glass-card p-6">
            <div class="flex items-start gap-4">
              <!-- 警告圖示 -->
              <div v-if="alertType === 'warning'" class="w-10 h-10 rounded-full bg-warning-500/20 flex items-center justify-center flex-shrink-0">
                <Icon icon="warning" size="xl" class="text-warning-500" />
              </div>
              <!-- 錯誤圖示 -->
              <div v-else-if="alertType === 'error'" class="w-10 h-10 rounded-full bg-critical-500/20 flex items-center justify-center flex-shrink-0">
                <Icon icon="error" size="xl" class="text-critical-500" />
              </div>
              <!-- 資訊圖示 -->
              <div v-else class="w-10 h-10 rounded-full bg-primary-500/20 flex items-center justify-center flex-shrink-0">
                <Icon icon="info" size="xl" class="text-primary-500" />
              </div>

              <div class="flex-1">
                <h4 v-if="alertTitle" class="font-medium text-white mb-2">{{ alertTitle }}</h4>
                <p class="text-sm text-slate-400">{{ alertMessage }}</p>
              </div>
            </div>

            <div class="flex gap-3 mt-6">
              <button
                @click="handleAlertConfirm"
                class="flex-1 py-2.5 rounded-xl font-medium transition-all"
                :class="alertType === 'error' ? 'btn-critical' : 'btn-primary'"
              >
                {{ $t('common.confirm') }}
              </button>
            </div>
          </div>
        </div>

        <!-- 錯誤訊息 -->
        <div v-else-if="error" class="p-8 text-center">
          <div class="text-critical-500 mb-2">
            <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <p class="text-slate-300 mb-4">{{ error }}</p>
          <button @click="fetchData" class="btn-primary px-4 py-2 text-sm">
            {{ $t('common.retry') }}
          </button>
        </div>

        <!-- ScheduleRuleForm 表單組件 -->
        <div v-show="!dataLoading && !error" class="p-4">
          <ScheduleRuleForm
            ref="formRef"
            :editing-rule="editingRule"
            :update-mode="updateMode"
            :validation-loading="validationLoading"
            @cancel="handleClose"
            @submit="handleFormSubmit"
          />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import ScheduleRuleForm from './ScheduleRuleForm.vue'
import ConflictAlert from './ConflictAlert.vue'
import Icon from '~/components/base/Icon.vue'

// Alert composable
const { confirm: alertConfirm } = useAlert()

// Props
const props = defineProps<{
  editingRule?: any | null
  updateMode?: string
}>()

const emit = defineEmits<{
  close: []
  submit: [formData: any, updateMode: string]
  created: []
}>()

const loading = ref(false)
const dataLoading = ref(true)
const error = ref<string | null>(null)
const conflictError = ref<string | null>(null)
const conflictErrors = ref<string[]>([])
const conflictType = ref<'hard' | 'buffer'>('hard')
const showAlert = ref(false)
const alertTitle = ref('')
const alertMessage = ref('')
const alertType = ref<'info' | 'warning' | 'error'>('info')
const showCreateModal = computed(() => !props.editingRule)
const showEditModal = computed(() => !!props.editingRule)
const formRef = ref<InstanceType<typeof ScheduleRuleForm> | null>(null)

// 驗證相關
const validationLoading = ref(false)
const { quickValidate, formatConflictMessages, hasHardConflicts } = useSchedulingValidation()
const { getCenterId } = useCenterId()

const fetchData = async () => {
  dataLoading.value = true
  error.value = null

  try {
    const { fetchAllResources } = useResourceCache()
    await fetchAllResources()
  } catch (err: any) {
    console.error('Failed to fetch data:', err)
    error.value = err.message || $t('validation.loadDataFailed')
  } finally {
    dataLoading.value = false
  }
}

// 清除衝突錯誤
const clearConflictError = () => {
  conflictError.value = null
  conflictErrors.value = []
  conflictType.value = 'hard'
}

// 自定義 Alert 函數
const customAlert = (message: string, title?: string, type: 'info' | 'warning' | 'error' = 'info') => {
  alertTitle.value = title || (type === 'error' ? '操作失敗' : type === 'warning' ? '提醒' : '提示')
  alertMessage.value = message
  alertType.value = type
  showAlert.value = true
}

// Alert 確認處理
const handleAlertConfirm = () => {
  showAlert.value = false
}

const handleClose = () => {
  emit('close')
}

/**
 * 執行衝突驗證
 *
 * @param formData 表單資料
 * @returns 驗證通過返回 true，有硬衝突返回 false，有緩衝衝突返回 'buffer_warning'
 */
const validateScheduleConflicts = async (
  formData: Record<string, unknown>
): Promise<boolean | 'buffer_warning'> => {
  // 如果沒有選擇老師或教室，跳過衝突檢查
  if (!formData.teacher_id && !formData.room_id) {
    return true
  }

  validationLoading.value = true

  try {
    const centerId = getCenterId.value

    if (!centerId) {
      console.warn('無法取得中心 ID，跳過驗證')
      return true
    }

    const weekdays = formData.weekdays as number[]

    // 對每個選擇的星期進行驗證
    for (const weekday of weekdays) {
      const result = await quickValidate({
        center_id: centerId,
        teacher_id: formData.teacher_id as number | undefined,
        room_id: (formData.room_id as number) || 0,
        date: (formData.start_date as string)?.split(/[T ]/)[0],
        start_time: formData.start_time as string,
        end_time: formData.end_time as string,
        rule_id: props.editingRule?.id,
        course_id: parseInt(formData.offering_id as string),
      })

      // 確保 result 存在且有 conflicts 陣列
      if (!result || !result.conflicts) {
        console.warn('驗證結果格式異常，跳過衝突檢查')
        continue
      }

      // 獲取衝突列表（後端可能返回 hard_conflicts 或只在 conflicts 中）
      const allConflicts = result.conflicts || []
      const hardConflicts = result.hard_conflicts || allConflicts.filter(c =>
        c.type === 'OVERLAP' || c.type === 'TEACHER_OVERLAP' || c.type === 'ROOM_OVERLAP'
      )
      const bufferConflicts = result.buffer_conflicts || allConflicts.filter(c =>
        c.type === 'TEACHER_BUFFER' || c.type === 'ROOM_BUFFER'
      )

      // 如果有硬衝突，直接顯示錯誤
      if (hasHardConflicts(hardConflicts)) {
        conflictError.value = $t('validation.overlapConflict')
        conflictErrors.value = formatConflictMessages(hardConflicts)
        conflictType.value = 'hard'
        return false
      }

      // 如果有緩衝衝突，顯示警告
      if (bufferConflicts.length > 0) {
        conflictError.value = $t('validation.bufferInsufficient')
        conflictErrors.value = formatConflictMessages(bufferConflicts)
        conflictType.value = 'buffer'
        // 返回緩衝警告，不阻擋提交但需要使用者確認
        return 'buffer_warning'
      }
    }

    return true
  } catch (error: any) {
    console.error('驗證失敗:', error)
    // 驗證失敗時允許提交，讓後端處理
    return true
  } finally {
    validationLoading.value = false
  }
}

/**
 * 處理表單提交
 *
 * 包含驗證邏輯與 API 呼叫
 */
const handleFormSubmit = async (formData: Record<string, unknown>, updateMode: string) => {
  validationLoading.value = true

  try {
    const centerId = getCenterId.value

    // 對每個選擇的星期進行驗證
    let hasHardConflict = false
    let hasBufferWarning = false
    let bufferConflictMessages: string[] = []

    const weekdays = formData.weekdays as number[]

    for (const weekday of weekdays) {
      const result = await quickValidate({
        center_id: centerId || 1,
        teacher_id: formData.teacher_id as number | undefined,
        room_id: (formData.room_id as number) || 0,
        date: (formData.start_date as string)?.split(/[T ]/)[0],
        start_time: formData.start_time as string,
        end_time: formData.end_time as string,
        rule_id: props.editingRule?.id,
        course_id: parseInt(formData.offering_id as string),
      })

      // 確保 result 存在且有 conflicts 陣列
      if (!result || !result.conflicts) {
        console.warn('驗證結果格式異常，跳過衝突檢查')
        continue
      }

      // 獲取衝突列表（後端可能返回 hard_conflicts 或只在 conflicts 中）
      const allConflicts = result.conflicts || []
      const hardConflicts = result.hard_conflicts || allConflicts.filter(c =>
        c.type === 'OVERLAP' || c.type === 'TEACHER_OVERLAP' || c.type === 'ROOM_OVERLAP'
      )
      const bufferConflicts = result.buffer_conflicts || allConflicts.filter(c =>
        c.type === 'TEACHER_BUFFER' || c.type === 'ROOM_BUFFER'
      )

      // 如果有硬衝突，顯示錯誤
      if (hasHardConflicts(hardConflicts)) {
        conflictError.value = '排課時間與現有規則衝突'
        conflictErrors.value = formatConflictMessages(hardConflicts)
        conflictType.value = 'hard'
        hasHardConflict = true
        break
      }

      // 如果有緩衝衝突，記錄下來
      if (bufferConflicts.length > 0) {
        hasBufferWarning = true
        bufferConflictMessages = formatConflictMessages(bufferConflicts)
        conflictType.value = 'buffer'
      }

    }

    // 有硬衝突，停止提交並顯示錯誤通知
    if (hasHardConflict) {
      validationLoading.value = false
      // 使用 alertConfirm 確保使用者看到衝突訊息
      await alertConfirm(
        `排課時間與現有規則衝突：\n${conflictErrors.value.join('\n')}`,
        '時間衝突',
        'error'
      )
      // 清除衝突狀態，讓使用者可以修改表單
      clearConflictError()
      return
    }

    // 有緩衝衝突，詢問是否覆寫
    if (hasBufferWarning) {
      const shouldOverride = await alertConfirm(
        `${$t('validation.bufferWarning')}:\n${bufferConflictMessages.join('\n')}\n\n${$t('validation.confirmOverride')}`
      )

      if (!shouldOverride) {
        validationLoading.value = false
        return
      }
    }

    // 驗證全部通過，顯示成功通知並等待確認
    if (!hasHardConflict && !hasBufferWarning) {
      const shouldContinue = await alertConfirm(
        '排課時間驗證通過，沒有發現時間衝突。\n\n是否繼續提交？',
        '驗證成功',
        'info'
      )

      if (!shouldContinue) {
        validationLoading.value = false
        return
      }
    }

    // 驗證通過，繼續提交
    validationLoading.value = false

  } catch (error: any) {
    console.error('驗證失敗:', error)
    validationLoading.value = false

    // 驗證失敗時，顯示錯誤訊息並停止提交
    // 使用 alertConfirm 確保使用者看到錯誤訊息
    await alertConfirm(
      error.response?.data?.message || error.message || '驗證失敗，請稍後再試',
      '驗證失敗',
      'error'
    )
    return
  }

  loading.value = true

  try {
    // 準備提交資料
    const submitData: Record<string, unknown> = {
      name: formData.name,
      offering_id: parseInt(formData.offering_id as string),
      start_time: formData.start_time,
      end_time: formData.end_time,
      duration: formData.duration,
      weekdays: formData.weekdays,
      status: formData.status, // 新增：課程狀態
      start_date: (formData.start_date as string)?.split(/[T ]/)[0],
      end_date: (formData.end_date as string)?.split(/[T ]/)[0] || null,
      suspended_dates: formData.suspended_dates || [],
    }

    // 只有當有選擇老師時才傳送
    if (formData.teacher_id) {
      submitData.teacher_id = formData.teacher_id
    }

    // 只有當有選擇教室時才傳送
    if (formData.room_id) {
      submitData.room_id = formData.room_id
    }

    // 編輯模式：處理日期欄位
    if (showEditModal.value) {
      // 如果日期為空，從 data 中移除，讓後端保留現有值
      if (!submitData.start_date) {
        delete submitData.start_date
      }
      if (!submitData.end_date) {
        delete submitData.end_date
      }
      // 發射表單資料給父元件處理
      emit('submit', submitData, updateMode)
      handleClose()
    } else {
      // 新增模式：直接呼叫 API
      const api = useApi()
      await api.post('/admin/rules', submitData)
      handleClose()
      // 通知父元件刷新列表
      emit('created')
    }
  } catch (error: any) {
    console.error('Failed to save schedule rule:', error)

    // 處理衝突錯誤
    if (error.response?.data?.code === 40002 || error.response?.data?.code === 'OVERLAP' || error.response?.data?.code === 20002) {
      conflictError.value = error.response?.data?.message || $t('validation.overlapConflict')
      conflictErrors.value = error.response?.data?.datas?.conflicts || []
      conflictType.value = 'hard'
    } else {
      // 使用 alertConfirm 確保使用者看到錯誤訊息
      await alertConfirm(
        error.response?.data?.message || error.message || $t('common.operationFailed'),
        showEditModal.value ? $t('common.operationFailed') : $t('common.operationFailed'),
        'error'
      )
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
