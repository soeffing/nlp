package api

import (
	"bytes"
	"container/list"
	"encoding/json"
	"github.com/soeffing/nlp/downloader"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type APIResponse struct {
	pages list.List
}

type StaticHTTPHandler struct{}
type APIHTTPHandler struct{}

func (h *StaticHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	staticHandler(w, r)
}

func (h *APIHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	apiHandler(w, r)
}

// Helper function to check for content in body
func AssertSee(t *testing.T, body, query string) {
	if strings.Index(body, query) == -1 {
		t.Errorf("Expected to see %s in\n%s", query, body)
	}
}

func TestStaticHandlerAvailability(t *testing.T) {
	handler := &StaticHTTPHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/static/uli")
	if err != nil {
		t.Fatal(err)
	}

	expected := 200

	if expected != resp.StatusCode {
		t.Fatalf("Endpoint does return non 200 response code: %d\n", resp.StatusCode)
	}
}

func TestStaticHandlerTemplateOutput(t *testing.T) {
	handler := &StaticHTTPHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/static/uli")
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	AssertSee(t, string(body), "Hello, uli")
}

func TestApiHandlerAvailability(t *testing.T) {
	handler := &APIHTTPHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	var urls []string
	urls = append(urls, "https://blog.golang.org/go-maps-in-action")
	params := apiRequestParams{urls, "action"}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(params)
	res, err := http.Post(server.URL+"/api/downloader", "application/json; charset=utf-8", buffer)

	if err != nil {
		t.Fatal(err)
	}

	expected := 200

	if expected != res.StatusCode {
		t.Fatalf("Endpoint does return non 200 response code: %d\n", res.StatusCode)
	}
}

func TestApiHandlerResponse(t *testing.T) {
	handler := &APIHTTPHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	var urls []string
	urls = append(urls, "https://blog.golang.org/go-maps-in-action")

	params := apiRequestParams{urls, "action"}

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(params)
	res, _ := http.Post(server.URL+"/api/downloader", "application/json; charset=utf-8", buffer)
	defer res.Body.Close()

	var pages []downloader.Page
	data, _ := ioutil.ReadAll(res.Body)

	jsonErr := json.Unmarshal(data, &pages)

	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	expected := 1
	actual := len(pages)

	if expected != actual {
		t.Fatalf("Downloader does not return pages")
	}
}
