package services

import (
	"fmt"
	"timeLedger/app/models"
)

// LineBotTemplateService Flex Message 範本服務
type LineBotTemplateService interface {
	// 取得歡迎訊息範本
	GetWelcomeTeacherTemplate(teacher *models.Teacher, centerName string) interface{}
	GetWelcomeAdminTemplate(admin *models.AdminUser, centerName string) interface{}

	// 取得例外通知範本
	GetExceptionSubmitTemplate(exception *models.ScheduleException, teacherName string, centerName string) interface{}
	GetExceptionApproveTemplate(exception *models.ScheduleException, teacherName string) interface{}
	GetExceptionRejectTemplate(exception *models.ScheduleException, teacherName string, reason string) interface{}

	// 取得邀請通知範本
	GetInvitationAcceptedTemplate(teacher *models.Teacher, centerName string, role string) interface{}
}

// LineBotTemplateServiceImpl Flex Message 範本服務實現
type LineBotTemplateServiceImpl struct {
	baseURL string // 前端網站 URL
}

func NewLineBotTemplateService(baseURL string) LineBotTemplateService {
	return &LineBotTemplateServiceImpl{
		baseURL: baseURL,
	}
}

// GetWelcomeTeacherTemplate 老師歡迎訊息範本
func (s *LineBotTemplateServiceImpl) GetWelcomeTeacherTemplate(teacher *models.Teacher, centerName string) interface{} {
	bindURL := fmt.Sprintf("%s/teacher/bind?teacher_id=%d", s.baseURL, teacher.ID)

	return map[string]interface{}{
		"type": "bubble",
		"hero": map[string]interface{}{
			"type":        "image",
			"url":         "https://timeledger.example.com/images/welcome-teacher.png",
			"size":        "full",
			"aspectRatio": "20:13",
		},
		"body": map[string]interface{}{
			"type": "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "👋 歡迎加入 TimeLedger！",
					"weight": "bold",
					"size": "lg",
				},
				map[string]interface{}{
					"type": "text",
					"text": " ",
					"size": "sm",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("您的中心：%s", centerName),
					"size": "md",
					"color": "#666666",
				},
				map[string]interface{}{
					"type": "text",
					"text": " ",
					"size": "sm",
				},
				map[string]interface{}{
					"type": "text",
					"text": "點擊下方按鈕完成綁定，即可使用：",
					"size": "sm",
					"color": "#999999",
				},
				map[string]interface{}{
					"type": "text",
					"text": "✅ 查看課表",
					"weight": "bold",
					"color": "#4CAF50",
					"size": "sm",
				},
				map[string]interface{}{
					"type": "text",
					"text": "✅ 提交例外申請",
					"weight": "bold",
					"color": "#4CAF50",
					"size": "sm",
				},
				map[string]interface{}{
					"type": "text",
					"text": "✅ 管理私人行程",
					"weight": "bold",
					"color": "#4CAF50",
					"size": "sm",
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "立即綁定",
						"uri":   bindURL,
					},
				},
			},
		},
	}
}

// GetWelcomeAdminTemplate 管理員歡迎訊息範本
func (s *LineBotTemplateServiceImpl) GetWelcomeAdminTemplate(admin *models.AdminUser, centerName string) interface{} {
	bindURL := fmt.Sprintf("%s/admin/line-bind", s.baseURL)

	roleText := "中心管理員"
	if admin.Role == "OWNER" {
		roleText = "中心擁有者"
	}

	return map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type": "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "🎉 歡迎使用 TimeLedger！",
					"weight": "bold",
					"size": "lg",
				},
				map[string]interface{}{
					"type": "text",
					"text": " ",
					"size": "sm",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("您的中心：%s", centerName),
					"size": "md",
					"color": "#666666",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("您的角色：%s", roleText),
					"size": "md",
					"color": "#666666",
				},
				map[string]interface{}{
					"type": "separator",
					"margin": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": "🔔 及時通知功能",
					"weight": "bold",
					"margin": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": "綁定 LINE 帳號後，當老師提交例外申請時，\n您會立即收到通知！",
					"size": "sm",
					"color": "#999999",
					"wrap": true,
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box",
			"layout": "horizontal",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "立即綁定",
						"uri":   bindURL,
					},
				},
				map[string]interface{}{
					"type":   "button",
					"style":  "secondary",
					"action": map[string]interface{}{
						"type":  "message",
						"label": "稍後再說",
						"text":  "稍後綁定",
					},
				},
			},
		},
		"quickReply": map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{
					"type": "action",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "🔗 前往綁定",
						"uri":   bindURL,
					},
				},
				map[string]interface{}{
					"type": "action",
					"action": map[string]interface{}{
						"type":  "message",
						"label": "❓ 了解更多",
						"text":  "了解更多",
					},
				},
			},
		},
	}
}

// GetExceptionSubmitTemplate 例外申請通知範本（發給管理員）
func (s *LineBotTemplateServiceImpl) GetExceptionSubmitTemplate(exception *models.ScheduleException, teacherName string, centerName string) interface{} {
	adminURL := fmt.Sprintf("%s/admin/exceptions/%d", s.baseURL, exception.ID)

	// 根據類型顯示不同標題
	typeTitle := "例外申請"
	switch exception.ExceptionType {
	case "LEAVE":
		typeTitle = "請假申請"
	case "RESCHEDULE":
		typeTitle = "調課申請"
	case "SWAP":
		typeTitle = "代課申請"
	case "CANCEL":
		typeTitle = "取消課程"
	}

	return map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type": "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "🔔 新的" + typeTitle,
					"weight": "bold",
					"size": "lg",
				},
				map[string]interface{}{
					"type": "text",
					"text": "━━━━━━━━━━━━━━",
					"size": "xs",
					"color": "#CCCCCC",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("👤 申請人：%s 老師", teacherName),
					"size": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("📅 日期：%s", exception.GetDate().Format("2006/01/02 (Mon)")),
					"size": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("🕐 時間：%s", exception.GetTimeRange()),
					"size": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": "━━━━━━━━━━━━━━",
					"size": "xs",
					"color": "#CCCCCC",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("📝 原因：%s", exception.Reason),
					"size": "sm",
					"wrap": true,
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box",
			"layout": "horizontal",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "前往處理",
						"uri":   adminURL,
					},
				},
			},
		},
	}
}

// GetExceptionApproveTemplate 例外核准通知範本（發給老師）
func (s *LineBotTemplateServiceImpl) GetExceptionApproveTemplate(exception *models.ScheduleException, teacherName string) interface{} {
	teacherURL := fmt.Sprintf("%s/teacher/exceptions/%d", s.baseURL, exception.ID)

	// 根據類型顯示不同標題
	typeTitle := "例外申請"
	switch exception.ExceptionType {
	case "LEAVE":
		typeTitle = "請假申請"
	case "RESCHEDULE":
		typeTitle = "調課申請"
	case "SWAP":
		typeTitle = "代課申請"
	case "CANCEL":
		typeTitle = "取消課程"
	}

	return map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type": "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "✅ " + typeTitle + "已核准",
					"weight": "bold",
					"size": "lg",
					"color": "#4CAF50",
				},
				map[string]interface{}{
					"type": "text",
					"text": "━━━━━━━━━━━━━━",
					"size": "xs",
					"color": "#CCCCCC",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("📅 日期：%s", exception.GetDate().Format("2006/01/02 (Mon)")),
					"size": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("🕐 時間：%s", exception.GetTimeRange()),
					"size": "md",
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box",
			"layout": "horizontal",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "查看詳情",
						"uri":   teacherURL,
					},
				},
			},
		},
	}
}

// GetExceptionRejectTemplate 例外拒絕通知範本（發給老師）
func (s *LineBotTemplateServiceImpl) GetExceptionRejectTemplate(exception *models.ScheduleException, teacherName string, reason string) interface{} {
	teacherURL := fmt.Sprintf("%s/teacher/exceptions/%d", s.baseURL, exception.ID)

	// 根據類型顯示不同標題
	typeTitle := "例外申請"
	switch exception.ExceptionType {
	case "LEAVE":
		typeTitle = "請假申請"
	case "RESCHEDULE":
		typeTitle = "調課申請"
	case "SWAP":
		typeTitle = "代課申請"
	case "CANCEL":
		typeTitle = "取消課程"
	}

	rejectReason := reason
	if rejectReason == "" {
		rejectReason = "未說明原因"
	}

	return map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type": "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "❌ " + typeTitle + "已拒絕",
					"weight": "bold",
					"size": "lg",
					"color": "#F44336",
				},
				map[string]interface{}{
					"type": "text",
					"text": "━━━━━━━━━━━━━━",
					"size": "xs",
					"color": "#CCCCCC",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("📅 日期：%s", exception.GetDate().Format("2006/01/02 (Mon)")),
					"size": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("🕐 時間：%s", exception.GetTimeRange()),
					"size": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": "━━━━━━━━━━━━━━",
					"size": "xs",
					"color": "#CCCCCC",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("📝 拒絕原因：%s", rejectReason),
					"size": "sm",
					"wrap": true,
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box",
			"layout": "horizontal",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "secondary",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "查看詳情",
						"uri":   teacherURL,
					},
				},
			},
		},
	}
}

// GetInvitationAcceptedTemplate 邀請接受通知範本（發給管理員）
func (s *LineBotTemplateServiceImpl) GetInvitationAcceptedTemplate(teacher *models.Teacher, centerName string, role string) interface{} {
	adminURL := fmt.Sprintf("%s/admin/teachers", s.baseURL)

	// 角色顯示文字
	roleText := "老師"
	switch role {
	case "SUBSTITUTE":
		roleText = "代課老師"
	case "TEACHER":
		roleText = "正職老師"
	}

	return map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type": "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "🎉 新成員加入！",
					"weight": "bold",
					"size": "lg",
					"color": "#4CAF50",
				},
				map[string]interface{}{
					"type": "text",
					"text": "━━━━━━━━━━━━━━",
					"size": "xs",
					"color": "#CCCCCC",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("👤 新成員：%s", teacher.Name),
					"size": "md",
					"weight": "bold",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("🏢 中心：%s", centerName),
					"size": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("📋 角色：%s", roleText),
					"size": "md",
				},
				map[string]interface{}{
					"type": "text",
					"text": "━━━━━━━━━━━━━━",
					"size": "xs",
					"color": "#CCCCCC",
				},
				map[string]interface{}{
					"type": "text",
					"text": "✅ 歡迎新老師加入！",
					"size": "sm",
					"color": "#666666",
				},
			},
		},
		"footer": map[string]interface{}{
			"type": "box",
			"layout": "horizontal",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "查看成員",
						"uri":   adminURL,
					},
				},
			},
		},
	}
}
