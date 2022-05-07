package model

import "time"

type Post struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	Title     string    `json:"name"`
	Text      string    `json:"username"`
	Img       string    `json:"email"`
	Link      string    `json:"password"`
	CreatedAt time.Time `json:"date"`
}
