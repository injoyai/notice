package telegram

import (
	"errors"
	"github.com/injoyai/notice/pkg/push"
	"gopkg.in/telebot.v4"
	"time"
)

func New(token string) (*Telegram, error) {
	if token == "" {
		return &Telegram{}, errors.New("token不能为空")
	}
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: time.Second * 3},
	})
	if err != nil {
		return nil, err
	}
	return &Telegram{
		Bot: bot,
	}, nil
}

type Telegram struct {
	Bot *telebot.Bot
}

func (this *Telegram) Name() string { return "Telegram" }

func (this *Telegram) Types() []string { return []string{push.TypeTelegram} }

func (this *Telegram) Push(msg *push.Message) error {
	if this.Bot == nil {
		return errors.New("telegram未初始化")
	}
	_, err := this.Bot.Send(ID(msg.Target), msg.Content)
	return err
}

type ID string

func (id ID) Recipient() string { return string(id) }
