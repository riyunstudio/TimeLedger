export interface User {
  id: number
  email?: string
  name: string
  avatar_url?: string
  bio?: string
}

export interface Teacher extends User {
  line_user_id: string
  is_open_to_hiring: boolean
  city?: string
  district?: string
  public_contact_info?: PublicContactInfo
  skills?: TeacherSkill[]
  certificates?: TeacherCertificate[]
  personal_hashtags?: PersonalHashtag[]
}

export interface PublicContactInfo {
  instagram?: string
  youtube?: string
  website?: string
  other?: string
}

export interface PersonalHashtag {
  id: number
  hashtag_id: number
  name: string
}

export interface TeacherSkill {
  id: number
  teacher_id: number
  skill_name: string
  hashtags?: TeacherSkillHashtag[]
}

export interface TeacherCertificate {
  id: number
  teacher_id: number
  certificate_name: string
  issued_by?: string
  issued_date?: string
  file_url?: string
}

export interface TeacherPersonalHashtag {
  id: number
  teacher_id: number
  hashtag_id: number
  hashtag?: Hashtag
}

export interface TeacherSkillHashtag {
  id: number
  teacher_skill_id: number
  hashtag_id: number
  hashtag?: Hashtag
}

export interface Hashtag {
  id: number
  name: string
  usage_count: number
}

export interface AdminUser extends User {
  role: 'ADMIN' | 'CENTER_ADMIN'
  center_id?: number
}

export interface Center {
  id: number
  name: string
  plan_level: 'STARTER' | 'GROWTH' | 'PRO'
  settings: CenterSettings
  created_at: string
  updated_at: string
}

export interface CenterSettings {
  allow_public_register: boolean
  default_language: string
}

export interface Course {
  id: number
  center_id: number
  name: string
  teacher_buffer_min: number
  room_buffer_min: number
  created_at: string
  updated_at: string
}

export interface Offering {
  id: number
  center_id: number
  course_id: number
  default_room_id?: number
  default_teacher_id?: number
  allow_buffer_override: boolean
  created_at: string
  updated_at: string
}

export interface Room {
  id: number
  center_id: number
  name: string
  capacity: number
  created_at: string
  updated_at: string
}

export interface ScheduleRule {
  id: number
  center_id: number
  offering_id: number
  teacher_id?: number
  room_id: number
  weekday: number
  start_time: string
  end_time: string
  effective_range: DateRange
  created_at: string
  updated_at: string
  exceptions?: ScheduleException[]
}

export interface ScheduleException {
  id: number
  center_id: number
  rule_id: number
  teacher_id: number
  original_date: string
  type: 'CANCEL' | 'RESCHEDULE'
  status: 'PENDING' | 'APPROVED' | 'REJECTED' | 'REVOKED'
  new_start_at?: string
  new_end_at?: string
  new_teacher_id?: number
  reason: string
  created_at: string
  updated_at: string
}

export interface PersonalEvent {
  id: number | string  // string when expanded from recurrence (format: "originalId_date")
  originalId?: number  // Original ID for API calls
  teacher_id: number
  title: string
  start_at: string
  end_at: string
  recurrence_rule?: RecurrenceRule
  color: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface DateRange {
  start_date: string
  end_date: string
}

export interface RecurrenceRule {
  frequency: 'NONE' | 'DAILY' | 'WEEKLY' | 'BIWEEKLY' | 'MONTHLY'
  interval: number
  end_date?: string
}

export interface SessionNote {
  id: number
  center_id: number
  rule_id: number
  session_date: string
  content: string
  prep_note: string
  created_at: string
  updated_at: string
}

export interface Notification {
  id: number
  user_id: number
  user_type: 'ADMIN' | 'TEACHER'
  center_id?: number
  title: string
  message: string
  type: 'SCHEDULE' | 'EXCEPTION' | 'REVIEW' | 'GENERAL' | 'APPROVAL' | 'CENTER_INVITE'
  is_read: boolean
  read_at?: string
  created_at: string
}

export interface CenterMembership {
  id: number
  center_id: number
  center_name?: string
  teacher_id: number
  status: 'ACTIVE' | 'INACTIVE' | 'INVITED'
}

export interface AuthResponse {
  token: string
  refresh_token: string
  user?: User
  teacher?: Teacher
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface ValidationResult {
  valid: boolean
  conflicts: ValidationConflict[]
}

export interface ValidationConflict {
  type: 'OVERLAP' | 'TEACHER_OVERLAP' | 'ROOM_OVERLAP'
  message: string
  details?: string
}

export interface MatchScore {
  teacher_id: number
  teacher_name: string
  match_score: number
  skill_match: number
  rating: number
  notes?: string
}

export interface ScheduleCell {
  date: string
  time: string
  items: (ScheduleRule | PersonalEvent)[]
  has_conflict: boolean
}

export interface WeekSchedule {
  week_start: string
  week_end: string
  days: DaySchedule[]
}

export interface DaySchedule {
  date: string
  day_of_week: number
  items: ScheduleItem[]
}

export interface ScheduleItem {
  type: 'SCHEDULE_RULE' | 'PERSONAL_EVENT' | 'CENTER_SESSION'
  id: number | string
  title: string
  start_time: string
  end_time: string
  color?: string
  status?: string
  center_name?: string
  data?: ScheduleRule | PersonalEvent
  date?: string
  room_id?: number
  teacher_id?: number
  center_id?: number
}
