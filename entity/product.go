package entity

import "github.com/google/uuid"

type Product struct {
	Id                 uuid.UUID           `gorm:"primaryKey;column:product_id;type:varchar(36)"`
	Name               string              `gorm:"index;column:name;type:varchar(100)"`
	Price              int64               `gorm:"column:price"`
	Quantity           int32               `gorm:"column:quantity"`
	Category           string              `gorm:"column:category;type:varchar(50)"`
	Description        string              `gorm:"column:description;type:varchar(255)"`
	ImagePath          string              `gorm:"column:image_path;type:varchar(255)"`
	UserID             string              `gorm:"column:user_id;type:varchar(100);index"`
	Owner              User                `gorm:"foreignKey:UserID;references:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TransactionDetails []TransactionDetail `gorm:"ForeignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Product) TableName() string {
	return "tb_product"
}
