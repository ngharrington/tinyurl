package tinyurl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type HandleFunc = func(http.ResponseWriter, *http.Request)

func urlIsValid(uri string) (bool, error) {
	u, err := url.Parse(uri)
	if u.Scheme == "" {
		u.Scheme = "https"
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return false, nil
	}
	if err != nil {
		return false, errors.New("error parsing url")
	}

	return true, nil
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
	var u TinyRequestPayload
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		panic("oh no")
	}
	isValid, err := urlIsValid(u.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := a.store.Store(u.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("generating small url"))
		return
	}
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
