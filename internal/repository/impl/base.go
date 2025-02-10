package impl

import (
	"context"
	"errors"
	"reflect"
	pkgErrors "todo/pkg/errors"

	"gorm.io/gorm"
)

// EntityType 实体类型
type EntityType string

const (
	EntityUser     EntityType = "user"
	EntityTodo     EntityType = "todo"
	EntityCategory EntityType = "category"
	EntityReminder EntityType = "reminder"
)

// BaseRepository 基础仓储实现
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓储实例
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// Transaction 事务处理
func (r *BaseRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, "tx", tx))
	})
}

// GetDB 获取数据库连接
func (r *BaseRepository) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}

// Create 通用创建方法
func (r *BaseRepository) Create(ctx context.Context, model interface{}) error {
	return r.GetDB(ctx).Create(model).Error
}

// Update 通用更新方法
func (r *BaseRepository) Update(ctx context.Context, model interface{}) error {
	return r.GetDB(ctx).Save(model).Error
}

// Delete 通用删除方法
func (r *BaseRepository) Delete(ctx context.Context, model interface{}) error {
	// 获取模型的主键值
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// 获取主键字段
	idField := value.FieldByName("ID")
	if !idField.IsValid() {
		return errors.New("模型缺少 ID 字段")
	}

	// 使用主键作为删除条件
	return r.GetDB(ctx).Where("id = ?", idField.Interface()).Delete(model).Error
}

// GetByID 通用根据ID获取方法
func (r *BaseRepository) GetByID(ctx context.Context, id uint, model interface{}) error {
	return r.GetDB(ctx).First(model, id).Error
}

// List 通用列表查询方法
func (r *BaseRepository) List(ctx context.Context, offset, limit int, models interface{}, conditions ...interface{}) error {
	db := r.GetDB(ctx)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	return db.Offset(offset).Limit(limit).Find(models).Error
}

// Count 通用计数方法
func (r *BaseRepository) Count(ctx context.Context, model interface{}, conditions ...interface{}) (int64, error) {
	var count int64
	db := r.GetDB(ctx)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	err := db.Model(model).Count(&count).Error
	return count, err
}

// handleError 通用错误处理函数
func (r *BaseRepository) handleError(err error, entityType EntityType) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		switch entityType {
		case EntityUser:
			return pkgErrors.ErrUserNotFound
		case EntityTodo:
			return pkgErrors.ErrTodoNotFound
		case EntityCategory:
			return pkgErrors.ErrCategoryNotFound
		case EntityReminder:
			return pkgErrors.ErrReminderNotFound
		default:
			return pkgErrors.ErrNotFound
		}
	}

	return err
}

// FindByID 通用查询方法
func (r *BaseRepository) FindByID(ctx context.Context, id uint, entityType EntityType) (interface{}, error) {
	var entity interface{}

	err := r.db.WithContext(ctx).First(&entity, id).Error
	return entity, r.handleError(err, entityType)
}

// GetByField 通用字段查询方法
func (r *BaseRepository) GetByField(ctx context.Context, model interface{}, field string, value interface{}, entityType EntityType) error {
	err := r.GetDB(ctx).Where(field+" = ?", value).First(model).Error
	return r.handleError(err, entityType)
}
