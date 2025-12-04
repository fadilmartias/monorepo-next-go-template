package repositories

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/responses"
	"gorm.io/gorm"
)

// BaseRepository adalah interface generik untuk operasi CRUD
type BaseRepository[T any] interface {
	WithTx(tx *gorm.DB) BaseRepository[T]
	WithLimit(limit int) QueryOption[T]
	WithOffset(offset int) QueryOption[T]
	WithPreload(field string, fn ...func(*gorm.DB) *gorm.DB) QueryOption[T]
	WithOrder(order string) QueryOption[T]
	WithCondition(query any, args ...any) QueryOption[T]
	GetAll(page int, pageSize int, opts ...QueryOption[T]) ([]T, *responses.Pagination, error)
	FindByID(id any, opts ...QueryOption[T]) (*T, error)
	Create(model *T) error
	Save(model *T) error
	UpdateByID(id string, updates map[string]any) (int64, error)
	UpdateByModel(model *T) (int64, error)
	UpdateWithModel(model *T, updates map[string]any) (int64, error)
	UpdateWithExistingModel(existing *T, updated T) (int64, error)
	DeleteByID(id string) (int64, error)
	DeleteByModel(model *T) (int64, error)
}

// baseRepository adalah implementasi BaseRepository dengan Go Generics
type baseRepository[T any] struct {
	DB *gorm.DB
}

// QueryOption definisi fungsi untuk modifikasi query
type QueryOption[T any] func(*gorm.DB) *gorm.DB

// NewBaseRepository membuat instance baru dari baseRepository
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{DB: db}
}

func (r *baseRepository[T]) WithTx(tx *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{DB: tx}
}

// ---- Helper options ----

// WithLimit menambahkan LIMIT
func (r *baseRepository[T]) WithLimit(limit int) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		if limit > 0 {
			return db.Limit(limit)
		}
		return db
	}
}

// WithOffset menambahkan OFFSET
func (r *baseRepository[T]) WithOffset(offset int) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		if offset > 0 {
			return db.Offset(offset)
		}
		return db
	}
}

// WithPreload preload relasi
func (r *baseRepository[T]) WithPreload(field string, fn ...func(*gorm.DB) *gorm.DB) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		if len(fn) > 0 {
			return db.Preload(field, fn[0])
		}
		return db.Preload(field)
	}
}

// WithOrder untuk sorting
func (r *baseRepository[T]) WithOrder(order string) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		if order != "" {
			return db.Order(order)
		}
		return db
	}
}

// WithCondition untuk where custom
func (r *baseRepository[T]) WithCondition(query any, args ...any) QueryOption[T] {
	return func(db *gorm.DB) *gorm.DB {
		if query != nil {
			return db.Where(query, args...)
		}
		return db
	}
}

// Create menyimpan model baru ke database
func (r *baseRepository[T]) Create(model *T) error {
	return r.DB.Create(model).Error
}

func (r *baseRepository[T]) CreateInBatches(models *[]T, batchSize int) error {
	return r.DB.CreateInBatches(models, batchSize).Error
}

// Save menyimpan model baru ke database
func (r *baseRepository[T]) Save(model *T) error {
	return r.DB.Save(model).Error
}

// FindByID mencari model berdasarkan ID
func (r *baseRepository[T]) FindByID(id any, opts ...QueryOption[T]) (*T, error) {
	db := r.DB

	// apply semua opsi ke query
	for _, opt := range opts {
		db = opt(db)
	}

	var model T
	err := db.Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// GetAll dengan opsi (limit, offset, preload, filter, dll)
func (r *baseRepository[T]) GetAll(page int, pageSize int, opts ...QueryOption[T]) ([]T, *responses.Pagination, error) {
	db := r.DB

	// apply semua opsi ke query
	for _, opt := range opts {
		db = opt(db)
	}

	// hitung total items dulu
	var totalItems int64
	if err := db.Model(new(T)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	// pagination
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var models []T
	if err := db.Limit(pageSize).Offset(offset).Find(&models).Error; err != nil {
		return nil, nil, err
	}

	totalPages := (totalItems + int64(pageSize) - 1) / int64(pageSize) // ceil division
	from := offset + 1
	to := offset + len(models)

	pagination := &responses.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
		HasMore:    page < int(totalPages),
		From:       from,
		To:         to,
	}

	return models, pagination, nil
}

func (r *baseRepository[T]) UpdateByID(id string, updates map[string]any) (int64, error) {
	tx := r.DB.Model(new(T)).Where("id = ?", id).Updates(updates)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return tx.RowsAffected, nil
}

func (r *baseRepository[T]) UpdateByModel(model *T) (int64, error) {
	tx := r.DB.Save(model)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return tx.RowsAffected, nil
}

func (r *baseRepository[T]) UpdateWithModel(model *T, updates map[string]any) (int64, error) {
	tx := r.DB.Model(model).Updates(updates)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return tx.RowsAffected, nil
}

func (r *baseRepository[T]) UpdateWithExistingModel(existing *T, updated T) (int64, error) {
	tx := r.DB.Model(existing).Updates(updated)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return tx.RowsAffected, nil
}

func (r *baseRepository[T]) DeleteByID(id string) (int64, error) {
	tx := r.DB.Where("id = ?", id).Delete(new(T))
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return tx.RowsAffected, nil
}

func (r *baseRepository[T]) DeleteByModel(model *T) (int64, error) {
	tx := r.DB.Delete(model)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return tx.RowsAffected, nil
}
