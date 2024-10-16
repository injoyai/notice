package output

import "github.com/injoyai/base/chans"

var Trunk = chans.NewTrunk(20)

// Message 消息格式
type Message struct {
	Output  []string       `json:"output"`  //输出方式(wechat:group:群名)
	Title   string         `json:"title"`   //消息标题
	Content string         `json:"content"` //消息内容
	Param   map[string]any `json:"param"`   //其它参数

	Back func(out string, err error) `json:"-"` //回调
}
