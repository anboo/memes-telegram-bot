package user_source

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
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, us UserSource) error {
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&us).Error
	if err != nil {
		return errors.Wrap(err, "fail create user source")
	}
	return nil
}

func (r *Repository) ByUserId(ctx context.Context, userId string) (res []UserSource, err error) {
	err = r.db.WithContext(ctx).Find(&res, "user_id = ?", userId).Error
	if err != nil {
		return res, errors.Wrap(err, "find user source by user id")
	}
	return res, nil
}
