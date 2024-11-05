package wechat

import (
	"context"
	"errors"
	"github.com/eatmoreapple/openwechat"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/output"
	"io"
	"os"
	"path/filepath"
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

func Init(dir string) (err error) {

	HotLoginFilename = filepath.Join(dir, "data/cache/wechat_hot_login")

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

	//获取当前用户信息
	Self, err = Client.GetCurrentUser()
	if err != nil {
		return err
	}

	output.Trunk.Subscribe(func(ctx context.Context, msg *output.Message) {
		msg.Listen(map[string]func(name string, msg *output.Message) error{
			output.TypeWechatGroup: func(name string, msg *output.Message) error {
				//给群组发送消息
				logs.Tracef("给群组[%s]发送消息[%s]\n", name, msg.Content)
				defer func() {
					if err != nil {
						logs.Warnf("给群组[%s]发送消息错误： %v\n", name, err)
					}
				}()

				mu.RLock()
				group, ok := Groups[name]
				mu.RUnlock()
				if !ok {
					groups, err := Self.Groups(true)
					if err != nil {
						return err
					}

					mu.Lock()
					Groups = map[string]*openwechat.Group{}
					for _, v := range groups {
						Groups[v.NickName] = v
						if v.NickName == name {
							group = v
						}
					}
					mu.Unlock()
				}
				if group == nil {
					return errors.New("群组不存在")
				}

				_, err = group.SendText(msg.Content)
				return err
			},
			output.TypeWechatFriend: func(name string, msg *output.Message) (err error) {
				//给好友发送消息
				logs.Tracef("给好友[%s]发送消息[%s]\n", name, msg.Content)
				defer func() {
					if err != nil {
						logs.Warnf("给好友[%s]发送消息错误： %v\n", name, err)
					}
				}()

				mu.RLock()
				friend, ok := Friends[name]
				mu.RUnlock()
				if !ok {
					friends, err := Self.Friends()
					if err != nil {
						return err
					}

					mu.Lock()
					Friends = map[string]*openwechat.Friend{}
					for _, v := range friends {
						if len(v.RemarkName) > 0 {
							Friends[v.RemarkName] = v
							if v.RemarkName == name {
								friend = v
							}
							continue
						}
						Friends[v.NickName] = v
						if v.NickName == name {
							friend = v
						}
					}
					mu.Unlock()
				}

				if friend == nil {
					return errors.New("好友不存在")
				}

				_, err = friend.SendText(msg.Content)
				return err
			},
		})

	})

	return nil
}
