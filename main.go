package main

import (
	"final-project/handler"
	"final-project/repo"
	"net/http"
)

var PORT = ":8080"

func main() {
	repo.ConnectDB()
	//http.HandleFunc("/", repo.OutputHandler)
	http.HandleFunc("/login", repo.Login)
	http.HandleFunc("/register", handler.Login2)
	http.ListenAndServe(PORT, nil)
	defer repo.DB.Close()

}
