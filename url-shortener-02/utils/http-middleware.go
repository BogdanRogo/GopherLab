package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

// WithMetrics wraps the handler and calculates metrics
func WithMetrics(next http.HandlerFunc) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		// log.Printf("Logged connection from %s", req.RemoteAddr)
		// log.Printf("--> %s %s", req.Method, req.URL.Path)

		startTime := time.Now()
		lwr := newLoggingResponseWriter(wr)
		next.ServeHTTP(lwr, req)

		statusCode := lwr.statusCode
		// log.Printf("<-- %d %s", statusCode, http.StatusText(statusCode))

		ctx, err := tag.New(context.Background(),
			tag.Insert(HTTPMethod, req.Method),
			tag.Insert(HTTPHandler, req.URL.Path),
			tag.Insert(HTTPStatus, fmt.Sprintf("%v", statusCode)),
		)
		if err != nil {
			log.Printf(err.Error())
		}

		stats.Record(ctx, MLatencyMs.M(time.Now().Sub(startTime).Seconds()))
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// func WithLogging(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("Logged connection from %s", r.RemoteAddr)
// 		next.ServeHTTP(w, r)
// 	}
// }

// func WithTracing(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("Tracing request for %s", r.RequestURI)
// 		next.ServeHTTP(w, r)
// 	}
// }
