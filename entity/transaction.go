package entity

import "github.com/google/uuid"

type Transaction struct {
	Id                 uuid.UUID           `gorm:"primaryKey;column:transaction_id;type:varchar(36)"`
	TotalPrice         int64               `gorm:"column:total_price"`
	UserID             string              `gorm:"column:user_id;type:varchar(100);index"`
	Status             string              `gorm:"column:status;type:varchar(20);default:'proses'"`
	User               User                `gorm:"foreignKey:UserID;references:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TransactionDetails []TransactionDetail `gorm:"ForeignKey:TransactionId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Transaction) TableName() string {
	return "tb_transaction"
}
