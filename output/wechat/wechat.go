package wechat

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/output"
	"io"
	"os"
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
	Client.MessageHandler = DealMessage

	// 注册登陆二维码回调
	Client.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if !oss.Exists(HotLoginFilename) {
		os.Create(HotLoginFilename)
	}
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
			if name, ok := strings.CutPrefix(out, output.TypeWechatGroup+":"); ok {
				//给群组发送消息
				logs.Tracef("给群组[%s]发送消息[%s]\n", name, msg.Content)

				mu.RLock()
				group, ok := Groups[name]
				mu.RUnlock()
				if !ok {
					groups, err := Self.Groups(true)
					if err != nil {
						logs.Warnf("给群组[%s]发送消息错误： %v\n", name, err)
						return
					}

					mu.Lock()
					Groups = map[string]*openwechat.Group{}
					for _, v := range groups {
						logs.Debug(v.NickName)
						Groups[v.NickName] = v
						if v.NickName == name {
							group = v
						}
					}
					mu.Unlock()
				}

				if group == nil {
					logs.Warnf("给好友[%s]发送消息错误： 群组不存在\n", name)
					return
				}

				_, err := group.SendText(msg.Content)
				if err != nil {
					logs.Warnf("给群组[%s]发送消息错误： %v\n", name, err)
				}

			} else if name, ok := strings.CutPrefix(out, output.TypeWechatFriend+":"); ok {
				//给好友发送信息
				logs.Tracef("给好友[%s]发送消息[%s]\n", name, msg.Content)

				mu.RLock()
				friend, ok := Friends[name]
				mu.RUnlock()
				if !ok {
					friends, err := Self.Friends()
					if err != nil {
						logs.Warnf("给好友[%s]发送消息错误： %v\n", name, err)
						return
					}

					mu.Lock()
					Friends = map[string]*openwechat.Friend{}
					for _, v := range friends {
						Friends[v.NickName] = v
						if v.NickName == name {
							friend = v
						}
					}
					mu.Unlock()
				}

				if friend == nil {
					logs.Warnf("给好友[%s]发送消息错误： 好友不存在\n", name)
					return
				}

				_, err := friend.SendText(msg.Content)
				if err != nil {
					logs.Warnf("给好友[%s]发送消息错误： %v\n", name, err)
				}

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
