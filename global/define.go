package global

type contextKey string

// gRPC / Gin 客戶端請求時間
const RequestTimeKey = contextKey("RequestTime")

// Gin 客戶端請求參數
const ArgsBodyKey = contextKey("ArgsBody")

// Gin 客戶端請求 Session ID
const SidKey = contextKey("Sid")

// gRPC 客戶端請求 Metedata (Headers)
const MeteDataKey = contextKey("MD")

// gRPC / Gin 客戶端請求 Trade ID
const TraceIDKey = contextKey("Tid")

// gRPC / Gin 客戶端請求 IP
const XForwardForKey = contextKey("x-forwarded-for")

// Gin 客戶端 CTX 處理結果
const RetKey = contextKey("Ret")
