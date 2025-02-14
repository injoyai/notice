package telegram

import (
	"errors"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/push"
	"time"
)

func New(token, proxy, defaultChatID string) (*Telegram, error) {
	if token == "" {
		return &Telegram{}, errors.New("token不能为空")
	}
	c := http.NewClient()
	c.SetTimeout(0)
	c.SetProxy(proxy)
	bot, err := api.NewBotAPIWithClient("7682330717:AAEEJkwa5tIjB21c3ieb7SNKoppyqvkKuwI", api.APIEndpoint, c.Client)
	if err != nil {
		return nil, err
	}
	go func() {

		u := api.NewUpdate(0)
		u.Timeout = 60
		for update := range bot.GetUpdatesChan(u) {
			if update.Message != nil { // 检查是否有新消息

				switch update.Message.Text {
				case "/chat_id", "/chatid":
					// 获取用户的 chat_id
					msg := api.NewMessage(update.Message.Chat.ID, "Hello! Your chat_id is: "+conv.String(update.Message.Chat.ID))
					msg.ReplyToMessageID = update.Message.MessageID
					_, err := bot.Send(msg)
					logs.PrintErr(err)

				case "/time":
					msg := api.NewMessage(update.Message.Chat.ID, time.Now().Format("15:04:05"))
					msg.ReplyToMessageID = update.Message.MessageID
					_, err := bot.Send(msg)
					logs.PrintErr(err)

				}
			}
		}
	}()
	return &Telegram{
		Bot:           bot,
		DefaultChatID: defaultChatID,
	}, nil
}

type Telegram struct {
	Bot           *api.BotAPI
	DefaultChatID string //默认消息id
}

func (this *Telegram) Name() string { return "Telegram" }

func (this *Telegram) Types() []string { return []string{push.TypeTelegram} }

func (this *Telegram) Push(msg *push.Message) error {
	if this.Bot == nil {
		return errors.New("telegram未初始化")
	}
	if msg.Target == "" {
		msg.Target = this.DefaultChatID
	}
	_, err := this.Bot.Send(api.NewMessage(conv.Int64(msg.Target), msg.Content))
	return err
}
