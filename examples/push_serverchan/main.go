package main

import (
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/mux"
	"github.com/injoyai/notice/pkg/middle"
	"github.com/injoyai/notice/pkg/push"
	"github.com/injoyai/notice/pkg/push/serverchan"
)

func main() {

	push.Manager.Register(
		serverchan.New("xxx"),
	)

	push.Manager.Use(
		middle.NewAuth(), //权限
		middle.NewLog(),  //日志
	)

	s := mux.New().SetPort(8080)
	s.POST("/send", func(r *mux.Request) {
		err := push.Manager.Handler(r.Request, push.NewUser("无"))
		in.Err(err)
	})
	s.Run()

}
