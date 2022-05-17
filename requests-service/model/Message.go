package model

type Message struct {
	ID         uint   `json:"id"`
	Text       string `json:"text" gorm:"not null"`
	SenderId   uint   `json:"senderId" gorm:"not null"`
	ReceiverId uint   `json:"receiverId" gorm:"not null"`
}
