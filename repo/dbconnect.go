package repo

import (
	"database/sql"
	"encoding/json"
	"final-project/entity"
	"fmt"
	"net/http"

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

func GetDB() *sql.DB {
	return DB
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

// func checkCount(rows *sql.Rows) (count int) {
// 	for rows.Next() {
// 		err := rows.Scan(&count)
// 		checkErr(err)
// 	}
// 	return count
// }

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Println("hehe")
		http.ServeFile(w, r, "view/register.html")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	age := r.FormValue("age")
	fmt.Println(age)
	temppassword := "aa"
	fmt.Println(username, password)
	//deskripsi dan compare password
	var password_tes = bcrypt.CompareHashAndPassword([]byte(temppassword), []byte(password))
	fmt.Println(password_tes)
}
func UsernameQuery(w http.ResponseWriter, r *http.Request, username string) {
	sqlStatement := "Select * from users where username = $1"
	err := DB.QueryRow(sqlStatement, username).Scan(&username)
	//handler.CheckErr(err)
	if err != nil {
		panic(err)
	}
}

// func login2(w http.ResponseWriter, r *http.Request) {
// 	session := sessions.Start(w, r)
// 	if len(session.GetString("username")) != 0 && checkErr(w, r, err) {
// 		http.Redirect(w, r, "/", 302)
// 	}
// 	if r.Method != "POST" {
// 		http.ServeFile(w, r, "views/login.html")
// 		return
// 	}
// 	username := r.FormValue("username")
// 	password := r.FormValue("password")

// 	users := QueryUser(username)

// 	//deskripsi dan compare password
// 	var password_tes = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))

// 	if password_tes == nil {
// 		//login success
// 		session := sessions.Start(w, r)
// 		session.Set("username", users.Username)
// 		session.Set("name", users.FirstName)
// 		http.Redirect(w, r, "/", 302)
// 	} else {
// 		//login failed
// 		http.Redirect(w, r, "/login", 302)
// 	}

// }

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
