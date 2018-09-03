package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"hash/adler32"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	storageService    = "http://localhost:8080" // @todo --> use ENV vars
	storageServiceSet = "/set-key"              // POST json
	storageServiceGet = "/get-key/"             // GET /get-key/2600343750
)

type storageStruct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func hashURL(url string) uint32 {
	// returns the hash of the url
	const Size = 4
	// return crc32.ChecksumIEEE([]byte(url)) // CRC
	return adler32.Checksum([]byte(url)) //ADLER
}

func newHTTPClient() *http.Client {
	// preconfigured HTTP client
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout: 1 * time.Second,
	}
}

// func checkErr(writer http.ResponseWriter, err error, message string, statusCode int) {
// 	if err != nil {
// 		log.Println(err.Error())
// 		http.Error(writer, message, statusCode)
// 		return
// 	}
// }

func shortHandler(wr http.ResponseWriter, req *http.Request) {

	urls, ok := req.URL.Query()["url"] // Get a copy of the query values.
	if !ok || len(urls[0]) < 1 {
		http.Error(wr, "error missing url", http.StatusBadRequest)
		return
	}

	url, err := url.ParseRequestURI(urls[0])
	if err != nil {
		http.Error(wr, "error does not seem to be an url", http.StatusBadRequest)
		return
	}

	ss := storageStruct{
		Key:   fmt.Sprint(hashURL(url.String())),
		Value: url.String(),
	}
	ssJSON, _ := json.Marshal(ss)
	log.Printf("%v\n", string(ssJSON))

	storageServiceReq, err := http.NewRequest(http.MethodPost, storageService+storageServiceSet, bytes.NewBuffer(ssJSON))
	if err != nil {
		http.Error(wr, "Oops... something unknown happened :)", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	client := newHTTPClient()
	resp, err := client.Do(storageServiceReq)
	if err != nil {
		http.Error(wr, "Oops... could not contact backing service", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	if resp.StatusCode == http.StatusOK {
		wr.WriteHeader(http.StatusOK)
		wr.Write([]byte(req.Host + "/r/" + ss.Key))
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
	reqURL := storageService + storageServiceGet + key
	// fmt.Printf("%v, %v \n", key, reqURL)

	storageServiceReq, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		http.Error(wr, "Oops... something unknown happened :)", http.StatusInternalServerError)
		log.Printf(err.Error())
		return
	}

	client := newHTTPClient()
	resp, err := client.Do(storageServiceReq)
	if err != nil {
		log.Println(err.Error())
		http.Error(wr, "Oops... could not contact backing service", http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(wr, "Ooops... could not contact backing service", http.StatusInternalServerError)
		return
	}
	// log.Printf("%v", string(body))

	var message map[string]interface{}
	_ = json.Unmarshal(body, &message) // handle error

	redirectURL := message["value"].(string) // .(string) type assertion
	if redirectURL == "" {
		log.Println("received key value is empty")
		return
	}
	http.Redirect(wr, req, redirectURL, http.StatusMovedPermanently)

}

func main() {
	var addr = flag.String("addr", ":8081", "The addr of the application.")
	// /short?url="https://medium.com/metrosystemsro/gitops-with-weave-flux-40997e929254"

	mux := http.NewServeMux()
	mux.HandleFunc("/short", shortHandler)
	mux.HandleFunc("/r/", redirectHandler)

	log.Println("Starting application on", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
