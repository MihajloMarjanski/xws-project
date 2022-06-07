package model

import "time"

type User struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	UserName    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone"`
	DateOfBirth time.Time `json:"date"`
	Biography   string    `json:"biography"`
}
