package entity

type Person struct {
	Id        string
	Age       int
	Email     string
	Password  string
	Username  string
	CreatedAt string
	UpdatedAt string
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
