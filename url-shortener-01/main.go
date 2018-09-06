package main

import (
	"log"
	"net/http"
)

// requirements:
// - handlers:
//   - main (root)  - awesome html + js page
//   - short        - get the url
//					- validate the url
//					- hash it
//					- send it to storage service (redis-service )
//   - redirect     - read the key from /r/{key}
// 					- query the storage and get the url
// 					- reply with status code moved permanently and send the url to the user
//

func healthHandler(wr http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		wr.Header().Set("Access-Control-Allow-Origin", origin)
	}
	wr.Header().Set("Access-Control-Allow-Methods", "GET")
	wr.Header().Set("Content-Type", "application/json")
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte(`{"Status": "OK"}`))
	// io.WriteString(wr, `{"Status": OK}`)
}

// shortener handler
func shortenerHandler(wr http.ResponseWriter, req *http.Request) {
	log.Println("shortenerHandler called")
	wr.WriteHeader(http.StatusOK)
}

// redirect handler
func redirectHandler(wr http.ResponseWriter, req *http.Request) {
	log.Println("shortenerHandler called")
	wr.WriteHeader(http.StatusOK)
}

// main function
func main() {

	var appPort = ":8081"

	// we will use default go package: https://golang.org/pkg/net/http/#ServeMux
	// and crete a ne multiplexer http server
	mux := http.NewServeMux()

	// root handler will redirect user to www folder and act as a classic http server
	mux.Handle("/", http.FileServer(http.Dir("./www")))

	// shortener handler will not allow child routes as it does not have a final '/'
	mux.HandleFunc("/s", shortenerHandler)

	// redirect handler will allow child routes as it does have a final '/'
	mux.HandleFunc("/r/", redirectHandler)

	// health checks handlers
	// https://reliability.metrosystems.net/docs/dev-best-practices/kubernetes/health-checks/
	mux.HandleFunc("/.well-known/live", healthHandler)
	mux.HandleFunc("/.well-known/ready", healthHandler)

	// start the http server
	log.Println("Starting application on", appPort)
	if err := http.ListenAndServe(appPort, mux); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
