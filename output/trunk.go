package output

import (
	"context"
	"github.com/injoyai/base/chans"
)

var Trunk = &trunk{chans.NewTrunk(20)}

type trunk struct {
	*chans.Trunk
}

func (this *trunk) Subscribe(f func(ctx context.Context, msg *Message)) {
	this.Trunk.Subscribe(func(ctx context.Context, data interface{}) {
		f(ctx, data.(*Message))
	})
}
