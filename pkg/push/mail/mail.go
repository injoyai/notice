package mail

import (
	"crypto/tls"
	"github.com/injoyai/conv"
	"github.com/injoyai/notice/pkg/push"
	"gopkg.in/gomail.v2"
	"strings"
)

func New(cfg *Config) *Mail {
	if len(cfg.Host) == 0 {
		cfg.Host = "smtp.qq.com"
	}
	if cfg.Port == 0 {
		cfg.Port = 25
	}
	dial := gomail.NewDialer(
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
	)
	dial.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &Mail{
		Dialer: dial,
	}
}

type Mail struct {
	*gomail.Dialer
}

func (this *Mail) Name() string {
	return "邮件"
}

func (this *Mail) Types() []string {
	return []string{push.TypeMail}
}

func (this *Mail) Push(msg *push.Message) error {
	m := gomail.NewMessage()
	m.SetHeader("From", this.Username) // 发件人
	//m.SetHeader("From", "alias"+"<"+userName+">") // 增加发件人别名

	m.SetHeader("To", strings.Split(msg.Target, ",")...) // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	m.SetHeader("Subject", msg.Title)                    // 邮件主题
	m.SetBody("text/html", msg.Content)

	// 抄送，可以多个
	if cc := conv.String(msg.Param["copyTo"]); len(cc) > 0 {
		m.SetHeader("Cc", strings.Split(cc, ",")...)
	}

	// 暗送，可以多个
	if bcc := conv.String(msg.Param["darkTo"]); len(bcc) > 0 {
		m.SetHeader("Bcc", strings.Split(bcc, ",")...)
	}

	// text/html 的意思是将文件的 content-type 设置为 text/html 的形式，浏览器在获取到这种文件时会自动调用html的解析器对文件进行相应的处理。
	// 可以通过 text/html 处理文本格式进行特殊处理，如换行、缩进、加粗等等
	// text/plain的意思是将文件设置为纯文本的形式，浏览器在获取到这种文件时并不会对其进行处理
	// m.SetBody("text/plain", "纯文本")
	// m.Attach("test.sh")   // 附件文件，可以是文件，照片，视频等等
	// m.Attach("lolcatVideo.mp4") // 视频
	// m.Attach("lolcat.jpg") // 照片

	return this.DialAndSend(m)
}

type Config struct {
	// QQ 邮箱：
	// SMTP 服务器地址：smtp.qq.com（SSL协议端口：465/994 | 非SSL协议端口：25）
	// 163 邮箱：
	// SMTP 服务器地址：smtp.163.com（端口：25）
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	// 如果是网易邮箱 pass填密码，qq邮箱填授权码
	Password string `json:"password"`
}
