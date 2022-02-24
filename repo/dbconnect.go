package repo

import (
	"database/sql"
	"encoding/json"
	"final-project/entity"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	sqlStatement := "select * from user"
	rows, err := DB.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	var person = entity.Person{}
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")
		ages := r.FormValue("age")
		age, _ := strconv.Atoi(ages)
		person = entity.Person{
			Username: username,
			Password: password,
			Email:    email,
			Age:      age,
		}
	}

	for rows.Next() {
		var email = entity.Person{}
		err = rows.Scan(email.Email)
		fmt.Println(email.Email)
		if email.Email == person.Email {
			panic(err)
		}
	}
	sqlstatement := "insert into user(username,password,email,age,id) value ($1,$2,$3,$4,$5) returning *"
	err = DB.QueryRow(sqlstatement, person.Username, person.Password, person.Email, person.Age, 1).Scan(&person.Username, &person.Password, &person.Email, &person.Age, 1)
	if err != nil {
		panic(err)
	}
}
func GetEmail(w http.ResponseWriter, r *http.Request) {

}
func Tmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person = entity.Person{}
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")
		tempAges := r.FormValue("age")
		age, _ := strconv.Atoi(tempAges)
		id := r.FormValue("id")
		createdAt := time.Now().String()
		updatedAt := time.Now().String()
		person = entity.Person{
			Username:  username,
			Password:  password,
			Email:     email,
			Age:       age,
			Id:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}
	sqlStatement := "insert into user(username,password,email,age,id,created_at,updated_at) Values($1,$2,$3,$4,$5,$6,$7) returning*"
	err := DB.QueryRow(sqlStatement, person.Username, person.Password, person.Email, person.Age, person.Id, person.CreatedAt, person.UpdatedAt).
		Scan(&person.Username, &person.Password, &person.Email, &person.Age, &person.Id, &person.CreatedAt, &person.UpdatedAt)
	if err != nil {
		panic(err)
	}
}
func InsertUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person = entity.Person{}
	if r.Method == "POST" {
		username := r.FormValue("username")
		person = entity.Person{
			Username: username,
		}
	}
	sqlStatement := "insert into user(username) values ($1) returning*"
	err := DB.QueryRow(sqlStatement, person.Username)
	if err != nil {
		panic(err)
	}
	fmt.Print("hehe")
}

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func OutputHandler(w http.ResponseWriter, r *http.Request) {
	myvar := map[string]interface{}{"MyVar": "hehe"}
	outputHTML(w, "register.html", myvar)
}
func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//var order = orders{}
	var result = []entity.Person{}
	sqlStatement := "select * from user"
	rows, err := DB.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var order = entity.Person{}

		err = rows.Scan(&order.Id)
		if err != nil {
			panic(err)
		}
		result = append(result, order)
	}
	fmt.Println("Order Data", result)
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(result)
		return
	}

	http.Error(w, "Invalid Method", http.StatusBadRequest)
}
func Temp(w http.ResponseWriter, r *http.Request) {
	sqlStatement := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name;"
	rows, err := DB.Query(sqlStatement)
	fmt.Fprint(w, err, rows)
}
