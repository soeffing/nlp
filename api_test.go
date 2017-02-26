package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/soeffing/nlp/downloader"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StaticHTTPHandler struct{}
type DownloadHTTPHandler struct{}

func (h *StaticHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	staticHandler(w, r)
}

func (h *DownloadHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	downloadHandler(w, r)
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

func TestDownloadHandlerAvailability(t *testing.T) {
	handler := &DownloadHTTPHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	var urls []string
	urls = append(urls, "https://blog.golang.org/go-maps-in-action")
	params := downloadRequestParams{urls}
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

func TestDownloadHandlerResponse(t *testing.T) {
	handler := &DownloadHTTPHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	var urls []string
	urls = append(urls, "https://blog.golang.org/go-maps-in-action")

	params := downloadRequestParams{urls}

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

func TestDownloadHandlerResponseWithHTTPMock(t *testing.T) {
	// Setup the server
	handler := &DownloadHTTPHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	// Setup the
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text")
		fmt.Fprint(w, "fake text")
	}))

	defer ts.Close()

	var urls []string
	urls = append(urls, ts.URL)

	params := downloadRequestParams{urls}

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

	// Double check if http mock is actually working
	expectedContent := "fake text"
	actualContent := pages[0].Content
	if strings.Compare(expectedContent, actualContent) != 0 {
		t.Fatalf("Http mock not working")
	}

	expected := 1
	actual := len(pages)

	if expected != actual {
		t.Fatalf("Downloader API endpoint does not return pages")
	}
}

// TODO: find out how I can mock the DBpedia call inside the sparql package
func TestSparqlHandlerAvailability(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/sparql", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Add "term" request parameter
	// TODO: can this be done more elegantly?
	q := req.URL.Query()
	q.Add("term", "Bitcoin")
	req.URL.RawQuery = q.Encode()

	// Setup response recoder and handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sparqlHandler)

	// I can call ServeHttp because our handler satisfies http.Handler
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRakeHandlerAvailability(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/rake", nil)

	if err != nil {
		t.Fatal(err)
	}

	// TODO: can this be done more elegantly?
	q := req.URL.Query()
	q.Add("text", "This is a tiny test. A test for rake, a keyword extraction technology.")
	req.URL.RawQuery = q.Encode()

	// Setup response recoder and handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rakeHandler)

	// I can call ServeHttp because our handler satisfies http.Handler
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
