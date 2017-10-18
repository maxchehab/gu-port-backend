package main

import (
	"github.com/siddontang/go-mysql/client"
	"os/exec"
	"errors"
	"strconv"
)

/*********************
 * Standard way to connect to an SQL database.
 *********************/
 func SQLConnect() (*client.Conn, error) {
	return client.Connect(SQLURL, SQLUSER, SQLPASS, DBNAME) // url, user, pass, db
}

/*********************
 * GetPage returns a Page object given an unique PageID and UserID
 *********************/
func GetPage(pageID string, userID string) (Page, error) {
	conn, err := SQLConnect()
	if err != nil {
		return Page{}, err
	}

	query := `SELECT * FROM pages WHERE pageID=` + pageID + ` AND userID=` + userID

	results, err := conn.Execute(query)
	if err != nil {
		return Page{}, err
	}

	name, _ := results.GetString(0, 2)
	body, _ := results.GetString(0, 4)
	author, _ := results.GetString(0, 3)

	conn.Close()

	return Page{
		Name:   name,
		Body:   body,
		Author: author,
		PageID: pageID,
		UserID: userID,
	}, nil
}

/*********************
 * GetBook returns a collection of constructed with the pagination parameters
 * of offset and countcreates a hash with a baked salt from a given password.
 *********************/
func GetBook(offset int, count int) (Book, error) {
	conn, err := SQLConnect()
	if err != nil {
		return Book{}, err
	}

	query := `SELECT * FROM pages LIMIT ` + strconv.Itoa(offset) + `,` + strconv.Itoa(count)
	results, err := conn.Execute(query)
	if err != nil {
		return Book{}, err
	}

	book := Book{
		Pages:  []Page{},
		Count:  count,
		Offset: offset,
		Total:  len(results.Values),
	}

	for i := 0; i < len(results.Values); i++ {
		name, _ := results.GetString(i, 2)
		author, _ := results.GetString(i, 3)
		body, _ := results.GetString(i, 4)
		pageID, _ := results.GetString(i, 0)
		userID, _ := results.GetString(i, 1)
		page := Page{
			Name:   name,
			Author: author,
			Body:   body,
			PageID: pageID,
			UserID: userID,
		}
		book.Pages = append(book.Pages, page)
	}

	return book, nil
}

/*********************
 * Register will append a new user to the `users` database and invalidate
 * the accessCode provided.
 *********************/
func Register(username string, email string, accessCode string, password string) (string, error) {
	password = GenerateHash(password)
	session_b, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}
	session := string(StripSpaces(string(session_b)))

	conn, err := SQLConnect()
	if err != nil {
		return "", err
	}

	query := `INSERT INTO users (userID, username, email, access_code, password, session)
				VALUES (NULL, '` + username + `', '` + email + `', '` + accessCode + `', '` + password + `', '` + session + `')`
	_, err = conn.Execute(query)
	if err != nil {
		return "", err
	}

	query = `UPDATE access_codes SET valid=0 WHERE access_code='` + accessCode + `' LIMIT 1`
	_, err = conn.Execute(query)

	if err != nil {
		return "", err
	}
	return session, nil
}

/*********************
 * Login checks if the provided key and password are contained in the database.
 * If successfuly authenticated, a new session will be appended to the user.
 * This new session is returned.
 *********************/
func Login(key string, password string) (bool, string, error) {
	conn, err := SQLConnect()
	if err != nil {
		return false, "", err
	}

     query := `SELECT userID, password, session FROM users WHERE username='` + key + `' OR email='` + key + `'`
     results, err := conn.Execute(query)
	if err != nil {
		return false, "", err
	}

     userID := int64(-1)
     for i := 0; i < len(results.Values); i++{
          hashedPassword, _ := results.GetString(i,1)
		session, _ := results.GetString(i,2)

          if CompareHash(hashedPassword, password) || session == password{
               userID, _ = results.GetInt(i, 0)
               break;
          }
     }

     if(userID >= 0){
		session_b, err := exec.Command("uuidgen").Output()
		if err != nil {
			return false, "", err
		}
		session := string(StripSpaces(string(session_b)))
		query := `UPDATE users SET session='` + session + `' WHERE userID='` + strconv.Itoa(int(userID)) + `' LIMIT 1`
          _, err = conn.Execute(query)
		if err != nil {
			return false, "", err
		}

          return true, session, nil
     }else{
          return false, "", nil
     }
}

/*********************
 * FieldExists checks if a field of a certain value in a table exists.
 *********************/
func FieldExists(field string, value string, table string)(bool, error){
	conn, err := SQLConnect()
	if err != nil {
		return false, err
	}

	query := `SELECT ` + field + ` FROM ` + table + ` WHERE ` + field + `='` + value + `'`
	results, err := conn.Execute(query)
	if err != nil{
		return false, err
	}
	return (len(results.Values) > 0), nil

}

/*********************
 * AccessCodeValid checks if an accessCode exists and is valid within
 * the `access_codes` database.
 *********************/
func AccessCodeValid(accessCode string) (bool, error){
	conn, err := SQLConnect()
	if err != nil {
		return false, err
	}

	query := `SELECT access_code FROM access_codes WHERE access_code ='` + accessCode + `' AND valid=1`
	results, err := conn.Execute(query)
	if err != nil{
		return false, err
	}
	return (len(results.Values) > 0), nil
}

/*********************
 * AddPage adds a page to a database.
 *********************/

func AddPage(username string, name string, body string) (error){
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}else if userID == "-1"{
		return errors.New("Invalid userID {-1}")
	}

	conn, err := SQLConnect()
	query := `INSERT INTO pages (pageID, userID, name, author, body)
			VALUES (NULL, '` + userID + `', '` + name +`', '` + username + `', '` + body + `')`
	_, err = conn.Execute(query)
	return err
}

/*********************
 * Given a username, GetUserID returns the userID as a string
 * or "-1" if the username is not found.
 *********************/
func GetUserID(username string) (string, error){
	conn, err := SQLConnect()
	if err != nil {
		return "-1", err
	}
	query := `SELECT userID FROM users WHERE username='` + username + `'`
	results, err := conn.Execute(query)
	if err != nil {
		return "-1", err
	}

	if len(results.Values) > 0{
		i, err := results.GetString(0,0)
		if err != nil {
			return "-1", err
		}else{
			return i, nil
		}
	}else{
		return "-1", nil
	}
}
