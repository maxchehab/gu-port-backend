package main

import "github.com/siddontang/go-mysql/client"

// The standard way to connect to the MySQL client
func SQLConnect() (*client.Conn, error) {
	return client.Connect(SQLUrl, "gu-port", "gu-port", DBName) // url, user, pass, db
}
