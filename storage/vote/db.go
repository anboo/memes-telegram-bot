package vote

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, v Vote) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(
			clause.OnConflict{
				UpdateAll: true,
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "mem_id"}},
			}).Create(&v).Error
	})
	if err != nil {
		return errors.Wrap(err, "save vote")
	}
	return nil
}
