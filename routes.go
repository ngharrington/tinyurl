package tinyurl

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type HandleFunc = func(http.ResponseWriter, *http.Request)

func urlIsValid(uri string) bool {
	_, err := url.Parse(uri)
	return err == nil
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

func MakeHandleFunction(a App, fn func(App, http.ResponseWriter, *http.Request)) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(a, w, r)
	}
}

func HandleRoot(a App, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(a.version))
}

func HandleTinyfication(a App, w http.ResponseWriter, r *http.Request) {
	var url TinyRequestPayload
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		panic("oh no")
	}
	if !urlIsValid(url.Url) {
		w.WriteHeader(http.StatusBadRequest)
	}
	id, _ := a.store.Store(url.Url)
	code := encodeIdAsString(id)
	codeResponse := TinyRequestResponse{Code: code}
	json.NewEncoder(w).Encode(codeResponse)
}

func HandleRedirect(a App, w http.ResponseWriter, r *http.Request) {
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
