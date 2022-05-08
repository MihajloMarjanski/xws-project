package model

import "time"

type Connection struct{
	Date time.Time `json:"date"`
	UserOne uint `json:"userOne" gorm:"primaryKey;autoIncrement:false"`
	UserTwo uint `json:"userTwo" gorm:"primaryKey;autoIncrement:false"`
}