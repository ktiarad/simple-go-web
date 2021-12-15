package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	id   int
	name string
	city string
	img  string
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang1")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func insertRow(data map[string]string) {
	db, err := connect()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users(`name`,`city`,`img`) VALUES (?, ?, ?)", data["name"], data["city"], data["fileName"])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Insert success!")
}
