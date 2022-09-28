package dtos

import "time"

type AddSkillEmployeeDTO struct {
	EmployeeID    int64 `json:"employeeID" `
	SkillID       int64 `json:"skillID"`
	LevelRatingId int64 `json:"levelRatingId"`
	Experience    int64 `json:"experience"`
}

type SkillsByEmployeeDTO struct {
	Skillid    int64  `json:"skillID" `
	Skill      string `json:"skill"`
	Levelid    int64  `json:"levelId"`
	Level      string `json:"level"`
	Experience int64  `json:"experience"`
}

type UserLogin struct {
	Email      string    `json:"email"`
	Expiration time.Time `json:"expiration"`
	Token      string    `json:"token"`
}

type EmployeeLevel struct {
	// gorm.Model

	Id          int64     `json:"id" gorm:"primaryKey;autoIncrement:true;constraint:OnDelete:CASCADE;"`
	FirstName   string    `json:"firstName" gorm:"type:varchar(50)"`
	LastName    string    `json:"lastName" gorm:"type:varchar(50)"`
	DOJ         time.Time `json:"doj"`
	Designation string    `json:"designation" gorm:"type:varchar(50)"`
	Email       string    `json:"email" gorm:"type:varchar(150)"`
	Skill       string    `json:"skill" gorm:"type:varchar(150)"`
	Levelname   string    `json:"levelName" gorm:"type:varchar(150)"`
}
