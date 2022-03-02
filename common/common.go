package common

import (
	"final-project/entity"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

func EncryptPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
func AuthenticateJWT(w http.ResponseWriter, r *http.Request) bool {
	tmp := true
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(1)
			return tmp == false
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(2)
		return tmp == false
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &entity.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(3)
			return tmp == false
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(4)
		return tmp == false
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(5)
		return tmp == false
	}
	return tmp
}
func GetJWTData(w http.ResponseWriter, r *http.Request) entity.Claims {
	c, _ := r.Cookie("token")
	tknStr := c.Value
	claims := &entity.Claims{}
	_, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	CheckErr(err)
	return *claims
}
