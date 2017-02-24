package sparql

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testResults = `
{
   "head": {
       "link": [],
       "vars": [ "abstract", "label" ]
       },
   "results": {
       "bindings": [
               {
                   "abstract" : { "type": "literal", "value": "fake abstract" },
                   "label" : { "type": "literal", "value": "fake label" }
               }
           ]
       }
}`

// TODO: check if there is a cleaner way to
// statusHandler is an http.Handler that writes an empty response using itself
// as the response status code
type statusHandler int

func (h *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(*h))
}

func TestGetLabelAbstractByTermSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/sparql-results+json")
		fmt.Fprintln(w, testResults)
	}))
	defer ts.Close()

	dbpediaURL = ts.URL

	term := "Bitcoin"
	data, err := GetLabelAbstractByTerm(term)

	expectedLabel := "fake label"
	expectedAbstract := "fake abstract"

	if data[0]["label"] != expectedLabel || data[0]["abstract"] != expectedAbstract || err != nil {
		t.Fatalf("GetLabelAbstractByTerm does not return corrcect results")
	}

}

func TestGetLabelAbstractByTermFailure(t *testing.T) {
	status := statusHandler(http.StatusBadRequest)
	ts := httptest.NewServer(&status)

	// ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	dbpediaURL = ts.URL

	term := "Bitcoin"
	_, err := GetLabelAbstractByTerm(term)

	if err == nil {
		t.Fatalf("GetLabelAbstractByTerm does not handle errors from DBpedia correctly")
	}
}

func TestGetLabelAbstractByTermWithEmptyString(t *testing.T) {
	status := statusHandler(http.StatusBadRequest)
	ts := httptest.NewServer(&status)

	// ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	dbpediaURL = ts.URL

	term := ""
	_, err := GetLabelAbstractByTerm(term)

	if err == nil {
		t.Fatalf("GetLabelAbstractByTerm does not return an error when empty string is passed in")
	}
}
