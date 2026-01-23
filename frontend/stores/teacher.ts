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
  const isMock = ref(false)

  const checkMockMode = () => {
    if (typeof window === 'undefined') return false
    return localStorage.getItem('timeledger_mock_mode') === 'true'
  }

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

  const loadMockCenters = () => {
    centers.value = [
      {
        id: 1,
        center_id: 1,
        teacher_id: 1,
        center: {
          id: 1,
          name: '音樂教室',
          plan_level: 'PRO' as const,
          settings: { allow_public_register: true, default_language: 'zh-TW' },
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        status: 'ACTIVE' as const,
      },
      {
        id: 2,
        center_id: 2,
        teacher_id: 1,
        center: {
          id: 2,
          name: '藝術中心',
          plan_level: 'GROWTH' as const,
          settings: { allow_public_register: true, default_language: 'zh-TW' },
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        status: 'INVITED' as const,
      },
    ]
    currentCenter.value = centers.value[0].center || null
    isMock.value = true
    if (typeof window !== 'undefined') {
      localStorage.setItem('timeledger_mock_mode', 'true')
    }
  }

  const loadMockSchedule = () => {
    const today = new Date()
    const weekStartDate = getWeekStart(today)
    
    const mockDays = [
      {
        date: getOffsetDate(weekStartDate, 0),
        items: [
          { id: 1, type: 'SCHEDULE_RULE' as const, title: '鋼琴基礎', start_time: '10:00', end_time: '11:00', status: 'APPROVED' as const, color: '#10B981', center_name: '音樂教室', data: { id: 1, center_id: 1, offering_id: 1, teacher_id: 1, room_id: 1, weekday: 1, start_time: '10:00', end_time: '11:00', effective_range: { start: '2026-01-01', end: '2026-12-31' }, created_at: '', updated_at: '' } },
          { id: 2, type: 'SCHEDULE_RULE' as const, title: '鋼琴進階', start_time: '14:00', end_time: '15:00', status: 'APPROVED' as const, color: '#6366F1', center_name: '音樂教室', data: { id: 2, center_id: 1, offering_id: 2, teacher_id: 1, room_id: 1, weekday: 1, start_time: '14:00', end_time: '15:00', effective_range: { start: '2026-01-01', end: '2026-12-31' }, created_at: '', updated_at: '' } },
        ],
      },
      {
        date: getOffsetDate(weekStartDate, 1),
        items: [
          { id: 3, type: 'SCHEDULE_RULE' as const, title: '小提琴入門', start_time: '11:00', end_time: '12:00', status: 'APPROVED' as const, color: '#A855F7', center_name: '藝術中心', data: { id: 3, center_id: 2, offering_id: 3, teacher_id: 2, room_id: 2, weekday: 2, start_time: '11:00', end_time: '12:00', effective_range: { start: '2026-01-01', end: '2026-12-31' }, created_at: '', updated_at: '' } },
        ],
      },
      {
        date: getOffsetDate(weekStartDate, 2),
        items: [],
      },
      {
        date: getOffsetDate(weekStartDate, 3),
        items: [
          { id: 4, type: 'PERSONAL_EVENT' as const, title: '休息時間', start_time: '12:00', end_time: '13:00', status: 'APPROVED' as const, color: '#F59E0B', center_name: '', data: { id: 1, teacher_id: 1, title: '休息時間', start_at: '', end_at: '', created_at: '', updated_at: '' } },
        ],
      },
      {
        date: getOffsetDate(weekStartDate, 4),
        items: [
          { id: 5, type: 'SCHEDULE_RULE' as const, title: '樂理課程', start_time: '15:00', end_time: '16:00', status: 'PENDING' as const, color: '#EC4899', center_name: '音樂教室', data: { id: 5, center_id: 1, offering_id: 4, teacher_id: 1, room_id: 1, weekday: 5, start_time: '15:00', end_time: '16:00', effective_range: { start: '2026-01-01', end: '2026-12-31' }, created_at: '', updated_at: '' } },
        ],
      },
      {
        date: getOffsetDate(weekStartDate, 5),
        items: [],
      },
      {
        date: getOffsetDate(weekStartDate, 6),
        items: [
          { id: 6, type: 'SCHEDULE_RULE' as const, title: '鋼琴基礎', start_time: '09:00', end_time: '10:00', status: 'APPROVED' as const, color: '#10B981', center_name: '音樂教室', data: { id: 6, center_id: 1, offering_id: 1, teacher_id: 1, room_id: 1, weekday: 6, start_time: '09:00', end_time: '10:00', effective_range: { start: '2026-01-01', end: '2026-12-31' }, created_at: '', updated_at: '' } },
        ],
      },
    ]
    
    schedule.value = { days: mockDays } as WeekSchedule
  }

  const getOffsetDate = (date: Date, offset: number): string => {
    const d = new Date(date)
    d.setDate(d.getDate() + offset)
    return d.toISOString().split('T')[0]
  }

  const fetchCenters = async () => {
    if (isMock.value || checkMockMode()) return

    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: CenterMembership[] }>('/teacher/me/centers')
      centers.value = response.datas || []
      if (centers.value.length > 0 && !currentCenter.value) {
        currentCenter.value = centers.value[0].center || null
      }
    } catch (error) {
      console.error('Failed to fetch centers:', error)
    }
  }

  const changeWeek = (delta: number) => {
    if (!weekStart.value) return
    const newStart = new Date(weekStart.value)
    newStart.setDate(newStart.getDate() + (delta * 7))
    weekStart.value = getWeekStart(newStart)
  }

  const fetchSchedule = async () => {
    if (isMock.value || checkMockMode()) {
      loadMockSchedule()
      return
    }

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
    if (isMock.value || checkMockMode()) {
      exceptions.value = [
        {
          id: 1,
          center_id: 1,
          rule_id: 1,
          teacher_id: 1,
          original_date: '2026-01-25',
          type: 'CANCEL',
          status: 'PENDING',
          reason: '身體不適',
          created_at: '2026-01-20T10:00:00Z',
          updated_at: '2026-01-20T10:00:00Z',
        },
      ]
      return
    }

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
    if (isMock.value || checkMockMode()) {
      sessionNote.value = {
        id: 1,
        center_id: 0,
        rule_id: ruleId,
        date: sessionDate,
        content: '上次教到第3頁',
        prep_note: '記得帶講義',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      return sessionNote.value
    }

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
    if (isMock.value || checkMockMode()) {
      personalEvents.value = [
        {
          id: 1,
          teacher_id: 1,
          title: '休息時間',
          start_at: '2026-01-23T12:00:00Z',
          end_at: '2026-01-23T13:00:00Z',
          color: '#F59E0B',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
      ]
      return
    }

    if (!weekStart.value || !weekEnd.value) return

    try {
      const api = useApi()
      const response = await api.get<{ code: number; message: string; datas: PersonalEvent[] }>(
        `/teacher/me/personal-events?from=${formatDate(weekStart.value)}&to=${formatDate(weekEnd.value)}`
      )
      personalEvents.value = response.datas || []
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
    const response = await api.post<{ code: number; message: string; datas: PersonalEvent }>('/teacher/me/personal-events', data)
    personalEvents.value.push(response.datas)
    return response.datas
  }

  const deletePersonalEvent = async (eventId: number) => {
    const api = useApi()
    await api.delete(`/teacher/me/personal-events/${eventId}`)
    personalEvents.value = personalEvents.value.filter(e => e.id !== eventId)
  }

  const skills = ref<TeacherSkill[]>([])

  const fetchSkills = async () => {
    if (isMock.value || checkMockMode()) {
      skills.value = [
        { id: 1, teacher_id: 1, skill_name: '鋼琴', level: 'Advanced', hashtags: [] },
        { id: 2, teacher_id: 1, skill_name: '小提琴', level: 'Intermediate', hashtags: [] },
      ]
      return
    }

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
    if (isMock.value || checkMockMode()) {
      certificates.value = [
        { id: 1, teacher_id: 1, certificate_name: '鋼琴檢定證書', issued_by: '音樂協會', issued_date: '2023-06-15' },
      ]
      return
    }

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
    if (isMock.value || checkMockMode()) {
      profile.value = {
        id: 1,
        name: '老師1',
        email: 'teacher1@example.com',
        line_user_id: 'LINE_USER_001',
        is_open_to_hiring: true,
        bio: '資深鋼琴教師',
        skills: [],
        certificates: [],
        personal_hashtags: [],
      }
      return
    }

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
    if (isMock.value || checkMockMode()) {
      await fetchSchedule()
      return
    }

    const api = useApi()
    const endpoint = data.item_type === 'PERSONAL_EVENT'
      ? `/teacher/me/personal-events/${data.item_id}`
      : `/teacher/schedule/${data.item_id}/move`

    const body: any = {
      new_date: data.new_date,
      new_start_time: data.new_start_time,
      new_end_time: data.new_end_time,
    }

    if (data.item_type === 'PERSONAL_EVENT') {
      body.update_mode = data.update_mode || 'SINGLE'
      body.start_at = `${data.new_date}T${data.new_start_time}:00`
      body.end_at = `${data.new_date}T${data.new_end_time}:00`
    }

    if (data.item_type === 'PERSONAL_EVENT') {
      await api.patch(endpoint, body)
    } else {
      await api.post(endpoint, body)
    }
  }

  return {
    centers,
    currentCenter,
    schedule,
    exceptions,
    sessionNote,
    isMock,
    weekStart,
    weekEnd,
    weekLabel,
    personalEvents,
    skills,
    certificates,
    profile,
    notifications,
    loadMockCenters,
    loadMockSchedule,
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
