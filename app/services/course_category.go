package services

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"
)

// CourseCategoryService 課程類別 Service
type CourseCategoryService struct {
	BaseService
	app            *app.App
	categoryRepo   *repositories.CourseCategoryRepository
}

// NewCourseCategoryService 建立 CourseCategoryService 實例
func NewCourseCategoryService(appInstance *app.App) *CourseCategoryService {
	baseSvc := NewBaseService(appInstance, "CourseCategoryService")
	return &CourseCategoryService{
		BaseService:    *baseSvc,
		app:            appInstance,
		categoryRepo:   repositories.NewCourseCategoryRepository(appInstance),
	}
}

// GetCategories 取得課程類別列表
func (s *CourseCategoryService) GetCategories(ctx context.Context, centerID uint) ([]models.CourseCategory, *errInfos.Res, error) {
	categories, err := s.categoryRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	return categories, nil, nil
}

// CreateCategory 建立課程類別
func (s *CourseCategoryService) CreateCategory(ctx context.Context, centerID uint, name string) (*models.CourseCategory, *errInfos.Res, error) {
	// 檢查是否已存在
	existing, err := s.categoryRepo.GetByNameAndCenterID(ctx, centerID, name)
	if err == nil && existing != nil {
		return nil, s.app.Err.New(errInfos.DUPLICATE), nil
	}

	category := &models.CourseCategory{
		CenterID: centerID,
		Name:     name,
	}

	created, err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return created, nil, nil
}

// UpdateCategory 更新課程類別
func (s *CourseCategoryService) UpdateCategory(ctx context.Context, centerID, categoryID uint, name string) (*models.CourseCategory, *errInfos.Res, error) {
	// 取得現有類別
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 檢查權限
	if category.CenterID != centerID {
		return nil, s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 檢查新名稱是否與其他類別重複
	existing, err := s.categoryRepo.GetByNameAndCenterID(ctx, centerID, name)
	if err == nil && existing != nil && existing.ID != categoryID {
		return nil, s.app.Err.New(errInfos.DUPLICATE), nil
	}

	// 更新
	category.Name = name
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return category, nil, nil
}

// DeleteCategory 刪除課程類別
func (s *CourseCategoryService) DeleteCategory(ctx context.Context, centerID, categoryID uint) (*errInfos.Res, error) {
	// 取得現有類別
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 檢查權限
	if category.CenterID != centerID {
		return s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 刪除
	if err := s.categoryRepo.DeleteByID(ctx, categoryID); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return nil, nil
}
