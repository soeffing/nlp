package api

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strings"
)

// Greeting is simple data structure for static page
type Greeting struct {
	Message string
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: find better way to parse URL params
	msg := strings.Split(r.URL.Path, "/")[1]
	data := &Greeting{Message: msg}
	t, _ := template.ParseFiles("tmpl/static.html")
	t.Execute(w, data)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/static/{greeting}", staticHandler)

	http.ListenAndServe(":8080", nil)
}
