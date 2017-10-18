package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*********************
 * PageHandler takes a unique UserID and PageID and returns a Page that mathes the criteria.
 * The UserID and PageId must be castable as an integer.
 *********************/
func PageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	userID := vars["userID"]
	pageID := vars["pageID"]
	if !IsInt(userID, pageID) {
		return
	}

	page, err := GetPage(pageID, userID)

	if err != nil {
		panic(err)
	} else if err := json.NewEncoder(w).Encode(page); err != nil {
		panic(err)
	}
}

/*********************
 * BookHandler takes a count and offset parameter.
 * Count is defined as the ammount of pages returned with the Book.
 * Offset the position in the database where pages will begin to be read.
 * Count and offset must be castable as an integer.
 *********************/
func BookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)

	if !IsInt(vars["count"], vars["offset"]) {
		return
	}

	count, _ := strconv.Atoi(vars["count"])
	offset, _ := strconv.Atoi(vars["offset"])
	offset *= count

	book, err := GetBook(offset, count)
	if err != nil {

	} else if err := json.NewEncoder(w).Encode(book); err != nil {
		panic(err)
	}
}

/*********************
 * RegisterHandler takes a username, email, password, and accessCode.
 * If the parameters do not reach the standards of FieldCheck() then
 * a custom error object (Params) will be returned explaining the errors.
 * If the paramters reach the standards of FieldCheck() the user will be
 * registered and a unique session string will be returned.
 *********************/
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	accessCode := r.Form.Get("access_code")

	params := FieldCheck(username, email, password, accessCode)
	if params.Valid {
		session, err := Register(username, email, accessCode, password)
		if err != nil {
			fmt.Fprintln(w, `{"valid":false}`)
			panic(err)
		} else {
			fmt.Fprintln(w, `{"session": "`+string(session)+`", "valid":true}`)
		}
	} else {
		if err := json.NewEncoder(w).Encode(params); err != nil {
			panic(err)
		}
	}
}

/*********************
 * LoginHandler takes a key and password. The key can be a username, email,
 * or session. If the database does not contain any matches then an invalid
 * JSON message will be returned. If there is a match in the database, a
 * new session will be attached to the user and returned as a JSON object.
 *********************/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	key := r.Form.Get("key")
	password := r.Form.Get("password")

	valid, session, err := Login(key, password)
	if err != nil {
		panic(err)
	}

	if valid {
		fmt.Fprintln(w, `{"session": "`+string(session)+`", "valid":true}`)
	} else {
		fmt.Fprintln(w, `{"valid":false}`)
	}
}

/*********************
 * UploadHandler recieves a username and session cookie and a name and body form.
 * It will fail if the user is not authenticated or if their is a MySQL error.
 *********************/
func UploadHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()

	name := r.Form.Get("name")
	body := r.Form.Get("body")
	username_cookie, err := r.Cookie("username")
	if err != nil {
		panic(err)
	}
	username := username_cookie.Value
	session_cookie, err := r.Cookie("session")
	if err != nil {
		panic(err)
	}
	session := session_cookie.Value

	valid, session, err := Login(username, session)
	if err != nil {
		panic(err)
	}

	if valid {
		err := AddPage(username, name, body)
		if err != nil{
			panic(err)
		}else{
			fmt.Fprintln(w, `{"session": "`+string(session)+`", "valid":true}`)
		}
	} else {
		fmt.Fprintln(w, `{"valid":false}`)
	}

}
