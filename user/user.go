package user

import (
	"errors"
	"github.com/injoyai/base/g"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/minidb"
	"github.com/injoyai/notice/oem"
	"time"
)

var (
	DB = minidb.New("./data/database/", "default", func(db *minidb.DB) {
		minidb.WithID("ID")
	})
	Cache = maps.NewSafe()
)

func Init() error {
	if err := DB.Sync(new(User)); err != nil {
		return err
	}
	data, err := All()
	if err != nil {
		return err
	}
	for _, v := range data {
		Cache.Set(v.Username, v)
	}
	return nil
}

type LoginReq struct {
	Username  string `json:"username"`  //用户名
	Signal    string `json:"signal"`    //签名
	Timestamp int64  `json:"timestamp"` //时间戳
}

type User struct {
	ID       string   `json:"id"`       //主键
	Name     string   `json:"name"`     //名称
	Username string   `json:"username"` //账号
	Password string   `json:"password"` //密码
	Limit    []string `json:"limit"`    //消息推送限制
}

func (this *User) LimitMap() map[string]struct{} {
	m := make(map[string]struct{})
	for _, v := range this.Limit {
		m[v] = struct{}{}
	}
	return m
}

func CheckToken(token string) (*User, error) {
	if oem.IsSuperToken(token) {
		return &User{Username: "super"}, nil
	}
	username, err := oem.GetToken(token)
	if err != nil {
		return nil, err
	}
	return GetByCache(username)
}

func GetByCache(username string) (*User, error) {
	val, ok := Cache.Get(username)
	if !ok {
		return nil, errors.New("用户不存在")
	}
	return val.(*User), nil
}

func Login(req *LoginReq) (string, error) {

	//判断用户是否存在
	user, err := GetByCache(req.Username)
	if err != nil {
		return "", err
	}

	signal := oem.Signal(user.Username, user.Password, time.Unix(req.Timestamp, 0))

	if req.Signal != signal {
		return "", errors.New("验证失败")
	}

	token := g.RandString(16)
	err = oem.SetToken(token, user.Username, time.Hour*24*3)

	return token, err
}

func All() ([]*User, error) {
	data := []*User(nil)
	err := DB.Find(&data)
	return data, err
}

//func List() ([]*User, error) {
//	data := []*User{}
//
//}
