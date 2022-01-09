package telegram

import (
	"fmt"
	"github.com/MagnetosUA/TikTok-Music-Bot/pkg/shazam"
	"github.com/MagnetosUA/TikTok-Music-Bot/pkg/tiktok"
	zk2 "github.com/MagnetosUA/TikTok-Music-Bot/pkg/zk"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: bot}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(tiktok.NewTikTok(), updates)

	return nil
}

func (b *Bot) handleUpdates(tk *tiktok.TikTok, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		var message string

		message, err := tk.GetTikTokMessage(update.Message.Text)
		if err != nil {
			return
		}

		songName, err := tk.GetSongName(shazam.NewShazam(shazam.ShazamResourceURL), message)
		if err != nil {
			return
		}

		if songName == "" {
			songName = "Bad request!"
		}

		var messages []string

		var zk = zk2.NewZK()

		if len(songName) < 2 {
			songName = "Song is not found"
		} else {
			messages, err = zk.DownloadFullSongZK(songName)
			if err != nil {
				return
			}
		}

		fmt.Println("SOUND NAME: ", songName)

		fmt.Sprintf("MESSAGES.len: %v", len(messages))

		if messages != nil {
			for i, mes := range messages {
				if i >= 2 {
					break
				}

				fmt.Println("SOUND URL: ", mes)

				textMessage := tgbotapi.NewMessage(update.Message.Chat.ID, songName)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, mes+"\n")

				b.bot.Send(textMessage)
				b.bot.Send(msg)

				break
			}
		} else {
			textMessage := tgbotapi.NewMessage(update.Message.Chat.ID, songName)
			b.bot.Send(textMessage)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
