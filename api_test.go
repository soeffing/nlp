package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StaticHTTPHandler struct{}

func (h *StaticHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	staticHandler(w, r)
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

	resp, err := http.Get(server.URL + "/uli")
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

	resp, err := http.Get(server.URL + "/uli")
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	AssertSee(t, string(body), "Hello, uli")
}
