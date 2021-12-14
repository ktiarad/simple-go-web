package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/form", handlerForm)
	http.HandleFunc("/process", routeSubmitPost)
	http.HandleFunc("/download", handlerDownload)
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(http.Dir("assets")))))

	var address = "localhost:9000"
	fmt.Printf("Server started at %s\n", address)
	http.ListenAndServe(address, nil)
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles(
		"views/index.html",
		"views/_header.html",
	))
	var data = map[string]interface{}{
		"title": "Index",
	}

	var err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
