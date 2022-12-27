package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
)

type UrlStore interface {
	GetById(int) (string, error)
	Store(string) int
}

type InMemoryUrlStore struct {
	data  []string
	len   int
	mutex sync.Mutex
}

type TinyRequestPayload struct {
	Url string
}

type TinyRequestResponse struct {
	Code string
}

func (s *InMemoryUrlStore) Store(url string) int {
	s.mutex.Lock()
	s.data = append(s.data, url)
	idx := len(s.data)
	s.len = s.len + 1
	s.mutex.Unlock()
	return idx
}

func urlIsValid(uri string) bool {
	_, err := url.Parse(uri)
	return err == nil
}

func (s *InMemoryUrlStore) GetById(id int) (string, error) {
	if id > s.len {
		return "", errors.New("record does not exist")
	}
	return s.data[id-1], nil
}

func NewInMemoryUrlStore() *InMemoryUrlStore {
	return &InMemoryUrlStore{data: make([]string, 0), len: 0}
}

func encodeIdAsString(id int) string {
	var BASE int = 62
	digits := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortUrl := ""

	for id > 0 {
		shortUrl = string(digits[id%62]) + shortUrl

		// id is an int, thus this is floor division
		id = id / BASE
	}
	return shortUrl
}

func decodeStringToId(s string) int {
	var id int = 0
	for _, c := range s {
		if c >= rune('a') && c <= rune('z') {
			id = id*62 + int(c) - int('a')
		} else if c >= rune('A') && c <= rune('Z') {
			id = id*62 + int(c) - int('Z') + 26
		} else {
			id = id*62 + int(c) - int('0') + 52
		}
	}
	return id
}

type App struct {
	store   UrlStore
	version string
}

// func handleRoot(w http.ResponseWriter, r *http.Request) {
func handleRoot(a App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(a.version))
	}
}

func handleTinyfication(a App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var url TinyRequestPayload
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			panic("oh no")
		}
		if !urlIsValid(url.Url) {
			w.WriteHeader(http.StatusBadRequest)
		}
		id := a.store.Store(url.Url)
		code := encodeIdAsString(id)
		codeResponse := TinyRequestResponse{Code: code}
		json.NewEncoder(w).Encode(codeResponse)
	}
}

func handleRedirect(a App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		code := params["code"]
		id := decodeStringToId(code)
		fmt.Printf("%d\n", id)
		url, err := a.store.GetById(id)
		if err != nil {
			log.Panic(err)
			fmt.Print("error")
			w.Write([]byte("not found"))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func main() {
	store := NewInMemoryUrlStore()
	app := App{version: "hello world v1", store: store}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiny/{code}", handleRedirect(app)).Methods("GET")
	router.HandleFunc("/tiny", handleTinyfication(app)).Methods("POST")
	router.HandleFunc("/", handleRoot(app)).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))

}
