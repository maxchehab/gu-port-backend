package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/siddontang/go-mysql/client"
	"golang.org/x/crypto/bcrypt"
)

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

	conn, _ := client.Connect("104.236.141.69:3306", "gu-port", "gu-port", "gu-port")

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
