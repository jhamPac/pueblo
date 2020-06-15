package main

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/jhampac/pueblo/node"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

func main() {
	if err := node.Run(basePath + "/database/.db"); err != nil {
		log.Fatalf("An error as occured: %v", err)
	}

}
