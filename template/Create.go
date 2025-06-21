package main

import (
	"fmt"
	"log"
	"github.com/ewriq/pouch" 
)

func main() {
	fmt.Println("--- Create For MySql Container ---")
	opt := pouch.CreateOptions{
		Name:  "my-nginx-dssdsdsdsd",
		Image: "mysql:8.0",
		Port:  "8080", 
	}
	id, err := pouch.Create(opt)
	if err != nil {
		log.Printf("Error for creating container: %v", err)
	} else {
		fmt.Printf("Container ID: %s\n", id)
	}
	fmt.Println("---")
}