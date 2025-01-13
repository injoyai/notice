package wechat

import (
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/logs"
	"github.com/injoyai/notice/pkg/push"
	"github.com/skip2/go-qrcode"
	"io"
	"os"
	"path/filepath"
	"runtime"
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
	// 注册登陆二维码回调 openwechat.PrintlnQrcodeUrl
	w.Client.UUIDCallback = w.PrintQrcode
	// 创建文件
	oss.NewNotExist(HotLoginFilename, "")
	// 登陆
	err := w.Client.HotLogin(openwechat.NewFileHotReloadStorage(HotLoginFilename))
	if err != nil {
		switch err.Error() {
		case io.EOF.Error():
			err = w.Client.Login()
		case "invalid storage", "cookie invalid", "failed login check":
			f, er := os.Create(HotLoginFilename)
			if er == nil {
				f.Truncate(0)
				f.Close()
			}
			err = w.Client.HotLogin(openwechat.NewFileHotReloadStorage(HotLoginFilename))
			logs.Debug(err)
		}
	}
	if err != nil {
		return nil, err
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

func (this *Wechat) Name() string {
	return "微信"
}

func (this *Wechat) Types() []string {
	return []string{push.TypeWechatGroup, push.TypeWechatFriend}
}

func (this *Wechat) Push(msg *push.Message) (err error) {
	if this.Client == nil || this.Self == nil {
		return errors.New("初始化失败")
	}

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

func (this *Wechat) PrintQrcode(uuid string) {
	switch runtime.GOOS {
	case "windows":
		openwechat.PrintlnQrcodeUrl(uuid)
		//return
	}

	url := "https://login.weixin.qq.com/l/" + uuid
	//fmt.Println(url)
	qr, err := qrcode.New(url, qrcode.Medium)
	if err == nil {
		s := ""
		for _, line := range qr.Bitmap() {
			for _, r := range line {
				if r {
					s += "   "
				} else {
					s += "███"
				}
			}
			s += "\n"
		}
		fmt.Println(s)
	}

}
