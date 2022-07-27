package model

type UsernameWithRequestId struct {
	ReceiverId int64  `json:"receiverId"`
	SenderId   int64  `json:"senderId"`
	UserName   string `json:"username"`
}
