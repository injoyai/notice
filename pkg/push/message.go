package push

// Message 消息格式
type Message struct {
	ID      string         `json:"id"`              //消息id
	Method  string         `json:"method"`          //推送方式(wechat:group)
	Target  string         `json:"target"`          //推送目标(群名)
	Type    string         `json:"type,omitempty"`  //消息类型,默认文本,视频,图片,语音
	Title   string         `json:"title,omitempty"` //消息标题
	Content string         `json:"content"`         //消息内容
	Param   map[string]any `json:"param,omitempty"` //其它参数,可选
	Time    int64          `json:"time"`            //消息时间戳,可选
}
