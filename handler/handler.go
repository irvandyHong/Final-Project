package handler

import (
	"final-project/common"
	"final-project/repo"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"
)

var i, j int
var t *template.Template
var ta *template.Template
var err error

func Welcome(w http.ResponseWriter, r *http.Request) {
	//i := 0
	var t *template.Template
	//var err error
	common.AuthenticateJWT(w, r)
	if common.AuthenticateJWT(w, r) == false {
		fmt.Println("error")
	} else {
		tmp(w, r)
	}

	if r.Method == http.MethodPost {

		url := repo.GetPhotoData(w, r, i)
		i++
		ta, err := template.ParseFiles("view/home.html")
		myvar := map[string]interface{}{"url1": url.PhotoUrl, "title": url.Title, "caption": url.Caption}
		ta.Execute(w, myvar)
		t.Execute(os.Stdout, myvar)
		common.CheckErr(err)
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
}
func tmp(w http.ResponseWriter, r *http.Request) {
	url := repo.GetPhotoData(w, r, i)
	t, err = template.ParseFiles("view/home.html")
	myvar := map[string]interface{}{"url1": url.PhotoUrl, "title": url.Title, "caption": url.Caption}
	t.Execute(w, myvar)
	common.CheckErr(err)
}
func UploadFile(w http.ResponseWriter, r *http.Request) {
	if common.AuthenticateJWT(w, r) {
		if r.Method != "POST" {
			http.ServeFile(w, r, "view/upload.html")
			fmt.Fprint(w, " <img src=\"/temp-images/upload-092252896.png\" alt=\"t\">")
			return
		}
		fmt.Println("File Upload Endpoint Hit")

		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(10 << 20)
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
		if err != nil {
			fmt.Println(err)
		}
		title := r.FormValue("title")
		caption := r.FormValue("caption")
		repo.PhotoQuery(w, r, tempFile.Name(), title, caption)
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		// write this byte array to our temporary file
		tempFile.Write(fileBytes)
		// return that we have successfully uploaded our file!
		fmt.Fprintf(w, "Successfully Uploaded File\n")
	} else {
	}

}
func Profile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	common.AuthenticateJWT(w, r)
	//var t template.Template
	if common.AuthenticateJWT(w, r) == false {
		fmt.Println("error")
	} else {
		tmp(w, r)
	}
	url := repo.GetUserData(w, r)
	ta, err = template.ParseFiles("view/profile.html")
	myvar := map[string]interface{}{"id": url.Id, "username": url.Username, "email": url.Email, "age": url.Age, "updated_at": url.UpdatedAt}
	ta.Execute(w, myvar)
	common.CheckErr(err)
	//}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		fmt.Println("called", username, email)
		repo.EditUserData(w, r, username, email)
	}
	// if r.Method == http.MethodPut {
	// 	//common.AuthenticateJWT(w, r)
	// 	username := r.FormValue("username")
	// 	email := r.FormValue("email")
	// 	repo.EditUserData(w, r, username, email)
	// }
	// if r.Method == http.MethodGet {
	// 	repo.DeleteUserData(w, r)
	// }
	// Finally, return the welcome message to the user, along with their
	// username given in the token
}
func MyPhoto(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	common.AuthenticateJWT(w, r)
	claims := common.GetJWTData(w, r)
	if common.AuthenticateJWT(w, r) == false {
		fmt.Println("error")
	} else {
		url := repo.GetUserPhoto(w, r, claims.Id, j)
		t, err = template.ParseFiles("view/myphoto.html")
		myvar := map[string]interface{}{"url1": url.PhotoUrl, "title": url.Title, "caption": url.Caption}
		t.Execute(w, myvar)
		common.CheckErr(err)
	}
	if r.Method == http.MethodGet {
		url := repo.GetUserPhoto(w, r, claims.Id, j)
		j++
		ta, err := template.ParseFiles("view/myphoto.html")
		myvar := map[string]interface{}{"url1": url.PhotoUrl, "title": url.Title, "caption": url.Caption}
		ta.Execute(w, myvar)
		t.Execute(os.Stdout, myvar)
		common.CheckErr(err)
	} else if r.Method == http.MethodPost {
		url := repo.GetUserPhoto(w, r, claims.Id, j)
		repo.DeletePhoto(w, r, url.PhotoUrl)
		fmt.Println("success")
	}
}
func EditPhoto(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var t *template.Template
	common.AuthenticateJWT(w, r)
	claims := common.GetJWTData(w, r)
	if common.AuthenticateJWT(w, r) == false {
		fmt.Println("error")
	} else {
		url := repo.GetUserPhoto(w, r, claims.Id, j)
		t, err = template.ParseFiles("view/editphoto.html")
		myvar := map[string]interface{}{"url1": url.PhotoUrl, "title": url.Title, "caption": url.Caption}
		t.Execute(w, myvar)
		common.CheckErr(err)
	}
	if r.Method == http.MethodGet {
		url := repo.GetUserPhoto(w, r, claims.Id, j)
		j++
		ta, err := template.ParseFiles("view/editphoto.html")
		myvar := map[string]interface{}{"url1": url.PhotoUrl, "title": url.Title, "caption": url.Caption}
		ta.Execute(w, myvar)
		t.Execute(os.Stdout, myvar)
		common.CheckErr(err)
	} else if r.Method == http.MethodPost {
		title := r.FormValue("title")
		caption := r.FormValue("caption")
		fmt.Println("title :", title, caption)
		url := repo.GetUserPhoto(w, r, claims.Id, j)
		repo.EditPhoto(w, r, title, caption, url.PhotoUrl)
	}
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	http.ServeFile(w, r, "view/comment.html")

	claims := common.GetJWTData(w, r)
	id := repo.GetCommentId(w, r)
	message := r.FormValue("comment")
	photo_id := i
	user_id := claims.Id
	created_at := time.Now()
	fmt.Println(message)
	if r.Method == http.MethodPost {
		repo.PostComment(w, r, id, photo_id, user_id, created_at, message)
	}

}
