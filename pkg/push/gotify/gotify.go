package gotify

import (
	"fmt"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/notice/pkg/push"
)

func New(host, token string, priority int) *Gotify {
	return &Gotify{
		Host:     host,
		Token:    token,
		Priority: priority,
	}
}

type Gotify struct {
	Host     string
	Token    string
	Priority int
}

func (this *Gotify) Types() []string {
	return []string{push.TypeGotify}
}

func (this *Gotify) Push(msg *push.Message) error {
	return http.Url(fmt.Sprintf("%s/message?token=%s", this.Host, this.Token)).SetBody(g.Map{
		"title":    msg.Title,
		"message":  msg.Content,
		"priority": this.Priority,
	}).Post().Err()
}
