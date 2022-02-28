package handler

import (
	"final-project/entity"
	"final-project/repo"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var DB = repo.GetDB()

func RegisterHandler() {

}
func Register(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	var person = entity.Person{}
	if r.Method != "POST" {
		fmt.Println("hehe")
		http.ServeFile(w, r, "view/register.html")
		return
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		EncryptPassword := EncryptPassword(password)
		email := r.FormValue("email")
		age := r.FormValue("age")
		id := GetUserId(w, r)
		createdAt := time.Now()
		updatedAt := time.Now()
		person = entity.Person{
			Id:        id,
			Username:  username,
			Password:  EncryptPassword,
			Email:     email,
			Age:       age,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}
	sqlStatement := "insert Into users(username,password,email,age,id,created_at,updated_at) Values($1,$2,$3,$4,$5,$6,$7) returning*"
	err := DB.QueryRow(sqlStatement, person.Username, person.Password, person.Email, person.Age, person.Id, person.CreatedAt, person.UpdatedAt).
		Scan(&person.Username, &person.Password, &person.Email, &person.Age, &person.Id, &person.CreatedAt, &person.UpdatedAt)
	CheckErr(err)
}
func Login2(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Println("hehe")
		http.ServeFile(w, r, "view/register.html")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	temppassword := "aa"
	fmt.Println(username, password)
	//deskripsi dan compare password
	var password_tes = bcrypt.CompareHashAndPassword([]byte(temppassword), []byte(password))
	fmt.Println(password_tes)
}
func GetUserId(w http.ResponseWriter, r *http.Request) (count int) {
	i := 1
	w.Header().Set("Content-Type", "application/json")
	SqlStatement := "select * from user"
	rows, err := DB.Query(SqlStatement)

	for rows.Next() {
		if err != nil {
			err := rows.Scan(&count)
			CheckErr(err)
		}
	}
	fmt.Println(i)
	return count
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
func EncryptPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
