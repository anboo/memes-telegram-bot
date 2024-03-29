package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"memes-bot/cmd"
	"memes-bot/handler"
	"memes-bot/handler/choose_age"
	"memes-bot/handler/choose_sex"
	"memes-bot/handler/config_source"
	"memes-bot/handler/send_mem"
	vote_handler "memes-bot/handler/vote"
	"memes-bot/handler/welcome"
	"memes-bot/importer"
	"memes-bot/resource"
	"memes-bot/storage/mem"
	"memes-bot/storage/user"
	"memes-bot/storage/user_source"
	"memes-bot/storage/vote"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	res, err := resource.Init()
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(res.Env.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	l := zerolog.New(output)
	vk := importer.NewVkImporter(res.Env.VkAccessToken, &l)

	memesRepository := mem.NewRepository(res.DB)
	usersRepository := user.NewRepository(res.DB)
	voteRepository := vote.NewRepository(res.DB)
	userSourceRepository := user_source.NewRepository(res.DB)

	router := handler.NewRouter(
		&l,
		welcome.NewHandler(bot),
		vote_handler.NewHandler(bot, memesRepository, voteRepository, &l),
		choose_sex.NewHandler(bot, usersRepository),
		choose_age.NewHandler(bot, usersRepository),
		config_source.NewHandler(bot, userSourceRepository, res.Env.VkGroups),
		send_mem.NewHandler(bot, memesRepository, &l),
	)

	collector := importer.NewCollector(memesRepository, []importer.Importer{vk})

	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			l.Fatal().Err(err).Msg("cannot start monitoring server")
		}
	}()

	switch arg {
	case "import_memes":
		err := cmd.NewImportMemesCmd(collector).Execute(ctx)
		if err != nil {
			panic(err)
		}
	default:
		err := cmd.NewStartBotCmd(bot, usersRepository, router, &l).Execute(ctx)
		if err != nil {
			panic(err)
		}
	}
}
