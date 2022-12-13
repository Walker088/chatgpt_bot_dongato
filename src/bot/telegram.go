package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Walker088/chatgpt_bot_dongato/src/chatgpt"
	"github.com/Walker088/chatgpt_bot_dongato/src/config"
	"github.com/Walker088/chatgpt_bot_dongato/src/utils"
)

type Bot struct {
	AllowedUsers []int64
	Api          *tgbotapi.BotAPI
	Timeout      int
	engine       chatgpt.ChatGPT
}

func NewBot(engine chatgpt.ChatGPT, cfg *config.AppConfig) (*Bot, error) {
	botapi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorized on account %s", botapi.Self.UserName)
	botapi.Debug = cfg.BotDebug

	return &Bot{
		AllowedUsers: cfg.AllowdUsers,
		Api:          botapi,
		Timeout:      cfg.BotTimeout,
		engine:       engine,
	}, nil
}

func (b *Bot) GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.Timeout
	return b.Api.GetUpdatesChan(u)
}

func (b *Bot) HandleUpdate(update tgbotapi.Update) {
	var (
		updateText      = update.Message.Text
		updateChatID    = update.Message.Chat.ID
		updateMessageID = update.Message.MessageID
		updateUserID    = update.Message.From.ID
	)
	if len(b.AllowedUsers) != 0 && !utils.Contains(b.AllowedUsers, updateUserID) {
		log.Printf("User %d is not allowed to use the bot", updateUserID)
		var reply = "You are not allowed to use the bot, please contact the provider for further info ヾ(=ﾟ･ﾟ=)ﾉ"
		replyMsg := tgbotapi.NewMessage(updateChatID, reply)
		_, _ = b.Api.Send(replyMsg)
	} else {
		var reply, _ = b.engine.AskQuestion(updateText)
		replyMsg := tgbotapi.NewMessage(updateChatID, string(reply))
		replyMsg.ReplyToMessageID = updateMessageID
		_, _ = b.Api.Send(replyMsg)
	}
}

func (b *Bot) Stop() {
	b.Api.StopReceivingUpdates()
}
