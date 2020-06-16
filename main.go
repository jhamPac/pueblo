package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// import (
// 	"log"
// 	"path/filepath"
// 	"runtime"

// 	"github.com/jhampac/pueblo/node"
// )

// var (
// 	_, b, _, _ = runtime.Caller(0)
// 	basePath   = filepath.Dir(b)
// )

// func main() {
// 	if err := node.Run(basePath + "/database/.db"); err != nil {
// 		log.Fatalf("An error as occured: %v", err)
// 	}

// }

var secret string

func main() {
	secret = os.Getenv("SECRET")
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "The secret is %v", secret)
}
