package servers

import (
	"akali/app"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	app     *app.App
	engine  *gin.Engine
	routes  []route
	action  actions
	timeout time.Duration

	srv *http.Server   // Shutdown時直接使用
	wg  sync.WaitGroup // For 優雅退出
}

func Initialize(app *app.App) *Server {
	r := gin.New()
	r.Use(gin.Recovery())

	return &Server{
		app: app, engine: r,
		timeout: 5 * time.Second,
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
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("Server start error, Err: %v", err))
		}
	}()
}

func (s *Server) Stop() {
	log.Println("Stopping HTTP server...")

	// 1. 使用 Shutdown 停止 HTTP Server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown error: %v", err)
	}

	// 2. 等待所有請求完成
	s.wg.Wait()
	log.Println("Graceful Stop complete")
}
