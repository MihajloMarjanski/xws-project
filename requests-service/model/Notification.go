package model

import "time"

type Notification struct {
	ID         uint      `json:"id"`
	ReceiverId uint      `json:"receiverId" gorm:"not null"`
	Text       string    `json:"text" gorm:"not null"`
	Date       time.Time `json:"date" gorm:"not null"`
}
