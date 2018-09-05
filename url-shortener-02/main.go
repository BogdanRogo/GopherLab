package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"

	"git.metrosystems.net/reliability-engineering/reliability-sandbox/GopherLab/url-shortener/storage"
	"git.metrosystems.net/reliability-engineering/reliability-sandbox/GopherLab/url-shortener/utils"
)

var (
	// StoreConfig is the backend service config
	StoreConfig storage.Config
)

func init() {
	StoreConfig.Addr = "http://localhost:8080"
	StoreConfig.Set = "/set-key"
	StoreConfig.Get = "/get-key/"
}

func shortHandler(wr http.ResponseWriter, req *http.Request) {

	urls, ok := req.URL.Query()["url"] // Get a copy of the queried value.
	if !ok || len(urls[0]) < 1 {
		http.Error(wr, utils.ReturnError("missing url"), http.StatusBadRequest)
		return
	}

	url, err := url.ParseRequestURI(urls[0])
	if err != nil {
		http.Error(wr, utils.ReturnError("failed to parse URL"), http.StatusBadRequest)
		return
	}

	urlHash := utils.DataHash(url.String())
	ssJSON, err := StoreConfig.NewStorageKey(urlHash, url.String())
	if err != nil {
		log.Printf(err.Error())
		http.Error(wr, utils.ReturnError("Oops... JSONs"), http.StatusInternalServerError)
		return
	}
	log.Printf("%v", string(ssJSON))
	ok, err = StoreConfig.StorageSet(ssJSON)
	if err != nil {
		log.Printf(err.Error())
		http.Error(wr, utils.ReturnError("Oops... could not contact backing service"), http.StatusInternalServerError)
		return
	}

	if ok {
		wr.WriteHeader(http.StatusOK)
		wr.Write(utils.ReturnURL(req.Host + "/r/" + urlHash))
	}
}

func redirectHandler(wr http.ResponseWriter, req *http.Request) {
	// fmt.Println(req.URL.Path)
	p := strings.Split(req.URL.Path, "/")[1:] // get the keys from 1 to n

	if len(p) < 2 {
		http.Error(wr, "missing key", http.StatusNotFound)
		log.Printf("Key not found in url path")
		return
	}
	key := p[1]
	storageData, err := StoreConfig.StorageGet(key)
	if err != nil {
		log.Printf(err.Error())
		http.Error(wr, utils.ReturnError("Oops... Backing services"), http.StatusInternalServerError)
	}
	redirectURL, _ := StoreConfig.DecodeStorageData(storageData)
	if err != nil {
		log.Printf(err.Error())
		http.Error(wr, utils.ReturnError("Oops... url not in our DB"), http.StatusBadRequest)
	}

	http.Redirect(wr, req, redirectURL, http.StatusMovedPermanently)
}

func main() {

	var addr = flag.String("addr", ":8081", "The addr of the application.")
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./html")))
	mux.Handle("/metrics", utils.OCPrometheusExporter())

	mux.HandleFunc("/short", utils.WithMetrics(shortHandler))
	mux.HandleFunc("/r/", redirectHandler)

	log.Println("Starting application on", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
