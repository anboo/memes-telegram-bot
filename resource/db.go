package resource

import (
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (r *Resource) initDb() error {
	db, err := gorm.Open(postgres.Open(r.Env.DsnDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return errors.Wrap(err, "try init db")
	}
	r.DB = db

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = goose.Up(sqlDB, "./migrations")
	if err != nil {
		return errors.Wrap(err, "run migrations")
	}

	return nil
}
