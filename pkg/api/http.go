package api

import (
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/mux"
	"github.com/injoyai/notice/pkg/push"
	"github.com/injoyai/notice/pkg/user"
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
		g.ALL("/notice", func(r *mux.Request) {
			u := r.GetCache("user").Val().(*user.User)
			msg := &push.Message{}
			r.Parse(msg)
			//加入发送队列
			err := push.Manager.Push(u, msg)
			in.Err(err)
		})

		//查询用户列表
		g.GET("/user/all", func(r *mux.Request) {
			data, err := user.GetAll()
			in.CheckErr(err)
			in.Succ(data)
		})

		//添加/修改用户
		g.POST("/user", func(r *mux.Request) {
			req := new(user.User)
			r.Parse(req)
			err := user.Create(req)
			in.Err(err)
		})

		//删除用户
		g.DELETE("/user", func(r *mux.Request) {
			username := r.GetString("username")
			err := user.Del(username)
			in.Err(err)
		})

	})

	return s.Run()
}
