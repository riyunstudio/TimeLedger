package global

type ContextKey string

// gRPC / Gin 客戶端請求時間
const RequestTimeKey = ContextKey("RequestTime")

// Gin 客戶端請求參數
const ArgsBodyKey = ContextKey("ArgsBody")

// Gin 客戶端請求 Session ID
const SidKey = ContextKey("Sid")

// gRPC 客戶端請求 Metedata (Headers)
const MeteDataKey = ContextKey("MD")

// gRPC / Gin 客戶端請求 Trade ID
const TraceIDKey = ContextKey("Tid")

// gRPC / Gin 客戶端請求 IP
const XForwardForKey = ContextKey("x-forwarded-for")

// Gin 客戶端 CTX 處理結果
const RetKey = ContextKey("Ret")

// Auth Keys
const UserIDKey = ContextKey("user_id")
const UserTypeKey = ContextKey("user_type")
const CenterIDKey = ContextKey("center_id")
const LineUserIDKey = ContextKey("line_user_id")

// HTTP Status Codes
const UNAUTHORIZED = 401
const FORBIDDEN = 403
const BAD_REQUEST = 400
const INTERNAL_SERVER_ERROR = 500
