package app

import (
	"os"

	"github.com/Corray333/quiz_bot/internal/storage"
	"github.com/Corray333/quiz_bot/internal/telegram"
)

type App struct {
	tg *telegram.TelegramClient
}

func New() *App {
	store := storage.New()

	return &App{
		tg: telegram.NewClient(os.Getenv("BOT_TOKEN"), store),
	}
}

func (app *App) Run() {
	app.tg.Run()
}
