package jaguar_test

// Tests for Jaguar
// Uses net/http/httptest for server stubbing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/automattic/go/jaguar"
)

// This tests basic GET request
func TestGetBasic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hola mundo")
	}))
	defer ts.Close()

	j := jaguar.New()
	resp, err := j.Url(ts.URL).Send()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if resp.String() != "hola mundo" {
		t.Errorf("Unexpected result: %v", resp.String())
	}
}

// This tests GET request with passing in a parameter
func TestGetParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.FormValue("p"))
	}))
	defer ts.Close()

	j := jaguar.New()
	j.Params.Add("p", "hello")
	resp, err := j.Url(ts.URL).Send()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if resp.String() != "hello" {
		t.Errorf("Unexpected result: %v", resp.String())
	}
}

// This tests POST request with passing in a parameter
func TestPostParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, r.FormValue("p"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer ts.Close()

	j := jaguar.New()
	j.Params.Add("p", "hello")
	resp, err := j.Post(ts.URL).Send()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Invalid Status Code: %v", resp.StatusCode)
	}

	if resp.String() != "hello" {
		t.Errorf("Unexpected result: %v", resp.String())
	}
}
