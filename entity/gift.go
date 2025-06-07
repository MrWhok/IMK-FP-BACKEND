package entity

import "github.com/google/uuid"

type Gift struct {
	Id         uuid.UUID `gorm:"primaryKey;column:gift_id;type:varchar(36)"`
	Name       string    `gorm:"index;column:name;type:varchar(100)"`
	PointPrice int64     `gorm:"column:point_price"`
	Quantity   int32     `gorm:"column:quantity"`
	ImagePath  string    `gorm:"column:image_path;type:varchar(255)"`
}

func (Gift) TableName() string {
	return "tb_gift"
}
