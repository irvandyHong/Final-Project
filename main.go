package main

import (
	"final-project/repo"
	"net/http"
)

var PORT = ":8080"

func main() {
	repo.ConnectDB()
	//http.HandleFunc("/", repo.OutputHandler)
	http.HandleFunc("/register", repo.Register)
	http.HandleFunc("/login", repo.Login)
	//http.HandleFunc("/register", repo.InsertUser)
	//http.HandleFunc("/get", repo.Temp)
	http.ListenAndServe(PORT, nil)
	defer repo.DB.Close()

}
