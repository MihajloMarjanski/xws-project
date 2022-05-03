package model

import "time"

type User struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name" gorm:"not null"`
	UserName    string       `json:"username" gorm:"unique;not null"`
	Email       string       `json:"email" gorm:"unique;not null"`
	Password    string       `json:"password" gorm:"not null"`
	Gender      Gender       `json:"gender"`
	PhoneNumber string       `json:"phone"`
	DateOfBirth time.Time    `json:"date"`
	Biography   string       `json:"biography"`
	Interests   []Interest   `json:"interests"`
	Experiences []Experience `json:"experiences"`
}

type UserDTO struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name" gorm:"not null"`
	UserName    string       `json:"username" gorm:"unique;not null"`
	Email       string       `json:"email" gorm:"unique;not null"`
	Gender      Gender       `json:"gender"`
	PhoneNumber string       `json:"phone"`
	DateOfBirth time.Time    `json:"date"`
	Biography   string       `json:"biography"`
	Interests   []Interest   `json:"interests"`
	Experiences []Experience `json:"experiences"`
}

type Interest struct {
	ID       uint   `json:"id"`
	Interest string `json:"interest" gorm:"not null"`
	UserID   uint   `json:"user" gorm:"not null"`
}

type Experience struct {
	ID       uint      `json:"id"`
	Company  string    `json:"company" gorm:"not null"`
	Position string    `json:"postiion" gorm:"not null"`
	From     time.Time `json:"from" gorm:"not null"`
	Until    time.Time `json:"until" gorm:"not null"`
	UserID   uint      `json:"user" gorm:"not null"`
}

type Gender string

const (
	Male   Gender = "male"
	Female        = "female"
)

type ResponseId struct {
	Id int `json:"id"`
}
