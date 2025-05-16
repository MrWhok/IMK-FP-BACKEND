package entity

type User struct {
	Username  string     `gorm:"primaryKey;column:username;type:varchar(100)"`
	Password  string     `gorm:"column:password;type:varchar(200)"`
	FirstName string     `gorm:"column:first_name;type:varchar(100)"`
	LastName  string     `gorm:"column:last_name;type:varchar(100)"`
	Email     string     `gorm:"column:email;type:varchar(100)"`
	Phone     string     `gorm:"column:phone;type:varchar(100)"`
	Address   string     `gorm:"column:address;type:varchar(255)"`
	IsActive  bool       `gorm:"column:is_active;type:boolean"`
	Points    int32      `gorm:"column:points;type:int;default:0"`
	UserRoles []UserRole `gorm:"ForeignKey:Username;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "tb_user"
}
