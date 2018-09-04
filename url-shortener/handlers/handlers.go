package handlers

import (
	"git.metrosystems.net/reliability-engineering/reliability-sandbox/GopherLab/url-shortener/storage"
)

var (
	// StoreConfig is the backend service config
	StoreConfig storage.StorageConfig
)

func init() {
	StoreConfig.Addr = "http://localhost:8080"
	StoreConfig.Set = "/set-key"
	StoreConfig.Get = "/get-key/"
}
