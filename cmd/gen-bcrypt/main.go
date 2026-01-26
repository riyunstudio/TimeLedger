package main

import (
	"fmt"
	"os"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: go run gen-bcrypt/main.go <password>")
		os.Exit(1)
	}

	password := args[0]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(hash))
}
