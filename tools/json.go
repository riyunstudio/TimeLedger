package tools

import (
	"encoding/json"
	"errors"
)

// MergeJSONStrings 將原本的 JSON array 字串與新的 string slice 合併，並去重
func (tl *Tools) MergeJSONStrings(oldJSON string, newItems []string) (string, error) {
	if len(newItems) == 0 {
		// 沒有新的資料，直接回傳原本的 JSON
		if oldJSON == "" {
			return "[]", nil
		}
		return oldJSON, nil
	}

	// 先解析舊的 JSON
	var existing []string
	if len(oldJSON) > 0 {
		if err := json.Unmarshal([]byte(oldJSON), &existing); err != nil {
			return "", errors.New("failed to decode old JSON: " + err.Error())
		}
	}

	// 用 map 去重
	itemMap := make(map[string]struct{})
	for _, item := range existing {
		itemMap[item] = struct{}{}
	}
	for _, item := range newItems {
		if _, exists := itemMap[item]; !exists {
			existing = append(existing, item)
			itemMap[item] = struct{}{}
		}
	}

	// 轉回 JSON
	result, err := json.Marshal(existing)
	if err != nil {
		return "", errors.New("failed to encode merged JSON: " + err.Error())
	}

	return string(result), nil
}
