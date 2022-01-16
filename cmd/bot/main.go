package main

import (
	config2 "github.com/MagnetosUA/TikTok-Music-Bot/pkg/config"
	"github.com/MagnetosUA/TikTok-Music-Bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {

	config, err := config2.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot)
	if err := telegramBot.Start(config); err != nil {
		log.Fatal(err)
	}
}
