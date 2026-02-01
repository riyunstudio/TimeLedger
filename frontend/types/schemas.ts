/**
 * Zod Schema 定義
 *
 * 使用 Zod 進行 runtime 資料驗證，確保 API 回應符合預期結構
 *
 * 使用方式：
 * ```typescript
 * import { ScheduleRuleSchema, TeacherProfileSchema } from '~/types/schemas'
 *
 * // 驗證 API 回應資料
 * const result = ScheduleRuleSchema.safeParse(apiData)
 * if (!result.success) {
 *   console.warn('Schema 驗證失敗:', result.error.issues)
 * }
 * ```
 */

import { z } from 'zod'

// ==================== 基礎類型 Schema ====================

/** ID 驗證 */
export const IDSchema = z.number().int().positive()

/** 時間戳記驗證 (允許不同格式的字串) */
export const TimestampSchema = z.string()

/** 日期字串驗證 (YYYY-MM-DD 格式，允許帶時間或不完整格式) */
export const DateStringSchema = z.string().transform((val) => val.split(' ')[0].split('T')[0])

/** 分頁排序方向驗證 */
export const SortOrderSchema = z.enum(['ASC', 'DESC'])

// ==================== 分頁相關 Schema ====================

/**
 * 分頁查詢參數 Schema
 */
export const PaginationParamsSchema = z.object({
  /** 頁碼 (從 1 開始) */
  page: z.number().int().positive().default(1),
  /** 每頁筆數 (預設 20，最大 100) */
  limit: z.number().int().positive().max(100).default(20),
  /** 排序欄位 */
  sort_by: z.string().optional(),
  /** 排序方向 */
  sort_order: SortOrderSchema.default('ASC'),
})

/**
 * 分頁結果 Schema
 */
export const PaginationResultSchema = z.object({
  /** 當前頁碼 */
  page: z.number().int().positive(),
  /** 每頁筆數 */
  limit: z.number().int().positive(),
  /** 總筆數 */
  total: z.number().int().nonnegative(),
  /** 總頁數 */
  total_pages: z.number().int().nonnegative(),
  /** 是否有下一頁 */
  has_next: z.boolean(),
  /** 是否有上一頁 */
  has_prev: z.boolean(),
})

// ==================== 日期範圍 Schema ====================

/**
 * 日期範圍 Schema
 */
export const DateRangeSchema = z.object({
  /** 開始日期 (YYYY-MM-DD) */
  start_date: DateStringSchema,
  /** 結束日期 (YYYY-MM-DD) */
  end_date: DateStringSchema,
})

/**
 * 日期時間範圍 Schema
 */
export const DateTimeRangeSchema = z.object({
  /** 開始時間 (ISO 8601) */
  start_at: TimestampSchema,
  /** 結束時間 (ISO 8601) */
  end_at: TimestampSchema,
})

// ==================== 循環規則 Schema ====================

/**
 * 循環頻率類型 Schema
 */
export const RecurrenceFrequencySchema = z.enum([
  'NONE',
  'DAILY',
  'WEEKLY',
  'BIWEEKLY',
  'MONTHLY',
])

/**
 * 循環規則 Schema
 *
 * 對應後端 RecurrenceRule 模型
 */
export const RecurrenceRuleSchema = z.object({
  /** 循環類型 */
  type: RecurrenceFrequencySchema,
  /** 循環間隔 (如每 2 週為 interval: 2) */
  interval: z.number().int().positive().default(1),
  /** 結束日期 (可選，不設限則持續循環) */
  end_date: DateStringSchema.optional(),
  /** 循環次數上限 (可選) */
  count: z.number().int().positive().optional(),
  /** 結束日期 ISO 字串 (可選) */
  until: z.string().optional(),
  /** 星期幾循環 (用於 WEEKLY, BIWEEKLY) */
  weekdays: z.array(z.number().min(0).max(6)).optional(),
})

// ==================== 驗證結果 Schema ====================

/**
 * 驗證衝突類型 Schema
 */
export const ValidationConflictTypeSchema = z.enum([
  'OVERLAP',
  'TEACHER_OVERLAP',
  'ROOM_OVERLAP',
  'BUFFER',
])

/**
 * 驗證衝突詳細資訊 Schema
 */
export const ValidationConflictSchema = z.object({
  /** 衝突類型 */
  type: ValidationConflictTypeSchema,
  /** 衝突訊息 */
  message: z.string(),
  /** 額外詳細資訊 */
  details: z.string().optional(),
})

/**
 * 驗證結果 Schema
 */
export const ValidationResultSchema = z.object({
  /** 是否有效 */
  valid: z.boolean(),
  /** 衝突列表 */
  conflicts: z.array(ValidationConflictSchema),
})

// ==================== 課表規則相關 Schema ====================

/**
 * 課表規則 Schema
 *
 * 對應後端 ScheduleRule 模型 / resources.ScheduleRuleResource
 */
export const ScheduleRuleSchema = z.object({
  /** 規則 ID */
  id: IDSchema,
  /** 所屬中心 ID */
  center_id: IDSchema,
  /** 關聯開課 ID */
  offering_id: IDSchema,
  /** 開課名稱 (可選) */
  offering_name: z.string().optional(),
  /** 指定教師 ID (可選) */
  teacher_id: IDSchema.nullable().optional(),
  /** 教師名稱 (可選) */
  teacher_name: z.string().nullable().optional(),
  /** 教室 ID */
  room_id: IDSchema,
  /** 教室名稱 (可選) */
  room_name: z.string().nullable().optional(),
  /** 星期幾 (0-6, 0 為週日) */
  weekday: z.number().int().min(0).max(6),
  /** 開始時間 (HH:mm 或 HH:mm:ss) */
  start_time: z.string(),
  /** 結束時間 (HH:mm 或 HH:mm:ss) */
  end_time: z.string(),
  /** 課程時長 (分鐘) */
  duration: z.number().int().positive().optional(),
  /** 有效範圍 */
  effective_range: DateRangeSchema.optional(),
  /** 循環規則 (可選) */
  recurrence_rule: RecurrenceRuleSchema.optional(),
  /** 是否為跨日課程 */
  is_cross_day: z.boolean().optional(),
  /** 是否已鎖定 */
  is_locked: z.boolean().optional(),
  /** 鎖定時間 (可選) */
  lock_at: TimestampSchema.optional(),
  /** 鎖定原因 (可選) */
  lock_reason: z.string().nullable().optional(),
  /** 例外列表 (可選) */
  exceptions: z.array(z.any()).optional(),
  /** 建立時間 */
  created_at: TimestampSchema,
  /** 更新時間 */
  updated_at: TimestampSchema,
})

/**
 * 課表規則列表 Schema (含分頁)
 */
export const ScheduleRuleListSchema = z.object({
  /** 資料陣列 */
  data: z.array(ScheduleRuleSchema),
  /** 分頁資訊 */
  pagination: PaginationResultSchema,
})

/**
 * 新增課表規則請求 Schema
 */
export const CreateScheduleRuleSchema = z.object({
  /** 開課 ID */
  offering_id: IDSchema,
  /** 指定教師 ID (可選) */
  teacher_id: IDSchema.optional(),
  /** 教室 ID */
  room_id: IDSchema,
  /** 星期幾 (0-6) */
  weekday: z.number().int().min(0).max(6),
  /** 開始時間 */
  start_time: z.string(),
  /** 結束時間 */
  end_time: z.string(),
  /** 有效開始日期 */
  effective_start_date: DateStringSchema,
  /** 有效結束日期 */
  effective_end_date: DateStringSchema,
  /** 循環規則 (可選) */
  recurrence_rule: RecurrenceRuleSchema.optional(),
})

/**
 * 更新課表規則請求 Schema
 */
export const UpdateScheduleRuleSchema = z.object({
  /** 指定教師 ID (可選) */
  teacher_id: IDSchema.nullable().optional(),
  /** 教室 ID */
  room_id: IDSchema.optional(),
  /** 開始時間 */
  start_time: z.string().optional(),
  /** 結束時間 */
  end_time: z.string().optional(),
  /** 有效開始日期 */
  effective_start_date: DateStringSchema.optional(),
  /** 有效結束日期 */
  effective_end_date: DateStringSchema.optional(),
})

// ==================== 教師相關 Schema ====================

/**
 * 技能分類 Schema
 */
export const SkillCategorySchema = z.enum([
  'MUSIC',
  'ART',
  'DANCE',
  'LANGUAGE',
  'SPORTS',
  'OTHER',
])

/**
 * 標籤 Schema
 */
export const HashtagSchema = z.object({
  /** 標籤 ID */
  id: IDSchema,
  /** 標籤名稱 (含 # 符號) */
  name: z.string().startsWith('#'),
  /** 使用次數 */
  usage_count: z.number().int().nonnegative(),
})

/**
 * 教師技能 Schema
 */
export const TeacherSkillSchema = z.object({
  /** 技能 ID */
  id: IDSchema,
  /** 所屬教師 ID */
  teacher_id: IDSchema,
  /** 技能分類 */
  category: SkillCategorySchema.or(z.string()),
  /** 技能名稱 */
  skill_name: z.string(),
  /** 技能標籤 */
  hashtags: z.array(HashtagSchema).optional(),
})

/**
 * 教師證照 Schema
 *
 * 對應後端 TeacherCertificate 模型 / AdminCertificateResponse
 */
export const TeacherCertificateSchema = z.object({
  /** 證照 ID */
  id: IDSchema,
  /** 證照名稱 */
  name: z.string(),
  /** 證照圖片 URL */
  file_url: z.string().url().optional(),
  /** 發證日期 */
  issued_at: TimestampSchema,
  /** 是否已驗證 */
  is_verified: z.boolean().optional(),
  /** 建立時間 */
  created_at: TimestampSchema,
  /** 更新時間 */
  updated_at: TimestampSchema,
})

/**
 * 個人標籤 Schema
 *
 * 對應後端 PersonalHashtag / TeacherPersonalHashtagResource
 */
export const PersonalHashtagSchema = z.object({
  /** 標籤 ID */
  hashtag_id: IDSchema,
  /** 標籤名稱 */
  name: z.string(),
  /** 排序順序 */
  sort_order: z.number().int().optional(),
})

/**
 * 教師基礎資料 Schema
 */
export const TeacherSchema = z.object({
  /** 教師 ID */
  id: IDSchema,
  /** 名稱 */
  name: z.string().min(1),
  /** 電子郵件 */
  email: z.string().email().optional(),
  /** 電話 */
  phone: z.string().optional(),
  /** LINE 用戶 ID (帳號綁定，不可解除) */
  line_user_id: z.string(),
  /** 縣市 */
  city: z.string().optional(),
  /** 區域 */
  district: z.string().optional(),
  /** 技能列表 */
  skills: z.array(TeacherSkillSchema).optional(),
  /** 證照列表 */
  certificates: z.array(TeacherCertificateSchema).optional(),
  /** 個人標籤 */
  personal_hashtags: z.array(PersonalHashtagSchema).optional(),
  /** 是否開放應徵 (人才庫) */
  is_open_to_hiring: z.boolean(),
  /** 是否已激活 */
  is_active: z.boolean(),
  /** 邀請時間 */
  invited_at: TimestampSchema.optional(),
  /** 激活時間 */
  activated_at: TimestampSchema.optional(),
  /** 建立時間 */
  created_at: TimestampSchema,
  /** 更新時間 */
  updated_at: TimestampSchema,
})

/**
 * 教師中心關聯 Schema
 */
export const TeacherCenterMembershipSchema = z.object({
  /** 關聯 ID */
  id: IDSchema,
  /** 中心 ID */
  center_id: IDSchema,
  /** 中心名稱 */
  center_name: z.string().optional(),
  /** 教師 ID */
  teacher_id: IDSchema,
  /** 狀態 */
  status: z.enum(['ACTIVE', 'INACTIVE', 'INVITED']),
  /** 加入時間 */
  joined_at: TimestampSchema.optional(),
  /** 邀請時間 */
  invited_at: TimestampSchema.optional(),
})

/**
 * 教師完整設定檔 Schema
 *
 * 對應後端 TeacherProfile 模型
 */
export const TeacherProfileSchema = z.object({
  /** 教師基本資訊 */
  profile: TeacherSchema,
  /** 技能列表 */
  skills: z.array(TeacherSkillSchema),
  /** 證照列表 */
  certificates: z.array(TeacherCertificateSchema),
  /** 個人標籤 */
  hashtags: z.array(PersonalHashtagSchema),
  /** 加入的中心列表 */
  centers: z.array(TeacherCenterMembershipSchema),
})

/**
 * 教師設定檔 API 回應 Schema
 */
export const TeacherProfileResponseSchema = z.object({
  /** 錯誤碼 */
  code: z.string(),
  /** 訊息描述 */
  message: z.string(),
  /** 資料 */
  data: TeacherProfileSchema.optional(),
  /** 多筆資料 (部分 API 使用) */
  datas: TeacherProfileSchema.optional(),
})

// ==================== 教師課表項目 Schema ====================

/**
 * 教師課表項目 Schema
 *
 * 對應 useScheduleStore 中的 TeacherScheduleItem 介面
 * 用於 fetchSchedule API 回應驗證
 */
export const TeacherScheduleItemSchema = z.object({
  /** 項目 ID */
  id: z.string(),
  /** 項目類型 (SCHEDULE_RULE, PERSONAL_EVENT, CENTER_SESSION) */
  type: z.string(),
  /** 標題 */
  title: z.string(),
  /** 日期 (YYYY-MM-DD) */
  date: DateStringSchema,
  /** 開始時間 */
  start_time: z.string(),
  /** 結束時間 */
  end_time: z.string(),
  /** 教室 ID */
  room_id: IDSchema.optional(),
  /** 教師 ID (可選) */
  teacher_id: IDSchema.optional(),
  /** 中心 ID */
  center_id: IDSchema.optional(),
  /** 中心名稱 (可選) */
  center_name: z.string().optional(),
  /** 狀態 */
  status: z.string().optional(),
  /** 規則 ID (可選) */
  rule_id: IDSchema.optional(),
  /** 原始資料 (可選) */
  data: z.any().optional(),
  /** 是否為跨日課程的一部分 (可選) */
  is_cross_day_part: z.boolean().optional(),
})

/**
 * 教師課表 API 回應 Schema
 */
export const TeacherScheduleResponseSchema = z.object({
  /** 錯誤碼 */
  code: z.number(),
  /** 訊息描述 */
  message: z.string(),
  /** 課表資料 */
  datas: z.array(TeacherScheduleItemSchema).optional(),
})

/**
 * 教師課表純資料 Schema (用於 useApi 驗證內層 datas)
 */
export const TeacherScheduleDataSchema = z.array(TeacherScheduleItemSchema)

// ==================== 例外申請 Schema ====================

/**
 * 例外類型 Schema
 */
export const ExceptionTypeSchema = z.enum([
  'CANCEL',
  'RESCHEDULE',
  'REPLACE_TEACHER',
])

/**
 * 審核狀態 Schema
 */
export const ReviewStatusSchema = z.enum([
  'PENDING',
  'APPROVED',
  'REJECTED',
  'REVOKED',
  'CANCELLED',
])

/**
 * 例外申請 Schema
 *
 * 對應後端 ScheduleException 模型
 * 用於 fetchExceptions API 回應驗證
 */
export const ScheduleExceptionSchema = z.object({
  /** 例外 ID */
  id: IDSchema,
  /** 所屬中心 ID */
  center_id: IDSchema,
  /** 關聯規則 ID */
  rule_id: IDSchema,
  /** 教師 ID */
  teacher_id: IDSchema.nullable().optional(),
  /** 教師名稱 (可選) */
  teacher_name: z.string().nullable().optional(),
  /** 原始日期 (YYYY-MM-DD) */
  original_date: z.string().optional(),
  /** 例外類型 */
  type: z.string().optional(),
  /** 狀態 */
  status: ReviewStatusSchema.or(z.string()),
  /** 新開始時間 (可選) */
  new_start_at: z.string().nullable().optional(),
  /** 新結束時間 (可選) */
  new_end_at: z.string().nullable().optional(),
  /** 替換教師 ID (可選) */
  new_teacher_id: IDSchema.nullable().optional(),
  /** 替換教師名稱 (可選) */
  new_teacher_name: z.string().nullable().optional(),
  /** 原因 */
  reason: z.string().optional(),
  /** 審核者 ID (可選) */
  reviewed_by: IDSchema.nullable().optional(),
  /** 審核者名稱 (可選) */
  reviewed_by_name: z.string().nullable().optional(),
  /** 審核時間 (可選) */
  reviewed_at: z.string().nullable().optional(),
  /** 審核意見 (可選) */
  review_note: z.string().nullable().optional(),
  /** 建立時間 */
  created_at: TimestampSchema.optional(),
  /** 更新時間 */
  updated_at: TimestampSchema.optional(),
})

/**
 * 例外申請列表 API 回應 Schema
 */
export const ScheduleExceptionListResponseSchema = z.object({
  /** 錯誤碼 */
  code: z.number(),
  /** 訊息描述 */
  message: z.string(),
  /** 例外資料列表 */
  datas: z.array(ScheduleExceptionSchema).optional(),
})

/**
 * 例外申請列表純資料 Schema (用於 useApi 驗證內層 datas)
 */
export const ScheduleExceptionListDataSchema = z.array(ScheduleExceptionSchema)

// ==================== 類型匯出 ====================

export type ScheduleRule = z.infer<typeof ScheduleRuleSchema>
export type TeacherProfile = z.infer<typeof TeacherProfileSchema>
export type ValidationResult = z.infer<typeof ValidationResultSchema>
export type RecurrenceRule = z.infer<typeof RecurrenceRuleSchema>
export type TeacherScheduleItem = z.infer<typeof TeacherScheduleItemSchema>
export type ScheduleException = z.infer<typeof ScheduleExceptionSchema>
