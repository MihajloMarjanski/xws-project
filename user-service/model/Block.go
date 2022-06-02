package model

type Block struct {
	UserId      uint `json:"userId" gorm:"primaryKey;autoIncrement:false"`
	BlockedUser uint `json:"blockedUserId" gorm:"primaryKey;autoIncrement:false"`
}
