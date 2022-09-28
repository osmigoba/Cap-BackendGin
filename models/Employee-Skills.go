package models

// import (
// 	"gorm.io/gorm"
// )

type EmployeeSkills struct {
	// gorm.Model

	Employee   *Employee `json:"employee" gorm:"primaryKey:EmployeeID;constraint:OnDelete:CASCADE"`
	EmployeeID int64     `json:"employee_Id"`

	Skill   *Skill `json:"skill" gorm:"primaryKey:SkillID;constraint:OnDelete:CASCADE"`
	SkillID int64  `json:"skill_ID" `

	LevelRating *Level `json:"levelRating" gorm:"foreignKey:LevelID;constraint:OnDelete:CASCADE"`
	LevelID     int64  `json:"levelRating_ID"`
	Experience  int64  `json:"experience"`
}
