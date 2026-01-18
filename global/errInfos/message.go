package errInfos

// 訊息表
var messagesMap = map[ErrCode]message{
	SYSTEM_ERROR:          {EN: "System error."},
	PARAMS_VALIDATE_ERROR: {EN: "Request payload validate fail."},
	JSON_ENCODE_ERROR:     {EN: "Json encode error."},
	JSON_DECODE_ERROR:     {EN: "Json decode error."},
	JSON_PROCESS_ERROR:    {EN: "Json process error."},
	FORMAT_RESOURCE_ERROR: {EN: "Format resource error."},
	REQUEST_TIMEOUT:       {EN: "Request timeout."},
	SQL_ERROR:             {EN: "SQL syntax error."},
	TX_ERROR:              {EN: "SQL transaction error."},
}
