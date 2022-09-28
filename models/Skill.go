package models

import "time"

type Skill struct {
	Id          int64     `json:"id" gorm:"primaryKey;autoIncrement:true;constraint:OnDelete:CASCADE;"`
	Skill       string    `json:"skill" gorm:"type:varchar(50)"`
	CreatedDate time.Time `json:"createdDate"`
}
