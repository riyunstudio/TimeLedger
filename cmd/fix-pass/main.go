package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dsn := "root:timeledger_root_2026@tcp(127.0.0.1:3306)/timeledger?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Generate correct hash for admin123
	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated hash: %s\n", string(hash))

	// Update database
	result, err := db.Exec("UPDATE admin_users SET password_hash = ?", string(hash))
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Rows affected: %d\n", rowsAffected)

	// Verify
	var storedHash string
	db.QueryRow("SELECT password_hash FROM admin_users WHERE email = ?", "admin@timeledger.com").Scan(&storedHash)
	fmt.Printf("Stored hash: %s\n", storedHash)

	// Verify password matches
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Println("Password verification: FAILED")
	} else {
		fmt.Println("Password verification: SUCCESS")
	}
}
