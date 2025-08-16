package main

import (
	"fmt"
	"log"
	"os"

	"github.com/scopophobic/scopoBlog/internal/services"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run cmd/generate-password/main.go <password>")
		os.Exit(1)
	}

	password := os.Args[1]

	hash, err := services.GeneratePasswordHash(password)
	if err != nil {
		log.Fatalf("Error generating hash: %v", err)
	}

	fmt.Printf("Password hash: %s\n", hash)
	fmt.Println("Copy this hash to your config.yaml file")
}
