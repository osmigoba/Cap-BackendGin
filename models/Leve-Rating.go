package models

// import (
// 	"gorm.io/gorm"
// )

type Level struct {
	Id   int64  `json:"id" gorm:"primaryKey;autoIncrement:true;constraint:OnDelete:CASCADE;"`
	Name string `json:"name" gorm:"type:varchar(50)"`
}
