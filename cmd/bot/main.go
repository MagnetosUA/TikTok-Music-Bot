package main

import (
	"github.com/MagnetosUA/TikTok-Music-Bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	//"io/ioutil"
	"log"
)

type ResponseTikTok struct {
	Html string `json:"html"`
}

func main() {

	bot, err := tgbotapi.NewBotAPI("1675374273:AAHCAYwVhhka8Qf-bYWFGRMViV5t2eZcPAE")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot)
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
