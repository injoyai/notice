package pushplus

import (
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/pkg/push"
)

var url = "http://www.pushplus.plus/send"

func New(token string) *PushPlus {
	return &PushPlus{Token: token}
}

type PushPlus struct {
	Token string
}

func (this *PushPlus) Types() []string {
	return []string{push.TypePushPlus}
}

func (this *PushPlus) Push(msg *push.Message) error {
	return http.Url(url).SetBody(g.Map{
		"token":   this.Token,
		"title":   msg.Title,
		"content": msg.Content,
	}).Post().Err()
}
