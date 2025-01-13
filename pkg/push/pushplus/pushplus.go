package pushplus

import (
	"errors"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/pkg/push"
)

const url = "http://www.pushplus.plus/send"

func New(token string, client ...*http.Client) *PushPlus {
	p := &PushPlus{
		Token:  token,
		client: http.DefaultClient,
	}
	if len(client) > 0 && client[0] != nil {
		p.client = client[0]
	}
	return p
}

type PushPlus struct {
	Token  string
	client *http.Client
}

func (this *PushPlus) Name() string {
	return "推送加"
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
