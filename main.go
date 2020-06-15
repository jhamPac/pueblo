package main

import (
	"log"

	"github.com/jhampac/pueblo/node"
)

func main() {
	if err := node.Run("pwd"); err != nil {
		log.Fatalf("An error as occured: %v", err)
	}
}
