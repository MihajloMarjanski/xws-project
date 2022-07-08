package model

type Connection struct {
	UserOne uint `json:"userOne"`
	UserTwo uint `json:"userTwo"`
}

type User struct {
	UserId string
}
