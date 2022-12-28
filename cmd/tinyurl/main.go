package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ngharrington/tinyurl"
	"github.com/ngharrington/tinyurl/store"
	"github.com/rs/cors"
)

func main() {
	store, _ := store.NewSqliteUrlStore("./db.sqlite")
	app := tinyurl.NewApp("v0", store)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
	})
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiny/{code}", tinyurl.MakeHandleFunction(app, tinyurl.HandleRedirect)).Methods("GET")
	router.HandleFunc("/tiny", tinyurl.MakeHandleFunction(app, tinyurl.HandleTinyfication)).Methods("POST")
	router.HandleFunc("/", tinyurl.MakeHandleFunction(app, tinyurl.HandleRoot)).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", c.Handler(router)))

}
