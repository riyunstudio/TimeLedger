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

// 資源鎖定與衝突類 (11) - 用於交易衝突和資源鎖定場景
const (
	ERR_RESOURCE_LOCKED     ErrCode = 110001 // 資源被鎖定（如另一筆交易正在修改）
	ERR_CONCURRENT_MODIFIED ErrCode = 110002 // 並發修改衝突
	ERR_TX_FAILED           ErrCode = 110003 // 交易執行失敗
)

// 交易執行類錯誤 (12) - 用於交易過程中的各種錯誤場景
const (
	ERR_TX_TIMEOUT              ErrCode = 120001 // 交易執行超時
	ERR_DEADLOCK_DETECTED       ErrCode = 120002 // 偵測到資料庫死鎖
	ERR_PARTIAL_COMPLETION      ErrCode = 120003 // 交易部分完成（部分操作成功）
	ERR_ROLLBACK_FAILED         ErrCode = 120004 // 交易回滾失敗
	ERR_CONSTRAINT_VIOLATION    ErrCode = 120005 // 違反資料庫約束
	ERR_FOREIGN_KEY_VIOLATION   ErrCode = 120006 // 外鍵約束違反
	ERR_UNIQUE_VIOLATION        ErrCode = 120007 // 唯一約束違反
	ERR_CHECK_VIOLATION         ErrCode = 120008 // CHECK 約束違反
	ERR_SERIALIZATION_FAILURE   ErrCode = 120009 // 序列化失敗（並發衝突）
	ERR_LOCK_WAIT_TIMEOUT       ErrCode = 120010 // 鎖等待超時
	ERR_SHARE_LOCK_FAILED       ErrCode = 120011 // 共享鎖獲取失敗
	ERR_EXCLUSIVE_LOCK_FAILED   ErrCode = 120012 // 排他鎖獲取失敗
)

// 排課業務驗證類錯誤 (50)
const (
	SCHED_TEACHER_REQUIRED       ErrCode = 50009 // 必須指定老師
	SCHED_ROOM_REQUIRED          ErrCode = 50010 // 必須指定教室
	SCHED_OFFERING_NOT_FOUND     ErrCode = 50011 // 班別不存在
	SCHED_COURSE_NOT_FOUND       ErrCode = 50012 // 課程模板不存在
	SCHED_INVALID_WEEKDAY        ErrCode = 50013 // 無效的星期幾
	SCHED_INVALID_DURATION       ErrCode = 50014 // 無效的課程時長
	SCHED_START_AFTER_END        ErrCode = 50015 // 開始時間晚於結束時間
	SCHED_INVALID_DATE_FORMAT    ErrCode = 50016 // 無效的日期格式
	SCHED_END_BEFORE_START       ErrCode = 50017 // 結束日期早於開始日期
	SCHED_DURATION_EXCEEDS_LIMIT ErrCode = 50018 // 課程時長超過限制
)

// 例外審核業務類錯誤 (60)
const (
	EXCEPTION_DEADLINE_EXCEEDED       ErrCode = 60006 // 超過例外申請截止日
	EXCEPTION_SELF_REVIEW_FORBIDDEN   ErrCode = 60007 // 不能審核自己提交的申請
	EXCEPTION_ALREADY_PROCESSED       ErrCode = 60008 // 例外已被處理
	EXCEPTION_RESCHEDULE_CONFLICT     ErrCode = 60009 // 調課時間與現有排程衝突
	EXCEPTION_REPLACE_TEACHER_INVALID ErrCode = 60010 // 代課老師無效
	EXCEPTION_CANCEL_DEADLINE_PASSED  ErrCode = 60011 // 停課截止日已過
	EXCEPTION_RESCHEDULE_NO_NEW_TIME  ErrCode = 60012 // 調課必須提供新時間
)

// 循環編輯業務類錯誤 (60)
const (
	RECURRENCE_EDIT_MODE_INVALID    ErrCode = 60013 // 無效的編輯模式
	RECURRENCE_NO_AFFECTED_SESSIONS ErrCode = 60014 // 沒有受影響的場次
	RECURRENCE_FUTURE_WITH_EDIT_DATE ErrCode = 60015 // FUTURE 模式必須指定編輯日期
	RECURRENCE_EDIT_DATE_REQUIRED   ErrCode = 60016 // 編輯日期為必填
	RECURRENCE_DELETE_CONFIRM       ErrCode = 60017 // 刪除操作需要確認
	RECURRENCE_BATCH_LIMIT_EXCEEDED ErrCode = 60018 // 批量操作超過限制
)
