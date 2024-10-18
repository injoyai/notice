package tcp

import (
	"encoding/json"
	"github.com/injoyai/conv"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/server"
	"github.com/injoyai/ios/server/listen"
	"github.com/injoyai/notice/output"
	"github.com/injoyai/notice/user"
	"time"
)

var Server *server.Server

func Init(port int, dealMessage func(c *client.Client, msg ios.Acker)) (err error) {
	Server, err = listen.TCP(port, func(s *server.Server) {

		//客户端连接事件
		s.Event.OnConnected = func(s *server.Server, c *client.Client) error {
			c.SetReadTimeout(time.Second)
			bs, err := c.ReadMessage()
			if err != nil {
				return err
			}
			c.SetReadTimeout(0)
			req := new(user.LoginReq)
			if err := json.Unmarshal(bs, req); err != nil {
				return err
			}

			_, err = user.Login(req)
			if er := c.WriteAny(output.Resp{
				ID:      req.ID,
				Success: err == nil,
				Message: conv.String(err),
			}); er != nil {
				return er
			}
			if err != nil {
				return err
			}

			c.SetKey(req.Username)
			return nil
		}
		s.SetClientOption(func(c *client.Client) {
			c.Event.OnDealMessage = dealMessage
		})

	})
	if err != nil {
		return err
	}
	return Server.Run()
}
