package model

import "time"

type User struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" gorm:"not null"`
	UserName    string    `json:"username" gorm:"unique;not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Password    string    `json:"password" gorm:"not null"`
	Gender      Gender    `json:"gender"`
	PhoneNumber string    `json:"phone"`
	DateOfBirth time.Time `json:"date"`
	Biography   string    `json:"biography"`
}
type UserDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" gorm:"not null"`
	UserName    string    `json:"username" gorm:"unique;not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Gender      Gender    `json:"gender"`
	PhoneNumber string    `json:"phone"`
	DateOfBirth time.Time `json:"date"`
	Biography   string    `json:"biography"`
}

type Gender string

const (
	Male   Gender = "male"
	Female        = "female"
)

type ResponseId struct {
	Id int `json:"id"`
}
