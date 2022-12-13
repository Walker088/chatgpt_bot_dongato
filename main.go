package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Walker088/chatgpt_bot_dongato/src/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, reply string) {
	//text := update.Message.Text
	chatID := update.Message.Chat.ID
	replyMsg := tgbotapi.NewMessage(chatID, reply)
	replyMsg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(replyMsg)
}

func main() {
	cfg := config.GetAppConfig()

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	osSignal := make(chan os.Signal, 2)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-osSignal
		bot.StopReceivingUpdates()
		os.Exit(0)
	}()

	for update := range updates {
		go func(u tgbotapi.Update) {
			var reply = "I am Don Gato, and I am having trouble to get responses from chatgpt"
			handleUpdate(bot, u, reply)
		}(update)
	}
}
