package servers

import (
	"akali/global"
	"akali/global/errInfos"
	"context"
	"net/http"
	"time"

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

// func (s *Server) TimeoutMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// 建立有超時的 context
// 		ctx, cancel := context.WithTimeout(c.Request.Context(), s.timeout)
// 		defer cancel()

// 		// 用 channel 來控制 handler 是否完成
// 		done := make(chan struct{})

// 		c.Request = c.Request.WithContext(ctx) // 替換原始 context

// 		go func() {
// 			c.Next()
// 			close(done)
// 		}()

// 		select {
// 		case <-done:
// 			return
// 		case <-ctx.Done():
// 			// 超時或 client 取消
// 			eInfo := errInfos.New(errInfos.REQUEST_TIMEOUT)
// 			c.JSON(http.StatusGatewayTimeout, global.ApiResponse{
// 				Code:    eInfo.Code,
// 				Message: eInfo.Msg,
// 			})
// 			c.Abort()
// 			return
// 		}
// 	}
// }

func (s *Server) TimeoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), s.timeout)
		defer cancel()

		done := make(chan struct{})

		// 替換原本的 request context
		c.Request = c.Request.WithContext(ctx)

		go func() {
			defer func() {
				close(done)
			}()
			c.Next()
		}()

		select {
		case <-done:
			// handler 執行完成
			return
		case <-time.After(s.timeout):
			// 超時回傳錯誤 (不等待)
			eInfo := errInfos.New(errInfos.REQUEST_TIMEOUT)
			c.JSON(http.StatusGatewayTimeout, global.ApiResponse{
				Code:    eInfo.Code,
				Message: eInfo.Msg,
			})
			c.Abort()
			return
		case <-ctx.Done():
			// client 取消
			eInfo := errInfos.New(errInfos.REQUEST_TIMEOUT)
			c.JSON(http.StatusGatewayTimeout, global.ApiResponse{
				Code:    eInfo.Code,
				Message: eInfo.Msg,
			})
			c.Abort()
			return
		}
	}
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
