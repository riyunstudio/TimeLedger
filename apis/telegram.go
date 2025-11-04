package apis

import (
	"akali/libs"
	"fmt"
)

func (api *Api) PostTelegramMessage(tid string, text string) (err error) {
	req := libs.CurlInit(api.Tools)
	req.SetHttps().NewRequest("api.telegram.org", "443", fmt.Sprintf("/bot%s/sendMessage", api.Env.TelegramBotToken))
	req.SetQueries(map[string]any{
		"chat_id": api.Env.TelegramChatID,
		"text":    text,
	})
	req.SetTraceID(tid)
	_, err = req.Post()

	return
}
