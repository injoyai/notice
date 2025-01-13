package user

import (
	"errors"
	"github.com/injoyai/base/g"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/database/sqlite"
	"github.com/injoyai/goutil/database/xorms"
	"path/filepath"
	"time"
)

const (
	Admin = "admin"
	All   = "all"
)

type IUser interface {
	GetID() string
	GetName() string
}

var (
	DB    *xorms.Engine
	Cache = maps.NewSafe()
)

func Init(dir string) (err error) {
	DB, err = sqlite.NewXorm(filepath.Join(dir, "database/user.db"))
	if err != nil {
		return
	}
	initToken()
	if err = DB.Sync(new(User)); err != nil {
		return err
	}
	data, err := GetAll()
	if err != nil {
		return err
	}
	for _, v := range data {
		Cache.Set(v.Username, v)
	}
	return nil
}

type LoginReq struct {
	ID        string `json:"id"`        //消息id
	Username  string `json:"username"`  //用户名
	Signal    string `json:"signal"`    //签名
	Timestamp int64  `json:"timestamp"` //时间戳
}

type User struct {
	ID       int64    `json:"id"`       //主键
	Name     string   `json:"name"`     //名称
	Username string   `json:"username"` //账号
	Password string   `json:"password"` //密码
	Limit    []string `json:"limit"`    //消息推送限制
}

func (this *User) GetID() string {
	return conv.String(this.ID)
}

func (this *User) GetName() string {
	return this.Name
}

func (this *User) Limits(method string) bool {
	if len(this.Limit) == 1 && this.Limit[0] == All {
		return true
	}
	for _, v := range this.Limit {
		if v == method {
			return true
		}
	}
	return false
}

func CheckToken(token string) (*User, error) {
	if Token.IsSuper(token) {
		return &User{Username: Admin, Limit: []string{All}}, nil
	}
	username, err := Token.Get(token)
	if err != nil {
		return nil, err
	}
	if len(username) == 0 {
		return nil, errors.New("token无效")
	}
	u, err := GetByCache(username)
	if err != nil {
		if err.Error() == "用户不存在" {
			return nil, errors.New("token无效")
		}
		return nil, err
	}
	return u, nil
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

	signal := Signal(user.Username, user.Password, time.Unix(req.Timestamp, 0))

	if req.Signal != signal {
		return "", errors.New("验证失败")
	}

	token := g.RandString(16)
	err = Token.Set(token, user.Username, time.Hour*24*3)

	return token, err
}

func GetAll() ([]*User, error) {
	data := []*User(nil)
	err := DB.Find(&data)
	return data, err
}

func Create(user *User) error {
	if user.Username == "" {
		return errors.New("用户名不能为空")
	}
	if user.Password == "" {
		return errors.New("密码不能为空")
	}
	u, err := GetByCache(user.Username)
	if err != nil && err.Error() != "用户不存在" {
		return err
	}
	if u == nil {
		_, err = DB.Insert(user)
		if err == nil {
			Cache.Set(user.Username, user)
		}
		return err
	}
	_, err = DB.Where("Username=?", user.Username).Update(user)
	if err == nil {
		Cache.Set(user.Username, user)
	}
	return err
}

func Del(username string) error {
	_, err := DB.Where("Username=?", username).Delete(&User{})
	if err == nil {
		Cache.Del(username)
	}
	return err
}
