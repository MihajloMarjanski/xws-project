package model

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
}

type ResponseId struct {
	Id int `json:"id"`
}
