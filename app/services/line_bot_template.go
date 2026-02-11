package services

import (
	"fmt"
	"time"
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

	// 取得行程聚合範本
	GenerateAgendaFlex(agendaItems []AgendaItem, targetDate time.Time, userName string) interface{}

	// 取得廣播訊息範本
	GetBroadcastTemplate(centerName string, title string, message string, warning string, actionLabel string, actionURL string) interface{}
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
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "👋 歡迎加入 TimeLedger！",
					"weight": "bold",
					"size":   "lg",
				},
				map[string]interface{}{
					"type": "text",
					"text": " ",
					"size": "sm",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  fmt.Sprintf("您的中心：%s", centerName),
					"size":  "md",
					"color": "#666666",
				},
				map[string]interface{}{
					"type": "text",
					"text": " ",
					"size": "sm",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "點擊下方按鈕完成綁定，即可使用：",
					"size":  "sm",
					"color": "#999999",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "✅ 查看課表",
					"weight": "bold",
					"color":  "#4CAF50",
					"size":   "sm",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "✅ 提交例外申請",
					"weight": "bold",
					"color":  "#4CAF50",
					"size":   "sm",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "✅ 管理私人行程",
					"weight": "bold",
					"color":  "#4CAF50",
					"size":   "sm",
				},
			},
		},
		"footer": map[string]interface{}{
			"type":   "box",
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
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "🎉 歡迎使用 TimeLedger！",
					"weight": "bold",
					"size":   "lg",
				},
				map[string]interface{}{
					"type": "text",
					"text": " ",
					"size": "sm",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  fmt.Sprintf("您的中心：%s", centerName),
					"size":  "md",
					"color": "#666666",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  fmt.Sprintf("您的角色：%s", roleText),
					"size":  "md",
					"color": "#666666",
				},
				map[string]interface{}{
					"type":   "separator",
					"margin": "md",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   "🔔 及時通知功能",
					"weight": "bold",
					"margin": "md",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "綁定 LINE 帳號後，當老師提交例外申請時，\n您會立即收到通知！",
					"size":  "sm",
					"color": "#999999",
					"wrap":  true,
				},
			},
		},
		"footer": map[string]interface{}{
			"type":   "box",
			"layout": "horizontal",
			"contents": []interface{}{
				map[string]interface{}{
					"type":  "button",
					"style": "primary",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "立即綁定",
						"uri":   bindURL,
					},
				},
				map[string]interface{}{
					"type":  "button",
					"style": "secondary",
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
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "🔔 新的" + typeTitle,
					"weight": "bold",
					"size":   "lg",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "━━━━━━━━━━━━━━",
					"size":  "xs",
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
					"type":  "text",
					"text":  "━━━━━━━━━━━━━━",
					"size":  "xs",
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
			"type":   "box",
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
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "✅ " + typeTitle + "已核准",
					"weight": "bold",
					"size":   "lg",
					"color":  "#4CAF50",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "━━━━━━━━━━━━━━",
					"size":  "xs",
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
			"type":   "box",
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
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "❌ " + typeTitle + "已拒絕",
					"weight": "bold",
					"size":   "lg",
					"color":  "#F44336",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "━━━━━━━━━━━━━━",
					"size":  "xs",
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
					"type":  "text",
					"text":  "━━━━━━━━━━━━━━",
					"size":  "xs",
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
			"type":   "box",
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
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "text",
					"text":   "🎉 新成員加入！",
					"weight": "bold",
					"size":   "lg",
					"color":  "#4CAF50",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "━━━━━━━━━━━━━━",
					"size":  "xs",
					"color": "#CCCCCC",
				},
				map[string]interface{}{
					"type":   "text",
					"text":   fmt.Sprintf("👤 新成員：%s", teacher.Name),
					"size":   "md",
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
					"type":  "text",
					"text":  "━━━━━━━━━━━━━━",
					"size":  "xs",
					"color": "#CCCCCC",
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "✅ 歡迎新老師加入！",
					"size":  "sm",
					"color": "#666666",
				},
			},
		},
		"footer": map[string]interface{}{
			"type":   "box",
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

// GenerateAgendaFlex 行程聚合 Flex Message 範本
// 支援顯示多筆行程列表，中心課程使用藍色系，個人行程使用紫色系
func (s *LineBotTemplateServiceImpl) GenerateAgendaFlex(agendaItems []AgendaItem, targetDate time.Time, userName string) interface{} {
	// 格式化日期
	dateStr := targetDate.Format("2006年1月2日")
	weekdayStr := targetDate.Format("Mon")
	weekdayMap := map[string]string{
		"Monday":    "週一",
		"Tuesday":   "週二",
		"Wednesday": "週三",
		"Thursday":  "週四",
		"Friday":   "週五",
		"Saturday":  "週六",
		"Sunday":    "週日",
	}
	weekdayTW := weekdayMap[weekdayStr]
	if weekdayTW == "" {
		weekdayTW = weekdayStr
	}

	// 構建行程列表內容
	var agendaContents []interface{}

	// 日期標題
	agendaContents = append(agendaContents, map[string]interface{}{
		"type": "text",
		"text": fmt.Sprintf("📅 %s (%s)", dateStr, weekdayTW),
		"weight": "bold",
		"size":   "lg",
		"align":  "center",
	})

	// 如果沒有行程
	if len(agendaItems) == 0 {
		agendaContents = append(agendaContents, map[string]interface{}{
			"type":   "text",
			"text":   "🎉 今天沒有行程",
			"size":   "md",
			"color":  "#666666",
			"align":  "center",
			"margin": "md",
		})
	} else {
		// 分隔線
		agendaContents = append(agendaContents, map[string]interface{}{
			"type":  "separator",
			"margin": "md",
		})

		// 遍歷所有行程項目
		for _, item := range agendaItems {
			// 根據來源類型設定顏色
			var icon, color, bgColor string
			if item.SourceType == AgendaSourceTypeCenter {
				icon = "🏢"
				color = "#1E88E5" // 藍色系
				bgColor = "#E3F2FD"
			} else {
				icon = "📌"
				color = "#9C27B0" // 紫色系
				bgColor = "#F3E5F5"
			}

			// 行程項目
			itemBox := map[string]interface{}{
				"type": "box",
				"layout": "horizontal",
				"margin": "sm",
				"paddingAll": "8px",
				"backgroundColor": bgColor,
				"cornerRadius": "8px",
				"contents": []interface{}{
					// 時間
					map[string]interface{}{
						"type": "text",
						"text": item.Time,
						"size": "md",
						"weight": "bold",
						"color": color,
						"flex": 0,
						"align": "center",
						"minWidth": "60px",
					},
					// 分隔線
					map[string]interface{}{
						"type":  "separator",
						"color": color,
						"margin": "xs",
					},
					// 標題和來源
					map[string]interface{}{
						"type":   "box",
						"layout": "vertical",
						"flex":   1,
						"contents": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": item.Title,
								"size": "md",
								"weight": "bold",
								"color": "#333333",
								"wrap": true,
							},
							map[string]interface{}{
								"type": "text",
								"text": fmt.Sprintf("%s %s", icon, item.SourceName),
								"size": "xs",
								"color": "#888888",
								"margin": "xs",
							},
						},
					},
				},
			}
			agendaContents = append(agendaContents, itemBox)
		}
	}

	// 統計資訊
	if len(agendaItems) > 0 {
		agendaContents = append(agendaContents, map[string]interface{}{
			"type":  "separator",
			"margin": "md",
		})
		agendaContents = append(agendaContents, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("📊 共 %d 筆行程", len(agendaItems)),
			"size":  "sm",
			"color": "#999999",
			"align": "end",
			"margin": "sm",
		})
	}

	homeURL := fmt.Sprintf("%s", s.baseURL)

	return map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				// 用戶歡迎標題
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("👋 %s 的今日行程", userName),
					"weight": "bold",
					"size":   "xl",
					"align":  "center",
					"margin": "md",
				},
				// 分隔線
				map[string]interface{}{
					"type":  "separator",
					"margin": "md",
				},
				// 行程列表
				map[string]interface{}{
					"type":   "box",
					"layout": "vertical",
					"margin": "md",
					"contents": agendaContents,
				},
			},
		},
		"footer": map[string]interface{}{
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": "📱 進入系統首頁",
						"uri":   homeURL,
					},
				},
				map[string]interface{}{
					"type":  "text",
					"text":  "按鈕無法點擊？請直接複製連結",
					"size":  "xs",
					"color": "#AAAAAA",
					"align": "center",
					"margin": "sm",
				},
			},
		},
		"styles": map[string]interface{}{
			"footer": map[string]interface{}{
				"separator": true,
			},
		},
	}
}

// GetBroadcastTemplate 廣播訊息 Flex Message 範本
func (s *LineBotTemplateServiceImpl) GetBroadcastTemplate(centerName string, title string, message string, warning string, actionLabel string, actionURL string) interface{} {
	// 構建內容列表
	contents := []interface{}{
		// 標題
		map[string]interface{}{
			"type":   "text",
			"text":   title,
			"weight": "bold",
			"size":   "lg",
		},
		// 分隔線
		map[string]interface{}{
			"type":  "text",
			"text":  "━━━━━━━━━━━━━━━━",
			"size":  "xs",
			"color": "#CCCCCC",
		},
	}

	// 添加中心名稱
	contents = append(contents, map[string]interface{}{
		"type":  "text",
		"text":  fmt.Sprintf("🏢 來自：%s", centerName),
		"size":  "md",
		"color": "#666666",
	})

	// 添加訊息內容
	contents = append(contents, map[string]interface{}{
		"type":  "text",
		"text":  "━━━━━━━━━━━━━━━━",
		"size":  "xs",
		"color": "#CCCCCC",
	})
	contents = append(contents, map[string]interface{}{
		"type":  "text",
		"text":  message,
		"size":  "md",
		"wrap":  true,
		"margin": "md",
	})

	// 如果有警告訊息，添加警告區塊
	if warning != "" {
		contents = append(contents, []interface{}{
			map[string]interface{}{
				"type":  "separator",
				"margin": "md",
			},
			map[string]interface{}{
				"type":   "box",
				"layout": "vertical",
				"margin": "md",
				"paddingAll": "12px",
				"backgroundColor": "#FFF8E1",
				"cornerRadius": "8px",
				"contents": []interface{}{
					map[string]interface{}{
						"type":  "text",
						"text":  fmt.Sprintf("⚠️ %s", warning),
						"size":  "sm",
						"color": "#F57C00",
						"wrap":  true,
					},
				},
			},
		}...)
	}

	// 構建 Flex Message
	flexMessage := map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type":   "box",
			"layout": "vertical",
			"contents": contents,
		},
	}

	// 如果有動作按鈕，添加 footer
	if actionLabel != "" && actionURL != "" {
		flexMessage["footer"] = map[string]interface{}{
			"type":   "box",
			"layout": "vertical",
			"contents": []interface{}{
				map[string]interface{}{
					"type":   "button",
					"style":  "primary",
					"height": "sm",
					"action": map[string]interface{}{
						"type":  "uri",
						"label": actionLabel,
						"uri":   actionURL,
					},
				},
			},
		}
	}

	return flexMessage
}
