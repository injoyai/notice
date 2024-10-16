package user

const (
	TypeTCP  = "tcp"
	TypeHTTP = "http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (this *User) Login(Type string) {

}
