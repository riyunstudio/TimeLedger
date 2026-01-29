package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"timeLedger/app"
	"timeLedger/app/services"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	app  *app.App
	auth services.AuthService
}

func NewAuthMiddleware(app *app.App, auth services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		app:  app,
		auth: auth,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, global.ApiResponse{
				Code:    global.UNAUTHORIZED,
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, global.ApiResponse{
				Code:    global.UNAUTHORIZED,
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := m.auth.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, global.ApiResponse{
				Code:    global.UNAUTHORIZED,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(global.UserIDKey, claims.UserID)
		c.Set(global.UserTypeKey, claims.UserType)
		c.Set(global.CenterIDKey, uint(claims.CenterID))
		c.Set(global.LineUserIDKey, claims.LineUserID)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get(global.UserTypeKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, global.ApiResponse{
				Code:    global.UNAUTHORIZED,
				Message: "Authentication required",
			})
			c.Abort()
			return
		}

		if userType != "ADMIN" && userType != "OWNER" {
			c.JSON(http.StatusForbidden, global.ApiResponse{
				Code:    global.FORBIDDEN,
				Message: "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RequireTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get(global.UserTypeKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, global.ApiResponse{
				Code:    global.UNAUTHORIZED,
				Message: "Authentication required",
			})
			c.Abort()
			return
		}

		if userType != "TEACHER" {
			c.JSON(http.StatusForbidden, global.ApiResponse{
				Code:    global.FORBIDDEN,
				Message: "Teacher access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RequireCenterAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get(global.UserTypeKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, global.ApiResponse{
				Code:    global.UNAUTHORIZED,
				Message: "Authentication required",
			})
			c.Abort()
			return
		}

		if userType != "ADMIN" && userType != "OWNER" {
			c.JSON(http.StatusForbidden, global.ApiResponse{
				Code:    global.FORBIDDEN,
				Message: "Admin access required",
			})
			c.Abort()
			return
		}

		centerIDFromToken, exists := c.Get(global.CenterIDKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, global.ApiResponse{
				Code:    global.UNAUTHORIZED,
				Message: "Center ID required",
			})
			c.Abort()
			return
		}

		centerIDTokenUint := centerIDFromToken.(uint)

		// If center_id from token is valid (non-zero), use it directly
		if centerIDTokenUint != 0 {
			c.Set(global.CenterIDKey, centerIDTokenUint)
			c.Next()
			return
		}

		// If center_id from token is 0 (super admin), require param
		centerIDParam := c.Param("id")
		if centerIDParam == "" {
			centerIDParam = c.Query("center_id")
		}
		if centerIDParam == "" {
			body, err := c.GetRawData()
			if err == nil && len(body) > 0 {
				var jsonBody map[string]interface{}
				if err := json.Unmarshal(body, &jsonBody); err == nil {
					if cid, ok := jsonBody["center_id"]; ok {
						switch v := cid.(type) {
						case float64:
							centerIDParam = fmt.Sprintf("%d", int(v))
						case int:
							centerIDParam = fmt.Sprintf("%d", v)
						case string:
							centerIDParam = v
						}
					}
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		if centerIDParam == "" {
			c.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Center ID required",
			})
			c.Abort()
			return
		}

		var centerIDParamUint uint
		if _, err := fmt.Sscanf(centerIDParam, "%d", &centerIDParamUint); err != nil {
			c.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid center ID format",
			})
			c.Abort()
			return
		}

		c.Set(global.CenterIDKey, centerIDParamUint)
		c.Next()
	}
}
