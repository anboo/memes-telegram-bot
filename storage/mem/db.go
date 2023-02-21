package mem

import (
	"context"

	"memes-bot/storage/user"

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

func (r *Repository) UpsertMem(ctx context.Context, mem Mem) (Mem, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var m Mem
		err := tx.Where("id = ?", mem.ID).First(&m).Error

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			err = tx.Clauses(clause.OnConflict{UpdateAll: true, Columns: []clause.Column{{Name: "external_id"}, {Name: "source"}}}).Create(mem).Error
			if err != nil {
				return errors.Wrap(err, "fail create")
			}
		case err != nil:
			return errors.Wrap(err, "fail fetch")
		default:
			err = tx.Where("id = ?", mem.ID).Updates(mem).Error
			if err != nil {
				return errors.Wrap(err, "fail update")
			}
		}
		return nil
	})

	return mem, err
}

func (r *Repository) ReserveNewMem(ctx context.Context, user user.User, mem Mem) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&ReservedMemUser{
			MemID:  mem.ID,
			UserID: user.ID,
		}).Error
	})

	if err != nil {
		return errors.Wrap(err, "reserve mem")
	}

	return nil
}
