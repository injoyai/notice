package pushplus

import (
	"errors"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/pkg/push"
)

const url = "http://www.pushplus.plus/send"

func Push(token string, title, content string) error {
	return New(token).Push(push.NewMessage(title, content))
}

func New(token string, client ...*http.Client) *PushPlus {
	return &PushPlus{
		DefaultToken: token,
		client:       conv.Default(http.DefaultClient, client...),
	}
}

type PushPlus struct {
	DefaultToken string
	client       *http.Client
}

func (this *PushPlus) Name() string {
	return "推送加"
}

func (this *PushPlus) Types() []string {
	return []string{push.TypePushPlus}
}

func (this *PushPlus) Push(msg *push.Message) error {
	token := conv.Select(msg.Target != "", msg.Target, this.DefaultToken)
	if token == "" {
		return errors.New("无效的PushPlus推送Token")
	}
	return this.client.Url(url).SetBody(g.Map{
		"token":   token,
		"title":   msg.Title,
		"content": msg.Content,
	}).Post().Err()
}
