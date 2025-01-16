package middle

import (
	"errors"
	"fmt"
	"github.com/injoyai/base/maps/wait/v2"
	"github.com/injoyai/notice/pkg/push"
	"time"
)

/*
消息队列中间件
*/

func NewQueue(limit, _cap int, timeout time.Duration) *Queue {
	q := &Queue{
		c:    make(chan func() error, _cap),
		wait: wait.New(timeout),
	}

	for i := 0; i < limit; i++ {
		go func() {
			for {
				select {
				case f := <-q.c:
					err := f()
					q.wait.Done(fmt.Sprintf("%p", f), err)
				}
			}
		}()
	}

	return q
}

type Queue struct {
	c    chan func() error
	wait *wait.Entity
}

func (this *Queue) Handler(msg *push.Message, u push.User, f func() error) error {
	select {
	case this.c <- f:
		_, err := this.wait.Wait(fmt.Sprintf("%p", f))
		return err
	case <-time.After(time.Second):
		return errors.New("消息队列已满")
	}
}
