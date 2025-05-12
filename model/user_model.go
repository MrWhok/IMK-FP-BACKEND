package model

type UserModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Roles    []string `json:"roles"`
	Points   int32 `json:"points"`
}
