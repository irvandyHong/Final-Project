package main

import (
	"final-project/repo"
	"net/http"
)

var PORT = ":8081"

func main() {
	repo.ConnectDB()
	http.HandleFunc("/", repo.OutputHandler)
	http.HandleFunc("/register", repo.InsertUser)
	http.HandleFunc("/get", repo.Temp)
	http.ListenAndServe(PORT, nil)
	defer repo.DB.Close()
}
