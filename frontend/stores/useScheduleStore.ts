import { defineStore } from 'pinia'
import type { CenterMembership, Center, WeekSchedule, ScheduleException, SessionNote, PersonalEvent, RecurrenceRule, Invitation } from '~/types'
import {
  TeacherScheduleDataSchema,
  ScheduleExceptionListDataSchema,
} from '~/types/schemas'
import { formatDateToString } from '~/composables/useTaiwanTime'
import { expandRecurrenceEvents } from '~/composables/useRecurrence'
import { withLoading } from '~/utils/loadingHelper'

export interface TeacherScheduleItem {
  id: string
  type: string
  title: string
  date: string
  start_time: string
  end_time: string
  room_id: number
  teacher_id?: number
  center_id: number
  center_name?: string
  status: string
  rule_id?: number
  data?: any
  is_cross_day_part?: boolean
}

export const useScheduleStore = defineStore('schedule', () => {
  // 資料狀態
  const centers = ref<CenterMembership[]>([])
  const currentCenter = ref<Center | null>(null)
  const schedule = ref<WeekSchedule | null>(null)
  const exceptions = ref<ScheduleException[]>([])
  const personalEvents = ref<PersonalEvent[]>([])
  const sessionNote = ref<SessionNote | null>(null)
  const invitations = ref<Invitation[]>([])
  const pendingInvitationsCount = ref(0)
  const subscriptionUrl = ref<string | null>(null)

  // Loading 狀態
  const isLoading = ref(false)
  const isFetching = ref(false)
  const isCreating = ref(false)
  const isUpdating = ref(false)
  const isDeleting = ref(false)
  const isCreatingEvent = ref(false)
  const isUpdatingEvent = ref(false)
  const isDeletingEvent = ref(false)
  const isCreatingException = ref(false)
  const isRevokingException = ref(false)
  const isSavingNote = ref(false)
  const isRespondingInvitation = ref(false)
  const isCreatingSubscription = ref(false)
  const isDeletingSubscription = ref(false)
  const isDownloadingImage = ref(false)

  // 日期週次相關
  const getWeekStart = (date: Date): Date => {
    const d = new Date(date)
    const day = d.getDay()
    const diff = d.getDate() - day + (day === 0 ? -6 : 1)
    return new Date(d.setDate(diff))
  }

  const weekStart = ref<Date>(getWeekStart(new Date()))

  const weekEnd = computed(() => {
    if (!weekStart.value) return null
    const end = new Date(weekStart.value)
    end.setDate(end.getDate() + 6)
    return end
  })

  const weekLabel = computed(() => {
    if (!weekStart.value || !weekEnd.value) return ''
    const start = weekStart.value.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric' })
    const end = weekEnd.value.toLocaleDateString('zh-TW', { month: 'long', day: 'numeric', year: 'numeric' })
    return `${start} - ${end}`
  })

  const changeWeek = (delta: number) => {
    if (!weekStart.value) return
    // 直接加減 7 天，不再強制對齊週一，保持選中的日期為該週第一天
    const newStart = new Date(weekStart.value)
    newStart.setDate(newStart.getDate() + (delta * 7))
    weekStart.value = newStart
  }

  const formatDate = (date: Date): string => {
    return formatDateToString(date)
  }

  // 台灣時區 datetime 格式化（用於 API 傳送）
  const formatDateTimeForApi = (date: Date): string => {
    const pad = (n: number) => String(n).padStart(2, '0')
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}+08:00`
  }

  // 中心相關
  const fetchCenters = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const membershipList = await api.get<CenterMembership[]>('/teacher/me/centers')
        // 處理 null response 或空陣列的情況
        if (!membershipList || membershipList.length === 0) {
          centers.value = []
          return
        }
        centers.value = membershipList
        if (centers.value.length > 0 && !currentCenter.value && centers.value[0].center_id) {
          currentCenter.value = { id: centers.value[0].center_id, name: centers.value[0].center_name || '' } as any
        }
      } catch (error) {
        console.error('Failed to fetch centers:', error)
        centers.value = [] // 確保即使失敗也設定為空陣列
        throw error
      }
    })
  }

  // 課表相關
  const fetchSchedule = async () => {
    if (!weekStart.value || !weekEnd.value) return

    return withLoading(isLoading, async () => {
      try {
        const api = useApi()
        const response = await api.get<any[]>(
          `/teacher/me/schedule?from=${formatDate(weekStart.value!)}&to=${formatDate(weekEnd.value!)}`,
          undefined,
          undefined,
          TeacherScheduleDataSchema
        )
        schedule.value = transformToWeekSchedule(response || [])
      } catch (error) {
        console.error('Failed to fetch schedule:', error)
        throw error
      }
    })
  }

  const transformToWeekSchedule = (items: TeacherScheduleItem[]): WeekSchedule => {
    const daysMap = new Map<string, TeacherScheduleItem[]>()

    items.forEach(item => {
      const date = item.date
      if (!daysMap.has(date)) {
        daysMap.set(date, [])
      }
      daysMap.get(date)!.push(item)
    })

    const days = []
    const start = new Date(weekStart.value!)
    const end = new Date(start)
    end.setDate(end.getDate() + 6)

    for (let i = 0; i < 7; i++) {
      const d = new Date(start)
      d.setDate(d.getDate() + i)
      const dateStr = formatDateToString(d)
      const dayOfWeek = d.getDay()

      const dayItems: any[] = (daysMap.get(dateStr) || []).map(item => ({
        type: item.type as 'SCHEDULE_RULE' | 'PERSONAL_EVENT' | 'CENTER_SESSION',
        id: item.id,
        title: item.title,
        start_time: item.start_time,
        end_time: item.end_time,
        color: item.status === 'PENDING_CANCEL' ? '#F59E0B' : undefined,
        status: item.status,
        center_name: item.center_name,
        data: item.data,
        date: item.date,
        room_id: item.room_id,
        center_id: item.center_id,
        rule_id: item.rule_id,
        is_cross_day_part: item.is_cross_day_part,
      }))

      days.push({
        date: dateStr,
        day_of_week: dayOfWeek,
        items: dayItems,
      })
    }

    return {
      week_start: formatDateToString(start),
      week_end: formatDateToString(end),
      days
    } as WeekSchedule
  }

  // 例外申請相關
  const fetchExceptions = async (status?: string) => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const endpoint = status
          ? `/teacher/exceptions?status=${status}`
          : '/teacher/exceptions'
        const response = await api.get<any[]>(
          endpoint,
          undefined,
          undefined,
          ScheduleExceptionListDataSchema
        )
        exceptions.value = response || []
      } catch (error) {
        console.error('Failed to fetch exceptions:', error)
        throw error
      }
    })
  }

  const createException = async (data: {
    center_id: number
    rule_id: number
    original_date: string
    type: 'CANCEL' | 'RESCHEDULE' | 'REPLACE_TEACHER'
    new_start_at?: string
    new_end_at?: string
    new_teacher_id?: number
    new_teacher_name?: string
    reason: string
  }) => {
    return withLoading(isCreatingException, async () => {
      const api = useApi()
      const newException = await api.post<ScheduleException>('/teacher/exceptions', data)
      exceptions.value.unshift(newException)
      return newException
    })
  }

  const revokeException = async (exceptionId: number) => {
    return withLoading(isRevokingException, async () => {
      const api = useApi()
      await api.post(`/teacher/exceptions/${exceptionId}/revoke`, {})
      await fetchExceptions()
    })
  }

  // 課堂筆記相關
  const fetchSessionNote = async (ruleId: number, sessionDate: string) => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const response = await api.get<{ note: SessionNote; is_new: boolean }>(
          `/teacher/sessions/note?rule_id=${ruleId}&session_date=${sessionDate}`
        )
        sessionNote.value = response?.note || null
        return sessionNote.value
      } catch (error) {
        console.error('Failed to fetch session note:', error)
        throw error
      }
    })
  }

  const saveSessionNote = async (ruleId: number, sessionDate: string, content: string, prepNote: string) => {
    return withLoading(isSavingNote, async () => {
      try {
        const api = useApi()
        const response = await api.put<SessionNote>('/teacher/sessions/note', {
          rule_id: ruleId,
          session_date: sessionDate,
          content,
          prep_note: prepNote,
        })
        sessionNote.value = response || null
        return sessionNote.value
      } catch (error) {
        console.error('Failed to save session note:', error)
        throw error
      }
    })
  }

  // 個人行程相關
  const fetchPersonalEvents = async () => {
    if (!weekStart.value || !weekEnd.value) return

    return withLoading(isLoading, async () => {
      try {
        const api = useApi()
        const events = await api.get<PersonalEvent[]>(
          `/teacher/me/personal-events?from=${formatDate(weekStart.value!)}&to=${formatDate(weekEnd.value!)}`
        ) || []

        const currentWeekStart = new Date(weekStart.value!)
        const currentWeekEnd = new Date(weekEnd.value!)
        currentWeekEnd.setHours(23, 59, 59, 999)

        // 使用 composable 展開循環事件
        const expandedEvents = expandRecurrenceEvents(events, currentWeekStart, currentWeekEnd)

        personalEvents.value = expandedEvents.filter(event => {
          const eventStart = new Date(event.start_at)
          return eventStart >= currentWeekStart && eventStart <= currentWeekEnd
        }) as any[]
      } catch (error) {
        console.error('Failed to fetch personal events:', error)
        throw error
      }
    })
  }

  const createPersonalEvent = async (data: {
    title: string
    start_at: string
    end_at: string
    is_all_day?: boolean
    color_hex?: string
    recurrence_rule?: RecurrenceRule
  }) => {
    return withLoading(isCreatingEvent, async () => {
      const api = useApi()
      const newEvent = await api.post<PersonalEvent>('/teacher/me/personal-events', data)
      if (!newEvent) {
        throw new Error('建立個人行程失敗：無法取得回傳資料')
      }
      if (data.recurrence_rule && !newEvent.recurrence_rule) {
        newEvent.recurrence_rule = data.recurrence_rule
      }
      personalEvents.value.push(newEvent)
      return newEvent
    })
  }

  const deletePersonalEvent = async (eventId: number | string) => {
    return withLoading(isDeletingEvent, async () => {
      const api = useApi()
      const originalId = typeof eventId === 'string' && eventId.includes('_')
        ? parseInt(eventId.split('_')[0])
        : eventId
      await api.delete(`/teacher/me/personal-events/${originalId}`)
      personalEvents.value = personalEvents.value.filter(e => e.id !== eventId)
    })
  }

  const updatePersonalEvent = async (eventId: number | string, data: {
    title: string
    start_at: string
    end_at: string
    color_hex?: string
    recurrence_rule?: RecurrenceRule
  }) => {
    return withLoading(isUpdatingEvent, async () => {
      const api = useApi()
      const originalId = typeof eventId === 'string' && eventId.includes('_')
        ? parseInt(eventId.split('_')[0])
        : eventId
      const updateData = { ...data, update_mode: 'SINGLE' }
      const updatedEvent = await api.patch<PersonalEvent>(`/teacher/me/personal-events/${originalId}`, updateData)
      if (!updatedEvent) {
        throw new Error('更新個人行程失敗：無法取得回傳資料')
      }
      const index = personalEvents.value.findIndex(e => e.id === eventId)
      if (index !== -1) {
        personalEvents.value[index] = updatedEvent
      }
      return updatedEvent
    })
  }

  // 邀請相關
  const fetchInvitations = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const invitationsList = await api.get<Invitation[]>('/teacher/me/invitations')
        invitations.value = invitationsList || []
      } catch (error) {
        console.error('Failed to fetch invitations:', error)
        throw error
      }
    })
  }

  const respondToInvitation = async (invitationId: number, action: 'ACCEPT' | 'REJECT') => {
    return withLoading(isRespondingInvitation, async () => {
      try {
        const api = useApi()
        await api.post('/teacher/me/invitations/respond', {
          invitation_id: invitationId,
          response: action,
        })
        await fetchInvitations()
        await fetchPendingCount()
      } catch (error) {
        console.error('Failed to respond to invitation:', error)
        throw error
      }
    })
  }

  const fetchPendingCount = async () => {
    try {
      const api = useApi()
      const response = await api.get<{ count: number }>('/teacher/me/invitations/pending-count')
      pendingInvitationsCount.value = response?.count || 0
    } catch (error) {
      console.error('Failed to fetch pending invitations count:', error)
    }
  }

  // 調動課表
  const moveScheduleItem = async (data: {
    item_id: number
    item_type: 'SCHEDULE_RULE' | 'PERSONAL_EVENT' | 'CENTER_SESSION'
    center_id: number
    new_date: string
    new_start_time: string
    new_end_time: string
    update_mode?: 'SINGLE' | 'FUTURE' | 'ALL'
  }) => {
    const api = useApi()

    if (data.item_type === 'PERSONAL_EVENT') {
      await api.patch(`/teacher/me/personal-events/${data.item_id}`, {
        start_at: `${data.new_date}T${data.new_start_time}:00`,
        end_at: `${data.new_date}T${data.new_end_time}:00`,
        update_mode: data.update_mode || 'SINGLE',
      })
    } else {
      await api.post(`/teacher/scheduling/edit-recurring`, {
        rule_id: data.item_id,
        edit_date: data.new_date,
        mode: data.update_mode || 'SINGLE',
        new_start_time: data.new_start_time,
        new_end_time: data.new_end_time,
      })
    }
  }

  // 訂閱相關
  const createSubscription = async () => {
    return withLoading(isCreatingSubscription, async () => {
      try {
        const api = useApi()
        const response = await api.post<{ subscription_url: string }>(
          '/teacher/me/schedule/subscription',
          {}
        )
        subscriptionUrl.value = response?.subscription_url || null
        return subscriptionUrl.value
      } catch (error) {
        console.error('Failed to create subscription:', error)
        throw error
      }
    })
  }

  const deleteSubscription = async () => {
    return withLoading(isDeletingSubscription, async () => {
      try {
        const api = useApi()
        await api.delete('/teacher/me/schedule/subscription')
        subscriptionUrl.value = null
      } catch (error) {
        console.error('Failed to delete subscription:', error)
        throw error
      }
    })
  }

  // 下載課表圖片
  const downloadImage = async (startDate?: string, endDate?: string) => {
    return withLoading(isDownloadingImage, async () => {
      try {
        const api = useApi()
        // 支援傳遞日期範圍參數
        const dateParams = startDate && endDate
          ? `?start_date=${startDate}&end_date=${endDate}`
          : ''
        const response = await api.raw<Blob>(`/teacher/me/schedule/image${dateParams}`)

        // 建立下載連結
        const url = URL.createObjectURL(response)
        const link = document.createElement('a')
        link.href = url

        // 產生檔案名稱
        const today = new Date().toISOString().split('T')[0]
        link.download = `課表-${today}.png`

        // 觸發下載
        link.click()

        // 清理
        URL.revokeObjectURL(url)
      } catch (error) {
        console.error('Failed to download image:', error)
        throw error
      }
    })
  }

  return {
    // 資料狀態
    centers,
    currentCenter,
    schedule,
    exceptions,
    personalEvents,
    sessionNote,
    invitations,
    pendingInvitationsCount,
    subscriptionUrl,
    weekStart,
    weekEnd,
    weekLabel,

    // Loading 狀態
    isLoading,
    isFetching,
    isCreating,
    isUpdating,
    isDeleting,
    isCreatingEvent,
    isUpdatingEvent,
    isDeletingEvent,
    isCreatingException,
    isRevokingException,
    isSavingNote,
    isRespondingInvitation,
    isCreatingSubscription,
    isDeletingSubscription,
    isDownloadingImage,

    // 方法
    fetchCenters,
    changeWeek,
    fetchSchedule,
    fetchExceptions,
    createException,
    revokeException,
    fetchSessionNote,
    saveSessionNote,
    formatDate,
    fetchPersonalEvents,
    createPersonalEvent,
    updatePersonalEvent,
    deletePersonalEvent,
    fetchInvitations,
    respondToInvitation,
    fetchPendingCount,
    moveScheduleItem,
    createSubscription,
    deleteSubscription,
    downloadImage,
  }
})
