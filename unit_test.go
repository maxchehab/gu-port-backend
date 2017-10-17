package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


/*********************
 * Helpers.go unit testing
 *********************/
func Test_IsInt(t *testing.T) {
	assert.Equal(t, true, IsInt("142134123"))
	assert.Equal(t, false, IsInt("a"))
	assert.Equal(t, true, IsInt("1", "-1", "0"))
	assert.Equal(t, false, IsInt("not int", "-1", "0"))
}

func Test_StripSpaces(t *testing.T) {
	assert.Equal(t, "helloworld", StripSpaces("hello world"))
	assert.Equal(t, "hello", StripSpaces("h e l l o"))
	assert.Equal(t, "helloworld", StripSpaces("helloworld"))
}

func TestAccessCode_Hashing(t *testing.T) {
	password := "password"
	hashedPassword := GenerateHash(password)

	assert.Equal(t, true, CompareHash(hashedPassword, "password"))
	assert.Equal(t, false, CompareHash(hashedPassword, "test"))
	assert.Equal(t, false, CompareHash(hashedPassword, ""))
}

func Test_FieldCheck(t *testing.T) {
	params := FieldCheck("test","test@gmail.com","","test");
	assert.Equal(t, false, params.Valid)
	assert.Equal(t, false, params.Username.Valid)
	assert.Equal(t, "Username taken.", params.Username.Message[0])
	assert.Equal(t, false, params.Email.Valid)
	assert.Equal(t, "Email taken.", params.Email.Message[0])
	assert.Equal(t, "Username taken.", params.Username.Message[0])
	assert.Equal(t, false, params.Email.Valid)
	assert.Equal(t, "Email taken.", params.Email.Message[0])
	assert.Equal(t, false, params.Password.Valid)
	assert.Equal(t, "Password cannot be empty.", params.Password.Message[0])
	assert.Equal(t, false, params.AccessCode.Valid)
	assert.Equal(t, "Access code invalid.", params.AccessCode.Message[0])

	params = FieldCheck("","testgmail.com","","");
	assert.Equal(t, false, params.Valid)
	assert.Equal(t, false, params.Username.Valid)
	assert.Equal(t, "Username cannot be empty.", params.Username.Message[0])
	assert.Equal(t, false, params.Email.Valid)
	assert.Equal(t, "Invalid email.", params.Email.Message[0])
	assert.Equal(t, false, params.Password.Valid)
	assert.Equal(t, "Password cannot be empty.", params.Password.Message[0])
	assert.Equal(t, false, params.AccessCode.Valid)
	assert.Equal(t, "Access code invalid.", params.AccessCode.Message[0])

	params = FieldCheck("unique", "unique@gmail.com", "password", "unique")
	assert.Equal(t, true, params.Valid)
}
