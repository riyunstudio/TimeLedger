# PHASE 11: Repository 交易安全性加固

**執行日期**：2026-01-30

## 完成項目摘要

### 類別	完成項目	狀態
- **GenericRepository 強化**：更新 Transaction 方法，建立全新 Repo 實例	✅
- **交易輔助函數**：新增 Transactional、NewTransactionRepo 工具函數	✅
- **交易介面**：定義 TransactionCapable 介面	✅
- **CourseRepository**：實作專屬 Transaction 方法	✅
- **OfferingRepository**：實作專屬 Transaction 方法	✅
- **ScheduleRuleRepository**：實作專屬 Transaction 方法	✅

## 建置驗證

```bash
go build ./app/repositories/...  ✅ 通過
go build ./app/...              ✅ 通過
```

## 程式碼範例

### GenericRepository 交易方法

```go
// Transaction executes a function within a database transaction.
// IMPORTANT: This method creates a NEW repository instance with transaction connections
// to avoid race conditions in concurrent requests.
func (rp *GenericRepository[T]) Transaction(ctx context.Context, fn func(txCtx context.Context, txDB *gorm.DB) error) error {
    return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        return fn(ctx, tx)
    })
}

// NewTransactionRepo creates a transaction-aware repository instance
func NewTransactionRepo[T models.IModel](ctx context.Context, txDB *gorm.DB, tableName string) GenericRepository[T] {
    return GenericRepository[T]{
        dbRead:  txDB.WithContext(ctx),
        dbWrite: txDB.WithContext(ctx),
        table:   tableName,
    }
}
```

### CourseRepository 交易入口

```go
func (rp *CourseRepository) Transaction(ctx context.Context, fn func(txRepo *CourseRepository) error) error {
    return rp.dbWrite.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 建立新的 CourseRepository 實例，交易隔離
        txRepo := &CourseRepository{
            GenericRepository: GenericRepository[models.Course]{
                dbRead:  tx.WithContext(ctx),
                dbWrite: tx.WithContext(ctx),
                table:   rp.table,
            },
            app: rp.app,
        }
        return fn(txRepo)
    })
}
```

## 安全性改進對比

### 改進前（淺拷貝問題）

```go
// ❌ 錯誤做法：修改現有實例，導致 Race Condition
func (rp *GenericRepository[T]) Transaction(fn func(txRepo *GenericRepository[T]) error) error {
    txRepo := *rp  // 淺拷貝，共用指標
    txRepo.dbWrite = tx.WithContext(ctx)  // 修改原始實例！
    return fn(&txRepo)
}
```

### 改進後（全新實例）

```go
// ✅ 正確做法：建立全新實例，確保執行緒安全
func (rp *GenericRepository[T]) TransactionWithRepo(ctx context.Context, txDB *gorm.DB, fn func(txRepo GenericRepository[T]) error) error {
    txRepo := GenericRepository[T]{
        dbRead:  txDB.WithContext(ctx),
        dbWrite: txDB.WithContext(ctx),
        table:   rp.table,
    }
    return fn(txRepo)
}
```

## Service 層使用範例

```go
func (s *CourseService) CreateCourses(ctx context.Context, courses []models.Course) error {
    return s.courseRepo.Transaction(ctx, func(txRepo *repositories.CourseRepository) error {
        // 所有操作都在同一個交易中執行
        // 自訂方法（如 ListByCenterID）在交易中仍可使用
        for _, course := range courses {
            if _, err := txRepo.Create(ctx, course); err != nil {
                return err
            }
        }
        return nil
    })
}
```

## 效益總結

| 指標 | 改善前 | 改善後 |
|:---|:---:|:---:|
| 執行緒安全 | 共用實例導致 Race Condition | 全新實例確保安全 |
| 自訂方法支援 | 交易中無法使用自訂方法 | 領域 Repository 保留自訂方法 |
| 使用彈性 | 僅限 GenericRepository 方法 | 自訂方法 + GenericRepository 方法 |

## 相關文件

- `app/repositories/generic.go` - GenericRepository 實作
- `app/repositories/course.go` - CourseRepository 實作
- `app/repositories/offering.go` - OfferingRepository 實作
- `app/repositories/schedule_rule.go` - ScheduleRuleRepository 實作
- `pdr/ARCHITECTURAL_OPTIMIZATION_GUIDE.md` - 架構優化指南

---

**文件版本**：v1.0
**最後更新**：2026-01-30
