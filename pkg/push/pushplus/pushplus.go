package pushplus

import (
	"errors"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/pkg/push"
	"time"
)

var url = "http://www.pushplus.plus/send"

func New(token string, timeouts ...time.Duration) *PushPlus {
	return &PushPlus{
		Token:  token,
		client: http.NewClient().SetTimeout(conv.DefaultDuration(0, timeouts...)),
	}
}

type PushPlus struct {
	Token  string
	client *http.Client
}

func (this *PushPlus) Types() []string {
	return []string{push.TypePushPlus}
}

func (this *PushPlus) Push(msg *push.Message) error {
	if this.Token == "" {
		return errors.New("无效的PushPlus推送Token")
	}
	return this.client.Url(url).SetBody(g.Map{
		"token":   this.Token,
		"title":   msg.Title,
		"content": msg.Content,
	}).Post().Err()
}
