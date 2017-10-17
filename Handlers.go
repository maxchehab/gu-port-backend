package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func PageIndex(w http.ResponseWriter, r *http.Request) {
	pages := Pages{
		Page{Name: "Markdown numero uno"},
		Page{Name: "Markdown numer dos"},
	}

	if err := json.NewEncoder(w).Encode(pages); err != nil {
		panic(err)
	}
}

func PageShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	userID := vars["userID"]
	pageID := vars["pageID"]
	if !Sanitized(userID, pageID) {
		return
	}

	conn, _ := SQLConnect()

	conn.Ping()

	m, _ := conn.Execute(`SELECT * FROM pages WHERE pageID=` + pageID + ` AND userID=` + userID)

	name, _ := m.GetString(0, 2)
	body, _ := m.GetString(0, 4)
	author, _ := m.GetString(0, 3)
	page := Page{
		Name:   name,
		Body:   body,
		Author: author,
		PageID: pageID,
		UserID: userID,
	}

	if err := json.NewEncoder(w).Encode(page); err != nil {
		panic(err)
	}
}

func PagePagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)

	if !Sanitized(vars["count"], vars["offset"]) {
		return
	}

	count, _ := strconv.Atoi(vars["count"])
	offset, _ := strconv.Atoi(vars["offset"])
	offset *= count

	conn, _ := SQLConnect()

	conn.Ping()
	m, _ := conn.Execute(`SELECT * FROM pages LIMIT ` + strconv.Itoa(offset) + `,` + strconv.Itoa(count))

	type Book struct {
		Pages  []Page `json:"pages"`
		Count  int    `json:"count"`
		Offset int    `json:"offset"`
		Total  int    `json:"total"`
	}

	book := Book{
		Pages:  []Page{},
		Count:  count,
		Offset: offset,
		Total:  0,
	}

	for i := 0; i < count; i++ {
		name, _ := m.GetString(i, 2)
		author, _ := m.GetString(i, 3)
		body, _ := m.GetString(i, 4)
		pageID, _ := m.GetString(i, 0)
		userID, _ := m.GetString(i, 1)
		page := Page{
			Name:   name,
			Author: author,
			Body:   body,
			PageID: pageID,
			UserID: userID,
		}
		book.Pages = append(book.Pages, page)
	}

	l, _ := conn.Execute(`SELECT COUNT(pageID) FROM pages`)
	total, _ := l.GetString(0, 0)
	book.Total, _ = strconv.Atoi(total)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		panic(err)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	accessCode := r.Form.Get("access_code")

	params := FieldCheck(username, email, password, accessCode)
	if params.Valid {
		password = GenerateHash(password)
		session_b, _ := exec.Command("uuidgen").Output()
		session := stripSpaces(string(session_b))
		conn, _ := SQLConnect()
		conn.Ping()
		conn.Execute(`INSERT INTO users (userID, username, email, access_code, password, session)
					VALUES (NULL, '` + username + `', '` + email + `', '` + accessCode + `', '` + password + `', '` + string(session) + `')`)
		conn.Execute(`UPDATE access_codes SET valid=0 WHERE access_code='` + accessCode + `' LIMIT 1`)
		fmt.Fprintln(w, `{"session": "`+string(session)+`", "valid":true}`)

	} else {
		if err := json.NewEncoder(w).Encode(params); err != nil {
			panic(err)
		}
	}
}

type Validator struct {
	Valid   bool     `json:"valid"`
	Message []string `json:"message"`
}

type Params struct {
	Username   Validator `json:"username"`
	Email      Validator `json:"email"`
	Password   Validator `json:"password"`
	AccessCode Validator `json:"accessCode"`
	Valid      bool      `json:"valid"`
}

func FieldCheck(username string, email string, password string, accessCode string) Params {
	params := Params{
		Valid:      true,
		Username:   Validator{Valid: true},
		Email:      Validator{Valid: true},
		Password:   Validator{Valid: true},
		AccessCode: Validator{Valid: true},
	}

	conn, _ := SQLConnect()

	conn.Ping()
	uM, _ := conn.Execute(`SELECT COUNT(userID) FROM users WHERE username='` + username + `'`)

	userSize_s, _ := uM.GetString(0, 0)
	userSize, _ := strconv.Atoi(userSize_s)
	if userSize > 0 {
		params.Username.Valid = false
		params.Username.Message = append(params.Username.Message, "Username taken.")
		params.Valid = false
	}

	if len(username) == 0 {
		params.Username.Valid = false
		params.Username.Message = append(params.Username.Message, "Username cannot be empty.")
		params.Valid = false
	}

	eM, _ := conn.Execute(`SELECT COUNT(userID) FROM users WHERE email='` + email + `'`)

	emailSize_s, _ := eM.GetString(0, 0)
	emailSize, _ := strconv.Atoi(emailSize_s)
	if emailSize > 0 {
		params.Email.Valid = false
		params.Email.Message = append(params.Email.Message, "Email taken.")
		params.Valid = false
	}

	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !Re.MatchString(email) {
		params.Email.Valid = false
		params.Email.Message = append(params.Email.Message, "Invalid email.")
		params.Valid = false
	}

	if len(password) == 0 {
		params.Password.Valid = false
		params.Password.Message = append(params.Password.Message, "Password cannot be empty.")
		params.Valid = false
	}

	aM, _ := conn.Execute(`SELECT COUNT(access_code) FROM access_codes WHERE access_code ='` + accessCode + `' AND valid=1`)
	accessSize_s, _ := aM.GetString(0, 0)
	accessSize, _ := strconv.Atoi(accessSize_s)
	if accessSize == 0 {
		params.AccessCode.Valid = false
		params.AccessCode.Message = append(params.AccessCode.Message, "Access code invalid.")
		params.Valid = false
	}

	return params
}

func Sanitized(inputs ...string) bool {
	for _, i := range inputs {
		_, e := strconv.Atoi(i)
		if e != nil {
			return false
		}
	}
	return true
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

func GenerateHash(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func CompareHash(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return (err == nil)
}
