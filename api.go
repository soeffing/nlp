package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/soeffing/nlp/downloader"
	"github.com/soeffing/nlp/sparql"
	"html/template"
	"net/http"
	"strings"
)

// Greeting is simple data structure for static page
type Greeting struct {
	Message string
}

type downloadRequestParams struct {
	Urls []string
}

type sparqlRequestParams struct {
	Term string
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: find better way to parse URL params
	// mux.Vars not working with go 1.7
	msg := strings.Split(r.URL.Path, "/")[2]
	data := &Greeting{Message: msg}
	t, _ := template.ParseFiles("tmpl/static.html")
	t.Execute(w, data)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// set proper header
	// TODO: use some sort of pre-hook to set those
	w.Header().Set("Content-Type", "application/json")

	var params downloadRequestParams
	err := json.NewDecoder(r.Body).Decode(&params)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	downloader := downloader.New()
	downloader.Download(params.Urls)

	// Try out the json encoder
	// json.NewEncoder(w).Encode(&pages)
	jData, _ := json.Marshal(downloader.Pages)
	w.Write(jData)
}

func sparqlHandler(w http.ResponseWriter, r *http.Request) {
	// set proper header
	// TODO: use some sort of pre-hook to set those
	w.Header().Set("Content-Type", "application/json")

	term := r.URL.Query().Get("term")

	data, err := sparql.GetLabelAbstractByTerm(term)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jData, _ := json.Marshal(data)
	w.Write(jData)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/static/{greeting}", staticHandler)
	r.HandleFunc("/api/download", downloadHandler).Methods("POST")
	r.HandleFunc("/api/sparql", sparqlHandler).Methods("GET")

	http.ListenAndServe(":8080", nil)
}
