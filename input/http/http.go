package http

import (
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/mux"
	"github.com/injoyai/notice/output"
	"github.com/injoyai/notice/user"
)

func Init(port int) error {

	s := mux.New().SetPort(port)
	s.Group("/api", func(g *mux.Grouper) {

		//校验权限
		g.Middle(func(r *mux.Request) {
			token := r.GetHeader("Authorization")
			u, err := user.CheckToken(token)
			in.CheckErr(err)
			r.SetCache("user", u)
		})

		//登录
		g.POST("/login", func(r *mux.Request) {
			req := &user.LoginReq{}
			r.Parse(req)
			token, err := user.Login(req)
			in.CheckErr(err)
			in.Succ(token)
		})

		//发送消息
		g.POST("/notice", func(r *mux.Request) {
			u := r.GetCache("user").Val().(*user.User)
			msg := &output.Message{}
			r.Parse(msg)
			err := msg.Check(u.LimitMap())
			in.CheckErr(err)
			output.Trunk.Do(msg)
			in.Succ(nil)
		})

		//查询用户列表

	})
	return s.Run()
}
