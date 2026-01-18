package test

import (
	"fmt"
	"reflect"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/requests"
	"timeLedger/app/resources"
	"timeLedger/app/services"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	mockRedis "timeLedger/testing/redis"
	"timeLedger/testing/sqlite"

	_ "github.com/joho/godotenv/autoload"

	"testing"
)

func TestUserService_CreateAndGet(t *testing.T) {
	ctx := t.Context()

	// 初始化 Mock SQLite
	sqliteDB, err := sqlite.Initialize()
	if err != nil {
		panic(fmt.Sprintf("[流程測試錯誤] Err: %s", err.Error()))
	}

	// 建立資料表
	if err := sqliteDB.AutoMigrate([]any{
		&models.User{},
	}...); err != nil {
		panic(fmt.Sprintf("[流程測試錯誤] Err: %s", err.Error()))
	}

	// 初始化 Mock Redis
	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("[流程測試錯誤] Err: %s", err.Error()))
	}
	defer mr.Close()

	// 初始化 App
	app := app.Initialize()
	app.Mysql = &mysql.DB{WDB: sqliteDB, RDB: sqliteDB}
	app.Redis = &redis.Redis{DB0: rdb}
	s := services.NewUserService(app)

	// 自定義測試範例
	tests := []struct {
		name           string
		req            *requests.UserCreateRequest
		expectError    bool
		expectErrCode  int
		expectUserData models.User
	}{
		{
			name: "正常建立用戶",
			req: &requests.UserCreateRequest{
				Name: "Alice",
				Ips:  []string{"127.0.0.1"},
			},
			expectError: false,
			expectUserData: models.User{
				Name: "Alice",
			},
		},
		{
			name: "正常建立另一個用戶",
			req: &requests.UserCreateRequest{
				Name: "Bob",
				Ips:  []string{"192.168.0.1", "10.0.0.1"},
			},
			expectError: false,
			expectUserData: models.User{
				Name: "Bob",
			},
		},
		// {
		// 	name: "建立用戶失敗範例",
		// 	req: &requests.UserCreateRequest{
		// 		Name: "", // 假設 Name 不可為空
		// 		Ips:  []string{},
		// 	},
		// 	expectError:   true,
		// 	expectErrCode: 400, // 假設服務會回傳 400
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, eInfo, err := s.Create(ctx, tt.req)
			if tt.expectError {
				if err == nil {
					t.Fatalf("預期錯誤，但實際沒有錯誤")
				}
				if eInfo.Code != tt.expectErrCode {
					t.Fatalf("預期錯誤碼 %d, 但得到 %d", tt.expectErrCode, eInfo.Code)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			newUser := data.(*resources.UserCreateResource)

			if newUser.Name != tt.expectUserData.Name {
				t.Errorf("預期 Name %s, 但得到 %s", tt.expectUserData.Name, newUser.Name)
			}
			// 可以驗證其他欄位，例如 Ips 或其他自動欄位
			if !reflect.DeepEqual(newUser.Ips, tt.req.Ips) {
				t.Errorf("預期 Ips %v, 但得到 %v", tt.req.Ips, newUser.Ips)
			}

			data, eInfo, err = s.Get(ctx, &requests.UserGetRequest{ID: newUser.ID})
			if tt.expectError {
				if err == nil {
					t.Fatalf("預期錯誤，但實際沒有錯誤")
				}
				if eInfo.Code != tt.expectErrCode {
					t.Fatalf("預期錯誤碼 %d, 但得到 %d", tt.expectErrCode, eInfo.Code)
				}
				return
			}
			if err != nil {
				t.Fatalf("Get error: %v", err)
			}
			user := data.(*resources.UserGetResource)
			if user.Name != newUser.Name {
				t.Errorf("Get Name 不符合: want %s, got %s", newUser.Name, user.Name)
			}
			if !reflect.DeepEqual(user.Ips, newUser.Ips) {
				t.Errorf("Get Ips 不符合: want %v, got %v", newUser.Ips, user.Ips)
			}
		})
	}
}
