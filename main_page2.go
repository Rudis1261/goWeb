package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"time"
)

type appContext struct {
	db *sql.DB
}

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

func (c *appContext) authHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		//authToken := req.Header.Get("Authorization")
		/*user, err := getUser(c.db, authToken)

		if err != nil {
			http.Error(rw, http.StatusText(401), 401)
			return
		}
		*/
		context.Set(req, "user", {id: 1000})
		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

func (c *appContext) adminHandler(rw http.ResponseWriter, req *http.Request) {
	user := context.Get(req, "user")
	json.NewEncoder(rw).Encode(user)
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/serge")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	appC := appContext{db}

	requestHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)
	http.Handle("/", requestHandlers.ThenFunc(indexHandler))
	http.Handle("/about", requestHandlers.ThenFunc(aboutHandler))
	http.Handle("/admin", requestHandlers.Append(appC.authHandler).ThenFunc(appC.adminHandler))
	http.ListenAndServe(":8080", nil)
}
