package entity

import "time"

type Person struct {
	Id        int
	Age       string
	Email     string
	Password  string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Photo struct {
	Title    string
	Caption  string
	PhotoUrl string
}
type Comment struct {
	Message  string
	Photo_id int
}
type SocialMedia struct {
	Name           string
	SocialMediaUrl string
}
