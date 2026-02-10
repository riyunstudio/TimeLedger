package errInfos

var messagesMap = map[ErrCode]message{
	// 系統相關
	SYSTEM_ERROR:          {EN: "System error.", TW: "系統錯誤", CN: "系统错误"},
	PARAMS_VALIDATE_ERROR: {EN: "Invalid parameters", TW: "參數格式錯誤", CN: "参数格式错误"},
	JSON_ENCODE_ERROR:     {EN: "Json encode error.", TW: "JSON 編碼錯誤", CN: "JSON 编码错误"},
	JSON_DECODE_ERROR:     {EN: "Json decode error.", TW: "JSON 解碼錯誤", CN: "JSON 解码错误"},
	JSON_PROCESS_ERROR:    {EN: "Json process error.", TW: "JSON 處理錯誤", CN: "JSON 处理错误"},
	FORMAT_RESOURCE_ERROR: {EN: "Format resource error.", TW: "資源格式錯誤", CN: "资源格式错误"},
	REQUEST_TIMEOUT:       {EN: "Request timeout.", TW: "請求超時", CN: "请求超时"},
	NOT_IMPLEMENTED:       {EN: "Feature not implemented", TW: "功能尚未實作", CN: "功能尚未实现"},
	RATE_LIMIT_EXCEEDED:   {EN: "Rate limit exceeded", TW: "請求頻率過高，請稍後再試", CN: "请求频率过高，请稍后再试"},

	// 資料庫、快取相關
	SQL_ERROR: {EN: "Database operation failed", TW: "資料庫操作失敗", CN: "数据库操作失败"},
	TX_ERROR:  {EN: "SQL transaction error.", TW: "交易錯誤", CN: "交易错误"},

	// 權限與認證
	UNAUTHORIZED:   {EN: "Please login first", TW: "請先登入", CN: "请先登录"},
	FORBIDDEN:      {EN: "Permission denied", TW: "權限不足", CN: "权限不足"},
	TOKEN_EXPIRED:  {EN: "Token expired", TW: "Token 已過期", CN: "Token 已过期"},
	INVALID_TOKEN:  {EN: "Invalid token", TW: "無效的 Token", CN: "无效的 Token"},
	INVALID_INVITE: {EN: "Invalid or expired invite", TW: "邀請碼無效或已過期", CN: "邀请码无效或已过期"},

	// 業務資源類
	NOT_FOUND:          {EN: "Resource not found", TW: "找不到資源", CN: "找不到资源"},
	DUPLICATE:          {EN: "Resource already exists", TW: "資源已存在", CN: "资源已存在"},
	TAG_INVALID:        {EN: "Invalid hashtag format", TW: "標籤格式錯誤", CN: "标签格式错误"},
	LIMIT_EXCEEDED:     {EN: "Plan limit reached", TW: "超過方案配額", CN: "超过方案配额"},
	RESOURCE_IN_USE:    {EN: "Resource is in use", TW: "資源仍在使用中", CN: "资源仍在使用中"},
	COURSE_IN_USE:      {EN: "Course has active offerings", TW: "課程模板仍有關聯班別", CN: "课程模板仍有关联班别"},
	OFFERING_HAS_RULES: {EN: "Offering has schedule rules", TW: "班別仍有排課規則", CN: "班别仍有排课规则"},
	ROOM_IN_USE:        {EN: "Room has active schedules", TW: "教室仍有排課安排", CN: "教室仍有排课安排"},
	INVALID_STATUS:     {EN: "Invalid status transition", TW: "不允許的狀態轉換", CN: "不允许的状态转换"},

	// 排課核心類
	SCHED_OVERLAP:          {EN: "Time slot occupied", TW: "時段被佔用", CN: "时段被占用"},
	SCHED_BUFFER:           {EN: "Insufficient buffer time", TW: "緩衝時間不足", CN: "缓冲时间不足"},
	SCHED_PAST:             {EN: "Cannot book past time", TW: "不能排過去的時間", CN: "不能排过去的时间"},
	SCHED_LOCKED:           {EN: "Slot is locked by another", TW: "時段已被鎖定", CN: "时段已被锁定"},
	SCHED_CLOSED:           {EN: "Center is closed", TW: "非營業時間", CN: "非营业时间"},
	SCHED_INVALID_RANGE:    {EN: "Invalid date range", TW: "日期範圍錯誤", CN: "日期范围错误"},
	SCHED_RULE_CONFLICT:    {EN: "Rule conflict detected", TW: "規則衝突", CN: "规则冲突"},
	SCHED_EXCEPTION_EXISTS: {EN: "Exception already exists", TW: "該日期已有例外單", CN: "该日期已有例外单"},

	// 例外與審核類
	EXCEPTION_NOT_FOUND:           {EN: "Exception request not found", TW: "例外申請不存在", CN: "例外申请不存在"},
	EXCEPTION_INVALID_ACTION:      {EN: "Invalid action for current status", TW: "當前狀態不允許此操作", CN: "当前状态不允许此操作"},
	EXCEPTION_REVIEWED:            {EN: "Exception already reviewed", TW: "例外已審核過", CN: "例外已审核过"},
	EXCEPTION_REVOKED:             {EN: "Exception was revoked", TW: "例外已撤回", CN: "例外已撤回"},
	EXCEPTION_REJECT_SELF:         {EN: "Cannot reject own request", TW: "不能拒絕自己提交的申請", CN: "不能拒绝自己提交的申请"},
	EXCEPTION_DEADLINE_EXCEEDED:   {EN: "Deadline exceeded, application must be submitted 14 days in advance", TW: "已超過異動截止日（需提前 14 天申請）", CN: "已超过异动截止日（需提前 14 天申请）"},
	EXCEPTION_ALREADY_PROCESSED:    {EN: "Exception has already been processed", TW: "例外已被處理", CN: "例外已被处理"},
	EXCEPTION_RESCHEDULE_CONFLICT:  {EN: "Reschedule time conflicts with existing schedule", TW: "調課時間與現有排程衝突", CN: "调课时间与现有排程冲突"},
	EXCEPTION_CANCEL_DEADLINE_PASSED: {EN: "Cancellation deadline has passed", TW: "停課截止日已過", CN: "停课截止日已过"},

	// 檔案與媒體類
	FILE_TOO_LARGE:        {EN: "File size exceeds limit", TW: "檔案超過大小限制", CN: "文件超过大小限制"},
	FILE_TYPE_INVALID:     {EN: "Invalid file type", TW: "不支援的檔案類型", CN: "不支持的文件类型"},
	UPLOAD_FAILED:         {EN: "Upload failed", TW: "上傳失敗", CN: "上传失败"},
	CERTIFICATE_NOT_FOUND: {EN: "Certificate not found", TW: "證照不存在", CN: "证书不存在"},

	// 搜尋與媒合類
	TALENT_NOT_OPEN: {EN: "Talent search not available", TW: "該老師未開放搜尋", CN: "该老师未开放搜索"},

	// LINE Bot 與通知類
	LINE_ALREADY_BOUND:        {EN: "LINE account already bound", TW: "LINE 帳號已綁定", CN: "LINE 账号已绑定"},
	LINE_NOT_BOUND:            {EN: "LINE account not bound", TW: "LINE 帳號未綁定", CN: "LINE 账号未绑定"},
	LINE_BINDING_CODE_INVALID: {EN: "Invalid binding code", TW: "驗證碼無效", CN: "验证码无效"},
	LINE_BINDING_EXPIRED:      {EN: "Binding code expired", TW: "驗證碼已過期，請重新產生", CN: "验证码已过期，请重新产生"},
	LINE_NOTIFY_FAILED:        {EN: "Failed to send LINE notification", TW: "LINE 通知發送失敗", CN: "LINE 通知发送失败"},

	// 管理員類
	ADMIN_NOT_FOUND:           {EN: "Admin user not found", TW: "管理員不存在", CN: "管理员不存在"},
	ADMIN_EMAIL_EXISTS:        {EN: "Email already registered", TW: "Email 已被註冊", CN: "Email 已被注册"},
	PASSWORD_NOT_MATCH:        {EN: "Current password is incorrect", TW: "舊密碼錯誤", CN: "旧密码错误"},
	ADMIN_CANNOT_DISABLE_SELF: {EN: "Cannot disable yourself", TW: "不能停用自己的帳號", CN: "不能停用自己的账号"},

	// 資源鎖定與衝突類
	ERR_RESOURCE_LOCKED:     {EN: "Resource is locked by another operation", TW: "資源正在被其他操作修改，請稍後再試", CN: "资源正在被其他操作修改，请稍后再试"},
	ERR_CONCURRENT_MODIFIED: {EN: "Resource was modified by another request", TW: "資源已被其他請求修改，請重新整理後再試", CN: "资源已被其他请求修改，请重新整理后再试"},
	ERR_TX_FAILED:           {EN: "Transaction failed", TW: "交易執行失敗，請稍後再試", CN: "交易执行失败，请稍后再试"},
}
