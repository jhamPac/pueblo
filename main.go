package main

import (
	"fmt"

	"github.com/jhampac/pueblo/database"
)

func main() {
	t := database.NewAccount("5")
	fmt.Printf("the account %v", t)
}
