package main

import (
	"github.com/injoyai/notice"
	"github.com/injoyai/notice/pkg/push"
	"github.com/injoyai/notice/pkg/push/webhook"
	"log"
)

func main() {

	target := "test"

	w := webhook.New(map[string]*webhook.Config{
		target: {
			Url:    "http://127.0.0.1:8080/message",
			Method: "",
			Header: nil,
			Body:   `{"title":"${title}","content":"${content}"}`,
			Retry:  0,
			Client: notice.NewClient(),
		},
	})

	push.Manager.Register(w)

	err := push.Manager.Push(&push.Message{
		Method:  push.TypeWebhook,
		Target:  target,
		Title:   "标题",
		Content: "内容",
	}, nil)

	log.Println(err)

}
