package errInfos

// 訊息表
var messagesMap = map[int]message{
	int(SYSTEM_ERROR):          {EN: "System error."},
	int(PARAMS_VALIDATE_ERROR): {EN: "Request payload validate fail."},
	int(JSON_ENCODE_ERROR):     {EN: "Json encode error."},
	int(JSON_DECODE_ERROR):     {EN: "Json decode error."},
	int(JSON_PROCESS_ERROR):    {EN: "Json process error."},
	int(FORMAT_RESOURCE_ERROR): {EN: "Format resource error."},
	int(REQUEST_TIMEOUT):       {EN: "Request timeout."},

	int(SQL_ERROR): {EN: "SQL syntax error."},
	int(TX_ERROR):  {EN: "SQL transaction error."},
}
