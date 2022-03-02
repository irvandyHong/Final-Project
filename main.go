package main

import (
	"final-project/handler"
	"final-project/repo"
	"net/http"
)

var PORT = ":8080"

func main() {
	repo.ConnectDB()
	http.Handle("/temp-images/", http.StripPrefix("/temp-images", http.FileServer(http.Dir("./temp-images"))))
	http.HandleFunc("/", handler.Welcome)
	http.HandleFunc("/users", handler.Profile)
	http.HandleFunc("/users/login", repo.Login)
	http.HandleFunc("/users/register", repo.Register)
	http.HandleFunc("/users/delete", repo.DeleteUserData)
	http.HandleFunc("/users/photo", handler.MyPhoto)
	http.HandleFunc("/upload", handler.UploadFile)
	http.HandleFunc("/editphoto", handler.EditPhoto)
	http.HandleFunc("/favicon.ico", repo.DoNothing)
	http.HandleFunc("/comment", handler.CommentHandler)
	http.ListenAndServe(PORT, nil)
	defer repo.DB.Close()
}
