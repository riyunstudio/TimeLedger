package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/database/mysql"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"
	mockRedis "timeLedger/testing/redis"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCenterRepository_CRUD(t *testing.T) {
	ctx := context.Background()

	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/timeledger_test?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL init error: %s", err.Error()))
	}

	rdb, mr, err := mockRedis.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Redis init error: %s", err.Error()))
	}
	defer mr.Close()

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	appInstance := &app.App{
		Env:   nil,
		Err:   e,
		Tools: tool,
		Mysql: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
		Redis: &redis.Redis{DB0: rdb},
		Api:   nil,
		Rpc:   nil,
	}

	repo := repositories.NewCenterRepository(appInstance)

	testCenterName := fmt.Sprintf("Test Center %d", time.Now().UnixNano())

	center := models.Center{
		Name:      testCenterName,
		PlanLevel: "STARTER",
		Settings: models.CenterSettings{
			AllowPublicRegister: true,
			DefaultLanguage:     "zh-TW",
		},
		CreatedAt: time.Now(),
	}

	createdCenter, err := repo.Create(ctx, center)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdCenter.ID == 0 {
		t.Fatal("Created center ID should not be 0")
	}

	fetchedCenter, err := repo.GetByID(ctx, createdCenter.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if fetchedCenter.Name != testCenterName {
		t.Errorf("Expected name %s, got %s", testCenterName, fetchedCenter.Name)
	}

	centers, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	found := false
	for _, c := range centers {
		if c.ID == createdCenter.ID && c.Name == testCenterName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Created center not found in list")
	}
}
