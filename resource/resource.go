package resource

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Resource struct {
	Env Env
	DB  *gorm.DB
}

func Init() (*Resource, error) {
	var (
		r   = &Resource{}
		err error
	)

	if err = r.initEnv(); err != nil {
		return &Resource{}, errors.Wrap(err, "init env")
	}

	if err = r.initDb(); err != nil {
		return &Resource{}, errors.Wrap(err, "init env")
	}

	return r, nil
}
