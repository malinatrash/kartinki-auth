package kafka

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Secret   string `json:"secret"`
}
