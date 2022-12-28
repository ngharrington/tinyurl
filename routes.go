package tinyurl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type HandleFunc = func(http.ResponseWriter, *http.Request)

func cleanUrl(uri string) (string, error) {
	u, err := url.Parse(uri)
	fmt.Printf("host %s", u.Host)
	if err != nil {
		return "", errors.New("error parsing url")
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	return u.String(), nil
}

func urlIsValid(uri string) (bool, error) {
	u, err := url.Parse(uri)
	if u.Scheme == "" {
		u.Scheme = "https"
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return false, nil
	} else if u.Host == "localhost" {
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
	if len(shortUrl) < 6 {
		fmt.Println("here")
		diff := 6 - len(shortUrl)
		shortUrl = strings.Repeat(string(digits[0]), diff) + shortUrl
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
	var uri string
	err := json.NewDecoder(r.Body).Decode(&u)
	uri = u.Url
	if err != nil {
		panic("oh no")
	}
	uri, err = cleanUrl(uri)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	isValid, err := urlIsValid(uri)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := a.store.Store(uri)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusFound)
}
