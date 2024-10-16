package output

const (
	TypeTCP          = "tcp"
	TypeTCPDesktop   = "tcp:desktop"
	TypeTCPAndroid   = "tcp:android"
	TypeWechatGroup  = "wechat:group"
	TypeWechatFriend = "wechat:friend"
)

// Message 消息格式
type Message struct {
	Output  []string       `json:"output"`  //输出方式(wechat:group:群名)
	Title   string         `json:"title"`   //消息标题
	Content string         `json:"content"` //消息内容
	Param   map[string]any `json:"param"`   //其它参数
}

func (this *Message) Details() *Details {
	return &Details{
		Title:   this.Title,
		Content: this.Content,
		Param:   this.Param,
	}
}

type Details struct {
	Title   string         `json:"title"`   //消息标题
	Content string         `json:"content"` //消息内容
	Param   map[string]any `json:"param"`   //其它参数
}
