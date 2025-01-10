package local

import (
	"github.com/go-toast/toast"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/notice"
	"github.com/injoyai/notice/pkg/push"
	"syscall"
	"unsafe"
)

func New() *Local {
	return &Local{}
}

type Local struct{}

func (this *Local) Types() []string {
	return []string{
		push.TypeLocalVideo,
		push.TypeLocalPopup,
		push.TypeLocalNotice,
	}
}

func (this *Local) Push(msg *push.Message) (err error) {
	switch msg.Target {
	case push.TypeLocalPopup:

		//弹窗通知,会阻塞,等待用户关掉才能返回
		user32dll, err := syscall.LoadLibrary("user32.dll")
		if err != nil {
			return err
		}
		user32 := syscall.NewLazyDLL("user32.dll")
		MessageBoxW := user32.NewProc("MessageBoxW")
		_, _, err = MessageBoxW.Call(
			uintptr(0),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msg.Content))),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msg.Title))),
			uintptr(0),
		)
		defer syscall.FreeLibrary(user32dll)
		if err != nil && err.Error() == "The operation completed successfully." {
			return nil
		}

	case push.TypeLocalVideo:
		err = notice.DefaultVoice.Speak(msg.Content)

	case push.TypeLocalNotice:

		//右下角通知
		notification := toast.Notification{
			AppID:    conv.New(msg.Param["AppID"]).String("Microsoft.Windows.Shell.RunDialog"),
			Title:    msg.Title,
			Message:  msg.Content,
			Audio:    toast.Default,
			Duration: toast.Short,
		}
		err = notification.Push()

	}
	return
}
