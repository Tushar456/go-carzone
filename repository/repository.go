package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository is a generic repository providing basic CRUD operations.
type Repository[T any] struct {
	db *gorm.DB
}

// New creates a new generic repository.
func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

// Get finds a single record matching the given condition.
func (r *Repository[T]) Get(ctx context.Context, dest *T, conds ...interface{}) error {
	return r.db.WithContext(ctx).First(dest, conds...).Error
}

// GetWithPreload finds a single record with preloaded associations.
func (r *Repository[T]) GetWithPreload(ctx context.Context, dest *T, preloads []string, conds ...interface{}) error {
	query := r.db.WithContext(ctx)
	for _, p := range preloads {
		query = query.Preload(p)
	}
	return query.First(dest, conds...).Error
}

// Create inserts a new record into the database.
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update saves an existing record in the database.
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete removes a record from the database.
func (r *Repository[T]) Delete(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Delete(entity).Error
}

// Find finds records matching the given condition.
func (r *Repository[T]) Find(ctx context.Context, dest *[]T, conds ...interface{}) error {
	return r.db.WithContext(ctx).Find(dest, conds...).Error
}

// FindWithPreload finds records with preloaded associations.
func (r *Repository[T]) FindWithPreload(ctx context.Context, dest *[]T, preloads []string, conds ...interface{}) error {
	query := r.db.WithContext(ctx)
	for _, p := range preloads {
		query = query.Preload(p)
	}
	return query.Find(dest, conds...).Error
}
