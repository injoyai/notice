package plugin

import "github.com/injoyai/notice/pkg/push"

func New() *Plugin {
	return &Plugin{}
}

type Plugin struct{}

func (*Plugin) Types() []string {
	return []string{push.TypePlugin}
}

func (*Plugin) Push(msg *push.Message) error { return nil }
