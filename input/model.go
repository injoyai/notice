package input

// Message 消息格式
type Message struct {
	Token   string         `json:"token"`   //验证信息
	Output  []string       `json:"output"`  //输出方式(wechat:group:群名)
	Title   string         `json:"title"`   //消息标题
	Content string         `json:"content"` //消息内容
	Param   map[string]any `json:"param"`   //其它参数
}
