package model

import "github.com/dgrijalva/jwt-go"

type Connection struct {
	UserOne uint `json:"userOne" gorm:"primaryKey;autoIncrement:false"`
	UserTwo uint `json:"userTwo" gorm:"primaryKey;autoIncrement:false"`
}

type Claims struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
