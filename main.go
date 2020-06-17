package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

	var names []string

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		names = append(names, attrs.Name)
	}
	////////////////////////////////////////

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		indexHandler(w, names)
	})

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

func indexHandler(w io.Writer, fileNames []string) {
	s := strings.Join(fileNames, " ")
	fmt.Fprintf(w, "File name is: %s", s)
}
