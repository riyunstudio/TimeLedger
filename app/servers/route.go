package servers

import (
	"akali/app/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type route struct {
	Method      string
	Path        string
	Controller  func(*gin.Context)
	Middlewares []gin.HandlerFunc
}

// 載入路由
func (s *Server) LoadRoutes() {
	s.routes = []route{
		// 使用者相關
		{http.MethodGet, "/user", s.action.user.Get, []gin.HandlerFunc{}},
		{http.MethodPost, "/user", s.action.user.Create, []gin.HandlerFunc{}},
		{http.MethodPut, "/user", s.action.user.Update, []gin.HandlerFunc{}},
	}
}

// 註冊路由
func (s *Server) registerRoutes(r *gin.Engine) {
	// pod 健康檢查
	r.GET("/healthy", func(c *gin.Context) {
		c.String(http.StatusOK, "Healthy")
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 全局 Middleware
	r.Use(s.InitMiddleware())
	r.Use(s.RecoverMiddleware())
	r.Use(s.MainMiddleware())

	for _, rt := range s.routes {
		switch rt.Method {
		case http.MethodGet:
			r.GET(rt.Path, append(rt.Middlewares, rt.Controller)...)
		case http.MethodPost:
			r.POST(rt.Path, append(rt.Middlewares, rt.Controller)...)
		case http.MethodPut:
			r.PUT(rt.Path, append(rt.Middlewares, rt.Controller)...)
		case http.MethodDelete:
			r.DELETE(rt.Path, append(rt.Middlewares, rt.Controller)...)
		}
	}
}

type actions struct {
	user *controllers.UserController
}

// 建構控制器
func (s *Server) NewControllers() {
	controllers.NewBaseController(s.app)
	s.action.user = controllers.NewUserController(s.app)
}
