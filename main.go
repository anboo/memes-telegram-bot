package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"memes-bot/cmd"
	"memes-bot/handler"
	"memes-bot/handler/choose_sex"
	"memes-bot/handler/send_mem"
	vote_handler "memes-bot/handler/vote"
	"memes-bot/handler/welcome"
	"memes-bot/importer"
	"memes-bot/storage/mem"
	"memes-bot/storage/user"
	"memes-bot/storage/vote"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	//dsn := os.Getenv("DB_DSN")
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = goose.Up(sqlDB, "./migrations")
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	l := zerolog.New(os.Stdout)
	vk := importer.NewVkImporter(os.Getenv("VK_ACCESS_TOKEN"), &l)

	memesRepository := mem.NewRepository(db)
	usersRepository := user.NewRepository(db)
	voteRepository := vote.NewRepository(db)

	router := handler.NewRouter(
		welcome.NewHandler(bot),
		vote_handler.NewHandler(bot, memesRepository, voteRepository, &l),
		choose_sex.NewHandler(bot, usersRepository),
		send_mem.NewHandler(bot, memesRepository, &l),
	)

	collector := importer.NewCollector(memesRepository, []importer.Importer{vk})

	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

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
