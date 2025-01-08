package output

type Interface interface {
	Types() []string
	Push(msg *Message) error
}
