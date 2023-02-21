package mem

import (
	"context"
	"database/sql"

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

func (r *Repository) Find(ctx context.Context, id string) (Mem, error) {
	var res Mem

	err := r.db.WithContext(ctx).First(&res, "id = ?", id).Error
	if err != nil {
		return Mem{}, errors.Wrap(err, "try find mem rep")
	}

	return res, nil
}

func (r *Repository) FindRelevantMemForUser(ctx context.Context, u user.User) (Mem, error) {
	var res Mem

	err := r.db.Raw(`
		WITH my_memes AS (
			SELECT v1.mem_id, v1.vote FROM votes v1 WHERE v1.user_id = @userId
		), mutual_users AS (
			SELECT u.id AS user_id,
				   (
					   SELECT COUNT(*)
					   FROM votes v2
					   WHERE v2.user_id = u.id
						 AND v2.mem_id || '_' || v2.vote IN (SELECT mm.mem_id || '_' || mm.vote FROM my_memes mm)
				   ) as count_mutual_votes
			FROM users u
			WHERE u.id != @userId
			  AND EXISTS (
					SELECT v2.mem_id
					FROM votes v2
					WHERE v2.user_id = u.id
					  AND v2.mem_id IN (SELECT mm.mem_id FROM my_memes mm)
				)
			ORDER BY 2 DESC
		)
		SELECT m.*
		FROM votes v
				 JOIN memes m on m.id = v.mem_id
				 JOIN mutual_users mu ON mu.user_id = v.user_id
		WHERE v.user_id != @userId
		  AND v.vote > 0
		  AND m.id NOT IN (
			SELECT rmu.mem_id FROM reserved_mem_users rmu WHERE rmu.user_id = @userId
		)
		AND EXISTS (
			SELECT v2.mem_id
			FROM votes v2
			WHERE v2.user_id = v.user_id
			  AND v2.mem_id IN (SELECT mm.mem_id FROM my_memes mm)
		)
		ORDER BY mu.count_mutual_votes DESC, v.vote DESC, v.created_at DESC, m.id DESC
		LIMIT 1
	`, sql.Named("userId", u.ID)).Scan(&res).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Mem{}, errors.Wrap(err, "find relevant meme")
	}

	if res.ID != "" {
		return res, nil
	}

	err = r.db.WithContext(ctx).Order("RANDOM()").First(&res).Error
	if err != nil {
		return Mem{}, errors.Wrap(err, "find relevant mem rep")
	}
	return res, nil
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

func (r *Repository) UpdateRating(ctx context.Context, memId string, diff int) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(Mem{}).Where("id = ?", memId).UpdateColumn("rating", gorm.Expr("rating + ?", diff)).Error
	})
	if err != nil {
		return errors.Wrap(err, "update rating")
	}
	return nil
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
