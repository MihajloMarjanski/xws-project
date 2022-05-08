package model

type Connection struct {
	UserOne uint `json:"userOne" gorm:"primaryKey;autoIncrement:false"`
	UserTwo uint `json:"userTwo" gorm:"primaryKey;autoIncrement:false"`
}
