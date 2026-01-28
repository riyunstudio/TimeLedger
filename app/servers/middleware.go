package servers

import (
	"net/http"
	"timeLedger/global/errInfos"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

func (s *Server) InitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.setRequestTime(c)
		s.wg.Add(1)
		defer func() {
			s.wg.Done()
		}()
		c.Next()
	}
}

func (s *Server) RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				e := s.app.Tools.PanicParser(err)
				// 紀錄 TraceLog
				s.writePanicLog(c, e)
				// 回傳統一錯誤給 client
				c.JSON(http.StatusInternalServerError, gin.H{"code": errInfos.SYSTEM_ERROR})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RateLimitMiddleware 速率限制中介層
func (s *Server) RateLimitMiddleware() gin.HandlerFunc {
	rateLimiter := services.NewRateLimiterService(s.app)
	return services.RateLimitMiddleware(rateLimiter)
}

// 主要中介層
func (s *Server) MainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Post/Put 需要先把body備份寫入ctx
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			// 將body內容複製
			s.setBodyParams(c)
		}

		// 設置 traceID
		s.setTraceID(c)

		// 執行主程式 handler
		c.Next()

		// 紀錄 TraceLog
		s.writeApiLog(c)
	}
}

func (s *Server) IpMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
