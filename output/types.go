package output

import (
	"fmt"
	"strings"
)

var TypeMap = map[string]struct{}{
	TypeDesktopNotice: {},
	TypeDesktopVoice:  {},
	TypeDesktopPopup:  {},
	TypeAndroidNotice: {},
	TypeIosNotice:     {},
	TypeWechatGroup:   {},
	TypeWechatFriend:  {},
	TypeAliyunSMS:     {},
	TypeAliyunVoice:   {},
	TypeTencentSMS:    {},
	TypeTencentVoice:  {},
}

const (
	TypeTCP           = "tcp"            //tcp前缀
	TypeDesktopNotice = "desktop:notice" //桌面端通知
	TypeDesktopVoice  = "desktop:voice"  //桌面端语音
	TypeDesktopPopup  = "desktop:popup"  //桌面端弹窗
	TypeAndroidNotice = "android:notice" //安卓通知
	TypeIosNotice     = "ios:notice"     //苹果通知
	TypeWechatGroup   = "wechat:group"   //微信群聊
	TypeWechatFriend  = "wechat:friend"  //微信好友
	TypeAliyunSMS     = "aliyun:sms"     //阿里云短信
	TypeAliyunVoice   = "aliyun:voice"   //阿里云语音
	TypeTencentSMS    = "tencent:sms"    //腾讯云短信
	TypeTencentVoice  = "tencent:voice"  //腾讯云语音
)

const (
	Text  = "text"  //文本
	Video = "video" //视频
	Image = "image" //图片
	Audio = "audio" //语音
	File  = "file"  //文件
)

const (
	WinTypeNotice = "notice"
	WinTypeVoice  = "voice"
	WinTypePopup  = "popup"
)

// Message 消息格式
type Message struct {
	Output  []string       `json:"output"`          //输出方式(wechat:group:群名)
	Type    string         `json:"type,omitempty"`  //消息类型,默认文本,视频,图片,语音
	Title   string         `json:"title,omitempty"` //消息标题
	Content string         `json:"content"`         //消息内容
	Param   map[string]any `json:"param,omitempty"` //其它参数
	Time    int64          `json:"time"`            //
}

func (this *Message) Check(limit map[string]struct{}) error {
	for _, out := range this.Output {
		exist := false
		for k, _ := range limit { //TypeMap
			if strings.HasPrefix(out, k) {
				exist = true
			}
		}
		if !exist {
			return fmt.Errorf("输出方式[%s]不存在", out)
		}
	}
	return nil
}

func (this *Message) Details() *Details {
	return &Details{
		Title:   this.Title,
		Content: this.Content,
		Param:   this.Param,
	}
}

func (this *Message) On(prefix string, f func(name string, msg *Message)) {
	for _, out := range this.Output {
		if out == prefix {
			f("", this)
		} else if name, ok := strings.CutPrefix(out, prefix+":"); ok {
			f(name, this)
		}
	}
}

func (this *Message) Listen(m map[string]func(name string, msg *Message)) {
	for _, out := range this.Output {
		for prefix, f := range m {
			if name, ok := strings.CutPrefix(out, prefix+":"); ok {
				f(name, this)
			}
		}
	}
}

type Details struct {
	Title   string         `json:"title"`   //消息标题
	Content string         `json:"content"` //消息内容
	Param   map[string]any `json:"param"`   //其它参数
}
