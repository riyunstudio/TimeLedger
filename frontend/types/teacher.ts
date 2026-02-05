/**
 * 教師相關類型定義
 *
 * 包含教師個人檔案、技能、證照、標籤等相關類型
 */

import type { ID, Timestamp, ApiResponse, PaginationParams, PaginatedResponse, DateRange } from './api'

// ==================== 教師基礎類型 ====================

/**
 * 教師用戶
 */
export interface Teacher {
  /** 教師 ID */
  id: ID
  /** 名稱 */
  name: string
  /** 電子郵件 */
  email?: string
  /** 電話 */
  phone?: string
  /** LINE 用戶 ID (帳號綁定，不可解除) */
  line_user_id?: string
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 公開聯絡資訊 */
  public_contact_info?: PublicContactInfo
  /** 技能列表 */
  skills?: TeacherSkill[]
  /** 證照列表 */
  certificates?: TeacherCertificate[]
  /** 個人標籤 */
  personal_hashtags?: PersonalHashtag[]
  /** 是否開放應徵 (人才庫) */
  is_open_to_hiring: boolean
  /** 是否已激活 */
  is_active: boolean
  /** 是否為佔位老師 (待綁定) */
  is_placeholder?: boolean
  /** 邀請時間 */
  invited_at?: Timestamp
  /** 激活時間 */
  activated_at?: Timestamp
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

// ==================== 技能相關類型 ====================

/**
 * 技能分類配置
 */
export const SKILL_CATEGORIES = {
  MUSIC: { label: '音樂', color: 'bg-pink-500/20 text-pink-400 border-pink-500/30', icon: '🎵' },
  ART: { label: '美術', color: 'bg-purple-500/20 text-purple-400 border-purple-500/30', icon: '🎨' },
  DANCE: { label: '舞蹈', color: 'bg-orange-500/20 text-orange-400 border-orange-500/30', icon: '💃' },
  LANGUAGE: { label: '語言', color: 'bg-blue-500/20 text-blue-400 border-blue-500/30', icon: '🗣️' },
  SPORTS: { label: '運動', color: 'bg-green-500/20 text-green-400 border-green-500/30', icon: '⚽' },
  OTHER: { label: '其他', color: 'bg-slate-500/20 text-slate-400 border-slate-500/30', icon: '✨' },
} as const

/**
 * 技能分類類型
 */
export type SkillCategory = keyof typeof SKILL_CATEGORIES

/**
 * 教師技能
 */
export interface TeacherSkill {
  /** 技能 ID */
  id: ID
  /** 所屬教師 ID */
  teacher_id: ID
  /** 技能分類 */
  category: SkillCategory | string
  /** 技能名稱 */
  skill_name: string
  /** 技能標籤 */
  hashtags?: TeacherSkillHashtag[]
}

/**
 * 技能標籤關聯
 */
export interface TeacherSkillHashtag {
  /** 關聯 ID */
  id: ID
  /** 技能 ID */
  teacher_skill_id: ID
  /** 標籤 ID */
  hashtag_id: ID
  /** 標籤資訊 */
  hashtag?: Hashtag
}

/**
 * 新增技能請求
 */
export interface CreateSkillRequest {
  /** 技能分類 */
  category: SkillCategory | string
  /** 技能名稱 */
  skill_name: string
  /** 標籤 ID 陣列 */
  hashtag_ids?: ID[]
}

/**
 * 更新技能請求
 */
export interface UpdateSkillRequest {
  /** 技能分類 */
  category?: SkillCategory | string
  /** 技能名稱 */
  skill_name?: string
  /** 標籤名稱陣列 */
  hashtags?: string[]
}

// ==================== 證照相關類型 ====================

/**
 * 教師證照
 */
export interface TeacherCertificate {
  /** 證照 ID */
  id: ID
  /** 所屬教師 ID */
  teacher_id: ID
  /** 證照名稱 */
  certificate_name: string
  /** 發證單位 */
  issued_by?: string
  /** 發證日期 */
  issued_date?: string
  /** 證照圖片 URL */
  file_url?: string
  /** 建立時間 */
  created_at: Timestamp
  /** 更新時間 */
  updated_at: Timestamp
}

/**
 * 新增證照請求
 */
export interface CreateCertificateRequest {
  /** 證照名稱 */
  name: string
  /** 發證單位 */
  issued_by?: string
  /** 發證日期 (YYYY-MM-DD) */
  issued_at?: string
  /** 證照檔案 */
  file?: File
}

// ==================== 標籤相關類型 ====================

/**
 * 標籤
 */
export interface Hashtag {
  /** 標籤 ID */
  id: ID
  /** 標籤名稱 (含 # 符號) */
  name: string
  /** 使用次數 */
  usage_count: number
}

/**
 * 個人標籤
 */
export interface PersonalHashtag {
  /** 關聯 ID */
  id: ID
  /** 標籤 ID */
  hashtag_id: ID
  /** 標籤名稱 */
  name: string
}

/**
 * 教師個人標籤關聯
 */
export interface TeacherPersonalHashtag {
  /** 關聯 ID */
  id: ID
  /** 教師 ID */
  teacher_id: ID
  /** 標籤 ID */
  hashtag_id: ID
  /** 標籤資訊 */
  hashtag?: Hashtag
}

// ==================== 教師中心關聯 ====================

/**
 * 教師與中心關聯
 */
export interface TeacherCenterMembership {
  /** 關聯 ID */
  id: ID
  /** 中心 ID */
  center_id: ID
  /** 中心名稱 */
  center_name?: string
  /** 教師 ID */
  teacher_id: ID
  /** 狀態 */
  status: 'ACTIVE' | 'INACTIVE' | 'INVITED'
  /** 加入時間 */
  joined_at?: Timestamp
  /** 邀請時間 */
  invited_at?: Timestamp
}

// ==================== 教師邀請相關類型 ====================

/**
 * 教師收到的邀請
 */
export interface Invitation {
  /** 邀請 ID */
  id: number
  /** 中心 ID */
  center_id: number
  /** 中心名稱 */
  center_name: string
  /** 邀請類型 */
  invite_type: 'TALENT_POOL' | 'TEACHER' | 'MEMBER'
  /** 狀態 */
  status: 'PENDING' | 'ACCEPTED' | 'DECLINED' | 'EXPIRED'
  /** 邀請訊息 */
  message?: string
  /** 創建時間 */
  created_at: string
  /** 過期時間 */
  expires_at: string | null
  /** 回應時間 */
  responded_at?: string
  /** 中心標誌 URL */
  center_logo_url?: string
}

/**
 * 邀請請求回應
 */
export interface InvitationRespondRequest {
  /** 邀請 ID */
  invitation_id: number
  /** 回應動作 */
  response: 'ACCEPT' | 'REJECT'
}

/**
 * 待處理邀請數量回應
 */
export interface PendingInvitationCountResponse {
  /** 待處理邀請數量 */
  count: number
}

/**
 * 邀請教師請求
 */
export interface InviteTeacherRequest {
  /** 教師 ID */
  teacher_id?: ID
  /** 教師名稱 (若無 ID) */
  name?: string
  /** 教師郵箱 (若無 ID) */
  email?: string
  /** 訊息 */
  message?: string
}

/**
 * 邀請回應
 */
export interface InviteTeacherResponse {
  /** 邀請 ID */
  id: ID
  /** 邀請 Token */
  token: string
  /** 狀態 */
  status: 'PENDING' | 'ACCEPTED' | 'EXPIRED' | 'CANCELLED'
  /** 創建時間 */
  created_at: Timestamp
  /** 過期時間 */
  expires_at: Timestamp
}

// ==================== 教師列表相關類型 ====================

/**
 * 教師列表查詢參數
 */
export interface TeacherListParams extends PaginationParams {
  /** 搜尋關鍵字 */
  keyword?: string
  /** 縣市篩選 */
  city?: string
  /** 技能分類篩選 */
  category?: string
  /** 技能名稱篩選 */
  skill_name?: string
  /** 是否僅顯示開放應徵 */
  open_to_hiring?: boolean
  /** 是否僅顯示已激活 */
  is_active?: boolean
}

/**
 * 教師列表項目 (簡化版)
 */
export interface TeacherListItem {
  /** 教師 ID */
  id: ID
  /** 名稱 */
  name: string
  /** 電子郵件 */
  email?: string
  /** 頭像 URL */
  avatar_url?: string
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 是否開放應徵 */
  is_open_to_hiring: boolean
  /** 技能數量 */
  skill_count: number
  /** 證照數量 */
  certificate_count: number
  /** 平均評分 */
  average_rating?: number
  /** 狀態 */
  status: 'ACTIVE' | 'INACTIVE' | 'INVITED'
}

/**
 * 教師列表 API 回應
 */
export type TeacherListResponse = PaginatedResponse<TeacherListItem>

// ==================== 教師邀請列表 ====================

/**
 * 教師邀請列表查詢參數
 */
export interface TeacherInvitationListParams extends PaginationParams {
  /** 狀態篩選 */
  status?: 'PENDING' | 'ACCEPTED' | 'EXPIRED' | 'CANCELLED'
  /** 中心 ID */
  center_id?: ID
}

/**
 * 教師邀請列表項目
 */
export interface TeacherInvitationItem {
  /** 邀請 ID */
  id: ID
  /** 中心 ID */
  center_id: ID
  /** 中心名稱 */
  center_name: string
  /** 教師名稱 */
  teacher_name?: string
  /** 教師郵箱 */
  teacher_email?: string
  /** 狀態 */
  status: 'PENDING' | 'ACCEPTED' | 'EXPIRED' | 'CANCELLED'
  /** 發送時間 */
  sent_at: Timestamp
  /** 接受時間 */
  accepted_at?: Timestamp
  /** 過期時間 */
  expires_at: Timestamp
}

/**
 * 教師邀請列表回應
 */
export type TeacherInvitationListResponse = PaginatedResponse<TeacherInvitationItem>

// ==================== 教師設定檔回應 ====================

/**
 * 教師完整設定檔
 */
export interface TeacherProfile {
  /** 教師基本資訊 */
  profile: Teacher
  /** 技能列表 */
  skills: TeacherSkill[]
  /** 證照列表 */
  certificates: TeacherCertificate[]
  /** 個人標籤 */
  hashtags: PersonalHashtag[]
  /** 加入的中心列表 */
  centers: TeacherCenterMembership[]
}

/**
 * 教師設定檔 API 回應
 */
export type TeacherProfileResponse = ApiResponse<TeacherProfile>

// ==================== 公開聯絡資訊 ====================

/**
 * 公開聯絡資訊
 */
export interface PublicContactInfo {
  /** Instagram */
  instagram?: string
  /** YouTube */
  youtube?: string
  /** 個人網站 */
  website?: string
  /** 其他 */
  other?: string
}

/**
 * 更新教師設定檔請求
 */
export interface UpdateTeacherProfileRequest {
  /** 名稱 */
  name?: string
  /** 電話 */
  phone?: string
  /** 縣市 */
  city?: string
  /** 區域 */
  district?: string
  /** 公開聯絡資訊 */
  public_contact_info?: PublicContactInfo
  /** 個人簡介 */
  bio?: string
}

// ==================== 教師評價相關類型 ====================

/**
 * 教師評價
 */
export interface TeacherRating {
  /** 評價 ID */
  id: ID
  /** 教師 ID */
  teacher_id: ID
  /** 評價者 ID */
  reviewer_id: ID
  /** 評分 (1-5) */
  rating: number
  /** 評價內容 */
  comment?: string
  /** 中心 ID */
  center_id: ID
  /** 建立時間 */
  created_at: Timestamp
}

/**
 * 教師評價統計
 */
export interface TeacherRatingStats {
  /** 教師 ID */
  teacher_id: ID
  /** 評價數量 */
  count: number
  /** 平均評分 */
  average_rating: number
  /** 評分分佈 */
  distribution: {
    1: number
    2: number
    3: number
    4: number
    5: number
  }
}
