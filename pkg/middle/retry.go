package middle

import "github.com/injoyai/notice/pkg/push"

func NewRetry(count int) *Retry {
	return &Retry{Count: count}
}

type Retry struct {
	Count int
}

func (this *Retry) Handler(u push.User, msg *push.Message, next func() error) (err error) {
	err = next()
	for i := 0; err != nil && i < this.Count; i++ {
		err = next()
	}
	return
}
