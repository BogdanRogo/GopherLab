package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	storageService    = "http://localhost:8080" // @todo --> use ENV vars
	storageServiceSet = "/set-key"              // POST json
	storageServiceGet = "/get-key/"             // GET /get-key/2600343750
)

func main() {
	var addr = flag.String("addr", ":8081", "The addr of the application.")
	// /short?url="https://medium.com/metrosystemsro/gitops-with-weave-flux-40997e929254"
	mux := http.NewServeMux()
	mux.HandleFunc("/short", shortHandler)
	mux.HandleFunc("/r/", redirectHandler)
	// mux.HandleFunc("/", homepageHandler)
	mux.Handle("/", http.FileServer(http.Dir("./html")))

	log.Println("Starting application on", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
