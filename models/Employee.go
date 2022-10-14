package models

import (
	// "gorm.io/gorm"
	"time"
)

type Employee struct {
	// gorm.Model

	Id          int64     `json:"id" gorm:"primaryKey;autoIncrement:true;constraint:OnDelete:CASCADE;" `
	FirstName   string    `json:"firstName" gorm:"type:varchar(50)" binding:"required"`
	LastName    string    `json:"lastName" gorm:"type:varchar(50)" binding:"required"`
	DOJ         time.Time `json:"doj"`
	Designation string    `json:"designation" gorm:"type:varchar(50)" binding:"required"`
	Email       string    `json:"email" gorm:"type:varchar(150)" binding:"required,email"`
}
