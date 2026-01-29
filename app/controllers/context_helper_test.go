package controllers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

func TestContextHelper_UserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupCtx  func(*gin.Context)
		expectID  uint
		expectErr bool
	}{
		{
			name: "Valid user ID",
			setupCtx: func(c *gin.Context) {
				c.Set(global.UserIDKey, uint(123))
			},
			expectID:  123,
			expectErr: false,
		},
		{
			name: "Zero user ID",
			setupCtx: func(c *gin.Context) {
				c.Set(global.UserIDKey, uint(0))
			},
			expectID:  0,
			expectErr: true,
		},
		{
			name: "Missing user ID",
			setupCtx: func(c *gin.Context) {
				// 不設置任何值
			},
			expectID:  0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setupCtx(c)

			helper := NewContextHelper(c)
			uid, err := helper.UserID()

			if tt.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if uid != tt.expectID {
				t.Errorf("Expected user ID %d, got %d", tt.expectID, uid)
			}
		})
	}
}

func TestContextHelper_CenterID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupCtx  func(*gin.Context)
		expectID  uint
		expectErr bool
	}{
		{
			name: "Valid center ID",
			setupCtx: func(c *gin.Context) {
				c.Set(global.CenterIDKey, uint(456))
			},
			expectID:  456,
			expectErr: false,
		},
		{
			name: "Zero center ID",
			setupCtx: func(c *gin.Context) {
				c.Set(global.CenterIDKey, uint(0))
			},
			expectID:  0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setupCtx(c)

			helper := NewContextHelper(c)
			centerID, err := helper.CenterID()

			if tt.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if centerID != tt.expectID {
				t.Errorf("Expected center ID %d, got %d", tt.expectID, centerID)
			}
		})
	}
}

func TestContextHelper_UserType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		setupCtx   func(*gin.Context)
		expectType string
		expectOK   bool
	}{
		{
			name: "Admin user type",
			setupCtx: func(c *gin.Context) {
				c.Set(global.UserTypeKey, "ADMIN")
			},
			expectType: "ADMIN",
			expectOK:   true,
		},
		{
			name: "Teacher user type",
			setupCtx: func(c *gin.Context) {
				c.Set(global.UserTypeKey, "TEACHER")
			},
			expectType: "TEACHER",
			expectOK:   true,
		},
		{
			name: "Missing user type",
			setupCtx: func(c *gin.Context) {
				// 不設置任何值
			},
			expectType: "",
			expectOK:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setupCtx(c)

			helper := NewContextHelper(c)
			userType, ok := helper.UserType()

			if ok != tt.expectOK {
				t.Errorf("Expected ok=%v, got %v", tt.expectOK, ok)
			}
			if userType != tt.expectType {
				t.Errorf("Expected user type %s, got %s", tt.expectType, userType)
			}
		})
	}
}

func TestContextHelper_IsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		userType string
		expected bool
	}{
		{"Admin user", "ADMIN", true},
		{"Owner user", "OWNER", true},
		{"Teacher user", "TEACHER", false},
		{"Unknown user", "UNKNOWN", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set(global.UserTypeKey, tt.userType)

			helper := NewContextHelper(c)
			if helper.IsAdmin() != tt.expected {
				t.Errorf("Expected IsAdmin=%v for user type %s", tt.expected, tt.userType)
			}
		})
	}
}

func TestContextHelper_QueryString(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test?key=value", nil)

	helper := NewContextHelper(c)
	val, ok := helper.QueryString("key")

	if !ok {
		t.Error("Expected ok=true for existing key")
	}
	if val != "value" {
		t.Errorf("Expected 'value', got '%s'", val)
	}

	// 測試不存在的 key
	_, ok = helper.QueryString("missing")
	if ok {
		t.Error("Expected ok=false for missing key")
	}
}

func TestContextHelper_QueryStringOrDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test?key=value", nil)

	helper := NewContextHelper(c)

	// 存在的 key
	val := helper.QueryStringOrDefault("key", "default")
	if val != "value" {
		t.Errorf("Expected 'value', got '%s'", val)
	}

	// 不存在的 key
	val = helper.QueryStringOrDefault("missing", "default")
	if val != "default" {
		t.Errorf("Expected 'default', got '%s'", val)
	}
}

func TestContextHelper_QueryDate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		queryKey   string
		queryValue string
		expectErr  bool
	}{
		{
			name:       "Valid date",
			queryKey:   "date",
			queryValue: "2026-01-20",
			expectErr:  false,
		},
		{
			name:       "Missing date",
			queryKey:   "date",
			queryValue: "",
			expectErr:  true,
		},
		{
			name:       "Invalid date format",
			queryKey:   "date",
			queryValue: "20/01/2026",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			url := "/test"
			if tt.queryValue != "" {
				url += "?date=" + tt.queryValue
			}
			c.Request, _ = http.NewRequest("GET", url, nil)

			helper := NewContextHelper(c)
			_, err := helper.QueryDate("date")

			if tt.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestContextHelper_ParamUint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		paramValue string
		expectID   uint
		expectErr  bool
	}{
		{
			name:       "Valid uint",
			paramValue: "123",
			expectID:   123,
			expectErr:  false,
		},
		{
			name:       "Zero value",
			paramValue: "0",
			expectID:   0,
			expectErr:  true,
		},
		{
			name:       "Invalid format",
			paramValue: "abc",
			expectID:   0,
			expectErr:  true,
		},
		{
			name:       "Empty value",
			paramValue: "",
			expectID:   0,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}

			helper := NewContextHelper(c)
			id, err := helper.ParamUint("id")

			if tt.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if id != tt.expectID {
				t.Errorf("Expected ID %d, got %d", tt.expectID, id)
			}
		})
	}
}

func TestContextHelper_BindJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type TestRequest struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name      string
		body      string
		expectErr bool
	}{
		{
			name:      "Valid JSON",
			body:      `{"name":"test","age":25}`,
			expectErr: false,
		},
		{
			name:      "Invalid JSON",
			body:      `{invalid}`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/test", nil)
			c.Request.Body = mockBody(tt.body)

			helper := NewContextHelper(c)
			var req TestRequest
			err := helper.BindJSON(&req)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestContextHelper_ResponseMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helper := NewContextHelper(c)
		helper.Success(map[string]string{"key": "value"})

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Created response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helper := NewContextHelper(c)
		helper.Created(map[string]string{"id": "123"})

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("BadRequest response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helper := NewContextHelper(c)
		helper.BadRequest("test error")

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Unauthorized response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helper := NewContextHelper(c)
		helper.Unauthorized("auth required")

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("NotFound response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helper := NewContextHelper(c)
		helper.NotFound("resource not found")

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

// mockBody 創建模擬請求體
func mockBody(content string) io.ReadCloser {
	return &mockReadCloser{content: content}
}

type mockReadCloser struct {
	content string
	pos     int
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.content) {
		return 0, io.EOF
	}
	n = copy(p, m.content[m.pos:])
	m.pos += n
	return n, nil
}

func (m *mockReadCloser) Close() error {
	return nil
}
