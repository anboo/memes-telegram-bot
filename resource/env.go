package resource

import (
	"github.com/caarlos0/env/v7"
	"github.com/pkg/errors"
)

type Env struct {
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	VkAccessToken    string `env:"VK_ACCESS_TOKEN,required"`
	DsnDB            string `env:"DB_DSN,required"`
	VkGroups         string `env:"VK_GROUPS,required,default=borsch,agil_vk,fuck_humor,in.humour,dzenpub,mhk,dobriememes,dayvinchik,sciencemem"`
}

func (r *Resource) initEnv() error {
	r.Env = Env{}
	err := env.Parse(&r.Env)
	if err != nil {
		return errors.Wrap(err, "parse env")
	}
	return nil
}
