package model

type Request struct {
	SenderID   uint `json:"sender_id" gorm:"primaryKey;autoIncrement:false"`
	ReceiverID uint `json:"receiver_id" gorm:"primaryKey;autoIncrement:false"`
}
