package telegram

import (
	"github.com/injoyai/notice/pkg/push"
	"gopkg.in/telebot.v4"
)

func New(token string) (*Telegram, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10},
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

func (this *Telegram) Push(msg *push.Message) error {
	_, err := this.Bot.Send(ID(msg.Target), msg.Content)
	return err
}

type ID string

func (id ID) Recipient() string { return string(id) }
