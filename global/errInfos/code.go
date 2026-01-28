package errInfos

// 專案流水號 (2位) + 功能類型 (2位) + 流水號 (4位)
// 1-1-0001 -> 110001 -> SQL_ERROR
// 1-2-0001 -> 120001 -> USER_NOT_FOUND

type ErrCode int

// 系統相關 (1)
const (
	SYSTEM_ERROR          ErrCode = 10001
	PARAMS_VALIDATE_ERROR ErrCode = 10002
	JSON_ENCODE_ERROR     ErrCode = 10003
	JSON_DECODE_ERROR     ErrCode = 10004
	JSON_PROCESS_ERROR    ErrCode = 10005
	FORMAT_RESOURCE_ERROR ErrCode = 10006
	REQUEST_TIMEOUT       ErrCode = 10007
	NOT_IMPLEMENTED       ErrCode = 10008
	RATE_LIMIT_EXCEEDED   ErrCode = 10009
)

// 資料庫、快取相關 (2)
const (
	SQL_ERROR ErrCode = 20001
	TX_ERROR  ErrCode = 20002
)

// 權限與認證 (3)
const (
	UNAUTHORIZED   ErrCode = 30001
	FORBIDDEN      ErrCode = 30002
	TOKEN_EXPIRED  ErrCode = 30003
	INVALID_TOKEN  ErrCode = 30004
	INVALID_INVITE ErrCode = 30005
)

// 業務資源類 (4)
const (
	NOT_FOUND          ErrCode = 40001
	DUPLICATE          ErrCode = 40002
	TAG_INVALID        ErrCode = 40003
	LIMIT_EXCEEDED     ErrCode = 40004
	RESOURCE_IN_USE    ErrCode = 40005
	COURSE_IN_USE      ErrCode = 40006
	OFFERING_HAS_RULES ErrCode = 40007
	ROOM_IN_USE        ErrCode = 40008
	INVALID_STATUS     ErrCode = 40009
)

// 排課核心類 (5)
const (
	SCHED_OVERLAP          ErrCode = 50001
	SCHED_BUFFER           ErrCode = 50002
	SCHED_PAST             ErrCode = 50003
	SCHED_LOCKED           ErrCode = 50004
	SCHED_CLOSED           ErrCode = 50005
	SCHED_INVALID_RANGE    ErrCode = 50006
	SCHED_RULE_CONFLICT    ErrCode = 50007
	SCHED_EXCEPTION_EXISTS ErrCode = 50008
)

// 例外與審核類 (6)
const (
	EXCEPTION_NOT_FOUND      ErrCode = 60001
	EXCEPTION_INVALID_ACTION ErrCode = 60002
	EXCEPTION_REVIEWED       ErrCode = 60003
	EXCEPTION_REVOKED        ErrCode = 60004
	EXCEPTION_REJECT_SELF    ErrCode = 60005
)

// 檔案與媒體類 (7)
const (
	FILE_TOO_LARGE        ErrCode = 70001
	FILE_TYPE_INVALID     ErrCode = 70002
	UPLOAD_FAILED         ErrCode = 70003
	CERTIFICATE_NOT_FOUND ErrCode = 70004
)

// 搜尋與媒合類 (8)
const (
	TALENT_NOT_OPEN ErrCode = 80001
)

// LINE Bot 與通知類 (9)
const (
	LINE_ALREADY_BOUND        ErrCode = 90001
	LINE_NOT_BOUND            ErrCode = 90002
	LINE_BINDING_CODE_INVALID ErrCode = 90003
	LINE_BINDING_EXPIRED      ErrCode = 90004
	LINE_NOTIFY_FAILED        ErrCode = 90005
)

// 管理員類 (10)
const (
	ADMIN_NOT_FOUND           ErrCode = 100001
	ADMIN_EMAIL_EXISTS        ErrCode = 100002
	PASSWORD_NOT_MATCH        ErrCode = 100003
	ADMIN_CANNOT_DISABLE_SELF ErrCode = 100004
)
