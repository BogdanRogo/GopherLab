package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.metrosystems.net/reliability-engineering/reliability-sandbox/GopherLab/url-shortener/storage"
	"git.metrosystems.net/reliability-engineering/reliability-sandbox/GopherLab/url-shortener/utils"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// StoreConfig is the backend service config
	StoreConfig storage.StorageConfig
	// MLatencyMs The latency in milliseconds
	MLatencyMs = stats.Float64("urlshortener/latency", "The latency in milliseconds per short request", "ms")
	// MErrors Encounters the number of non EOF(end-of-file) errors.
	MErrors = stats.Int64("urlshortener/errors", "The number of errors encountered", "1")
	// HTTPMethod ...
	HTTPMethod, _ = tag.NewKey("method")
	// HTTPHandler ...
	HTTPHandler, _ = tag.NewKey("handler")
	ctx            context.Context
)

var (
	LatencyView = &view.View{
		Name:        "urlshortener/latency",
		Measure:     MLatencyMs,
		Description: "The distribution of the latencies",

		// Latency in buckets:
		// [>=5ms, >=10ms ... >=4s, >=6s]
		Aggregation: view.Distribution(5, 10, 15, 20, 25, 50, 75, 100, 250, 500, 750, 1000, 5000),
		TagKeys:     []tag.Key{HTTPMethod, HTTPHandler}}

	ErrorCountView = &view.View{
		Name:        "urlshortener/errors",
		Measure:     MErrors,
		Description: "The number of errors encountered",
		Aggregation: view.Count(),
	}
)

func init() {
	StoreConfig.Addr = "http://localhost:8080"
	StoreConfig.Set = "/set-key"
	StoreConfig.Get = "/get-key/"
}

func shortHandler(ctx context.Context, wr http.ResponseWriter, req *http.Request) {
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
	ctx = context.Background()
	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatal(err)
	}
	view.RegisterExporter(exporter)
	if err := view.Register(LatencyView, ErrorCountView); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}
	view.SetReportingPeriod(1 * time.Second)

	var addr = flag.String("addr", ":8081", "The addr of the application.")
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./html")))
	mux.Handle("/metrics", exporter)

	mux.HandleFunc("/short", func(wr http.ResponseWriter, req *http.Request) {
		ctx, err := tag.New(context.Background(),
			tag.Insert(HTTPMethod, req.Method),
			tag.Insert(HTTPHandler, "short"),
		)
		if err != nil {
			log.Printf(err.Error())
		}
		startTime := time.Now()
		shortHandler(ctx, wr, req)
		log.Printf("%v", wr.Header().Get("method"))
		stats.Record(ctx, MLatencyMs.M(time.Now().Sub(startTime).Seconds()))
	})

	mux.HandleFunc("/r/", redirectHandler)

	log.Println("Starting application on", *addr)
	if err := http.ListenAndServe(*addr, &ochttp.Handler{
		Handler:     mux,
		Propagation: &b3.HTTPFormat{},
	}); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
