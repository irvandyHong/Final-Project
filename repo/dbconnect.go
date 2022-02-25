package repo

import (
	"database/sql"
	"encoding/json"
	"final-project/entity"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "irvandy2"
	password = "koinworks"
	dbname   = "final"
	PORT     = ":8080"
)

var connString string
var DB *sql.DB
var err error
var student entity.Person

func ConnectDB() {
	connString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	DB, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("connected")
	//defer db.Close()
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result = []entity.Person{}
	SqlStatement := "select * from user"
	rows, err := DB.Query(SqlStatement)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var person = entity.Person{}
		err = rows.Scan(&person.Email)
		if err != nil {
			panic(err)
		}
		result = append(result, person)
	}
	fmt.Println("Order Data", result)
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(result)
		return
	}
	http.Error(w, "Invalid Method", http.StatusBadRequest)
}

func GetUserId(w http.ResponseWriter, r *http.Request) (count int) {
	i := 1
	w.Header().Set("Content-Type", "application/json")
	SqlStatement := "select * from user"
	rows, err := DB.Query(SqlStatement)

	for rows.Next() {
		if err != nil {
			err := rows.Scan(&count)
			checkErr(err)
		}
	}
	fmt.Println(i)
	return count
}
func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Println("hehe")
		http.ServeFile(w, r, "view/login.html")
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
func UsernameQuery(w http.ResponseWriter, r *http.Request, username string) {
	sqlStatement := "Select * from users where username = $1"
	err := DB.QueryRow(sqlStatement, username).Scan(&username)
	checkErr(err)

}
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	if err != nil {
		panic(err)
	}
}

func EncryptPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
// 	t, err := template.ParseFiles(filename)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	if err := t.Execute(w, data); err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// }
// func OutputHandler(w http.ResponseWriter, r *http.Request) {
// 	myvar := map[string]interface{}{"MyVar": "hehe"}
// 	outputHTML(w, "view/register.html", myvar)
// }
