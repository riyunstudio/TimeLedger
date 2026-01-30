package repositories

import (
	"context"
	"timeLedger/app/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GenericRepository[T] provides standard CRUD operations for any model type T.
// It is designed to reduce repetitive CRUD code across repositories while
// allowing each repository to still implement custom methods as needed.
//
// Usage:
//
//	type UserRepository struct {
//		GenericRepository[models.User]
//		app *app.App
//	}
//
//	func NewUserRepository(app *app.App) *UserRepository {
//		return &UserRepository{
//			GenericRepository: NewGenericRepository[models.User](app, app.MySQL.RDB, app.MySQL.WDB),
//			app: app,
//		}
//	}
type GenericRepository[T models.IModel] struct {
	dbRead  *gorm.DB
	dbWrite *gorm.DB
	table   string
}

// NewGenericRepository creates a new GenericRepository instance for the given model type T.
//
// Parameters:
//   - dbRead: Read-only database connection (uses RDB/slave)
//   - dbWrite: Write database connection (uses WDB/master)
//
// Returns a configured GenericRepository ready for CRUD operations.
func NewGenericRepository[T models.IModel](dbRead *gorm.DB, dbWrite *gorm.DB) GenericRepository[T] {
	var model T

	// Create db connections with context, handling nil gracefully
	var dbReadWithCtx, dbWriteWithCtx *gorm.DB
	if dbRead != nil {
		dbReadWithCtx = dbRead.WithContext(context.Background())
	}
	if dbWrite != nil {
		dbWriteWithCtx = dbWrite.WithContext(context.Background())
	}

	return GenericRepository[T]{
		dbRead:  dbReadWithCtx,
		dbWrite: dbWriteWithCtx,
		table:   model.TableName(),
	}
}

// GetDBRead returns the read database connection for custom queries.
func (rp *GenericRepository[T]) GetDBRead() *gorm.DB {
	return rp.dbRead
}

// GetDBWrite returns the write database connection for custom queries.
func (rp *GenericRepository[T]) GetDBWrite() *gorm.DB {
	return rp.dbWrite
}

// GetTableName returns the table name for this repository.
func (rp *GenericRepository[T]) GetTableName() string {
	return rp.table
}

// withContext returns a new GenericRepository with the given context.
// All subsequent operations will use this context for cancellation and timeouts.
func (rp *GenericRepository[T]) withContext(ctx context.Context) *GenericRepository[T] {
	return &GenericRepository[T]{
		dbRead:  rp.dbRead.WithContext(ctx),
		dbWrite: rp.dbWrite.WithContext(ctx),
		table:   rp.table,
	}
}

// GetByID retrieves a single record by its primary key ID.
//
// Returns:
//   - T: The found record
//   - error: gorm.ErrRecordNotFound if not found
func (rp *GenericRepository[T]) GetByID(ctx context.Context, id uint) (T, error) {
	var data T
	err := rp.withContext(ctx).dbRead.Table(rp.table).Where("id = ?", id).First(&data).Error
	return data, err
}

// GetByIDWithCenterScope retrieves a record with center_id scope check.
// This is useful for multi-tenant data isolation.
//
// Parameters:
//   - id: Record ID
//   - centerID: The center_id to scope the query to
//
// Returns:
//   - T: The found record
//   - error: gorm.ErrRecordNotFound if not found or center_id doesn't match
func (rp *GenericRepository[T]) GetByIDWithCenterScope(ctx context.Context, id, centerID uint) (T, error) {
	var data T
	err := rp.withContext(ctx).dbRead.Table(rp.table).
		Where("id = ? AND center_id = ?", id, centerID).
		First(&data).Error
	return data, err
}

// First retrieves the first record matching the given conditions.
// If no conditions provided, returns the first record in the table.
//
// Returns:
//   - T: The found record
//   - error: gorm.ErrRecordNotFound if not found
func (rp *GenericRepository[T]) First(ctx context.Context, conditions ...interface{}) (T, error) {
	var data T
	query := rp.withContext(ctx).dbRead.Table(rp.table)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.First(&data).Error
	return data, err
}

// FirstWithCenterScope retrieves the first record matching conditions with center_id scope.
func (rp *GenericRepository[T]) FirstWithCenterScope(ctx context.Context, centerID uint, conditions ...interface{}) (T, error) {
	var data T
	query := rp.withContext(ctx).dbRead.Table(rp.table).
		Where("center_id = ?", centerID)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.First(&data).Error
	return data, err
}

// Find retrieves all records matching the given conditions.
// If no conditions provided, returns all records in the table.
//
// Returns:
//   - []T: Slice of matching records (empty slice if none found)
//   - error: Any error encountered during query
func (rp *GenericRepository[T]) Find(ctx context.Context, conditions ...interface{}) ([]T, error) {
	var data []T
	query := rp.withContext(ctx).dbRead.Table(rp.table)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Find(&data).Error
	return data, err
}

// FindWithCenterScope retrieves all records matching conditions with center_id scope.
func (rp *GenericRepository[T]) FindWithCenterScope(ctx context.Context, centerID uint, conditions ...interface{}) ([]T, error) {
	var data []T
	query := rp.withContext(ctx).dbRead.Table(rp.table).
		Where("center_id = ?", centerID)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Find(&data).Error
	return data, err
}

// FindAll retrieves all records in the table.
//
// Returns:
//   - []T: Slice of all records
//   - error: Any error encountered during query
func (rp *GenericRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	return rp.Find(ctx)
}

// FindPaged retrieves records with pagination support.
//
// Parameters:
//   - ctx: Context
//   - page: Page number (1-indexed, default 1)
//   - limit: Number of records per page (default 20, max 100)
//   - orderBy: Ordering field (e.g., "created_at DESC")
//   - conditions: Additional WHERE conditions
//
// Returns:
//   - []T: Slice of records for the page
//   - total: Total count of matching records
//   - error: Any error encountered
func (rp *GenericRepository[T]) FindPaged(ctx context.Context, page, limit int, orderBy string, conditions ...interface{}) ([]T, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if orderBy == "" {
		orderBy = "id DESC"
	}

	var total int64
	var data []T

	// Count query
	countQuery := rp.withContext(ctx).dbRead.Table(rp.table)
	if len(conditions) > 0 {
		countQuery = countQuery.Where(conditions[0], conditions[1:]...)
	}
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Data query with pagination
	dataQuery := rp.withContext(ctx).dbRead.Table(rp.table)
	if len(conditions) > 0 {
		dataQuery = dataQuery.Where(conditions[0], conditions[1:]...)
	}
	err = dataQuery.Order(orderBy).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&data).Error

	return data, total, err
}

// Create inserts a new record into the database.
//
// Parameters:
//   - data: The record to create (pointer to T)
//
// Returns:
//   - T: The created record (with generated fields like ID, CreatedAt)
//   - error: Any error encountered
func (rp *GenericRepository[T]) Create(ctx context.Context, data T) (T, error) {
	err := rp.withContext(ctx).dbWrite.Table(rp.table).Create(&data).Error
	return data, err
}

// CreateBatch inserts multiple records in a single transaction.
//
// Parameters:
//   - data: Slice of records to create
//
// Returns:
//   - []T: The created records
//   - error: Any error encountered
func (rp *GenericRepository[T]) CreateBatch(ctx context.Context, data []T) ([]T, error) {
	if len(data) == 0 {
		return data, nil
	}
	err := rp.withContext(ctx).dbWrite.Table(rp.table).Create(&data).Error
	return data, err
}

// Update updates an existing record. The record must have an ID set.
//
// Parameters:
//   - data: The record to update (must have ID)
//
// Returns:
//   - error: Any error encountered
func (rp *GenericRepository[T]) Update(ctx context.Context, data T) error {
	return rp.withContext(ctx).dbWrite.Table(rp.table).Save(&data).Error
}

// UpdateFields updates specific fields of a record identified by ID.
//
// Parameters:
//   - id: Record ID
//   - fields: Map of field names to values to update
//
// Returns:
//   - error: Any error encountered
func (rp *GenericRepository[T]) UpdateFields(ctx context.Context, id uint, fields map[string]interface{}) error {
	return rp.withContext(ctx).dbWrite.Table(rp.table).
		Where("id = ?", id).
		Updates(fields).Error
}

// UpdateFieldsWithCenterScope updates fields with center_id scope for multi-tenant isolation.
func (rp *GenericRepository[T]) UpdateFieldsWithCenterScope(ctx context.Context, id, centerID uint, fields map[string]interface{}) error {
	return rp.withContext(ctx).dbWrite.Table(rp.table).
		Where("id = ? AND center_id = ?", id, centerID).
		Updates(fields).Error
}

// UpdateWhere updates all records matching conditions with the given fields.
//
// Parameters:
//   - conditions: WHERE conditions (query, args...)
//   - fields: Map of field names to values
//
// Returns:
//   - affected rows count
//   - error: Any error encountered
func (rp *GenericRepository[T]) UpdateWhere(ctx context.Context, conditions, fields interface{}) (int64, error) {
	result := rp.withContext(ctx).dbWrite.Table(rp.table).
		Where(conditions).
		Updates(fields)
	return result.RowsAffected, result.Error
}

// Upsert inserts a new record or updates an existing one based on unique constraints.
//
// Parameters:
//   - data: The record to upsert
//   - onConflict: Columns to use for conflict resolution (e.g., "email,phone")
//   - updateColumns: Columns to update on conflict (use gorm.Expr("column = excluded.column") for all)
//
// Returns:
//   - error: Any error encountered
func (rp *GenericRepository[T]) Upsert(ctx context.Context, data T, onConflict string, updateColumns []string) error {
	onConflictClause := clause.OnConflict{
		Columns:   []clause.Column{{Name: onConflict}},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}

	return rp.withContext(ctx).dbWrite.Table(rp.table).
		Clauses(onConflictClause).
		Create(&data).Error
}

// DeleteByID deletes a record by its ID (soft delete if enabled in model).
//
// Parameters:
//   - id: Record ID to delete
//
// Returns:
//   - error: Any error encountered
func (rp *GenericRepository[T]) DeleteByID(ctx context.Context, id uint) error {
	return rp.withContext(ctx).dbWrite.Table(rp.table).Where("id = ?", id).Delete(new(T)).Error
}

// DeleteByIDWithCenterScope deletes a record with center_id scope for multi-tenant isolation.
func (rp *GenericRepository[T]) DeleteByIDWithCenterScope(ctx context.Context, id, centerID uint) error {
	return rp.withContext(ctx).dbWrite.Table(rp.table).
		Where("id = ? AND center_id = ?", id, centerID).
		Delete(new(T)).Error
}

// DeleteWhere deletes all records matching conditions (soft delete if enabled).
//
// Parameters:
//   - conditions: WHERE conditions (query, args...)
//
// Returns:
//   - affected rows count
//   - error: Any error encountered
func (rp *GenericRepository[T]) DeleteWhere(ctx context.Context, conditions ...interface{}) (int64, error) {
	result := rp.withContext(ctx).dbWrite.Table(rp.table)
	if len(conditions) > 0 {
		result = result.Where(conditions[0], conditions[1:]...)
	}
	result = result.Delete(new(T))
	return result.RowsAffected, result.Error
}

// HardDeleteByID permanently deletes a record by ID (use with caution).
//
// Parameters:
//   - id: Record ID to permanently delete
//
// Returns:
//   - error: Any error encountered
func (rp *GenericRepository[T]) HardDeleteByID(ctx context.Context, id uint) error {
	return rp.withContext(ctx).dbWrite.Table(rp.table).Unscoped().Where("id = ?", id).Delete(new(T)).Error
}

// Exists checks if a record exists with the given conditions.
//
// Parameters:
//   - conditions: WHERE conditions (query, args...)
//
// Returns:
//   - bool: True if at least one record exists
//   - error: Any error encountered
func (rp *GenericRepository[T]) Exists(ctx context.Context, conditions ...interface{}) (bool, error) {
	var count int64
	query := rp.withContext(ctx).dbRead.Table(rp.table)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// ExistsByID checks if a record exists with the given ID.
//
// Returns:
//   - bool: True if record exists
//   - error: Any error encountered
func (rp *GenericRepository[T]) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return rp.Exists(ctx, "id = ?", id)
}

// Count returns the number of records matching the conditions.
//
// Parameters:
//   - conditions: WHERE conditions (query, args...)
//
// Returns:
//   - int64: Number of matching records
//   - error: Any error encountered
func (rp *GenericRepository[T]) Count(ctx context.Context, conditions ...interface{}) (int64, error) {
	var count int64
	query := rp.withContext(ctx).dbRead.Table(rp.table)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Count(&count).Error
	return count, err
}

// CountWithCenterScope returns the count of records with center_id scope.
func (rp *GenericRepository[T]) CountWithCenterScope(ctx context.Context, centerID uint, conditions ...interface{}) (int64, error) {
	var count int64
	query := rp.withContext(ctx).dbRead.Table(rp.table).
		Where("center_id = ?", centerID)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Count(&count).Error
	return count, err
}

// Pluck retrieves specific column values into a slice.
//
// Parameters:
//   - ctx: Context
//   - column: Column name to retrieve
//   - dest: Destination slice
//   - conditions: WHERE conditions (query, args...)
//
// Returns:
//   - error: Any error encountered
func (rp *GenericRepository[T]) Pluck(ctx context.Context, column string, dest interface{}, conditions ...interface{}) error {
	query := rp.withContext(ctx).dbRead.Table(rp.table)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.Pluck(column, dest).Error
}

// PluckWithCenterScope retrieves column values with center_id scope.
func (rp *GenericRepository[T]) PluckWithCenterScope(ctx context.Context, centerID uint, column string, dest interface{}, conditions ...interface{}) error {
	query := rp.withContext(ctx).dbRead.Table(rp.table).
		Where("center_id = ?", centerID)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.Pluck(column, dest).Error
}

// Transaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// If the function panics, the transaction is rolled back.
//
// IMPORTANT: This method creates a NEW repository instance with transaction connections
// to avoid race conditions in concurrent requests. Do NOT modify the original repository.
//
// Parameters:
//   - fn: Function to execute within the transaction (receives context and transaction DB)
//
// Returns:
//   - error: Any error from the function or transaction handling
//
// Usage Example (in Service layer):
//
//	txErr := rp.Transaction(ctx, func(txCtx context.Context, txDB *gorm.DB) error {
//	    // Use txDB directly for DB operations within the transaction
//	    if err := txDB.WithContext(txCtx).Table("courses").Create(&course).Error; err != nil {
//	        return err
//	    }
//	    if err := txDB.WithContext(txCtx).Table("audit_logs").Create(&auditLog).Error; err != nil {
//	        return err
//	    }
//	    return nil
//	})
//
// Note: For transactions that need to use repository methods, use TransactionWithRepo() instead.
// For domain repositories with custom methods, implement a dedicated Transaction method
// (see CourseRepository.Transaction as example).
func (rp *GenericRepository[T]) Transaction(ctx context.Context, fn func(txCtx context.Context, txDB *gorm.DB) error) error {
	return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Pass the transaction DB connection to the callback function
		return fn(ctx, tx)
	})
}

// TransactionCapable is an interface that repositories can implement
// to provide transaction-aware method access within transactions.
// This allows custom methods to be used with transaction connections.
//
// Usage:
//
//	func (rp *CourseRepository) Transaction(ctx context.Context, fn func(txRepo *CourseRepository) error) error {
//	    return rp.GenericRepository.TransactionWithRepo(ctx, txDB, func(txRepo GenericRepository[models.Course]) error {
//	        // Create domain repo wrapper for custom methods
//	        domainRepo := &CourseRepository{
//	            GenericRepository: txRepo,
//	            app:               rp.app,
//	        }
//	        return fn(domainRepo)
//	    })
//	}
type TransactionCapable[T models.IModel] interface {
	Transaction(ctx context.Context, fn func(txRepo interface{}) error) error
	TransactionWithRepo(ctx context.Context, txDB *gorm.DB, fn func(txRepo GenericRepository[T]) error) error
}

// TransactionWithRepo executes a function with a transaction-aware repository.
// This is the preferred pattern for transactions that need repository methods.
//
// Parameters:
//   - ctx: Context for the transaction
//   - txDB: The transaction DB connection (from WDB.Transaction)
//   - fn: Function that receives a transaction-aware generic repository and returns an error
//
// Returns:
//   - error: Any error from the function or transaction handling
//
// Usage Example:
//
//	result, err := rp.TransactionWithRepo(ctx, txDB, func(txRepo GenericRepository[models.Course]) error {
//	    // All operations using txRepo will be within the same transaction
//	    if _, err := txRepo.Create(ctx, data1); err != nil {
//	        return err
//	    }
//	    if _, err := txRepo.Create(ctx, data2); err != nil {
//	        return err
//	    }
//	    return nil
//	})
//
// Note: For domain repositories with custom methods, use the domain-specific
// Transaction method instead (see CourseRepository.Transaction).
func (rp *GenericRepository[T]) TransactionWithRepo(ctx context.Context, txDB *gorm.DB, fn func(txRepo GenericRepository[T]) error) error {
	// Create a new repository instance with transaction connections
	txRepo := GenericRepository[T]{
		dbRead:  txDB.WithContext(ctx),
		dbWrite: txDB.WithContext(ctx),
		table:   rp.table,
	}
	return fn(txRepo)
}

// WithTransaction creates a transaction-aware repository instance.
// This helper is used by domain repositories to implement their Transaction methods.
//
// Usage in domain repository:
//
//	func (rp *CourseRepository) Transaction(ctx context.Context, fn func(txRepo *CourseRepository) error) error {
//	    return rp.GenericRepository.Transactional(ctx, rp.dbWrite, func(txDB *gorm.DB) error {
//	        txCourseRepo := &CourseRepository{
//	            GenericRepository: GenericRepository[models.Course]{
//	                dbRead:  txDB.WithContext(ctx),
//	                dbWrite: txDB.WithContext(ctx),
//	                table:   rp.table,
//	            },
//	            app: rp.app,
//	        }
//	        return fn(txCourseRepo)
//	    })
//	}
func (rp *GenericRepository[T]) Transactional(ctx context.Context, txDB *gorm.DB, fn func(txDB *gorm.DB) error) error {
	return txDB.WithContext(ctx).Transaction(fn)
}

// NewTransactionRepo creates a new repository instance with transaction connections.
// This is a convenience function for creating transaction-aware repository instances
// to be used within Transaction callbacks.
//
// Parameters:
//   - ctx: Context for the transaction
//   - txDB: The transaction DB connection
//   - tableName: The table name for this repository
//
// Returns:
//   - GenericRepository[T]: A new repository instance configured for the transaction
//
// Usage:
//
//	func (rp *CourseRepository) Transaction(ctx context.Context, fn func(txRepo *CourseRepository) error) error {
//	    return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//	        txRepo := NewTransactionRepo[models.Course](ctx, tx, rp.table)
//	        domainRepo := &CourseRepository{
//	            GenericRepository: txRepo,
//	            app: rp.app,
//	        }
//	        return fn(domainRepo)
//	    })
//	}
func NewTransactionRepo[T models.IModel](ctx context.Context, txDB *gorm.DB, tableName string) GenericRepository[T] {
	return GenericRepository[T]{
		dbRead:  txDB.WithContext(ctx),
		dbWrite: txDB.WithContext(ctx),
		table:   tableName,
	}
}
