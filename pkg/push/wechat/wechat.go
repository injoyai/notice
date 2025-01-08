package wechat

import (
	"errors"
	"github.com/eatmoreapple/openwechat"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/notice/pkg/push"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func New(cacheDir string) (*Wechat, error) {
	HotLoginFilename := filepath.Join(cacheDir, "wechat_hot_login")
	w := &Wechat{
		Client:  openwechat.DefaultBot(openwechat.Desktop),
		Self:    nil,
		Groups:  make(map[string]*openwechat.Group),
		Friends: make(map[string]*openwechat.Friend),
		mu:      sync.RWMutex{},
	}
	// 注册消息处理函数
	w.Client.MessageHandler = w.DealMessage
	// 注册登陆二维码回调
	w.Client.UUIDCallback = openwechat.PrintlnQrcodeUrl
	// 登陆
	if !oss.Exists(HotLoginFilename) {
		os.Create(HotLoginFilename)
	}
	err := w.Client.HotLogin(openwechat.NewFileHotReloadStorage(HotLoginFilename))
	if err != nil {
		if err != io.EOF {
			if err.Error() == "invalid storage" || err.Error() == "failed login check" {
				os.Remove(HotLoginFilename)
			}
			return nil, err
		}
		if err := w.Client.Login(); err != nil {
			return nil, err
		}
	}

	//获取当前用户信息
	w.Self, err = w.Client.GetCurrentUser()
	if err != nil {
		return nil, err
	}

	return w, nil
}

type Wechat struct {
	Client  *openwechat.Bot
	Self    *openwechat.Self
	Groups  map[string]*openwechat.Group
	Friends map[string]*openwechat.Friend
	mu      sync.RWMutex
}

func (this *Wechat) Types() []string {
	return []string{push.TypeWechatGroup, push.TypeWechatFriend}
}

func (this *Wechat) Push(msg *push.Message) (err error) {

	switch msg.Method {
	case push.TypeWechatGroup:
		this.mu.RLock()
		group, ok := this.Groups[msg.Target]
		this.mu.RUnlock()
		if !ok {
			groups, err := this.Self.Groups(true)
			if err != nil {
				return err
			}
			this.mu.Lock()
			this.Groups = map[string]*openwechat.Group{}
			for _, v := range groups {
				this.Groups[v.NickName] = v
				if v.NickName == msg.Target {
					group = v
				}
			}
			this.mu.Unlock()
		}
		if group == nil {
			return errors.New("群组不存在")
		}
		_, err = group.SendText(msg.Content)

	case push.TypeWechatFriend:

		this.mu.RLock()
		friend, ok := this.Friends[msg.Target]
		this.mu.RUnlock()
		if !ok {
			friends, err := this.Self.Friends()
			if err != nil {
				return err
			}

			this.mu.Lock()
			this.Friends = map[string]*openwechat.Friend{}
			for _, v := range friends {
				this.Friends[v.NickName] = v
				if v.NickName == msg.Target {
					friend = v
				}
			}
			this.mu.Unlock()
		}

		if friend == nil {
			return errors.New("好友不存在")
		}

		_, err = friend.SendText(msg.Content)

	}

	return
}
