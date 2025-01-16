package push

type Pusher interface {
	Name() string            //推送名称
	Types() []string         //推送类型
	Push(msg *Message) error //推送消息
}

type User interface {
	GetID() string             //用户id
	GetName() string           //用户名称
	Limits(method string) bool //用户推送权限
}

type Middle interface {
	Handler(msg *Message, u User, next func() error) error
}

type Handler func(msg *Message, u User) error

type MiddleFunc func(msg *Message, u User, next func() error) error

func (f MiddleFunc) Handler(msg *Message, u User, next func() error) error { return f(msg, u, next) }
