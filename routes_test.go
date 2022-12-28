package tinyurl

import "testing"

func TestUrlIsValid(t *testing.T) {
	got, err := urlIsValid("www.google.com/blah")
	if err != nil {
		t.Error("Error validating url")
	}
	if !got {
		t.Error("should be valid")
	}

	got, err = urlIsValid("postgres://localhost")
	if err != nil {
		t.Error("error validating url")
	}
	if got {
		t.Error("should be invalid")
	}

	got, err = urlIsValid("https://localhost:3000")
	if err != nil {
		t.Error("error validating url")
	}
	if !got {
		t.Error("should be invalid")
	}
}

func TestCleanUrl(t *testing.T) {
	got, err := cleanUrl("amazon.com")
	if err != nil {
		t.Error(err)
	}
	if got != "https://amazon.com" {
		t.Error("url not properly cleaned")
	}
}
