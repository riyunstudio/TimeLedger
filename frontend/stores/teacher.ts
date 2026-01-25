import { defineStore } from 'pinia'
import type { CenterMembership, Center, WeekSchedule, ScheduleException, SessionNote, PersonalEvent, RecurrenceRule, TeacherSkill, TeacherCertificate, Teacher, Notification } from '~/types'

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
  data?: any
}

export const useTeacherStore = defineStore('teacher', () => {
  const centers = ref<CenterMembership[]>([])
  const currentCenter = ref<Center | null>(null)
  const schedule = ref<WeekSchedule | null>(null)
  const exceptions = ref<ScheduleException[]>([])

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
    const newStart = new Date(weekStart.value)
    newStart.setDate(newStart.getDate() + (delta * 7))
    weekStart.value = getWeekStart(newStart)
  }

  const fetchCenters = async () => {
    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: CenterMembership[] }>('/teacher/me/centers')
      centers.value = response.datas || []
      if (centers.value.length > 0 && !currentCenter.value && centers.value[0].center_id) {
        currentCenter.value = { id: centers.value[0].center_id, name: centers.value[0].center_name || '' } as any
      }
    } catch (error) {
      console.error('Failed to fetch centers:', error)
    }
  }

  const fetchSchedule = async () => {
    if (!weekStart.value || !weekEnd.value) return

    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: TeacherScheduleItem[] }>(
        `/teacher/me/schedule?from=${formatDate(weekStart.value)}&to=${formatDate(weekEnd.value)}`
      )
      schedule.value = transformToWeekSchedule(response.datas || [])
    } catch (error) {
      console.error('Failed to fetch schedule:', error)
    }
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
      const dateStr = d.toISOString().split('T')[0]
      const dayOfWeek = d.getDay()
      days.push({
        date: dateStr,
        day_of_week: dayOfWeek,
        items: daysMap.get(dateStr) || [],
      })
    }

    return {
      week_start: start.toISOString().split('T')[0],
      week_end: end.toISOString().split('T')[0],
      days
    } as WeekSchedule
  }

  const fetchExceptions = async (status?: string) => {
    try {
      const api = useApi()
      const endpoint = status
        ? `/teacher/exceptions?status=${status}`
        : '/teacher/exceptions'
      const response = await api.get<{ code: number; message: string; datas: ScheduleException[] }>(endpoint)
      exceptions.value = response.datas || []
    } catch (error) {
      console.error('Failed to fetch exceptions:', error)
    }
  }

  const createException = async (data: {
    center_id: number
    rule_id: number
    original_date: string
    type: 'CANCEL' | 'RESCHEDULE'
    new_start_at?: string
    new_end_at?: string
    new_teacher_id?: number
    reason: string
  }) => {
    const api = useApi()
    const response = await api.post<{ code: number; message: string; datas: ScheduleException }>('/teacher/exceptions', data)
    exceptions.value.unshift(response.datas)
    return response.datas
  }

  const revokeException = async (exceptionId: number) => {
    const api = useApi()
    await api.post(`/teacher/exceptions/${exceptionId}/revoke`, {})
    await fetchExceptions()
  }

  const sessionNote = ref<SessionNote | null>(null)

  const fetchSessionNote = async (ruleId: number, sessionDate: string) => {
    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: { note: SessionNote; is_new: boolean } }>(
        `/teacher/sessions/note?rule_id=${ruleId}&session_date=${sessionDate}`
      )
      sessionNote.value = response.datas?.note || null
      return sessionNote.value
    } catch (error) {
      console.error('Failed to fetch session note:', error)
      return null
    }
  }

  const saveSessionNote = async (ruleId: number, sessionDate: string, content: string, prepNote: string) => {
    try {
      const api = useApi()
      const response = await api.put<{ code: number; message: string; datas: SessionNote }>('/teacher/sessions/note', {
        rule_id: ruleId,
        session_date: sessionDate,
        content,
        prep_note: prepNote,
      })
      sessionNote.value = response.datas || null
      return sessionNote.value
    } catch (error) {
      console.error('Failed to save session note:', error)
      throw error
    }
  }

  const formatDate = (date: Date): string => {
    return date.toISOString().split('T')[0]
  }

  const personalEvents = ref<PersonalEvent[]>([])

  const fetchPersonalEvents = async () => {
    if (!weekStart.value || !weekEnd.value) return

    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: PersonalEvent[] }>(
        `/teacher/me/personal-events?from=${formatDate(weekStart.value)}&to=${formatDate(weekEnd.value)}`
      )
      const events = response.datas || []

      // Expand recurring events
      const currentWeekStart = new Date(weekStart.value!)
      const currentWeekEnd = new Date(weekEnd.value!)
      currentWeekEnd.setHours(23, 59, 59, 999)

      const expandedEvents: PersonalEvent[] = []

      events.forEach(event => {
        if (!event.recurrence_rule) {
          expandedEvents.push(event)
          return
        }

        const { frequency, interval = 1 } = event.recurrence_rule
        const startDate = new Date(event.start_at)
        const endDate = new Date(event.end_at)
        const duration = endDate.getTime() - startDate.getTime()

        let currentDate = new Date(startDate)

        while (currentDate <= currentWeekEnd) {
          if (currentDate >= currentWeekStart) {
            const instanceEnd = new Date(currentDate.getTime() + duration)
            expandedEvents.push({
              ...event,
              id: `${event.id}_${currentDate.toISOString().split('T')[0]}`,
              originalId: event.id, // 保留原始 ID 用於 API 調用
              start_at: currentDate.toISOString(),
              end_at: instanceEnd.toISOString(),
            })
          }

          // Advance based on recurrence frequency
          switch (frequency) {
            case 'DAILY':
              currentDate.setDate(currentDate.getDate() + interval)
              break
            case 'WEEKLY':
              currentDate.setDate(currentDate.getDate() + (7 * interval))
              break
            case 'BIWEEKLY':
              currentDate.setDate(currentDate.getDate() + (14 * interval))
              break
            case 'MONTHLY':
              currentDate.setMonth(currentDate.getMonth() + interval)
              break
            default:
              currentDate = new Date(currentWeekEnd.getTime() + 1) // Stop loop
          }
        }
      })

      personalEvents.value = expandedEvents.filter(event => {
        const eventStart = new Date(event.start_at)
        return eventStart >= currentWeekStart && eventStart <= currentWeekEnd
      })
    } catch (error) {
      console.error('Failed to fetch personal events:', error)
    }
  }

  const createPersonalEvent = async (data: {
    title: string
    start_at: string
    end_at: string
    is_all_day?: boolean
    color_hex?: string
    recurrence_rule?: RecurrenceRule
  }) => {
    const api = useApi()
    const response = await api.post<{
      code: number
      message: string
      datas: PersonalEvent
    }>('/teacher/me/personal-events', data)
    const newEvent = response.datas
    // Ensure recurrence_rule is included if sent
    if (data.recurrence_rule && !newEvent.recurrence_rule) {
      newEvent.recurrence_rule = data.recurrence_rule
    }
    personalEvents.value.push(newEvent)
    return newEvent
  }

  const deletePersonalEvent = async (eventId: number | string) => {
    const api = useApi()
    // 使用 originalId 進行 API 調用
    const originalId = typeof eventId === 'string' && eventId.includes('_')
      ? parseInt(eventId.split('_')[0])
      : eventId
    await api.delete(`/teacher/me/personal-events/${originalId}`)
    personalEvents.value = personalEvents.value.filter(e => e.id !== eventId)
  }

  const updatePersonalEvent = async (eventId: number | string, data: {
    title: string
    start_at: string
    end_at: string
    color_hex?: string
    recurrence_rule?: RecurrenceRule
  }) => {
    const api = useApi()
    // 使用 originalId 進行 API 調用
    const originalId = typeof eventId === 'string' && eventId.includes('_')
      ? parseInt(eventId.split('_')[0])
      : eventId
    const updateData = { ...data, update_mode: 'SINGLE' }
    const response = await api.patch<{ code: number; message: string; datas: PersonalEvent }>(`/teacher/me/personal-events/${originalId}`, updateData)
    const index = personalEvents.value.findIndex(e => e.id === eventId)
    if (index !== -1) {
      personalEvents.value[index] = response.datas
    }
    return response.datas
  }

  const skills = ref<TeacherSkill[]>([])

  const fetchSkills = async () => {
    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: TeacherSkill[] }>('/teacher/me/skills')
      skills.value = response.datas || []
    } catch (error) {
      console.error('Failed to fetch skills:', error)
    }
  }

  const createSkill = async (data: { category: string; skill_name: string; level: string; hashtag_ids?: number[] }) => {
    const api = useApi()
    const response = await api.post<{ code: number; message: string; datas: TeacherSkill }>('/teacher/me/skills', data)
    skills.value.push(response.datas)
    return response.datas
  }

  const deleteSkill = async (skillId: number) => {
    const api = useApi()
    await api.delete(`/teacher/me/skills/${skillId}`)
    skills.value = skills.value.filter(s => s.id !== skillId)
  }

  const certificates = ref<TeacherCertificate[]>([])

  const fetchCertificates = async () => {
    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: TeacherCertificate[] }>('/teacher/me/certificates')
      certificates.value = response.datas || []
    } catch (error) {
      console.error('Failed to fetch certificates:', error)
    }
  }

  const createCertificate = async (data: {
    name: string
    file_url?: string
    issued_at?: string
  }) => {
    const api = useApi()
    const response = await api.post<{ code: number; message: string; datas: TeacherCertificate }>('/teacher/me/certificates', data)
    certificates.value.push(response.datas)
    return response.datas
  }

  const deleteCertificate = async (certId: number) => {
    const api = useApi()
    await api.delete(`/teacher/me/certificates/${certId}`)
    certificates.value = certificates.value.filter(c => c.id !== certId)
  }

  const profile = ref<Teacher | null>(null)

  const fetchProfile = async () => {
    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: Teacher }>('/teacher/me/profile')
      profile.value = response.datas || null
      return response.datas
    } catch (error) {
      console.error('Failed to fetch profile:', error)
      return null
    }
  }

  const updateProfile = async (data: Partial<Teacher>) => {
    const api = useApi()
    const response = await api.put<{ code: number; message: string; datas: Teacher }>('/teacher/me/profile', data)
    profile.value = response.datas || null
    return response.datas
  }

  const notifications = ref<Notification[]>([])

  const fetchNotifications = async () => {
    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: Notification[] }>('/notifications')
      notifications.value = response.datas || []
    } catch (error) {
      console.error('Failed to fetch notifications:', error)
    }
  }

  const markNotificationRead = async (notificationId: number) => {
    try {
      const api = useApi()
      await api.post(`/notifications/${notificationId}/read`, {})
      const notification = notifications.value.find(n => n.id === notificationId)
      if (notification) {
        notification.is_read = true
        notification.read_at = new Date().toISOString()
      }
    } catch (error) {
      console.error('Failed to mark notification as read:', error)
    }
  }

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
      await api.put(`/teacher/me/personal-events/${data.item_id}`, {
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

  return {
    centers,
    currentCenter,
    schedule,
    exceptions,
    sessionNote,
    weekStart,
    weekEnd,
    weekLabel,
    personalEvents,
    skills,
    certificates,
    profile,
    notifications,
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
    fetchSkills,
    createSkill,
    deleteSkill,
    fetchCertificates,
    createCertificate,
    deleteCertificate,
    fetchProfile,
    updateProfile,
    fetchNotifications,
    markNotificationRead,
    moveScheduleItem,
  }
})
