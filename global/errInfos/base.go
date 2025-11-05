package errInfos

import (
	"fmt"
	"strconv"
)

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

var SUPORT_LANG = []string{LANG_EN, LANG_TW, LANG_CN}

type ErrInfo struct {
	appID int
}

func Initialize(appID int) *ErrInfo {
	return &ErrInfo{
		appID: appID,
	}
}

type Res struct {
	Code int
	Msg  string
}

// 透過錯誤代碼取得對應的訊息
func (e *ErrInfo) New(code int, lang ...string) *Res {
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

	codeStr := fmt.Sprintf("%d%d", e.appID, code)
	code, _ = strconv.Atoi(codeStr)

	switch selectedLang {
	case LANG_EN:
		return &Res{Code: code, Msg: msgs.EN}
	case LANG_TW:
		return &Res{Code: code, Msg: msgs.TW}
	case LANG_CN:
		return &Res{Code: code, Msg: msgs.CN}
	default:
		return &Res{Code: code, Msg: msgs.EN}
	}
}
