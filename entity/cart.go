package entity

type Cart struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;column:cart_id"`
	Username  string `gorm:"column:username;type:varchar(100);not null;unique"`
	ProductID string `gorm:"column:product_id;type:varchar(36);not null"`
	Quantity  int32  `gorm:"column:quantity"`

	Product Product `gorm:"foreignKey:ProductID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Cart) TableName() string {
	return "tb_cart"
}
