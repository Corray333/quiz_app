package app

import (
	"os"

	"github.com/Corray333/quiz/internal/repository"
	"github.com/Corray333/quiz/internal/service"
	"github.com/Corray333/quiz/internal/telegram"
	"github.com/Corray333/quiz/internal/transport"
)

type App struct {
	tg     *telegram.TelegramClient
	server *transport.Server
}

func New() *App {
	store := repository.New()

	service := service.NewService(store)

	return &App{
		tg:     telegram.NewClient(os.Getenv("BOT_TOKEN"), service),
		server: transport.NewServer(service),
	}
}

func (app *App) Run() {
	go app.tg.Run()
	app.server.Run()
}
