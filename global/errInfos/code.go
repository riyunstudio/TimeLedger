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
)

// 資料庫、快取相關 (2)
const (
	SQL_ERROR ErrCode = 20001
	TX_ERROR  ErrCode = 20002
)

// 其他相關 (3)
const ()

// 使用者相關 (4)
const (
	USER_NOT_FOUNT ErrCode = 40001
)
