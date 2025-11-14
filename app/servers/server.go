package servers

import (
	"akali/app"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	app    *app.App
	engine *gin.Engine
	routes []route
	action actions

	srv *http.Server   // Shutdown時直接使用
	wg  sync.WaitGroup // For 優雅退出
}

func Initialize(app *app.App) *Server {
	r := gin.New()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.Discard

	r.Use(gin.Recovery())

	return &Server{
		app: app, engine: r,
	}
}

func (s *Server) Start() {
	// 初始化控制器
	s.NewControllers()

	// 載入所有路由
	s.LoadRoutes()

	// 註冊所有路由
	s.registerRoutes(s.engine)

	s.srv = &http.Server{
		Addr:    ":" + s.app.Env.ServerPort,
		Handler: s.engine,
	}
	go func() {
		log.Println("Gin server started")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("gin server start error, Err: %v", err))
		}
	}()
}

func (s *Server) Stop() {
	log.Println("Gin server stopping...")

	// 1. 使用 Shutdown 停止 HTTP Server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Printf("Gin server shutdown error: %v", err)
	}

	// 2. 等待所有請求完成
	s.wg.Wait()
	log.Println("Gin server graceful stop complete")
}
