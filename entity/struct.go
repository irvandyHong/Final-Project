package entity

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Person struct {
	Id        int       `json:"id"`
	Age       string    `json:"age"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}
type Photo struct {
	Id        int
	Title     string
	Caption   string
	PhotoUrl  string
	User_Id   int
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}
type Comment struct {
	Id        int
	Message   string
	Photo_id  int
	CreatedAt time.Time
	UserId    int
}
type SocialMedia struct {
	Name           string
	SocialMediaUrl string
}
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
type Claims struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	jwt.StandardClaims
}
