package wechat

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/injoyai/notice/output"
	"io"
	"strings"
	"sync"
)

var (
	HotLoginFilename = "./data/cache/wechat_hot_login"
	Client           = openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	Self             *openwechat.Self
	Groups           = map[string]*openwechat.Group{}
	Friends          = map[string]*openwechat.Friend{}
	mu               sync.RWMutex
)

func Init() (err error) {

	// 注册消息处理函数
	Client.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}

	// 注册登陆二维码回调
	Client.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := Client.HotLogin(openwechat.NewFileHotReloadStorage(HotLoginFilename)); err != nil {
		if err != io.EOF {
			return err
		}
		if err := Client.Login(); err != nil {
			return err
		}
	}

	Self, err = Client.GetCurrentUser()
	if err != nil {
		return err
	}

	output.Trunk.Subscribe(func(ctx context.Context, data interface{}) {
		msg := data.(*output.Message)
		for _, out := range msg.Output {
			if name, ok := strings.CutPrefix(out, "wechat:group:"); ok {
				//给群组发送消息

				mu.RLock()
				group, ok := Groups[name]
				mu.RUnlock()
				if !ok {
					groups, err := Self.Groups()
					if err != nil {
						msg.Back(out, err)
						return
					}

					mu.Lock()
					clear(Groups)
					for _, v := range groups {
						Groups[v.NickName] = v
						if v.NickName == name {
							group = v
						}
					}
					mu.Unlock()
				}

				_, err := group.SendText(msg.Content)
				msg.Back(out, err)

			} else if name, ok := strings.CutPrefix(out, "wechat:friend:"); ok {
				//给好友发送信息

				mu.RLock()
				friend, ok := Friends[name]
				mu.RUnlock()
				if !ok {
					friends, err := Self.Friends()
					if err != nil {
						msg.Back(out, err)
						return
					}

					mu.Lock()
					clear(Friends)
					for _, v := range friends {
						Friends[v.NickName] = v
						if v.NickName == name {
							friend = v
						}
					}
					mu.Unlock()
				}

				_, err := friend.SendText(msg.Content)
				msg.Back(out, err)

			}
		}
	})

	return nil
}

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
