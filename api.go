package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/soeffing/nlp/downloader"
	"html/template"
	"net/http"
	"strings"
)

// Greeting is simple data structure for static page
type Greeting struct {
	Message string
}

type apiRequestParams struct {
	Urls   []string
	Action string
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: find better way to parse URL params
	// mux.Vars not working with go 1.7
	msg := strings.Split(r.URL.Path, "/")[2]
	data := &Greeting{Message: msg}
	t, _ := template.ParseFiles("tmpl/static.html")
	t.Execute(w, data)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// set proper header
	// TODO: use some sort of pre-hook to set those
	w.Header().Set("Content-Type", "application/json")

	action := strings.Split(r.URL.Path, "/")[2]
	var params apiRequestParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	if action == "downloader" {
		downloader := downloader.New()
		downloader.Download(params.Urls)

		// Try out the json encoder
		// json.NewEncoder(w).Encode(&pages)
		jData, _ := json.Marshal(downloader.Pages)
		w.Write(jData)
	} else if action == "parallel_downloader" {
		fmt.Println("To be implemented...")
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/static/{greeting}", staticHandler)
	r.HandleFunc("/api/{action}", apiHandler).Methods("POST")

	http.ListenAndServe(":8080", nil)
}
