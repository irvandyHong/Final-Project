package repo

import (
	"database/sql"
	"final-project/common"
	"final-project/entity"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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
var i int
var DB *sql.DB
var err error
var student entity.Person

func ConnectDB() {
	connString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	DB, err = sql.Open("postgres", connString)
	CheckErr(err)
	err = DB.Ping()
	CheckErr(err)
	fmt.Println("connected")
}

func GetDB() *sql.DB {
	return DB
}
func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != "POST" {
		http.ServeFile(w, r, "view/login.html")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := QueryUser(username)
	//deskripsi dan compare password
	var password_tes = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	CheckErr(password_tes)
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &entity.Claims{
		Username: username,
		Id:       user.Id,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	fmt.Println(tokenString)
	CheckErr(err)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	http.Redirect(w, r, "/upload", http.StatusSeeOther)
}

var jwtKey = []byte("my_secret_key")

//get user data for login
func QueryUser(username string) entity.Person {
	var users = entity.Person{}
	sqlstatement := "SELECT id, username,password FROM users WHERE username=$1"
	err = DB.QueryRow(sqlstatement, username).
		Scan(
			&users.Id,
			&users.Username,
			&users.Password,
		)
	fmt.Println(users.Id)
	return users
}
func Register(w http.ResponseWriter, r *http.Request) {
	var person = entity.Person{}
	r.ParseForm()
	if r.Method != "POST" {
		http.ServeFile(w, r, "view/register.html")
		return
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		EncryptPassword := common.EncryptPassword(password)
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
	DB.QueryRow(sqlStatement, person.Username, person.Password, person.Email, person.Age, person.Id, person.CreatedAt, person.UpdatedAt).
		Scan(&person.Username, &person.Password, &person.Email, &person.Age, &person.Id, &person.CreatedAt, &person.UpdatedAt)
	fmt.Println("reach")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

//count rows length for user id
func GetUserId(w http.ResponseWriter, r *http.Request) (count int) {
	var result entity.Person
	fmt.Println("reach")
	sqlStatement := "Select id from users order by users desc limit 1"
	rows, err := DB.Query(sqlStatement)
	CheckErr(err)
	for rows.Next() {
		rows.Scan(&result.Id)
	}
	return result.Id + 1
	// w.Header().Set("Content-Type", "application/json")
	// SqlStatement := "select * from users"
	// rows, err := DB.Query(SqlStatement)

	// for rows.Next() {
	// 	if err != nil {
	// 		err := rows.Scan(&count)
	// 		CheckErr(err)
	// 	}
	// 	count++
	// }
	// return count
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PhotoQuery(w http.ResponseWriter, r *http.Request, url, title, caption string) {
	c, _ := r.Cookie("token")
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &entity.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	_, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	common.CheckErr(err)
	r.ParseForm()
	var photo = entity.Photo{}
	id := GetPhotoId(w, r)
	photo_url := url
	user_id := claims.Id
	created_at := time.Now()
	updated_at := time.Now()
	photo = entity.Photo{
		Id:        id,
		Title:     title,
		Caption:   caption,
		PhotoUrl:  photo_url,
		User_Id:   user_id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}
	sqlStatement := "insert Into photo(id,title,caption,photo_url,user_id,created_at,updated_at) Values($1,$2,$3,$4,$5,$6,$7) returning* "
	DB.QueryRow(sqlStatement, photo.Id, photo.Title, photo.Caption, photo.PhotoUrl, photo.User_Id, photo.CreatedAt, photo.UpdatedAt).
		Scan(&photo.Id, &photo.Title, &photo.Caption, &photo.PhotoUrl, &photo.User_Id, &photo.CreatedAt, &photo.UpdatedAt)
}

//count rows length for user id
func GetPhotoId(w http.ResponseWriter, r *http.Request) (count int) {
	var result entity.Photo
	fmt.Println("reach")
	sqlStatement := "Select id from photo order by photo desc limit 1"
	rows, err := DB.Query(sqlStatement)
	CheckErr(err)
	for rows.Next() {
		rows.Scan(&result.Id)
	}
	return result.Id + 1
}

func GetPhotoData(w http.ResponseWriter, r *http.Request, num int) entity.Photo {
	count := 0
	var result []entity.Photo
	sqlStatement := "Select photo_url,title,caption from photo"
	tmp, err := DB.Query(sqlStatement)
	for tmp.Next() {
		var photo = entity.Photo{}
		err := tmp.Scan(&photo.PhotoUrl, &photo.Title, &photo.Caption)
		CheckErr(err)
		result = append(result, photo)
		count++
	}
	CheckErr(err)
	if count == 0 {
		temp := entity.Photo{
			PhotoUrl: "temp-images/",
			Title:    "empty",
			Caption:  "empty",
		}
		return temp
	} else if num < count {
		return result[num]
	} else {
		temp := entity.Photo{
			PhotoUrl: "temp-images/",
			Title:    "empty",
			Caption:  "empty",
		}
		return temp
	}
}
func GetUserData(w http.ResponseWriter, r *http.Request) entity.Person {
	c, _ := r.Cookie("token")
	tknStr := c.Value
	claims := &entity.Claims{}
	_, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	common.CheckErr(err)
	var result entity.Person
	username := claims.Username
	sqlStatement := "Select id,username,email,age,updated_at from users where username = $1"
	DB.QueryRow(sqlStatement, username).Scan(&result.Id, &result.Username, &result.Email, &result.Age, &result.UpdatedAt)
	return result
}
func EditUserData(w http.ResponseWriter, r *http.Request, username, email string) {
	claims := common.GetJWTData(w, r)
	sqlStatement := "update users set username = $1,email =$2 where id = $3"
	DB.QueryRow(sqlStatement, username, email, claims.Id)
}
func DeleteUserData(w http.ResponseWriter, r *http.Request) {
	claims := common.GetJWTData(w, r)
	sqlStatement := "Delete from users where id = $1"
	DB.QueryRow(sqlStatement, claims.Id)
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}
func GetUserPhoto(w http.ResponseWriter, r *http.Request, id, num int) entity.Photo {
	count := 0
	var result []entity.Photo
	sqlStatement := "Select photo_url,title,caption from photo where user_id =$1"
	tmp, err := DB.Query(sqlStatement, id)
	for tmp.Next() {
		var photo = entity.Photo{}
		err := tmp.Scan(&photo.PhotoUrl, &photo.Title, &photo.Caption)
		CheckErr(err)
		result = append(result, photo)
		count++
	}
	CheckErr(err)
	if count == 0 {
		temp := entity.Photo{
			PhotoUrl: "temp-images/",
			Title:    "empty",
			Caption:  "empty",
		}
		return temp
	} else if num < count {
		return result[num]
	} else {
		temp := entity.Photo{
			PhotoUrl: "temp-images/",
			Title:    "empty",
			Caption:  "empty",
		}
		return temp
	}
}
func DeletePhoto(w http.ResponseWriter, r *http.Request, id string) {
	sqlStatement := "Delete from photo where photo_url = $1"
	DB.QueryRow(sqlStatement, id)
	fmt.Println("called")
}
func EditPhoto(w http.ResponseWriter, r *http.Request, title, caption, id string) {
	fmt.Println(id)
	sqlStatement := "update photo set title = $1,caption =$2 where photo_url = $3"
	fmt.Println("called")
	DB.QueryRow(sqlStatement, title, caption, id)
}
func DoNothing(w http.ResponseWriter, r *http.Request) {}
func GetCommentId(w http.ResponseWriter, r *http.Request) int {
	var result entity.Comment
	fmt.Println("reach")
	sqlStatement := "Select id from comment order by comment desc limit 1"
	rows, err := DB.Query(sqlStatement)
	CheckErr(err)
	for rows.Next() {
		rows.Scan(&result.Id)
	}
	return result.Id + 1
}
func PostComment(w http.ResponseWriter, r *http.Request, id, photo, user int, time time.Time, message string) {
	comment := entity.Comment{
		Id:        id,
		Message:   message,
		CreatedAt: time,
		Photo_id:  photo,
		UserId:    user,
	}
	sqlStatement := "insert into comment (id,message,created_at,photo_id,user_id) values($1,$2,$3,$4,$5) returning*"
	fmt.Println(comment.Id, comment.Message, comment.CreatedAt, comment.Photo_id, comment.UserId)
	DB.QueryRow(sqlStatement, comment.Id, comment.Message, comment.CreatedAt, comment.Photo_id, comment.UserId).
		Scan(&comment.Id, &comment.Message, &comment.CreatedAt, &comment.Photo_id, &comment.UserId)
}
func GetComment(w http.ResponseWriter, r *http.Request, num int) entity.Comment {
	result := []entity.Comment{}
	count := 0
	sqlStatement := "select * from comment"
	rows, _ := DB.Query(sqlStatement)
	for rows.Next() {
		comment := entity.Comment{}
		rows.Scan(&comment.Id, &comment.Message, &comment.CreatedAt, &comment.UserId, &comment.Photo_id)
		result = append(result, comment)
		count++
	}
	if count == 0 {
		temp := entity.Comment{
			Message: "empty",
		}
		return temp
	} else if num < count {
		return result[num]
	} else {
		temp := entity.Comment{
			Message: "empty",
		}
		return temp
	}
}
