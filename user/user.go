package user

//import (
//	"github.com/injoyai/minidb"
//)
//
//var (
//	DB = minidb.New
//)

func Init() error {

	return nil
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (this *User) Login(Type string) {

}

func Login(req *LoginReq) (string, error) {

	return "token", nil
}

//func List() ([]*User, error) {
//	data := []*User{}
//
//}
