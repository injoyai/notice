package pushplus

import (
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/output"
)

var url = "http://www.pushplus.plus/send"

func New(token string) *PushPlus {
	return &PushPlus{Token: token}
}

type PushPlus struct {
	Token string
}

func (this *PushPlus) Types() []string {
	return []string{output.TypePushPlus}
}

func (this *PushPlus) Push(msg *output.Message) error {
	return http.Url(url).SetBody(g.Map{
		"token":   this.Token,
		"title":   msg.Title,
		"content": msg.Content,
	}).Post().Err()
}
