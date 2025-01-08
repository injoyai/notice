package push

const (
	TypeAll           = "all"            //全部
	TypeTCP           = "tcp"            //tcp前缀
	TypeDesktopNotice = "desktop:notice" //桌面端通知
	TypeDesktopVoice  = "desktop:voice"  //桌面端语音
	TypeDesktopPopup  = "desktop:popup"  //桌面端弹窗
	TypeWechatGroup   = "wechat:group"   //微信群聊
	TypeWechatFriend  = "wechat:friend"  //微信好友
	TypeMail          = "mail"           //邮件
	TypeAliyunSMS     = "aliyun:sms"     //阿里云短信
	TypeAliyunVoice   = "aliyun:voice"   //阿里云语音
	TypeTencentSMS    = "tencent:sms"    //腾讯云短信
	TypeTencentVoice  = "tencent:voice"  //腾讯云语音
	TypePushPlus      = "pushplus"       //pushplus
	TypeGotify        = "gotify"         //gotify
	TypeWebhook       = "webhook"        //webhook
	TypePlugin        = "plugin"         //插件
	TypeScript        = "script"         //脚本
)

const (
	WinTypeNotice = "notice"
	WinTypeVoice  = "voice"
	WinTypePopup  = "popup"
)

const (
	Text  = "text"  //文本
	Video = "video" //视频
	Image = "image" //图片
	Audio = "audio" //语音
	File  = "file"  //文件
)
