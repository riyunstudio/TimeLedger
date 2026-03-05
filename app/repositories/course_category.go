package repositories

import (
	"context"

	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

// CourseCategoryRepository 課程類別 Repository
type CourseCategoryRepository struct {
	app        *app.App
	dbRead     *gorm.DB
	dbWrite    *gorm.DB
}

// NewCourseCategoryRepository 建立 CourseCategoryRepository 實例
func NewCourseCategoryRepository(appInstance *app.App) *CourseCategoryRepository {
	return &CourseCategoryRepository{
		app:        appInstance,
		dbRead:     appInstance.MySQL.RDB,
		dbWrite:    appInstance.MySQL.WDB,
	}
}

// GetDBRead 取得讀取資料庫連線
func (r *CourseCategoryRepository) GetDBRead() *gorm.DB {
	return r.dbRead
}

// GetDBWrite 取得寫入資料庫連線
func (r *CourseCategoryRepository) GetDBWrite() *gorm.DB {
	return r.dbWrite
}

// ListByCenterID 取得某中心的課程類別列表
func (r *CourseCategoryRepository) ListByCenterID(ctx context.Context, centerID uint) ([]models.CourseCategory, error) {
	var categories []models.CourseCategory
	err := r.dbRead.WithContext(ctx).
		Where("center_id = ?", centerID).
		Order("name ASC").
		Find(&categories).Error
	return categories, err
}

// Create 建立課程類別
func (r *CourseCategoryRepository) Create(ctx context.Context, category *models.CourseCategory) (*models.CourseCategory, error) {
	err := r.dbWrite.WithContext(ctx).Create(category).Error
	return category, err
}

// Update 更新課程類別
func (r *CourseCategoryRepository) Update(ctx context.Context, category *models.CourseCategory) error {
	return r.dbWrite.WithContext(ctx).Save(category).Error
}

// DeleteByID 刪除課程類別
func (r *CourseCategoryRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.dbWrite.WithContext(ctx).Delete(&models.CourseCategory{}, id).Error
}

// GetByID 取得課程類別
func (r *CourseCategoryRepository) GetByID(ctx context.Context, id uint) (*models.CourseCategory, error) {
	var category models.CourseCategory
	err := r.dbRead.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByNameAndCenterID 根據名稱和中心ID取得課程類別
func (r *CourseCategoryRepository) GetByNameAndCenterID(ctx context.Context, centerID uint, name string) (*models.CourseCategory, error) {
	var category models.CourseCategory
	err := r.dbRead.WithContext(ctx).
		Where("center_id = ? AND name = ?", centerID, name).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// CountByCenterID 取得某中心的課程類別數量
func (r *CourseCategoryRepository) CountByCenterID(ctx context.Context, centerID uint) (int64, error) {
	var count int64
	err := r.dbRead.WithContext(ctx).
		Model(&models.CourseCategory{}).
		Where("center_id = ?", centerID).
		Count(&count).Error
	return count, err
}
