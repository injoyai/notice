package desktop

import (
	"encoding/json"
	"errors"
	"github.com/injoyai/conv"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/server"
	"github.com/injoyai/ios/server/listen"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/push"
	"github.com/injoyai/notice/pkg/user"
	"time"
)

func New(port int) (*Desktop, error) {
	s, err := listen.TCP(port, func(s *server.Server) {

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
			if er := c.WriteAny(Resp{
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
		return nil, err
	}
	go s.Run()

	return &Desktop{
		Server: s,
		f: func(name string, msg *push.Message, Type string) error {
			//给桌面端发送消息
			logs.Tracef("给桌面端[%s]发送消息[%s]\n", name, msg.Content)

			c := s.GetClient(name)
			if c == nil {
				logs.Warnf("给桌面端[%s]发送消息错误： 客户端不在线\n", name)
				return errors.New("客户端不在线")
			}

			return c.WriteAny(Req{
				Type:    Type,
				Title:   msg.Title,
				Content: msg.Content,
			})
		},
	}, nil
}

type Desktop struct {
	*server.Server
	f func(name string, msg *push.Message, Type string) error
}

func (this *Desktop) Types() []string {
	return []string{
		push.TypeDesktopNotice,
		push.TypeDesktopVoice,
		push.TypeDesktopPopup,
	}
}

func (this *Desktop) Push(msg *push.Message) (err error) {

	switch msg.Method {
	case push.TypeDesktopNotice:
		err = this.f(msg.Target, msg, push.WinTypeNotice)

	case push.TypeDesktopVoice:
		err = this.f(msg.Target, msg, push.WinTypeVoice)

	case push.TypeDesktopPopup:
		err = this.f(msg.Target, msg, push.WinTypePopup)

	}

	return err
}

func dealMessage(c *client.Client, msg ios.Acker) {

	data := new(push.Message)
	if err := json.Unmarshal(msg.Payload(), data); err != nil {
		logs.Warn("解析消息失败", err)
		return
	}

	u, err := user.GetByCache(c.GetKey())
	if err != nil {
		c.WriteAny(Resp{
			ID:      data.ID,
			Success: false,
			Message: err.Error(),
		})
		return
	}

	err = push.Manager.Push(u, data)

	c.WriteAny(Resp{
		ID:      data.ID,
		Success: err == nil,
		Message: conv.String(err),
	})

}

type Req struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Resp struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
