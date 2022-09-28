package models

type User struct {
	Id       int64  `json:"id" gorm:"primaryKey;autoIncrement:true;"`
	Email    string `json:"email" gorm:"unique;unique_index"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
