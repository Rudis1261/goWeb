package main

import (
	"fmt"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"time"
)

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Welcome!")
}

func aboutHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Welcome to the about page")
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		end := time.Now()
		log.Printf("[%s] %q %v\n", req.Method, req.URL.String(), end.Sub(start))
	}

	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %+v", err)
				http.Error(rw, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

func main() {
	requestHandlers := alice.New(loggingHandler, recoverHandler)
	http.Handle("/", requestHandlers.ThenFunc(indexHandler))
	http.Handle("/about", requestHandlers.ThenFunc(aboutHandler))
	http.ListenAndServe(":8080", nil)
}
