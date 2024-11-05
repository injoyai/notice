package output

import (
	"context"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/maps/wait/v2"
	"time"
)

var Trunk = &trunk{
	Trunk: chans.NewTrunk(20),
	wait:  wait.New(time.Second * 2),
}

type trunk struct {
	*chans.Trunk
	wait *wait.Entity
}

func (this *trunk) Subscribe(f func(ctx context.Context, msg *Message)) {
	this.Trunk.Subscribe(func(ctx context.Context, data interface{}) {
		f(ctx, data.(*Message))
	})
}

func (this *trunk) Do(msg *Message) (any, error) {
	msg.resp = func(any any) {
		this.wait.Done(msg.ID, any)
	}
	if err := this.Trunk.Do(msg); err != nil {
		return nil, err
	}
	return this.wait.Wait(msg.ID)
}
