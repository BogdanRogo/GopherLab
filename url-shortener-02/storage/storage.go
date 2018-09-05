package storage

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Config contains the config for the backend store service
type Config struct {
	Addr string //= "http://localhost:8080" // @todo --> use ENV vars
	Set  string //= "/set-key"              // POST json
	Get  string //= "/get-key/"             // GET /get-key/2600343750
}

type storageStruct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func newHTTPClient() *http.Client {
	// preconfigured HTTP client
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout: 1 * time.Second,
	}
}

// NewStorageKey ...
func (sc *Config) NewStorageKey(key, value string) ([]byte, error) {
	ss := storageStruct{
		Key:   key,
		Value: value,
	}
	storageStructJSON, err := json.Marshal(ss)
	if err != nil {
		return nil, err
	}
	return storageStructJSON, nil
}

// StorageSet ...
func (sc *Config) StorageSet(data []byte) (bool, error) {
	storageServiceReq, err := http.NewRequest(http.MethodPost, sc.Addr+sc.Set, bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}

	client := newHTTPClient()
	resp, err := client.Do(storageServiceReq)
	if err != nil {
		return false, err
	}
	log.Println(resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, fmt.Errorf("%v", resp.StatusCode)
}

// StorageGet ...
func (sc *Config) StorageGet(key string) ([]byte, error) {

	reqURL := sc.Addr + sc.Get + key
	// fmt.Printf("%v, %v \n", key, reqURL)

	storageServiceReq, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		log.Printf(err.Error())
		return nil, fmt.Errorf("Oops... something unknown happened :)")
	}

	client := newHTTPClient()
	resp, err := client.Do(storageServiceReq)
	if err != nil {
		log.Println(err.Error())
		return nil, fmt.Errorf("Oops... could not contact backing service")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, fmt.Errorf("Oops... could not contact backing service")
	}
	// log.Printf("%v", string(body))
	return body, nil
}

// DecodeStorageData ...
func (sc *Config) DecodeStorageData(data []byte) (string, error) {
	var message map[string]interface{}
	err := json.Unmarshal(data, &message) // handle error
	if err != nil {
		return "", err
	}

	redirectURL := message["value"].(string) // .(string) type assertion
	if redirectURL == "" {
		return "", fmt.Errorf("received key value is empty")
	}
	return redirectURL, nil
}
