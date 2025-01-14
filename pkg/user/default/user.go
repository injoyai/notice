package user

import (
	"errors"
	"github.com/injoyai/base/g"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/database/mysql"
	"github.com/injoyai/goutil/database/sqlite"
	"github.com/injoyai/goutil/database/xorms"
	"time"
)

const (
	DefaultUsername = "admin"            //默认管理员账号
	DefaultPassword = "admin"            //默认管理员密码
	All             = "all"              //所有权限
	Sqlite          = "sqlite"           //sqlite数据库
	Mysql           = "mysql"            //mysql数据库
	Redis           = "redis"            //redis缓存
	Memory          = "memory"           //内存缓存
	Filename        = "database/user.db" //
)

var (
	DB    *xorms.Engine
	Users = maps.NewSafe()
	Auth  *auth
)

func Init(conf *Config) (err error) {
	switch conf.Type {
	case Sqlite:
		DB, err = sqlite.NewXorm(conf.DSN)
	case Mysql:
		DB, err = mysql.NewXorm(conf.DSN)
	default:
		return errors.New("未知类型: " + conf.Type)
	}
	if err != nil {
		return
	}

	Auth = initToken(conf)

	if err = DB.Sync(new(User)); err != nil {
		return err
	}
	data, err := GetAll()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		_, err = DB.Insert(&User{Username: DefaultUsername, Password: DefaultPassword, Limit: []string{All}})
		if err != nil {
			return err
		}
	}
	for _, v := range data {
		Users.Set(v.Username, v)
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

func CheckToken(token string) (u *User, valid bool, err error) {
	if Auth.IsSuper(token) {
		return &User{Username: DefaultUsername, Limit: []string{All}}, true, nil
	}
	username, err := Auth.Cache.Get(token)
	if err != nil {
		return nil, false, err
	}
	if len(username) == 0 {
		return nil, false, nil
	}
	u, err = GetByCache(username)
	if err != nil {
		if err.Error() == "用户不存在" {
			return nil, false, nil
		}
		return nil, false, err
	}
	return u, true, nil
}

func GetByCache(username string) (*User, error) {
	val, ok := Users.Get(username)
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
	err = Auth.Cache.Set(token, user.Username, time.Hour*24*3)

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
			Users.Set(user.Username, user)
		}
		return err
	}
	_, err = DB.Where("Username=?", user.Username).Update(user)
	if err == nil {
		Users.Set(user.Username, user)
	}
	return err
}

func Del(username string) error {
	_, err := DB.Where("Username=?", username).Delete(&User{})
	if err == nil {
		Users.Del(username)
	}
	return err
}
