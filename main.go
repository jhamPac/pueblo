package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
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
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	bucketName := os.Getenv("STORAGE_BUCKET")

	bucket := client.Bucket(bucketName)
	query := &storage.Query{}
	it := bucket.Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		log.Println(attrs.Name)
	}
	////////////////////////////////////////

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
