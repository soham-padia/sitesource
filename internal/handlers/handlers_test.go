package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"phone", "/phone", "GET", []postData{}, http.StatusOK},
	{"pc", "/pc", "GET", []postData{}, http.StatusOK},
	{"laptop", "/laptop", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"docs", "/docs", "GET", []postData{}, http.StatusOK},
	{"donate", "/donate", "GET", []postData{}, http.StatusOK},
	{"login", "/login", "GET", []postData{}, http.StatusOK},
	{"log-json", "/login-json", "GET", []postData{}, http.StatusOK},
	{"reg", "/register", "GET", []postData{}, http.StatusOK},
	{"reg-summ", "/registration-summary", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
