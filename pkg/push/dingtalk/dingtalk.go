package dingtalk

import "github.com/injoyai/notice/pkg/push"

//"github.com/alibabacloud-go/dingtalk/robot_1_0"

func New() *DingTalk { return &DingTalk{} }

type DingTalk struct{}

func (this *DingTalk) Types() []string { return []string{push.TypeDingTalk} }

func (this *DingTalk) Push(msg *push.Message) (err error) { return }
