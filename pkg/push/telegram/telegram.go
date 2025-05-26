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

func Push(token, proxy, chatID, content string) error {
	t, err := _new(token, proxy, chatID)
	if err != nil {
		return err
	}
	return t.Push(push.NewContent(content))
}

func New(token, proxy, chatID string) (*Telegram, error) {
	t, err := _new(token, proxy, chatID)
	if err == nil {
		go t.listen()
	}
	return t, err
}

func _new(token, proxy, defaultChatID string) (*Telegram, error) {
	if token == "" {
		return &Telegram{}, errors.New("token不能为空")
	}
	c := http.NewClient()
	c.SetTimeout(0)
	c.SetProxy(proxy)
	bot, err := api.NewBotAPIWithClient(token, api.APIEndpoint, c.Client)
	if err != nil {
		return nil, err
	}
	return &Telegram{
		Bot:           bot,
		DefaultChatID: defaultChatID,
	}, err
}

type Telegram struct {
	Bot           *api.BotAPI
	DefaultChatID string //默认消息id
}

func (this *Telegram) listen() {
	bot := this.Bot
	u := api.NewUpdate(0)
	u.Timeout = 60
	for update := range bot.GetUpdatesChan(u) {
		if update.Message != nil { // 检查是否有新消息

			switch update.Message.Text {
			case "/chat_id", "/chatid", "chatID":
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
}

func (this *Telegram) Name() string { return "Telegram" }

func (this *Telegram) Types() []string { return []string{push.TypeTelegram} }

func (this *Telegram) Push(msg *push.Message) error {
	if this == nil || this.Bot == nil {
		return errors.New("telegram未初始化")
	}
	target := conv.Select(msg.Target != "", msg.Target, this.DefaultChatID)
	_, err := this.Bot.Send(api.NewMessage(conv.Int64(target), msg.Content))
	return err
}
