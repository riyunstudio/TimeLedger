package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
	"timeLedger/configs"
	"timeLedger/database/mysql"
	"timeLedger/global/errInfos"

	"gitlab.en.mcbwvx.com/frame/teemo/tools"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupTestApp 建立測試用的 App 實例
func setupTeacherProfileTestApp(t *testing.T) *app.App {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("MySQL init error: %s. Skipping test.", err.Error())
		return nil
	}

	// 檢查資料庫連線
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		t.Skipf("MySQL DB error: %s. Skipping test.", err.Error())
		return nil
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("MySQL ping error: %s. Skipping test.", err.Error())
		return nil
	}

	e := errInfos.Initialize(1)
	tool := tools.Initialize("Asia/Taipei")

	env := &configs.Env{
		JWTSecret:   "test-jwt-secret-key-for-testing-only",
		AppEnv:      "test",
		AppDebug:    true,
		AppTimezone: "Asia/Taipei",
	}

	appInstance := &app.App{
		Env:   env,
		Err:   e,
		Tools: tool,
		MySQL: &mysql.DB{WDB: mysqlDB, RDB: mysqlDB},
		Redis: nil,
		Api:   nil,
		Rpc:   nil,
	}

	return appInstance
}

// TestTeacherProfileService_GetProfile 測試取得老師個人資料
func TestTeacherProfileService_GetProfile(t *testing.T) {
	t.Run("GetProfile_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			// 清理測試資料
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestProfile%").Delete(&models.Teacher{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)

		// 建立測試資料
		teacher := models.Teacher{
			Name:              "TestProfile_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:             fmt.Sprintf("test%d@test.com", time.Now().UnixNano()),
			LineUserID:        "test-line-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Bio:               "Test bio for profile",
			City:              "台北市",
			District:          "大安區",
			PublicContactInfo: "0912345678",
			IsOpenToHiring:    true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 執行測試
		svc := services.NewTeacherProfileService(appInstance)
		profile, eInfo, err := svc.GetProfile(ctx, teacher.ID)

		// 驗證結果
		if err != nil {
			t.Fatalf("GetProfile 失敗: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("GetProfile 返回錯誤: %s", eInfo.Msg)
		}
		if profile == nil {
			t.Fatal("GetProfile 返回空結果")
		}

		// 驗證欄位
		if profile.ID != teacher.ID {
			t.Errorf("ID 不匹配: 預期 %d, 實際 %d", teacher.ID, profile.ID)
		}
		if profile.Name != teacher.Name {
			t.Errorf("Name 不匹配: 預期 %s, 實際 %s", teacher.Name, profile.Name)
		}
		if profile.Email != teacher.Email {
			t.Errorf("Email 不匹配: 預期 %s, 實際 %s", teacher.Email, profile.Email)
		}
		if profile.Bio != teacher.Bio {
			t.Errorf("Bio 不匹配: 預期 %s, 實際 %s", teacher.Bio, profile.Bio)
		}
		if profile.City != teacher.City {
			t.Errorf("City 不匹配: 預期 %s, 實際 %s", teacher.City, profile.City)
		}
		if profile.District != teacher.District {
			t.Errorf("District 不匹配: 預期 %s, 實際 %s", teacher.District, profile.District)
		}
		if profile.IsOpenToHiring != teacher.IsOpenToHiring {
			t.Errorf("IsOpenToHiring 不匹配: 預期 %v, 實際 %v", teacher.IsOpenToHiring, profile.IsOpenToHiring)
		}
	})

	t.Run("GetProfile_TeacherNotFound", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}

		ctx := context.Background()
		svc := services.NewTeacherProfileService(appInstance)

		// 執行測試（使用不存在的 ID）
		profile, eInfo, err := svc.GetProfile(ctx, 99999999)

		// 驗證結果
		if err == nil {
			t.Fatal("預期錯誤但返回 nil")
		}
		if eInfo == nil {
			t.Fatal("預期錯誤資訊但返回 nil")
		}
		if profile != nil {
			t.Fatal("預期空結果但返回非空")
		}
	})
}

// TestTeacherProfileService_UpdateProfile 測試更新老師個人資料
func TestTeacherProfileService_UpdateProfile(t *testing.T) {
	t.Run("UpdateProfile_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestUpdateProfile%").Delete(&models.Teacher{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)

		// 建立測試資料
		teacher := models.Teacher{
			Name:              "TestUpdateProfile_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:             fmt.Sprintf("update%d@test.com", time.Now().UnixNano()),
			LineUserID:        "test-line-update-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Bio:               "Original bio",
			City:              "台北市",
			District:          "信義區",
			PublicContactInfo: "0911111111",
			IsOpenToHiring:    false,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)

		// 更新資料
		req := &services.UpdateProfileRequest{
			Bio:               "Updated bio",
			City:              "新北市",
			District:          "板橋區",
			PublicContactInfo: "0922222222",
			IsOpenToHiring:    true,
		}
		profile, eInfo, err := svc.UpdateProfile(ctx, teacher.ID, req)

		// 驗證結果
		if err != nil {
			t.Fatalf("UpdateProfile 失敗: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("UpdateProfile 返回錯誤: %s", eInfo.Msg)
		}
		if profile == nil {
			t.Fatal("UpdateProfile 返回空結果")
		}

		// 驗證更新後的欄位
		if profile.Bio != req.Bio {
			t.Errorf("Bio 不匹配: 預期 %s, 實際 %s", req.Bio, profile.Bio)
		}
		if profile.City != req.City {
			t.Errorf("City 不匹配: 預期 %s, 實際 %s", req.City, profile.City)
		}
		if profile.District != req.District {
			t.Errorf("District 不匹配: 預期 %s, 實際 %s", req.District, profile.District)
		}
		if profile.IsOpenToHiring != req.IsOpenToHiring {
			t.Errorf("IsOpenToHiring 不匹配: 預期 %v, 實際 %v", req.IsOpenToHiring, profile.IsOpenToHiring)
		}

		// 驗證資料庫中的實際資料
		updatedTeacher, _ := teacherRepo.GetByID(ctx, teacher.ID)
		if updatedTeacher.Bio != req.Bio {
			t.Errorf("資料庫 Bio 不匹配: 預期 %s, 實際 %s", req.Bio, updatedTeacher.Bio)
		}
	})

	t.Run("UpdateProfile_PartialUpdate", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestPartialUpdate%").Delete(&models.Teacher{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)

		// 建立測試資料
		teacher := models.Teacher{
			Name:              "TestPartialUpdate_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:             fmt.Sprintf("partial%d@test.com", time.Now().UnixNano()),
			LineUserID:        "test-line-partial-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Bio:               "Original bio",
			City:              "台北市",
			District:          "中山區",
			PublicContactInfo: "0933333333",
			IsOpenToHiring:    false,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)

		// 只更新 Bio，其他欄位保留
		req := &services.UpdateProfileRequest{
			Bio:            "Only update bio",
			City:           "", // 空的應該跳過更新
			District:       "",
			IsOpenToHiring: false,
		}
		profile, _, err := svc.UpdateProfile(ctx, teacher.ID, req)

		if err != nil {
			t.Fatalf("UpdateProfile 失敗: %v", err)
		}

		// 驗證 Bio 更新
		if profile.Bio != req.Bio {
			t.Errorf("Bio 不匹配: 預期 %s, 實際 %s", req.Bio, profile.Bio)
		}
		// 驗證 City 保持不變
		if profile.City != teacher.City {
			t.Errorf("City 不應該改變: 預期 %s, 實際 %s", teacher.City, profile.City)
		}
	})

	t.Run("UpdateProfile_TeacherNotFound", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}

		ctx := context.Background()
		svc := services.NewTeacherProfileService(appInstance)

		req := &services.UpdateProfileRequest{
			Bio: "Test bio",
		}
		_, eInfo, err := svc.UpdateProfile(ctx, 99999999, req)

		if err == nil {
			t.Fatal("預期錯誤但返回 nil")
		}
		if eInfo == nil {
			t.Fatal("預期錯誤資訊但返回 nil")
		}
	})
}

// TestTeacherProfileService_SkillCRUD 測試技能 CRUD 操作
func TestTeacherProfileService_SkillCRUD(t *testing.T) {
	t.Run("CreateSkill_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestSkillCRUD%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestSkillCRUD%").Delete(&models.TeacherSkill{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)

		// 建立測試老師
		teacher := models.Teacher{
			Name:       "TestSkillCRUD_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("skill%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-skill-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)

		// 建立技能
		req := &services.CreateSkillRequest{
			Category:   "運動",
			SkillName:  "瑜珈",
			Level:      "中級",
			HashtagIDs: nil,
		}
		skill, eInfo, err := svc.CreateSkill(ctx, teacher.ID, req)

		if err != nil {
			t.Fatalf("CreateSkill 失敗: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("CreateSkill 返回錯誤: %s", eInfo.Msg)
		}
		if skill == nil {
			t.Fatal("CreateSkill 返回空結果")
		}

		// 驗證欄位
		if skill.TeacherID != teacher.ID {
			t.Errorf("TeacherID 不匹配: 預期 %d, 實際 %d", teacher.ID, skill.TeacherID)
		}
		if skill.Category != req.Category {
			t.Errorf("Category 不匹配: 預期 %s, 實際 %s", req.Category, skill.Category)
		}
		if skill.SkillName != req.SkillName {
			t.Errorf("SkillName 不匹配: 預期 %s, 實際 %s", req.SkillName, skill.SkillName)
		}
		if skill.Level != req.Level {
			t.Errorf("Level 不匹配: 預期 %s, 實際 %s", req.Level, skill.Level)
		}
	})

	t.Run("GetSkills_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestGetSkills%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestGetSkills%").Delete(&models.TeacherSkill{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		skillRepo := repositories.NewTeacherSkillRepository(appInstance)

		// 建立測試老師
		teacher := models.Teacher{
			Name:       "TestGetSkills_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("getskills%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-getskills-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立多個技能
		for i := 0; i < 3; i++ {
			skillRepo.Create(ctx, models.TeacherSkill{
				TeacherID: teacher.ID,
				Category:  "運動",
				SkillName: fmt.Sprintf("技能%d", i+1),
				Level:     "初一級",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}

		svc := services.NewTeacherProfileService(appInstance)
		skills, eInfo, err := svc.GetSkills(ctx, teacher.ID)

		if err != nil {
			t.Fatalf("GetSkills 失敗: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("GetSkills 返回錯誤: %s", eInfo.Msg)
		}
		if len(skills) != 3 {
			t.Errorf("預期 3 個技能，實際 %d 個", len(skills))
		}
	})

	t.Run("UpdateSkill_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestUpdateSkill%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestUpdateSkill%").Delete(&models.TeacherSkill{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		skillRepo := repositories.NewTeacherSkillRepository(appInstance)

		// 建立測試老師和技能
		teacher := models.Teacher{
			Name:       "TestUpdateSkill_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("updateskill%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-updateskill-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		skill := models.TeacherSkill{
			TeacherID: teacher.ID,
			Category:  "運動",
			SkillName: "原始技能",
			Level:     "初一級",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		skill, err = skillRepo.Create(ctx, skill)
		if err != nil {
			t.Fatalf("建立測試技能失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)

		// 更新技能
		req := &services.UpdateSkillRequest{
			Category:  "舞蹈",
			SkillName: "更新後的技能",
		}
		updatedSkill, eInfo, err := svc.UpdateSkill(ctx, skill.ID, teacher.ID, req)

		if err != nil {
			t.Fatalf("UpdateSkill 失敗: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("UpdateSkill 返回錯誤: %s", eInfo.Msg)
		}

		// 驗證更新後的欄位
		if updatedSkill.Category != req.Category {
			t.Errorf("Category 不匹配: 預期 %s, 實際 %s", req.Category, updatedSkill.Category)
		}
		if updatedSkill.SkillName != req.SkillName {
			t.Errorf("SkillName 不匹配: 預期 %s, 實際 %s", req.SkillName, updatedSkill.SkillName)
		}
	})

	t.Run("UpdateSkill_Forbidden_DifferentTeacher", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestUpdateSkillForbid%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestUpdateSkillForbid%").Delete(&models.TeacherSkill{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		skillRepo := repositories.NewTeacherSkillRepository(appInstance)

		// 建立測試老師和技能
		teacher := models.Teacher{
			Name:       "TestUpdateSkillForbid_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("updateforbid%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-updateforbid-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		skill := models.TeacherSkill{
			TeacherID: teacher.ID,
			Category:  "運動",
			SkillName: "私人技能",
			Level:     "中級",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		skill, err = skillRepo.Create(ctx, skill)
		if err != nil {
			t.Fatalf("建立測試技能失敗: %v", err)
		}

		// 調試輸出
		t.Logf("skill.ID after Create: %d", skill.ID)

		svc := services.NewTeacherProfileService(appInstance)

		// 嘗試用不同老師 ID 更新
		req := &services.UpdateSkillRequest{
			Category:  "測試",
			SkillName: "被拒絕的更新",
		}
		_, eInfo, err := svc.UpdateSkill(ctx, skill.ID, 99999999, req)

		// FORBIDDEN 情況下 eInfo 有值，err 為 nil
		if eInfo == nil {
			t.Fatal("預期 eInfo 有值但返回 nil")
		}
		if eInfo.Code != 130002 { // FORBIDDEN (30002) with appID=1 prefix
			t.Errorf("預期 FORBIDDEN 錯誤，實際為 %d", eInfo.Code)
		}
	})

	t.Run("UpdateSkill_NotFound", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}

		ctx := context.Background()
		svc := services.NewTeacherProfileService(appInstance)

		req := &services.UpdateSkillRequest{
			Category:  "測試",
			SkillName: "不存在的技能",
		}
		_, eInfo, err := svc.UpdateSkill(ctx, 99999999, 1, req)

		if err == nil {
			t.Fatal("預期錯誤但返回 nil")
		}
		if eInfo == nil {
			t.Fatal("預期錯誤資訊但返回 nil")
		}
	})

	t.Run("DeleteSkill_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestDeleteSkill%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestDeleteSkill%").Delete(&models.TeacherSkill{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		skillRepo := repositories.NewTeacherSkillRepository(appInstance)

		// 建立測試老師和技能
		teacher := models.Teacher{
			Name:       "TestDeleteSkill_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("deleteskill%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-deleteskill-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		skill := models.TeacherSkill{
			TeacherID: teacher.ID,
			Category:  "運動",
			SkillName: "將被刪除的技能",
			Level:     "中級",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		skill, err = skillRepo.Create(ctx, skill)
		if err != nil {
			t.Fatalf("建立測試技能失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)

		// 刪除技能
		errInfo := svc.DeleteSkill(ctx, skill.ID, teacher.ID)

		if errInfo != nil {
			t.Fatalf("DeleteSkill 返回錯誤: %s", errInfo.Msg)
		}

		// 驗證技能已被刪除
		_, err = skillRepo.GetByID(ctx, skill.ID)
		if err == nil {
			t.Error("技能應該已被刪除，但仍能找到")
		}
	})

	t.Run("DeleteSkill_Forbidden_DifferentTeacher", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestDeleteSkillForbid%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestDeleteSkillForbid%").Delete(&models.TeacherSkill{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		skillRepo := repositories.NewTeacherSkillRepository(appInstance)

		// 建立測試老師和技能
		teacher := models.Teacher{
			Name:       "TestDeleteSkillForbid_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("deleteforbid%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-deleteforbid-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		skill := models.TeacherSkill{
			TeacherID: teacher.ID,
			Category:  "運動",
			SkillName: "私人技能-不可刪除",
			Level:     "中級",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		skill, err = skillRepo.Create(ctx, skill)
		if err != nil {
			t.Fatalf("建立測試技能失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)

		// 嘗試用不同老師 ID 刪除
		errInfo := svc.DeleteSkill(ctx, skill.ID, 99999999)

		if errInfo == nil {
			t.Fatal("預期錯誤但返回 nil")
		}
		if errInfo.Code != 130002 { // FORBIDDEN (30002) with appID=1 prefix
			t.Errorf("預期 FORBIDDEN 錯誤，實際為 %d", errInfo.Code)
		}
	})
}

// TestTeacherProfileService_CertificateCRUD 測試證照 CRUD 操作
func TestTeacherProfileService_CertificateCRUD(t *testing.T) {
	t.Run("CreateCertificate_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestCertCRUD%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestCertCRUD%").Delete(&models.TeacherCertificate{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)

		// 建立測試老師
		teacher := models.Teacher{
			Name:       "TestCertCRUD_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("cert%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-cert-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)

		// 建立證照
		req := &services.CreateCertificateRequest{
			Name:     "瑜珈師資證照",
			FileURL:  "https://example.com/cert.pdf",
			IssuedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		}
		cert, eInfo, err := svc.CreateCertificate(ctx, teacher.ID, req)

		if err != nil {
			t.Fatalf("CreateCertificate 失敗: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("CreateCertificate 返回錯誤: %s", eInfo.Msg)
		}
		if cert == nil {
			t.Fatal("CreateCertificate 返回空結果")
		}

		// 驗證欄位
		if cert.TeacherID != teacher.ID {
			t.Errorf("TeacherID 不匹配: 預期 %d, 實際 %d", teacher.ID, cert.TeacherID)
		}
		if cert.Name != req.Name {
			t.Errorf("Name 不匹配: 預期 %s, 實際 %s", req.Name, cert.Name)
		}
		if cert.FileURL != req.FileURL {
			t.Errorf("FileURL 不匹配: 預期 %s, 實際 %s", req.FileURL, cert.FileURL)
		}
	})

	t.Run("GetCertificates_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestGetCerts%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestGetCerts%").Delete(&models.TeacherCertificate{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		certRepo := repositories.NewTeacherCertificateRepository(appInstance)

		// 建立測試老師
		teacher := models.Teacher{
			Name:       "TestGetCerts_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("getcerts%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-getcerts-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立多個證照
		for i := 0; i < 2; i++ {
			certRepo.Create(ctx, models.TeacherCertificate{
				TeacherID: teacher.ID,
				Name:      fmt.Sprintf("證照 %d", i+1),
				FileURL:   fmt.Sprintf("https://example.com/cert%d.pdf", i+1),
				IssuedAt:  time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}

		svc := services.NewTeacherProfileService(appInstance)
		certs, eInfo, err := svc.GetCertificates(ctx, teacher.ID)

		if err != nil {
			t.Fatalf("GetCertificates 失敗: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("GetCertificates 返回錯誤: %s", eInfo.Msg)
		}
		if len(certs) != 2 {
			t.Errorf("預期 2 個證照，實際 %d 個", len(certs))
		}
	})

	t.Run("DeleteCertificate_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestDeleteCert%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestDeleteCert%").Delete(&models.TeacherCertificate{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		certRepo := repositories.NewTeacherCertificateRepository(appInstance)

		// 建立測試老師和證照
		teacher := models.Teacher{
			Name:       "TestDeleteCert_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("deletecert%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-deletecert-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		cert := models.TeacherCertificate{
			TeacherID: teacher.ID,
			Name:      "將被刪除的證照",
			FileURL:   "https://example.com/cert-delete.pdf",
			IssuedAt:  time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		cert, err = certRepo.Create(ctx, cert)
		if err != nil {
			t.Fatalf("建立測試證照失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)

		// 刪除證照
		errInfo := svc.DeleteCertificate(ctx, cert.ID, teacher.ID)

		if errInfo != nil {
			t.Fatalf("DeleteCertificate 返回錯誤: %s", errInfo.Msg)
		}

		// 驗證證照已被刪除
		_, err = certRepo.GetByID(ctx, cert.ID)
		if err == nil {
			t.Error("證照應該已被刪除，但仍能找到")
		}
	})

	t.Run("DeleteCertificate_Forbidden_DifferentTeacher", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestDeleteCertForbid%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestDeleteCertForbid%").Delete(&models.TeacherCertificate{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		certRepo := repositories.NewTeacherCertificateRepository(appInstance)

		// 建立測試老師和證照
		teacher := models.Teacher{
			Name:       "TestDeleteCertForbid_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("deletecertforbid%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-deletecertforbid-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		cert := models.TeacherCertificate{
			TeacherID: teacher.ID,
			Name:      "私人證照-不可刪除",
			FileURL:   "https://example.com/private-cert.pdf",
			IssuedAt:  time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		cert, err = certRepo.Create(ctx, cert)
		if err != nil {
			t.Fatalf("建立測試證照失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)

		// 嘗試用不同老師 ID 刪除
		errInfo := svc.DeleteCertificate(ctx, cert.ID, 99999999)

		if errInfo == nil {
			t.Fatal("預期錯誤但返回 nil")
		}
		if errInfo.Code != 130002 { // FORBIDDEN (30002) with appID=1 prefix
			t.Errorf("預期 FORBIDDEN 錯誤，實際為 %d", errInfo.Code)
		}
	})

	t.Run("DeleteCertificate_NotFound", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}

		ctx := context.Background()
		svc := services.NewTeacherProfileService(appInstance)

		errInfo := svc.DeleteCertificate(ctx, 99999999, 1)

		if errInfo == nil {
			t.Fatal("預期錯誤但返回 nil")
		}
	})
}

// TestTeacherProfileService_GetCenters 測試取得老師加入的中心列表
func TestTeacherProfileService_GetCenters(t *testing.T) {
	t.Run("GetCenters_Success", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestGetCenters%").Delete(&models.Teacher{})
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestCenter%").Delete(&models.Center{})
			appInstance.MySQL.RDB.Where("teacher_id IN (SELECT id FROM teachers WHERE name LIKE ?)", "TestGetCenters%").Delete(&models.CenterMembership{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)
		centerRepo := repositories.NewCenterRepository(appInstance)
		membershipRepo := repositories.NewCenterMembershipRepository(appInstance)

		// 建立測試老師
		teacher := models.Teacher{
			Name:       "TestGetCenters_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("getcenters%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-getcenters-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		// 建立測試中心
		center := models.Center{
			Name:      "TestCenter_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt: time.Now(),
		}
		center, err = centerRepo.Create(ctx, center)
		if err != nil {
			t.Fatalf("建立測試中心失敗: %v", err)
		}

		// 建立會籍
		membership := models.CenterMembership{
			TeacherID: teacher.ID,
			CenterID:  center.ID,
			Status:    "ACTIVE",
			CreatedAt: time.Now(),
		}
		membership, err = membershipRepo.Create(ctx, membership)
		if err != nil {
			t.Fatalf("建立測試會籍失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)
		centers, eInfo, err := svc.GetCenters(ctx, teacher.ID)

		if err != nil {
			t.Fatalf("GetCenters 失敗: %v", err)
		}
		if eInfo != nil {
			t.Fatalf("GetCenters 返回錯誤: %s", eInfo.Msg)
		}
		if len(centers) != 1 {
			t.Errorf("預期 1 個中心，實際 %d 個", len(centers))
		}
		if centers[0].CenterID != center.ID {
			t.Errorf("CenterID 不匹配: 預期 %d, 實際 %d", center.ID, centers[0].CenterID)
		}
		if centers[0].CenterName != center.Name {
			t.Errorf("CenterName 不匹配: 預期 %s, 實際 %s", center.Name, centers[0].CenterName)
		}
	})

	t.Run("GetCenters_NoMembership", func(t *testing.T) {
		appInstance := setupTeacherProfileTestApp(t)
		if appInstance == nil {
			return
		}
		defer func() {
			appInstance.MySQL.RDB.Where("name LIKE ?", "TestNoMembership%").Delete(&models.Teacher{})
		}()

		ctx := context.Background()
		teacherRepo := repositories.NewTeacherRepository(appInstance)

		// 建立測試老師（沒有會籍）
		teacher := models.Teacher{
			Name:       "TestNoMembership_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			Email:      fmt.Sprintf("nomembership%d@test.com", time.Now().UnixNano()),
			LineUserID: "test-line-nomembership-" + fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teacher, err := teacherRepo.Create(ctx, teacher)
		if err != nil {
			t.Fatalf("建立測試老師失敗: %v", err)
		}

		svc := services.NewTeacherProfileService(appInstance)
		centers, _, err := svc.GetCenters(ctx, teacher.ID)

		if err != nil {
			t.Fatalf("GetCenters 失敗: %v", err)
		}
		if len(centers) != 0 {
			t.Errorf("預期 0 個中心，實際 %d 個", len(centers))
		}
	})
}
