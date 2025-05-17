package entity

type Cart struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;column:cart_id"`
	Username string `gorm:"column:username;type:varchar(100);not null;unique"`

	// User  User       `gorm:"foreignKey:Username;references:Username;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Items []CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Cart) TableName() string {
	return "tb_cart"
}

type CartItem struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;column:item_id"`
	CartID    uint   `gorm:"column:cart_id;not null;index"`
	ProductID string `gorm:"column:product_id;type:varchar(36);not null"`
	Quantity  int32  `gorm:"column:quantity"`

	Product Product `gorm:"foreignKey:ProductID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CartItem) TableName() string {
	return "tb_cart_item"
}
