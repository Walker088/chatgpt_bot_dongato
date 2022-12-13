package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	b "github.com/Walker088/chatgpt_bot_dongato/src/bot"
	"github.com/Walker088/chatgpt_bot_dongato/src/config"
)

func main() {
	cfg := config.GetAppConfig()
	bot, err := b.NewBot(cfg)
	if err != nil {
		log.Fatalf("Couldn't start Telegram bot: %v", err)
	}

	osSignal := make(chan os.Signal, 2)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-osSignal
		bot.Stop()
		os.Exit(0)
	}()

	for update := range bot.GetUpdatesChan() {
		go bot.HandleUpdate(update)
	}
}
