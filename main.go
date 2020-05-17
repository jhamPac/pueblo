package main

import (
	"fmt"
	"log"

	"github.com/jhampac/pueblo/database"
)

func main() {
	s, err := database.NewStateFromDisk()
	if err != nil {
		log.Fatal("Error in starting the state")
	}

	fmt.Printf("the account %v", s)
}
