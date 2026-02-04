package controllers

import (
	"path/filepath"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/resources"
	"timeLedger/app/services"
	"timeLedger/libs"

	"github.com/gin-gonic/gin"
)

// TeacherProfileController 老師個人檔案相關 API
type TeacherProfileController struct {
	BaseController
	app            *app.App
	profileService *services.TeacherProfileService
	hashtagRepo    *repositories.HashtagRepository
	r2Storage      *libs.R2StorageService
	localStorage   *libs.LocalStorageService
}

func NewTeacherProfileController(app *app.App) *TeacherProfileController {
	r2Storage, _ := libs.NewR2StorageService(app.Env)
	localStorage := libs.NewLocalStorageService("./uploads/certificates")

	return &TeacherProfileController{
		app:            app,
		profileService: services.NewTeacherProfileService(app),
		hashtagRepo:    repositories.NewHashtagRepository(app),
		r2Storage:      r2Storage,
		localStorage:   localStorage,
	}
}

// GetProfile 取得老師個人資料
// @Summary 取得老師個人資料
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=resources.TeacherProfileResource}
// @Router /api/v1/teacher/me/profile [get]
func (ctl *TeacherProfileController) GetProfile(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	profile, errInfo, err := ctl.profileService.GetProfile(ctx, teacherID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(profile)
}

// SearchHashtags 搜尋標籤
// @Summary 搜尋標籤
// @Tags Hashtag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "搜尋關鍵字"
// @Success 200 {object} global.ApiResponse{data=[]resources.HashtagResource}
// @Router /api/v1/hashtags/search [get]
func (ctl *TeacherProfileController) SearchHashtags(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	query := helper.QueryStringOrDefault("q", "")
	if query == "" {
		helper.BadRequest("Query parameter 'q' is required")
		return
	}

	query = strings.TrimPrefix(query, "#")

	hashtags, err := ctl.hashtagRepo.Search(ctx, query)
	if err != nil {
		helper.InternalError("Failed to search hashtags")
		return
	}

	if hashtags == nil {
		hashtags = []models.Hashtag{}
	}

	var hashtagResources []resources.HashtagResource
	for _, h := range hashtags {
		hashtagResources = append(hashtagResources, resources.HashtagResource{
			ID:         h.ID,
			Name:       h.Name,
			UsageCount: h.UsageCount,
		})
	}

	helper.Success(hashtagResources)
}

// CreateHashtag 建立新標籤
// @Summary 建立新標籤
// @Tags Hashtag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateHashtagRequest true "標籤資訊"
// @Success 200 {object} global.ApiResponse{data=models.Hashtag}
// @Router /api/v1/hashtags [post]
func (ctl *TeacherProfileController) CreateHashtag(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	var req CreateHashtagRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	name := req.Name
	if !strings.HasPrefix(name, "#") {
		name = "#" + name
	}

	existing, err := ctl.hashtagRepo.GetByName(ctx, name)
	if err == nil && existing != nil {
		helper.Success(existing)
		return
	}

	hashtag := models.Hashtag{
		Name:       name,
		UsageCount: 1,
	}
	_, err = ctl.hashtagRepo.Create(ctx, hashtag)
	if err != nil {
		helper.InternalError("Failed to create hashtag")
		return
	}

	helper.Success(hashtag)
}

// UpdateProfile 更新老師個人資料
// @Summary 更新老師個人資料
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateTeacherProfileRequest true "個人資料"
// @Success 200 {object} global.ApiResponse{data=resources.TeacherProfileResource}
// @Router /api/v1/teacher/me/profile [put]
func (ctl *TeacherProfileController) UpdateProfile(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	var req UpdateTeacherProfileRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	profile, errInfo, err := ctl.profileService.UpdateProfile(ctx, teacherID, &services.UpdateProfileRequest{
		Name:              req.Name,
		Bio:               req.Bio,
		City:              req.City,
		District:          req.District,
		PublicContactInfo: req.PublicContactInfo,
		IsOpenToHiring:    req.IsOpenToHiring,
		PersonalHashtags:  req.PersonalHashtags,
	})
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(profile)
}

// GetCenters 取得老師已加入的中心列表
// @Summary 取得老師已加入的中心列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.CenterMembershipResource}
// @Router /api/v1/teacher/me/centers [get]
func (ctl *TeacherProfileController) GetCenters(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	centers, errInfo, err := ctl.profileService.GetCenters(ctx, teacherID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(centers)
}

// GetSkills 取得老師技能列表
// @Summary 取得老師技能列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.TeacherSkill}
// @Router /api/v1/teacher/me/skills [get]
func (ctl *TeacherProfileController) GetSkills(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	skills, errInfo, err := ctl.profileService.GetSkills(ctx, teacherID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(skills)
}

// CreateSkill 新增老師技能
// @Summary 新增老師技能
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSkillRequest true "技能資訊"
// @Success 200 {object} global.ApiResponse{data=models.TeacherSkill}
// @Router /api/v1/teacher/me/skills [post]
func (ctl *TeacherProfileController) CreateSkill(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	var req CreateSkillRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	skill, errInfo, err := ctl.profileService.CreateSkill(ctx, teacherID, &services.CreateSkillRequest{
		Category:   req.Category,
		SkillName:  req.SkillName,
		Level:      req.Level,
		HashtagIDs: req.HashtagIDs,
	})
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(skill)
}

// DeleteSkill 刪除老師技能
// @Summary 刪除老師技能
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "技能ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/skills/{id} [delete]
func (ctl *TeacherProfileController) DeleteSkill(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	skillID := helper.MustParamUint("id")
	if skillID == 0 {
		return
	}

	errInfo := ctl.profileService.DeleteSkill(ctx, skillID, teacherID)
	if errInfo != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(nil)
}

// UpdateSkill 更新老師技能
// @Summary 更新老師技能
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "技能ID"
// @Param request body UpdateSkillRequest true "技能資訊"
// @Success 200 {object} global.ApiResponse{data=models.TeacherSkill}
// @Router /api/v1/teacher/me/skills/{id} [put]
func (ctl *TeacherProfileController) UpdateSkill(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	skillID := helper.MustParamUint("id")
	if skillID == 0 {
		return
	}

	var req UpdateSkillRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	skill, errInfo, err := ctl.profileService.UpdateSkill(ctx, skillID, teacherID, &services.UpdateSkillRequest{
		Category:  req.Category,
		SkillName: req.SkillName,
		Hashtags:  req.Hashtags,
	})
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(skill)
}

// GetCertificates 取得老師證照列表
// @Summary 取得老師證照列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.TeacherCertificate}
// @Router /api/v1/teacher/me/certificates [get]
func (ctl *TeacherProfileController) GetCertificates(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	certificates, errInfo, err := ctl.profileService.GetCertificates(ctx, teacherID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(certificates)
}

// CreateCertificate 新增老師證照
// @Summary 新增老師證照
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCertificateRequest true "證照資訊"
// @Success 200 {object} global.ApiResponse{data=models.TeacherCertificate}
// @Router /api/v1/teacher/me/certificates [post]
func (ctl *TeacherProfileController) CreateCertificate(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	var req CreateCertificateRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	certificate, errInfo, err := ctl.profileService.CreateCertificate(ctx, teacherID, &services.CreateCertificateRequest{
		Name:     req.Name,
		FileURL:  req.FileURL,
		IssuedAt: req.IssuedAt,
	})
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(certificate)
}

// DeleteCertificate 刪除老師證照
// @Summary 刪除老師證照
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "證照ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/certificates/{id} [delete]
func (ctl *TeacherProfileController) DeleteCertificate(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	certificateID := helper.MustParamUint("id")
	if certificateID == 0 {
		return
	}

	errInfo := ctl.profileService.DeleteCertificate(ctx, certificateID, teacherID)
	if errInfo != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(nil)
}

// UploadCertificateFile 上傳證照檔案
// @Summary 上傳證照檔案
// @Tags Teacher
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "證照檔案"
// @Success 200 {object} global.ApiResponse{data=UploadFileResponse}
// @Router /api/v1/teacher/me/certificates/upload [post]
func (ctl *TeacherProfileController) UploadCertificateFile(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		helper.BadRequest("No file uploaded: " + err.Error())
		return
	}

	maxSize := 10 * 1024 * 1024
	if file.Size > int64(maxSize) {
		helper.BadRequest("File size exceeds maximum limit (10MB)")
		return
	}

	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		helper.BadRequest("Invalid file type. Allowed: jpg, jpeg, png, pdf")
		return
	}

	src, err := file.Open()
	if err != nil {
		helper.InternalError("Failed to open file: " + err.Error())
		return
	}
	defer src.Close()

	contentType := libs.GetContentType(file.Filename)

	var fileURL string
	if ctl.r2Storage != nil && ctl.r2Storage.IsEnabled() {
		fileURL, err = ctl.r2Storage.UploadFile(ctx, src, file.Filename, contentType)
		if err != nil {
			helper.InternalError("Failed to upload to R2: " + err.Error())
			return
		}
	} else {
		fileURL, err = ctl.localStorage.UploadFile(ctx, src, file.Filename, contentType)
		if err != nil {
			helper.InternalError("Failed to save file locally: " + err.Error())
			return
		}
	}

	helper.Success(UploadFileResponse{
		FileURL:  fileURL,
		FileName: file.Filename,
		FileSize: file.Size,
	})
}

// UploadFileResponse 上傳檔案回應結構
type UploadFileResponse struct {
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

// ==================== Request Types ====================

type UpdateTeacherProfileRequest struct {
	Name              string   `json:"name"`
	Bio               string   `json:"bio"`
	City              string   `json:"city"`
	District          string   `json:"district"`
	PublicContactInfo string   `json:"public_contact_info"`
	IsOpenToHiring    bool     `json:"is_open_to_hiring"`
	PersonalHashtags  []string `json:"personal_hashtags"`
}

type CreateSkillRequest struct {
	Category   string `json:"category" binding:"required"`
	SkillName  string `json:"skill_name" binding:"required"`
	Level      string `json:"level"`
	HashtagIDs []uint `json:"hashtag_ids"`
}

type UpdateSkillRequest struct {
	Category  string   `json:"category" binding:"required"`
	SkillName string   `json:"skill_name" binding:"required"`
	Hashtags  []string `json:"hashtags"`
}

type CreateHashtagRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateCertificateRequest struct {
	Name     string    `json:"name" binding:"required"`
	FileURL  string    `json:"file_url" binding:"required"`
	IssuedAt time.Time `json:"issued_at" binding:"required"`
}
