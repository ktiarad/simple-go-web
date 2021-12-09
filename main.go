package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", handleIndex)

	fmt.Println("Server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var filePath = path.Join("views", "index.html")
	var tmpl, err = template.ParseFiles(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data = map[string]interface{}{
		"title": "Gogogo",
		"name":  "Tiara",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
