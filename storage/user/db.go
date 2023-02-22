package user

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

func (r *Repository) ByTelegramID(ctx context.Context, telegramID string) (User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return User{}, errors.Wrap(err, "try find user by telegram id")
	}
	return user, nil
}

func (r *Repository) Upsert(ctx context.Context, u User) (User, bool, error) {
	var isUpdated bool

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user, err := r.ByTelegramID(ctx, u.TelegramID)

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			err = tx.Clauses(clause.OnConflict{UpdateAll: true, Columns: []clause.Column{{Name: "id"}}}).Create(u).Error
			if err != nil {
				return errors.Wrap(err, "fail create")
			}
		case err != nil:
			return errors.Wrap(err, "fail fetch")
		default:
			err = tx.Model(user).Updates(u).Error
			if err != nil {
				return errors.Wrap(err, "fail update")
			}
			isUpdated = true
		}
		return nil
	})

	return u, isUpdated, err
}
