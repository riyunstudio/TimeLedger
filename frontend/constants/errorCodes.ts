/**
 * 錯誤碼對照表
 *
 * 定義所有後端 API 錯誤碼對應的使用者友善訊息
 * 錯誤碼定義於 global/errInfos/code.go
 */

// ==================== 成功相關 ====================

export const SUCCESS_CODES = {
  SUCCESS: '操作成功',
  '0': '操作成功', // API 成功回應使用 code: "0"
} as const

// ==================== 系統錯誤 (11xxxx) ====================

export const SYSTEM_ERROR_CODES = {
  SYSTEM_ERROR: '系統錯誤，請稍後再試',
  APP_ERROR: '應用程式錯誤',
  CONFIG_ERROR: '配置錯誤',
  THIRD_PARTY_ERROR: '第三方服務錯誤',
  UNKNOWN_ERROR: '發生未知錯誤',
} as const

// ==================== 資料庫錯誤 (12xxxx) ====================

export const DATABASE_ERROR_CODES = {
  SQL_ERROR: '資料庫操作錯誤，請稍後再試',
  NOT_FOUND: '找不到請求的資源',
  CONFLICT: '資料衝突，請檢查後重試',
  DUPLICATE_ENTRY: '資料已存在',
  FOREIGN_KEY_VIOLATION: '關聯資料不存在',
  CONSTRAINT_VIOLATION: '資料驗證失敗',
  TRANSACTION_ERROR: '交易處理失敗',
  DATA_TOO_LONG: '資料長度超過限制',
  INVALID_DATA: '無效的資料格式',
} as const

// ==================== 權限錯誤 (13xxxx) ====================

export const PERMISSION_ERROR_CODES = {
  UNAUTHORIZED: '請先登入再進行操作',
  FORBIDDEN: '您沒有權限執行此操作',
  TOKEN_EXPIRED: '登入已過期，請重新登入',
  TOKEN_INVALID: '登入驗證失敗，請重新登入',
  ACCESS_DENIED: '拒絕存取',
  ROLE_NOT_ALLOWED: '您的角色無法執行此操作',
  CENTER_NOT_FOUND: '找不到指定的中心',
  CENTER_ACCESS_DENIED: '您沒有存取此中心的權限',
} as const

// ==================== 驗證錯誤 (11xxxx) ====================

export const VALIDATION_ERROR_CODES = {
  VALIDATION_ERROR: '輸入資料驗證失敗，請檢查後重試',
  PARAMETER_MISSING: '缺少必要參數',
  PARAMETER_INVALID: '參數格式錯誤',
  PARAMETER_TOO_LONG: '參數長度超過限制',
  DATE_FORMAT_INVALID: '日期格式錯誤',
  TIME_FORMAT_INVALID: '時間格式錯誤',
  EMAIL_FORMAT_INVALID: '電子郵件格式錯誤',
  PHONE_FORMAT_INVALID: '電話格式錯誤',
  ENUM_VALUE_INVALID: '不支援的選項值',
  CUSTOM_VALIDATION_FAILED: '自訂驗證失敗',
} as const

// ==================== 業務邏輯錯誤 (14xxxx) ====================

export const BUSINESS_ERROR_CODES = {
  SCHEDULE_CONFLICT: '排課時間衝突，請選擇其他時段',
  TEACHER_OVERLAP: '老師時段衝突',
  ROOM_OVERLAP: '教室時段衝突',
  BUFFER_VIOLATION: '緩衝時間不足',
  CANNOT_DELETE_IN_USE: '無法刪除，因為已被使用',
  CANNOT_UPDATE_FINALIZED: '無法修改已確定的資料',
  INVALID_STATUS_TRANSITION: '無效的狀態轉換',
  ALREADY_EXISTS: '資料已存在',
  NOT_APPLICABLE: '不適用于此操作',
  RATE_LIMIT_EXCEEDED: '請求過於頻繁，請稍後再試',
  RESOURCE_BUSY: '資源忙碌中，請稍後再試',
} as const

// ==================== 例外審核錯誤 (16xxxx) ====================

export const EXCEPTION_ERROR_CODES = {
  EXCEPTION_NOT_FOUND: '找不到例外申請',
  EXCEPTION_ALREADY_REVIEWED: '此例外申請已處理',
  EXCEPTION_CANNOT_CANCEL: '無法取消，已超過可取消時間',
  EXCEPTION_REVIEW_CONFLICT: '審核時間與其他排課衝突',
  EXCEPTION_REVIEW_BUFFER_VIOLATION: '審核緩衝時間不足',
  EXCEPTION_DEADLINE_EXCEEDED: '已超過異動截止日',
  EXCEPTION_SELF_REVIEW_FORBIDDEN: '不能審核自己提交的申請',
  EXCEPTION_ALREADY_PROCESSED: '例外已被處理',
  EXCEPTION_RESCHEDULE_CONFLICT: '調課時間與現有排程衝突',
  EXCEPTION_CANCEL_DEADLINE_PASSED: '停課截止日已過',
} as const

// ==================== 教師相關錯誤 (14xxxx) ====================

export const TEACHER_ERROR_CODES = {
  TEACHER_NOT_FOUND: '找不到指定的老師',
  TEACHER_ALREADY_INVITED: '此老師已被邀請',
  TEACHER_ALREADY_MEMBER: '此老師已是中心成員',
  TEACHER_NOT_IN_CENTER: '老師不屬於此中心',
  TEACHER_INACTIVE: '老師帳號已停用',
  TEACHER_NOT_ACTIVATED: '老師帳號尚未激活',
  SKILL_NOT_FOUND: '找不到指定的技能',
  CERTIFICATE_NOT_FOUND: '找不到指定的證照',
} as const

// ==================== 課程相關錯誤 (14xxxx) ====================

export const COURSE_ERROR_CODES = {
  COURSE_NOT_FOUND: '找不到指定的課程',
  COURSE_ALREADY_EXISTS: '課程已存在',
  COURSE_CANNOT_DELETE: '無法刪除課程，因為已有排課記錄',
  OFFERING_NOT_FOUND: '找不到指定的開課',
  OFFERING_FULL: '此開課人數已滿',
  OFFERING_NOT_AVAILABLE: '此開課已停止報名',
} as const

// ==================== 教室相關錯誤 (14xxxx) ====================

export const ROOM_ERROR_CODES = {
  ROOM_NOT_FOUND: '找不到指定的教室',
  ROOM_ALREADY_EXISTS: '教室已存在',
  ROOM_CAPACITY_EXCEEDED: '超過教室容納人數',
  ROOM_NOT_AVAILABLE: '教室不可用',
  ROOM_MAINTENANCE: '教室維護中',
} as const

// ==================== 通知相關錯誤 (19xxxx) ====================

export const NOTIFICATION_ERROR_CODES = {
  NOTIFICATION_NOT_FOUND: '找不到通知',
  NOTIFICATION_SEND_FAILED: '通知發送失敗',
  LINE_BIND_FAILED: 'LINE 綁定失敗',
  LINE_UNBIND_FAILED: 'LINE 解除綁定失敗',
  LINE_NOTIFY_FAILED: 'LINE 通知發送失敗',
} as const

// ==================== 整合所有錯誤碼 ====================

/**
 * 數字錯誤碼到字符串錯誤碼的映射
 * 後端返回數字格式錯誤碼（如 160006），需要映射到前端的字符串格式
 *
 * 錯誤碼定義規則（後端 global/errInfos/code.go）：
 * - 第1位：專案前綴 (1 = TimeLedger)
 * - 第2-3位：功能類型 (10=系統, 20=DB, 30=權限, 50=排課, 60=例外審核, ...)
 * - 後3-4位：流水號
 *
 * 範例：
 * - 110001 = SYSTEM_ERROR（系統錯誤）
 * - 150001 = SCHED_OVERLAP（時段被佔用）
 * - 150002 = SCHED_BUFFER（緩衝時間不足）
 * - 160001 = EXCEPTION_NOT_FOUND（例外申請不存在）
 * - 160006 = EXCEPTION_DEADLINE_EXCEEDED（超過例外申請截止日）
 * - 160008 = EXCEPTION_ALREADY_PROCESSED（例外已被處理）
 * - 160009 = EXCEPTION_RESCHEDULE_CONFLICT（調課時間衝突）
 * - 160011 = EXCEPTION_CANCEL_DEADLINE_PASSED（停課截止日已過）
 */
export const NUMERIC_ERROR_CODE_MAP: Record<number, string> = {
  // 排課類錯誤 (150xxx)
  150001: 'SCHED_OVERLAP',
  150002: 'SCHED_BUFFER',
  150003: 'SCHED_PAST',
  150004: 'SCHED_LOCKED',
  150005: 'SCHED_CLOSED',
  150006: 'SCHED_INVALID_RANGE',
  150007: 'SCHED_RULE_CONFLICT',
  150008: 'SCHED_EXCEPTION_EXISTS',

  // 排課業務驗證錯誤 (150xxx)
  150009: 'SCHED_TEACHER_REQUIRED',
  150010: 'SCHED_ROOM_REQUIRED',
  150011: 'SCHED_OFFERING_NOT_FOUND',
  150012: 'SCHED_COURSE_NOT_FOUND',
  150013: 'SCHED_INVALID_WEEKDAY',
  150014: 'SCHED_INVALID_DURATION',
  150015: 'SCHED_START_AFTER_END',
  150016: 'SCHED_INVALID_DATE_FORMAT',
  150017: 'SCHED_END_BEFORE_START',
  150018: 'SCHED_DURATION_EXCEEDS_LIMIT',

  // 例外審核類錯誤 (160xxx)
  160001: 'EXCEPTION_NOT_FOUND',
  160002: 'EXCEPTION_INVALID_ACTION',
  160003: 'EXCEPTION_REVIEWED',
  160004: 'EXCEPTION_REVOKED',
  160005: 'EXCEPTION_REJECT_SELF',

  // 例外審核業務類錯誤 (160xxx)
  160006: 'EXCEPTION_DEADLINE_EXCEEDED',
  160007: 'EXCEPTION_SELF_REVIEW_FORBIDDEN',
  160008: 'EXCEPTION_ALREADY_PROCESSED',
  160009: 'EXCEPTION_RESCHEDULE_CONFLICT',
  160010: 'EXCEPTION_REPLACE_TEACHER_INVALID',
  160011: 'EXCEPTION_CANCEL_DEADLINE_PASSED',
  160012: 'EXCEPTION_RESCHEDULE_NO_NEW_TIME',

  // 循環編輯類錯誤 (160xxx)
  160013: 'RECURRENCE_EDIT_MODE_INVALID',
  160014: 'RECURRENCE_NO_AFFECTED_SESSIONS',
  160015: 'RECURRENCE_FUTURE_WITH_EDIT_DATE',
  160016: 'RECURRENCE_EDIT_DATE_REQUIRED',
  160017: 'RECURRENCE_DELETE_CONFIRM',
  160018: 'RECURRENCE_BATCH_LIMIT_EXCEEDED',

  // 系統類錯誤 (110xxx)
  110001: 'SYSTEM_ERROR',
  110002: 'PARAMS_VALIDATE_ERROR',
  110003: 'JSON_ENCODE_ERROR',
  110004: 'JSON_DECODE_ERROR',
  110009: 'RATE_LIMIT_EXCEEDED',

  // 資料庫類錯誤 (120xxx)
  120001: 'SQL_ERROR',
  120002: 'TX_ERROR',

  // 權限類錯誤 (130xxx)
  130001: 'UNAUTHORIZED',
  130002: 'FORBIDDEN',
  130003: 'TOKEN_EXPIRED',
  130004: 'INVALID_TOKEN',
  130005: 'INVALID_INVITE',

  // 業務資源類錯誤 (140xxx)
  140001: 'NOT_FOUND',
  140002: 'DUPLICATE',
  140003: 'TAG_INVALID',
  140004: 'LIMIT_EXCEEDED',
  140005: 'RESOURCE_IN_USE',
  140006: 'COURSE_IN_USE',
  140007: 'OFFERING_HAS_RULES',
  140008: 'ROOM_IN_USE',
  140009: 'INVALID_STATUS',
  140010: 'TEACHER_NOT_REGISTERED',

  // 檔案類錯誤 (170xxx)
  170001: 'FILE_TOO_LARGE',
  170002: 'FILE_TYPE_INVALID',
  170003: 'UPLOAD_FAILED',
  170004: 'CERTIFICATE_NOT_FOUND',

  // 搜尋媒合類錯誤 (180xxx)
  180001: 'TALENT_NOT_OPEN',

  // LINE Bot 類錯誤 (190xxx)
  190001: 'LINE_ALREADY_BOUND',
  190002: 'LINE_NOT_BOUND',
  190003: 'LINE_BINDING_CODE_INVALID',
  190004: 'LINE_BINDING_EXPIRED',
  190005: 'LINE_NOTIFY_FAILED',

  // 管理員類錯誤 (1100xxx)
  1100001: 'ADMIN_NOT_FOUND',
  1100002: 'ADMIN_EMAIL_EXISTS',
  1100003: 'PASSWORD_NOT_MATCH',
  1100004: 'ADMIN_CANNOT_DISABLE_SELF',
}

/**
 * 所有錯誤碼的聯合類型
 */
export type ErrorCode =
  | keyof typeof SUCCESS_CODES
  | keyof typeof SYSTEM_ERROR_CODES
  | keyof typeof DATABASE_ERROR_CODES
  | keyof typeof PERMISSION_ERROR_CODES
  | keyof typeof VALIDATION_ERROR_CODES
  | keyof typeof BUSINESS_ERROR_CODES
  | keyof typeof EXCEPTION_ERROR_CODES
  | keyof typeof TEACHER_ERROR_CODES
  | keyof typeof COURSE_ERROR_CODES
  | keyof typeof ROOM_ERROR_CODES
  | keyof typeof NOTIFICATION_ERROR_CODES

/**
 * 完整的錯誤碼對照表
 */
export const ERROR_MESSAGES: Record<string, string> = {
  ...SUCCESS_CODES,
  ...SYSTEM_ERROR_CODES,
  ...DATABASE_ERROR_CODES,
  ...PERMISSION_ERROR_CODES,
  ...VALIDATION_ERROR_CODES,
  ...BUSINESS_ERROR_CODES,
  ...EXCEPTION_ERROR_CODES,
  ...TEACHER_ERROR_CODES,
  ...COURSE_ERROR_CODES,
  ...ROOM_ERROR_CODES,
  ...NOTIFICATION_ERROR_CODES,
}

/**
 * 檢查錯誤碼是否為成功
 */
export const isSuccessCode = (code: string | number): boolean => {
  return code === 0 || code === '0' || code === 'SUCCESS'
}

/**
 * 檢查錯誤碼是否為權限相關
 */
export const isPermissionError = (code: string): boolean => {
  return Object.keys(PERMISSION_ERROR_CODES).includes(code)
}

/**
 * 檢查錯誤碼是否為驗證相關
 */
export const isValidationError = (code: string): boolean => {
  return Object.keys(VALIDATION_ERROR_CODES).includes(code)
}

/**
 * 檢查錯誤碼是否為需要登入的錯誤
 */
export const isUnauthorizedError = (code: string): boolean => {
  return ['UNAUTHORIZED', 'TOKEN_EXPIRED', 'TOKEN_INVALID'].includes(code)
}

/**
 * 取得 HTTP 狀態碼對應的錯誤碼
 */
export const httpStatusToErrorCode: Record<number, string> = {
  400: 'VALIDATION_ERROR',
  401: 'UNAUTHORIZED',
  403: 'FORBIDDEN',
  404: 'NOT_FOUND',
  409: 'CONFLICT',
  422: 'VALIDATION_ERROR',
  429: 'RATE_LIMIT_EXCEEDED',
  500: 'SYSTEM_ERROR',
  502: 'SYSTEM_ERROR',
  503: 'SYSTEM_ERROR',
}
