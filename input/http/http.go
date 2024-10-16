package http

import (
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/mux"
	"github.com/injoyai/notice/output"
)

var Server = mux.New()

func Init(port int) error {

	s := Server.SetPort(port)
	s.Group("/api", func(g *mux.Grouper) {

		g.Middle(func(r *mux.Request) {
			token := r.GetHeader("Authorization")
			_ = token
		})

		g.POST("/notice", func(r *mux.Request) {
			msg := &output.Message{}
			r.Parse(msg)
			output.Trunk.Do(&output.Message{})
			in.Succ(nil)
		})

	})
	return s.Run()
}
