package services

import (
	"akali/global"
	"context"

	"github.com/gin-gonic/gin"
)

type BaseService struct{}

// 把Gin的ctx取指定內容寫入context ctx, 只有Gin Api操作SQL時會需要
func (s *BaseService) dbCtx(ctx *gin.Context) (dbCtx context.Context) {
	if val, exists := ctx.Get(global.TraceIDKey); exists {
		if tid, ok := val.(string); ok {
			dbCtx = context.WithValue(context.Background(), global.TraceIDKey, tid)
			return
		}
	}
	dbCtx = context.WithValue(context.Background(), global.TraceIDKey, "")
	return
}
