package model

type UserModel struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	Points   int32    `json:"points"`
}

type UserCreateModel struct {
	Username  string   `json:"username" validate:"required"`
	Password  string   `json:"password" validate:"required"`
	FirstName string   `json:"first_name" validate:"required"`
	LastName  string   `json:"last_name" validate:"required"`
	Email     string   `json:"email" validate:"required,email"`
	Phone     string   `json:"phone" validate:"required"`
	Address   string   `json:"address" validate:"required"`
	Roles     []string `json:"roles" validate:"required"`
	Points    int32    `json:"points" validate:"required"`
}

type UserLeaderboardModel struct {
	Rank      int    `json:"rank"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Points    int32  `json:"points"`
}
