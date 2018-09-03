package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ContactDetails struct {
	Url string
}

// func checkErr(writer http.ResponseWriter, err error, message string, statusCode int) {
// 	if err != nil {
// 		log.Println(err.Error())
// 		http.Error(writer, message, statusCode)
// 		return
// 	}
// }

func shortHandler(wr http.ResponseWriter, req *http.Request) {

	urls, ok := req.URL.Query()["url"] // Get a copy of the queried value.
	if !ok || len(urls[0]) < 1 {
		http.Error(wr, ReturnError("missing url"), http.StatusBadRequest)
		return
	}

	url, err := url.ParseRequestURI(urls[0])
	if err != nil {
		http.Error(wr, ReturnError("failed to parse URL"), http.StatusBadRequest)
		return
	}

	urlHash := dataHash(url.String())
	ssJSON, err := NewStorageKey(urlHash, url.String())
	if err != nil {
		log.Printf(err.Error())
		http.Error(wr, ReturnError("Oops... JSONs"), http.StatusInternalServerError)
		return
	}

	ok, err = StorageSet(ssJSON)
	if err != nil {
		log.Printf(err.Error())
		http.Error(wr, ReturnError("Oops... could not contact backing service"), http.StatusInternalServerError)
		return
	}

	if ok {
		wr.WriteHeader(http.StatusOK)
		wr.Write(ReturnURL(req.Host + "/r/" + urlHash))
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
	storageData, err := StorageGet(key)
	if err != nil {
		log.Printf(err.Error())
		http.Error(wr, ReturnError("Oops... Backing services"), http.StatusInternalServerError)
	}
	redirectURL, _ := DecodeStorageData(storageData)
	if err != nil {
		log.Printf(err.Error())
		http.Error(wr, ReturnError("Oops... url not in our DB"), http.StatusBadRequest)
	}

	http.Redirect(wr, req, redirectURL, http.StatusMovedPermanently)
}

// func homepageHandler(wr http.ResponseWriter, req *http.Request) {
// 	t, _ := template.ParseFiles("templates/index.html")
// 	details := ContactDetails{
// 		Url: req.FormValue("url"),
// 	}
// 	_ = details
// 	t.Execute(wr, "Hello")
// }
