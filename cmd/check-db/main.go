package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id int
	var email, hash string
	err = db.QueryRow("SELECT id, email, password_hash FROM admin_users WHERE email = ?", "admin@timeledger.com").Scan(&id, &email, &hash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %d\n", id)
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Hash: %s\n", hash)
}
