package errInfos

// 語系設定
const (
	LANG_EN = "en"
	LANG_TW = "tw"
	LANG_CN = "cn"
)

type message struct {
	EN string
	TW string
	CN string
}

type Err struct {
	Code int
	Msg  string
}

var SUPORT_LANG = []string{LANG_EN, LANG_TW, LANG_CN}

// 透過錯誤代碼取得對應的訊息
func New(code int, lang ...string) *Err {
	msgs, exists := messagesMap[code]

	if !exists {
		// 未定義訊息
		msgs = messagesMap[-1]
	}

	// 如果沒傳 lang 參數，預設 EN
	selectedLang := LANG_EN
	if len(lang) > 0 && lang[0] != "" {
		selectedLang = lang[0]
	}

	switch selectedLang {
	case LANG_EN:
		return &Err{Code: code, Msg: msgs.EN}
	case LANG_TW:
		return &Err{Code: code, Msg: msgs.TW}
	case LANG_CN:
		return &Err{Code: code, Msg: msgs.CN}
	default:
		return &Err{Code: code, Msg: msgs.EN}
	}
}
